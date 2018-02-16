package logger

import (
    "golang.org/x/sys/windows"
    "syscall"
    "fmt"
)

type winLogger struct {
    dll windows.Handle

    keyCallback uintptr
    pop uintptr
    wait uintptr
}

func GetLogger() (Logger, error) {
    return &winLogger{}, nil
}

func (wl *winLogger) Setup() error {
    var err error

    wl.dll, err = windows.LoadLibrary("logger.dll")
    if err != nil {
        return err
    }

    // TODO Get 'SetWindowsHookEx' and 'UnhookWindowsHookEx' and set the hook

    wl.keyCallback, err = windows.GetProcAddress(wl.dll, "keyCallback")
    if err != nil {
        return err
    }

    wl.pop, err = windows.GetProcAddress(wl.dll, "pop")
    if err != nil {
        return err
    }

    wl.wait, err = windows.GetProcAddress(wl.dll, "wait")
    return err
}

func (wl *winLogger) Clean() {

    if wl.dll != 0 {
        windows.FreeLibrary(wl.dll)
        wl.dll = 0
    }
}

func (wl *winLogger) Wait() error {
    _, _, err := syscall.Syscall(uintptr(wl.wait), 0, 0, 0, 0)
    return err
}

func (wl *winLogger) Pop() (uint16, error) {
    r1, _, err := syscall.Syscall(uintptr(wl.pop), 0, 0, 0, 0)
    return uint16(r1 & 0xFFFF) err
}

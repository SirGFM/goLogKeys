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
// Flags that a key was just released
const keyReleasedMask uint8 = (1 << 7)

// Checks if a syscall Call (of type syscal.Errno) actually returned an error.
func checkError(err error) error {
    werr := err.(syscall.Errno)
    if werr == 0 {
        return nil
    }
    return err
}

// Returns this platform's logger.
func GetLogger() (Logger, error) {
    return &winLogger{}, nil
}

// Initializes (loads libraries or whatever) everything required by the interface
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

    return nil
}

// Starts logging keys.
func (wl *winLogger) Start() error {
    var err error

    // Set the hook
    wl.hook1, wl.hook2, err = wl.setHook.Call(wh_KEYBOARD_LL, wl.keyCallback.Addr(),
        wl.klDll.Handle(), 0)
    if err = checkError(err); err != nil {
        return err
    }

    wl.done = true
    return nil
}

// Stops logging keys and clean up everything. Should be defer'ed before calling Setup().
func (wl *winLogger) Clean() {
    if wl.done {
        _, _, err := wl.clearHook.Call(wl.hook1, wl.hook2)
        err = checkError(err)
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

// Keeps the application busy, waiting to be closed.
func (wl *winLogger) Wait() error {
    _, _, err := wl.wait.Call()
    err = checkError(err)
    return err
}

// Removes the oldest key from the FIFO.
func (wl *winLogger) Pop() (k Key, ks KeyState, err error) {
    var r1 uintptr

    r1, _, err = wl.pop.Call()
    err = checkError(err)

    k = winVKeyToLoggerKey[uint8(r1 & 0xFF)]
    if (uint8(r1 >> 8) & keyReleasedMask) == 0 {
        ks = Pressed
    }

    return
}

// Map all recognized keys to their internal value. The values were taken from
// https://msdn.microsoft.com/en-US/library/windows/desktop/dd375731
var winVKeyToLoggerKey = map[uint8]Key {
    0x08: Backspace,
    0x09: Tab,
    0x0D: Return,
    0x1B: Esc,
    0x20: Space,
    0x21: PageUp,
    0x22: PageDown,
    0x23: End,
    0x24: Home,
    0x25: Left,
    0x26: Up,
    0x27: Right,
    0x28: Down,
    0x2D: Insert,
    0x2E: Delete,
    0x30: Key0,
    0x31: Key1,
    0x32: Key2,
    0x33: Key3,
    0x34: Key4,
    0x35: Key5,
    0x36: Key6,
    0x37: Key7,
    0x38: Key8,
    0x39: Key9,
    0x41: A,
    0x42: B,
    0x43: C,
    0x44: D,
    0x45: E,
    0x46: F,
    0x47: G,
    0x48: H,
    0x49: I,
    0x4A: J,
    0x4B: K,
    0x4C: L,
    0x4D: M,
    0x4E: N,
    0x4F: O,
    0x50: P,
    0x51: Q,
    0x52: R,
    0x53: S,
    0x54: T,
    0x55: U,
    0x56: V,
    0x57: W,
    0x58: X,
    0x59: Y,
    0x5A: Z,
    0x60: Numpad0,
    0x61: Numpad1,
    0x62: Numpad2,
    0x63: Numpad3,
    0x64: Numpad4,
    0x65: Numpad5,
    0x66: Numpad6,
    0x67: Numpad7,
    0x68: Numpad8,
    0x69: Numpad9,
    0x70: F1,
    0x71: F2,
    0x72: F3,
    0x73: F4,
    0x74: F5,
    0x75: F6,
    0x76: F7,
    0x77: F8,
    0x78: F9,
    0x79: F10,
    0x7A: F11,
    0x7B: F12,
    0xA0: LShift,
    0xA1: RShift,
    0xA2: LCtrl,
    0xA3: RCtrl,
    0xA4: LAlt,
    0xA5: RAlt,
    0xBC: Comma,
    0xBE: Period,
}

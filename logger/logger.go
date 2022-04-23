package logger

// Possible states of a given key
type KeyState uint8
const (
    Released KeyState = 0x00
    Pressed KeyState = 0x01
)

// Keys recognized and returned by the logger
type Key uint8
const (
    _Nil Key = iota
    Backspace
    Tab
    Return
    Esc
    Space
    PageUp
    PageDown
    End
    Home
    Left
    Up
    Right
    Down
    Insert
    Delete
    Key0
    Key1
    Key2
    Key3
    Key4
    Key5
    Key6
    Key7
    Key8
    Key9
    A
    B
    C
    D
    E
    F
    G
    H
    I
    J
    K
    L
    M
    N
    O
    P
    Q
    R
    S
    T
    U
    V
    W
    X
    Y
    Z
    Numpad0
    Numpad1
    Numpad2
    Numpad3
    Numpad4
    Numpad5
    Numpad6
    Numpad7
    Numpad8
    Numpad9
    F1
    F2
    F3
    F4
    F5
    F6
    F7
    F8
    F9
    F10
    F11
    F12
    LShift
    RShift
    LCtrl
    RCtrl
    LAlt
    RAlt
    Comma
    Period
    KeyCount
)

// Simple map to convert a string to a Key
var _getKey = map[string]Key {
    "backspace": Backspace,
    "tab": Tab,
    "return": Return,
    "esc": Esc,
    "space": Space,
    "pageup": PageUp,
    "pagedown": PageDown,
    "end": End,
    "home": Home,
    "left": Left,
    "up": Up,
    "right": Right,
    "down": Down,
    "insert": Insert,
    "delete": Delete,
    "0": Key0,
    "1": Key1,
    "2": Key2,
    "3": Key3,
    "4": Key4,
    "5": Key5,
    "6": Key6,
    "7": Key7,
    "8": Key8,
    "9": Key9,
    "a": A,
    "b": B,
    "c": C,
    "d": D,
    "e": E,
    "f": F,
    "g": G,
    "h": H,
    "i": I,
    "j": J,
    "k": K,
    "l": L,
    "m": M,
    "n": N,
    "o": O,
    "p": P,
    "q": Q,
    "r": R,
    "s": S,
    "t": T,
    "u": U,
    "v": V,
    "w": W,
    "x": X,
    "y": Y,
    "z": Z,
    "num0": Numpad0,
    "num1": Numpad1,
    "num2": Numpad2,
    "num3": Numpad3,
    "num4": Numpad4,
    "num5": Numpad5,
    "num6": Numpad6,
    "num7": Numpad7,
    "num8": Numpad8,
    "num9": Numpad9,
    "f1": F1,
    "f2": F2,
    "f3": F3,
    "f4": F4,
    "f5": F5,
    "f6": F6,
    "f7": F7,
    "f8": F8,
    "f9": F9,
    "f10": F10,
    "f11": F11,
    "f12": F12,
    "lshift": LShift,
    "rshift": RShift,
    "lctrl": LCtrl,
    "rctrl": RCtrl,
    "lalt": LAlt,
    "ralt": RAlt,
    "comma": Comma,
    "period": Period,
}

// GetKey retrieves the key from its name
func GetKey(key string) Key {
    return _getKey[key]
}

// A key logging interface. Each platform should implement its own 'GetLogger()',
// which returns a Logger.
type Logger interface {
    // Initializes (loads libraries or whatever) everything required by the interface.
    // The application SHALL NOT start logging keys after calling this.
    Setup() error
    // Starts logging keys.
    Start() error
    // Stops logging keys and clean up everything. Should be defer'ed before calling Setup().
    Clean()
    // Keeps the application busy, waiting to be closed.
    Wait() error
    // Removes the oldest key from the FIFO.
    Pop() (Key, KeyState, error)
    // Removes various keys from the FIFO. Arrays of keys and states may be supplied
    // to avoid alloc'ing memory.
    PopMulti(inKeys []Key, inStates []KeyState) (keys []Key,
        states []KeyState, err error)
}

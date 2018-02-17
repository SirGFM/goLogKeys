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
    _ Key = iota
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
)

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

// +build !dummy,linux

package logger

import (
    "encoding/binary"
    "errors"
    "io"
    "io/ioutil"
    "os"
    "os/exec"
    "os/signal"
    "syscall"
)

const buffer_len = 128

var Empty = errors.New("Empty buffer")

type xlibpyLogger struct {
    buf *keyBuffer
    cmd *exec.Cmd
    tmpfile string
    stdout io.Reader
    running bool
    wait chan struct{}
}

// Returns this platform's logger.
func GetLogger() (Logger, error) {
    return &xlibpyLogger{}, nil
}

// Initializes (loads libraries or whatever) everything required by the interface.
// The application SHALL NOT start logging keys after calling this.
func (xl *xlibpyLogger) Setup() error {
    // Create a temporary file with the xlib logger
    f, err := ioutil.TempFile("", "py_xlib_logger")
    if err != nil {
        return err
    }

    xl.tmpfile = f.Name()

    data := []byte(xlib_py_script)
    for total := len(data); total > 0; {
        n, err := f.Write(data)
        if err != nil {
            f.Close()
            return err
        }

        total -= n
        data = data[n:]
    }
    err = f.Close()
    if err != nil {
        return err
    }

    // Configure the command
    xl.cmd = exec.Command("python3", xl.tmpfile)
    xl.stdout, err = xl.cmd.StdoutPipe()
    if err != nil {
        return err
    }
    xl.running = false
    xl.wait = make(chan struct{}, )

    xl.buf = newKeyBuffer(buffer_len)

    return nil
}

func (xl *xlibpyLogger) run() {
    read := make([]byte, 4)

    for xl.running {
        n, err := xl.stdout.Read(read)
        if n != len(read) || err != nil {
            break
        }

        if read[1] == 0xff {
            // u32 with either -2, -3 or -4
            continue
        }
        s := KeyState(read[0])
        v := binary.BigEndian.Uint16(read[2:])

        k, ok := keycodeToLoggerKey[v]
        if ok {
            xl.buf.push(k, s)
        }
    }

    xl.wait <- struct{}{}
}

// Starts logging keys.
func (xl *xlibpyLogger) Start() error {
    err := xl.cmd.Start()
    if err != nil {
        return err
    }
    xl.running = true
    go xl.run()
    return nil
}

// Stops logging keys and clean up everything. Should be defer'ed before calling Setup().
func (xl *xlibpyLogger) Clean() {
    if xl.tmpfile != "" {
        os.Remove(xl.tmpfile)
    }
    if xl.cmd != nil && xl.cmd.Process != nil {
        xl.running = false
        xl.cmd.Process.Signal(syscall.SIGUSR1)
        _ = <-xl.wait
        xl.cmd.Wait()
    }
}

// Keeps the application busy, waiting to be closed.
func (xl *xlibpyLogger) Wait() error {
    // Install (possibly a repeated) SIGINT handler for killing the lib
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    _ = <-c
    return nil
}

// Removes the oldest key from the FIFO.
func (xl *xlibpyLogger) Pop() (Key, KeyState, error) {
    var err error
    k, s := xl.buf.pop()
    if k == _Nil {
        err = Empty
    }
    return k, s, err
}

func min(a, b int) int {
    if a < b {
        return a
    } else {
        return b
    }
}

// Removes various keys from the FIFO. Arrays of keys and states may be supplied
// to avoid alloc'ing memory.
func (xl *xlibpyLogger) PopMulti(inKeys []Key, inStates []KeyState) (keys []Key,
    states []KeyState, err error) {

    if len(inKeys) == 0 {
        keys = make([]Key, buffer_len)
        states = make([]KeyState, buffer_len)
    } else {
        keys = inKeys
        states = inStates
    }

    count := xl.buf.count_continuous()
    if count < len(keys) {
        xl.buf.fast_pop(keys[:count], states[:count])

        rem := xl.buf.count_continuous()
        rem = min(rem, len(keys[count:]))
        if rem > 0 {
            xl.buf.fast_pop(keys[count:count+rem], states[count:count+rem])
        }

        keys = keys[:count+rem]
        states = states[:count+rem]
    } else {
        xl.buf.fast_pop(keys, states)
    }

    return
}

var keycodeToLoggerKey = map[uint16]Key {
    14: Backspace,
    15: Tab,
    28: Return,
    1: Esc,
    57: Space,
    104: PageUp,
    109: PageDown,
    107: End,
    102: Home,
    105: Left,
    103: Up,
    106: Right,
    108: Down,
    110: Insert,
    111: Delete,
    11: Key0,
    2: Key1,
    3: Key2,
    4: Key3,
    5: Key4,
    6: Key5,
    7: Key6,
    8: Key7,
    9: Key8,
    10: Key9,
    30: A,
    48: B,
    46: C,
    32: D,
    18: E,
    33: F,
    34: G,
    35: H,
    23: I,
    36: J,
    37: K,
    38: L,
    50: M,
    49: N,
    24: O,
    25: P,
    16: Q,
    19: R,
    31: S,
    20: T,
    22: U,
    47: V,
    17: W,
    45: X,
    21: Y,
    44: Z,
    82: Numpad0,
    79: Numpad1,
    80: Numpad2,
    81: Numpad3,
    75: Numpad4,
    76: Numpad5,
    77: Numpad6,
    71: Numpad7,
    72: Numpad8,
    73: Numpad9,
    59: F1,
    60: F2,
    61: F3,
    62: F4,
    63: F5,
    64: F6,
    65: F7,
    66: F8,
    67: F9,
    68: F10,
    87: F11,
    88: F12,
    42: LShift,
    54: RShift,
    29: LCtrl,
    97: RCtrl,
    56: LAlt,
    100: RAlt,
    51: Comma,
    52: Period,
}

// +build dummy
package logger

import (
    "math/rand"
    "os"
    "os/signal"
    "syscall"
    "time"
)

type dummyLogger struct {
}

// Returns this platform's logger.
func GetLogger() (Logger, error) {
    return &dummyLogger{}, nil
}

// Initializes (loads libraries or whatever) everything required by the interface
func (dl *dummyLogger) Setup() error {
    return nil
}

// Starts logging keys.
func (wl *dummyLogger) Start() error {
    return nil
}

// Stops logging keys and clean up everything. Should be defer'ed before calling Setup().
func (wl *dummyLogger) Clean() {
}

// Keeps the application busy, waiting to be closed.
func (wl *dummyLogger) Wait() error {
    wait := make(chan os.Signal, 1)
    signal.Notify(wait, syscall.SIGUSR1)
    <-wait
    return nil
}

// Removes the oldest key from the FIFO.
func (wl *dummyLogger) Pop() (k Key, ks KeyState, err error) {
    k = Key(rand.Int31n(int32(_Max)))
    ks = KeyState(rand.Int31n(2))
    return
}

// Removes various keys from the FIFO. Arrays of keys and states may be supplied
// to avoid alloc'ing memory.
func (wl *dummyLogger) PopMulti(inKeys []Key, inStates []KeyState) (keys []Key,
    states []KeyState, err error) {

    // Wait another "frame"
    time.Sleep(16666 * time.Microsecond)

    // Generate a new dummy input
    num := int(rand.Int31n(65))

    for i := 0; i < num; i++ {
        k, ks, _ := wl.Pop()

        keys = append(keys, k)
        states = append(states, ks)
    }

    return
}

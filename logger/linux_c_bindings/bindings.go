package linux_c_bindings

// #cgo CFLAGS: -DNO_TEST_MAIN
// #include <errno.h>
// #include <stddef.h>
// #include <stdint.h>
// #include <stdlib.h>
// #include "bindings.h"
import "C"

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"unsafe"
)

type Event struct {
	// The event's key.
	Key uint16
	// The key's new state.
	Pressed uint8
}

type bindings struct {
	// The internal context to interface with the C library.
	ctx unsafe.Pointer
	// Whether the goroutine is running.
	running bool
	// WaitGroup for waiting for the goroutine.
	wg sync.WaitGroup
	// Events generated in the goroutine.
	keys chan Event
}

type Bindings interface {
	// Close releases every resource used by the internal library.
	Close() error

	// Run starts a goroutine that keeps monitoring for keyboard inputs.
	Run()

	// Pop blocks until an event is available, then reads it.
	Pop() Event
}

// getInputDevices lists every inputs device in the system.
func getInputDevices() ([]string, error) {
	var devices []string

	dirs, err := os.ReadDir("/dev/input")
	if err != nil {
		return nil, errors.Join(ErrReadInputs, err)
	}

	for _, dir := range dirs {
		if dir.IsDir() || (dir.Type() & fs.ModeCharDevice) != fs.ModeCharDevice {
			continue
		}
		devices = append(devices, dir.Name())
	}

	return devices, nil
}

// New initializes a new Bindings.
func New(bufSize int) (_ Bindings, rerr error) {
	devices, err := getInputDevices()
	if err != nil {
		return nil, err
	}

	ctx := C.new_context(C.int(bufSize))
	if ctx == C.NULL {
		return nil, ErrInitializeContext
	}
	defer func() {
		if rerr != nil {
			C.close_context(ctx)
		}
	}()

	for _, file := range devices {
		path := filepath.Join("/dev/input", file)

		cs := C.CString(path)
		rv := C.open_keyboard_fd(ctx, cs)
		C.free(unsafe.Pointer(cs))
		if rv == -1 {
			log.Printf("logger/linux_c_bindings: failed to open '%s'", path)
			return nil, ErrOpenDevice
		}
	}

	numFd := int(C.get_num_fd(ctx))
	if numFd == 0 {
		return nil, ErrNoDevice
	}

	bindings := bindings{
		ctx:  ctx,
		keys: make(chan Event, bufSize * 10 * numFd),
	}

	return &bindings, nil
}

// Close implements Bindings.
func (b *bindings) Close() error {
	b.running = false
	b.wg.Wait()

	if b.ctx != nil {
		close(b.keys)
		C.close_context(b.ctx)
		b.ctx = nil
	}

	return nil
}

// Run implements Bindings.
func (b *bindings) Run() {
	if b.running {
		return
	}

	defer b.wg.Done()

	b.running = true
	b.wg.Add(1)

	for b.running {
		ready, errno := C.wait_events(b.ctx);
		if ready == -1 {
			log.Printf("logger/linux_c_bindings: failed to wait for keyboard events (%v)", errno)
			continue
		}

		for i := C.int(0); i < C.get_num_fd(b.ctx) && C.get_ready(b.ctx) > 0; i++ {
			var rawEvents *C.struct_event

			tmp, errno := C.read_events(b.ctx, i, &rawEvents)
			if tmp == -1 {
				log.Printf("logger/linux_c_bindings: failed to get events (%v)", errno)
				continue
			}

			num := int(tmp)
			events := unsafe.Slice(rawEvents, num)

			for j := 0; j < num; j++ {
				b.keys <- Event{
					Key:     uint16(events[j].key),
					Pressed: uint8(events[j].status),
				}
			}
		}
	}
}

// Pop implements Bindings.
func (b *bindings) Pop() Event {
	return <-b.keys
}

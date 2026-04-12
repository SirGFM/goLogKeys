package linux_c_bindings

type errCode int

const (
	// Failed to read the input devices.
	ErrReadInputs errCode = iota
	// Failed to initialize the context.
	ErrInitializeContext
	// Failed to open an input device.
	ErrOpenDevice
	// No device was opened.
	ErrNoDevice
)

func (e errCode) Error() string {
	switch e {
	case ErrReadInputs:
		return "failed to read the input devices"
	case ErrInitializeContext:
		return "failed to initialize the context"
	case ErrOpenDevice:
		return "failed to open an input device"
	case ErrNoDevice:
		return "no device was opened"

	default:
		return "unknown error"
	}
}

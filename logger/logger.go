package logger

type Logger interface {
    Setup() error
    Clean()
    Wait() error
    Pop() (uint16, error)
}

package common

// Logger simple interface for logger
type Logger interface {
	Info(string)
	Warning(string)
}

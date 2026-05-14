package interfaces

import "io"

type LoggerInterface interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Warn(msg string, args ...any)
	Debug(msg string, args ...any)
	Log(level int8, msg string, args ...any)
	Fatal(msg string, args ...any)
	GetLogWriter() io.Writer
}

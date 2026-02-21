package logger

import (
	"fmt"
	"io"
	"log"
)

// Logger is a thin wrapper over log.Logger to keep dependencies minimal.
type Logger struct {
	std *log.Logger
}

// New creates a new logger writing to the provided writer.
func New(writer io.Writer) Logger {
	return Logger{
		std: log.New(writer, "", log.LstdFlags),
	}
}

// Info logs informational messages.
func (l Logger) Info(msg string, args ...any) {
	l.std.Printf("INFO "+msg, args...)
}

// Error logs error messages.
func (l Logger) Error(msg string, args ...any) {
	l.std.Printf("ERROR "+msg, args...)
}

// Sub returns a child logger with prefix.
func (l Logger) Sub(prefix string) Logger {
	return Logger{std: log.New(l.std.Writer(), fmt.Sprintf("[%s] ", prefix), log.LstdFlags)}
}

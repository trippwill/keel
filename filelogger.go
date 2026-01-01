package keel

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/trippwill/keel/engine"
)

// FileLogger writes render log events to the provided writer.
// It is safe for concurrent use.
type FileLogger struct {
	mu sync.Mutex
	w  io.Writer
}

// NewFileLogger returns a FileLogger that writes to the given writer.
func NewFileLogger(w io.Writer) *FileLogger {
	return &FileLogger{w: w}
}

// NewFileLoggerPath returns a FileLogger writing to the given path and the
// opened file. Call Close on the file when done.
func NewFileLoggerPath(path string) (*FileLogger, *os.File, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, nil, err
	}
	return NewFileLogger(file), file, nil
}

// Log writes a single log entry.
func (l *FileLogger) Log(event engine.LogEvent, path, msg string) {
	if l == nil || l.w == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	_, _ = fmt.Fprintf(l.w, "%s\t%s\t%s\n", event.String(), path, msg)
}

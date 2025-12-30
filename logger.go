package keel

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// LogEvent identifies the source of a render log message.
type LogEvent uint8

const (
	LogEventStackAlloc LogEvent = iota
	LogEventFrameRender
	LogEventRenderError
)

// LogEventFormats provides default formats for log messages.
// The map entries should be treated as read-only.
var LogEventFormats = map[LogEvent]string{
	LogEventStackAlloc:  "axis=%s total=%d slots=%d sizes=%v required=%d",
	LogEventFrameRender: "id=%v alloc=%dx%d frame=%dx%d content=%dx%d fit=%v",
	LogEventRenderError: "stage=%s err=%v",
}

func (e LogEvent) String() string {
	switch e {
	case LogEventStackAlloc:
		return "stack.alloc"
	case LogEventFrameRender:
		return "frame.render"
	case LogEventRenderError:
		return "render.error"
	default:
		return fmt.Sprintf("LogEvent(%d)", e)
	}
}

// LoggerFunc receives a render event, a slash-delimited slot path, and a formatted message.
type LoggerFunc func(event LogEvent, path, msg string)

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
func (l *FileLogger) Log(event LogEvent, path, msg string) {
	if l == nil || l.w == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	_, _ = fmt.Fprintf(l.w, "%s\t%s\t%s\n", event.String(), path, msg)
}

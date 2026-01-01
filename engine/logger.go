package engine

import (
	"fmt"
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

func (f *LoggerFunc) LogEvent(path string, event LogEvent, args ...any) {
	if f == nil || *f == nil {
		return
	}
	msgFormat, ok := LogEventFormats[event]
	if !ok {
		msgFormat = "event=%v"
		args = []any{event}
	}
	(*f)(event, path, fmt.Sprintf(msgFormat, args...))
}

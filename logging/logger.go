package logging

import (
	"context"
	"log/slog"
)

// Event identifies the source of a render log entry.
type Event string

const (
	EventStackAlloc  Event = "stack.alloc"
	EventFrameRender Event = "frame.render"
	EventRenderError Event = "render.error"
)

// LogEvent logs a structured render event to the provided logger.
func LogEvent(logger *slog.Logger, level slog.Level, event Event, path string, attrs ...slog.Attr) {
	if logger == nil {
		return
	}
	if !logger.Enabled(context.Background(), level) {
		return
	}
	if path != "" {
		attrs = append(attrs, slog.String("path", path))
	}
	attrs = append(attrs, slog.String("event", string(event)))
	logger.LogAttrs(context.Background(), level, "keel.render", attrs...)
}

// LogError logs a render error event at error level.
func LogError(logger *slog.Logger, path string, stage string, err error) {
	if err == nil {
		return
	}
	LogEvent(
		logger,
		slog.LevelError,
		EventRenderError,
		path,
		slog.String("stage", stage),
		slog.Any("err", err),
	)
}

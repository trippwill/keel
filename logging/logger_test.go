package logging

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"
)

type logEntry struct {
	level slog.Level
	msg   string
	attrs map[string]any
}

type captureHandler struct {
	entries *[]logEntry
	attrs   []slog.Attr
	groups  []string
}

func newCaptureHandler() (*captureHandler, *[]logEntry) {
	entries := []logEntry{}
	return &captureHandler{entries: &entries}, &entries
}

func (h *captureHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *captureHandler) Handle(_ context.Context, record slog.Record) error {
	attrs := map[string]any{}
	groupPrefix := strings.Join(h.groups, ".")
	addAttr := func(attr slog.Attr) {
		key := attr.Key
		if groupPrefix != "" {
			key = groupPrefix + "." + key
		}
		attrs[key] = attr.Value.Any()
	}

	for _, attr := range h.attrs {
		addAttr(attr)
	}
	record.Attrs(func(attr slog.Attr) bool {
		addAttr(attr)
		return true
	})

	*h.entries = append(*h.entries, logEntry{
		level: record.Level,
		msg:   record.Message,
		attrs: attrs,
	})
	return nil
}

func (h *captureHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	nextAttrs := append([]slog.Attr{}, h.attrs...)
	nextAttrs = append(nextAttrs, attrs...)
	return &captureHandler{
		entries: h.entries,
		attrs:   nextAttrs,
		groups:  append([]string{}, h.groups...),
	}
}

func (h *captureHandler) WithGroup(name string) slog.Handler {
	nextGroups := append([]string{}, h.groups...)
	nextGroups = append(nextGroups, name)
	return &captureHandler{
		entries: h.entries,
		attrs:   append([]slog.Attr{}, h.attrs...),
		groups:  nextGroups,
	}
}

func TestLogEvent(t *testing.T) {
	handler, entries := newCaptureHandler()
	logger := slog.New(handler)

	LogEvent(logger, slog.LevelDebug, EventStackAlloc, "/0", slog.String("axis", "Horizontal"))

	if len(*entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(*entries))
	}
	entry := (*entries)[0]
	if entry.msg != "keel.render" {
		t.Fatalf("unexpected message: %q", entry.msg)
	}
	if entry.level != slog.LevelDebug {
		t.Fatalf("unexpected level: %v", entry.level)
	}
	if entry.attrs["event"] != string(EventStackAlloc) {
		t.Fatalf("unexpected event attr: %v", entry.attrs["event"])
	}
	if entry.attrs["path"] != "/0" {
		t.Fatalf("unexpected path attr: %v", entry.attrs["path"])
	}
	if entry.attrs["axis"] != "Horizontal" {
		t.Fatalf("unexpected axis attr: %v", entry.attrs["axis"])
	}
}

func TestLogError(t *testing.T) {
	handler, entries := newCaptureHandler()
	logger := slog.New(handler)

	err := errors.New("boom")
	LogError(logger, "/1", "frame.content", err)

	if len(*entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(*entries))
	}
	entry := (*entries)[0]
	if entry.attrs["event"] != string(EventRenderError) {
		t.Fatalf("unexpected event attr: %v", entry.attrs["event"])
	}
	if entry.attrs["stage"] != "frame.content" {
		t.Fatalf("unexpected stage attr: %v", entry.attrs["stage"])
	}
	if gotErr, ok := entry.attrs["err"].(error); !ok || gotErr.Error() != "boom" {
		t.Fatalf("unexpected err attr: %v", entry.attrs["err"])
	}
}

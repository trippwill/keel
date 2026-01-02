package logging

import "testing"

func TestLogEventString(t *testing.T) {
	if LogEventStackAlloc.String() != "stack.alloc" {
		t.Fatalf("unexpected string")
	}
	if LogEventFrameRender.String() != "frame.render" {
		t.Fatalf("unexpected string")
	}
	if LogEventRenderError.String() != "render.error" {
		t.Fatalf("unexpected string")
	}
	if LogEvent(99).String() != "LogEvent(99)" {
		t.Fatalf("unexpected string for unknown")
	}
}

func TestLoggerFuncLogEvent(t *testing.T) {
	var gotEvent LogEvent
	var gotPath string
	var gotMsg string
	logger := LoggerFunc(func(event LogEvent, path, msg string) {
		gotEvent = event
		gotPath = path
		gotMsg = msg
	})
	logger.LogEvent("/", LogEventStackAlloc, "h", 10, 2, []int{3, 7}, 10)
	if gotEvent != LogEventStackAlloc {
		t.Fatalf("unexpected event")
	}
	if gotPath != "/" {
		t.Fatalf("unexpected path")
	}
	if gotMsg == "" {
		t.Fatalf("expected message")
	}
}

func TestLoggerFuncUnknownEventFormat(t *testing.T) {
	var gotMsg string
	logger := LoggerFunc(func(event LogEvent, path, msg string) {
		gotMsg = msg
	})
	logger.LogEvent("/", LogEvent(99))
	if gotMsg != "event=LogEvent(99)" {
		t.Fatalf("unexpected message: %q", gotMsg)
	}
}

func TestLoggerFuncNilSafe(t *testing.T) {
	var logger LoggerFunc
	logger.LogEvent("/", LogEventStackAlloc, "x")
}

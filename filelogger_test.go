package keel

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/trippwill/keel/engine"
)

func TestFileLoggerWritesEntry(t *testing.T) {
	var buf bytes.Buffer
	logger := NewFileLogger(&buf)
	logger.Log(engine.LogEventFrameRender, "0/1", "ok")

	got := buf.String()
	if got == "" {
		t.Fatalf("expected output")
	}
	if !strings.HasPrefix(got, "frame.render\t0/1\tok\n") {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestFileLoggerNilWriterNoPanic(t *testing.T) {
	var logger *FileLogger
	logger.Log(engine.LogEventRenderError, "", "nope")
}

func TestNewFileLoggerPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "keel.log")

	logger, file, err := NewFileLoggerPath(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	logger.Log(engine.LogEventRenderError, "/", "oops")
	if err := file.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("unexpected read error: %v", err)
	}
	if !strings.HasPrefix(string(data), "render.error\t/\toops\n") {
		t.Fatalf("unexpected output: %q", string(data))
	}
}

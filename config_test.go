package keel

import (
	"testing"

	"github.com/trippwill/keel/logging"
)

func TestConfigNilReceiver(t *testing.T) {
	var config *Config
	if config.Logger() != nil {
		t.Fatalf("expected nil logger")
	}
	if config.Debug() {
		t.Fatalf("expected debug false")
	}
	config.SetLogger(logging.LoggerFunc(func(logging.LogEvent, string, string) {}))
	config.SetDebug(true)
}

func TestConfigSetters(t *testing.T) {
	config := NewConfig()
	called := false
	logger := logging.LoggerFunc(func(logging.LogEvent, string, string) { called = true })
	config.SetLogger(logger)
	config.SetDebug(true)
	if config.Logger() == nil {
		t.Fatalf("expected logger")
	}
	if !config.Debug() {
		t.Fatalf("expected debug true")
	}
	config.Logger().LogEvent("/", logging.LogEventRenderError, "x")
	if !called {
		t.Fatalf("expected logger called")
	}
}

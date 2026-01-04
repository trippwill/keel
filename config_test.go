package keel

import (
	"io"
	"log/slog"
	"testing"
)

func TestConfigNilReceiver(t *testing.T) {
	var config *Config
	if config.Logger() != nil {
		t.Fatalf("expected nil logger")
	}
	if config.Debug() {
		t.Fatalf("expected debug false")
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
	config.SetLogger(logger)
	config.SetDebug(true)
}

func TestConfigSetters(t *testing.T) {
	config := NewConfig()
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
	config.SetLogger(logger)
	config.SetDebug(true)
	if config.Logger() != logger {
		t.Fatalf("expected logger")
	}
	if !config.Debug() {
		t.Fatalf("expected debug true")
	}
}

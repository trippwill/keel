package keel

import "log/slog"

// Config stores shared render settings like logging and debug state.
// It is safe to share a single config across multiple renderers.
type Config struct {
	logger *slog.Logger
	debug  bool
}

// NewConfig returns a new renderer configuration with the default settings.
func NewConfig() *Config {
	return &Config{}
}

// Logger returns the configured logger, if any.
func (c *Config) Logger() *slog.Logger {
	if c == nil {
		return nil
	}
	return c.logger
}

// SetLogger sets the render logger.
func (c *Config) SetLogger(logger *slog.Logger) {
	if c == nil {
		return
	}
	c.logger = logger
}

// Debug reports whether debug rendering is enabled.
func (c *Config) Debug() bool {
	if c == nil {
		return false
	}
	return c.debug
}

// SetDebug sets debug rendering.
func (c *Config) SetDebug(debug bool) {
	if c == nil {
		return
	}
	c.debug = debug
}

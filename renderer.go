package keel

import (
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/trippwill/keel/engine"
)

// StyleProvider returns a style for the given frame ID. Nil means "no style".
// Returned styles are treated as immutable and are safe to cache; the renderer
// copies them before mutation.
type StyleProvider[KID KeelID] func(id KID) *gloss.Style

// ContentProvider returns content for the given frame allocation.
// Providers should respect ContentWidth/ContentHeight.
// FitMode will be applied after content is retrieved.
type ContentProvider[KID KeelID] func(id KID, info FrameInfo) (string, error)

// RenderConfig stores shared render settings like logging and debug state.
// It is safe to share a single config across multiple renderers.
type RenderConfig struct {
	logger engine.LoggerFunc
	debug  bool
}

// NewRenderConfig returns a new render config with default settings.
func NewRenderConfig() *RenderConfig {
	return &RenderConfig{}
}

// Logger returns the configured logger, if any.
func (c *RenderConfig) Logger() engine.LoggerFunc {
	if c == nil {
		return nil
	}
	return c.logger
}

// SetLogger sets the render logger.
func (c *RenderConfig) SetLogger(logger engine.LoggerFunc) {
	if c == nil {
		return
	}
	c.logger = logger
}

// Debug reports whether debug rendering is enabled.
func (c *RenderConfig) Debug() bool {
	if c == nil {
		return false
	}
	return c.debug
}

// SetDebug toggles debug rendering.
func (c *RenderConfig) SetDebug(debug bool) {
	if c == nil {
		return
	}
	c.debug = debug
}

// Renderer owns render providers and uses a shared config for logging/debugging.
type Renderer[KID KeelID] struct {
	config    *RenderConfig
	spec      Spec
	style     StyleProvider[KID]
	content   ContentProvider[KID]
	layout    engine.Layout[KID]
	last      Size
	hasLayout bool
}

// NewRenderer returns a renderer for the given spec with a fresh config.
func NewRenderer[KID KeelID](spec Spec, styleProvider StyleProvider[KID], contentProvider ContentProvider[KID]) *Renderer[KID] {
	return &Renderer[KID]{
		config:  NewRenderConfig(),
		spec:    spec,
		style:   styleProvider,
		content: contentProvider,
	}
}

// NewRendererWithConfig returns a renderer for the given spec using the provided config.
func NewRendererWithConfig[KID KeelID](config *RenderConfig, spec Spec, styleProvider StyleProvider[KID], contentProvider ContentProvider[KID]) *Renderer[KID] {
	if config == nil {
		config = NewRenderConfig()
	}
	return &Renderer[KID]{
		config:  config,
		spec:    spec,
		style:   styleProvider,
		content: contentProvider,
	}
}

// Config returns the renderer's config, allocating one if needed.
func (r *Renderer[KID]) Config() *RenderConfig {
	if r == nil {
		return nil
	}
	if r.config == nil {
		r.config = NewRenderConfig()
	}
	return r.config
}

// SetConfig replaces the renderer config.
// Invalidates cached layout state.
func (r *Renderer[KID]) SetConfig(config *RenderConfig) {
	if r == nil {
		return
	}
	if config == nil {
		config = NewRenderConfig()
	}
	r.config = config
	r.Invalidate()
}

// SetStyleProvider replaces the renderer style provider.
func (r *Renderer[KID]) SetStyleProvider(p StyleProvider[KID]) {
	if r == nil {
		return
	}
	r.style = p
}

// SetContentProvider replaces the renderer content provider.
func (r *Renderer[KID]) SetContentProvider(p ContentProvider[KID]) {
	if r == nil {
		return
	}
	r.content = p
}

// Invalidate clears cached layout state.
func (r *Renderer[KID]) Invalidate() {
	if r == nil {
		return
	}
	r.hasLayout = false
}

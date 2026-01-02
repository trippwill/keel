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

// Renderer owns render providers and uses a shared config for logging/debugging.
type Renderer[KID KeelID] struct {
	config    *Config
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
		config:  NewConfig(),
		spec:    spec,
		style:   styleProvider,
		content: contentProvider,
	}
}

// NewRendererWithConfig returns a renderer for the given spec using the provided config.
func NewRendererWithConfig[KID KeelID](config *Config, spec Spec, styleProvider StyleProvider[KID], contentProvider ContentProvider[KID]) *Renderer[KID] {
	if config == nil {
		config = NewConfig()
	}
	return &Renderer[KID]{
		config:  config,
		spec:    spec,
		style:   styleProvider,
		content: contentProvider,
	}
}

// Config returns the renderer's config, allocating one if needed.
func (r *Renderer[KID]) Config() *Config {
	if r == nil {
		return nil
	}
	if r.config == nil {
		r.config = NewConfig()
	}
	return r.config
}

// SetConfig replaces the renderer config.
// Invalidates cached layout state.
func (r *Renderer[KID]) SetConfig(config *Config) {
	if r == nil {
		return
	}
	if config == nil {
		config = NewConfig()
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

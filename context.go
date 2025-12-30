package keel

import (
	gloss "github.com/charmbracelet/lipgloss"
)

// KeelID is a comparable type used as a stable identifier for frames and resources.
type KeelID comparable

// StyleProvider returns a style for the given frame ID. Nil means "no style".
// Returned styles are treated as immutable and are safe to cache; the renderer
// copies them before mutation.
type StyleProvider[KID KeelID] func(id KID) *gloss.Style

// ContentProvider returns content for the given frame allocation.
// Providers should respect ContentWidth/ContentHeight.
// FitMode will be applied after content is retrieved.
type ContentProvider[KID KeelID] func(id KID, info RenderInfo) (string, error)

// Context provides rendering inputs for a render pass, including the
// content/style providers used by frames.
type Context[KID KeelID] struct {
	// Provides style information for a [FrameSpec].
	StyleProvider StyleProvider[KID]

	// Provides string content for a [FrameSpec].
	ContentProvider ContentProvider[KID]

	// Optional logger for render events.
	Logger LoggerFunc
}

func NewContext[KID KeelID](styleProvider StyleProvider[KID], contentProvider ContentProvider[KID]) Context[KID] {
	return Context[KID]{
		StyleProvider:   styleProvider,
		ContentProvider: contentProvider,
		Logger:          nil,
	}
}

func DefaultContext[KID KeelID]() Context[KID] {
	return Context[KID]{}
}

// WithStyleProvider returns a copy of the context with the given style provider.
func (c Context[KID]) WithStyleProvider(p StyleProvider[KID]) Context[KID] {
	return Context[KID]{
		StyleProvider:   p,
		ContentProvider: c.ContentProvider,
		Logger:          c.Logger,
	}
}

// WithContentProvider returns a copy of the context with the given content provider.
func (c Context[KID]) WithContentProvider(p ContentProvider[KID]) Context[KID] {
	return Context[KID]{
		StyleProvider:   c.StyleProvider,
		ContentProvider: p,
		Logger:          c.Logger,
	}
}

// WithLogger returns a copy of the context with the given logger.
func (c Context[KID]) WithLogger(logger LoggerFunc) Context[KID] {
	return Context[KID]{
		StyleProvider:   c.StyleProvider,
		ContentProvider: c.ContentProvider,
		Logger:          logger,
	}
}

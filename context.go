package keel

import (
	gloss "github.com/charmbracelet/lipgloss"
)

// KeelID is a comparable type used as a stable identifier for blocks and resources.
type KeelID comparable

// StyleProvider returns a style for the given block ID. Nil means "no style".
// Returned styles are treated as immutable and are safe to cache; the renderer
// copies them before mutation.
type StyleProvider[KID KeelID] func(id KID) *gloss.Style

// ContentProvider returns content for the given block allocation.
// Providers should respect ContentWidth/ContentHeight.
// FitMode will be applied after content is retrieved.
type ContentProvider[KID KeelID] func(id KID, info RenderInfo) (string, error)

// Context provides rendering inputs for a render pass, including allocation
// size and the content/style providers used by blocks.
type Context[KID KeelID] struct {
	// Width and Height define the total space for rendering.
	Width, Height int

	// Provides style information for a [Block].
	StyleProvider StyleProvider[KID]

	// Provides string content for a [Block].
	ContentProvider ContentProvider[KID]

	// Optional logger for render events.
	Logger LoggerFunc
}

func NewContext[KID KeelID](width, height int, styleProvider StyleProvider[KID], contentProvider ContentProvider[KID]) Context[KID] {
	return Context[KID]{
		Width:           width,
		Height:          height,
		StyleProvider:   styleProvider,
		ContentProvider: contentProvider,
		Logger:          nil,
	}
}

// WithSize returns a copy of the context with updated dimensions.
func (c Context[KID]) WithSize(width, height int) Context[KID] {
	return Context[KID]{
		Width:           width,
		Height:          height,
		StyleProvider:   c.StyleProvider,
		ContentProvider: c.ContentProvider,
		Logger:          c.Logger,
	}
}

// WithWidth returns a copy of the context with an updated width.
func (c Context[KID]) WithWidth(width int) Context[KID] {
	return Context[KID]{
		Width:           width,
		Height:          c.Height,
		StyleProvider:   c.StyleProvider,
		ContentProvider: c.ContentProvider,
		Logger:          c.Logger,
	}
}

// WithHeight returns a copy of the context with an updated height.
func (c Context[KID]) WithHeight(height int) Context[KID] {
	return Context[KID]{
		Width:           c.Width,
		Height:          height,
		StyleProvider:   c.StyleProvider,
		ContentProvider: c.ContentProvider,
		Logger:          c.Logger,
	}
}

// WithStyleProvider returns a copy of the context with the given style provider.
func (c Context[KID]) WithStyleProvider(p StyleProvider[KID]) Context[KID] {
	return Context[KID]{
		Width:           c.Width,
		Height:          c.Height,
		StyleProvider:   p,
		ContentProvider: c.ContentProvider,
		Logger:          c.Logger,
	}
}

// WithContentProvider returns a copy of the context with the given content provider.
func (c Context[KID]) WithContentProvider(p ContentProvider[KID]) Context[KID] {
	return Context[KID]{
		Width:           c.Width,
		Height:          c.Height,
		StyleProvider:   c.StyleProvider,
		ContentProvider: p,
		Logger:          c.Logger,
	}
}

// WithLogger returns a copy of the context with the given logger.
func (c Context[KID]) WithLogger(logger LoggerFunc) Context[KID] {
	return Context[KID]{
		Width:           c.Width,
		Height:          c.Height,
		StyleProvider:   c.StyleProvider,
		ContentProvider: c.ContentProvider,
		Logger:          logger,
	}
}

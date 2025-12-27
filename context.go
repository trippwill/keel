package keel

import gloss "github.com/charmbracelet/lipgloss"

// KeelID is a comparable type used as an identifier for Renderables and resources.
type KeelID comparable

// StyleProvider returns a style for the given ID. Nil means "no style".
// Returned styles are treated as immutable and may be cached.
type StyleProvider[KID KeelID] func(id KID) *gloss.Style

// ContentProvider returns content for the given render allocation.
type ContentProvider[KID KeelID] func(info RenderInfo[KID]) (string, error)

// RenderInfo describes the allocated space for a leaf Renderable.
type RenderInfo[KID KeelID] struct {
	ID                          KID
	Width, Height               int
	ContentWidth, ContentHeight int
	FrameWidth, FrameHeight     int
	Clip                        ClipConstraint
}

// Context provides rendering context for a Renderable.
type Context[KID KeelID] struct {
	// Width and Height define the total space for rendering.
	Width, Height int

	// Provides style information for a [Block].
	StyleProvider StyleProvider[KID]

	// Provides string content for a [Block].
	ContentProvider ContentProvider[KID]
}

// WithSize returns a copy of the context with updated dimensions.
func (c Context[KID]) WithSize(width, height int) Context[KID] {
	return Context[KID]{
		Width:           width,
		Height:          height,
		StyleProvider:   c.StyleProvider,
		ContentProvider: c.ContentProvider,
	}
}

// WithWidth returns a copy of the context with an updated width.
func (c Context[KID]) WithWidth(width int) Context[KID] {
	return Context[KID]{
		Width:           width,
		Height:          c.Height,
		StyleProvider:   c.StyleProvider,
		ContentProvider: c.ContentProvider,
	}
}

// WithHeight returns a copy of the context with an updated height.
func (c Context[KID]) WithHeight(height int) Context[KID] {
	return Context[KID]{
		Width:           c.Width,
		Height:          height,
		StyleProvider:   c.StyleProvider,
		ContentProvider: c.ContentProvider,
	}
}

// WithStyleProvider returns a copy of the context with the given style provider.
func (c Context[KID]) WithStyleProvider(p StyleProvider[KID]) Context[KID] {
	return Context[KID]{
		Width:           c.Width,
		Height:          c.Height,
		StyleProvider:   p,
		ContentProvider: c.ContentProvider,
	}
}

// WithContentProvider returns a copy of the context with the given content provider.
func (c Context[KID]) WithContentProvider(p ContentProvider[KID]) Context[KID] {
	return Context[KID]{
		Width:           c.Width,
		Height:          c.Height,
		StyleProvider:   c.StyleProvider,
		ContentProvider: p,
	}
}

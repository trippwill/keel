//go:generate stringer -type=NodeKind
package keel

import gloss "github.com/charmbracelet/lipgloss"

// KeelID is a comparable type used as an identifier for Renderables and resources.
type KeelID comparable

// NodeKind describes whether a node is a container or leaf content.
type NodeKind uint8

const (
	NodeContent NodeKind = iota
	NodeContainer
)

// StyleProvider returns a style for the given ID and node kind. Nil means "no style".
type StyleProvider[KID KeelID] func(id KID, kind NodeKind) *gloss.Style

// ContentProvider returns content for the given ID.
type ContentProvider[KID KeelID] func(id KID) (string, error)

// Context provides rendering context for a Renderable.
type Context[KID KeelID] struct {
	// Width and Height define the space for rendering.
	Width, Height int

	// Provides style information for a given KID.
	StyleProvider StyleProvider[KID]

	// Provides string content for a given KID.
	// Not called for containers; only for leaf Renderables.
	ContentProvider ContentProvider[KID]
}

func (c Context[KID]) WithSize(width, height int) Context[KID] {
	return Context[KID]{
		Width:           width,
		Height:          height,
		StyleProvider:   c.StyleProvider,
		ContentProvider: c.ContentProvider,
	}
}

func (c Context[KID]) WithWidth(width int) Context[KID] {
	return Context[KID]{
		Width:           width,
		Height:          c.Height,
		StyleProvider:   c.StyleProvider,
		ContentProvider: c.ContentProvider,
	}
}

func (c Context[KID]) WithHeight(height int) Context[KID] {
	return Context[KID]{
		Width:           c.Width,
		Height:          height,
		StyleProvider:   c.StyleProvider,
		ContentProvider: c.ContentProvider,
	}
}

func (c Context[KID]) WithStyleProvider(p StyleProvider[KID]) Context[KID] {
	return Context[KID]{
		Width:           c.Width,
		Height:          c.Height,
		StyleProvider:   p,
		ContentProvider: c.ContentProvider,
	}
}

func (c Context[KID]) WithContentProvider(p ContentProvider[KID]) Context[KID] {
	return Context[KID]{
		Width:           c.Width,
		Height:          c.Height,
		StyleProvider:   c.StyleProvider,
		ContentProvider: p,
	}
}

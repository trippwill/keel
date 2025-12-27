package keel

// PanelSpec is the leaf Renderable that displays content.
type PanelSpec[KID KeelID] struct {
	ExtentConstraint // Size specification for the panel
	ClipConstraint   // Maximum content size for the panel
	id               KID
}

// GetID implements [Block].
func (p *PanelSpec[KID]) GetID() KID {
	return p.id
}

var _ Block[string] = (*PanelSpec[string])(nil)

// Panel creates a new PanelSpec with the given ID.
func Panel[KID KeelID](extent ExtentConstraint, id KID) *PanelSpec[KID] {
	return PanelClip(extent, ClipConstraint{}, id)
}

// PanelClip creates a new PanelSpec with the given ID and content clip.
func PanelClip[KID KeelID](extent ExtentConstraint, clip ClipConstraint, id KID) *PanelSpec[KID] {
	return &PanelSpec[KID]{
		ExtentConstraint: extent,
		ClipConstraint:   clip,
		id:               id,
	}
}

package keel

// PanelSpec is the leaf Renderable that displays content.
type PanelSpec[KID KeelID] struct {
	id KID
}

var _ Renderable[string] = (*PanelSpec[string])(nil)

// Panel creates a new PanelSpec with the given ID.
func Panel[KID KeelID](id KID) *PanelSpec[KID] {
	return &PanelSpec[KID]{id: id}
}

func (p *PanelSpec[KID]) GetID() KID { return p.id }

// Render implements [Renderable].
func (p *PanelSpec[KID]) Render(ctx Context[KID]) (string, error) {
	return renderPanel(p, ctx)
}

package keel

import "testing"

func TestResolveBuildsRects(t *testing.T) {
	layout := Row(FlexUnit(),
		Panel(Fixed(3), "a"),
		Col(FlexUnit(),
			Panel(Fixed(2), "b"),
			Panel(FlexUnit(), "c"),
		),
	)

	ctx := Context[string]{}
	size := Size{Width: 10, Height: 5}
	resolved, err := Resolve[string](ctx, layout, size)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resolved.Root.Rect.Width != 10 || resolved.Root.Rect.Height != 5 {
		t.Fatalf("expected root 10x5, got %dx%d", resolved.Root.Rect.Width, resolved.Root.Rect.Height)
	}

	if len(resolved.Root.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(resolved.Root.Children))
	}

	left := resolved.Root.Children[0]
	right := resolved.Root.Children[1]

	if left.Rect.X != 0 || left.Rect.Width != 3 || left.Rect.Height != 5 {
		t.Fatalf("unexpected left rect: %+v", left.Rect)
	}
	if right.Rect.X != 3 || right.Rect.Width != 7 || right.Rect.Height != 5 {
		t.Fatalf("unexpected right rect: %+v", right.Rect)
	}

	if len(right.Children) != 2 {
		t.Fatalf("expected 2 children under right, got %d", len(right.Children))
	}

	top := right.Children[0]
	bottom := right.Children[1]
	if top.Rect.Y != 0 || top.Rect.Height != 2 {
		t.Fatalf("unexpected top rect: %+v", top.Rect)
	}
	if bottom.Rect.Y != 2 || bottom.Rect.Height != 3 {
		t.Fatalf("unexpected bottom rect: %+v", bottom.Rect)
	}
}

func TestRenderResolvedMatchesRender(t *testing.T) {
	layout := Row(FlexUnit(),
		Panel(Fixed(2), "a"),
		Panel(FlexUnit(), "b"),
	)

	ctx := Context[string]{
		ContentProvider: func(id string, _ RenderInfo) (string, error) {
			switch id {
			case "a":
				return "aa", nil
			case "b":
				return "bbb", nil
			default:
				return "", &UnknownBlockIDError{ID: id}
			}
		},
	}
	size := Size{Width: 5, Height: 1}

	want, err := Render(ctx, layout, size)
	if err != nil {
		t.Fatalf("unexpected render error: %v", err)
	}

	resolved, err := Resolve[string](ctx, layout, size)
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}

	got, err := RenderResolved(ctx, resolved)
	if err != nil {
		t.Fatalf("unexpected resolved render error: %v", err)
	}
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

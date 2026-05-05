package panels

type HeadersPanel struct{}

func (p HeadersPanel) View(width, height int, focused bool) string {
	return box(width, height, focused).Render(labelSt.Render("HEADERS"))
}

package panels

type URLPanel struct{}

func (p URLPanel) View(width, height int, focused bool) string {
	return box(width, height, focused).Render(labelSt.Render("URL"))
}

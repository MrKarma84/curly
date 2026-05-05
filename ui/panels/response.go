package panels

type ResponsePanel struct{}

func (p ResponsePanel) View(width, height int, focused bool) string {
	return box(width, height, focused).Render(labelSt.Render("RESPONSE"))
}

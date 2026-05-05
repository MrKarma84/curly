package panels

type MethodPanel struct{}

func (p MethodPanel) View(width, height int, focused bool) string {
	return box(width, height, focused).Render(labelSt.Render("METHOD"))
}

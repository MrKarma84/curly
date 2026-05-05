package panels

import "github.com/charmbracelet/lipgloss"

var (
	activeBorder   = lipgloss.Color("#7C3AED")
	inactiveBorder = lipgloss.Color("#4B5563")
	labelColor     = lipgloss.Color("#9CA3AF")
)

var labelSt = lipgloss.NewStyle().
	Foreground(labelColor).
	Bold(true)

func box(width, height int, focused bool) lipgloss.Style {
	border := inactiveBorder
	if focused {
		border = activeBorder
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(border).
		Width(max(0, width-2)).
		Height(max(0, height-2))
}

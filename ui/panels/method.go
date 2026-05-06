package panels

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var methodColors = map[string]lipgloss.Color{
	"GET":    "#10B981",
	"POST":   "#F59E0B",
	"PUT":    "#3B82F6",
	"PATCH":  "#F97316",
	"DELETE": "#EF4444",
}

type MethodPanel struct {
	methods  []string
	selected int
}

func NewMethodPanel() MethodPanel {
	return MethodPanel{
		methods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
	}
}

func (p MethodPanel) Selected() string {
	return p.methods[p.selected]
}

func (p MethodPanel) SetSelected(method string) MethodPanel {
	for i, m := range p.methods {
		if m == method {
			p.selected = i
			return p
		}
	}
	return p
}

func (p MethodPanel) Update(msg tea.KeyMsg) MethodPanel {
	switch msg.String() {
	case "up", "k":
		if p.selected > 0 {
			p.selected--
		}
	case "down", "j":
		if p.selected < len(p.methods)-1 {
			p.selected++
		}
	}
	return p
}

func (p MethodPanel) View(width, height int, focused bool) string {
	rows := []string{labelSt.Render("METHOD"), ""}

	for i, method := range p.methods {
		color := methodColors[method]
		var row string

		if i == p.selected {
			cursor := "  "
			if focused {
				cursor = "► "
			}
			row = lipgloss.NewStyle().Foreground(color).Bold(focused).Render(cursor + method)
		} else {
			row = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render("  " + method)
		}
		rows = append(rows, row)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return box(width, height, focused).Render(content)
}

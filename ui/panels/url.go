package panels

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type URLPanel struct {
	input textinput.Model
}

func NewURLPanel() URLPanel {
	ti := textinput.New()
	ti.Placeholder = "https://api.example.com/users"
	ti.CharLimit = 500
	return URLPanel{input: ti}
}

func (p URLPanel) Value() string {
	return p.input.Value()
}

func (p URLPanel) SetValue(url string) URLPanel {
	p.input.SetValue(url)
	return p
}

func (p URLPanel) Focus() (URLPanel, tea.Cmd) {
	cmd := p.input.Focus()
	return p, cmd
}

func (p URLPanel) Blur() URLPanel {
	p.input.Blur()
	return p
}

func (p URLPanel) Update(msg tea.Msg) (URLPanel, tea.Cmd) {
	var cmd tea.Cmd
	p.input, cmd = p.input.Update(msg)
	return p, cmd
}

func (p URLPanel) View(width, height int, focused bool) string {
	p.input.Width = width - 4
	content := lipgloss.JoinVertical(lipgloss.Left,
		labelSt.Render("URL"),
		p.input.View(),
	)
	return box(width, height, focused).Render(content)
}

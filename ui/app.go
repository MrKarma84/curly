package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED")).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Italic(true)

	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#374151")).
			MarginTop(2)
)

type Model struct {
	width  int
	height int
}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	title := titleStyle.Render("curly")
	subtitle := subtitleStyle.Render("A keyboard-driven TUI HTTP client")
	hint := hintStyle.Render("Press q to quit")

	content := lipgloss.JoinVertical(lipgloss.Center, title, subtitle, hint)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

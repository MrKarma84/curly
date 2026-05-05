package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/MrKarma84/curly/ui/panels"
)

const (
	panelMethod = iota
	panelURL
	panelHeaders
	panelResponse
	panelCount
)

const methodWidth = 14
const topHeight = 3

var hintSt = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))

type Model struct {
	width    int
	height   int
	focused  int
	method   panels.MethodPanel
	url      panels.URLPanel
	headers  panels.HeadersPanel
	response panels.ResponsePanel
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
		case "tab":
			m.focused = (m.focused + 1) % panelCount
		case "shift+tab":
			m.focused = (m.focused - 1 + panelCount) % panelCount
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	bottomHeight := m.height - topHeight - 1
	headersWidth := m.width / 2
	responseWidth := m.width - headersWidth

	topRow := lipgloss.JoinHorizontal(lipgloss.Top,
		m.method.View(methodWidth, topHeight, m.focused == panelMethod),
		m.url.View(m.width-methodWidth, topHeight, m.focused == panelURL),
	)

	bottomRow := lipgloss.JoinHorizontal(lipgloss.Top,
		m.headers.View(headersWidth, bottomHeight, m.focused == panelHeaders),
		m.response.View(responseWidth, bottomHeight, m.focused == panelResponse),
	)

	hint := hintSt.Render("Tab · next panel   Shift+Tab · prev   q · quit")

	return lipgloss.JoinVertical(lipgloss.Left, topRow, bottomRow, hint)
}

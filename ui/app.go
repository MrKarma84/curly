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

const (
	methodWidth = 14
	urlHeight   = 3
)

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
	return Model{
		method: panels.NewMethodPanel(),
	}
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
		default:
			if m.focused == panelMethod {
				m.method = m.method.Update(msg)
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	rightWidth := m.width - methodWidth
	panelHeight := m.height - 1 // -1 for hint line
	bottomHeight := panelHeight - urlHeight
	headersWidth := rightWidth / 2
	responseWidth := rightWidth - headersWidth

	methodView := m.method.View(methodWidth, panelHeight, m.focused == panelMethod)

	urlView := m.url.View(rightWidth, urlHeight, m.focused == panelURL)
	headersView := m.headers.View(headersWidth, bottomHeight, m.focused == panelHeaders)
	responseView := m.response.View(responseWidth, bottomHeight, m.focused == panelResponse)

	rightCol := lipgloss.JoinVertical(lipgloss.Left,
		urlView,
		lipgloss.JoinHorizontal(lipgloss.Top, headersView, responseView),
	)

	hint := hintSt.Render("Tab · next panel   ↑↓ · select method   q · quit")

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, methodView, rightCol),
		hint,
	)
}

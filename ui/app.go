package ui

import (
	"github.com/MrKarma84/curly/httpclient"
	"github.com/MrKarma84/curly/ui/panels"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

type ResponseMsg httpclient.Response

func doRequest(method, url string, headers map[string]string) tea.Cmd {
	return func() tea.Msg {
		return ResponseMsg(httpclient.Send(method, url, headers))
	}
}

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
		method:  panels.NewMethodPanel(),
		url:     panels.NewURLPanel(),
		headers: panels.NewHeadersPanel(),
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

	case ResponseMsg:
		m.response = m.response.SetResponse(httpclient.Response(msg))

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "ctrl+r":
			url := m.url.Value()
			if url == "" {
				return m, nil
			}
			m.response = m.response.SetLoading()
			return m, doRequest(m.method.Selected(), url, m.headers.Headers())

		case "tab":
			if m.focused == panelHeaders && m.headers.IsEditing() {
				var cmd tea.Cmd
				m.headers, cmd = m.headers.Update(msg)
				return m, cmd
			}
			if m.focused == panelURL {
				m.url = m.url.Blur()
			}
			m.focused = (m.focused + 1) % panelCount
			if m.focused == panelURL {
				var cmd tea.Cmd
				m.url, cmd = m.url.Focus()
				return m, cmd
			}

		case "shift+tab":
			if m.focused == panelHeaders && m.headers.IsEditing() {
				var cmd tea.Cmd
				m.headers, cmd = m.headers.Update(msg)
				return m, cmd
			}
			if m.focused == panelURL {
				m.url = m.url.Blur()
			}
			m.focused = (m.focused - 1 + panelCount) % panelCount
			if m.focused == panelURL {
				var cmd tea.Cmd
				m.url, cmd = m.url.Focus()
				return m, cmd
			}

		default:
			switch m.focused {
			case panelMethod:
				m.method = m.method.Update(msg)
			case panelURL:
				var cmd tea.Cmd
				m.url, cmd = m.url.Update(msg)
				return m, cmd
			case panelHeaders:
				var cmd tea.Cmd
				m.headers, cmd = m.headers.Update(msg)
				return m, cmd
			case panelResponse:
				var cmd tea.Cmd
				m.response, cmd = m.response.Update(msg)
				return m, cmd
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
	panelHeight := m.height - 1
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

	hint := hintSt.Render("Tab · next panel   Ctrl+R · send   ↑↓ · select method / scroll   q · quit")

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, methodView, rightCol),
		hint,
	)
}

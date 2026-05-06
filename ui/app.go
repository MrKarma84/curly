package ui

import (
	"fmt"
	"time"

	"github.com/MrKarma84/curly/history"
	"github.com/MrKarma84/curly/httpclient"
	"github.com/MrKarma84/curly/ui/panels"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	panelMethod = iota
	panelURL
	panelHeaders
	panelBody
	panelResponse
	panelCount
)

const (
	methodWidth = 14
	urlHeight   = 3
)

var hintSt = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))

type ResponseMsg httpclient.Response
type InferSchemaMsg string

func doRequest(method, url, body string, headers map[string]string) tea.Cmd {
	return func() tea.Msg {
		return ResponseMsg(httpclient.Send(method, url, body, headers))
	}
}

func inferSchema(url string) tea.Cmd {
	return func() tea.Msg {
		resp := httpclient.Send("GET", url, "", nil)
		if resp.Err != "" || resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil
		}
		return InferSchemaMsg(resp.Body)
	}
}

type Model struct {
	width    int
	height   int
	focused  int
	method   panels.MethodPanel
	url      panels.URLPanel
	headers  panels.HeadersPanel
	body     panels.BodyPanel
	response panels.ResponsePanel

	store      *history.Store
	histIdx    int // -1 = live; ≥0 = browsing history (0 = most recent)
	liveMethod string
	liveURL    string
	liveHdrs   map[string]string
	liveBody   string
}

func New() Model {
	return Model{
		method:  panels.NewMethodPanel(),
		url:     panels.NewURLPanel(),
		headers: panels.NewHeadersPanel(),
		body:    panels.NewBodyPanel(),
		store:   history.New(),
		histIdx: -1,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) bodyActive() bool {
	return m.body.IsActive(m.method.Selected())
}

func (m Model) nextPanel() int {
	next := (m.focused + 1) % panelCount
	if next == panelBody && !m.bodyActive() {
		next = (next + 1) % panelCount
	}
	return next
}

func (m Model) prevPanel() int {
	prev := (m.focused - 1 + panelCount) % panelCount
	if prev == panelBody && !m.bodyActive() {
		prev = (prev - 1 + panelCount) % panelCount
	}
	return prev
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case ResponseMsg:
		m.response = m.response.SetResponse(httpclient.Response(msg))

	case InferSchemaMsg:
		m.body = m.body.InferFrom(string(msg))

	case tea.KeyMsg:
		// Let editing panels intercept Tab/Shift+Tab
		if (m.focused == panelHeaders && m.headers.IsEditing()) ||
			(m.focused == panelBody && m.body.IsEditing()) {
			return m.updateFocused(msg)
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "ctrl+r":
			url := m.url.Value()
			if url == "" {
				return m, nil
			}
			m.response = m.response.SetLoading()
			body := ""
			if m.bodyActive() {
				body = m.body.Body()
			}
			m.histIdx = -1
			m.store.Add(history.Entry{
				Timestamp: time.Now(),
				Method:    m.method.Selected(),
				URL:       url,
				Headers:   m.headers.Headers(),
				Body:      body,
			})
			return m, doRequest(m.method.Selected(), url, body, m.headers.Headers())

		case "alt+up":
			if m.store.Len() == 0 {
				return m, nil
			}
			if m.histIdx == -1 {
				m.liveMethod = m.method.Selected()
				m.liveURL = m.url.Value()
				m.liveHdrs = m.headers.Headers()
				m.liveBody = ""
				if m.bodyActive() {
					m.liveBody = m.body.Body()
				}
				m.histIdx = 0
			} else if m.histIdx < m.store.Len()-1 {
				m.histIdx++
			} else {
				return m, nil
			}
			return m.applyHistoryEntry(m.histIdx)

		case "alt+down":
			if m.histIdx == -1 {
				return m, nil
			}
			if m.histIdx > 0 {
				m.histIdx--
				return m.applyHistoryEntry(m.histIdx)
			}
			m.histIdx = -1
			m.method = m.method.SetSelected(m.liveMethod)
			m.url = m.url.SetValue(m.liveURL)
			m.headers = panels.NewHeadersPanelFrom(m.liveHdrs)
			m.body = panels.NewBodyPanel()
			if m.liveBody != "" {
				m.body = m.body.InferFrom(m.liveBody)
			}
			return m, nil

		case "i":
			if m.focused == panelBody && m.bodyActive() && m.url.Value() != "" {
				return m, inferSchema(m.url.Value())
			}

		case "tab":
			if m.focused == panelURL {
				m.url = m.url.Blur()
			}
			m.focused = m.nextPanel()
			if m.focused == panelURL {
				var cmd tea.Cmd
				m.url, cmd = m.url.Focus()
				return m, cmd
			}

		case "shift+tab":
			if m.focused == panelURL {
				m.url = m.url.Blur()
			}
			m.focused = m.prevPanel()
			if m.focused == panelURL {
				var cmd tea.Cmd
				m.url, cmd = m.url.Focus()
				return m, cmd
			}

		default:
			return m.updateFocused(msg)
		}
	}

	return m, nil
}

func (m Model) applyHistoryEntry(idx int) (tea.Model, tea.Cmd) {
	e := m.store.Get(idx)
	m.method = m.method.SetSelected(e.Method)
	m.url = m.url.SetValue(e.URL)
	m.headers = panels.NewHeadersPanelFrom(e.Headers)
	m.body = panels.NewBodyPanel()
	if e.Body != "" {
		m.body = m.body.InferFrom(e.Body)
	}
	return m, nil
}

func (m Model) updateFocused(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.focused {
	case panelMethod:
		m.method = m.method.Update(msg.(tea.KeyMsg))
	case panelURL:
		m.url, cmd = m.url.Update(msg)
	case panelHeaders:
		m.headers, cmd = m.headers.Update(msg)
	case panelBody:
		m.body, cmd = m.body.Update(msg)
	case panelResponse:
		m.response, cmd = m.response.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	rightWidth := m.width - methodWidth
	panelHeight := m.height - 1
	bottomHeight := panelHeight - urlHeight
	leftWidth := rightWidth / 2
	responseWidth := rightWidth - leftWidth
	headersHeight := bottomHeight / 2
	bodyHeight := bottomHeight - headersHeight

	methodView := m.method.View(methodWidth, panelHeight, m.focused == panelMethod)
	urlView := m.url.View(rightWidth, urlHeight, m.focused == panelURL)
	headersView := m.headers.View(leftWidth, headersHeight, m.focused == panelHeaders)
	bodyView := m.body.View(leftWidth, bodyHeight, m.focused == panelBody, m.bodyActive())
	responseView := m.response.View(responseWidth, bottomHeight, m.focused == panelResponse)

	leftCol := lipgloss.JoinVertical(lipgloss.Left, headersView, bodyView)
	rightCol := lipgloss.JoinVertical(lipgloss.Left,
		urlView,
		lipgloss.JoinHorizontal(lipgloss.Top, leftCol, responseView),
	)

	hintText := "Tab · next panel   Ctrl+R · send   Alt+↑/↓ · history   i · infer schema (body)   q · quit"
	if m.histIdx >= 0 {
		hintText = fmt.Sprintf("[history %d/%d]   Alt+↑ · older   Alt+↓ · newer / back to live", m.histIdx+1, m.store.Len())
	}
	hint := hintSt.Render(hintText)

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, methodView, rightCol),
		hint,
	)
}

package panels

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/MrKarma84/curly/httpclient"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	statusOKSt  = lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981")).Bold(true)
	statusErrSt = lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444")).Bold(true)
	durationSt  = lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF"))

	jsonKeySt   = lipgloss.NewStyle().Foreground(lipgloss.Color("#60A5FA"))
	jsonStrSt   = lipgloss.NewStyle().Foreground(lipgloss.Color("#34D399"))
	jsonNumSt   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FBBF24"))
	jsonBoolSt  = lipgloss.NewStyle().Foreground(lipgloss.Color("#F87171"))
	jsonPunctSt = lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF"))
)

type ResponsePanel struct {
	viewport viewport.Model
	status   string
	duration string
	loading  bool
	ready    bool
}

func (p ResponsePanel) SetLoading() ResponsePanel {
	p.loading = true
	p.ready = false
	p.status = ""
	p.viewport.SetContent(durationSt.Render("Sending request..."))
	return p
}

func (p ResponsePanel) SetResponse(resp httpclient.Response) ResponsePanel {
	p.loading = false
	p.ready = true

	if resp.Err != "" {
		p.status = statusErrSt.Render("Error")
		p.duration = ""
		p.viewport.SetContent(statusErrSt.Render(resp.Err))
		return p
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		p.status = statusOKSt.Render(resp.Status)
	} else {
		p.status = statusErrSt.Render(resp.Status)
	}
	p.duration = durationSt.Render(fmt.Sprintf("%dms", resp.Duration.Milliseconds()))

	p.viewport.SetContent(colorizeJSON(resp.Body))
	p.viewport.GotoTop()
	return p
}

func (p ResponsePanel) Update(msg tea.Msg) (ResponsePanel, tea.Cmd) {
	var cmd tea.Cmd
	p.viewport, cmd = p.viewport.Update(msg)
	return p, cmd
}

func (p ResponsePanel) View(width, height int, focused bool) string {
	innerW := width - 2
	innerH := height - 2

	p.viewport.Width = innerW
	p.viewport.Height = max(0, innerH-2) // -2 for header + gap

	var header string
	if p.status != "" {
		header = lipgloss.JoinHorizontal(lipgloss.Left, p.status, "  ", p.duration)
	} else {
		header = labelSt.Render("RESPONSE")
	}

	var body string
	if !p.ready && !p.loading {
		body = durationSt.Render("Press Ctrl+R to send a request...")
	} else {
		body = p.viewport.View()
	}

	content := lipgloss.JoinVertical(lipgloss.Left, header, body)
	return box(width, height, focused).Render(content)
}

func colorizeJSON(raw string) string {
	var obj any
	if err := json.Unmarshal([]byte(raw), &obj); err != nil {
		return raw
	}
	pretty, _ := json.MarshalIndent(obj, "", "  ")
	lines := strings.Split(string(pretty), "\n")
	result := make([]string, len(lines))
	for i, line := range lines {
		result[i] = colorizeLine(line)
	}
	return strings.Join(result, "\n")
}

func colorizeLine(line string) string {
	trimmed := strings.TrimLeft(line, " \t")
	indent := line[:len(line)-len(trimmed)]

	if trimmed == "{" || trimmed == "}" || trimmed == "}," ||
		trimmed == "[" || trimmed == "]" || trimmed == "]," {
		return indent + jsonPunctSt.Render(trimmed)
	}

	if len(trimmed) > 0 && trimmed[0] == '"' {
		if idx := strings.Index(trimmed, `": `); idx != -1 {
			key := trimmed[:idx+1]
			rest := trimmed[idx+3:]
			trailing := ""
			if strings.HasSuffix(rest, ",") {
				trailing = ","
				rest = rest[:len(rest)-1]
			}
			return indent +
				jsonKeySt.Render(key) + jsonPunctSt.Render(": ") +
				colorizeValue(rest) + jsonPunctSt.Render(trailing)
		}
	}

	trailing := ""
	if strings.HasSuffix(trimmed, ",") {
		trailing = ","
		trimmed = trimmed[:len(trimmed)-1]
	}
	return indent + colorizeValue(trimmed) + jsonPunctSt.Render(trailing)
}

func colorizeValue(v string) string {
	switch {
	case v == "true" || v == "false" || v == "null":
		return jsonBoolSt.Render(v)
	case strings.HasPrefix(v, `"`):
		return jsonStrSt.Render(v)
	case len(v) > 0 && (v[0] == '{' || v[0] == '['):
		return jsonPunctSt.Render(v)
	default:
		if _, err := strconv.ParseFloat(v, 64); err == nil {
			return jsonNumSt.Render(v)
		}
		return v
	}
}

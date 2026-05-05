package panels

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	headerKeySt = lipgloss.NewStyle().Foreground(lipgloss.Color("#60A5FA"))
	headerValSt = lipgloss.NewStyle().Foreground(lipgloss.Color("#D1D5DB"))
	selectedSt  = lipgloss.NewStyle().Background(lipgloss.Color("#1F2937")).Foreground(lipgloss.Color("#F9FAFB"))
	dimSt       = lipgloss.NewStyle().Foreground(lipgloss.Color("#4B5563"))
)

type Header struct {
	Key   string
	Value string
}

type HeadersPanel struct {
	headers    []Header
	selected   int
	editing    bool
	editingVal bool
	isNew      bool
	keyInput   textinput.Model
	valInput   textinput.Model
}

func NewHeadersPanel() HeadersPanel {
	ki := textinput.New()
	ki.Placeholder = "Header-Name"
	ki.CharLimit = 100

	vi := textinput.New()
	vi.Placeholder = "value"
	vi.CharLimit = 500

	return HeadersPanel{keyInput: ki, valInput: vi}
}

func (p HeadersPanel) IsEditing() bool {
	return p.editing
}

func (p HeadersPanel) Headers() map[string]string {
	h := make(map[string]string, len(p.headers))
	for _, hdr := range p.headers {
		if hdr.Key != "" {
			h[hdr.Key] = hdr.Value
		}
	}
	return h
}

func (p HeadersPanel) startEditing() (HeadersPanel, tea.Cmd) {
	h := p.headers[p.selected]
	p.keyInput.SetValue(h.Key)
	p.valInput.SetValue(h.Value)
	p.editing = true
	p.editingVal = false
	p.valInput.Blur()
	return p, p.keyInput.Focus()
}

func (p HeadersPanel) Update(msg tea.Msg) (HeadersPanel, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return p, nil
	}

	if p.editing {
		switch key.String() {
		case "tab", "shift+tab":
			p.editingVal = !p.editingVal
			if p.editingVal {
				p.keyInput.Blur()
				return p, p.valInput.Focus()
			}
			p.valInput.Blur()
			return p, p.keyInput.Focus()

		case "enter":
			p.headers[p.selected] = Header{
				Key:   p.keyInput.Value(),
				Value: p.valInput.Value(),
			}
			p.editing = false
			p.keyInput.Blur()
			p.valInput.Blur()
			return p, nil

		case "esc":
			if p.isNew {
				p.headers = append(p.headers[:p.selected], p.headers[p.selected+1:]...)
				if p.selected >= len(p.headers) && p.selected > 0 {
					p.selected--
				}
			}
			p.editing = false
			p.keyInput.Blur()
			p.valInput.Blur()
			return p, nil

		default:
			var cmd tea.Cmd
			if p.editingVal {
				p.valInput, cmd = p.valInput.Update(msg)
			} else {
				p.keyInput, cmd = p.keyInput.Update(msg)
			}
			return p, cmd
		}
	}

	// Normal mode
	switch key.String() {
	case "up", "k":
		if p.selected > 0 {
			p.selected--
		}
	case "down", "j":
		if p.selected < len(p.headers)-1 {
			p.selected++
		}
	case "a":
		p.headers = append(p.headers, Header{})
		p.selected = len(p.headers) - 1
		p.isNew = true
		return p.startEditing()
	case "enter":
		if len(p.headers) > 0 {
			p.isNew = false
			return p.startEditing()
		}
	case "d", "x":
		if len(p.headers) > 0 {
			p.headers = append(p.headers[:p.selected], p.headers[p.selected+1:]...)
			if p.selected >= len(p.headers) && p.selected > 0 {
				p.selected--
			}
		}
	}

	return p, nil
}

func (p HeadersPanel) View(width, height int, focused bool) string {
	innerW := width - 2
	var rows []string

	rows = append(rows, labelSt.Render("HEADERS"))

	if len(p.headers) == 0 {
		rows = append(rows, dimSt.Render("No headers"))
	} else {
		for i, h := range p.headers {
			row := headerKeySt.Render(h.Key) + ": " + headerValSt.Render(h.Value)
			if i == p.selected && focused && !p.editing {
				row = selectedSt.Width(innerW).Render(row)
			}
			rows = append(rows, row)
		}
	}

	if p.editing {
		p.keyInput.Width = innerW - 8
		p.valInput.Width = innerW - 8
		rows = append(rows, "")
		rows = append(rows, "Key:   "+p.keyInput.View())
		rows = append(rows, "Value: "+p.valInput.View())
		rows = append(rows, dimSt.Render("Tab·switch  Enter·save  Esc·cancel"))
	} else if focused {
		rows = append(rows, "")
		rows = append(rows, dimSt.Render("a·add  d·del  Enter·edit"))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return box(width, height, focused).Render(content)
}

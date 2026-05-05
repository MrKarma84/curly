package panels

import (
	"encoding/json"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var bodyMethods = map[string]bool{"POST": true, "PUT": true, "PATCH": true}

type BodyField struct {
	Key   string
	Value string
}

type BodyPanel struct {
	fields     []BodyField
	selected   int
	editing    bool
	editingVal bool
	isNew      bool
	keyInput   textinput.Model
	valInput   textinput.Model
}

func NewBodyPanel() BodyPanel {
	ki := textinput.New()
	ki.Placeholder = "field"
	ki.CharLimit = 100

	vi := textinput.New()
	vi.Placeholder = "value"
	vi.CharLimit = 500

	return BodyPanel{keyInput: ki, valInput: vi}
}

func (p BodyPanel) IsActive(method string) bool {
	return bodyMethods[method]
}

func (p BodyPanel) IsEditing() bool {
	return p.editing
}

// Body returns the JSON representation of all fields.
func (p BodyPanel) Body() string {
	if len(p.fields) == 0 {
		return ""
	}
	obj := make(map[string]any, len(p.fields))
	for _, f := range p.fields {
		if f.Key != "" {
			obj[f.Key] = parseBodyValue(f.Value)
		}
	}
	b, _ := json.Marshal(obj)
	return string(b)
}

// InferFrom populates fields from the top-level keys of a JSON object.
func (p BodyPanel) InferFrom(jsonStr string) BodyPanel {
	var obj map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
		return p
	}
	p.fields = nil
	for k, v := range obj {
		p.fields = append(p.fields, BodyField{
			Key:   k,
			Value: formatBodyValue(v),
		})
	}
	p.selected = 0
	return p
}

func (p BodyPanel) startEditing() (BodyPanel, tea.Cmd) {
	h := p.fields[p.selected]
	p.keyInput.SetValue(h.Key)
	p.valInput.SetValue(h.Value)
	p.editing = true
	p.editingVal = false
	p.valInput.Blur()
	return p, p.keyInput.Focus()
}

func (p BodyPanel) Update(msg tea.Msg) (BodyPanel, tea.Cmd) {
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
			p.fields[p.selected] = BodyField{
				Key:   p.keyInput.Value(),
				Value: p.valInput.Value(),
			}
			p.editing = false
			p.keyInput.Blur()
			p.valInput.Blur()
			return p, nil

		case "esc":
			if p.isNew {
				p.fields = append(p.fields[:p.selected], p.fields[p.selected+1:]...)
				if p.selected >= len(p.fields) && p.selected > 0 {
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

	switch key.String() {
	case "up", "k":
		if p.selected > 0 {
			p.selected--
		}
	case "down", "j":
		if p.selected < len(p.fields)-1 {
			p.selected++
		}
	case "a":
		p.fields = append(p.fields, BodyField{})
		p.selected = len(p.fields) - 1
		p.isNew = true
		return p.startEditing()
	case "enter":
		if len(p.fields) > 0 {
			p.isNew = false
			return p.startEditing()
		}
	case "d", "x":
		if len(p.fields) > 0 {
			p.fields = append(p.fields[:p.selected], p.fields[p.selected+1:]...)
			if p.selected >= len(p.fields) && p.selected > 0 {
				p.selected--
			}
		}
	}

	return p, nil
}

func (p BodyPanel) View(width, height int, focused, active bool) string {
	innerW := width - 2
	var rows []string

	if !active {
		rows = append(rows, labelSt.Render("BODY"))
		rows = append(rows, dimSt.Render("Not available for this method"))
		content := lipgloss.JoinVertical(lipgloss.Left, rows...)
		return box(width, height, false).Render(content)
	}

	rows = append(rows, labelSt.Render("BODY"))

	if len(p.fields) == 0 {
		rows = append(rows, dimSt.Render("No fields"))
	} else {
		for i, f := range p.fields {
			typ := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render(valueType(f.Value))
			row := headerKeySt.Render(f.Key) + ": " + headerValSt.Render(f.Value) + " " + typ
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
		rows = append(rows, dimSt.Render("a·add  d·del  Enter·edit  i·infer schema"))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return box(width, height, focused).Render(content)
}

func parseBodyValue(v string) any {
	if v == "true" {
		return true
	}
	if v == "false" {
		return false
	}
	if v == "null" {
		return nil
	}
	if n, err := strconv.ParseFloat(v, 64); err == nil {
		return n
	}
	return v
}

func formatBodyValue(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		if val {
			return "true"
		}
		return "false"
	case nil:
		return "null"
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

func valueType(v string) string {
	if v == "true" || v == "false" {
		return "bool"
	}
	if v == "null" {
		return "null"
	}
	if _, err := strconv.ParseFloat(v, 64); err == nil {
		return "number"
	}
	return "string"
}

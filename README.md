# curly

A keyboard-driven TUI HTTP client — a terminal alternative to Insomnia/Postman.

```
╭────────╮╭──────────────────────────────────────────╮
│ METHOD ││ URL                                      │
│        │╰──────────────────────────────────────────╯
│ ► GET  │╭─────────────────────╮╭───────────────────╮
│   POST ││ HEADERS             ││ RESPONSE          │
│   PUT  ││                     ││                   │
│   PATCH│╰─────────────────────╯╰───────────────────╯
│ DELETE │
╰────────╯
```

**Key differentiators:** request chaining · replay diff · Postman/Insomnia/Bruno import

---

## Quick start

**Requirements:** [Go 1.21+](https://go.dev/dl/)

```bash
git clone https://github.com/MrKarma84/curly.git
cd curly
go run .
```

---

## Keyboard shortcuts

| Key | Action |
|-----|--------|
| `Tab` / `Shift+Tab` | Move between panels |
| `↑` `↓` | Navigate lists |
| `Enter` | Select / confirm |
| `Ctrl+R` | Send request |
| `Ctrl+S` | Save to collection |
| `Ctrl+N` | New request |
| `Ctrl+P` | History — previous request |
| `Ctrl+D` | Replay diff |
| `Ctrl+W` | Watch mode |
| `Ctrl+L` | Chain request |
| `?` | Help |
| `q` / `Ctrl+C` | Quit |

---

## Project structure

```
curly/
├── main.go               # Entry point — starts the TUI program
├── ui/
│   ├── app.go            # Main Bubble Tea model (the "brain" of the UI)
│   └── panels/
│       ├── panel.go      # Shared styles and helpers for all panels
│       ├── method.go     # HTTP method selector (GET, POST, PUT…)
│       ├── url.go        # URL input field
│       ├── headers.go    # Request headers editor
│       └── response.go   # Response display
├── go.mod                # Go module definition (like package.json in Node)
└── go.sum                # Dependency checksums (auto-generated, don't edit)
```

---

## Go concepts — explained for beginners

This section explains the Go concepts introduced at each step of the project.
If you're new to Go, read this alongside the code.

### Packages

Go organizes code into **packages**. Every `.go` file starts with `package <name>`.

```go
package main   // the entry point — Go looks for this to run the program
package ui     // the ui package — groups all UI-related code
package panels // the panels sub-package
```

To use code from another package, you **import** it:

```go
import "github.com/MrKarma84/curly/ui/panels"

// now you can use panels.MethodPanel, panels.URLPanel, etc.
```

### Structs

A **struct** is a group of related fields — similar to a class in Python/JS,
but without inheritance.

```go
type Model struct {
    width   int    // terminal width in characters
    height  int    // terminal height in characters
    focused int    // index of the currently active panel (0, 1, 2, 3)
}
```

You create a struct with:
```go
m := Model{width: 80, height: 24, focused: 0}
// or using the New() constructor:
m := ui.New()
```

### Methods on structs

In Go, you attach functions to structs using a **receiver**:

```go
//           ↓ receiver: "this method belongs to Model"
func (m Model) View() string {
    // m is the struct instance, like `self` in Python
    return "hello"
}
```

### The Bubble Tea pattern (Init / Update / View)

Bubble Tea uses the **Elm architecture** — a simple loop:

```
User input → Update() → new state → View() → rendered screen
                ↑                                    |
                └────────────────────────────────────┘
```

Every Bubble Tea model must implement 3 methods:

```go
// Init — runs once at startup, returns an optional command
func (m Model) Init() tea.Cmd { return nil }

// Update — receives a message (key press, window resize…) and returns
//           a new model + an optional command to run next
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { ... }

// View — converts the current state into a string to display
func (m Model) View() string { ... }
```

### iota — auto-incrementing constants

`iota` generates sequential integers automatically:

```go
const (
    panelMethod  = iota // = 0
    panelURL            // = 1
    panelHeaders        // = 2
    panelResponse       // = 3
    panelCount          // = 4  ← used for the Tab cycle modulo
)
```

Tab cycling uses modulo `%` to wrap around:
```go
m.focused = (m.focused + 1) % panelCount
// 0 → 1 → 2 → 3 → 0 → 1 → …
```

### Goroutines & async commands (Bubble Tea)

In a TUI, you can't block the UI thread to wait for an HTTP response.
Bubble Tea solves this with **commands** (`tea.Cmd`) — functions that run in the background
and send a message when done:

```go
// A Cmd is just a function that returns a Msg
func doRequest(method, url string) tea.Cmd {
    return func() tea.Msg {           // ← runs in a goroutine automatically
        resp := httpclient.Send(...)  // blocks here, but UI stays responsive
        return ResponseMsg(resp)      // sends the result back to Update()
    }
}

// Triggered by Ctrl+R in Update():
return m, doRequest(m.method.Selected(), m.url.Value())

// When the goroutine finishes, Update() receives:
case ResponseMsg:
    m.response = m.response.SetResponse(httpclient.Response(msg))
```

### Type aliases

`type ResponseMsg httpclient.Response` creates a **distinct type** from `httpclient.Response`.
This lets Bubble Tea route the message to the right `case` in `Update()`:

```go
case ResponseMsg:   // ✅ only matches ResponseMsg
    // handle it

case httpclient.Response:  // ← would never match — different type
```

### defer

`defer` schedules a function call to run when the enclosing function returns,
no matter what (even on error). Used to close resources:

```go
resp, err := client.Do(req)
defer resp.Body.Close()  // always runs when Send() returns
                         // without this, the connection leaks
```

### Slices

A **slice** is an ordered, resizable list — the most common collection in Go:

```go
methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
//          ↑ type: slice of strings

methods[0]        // "GET"   — access by index
len(methods)      // 5       — number of elements

// loop over a slice with range:
for i, method := range methods {
    // i = index (0, 1, 2…), method = value ("GET", "POST"…)
}
```

### Maps

A **map** is a key → value store, like a dictionary in Python or an object in JS:

```go
var methodColors = map[string]lipgloss.Color{
//                     ↑ key type   ↑ value type
    "GET":    "#10B981",  // green
    "POST":   "#F59E0B",  // yellow
    "DELETE": "#EF4444",  // red
}

color := methodColors["GET"]  // look up a value by key
```

### Immutability in Bubble Tea

In Bubble Tea, `Update` always returns a **new** model instead of modifying the current one.
This is why `MethodPanel.Update` returns a `MethodPanel`:

```go
// ✅ correct Bubble Tea style — return a new copy
func (p MethodPanel) Update(msg tea.KeyMsg) MethodPanel {
    p.selected = newIndex   // modifies the copy, not the original
    return p
}

// ❌ would not work — Bubble Tea models are values, not pointers
func (p *MethodPanel) Update(msg tea.KeyMsg) {
    p.selected = newIndex
}
```

### Lip Gloss — terminal styling

Lip Gloss lets you style terminal output like CSS:

```go
style := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).   // rounded box border
    BorderForeground(lipgloss.Color("#7C3AED")). // purple
    Width(40).                          // inner content width
    Height(10)                          // inner content height

output := style.Render("hello")        // returns a styled string
```

---

## Development roadmap

| Step | Feature | Status |
|------|---------|--------|
| 1 | Scaffolding & Hello World TUI | ✅ done |
| 2 | Basic layout + panel navigation | ✅ done |
| 3 | HTTP method selector | ✅ done |
| 4 | URL input + send GET request | ✅ done |
| 5 | Headers editor | 🔜 next |
| 6 | Body + schema detection | ⬜ |
| 7 | Navigable history | ⬜ |
| 8 | Replay & diff | ⬜ |
| 9 | Collections | ⬜ |
| 10 | Environment variables | ⬜ |
| 11 | Request chaining | ⬜ |
| 12 | Watch mode | ⬜ |
| 13 | Postman / Insomnia / Bruno import | ⬜ |
| 14 | Polish & release | ⬜ |

---

## License

MIT — see [LICENSE](LICENSE)

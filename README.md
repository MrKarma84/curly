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
| `Alt+↑` | History — go to older request |
| `Alt+↓` | History — go to newer request / back to live |
| `i` | Infer body schema from URL (in BODY panel) |
| `Ctrl+S` | Save to collection _(coming soon)_ |
| `Ctrl+D` | Replay diff _(coming soon)_ |
| `Ctrl+W` | Watch mode _(coming soon)_ |
| `Ctrl+L` | Chain request _(coming soon)_ |
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
│       ├── body.go       # Request body editor with schema inference
│       └── response.go   # Response display
├── history/
│   └── store.go          # Request history (persisted to ~/.curly/history.json)
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

### Goroutines & async commands (Bubble Tea)

In a TUI, you can't block the UI thread to wait for an HTTP response.
Bubble Tea solves this with **commands** (`tea.Cmd`) — functions that run in the background
and send a message when done:

```go
func doRequest(method, url string, headers map[string]string) tea.Cmd {
    return func() tea.Msg {           // ← runs in a goroutine automatically
        resp := httpclient.Send(...)  // blocks here, but UI stays responsive
        return ResponseMsg(resp)      // sends the result back to Update()
    }
}

// When the goroutine finishes, Update() receives:
case ResponseMsg:
    m.response = m.response.SetResponse(httpclient.Response(msg))
```

### defer

`defer` schedules a call to run when the enclosing function returns, even on error:

```go
resp, err := client.Do(req)
defer resp.Body.Close()  // always runs — without this, the connection leaks
```

### Type assertion

Extracting a concrete type from an interface:

```go
key, ok := msg.(tea.KeyMsg)  // "is msg a tea.KeyMsg?"
if !ok {
    return p, nil  // nope — ignore it
}
// key is now usable as tea.KeyMsg
```

The two-value form (`value, ok`) never panics. The single-value form panics if the type doesn't match.

### Slice manipulation

Removing an element at index `i` — the standard Go pattern:

```go
s = append(s[:i], s[i+1:]...)
// s[:i]   → everything before i
// s[i+1:] → everything after i
// ...     → unpacks the slice as individual arguments to append()
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

### Reading and writing files — `os` and `filepath`

The history is persisted in `~/.curly/history.json` between sessions.
`os` is Go's standard library for interacting with the operating system:

```go
// Find the user's home directory (~)
home, err := os.UserHomeDir()   // "/home/alice"

// Create all missing directories in a path (like mkdir -p)
os.MkdirAll("/home/alice/.curly", 0o700)
//                                 ↑ Unix permissions: owner read+write+exec only

// Write a file (creates or overwrites)
os.WriteFile(path, data, 0o600)
//                         ↑ owner read+write only (no exec — it's data, not a script)

// Read a file into a []byte slice
data, err := os.ReadFile(path)

// Check for a specific error — "does the file exist?"
if os.IsNotExist(err) {
    // first run — no history yet, start fresh
}
```

`filepath.Join` builds paths safely for the current OS (handles `/` vs `\`):

```go
path := filepath.Join(home, ".curly", "history.json")
// → "/home/alice/.curly/history.json"  on Linux/macOS
// → "C:\Users\Alice\.curly\history.json"  on Windows
```

### JSON serialisation — `encoding/json`

Go can convert any struct to JSON and back automatically, using **struct tags**:

```go
type Entry struct {
    Timestamp time.Time         `json:"timestamp"`
    Method    string            `json:"method"`
    URL       string            `json:"url"`
    Headers   map[string]string `json:"headers,omitempty"`
    //                                          ↑ omit this field if the map is empty
    Body      string            `json:"body,omitempty"`
}
```

The backtick strings after field types are **struct tags** — metadata read at runtime by `encoding/json`:

```go
// Struct → JSON bytes
data, err := json.MarshalIndent(store, "", "  ")
// {"entries": [{"timestamp": "2026-05-06T…", "method": "GET", …}]}

// JSON bytes → Struct (fills in the fields automatically)
err := json.Unmarshal(data, &store)
```

### Pointers — when you need to share state

Until now, every panel was a **value** — passed by copy in Bubble Tea.
The history `Store` is different: it must be **shared** across all copies of `Model`
(otherwise saving in one copy wouldn't be visible in the next).

```go
type Model struct {
    // ...
    store *history.Store  // ← pointer: all copies of Model point to the same Store
}

// Calling Add() through a pointer modifies the original:
m.store.Add(entry)  // ✅ the real store is updated
```

Rule of thumb: use a pointer (`*T`) when you need mutations to be visible everywhere,
use a value (`T`) when you want an independent copy.

---

## Development roadmap

| Step | Feature | Status |
|------|---------|--------|
| 1 | Scaffolding & Hello World TUI | ✅ done |
| 2 | Basic layout + panel navigation | ✅ done |
| 3 | HTTP method selector | ✅ done |
| 4 | URL input + send GET request | ✅ done |
| 5 | Headers editor | ✅ done |
| 6 | Body + schema detection | ✅ done |
| 7 | Navigable history | ✅ done |
| 8 | Replay & diff | 🔜 next |
| 9 | Collections | ⬜ |
| 10 | Environment variables | ⬜ |
| 11 | Request chaining | ⬜ |
| 12 | Watch mode | ⬜ |
| 13 | Postman / Insomnia / Bruno import | ⬜ |
| 14 | Polish & release | ⬜ |

---

## License

MIT — see [LICENSE](LICENSE)

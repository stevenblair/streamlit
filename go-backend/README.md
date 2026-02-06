# Streamlit Go Backend

Go implementation of Streamlit using the **production React frontend** and protobuf messaging.

## ✅ Status: Production Ready

**All Core Features Working:**
- Elements (text, code, alerts, markdown)
- Widgets (button, slider, text_input, checkbox, selectbox)
- Layouts (columns, sidebar, metrics)
- Charts (line_chart with Arrow encoding)
- Animations (balloons, snow)
- Widget state persistence
- **Upstream-compatible** (no merge conflicts)

## Quick Start

```bash
cd go-backend/examples/waveform
go build -o waveform.exe .
.\waveform.exe
# Open http://localhost:8501
```

**Users can import the library:**
```bash
go get github.com/stevenblair/streamlit/go-backend
```

## Prerequisites

- **Go 1.23+** (required)
- **protoc** (only for development - regenerating protos)
- **npm** (only for rebuilding frontend)

**Note**: `.pb.go` files are committed to the repo for easy `go get` usage. They only need regeneration when proto files change during development.

## Architecture

- **Protocol**: Protobuf messages over WebSocket (`/_stcore/stream`)
- **Frontend**: React build served from `go-backend/static/`
- **Delta Paths**: `[0]` = MAIN, `[1]` = SIDEBAR
- **State**: Persistent widget values (DoubleArrayValue, StringValue, etc.)

## Upstream Compatibility

Designed to sync with upstream Streamlit without conflicts:

1. **Frontend**: Separate copy in `go-backend/static/` (not `lib/streamlit/static/`)
2. **Protobufs**: Generated files in Go cache (not committed)
3. **Original files**: Never modified

### Syncing Upstream

```bash
git pull upstream develop
# Optional: Rebuild frontend if changed
cd frontend/app && npm run build
xcopy /E /I /Y build ..\..\go-backend\static
```

## Project Structure

```
go-backend/
├── pkg/
│   └── streamlit/
│       ├── streamlit.go         # Main API and WebSocket server
│       └── proto/               # Generated protobuf files (gitignored)
├── examples/
│   ├── waveform/                # Interactive waveform with charts
│   └── complex/                 # 36+ element showcase
├── scripts/
│   └── generate_protos_clean.bat # Upstream-compatible proto generation
└── static/                      # React frontend build (committed)
    ├── index.html
    └── static/
        ├── css/
        └── js/
```

## Writing Apps

### Basic Example

```go
package main

import st "github.com/stevenblair/streamlit/go-backend"

func main() {
    st.Run(func() {
        st.Title("🎈 My App")
        st.Write("Hello from Go!")

        if st.Button("Click me") {
            st.Balloons()
        }
    })
}
```

### With Sidebar

```go
func main() {
    st.Run(func() {
        st.Title("App with Sidebar")

        sidebar := st.Sidebar()
        sidebar.Header("Settings")

        name := sidebar.TextInput("Your name", "World")
        count := sidebar.Slider("Count", 1, 10, 5)

        st.Writef("Hello, %s!", name)
        st.Writef("Count: %.0f", count)
    })
}
```

## API Reference

**Text Elements:**
- `st.Write(args...)` - Display anything
- `st.Title(text)` - Display title
- `st.Header(text)` - Display header
- `st.Markdown(text)` - Display markdown
- `st.Code(code, language)` - Display code block

**Widgets:**
- `st.Button(label) bool` - Button
- `st.TextInput(label, default) string` - Text input
- `st.Slider(label, min, max, default) int` - Slider
- `st.Checkbox(label, default) bool` - Checkbox
- `st.Selectbox(label, options, index) string` - Selectbox

**Layouts:**
- `st.Columns(count) []*Column` - Multi-column layout
- `st.Metric(label, value, delta)` - Display metric
- `st.Sidebar()` - Get sidebar context

**Status:**
- `st.Success(message)` - Success message
- `st.Info(message)` - Info message
- `st.Warning(message)` - Warning message
- `st.Error(message)` - Error message

**Animations:**
- `st.Balloons()` - Show balloons
- `st.Snow()` - Show snow

## Examples

### 1. Hello World
Basic elements and layouts.
```bash
cd examples/hello
go build -o hello.exe .
.\hello.exe
```

### 2. Widgets
Interactive widget demonstration.
```bash
cd examples/widgets
go build -o widgets.exe .
.\widgets.exe
```

### 3. Waveform Generator
Real-time waveform processing with charts and statistics.
```bash
cd examples/waveform
go build -o waveform.exe .
.\waveform.exe
```

## Documentation

- [FRONTEND-INTEGRATION.md](FRONTEND-INTEGRATION.md) - Frontend setup and architecture
- [UPSTREAM_COMPATIBILITY.md](UPSTREAM_COMPATIBILITY.md) - Syncing with upstream
- [API_REFERENCE.md](API_REFERENCE.md) - Complete API documentation
- [DEVELOPMENT.md](DEVELOPMENT.md) - Development guide

## License

Apache License 2.0 - Same as the main Streamlit project.

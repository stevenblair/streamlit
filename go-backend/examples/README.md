# Examples Directory

This directory contains example Streamlit Go applications.

Each example is in its own subdirectory to keep them organized and avoid IDE conflicts.

## Directory Structure

```
examples/
├── hello/          # Basic hello world example
│   └── main.go
├── widgets/        # Interactive widgets demo
│   └── main.go
└── README.md
```

## Building Examples

```bash
# Using build scripts (recommended)
./scripts/build.sh      # Unix/macOS
scripts\build.bat       # Windows

# Or manually
go build -o bin/hello ./examples/hello
go build -o bin/widgets ./examples/widgets
```

## Running Examples

```bash
# Start the server first
./bin/streamlit-go &

# Then run any example
./bin/hello
./bin/widgets
```

## Available Examples

### hello/
Basic "Hello World" app demonstrating:
- Text elements (Title, Header, Write)
- Markdown and Code blocks
- Status messages
- Layouts (Columns, Metrics)
- Animations (Balloons)

### widgets/
Interactive widgets demonstration:
- TextInput
- Slider
- Selectbox
- Checkbox
- Button
- Metrics

## Creating Your Own App

Create a new .go file:

```go
package main

import st "github.com/stevenblair/streamlit/go-backend/pkg/streamlit"

func main() {
    st.Run(app)
}

func app() {
    st.Title("My App")
    st.Write("Hello, World!")
}
```

Build and run:
```bash
go build -o myapp myapp.go
./myapp
```

## Note

Each example is in its own directory so VS Code and other IDEs can analyze them as separate packages without conflicts.

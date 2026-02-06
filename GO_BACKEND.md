# Streamlit Go Backend

The `go-backend/` directory contains a native Go implementation of Streamlit. **Write Streamlit apps in Go** instead of Python!

## Quick Links

📖 **[Full Documentation](go-backend/INDEX.md)** - Complete overview
🚀 **[Quick Start Guide](go-backend/QUICKSTART.md)** - Get running fast
📚 **[API Reference](go-backend/API_REFERENCE.md)** - Complete API docs
🔗 **[Integration Guide](go-backend/INTEGRATION.md)** - How it fits in

## What is it?

A native Go API for Streamlit that:
- ✅ Write apps in Go (not Python!)
- ✅ Use familiar API: `st.Write()`, `st.Button()`, etc.
- ✅ Compile to single binary
- ✅ Type-safe development
- ✅ Reuses same frontend
- ⚠️ Experimental (not production-ready)

## Quick Start

```bash
# Navigate to go-backend
cd go-backend

# Setup
./scripts/setup.sh && ./scripts/proto.sh

# Build
./scripts/build.sh

# Start server
./bin/streamlit-go &

# Run example
./bin/hello

# Open browser to http://localhost:8501
```

## Example App

```go
package main

import st "github.com/stevenblair/streamlit/go-backend/pkg/streamlit"

func main() {
    st.Run(app)
}

func app() {
    st.Title("🎈 Hello Streamlit")

    name := st.TextInput("Your name", "World")
    st.Writef("Hello, %s!", name)

    if st.Button("Celebrate") {
        st.Balloons()
    }
}
```

Build and run:
```bash
go build -o myapp myapp.go
./myapp  # Server must be running
```

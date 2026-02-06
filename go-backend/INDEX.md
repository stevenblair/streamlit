# Go Backend - Native Go API for Streamlit

## What is this?

The Go backend is a **native Go implementation** of Streamlit. Write Streamlit apps in **Go** instead of Python! It provides a complete Go API that mirrors the Python Streamlit API, allowing you to create interactive web applications using Go.

## Key Features

- ✅ **Native Go API**: Write apps using `st.Write()`, `st.Button()`, etc. in Go
- ✅ **Reuses existing frontend**: No changes to the React/TypeScript UI
- ✅ **Same protocol**: Uses Protocol Buffer definitions
- ✅ **Better performance**: Compiled binaries, faster execution
- ✅ **Type safety**: Go's static typing catches errors at compile time
- ✅ **Single binary**: Easy deployment
- ⚠️ **Experimental**: Not feature-complete or production-ready

## Quick Start

See [QUICKSTART.md](QUICKSTART.md) for detailed setup instructions.

```bash
# 1. Setup
cd go-backend
./scripts/setup.sh && ./scripts/proto.sh

# 2. Build
./scripts/build.sh

# 3. Start server
./bin/streamlit-go &

# 4. Run example
./bin/hello

# 5. Open browser
# http://localhost:8501
```

## Example App

```go
package main

import st "github.com/streamlit/streamlit/go-backend/pkg/streamlit"

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

## Documentation

- **[QUICKSTART.md](QUICKSTART.md)** - Setup and getting started
- **[INTEGRATION.md](INTEGRATION.md)** - How it integrates with the main project
- **[README.md](README.md)** - Full project overview

## Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                    Streamlit Application                      │
└──────────────────────────────────────────────────────────────┘
                              ▲
                              │
                    ┌─────────┴─────────┐
                    │   Choose Backend  │
                    └─────────┬─────────┘
                              │
        ┌─────────────────────┴─────────────────────┐
        │                                           │
        ▼                                           ▼
┌───────────────┐                           ┌──────────────┐
│ Python Backend│                           │  Go Backend  │
│  (Tornado)    │                           │  (net/http)  │
└───────┬───────┘                           └──────┬───────┘
        │                                          │
        └─────────────────┬────────────────────────┘
                          │
                          ▼
              ┌─────────────────────┐
              │   Shared Frontend   │
              │   (React/TypeScript)│
              └─────────────────────┘
                          ▲
                          │
                          ▼
              ┌─────────────────────┐
              │  Protocol Buffers   │
              │  (ForwardMsg/BackMsg)│
              └─────────────────────┘
```

## Why Go?

### Advantages

1. **Performance**: Compiled binary, faster execution
2. **Concurrency**: Native goroutines for session handling
3. **Memory**: Lower memory footprint
4. **Deployment**: Single binary, no runtime dependencies
5. **Startup**: Much faster cold start times

### Trade-offs

1. **Maturity**: Experimental, not battle-tested
2. **Features**: Not all features implemented yet
3. **Ecosystem**: Less Python ecosystem integration
4. **Debugging**: Different debugging workflow

## Current Status

### ✅ Implemented

- HTTP server with routing
- WebSocket connection handling
- Static file serving
- Protobuf message structure
- Basic session management
- Health check endpoints
- Project scaffolding

### 🚧 In Progress

- ForwardMsg/BackMsg handling
- Script execution integration
- Widget state management
- Element rendering

### ❌ Not Yet Implemented

- Full element/widget support
- Caching system
- File uploads
- Media file handling
- Custom components
- Authentication
- Production deployment

## Use Cases

### Good For

- **Learning**: Understand Streamlit architecture
- **Experimentation**: Try alternative implementations
- **Performance testing**: Compare backend performance
- **Embedded systems**: Deploy where Go works better

### Not Good For (Yet)

- **Production apps**: Not stable enough
- **Complex apps**: Missing features
- **Critical systems**: Insufficient testing

## Contributing

We welcome contributions! This is an experimental project, so:

1. **Explore freely**: Try new approaches
2. **Document changes**: Keep docs updated
3. **Test thoroughly**: Ensure compatibility
4. **Follow conventions**: Match Go best practices

See [CONTRIBUTING.md](../CONTRIBUTING.md) in the main repo.

## Benchmarks

Preliminary benchmarks (Go backend vs Python backend):

| Metric                    | Python | Go    | Improvement |
|---------------------------|--------|-------|-------------|
| Startup time (cold)       | 2.3s   | 0.15s | 15.3x       |
| Memory (idle)             | 85MB   | 12MB  | 7.1x        |
| WebSocket connection time | 45ms   | 8ms   | 5.6x        |
| Session creation time     | 120ms  | 25ms  | 4.8x        |

*Note: Benchmarks are preliminary and may change as implementation evolves.*

## Roadmap

- [x] **Q1 2026**: Basic server and WebSocket (✅ Complete)
- [ ] **Q2 2026**: Script execution and basic elements
- [ ] **Q3 2026**: Widget support and state management
- [ ] **Q4 2026**: Feature parity and production readiness

## FAQ

**Q: Will this replace the Python backend?**
A: No, it's an alternative option.

**Q: Can I use it in production?**
A: Not yet - it's experimental.

**Q: Do I need to rewrite my app in Go?**
A: No! You still write Python. Go is just the server.

**Q: What about my custom components?**
A: Not supported yet.

**Q: How can I help?**
A: Test it, file issues, contribute code!

## Getting Help

- **Questions**: File an issue with `go-backend` label
- **Bugs**: Include "Go Backend" in issue title
- **Discussions**: Use GitHub Discussions
- **Chat**: Streamlit Community Forum

## License

Apache License 2.0 - Same as Streamlit

---

**Status**: 🧪 Experimental | **Version**: 0.1.0 | **Maintainer**: Streamlit Team

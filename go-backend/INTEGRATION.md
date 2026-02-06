# Go Backend Integration Guide

This document explains how the Go backend integrates with the existing Streamlit architecture and how to use both backends together.

## Architecture Integration

### Shared Components

The Go backend **reuses** the following from the main Streamlit project:

1. **Protocol Buffers** (`../proto/streamlit/proto/`)
   - All `.proto` files are shared
   - Go backend generates Go code from the same definitions
   - Ensures protocol compatibility between backends

2. **Frontend Assets** (`../frontend/build/`)
   - Same React/TypeScript frontend
   - No frontend changes needed
   - Frontend is agnostic to backend implementation

3. **Configuration Schema**
   - Respects the same config structure
   - Compatible with `.streamlit/config.toml`

### Independent Components

The Go backend implements **independently**:

1. **HTTP/WebSocket Server**
   - Go's net/http package (instead of Tornado)
   - gorilla/websocket (instead of Tornado WebSocket)

2. **Runtime and Session Management**
   - Native Go concurrency (goroutines)
   - Separate session handling logic

3. **Script Execution** (future)
   - Will integrate with Python scripts
   - Uses subprocess for Python execution

## Using Both Backends

You can have both backends installed and switch between them:

### Python Backend (Default)

```bash
# From repo root
uv run streamlit run your_app.py
```

Runs on: `http://localhost:8501` (default)

### Go Backend (Experimental)

```bash
# From go-backend/
./streamlit-go your_app.py

# Or with custom port to run alongside Python
./streamlit-go --port 8502 your_app.py
```

Runs on: `http://localhost:8502` (if port specified)

### Running Both Simultaneously

```bash
# Terminal 1: Python backend
cd /path/to/streamlit
uv run streamlit run app.py --server.port 8501

# Terminal 2: Go backend
cd /path/to/streamlit/go-backend
./streamlit-go app.py --port 8502
```

Now you can compare:
- Python: `http://localhost:8501`
- Go: `http://localhost:8502`

## Development Workflow

### When Working on Frontend

```bash
# 1. Make frontend changes in ../frontend/

# 2. Rebuild frontend
cd ../frontend
npm run build

# 3. Restart either backend
# Python:
uv run streamlit run app.py

# Go:
cd ../go-backend
./streamlit-go app.py
```

The changes will be reflected in both backends since they serve the same build artifacts.

### When Working on Protocol

```bash
# 1. Modify .proto files in ../proto/

# 2. Regenerate for Python
cd /path/to/streamlit
make protobuf

# 3. Regenerate for Go
cd go-backend
make proto

# 4. Update backend code to handle new messages
# - Python: ../lib/streamlit/
# - Go: internal/
```

### When Working on Go Backend

```bash
# 1. Make changes to Go code

# 2. Rebuild
cd go-backend
make build

# 3. Test
./streamlit-go examples/hello.py

# 4. Run tests
make test
```

## Message Flow Comparison

### Python Backend

```
Browser ←WebSocket→ Tornado Server ←→ Python Runtime
   ↓                      ↓                    ↓
ForwardMsg          BrowserWSHandler      AppSession
BackMsg             (Tornado)             (Python)
```

### Go Backend

```
Browser ←WebSocket→ Go HTTP Server ←→ Go Runtime
   ↓                      ↓                 ↓
ForwardMsg         WebSocketSession     Runtime
BackMsg            (gorilla/ws)         (Go)
```

Both use the **same Protocol Buffer messages** (ForwardMsg/BackMsg).

## File Organization

```
streamlit/                          # Main repo
├── proto/                          # Shared protobuf definitions
│   └── streamlit/proto/           # ← Both backends use these
├── frontend/                       # Shared frontend
│   └── build/                     # ← Both backends serve this
├── lib/                           # Python backend
│   └── streamlit/
│       ├── web/server/            # Python WebSocket server
│       └── runtime/               # Python runtime
└── go-backend/                    # Go backend
    ├── proto/                     # Generated Go protobuf code
    ├── internal/
    │   ├── server/                # Go WebSocket server
    │   └── runtime/               # Go runtime
    └── cmd/streamlit/             # Go CLI
```

## Configuration Compatibility

Both backends can read the same configuration:

**.streamlit/config.toml**
```toml
[server]
port = 8501
headless = false
address = "localhost"

[browser]
serverAddress = "localhost"
```

### Python Backend
Reads via: `streamlit.config`

### Go Backend
Reads via: Command-line flags or environment variables (future: will support config.toml)

## Testing Integration

### Test Strategy

1. **Protocol Tests**: Ensure both backends generate compatible messages
2. **Frontend Tests**: Same tests work with both backends
3. **E2E Tests**: Can run against either backend

### Running E2E Tests Against Go Backend

```bash
# Start Go backend
cd go-backend
./streamlit-go --port 8502 examples/hello.py &

# Run E2E tests pointing to Go backend
cd ../e2e_playwright
STREAMLIT_SERVER_URL=http://localhost:8502 uv run pytest basic_app_test.py
```

## Performance Comparison

| Aspect               | Python Backend | Go Backend      |
|---------------------|----------------|-----------------|
| Startup Time        | ~2-3s          | ~100-200ms      |
| Memory (idle)       | ~50-100MB      | ~10-20MB        |
| Concurrent Sessions | Good           | Excellent       |
| CPU Usage           | Moderate       | Low             |
| Protobuf Ser/Deser  | Good           | Faster          |

**Note**: Go backend performance will improve as implementation matures.

## Migration Path

For apps wanting to use the Go backend:

### 1. No Changes Needed (Basic Apps)

Most simple apps work without modification:
```python
import streamlit as st
st.write("Hello World")
```

### 2. Check Compatibility (Advanced Features)

Some features may not be implemented yet:
- Custom components
- Advanced caching
- Specific integrations

### 3. Test Thoroughly

```bash
# Test with Go backend
./streamlit-go your_app.py

# Compare with Python backend
uv run streamlit run your_app.py
```

## Debugging Both Backends

### Python Backend Debugging

```bash
# Verbose logging
uv run streamlit run app.py --logger.level debug

# Check logs
tail -f ~/.streamlit/logs/
```

### Go Backend Debugging

```bash
# Build with debug symbols
cd go-backend
go build -gcflags="all=-N -l" -o streamlit-go-debug cmd/streamlit/main.go

# Run with verbose output
./streamlit-go-debug app.py

# Use delve debugger
dlv exec ./streamlit-go-debug -- app.py
```

## Contributing to Go Backend

When contributing:

1. **Keep protocol compatibility**: Don't break protobuf messages
2. **Test with existing frontend**: Ensure frontend works unchanged
3. **Document differences**: Note any behavioral differences
4. **Follow Go conventions**: Use Go idioms and patterns

## FAQ

### Q: Will the Go backend replace the Python backend?

**A**: No, this is an experimental alternative. The Python backend remains the primary implementation.

### Q: Can I use Go backend in production?

**A**: Not recommended yet. This is experimental and not feature-complete.

### Q: Do I need to write Go code to use it?

**A**: No! The Go backend runs Python scripts. You write Python code as usual.

### Q: What about custom components?

**A**: Custom components are not fully supported yet in the Go backend.

### Q: Can I deploy apps using Go backend?

**A**: Deployment support is planned but not ready yet.

## Roadmap

### Phase 1: Core Foundation ✅
- [x] Project structure
- [x] WebSocket server
- [x] Protobuf integration
- [x] Static file serving
- [x] Basic runtime

### Phase 2: Script Execution 🚧
- [ ] Python script integration
- [ ] Element rendering
- [ ] Widget state management
- [ ] Session persistence

### Phase 3: Feature Parity
- [ ] All basic elements
- [ ] All widgets
- [ ] Caching
- [ ] File uploads
- [ ] Media files

### Phase 4: Advanced Features
- [ ] Custom components
- [ ] Authentication
- [ ] Deployment
- [ ] Performance optimization

### Phase 5: Production Ready
- [ ] Comprehensive testing
- [ ] Documentation
- [ ] Deployment guides
- [ ] Migration tools

## Resources

- **Main Repo**: [github.com/streamlit/streamlit](https://github.com/streamlit/streamlit)
- **Go Backend**: `streamlit/go-backend/` (this directory)
- **Protocol Buffers**: `streamlit/proto/`
- **Frontend**: `streamlit/frontend/`

For questions or issues, please file an issue on the main Streamlit repo with the `go-backend` label.

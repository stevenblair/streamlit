# Streamlit Go Backend - Quick Start Guide

## Overview

This directory contains an experimental Go implementation of the Streamlit backend. It reuses the existing frontend assets and protobuf protocol definitions to provide an alternative server implementation.

## Architecture Overview

```
┌─────────────┐          WebSocket           ┌──────────────┐
│   Browser   │ ◄───────────────────────────► │  Go Backend  │
│  (Frontend) │    Protocol Buffers (proto)   │   (Server)   │
└─────────────┘                                └──────────────┘
      │                                               │
      │ Serves static assets                         │
      └───────────────────────────────────────────────┘
```

### Key Components

1. **cmd/streamlit/main.go** - Entry point and CLI
2. **internal/server/** - HTTP/WebSocket server implementation
3. **internal/runtime/** - Script execution and session management
4. **proto/** - Generated protobuf code (auto-generated from `../proto/`)

## Quick Start

### Prerequisites

```bash
# Install Go 1.22+
# Download from: https://go.dev/dl/

# Install protoc (Protocol Buffer compiler)
# Windows: Download from https://github.com/protocolbuffers/protobuf/releases
# macOS: brew install protobuf
# Linux: apt-get install protobuf-compiler
```

### Setup

```bash
# 1. Navigate to go-backend directory
cd go-backend

# 2. Run setup script
./scripts/setup.sh      # Unix/macOS
scripts\setup.bat       # Windows

# 3. Generate protobuf code
./scripts/proto.sh      # Unix/macOS
scripts\proto.bat       # Windows

# 4. Build the frontend (from repo root)
cd ../frontend
npm install && npm run build
cd ../go-backend

# 5. Build everything
./scripts/build.sh      # Unix/macOS
scripts\build.bat       # Windows
```

### Running

```bash
# Start the server
./bin/streamlit-go

# In another terminal, run an example app
./bin/hello

# Or on Windows
.\bin\streamlit-go.exe
.\bin\hello.exe

# Then open your browser to:
# http://localhost:8501
```

### Using the Makefile

```bash
make setup      # Install dependencies
make proto      # Generate protobuf code
make build      # Build the binary
make run        # Build and run example app
make test       # Run tests
make clean      # Clean build artifacts
```

## Project Status

### ✅ Implemented

- Basic HTTP server with routing
- WebSocket connection handling
- Static file serving for frontend assets
- Health check endpoint
- Host configuration endpoint
- Project structure and build system
- Protobuf generation pipeline

### 🚧 In Progress / TODO

- Full protobuf message handling (ForwardMsg/BackMsg)
- Script execution integration
- Widget state management
- Session persistence
- Cache management
- Media file handling
- File upload support
- Custom component support

### ❌ Not Yet Implemented

- Full API parity with Python backend
- Authentication/authorization
- Deployment configurations
- Performance optimizations
- Comprehensive test coverage

## Development

### Directory Structure

```
go-backend/
├── cmd/
│   └── streamlit/              # Main application
│       └── main.go            # Entry point
├── internal/
│   ├── runtime/                # Runtime and execution
│   │   └── runtime.go         # Script execution logic
│   └── server/                 # HTTP/WebSocket server
│       ├── server.go          # Main server
│       └── websocket.go       # WebSocket handler
├── examples/                   # Example Streamlit apps
│   ├── hello.py               # Basic example
│   └── widgets.py             # Widget examples
├── proto/                      # Generated protobuf (auto-generated)
├── scripts/                    # Build scripts
│   ├── generate_protos.sh     # Unix protobuf generation
│   └── generate_protos.bat    # Windows protobuf generation
├── go.mod                      # Go module definition
├── Makefile                    # Build automation
└── README.md                   # This file
```

### Adding New Features

1. **Adding Protobuf Support**: Regenerate proto files after changes
   ```bash
   make proto
   ```

2. **Adding HTTP Endpoints**: Edit `internal/server/server.go`

3. **Handling New Message Types**: Edit `internal/server/websocket.go`

4. **Script Execution Logic**: Edit `internal/runtime/runtime.go`

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test ./internal/server -v
```

### Debugging

```bash
# Build with debug symbols
go build -gcflags="all=-N -l" -o streamlit-go-debug cmd/streamlit/main.go

# Run with verbose logging
STREAMLIT_LOG_LEVEL=debug ./streamlit-go examples/hello.py
```

## Communication Protocol

The Go backend uses the same Protocol Buffer definitions as the Python backend:

- **ForwardMsg**: Server → Browser messages (render elements, updates, etc.)
- **BackMsg**: Browser → Server messages (user interactions, reruns, etc.)

See `../proto/streamlit/proto/` for complete protocol definitions.

## Differences from Python Backend

### Advantages
- **Performance**: Compiled binary, faster startup
- **Concurrency**: Native goroutines for handling multiple sessions
- **Deployment**: Single binary, no Python runtime needed
- **Resource Usage**: Lower memory footprint

### Limitations
- **Script Execution**: Currently spawns Python subprocess (not native)
- **Feature Parity**: Not all features implemented yet
- **Ecosystem**: Less integration with Python ecosystem
- **Maturity**: Experimental, not production-ready

## Contributing

This is an experimental project. To contribute:

1. Follow the main Streamlit contribution guidelines
2. Ensure code passes `go vet` and `go fmt`
3. Add tests for new functionality
4. Update documentation

## Troubleshooting

### Frontend not found
```
Error: frontend build directory not found
```
**Solution**: Build the frontend first:
```bash
cd ../frontend && npm install && npm run build
```

### Protobuf errors
```
Error: proto files not found or compilation failed
```
**Solution**: Ensure protoc is installed and regenerate:
```bash
make proto
```

### Port already in use
```
Error: bind: address already in use
```
**Solution**: Use a different port:
```bash
./streamlit-go --port 8502 examples/hello.py
```

## Resources

- [Streamlit Documentation](https://docs.streamlit.io)
- [Go Documentation](https://go.dev/doc)
- [Protocol Buffers](https://protobuf.dev)
- [Gorilla WebSocket](https://github.com/gorilla/websocket)

## License

Apache License 2.0 - Same as the main Streamlit project.

# Windows Setup Guide for Streamlit Go Backend

This guide provides Windows-specific instructions for setting up and running the Streamlit Go backend.

## Prerequisites

### 1. Install Go

Download and install Go 1.22 or later from [go.dev/dl](https://go.dev/dl/)

```powershell
# Verify installation
go version
# Should output: go version go1.22.x windows/amd64
```

### 2. Install Protocol Buffers Compiler

**Option A: Using Chocolatey (Recommended)**
```powershell
choco install protoc
```

**Option B: Manual Installation**
1. Download from [github.com/protocolbuffers/protobuf/releases](https://github.com/protocolbuffers/protobuf/releases)
2. Look for `protoc-XX.X-win64.zip`
3. Extract to `C:\protoc`
4. Add `C:\protoc\bin` to your PATH

**Option C: Using Scoop**
```powershell
scoop install protobuf
```

Verify installation:
```powershell
protoc --version
# Should output: libprotoc X.XX.X
```

### 3. Install Python (for running scripts)

Download from [python.org](https://python.org) or use:

```powershell
# Using Chocolatey
choco install python

# Using Scoop
scoop install python
```

## Quick Setup

### 1. Navigate to go-backend

```powershell
cd C:\Users\YourName\streamlit\go-backend
```

### 2. Download Go Dependencies

```powershell
go mod download
```

This will download:
- `github.com/gorilla/websocket`
- `google.golang.org/protobuf`

### 3. Generate Protobuf Code

```powershell
# Install protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# Generate proto files
.\scripts\generate_protos.bat
```

Expected output:
```
Generating Go protobuf code from Streamlit proto files...
Generating protobuf files...
Generating: Alert.proto
Generating: AppPage.proto
...
Protobuf generation complete!
```

### 4. Build the Go Backend

```powershell
# Build using go build
go build -o streamlit-go.exe .\cmd\streamlit\main.go

# Or use the binary output directory
go build -o bin\streamlit-go.exe .\cmd\streamlit\main.go
```

Expected output: `streamlit-go.exe` binary created

### 5. Build the Frontend

The Go backend serves the compiled frontend assets, so build them first:

```powershell
# Navigate to frontend directory
cd ..\frontend

# Install dependencies (first time only)
npm install

# Build the frontend
npm run build

# Return to go-backend
cd ..\go-backend
```

### 6. Run the Example

```powershell
.\streamlit-go.exe examples\hello.py
```

Expected output:
```
Starting Streamlit Go backend...
Script: C:\Users\YourName\streamlit\go-backend\examples\hello.py
Server: http://localhost:8501

You can now view your Streamlit app in your browser.

  Local URL: http://localhost:8501
```

### 7. Open in Browser

Open your browser to: http://localhost:8501

## Troubleshooting

### Error: "protoc: command not found"

**Solution**: Add protoc to your PATH

```powershell
# Check current PATH
$env:Path

# Add protoc to PATH (temporary)
$env:Path += ";C:\protoc\bin"

# Add permanently via System Properties > Environment Variables
# Or use PowerShell (requires admin):
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\protoc\bin", "Machine")
```

### Error: "frontend build directory not found"

**Solution**: Build the frontend first

```powershell
cd ..\frontend
npm install
npm run build
cd ..\go-backend
```

### Error: "cannot find package"

**Solution**: Download dependencies

```powershell
go mod download
go mod tidy
```

### Error: "Port 8501 already in use"

**Solution**: Use a different port

```powershell
.\streamlit-go.exe --port 8502 examples\hello.py
```

### Error: "protoc-gen-go: program not found"

**Solution**: Install and add to PATH

```powershell
# Install
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# The binary is in %USERPROFILE%\go\bin
# Add to PATH:
$env:Path += ";$env:USERPROFILE\go\bin"
```

## Using PowerShell Scripts

Instead of using Makefile (which requires Unix tools), use PowerShell directly:

### Setup Script

Create `setup.ps1`:
```powershell
# Download dependencies
go mod download

# Install protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

Write-Host "Setup complete!" -ForegroundColor Green
```

Run:
```powershell
.\setup.ps1
```

### Build Script

Create `build.ps1`:
```powershell
# Generate protobuf files
.\scripts\generate_protos.bat

# Build the binary
go build -o bin\streamlit-go.exe .\cmd\streamlit\main.go

Write-Host "Build complete!" -ForegroundColor Green
```

Run:
```powershell
.\build.ps1
```

## Development Workflow on Windows

### 1. Make Code Changes

Edit Go files in `internal/` or `cmd/` using your favorite editor:
- Visual Studio Code (recommended)
- GoLand
- Notepad++

### 2. Rebuild

```powershell
go build -o streamlit-go.exe .\cmd\streamlit\main.go
```

### 3. Test

```powershell
.\streamlit-go.exe examples\test.py
```

### 4. Debug

Using Visual Studio Code:

Create `.vscode\launch.json`:
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Streamlit Go",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/cmd/streamlit",
            "args": ["examples/hello.py"]
        }
    ]
}
```

Press F5 to debug.

## Using with Make (Optional)

If you have Make installed (via Chocolatey, Scoop, or WSL):

```powershell
# Install Make
choco install make

# Use Makefile targets
make setup
make proto
make build
make test
```

## Performance Tips

### 1. Use Go Build Cache

Go caches builds automatically. To see cache stats:

```powershell
go env GOCACHE
```

### 2. Build with Optimizations

```powershell
# Release build (smaller, faster)
go build -ldflags="-s -w" -o streamlit-go.exe .\cmd\streamlit\main.go
```

### 3. Cross-Compile

Build for other platforms from Windows:

```powershell
# Build for Linux
$env:GOOS="linux"; $env:GOARCH="amd64"
go build -o streamlit-go-linux .\cmd\streamlit\main.go

# Build for macOS
$env:GOOS="darwin"; $env:GOARCH="amd64"
go build -o streamlit-go-macos .\cmd\streamlit\main.go

# Reset to Windows
$env:GOOS="windows"; $env:GOARCH="amd64"
```

## Running Tests

```powershell
# Run all tests
go test .\...

# Run with verbose output
go test -v .\...

# Run specific test
go test -v .\internal\server -run TestHealthHandler
```

## IDE Setup

### Visual Studio Code

1. Install Go extension: `ms-vscode.go`
2. Open go-backend folder in VS Code
3. Install recommended tools when prompted

**Recommended VS Code Extensions:**
- Go (golang.go)
- Protocol Buffers (zxh404.vscode-proto3)
- REST Client (humao.rest-client)

### GoLand

1. Open `go-backend` as a project
2. GoLand will auto-detect go.mod
3. Enable Go Modules in Settings

## Common Commands

```powershell
# Download dependencies
go mod download

# Update dependencies
go get -u ./...
go mod tidy

# Format code
go fmt ./...

# Run linter
go vet ./...

# Check for issues
go mod verify

# Clean build cache
go clean -cache

# Show dependencies
go list -m all
```

## Next Steps

1. Read [QUICKSTART.md](QUICKSTART.md) for detailed usage
2. Check [DEVELOPMENT.md](DEVELOPMENT.md) for development guidelines
3. Explore [examples/](examples/) for sample apps
4. Join the Streamlit community for support

## Resources

- [Go Documentation](https://go.dev/doc/)
- [Go by Example](https://gobyexample.com/)
- [Protocol Buffers Go Tutorial](https://protobuf.dev/getting-started/gotutorial/)
- [Streamlit Documentation](https://docs.streamlit.io)

---

**Need Help?** File an issue with the `go-backend` label on the main Streamlit repository.

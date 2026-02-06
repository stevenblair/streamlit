# Development Notes - Streamlit Go Backend

## Implementation Notes

### Design Decisions

1. **Gorilla WebSocket**: Chose gorilla/websocket over alternatives because:
   - Well-maintained and widely used
   - Good documentation
   - Compatible with existing protocol

2. **Internal Package Structure**: Following Go conventions:
   - `internal/` for private packages
   - `cmd/` for executables
   - `pkg/` avoided (not needed yet)

3. **Protobuf Generation**: Using `protoc-gen-go`:
   - Generate from existing `.proto` files
   - Keep in sync with Python backend
   - Auto-generate, don't commit to git

### Current Limitations

1. **Script Execution**: Currently not implemented
   - Will use subprocess to run Python
   - Need to capture stdout/stderr
   - Parse Python API calls

2. **Message Handling**: Basic structure only
   - BackMsg parsing incomplete
   - ForwardMsg generation incomplete
   - Need full message type handlers

3. **State Management**: Simplified
   - No session persistence yet
   - Widget state not tracked
   - Cache not implemented

## Next Steps

### Immediate (Phase 1)

- [ ] Implement BackMsg unmarshaling
- [ ] Implement ForwardMsg marshaling
- [ ] Add basic element rendering
- [ ] Handle rerun_script message
- [ ] Add session ID generation

### Short-term (Phase 2)

- [ ] Python script execution via subprocess
- [ ] Capture Python output
- [ ] Parse st.* API calls
- [ ] Convert to ForwardMsgs
- [ ] Send to frontend

### Medium-term (Phase 3)

- [ ] Widget state management
- [ ] Session persistence
- [ ] File upload handling
- [ ] Media file support
- [ ] Cache implementation

### Long-term (Phase 4)

- [ ] Custom component support
- [ ] Authentication integration
- [ ] Production deployment guides
- [ ] Performance optimizations
- [ ] Comprehensive testing

## Technical Challenges

### 1. Script Execution

**Challenge**: How to run Python scripts from Go?

**Options**:
- A) Subprocess: Simple but limited
- B) Embedded Python: Complex but powerful
- C) RPC to Python: Clean separation

**Current Approach**: Planning to use subprocess initially

### 2. API Call Capture

**Challenge**: How to capture st.* calls from Python?

**Options**:
- A) Intercept stdout/serialized calls
- B) Use Python C API
- C) Custom Python backend module

**Current Approach**: Need to decide

### 3. State Synchronization

**Challenge**: Keep widget state consistent

**Options**:
- A) In-memory maps (simple)
- B) Persistent storage (reliable)
- C) Hybrid approach

**Current Approach**: In-memory for now

## Testing Strategy

### Unit Tests

```go
// Example: Test WebSocket handler
func TestWebSocketSession(t *testing.T) {
    // Create mock connection
    // Send test messages
    // Verify responses
}
```

### Integration Tests

```bash
# Start Go backend
./streamlit-go test_app.py &

# Run Python E2E tests
cd ../e2e_playwright
pytest basic_app_test.py
```

### Performance Tests

```go
func BenchmarkSessionCreation(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Create session
        // Measure time
    }
}
```

## Code Organization

### Package Dependencies

```
cmd/streamlit
    └─> internal/server
            └─> internal/runtime
                    └─> proto (generated)
```

### Message Flow

```
Browser
    │
    ▼ BackMsg (binary protobuf)
WebSocketSession.handleBackMsg()
    │
    ▼ Unmarshal
BackMsg struct
    │
    ▼ Route based on type
Runtime.RunScript() / StopScript() / etc.
    │
    ▼ Execute action
Generate ForwardMsg
    │
    ▼ Marshal
ForwardMsg (binary protobuf)
    │
    ▼ Send over WebSocket
Browser
```

## Performance Considerations

### Memory Usage

- Use sync.Pool for protobuf messages
- Reuse byte buffers
- Limit message queue sizes

### Concurrency

- One goroutine per WebSocket connection
- Shared runtime across sessions
- Use channels for communication

### Startup Time

- Lazy load resources
- Parallel initialization
- Cache compiled assets

## Debugging Tips

### Enable Verbose Logging

```go
// In websocket.go
log.Printf("DEBUG: Received %d bytes: %v", len(data), data)
```

### Use Delve Debugger

```bash
dlv debug ./cmd/streamlit -- examples/hello.py
(dlv) break internal/server.(*WebSocketSession).handleBackMsg
(dlv) continue
```

### Inspect WebSocket Traffic

```javascript
// Browser console
window.addEventListener('message', (e) => {
    console.log('WebSocket:', e.data);
});
```

## Resources Used

### Go Libraries

- `github.com/gorilla/websocket` - WebSocket handling
- `google.golang.org/protobuf` - Protocol Buffers
- Standard library: `net/http`, `context`, `sync`

### Documentation

- [Tornado WebSocket Docs](https://www.tornadoweb.org/en/stable/websocket.html)
- [gorilla/websocket Examples](https://github.com/gorilla/websocket/tree/master/examples)
- [Protocol Buffers Go Tutorial](https://protobuf.dev/getting-started/gotutorial/)

## Known Issues

1. **Protobuf import errors**: Need to run `make proto` first
2. **Frontend not found**: Need to build frontend with `make frontend-fast`
3. **Port conflicts**: Use `--port` flag to specify different port

## Future Optimizations

1. **Connection pooling**: Reuse HTTP connections
2. **Message batching**: Combine multiple ForwardMsgs
3. **Compression**: Enable WebSocket compression
4. **Caching**: Cache static assets and responses
5. **Load balancing**: Support multiple Go backend instances

## Comparison with Python Backend

### Similarities

- Same protocol (protobuf)
- Same frontend
- Same configuration schema
- Same user-facing API

### Differences

| Aspect          | Python               | Go                    |
|-----------------|----------------------|-----------------------|
| Server          | Tornado              | net/http              |
| WebSocket       | Tornado WebSocket    | gorilla/websocket     |
| Concurrency     | asyncio              | goroutines            |
| Type System     | Dynamic              | Static                |
| Performance     | Good                 | Better (expected)     |
| Dependencies    | Many                 | Minimal               |
| Binary Size     | N/A (interpreted)    | ~10-20MB              |

## Git Workflow

```bash
# Create feature branch
git checkout -b go-backend/feature-name

# Make changes in go-backend/

# Test
cd go-backend
make test

# Commit
git add go-backend/
git commit -m "go-backend: Add feature X"

# Push and create PR
git push origin go-backend/feature-name
```

## Release Process (Future)

1. Tag release: `go-backend-v0.1.0`
2. Build binaries for all platforms
3. Create GitHub release
4. Update documentation
5. Announce in community

---

**Last Updated**: 2026-02-06
**Status**: Active Development

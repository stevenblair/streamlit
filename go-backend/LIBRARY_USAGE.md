# Using as a Go Library

## Current Status: ✅ Ready for `go get`

Protobuf files (`.pb.go`) are now **committed** to the repository, making this library immediately usable:

```bash
go get github.com/stevenblair/streamlit/go-backend
```

```go
import st "github.com/stevenblair/streamlit/go-backend"

func main() {
    st.Run(func() {
        st.Title("Hello!")
    })
}
```

## Development

**When to regenerate `.pb.go` files:**
- Proto files change in upstream Streamlit
- You modify proto definitions

**How to regenerate:**
```bash
cd go-backend
scripts\generate_protos_clean.bat
git add proto/*.pb.go
git commit -m "chore: regenerate protobuf files"
```

## Trade-offs

**Committed .pb.go files:**
- ✅ Users can `go get` immediately (no protoc required)
- ✅ Zero setup for library usage
- ⚠️ Adds ~500KB to repo
- ⚠️ Generated files in git

**Still upstream compatible:**
- Original proto files never modified
- `.pb.go` files only in `go-backend/proto/`
- No conflicts with `proto/streamlit/proto/*.proto`

# Upstream Compatibility

## Problem Solved

The Go backend previously modified 87+ upstream files (proto files and `lib/streamlit/static/index.html`), causing merge conflicts.

## Solution

**Never modify upstream files:**

1. **Frontend**: Separate copy in `go-backend/static/`
2. **Protobufs**: Generated files in Go cache (gitignored)
3. **Original files**: Remain pristine

## File Locations

| Type | Upstream | Go Backend | Status |
|------|----------|------------|--------|
| Proto definitions | `proto/streamlit/proto/*.proto` | Same (read-only) | Never modified |
| Generated `.pb.go` | N/A | Go module cache | Generated, cached |
| React source | `frontend/app/` | Same (read-only) | Never modified |
| React build | `lib/streamlit/static/` | `go-backend/static/` | Separate copy |

## Syncing with Upstream

```bash
# Pull upstream changes (no conflicts)
git checkout develop
git pull upstream develop

# Optional: Rebuild frontend if changed
cd frontend/app
npm run build
xcopy /E /I /Y build ..\..\go-backend\static
```

## Verification

```bash
# Check git status - should NOT show proto or index.html
git status

# Verify proto files are pristine
git diff proto/streamlit/proto/Alert.proto
# Output: (nothing)

# Verify index.html is pristine
git diff lib/streamlit/static/index.html
# Output: (nothing)
```

## Benefits

✅ No merge conflicts when pulling upstream
✅ Automatic pickup of new proto files
✅ Frontend updates work seamlessly
✅ Clean git history

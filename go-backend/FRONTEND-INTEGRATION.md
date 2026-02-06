# Frontend Integration

## Status: ✅ Complete

The Go backend serves the **production React frontend** with full protobuf messaging.

## Quick Start

```bash
cd go-backend/examples/waveform
go build -o waveform.exe .
.\waveform.exe
```

Open http://localhost:8501

## Architecture

```
┌─────────────────────┐
│   App Binary        │
│  (waveform.exe)     │
│                     │
│  ┌───────────────┐  │
│  │  st.Run()     │  │
│  │  - HTTP       │◄─┼─── http://localhost:8501/
│  │  - WebSocket  │◄─┼─── ws://localhost:8501/_stcore/stream
│  └───────────────┘  │
└─────────────────────┘
         │
         │ serves
         ▼
┌──────────────────────┐
│  go-backend/static/  │
│  - index.html        │
│  - static/js/...     │
│  - static/css/...    │
└──────────────────────┘
```

## Frontend Path Resolution

The server checks these locations (in order):
1. `../../static/` (from examples/waveform/)
2. `../../go-backend/static/`
3. `../static/`
4. `../../lib/streamlit/static/` (upstream fallback)

First match is used.

## Rebuilding Frontend (Optional)

Frontend is already built and committed. Rebuild only when syncing upstream changes:

```bash
cd frontend/app
npm install
npm run build
xcopy /E /I /Y build ..\..\go-backend\static
```

## Working Features

✅ All elements (text, headings, code, alerts, markdown)
✅ All widgets (button, slider, text_input, checkbox, selectbox)
✅ Layouts (columns, sidebar, metrics)
✅ Charts (line_chart with Arrow encoding)
✅ Animations (balloons, snow)
✅ Widget state persistence

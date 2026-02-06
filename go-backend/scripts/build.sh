#!/bin/bash
# Build script for Streamlit Go backend

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Building Streamlit Go backend...${NC}"

# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
GO_BACKEND_DIR="$(dirname "$SCRIPT_DIR")"

# Create bin directory
mkdir -p "$GO_BACKEND_DIR/bin"

# Build the main binary
echo "Building streamlit-go..."
cd "$GO_BACKEND_DIR"
go build -o bin/streamlit-go ./cmd/streamlit

# Build example applications
echo "Building example apps..."
go build -o bin/hello ./examples/hello
go build -o bin/widgets ./examples/widgets
go build -o bin/waveform ./examples/waveform

echo -e "${GREEN}✓ Build complete!${NC}"
echo ""
echo "Binaries created:"
echo "  bin/streamlit-go    # Main server"
echo "  bin/hello           # Hello example"
echo "  bin/widgets         # Widgets example"
echo "  bin/waveform        # Waveform example"
echo ""
echo "To run:"
echo "  ./bin/streamlit-go   # Start server with default app"
echo "  ./bin/hello          # Run hello example"
echo "  ./bin/widgets        # Run widgets example"
echo "  ./bin/waveform       # Run waveform example"

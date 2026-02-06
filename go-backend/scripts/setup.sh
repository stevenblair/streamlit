#!/bin/bash
# Setup script for Streamlit Go backend

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Setting up Streamlit Go backend...${NC}"

# Check Go installation
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.22 or later."
    echo "Download from: https://go.dev/dl/"
    exit 1
fi

echo "✓ Go version: $(go version)"

# Check protoc installation
if ! command -v protoc &> /dev/null; then
    echo "Warning: protoc compiler not found."
    echo "  macOS: brew install protobuf"
    echo "  Linux: apt-get install protobuf-compiler"
    echo "  Or download from: https://github.com/protocolbuffers/protobuf/releases"
    echo ""
fi

# Download Go dependencies
echo -e "${BLUE}Downloading Go modules...${NC}"
go mod download

# Install protoc-gen-go
echo -e "${BLUE}Installing protoc-gen-go...${NC}"
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

echo -e "${GREEN}✓ Setup complete!${NC}"
echo ""
echo "Next steps:"
echo "  1. Run: ./scripts/proto.sh     # Generate protobuf code"
echo "  2. Run: ./scripts/build.sh     # Build the binary"
echo "  3. Run: ./bin/streamlit-go     # Run an example"

#!/bin/bash
# Script to generate Go protobuf code from the Streamlit proto files

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Generating Go protobuf code from Streamlit proto files...${NC}"

# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
GO_BACKEND_DIR="$(dirname "$SCRIPT_DIR")"
REPO_ROOT="$(dirname "$GO_BACKEND_DIR")"
PROTO_SRC_DIR="$REPO_ROOT/proto"
PROTO_OUT_DIR="$GO_BACKEND_DIR/proto"

# Create output directory
mkdir -p "$PROTO_OUT_DIR"

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc compiler not found. Please install protobuf compiler."
    echo "  macOS: brew install protobuf"
    echo "  Linux: apt-get install protobuf-compiler"
    echo "  Or download from: https://github.com/protocolbuffers/protobuf/releases"
    exit 1
fi

# Check if protoc-gen-go is installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "Error: protoc-gen-go not found. Installing..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

# Find all proto files
PROTO_FILES=$(find "$PROTO_SRC_DIR" -name "*.proto")

if [ -z "$PROTO_FILES" ]; then
    echo "Error: No proto files found in $PROTO_SRC_DIR"
    exit 1
fi

echo -e "${BLUE}Found $(echo "$PROTO_FILES" | wc -l) proto files${NC}"

# Generate Go code for each proto file
for proto_file in $PROTO_FILES; do
    echo "Generating: $(basename $proto_file)"
    protoc \
        --proto_path="$PROTO_SRC_DIR" \
        --go_out="$PROTO_OUT_DIR" \
        --go_opt=paths=source_relative \
        "$proto_file"
done

echo -e "${GREEN}✓ Protobuf generation complete!${NC}"
echo -e "Generated files in: ${PROTO_OUT_DIR}"

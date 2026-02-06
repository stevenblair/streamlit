#!/bin/bash
# Test script for Streamlit Go backend

set -e

BLUE='\033[0;34m'
GREEN='\033[0;32m'
NC='\033[0m'

echo -e "${BLUE}Running tests for Streamlit Go backend...${NC}"

# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
GO_BACKEND_DIR="$(dirname "$SCRIPT_DIR")"

cd "$GO_BACKEND_DIR"

# Run all tests
echo "Running unit tests..."
go test -v ./...

echo -e "${GREEN}✓ All tests passed!${NC}"

@echo off
REM Setup script for Streamlit Go backend (Windows)

setlocal enabledelayedexpansion

echo Setting up Streamlit Go backend...
echo.

REM Check Go installation
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Error: Go is not installed. Please install Go 1.22 or later.
    echo Download from: https://go.dev/dl/
    exit /b 1
)

echo Go version:
go version
echo.

REM Check protoc installation
where protoc >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Warning: protoc compiler not found.
    echo Download from: https://github.com/protocolbuffers/protobuf/releases
    echo.
)

REM Download Go dependencies
echo Downloading Go modules...
go mod download

REM Install protoc-gen-go
echo Installing protoc-gen-go...
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

echo.
echo Setup complete!
echo.
echo Next steps:
echo   1. Run: scripts\proto.bat      # Generate protobuf code
echo   2. Run: scripts\build.bat      # Build the binary
echo   3. Run: bin\streamlit-go.exe   # Run an example

endlocal

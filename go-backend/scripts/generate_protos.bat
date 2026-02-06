@echo off
REM Script to generate Go protobuf code from the Streamlit proto files (Windows)

setlocal enabledelayedexpansion

echo Generating Go protobuf code from Streamlit proto files...

REM Get directories
set "SCRIPT_DIR=%~dp0"
set "GO_BACKEND_DIR=%SCRIPT_DIR%.."
set "REPO_ROOT=%GO_BACKEND_DIR%\.."
set "PROTO_SRC_DIR=%REPO_ROOT%\proto"
set "PROTO_OUT_DIR=%GO_BACKEND_DIR%\proto"

REM Create output directory
if not exist "%PROTO_OUT_DIR%" mkdir "%PROTO_OUT_DIR%"

REM Check if protoc is installed
where protoc >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Error: protoc compiler not found. Please install protobuf compiler.
    echo Download from: https://github.com/protocolbuffers/protobuf/releases
    exit /b 1
)

REM Check if protoc-gen-go is installed
where protoc-gen-go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Error: protoc-gen-go not found. Installing...
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
)

echo Generating protobuf files...

REM Generate Go code for all proto files
for /r "%PROTO_SRC_DIR%" %%f in (*.proto) do (
    echo Generating: %%~nxf
    protoc --proto_path="%REPO_ROOT%" --go_out="%PROTO_OUT_DIR%" --go_opt=paths=source_relative "%%f"
)

echo.
echo Protobuf generation complete!
echo Generated files in: %PROTO_OUT_DIR%

endlocal

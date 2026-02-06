@echo off
REM Build script for Streamlit Go backend (Windows)

setlocal enabledelayedexpansion

echo Building Streamlit Go backend...
echo.

REM Get script directory
set "SCRIPT_DIR=%~dp0"
set "GO_BACKEND_DIR=%SCRIPT_DIR%.."

REM Create bin directory
if not exist "%GO_BACKEND_DIR%\bin" mkdir "%GO_BACKEND_DIR%\bin"

REM Build the main binary
echo Building streamlit-go.exe...
cd "%GO_BACKEND_DIR%"
go build -o bin\streamlit-go.exe .\cmd\streamlit

REM Build example applications
echo Building example apps...
go build -o bin\hello.exe .\examples\hello
go build -o bin\widgets.exe .\examples\widgets
go build -o bin\waveform.exe .\examples\waveform
echo.
echo Build complete!
echo.
echo Binaries created:
echo   bin\streamlit-go.exe    # Main server
echo   bin\hello.exe           # Hello example
echo   bin\widgets.exe         # Widgets example
echo   bin\waveform.exe        # Waveform example
echo.
echo To run:
echo   .\bin\streamlit-go.exe     # Start server with default app
echo   .\bin\hello.exe            # Run hello example
echo   .\bin\widgets.exe          # Run widgets example
echo   .\bin\waveform.exe         # Run waveform example

endlocal

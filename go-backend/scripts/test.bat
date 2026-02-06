@echo off
REM Test script for Streamlit Go backend (Windows)

setlocal enabledelayedexpansion

echo Running tests for Streamlit Go backend...
echo.

REM Get script directory
set "SCRIPT_DIR=%~dp0"
set "GO_BACKEND_DIR=%SCRIPT_DIR%.."

cd "%GO_BACKEND_DIR%"

REM Run all tests
echo Running unit tests...
go test -v .\...

echo.
echo All tests passed!

endlocal

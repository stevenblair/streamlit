@echo off
REM Script to generate Go protobuf code WITHOUT modifying original proto files

setlocal enabledelayedexpansion

echo Generating Go protobuf code (upstream-compatible version)...

REM Get directories
set "SCRIPT_DIR=%~dp0"
set "GO_BACKEND_DIR=%SCRIPT_DIR%.."
set "REPO_ROOT=%GO_BACKEND_DIR%\.."
set "PROTO_SRC_DIR=%REPO_ROOT%\proto\streamlit\proto"
set "PROTO_OUT_DIR=%GO_BACKEND_DIR%\proto"
set "TEMP_PROTO_DIR=%GO_BACKEND_DIR%\temp_proto"

REM Clean up previous temp and output
if exist "%TEMP_PROTO_DIR%" rmdir /s /q "%TEMP_PROTO_DIR%"

REM Create directories
mkdir "%TEMP_PROTO_DIR%\streamlit\proto"
mkdir "%PROTO_OUT_DIR%"

echo.
echo Step 1: Copying proto files to temp directory...
xcopy "%PROTO_SRC_DIR%\*" "%TEMP_PROTO_DIR%\streamlit\proto\" /I /Q

echo.
echo Step 2: Adding go_package option to temp proto files...
powershell -ExecutionPolicy Bypass -File "%SCRIPT_DIR%add_go_package.ps1" -TempProtoDir "%TEMP_PROTO_DIR%"

echo.
echo Step 3: Generating Go code from temp proto files...
cd "%TEMP_PROTO_DIR%"
for %%f in (streamlit\proto\*.proto) do (
    echo Generating: %%~nxf
    protoc --proto_path=. --go_out="%PROTO_OUT_DIR%" --go_opt=module=github.com/streamlit/streamlit/go-backend/proto "%%f"
)
cd "%SCRIPT_DIR%"

echo.
echo Step 4: Cleaning up temp directory...
rmdir /s /q "%TEMP_PROTO_DIR%"

echo.
echo Protobuf generation complete!
echo - Original proto files: UNCHANGED
echo - Generated .pb.go files: %PROTO_OUT_DIR%
echo.
echo This approach is upstream-compatible!

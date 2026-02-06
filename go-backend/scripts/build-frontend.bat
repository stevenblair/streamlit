@echo off
REM Build the Streamlit frontend and copy to static directory

echo Building Streamlit frontend...
echo.

cd ..\frontend

REM Check if node_modules exists
if not exist "node_modules" (
    echo Installing dependencies with corepack...
    corepack enable
    corepack yarn install
)

REM Build the frontend
echo Building frontend packages...
corepack yarn workspaces foreach --recursive --topological --parallel --from @streamlit/app --exclude @streamlit/lib run build

REM Copy build to static directory
echo.
echo Copying build to lib/streamlit/static...

if not exist "..\lib\streamlit\static" mkdir "..\lib\streamlit\static"

REM Use robocopy instead of rsync on Windows
robocopy "app\build" "..\lib\streamlit\static" /E /PURGE /XD reports /NFL /NDL /NJH /NJS

REM Move manifest.json
if exist "..\lib\streamlit\static\.vite\manifest.json" (
    move /Y "..\lib\streamlit\static\.vite\manifest.json" "..\lib\streamlit\static\manifest.json"
)

echo.
echo Frontend build complete!
echo Static files available at: lib\streamlit\static\

cd ..\go-backend

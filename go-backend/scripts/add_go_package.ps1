param(
    [string]$TempProtoDir
)

$files = Get-ChildItem -Path "$TempProtoDir\streamlit\proto" -Filter *.proto

foreach ($file in $files) {
    $content = Get-Content $file.FullName -Raw

    if ($content -notmatch 'option go_package') {
        $lines = $content -split "`n"
        $newLines = @()
        $inserted = $false

        # Determine go_package based on filename
        $goPackage = if ($file.Name -eq 'openmetrics_data_model.proto') {
            'github.com/streamlit/streamlit/go-backend/proto/openmetrics'
        } else {
            'github.com/streamlit/streamlit/go-backend/proto'
        }

        foreach ($line in $lines) {
            $newLines += $line
            if ($line -match 'syntax\s*=\s*"proto3"' -and -not $inserted) {
                $newLines += ''
                $newLines += "option go_package = `"$goPackage`";"
                $inserted = $true
            }
        }

        $newContent = $newLines -join "`n"
        Set-Content $file.FullName -Value $newContent -NoNewline
        Write-Host "Added go_package to $($file.Name) -> $goPackage"
    }
}

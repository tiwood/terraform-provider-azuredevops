[CmdletBinding()]
param (
)

$script:PSDefaultParameterValues = @{
    '*:Confirm'           = $false
    '*:ErrorAction'       = 'Stop'
}

. (Join-Path -Path $PSScriptRoot -ChildPath 'commons.ps1' -Resolve)

Write-Host "Executing unit tests"
Push-Location -Path $SOURCE_DIR
try {
    go test -v ./... 
    if ($LASTEXITCODE) {
        throw "Build finished in error due to failed tests"
    } 
}
finally {
    Pop-Location
}

#!/usr/bin/env pwsh
[CmdletBinding()]
param (
    [Parameter()]
    [ValidateNotNullOrEmpty()]
    [string]
    $TestFilter,

    [Parameter()]
    [ValidateNotNullOrEmpty()]
    [string[]]
    $Tag = 'all'
)

$script:PSDefaultParameterValues = @{
    '*:Confirm'           = $false
    '*:ErrorAction'       = 'Stop'
}

. (Join-Path -Path $PSScriptRoot -ChildPath 'commons.ps1' -Resolve)

Write-Host "Executing unit tests"
Push-Location -Path $SOURCE_DIR
try {
    $argv = @(
        'test',
        '-v',
        '-mod=vendor'
    )
    if ($TestFilter) {
        $argv += @('-run', $TestFilter)
    }
    if ($Tag -and 0 -lt $Tag.Length) {
        $argv += @('-tags', [string]::Join(' ', $Tag))
    }
    go @argv ./...
    if ($LASTEXITCODE) {
        throw "Build finished in error due to failed tests"
    }
}
finally {
    Pop-Location
}

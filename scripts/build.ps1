[CmdletBinding()]
param (
    [switch]
    $SkipTests
)

$script:PSDefaultParameterValues = @{
    '*:Confirm'           = $false
    '*:ErrorAction'       = 'Stop'
}

. (Join-Path -Path $PSScriptRoot -ChildPath 'commons.ps1' -Resolve)

function clean() {
    Write-Host "Cleaning $BUILD_DIR"
    if (Test-Path -Path $BUILD_DIR) {
        Remove-Item -Recurse -Force -Path $BUILD_DIR
    }
    $null = New-Item -ItemType Container -Path $BUILD_DIR
}

function compile() {
    $NAME = Get-Content -Raw -Path $PROVIDER_NAME_FILE
    $VERSION = Get-Content -Raw -Path $PROVIDER_VERSION_FILE

    $BUILD_ARTIFACT="terraform-provider-${NAME}_v${VERSION}"

    Write-Host "Attempting to build $BUILD_ARTIFACT"
    Push-Location -Path $SOURCE_DIR
    try {
        go mod download 
        if ($LASTEXITCODE) {
            throw "Failed to download modules"
        }
        go build -o "$BUILD_DIR/$BUILD_ARTIFACT"
        if ($LASTEXITCODE) {
            throw "Build failed"
        }
    }
    finally {
        Pop-Location
    }
}

function clean_and_build() {
    clean
    compile
    if (-not $SkipTests) {
        & (Join-Path -Path $PSScriptRoot -ChildPath 'unittest.ps1' -Resolve)
    }
    Write-Host "Build finished successfully"
}

clean_and_build

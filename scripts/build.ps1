. (Join-Path -Path $PSScriptRoot -ChildPath 'commons.ps1')

function clean() {
    Write-Host "Cleaning $BUILD_DIR"
    if (Test-Path -Path $BUILD_DIR) {
        Remove-Item -Recurse -Force -Path $BUILD_DIR
    }
    $null = New-Item -ItemType Container -Path $BUILD_DIR
}

function compile() {
    $NAME=Get-Content -Raw -Path $PROVIDER_NAME_FILE
    $VERSION=Get-Content -Raw -Path $PROVIDER_VERSION_FILE

    $BUILD_ARTIFACT="terraform-provider-${NAME}_v${VERSION}"

    Write-Host "Attempting to build $BUILD_ARTIFACT"
    Push-Location -Path $SOURCE_DIR
    try {
        go mod download 
        go build -o "$BUILD_DIR/$BUILD_ARTIFACT"
    }
    finally {
        Pop-Location
    }
}

function clean_and_build() {
    clean
    #$(dirname $0)/unittest.sh
    compile
    Write-Host "Build finished successfully"
}

clean_and_build

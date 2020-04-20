[CmdletBinding()]
param (
    [Parameter()]
    [ValidateNotNullOrEmpty()]
    [string]
    $TestFilter = '^TestAcc',

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

Write-Host "Executing acceptance tests"
Push-Location -Path $SOURCE_DIR
try {
    # This is similar to the unit test command aside from the following:
    #   - TF_ACC=1 is a flag that will enable the acceptance tests. This flag is
    #     documented here:
    #       https://www.terraform.io/docs/extend/testing/acceptance-tests/index.html#running-acceptance-tests
    #
    #   - A `-run` parameter is used to target *only* tests starting with `TestAcc`. This prefix is
    #     recommended by Hashicorp and is documented here:
    #       https://www.terraform.io/docs/extend/testing/acceptance-tests/index.html#test-files
    #
    # Using build tags as test filter: https://stackoverflow.com/a/24036237
    $env:TF_ACC=1
    $argv = @(
        'test',
        '-mod=vendor',
        '-v'
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

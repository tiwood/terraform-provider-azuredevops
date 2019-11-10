[CmdletBinding()]
param (
    [Parameter()]
    [string]
    $PluginsDirectory = [IO.Path]::Combine($HOME, '.terraform.d', 'plugins')
)

$script:PSDefaultParameterValues = @{
    '*:Confirm'           = $false
    '*:ErrorAction'       = 'Stop'
}

. (Join-Path -Path $PSScriptRoot -ChildPath 'commons.ps1' -Resolve)

if (Test-Path -Path $PluginsDirectory) {
    Write-Verbose -Message "Terraform Plugins directory [$PluginsDirectory] already exists"
}
else {
    Write-Verbose -Message "Creating Terraform Plugins directory [$PluginsDirectory]"
    $null = New-Item -Path $PluginsDirectory -ItemType Directory
}

Write-Host "Installing provider to $PluginsDirectory"
Copy-Item -Path (Join-Path -Path $BUILD_DIR -ChildPath '*') -Destination $PluginsDirectory -Force

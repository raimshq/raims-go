[CmdletBinding()]
param(
    [string]$PythonPath = 'python',
    [switch]$Race
)

$ErrorActionPreference = 'Stop'
$sdkRoot = [System.IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..'))
$workspaceRoot = [System.IO.Directory]::GetParent($sdkRoot).FullName
$validator = Join-Path $workspaceRoot 'Developer-Toolkit\scripts\release-preflight\go_sdk.py'
if (-not [System.IO.File]::Exists($validator)) {
    throw 'O validador canônico do Developer Toolkit não está disponível no workspace. Consulte CONTRIBUTING.md.'
}
$arguments = @($validator, '--source', $sdkRoot, '--workspace', $workspaceRoot)
if ($Race) { $arguments += '--race' }
& $PythonPath @arguments
if ($LASTEXITCODE -ne 0) {
    throw 'A validação do SDK falhou. Consulte o relatório informado pelo validador.'
}

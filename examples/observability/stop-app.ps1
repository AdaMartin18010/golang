$ErrorActionPreference = 'Stop'

$scriptRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $scriptRoot

$pidFile = Join-Path $scriptRoot ".app.pid"

if (!(Test-Path $pidFile)) {
    Write-Host ".app.pid not found. nothing to stop." -ForegroundColor Yellow
    exit 0
}

$procId = Get-Content $pidFile | Select-Object -First 1
if (-not $procId) { Write-Host "empty pid file." -ForegroundColor Yellow; exit 0 }

try {
    Stop-Process -Id $procId -Force -ErrorAction Stop
    Remove-Item $pidFile -Force -ErrorAction SilentlyContinue
    Write-Host "stopped app. pid=$procId" -ForegroundColor Green
} catch {
    Write-Host "process not found. pid=$procId" -ForegroundColor Yellow
    Remove-Item $pidFile -Force -ErrorAction SilentlyContinue
}


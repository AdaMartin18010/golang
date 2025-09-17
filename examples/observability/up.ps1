$ErrorActionPreference = 'Stop'
$scriptRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $scriptRoot

Write-Host "docker compose up -d ..." -ForegroundColor Cyan
docker compose up -d
docker compose ps


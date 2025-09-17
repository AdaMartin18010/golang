$ErrorActionPreference = 'Stop'

$scriptRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $scriptRoot

$env:OTEL_SERVICE_NAME = if ([string]::IsNullOrEmpty($env:OTEL_SERVICE_NAME)) { 'example-observability-app' } else { $env:OTEL_SERVICE_NAME }
$env:OTEL_ENV = if ([string]::IsNullOrEmpty($env:OTEL_ENV)) { 'dev' } else { $env:OTEL_ENV }
# 如果未显式设置导出端点，默认指向本机 Collector gRPC 4317，并启用 insecure
if ([string]::IsNullOrEmpty($env:OTEL_EXPORTER_OTLP_ENDPOINT)) { $env:OTEL_EXPORTER_OTLP_ENDPOINT = 'localhost:4317' }
if ([string]::IsNullOrEmpty($env:OTEL_EXPORTER_OTLP_INSECURE)) { $env:OTEL_EXPORTER_OTLP_INSECURE = 'true' }

$pidFile = Join-Path $scriptRoot ".app.pid"

function Test-ProcessAlive($procId) {
    try { $null = Get-Process -Id $procId -ErrorAction Stop; return $true } catch { return $false }
}

if (Test-Path $pidFile) {
    $procId = Get-Content $pidFile | Select-Object -First 1
    if ($procId -and (Test-ProcessAlive $procId)) {
        Write-Host "app already running. pid=$procId" -ForegroundColor Yellow
        exit 0
    } else {
        Remove-Item $pidFile -Force -ErrorAction SilentlyContinue
    }
}

Write-Host "starting app (go run ./app) ..." -ForegroundColor Cyan
$p = Start-Process -FilePath 'go' -ArgumentList 'run ./app' -WorkingDirectory $scriptRoot -WindowStyle Hidden -PassThru
Set-Content -Path $pidFile -Value $p.Id -Encoding ASCII
Write-Host "started. pid=$($p.Id)" -ForegroundColor Green


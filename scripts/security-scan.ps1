# 安全扫描脚本 (PowerShell)
# 功能：运行多种安全扫描工具

Write-Host "开始安全扫描..." -ForegroundColor Green

# 检查工具是否安装
function Test-Tool {
    param([string]$ToolName)
    
    $tool = Get-Command $ToolName -ErrorAction SilentlyContinue
    if (-not $tool) {
        Write-Host "警告: $ToolName 未安装，跳过" -ForegroundColor Yellow
        return $false
    }
    return $true
}

# 1. Gosec 安全扫描
if (Test-Tool "gosec") {
    Write-Host "运行 Gosec 安全扫描..." -ForegroundColor Yellow
    gosec -fmt json -out gosec-report.json ./...
    gosec ./...
}

# 2. Trivy 漏洞扫描
if (Test-Tool "trivy") {
    Write-Host "运行 Trivy 漏洞扫描..." -ForegroundColor Yellow
    trivy fs --format json --output trivy-report.json .
    trivy fs .
}

# 3. 依赖漏洞扫描
if (Test-Tool "govulncheck") {
    Write-Host "运行 Go 漏洞检查..." -ForegroundColor Yellow
    govulncheck ./...
}

# 4. 检查硬编码密钥
Write-Host "检查硬编码密钥..." -ForegroundColor Yellow
$passwordMatches = Select-String -Path "*.go" -Pattern "password\s*=" -Exclude "*_test.go" | Where-Object { $_.Line -notmatch "test|example" }
if ($passwordMatches) {
    Write-Host "警告: 发现可能的硬编码密码" -ForegroundColor Red
}

$secretMatches = Select-String -Path "*.go" -Pattern "secret\s*=" -Exclude "*_test.go" | Where-Object { $_.Line -notmatch "test|example" }
if ($secretMatches) {
    Write-Host "警告: 发现可能的硬编码密钥" -ForegroundColor Red
}

# 5. 检查敏感信息
Write-Host "检查敏感信息..." -ForegroundColor Yellow
$apiKeyMatches = Select-String -Path "*.go" -Pattern "api[_-]key" -CaseSensitive:$false -Exclude "*_test.go"
if ($apiKeyMatches) {
    Write-Host "警告: 发现可能的 API 密钥" -ForegroundColor Red
}

$privateKeyMatches = Select-String -Path "*.go" -Pattern "private[_-]key" -CaseSensitive:$false -Exclude "*_test.go"
if ($privateKeyMatches) {
    Write-Host "警告: 发现可能的私钥" -ForegroundColor Red
}

Write-Host "安全扫描完成!" -ForegroundColor Green


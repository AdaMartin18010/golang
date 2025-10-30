# Markdown 全面修复脚本
# 依次执行所有修复操作

param (
    [string]$Path = "docs"
)

Write-Host "=" * 80 -ForegroundColor Cyan
Write-Host "🔧 Markdown 全面修复工具" -ForegroundColor Cyan
Write-Host "=" * 80 -ForegroundColor Cyan
Write-Host ""

$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
$totalFixed = 0

# 1. 修复重复版本信息
Write-Host "1️⃣  修复重复版本信息块..." -ForegroundColor Yellow
Write-Host "-" * 80 -ForegroundColor Gray
& "$scriptPath\fix-all-duplicates.ps1" -Path $Path
Write-Host ""

# 2. 修复目录链接
Write-Host "2️⃣  修复目录链接..." -ForegroundColor Yellow
Write-Host "-" * 80 -ForegroundColor Gray
& "$scriptPath\fix-toc-links.ps1" -Path $Path
Write-Host ""

# 3. 修复基础格式问题
Write-Host "3️⃣  修复基础格式问题..." -ForegroundColor Yellow
Write-Host "-" * 80 -ForegroundColor Gray
& "$scriptPath\check-markdown-basic.ps1" -Path $Path -Fix
Write-Host ""

# 4. 使用 markdownlint 自动修复（如果可用）
Write-Host "4️⃣  使用 markdownlint 自动修复..." -ForegroundColor Yellow
Write-Host "-" * 80 -ForegroundColor Gray
try {
    $null = Get-Command markdownlint -ErrorAction Stop
    & "$scriptPath\check-markdown-format.ps1" -Path $Path -Fix
} catch {
    Write-Host "⚠️  markdownlint-cli 未安装，跳过此步骤" -ForegroundColor Yellow
    Write-Host "   安装命令: npm install -g markdownlint-cli" -ForegroundColor Gray
}
Write-Host ""

# 5. 最终验证
Write-Host "5️⃣  最终验证..." -ForegroundColor Yellow
Write-Host "-" * 80 -ForegroundColor Gray
& "$scriptPath\check-markdown-basic.ps1" -Path $Path
Write-Host ""

Write-Host "=" * 80 -ForegroundColor Green
Write-Host "✅ 全面修复完成！" -ForegroundColor Green
Write-Host "=" * 80 -ForegroundColor Green
Write-Host ""
Write-Host "📝 建议: 请查看修复结果，确认所有更改符合预期" -ForegroundColor Cyan

# Markdown 格式检查脚本
# 使用 markdownlint-cli 检查所有 Markdown 文件

param (
    [string]$Path = "docs",
    [switch]$Fix = $false,
    [switch]$Verbose = $false
)

Write-Host "=" * 80 -ForegroundColor Cyan
Write-Host "📋 Markdown 格式检查工具" -ForegroundColor Cyan
Write-Host "=" * 80 -ForegroundColor Cyan
Write-Host ""

# 检查是否安装了 markdownlint-cli
$hasMarkdownlint = $false
try {
    $null = Get-Command markdownlint -ErrorAction Stop
    $hasMarkdownlint = $true
    Write-Host "✅ 检测到 markdownlint-cli" -ForegroundColor Green
} catch {
    Write-Host "⚠️  未检测到 markdownlint-cli" -ForegroundColor Yellow
}

# 如果没有安装，提供安装指南
if (-not $hasMarkdownlint) {
    Write-Host ""
    Write-Host "📦 请先安装 markdownlint-cli:" -ForegroundColor Yellow
    Write-Host "   npm install -g markdownlint-cli" -ForegroundColor White
    Write-Host ""
    Write-Host "   或使用 pnpm:" -ForegroundColor White
    Write-Host "   pnpm add -g markdownlint-cli" -ForegroundColor White
    Write-Host ""
    Write-Host "   或使用 yarn:" -ForegroundColor White
    Write-Host "   yarn global add markdownlint-cli" -ForegroundColor White
    Write-Host ""

    # 提供备用方案：手动检查
    Write-Host "💡 正在使用内置规则进行基础检查..." -ForegroundColor Cyan
    Write-Host ""

    # 执行基础检查
    & "$PSScriptRoot\check-markdown-basic.ps1" -Path $Path
    exit 0
}

# 使用 markdownlint-cli 进行检查
Write-Host ""
Write-Host "🔍 扫描路径: $Path" -ForegroundColor Cyan
Write-Host ""

$configFile = ".markdownlint.json"
if (-not (Test-Path $configFile)) {
    $configFile = ".markdownlint.jsonc"
}

if (Test-Path $configFile) {
    Write-Host "⚙️  使用配置文件: $configFile" -ForegroundColor Cyan
} else {
    Write-Host "⚠️  未找到配置文件，使用默认规则" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "-" * 80 -ForegroundColor Gray
Write-Host ""

# 构建命令
$cmd = "markdownlint"
$args = @()

if ($Fix) {
    $args += "--fix"
    Write-Host "🔧 自动修复模式已启用" -ForegroundColor Green
    Write-Host ""
}

if (Test-Path $configFile) {
    $args += "--config", $configFile
}

$args += "$Path/**/*.md"

# 执行检查
try {
    if ($Verbose) {
        Write-Host "执行命令: $cmd $($args -join ' ')" -ForegroundColor Gray
        Write-Host ""
    }

    $output = & $cmd $args 2>&1
    $exitCode = $LASTEXITCODE

    if ($exitCode -eq 0) {
        Write-Host ""
        Write-Host "=" * 80 -ForegroundColor Green
        Write-Host "✅ 所有 Markdown 文件格式正确！" -ForegroundColor Green
        Write-Host "=" * 80 -ForegroundColor Green
    } else {
        Write-Host $output
        Write-Host ""
        Write-Host "=" * 80 -ForegroundColor Yellow
        Write-Host "⚠️  发现格式问题，请查看上面的详细信息" -ForegroundColor Yellow
        if (-not $Fix) {
            Write-Host "💡 提示: 使用 -Fix 参数自动修复可修复的问题" -ForegroundColor Cyan
            Write-Host "   示例: .\scripts\check-markdown-format.ps1 -Fix" -ForegroundColor White
        }
        Write-Host "=" * 80 -ForegroundColor Yellow
        exit $exitCode
    }
} catch {
    Write-Host "❌ 执行检查时出错: $_" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "📊 检查完成时间: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" -ForegroundColor Gray

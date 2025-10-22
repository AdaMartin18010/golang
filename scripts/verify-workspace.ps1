# Go 1.25.3 Workspace 环境验证脚本
# 用于快速验证 Workspace 配置是否正确

Write-Host "==================================" -ForegroundColor Cyan
Write-Host "Go 1.25.3 Workspace 环境验证" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""

$allPassed = $true

# 1. 检查 Go 版本
Write-Host "【1】检查 Go 版本..." -ForegroundColor Yellow
try {
    $goVersion = go version
    Write-Host "    ✅ $goVersion" -ForegroundColor Green
    
    if ($goVersion -match "go1\.25\.3") {
        Write-Host "    ✅ Go 1.25.3 版本正确" -ForegroundColor Green
    } else {
        Write-Host "    ⚠️  警告: Go 版本不是 1.25.3" -ForegroundColor Yellow
    }
} catch {
    Write-Host "    ❌ Go 未安装或未在 PATH 中" -ForegroundColor Red
    $allPassed = $false
}
Write-Host ""

# 2. 检查 go.work 文件
Write-Host "【2】检查 go.work 文件..." -ForegroundColor Yellow
if (Test-Path "go.work") {
    Write-Host "    ✅ go.work 文件存在" -ForegroundColor Green
    
    # 显示 go.work 内容
    Write-Host "    📄 内容预览:" -ForegroundColor Cyan
    $content = Get-Content "go.work" -First 10
    $content | ForEach-Object { Write-Host "       $_" -ForegroundColor Gray }
} else {
    Write-Host "    ❌ go.work 文件不存在" -ForegroundColor Red
    $allPassed = $false
}
Write-Host ""

# 3. 检查 examples/go.mod
Write-Host "【3】检查 examples/go.mod..." -ForegroundColor Yellow
if (Test-Path "examples/go.mod") {
    Write-Host "    ✅ examples/go.mod 存在" -ForegroundColor Green
    
    $modContent = Get-Content "examples/go.mod" -First 5
    $goLine = $modContent | Where-Object { $_ -match "^go " }
    if ($goLine) {
        Write-Host "    📄 $goLine" -ForegroundColor Gray
    }
} else {
    Write-Host "    ❌ examples/go.mod 不存在" -ForegroundColor Red
    $allPassed = $false
}
Write-Host ""

# 4. 执行 go work sync
Write-Host "【4】执行 go work sync..." -ForegroundColor Yellow
try {
    $syncResult = go work sync 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "    ✅ go work sync 执行成功" -ForegroundColor Green
    } else {
        Write-Host "    ❌ go work sync 失败" -ForegroundColor Red
        Write-Host "    错误: $syncResult" -ForegroundColor Red
        $allPassed = $false
    }
} catch {
    Write-Host "    ❌ go work sync 执行出错: $_" -ForegroundColor Red
    $allPassed = $false
}
Write-Host ""

# 5. 检查模块列表
Write-Host "【5】检查模块列表..." -ForegroundColor Yellow
try {
    $modules = go list -m
    if ($modules) {
        Write-Host "    ✅ 主模块: $modules" -ForegroundColor Green
    } else {
        Write-Host "    ⚠️  无法获取模块列表" -ForegroundColor Yellow
    }
} catch {
    Write-Host "    ❌ 无法列出模块: $_" -ForegroundColor Red
    $allPassed = $false
}
Write-Host ""

# 6. 检查关键文档
Write-Host "【6】检查关键文档..." -ForegroundColor Yellow
$keyDocs = @(
    "START_HERE.md",
    "🚀-立即开始-3分钟上手.md",
    "📌-项目状态总览.md",
    "📚-Workspace文档索引.md"
)

$docsFound = 0
foreach ($doc in $keyDocs) {
    if (Test-Path $doc) {
        $docsFound++
    }
}

if ($docsFound -eq $keyDocs.Count) {
    Write-Host "    ✅ 所有关键文档存在 ($docsFound/$($keyDocs.Count))" -ForegroundColor Green
} else {
    Write-Host "    ⚠️  部分文档缺失 ($docsFound/$($keyDocs.Count))" -ForegroundColor Yellow
}
Write-Host ""

# 7. 检查迁移脚本
Write-Host "【7】检查迁移脚本..." -ForegroundColor Yellow
if (Test-Path "scripts/migrate-to-workspace.ps1") {
    Write-Host "    ✅ 迁移脚本存在" -ForegroundColor Green
} else {
    Write-Host "    ⚠️  迁移脚本不存在" -ForegroundColor Yellow
}
Write-Host ""

# 最终总结
Write-Host "==================================" -ForegroundColor Cyan
Write-Host "验证总结" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan

if ($allPassed) {
    Write-Host "✅ 所有关键验证通过！" -ForegroundColor Green
    Write-Host ""
    Write-Host "🚀 你现在可以开始使用 Workspace 了！" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "快速开始：" -ForegroundColor Yellow
    Write-Host "  1. go work sync" -ForegroundColor White
    Write-Host "  2. go test ./examples/..." -ForegroundColor White
    Write-Host ""
    Write-Host "查看文档：" -ForegroundColor Yellow
    Write-Host "  - START_HERE.md" -ForegroundColor White
    Write-Host "  - 🚀-立即开始-3分钟上手.md" -ForegroundColor White
} else {
    Write-Host "⚠️  部分验证未通过" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "请检查上述错误项并修复" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "获取帮助：" -ForegroundColor Yellow
    Write-Host "  - 查看 START_HERE.md" -ForegroundColor White
    Write-Host "  - 查看 📌-项目状态总览.md" -ForegroundColor White
}

Write-Host ""
Write-Host "验证脚本执行完成" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan


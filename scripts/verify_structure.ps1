#!/usr/bin/env pwsh
<#
.SYNOPSIS
    验证项目结构是否符合重组规范
.DESCRIPTION
    检查文档代码分离、目录职责等规范是否被正确遵守
.EXAMPLE
    .\verify_structure.ps1
#>

$ErrorActionPreference = "Stop"
$WarningPreference = "Continue"

Write-Host "🔍 开始验证项目结构..." -ForegroundColor Cyan
Write-Host ""

# 统计变量
$script:ErrorCount = 0
$script:WarningCount = 0
$script:PassCount = 0

# 辅助函数
function Test-Rule {
    param(
        [string]$Name,
        [scriptblock]$Test,
        [string]$ErrorMsg,
        [string]$PassMsg
    )
    
    Write-Host "➤ $Name..." -NoNewline
    
    try {
        $result = & $Test
        if ($result) {
            Write-Host " ✅" -ForegroundColor Green
            if ($PassMsg) {
                Write-Host "  └─ $PassMsg" -ForegroundColor Gray
            }
            $script:PassCount++
        } else {
            Write-Host " ❌" -ForegroundColor Red
            Write-Host "  └─ $ErrorMsg" -ForegroundColor Yellow
            $script:ErrorCount++
        }
    } catch {
        Write-Host " ⚠️" -ForegroundColor Yellow
        Write-Host "  └─ 检查失败: $_" -ForegroundColor Yellow
        $script:WarningCount++
    }
}

Write-Host "=" * 60
Write-Host "📋 规则1: 文档代码分离" -ForegroundColor Blue
Write-Host "=" * 60

# 检查 docs/ 目录是否有代码文件
Test-Rule -Name "docs/ 目录无 .go 文件" -Test {
    $goFiles = Get-ChildItem -Path "docs" -Recurse -Filter "*.go" -ErrorAction SilentlyContinue
    $goFiles.Count -eq 0
} -ErrorMsg "发现 $($goFiles.Count) 个 .go 文件，应该移至 examples/" -PassMsg "docs/ 目录纯文档 ✓"

Test-Rule -Name "docs/ 目录无 go.mod 文件" -Test {
    $modFiles = Get-ChildItem -Path "docs" -Recurse -Filter "go.mod" -ErrorAction SilentlyContinue
    $modFiles.Count -eq 0
} -ErrorMsg "发现 $($modFiles.Count) 个 go.mod 文件" -PassMsg "无 go.mod 文件 ✓"

Test-Rule -Name "docs/ 目录无可执行文件" -Test {
    $exeFiles = Get-ChildItem -Path "docs" -Recurse -Filter "*.exe" -ErrorAction SilentlyContinue
    $exeFiles.Count -eq 0
} -ErrorMsg "发现 $($exeFiles.Count) 个可执行文件" -PassMsg "无可执行文件 ✓"

Write-Host ""
Write-Host "=" * 60
Write-Host "📋 规则2: 根目录清洁" -ForegroundColor Blue
Write-Host "=" * 60

# 检查根目录临时文件
Test-Rule -Name "根目录无 Phase 报告" -Test {
    $phaseFiles = Get-ChildItem -Path "." -Filter "Phase-*.md" -ErrorAction SilentlyContinue
    $phaseFiles.Count -eq 0
} -ErrorMsg "发现 $($phaseFiles.Count) 个 Phase 报告，应移至 reports/phase-reports/" -PassMsg "无 Phase 报告 ✓"

Test-Rule -Name "根目录无临时报告文件" -Test {
    $tempReports = Get-ChildItem -Path "." -Filter "*报告*.md" -ErrorAction SilentlyContinue
    $tempReports = $tempReports | Where-Object { $_.Name -notlike "RESTRUCTURE.md" }
    $tempReports.Count -eq 0
} -ErrorMsg "发现 $($tempReports.Count) 个临时报告" -PassMsg "无临时报告 ✓"

Test-Rule -Name "根目录文档数量合理" -Test {
    $rootMdFiles = Get-ChildItem -Path "." -Filter "*.md" -File -ErrorAction SilentlyContinue
    $count = $rootMdFiles.Count
    # 应该在 10-15 个之间
    ($count -ge 8) -and ($count -le 20)
} -ErrorMsg "根目录有 $count 个 .md 文件，建议保持在 8-20 个" -PassMsg "文件数量: $count 个 ✓"

Write-Host ""
Write-Host "=" * 60
Write-Host "📋 规则3: 目录职责" -ForegroundColor Blue
Write-Host "=" * 60

# 检查必要目录存在
Test-Rule -Name "存在 docs/ 目录" -Test {
    Test-Path "docs"
} -ErrorMsg "缺少 docs/ 目录" -PassMsg "存在 ✓"

Test-Rule -Name "存在 examples/ 目录" -Test {
    Test-Path "examples"
} -ErrorMsg "缺少 examples/ 目录" -PassMsg "存在 ✓"

Test-Rule -Name "存在 reports/ 目录" -Test {
    Test-Path "reports"
} -ErrorMsg "缺少 reports/ 目录" -PassMsg "存在 ✓"

Test-Rule -Name "存在 archive/ 目录" -Test {
    Test-Path "archive"
} -ErrorMsg "缺少 archive/ 目录" -PassMsg "存在 ✓"

Test-Rule -Name "存在 scripts/ 目录" -Test {
    Test-Path "scripts"
} -ErrorMsg "缺少 scripts/ 目录" -PassMsg "存在 ✓"

Write-Host ""
Write-Host "=" * 60
Write-Host "📋 规则4: 关键文件" -ForegroundColor Blue
Write-Host "=" * 60

$keyFiles = @(
    "README.md",
    "RESTRUCTURE.md",
    "MIGRATION_GUIDE.md",
    "CONTRIBUTING.md",
    "FAQ.md",
    "LICENSE"
)

foreach ($file in $keyFiles) {
    Test-Rule -Name "存在 $file" -Test {
        Test-Path $file
    } -ErrorMsg "缺少 $file" -PassMsg "存在 ✓"
}

Write-Host ""
Write-Host "=" * 60
Write-Host "📋 规则5: examples/ 结构" -ForegroundColor Blue
Write-Host "=" * 60

$exampleDirs = @(
    "advanced",
    "concurrency",
    "go125",
    "modern-features",
    "testing-framework"
)

foreach ($dir in $exampleDirs) {
    Test-Rule -Name "存在 examples/$dir/" -Test {
        Test-Path "examples/$dir"
    } -ErrorMsg "缺少 examples/$dir/ 目录" -PassMsg "存在 ✓"
}

Write-Host ""
Write-Host "=" * 60
Write-Host "📋 规则6: 代码质量" -ForegroundColor Blue
Write-Host "=" * 60

Test-Rule -Name "examples/ 中代码可编译" -Test {
    Push-Location "examples"
    try {
        $result = go build ./... 2>&1
        $LASTEXITCODE -eq 0
    } finally {
        Pop-Location
    }
} -ErrorMsg "代码编译失败" -PassMsg "编译通过 ✓"

Write-Host ""
Write-Host "=" * 60
Write-Host "📊 验证结果统计" -ForegroundColor Magenta
Write-Host "=" * 60

$total = $script:PassCount + $script:ErrorCount + $script:WarningCount

Write-Host ""
Write-Host "通过: $script:PassCount / $total" -ForegroundColor Green
Write-Host "失败: $script:ErrorCount / $total" -ForegroundColor Red
Write-Host "警告: $script:WarningCount / $total" -ForegroundColor Yellow
Write-Host ""

if ($script:ErrorCount -eq 0) {
    Write-Host "✅ 项目结构验证通过！" -ForegroundColor Green
    Write-Host "项目结构符合重组规范。" -ForegroundColor Gray
    exit 0
} else {
    Write-Host "❌ 项目结构验证失败！" -ForegroundColor Red
    Write-Host "发现 $script:ErrorCount 个问题需要修复。" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "请参考 RESTRUCTURE.md 了解项目结构规范。" -ForegroundColor Gray
    exit 1
}


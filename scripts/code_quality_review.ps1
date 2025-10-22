<#
.SYNOPSIS
    Comprehensive code quality review script for the Golang project.

.DESCRIPTION
    This script performs a complete code quality audit including:
    - Code formatting checks (go fmt)
    - Code quality analysis (go vet)
    - Linter checks (golangci-lint)
    - Test coverage analysis
    - Code complexity metrics
    - Generates detailed quality reports

.PARAMETER OutputDir
    Directory where quality reports will be saved.
    Defaults to "reports/quality".

.PARAMETER Fix
    If specified, automatically fix formatting issues.

.PARAMETER Detailed
    If specified, generates detailed reports with all findings.

.EXAMPLE
    .\scripts\code_quality_review.ps1
    Run basic quality review and save report to default location.

.EXAMPLE
    .\scripts\code_quality_review.ps1 -Fix -Detailed
    Run detailed review and automatically fix formatting issues.

.NOTES
    Requires Go 1.25.3+ and golangci-lint to be installed.
#>

[CmdletBinding()]
param (
    [string]$OutputDir = "reports/quality",
    [switch]$Fix,
    [switch]$Detailed
)

# Color functions
function Write-ColorOutput {
    param(
        [string]$Message,
        [string]$Color = "White"
    )
    Write-Host $Message -ForegroundColor $Color
}

function Write-Section {
    param([string]$Title)
    Write-Host ""
    Write-Host "═══════════════════════════════════════════════════════" -ForegroundColor Cyan
    Write-Host " $Title" -ForegroundColor Yellow
    Write-Host "═══════════════════════════════════════════════════════" -ForegroundColor Cyan
}

function Write-Success {
    param([string]$Message)
    Write-ColorOutput "✅ $Message" "Green"
}

function Write-Warning {
    param([string]$Message)
    Write-ColorOutput "⚠️  $Message" "Yellow"
}

function Write-Error-Custom {
    param([string]$Message)
    Write-ColorOutput "❌ $Message" "Red"
}

function Write-Info {
    param([string]$Message)
    Write-ColorOutput "ℹ️  $Message" "Cyan"
}

# Initialize
$ErrorActionPreference = "Continue"
$StartTime = Get-Date
$ReportFile = Join-Path $OutputDir "quality-review-$(Get-Date -Format 'yyyy-MM-dd-HHmmss').md"
$IssueCount = 0
$WarningCount = 0
$PassCount = 0

# Create output directory
if (-not (Test-Path $OutputDir)) {
    New-Item -Path $OutputDir -ItemType Directory -Force | Out-Null
    Write-Info "Created output directory: $OutputDir"
}

# Initialize report
$Report = @"
# 🔍 代码质量审查报告

> **生成时间**: $(Get-Date -Format "yyyy-MM-dd HH:mm:ss")  
> **项目版本**: v2.0.0

---

## 📋 审查概要

"@

Write-Section "🚀 开始代码质量审查"
Write-Info "项目路径: $(Get-Location)"
Write-Info "报告将保存到: $ReportFile"

# ============================================================================
# 1. 代码格式检查 (go fmt)
# ============================================================================
Write-Section "1️⃣  代码格式检查 (go fmt)"

$Report += @"

### 1. 代码格式检查

**工具**: go fmt

"@

try {
    if ($Fix) {
        Write-Info "正在自动修复格式问题..."
        $FmtOutput = go fmt ./...
    } else {
        Write-Info "正在检查代码格式..."
        $FmtOutput = gofmt -l .
    }
    
    if ($LASTEXITCODE -eq 0 -and [string]::IsNullOrWhiteSpace($FmtOutput)) {
        Write-Success "所有文件格式正确"
        $Report += "**结果**: ✅ 通过 - 所有文件格式正确`n`n"
        $PassCount++
    } else {
        $UnformattedFiles = $FmtOutput -split "`n" | Where-Object { $_ -ne "" }
        $Count = $UnformattedFiles.Count
        if ($Fix) {
            Write-Warning "已修复 $Count 个文件的格式"
            $Report += "**结果**: ⚠️  已修复 - 修复了 $Count 个文件的格式`n`n"
        } else {
            Write-Warning "发现 $Count 个文件需要格式化"
            $Report += "**结果**: ⚠️  需要修复 - $Count 个文件需要格式化`n`n"
            $Report += "**文件列表**:`n````n$($UnformattedFiles -join "`n")`n```

`n`n"
        }
        $WarningCount++
    }
} catch {
    Write-Error-Custom "go fmt 检查失败: $_"
    $Report += "**结果**: ❌ 失败 - 检查过程出错`n`n"
    $IssueCount++
}

# ============================================================================
# 2. 代码质量分析 (go vet)
# ============================================================================
Write-Section "2️⃣  代码质量分析 (go vet)"

$Report += @"
### 2. 代码质量分析

**工具**: go vet

"@

try {
    Write-Info "正在运行 go vet..."
    $VetOutput = go vet ./... 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "go vet 检查通过"
        $Report += "**结果**: ✅ 通过 - 未发现问题`n`n"
        $PassCount++
    } else {
        Write-Warning "go vet 发现问题"
        $Report += "**结果**: ⚠️  发现问题`n`n"
        $Report += "**详细信息**:`n````n$VetOutput`n```

`n`n"
        $WarningCount++
    }
} catch {
    Write-Error-Custom "go vet 检查失败: $_"
    $Report += "**结果**: ❌ 失败 - 检查过程出错`n`n"
    $IssueCount++
}

# ============================================================================
# 3. Linter检查 (golangci-lint)
# ============================================================================
Write-Section "3️⃣  Linter检查 (golangci-lint)"

$Report += @"
### 3. Linter检查

**工具**: golangci-lint

"@

try {
    # Check if golangci-lint is installed
    $LintInstalled = Get-Command golangci-lint -ErrorAction SilentlyContinue
    
    if ($LintInstalled) {
        Write-Info "正在运行 golangci-lint..."
        $LintOutput = golangci-lint run ./... 2>&1
        
        if ($LASTEXITCODE -eq 0 -and [string]::IsNullOrWhiteSpace($LintOutput)) {
            Write-Success "golangci-lint 检查通过"
            $Report += "**结果**: ✅ 通过 - 未发现问题`n`n"
            $PassCount++
        } else {
            Write-Warning "golangci-lint 发现问题"
            $IssueLines = ($LintOutput -split "`n" | Where-Object { $_ -match ":" }).Count
            $Report += "**结果**: ⚠️  发现 $IssueLines 个问题`n`n"
            
            if ($Detailed) {
                $Report += "**详细信息**:`n```

$LintOutput`n```
`n`n"
            } else {
                $TopIssues = ($LintOutput -split "`n" | Select-Object -First 20) -join "`n"
                $Report += "**主要问题** (前20个):`n```
$TopIssues`n```
`n`n"
            }
            $WarningCount++
        }
    } else {
        Write-Warning "golangci-lint 未安装，跳过检查"
        $Report += "**结果**: ⚠️  跳过 - golangci-lint 未安装`n`n"
        $Report += "**安装命令**: `curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b `$(go env GOPATH)/bin`

`n`n"
        $WarningCount++
    }
} catch {
    Write-Error-Custom "golangci-lint 检查失败: $_"
    $Report += "**结果**: ❌ 失败 - 检查过程出错`n`n"
    $IssueCount++
}

# ============================================================================
# 4. 测试覆盖率分析
# ============================================================================
Write-Section "4️⃣  测试覆盖率分析"

$Report += @"
### 4. 测试覆盖率分析

**工具**: go test -cover

"@

try {
    Write-Info "正在运行测试并生成覆盖率报告..."
    $CoverageFile = Join-Path $OutputDir "coverage.out"
    $CoverageHtml = Join-Path $OutputDir "coverage.html"
    
    # Run tests with coverage
    $TestOutput = go test -coverprofile=$CoverageFile ./... 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        # Generate coverage report
        $CoverageReport = go tool cover -func=$CoverageFile
        
        # Calculate average coverage
        $CoverageLines = $CoverageReport -split "`n" | Where-Object { $_ -match "total:" }
        if ($CoverageLines) {
            $TotalCoverage = $CoverageLines[0] -replace ".*total:.*?([0-9.]+)%.*", '$1'
            $CoveragePct = [double]$TotalCoverage
            
            Write-Info "总体覆盖率: $CoveragePct%"
            
            # Generate HTML report
            go tool cover -html=$CoverageFile -o $CoverageHtml
            
            if ($CoveragePct -ge 80) {
                Write-Success "测试覆盖率良好: $CoveragePct%"
                $Report += "**结果**: ✅ 优秀 - 覆盖率 $CoveragePct% (≥80%)`n`n"
                $PassCount++
            } elseif ($CoveragePct -ge 60) {
                Write-Warning "测试覆盖率中等: $CoveragePct%"
                $Report += "**结果**: ⚠️  中等 - 覆盖率 $CoveragePct% (60-80%)`n`n"
                $WarningCount++
            } else {
                Write-Warning "测试覆盖率较低: $CoveragePct%"
                $Report += "**结果**: ⚠️  偏低 - 覆盖率 $CoveragePct% (<60%)`n`n"
                $WarningCount++
            }
            
            $Report += "**覆盖率报告**: [coverage.html]($CoverageHtml)`n`n"
            
            # Module-level coverage
            $Report += "**模块覆盖率**:`n`n"
            $ModuleCoverage = $CoverageReport -split "`n" | Where-Object { $_ -match "%" -and $_ -notmatch "total:" } | Select-Object -First 10
            $Report += "```
$($ModuleCoverage -join "`n")`n```
`n`n"
        } else {
            Write-Warning "无法解析覆盖率数据"
            $Report += "**结果**: ⚠️  无法解析覆盖率数据`n`n"
            $WarningCount++
        }
    } else {
        Write-Warning "部分测试失败"
        $Report += "**结果**: ⚠️  部分测试失败`n`n"
        $Report += "**测试输出**:`n```
$TestOutput`n```
`n`n"
        $WarningCount++
    }
} catch {
    Write-Error-Custom "测试覆盖率分析失败: $_"
    $Report += "**结果**: ❌ 失败 - 分析过程出错`n`n"
    $IssueCount++
}

# ============================================================================
# 5. 模块完整性检查
# ============================================================================
Write-Section "5️⃣  模块完整性检查"

$Report += @"
### 5. 模块完整性检查

**工具**: go mod verify, go mod tidy

"@

try {
    Write-Info "正在检查模块完整性..."
    
    # Check go mod verify
    $VerifyOutput = go mod verify 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Success "模块完整性验证通过"
        $Report += "**go mod verify**: ✅ 通过`n`n"
        $PassCount++
    } else {
        Write-Warning "模块完整性验证失败"
        $Report += "**go mod verify**: ⚠️  失败`n`n"
        $Report += "```
$VerifyOutput`n```
`n`n"
        $WarningCount++
    }
    
    # Check if go mod tidy would make changes
    Write-Info "正在检查 go.mod 是否需要整理..."
    $TidyCheck = go mod tidy -diff 2>&1
    if ([string]::IsNullOrWhiteSpace($TidyCheck)) {
        Write-Success "go.mod 已是最新状态"
        $Report += "**go mod tidy**: ✅ 无需更改`n`n"
        $PassCount++
    } else {
        Write-Warning "go.mod 需要整理"
        $Report += "**go mod tidy**: ⚠️  需要运行 'go mod tidy'`n`n"
        $WarningCount++
    }
} catch {
    Write-Error-Custom "模块完整性检查失败: $_"
    $Report += "**结果**: ❌ 失败 - 检查过程出错`n`n"
    $IssueCount++
}

# ============================================================================
# 6. 代码统计
# ============================================================================
Write-Section "6️⃣  代码统计"

$Report += @"
### 6. 代码统计

"@

try {
    Write-Info "正在统计代码行数..."
    
    # Count Go files
    $GoFiles = Get-ChildItem -Path . -Filter "*.go" -Recurse | Where-Object { $_.FullName -notmatch "\\vendor\\" }
    $GoFileCount = $GoFiles.Count
    
    # Count lines of code
    $TotalLines = 0
    $GoFiles | ForEach-Object {
        $TotalLines += (Get-Content $_.FullName | Measure-Object -Line).Lines
    }
    
    # Count test files
    $TestFiles = $GoFiles | Where-Object { $_.Name -match "_test\.go$" }
    $TestFileCount = $TestFiles.Count
    
    # Count packages
    $Packages = Get-ChildItem -Path . -Filter "*.go" -Recurse | ForEach-Object { Split-Path $_.DirectoryName -Leaf } | Sort-Object -Unique
    $PackageCount = $Packages.Count
    
    Write-Success "统计完成"
    
    $Report += @"
| 指标 | 数量 |
|------|------|
| Go文件数 | $GoFileCount |
| 代码总行数 | $TotalLines |
| 测试文件数 | $TestFileCount |
| 包数量 | $PackageCount |
| 平均每文件行数 | $([math]::Round($TotalLines / $GoFileCount, 2)) |

"@
    $PassCount++
} catch {
    Write-Error-Custom "代码统计失败: $_"
    $Report += "**结果**: ❌ 失败 - 统计过程出错`n`n"
    $IssueCount++
}

# ============================================================================
# 7. 生成总结
# ============================================================================
Write-Section "📊 生成审查总结"

$EndTime = Get-Date
$Duration = $EndTime - $StartTime

# Calculate quality score
$TotalChecks = $PassCount + $WarningCount + $IssueCount
$QualityScore = if ($TotalChecks -gt 0) { 
    [math]::Round(($PassCount / $TotalChecks) * 100, 2) 
} else { 
    0 
}

# Quality rating
$Rating = if ($QualityScore -ge 90) { "优秀 🌟🌟🌟🌟🌟" }
          elseif ($QualityScore -ge 75) { "良好 🌟🌟🌟🌟" }
          elseif ($QualityScore -ge 60) { "中等 🌟🌟🌟" }
          else { "需要改进 🌟🌟" }

# Add summary to report
$SummaryTable = @"
| 检查项 | 通过 | 警告 | 失败 |
|--------|------|------|------|
| 数量 | $PassCount | $WarningCount | $IssueCount |

**质量评分**: $QualityScore / 100

**评级**: $Rating

**审查时长**: $($Duration.TotalSeconds) 秒

---

## 🎯 改进建议

"@

$Report = $Report.Insert($Report.IndexOf("---") + 4, $SummaryTable)

# Add recommendations
$Recommendations = @()

if ($WarningCount -gt 0 -or $IssueCount -gt 0) {
    if ($Report -match "需要格式化") {
        $Recommendations += "- 运行 `go fmt ./...` 修复代码格式"
    }
    if ($Report -match "go vet 发现问题") {
        $Recommendations += "- 修复 `go vet` 发现的代码问题"
    }
    if ($Report -match "golangci-lint 发现") {
        $Recommendations += "- 修复 `golangci-lint` 发现的代码质量问题"
    }
    if ($Report -match "覆盖率.*[0-5][0-9]") {
        $Recommendations += "- 提高测试覆盖率至60%以上"
    }
    if ($Report -match "go mod tidy") {
        $Recommendations += "- 运行 `go mod tidy` 整理模块依赖"
    }
}

if ($Recommendations.Count -gt 0) {
    $Report += ($Recommendations -join "`n") + "`n`n"
} else {
    $Report += "✅ 代码质量良好，无需特别改进`n`n"
}

$Report += @"
---

## 📚 相关资源

- [代码规范](CONTRIBUTING.md)
- [安全指南](SECURITY.md)
- [测试指南](docs/README.md)

---

**审查完成时间**: $(Get-Date -Format "yyyy-MM-dd HH:mm:ss")
"@

# Save report
$Report | Out-File -FilePath $ReportFile -Encoding UTF8

# Display summary
Write-Section "✅ 审查完成"
Write-Host ""
Write-Host "📊 审查结果汇总:" -ForegroundColor Cyan
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-ColorOutput "  ✅ 通过: $PassCount" "Green"
Write-ColorOutput "  ⚠️  警告: $WarningCount" "Yellow"
Write-ColorOutput "  ❌ 失败: $IssueCount" "Red"
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-ColorOutput "  📊 质量评分: $QualityScore / 100" "Cyan"
Write-ColorOutput "  🏆 评级: $Rating" "Cyan"
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host ""
Write-Success "报告已保存到: $ReportFile"
Write-Host ""

# Return exit code based on issues
if ($IssueCount -gt 0) {
    exit 1
} else {
    exit 0
}


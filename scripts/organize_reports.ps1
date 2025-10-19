# 整理报告文件脚本
# 将根目录的报告文件移动到 reports/ 目录

$rootDir = "E:\_src\golang"
$reportsDir = "$rootDir\reports"

# 创建子目录
$subDirs = @(
    "$reportsDir\phase-reports",
    "$reportsDir\daily-summaries",
    "$reportsDir\analysis-reports",
    "$reportsDir\code-quality",
    "$reportsDir\archive"
)

foreach ($dir in $subDirs) {
    if (-not (Test-Path $dir)) {
        New-Item -ItemType Directory -Path $dir -Force | Out-Null
        Write-Host "✓ 创建目录: $dir"
    }
}

# 定义文件移动映射
$fileMappings = @{
    "phase-reports" = @(
        "*Phase-*.md",
        "*阶段*.md",
        "*里程碑*.md"
    )
    "daily-summaries" = @(
        "*今日工作*.md",
        "*工作总结*.md",
        "*工作完成*.md"
    )
    "analysis-reports" = @(
        "*分析*.md",
        "*评估*.md",
        "*评价*.md",
        "*对标*.md"
    )
    "code-quality" = @(
        "code_quality_report*.md",
        "*代码*.md",
        "*格式化*.md",
        "*验证*.md"
    )
    "archive" = @(
        "*2025*.md",
        "*推进*.md",
        "*执行*.md",
        "*清单*.md",
        "*计划*.md",
        "*总结*.md",
        "*简报*.md",
        "*报告*.md",
        "*指南*.md"
    )
}

# 移动文件函数
function Move-ReportFiles {
    param (
        [string]$targetSubDir,
        [string[]]$patterns
    )
    
    $movedCount = 0
    foreach ($pattern in $patterns) {
        $files = Get-ChildItem -Path $rootDir -Filter $pattern -File -ErrorAction SilentlyContinue
        foreach ($file in $files) {
            # 跳过已经在reports目录的文件
            if ($file.FullName -like "$reportsDir\*") {
                continue
            }
            
            # 跳过特殊文件
            if ($file.Name -in @("README.md", "README_EN.md", "CONTRIBUTING.md", "CONTRIBUTING_EN.md", 
                                  "CHANGELOG.md", "LICENSE", "FAQ.md", "EXAMPLES.md", "EXAMPLES_EN.md",
                                  "QUICK_START.md", "QUICK_START_EN.md", "CODE_OF_CONDUCT.md",
                                  "PROJECT_STRUCTURE_NEW.md", "RELEASE_NOTES.md", "RELEASE_v2.0.0.md")) {
                continue
            }
            
            $targetPath = "$reportsDir\$targetSubDir\$($file.Name)"
            
            # 如果目标文件已存在，跳过
            if (Test-Path $targetPath) {
                Write-Host "⊙ 跳过（已存在）: $($file.Name)"
                continue
            }
            
            try {
                Move-Item -Path $file.FullName -Destination $targetPath -Force
                Write-Host "→ 移动: $($file.Name) -> $targetSubDir/"
                $movedCount++
            } catch {
                Write-Host "✗ 失败: $($file.Name) - $($_.Exception.Message)" -ForegroundColor Red
            }
        }
    }
    return $movedCount
}

Write-Host "`n=== 开始整理报告文件 ===`n" -ForegroundColor Cyan

$totalMoved = 0

# 按优先级移动文件（先具体后通用，避免重复移动）
Write-Host "📁 移动 Phase 报告..." -ForegroundColor Yellow
$totalMoved += Move-ReportFiles "phase-reports" $fileMappings["phase-reports"]

Write-Host "`n📁 移动每日工作总结..." -ForegroundColor Yellow
$totalMoved += Move-ReportFiles "daily-summaries" $fileMappings["daily-summaries"]

Write-Host "`n📁 移动分析报告..." -ForegroundColor Yellow
$totalMoved += Move-ReportFiles "analysis-reports" $fileMappings["analysis-reports"]

Write-Host "`n📁 移动代码质量报告..." -ForegroundColor Yellow
$totalMoved += Move-ReportFiles "code-quality" $fileMappings["code-quality"]

Write-Host "`n📁 移动其他报告到归档..." -ForegroundColor Yellow
$totalMoved += Move-ReportFiles "archive" $fileMappings["archive"]

Write-Host "`n=== 整理完成 ===" -ForegroundColor Green
Write-Host "✓ 总共移动了 $totalMoved 个文件" -ForegroundColor Green

# 生成报告
Write-Host "`n📊 生成整理报告..." -ForegroundColor Yellow
$reportContent = @"
# 报告文件整理记录

> **整理日期**: $(Get-Date -Format "yyyy-MM-dd HH:mm:ss")  
> **总移动文件数**: $totalMoved

## 📁 目录结构

``````text
reports/
├── phase-reports/      # Phase 1-5 阶段报告
├── daily-summaries/    # 每日工作总结
├── analysis-reports/   # 分析和评估报告
├── code-quality/       # 代码质量报告
└── archive/           # 历史归档文件
``````

## 📋 文件分类

### Phase 报告
- 所有 Phase-*.md 文件
- 阶段性总结和里程碑报告

### 每日工作总结
- 今日工作相关的总结文件
- 工作完成记录

### 分析报告
- 项目分析和评估
- 生态对标分析

### 代码质量报告
- code_quality_report*.md
- 代码验证和格式化报告

### 归档文件
- 历史计划和总结
- 推进记录和简报

## 🔍 查找文件

如需查找特定报告，请查看对应子目录。

---

**整理完成日期**: $(Get-Date -Format "yyyy年MM月dd日")
"@

$reportPath = "$reportsDir\ORGANIZATION_REPORT.md"
$reportContent | Out-File -FilePath $reportPath -Encoding UTF8
Write-Host "✓ 整理报告已生成: $reportPath" -ForegroundColor Green

Write-Host "`n✨ 根目录现在更加整洁！" -ForegroundColor Cyan


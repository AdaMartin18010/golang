# 清理无关文件脚本
# 功能: 清理与本项目内容无关的文件
# 日期: 2025-11-11

param(
    [switch]$DryRun = $false,
    [switch]$Verbose = $false
)

Write-Host "🧹 清理无关文件脚本" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

# 创建归档目录
$archiveDir = "archive/documentation-cleanup-2025-11-11"
if (-not $DryRun) {
    New-Item -ItemType Directory -Path $archiveDir -Force | Out-Null
    Write-Host "📁 创建归档目录: $archiveDir" -ForegroundColor Green
}

# 需要清理的文件模式
$patternsToClean = @(
    # 根目录临时报告（2025-11-11）
    "*-2025-11-11.md",
    # 历史报告（2025-10-28, 2025-10-30）
    "*-2025-10-28.md",
    "*-2025-10-30.md",
    # docs目录格式梳理报告
    "docs/*-2025-10-29.md",
    # 备份文件
    "*.bak"
)

# 需要保留的核心文件
$coreFiles = @(
    "README.md",
    "README_EN.md",
    "README-MARKDOWN-TOOLS.md",
    "CONTRIBUTING.md",
    "CONTRIBUTING_EN.md",
    "CODE_OF_CONDUCT.md",
    "SECURITY.md",
    "LICENSE",
    "CHANGELOG.md",
    "go.work",
    "go.work.sum",
    "codecov.yml",
    "lychee.toml",
    "cspell.json"
)

$filesToClean = @()
$filesToKeep = @()

# 扫描根目录
Write-Host "📁 扫描根目录..." -ForegroundColor Cyan
$rootFiles = Get-ChildItem -Path . -Filter "*.md" -File

foreach ($file in $rootFiles) {
    $fileName = $file.Name
    
    # 检查是否是核心文件
    $isCore = $false
    foreach ($core in $coreFiles) {
        if ($fileName -eq $core) {
            $isCore = $true
            break
        }
    }
    
    if ($isCore) {
        $filesToKeep += $file
        if ($Verbose) {
            Write-Host "  ✓ 保留: $fileName" -ForegroundColor Gray
        }
    } elseif ($fileName -match '-2025-11-11\.md$' -or 
              $fileName -match '-2025-10-28\.md$' -or 
              $fileName -match '-2025-10-30\.md$') {
        $filesToClean += $file
    }
}

# 扫描docs目录
Write-Host "📁 扫描docs目录..." -ForegroundColor Cyan
$docsFiles = Get-ChildItem -Path "docs" -Filter "*-2025-10-29.md" -File
foreach ($file in $docsFiles) {
    $filesToClean += $file
}

# 扫描备份文件
Write-Host "📁 扫描备份文件..." -ForegroundColor Cyan
$bakFiles = Get-ChildItem -Path . -Filter "*.bak" -Recurse -File
foreach ($file in $bakFiles) {
    $filesToClean += $file
}

# 显示统计
Write-Host ""
Write-Host "📊 清理统计" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host "  保留文件: $($filesToKeep.Count)" -ForegroundColor Green
Write-Host "  清理文件: $($filesToClean.Count)" -ForegroundColor Yellow
Write-Host ""

# 显示要清理的文件
if ($filesToClean.Count -gt 0) {
    Write-Host "📋 待清理文件列表:" -ForegroundColor Yellow
    foreach ($file in $filesToClean) {
        Write-Host "  - $($file.FullName)" -ForegroundColor Gray
    }
    Write-Host ""
}

# 执行清理
if ($filesToClean.Count -gt 0) {
    if ($DryRun) {
        Write-Host "⚠ 预览模式：以下文件将被移动到归档目录" -ForegroundColor Yellow
        Write-Host "  归档目录: $archiveDir" -ForegroundColor Yellow
    } else {
        Write-Host "🔄 开始清理..." -ForegroundColor Cyan
        
        foreach ($file in $filesToClean) {
            try {
                $relativePath = $file.FullName.Replace((Get-Location).Path + "\", "")
                $targetPath = Join-Path $archiveDir $file.Name
                
                # 如果目标文件已存在，添加序号
                $counter = 1
                $originalTarget = $targetPath
                while (Test-Path $targetPath) {
                    $nameWithoutExt = [System.IO.Path]::GetFileNameWithoutExtension($file.Name)
                    $ext = [System.IO.Path]::GetExtension($file.Name)
                    $targetPath = Join-Path $archiveDir "$nameWithoutExt-$counter$ext"
                    $counter++
                }
                
                Move-Item -Path $file.FullName -Destination $targetPath -Force
                Write-Host "  ✓ 已移动: $($file.Name)" -ForegroundColor Green
            } catch {
                Write-Host "  ❌ 错误: $($file.Name) - $_" -ForegroundColor Red
            }
        }
        
        Write-Host ""
        Write-Host "✅ 清理完成" -ForegroundColor Green
    }
} else {
    Write-Host "✅ 没有需要清理的文件" -ForegroundColor Green
}

# 创建归档说明
if (-not $DryRun -and $filesToClean.Count -gt 0) {
    $readmeContent = @"
# 文档清理归档说明

**归档日期**: 2025年11月11日
**归档原因**: 清理与本项目内容无关的临时报告文档

## 📋 归档内容

### 归档文件数

- **临时报告文档**: $($filesToClean.Count) 个
- **归档位置**: archive/documentation-cleanup-2025-11-11/

### 归档文件类型

1. **文档梳理工作报告** (2025-11-11)
   - 格式、目录、结构梳理报告
   - 内容语义梳理报告
   - 文件夹结构分析报告
   - 链接梳理报告

2. **历史报告** (2025-10-28, 2025-10-30)
   - 项目归档报告
   - 文档优化报告
   - 格式梳理报告

3. **文档格式梳理报告** (2025-10-29)
   - docs/目录下的格式梳理报告

4. **备份文件**
   - .bak 文件

## 🎯 归档原则

这些文件是文档梳理工作过程中产生的临时报告，已完成其历史使命，现归档保存。

## 📚 查看归档文件

如需查看归档文件，请访问：
\`archive/documentation-cleanup-2025-11-11/\`

---

**归档时间**: 2025年11月11日
"@
    
    $readmePath = Join-Path $archiveDir "README.md"
    [System.IO.File]::WriteAllText($readmePath, $readmeContent, [System.Text.Encoding]::UTF8)
    Write-Host "📝 已创建归档说明: $readmePath" -ForegroundColor Green
}

Write-Host ""


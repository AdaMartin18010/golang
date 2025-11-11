# 统一目录格式脚本
# 功能: 统一所有文档的目录格式
# 日期: 2025-11-11

param(
    [string]$DocsPath = "docs",
    [switch]$DryRun = $false,
    [switch]$Verbose = $false
)

Write-Host "📋 统一目录格式脚本" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

# 标准目录标题
$StandardTocTitle = "## 📋 目录"

# 检测目录标题
function Get-TocTitle {
    param([string]$Content)

    $lines = $Content -split "`n"

    foreach ($line in $lines) {
        $trimmed = $line.Trim()
        if ($trimmed -match '^##+\s+.*目录') {
            return $trimmed
        }
    }

    return $null
}

# 统一目录标题
function Standardize-TocTitle {
    param([string]$Content)

    $lines = $Content -split "`n"
    $newLines = @()
    $changed = $false

    foreach ($line in $lines) {
        $trimmed = $line.Trim()

        # 检测目录标题并替换
        if ($trimmed -match '^##+\s+.*目录' -and $trimmed -ne $StandardTocTitle) {
            $newLines += $StandardTocTitle
            $changed = $true
        } else {
            $newLines += $line
        }
    }

    if ($changed) {
        return $newLines -join "`n"
    }

    return $Content
}

# 主处理逻辑
$mdFiles = Get-ChildItem -Path $DocsPath -Filter "*.md" -Recurse
$totalFiles = $mdFiles.Count
$processedFiles = 0
$standardizedFiles = 0

Write-Host "📁 找到 $totalFiles 个Markdown文件" -ForegroundColor Cyan
Write-Host ""

foreach ($file in $mdFiles) {
    $processedFiles++

    try {
        $content = Get-Content $file.FullName -Raw -Encoding UTF8
        $tocTitle = Get-TocTitle -Content $content

        if ($tocTitle -and $tocTitle -ne $StandardTocTitle) {
            Write-Host "[$processedFiles/$totalFiles] 📋 $($file.Name)" -ForegroundColor Cyan
            Write-Host "  当前: $tocTitle" -ForegroundColor Yellow
            Write-Host "  标准: $StandardTocTitle" -ForegroundColor Green

            if (-not $DryRun) {
                $newContent = Standardize-TocTitle -Content $content
                [System.IO.File]::WriteAllText($file.FullName, $newContent, [System.Text.Encoding]::UTF8)
                Write-Host "  ✅ 已标准化" -ForegroundColor Green
            } else {
                Write-Host "  ⚠ 预览模式（未实际修改）" -ForegroundColor Yellow
            }

            $standardizedFiles++
        } else {
            if ($Verbose) {
                if ($tocTitle) {
                    Write-Host "[$processedFiles/$totalFiles] ✓ $($file.Name) - 已是标准格式" -ForegroundColor Gray
                } else {
                    Write-Host "[$processedFiles/$totalFiles] - $($file.Name) - 无目录" -ForegroundColor Gray
                }
            }
        }
    } catch {
        Write-Host "  ❌ 错误: $_" -ForegroundColor Red
    }
}

# 总结
Write-Host ""
Write-Host "================================" -ForegroundColor Cyan
Write-Host "📊 标准化总结" -ForegroundColor Cyan
Write-Host "  处理文件: $processedFiles/$totalFiles" -ForegroundColor White
Write-Host "  已标准化: $standardizedFiles" -ForegroundColor White

if ($DryRun) {
    Write-Host "  ⚠ 这是预览模式，未实际修改文件" -ForegroundColor Yellow
    Write-Host "  使用 -DryRun:`$false 来实际标准化" -ForegroundColor Yellow
} else {
    Write-Host "  ✅ 文件已更新" -ForegroundColor Green
}

Write-Host ""

# 统一标题格式脚本
# 功能: 统一文档标题格式，移除上下文前缀
# 日期: 2025-11-11

param(
    [string]$DocsPath = "docs",
    [switch]$DryRun = $false,
    [switch]$Verbose = $false
)

Write-Host "📝 统一标题格式脚本" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

# 需要移除的上下文前缀模式
$ContextPrefixes = @(
    '^Go语言基础\s*-\s*',
    '^Go基础\s*-\s*',
    '^Go\s*-\s*',
    '^Go语言\s*-\s*',
    '^Go 1\.25\.3\s*-\s*',
    '^Go-1\.25\.3\s*-\s*'
)

# 清理标题中的上下文前缀
function Clean-Title {
    param([string]$Title)

    $cleaned = $Title

    foreach ($prefix in $ContextPrefixes) {
        if ($cleaned -match $prefix) {
            $cleaned = $cleaned -replace $prefix, ''
            break
        }
    }

    # 清理首尾空格
    $cleaned = $cleaned.Trim()

    return $cleaned
}

# 统一文档标题
function Standardize-DocumentTitle {
    param([string]$Content)

    $lines = $Content -split "`n"
    $newLines = @()
    $changed = $false

    foreach ($line in $lines) {
        # 检测一级标题（文档标题）
        if ($line -match '^#\s+(.+)$') {
            $originalTitle = $matches[1]
            $cleanedTitle = Clean-Title $originalTitle

            if ($cleanedTitle -ne $originalTitle) {
                $newLines += "# $cleanedTitle"
                $changed = $true

                if ($Verbose) {
                    Write-Host "  📝 标题: '$originalTitle' → '$cleanedTitle'" -ForegroundColor Yellow
                }
            } else {
                $newLines += $line
            }
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
        $newContent = Standardize-DocumentTitle -Content $content

        if ($newContent -ne $content) {
            Write-Host "[$processedFiles/$totalFiles] 📝 $($file.Name)" -ForegroundColor Cyan

            if (-not $DryRun) {
                [System.IO.File]::WriteAllText($file.FullName, $newContent, [System.Text.Encoding]::UTF8)
                Write-Host "  ✅ 已标准化" -ForegroundColor Green
            } else {
                Write-Host "  ⚠ 预览模式（未实际修改）" -ForegroundColor Yellow
            }

            $standardizedFiles++
        } else {
            if ($Verbose) {
                Write-Host "[$processedFiles/$totalFiles] ✓ $($file.Name) - 无需修改" -ForegroundColor Gray
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

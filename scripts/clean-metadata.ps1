# 清理元数据额外字段脚本
# 功能: 移除非标准元数据字段，统一为标准格式
# 日期: 2025-11-11

param(
    [string]$DocsPath = "docs",
    [switch]$DryRun = $false,
    [switch]$Verbose = $false
)

Write-Host "🧹 清理元数据额外字段脚本" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

# 标准元数据格式（只保留这三行）
$StandardMetadata = @"
**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3
"@

# 需要移除的额外字段
$ExtraFields = @(
    '\*\*字数\*\*',
    '\*\*代码示例\*\*',
    '\*\*实战案例\*\*',
    '\*\*适用人群\*\*',
    '\*\*问题数\*\*',
    '\*\*文档类型\*\*',
    '\*\*Go版本\*\*',
    '\*\*难度\*\*',
    '\*\*预计阅读\*\*',
    '\*\*基准日期\*\*'
)

# 清理元数据
function Clean-Metadata {
    param([string]$Content)

    $lines = $Content -split "`n"
    $newLines = @()
    $inMetadata = $false
    $metadataStart = -1
    $metadataEnd = -1
    $hasStandardMetadata = $false

    # 查找元数据区域（标题后到第二个分隔线或目录前）
    $firstSeparator = -1
    $secondSeparator = -1

    for ($i = 0; $i -lt $lines.Count; $i++) {
        $line = $lines[$i].Trim()

        # 检测标题
        if ($line -match '^#+\s+' -and $metadataStart -eq -1) {
            $metadataStart = $i + 1
            continue
        }

        if ($metadataStart -gt 0) {
            # 检测第一个分隔线
            if ($line -eq '---' -and $firstSeparator -eq -1) {
                $firstSeparator = $i
                continue
            }

            # 检测第二个分隔线或目录标题
            if ($firstSeparator -gt 0) {
                if ($line -eq '---' -and $secondSeparator -eq -1) {
                    $secondSeparator = $i
                    $metadataEnd = $i
                    break
                }
                if ($line -match '^##+\s+.*目录' -and $secondSeparator -eq -1) {
                    $metadataEnd = $i
                    break
                }
            }

            # 检测标准元数据字段
            if ($line -match '\*\*版本\*\*' -or
                $line -match '\*\*更新日期\*\*' -or
                $line -match '\*\*适用于\*\*') {
                $hasStandardMetadata = $true
            }
        }
    }

    # 如果没有找到结束位置，使用第一个分隔线后20行作为范围
    if ($metadataEnd -eq -1 -and $firstSeparator -gt 0) {
        $metadataEnd = [Math]::Min($firstSeparator + 20, $lines.Count)
    }

    # 如果没有标准元数据，不处理
    if (-not $hasStandardMetadata) {
        return $Content
    }

    # 重建内容
    $changed = $false
    for ($i = 0; $i -lt $lines.Count; $i++) {
        # 跳过元数据区域中的额外字段
        if ($metadataStart -gt 0 -and
            $i -ge $metadataStart -and
            $i -lt $metadataEnd) {

            $line = $lines[$i]
            $shouldRemove = $false

            # 检查是否是额外字段
            foreach ($field in $ExtraFields) {
                if ($line -match $field) {
                    $shouldRemove = $true
                    $changed = $true
                    break
                }
            }

            # 如果是标准字段或空行，保留
            if (-not $shouldRemove -and
                ($line -match '\*\*版本\*\*' -or
                 $line -match '\*\*更新日期\*\*' -or
                 $line -match '\*\*适用于\*\*' -or
                 $line.Trim() -eq '')) {
                $newLines += $line
            }
        } else {
            $newLines += $lines[$i]
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
$cleanedFiles = 0

Write-Host "📁 找到 $totalFiles 个Markdown文件" -ForegroundColor Cyan
Write-Host ""

foreach ($file in $mdFiles) {
    $processedFiles++

    try {
        $content = Get-Content $file.FullName -Raw -Encoding UTF8
        $newContent = Clean-Metadata -Content $content

        if ($newContent -ne $content) {
            Write-Host "[$processedFiles/$totalFiles] 🧹 $($file.Name)" -ForegroundColor Cyan

            if (-not $DryRun) {
                [System.IO.File]::WriteAllText($file.FullName, $newContent, [System.Text.Encoding]::UTF8)
                Write-Host "  ✅ 已清理" -ForegroundColor Green
            } else {
                Write-Host "  ⚠ 预览模式（未实际修改）" -ForegroundColor Yellow
            }

            $cleanedFiles++
        } else {
            if ($Verbose) {
                Write-Host "[$processedFiles/$totalFiles] ✓ $($file.Name) - 无需清理" -ForegroundColor Gray
            }
        }
    } catch {
        Write-Host "  ❌ 错误: $_" -ForegroundColor Red
    }
}

# 总结
Write-Host ""
Write-Host "================================" -ForegroundColor Cyan
Write-Host "📊 清理总结" -ForegroundColor Cyan
Write-Host "  处理文件: $processedFiles/$totalFiles" -ForegroundColor White
Write-Host "  已清理: $cleanedFiles" -ForegroundColor White

if ($DryRun) {
    Write-Host "  ⚠ 这是预览模式，未实际修改文件" -ForegroundColor Yellow
    Write-Host "  使用 -DryRun:`$false 来实际清理" -ForegroundColor Yellow
} else {
    Write-Host "  ✅ 文件已更新" -ForegroundColor Green
}

Write-Host ""

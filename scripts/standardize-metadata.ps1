# 文档元数据标准化脚本
# 功能: 统一Markdown文档的元数据格式
# 日期: 2025-11-11

param(
    [string]$DocsPath = "docs",
    [switch]$DryRun = $false,
    [switch]$Verbose = $false
)

Write-Host "📝 文档元数据标准化脚本" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

# 标准元数据格式
$StandardMetadata = @"
**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3
"@

# 检测元数据格式
function Get-MetadataFormat {
    param([string]$Content)

    $metadata = @{
        HasVersion = $false
        HasDate = $false
        HasGoVersion = $false
        Format = "unknown"
        StartLine = -1
        EndLine = -1
    }

    $lines = $Content -split "`n"
    $inMetadata = $false
    $metadataStart = -1

    $maxLines = if ($lines.Count -lt 20) { $lines.Count } else { 20 }
    for ($i = 0; $i -lt $maxLines; $i++) {
        $line = $lines[$i].Trim()

        # 检测元数据开始（通常在标题后）
        if ($line -match '^#+\s+') {
            $inMetadata = $true
            $metadataStart = $i + 1
            continue
        }

        if ($inMetadata) {
            # 检测分隔线（---）
            if ($line -eq '---') {
                $metadata.EndLine = $i
                break
            }

            # 检测元数据字段
            if ($line -match '\*\*版本\*\*' -or $line -match '\*\*Version\*\*') {
                $metadata.HasVersion = $true
            }
            if ($line -match '\*\*更新日期\*\*' -or $line -match '\*\*Last Updated\*\*' -or $line -match '\*\*基准日期\*\*') {
                $metadata.HasDate = $true
            }
            if ($line -match '\*\*适用于\*\*' -or $line -match '\*\*Go版本\*\*' -or $line -match '\*\*Go Version\*\*') {
                $metadata.HasGoVersion = $true
            }

            # 如果遇到空行且已有元数据，结束
            if ($line -eq '' -and ($metadata.HasVersion -or $metadata.HasDate)) {
                $metadata.EndLine = $i
                break
            }
        }
    }

    $metadata.StartLine = $metadataStart
    if ($metadata.EndLine -eq -1) {
        $metadata.EndLine = $metadataStart + 5
    }

    # 判断格式类型
    if ($metadata.HasVersion -and $metadata.HasDate -and $metadata.HasGoVersion) {
        $metadata.Format = "standard"
    } elseif ($metadata.HasVersion -and $metadata.HasDate) {
        $metadata.Format = "simplified"
    } elseif ($metadata.HasVersion -or $metadata.HasDate) {
        $metadata.Format = "partial"
    } else {
        $metadata.Format = "missing"
    }

    return $metadata
}

# 标准化元数据
function Standardize-Metadata {
    param(
        [string]$Content,
        [hashtable]$MetadataInfo
    )

    $lines = $Content -split "`n"
    $newLines = @()
    $inserted = $false

    for ($i = 0; $i -lt $lines.Count; $i++) {
        # 找到标题后的位置
        if (-not $inserted -and $lines[$i] -match '^#+\s+') {
            $newLines += $lines[$i]

            # 添加空行
            if ($i + 1 -lt $lines.Count -and $lines[$i + 1].Trim() -ne '') {
                $newLines += ""
            }

            # 插入标准元数据
            $newLines += $StandardMetadata
            $newLines += ""
            $newLines += "---"
            $newLines += ""
            $inserted = $true

            # 跳过旧的元数据行
            if ($MetadataInfo.StartLine -gt 0) {
                $skipUntil = $MetadataInfo.EndLine
                if ($skipUntil -gt $i) {
                    $i = $skipUntil
                    continue
                }
            }
        } else {
            # 跳过旧元数据范围内的行
            if ($MetadataInfo.StartLine -gt 0 -and
                $i -ge $MetadataInfo.StartLine -and
                $i -le $MetadataInfo.EndLine -and
                -not $inserted) {
                continue
            }

            $newLines += $lines[$i]
        }
    }

    return $newLines -join "`n"
}

# 主处理逻辑
$mdFiles = Get-ChildItem -Path $DocsPath -Filter "*.md" -Recurse
$totalFiles = $mdFiles.Count
$processedFiles = 0
$standardizedFiles = 0
$skippedFiles = 0

Write-Host "📁 找到 $totalFiles 个Markdown文件" -ForegroundColor Cyan
Write-Host ""

foreach ($file in $mdFiles) {
    $processedFiles++

    try {
        $content = Get-Content $file.FullName -Raw -Encoding UTF8
        $metadataInfo = Get-MetadataFormat -Content $content

        if ($metadataInfo.Format -eq "standard") {
            if ($Verbose) {
                Write-Host "[$processedFiles/$totalFiles] ✓ $($file.Name) - 已是标准格式" -ForegroundColor Gray
            }
            $skippedFiles++
            continue
        }

        Write-Host "[$processedFiles/$totalFiles] 📝 $($file.Name)" -ForegroundColor Cyan
        Write-Host "  格式: $($metadataInfo.Format)" -ForegroundColor Yellow

        if (-not $DryRun) {
            $newContent = Standardize-Metadata -Content $content -MetadataInfo $metadataInfo
            [System.IO.File]::WriteAllText($file.FullName, $newContent, [System.Text.Encoding]::UTF8)
            Write-Host "  ✅ 已标准化" -ForegroundColor Green
        } else {
            Write-Host "  ⚠ 预览模式（未实际修改）" -ForegroundColor Yellow
        }

        $standardizedFiles++
    } catch {
        Write-Host "  ❌ 错误: $_" -ForegroundColor Red
    }

    Write-Host ""
}

# 总结
Write-Host "================================" -ForegroundColor Cyan
Write-Host "📊 标准化总结" -ForegroundColor Cyan
Write-Host "  处理文件: $processedFiles/$totalFiles" -ForegroundColor White
Write-Host "  已标准化: $standardizedFiles" -ForegroundColor White
Write-Host "  已跳过: $skippedFiles (已是标准格式)" -ForegroundColor White

if ($DryRun) {
    Write-Host "  ⚠ 这是预览模式，未实际修改文件" -ForegroundColor Yellow
    Write-Host "  使用 -DryRun:`$false 来实际标准化" -ForegroundColor Yellow
} else {
    Write-Host "  ✅ 文件已更新" -ForegroundColor Green
}

Write-Host ""

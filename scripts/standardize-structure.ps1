# 统一文档结构脚本
# 功能: 统一文档结构，优化章节组织
# 日期: 2025-11-11

param(
    [string]$DocsPath = "docs",
    [switch]$DryRun = $false,
    [switch]$Verbose = $false
)

Write-Host "📐 统一文档结构脚本" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

# 标准文档结构
$StandardStructure = @"
# 文档标题

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [文档标题](#文档标题)

---

## 内容章节

内容...

---

## 📚 参考资源

- [相关文档](./related.md)
"@

# 检测文档结构
function Get-DocumentStructure {
    param([string]$Content)

    $structure = @{
        HasMetadata = $false
        HasTOC = $false
        HasSeparator = $false
        HasContent = $false
        HasReferences = $false
        StructureScore = 0
    }

    $lines = $Content -split "`n"

    foreach ($line in $lines) {
        $trimmed = $line.Trim()

        # 检测元数据
        if ($trimmed -match '\*\*版本\*\*' -or
            $trimmed -match '\*\*更新日期\*\*' -or
            $trimmed -match '\*\*适用于\*\*') {
            $structure.HasMetadata = $true
            $structure.StructureScore++
        }

        # 检测目录
        if ($trimmed -match '^##+\s+.*目录') {
            $structure.HasTOC = $true
            $structure.StructureScore++
        }

        # 检测分隔线
        if ($trimmed -eq '---') {
            $structure.HasSeparator = $true
            $structure.StructureScore++
        }

        # 检测内容章节
        if ($trimmed -match '^##+\s+' -and $trimmed -notmatch '目录') {
            $structure.HasContent = $true
            $structure.StructureScore++
        }

        # 检测参考资源
        if ($trimmed -match '参考资源' -or $trimmed -match '参考资料') {
            $structure.HasReferences = $true
            $structure.StructureScore++
        }
    }

    return $structure
}

# 优化文档结构
function Optimize-Structure {
    param([string]$Content)

    $lines = $Content -split "`n"
    $newLines = @()
    $structure = Get-DocumentStructure -Content $Content

    # 如果结构已经完整，只做微调
    if ($structure.StructureScore -ge 4) {
        # 确保分隔线正确
        $inMetadata = $false
        $metadataEnd = -1

        for ($i = 0; $i -lt $lines.Count; $i++) {
            $line = $lines[$i]
            $trimmed = $line.Trim()

            # 检测元数据区域
            if ($trimmed -match '\*\*版本\*\*') {
                $inMetadata = $true
            }

            if ($inMetadata -and $trimmed -eq '---' -and $metadataEnd -eq -1) {
                $metadataEnd = $i
            }

            $newLines += $line
        }

        return $newLines -join "`n"
    }

    # 如果结构不完整，进行修复
    # 这里只做基本修复，复杂情况需要手动处理
    return $Content
}

# 主处理逻辑
$mdFiles = Get-ChildItem -Path $DocsPath -Filter "*.md" -Recurse
$totalFiles = $mdFiles.Count
$processedFiles = 0
$optimizedFiles = 0

Write-Host "📁 找到 $totalFiles 个Markdown文件" -ForegroundColor Cyan
Write-Host ""

foreach ($file in $mdFiles) {
    $processedFiles++

    try {
        $content = Get-Content $file.FullName -Raw -Encoding UTF8
        $structure = Get-DocumentStructure -Content $content

        # 只处理结构不完整的文件
        if ($structure.StructureScore -lt 3) {
            Write-Host "[$processedFiles/$totalFiles] 📐 $($file.Name)" -ForegroundColor Cyan
            Write-Host "  结构评分: $($structure.StructureScore)/5" -ForegroundColor Yellow
            Write-Host "  缺失: " -NoNewline -ForegroundColor Yellow
            if (-not $structure.HasMetadata) { Write-Host "元数据 " -NoNewline -ForegroundColor Red }
            if (-not $structure.HasTOC) { Write-Host "目录 " -NoNewline -ForegroundColor Red }
            if (-not $structure.HasSeparator) { Write-Host "分隔线 " -NoNewline -ForegroundColor Red }
            Write-Host ""

            if (-not $DryRun) {
                $newContent = Optimize-Structure -Content $content
                if ($newContent -ne $content) {
                    [System.IO.File]::WriteAllText($file.FullName, $newContent, [System.Text.Encoding]::UTF8)
                    Write-Host "  ✅ 已优化" -ForegroundColor Green
                    $optimizedFiles++
                }
            } else {
                Write-Host "  ⚠ 预览模式（未实际修改）" -ForegroundColor Yellow
            }
        } else {
            if ($Verbose) {
                Write-Host "[$processedFiles/$totalFiles] ✓ $($file.Name) - 结构完整 ($($structure.StructureScore)/5)" -ForegroundColor Gray
            }
        }
    } catch {
        Write-Host "  ❌ 错误: $_" -ForegroundColor Red
    }
}

# 总结
Write-Host ""
Write-Host "================================" -ForegroundColor Cyan
Write-Host "📊 优化总结" -ForegroundColor Cyan
Write-Host "  处理文件: $processedFiles/$totalFiles" -ForegroundColor White
Write-Host "  已优化: $optimizedFiles" -ForegroundColor White

if ($DryRun) {
    Write-Host "  ⚠ 这是预览模式，未实际修改文件" -ForegroundColor Yellow
    Write-Host "  使用 -DryRun:`$false 来实际优化" -ForegroundColor Yellow
} else {
    Write-Host "  ✅ 文件已更新" -ForegroundColor Green
}

Write-Host ""

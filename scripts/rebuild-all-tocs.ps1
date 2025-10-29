# 重建所有文档的目录
# 基于文档中的实际标题生成正确的 TOC

param(
    [string]$Path = "docs",
    [switch]$DryRun,
    [switch]$Verbose
)

$ErrorActionPreference = "Stop"
$stats = @{
    FilesProcessed = 0
    TocsRebuilt = 0
    Errors = 0
}

Write-Host "📝 重建所有文档目录...`n" -ForegroundColor Cyan

# GitHub anchor 生成规则（精确版本）
function Get-GitHubAnchor {
    param([string]$Heading)
    
    if ([string]::IsNullOrWhiteSpace($Heading)) {
        return ""
    }
    
    # 1. 移除markdown格式符号
    $anchor = $Heading -replace '\*\*', '' -replace '`', '' -replace '\[', '' -replace '\]', ''
    
    # 2. 移除常见emoji（GitHub会移除它们）
    $emojis = @('📋', '🎯', '✅', '❓', '🎉', '📚', '📝', '🔍', '🚨', '🆕', '🔗', '📊', 
                '💻', '🔧', '⚠️', '📖', '🎊', '📑', '🏆', '✨', '⭐', '🔥', '💡', 
                '📈', '📉', '🛠️', '🚀', '💪', '🌟', '⚡', '🎨', '🔄', '⚙️', '📦')
    foreach ($emoji in $emojis) {
        $anchor = $anchor -replace [regex]::Escape($emoji), ''
    }
    
    # 3. 移除其他特殊字符（保留中文、英文、数字、空格、连字符、点号）
    $anchor = $anchor -replace '[^a-z0-9\s\-.\u4e00-\u9fa5]', ''
    
    # 4. 转小写
    $anchor = $anchor.ToLower()
    
    # 5. trim首尾空格
    $anchor = $anchor.Trim()
    
    # 6. 空格替换为连字符
    $anchor = $anchor -replace '\s+', '-'
    
    # 7. 多个连字符替换为单个
    $anchor = $anchor -replace '-+', '-'
    
    # 8. 移除首尾连字符
    $anchor = $anchor.Trim('-')
    
    return $anchor
}

# 解析文档标题结构
function Get-DocumentHeadings {
    param([string]$Content)
    
    $headings = @()
    $lines = $Content -split "`r?`n"
    
    for ($i = 0; $i -lt $lines.Count; $i++) {
        $line = $lines[$i]
        if ($line -match '^(#{2,6})\s+(.+)$') {
            $level = $matches[1].Length
            $text = $matches[2].Trim()
            $anchor = Get-GitHubAnchor $text
            
            $headings += [PSCustomObject]@{
                Level = $level
                Text = $text
                Anchor = $anchor
                LineNumber = $i + 1
            }
        }
    }
    
    return $headings
}

# 生成 TOC
function Build-TOC {
    param([array]$Headings)
    
    $toc = @()
    $toc += "## 📋 目录"
    $toc += ""
    
    foreach ($heading in $Headings) {
        # 跳过"📋 目录"本身
        if ($heading.Text -match '📋\s*目录' -or $heading.Text -eq '目录') {
            continue
        }
        
        # 计算缩进级别（## = 0, ### = 1, #### = 2, etc.）
        $indent = $heading.Level - 2
        if ($indent < 0) { $indent = 0 }
        
        # 生成缩进
        $prefix = "  " * $indent
        
        # 生成链接
        $link = "- [$($heading.Text)](#$($heading.Anchor))"
        $toc += "$prefix$link"
    }
    
    $toc += ""
    return $toc -join "`n"
}

# 处理文件
$files = Get-ChildItem -Path $Path -Recurse -Filter "*.md" -File

foreach ($file in $files) {
    try {
        $stats.FilesProcessed++
        $content = Get-Content $file.FullName -Raw -Encoding UTF8
        
        if ([string]::IsNullOrWhiteSpace($content)) {
            continue
        }
        
        # 跳过某些文件类型
        if ($file.Name -match '^(🎉|🔍|📊|📝|🎯|🚨|🔧|README|CHANGELOG|LICENSE)') {
            if ($Verbose) {
                Write-Host "⏭️  跳过: $($file.Name)" -ForegroundColor Gray
            }
            continue
        }
        
        # 检查是否有目录
        if ($content -notmatch '##\s+📋\s*目录|##\s+目录') {
            continue
        }
        
        # 解析标题
        $headings = Get-DocumentHeadings $content
        
        if ($headings.Count -eq 0) {
            continue
        }
        
        # 生成新 TOC
        $newToc = Build-TOC $headings
        
        # 替换旧 TOC
        # 查找 TOC 的开始和结束
        if ($content -match '(?s)(##\s+📋\s*目录|##\s+目录).*?\n\n(##\s+[^#])') {
            $tocStart = $matches[0].IndexOf($matches[1])
            $tocEnd = $matches[0].IndexOf($matches[2])
            
            # 提取 TOC 前后的内容
            $beforeToc = $content.Substring(0, $tocStart)
            $afterToc = $content.Substring($tocStart)
            
            # 找到下一个 ## 标题的位置
            if ($afterToc -match '\n\n##\s+[^#]') {
                $nextHeadingPos = $afterToc.IndexOf($matches[0]) + 2  # +2 for the \n\n
                $beforeNextHeading = $afterToc.Substring(0, $nextHeadingPos)
                $afterNextHeading = $afterToc.Substring($nextHeadingPos)
                
                # 构建新内容
                $newContent = $beforeToc + $newToc + "`n" + $afterNextHeading
                
                if ($newContent -ne $content) {
                    Write-Host "✓ $($file.FullName -replace [regex]::Escape($PWD), '.')" -ForegroundColor Green
                    Write-Host "  重建 TOC: $($headings.Count) 个标题" -ForegroundColor Gray
                    
                    $stats.TocsRebuilt++
                    
                    if (-not $DryRun) {
                        Set-Content -Path $file.FullName -Value $newContent -Encoding UTF8 -NoNewline
                    }
                }
            }
        }
        
    } catch {
        $stats.Errors++
        Write-Host "✗ $($file.Name): $($_.Exception.Message)" -ForegroundColor Red
    }
}

# 显示统计
Write-Host "`n" + ("="*70) -ForegroundColor Cyan
Write-Host "📊 重建统计:" -ForegroundColor Cyan
Write-Host "  处理文件: $($stats.FilesProcessed)" -ForegroundColor White
Write-Host "  重建TOC:  $($stats.TocsRebuilt)" -ForegroundColor Green
Write-Host "  错误数量: $($stats.Errors)" -ForegroundColor $(if($stats.Errors -eq 0){'Green'}else{'Red'})
Write-Host ("="*70) -ForegroundColor Cyan

if ($DryRun) {
    Write-Host "`n⚠️  试运行模式 - 未实际修改文件" -ForegroundColor Yellow
} else {
    Write-Host "`n✅ 重建完成！" -ForegroundColor Green
}


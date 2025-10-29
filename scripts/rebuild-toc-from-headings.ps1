# 从实际标题重建TOC
# 读取文档中的所有标题，生成正确的TOC

param(
    [string]$FilePath,
    [switch]$DryRun
)

$ErrorActionPreference = "Stop"

if (-not $FilePath) {
    Write-Host "用法: .\rebuild-toc-from-headings.ps1 -FilePath <文件路径>" -ForegroundColor Yellow
    exit 1
}

Write-Host "🔧 重建TOC: $FilePath`n" -ForegroundColor Cyan

# GitHub anchor 生成规则（精确版）
function Get-GitHubAnchor {
    param([string]$Heading)
    
    if ([string]::IsNullOrWhiteSpace($Heading)) {
        return ""
    }
    
    # 1. 移除markdown格式符号
    $anchor = $Heading -replace '\*\*', '' -replace '`', '' -replace '\[', '' -replace '\]', ''
    
    # 2. 移除所有emoji 
    # 列出文档中可能出现的所有emoji
    $emojis = '📋🎯✅❓🎉📚📝🔍🚨🆕🔗📊💻🔧⚠️📖🎊📑🏆✨⭐🔥💡📈📉🛠️🚀💪🌟⚡🎨🔄⚙️📦🌐🏗️☁️🔬🎖️📁🗂️📂📄📃📜📋🗃️🧩🎪🎭🎬🎮🎯🎲🎰🎳🏀🏈⚾🥎🏐🏉🎾🏸🏑🏒🏓🏏🥍🏹🥊🥋🥅⛳🚩🎌'
    $emojiArray = $emojis.ToCharArray()
    foreach ($emoji in $emojiArray) {
        $anchor = $anchor -replace [regex]::Escape($emoji.ToString()), ''
    }
    
    # 3. 移除其他特殊字符
    # GitHub规则：保留中文、英文、数字、连字符
    # 移除: 冒号、括号(中英文)、斜杠、点号等
    $anchor = $anchor -replace '：', '' -replace ':', ''
    $anchor = $anchor -replace '（', '' -replace '）', ''
    $anchor = $anchor -replace '\(', '' -replace '\)', ''
    $anchor = $anchor -replace '/', '' -replace '\\', ''
    $anchor = $anchor -replace '\.', ''
    $anchor = $anchor -replace ',', '' -replace '，', ''
    $anchor = $anchor -replace '、', ''
    $anchor = $anchor -replace '？', '' -replace '\?', ''
    $anchor = $anchor -replace '！', '' -replace '!', ''
    $anchor = $anchor -replace '"', '' -replace '"', '' -replace '"', ''
    $anchor = $anchor -replace ''', '' -replace ''', '' -replace "'", ''
    $anchor = $anchor -replace '\+', ''
    
    # 4. trim 空格
    $anchor = $anchor.Trim()
    
    # 5. 转小写（对英文生效）
    $anchor = $anchor.ToLower()
    
    # 6. 空格替换为连字符
    $anchor = $anchor -replace '\s+', '-'
    
    # 7. 多个连字符替换为单个
    $anchor = $anchor -replace '-+', '-'
    
    # 8. 移除首尾连字符
    $anchor = $anchor.Trim('-')
    
    return $anchor
}

# 读取文件
$content = Get-Content $FilePath -Raw -Encoding UTF8

# 提取所有标题（##到###）
$headings = @()
$lines = $content -split "`r?`n"

for ($i = 0; $i -lt $lines.Count; $i++) {
    $line = $lines[$i]
    
    # 匹配 ## 或 ### 标题（不包括 # 单个井号）
    if ($line -match '^(#{2,3})\s+(.+)$') {
        $level = $matches[1].Length
        $text = $matches[2].Trim()
        
        # 跳过 "📋 目录" 标题本身
        if ($text -match '📋\s*目录' -or $text -eq '目录') {
            continue
        }
        
        $anchor = Get-GitHubAnchor $text
        
        $headings += [PSCustomObject]@{
            Level = $level
            Text = $text
            Anchor = $anchor
            LineNumber = $i + 1
        }
    }
}

Write-Host "找到 $($headings.Count) 个标题`n" -ForegroundColor Green

# 生成新的TOC
$newToc = @()
$newToc += "## 📋 目录"
$newToc += ""

foreach ($heading in $headings) {
    # 计算缩进（## = 0, ### = 1）
    $indent = $heading.Level - 2
    $prefix = "  " * $indent
    
    # 生成链接
    $link = "- [$($heading.Text)](#$($heading.Anchor))"
    $newToc += "$prefix$link"
}

$newToc += ""

Write-Host "生成的TOC:" -ForegroundColor Cyan
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray
$newToc | Select-Object -First 30 | ForEach-Object { Write-Host $_ -ForegroundColor White }
if ($newToc.Count -gt 30) {
    Write-Host "... (还有 $($newToc.Count - 30) 行)" -ForegroundColor Gray
}
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray

if (-not $DryRun) {
    # 替换原文件中的TOC
    # 找到TOC的开始和结束位置
    $tocStart = $content.IndexOf("## 📋 目录")
    if ($tocStart -eq -1) {
        $tocStart = $content.IndexOf("## 目录")
    }
    
    if ($tocStart -ge 0) {
        # 跳过TOC标题行，从下一行开始找
        $afterTocStart = $tocStart + 10  # "## 📋 目录"的长度
        $afterToc = $content.Substring($afterTocStart)
        
        # 找到下一个 ## (不是###) 的位置
        if ($afterToc -match '(?m)^##\s+[^#]') {
            $nextHeadingRelPos = $afterToc.IndexOf($matches[0])
            $tocEnd = $afterTocStart + $nextHeadingRelPos
            
            # 构建新内容
            $beforeToc = $content.Substring(0, $tocStart)
            $afterTocContent = $content.Substring($tocEnd)
            
            $newContent = $beforeToc + ($newToc -join "`n") + "`n`n" + $afterTocContent
            
            # 保存
            Set-Content -Path $FilePath -Value $newContent -Encoding UTF8 -NoNewline
            
            Write-Host "`n✅ TOC已更新！" -ForegroundColor Green
        } else {
            Write-Host "`n⚠️  未找到TOC结束位置" -ForegroundColor Yellow
        }
    } else {
        Write-Host "`n⚠️  未找到TOC" -ForegroundColor Yellow
    }
} else {
    Write-Host "`n⚠️  试运行模式 - 未修改文件" -ForegroundColor Yellow
}


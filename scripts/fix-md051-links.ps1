# 修复 MD051 链接错误
# 自动修复文档中所有无效的内部链接

param(
    [string]$Path = "docs",
    [switch]$DryRun,
    [switch]$Verbose
)

$ErrorActionPreference = "Stop"
$stats = @{
    FilesProcessed = 0
    LinksFixed = 0
    Errors = 0
}

Write-Host "🔗 修复 MD051 链接错误...`n" -ForegroundColor Cyan

# GitHub anchor 生成规则
function Get-GitHubAnchor {
    param([string]$Heading)
    
    if ([string]::IsNullOrWhiteSpace($Heading)) {
        return ""
    }
    
    # 1. 移除markdown格式符号
    $anchor = $Heading -replace '\*\*', '' -replace '`', '' -replace '\[', '' -replace '\]', ''
    
    # 2. 移除常见emoji（使用字符串操作）
    $emojis = @('📋', '🎯', '✅', '❓', '🎉', '📚', '📝', '🔍', '🚨', '🆕', '🔗', '📊', 
                '💻', '🔧', '⚠️', '📖', '🎊', '📑', '🏆', '✨', '⭐', '🔥', '💡', 
                '📈', '📉', '🛠️', '🚀', '💪', '🌟', '⚡', '🎨')
    foreach ($emoji in $emojis) {
        $anchor = $anchor.Replace($emoji, '')
    }
    
    # 3. 移除其他特殊符号（保留中文、英文、数字、空格、连字符）
    # 只保留：字母、数字、中文、空格、连字符
    $result = ""
    foreach ($char in $anchor.ToCharArray()) {
        if ($char -match '[a-zA-Z0-9\s\-\u4e00-\u9fa5]') {
            $result += $char
        }
    }
    $anchor = $result
    
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

# 提取文档中所有标题及其anchor
function Get-DocumentHeadings {
    param([string]$Content)
    
    $headings = @{}
    $lines = $Content -split "`n"
    
    foreach ($line in $lines) {
        # 匹配标题行: ## Title 或 ### Title
        if ($line -match '^(#{1,6})\s+(.+)$') {
            $headingText = $Matches[2].Trim()
            $anchor = Get-GitHubAnchor $headingText
            
            if ($anchor) {
                # 如果anchor已存在，GitHub会添加数字后缀
                $originalAnchor = $anchor
                $counter = 1
                while ($headings.ContainsKey($anchor)) {
                    $anchor = "$originalAnchor-$counter"
                    $counter++
                }
                
                $headings[$anchor] = $headingText
            }
        }
    }
    
    return $headings
}

# 修复文档中的链接
function Fix-DocumentLinks {
    param(
        [string]$Content,
        [hashtable]$Headings
    )
    
    $fixedContent = $Content
    $fixCount = 0
    
    # 匹配所有内部链接: [text](#anchor)
    $linkPattern = '\[([^\]]+)\]\(#([^\)]+)\)'
    $matches = [regex]::Matches($Content, $linkPattern)
    
    foreach ($match in $matches) {
        $linkText = $match.Groups[1].Value
        $currentAnchor = $match.Groups[2].Value.Trim()
        
        # 检查这个anchor是否存在
        if (-not $Headings.ContainsKey($currentAnchor)) {
            # 尝试从链接文本生成正确的anchor
            $expectedAnchor = Get-GitHubAnchor $linkText
            
            # 如果生成的anchor存在于文档中
            if ($expectedAnchor -and $Headings.ContainsKey($expectedAnchor)) {
                $oldLink = $match.Value
                $newLink = "[$linkText](#$expectedAnchor)"
                $fixedContent = $fixedContent.Replace($oldLink, $newLink)
                $fixCount++
                
                if ($Verbose) {
                    Write-Host "  修复: $oldLink -> $newLink" -ForegroundColor Gray
                }
            }
            # 否则尝试模糊匹配
            else {
                # 尝试找到最相似的anchor
                $bestMatch = $null
                $bestScore = 0
                
                foreach ($anchor in $Headings.Keys) {
                    # 简单的相似度计算：检查包含关系
                    if ($currentAnchor -like "*$anchor*" -or $anchor -like "*$currentAnchor*") {
                        $score = [Math]::Min($currentAnchor.Length, $anchor.Length)
                        if ($score -gt $bestScore) {
                            $bestScore = $score
                            $bestMatch = $anchor
                        }
                    }
                }
                
                if ($bestMatch) {
                    $oldLink = $match.Value
                    $newLink = "[$linkText](#$bestMatch)"
                    $fixedContent = $fixedContent.Replace($oldLink, $newLink)
                    $fixCount++
                    
                    if ($Verbose) {
                        Write-Host "  修复(模糊): $oldLink -> $newLink" -ForegroundColor Yellow
                    }
                }
            }
        }
    }
    
    return @{
        Content = $fixedContent
        FixCount = $fixCount
    }
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
        
        # 获取文档中所有标题
        $headings = Get-DocumentHeadings $content
        
        if ($headings.Count -eq 0) {
            continue
        }
        
        # 修复链接
        $result = Fix-DocumentLinks $content $headings
        
        if ($result.FixCount -gt 0) {
            $stats.LinksFixed += $result.FixCount
            
            Write-Host "✓ $($file.FullName -replace [regex]::Escape($PWD), '.')" -ForegroundColor Green
            Write-Host "  修复链接: $($result.FixCount)个" -ForegroundColor Gray
            
            if (-not $DryRun) {
                Set-Content -Path $file.FullName -Value $result.Content -Encoding UTF8 -NoNewline
            }
        }
        
    } catch {
        $stats.Errors++
        Write-Host "✗ $($file.Name): $($_.Exception.Message)" -ForegroundColor Red
    }
}

# 显示统计
Write-Host "`n" + ("="*70) -ForegroundColor Cyan
Write-Host "📊 修复统计:" -ForegroundColor Cyan
Write-Host "  处理文件: $($stats.FilesProcessed)" -ForegroundColor White
Write-Host "  修复链接: $($stats.LinksFixed)" -ForegroundColor Green
Write-Host "  错误数量: $($stats.Errors)" -ForegroundColor $(if($stats.Errors -eq 0){'Green'}else{'Red'})
Write-Host ("="*70) -ForegroundColor Cyan

if ($DryRun) {
    Write-Host "`n⚠️  试运行模式 - 未实际修改文件" -ForegroundColor Yellow
}


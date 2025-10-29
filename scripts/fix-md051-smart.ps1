# 智能修复 MD051 链接错误
# 使用实际标题匹配和模糊搜索

param(
    [string]$Path = "docs",
    [switch]$DryRun,
    [switch]$Verbose
)

$ErrorActionPreference = "Stop"
$stats = @{
    FilesProcessed = 0
    LinksFixed = 0
    LinksSkipped = 0
    Errors = 0
}

Write-Host "🔗 智能修复 MD051 链接错误...`n" -ForegroundColor Cyan

# GitHub anchor 生成规则
function Get-GitHubAnchor {
    param([string]$Heading)
    
    if ([string]::IsNullOrWhiteSpace($Heading)) {
        return ""
    }
    
    # 1. 移除markdown格式符号
    $anchor = $Heading -replace '\*\*', '' -replace '`', '' -replace '\[', '' -replace '\]', ''
    
    # 2. 移除emoji
    $emojis = @('📋', '🎯', '✅', '❓', '🎉', '📚', '📝', '🔍', '🚨', '🆕', '🔗', '📊', 
                '💻', '🔧', '⚠️', '📖', '🎊', '📑', '🏆', '✨', '⭐', '🔥', '💡', 
                '📈', '📉', '🛠️', '🚀', '💪', '🌟', '⚡', '🎨', '🔄', '⚙️', '📦')
    foreach ($emoji in $emojis) {
        $anchor = $anchor.Replace($emoji, '')
    }
    
    # 3. 只保留：字母、数字、中文、空格、连字符
    $result = ""
    foreach ($char in $anchor.ToCharArray()) {
        if ($char -match '[a-zA-Z0-9\s\-\u4e00-\u9fa5\.]') {
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

# 提取文档中所有标题
function Get-DocumentHeadings {
    param([string]$Content)
    
    $headings = @{}
    $headingsList = @()
    $lines = $Content -split "`n"
    
    foreach ($line in $lines) {
        # 匹配标题行: ## Title 或 ### Title
        if ($line -match '^(#{1,6})\s+(.+)$') {
            $level = $Matches[1].Length
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
                $headingsList += @{
                    Anchor = $anchor
                    Text = $headingText
                    Level = $level
                }
            }
        }
    }
    
    return @{
        Map = $headings
        List = $headingsList
    }
}

# 计算两个字符串的相似度（简单版本）
function Get-StringSimilarity {
    param([string]$str1, [string]$str2)
    
    $str1 = $str1.ToLower()
    $str2 = $str2.ToLower()
    
    # 完全匹配
    if ($str1 -eq $str2) {
        return 100
    }
    
    # 包含关系
    if ($str1.Contains($str2) -or $str2.Contains($str1)) {
        $minLen = [Math]::Min($str1.Length, $str2.Length)
        $maxLen = [Math]::Max($str1.Length, $str2.Length)
        return [int](($minLen / $maxLen) * 80)
    }
    
    # 开头匹配
    $commonPrefix = 0
    for ($i = 0; $i -lt [Math]::Min($str1.Length, $str2.Length); $i++) {
        if ($str1[$i] -eq $str2[$i]) {
            $commonPrefix++
        } else {
            break
        }
    }
    
    if ($commonPrefix -gt 3) {
        return [int](($commonPrefix / [Math]::Max($str1.Length, $str2.Length)) * 60)
    }
    
    return 0
}

# 查找最佳匹配的标题
function Find-BestMatchingHeading {
    param(
        [string]$LinkText,
        [string]$CurrentAnchor,
        [array]$Headings
    )
    
    $bestMatch = $null
    $bestScore = 0
    
    foreach ($heading in $Headings) {
        # 尝试多种匹配策略
        $scores = @()
        
        # 策略1: anchor直接匹配
        $expectedAnchor = Get-GitHubAnchor $LinkText
        if ($expectedAnchor -eq $heading.Anchor) {
            return $heading.Anchor
        }
        
        # 策略2: 链接文本与标题文本相似度
        $textScore = Get-StringSimilarity $LinkText $heading.Text
        $scores += $textScore
        
        # 策略3: 当前anchor与目标anchor相似度
        $anchorScore = Get-StringSimilarity $CurrentAnchor $heading.Anchor
        $scores += $anchorScore
        
        # 策略4: 关键词匹配
        $linkWords = $LinkText -split '[\s\-]' | Where-Object { $_.Length -gt 2 }
        $headingWords = $heading.Text -split '[\s\-]' | Where-Object { $_.Length -gt 2 }
        $commonWords = ($linkWords | Where-Object { $headingWords -contains $_ }).Count
        if ($linkWords.Count -gt 0) {
            $keywordScore = [int](($commonWords / $linkWords.Count) * 70)
            $scores += $keywordScore
        }
        
        # 综合得分
        $totalScore = ($scores | Measure-Object -Average).Average
        
        if ($totalScore -gt $bestScore) {
            $bestScore = $totalScore
            $bestMatch = $heading
        }
    }
    
    # 只有当相似度足够高时才返回匹配
    if ($bestScore -ge 40) {
        return $bestMatch.Anchor
    }
    
    return $null
}

# 修复文档中的链接
function Fix-DocumentLinks {
    param(
        [string]$Content,
        [hashtable]$HeadingsMap,
        [array]$HeadingsList,
        [string]$FilePath
    )
    
    $fixedContent = $Content
    $fixCount = 0
    $skipCount = 0
    
    # 匹配所有内部链接: [text](#anchor)
    $linkPattern = '\[([^\]]+)\]\(#([^\)]+)\)'
    $matches = [regex]::Matches($Content, $linkPattern)
    
    foreach ($match in $matches) {
        $linkText = $match.Groups[1].Value
        $currentAnchor = $match.Groups[2].Value.Trim()
        
        # 检查这个anchor是否有效
        if (-not $HeadingsMap.ContainsKey($currentAnchor)) {
            # 尝试找到最佳匹配
            $bestAnchor = Find-BestMatchingHeading $linkText $currentAnchor $HeadingsList
            
            if ($bestAnchor) {
                $oldLink = $match.Value
                $newLink = "[$linkText](#$bestAnchor)"
                
                # 只替换当前这一个匹配项
                $index = $fixedContent.IndexOf($oldLink)
                if ($index -ge 0) {
                    $fixedContent = $fixedContent.Substring(0, $index) + $newLink + $fixedContent.Substring($index + $oldLink.Length)
                    $fixCount++
                    
                    if ($Verbose) {
                        Write-Host "  ✓ 修复: [$linkText](#$currentAnchor) -> [$linkText](#$bestAnchor)" -ForegroundColor Green
                        Write-Host "    标题: $($HeadingsMap[$bestAnchor])" -ForegroundColor Gray
                    }
                }
            } else {
                $skipCount++
                if ($Verbose) {
                    Write-Host "  ⚠ 跳过: [$linkText](#$currentAnchor) - 未找到匹配标题" -ForegroundColor Yellow
                }
            }
        }
    }
    
    return @{
        Content = $fixedContent
        FixCount = $fixCount
        SkipCount = $skipCount
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
        $headingData = Get-DocumentHeadings $content
        
        if ($headingData.Map.Count -eq 0) {
            continue
        }
        
        # 修复链接
        $result = Fix-DocumentLinks $content $headingData.Map $headingData.List $file.FullName
        
        if ($result.FixCount -gt 0 -or $result.SkipCount -gt 0) {
            Write-Host "✓ $($file.FullName -replace [regex]::Escape($PWD), '.')" -ForegroundColor Cyan
            if ($result.FixCount -gt 0) {
                Write-Host "  修复: $($result.FixCount)个" -ForegroundColor Green
            }
            if ($result.SkipCount -gt 0) {
                Write-Host "  跳过: $($result.SkipCount)个" -ForegroundColor Yellow
            }
            
            $stats.LinksFixed += $result.FixCount
            $stats.LinksSkipped += $result.SkipCount
            
            if (-not $DryRun -and $result.FixCount -gt 0) {
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
Write-Host "  跳过链接: $($stats.LinksSkipped)" -ForegroundColor Yellow
Write-Host "  错误数量: $($stats.Errors)" -ForegroundColor $(if($stats.Errors -eq 0){'Green'}else{'Red'})
Write-Host ("="*70) -ForegroundColor Cyan

if ($DryRun) {
    Write-Host "`n⚠️  试运行模式 - 未实际修改文件" -ForegroundColor Yellow
} else {
    Write-Host "`n✅ 修复完成！" -ForegroundColor Green
}


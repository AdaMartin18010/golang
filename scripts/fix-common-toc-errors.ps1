# 修复常见的 TOC 错误模式
# 处理链接文本与anchor不匹配的情况

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

Write-Host "🔧 修复常见 TOC 错误...`n" -ForegroundColor Cyan

# GitHub anchor 生成规则
function Get-GitHubAnchor {
    param([string]$Text)
    
    if ([string]::IsNullOrWhiteSpace($Text)) {
        return ""
    }
    
    # 移除markdown格式
    $anchor = $Text -replace '\*\*', '' -replace '`', '' -replace '\[', '' -replace '\]', ''
    
    # 移除emoji
    $emojis = @('📋', '🎯', '✅', '❓', '🎉', '📚', '📝', '🔍', '🚨', '🆕', '🔗', '📊', 
                '💻', '🔧', '⚠️', '📖', '🎊', '📑', '🏆', '✨', '⭐', '🔥', '💡', 
                '📈', '📉', '🛠️', '🚀', '💪', '🌟', '⚡', '🎨', '🔄', '⚙️', '📦')
    foreach ($emoji in $emojis) {
        $anchor = $anchor -replace [regex]::Escape($emoji), ''
    }
    
    # 移除特殊字符（保留中文、英文、数字、空格、连字符、点号）
    $anchor = $anchor -replace '[^a-z0-9\s\-.\u4e00-\u9fa5]', ''
    
    # 转小写、trim、替换空格为连字符
    $anchor = $anchor.ToLower().Trim()
    $anchor = $anchor -replace '\s+', '-'
    $anchor = $anchor -replace '-+', '-'
    $anchor = $anchor.Trim('-')
    
    return $anchor
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
        
        $originalContent = $content
        $fixCount = 0
        
        # 提取所有链接: [text](#anchor)
        $linkPattern = '\[([^\]]+)\]\(#([^\)]+)\)'
        $links = [regex]::Matches($content, $linkPattern)
        
        foreach ($link in $links) {
            $linkText = $link.Groups[1].Value
            $currentAnchor = $link.Groups[2].Value
            
            # 生成基于链接文本的正确anchor
            $correctAnchor = Get-GitHubAnchor $linkText
            
            # 如果当前anchor与正确anchor不同，则替换
            if ($correctAnchor -and $currentAnchor -ne $correctAnchor) {
                $oldLink = $link.Value
                $newLink = "[$linkText](#$correctAnchor)"
                
                # 使用精确替换（只替换这一个）
                $startIndex = $link.Index
                $length = $link.Length
                $before = $content.Substring(0, $startIndex)
                $after = $content.Substring($startIndex + $length)
                $content = $before + $newLink + $after
                
                $fixCount++
                
                if ($Verbose) {
                    Write-Host "    $oldLink" -ForegroundColor Gray
                    Write-Host " -> $newLink" -ForegroundColor Green
                }
            }
        }
        
        if ($fixCount -gt 0) {
            Write-Host "✓ $($file.FullName -replace [regex]::Escape($PWD), '.')" -ForegroundColor Green
            Write-Host "  修复: $fixCount 个链接" -ForegroundColor Gray
            
            $stats.LinksFixed += $fixCount
            
            if (-not $DryRun) {
                Set-Content -Path $file.FullName -Value $content -Encoding UTF8 -NoNewline
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
} else {
    Write-Host "`n✅ 修复完成！" -ForegroundColor Green
}


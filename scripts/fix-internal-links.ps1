# 修复文档内部链接
# 根据GitHub的anchor规则自动生成正确的链接

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

Write-Host "🔗 开始修复内部链接..." -ForegroundColor Cyan
Write-Host "工作目录: $Path"
Write-Host "模式: $(if($DryRun){'试运行'}else{'实际修复'})`n"

function ConvertTo-GitHubAnchor {
    param([string]$Heading)
    
    # GitHub anchor生成规则：
    # 1. 转小写
    # 2. 移除emoji和特殊符号（保留字母、数字、中文、连字符、空格）
    # 3. 空格和多个连字符替换为单个连字符
    # 4. 移除首尾连字符
    
    $anchor = $Heading.ToLower()
    
    # 移除markdown格式符号
    $anchor = $anchor -replace '\*\*', ''
    $anchor = $anchor -replace '`', ''
    
    # 移除emoji和特殊符号（保留中文、字母、数字、空格、连字符）
    $anchor = $anchor -replace '[^\w\s\-\u4e00-\u9fa5]', ''
    
    # 空格替换为连字符
    $anchor = $anchor -replace '\s+', '-'
    
    # 多个连字符替换为单个
    $anchor = $anchor -replace '\-+', '-'
    
    # 移除首尾连字符
    $anchor = $anchor.Trim('-')
    
    return $anchor
}

function Fix-FileLinks {
    param($FilePath)
    
    try {
        $content = Get-Content $FilePath -Raw -Encoding UTF8
        $originalContent = $content
        
        # 提取所有标题（支持多级标题）
        $headings = [regex]::Matches($content, '(?m)^(#{1,6})\s+(.+)$')
        
        # 构建标题到anchor的映射
        $headingMap = @{}
        foreach ($heading in $headings) {
            $headingText = $heading.Groups[2].Value.Trim()
            $correctAnchor = ConvertTo-GitHubAnchor -Heading $headingText
            $headingMap[$headingText] = $correctAnchor
        }
        
        if ($headingMap.Count -eq 0) {
            return @{ Modified = $false; FixedCount = 0 }
        }
        
        $fixedCount = 0
        
        # 查找所有内部链接
        $linkPattern = '\[([^\]]+)\]\(#([^\)]+)\)'
        $matches = [regex]::Matches($content, $linkPattern)
        
        # 从后往前替换（避免位置偏移）
        for ($i = $matches.Count - 1; $i -ge 0; $i--) {
            $match = $matches[$i]
            $linkText = $match.Groups[1].Value
            $oldAnchor = $match.Groups[2].Value
            
            # 尝试从链接文本推断正确的标题
            # 移除链接文本中的编号、emoji等，尝试匹配标题
            $cleanLinkText = $linkText -replace '^\d+\.?\s*', '' # 移除前导数字
            $cleanLinkText = $cleanLinkText -replace '^[^\w\u4e00-\u9fa5]+', '' # 移除前导符号
            
            $newAnchor = $null
            
            # 方法1：精确匹配标题文本
            foreach ($heading in $headingMap.Keys) {
                if ($heading -eq $linkText -or $heading.Contains($cleanLinkText)) {
                    $newAnchor = $headingMap[$heading]
                    break
                }
            }
            
            # 方法2：如果没找到，尝试生成anchor看是否匹配
            if (-not $newAnchor) {
                $generatedAnchor = ConvertTo-GitHubAnchor -Heading $linkText
                if ($headingMap.Values -contains $generatedAnchor) {
                    $newAnchor = $generatedAnchor
                }
            }
            
            # 如果找到了更好的anchor并且与当前不同
            if ($newAnchor -and $newAnchor -ne $oldAnchor) {
                $oldLink = $match.Value
                $newLink = "[$linkText](#$newAnchor)"
                
                $startIndex = $match.Index
                $length = $match.Length
                
                $content = $content.Substring(0, $startIndex) + $newLink + $content.Substring($startIndex + $length)
                $fixedCount++
                
                if ($Verbose) {
                    Write-Host "    修复: [$linkText](#$oldAnchor) → [$linkText](#$newAnchor)" -ForegroundColor Gray
                }
            }
        }
        
        $modified = ($content -ne $originalContent)
        
        if ($modified -and -not $DryRun) {
            Set-Content -Path $FilePath -Value $content -NoNewline -Encoding UTF8
        }
        
        return @{ Modified = $modified; FixedCount = $fixedCount }
        
    } catch {
        Write-Host "  ✗ 错误: $($_.Exception.Message)" -ForegroundColor Red
        $stats.Errors++
        return @{ Modified = $false; FixedCount = 0 }
    }
}

# 主执行
$files = Get-ChildItem -Path $Path -Recurse -Filter "*.md" -File
$totalFiles = $files.Count

Write-Host "找到 $totalFiles 个Markdown文件`n" -ForegroundColor Green

$progress = 0
foreach ($file in $files) {
    $progress++
    $percentComplete = [math]::Round(($progress / $totalFiles) * 100)
    
    Write-Progress -Activity "修复链接" -Status "$progress/$totalFiles" -PercentComplete $percentComplete
    
    $stats.FilesProcessed++
    
    $result = Fix-FileLinks -FilePath $file.FullName
    
    if ($result.Modified) {
        $stats.LinksFixed += $result.FixedCount
        if ($Verbose) {
            Write-Host "✓ $($file.Name): 修复 $($result.FixedCount) 个链接" -ForegroundColor Green
        } else {
            Write-Host "✓ $($file.Name)" -ForegroundColor Green
        }
    }
}

Write-Progress -Activity "修复链接" -Completed

# 结果报告
Write-Host "`n" + ("="*60) -ForegroundColor Cyan
Write-Host "📊 链接修复统计" -ForegroundColor Cyan
Write-Host ("="*60) -ForegroundColor Cyan

Write-Host "`n文件处理:"
Write-Host "  ✓ 已处理: $($stats.FilesProcessed)" -ForegroundColor Green
Write-Host "  ⚠ 错误: $($stats.Errors)" -ForegroundColor $(if($stats.Errors -gt 0){'Red'}else{'Gray'})

Write-Host "`n修复详情:"
Write-Host "  🔗 链接修复: $($stats.LinksFixed) 个" -ForegroundColor Yellow

if ($DryRun) {
    Write-Host "`n⚠️  这是试运行，未实际修改文件" -ForegroundColor Yellow
} else {
    Write-Host "`n✅ 修复已完成！" -ForegroundColor Green
}

Write-Host "`n" + ("="*60) -ForegroundColor Cyan

return $stats


# 修复Markdown文档中的目录链接
# Markdown锚点规则：
# 1. 转为小写
# 2. 空格转为 -
# 3. 移除大多数特殊字符
# 4. 中文字符保留
# 5. emoji和特殊符号移除

$ErrorActionPreference = "Stop"
$filesFixed = 0
$totalIssues = 0

function ConvertTo-AnchorLink {
    param([string]$text)
    
    # 移除markdown格式符号
    $cleaned = $text -replace '\[([^\]]+)\]\([^\)]+\)', '$1'  # 移除链接
    $cleaned = $cleaned -replace '\*\*|\*|`|~~', ''  # 移除粗体、斜体、代码、删除线
    
    # 移除emoji和特殊符号（保留中文、英文、数字、空格、-）
    $cleaned = $cleaned -replace '[^\p{L}\p{N}\s\-\.+]', ''
    
    # 转为小写（仅英文）
    $cleaned = $cleaned.ToLower()
    
    # 空格和点转为连字符
    $cleaned = $cleaned -replace '[\s\.]+', '-'
    
    # 移除开头和结尾的连字符
    $cleaned = $cleaned.Trim('-')
    
    # 多个连续连字符合并为一个
    $cleaned = $cleaned -replace '-+', '-'
    
    return $cleaned
}

function Get-Headings {
    param([string[]]$lines)
    
    $headings = @{}
    $inCodeBlock = $false
    
    for ($i = 0; $i -lt $lines.Count; $i++) {
        $line = $lines[$i]
        
        # 检查是否在代码块中
        if ($line -match '^```') {
            $inCodeBlock = -not $inCodeBlock
            continue
        }
        
        if ($inCodeBlock) {
            continue
        }
        
        # 匹配标题
        if ($line -match '^(#{1,6})\s+(.+)$') {
            $level = $matches[1].Length
            $title = $matches[2].Trim()
            $anchor = ConvertTo-AnchorLink $title
            
            if (-not $headings.ContainsKey($anchor)) {
                $headings[$anchor] = @{
                    Title = $title
                    Level = $level
                    Line = $i
                }
            }
        }
    }
    
    return $headings
}

function Get-TOCLinks {
    param([string[]]$lines)
    
    $tocLinks = @()
    $inTOC = $false
    $inCodeBlock = $false
    
    for ($i = 0; $i -lt $lines.Count; $i++) {
        $line = $lines[$i]
        
        # 检查代码块
        if ($line -match '^```') {
            $inCodeBlock = -not $inCodeBlock
            continue
        }
        
        if ($inCodeBlock) {
            continue
        }
        
        # 检测目录开始
        if ($line -match '^##\s+(📋\s*)?目录$' -or $line -match '^##\s+Table of Contents$') {
            $inTOC = $true
            continue
        }
        
        # 检测目录结束（遇到下一个二级标题或分隔符）
        if ($inTOC -and ($line -match '^##\s+' -or $line -eq '---')) {
            $inTOC = $false
        }
        
        # 提取目录中的链接
        if ($inTOC -and $line -match '^\s*-\s+\[([^\]]+)\]\(#([^\)]+)\)') {
            $tocLinks += @{
                Text = $matches[1]
                Link = $matches[2]
                Line = $i
                OriginalLine = $line
            }
        }
    }
    
    return $tocLinks
}

Write-Host "🔍 扫描所有Markdown文件..." -ForegroundColor Cyan

$mdFiles = Get-ChildItem -Path "docs" -Filter "*.md" -Recurse -File | Where-Object {
    $_.FullName -notlike "*\node_modules\*"
}

Write-Host "📝 找到 $($mdFiles.Count) 个文件`n" -ForegroundColor Green

foreach ($file in $mdFiles) {
    try {
        $content = Get-Content -Path $file.FullName -Raw -Encoding UTF8
        if (-not $content) {
            continue
        }
        
        $lines = $content -split "`r?`n"
        
        # 获取所有标题及其锚点
        $headings = Get-Headings $lines
        
        # 获取目录中的所有链接
        $tocLinks = Get-TOCLinks $lines
        
        if ($tocLinks.Count -eq 0) {
            continue
        }
        
        $fileIssues = 0
        $changed = $false
        $newLines = $lines.Clone()
        
        foreach ($tocLink in $tocLinks) {
            $currentLink = $tocLink.Link
            
            # 检查链接是否存在对应的标题
            if (-not $headings.ContainsKey($currentLink)) {
                # 尝试从文本生成正确的锚点
                $expectedAnchor = ConvertTo-AnchorLink $tocLink.Text
                
                if ($headings.ContainsKey($expectedAnchor)) {
                    # 找到匹配的标题，修复链接
                    $oldLine = $tocLink.OriginalLine
                    $newLine = $oldLine -replace "\(#$currentLink\)", "(#$expectedAnchor)"
                    
                    $newLines[$tocLink.Line] = $newLine
                    $fileIssues++
                    $changed = $true
                    
                    if ($fileIssues -eq 1) {
                        Write-Host "  修复: $($file.Name)" -ForegroundColor Yellow
                    }
                    Write-Host "    [$($tocLink.Text)] #$currentLink -> #$expectedAnchor" -ForegroundColor Gray
                }
                else {
                    # 找不到匹配的标题，报告问题
                    if ($fileIssues -eq 0) {
                        Write-Host "  ⚠️  $($file.Name)" -ForegroundColor Yellow
                    }
                    Write-Host "    未找到标题: [$($tocLink.Text)] (#$currentLink)" -ForegroundColor Red
                    $fileIssues++
                }
            }
        }
        
        if ($changed) {
            $newContent = $newLines -join "`n"
            Set-Content -Path $file.FullName -Value $newContent -Encoding UTF8 -NoNewline
            $filesFixed++
            $totalIssues += $fileIssues
            Write-Host "  ✅ 已修复 $fileIssues 个链接" -ForegroundColor Green
        }
        elseif ($fileIssues -gt 0) {
            Write-Host "  ⚠️  $fileIssues 个链接无法自动修复" -ForegroundColor Yellow
        }
    }
    catch {
        Write-Host "  ❌ 错误: $($file.Name) - $($_.Exception.Message)" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "=" * 60 -ForegroundColor Cyan
Write-Host "✨ 目录链接修复完成!" -ForegroundColor Green
Write-Host "📊 统计:" -ForegroundColor Cyan
Write-Host "  - 扫描文件: $($mdFiles.Count)" -ForegroundColor White
Write-Host "  - 修复文件: $filesFixed" -ForegroundColor Green
Write-Host "  - 修复链接: $totalIssues" -ForegroundColor Green
Write-Host "=" * 60 -ForegroundColor Cyan


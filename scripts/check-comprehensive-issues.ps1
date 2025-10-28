# 全面检查指定文件的各种问题

$files = @(
    'docs/fundamentals/language/00-Go-1.25.3形式化理论体系/08-学习路线图.md',
    'docs/fundamentals/language/00-Go-1.25.3核心机制完整解析/README.md',
    'docs/fundamentals/language/01-语法基础/00-概念定义体系.md'
)

$allIssues = @{
    broken_links = @()
    format_issues = @()
    reference_issues = @()
    consistency_issues = @()
}

Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "🔍 全面检查3个文件的问题" -ForegroundColor Yellow
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host ""

foreach ($file in $files) {
    if (-not (Test-Path $file)) {
        Write-Host "⚠️ 文件不存在: $file" -ForegroundColor Red
        continue
    }
    
    $content = Get-Content $file -Raw -Encoding UTF8
    $lines = Get-Content $file -Encoding UTF8
    $dir = Split-Path $file -Parent
    
    Write-Host "检查文件: $file" -ForegroundColor Green
    
    # 检查1: 本地文件链接
    $localLinks = [regex]::Matches($content, '\[([^\]]+)\]\(([^)]+\.md)\)')
    foreach ($match in $localLinks) {
        $linkText = $match.Groups[1].Value
        $linkPath = $match.Groups[2].Value
        
        # 跳过外部链接和锚点
        if ($linkPath -match '^https?://' -or $linkPath -match '^#') {
            continue
        }
        
        # 计算完整路径
        $fullPath = Join-Path $dir $linkPath
        $fullPath = $fullPath -replace '\\', '/'
        
        if (-not (Test-Path $fullPath)) {
            $allIssues.broken_links += [PSCustomObject]@{
                File = $file
                LinkText = $linkText
                LinkPath = $linkPath
                ResolvedPath = $fullPath
            }
            Write-Host "  ❌ 断链: [$linkText]($linkPath)" -ForegroundColor Red
        }
    }
    
    # 检查2: 目录格式一致性
    if ($content -match '##\s+📚\s+目录' -or $content -match '##\s+📖\s+目录') {
        $allIssues.format_issues += [PSCustomObject]@{
            File = $file
            Issue = "目录标题使用了非标准emoji（应为📋）"
        }
        Write-Host "  ❌ 目录标题格式不一致" -ForegroundColor Red
    }
    
    # 检查3: 章节编号连续性
    $chapterPattern = '^###?\s+(\d+)[\.、]'
    $prevNum = 0
    $lineNum = 0
    $inToc = $false
    
    foreach ($line in $lines) {
        $lineNum++
        
        # 跳过目录区域
        if ($line -match '^##\s+📋\s+目录') {
            $inToc = $true
        }
        if ($inToc -and $line -match '^##\s+[^📋]') {
            $inToc = $false
        }
        if ($inToc) {
            continue
        }
        
        if ($line -match $chapterPattern) {
            $num = [int]$matches[1]
            if ($prevNum -gt 0 -and $num -ne $prevNum + 1 -and $num -ne 1) {
                $allIssues.consistency_issues += [PSCustomObject]@{
                    File = $file
                    Line = $lineNum
                    Issue = "章节编号从 $prevNum 跳到 $num"
                }
                Write-Host "  ⚠️ 章节编号跳跃: Line $lineNum (从 $prevNum 到 $num)" -ForegroundColor Yellow
            }
            $prevNum = $num
        }
    }
    
    # 检查4: 多余的空行（超过2个连续空行）
    $emptyLineCount = 0
    $lineNum = 0
    foreach ($line in $lines) {
        $lineNum++
        if ($line.Trim() -eq '') {
            $emptyLineCount++
            if ($emptyLineCount -gt 2) {
                $allIssues.format_issues += [PSCustomObject]@{
                    File = $file
                    Line = $lineNum
                    Issue = "超过2个连续空行"
                }
            }
        } else {
            $emptyLineCount = 0
        }
    }
    
    # 检查5: 行尾空格
    $lineNum = 0
    $trailingSpaceCount = 0
    foreach ($line in $lines) {
        $lineNum++
        if ($line -match '\s+$') {
            $trailingSpaceCount++
        }
    }
    if ($trailingSpaceCount -gt 0) {
        $allIssues.format_issues += [PSCustomObject]@{
            File = $file
            Issue = "$trailingSpaceCount 行有行尾空格"
        }
        Write-Host "  ⚠️ $trailingSpaceCount 行有行尾空格" -ForegroundColor Yellow
    }
    
    Write-Host ""
}

Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "📊 检查结果汇总" -ForegroundColor Yellow
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host ""

$totalIssues = $allIssues.broken_links.Count + $allIssues.format_issues.Count + 
               $allIssues.reference_issues.Count + $allIssues.consistency_issues.Count

if ($totalIssues -eq 0) {
    Write-Host "✅ 未发现问题！" -ForegroundColor Green
} else {
    Write-Host "发现 $totalIssues 个问题:" -ForegroundColor Red
    Write-Host ""
    
    if ($allIssues.broken_links.Count -gt 0) {
        Write-Host "断链 ($($allIssues.broken_links.Count)个):" -ForegroundColor Red
        $allIssues.broken_links | Format-Table -AutoSize
    }
    
    if ($allIssues.format_issues.Count -gt 0) {
        Write-Host "格式问题 ($($allIssues.format_issues.Count)个):" -ForegroundColor Yellow
        $allIssues.format_issues | Format-Table -AutoSize
    }
    
    if ($allIssues.consistency_issues.Count -gt 0) {
        Write-Host "一致性问题 ($($allIssues.consistency_issues.Count)个):" -ForegroundColor Yellow
        $allIssues.consistency_issues | Format-Table -AutoSize
    }
    
    # 保存到JSON
    $allIssues | ConvertTo-Json -Depth 3 | Out-File 'comprehensive-issues.json' -Encoding UTF8
    Write-Host "问题已保存到: comprehensive-issues.json" -ForegroundColor Cyan
}

Write-Host ""


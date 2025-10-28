# 扫描指定3个文件的Markdown问题

$files = @(
    'docs/fundamentals/language/00-Go-1.25.3形式化理论体系/08-学习路线图.md',
    'docs/fundamentals/language/00-Go-1.25.3核心机制完整解析/README.md',
    'docs/fundamentals/language/01-语法基础/00-概念定义体系.md'
)

$issues = @()

Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "🔍 扫描指定的3个文件中的Markdown问题" -ForegroundColor Yellow
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host ""

foreach ($file in $files) {
    if (Test-Path $file) {
        $content = Get-Content $file -Raw -Encoding UTF8
        $lines = Get-Content $file -Encoding UTF8
        
        Write-Host "检查文件: $file" -ForegroundColor Green
        
        # 检查1: 多余的破折号或特殊字符
        if ($content -match '\*\*.*\*\*-\s*$') {
            $issues += [PSCustomObject]@{
                File = $file
                Type = '多余破折号'
                Line = ($content -split '\r?\n').Count
                Issue = '文档末尾有多余的破折号'
            }
            Write-Host "  ❌ 发现多余破折号" -ForegroundColor Red
        }
        
        # 检查2: HTML注释
        if ($content -match '<!--') {
            $lineNum = 0
            foreach ($line in $lines) {
                $lineNum++
                if ($line -match '<!--') {
                    $trimmedLine = $line.Trim()
                    $issues += [PSCustomObject]@{
                        File = $file
                        Type = 'HTML注释'
                        Line = $lineNum
                        Issue = "Line ${lineNum}: $trimmedLine"
                    }
                }
            }
            Write-Host "  ❌ 发现HTML注释" -ForegroundColor Red
        }
        
        # 检查3: 目录格式问题
        $inToc = $false
        $tocFormat = ''
        $lineNum = 0
        $mixedFormat = $false
        
        foreach ($line in $lines) {
            $lineNum++
            if ($line -match '^##\s+📋\s+目录\s*$') {
                $inToc = $true
                continue
            }
            if ($inToc -and $line -match '^##\s+[^📋]') {
                $inToc = $false
            }
            if ($inToc -and $line -match '^-\s+\[') {
                if (-not $tocFormat) {
                    $tocFormat = 'list'
                } elseif ($tocFormat -ne 'list') {
                    $mixedFormat = $true
                }
            } elseif ($inToc -and $line -match '^\d+\.') {
                if (-not $tocFormat) {
                    $tocFormat = 'numbered'
                } elseif ($tocFormat -ne 'numbered') {
                    $mixedFormat = $true
                }
            }
        }
        
        if ($mixedFormat) {
            $issues += [PSCustomObject]@{
                File = $file
                Type = '混合TOC格式'
                Line = 0
                Issue = '目录使用了混合格式（列表+编号）'
            }
            Write-Host "  ❌ 发现混合TOC格式" -ForegroundColor Red
        }
        
        # 检查4: 章节编号跳跃
        $chapterNumbers = @()
        $lineNum = 0
        foreach ($line in $lines) {
            $lineNum++
            if ($line -match '^###?\s+(\d+)[\.、]\s+') {
                $chapterNumbers += [PSCustomObject]@{
                    Number = [int]$matches[1]
                    Line = $lineNum
                }
            }
        }
        
        for ($i = 1; $i -lt $chapterNumbers.Count; $i++) {
            if ($chapterNumbers[$i].Number -ne $chapterNumbers[$i-1].Number + 1 -and 
                $chapterNumbers[$i].Number -ne $chapterNumbers[$i-1].Number -and
                $chapterNumbers[$i].Number -ne 1) {
                $prevNum = $chapterNumbers[$i-1].Number
                $currNum = $chapterNumbers[$i].Number
                $issues += [PSCustomObject]@{
                    File = $file
                    Type = '章节编号跳跃'
                    Line = $chapterNumbers[$i].Line
                    Issue = "从 $prevNum 跳到 $currNum"
                }
            }
        }
        
        Write-Host ""
    }
}

Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "📊 扫描结果统计" -ForegroundColor Yellow
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan

if ($issues.Count -eq 0) {
    Write-Host "✅ 未发现问题！" -ForegroundColor Green
} else {
    Write-Host "❌ 发现 $($issues.Count) 个问题" -ForegroundColor Red
    Write-Host ""
    $issues | Format-Table -AutoSize
    
    # 保存到JSON
    $issues | ConvertTo-Json -Depth 3 | Out-File 'markdown-issues-3files.json' -Encoding UTF8
    Write-Host "问题已保存到: markdown-issues-3files.json" -ForegroundColor Yellow
}

Write-Host ""


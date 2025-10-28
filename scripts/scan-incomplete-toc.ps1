# 扫描不完整的目录结构
Write-Output "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
Write-Output "🔍 扫描目录结构问题"
Write-Output "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
Write-Output ""

$docsFiles = Get-ChildItem -Path "docs" -Filter "*.md" -Recurse

$issues = @{
    "simplified_toc" = @()      # 简化版目录（只有1层）
    "incomplete_toc" = @()      # 不完整目录（缺少子标题）
    "wrong_format" = @()        # 格式错误
}

$total = 0
$checked = 0

foreach ($file in $docsFiles) {
    $checked++
    $content = Get-Content $file.FullName -Raw -Encoding UTF8
    
    # 检查是否有目录
    if ($content -notmatch '##\s*📋\s*目录') {
        continue
    }
    
    $total++
    
    # 提取目录部分
    if ($content -match '(?s)##\s*📋\s*目录\s*\n(.*?)\n##') {
        $tocSection = $matches[1]
        
        # 统计目录行数和层级
        $tocLines = ($tocSection -split '\n' | Where-Object { $_ -match '^\s*-\s*\[' }).Count
        $hasNested = $tocSection -match '^\s{2,}-\s*\['
        
        # 提取文档中所有标题
        $allHeadings = [regex]::Matches($content, '^(#{2,6})\s+(.+)$', [System.Text.RegularExpressions.RegexOptions]::Multiline)
        $headingCount = ($allHeadings | Where-Object { $_.Groups[2].Value -notmatch '📋\s*目录' }).Count
        
        # 判断问题类型
        if ($tocLines -le 5 -and $headingCount -gt 10) {
            # 简化版目录：目录项太少，但文档标题很多
            $issues["simplified_toc"] += [PSCustomObject]@{
                File = $file.FullName -replace [regex]::Escape($PWD.Path + '\'), ''
                TocLines = $tocLines
                HeadingCount = $headingCount
            }
        }
        elseif (-not $hasNested -and $headingCount -gt 15) {
            # 无嵌套但标题很多
            $issues["incomplete_toc"] += [PSCustomObject]@{
                File = $file.FullName -replace [regex]::Escape($PWD.Path + '\'), ''
                TocLines = $tocLines
                HeadingCount = $headingCount
            }
        }
        elseif ($tocLines -lt ($headingCount * 0.5)) {
            # 目录行数少于标题的50%
            $issues["incomplete_toc"] += [PSCustomObject]@{
                File = $file.FullName -replace [regex]::Escape($PWD.Path + '\'), ''
                TocLines = $tocLines
                HeadingCount = $headingCount
            }
        }
    }
}

Write-Output "扫描完成！"
Write-Output ""
Write-Output "📊 统计结果:"
Write-Output "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
Write-Output "总文件数: $checked"
Write-Output "有目录文件: $total"
Write-Output ""
Write-Output "问题分类:"
Write-Output "  简化版目录: $($issues['simplified_toc'].Count) 个"
Write-Output "  不完整目录: $($issues['incomplete_toc'].Count) 个"
Write-Output ""

if ($issues["simplified_toc"].Count -gt 0) {
    Write-Output "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    Write-Output "📋 简化版目录 ($($issues['simplified_toc'].Count) 个):"
    Write-Output "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    foreach ($item in $issues["simplified_toc"] | Sort-Object File) {
        Write-Output "  $($item.File)"
        Write-Output "    目录项: $($item.TocLines) | 标题数: $($item.HeadingCount)"
    }
    Write-Output ""
}

if ($issues["incomplete_toc"].Count -gt 0) {
    Write-Output "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    Write-Output "📋 不完整目录 ($($issues['incomplete_toc'].Count) 个):"
    Write-Output "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    foreach ($item in $issues["incomplete_toc"] | Sort-Object File) {
        Write-Output "  $($item.File)"
        Write-Output "    目录项: $($item.TocLines) | 标题数: $($item.HeadingCount)"
    }
    Write-Output ""
}

# 保存结果
$issues | ConvertTo-Json -Depth 10 | Out-File "toc-issues.json" -Encoding UTF8

Write-Output "结果已保存到: toc-issues.json"
Write-Output "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"


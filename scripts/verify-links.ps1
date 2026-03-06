# 验证文档内部链接
$docsPath = "docs"
$issues = @()

# 获取所有markdown文件
$allFiles = Get-ChildItem -Path $docsPath -Recurse -Filter "*.md" | Select-Object -ExpandProperty FullName

# 检查每个文件中的链接
foreach ($file in $allFiles) {
    $content = Get-Content $file -Raw -Encoding UTF8
    $dir = Split-Path $file -Parent
    
    # 匹配相对路径链接 [text](./path/to/file.md)
    $relativeLinks = [regex]::Matches($content, '\[([^\]]+)\]\(\.\/([^)]+)\)')
    foreach ($match in $relativeLinks) {
        $linkPath = $match.Groups[2].Value
        $fullLinkPath = Join-Path $dir $linkPath
        
        # 移除锚点
        $fullLinkPath = $fullLinkPath -replace '#.*$', ''
        
        if (-not (Test-Path $fullLinkPath)) {
            $issues += [PSCustomObject]@{
                Source = $file
                Link = $linkPath
                Type = "Relative"
            }
        }
    }
    
    # 匹配文档链接 [text](file.md)
    $fileLinks = [regex]::Matches($content, '\[([^\]]+)\]\(([^./][^)]+\.md)\)')
    foreach ($match in $fileLinks) {
        $linkPath = $match.Groups[2].Value
        $fullLinkPath = Join-Path $dir $linkPath
        
        if (-not (Test-Path $fullLinkPath)) {
            $issues += [PSCustomObject]@{
                Source = $file
                Link = $linkPath
                Type = "File"
            }
        }
    }
}

# 输出结果
Write-Host "=== 链接验证结果 ===" -ForegroundColor Cyan
if ($issues.Count -eq 0) {
    Write-Host "✓ 所有链接有效！" -ForegroundColor Green
} else {
    Write-Host "发现 $($issues.Count) 个无效链接:" -ForegroundColor Red
    $issues | Group-Object Link | ForEach-Object {
        Write-Host "`n  ✗ $($_.Name)" -ForegroundColor Red
        $_.Group | Select-Object -First 1 -ExpandProperty Source | ForEach-Object {
            Write-Host "    引用自: $_" -ForegroundColor Yellow
        }
    }
}

$issues

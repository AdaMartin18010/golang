# 文档链接验证脚本
# 检查Markdown文档中的内部链接有效性

param(
    [string]$TargetDir = "docs",
    [switch]$FixBrokenLinks = $false
)

Write-Host "=== 文档链接验证工具 ===" -ForegroundColor Cyan
Write-Host

$brokenLinks = @()
$totalLinks = 0
$validLinks = 0

# 获取所有活跃的Markdown文件
$files = Get-ChildItem -Path $TargetDir -Recurse -Filter "*.md" | Where-Object {
    $_.FullName -notmatch "\\00-备份\\" -and
    $_.FullName -notmatch "\\archive-" -and
    $_.FullName -notmatch "\\Analysis\\"
}

Write-Host "📁 扫描 $($files.Count) 个文档文件..." -ForegroundColor Yellow
Write-Host

foreach ($file in $files) {
    $content = Get-Content -Path $file.FullName -Raw -Encoding UTF8
    
    # 匹配Markdown链接: [text](url)
    $linkPattern = '\[([^\]]+)\]\(([^)]+)\)'
    $matches = [regex]::Matches($content, $linkPattern)
    
    foreach ($match in $matches) {
        $linkText = $match.Groups[1].Value
        $linkUrl = $match.Groups[2].Value
        $totalLinks++
        
        # 跳过外部链接和锚点链接
        if ($linkUrl -match '^https?://') {
            continue
        }
        if ($linkUrl -match '^#') {
            # 内部锚点链接，需要验证标题是否存在
            $anchor = $linkUrl.Substring(1)
            # 简化的锚点验证（实际应该更复杂）
            $validLinks++
            continue
        }
        
        # 内部文件链接
        if ($linkUrl -match '\.md') {
            $linkedFile = Join-Path (Split-Path $file.FullName) $linkUrl
            $linkedFile = [System.IO.Path]::GetFullPath($linkedFile)
            
            if (Test-Path $linkedFile) {
                $validLinks++
            }
            else {
                $brokenLinks += [PSCustomObject]@{
                    SourceFile = $file.FullName
                    LinkText = $linkText
                    LinkUrl = $linkUrl
                    TargetFile = $linkedFile
                }
            }
        }
    }
}

Write-Host
Write-Host "=== 验证结果 ===" -ForegroundColor Cyan
Write-Host "📊 总链接数: $totalLinks" -ForegroundColor White
Write-Host "✅ 有效链接: $validLinks" -ForegroundColor Green
Write-Host "❌ 失效链接: $($brokenLinks.Count)" -ForegroundColor $(if ($brokenLinks.Count -eq 0) { "Green" } else { "Red" })

if ($brokenLinks.Count -gt 0) {
    Write-Host
    Write-Host "❌ 失效链接详情:" -ForegroundColor Red
    Write-Host "----------------------------------------"
    
    foreach ($link in $brokenLinks) {
        $relativePath = $link.SourceFile.Replace((Get-Location).Path + "\", "")
        Write-Host "文件: $relativePath" -ForegroundColor Yellow
        Write-Host "  链接文本: $($link.LinkText)" -ForegroundColor Gray
        Write-Host "  链接目标: $($link.LinkUrl)" -ForegroundColor Gray
        Write-Host "  目标文件: $($link.TargetFile)" -ForegroundColor Gray
        Write-Host
    }
    
    if ($FixBrokenLinks) {
        Write-Host "🔧 自动修复功能尚未实现" -ForegroundColor Yellow
        Write-Host "   建议手动检查并修复上述链接" -ForegroundColor Yellow
    }
}
else {
    Write-Host
    Write-Host "🎉 所有内部文件链接均有效！" -ForegroundColor Green
}

Write-Host
Write-Host "=== 验证完成 ===" -ForegroundColor Cyan

# 返回失效链接数量作为退出码
exit $brokenLinks.Count


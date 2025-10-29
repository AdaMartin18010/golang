# 检查文档中的链接有效性
param(
    [string]$Path = "docs",
    [int]$MaxFiles = 50
)

$stats = @{
    Checked = 0
    BrokenLinks = 0
    BrokenFiles = @()
}

Write-Host "🔍 检查文档链接有效性..." -ForegroundColor Cyan

$files = Get-ChildItem -Path $Path -Recurse -Filter "*.md" -File | Select-Object -First $MaxFiles

foreach ($file in $files) {
    $stats.Checked++
    $content = Get-Content $file.FullName -Raw -Encoding UTF8
    
    # 提取所有内部链接 [text](#anchor)
    $internalLinks = [regex]::Matches($content, '\[([^\]]+)\]\(#([^\)]+)\)')
    
    # 提取所有anchor定义 (## heading)
    $anchors = [regex]::Matches($content, '(?m)^##+ (.+)$')
    
    # 将heading转换为GitHub风格的anchor
    $validAnchors = @{}
    foreach ($anchor in $anchors) {
        $heading = $anchor.Groups[1].Value
        # 转换为anchor格式：小写、空格替换为-、移除特殊字符（保留emoji和中文）
        $anchorId = $heading.ToLower() -replace '\s+', '-' -replace '[^\w\-\u4e00-\u9fa5]', ''
        $validAnchors[$anchorId] = $true
    }
    
    # 检查链接
    $fileBroken = @()
    foreach ($link in $internalLinks) {
        $linkText = $link.Groups[1].Value
        $linkAnchor = $link.Groups[2].Value
        
        if (-not $validAnchors.ContainsKey($linkAnchor)) {
            $fileBroken += "  ❌ [$linkText](#$linkAnchor)"
            $stats.BrokenLinks++
        }
    }
    
    if ($fileBroken.Count -gt 0) {
        Write-Host "`n📄 $($file.Name):" -ForegroundColor Yellow
        $fileBroken | ForEach-Object { Write-Host $_ -ForegroundColor Red }
        $stats.BrokenFiles += $file.FullName
    }
}

Write-Host "`n📊 统计:" -ForegroundColor Cyan
Write-Host "  检查文件: $($stats.Checked)" -ForegroundColor Gray
Write-Host "  失效链接: $($stats.BrokenLinks)" -ForegroundColor $(if($stats.BrokenLinks -gt 0){'Red'}else{'Green'})
Write-Host "  问题文件: $($stats.BrokenFiles.Count)" -ForegroundColor $(if($stats.BrokenFiles.Count -gt 0){'Red'}else{'Green'})

if ($stats.BrokenFiles.Count -gt 0) {
    Write-Host "`n问题文件列表:" -ForegroundColor Yellow
    $stats.BrokenFiles | ForEach-Object { Write-Host "  - $_" -ForegroundColor Gray }
}


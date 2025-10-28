# 验证修复结果

Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "✅ 验证修复结果" -ForegroundColor Green
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host ""

# 验证1: 破折号已删除
Write-Host "1. 验证08-学习路线图.md的破折号已修复:" -ForegroundColor Yellow
$content1 = Get-Content 'docs/fundamentals/language/00-Go-1.25.3形式化理论体系/08-学习路线图.md' -Raw -Encoding UTF8
if ($content1 -match '\*\*坚持学习[^*]+\*\*-') {
    Write-Host "  ❌ 仍有多余的破折号" -ForegroundColor Red
} else {
    Write-Host "  ✅ 破折号已删除" -ForegroundColor Green
}

# 验证2: 目录格式统一
Write-Host ""
Write-Host "2. 验证目录格式统一:" -ForegroundColor Yellow

$files = @(
    'docs/fundamentals/language/00-Go-1.25.3形式化理论体系/08-学习路线图.md',
    'docs/fundamentals/language/00-Go-1.25.3核心机制完整解析/README.md',
    'docs/fundamentals/language/01-语法基础/00-概念定义体系.md'
)

$correctCount = 0
foreach ($file in $files) {
    $c = Get-Content $file -Raw -Encoding UTF8
    $fileName = Split-Path $file -Leaf
    if ($c -match '##\s+📋\s+目录') {
        Write-Host "  ✅ $fileName - 使用标准📋目录" -ForegroundColor Green
        $correctCount++
    } else {
        Write-Host "  ❌ $fileName - 目录格式不标准" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "总计: $correctCount/$($files.Count) 个文件使用标准目录格式" -ForegroundColor Cyan

# 验证3: 没有HTML注释
Write-Host ""
Write-Host "3. 验证没有HTML注释:" -ForegroundColor Yellow
$htmlCommentCount = 0
foreach ($file in $files) {
    $c = Get-Content $file -Raw -Encoding UTF8
    if ($c -match '<!--') {
        $htmlCommentCount++
        $fileName = Split-Path $file -Leaf
        Write-Host "  ❌ $fileName - 包含HTML注释" -ForegroundColor Red
    }
}
if ($htmlCommentCount -eq 0) {
    Write-Host "  ✅ 所有文件都没有HTML注释" -ForegroundColor Green
}

Write-Host ""
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "✨ 验证完成！" -ForegroundColor Green
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host ""


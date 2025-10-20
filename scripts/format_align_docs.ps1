# 文档格式对齐脚本
# 自动对齐docs目录下所有Markdown文档的格式

param(
    [string]$TargetDir = "docs",
    [switch]$DryRun = $false
)

Write-Host "=== 文档格式对齐工具 ===" -ForegroundColor Cyan
Write-Host

$processedCount = 0
$skippedCount = 0
$errorCount = 0

# 获取所有需要处理的Markdown文件
$files = Get-ChildItem -Path $TargetDir -Recurse -Filter "*.md" | Where-Object {
    # 排除归档目录
    $_.FullName -notmatch "\\00-备份\\" -and
    $_.FullName -notmatch "\\archive-" -and
    $_.FullName -notmatch "\\Analysis\\"
}

Write-Host "📁 找到 $($files.Count) 个文档文件需要处理" -ForegroundColor Yellow
Write-Host

foreach ($file in $files) {
    try {
        Write-Host "处理: $($file.FullName)" -ForegroundColor Gray
        
        $content = Get-Content -Path $file.FullName -Raw -Encoding UTF8
        $modified = $false
        
        # 1. 检查是否有底部元信息
        if ($content -notmatch "文档维护者.*最后更新.*文档状态") {
            Write-Host "  ✓ 需要添加底部元信息" -ForegroundColor Yellow
            
            # 移除旧的元信息格式
            $content = $content -replace "(?ms)\*\*(?:模块)?维护者\*\*:.*?(?=\r?\n\r?\n|$)", ""
            $content = $content -replace "(?ms)\*\*最后更新\*\*:.*?(?=\r?\n|$)", ""
            $content = $content -replace "(?ms)\*\*(?:模块|文档)?状态\*\*:.*?(?=\r?\n|$)", ""
            
            # 添加标准元信息
            $metadata = @"

---

**文档维护者**: Go Documentation Team  
**最后更新**: 2025年10月20日  
**文档状态**: 完成  
**适用版本**: Go 1.21+
"@
            $content = $content.TrimEnd() + "`n" + $metadata + "`n"
            $modified = $true
        }
        
        # 2. 检查标题格式
        if ($content -match "^#+\s+\d+(\.\d+)+") {
            Write-Host "  ✓ 需要修正标题编号" -ForegroundColor Yellow
            $modified = $true
        }
        
        # 3. 检查是否有简介
        if ($content -notmatch "^#[^#].*\n\n>\s*\*\*简介\*\*:") {
            Write-Host "  ✓ 需要添加简介" -ForegroundColor Yellow
            $modified = $true
        }
        
        if ($modified -and -not $DryRun) {
            Set-Content -Path $file.FullName -Value $content -Encoding UTF8 -NoNewline
            Write-Host "  ✅ 已更新" -ForegroundColor Green
            $processedCount++
        }
        elseif ($modified -and $DryRun) {
            Write-Host "  🔍 [DryRun] 将会更新" -ForegroundColor Cyan
            $processedCount++
        }
        else {
            Write-Host "  ⏭️  无需更新" -ForegroundColor DarkGray
            $skippedCount++
        }
    }
    catch {
        Write-Host "  ❌ 错误: $_" -ForegroundColor Red
        $errorCount++
    }
}

Write-Host
Write-Host "=== 处理完成 ===" -ForegroundColor Cyan
Write-Host "✅ 已处理: $processedCount 个文件" -ForegroundColor Green
Write-Host "⏭️  跳过: $skippedCount 个文件" -ForegroundColor Yellow
Write-Host "❌ 错误: $errorCount 个文件" -ForegroundColor Red

if ($DryRun) {
    Write-Host
    Write-Host "这是模拟运行，没有实际修改文件。" -ForegroundColor Cyan
    Write-Host "移除 -DryRun 参数以实际执行修改。" -ForegroundColor Cyan
}


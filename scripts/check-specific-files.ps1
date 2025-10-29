# 检查特定文件的版本信息

$files = @(
    "docs\fundamentals\language\03-模块管理\00-概念定义体系.md",
    "docs\reference\00-概念定义体系.md",
    "docs\🎉-文档格式统一完成-2025-10-29.md",
    "docs\📊-文档格式梳理总结报告-2025-10-29.md"
)

foreach ($file in $files) {
    if (Test-Path $file) {
        Write-Host "`n=== $file ===" -ForegroundColor Cyan
        $content = Get-Content -Path $file -Raw -Encoding UTF8
        $matches = [regex]::Matches($content, '\*\*版本\*\*:')
        Write-Host "版本信息出现次数: $($matches.Count)" -ForegroundColor Yellow
        
        # 显示每次出现的位置（前50个字符）
        foreach ($match in $matches) {
            $start = [Math]::Max(0, $match.Index - 20)
            $length = [Math]::Min(70, $content.Length - $start)
            $context = $content.Substring($start, $length) -replace "`n", " " -replace "`r", ""
            Write-Host "  -> ...${context}..." -ForegroundColor Gray
        }
    } else {
        Write-Host "`n文件不存在: $file" -ForegroundColor Red
    }
}


# 全面修复3个文件的Markdown问题

$files = @(
    'docs/fundamentals/language/00-Go-1.25.3形式化理论体系/08-学习路线图.md',
    'docs/fundamentals/language/00-Go-1.25.3核心机制完整解析/README.md',
    'docs/fundamentals/language/01-语法基础/00-概念定义体系.md'
)

$totalFixed = 0

Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "🔧 全面修复3个文件的Markdown问题" -ForegroundColor Yellow
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host ""

foreach ($file in $files) {
    if (Test-Path $file) {
        $content = Get-Content $file -Raw -Encoding UTF8
        $originalContent = $content
        $fileFixed = 0
        
        Write-Host "处理文件: $file" -ForegroundColor Green
        
        # 修复1: 删除多余的破折号
        if ($content -match '\*\*[^*]+\*\*-\s*\n') {
            $content = $content -replace '(\*\*[^*]+\*\*)-\s*\n', '$1' + "`n"
            $fileFixed++
            Write-Host "  ✅ 修复: 删除多余的破折号" -ForegroundColor Yellow
        }
        
        # 修复2: 删除HTML注释（如果存在于TOC中）
        if ($content -match '<!--\s*TOC\s*START\s*-->') {
            $content = $content -replace '<!--\s*TOC\s*START\s*-->\r?\n?', ''
            $content = $content -replace '<!--\s*TOC\s*END\s*-->\r?\n?', ''
            $fileFixed++
            Write-Host "  ✅ 修复: 删除HTML注释TOC" -ForegroundColor Yellow
        }
        
        # 修复3: 统一目录标题格式为 "## 📋 目录"
        $content = $content -replace '##\s+📚\s+目录', '## 📋 目录'
        $content = $content -replace '##\s+📖\s+目录', '## 📋 目录'
        
        # 修复4: 确保目录后有空行
        $content = $content -replace '(##\s+📋\s+目录)\r?\n([^\r\n])', '$1' + "`n`n" + '$2'
        
        # 修复5: 删除多余的空行（超过2个连续空行）
        $content = $content -replace '(\r?\n\s*){3,}', "`n`n"
        
        # 修复6: 确保文件末尾有且仅有一个空行
        $content = $content.TrimEnd() + "`n"
        
        # 修复7: 删除行尾空格
        $lines = $content -split '\r?\n'
        $lines = $lines | ForEach-Object { $_.TrimEnd() }
        $content = $lines -join "`n"
        
        # 保存修改
        if ($content -ne $originalContent) {
            $content | Out-File $file -Encoding UTF8 -NoNewline
            $totalFixed++
            Write-Host "  📝 文件已更新" -ForegroundColor Green
        } else {
            Write-Host "  ✅ 文件无需修改" -ForegroundColor Cyan
        }
        
        Write-Host ""
    } else {
        Write-Host "⚠️ 文件不存在: $file" -ForegroundColor Red
    }
}

Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "✨ 修复完成" -ForegroundColor Yellow
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host ""
Write-Host "修复文件数: $totalFixed" -ForegroundColor Green
Write-Host ""


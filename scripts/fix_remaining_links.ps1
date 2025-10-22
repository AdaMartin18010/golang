# scripts/fix_remaining_links.ps1
# 修复剩余的失效链接

param (
    [string]$DocsPath = "docs-new"
)

Write-Host "=== 🔧 修复剩余失效链接 ===" -ForegroundColor Cyan
Write-Host ""

$fixedCount = 0

# 1. 修复Go版本特性READMEs中指向不存在文件的链接
Write-Host "1️⃣ 修复Go版本特性READMEs..." -ForegroundColor Yellow

$versionReadmes = @(
    "10-Go版本特性\01-Go-1.21特性\README.md",
    "10-Go版本特性\02-Go-1.22特性\README.md",
    "10-Go版本特性\03-Go-1.23特性\README.md",
    "10-Go版本特性\04-Go-1.24特性\README.md",
    "10-Go版本特性\05-Go-1.25特性\README.md"
)

foreach ($relPath in $versionReadmes) {
    $filePath = Join-Path $DocsPath $relPath
    
    if (Test-Path $filePath) {
        $content = Get-Content $filePath -Raw
        $originalContent = $content
        
        # 移除指向不存在的详细文件的链接
        # 匹配模式: [文本](./01-文件名.md) 或 [文本](./文件名.md)
        $content = $content -replace '\[([^\]]+)\]\(\./\d{2}-[^)]+\.md\)', ''
        $content = $content -replace '\[([^\]]+)\]\(\./[^)]+\.md\)', ''
        
        # 清理多余的空行
        $content = $content -replace '(\r?\n){3,}', "`n`n"
        
        if ($content -ne $originalContent) {
            Set-Content -Path $filePath -Value $content -Encoding UTF8
            Write-Host "  ✅ 已修复: $relPath" -ForegroundColor Green
            $fixedCount++
        }
    }
}

# 2. 修复代码片段误识别为链接的问题
Write-Host ""
Write-Host "2️⃣ 修复代码片段误识别..." -ForegroundColor Yellow

$codeLinkFiles = @(
    "01-语言基础\00-Go语言形式化语义与理论基础.md",
    "03-Web开发\00-HTTP编程深度实战指南.md",
    "07-性能优化\01-性能分析与pprof.md",
    "08-架构设计\01-创建型模式.md",
    "08-架构设计\03-行为型模式.md",
    "10-Go版本特性\01-Go-1.21特性\README.md",
    "10-Go版本特性\03-Go-1.23特性\README.md",
    "10-Go版本特性\04-实验性特性\README.md",
    "10-Go版本特性\05-实践应用\README.md",
    "FAQ.md"
)

$codeReplacements = @(
    @{ Pattern = '\[v T\]\(v T\)'; Replacement = '`v T`' }
    @{ Pattern = '\[initFunc func\(\]\(initFunc func\(\)\)'; Replacement = '`initFunc func(`' }
    @{ Pattern = '\[handler\]\(handler\)'; Replacement = '`handler`' }
    @{ Pattern = '\[params\]\(params\)'; Replacement = '`params`' }
    @{ Pattern = '\[x, y T\]\(x, y T\)'; Replacement = '`x, y T`' }
    @{ Pattern = '\[s \[\]T, f func\(T\]\(s \[\]T, f func\(T\)\)'; Replacement = '`s []T, f func(T`' }
    @{ Pattern = '\[m map\[K\]V\]\(m map\[K\]V\)'; Replacement = '`m map[K]V`' }
    @{ Pattern = '\[seq iter\.Seq\[T\], pred func\(T\]\(seq iter\.Seq\[T\], pred func\(T\)\)'; Replacement = '`seq iter.Seq[T], pred func(T`' }
    @{ Pattern = '\[seq iter\.Seq\[T\], fn func\(T\]\(seq iter\.Seq\[T\], fn func\(T\)\)'; Replacement = '`seq iter.Seq[T], fn func(T`' }
    @{ Pattern = '\[values \.\.\.T\]\(values \.\.\.T\)'; Replacement = '`values ...T`' }
    @{ Pattern = '\[slice \[\]T\]\(slice \[\]T\)'; Replacement = '`slice []T`' }
    @{ Pattern = '\[slice \[\]T, fn func\(T\]\(slice \[\]T, fn func\(T\)\)'; Replacement = '`slice []T, fn func(T`' }
    @{ Pattern = '\[k K, v V\]\(k K, v V\)'; Replacement = '`k K, v V`' }
    @{ Pattern = '\[items \[\]T\]\(items \[\]T\)'; Replacement = '`items []T`' }
)

foreach ($relPath in $codeLinkFiles) {
    $filePath = Join-Path $DocsPath $relPath
    
    if (Test-Path $filePath) {
        $content = Get-Content $filePath -Raw
        $originalContent = $content
        
        foreach ($rep in $codeReplacements) {
            $content = $content -replace $rep.Pattern, $rep.Replacement
        }
        
        if ($content -ne $originalContent) {
            Set-Content -Path $filePath -Value $content -Encoding UTF8
            Write-Host "  ✅ 已修复: $relPath" -ForegroundColor Green
            $fixedCount++
        }
    }
}

# 3. 修复其他特殊链接
Write-Host ""
Write-Host "3️⃣ 修复其他特殊链接..." -ForegroundColor Yellow

$otherFixes = @(
    @{ 
        File = "01-语言基础\01-语法基础\01-Hello-World.md"
        Old = "./README.md#11134-4-包和模块"
        New = "../README.md"
    }
)

foreach ($fix in $otherFixes) {
    $filePath = Join-Path $DocsPath $fix.File
    
    if (Test-Path $filePath) {
        $content = Get-Content $filePath -Raw
        
        if ($content.Contains($fix.Old)) {
            $content = $content.Replace($fix.Old, $fix.New)
            Set-Content -Path $filePath -Value $content -Encoding UTF8
            Write-Host "  ✅ 已修复: $($fix.File)" -ForegroundColor Green
            $fixedCount++
        }
    }
}

Write-Host ""
Write-Host "=== 完成 ===" -ForegroundColor Cyan
Write-Host "  已修复文件: $fixedCount" -ForegroundColor Green


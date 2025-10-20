# 文档格式对齐增强版脚本 v2.0
# 全面递归迭代对齐所有文档的格式，包括TOC、标题、元信息等

param(
    [string]$TargetDir = "docs",
    [switch]$DryRun = $false,
    [switch]$FixTOC = $true
)

Write-Host "=== 文档格式对齐工具 v2.0 ===" -ForegroundColor Cyan
Write-Host "🔧 增强功能：TOC格式修正、标题对齐、元信息统一" -ForegroundColor Yellow
Write-Host

$processedCount = 0
$skippedCount = 0
$errorCount = 0
$tocFixedCount = 0

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
        Write-Host "处理: $($file.Name)" -ForegroundColor Gray
        
        $content = Get-Content -Path $file.FullName -Raw -Encoding UTF8
        if ([string]::IsNullOrWhiteSpace($content)) {
            Write-Host "  ⏭️  空文件，跳过" -ForegroundColor DarkGray
            $skippedCount++
            continue
        }
        
        $modified = $false
        $fixes = @()
        
        # 1. 修正TOC链接格式（emoji后面应该是双横杠--）
        if ($FixTOC -and $content -match '<!-- TOC START -->') {
            $tocPattern = '\[([\d\.]+)\s+(📚|💻|🔧|📊|🧪|🎯|⚠️|🔍|🚀|🏗️|🛡️|✅|❌|📝|🔗|🏭|📋|💡|🎨|📖|🌟|⭐|🔥|💎|🎓|📢|🎊|🎉)\s+([^\]]+)\]\(#([\d\-]+)-([^\)]+)\)'
            
            # 修正：确保emoji后面的链接锚点格式正确
            $content = $content -replace '\]\(#(\d+)-�', '](#$1--�'
            $content = $content -replace '\]\(#(\d+)-💻\)', '](#$1--💻)'
            $content = $content -replace '\]\(#(\d+)-🔧\)', '](#$1--🔧)'
            $content = $content -replace '\]\(#(\d+)-📊\)', '](#$1--📊)'
            $content = $content -replace '\]\(#(\d+)-🧪\)', '](#$1--🧪)'
            $content = $content -replace '\]\(#(\d+)-🎯\)', '](#$1--🎯)'
            $content = $content -replace '\]\(#(\d+)-⚠️\)', '](#$1--⚠️)'
            $content = $content -replace '\]\(#(\d+)-🔍\)', '](#$1--🔍)'
            $content = $content -replace '\]\(#(\d+)-🚀\)', '](#$1--🚀)'
            $content = $content -replace '\]\(#(\d+)-🏗️\)', '](#$1--🏗️)'
            $content = $content -replace '\]\(#(\d+)-🛡️\)', '](#$1--🛡️)'
            $content = $content -replace '\]\(#(\d+)-📚\)', '](#$1--📚)'
            $content = $content -replace '\]\(#(\d+)-🏭\)', '](#$1--🏭)'
            $content = $content -replace '\]\(#(\d+)-📋\)', '](#$1--📋)'
            
            if ($content -ne (Get-Content -Path $file.FullName -Raw -Encoding UTF8)) {
                $fixes += "TOC链接格式"
                $tocFixedCount++
                $modified = $true
            }
        }
        
        # 2. 统一底部元信息格式
        # 移除所有旧格式的元信息
        $content = $content -replace '(?m)^---\s*\n\n\*\*(?:文档|模块)?维护者\*\*:.*?$', ''
        $content = $content -replace '(?m)^\*\*(?:文档|模块)?维护者\*\*:.*?$', ''
        $content = $content -replace '(?m)^\*\*最后更新\*\*:.*?$', ''
        $content = $content -replace '(?m)^\*\*(?:文档|模块)?状态\*\*:.*?$', ''
        $content = $content -replace '(?m)^\*\*适用版本\*\*:.*?$', ''
        
        # 清理多余的空行和分隔线
        $content = $content -replace '(?m)^---\s*\n\s*\n---\s*$', '---'
        $content = $content -replace '\n{3,}---\s*$', "`n`n---"
        
        # 添加标准元信息（如果还没有）
        if ($content -notmatch '文档维护者.*Go Documentation Team') {
            $content = $content.TrimEnd()
            
            # 移除末尾多余的分隔线
            $content = $content -replace '---\s*$', ''
            $content = $content.TrimEnd()
            
            $metadata = @"

---

**文档维护者**: Go Documentation Team  
**最后更新**: 2025年10月20日  
**文档状态**: 完成  
**适用版本**: Go 1.21+
"@
            $content = $content + "`n" + $metadata + "`n"
            $fixes += "标准元信息"
            $modified = $true
        }
        
        # 3. 确保简介部分格式正确
        if ($content -match '^#[^#]' -and $content -notmatch '^#[^#].*\n\n>\s*\*\*简介\*\*:') {
            # 如果没有简介，添加占位符（但不修改，避免覆盖已有内容）
            # 这里只是检测，不自动添加，因为简介应该手动编写
        }
        
        # 4. 移除标题中的多级编号前缀（如 9.1、6.1.1 等）
        $content = $content -replace '(?m)^(#{1,6})\s+\d+(\.\d+)+\s+', '$1 '
        
        if ($modified -and -not $DryRun) {
            Set-Content -Path $file.FullName -Value $content -Encoding UTF8 -NoNewline
            Write-Host "  ✅ 已更新 [$($fixes -join ', ')]" -ForegroundColor Green
            $processedCount++
        }
        elseif ($modified -and $DryRun) {
            Write-Host "  🔍 [DryRun] 将更新 [$($fixes -join ', ')]" -ForegroundColor Cyan
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
Write-Host "🔧 TOC修正: $tocFixedCount 个文件" -ForegroundColor Yellow
Write-Host "⏭️  跳过: $skippedCount 个文件" -ForegroundColor Yellow
Write-Host "❌ 错误: $errorCount 个文件" -ForegroundColor Red

if ($DryRun) {
    Write-Host
    Write-Host "这是模拟运行，没有实际修改文件。" -ForegroundColor Cyan
    Write-Host "移除 -DryRun 参数以实际执行修改。" -ForegroundColor Cyan
}

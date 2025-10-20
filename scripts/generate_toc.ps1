# TOC（目录）自动生成工具
# 为Markdown文档自动生成标准化的目录结构

param(
    [string]$TargetDir = "docs",
    [switch]$DryRun = $false,
    [switch]$Force = $false
)

Write-Host "=== TOC 自动生成工具 ===" -ForegroundColor Cyan
Write-Host

$processedCount = 0
$skippedCount = 0
$errorCount = 0

# 获取所有活跃的Markdown文件
$files = Get-ChildItem -Path $TargetDir -Recurse -Filter "*.md" | Where-Object {
    $_.FullName -notmatch "\\00-备份\\" -and
    $_.FullName -notmatch "\\archive-" -and
    $_.FullName -notmatch "\\Analysis\\"
}

Write-Host "📁 找到 $($files.Count) 个文档文件" -ForegroundColor Yellow
Write-Host

function Generate-TOC {
    param(
        [string]$Content
    )
    
    $toc = @()
    $lines = $content -split "`n"
    
    foreach ($line in $lines) {
        # 匹配标题行 (## 到 ####)
        if ($line -match '^(#{2,4})\s+(.+)$') {
            $level = $matches[1].Length
            $title = $matches[2].Trim()
            
            # 跳过可能的元信息标题
            if ($title -match '^\*\*') {
                continue
            }
            
            # 生成锚点链接（GitHub风格）
            $anchor = $title.ToLower()
            $anchor = $anchor -replace '[^\w\s\u4e00-\u9fa5-]', ''  # 保留字母、数字、中文和连字符
            $anchor = $anchor -replace '\s+', '-'  # 空格转连字符
            $anchor = $anchor -replace '--+', '-'  # 多个连字符合并
            
            # 计算缩进
            $indent = '  ' * ($level - 2)
            
            # 添加到TOC
            $tocLine = "$indent- [$title](#$anchor)"
            $toc += $tocLine
        }
    }
    
    if ($toc.Count -gt 0) {
        return "<!-- TOC START -->`n" + ($toc -join "`n") + "`n<!-- TOC END -->"
    }
    
    return $null
}

foreach ($file in $files) {
    try {
        Write-Host "处理: $($file.Name)" -ForegroundColor Gray
        
        $content = Get-Content -Path $file.FullName -Raw -Encoding UTF8
        
        if ([string]::IsNullOrWhiteSpace($content)) {
            Write-Host "  ⏭️  空文件，跳过" -ForegroundColor DarkGray
            $skippedCount++
            continue
        }
        
        # 检查是否已有TOC
        $hasTOC = $content -match '<!-- TOC START -->'
        
        if ($hasTOC -and -not $Force) {
            Write-Host "  ⏭️  已有TOC，跳过（使用-Force强制更新）" -ForegroundColor DarkGray
            $skippedCount++
            continue
        }
        
        # 生成新的TOC
        $newTOC = Generate-TOC -Content $content
        
        if ($null -eq $newTOC) {
            Write-Host "  ⏭️  无二级以上标题，无需TOC" -ForegroundColor DarkGray
            $skippedCount++
            continue
        }
        
        # 替换或插入TOC
        if ($hasTOC) {
            # 替换现有TOC
            $content = $content -replace '(?s)<!-- TOC START -->.*?<!-- TOC END -->', $newTOC
            $action = "更新"
        }
        else {
            # 在第一个标题后插入TOC (改进版)
            # 尝试多种模式
            $inserted = $false
            
            # 模式1: 标题 + 简介 + 两个换行
            if ($content -match '(?s)(^#[^#].*?\n\n>.*?\n\n)(.*)$') {
                $header = $matches[1]
                $rest = $matches[2]
                $content = $header + $newTOC + "`n`n" + $rest
                $inserted = $true
            }
            # 模式2: 标题 + 两个换行
            elseif ($content -match '(?s)(^#[^#].*?\n\n)(.*)$') {
                $header = $matches[1]
                $rest = $matches[2]
                $content = $header + $newTOC + "`n`n" + $rest
                $inserted = $true
            }
            # 模式3: 标题 + 一个换行 + 内容
            elseif ($content -match '(?s)(^#[^#].*?\n)(.*)$') {
                $header = $matches[1]
                $rest = $matches[2]
                $content = $header + "`n" + $newTOC + "`n`n" + $rest
                $inserted = $true
            }
            
            if ($inserted) {
                $action = "插入"
            }
            else {
                Write-Host "  ⚠️  无法确定TOC插入位置" -ForegroundColor Yellow
                $skippedCount++
                continue
            }
        }
        
        if (-not $DryRun) {
            Set-Content -Path $file.FullName -Value $content -Encoding UTF8 -NoNewline
            Write-Host "  ✅ 已${action}TOC" -ForegroundColor Green
            $processedCount++
        }
        else {
            Write-Host "  🔍 [DryRun] 将${action}TOC" -ForegroundColor Cyan
            $processedCount++
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


# 术语统一脚本
# 功能: 统一文档中的术语使用
# 日期: 2025-11-11

param(
    [string]$DocsPath = "docs",
    [switch]$DryRun = $false,
    [switch]$Verbose = $false
)

Write-Host "📚 术语统一脚本" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

# 标准术语表（需要统一为左侧的标准术语）
$TerminologyMap = @{
    # Goroutine相关
    'goroutine' = 'Goroutine'
    '协程' = 'Goroutine'
    'go routine' = 'Goroutine'

    # Channel相关
    'channel' = 'Channel'
    '通道' = 'Channel'
    'chan' = 'Channel'

    # Context相关
    'context' = 'Context'
    '上下文' = 'Context'

    # 其他常见术语
    'mutex' = 'Mutex'
    '互斥锁' = 'Mutex'
    'waitgroup' = 'WaitGroup'
    '等待组' = 'WaitGroup'
}

# 统一术语
function Unify-Terminology {
    param([string]$Content)

    $newContent = $Content
    $changeCount = 0

    foreach ($term in $TerminologyMap.Keys) {
        $standardTerm = $TerminologyMap[$term]

        # 使用单词边界匹配，避免部分替换
        $pattern = '\b' + [regex]::Escape($term) + '\b'
        $matches = [regex]::Matches($newContent, $pattern, [System.Text.RegularExpressions.RegexOptions]::IgnoreCase)

        if ($matches.Count -gt 0) {
            $newContent = $newContent -replace $pattern, $standardTerm
            $changeCount += $matches.Count

            if ($Verbose) {
                Write-Host "  📝 '$term' → '$standardTerm' ($($matches.Count)处)" -ForegroundColor Yellow
            }
        }
    }

    return @{
        Content = $newContent
        ChangeCount = $changeCount
        Changed = $changeCount -gt 0
    }
}

# 主处理逻辑
$mdFiles = Get-ChildItem -Path $DocsPath -Filter "*.md" -Recurse
$totalFiles = $mdFiles.Count
$processedFiles = 0
$unifiedFiles = 0
$totalChanges = 0

Write-Host "📁 找到 $totalFiles 个Markdown文件" -ForegroundColor Cyan
Write-Host ""

foreach ($file in $mdFiles) {
    $processedFiles++

    try {
        $content = Get-Content $file.FullName -Raw -Encoding UTF8
        $result = Unify-Terminology -Content $content

        if ($result.Changed) {
            Write-Host "[$processedFiles/$totalFiles] 📚 $($file.Name)" -ForegroundColor Cyan
            Write-Host "  修改: $($result.ChangeCount) 处" -ForegroundColor Yellow

            if (-not $DryRun) {
                [System.IO.File]::WriteAllText($file.FullName, $result.Content, [System.Text.Encoding]::UTF8)
                Write-Host "  ✅ 已统一" -ForegroundColor Green
            } else {
                Write-Host "  ⚠ 预览模式（未实际修改）" -ForegroundColor Yellow
            }

            $unifiedFiles++
            $totalChanges += $result.ChangeCount
        } else {
            if ($Verbose) {
                Write-Host "[$processedFiles/$totalFiles] ✓ $($file.Name) - 无需修改" -ForegroundColor Gray
            }
        }
    } catch {
        Write-Host "  ❌ 错误: $_" -ForegroundColor Red
    }
}

# 总结
Write-Host ""
Write-Host "================================" -ForegroundColor Cyan
Write-Host "📊 统一总结" -ForegroundColor Cyan
Write-Host "  处理文件: $processedFiles/$totalFiles" -ForegroundColor White
Write-Host "  已统一: $unifiedFiles" -ForegroundColor White
Write-Host "  总修改: $totalChanges 处" -ForegroundColor White

if ($DryRun) {
    Write-Host "  ⚠ 这是预览模式，未实际修改文件" -ForegroundColor Yellow
    Write-Host "  使用 -DryRun:`$false 来实际统一" -ForegroundColor Yellow
} else {
    Write-Host "  ✅ 文件已更新" -ForegroundColor Green
}

Write-Host ""

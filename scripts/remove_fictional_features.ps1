# 🗑️ 删除虚构特性脚本
# 删除不存在的Go特性相关文件和目录

Write-Host "🗑️  开始删除虚构特性..." -ForegroundColor Cyan
Write-Host ""

$removedItems = 0
$failedItems = 0

# 定义要删除的文件和目录
$itemsToRemove = @(
    # WaitGroup.Go() - 不存在的特性
    "docs\02-Go语言现代化\14-Go-1.23并发和网络\01-WaitGroup-Go方法.md",
    "docs\02-Go语言现代化\14-Go-1.23并发和网络\examples\waitgroup_go",
    
    # testing/synctest - 不存在的包
    "docs\02-Go语言现代化\14-Go-1.23并发和网络\02-testing-synctest包.md",
    "docs\02-Go语言现代化\14-Go-1.23并发和网络\examples\synctest",
    
    # go.mod ignore - 不存在的指令
    "docs\02-Go语言现代化\13-Go-1.23工具链增强\02-go-mod-ignore指令.md",
    "docs\02-Go语言现代化\13-Go-1.23工具链增强\examples\go_mod_ignore",
    
    # Greentea GC - 虚构的垃圾收集器
    "docs\02-Go语言现代化\12-Go-1.23运行时优化\01-greentea-GC垃圾收集器.md"
    
    # Swiss Tables - 不是Go标准实现
    # (保留文档但标注为研究性质)
)

Write-Host "📋 要删除的项目:" -ForegroundColor Yellow
foreach ($item in $itemsToRemove) {
    Write-Host "  - $item" -ForegroundColor Gray
}
Write-Host ""

Write-Host "⚠️  警告: 此操作将永久删除以上文件/目录！" -ForegroundColor Red
$confirm = Read-Host "确认删除吗？(yes/no)"

if ($confirm -ne "yes") {
    Write-Host "❌ 操作已取消" -ForegroundColor Yellow
    exit
}

Write-Host ""
Write-Host "🔥 开始删除..." -ForegroundColor Yellow

foreach ($item in $itemsToRemove) {
    $fullPath = Join-Path $PWD $item
    
    if (Test-Path $fullPath) {
        try {
            if (Test-Path $fullPath -PathType Container) {
                Remove-Item $fullPath -Recurse -Force
                Write-Host "  ✅ 已删除目录: $item" -ForegroundColor Green
            } else {
                Remove-Item $fullPath -Force
                Write-Host "  ✅ 已删除文件: $item" -ForegroundColor Green
            }
            $removedItems++
        }
        catch {
            Write-Host "  ❌ 删除失败: $item - $($_.Exception.Message)" -ForegroundColor Red
            $failedItems++
        }
    } else {
        Write-Host "  ℹ️  不存在: $item" -ForegroundColor Gray
    }
}

Write-Host ""
Write-Host "📊 删除统计:" -ForegroundColor Cyan
Write-Host "  成功删除: $removedItems" -ForegroundColor Green
if ($failedItems -gt 0) {
    Write-Host "  删除失败: $failedItems" -ForegroundColor Red
}
Write-Host ""

if ($removedItems -gt 0) {
    Write-Host "✅ 虚构特性清理完成！" -ForegroundColor Green
} else {
    Write-Host "ℹ️  没有项目需要删除" -ForegroundColor Yellow
}

# 生成删除报告
$reportFile = "特性删除报告-$(Get-Date -Format 'yyyy-MM-dd-HHmm').txt"
$report = @"
虚构特性删除报告
生成时间: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')

成功删除: $removedItems
删除失败: $failedItems

删除的项目:
$($itemsToRemove -join "`n")
"@

$report | Out-File -FilePath $reportFile -Encoding UTF8
Write-Host "📄 详细报告已保存到: $reportFile" -ForegroundColor Cyan


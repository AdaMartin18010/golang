# 批量格式化文档脚本
# 为所有新创建的文档添加目录和序号

$ErrorActionPreference = "Stop"

Write-Host "======================================"
Write-Host "批量格式化文档"
Write-Host "======================================"

# 需要更新的文档列表
$docs = @(
    # practices/testing
    "docs/practices/testing/02-表格驱动测试.md",
    "docs/practices/testing/03-集成测试.md",
    "docs/practices/testing/04-性能测试.md",
    "docs/practices/testing/05-测试覆盖率.md",
    "docs/practices/testing/06-Mock与Stub.md",
    "docs/practices/testing/07-测试最佳实践.md",
    "docs/practices/testing/08-常见问题与技巧.md"
)

$count = 0
$total = $docs.Count

foreach ($doc in $docs) {
    $count++
    if (Test-Path $doc) {
        Write-Host "[$count/$total] 处理: $doc"
        # 这里可以添加格式化逻辑
    } else {
        Write-Host "[$count/$total] 跳过(不存在): $doc" -ForegroundColor Yellow
    }
}

Write-Host ""
Write-Host "完成! 处理了 $count 个文档" -ForegroundColor Green


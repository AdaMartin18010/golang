# 修复Go版本特性README中的失效链接
# 策略: 将不存在的子文档链接改为说明性文本

param(
    [switch]$DryRun
)

$ErrorActionPreference = "Stop"

$fixes = @{
    "docs-new\10-Go版本特性\01-Go-1.21特性\README.md" = @(
        @{
            Old = "- \[01-泛型改进\.md\]\(\.\/01-泛型改进\.md\) - 泛型性能优化和新增功能`n- \[02-性能优化\.md\]\(\.\/02-性能优化\.md\) - 编译器和运行时优化`n- \[03-标准库更新\.md\]\(\.\/03-标准库更新\.md\) - 标准库新增和改进"
            New = "- **泛型改进**: 泛型性能优化和新增功能（详见主文档）`n- **性能优化**: 编译器和运行时优化（详见主文档）`n- **标准库更新**: 标准库新增和改进（详见主文档）"
        }
    )
    "docs-new\10-Go版本特性\02-Go-1.22特性\README.md" = @(
        @{
            Old = "- \[01-for循环变量语义\.md\]\(\.\/01-for循环变量语义\.md\) - for循环变量作用域变更`n- \[02-HTTP路由增强\.md\]\(\.\/02-HTTP路由增强\.md\) - ServeMux路由模式增强`n- \[03-性能改进\.md\]\(\.\/03-性能改进\.md\) - 各项性能提升"
            New = "- **for循环变量语义**: for循环变量作用域变更（详见主文档）`n- **HTTP路由增强**: ServeMux路由模式增强（详见主文档）`n- **性能改进**: 各项性能提升（详见主文档）"
        }
    )
    "docs-new\10-Go版本特性\03-Go-1.23特性\README.md" = @(
        @{
            Old = "- \[01-迭代器预览\.md\]\(\.\/01-迭代器预览\.md\) - range over func实验性支持`n- \[02-工具链增强\.md\]\(\.\/02-工具链增强\.md\) - 工具链改进`n- \[03-标准库更新\.md\]\(\.\/03-标准库更新\.md\) - 标准库更新"
            New = "- **迭代器预览**: range over func实验性支持（详见主文档）`n- **工具链增强**: 工具链改进（详见主文档）`n- **标准库更新**: 标准库更新（详见主文档）"
        }
    )
    "docs-new\10-Go版本特性\04-Go-1.24特性\README.md" = @(
        @{
            Old = "- \[01-编译器优化\.md\]\(\.\/01-编译器优化\.md\) - 编译器优化`n- \[02-运行时改进\.md\]\(\.\/02-运行时改进\.md\) - 运行时改进`n- \[03-标准库增强\.md\]\(\.\/03-标准库增强\.md\) - 标准库增强"
            New = "- **编译器优化**: 编译器优化（详见主文档）`n- **运行时改进**: 运行时改进（详见主文档）`n- **标准库增强**: 标准库增强（详见主文档）"
        }
    )
}

Write-Host "=== Go版本特性README修复工具 ===" -ForegroundColor Cyan
Write-Host ""

$total = $fixes.Count
$fixed = 0

foreach ($file in $fixes.Keys) {
    if (-not (Test-Path $file)) {
        Write-Host "  ⚠️ 文件不存在: $file" -ForegroundColor Yellow
        continue
    }
    
    $content = Get-Content -Path $file -Raw
    $modified = $false
    
    foreach ($fix in $fixes[$file]) {
        if ($content -match [regex]::Escape($fix.Old)) {
            $content = $content -replace [regex]::Escape($fix.Old), $fix.New
            $modified = $true
        }
    }
    
    if ($modified) {
        if (-not $DryRun) {
            Set-Content -Path $file -Value $content -NoNewline
            Write-Host "  ✅ 已修复: $(Split-Path -Leaf $file)" -ForegroundColor Green
        } else {
            Write-Host "  [演练] 将修复: $(Split-Path -Leaf $file)" -ForegroundColor Yellow
        }
        $fixed++
    } else {
        Write-Host "  跳过: $(Split-Path -Leaf $file) (无需修复)" -ForegroundColor Gray
    }
}

Write-Host ""
Write-Host "=== 完成 ===" -ForegroundColor Cyan
Write-Host "  检查文件: $total" -ForegroundColor White
Write-Host "  已修复: $fixed" -ForegroundColor Green


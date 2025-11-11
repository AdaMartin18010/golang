# 分析文件夹结构脚本
# 功能: 分析当前文件夹结构，识别问题
# 日期: 2025-11-11

param(
    [string]$DocsPath = "docs"
)

Write-Host "📁 文件夹结构分析脚本" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

# 分析文件夹结构
function Analyze-FolderStructure {
    param([string]$Path)

    $issues = @()
    $folders = Get-ChildItem -Path $Path -Directory -Recurse

    foreach ($folder in $folders) {
        $name = $folder.Name
        $depth = ($folder.FullName -replace [regex]::Escape($Path), '').Split('\').Count - 1

        # 检查命名规范
        if ($name -match '^[0-9]+-') {
            $issues += @{
                Type = "编号前缀"
                Path = $folder.FullName
                Name = $name
                Depth = $depth
            }
        }

        if ($name -match '[A-Z]') {
            $issues += @{
                Type = "包含大写字母"
                Path = $folder.FullName
                Name = $name
                Depth = $depth
            }
        }

        if ($name -match '[\u4e00-\u9fa5]') {
            $issues += @{
                Type = "包含中文"
                Path = $folder.FullName
                Name = $name
                Depth = $depth
            }
        }

        # 检查层级深度
        if ($depth -gt 4) {
            $issues += @{
                Type = "层级过深"
                Path = $folder.FullName
                Name = $name
                Depth = $depth
            }
        }
    }

    return $issues
}

# 主处理逻辑
Write-Host "📁 分析文件夹结构..." -ForegroundColor Cyan
Write-Host ""

$issues = Analyze-FolderStructure -Path $DocsPath

# 按类型分组
$grouped = $issues | Group-Object -Property Type

Write-Host "📊 问题统计" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

foreach ($group in $grouped) {
    Write-Host "$($group.Name): $($group.Count) 个" -ForegroundColor Yellow

    # 显示前5个示例
    $examples = $group.Group | Select-Object -First 5
    foreach ($issue in $examples) {
        Write-Host "  - $($issue.Name) (深度: $($issue.Depth))" -ForegroundColor Gray
        Write-Host "    路径: $($issue.Path)" -ForegroundColor DarkGray
    }

    if ($group.Count -gt 5) {
        Write-Host "  ... 还有 $($group.Count - 5) 个" -ForegroundColor DarkGray
    }

    Write-Host ""
}

# 总结
Write-Host "================================" -ForegroundColor Cyan
Write-Host "📊 分析总结" -ForegroundColor Cyan
Write-Host "  总问题数: $($issues.Count)" -ForegroundColor White
Write-Host "  问题类型: $($grouped.Count)" -ForegroundColor White
Write-Host ""

# PowerShell Script: 文档迁移助手
# 版本: v1.0
# 日期: 2025-10-22

param(
    [string]$SourceDir = "docs",
    [string]$TargetDir = "docs-new",
    [string]$MappingFile = "migration-mapping.json",
    [switch]$DryRun,
    [switch]$Verbose
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  文档迁移助手" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

if ($DryRun) {
    Write-Host "🔍 DryRun模式 - 仅预览，不实际执行" -ForegroundColor Yellow
    Write-Host ""
}

# 统计变量
$stats = @{
    TotalFiles = 0
    Migrated = 0
    Skipped = 0
    Failed = 0
}

# 迁移映射表（如果没有JSON文件，使用默认映射）
$defaultMapping = @{
    "01-语言基础" = "01-语言基础"
    "02-Web开发" = "03-Web开发"
    "03-Go新特性" = "10-Go版本特性"
    "05-微服务" = "05-微服务架构"
    "06-云原生" = "06-云原生与容器"
    "07-性能优化" = "07-性能优化"
    "08-架构设计" = "08-架构设计"
    "09-工程实践" = "09-工程实践"
    "10-进阶专题" = "11-高级专题"
    "11-行业应用" = "12-行业应用"
    "12-参考资料" = "13-参考资料"
}

# 加载映射配置
if (Test-Path $MappingFile) {
    Write-Host "📋 加载映射配置: $MappingFile" -ForegroundColor Green
    $mapping = Get-Content $MappingFile | ConvertFrom-Json -AsHashtable
} else {
    Write-Host "⚠️  未找到映射文件，使用默认映射" -ForegroundColor Yellow
    $mapping = $defaultMapping
}

# 函数: 迁移单个文件
function Migrate-File {
    param(
        [string]$SourcePath,
        [string]$TargetPath
    )
    
    try {
        if ($DryRun) {
            Write-Host "  [DryRun] $SourcePath -> $TargetPath" -ForegroundColor Gray
            return $true
        }
        
        # 确保目标目录存在
        $targetFolder = Split-Path $TargetPath -Parent
        if (!(Test-Path $targetFolder)) {
            New-Item -ItemType Directory -Path $targetFolder -Force | Out-Null
        }
        
        # 复制文件
        Copy-Item -Path $SourcePath -Destination $TargetPath -Force
        
        if ($Verbose) {
            Write-Host "  ✓ $SourcePath -> $TargetPath" -ForegroundColor Green
        }
        
        return $true
    } catch {
        Write-Host "  ✗ 迁移失败: $SourcePath" -ForegroundColor Red
        Write-Host "    错误: $_" -ForegroundColor Red
        return $false
    }
}

# 主迁移逻辑
Write-Host "🚀 开始迁移..." -ForegroundColor Yellow
Write-Host ""

foreach ($oldModule in $mapping.Keys) {
    $newModule = $mapping[$oldModule]
    $sourcePath = Join-Path $SourceDir $oldModule
    $targetPath = Join-Path $TargetDir $newModule
    
    if (!(Test-Path $sourcePath)) {
        Write-Host "⚠️  跳过不存在的模块: $oldModule" -ForegroundColor Yellow
        continue
    }
    
    Write-Host "📂 迁移模块: $oldModule -> $newModule" -ForegroundColor Cyan
    
    # 获取所有md文件
    $files = Get-ChildItem -Path $sourcePath -Filter "*.md" -Recurse
    $stats.TotalFiles += $files.Count
    
    foreach ($file in $files) {
        $relativePath = $file.FullName.Replace($sourcePath, "").TrimStart('\')
        $targetFilePath = Join-Path $targetPath $relativePath
        
        $result = Migrate-File -SourcePath $file.FullName -TargetPath $targetFilePath
        
        if ($result) {
            $stats.Migrated++
        } else {
            $stats.Failed++
        }
    }
    
    Write-Host "  完成: 迁移 $($files.Count) 个文件" -ForegroundColor Green
    Write-Host ""
}

# 输出统计
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "📊 迁移统计" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "总文件数: $($stats.TotalFiles)" -ForegroundColor White
Write-Host "已迁移:   $($stats.Migrated)" -ForegroundColor Green
Write-Host "跳过:     $($stats.Skipped)" -ForegroundColor Yellow
Write-Host "失败:     $($stats.Failed)" -ForegroundColor Red
Write-Host ""

if ($DryRun) {
    Write-Host "💡 这是DryRun预览，实际执行请移除 -DryRun 参数" -ForegroundColor Yellow
} else {
    Write-Host "✅ 迁移完成!" -ForegroundColor Green
}

Write-Host "========================================" -ForegroundColor Cyan

# 返回统计结果
return $stats


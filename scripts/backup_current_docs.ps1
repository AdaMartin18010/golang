# PowerShell Script: 备份当前文档
# 用途: 在重构前备份现有docs目录
# 版本: v1.0
# 日期: 2025-10-22

param(
    [string]$SourceDir = "docs",
    [string]$BackupDir = "docs-backup-$(Get-Date -Format 'yyyyMMdd-HHmmss')",
    [switch]$Compress
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  文档备份脚本" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 检查源目录
if (!(Test-Path $SourceDir)) {
    Write-Host "❌ 错误: 源目录不存在: $SourceDir" -ForegroundColor Red
    exit 1
}

Write-Host "📂 源目录: $SourceDir" -ForegroundColor Cyan
Write-Host "📦 备份目录: $BackupDir" -ForegroundColor Cyan
Write-Host ""

# 创建备份
Write-Host "⏳ 正在备份..." -ForegroundColor Yellow
try {
    Copy-Item -Path $SourceDir -Destination $BackupDir -Recurse -Force
    Write-Host "✅ 备份完成!" -ForegroundColor Green
    
    # 统计
    $fileCount = (Get-ChildItem -Path $BackupDir -Recurse -File).Count
    $dirCount = (Get-ChildItem -Path $BackupDir -Recurse -Directory).Count
    $size = (Get-ChildItem -Path $BackupDir -Recurse -File | Measure-Object -Property Length -Sum).Sum / 1MB
    
    Write-Host ""
    Write-Host "📊 备份统计:" -ForegroundColor Cyan
    Write-Host "   文件数: $fileCount" -ForegroundColor White
    Write-Host "   目录数: $dirCount" -ForegroundColor White
    Write-Host "   大小: $([Math]::Round($size, 2)) MB" -ForegroundColor White
    
    # 压缩
    if ($Compress) {
        Write-Host ""
        Write-Host "⏳ 正在压缩备份..." -ForegroundColor Yellow
        $zipPath = "$BackupDir.zip"
        Compress-Archive -Path $BackupDir -DestinationPath $zipPath -Force
        Write-Host "✅ 压缩完成: $zipPath" -ForegroundColor Green
        
        # 删除未压缩的备份
        Remove-Item -Path $BackupDir -Recurse -Force
        Write-Host "🗑️  已删除未压缩的备份目录" -ForegroundColor Gray
    }
    
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "✅ 备份成功完成!" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Cyan
    
} catch {
    Write-Host "❌ 备份失败: $_" -ForegroundColor Red
    exit 1
}


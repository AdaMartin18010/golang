# Go 1.26 文档重命名脚本
# 将文件名中的 1.25.3 改为 1.26，2025 改为 2026

$docsPath = "docs"

# 获取所有需要重命名的文件
$files = Get-ChildItem -Path $docsPath -Recurse | Where-Object { 
    $_.Name -match "1\.25\.3|1\.25\.1|Go1\.25\.3" -and !$_.PSIsContainer
}

Write-Host "找到 $($files.Count) 个需要重命名的文件" -ForegroundColor Yellow

$renameMap = @()

foreach ($file in $files) {
    $oldName = $file.Name
    $newName = $oldName -replace "1\.25\.3", "1.26" -replace "1\.25\.1", "1.26" -replace "Go1\.25\.3", "Go1.26"
    
    # 如果是版本文档，同时更新年份
    if ($newName -match "00-Go-1\.26" -and $newName -match "2025") {
        $newName = $newName -replace "2025", "2026"
    }
    
    if ($oldName -ne $newName) {
        $renameMap += [PSCustomObject]@{
            OldPath = $file.FullName
            NewPath = Join-Path $file.DirectoryName $newName
            OldName = $oldName
            NewName = $newName
        }
    }
}

# 显示预览
Write-Host "`n=== 重命名预览 ===" -ForegroundColor Cyan
$renameMap | ForEach-Object {
    Write-Host "$($_.OldName)" -ForegroundColor Red -NoNewline
    Write-Host " -> " -NoNewline
    Write-Host "$($_.NewName)" -ForegroundColor Green
}

# 执行重命名
Write-Host "`n=== 执行重命名 ===" -ForegroundColor Cyan
$count = 0
foreach ($item in $renameMap) {
    try {
        Rename-Item -Path $item.OldPath -NewName $item.NewName -Force
        Write-Host "✓ $($item.OldName) -> $($item.NewName)" -ForegroundColor Green
        $count++
    } catch {
        Write-Host "✗ 失败: $($item.OldName) - $_" -ForegroundColor Red
    }
}

Write-Host "`n完成! 成功重命名 $count 个文件" -ForegroundColor Yellow

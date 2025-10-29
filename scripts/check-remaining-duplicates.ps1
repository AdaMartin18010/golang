# 检查仍然有重复版本信息的文件

$ErrorActionPreference = "Stop"

Write-Host "🔍 检查仍有重复版本信息的文件...`n" -ForegroundColor Cyan

$mdFiles = Get-ChildItem -Path "docs" -Filter "*.md" -Recurse -File

$duplicates = @()

foreach ($file in $mdFiles) {
    try {
        $content = Get-Content -Path $file.FullName -Raw -Encoding UTF8
        if (-not $content) {
            continue
        }
        
        $matches = [regex]::Matches($content, '\*\*版本\*\*:')
        $count = $matches.Count
        
        if ($count -gt 1) {
            $duplicates += [PSCustomObject]@{
                File = $file.FullName.Replace((Get-Location).Path + "\", "")
                Count = $count
            }
        }
    }
    catch {
        Write-Host "错误: $($file.Name) - $($_.Exception.Message)" -ForegroundColor Red
    }
}

if ($duplicates.Count -gt 0) {
    Write-Host "❌ 发现 $($duplicates.Count) 个文件仍有重复版本信息:`n" -ForegroundColor Yellow
    $duplicates | Sort-Object Count -Descending | Format-Table -AutoSize
} else {
    Write-Host "✅ 没有发现重复的版本信息!" -ForegroundColor Green
}


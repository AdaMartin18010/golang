# 修复所有重复的版本信息 - 终极版本
# 处理所有可能的格式

$ErrorActionPreference = "Stop"
$filesFixed = 0

Write-Host "🔍 全面修复重复的版本信息...`n" -ForegroundColor Cyan

$mdFiles = Get-ChildItem -Path "docs" -Filter "*.md" -Recurse -File

foreach ($file in $mdFiles) {
    try {
        $content = Get-Content -Path $file.FullName -Raw -Encoding UTF8
        if (-not $content) {
            continue
        }
        
        $originalContent = $content
        $changed = $false
        
        # 计数版本信息出现次数
        $matches = [regex]::Matches($content, '\*\*版本\*\*:')
        $count = $matches.Count
        
        if ($count -le 1) {
            continue
        }
        
        # 找到所有版本块的位置
        $versionBlocks = @()
        $lines = $content -split "`r?`n"
        
        for ($i = 0; $i -lt $lines.Count; $i++) {
            if ($lines[$i] -match '^\*\*版本\*\*:') {
                # 找到版本块的结束位置（下一个---或下一个非空行）
                $endIndex = $i + 1
                while ($endIndex -lt $lines.Count) {
                    if ($lines[$endIndex] -eq '---') {
                        $endIndex++
                        break
                    }
                    if ($lines[$endIndex] -notmatch '^\*\*' -and $lines[$endIndex] -notmatch '^\s*$') {
                        break
                    }
                    $endIndex++
                }
                
                $versionBlocks += @{
                    Start = $i
                    End = $endIndex
                    Lines = $lines[$i..($endIndex-1)]
                }
                
                $i = $endIndex - 1
            }
        }
        
        if ($versionBlocks.Count -gt 1) {
            Write-Host "  处理: $($file.Name) (发现 $count 个版本块)" -ForegroundColor Yellow
            
            # 只保留第一个版本块
            $result = New-Object System.Collections.ArrayList
            $keepBlock = $versionBlocks[0]
            $skipRanges = @()
            
            # 收集要跳过的行范围
            for ($i = 1; $i -lt $versionBlocks.Count; $i++) {
                $skipRanges += @{
                    Start = $versionBlocks[$i].Start
                    End = $versionBlocks[$i].End
                }
            }
            
            # 重建内容
            for ($i = 0; $i -lt $lines.Count; $i++) {
                $skip = $false
                
                foreach ($range in $skipRanges) {
                    if ($i -ge $range.Start -and $i -lt $range.End) {
                        $skip = $true
                        break
                    }
                }
                
                if (-not $skip) {
                    [void]$result.Add($lines[$i])
                }
            }
            
            # 清理多余的连续分隔符和空行
            $cleanedResult = New-Object System.Collections.ArrayList
            $prevWasSeparator = $false
            $prevWasEmpty = $false
            
            for ($i = 0; $i -lt $result.Count; $i++) {
                $line = $result[$i]
                
                # 处理连续的---
                if ($line -eq '---') {
                    if ($prevWasSeparator) {
                        continue  # 跳过重复的分隔符
                    }
                    $prevWasSeparator = $true
                    $prevWasEmpty = $false
                    [void]$cleanedResult.Add($line)
                    continue
                }
                
                # 处理连续的空行（限制为最多2个）
                if ($line -match '^\s*$') {
                    if ($prevWasEmpty) {
                        # 已经有一个空行了,跳过
                        continue
                    }
                    $prevWasEmpty = $true
                    $prevWasSeparator = $false
                    [void]$cleanedResult.Add($line)
                    continue
                }
                
                # 普通行
                $prevWasSeparator = $false
                $prevWasEmpty = $false
                [void]$cleanedResult.Add($line)
            }
            
            $newContent = $cleanedResult -join "`n"
            
            if ($originalContent -ne $newContent) {
                Set-Content -Path $file.FullName -Value $newContent -Encoding UTF8 -NoNewline
                $filesFixed++
                $changed = $true
                Write-Host "  ✅ 已修复: $($file.Name)" -ForegroundColor Green
            }
        }
    }
    catch {
        Write-Host "  ❌ 错误: $($file.Name) - $($_.Exception.Message)" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "=" * 60 -ForegroundColor Cyan
Write-Host "✨ 全面修复完成!" -ForegroundColor Green
Write-Host "📊 统计:" -ForegroundColor Cyan
Write-Host "  - 修复文件数: $filesFixed" -ForegroundColor Green
Write-Host "=" * 60 -ForegroundColor Cyan


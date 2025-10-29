# 修复文档中重复的版本信息 - V2
# 更精确的处理方式

$ErrorActionPreference = "Stop"
$filesFixed = 0
$totalIssues = 0

Write-Host "🔍 扫描docs目录中的所有.md文件..." -ForegroundColor Cyan

# 获取所有Markdown文件
$mdFiles = Get-ChildItem -Path "docs" -Filter "*.md" -Recurse -File | Where-Object { 
    $_.FullName -notlike "*\node_modules\*" 
}

Write-Host "📝 找到 $($mdFiles.Count) 个Markdown文件`n" -ForegroundColor Green

foreach ($file in $mdFiles) {
    try {
        $content = Get-Content -Path $file.FullName -Raw -Encoding UTF8
        if (-not $content) {
            continue
        }
        
        $originalContent = $content
        $fileIssues = 0
        $changed = $false
        
        # 将内容按行处理
        $lines = $content -split "`r?`n"
        $result = New-Object System.Collections.ArrayList
        
        $i = 0
        $skipUntil = -1
        
        while ($i -lt $lines.Count) {
            # 如果在跳过范围内，直接跳过
            if ($i -lt $skipUntil) {
                $i++
                continue
            }
            
            $line = $lines[$i]
            
            # 检测版本信息块的开始
            if ($line -match '^\*\*版本\*\*:') {
                # 收集完整的版本信息块
                $blockStart = $i
                $blockLines = @()
                $j = $i
                
                # 读取版本块（版本+更新日期+适用于+分隔符）
                while ($j -lt $lines.Count) {
                    $blockLines += $lines[$j]
                    
                    # 如果遇到分隔符，检查后面是否还有另一个版本块
                    if ($lines[$j] -eq '---') {
                        # 跳过空行
                        $k = $j + 1
                        while ($k -lt $lines.Count -and $lines[$k] -match '^\s*$') {
                            $k++
                        }
                        
                        # 检查是否是重复的版本块
                        if ($k -lt $lines.Count -and $lines[$k] -match '^\*\*版本\*\*:') {
                            # 找到重复！跳过重复的版本块
                            $dupStart = $k
                            $fileIssues++
                            $changed = $true
                            
                            Write-Host "  发现重复版本块: $($file.Name) (行 $($blockStart+1) 和 行 $($dupStart+1))" -ForegroundColor Yellow
                            
                            # 跳过重复的版本块
                            while ($k -lt $lines.Count) {
                                if ($lines[$k] -eq '---') {
                                    $k++
                                    # 跳过分隔符后的空行
                                    while ($k -lt $lines.Count -and $lines[$k] -match '^\s*$') {
                                        $k++
                                    }
                                    break
                                }
                                $k++
                            }
                            
                            # 同时检查是否有多余的分隔符
                            if ($k -lt $lines.Count -and $lines[$k] -eq '---') {
                                $k++
                                $fileIssues++
                                Write-Host "  移除多余分隔符: $($file.Name)" -ForegroundColor Yellow
                            }
                            
                            $skipUntil = $k
                        }
                        
                        break
                    }
                    
                    $j++
                }
                
                # 添加第一个版本块
                foreach ($bl in $blockLines) {
                    [void]$result.Add($bl)
                }
                
                $i = $j + 1
                continue
            }
            
            # 普通行，直接添加
            [void]$result.Add($line)
            $i++
        }
        
        if ($changed) {
            # 重建内容
            $newContent = $result -join "`n"
            
            # 保存文件
            Set-Content -Path $file.FullName -Value $newContent -Encoding UTF8 -NoNewline
            $filesFixed++
            Write-Host "✅ 已修复: $($file.Name) ($fileIssues 个问题)" -ForegroundColor Green
        }
    }
    catch {
        Write-Host "❌ 处理失败: $($file.Name) - $($_.Exception.Message)" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "=" * 60 -ForegroundColor Cyan
Write-Host "✨ 修复完成!" -ForegroundColor Green
Write-Host "📊 统计:" -ForegroundColor Cyan
Write-Host "  - 扫描文件: $($mdFiles.Count)" -ForegroundColor White
Write-Host "  - 修复文件: $filesFixed" -ForegroundColor Green
Write-Host "  - 总共问题: $totalIssues" -ForegroundColor Green
Write-Host "=" * 60 -ForegroundColor Cyan


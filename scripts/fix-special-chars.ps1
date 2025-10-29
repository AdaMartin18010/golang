# 修复特殊字符anchor问题
# 处理 go.mod, Q:, 等特殊情况

param(
    [string]$Path = "docs",
    [switch]$DryRun,
    [switch]$Verbose
)

$ErrorActionPreference = "Stop"
$stats = @{
    FilesProcessed = 0
    LinksFixed = 0
    Errors = 0
}

Write-Host "🔧 修复特殊字符anchor...`n" -ForegroundColor Cyan

# 处理文件
$files = Get-ChildItem -Path $Path -Recurse -Filter "*.md" -File

foreach ($file in $files) {
    try {
        $stats.FilesProcessed++
        $content = Get-Content $file.FullName -Raw -Encoding UTF8
        
        if ([string]::IsNullOrWhiteSpace($content)) {
            continue
        }
        
        $originalContent = $content
        $fixCount = 0
        
        # 修复1: go.mod 文件 -> gomod-文件
        $pattern1 = '\(#go\.mod-'
        $replacement1 = '(#gomod-'
        if ($content -match $pattern1) {
            $content = $content -replace $pattern1, $replacement1
            $count = ([regex]::Matches($originalContent, $pattern1)).Count
            $fixCount += $count
            if ($Verbose) {
                Write-Host "  修复 go.mod: $count 个" -ForegroundColor Gray
            }
        }
        
        # 修复2: go.sum 文件 -> gosum-文件
        $pattern2 = '\(#go\.sum-'
        $replacement2 = '(#gosum-'
        if ($content -match $pattern2) {
            $content = $content -replace $pattern2, $replacement2
            $count = ([regex]::Matches($originalContent, $pattern2)).Count
            $fixCount += $count
            if ($Verbose) {
                Write-Host "  修复 go.sum: $count 个" -ForegroundColor Gray
            }
        }
        
        # 修复3: go.work 文件 -> gowork-文件
        $pattern3 = '\(#go\.work-'
        $replacement3 = '(#gowork-'
        if ($content -match $pattern3) {
            $content = $content -replace $pattern3, $replacement3
            $count = ([regex]::Matches($originalContent, $pattern3)).Count
            $fixCount += $count
            if ($Verbose) {
                Write-Host "  修复 go.work: $count 个" -ForegroundColor Gray
            }
        }
        
        # 修复4: Q: xxx -> q-xxx (移除冒号和空格)
        # (#q:-xxx) -> (#q-xxx)
        $pattern4 = '\(#q:-'
        $replacement4 = '(#q-'
        if ($content -match $pattern4) {
            $content = $content -replace $pattern4, $replacement4
            $count = ([regex]::Matches($originalContent, $pattern4)).Count
            $fixCount += $count
            if ($Verbose) {
                Write-Host "  修复 Q:: $count 个" -ForegroundColor Gray
            }
        }
        
        # 修复5: // indirect (反引号中的内容，点号会被保留但反引号会被移除)
        # (#xxx-`//indirect`-) -> (#xxx-indirect-)
        $pattern5 = '`//\s*indirect`'
        $pattern5_in_anchor = '\(#[^)]*`//\s*indirect`[^)]*\)'
        if ($content -match $pattern5_in_anchor) {
            # 找到包含这个模式的所有anchor
            $anchors = [regex]::Matches($content, $pattern5_in_anchor)
            foreach ($anchor in $anchors) {
                $oldAnchor = $anchor.Value
                $newAnchor = $oldAnchor -replace '`//\s*indirect`', 'indirect'
                $content = $content -replace [regex]::Escape($oldAnchor), $newAnchor
                $fixCount++
            }
            if ($Verbose) {
                Write-Host "  修复 //indirect: $($anchors.Count) 个" -ForegroundColor Gray
            }
        }
        
        # 修复6: 括号中的内容 (Module) -> module
        # (#模块-module) 已经是正确的，但如果是 (#模块-\(module\)) 就要修复
        $pattern6 = '\(#[^)]*\\\([^)]*\\\)[^)]*\)'
        if ($content -match $pattern6) {
            $anchors = [regex]::Matches($content, $pattern6)
            foreach ($anchor in $anchors) {
                $oldAnchor = $anchor.Value
                $newAnchor = $oldAnchor -replace '\\\(', '' -replace '\\\)', ''
                $content = $content -replace [regex]::Escape($oldAnchor), $newAnchor
                $fixCount++
            }
            if ($Verbose) {
                Write-Host "  修复括号: $($anchors.Count) 个" -ForegroundColor Gray
            }
        }
        
        # 修复7: vs -> vs (已经正确，但空格处理)
        # (#gopath-vs-go-modules) 已经正确
        
        # 修复8: Go 1.18+ -> go-118
        # (#使用-workspace-go-1.18+) -> (#使用-workspace-go-118)
        $pattern8 = '\(#([^)]*)-(\d+)\.(\d+)\+'
        $replacement8 = '(#$1-$2$3'
        if ($content -match $pattern8) {
            $content = $content -replace $pattern8, $replacement8
            $count = ([regex]::Matches($originalContent, $pattern8)).Count
            $fixCount += $count
            if ($Verbose) {
                Write-Host "  修复版本号+: $count 个" -ForegroundColor Gray
            }
        }
        
        # 修复9: (Go 1.18+) 末尾的情况
        # (#xxx-go-1.18+\)) -> (#xxx-go-118)
        $pattern9 = '-(\d+)\.(\d+)\+\)'
        $replacement9 = '-$1$2)'
        if ($content -match $pattern9) {
            $content = $content -replace $pattern9, $replacement9
            $count = ([regex]::Matches($originalContent, $pattern9)).Count
            $fixCount += $count
            if ($Verbose) {
                Write-Host "  修复版本号+末尾: $count 个" -ForegroundColor Gray
            }
        }
        
        # 修复10: HTTP/2, HTTP/3 -> http2, http3
        $pattern10 = '\(#([^)]*)http/(\d)'
        $replacement10 = '(#$1http$2'
        if ($content -match $pattern10) {
            $content = $content -replace $pattern10, $replacement10
            $count = ([regex]::Matches($originalContent, $pattern10)).Count
            $fixCount += $count
            if ($Verbose) {
                Write-Host "  修复 HTTP/x: $count 个" -ForegroundColor Gray
            }
        }
        
        # 修复11: CI/CD -> cicd
        $pattern11 = '\(#([^)]*)ci/cd'
        $replacement11 = '(#$1cicd'
        if ($content -match $pattern11) {
            $content = $content -replace $pattern11, $replacement11
            $count = ([regex]::Matches($originalContent, $pattern11)).Count
            $fixCount += $count
            if ($Verbose) {
                Write-Host "  修复 CI/CD: $count 个" -ForegroundColor Gray
            }
        }
        
        # 修复12: I/O -> io
        $pattern12 = '\(#([^)]*)i/o'
        $replacement12 = '(#$1io'
        if ($content -match $pattern12) {
            $content = $content -replace $pattern12, $replacement12
            $count = ([regex]::Matches($originalContent, $pattern12)).Count
            $fixCount += $count
            if ($Verbose) {
                Write-Host "  修复 I/O: $count 个" -ForegroundColor Gray
            }
        }
        
        if ($fixCount -gt 0) {
            Write-Host "✓ $($file.FullName -replace [regex]::Escape($PWD), '.')" -ForegroundColor Green
            Write-Host "  修复: $fixCount 个链接" -ForegroundColor Gray
            
            $stats.LinksFixed += $fixCount
            
            if (-not $DryRun) {
                Set-Content -Path $file.FullName -Value $content -Encoding UTF8 -NoNewline
            }
        }
        
    } catch {
        $stats.Errors++
        Write-Host "✗ $($file.Name): $($_.Exception.Message)" -ForegroundColor Red
    }
}

# 显示统计
Write-Host "`n" + ("="*70) -ForegroundColor Cyan
Write-Host "📊 修复统计:" -ForegroundColor Cyan
Write-Host "  处理文件: $($stats.FilesProcessed)" -ForegroundColor White
Write-Host "  修复链接: $($stats.LinksFixed)" -ForegroundColor Green
Write-Host "  错误数量: $($stats.Errors)" -ForegroundColor $(if($stats.Errors -eq 0){'Green'}else{'Red'})
Write-Host ("="*70) -ForegroundColor Cyan

if ($DryRun) {
    Write-Host "`n⚠️  试运行模式 - 未实际修改文件" -ForegroundColor Yellow
} else {
    Write-Host "`n✅ 修复完成！" -ForegroundColor Green
}


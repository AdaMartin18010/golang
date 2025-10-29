# 检查未修复的链接问题

$ErrorActionPreference = "Stop"

$problematicFiles = @(
    "docs\advanced\reference\99-完整术语表与索引.md",
    "docs\fundamentals\README.md",
    "docs\advanced\concurrency\07-并发模式实战深度指南.md",
    "docs\reference\DOCUMENT_STANDARD.md",
    "docs\reference\versions\01-Go-1.21特性\00-知识图谱.md",
    "docs\reference\versions\02-Go-1.22特性\00-知识图谱.md",
    "docs\reference\versions\03-Go-1.23特性\00-知识图谱.md",
    "docs\reference\versions\05-实践应用\00-知识图谱.md",
    "docs\projects\templates\00-项目模板说明.md",
    "docs\projects\templates\01-项目结构模板.md",
    "docs\projects\templates\03-Web应用模板.md",
    "docs\projects\templates\04-CLI工具模板.md",
    "docs\projects\templates\05-库项目模板.md",
    "docs\projects\templates\06-快速开始指南.md",
    "docs\reference\GO-ECOSYSTEM-2025.md"
)

foreach ($file in $problematicFiles) {
    if (Test-Path $file) {
        Write-Host "`n=====================================" -ForegroundColor Cyan
        Write-Host "📄 $file" -ForegroundColor Yellow
        Write-Host "=====================================" -ForegroundColor Cyan
        
        $content = Get-Content -Path $file -Raw -Encoding UTF8
        $lines = $content -split "`r?`n"
        
        # 查找目录部分
        $inTOC = $false
        $tocStart = -1
        $tocEnd = -1
        
        for ($i = 0; $i -lt $lines.Count; $i++) {
            if ($lines[$i] -match '^##\s+(📋\s*)?目录$') {
                $inTOC = $true
                $tocStart = $i
            }
            elseif ($inTOC -and $lines[$i] -match '^##\s+') {
                $tocEnd = $i
                break
            }
        }
        
        if ($tocStart -ge 0) {
            Write-Host "目录部分 (行 $($tocStart+1) 到 $($tocEnd+1)):" -ForegroundColor Green
            
            $endIdx = if ($tocEnd -gt 0) { $tocEnd } else { [Math]::Min($tocStart + 50, $lines.Count) }
            
            for ($i = $tocStart; $i -lt $endIdx -and $i -lt $lines.Count; $i++) {
                $line = $lines[$i]
                if ($line -match '\[([^\]]+)\]\(#([^\)]+)\)') {
                    Write-Host "  $($i+1): $line" -ForegroundColor Gray
                }
            }
        }
        else {
            Write-Host "⚠️  未找到目录部分" -ForegroundColor Red
        }
    }
    else {
        Write-Host "❌ 文件不存在: $file" -ForegroundColor Red
    }
}


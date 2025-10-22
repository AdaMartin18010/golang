# PowerShell Script: 调试助手
# 版本: v1.0
# 日期: 2025-10-22
# 功能: 诊断问题、收集信息、生成调试报告

param(
    [ValidateSet("env", "files", "links", "quality", "all")]
    [string]$Check = "all",
    [string]$TargetDir = "docs",
    [string]$OutputFile = "debug-report-$(Get-Date -Format 'yyyyMMdd-HHmmss').md"
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  调试助手" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$report = @()
$report += "# 🔍 调试报告"
$report += ""
$report += "> **生成时间**: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')"
$report += "> **检查类型**: $Check"
$report += "> **目标目录**: $TargetDir"
$report += ""
$report += "---"
$report += ""

# 函数: 检查环境
function Check-Environment {
    Write-Host "🔍 检查环境..." -ForegroundColor Yellow
    
    $envReport = @()
    $envReport += "## 🖥️ 环境信息"
    $envReport += ""
    
    # PowerShell版本
    $psVersion = $PSVersionTable.PSVersion
    $envReport += "### PowerShell"
    $envReport += "- **版本**: $($psVersion.Major).$($psVersion.Minor).$($psVersion.Build)"
    $envReport += "- **Edition**: $($PSVersionTable.PSEdition)"
    $envReport += "- **状态**: $(if ($psVersion.Major -ge 5) { '✅ 满足要求' } else { '❌ 需要5.0+' })"
    $envReport += ""
    
    # 执行策略
    $execPolicy = Get-ExecutionPolicy
    $envReport += "### 执行策略"
    $envReport += "- **当前策略**: $execPolicy"
    $envReport += "- **状态**: $(if ($execPolicy -ne 'Restricted') { '✅ 可执行脚本' } else { '❌ 受限制' })"
    $envReport += ""
    
    # Git
    try {
        $gitVersion = git --version 2>$null
        $envReport += "### Git"
        $envReport += "- **版本**: $gitVersion"
        $envReport += "- **状态**: ✅ 已安装"
    } catch {
        $envReport += "### Git"
        $envReport += "- **状态**: ❌ 未安装"
    }
    $envReport += ""
    
    # Go
    try {
        $goVersion = go version 2>$null
        $envReport += "### Go"
        $envReport += "- **版本**: $goVersion"
        $envReport += "- **状态**: ✅ 已安装"
    } catch {
        $envReport += "### Go"
        $envReport += "- **状态**: ⚠️ 未安装（可选）"
    }
    $envReport += ""
    
    # 磁盘空间
    $drive = (Get-Location).Drive
    $driveInfo = Get-PSDrive $drive.Name
    $freeGB = [math]::Round($driveInfo.Free / 1GB, 2)
    $usedGB = [math]::Round($driveInfo.Used / 1GB, 2)
    $totalGB = [math]::Round(($driveInfo.Used + $driveInfo.Free) / 1GB, 2)
    
    $envReport += "### 磁盘空间 ($($drive.Name):)"
    $envReport += "- **总空间**: $totalGB GB"
    $envReport += "- **已使用**: $usedGB GB"
    $envReport += "- **可用**: $freeGB GB"
    $envReport += "- **状态**: $(if ($freeGB -gt 10) { '✅ 充足' } elseif ($freeGB -gt 5) { '⚠️ 注意' } else { '❌ 不足' })"
    $envReport += ""
    
    # 工作目录
    $currentPath = Get-Location
    $envReport += "### 工作目录"
    $envReport += "- **路径**: ``$currentPath``"
    $envReport += "- **长度**: $($currentPath.Path.Length) 字符"
    $envReport += "- **状态**: $(if ($currentPath.Path.Length -lt 200) { '✅ 正常' } else { '⚠️ 路径过长' })"
    $envReport += ""
    
    return $envReport
}

# 函数: 检查文件
function Check-Files {
    Write-Host "🔍 检查文件..." -ForegroundColor Yellow
    
    $filesReport = @()
    $filesReport += "## 📂 文件统计"
    $filesReport += ""
    
    if (!(Test-Path $TargetDir)) {
        $filesReport += "❌ 目标目录不存在: $TargetDir"
        return $filesReport
    }
    
    # Markdown文件
    $mdFiles = Get-ChildItem -Path $TargetDir -Filter "*.md" -Recurse
    $filesReport += "### Markdown文件"
    $filesReport += "- **总数**: $($mdFiles.Count)"
    $filesReport += ""
    
    # 按大小分类
    $smallFiles = $mdFiles | Where-Object { $_.Length -lt 10KB }
    $mediumFiles = $mdFiles | Where-Object { $_.Length -ge 10KB -and $_.Length -lt 100KB }
    $largeFiles = $mdFiles | Where-Object { $_.Length -ge 100KB }
    
    $filesReport += "### 文件大小分布"
    $filesReport += "| 大小 | 数量 | 百分比 |"
    $filesReport += "|------|------|--------|"
    $filesReport += "| <10KB | $($smallFiles.Count) | $([math]::Round(($smallFiles.Count / $mdFiles.Count) * 100, 1))% |"
    $filesReport += "| 10-100KB | $($mediumFiles.Count) | $([math]::Round(($mediumFiles.Count / $mdFiles.Count) * 100, 1))% |"
    $filesReport += "| >100KB | $($largeFiles.Count) | $([math]::Round(($largeFiles.Count / $mdFiles.Count) * 100, 1))% |"
    $filesReport += ""
    
    # 大文件列表
    if ($largeFiles.Count -gt 0) {
        $filesReport += "### ⚠️ 大文件 (>100KB)"
        $largeFiles | Sort-Object Length -Descending | Select-Object -First 10 | ForEach-Object {
            $sizeKB = [math]::Round($_.Length / 1KB, 1)
            $relativePath = $_.FullName.Replace((Get-Location).Path, "").TrimStart('\')
            $filesReport += "- ``$relativePath`` ($sizeKB KB)"
        }
        $filesReport += ""
    }
    
    # README文件
    $readmeFiles = Get-ChildItem -Path $TargetDir -Filter "README.md" -Recurse
    $filesReport += "### README文件"
    $filesReport += "- **总数**: $($readmeFiles.Count)"
    $filesReport += "- **状态**: $(if ($readmeFiles.Count -ge 10) { '✅ 充足' } else { '⚠️ 可能缺失' })"
    $filesReport += ""
    
    # 空文件
    $emptyFiles = $mdFiles | Where-Object { $_.Length -eq 0 }
    if ($emptyFiles.Count -gt 0) {
        $filesReport += "### ❌ 空文件"
        $emptyFiles | ForEach-Object {
            $relativePath = $_.FullName.Replace((Get-Location).Path, "").TrimStart('\')
            $filesReport += "- ``$relativePath``"
        }
        $filesReport += ""
    }
    
    return $filesReport
}

# 函数: 检查链接
function Check-Links {
    Write-Host "🔍 检查链接..." -ForegroundColor Yellow
    
    $linksReport = @()
    $linksReport += "## 🔗 链接检查"
    $linksReport += ""
    
    if (!(Test-Path $TargetDir)) {
        $linksReport += "❌ 目标目录不存在: $TargetDir"
        return $linksReport
    }
    
    $totalLinks = 0
    $internalLinks = 0
    $externalLinks = 0
    $brokenLinks = 0
    $anchorLinks = 0
    
    $mdFiles = Get-ChildItem -Path $TargetDir -Filter "*.md" -Recurse
    
    foreach ($file in $mdFiles) {
        $content = Get-Content $file.FullName -Raw -Encoding UTF8
        $links = [regex]::Matches($content, '\[([^\]]+)\]\(([^\)]+)\)')
        
        foreach ($link in $links) {
            $totalLinks++
            $linkUrl = $link.Groups[2].Value
            
            if ($linkUrl -match "^#") {
                $anchorLinks++
            } elseif ($linkUrl -match "^https?://") {
                $externalLinks++
            } elseif ($linkUrl -match "^\.\.?/") {
                $internalLinks++
                
                # 检查内部链接
                $targetUrl = $linkUrl -replace '#.*$', ''
                $targetPath = Join-Path (Split-Path $file.FullName) $targetUrl
                $targetPath = [System.IO.Path]::GetFullPath($targetPath)
                
                if (!(Test-Path $targetPath)) {
                    $brokenLinks++
                }
            }
        }
    }
    
    $linksReport += "### 链接统计"
    $linksReport += "| 类型 | 数量 | 百分比 |"
    $linksReport += "|------|------|--------|"
    $linksReport += "| 总链接 | $totalLinks | 100% |"
    $linksReport += "| 内部链接 | $internalLinks | $([math]::Round(($internalLinks / $totalLinks) * 100, 1))% |"
    $linksReport += "| 外部链接 | $externalLinks | $([math]::Round(($externalLinks / $totalLinks) * 100, 1))% |"
    $linksReport += "| 锚点链接 | $anchorLinks | $([math]::Round(($anchorLinks / $totalLinks) * 100, 1))% |"
    $linksReport += "| 失效链接 | $brokenLinks | $([math]::Round(($brokenLinks / $totalLinks) * 100, 1))% |"
    $linksReport += ""
    
    $linksReport += "### 链接健康度"
    $healthRate = if ($totalLinks -gt 0) {
        [math]::Round((($totalLinks - $brokenLinks) / $totalLinks) * 100, 2)
    } else { 100 }
    
    $linksReport += "- **健康率**: $healthRate%"
    $linksReport += "- **状态**: $(if ($healthRate -ge 98) { '✅ 优秀' } elseif ($healthRate -ge 90) { '⚠️ 良好' } elseif ($healthRate -ge 80) { '⚠️ 一般' } else { '❌ 需要修复' })"
    $linksReport += ""
    
    if ($brokenLinks -gt 0) {
        $linksReport += "> 💡 **建议**: 运行 ``.\scripts\fix_links.ps1`` 修复失效链接"
        $linksReport += ""
    }
    
    return $linksReport
}

# 函数: 检查质量
function Check-Quality {
    Write-Host "🔍 检查质量..." -ForegroundColor Yellow
    
    $qualityReport = @()
    $qualityReport += "## ✅ 质量检查"
    $qualityReport += ""
    
    if (!(Test-Path $TargetDir)) {
        $qualityReport += "❌ 目标目录不存在: $TargetDir"
        return $qualityReport
    }
    
    $mdFiles = Get-ChildItem -Path $TargetDir -Filter "*.md" -Recurse
    
    $withMeta = 0
    $withCode = 0
    $withTitle = 0
    $withImages = 0
    
    foreach ($file in $mdFiles) {
        $content = Get-Content $file.FullName -Raw -Encoding UTF8
        
        # 检查元信息
        if ($content -match "维护者:|最后更新:|创建日期:") {
            $withMeta++
        }
        
        # 检查代码块
        if ($content -match "```") {
            $withCode++
        }
        
        # 检查标题
        if ($content -match "^# ") {
            $withTitle++
        }
        
        # 检查图片
        if ($content -match "!\[.*\]\(.*\)") {
            $withImages++
        }
    }
    
    $qualityReport += "### 内容质量"
    $qualityReport += "| 指标 | 数量 | 百分比 | 状态 |"
    $qualityReport += "|------|------|--------|------|"
    
    $metaRate = [math]::Round(($withMeta / $mdFiles.Count) * 100, 1)
    $qualityReport += "| 包含元信息 | $withMeta | $metaRate% | $(if ($metaRate -ge 80) { '✅' } elseif ($metaRate -ge 50) { '⚠️' } else { '❌' }) |"
    
    $codeRate = [math]::Round(($withCode / $mdFiles.Count) * 100, 1)
    $qualityReport += "| 包含代码示例 | $withCode | $codeRate% | $(if ($codeRate -ge 70) { '✅' } elseif ($codeRate -ge 50) { '⚠️' } else { '❌' }) |"
    
    $titleRate = [math]::Round(($withTitle / $mdFiles.Count) * 100, 1)
    $qualityReport += "| 包含标题 | $withTitle | $titleRate% | $(if ($titleRate -ge 95) { '✅' } elseif ($titleRate -ge 80) { '⚠️' } else { '❌' }) |"
    
    $imageRate = [math]::Round(($withImages / $mdFiles.Count) * 100, 1)
    $qualityReport += "| 包含图片 | $withImages | $imageRate% | $(if ($imageRate -ge 30) { '✅' } elseif ($imageRate -ge 15) { '⚠️' } else { '❌' }) |"
    $qualityReport += ""
    
    # 综合评分
    $overallScore = [math]::Round(($metaRate * 0.3 + $codeRate * 0.3 + $titleRate * 0.2 + $imageRate * 0.2), 1)
    $qualityReport += "### 综合评分"
    $qualityReport += "- **分数**: $overallScore / 100"
    $qualityReport += "- **等级**: $(if ($overallScore -ge 80) { '🌟🌟🌟 优秀' } elseif ($overallScore -ge 60) { '🌟🌟 良好' } elseif ($overallScore -ge 40) { '🌟 一般' } else { '❌ 需改进' })"
    $qualityReport += ""
    
    return $qualityReport
}

# 主执行逻辑
try {
    switch ($Check) {
        "env" {
            $report += Check-Environment
        }
        "files" {
            $report += Check-Files
        }
        "links" {
            $report += Check-Links
        }
        "quality" {
            $report += Check-Quality
        }
        "all" {
            $report += Check-Environment
            $report += ""
            $report += "---"
            $report += ""
            $report += Check-Files
            $report += ""
            $report += "---"
            $report += ""
            $report += Check-Links
            $report += ""
            $report += "---"
            $report += ""
            $report += Check-Quality
        }
    }
    
    $report += ""
    $report += "---"
    $report += ""
    $report += "## 📝 建议"
    $report += ""
    $report += "1. 运行 ``.\scripts\check_quality.ps1`` 进行详细质量检查"
    $report += "2. 运行 ``.\scripts\fix_links.ps1`` 修复失效链接"
    $report += "3. 运行 ``.\scripts\generate_statistics.ps1`` 生成详细统计"
    $report += ""
    $report += "---"
    $report += ""
    $report += "**生成工具**: debug_helper.ps1  "
    $report += "**版本**: v1.0"
    
    # 保存报告
    $report | Out-File -FilePath $OutputFile -Encoding UTF8
    
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "✅ 调试报告已生成" -ForegroundColor Green
    Write-Host "📄 文件: $OutputFile" -ForegroundColor Cyan
    Write-Host "========================================" -ForegroundColor Cyan
    
} catch {
    Write-Host "❌ 发生错误: $_" -ForegroundColor Red
    Write-Host "堆栈跟踪: $($_.ScriptStackTrace)" -ForegroundColor Gray
    exit 1
}


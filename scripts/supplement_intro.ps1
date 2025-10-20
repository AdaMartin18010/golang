# 文档简介补充工具
# 功能：分析并补充所有文档的简介部分

param(
    [switch]$DryRun,
    [switch]$Analyze
)

$targetDir = "docs"
$processedCount = 0
$needsIntroCount = 0
$goodIntroCount = 0
$results = @()

Write-Host "📚 开始分析文档简介..." -ForegroundColor Cyan
Write-Host ""

function Get-IntroQuality {
    param([string]$Content)
    
    # 检查是否有简介部分
    if ($Content -match '>\s*(?:📚\s*)?简介[:：]?\s*\n\n?>\s*(.+?)(?:\n\n|\n(?!>)|$)') {
        $intro = $matches[1].Trim()
        $wordCount = $intro.Length
        
        # 评估简介质量
        if ($wordCount -lt 50) {
            return @{ Status = "Weak"; Length = $wordCount; Content = $intro }
        }
        elseif ($wordCount -lt 150) {
            return @{ Status = "Good"; Length = $wordCount; Content = $intro }
        }
        else {
            return @{ Status = "Excellent"; Length = $wordCount; Content = $intro }
        }
    }
    else {
        return @{ Status = "Missing"; Length = 0; Content = "" }
    }
}

function Generate-Intro {
    param(
        [string]$FilePath,
        [string]$Content
    )
    
    # 提取文档标题
    if ($Content -match '^#\s+(.+)$') {
        $title = $matches[1].Trim()
    }
    else {
        $title = [System.IO.Path]::GetFileNameWithoutExtension($FilePath)
    }
    
    # 提取主要章节标题（前5个二级标题）
    $sections = @()
    $matches = [regex]::Matches($Content, '(?m)^##\s+(?:[\p{So}]\s*)?(.+)$')
    foreach ($match in $matches | Select-Object -First 5) {
        $sectionTitle = $match.Groups[1].Value.Trim()
        if ($sectionTitle -notmatch '^(目录|TOC|Table of Contents)$') {
            $sections += $sectionTitle
        }
    }
    
    # 判断文档类型
    $docType = "技术文档"
    if ($FilePath -match 'README\.md$') {
        $docType = "模块指南"
    }
    elseif ($title -match '实战|实践|案例') {
        $docType = "实战指南"
    }
    elseif ($title -match '深入|进阶|高级') {
        $docType = "进阶教程"
    }
    elseif ($title -match '基础|入门') {
        $docType = "基础教程"
    }
    
    # 生成简介模板
    $intro = "> 📚 **简介**`n>`n"
    
    if ($docType -eq "模块指南") {
        $intro += "> 本模块深入讲解$title，系统介绍相关概念、实践方法和最佳实践。"
    }
    else {
        $intro += "> 本文深入探讨$title，系统讲解其核心概念、技术原理和实践应用。"
    }
    
    if ($sections.Count -gt 0) {
        $intro += "内容涵盖"
        $intro += ($sections -join '、') + "等关键主题。"
    }
    
    $intro += "`n>`n> 通过本文，您将全面掌握相关技术要点，并能够在实际项目中应用这些知识。"
    
    return $intro
}

function Add-IntroToDocument {
    param(
        [string]$FilePath,
        [string]$Content
    )
    
    $newIntro = Generate-Intro -FilePath $FilePath -Content $Content
    
    # 查找插入位置（标题后）
    if ($Content -match '(?s)(^#\s+.+?\n\n)(.*)$') {
        $header = $matches[1]
        $rest = $matches[2]
        
        # 检查是否已有TOC
        if ($rest -match '^<!-- TOC START -->') {
            # TOC之前插入简介
            $newContent = $header + $newIntro + "`n`n" + $rest
        }
        else {
            # 直接插入
            $newContent = $header + $newIntro + "`n`n" + $rest
        }
        
        return $newContent
    }
    
    return $Content
}

function Enhance-WeakIntro {
    param(
        [string]$FilePath,
        [string]$Content,
        [string]$CurrentIntro
    )
    
    # 保留原有简介，但增强内容
    $enhancedIntro = Generate-Intro -FilePath $FilePath -Content $Content
    
    # 如果原简介太短，替换它
    if ($CurrentIntro.Length -lt 50) {
        $pattern = '(?s)>\s*(?:📚\s*)?简介[:：]?\s*\n\n?>\s*.+?(?:\n\n|\n(?!>))'
        $Content = $Content -replace $pattern, $enhancedIntro
    }
    
    return $Content
}

# 递归处理所有Markdown文件
Get-ChildItem -Path $targetDir -Filter "*.md" -Recurse | ForEach-Object {
    $file = $_
    $relativePath = $file.FullName.Replace((Get-Location).Path + "\", "")
    
    # 跳过归档目录
    if ($relativePath -match '[\\/](archive|00-备份|Analysis)[\\/]') {
        return
    }
    
    $content = Get-Content -Path $file.FullName -Raw -Encoding UTF8
    $processedCount++
    
    # 分析简介质量
    $quality = Get-IntroQuality -Content $content
    
    $result = [PSCustomObject]@{
        File = $relativePath
        Status = $quality.Status
        Length = $quality.Length
        NeedsWork = $false
    }
    
    if ($Analyze) {
        # 仅分析模式
        if ($quality.Status -eq "Missing") {
            $needsIntroCount++
            $result.NeedsWork = $true
            Write-Host "  ❌ $relativePath - 缺少简介" -ForegroundColor Red
        }
        elseif ($quality.Status -eq "Weak") {
            $needsIntroCount++
            $result.NeedsWork = $true
            Write-Host "  ⚠️  $relativePath - 简介过短 ($($quality.Length)字)" -ForegroundColor Yellow
        }
        else {
            $goodIntroCount++
            Write-Host "  ✅ $relativePath - 简介良好 ($($quality.Length)字)" -ForegroundColor Green
        }
    }
    else {
        # 补充模式
        $modified = $false
        $newContent = $content
        
        if ($quality.Status -eq "Missing") {
            Write-Host "  ➕ $relativePath - 添加简介" -ForegroundColor Cyan
            $newContent = Add-IntroToDocument -FilePath $file.FullName -Content $content
            $modified = $true
            $needsIntroCount++
        }
        elseif ($quality.Status -eq "Weak") {
            Write-Host "  ✏️  $relativePath - 增强简介" -ForegroundColor Yellow
            $newContent = Enhance-WeakIntro -FilePath $file.FullName -Content $content -CurrentIntro $quality.Content
            $modified = $true
            $needsIntroCount++
        }
        else {
            $goodIntroCount++
        }
        
        # 写入文件
        if ($modified -and -not $DryRun) {
            $newContent | Set-Content -Path $file.FullName -Encoding UTF8 -NoNewline
            Write-Host "    💾 已保存" -ForegroundColor Green
        }
        elseif ($modified -and $DryRun) {
            Write-Host "    🔍 [DryRun] 将会修改" -ForegroundColor Gray
        }
    }
    
    $results += $result
}

# 输出统计报告
Write-Host ""
Write-Host "=" * 60 -ForegroundColor Cyan
Write-Host "📊 处理统计" -ForegroundColor Cyan
Write-Host "=" * 60 -ForegroundColor Cyan
Write-Host "总文档数:        $processedCount"
Write-Host "需要补充/增强:   $needsIntroCount" -ForegroundColor $(if ($needsIntroCount -gt 0) { "Yellow" } else { "Green" })
Write-Host "简介良好:        $goodIntroCount" -ForegroundColor Green
Write-Host ""

if ($DryRun) {
    Write-Host "⚠️  这是预演模式，未实际修改文件" -ForegroundColor Yellow
}
elseif ($Analyze) {
    Write-Host "🔍 这是分析模式，未修改文件" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "💡 运行以下命令进行补充:" -ForegroundColor Green
    Write-Host "   .\scripts\supplement_intro.ps1" -ForegroundColor White
}

# 生成详细报告
if ($Analyze) {
    $reportPath = "reports/📝简介分析报告-$(Get-Date -Format 'yyyy-MM-dd').md"
    $report = @"
# 📝 文档简介质量分析报告

**生成时间**: $(Get-Date -Format 'yyyy年MM月dd日 HH:mm:ss')

## 📊 总体统计

- **总文档数**: $processedCount
- **简介优秀**: $($results | Where-Object { $_.Status -eq "Excellent" } | Measure-Object).Count
- **简介良好**: $($results | Where-Object { $_.Status -eq "Good" } | Measure-Object).Count
- **简介过短**: $($results | Where-Object { $_.Status -eq "Weak" } | Measure-Object).Count
- **缺少简介**: $($results | Where-Object { $_.Status -eq "Missing" } | Measure-Object).Count

## 🎯 需要补充的文档

### ❌ 缺少简介

"@
    
    $missingIntros = $results | Where-Object { $_.Status -eq "Missing" }
    if ($missingIntros.Count -eq 0) {
        $report += "`n无`n"
    }
    else {
        foreach ($item in $missingIntros) {
            $report += "- ``$($item.File)```n"
        }
    }
    
    $report += @"

### ⚠️ 简介过短

"@
    
    $weakIntros = $results | Where-Object { $_.Status -eq "Weak" }
    if ($weakIntros.Count -eq 0) {
        $report += "`n无`n"
    }
    else {
        foreach ($item in $weakIntros) {
            $report += "- ``$($item.File)`` ($($item.Length)字)`n"
        }
    }
    
    $report += @"

## ✅ 简介良好的文档

"@
    
    $goodIntros = $results | Where-Object { $_.Status -in @("Good", "Excellent") } | Select-Object -First 20
    foreach ($item in $goodIntros) {
        $report += "- ``$($item.File)`` ($($item.Length)字) ✨`n"
    }
    
    if (($results | Where-Object { $_.Status -in @("Good", "Excellent") }).Count -gt 20) {
        $report += "`n... 还有 $(($results | Where-Object { $_.Status -in @("Good", "Excellent") }).Count - 20) 个文档`n"
    }
    
    $report += @"

## 🎯 质量分布

``````
简介优秀 ($($results | Where-Object { $_.Status -eq "Excellent" } | Measure-Object).Count): $('█' * [Math]::Min(50, ($results | Where-Object { $_.Status -eq "Excellent" } | Measure-Object).Count))
简介良好 ($($results | Where-Object { $_.Status -eq "Good" } | Measure-Object).Count): $('█' * [Math]::Min(50, ($results | Where-Object { $_.Status -eq "Good" } | Measure-Object).Count))
简介过短 ($($results | Where-Object { $_.Status -eq "Weak" } | Measure-Object).Count): $('█' * [Math]::Min(50, ($results | Where-Object { $_.Status -eq "Weak" } | Measure-Object).Count))
缺少简介 ($($results | Where-Object { $_.Status -eq "Missing" } | Measure-Object).Count): $('█' * [Math]::Min(50, ($results | Where-Object { $_.Status -eq "Missing" } | Measure-Object).Count))
``````

## 🚀 下一步行动

1. 运行补充工具: ``.\scripts\supplement_intro.ps1``
2. 预演模式: ``.\scripts\supplement_intro.ps1 -DryRun``
3. 人工审查生成的简介
4. 提交改进后的文档

---

**报告生成者**: Go Documentation Team  
**工具版本**: v1.0  
**文档状态**: 分析完成
"@
    
    $report | Set-Content -Path $reportPath -Encoding UTF8
    Write-Host "📄 详细报告已生成: $reportPath" -ForegroundColor Green
}


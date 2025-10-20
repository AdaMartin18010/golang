# 失效链接自动修复工具
# 功能：分析并修复文档中的失效链接

param(
    [switch]$Analyze,    # 仅分析模式
    [switch]$DryRun,     # 预演模式
    [switch]$AutoFix     # 自动修复模式
)

$targetDir = "docs"
$brokenLinks = @()
$fixedLinks = @()
$cannotFixLinks = @()

Write-Host "🔧 开始链接修复流程..." -ForegroundColor Cyan
Write-Host ""

# 定义文件映射表（已知的重命名/移动）
$fileMapping = @{
    # Go 1.23 -> Go 1.25 映射
    "02-Go语言现代化/12-Go-1.23运行时优化" = "03-Go-1.25新特性/12-Go-1.25运行时优化"
    "01-greentea-GC垃圾收集器.md" = "README.md"  # 特性已整合
    "09-Go 1.23+微服务优化.md" = "09-Go 1.25.1微服务优化.md"
}

# 定义需要创建的占位文件
$placeholderFiles = @(
    "docs/01-语言基础/03-模块管理/02-go-mod文件详解.md",
    "docs/01-语言基础/03-模块管理/03-go-sum文件详解.md",
    "docs/01-语言基础/03-模块管理/04-语义化版本.md",
    "docs/01-语言基础/03-模块管理/06-依赖管理.md",
    "docs/02-Web开发/16-监控和日志.md"
)

function Find-BrokenLinks {
    param([string]$FilePath)
    
    $content = Get-Content -Path $FilePath -Raw -Encoding UTF8
    $relativePath = $FilePath.Replace((Get-Location).Path + "\", "")
    
    # 查找Markdown链接 [text](url)
    $linkPattern = '\[([^\]]+)\]\(([^\)]+)\)'
    $matches = [regex]::Matches($content, $linkPattern)
    
    $fileDir = Split-Path -Path $FilePath -Parent
    $broken = @()
    
    foreach ($match in $matches) {
        $linkText = $match.Groups[1].Value
        $linkUrl = $match.Groups[2].Value
        
        # 跳过外部链接、锚点、mailto等
        if ($linkUrl -match '^(http|https|#|mailto):') {
            continue
        }
        
        # 处理相对路径
        $targetPath = Join-Path -Path $fileDir -ChildPath $linkUrl
        $targetPath = [System.IO.Path]::GetFullPath($targetPath)
        
        # 移除URL中的锚点
        $targetPathNoAnchor = $targetPath -replace '#.*$', ''
        
        # 检查文件是否存在
        if (-not (Test-Path -Path $targetPathNoAnchor)) {
            $broken += [PSCustomObject]@{
                SourceFile = $relativePath
                LinkText = $linkText
                LinkUrl = $linkUrl
                TargetPath = $targetPath
                OriginalMatch = $match.Value
            }
        }
    }
    
    return $broken
}

function Get-SuggestedFix {
    param(
        [string]$SourceFile,
        [string]$BrokenLink
    )
    
    # 检查文件映射表
    foreach ($key in $fileMapping.Keys) {
        if ($BrokenLink -match [regex]::Escape($key)) {
            return $BrokenLink -replace [regex]::Escape($key), $fileMapping[$key]
        }
    }
    
    # 尝试在docs目录下搜索类似文件名
    $fileName = Split-Path -Path $BrokenLink -Leaf
    if ($fileName) {
        $similarFiles = Get-ChildItem -Path $targetDir -Filter "*$fileName*" -Recurse -File | 
                        Where-Object { $_.Extension -eq '.md' }
        
        if ($similarFiles.Count -eq 1) {
            # 找到唯一匹配，返回相对路径
            $sourceDir = Split-Path -Path $SourceFile -Parent
            $targetFile = $similarFiles[0].FullName
            
            # 计算相对路径
            $relativePath = Get-RelativePath -From $sourceDir -To $targetFile
            return $relativePath
        }
    }
    
    return $null
}

function Get-RelativePath {
    param(
        [string]$From,
        [string]$To
    )
    
    # 简化：返回相对于from的to路径
    $fromUri = New-Object System.Uri((Get-Item $From).FullName + "\")
    $toUri = New-Object System.Uri((Get-Item $To).FullName)
    
    $relativeUri = $fromUri.MakeRelativeUri($toUri)
    $relativePath = [System.Uri]::UnescapeDataString($relativeUri.ToString())
    
    return $relativePath -replace '/', '\'
}

function Create-PlaceholderFile {
    param([string]$FilePath)
    
    $fileName = [System.IO.Path]::GetFileNameWithoutExtension($FilePath)
    $content = @"
# $fileName

> 📚 **简介**
>
> 本文档正在编写中，即将完善。

---

## 占位符说明

本文档已被引用但尚未完成编写。

### 计划内容

- [ ] 核心概念介绍
- [ ] 实践示例
- [ ] 最佳实践
- [ ] 常见问题

---

**文档维护者**: Go Documentation Team  
**最后更新**: $(Get-Date -Format 'yyyy年MM月dd日')  
**文档状态**: 规划中  
**适用版本**: Go 1.21+
"@
    
    if (-not $DryRun) {
        $dir = Split-Path -Path $FilePath -Parent
        if (-not (Test-Path -Path $dir)) {
            New-Item -Path $dir -ItemType Directory -Force | Out-Null
        }
        $content | Set-Content -Path $FilePath -Encoding UTF8
        Write-Host "  ✅ 已创建: $FilePath" -ForegroundColor Green
    }
    else {
        Write-Host "  🔍 [DryRun] 将创建: $FilePath" -ForegroundColor Gray
    }
}

# ========== 主流程 ==========

Write-Host "📂 步骤 1: 扫描所有Markdown文件..." -ForegroundColor Yellow
$mdFiles = Get-ChildItem -Path $targetDir -Filter "*.md" -Recurse | 
           Where-Object { $_.FullName -notmatch '[\\/](archive|00-备份|Analysis)[\\/]' }

$totalFiles = $mdFiles.Count
$currentFile = 0

foreach ($file in $mdFiles) {
    $currentFile++
    Write-Progress -Activity "扫描文件" -Status "$currentFile / $totalFiles" -PercentComplete (($currentFile / $totalFiles) * 100)
    
    $broken = Find-BrokenLinks -FilePath $file.FullName
    if ($broken.Count -gt 0) {
        $brokenLinks += $broken
    }
}

Write-Progress -Activity "扫描文件" -Completed
Write-Host "✅ 扫描完成！" -ForegroundColor Green
Write-Host ""

# 统计
Write-Host "=" * 60 -ForegroundColor Cyan
Write-Host "📊 失效链接统计" -ForegroundColor Cyan
Write-Host "=" * 60 -ForegroundColor Cyan
Write-Host "总失效链接数: $($brokenLinks.Count)" -ForegroundColor Red
Write-Host ""

if ($brokenLinks.Count -eq 0) {
    Write-Host "🎉 太好了！没有发现失效链接！" -ForegroundColor Green
    exit 0
}

# 按源文件分组
$groupedLinks = $brokenLinks | Group-Object -Property SourceFile
Write-Host "受影响文件数: $($groupedLinks.Count)" -ForegroundColor Yellow
Write-Host ""

if ($Analyze) {
    # 仅分析模式
    Write-Host "🔍 失效链接详情:" -ForegroundColor Cyan
    Write-Host ""
    
    foreach ($group in $groupedLinks | Sort-Object Name) {
        Write-Host "📄 $($group.Name)" -ForegroundColor White
        foreach ($link in $group.Group) {
            Write-Host "   ❌ [$($link.LinkText)]($($link.LinkUrl))" -ForegroundColor Red
            
            # 尝试建议修复方案
            $suggestion = Get-SuggestedFix -SourceFile $link.SourceFile -BrokenLink $link.LinkUrl
            if ($suggestion) {
                Write-Host "      💡 建议: $suggestion" -ForegroundColor Yellow
            }
        }
        Write-Host ""
    }
    
    Write-Host "提示: 运行 -AutoFix 参数进行自动修复" -ForegroundColor Green
    exit 0
}

# ========== 自动修复流程 ==========

if ($AutoFix -or $DryRun) {
    Write-Host "🔧 步骤 2: 开始自动修复流程..." -ForegroundColor Yellow
    Write-Host ""
    
    # 策略 1: 创建占位文件
    Write-Host "策略 1: 创建缺失的重要文件" -ForegroundColor Cyan
    foreach ($filePath in $placeholderFiles) {
        if (-not (Test-Path -Path $filePath)) {
            Create-PlaceholderFile -FilePath $filePath
        }
    }
    Write-Host ""
    
    # 策略 2: 更新链接路径
    Write-Host "策略 2: 修复可自动更正的链接" -ForegroundColor Cyan
    
    $fileUpdates = @{}  # 文件 -> 需要的替换操作
    
    foreach ($link in $brokenLinks) {
        $suggestion = Get-SuggestedFix -SourceFile $link.SourceFile -BrokenLink $link.LinkUrl
        
        if ($suggestion) {
            # 可以自动修复
            $sourceFullPath = Join-Path -Path (Get-Location).Path -ChildPath $link.SourceFile
            
            if (-not $fileUpdates.ContainsKey($sourceFullPath)) {
                $fileUpdates[$sourceFullPath] = @()
            }
            
            $fileUpdates[$sourceFullPath] += @{
                Old = $link.OriginalMatch
                New = "[$($link.LinkText)]($suggestion)"
                LinkUrl = $link.LinkUrl
                Suggestion = $suggestion
            }
            
            $fixedLinks += $link
            Write-Host "  ✅ $($link.SourceFile)" -ForegroundColor Green
            Write-Host "     旧: $($link.LinkUrl)" -ForegroundColor Red
            Write-Host "     新: $suggestion" -ForegroundColor Green
        }
        else {
            $cannotFixLinks += $link
        }
    }
    
    # 应用文件更新
    if (-not $DryRun) {
        foreach ($file in $fileUpdates.Keys) {
            $content = Get-Content -Path $file -Raw -Encoding UTF8
            
            foreach ($update in $fileUpdates[$file]) {
                $content = $content -replace [regex]::Escape($update.Old), $update.New
            }
            
            $content | Set-Content -Path $file -Encoding UTF8 -NoNewline
        }
    }
    
    Write-Host ""
    Write-Host "=" * 60 -ForegroundColor Cyan
    Write-Host "📊 修复结果" -ForegroundColor Cyan
    Write-Host "=" * 60 -ForegroundColor Cyan
    Write-Host "✅ 已修复链接:   $($fixedLinks.Count)" -ForegroundColor Green
    Write-Host "⚠️  无法自动修复: $($cannotFixLinks.Count)" -ForegroundColor Yellow
    Write-Host ""
    
    if ($cannotFixLinks.Count -gt 0) {
        Write-Host "⚠️  需要人工处理的链接:" -ForegroundColor Yellow
        Write-Host ""
        
        foreach ($link in $cannotFixLinks | Select-Object -First 10) {
            Write-Host "  📄 $($link.SourceFile)" -ForegroundColor White
            Write-Host "     ❌ $($link.LinkUrl)" -ForegroundColor Red
        }
        
        if ($cannotFixLinks.Count -gt 10) {
            Write-Host ""
            Write-Host "  ... 还有 $($cannotFixLinks.Count - 10) 个链接" -ForegroundColor Gray
        }
    }
    
    if ($DryRun) {
        Write-Host ""
        Write-Host "⚠️  这是预演模式，未实际修改文件" -ForegroundColor Yellow
        Write-Host "💡 运行 -AutoFix 参数应用修复" -ForegroundColor Green
    }
}
else {
    Write-Host "💡 使用方式:" -ForegroundColor Green
    Write-Host "  -Analyze   : 分析失效链接" -ForegroundColor White
    Write-Host "  -DryRun    : 预览修复方案" -ForegroundColor White
    Write-Host "  -AutoFix   : 应用自动修复" -ForegroundColor White
}

Write-Host ""
Write-Host "✨ 完成！" -ForegroundColor Green


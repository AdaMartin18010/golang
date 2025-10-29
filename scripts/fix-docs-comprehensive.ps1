# 全面文档格式修复脚本
# 版本: 2.0 - 更全面的修复
# 日期: 2025-10-29

param(
    [string]$Path = "docs",
    [switch]$DryRun,
    [switch]$Verbose
)

$ErrorActionPreference = "Stop"
$stats = @{
    TitleFixed = 0
    MetadataFixed = 0
    DateFixed = 0
    FilesProcessed = 0
    Errors = 0
}

Write-Host "🚀 开始全面文档格式修复 v2.0..." -ForegroundColor Cyan
Write-Host "工作目录: $Path"
Write-Host "模式: $(if($DryRun){'试运行'}else{'实际修复'})`n"

function Fix-ComplexTitles {
    param($FilePath, $Content)
    
    $modified = $false
    $filename = Split-Path $FilePath -Leaf
    
    # 修复包含 "Go xxx -" 的复杂标题
    $patterns = @(
        @{ Pattern = '^# 📊 Go.* - .*'; Example = '# 📊 Go 1.25.3实战开发导航 - 2025' }
        @{ Pattern = '^# 📝 Go.* - .*'; Example = '# 📝 Go xxx - xxx' }
        @{ Pattern = '^# 📚 Go.* - .*'; Example = '# 📚 Go xxx - xxx' }
        @{ Pattern = '^# ❓ Go.* - .*'; Example = '# ❓ Go xxx - xxx' }
        @{ Pattern = '^# ✅ Go.* - .*'; Example = '# ✅ Go xxx - xxx' }
        @{ Pattern = '^# 🎯 Go.* - .*'; Example = '# 🎯 Go xxx - xxx' }
    )
    
    foreach ($p in $patterns) {
        if ($Content -match $p.Pattern) {
            # 提取标题
            if ($Content -match '^(# [^-]+) - (.+)$') {
                $newTitle = $matches[1].Trim()
                if ($Verbose) { 
                    Write-Host "  [复杂标题] $filename" -ForegroundColor Yellow 
                    Write-Host "    旧: $($matches[0])" -ForegroundColor Gray
                    Write-Host "    新: $newTitle" -ForegroundColor Green
                }
                $Content = $Content -replace '^# [^\r\n]+', $newTitle
                $modified = $true
                break
            }
        }
    }
    
    return @{ Modified = $modified; Content = $Content }
}

function Fix-Metadata {
    param($Content)
    
    $modified = $false
    
    # 统一日期格式为 2025-10-29
    if ($Content -match '\*\*更新日期\*\*:\s*2025-10-2[0-8]') {
        $Content = $Content -replace '\*\*更新日期\*\*:\s*2025-10-2[0-8]', '**更新日期**: 2025-10-29'
        $modified = $true
    }
    
    # 统一 "最后更新" 为 2025-10-29
    if ($Content -match '\*\*最后更新\*\*:\s*2025-10-2[0-8]') {
        $Content = $Content -replace '\*\*最后更新\*\*:\s*2025-10-2[0-8]', '**最后更新**: 2025-10-29'
        $modified = $true
    }
    
    return @{ Modified = $modified; Content = $Content }
}

function Fix-SpecialFiles {
    param($FilePath, $Content)
    
    $modified = $false
    $filename = Split-Path $FilePath -Leaf
    
    # 修复 "版本对比与选择指南" 的标题
    if ($filename -eq "00-版本对比与选择指南.md") {
        if ($Content -match '^# 📊 版本对比 - v1\.x vs v2\.0\.0') {
            $Content = $Content -replace '^# 📊 版本对比 - v1\.x vs v2\.0\.0', '# Go版本对比与选择指南'
            if ($Verbose) { Write-Host "  [特殊] 修复版本对比文档标题" -ForegroundColor Cyan }
            $modified = $true
        }
    }
    
    # 修复 "发布说明.md" 的复杂标题
    if ($filename -eq "发布说明.md") {
        if ($Content -match '^# 📋 .*发布说明 - ') {
            $Content = $Content -replace '^# 📋 .*发布说明 - .*', '# 发布说明'
            if ($Verbose) { Write-Host "  [特殊] 修复发布说明标题" -ForegroundColor Cyan }
            $modified = $true
        }
    }
    
    return @{ Modified = $modified; Content = $Content }
}

#region 主执行逻辑

try {
    $files = Get-ChildItem -Path $Path -Recurse -Filter "*.md" -File
    $totalFiles = $files.Count
    
    Write-Host "找到 $totalFiles 个Markdown文件`n" -ForegroundColor Green
    
    $progress = 0
    foreach ($file in $files) {
        $progress++
        $percentComplete = [math]::Round(($progress / $totalFiles) * 100)
        
        Write-Progress -Activity "处理文档" -Status "$progress/$totalFiles" -PercentComplete $percentComplete
        
        try {
            $stats.FilesProcessed++
            $content = Get-Content $file.FullName -Raw -Encoding UTF8
            $originalContent = $content
            $fileModified = $false
            
            # 1. 修复复杂标题
            $result = Fix-ComplexTitles -FilePath $file.FullName -Content $content
            if ($result.Modified) {
                $content = $result.Content
                $fileModified = $true
                $stats.TitleFixed++
            }
            
            # 2. 修复元数据
            $result = Fix-Metadata -Content $content
            if ($result.Modified) {
                $content = $result.Content
                $fileModified = $true
                $stats.MetadataFixed++
            }
            
            # 3. 修复特殊文件
            $result = Fix-SpecialFiles -FilePath $file.FullName -Content $content
            if ($result.Modified) {
                $content = $result.Content
                $fileModified = $true
                $stats.TitleFixed++
            }
            
            # 保存修改
            if ($fileModified) {
                if ($Verbose) {
                    Write-Host "✓ 修复: $($file.Name)" -ForegroundColor Green
                }
                
                if (-not $DryRun) {
                    Set-Content -Path $file.FullName -Value $content -NoNewline -Encoding UTF8
                }
            }
            
        } catch {
            $stats.Errors++
            Write-Host "✗ 错误: $($file.Name) - $($_.Exception.Message)" -ForegroundColor Red
        }
    }
    
    Write-Progress -Activity "处理文档" -Completed
    
} catch {
    Write-Host "`n❌ 发生错误: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

#endregion

#region 结果报告

Write-Host "`n" + ("="*60) -ForegroundColor Cyan
Write-Host "📊 修复统计报告 v2.0" -ForegroundColor Cyan
Write-Host ("="*60) -ForegroundColor Cyan

Write-Host "`n文件处理:"
Write-Host "  ✓ 已处理: $($stats.FilesProcessed)" -ForegroundColor Green
Write-Host "  ⚠ 错误: $($stats.Errors)" -ForegroundColor $(if($stats.Errors -gt 0){'Red'}else{'Gray'})

Write-Host "`n修复详情:"
Write-Host "  🏷️  标题修复: $($stats.TitleFixed) 个文件" -ForegroundColor Yellow
Write-Host "  📝 元数据修复: $($stats.MetadataFixed) 个文件" -ForegroundColor Yellow
Write-Host "  📋 总修复: $($stats.TitleFixed + $stats.MetadataFixed) 次" -ForegroundColor Yellow

if ($DryRun) {
    Write-Host "`n⚠️  这是试运行，未实际修改文件" -ForegroundColor Yellow
    Write-Host "执行实际修复请移除 -DryRun 参数" -ForegroundColor Yellow
} else {
    Write-Host "`n✅ 修复已完成！" -ForegroundColor Green
}

Write-Host "`n" + ("="*60) -ForegroundColor Cyan

#endregion

return $stats


# PowerShell脚本: 更新所有文档中的Go版本到1.25.3
# 用法: .\scripts\update_go_version.ps1 [-DryRun] [-TargetVersion "1.25.3"]

param(
    [switch]$DryRun,
    [string]$TargetVersion = "1.25.3"
)

$ErrorActionPreference = "Stop"
$docsRoot = Join-Path $PSScriptRoot "..\docs"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Go版本更新脚本 v1.0" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "目标版本: Go $TargetVersion" -ForegroundColor Yellow
Write-Host "文档根目录: $docsRoot" -ForegroundColor Yellow
Write-Host "运行模式: $(if ($DryRun) { 'Dry Run (预览)' } else { '实际更新' })" -ForegroundColor Yellow
Write-Host ""

# 获取所有Markdown文件（排除archive目录）
$allMarkdownFiles = Get-ChildItem -Path $docsRoot -Recurse -Include *.md | Where-Object {
    $_.FullName -notmatch "\\archive\\" -and
    $_.FullName -notmatch "\\00-备份\\"
}

Write-Host "📊 找到 $($allMarkdownFiles.Count) 个活跃文档" -ForegroundColor Green
Write-Host ""

$updatedFiles = @()
$updateStats = @{
    "文档元数据" = 0
    "go.mod版本" = 0
    "通用版本引用" = 0
    "版本范围" = 0
}

foreach ($file in $allMarkdownFiles) {
    $originalContent = Get-Content -Path $file.FullName -Raw -Encoding UTF8
    $newContent = $originalContent
    $fileModified = $false
    $changeLog = @()

    # 规则1: 更新文档元数据中的"适用版本"
    # 匹配: **适用版本**: Go 1.x+ 或 Go 1.xx+
    if ($originalContent -match '\*\*适用版本\*\*:\s*Go\s+1\.\d+\+') {
        $newContent = $newContent -replace '(\*\*适用版本\*\*:\s*)Go\s+1\.\d+\+', "`${1}Go $TargetVersion+"
        $changeLog += "更新适用版本 → Go $TargetVersion+"
        $updateStats["文档元数据"]++
        $fileModified = $true
    }

    # 规则2: 更新go.mod中的go版本声明
    # 匹配: go 1.21 或 go 1.22 等
    if ($originalContent -match 'go\s+1\.\d{1,2}\s*$' -or $originalContent -match 'go\s+1\.\d{1,2}\r?\n') {
        $newContent = $newContent -replace '(\r?\n|\A)go\s+1\.\d{1,2}(\r?\n|\Z)', "`${1}go $TargetVersion`${2}"
        $changeLog += "更新go.mod版本 → go $TargetVersion"
        $updateStats["go.mod版本"]++
        $fileModified = $true
    }

    # 规则3: 更新通用版本引用（在代码块外）
    # Go 1.21、Go 1.22、Go 1.23、Go 1.24、Go 1.25 → Go 1.25.3
    # 但保留性能对比报告中的历史数据（包含"vs"或"对比"的行）
    if ($originalContent -match 'Go\s+1\.(2[0-5])\s') {
        # 排除性能对比语境和图表中的版本引用
        $lines = $newContent -split "`r?`n"
        $inCodeBlock = $false
        $lineNumber = 0
        
        foreach ($line in $lines) {
            $lineNumber++
            
            # 跟踪代码块状态
            if ($line -match '^```') {
                $inCodeBlock = -not $inCodeBlock
            }
            
            # 跳过代码块、性能对比行、图表行
            if ($inCodeBlock -or 
                $line -match '(vs|对比|比较|提升|改善)' -or
                $line -match '[│┤├─╭╯▓█]' -or
                $line -match '^\s*\|.*\|' -or
                $file.Name -match '性能对比报告|FAQ') {
                continue
            }
            
            # 更新普通文本中的版本引用
            if ($line -match 'Go\s+1\.(2[0-4])\b' -and $line -notmatch 'Go\s+1\.25\.') {
                $newLine = $line -replace 'Go\s+1\.(2[0-4])\b', "Go $TargetVersion"
                if ($newLine -ne $line) {
                    $newContent = $newContent -replace [regex]::Escape($line), $newLine
                    if (-not $changeLog.Contains("更新通用版本引用")) {
                        $changeLog += "更新通用版本引用"
                        $updateStats["通用版本引用"]++
                    }
                    $fileModified = $true
                }
            }
        }
    }

    # 规则4: 更新版本范围引用
    # Go 1.21-1.24 → Go 1.21-1.25.3
    if ($originalContent -match 'Go\s+1\.\d+-1\.\d+') {
        $newContent = $newContent -replace 'Go\s+1\.(\d+)-1\.\d+', "Go 1.`${1}-$TargetVersion"
        $changeLog += "更新版本范围"
        $updateStats["版本范围"]++
        $fileModified = $true
    }

    # 保存修改
    if ($fileModified) {
        $relativePath = $file.FullName -replace [regex]::Escape($docsRoot), "docs"
        $updatedFiles += [PSCustomObject]@{
            Path = $relativePath
            Changes = $changeLog -join "; "
        }
        
        if (-not $DryRun) {
            Set-Content -Path $file.FullName -Value $newContent -Encoding UTF8 -NoNewline
        }
    }
}

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  📊 更新统计" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "更新类型统计:" -ForegroundColor Yellow
foreach ($type in $updateStats.Keys | Sort-Object) {
    $count = $updateStats[$type]
    if ($count -gt 0) {
        Write-Host "  • $type : $count 处" -ForegroundColor Green
    }
}
Write-Host ""

Write-Host "更新文件总数: $($updatedFiles.Count)" -ForegroundColor Green
Write-Host ""

if ($updatedFiles.Count -gt 0) {
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "  📝 更新文件详情 (前20个)" -ForegroundColor Cyan
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host ""
    
    $updatedFiles | Select-Object -First 20 | ForEach-Object {
        Write-Host "📄 $($_.Path)" -ForegroundColor Cyan
        Write-Host "   $($_.Changes)" -ForegroundColor DarkGray
        Write-Host ""
    }
    
    if ($updatedFiles.Count -gt 20) {
        Write-Host "... 还有 $($updatedFiles.Count - 20) 个文件" -ForegroundColor DarkGray
        Write-Host ""
    }
}

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  ✅ 完成" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

if ($DryRun) {
    Write-Host "⚠️  这是预览模式，文件未实际修改" -ForegroundColor Yellow
    Write-Host "💡 运行 -DryRun `$false 参数应用更改" -ForegroundColor Yellow
} else {
    Write-Host "✅ 已更新 $($updatedFiles.Count) 个文件到 Go $TargetVersion" -ForegroundColor Green
    Write-Host "💡 建议运行: git diff 查看变更" -ForegroundColor Yellow
}

Write-Host ""


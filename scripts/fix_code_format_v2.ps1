# 代码格式统一修复脚本 v2
# 功能：统一代码块格式、Mermaid图表格式、移除多余空行

param(
    [string]$TargetDir = "docs",
    [switch]$DryRun = $false
)

$ErrorActionPreference = "Stop"

# 统计变量
$stats = @{
    FilesScanned = 0
    FilesModified = 0
    Issues = @{
        TrailingEmptyLinesInCodeBlocks = 0
        ExcessiveNewlines = 0
        TrailingSpaces = 0
    }
}

function Write-Log {
    param([string]$Message, [string]$Level = "INFO")
    $timestamp = Get-Date -Format "HH:mm:ss"
    $color = switch ($Level) {
        "ERROR" { "Red" }
        "WARN"  { "Yellow" }
        "SUCCESS" { "Green" }
        default { "White" }
    }
    Write-Host "[$timestamp] $Level`: $Message" -ForegroundColor $color
}

function Process-MarkdownFile {
    param([string]$FilePath)
    
    $stats.FilesScanned++
    
    try {
        # 读取文件
        $content = Get-Content -Path $FilePath -Raw -Encoding UTF8
        if (-not $content) {
            return
        }
        
        $originalContent = $content
        $modified = $false
        
        # 1. 移除行尾空格
        $lines = $content -split '\r?\n'
        $cleanedLines = @()
        foreach ($line in $lines) {
            $trimmed = $line.TrimEnd()
            if ($trimmed -ne $line) {
                $modified = $true
                $stats.Issues.TrailingSpaces++
            }
            $cleanedLines += $trimmed
        }
        $content = $cleanedLines -join "`r`n"
        
        # 2. 移除代码块末尾的多余空行
        # 匹配 ``` 前面的连续空行
        $codeBlockPattern = '(\r?\n){2,}(\r?\n)```'
        if ($content -match $codeBlockPattern) {
            $content = $content -replace $codeBlockPattern, "`r`n```"
            $modified = $true
            $stats.Issues.TrailingEmptyLinesInCodeBlocks++
        }
        
        # 3. 修复连续3个以上空行 → 2个空行
        $excessiveNewlinePattern = '(\r?\n\s*\r?\n){3,}'
        if ($content -match $excessiveNewlinePattern) {
            $content = $content -replace $excessiveNewlinePattern, "`r`n`r`n"
            $modified = $true
            $stats.Issues.ExcessiveNewlines++
        }
        
        # 4. 确保文件末尾只有一个换行符
        $content = $content.TrimEnd() + "`r`n"
        
        # 保存文件
        if ($modified -or ($content -ne $originalContent)) {
            if (-not $DryRun) {
                [System.IO.File]::WriteAllText($FilePath, $content, [System.Text.UTF8Encoding]::new($false))
                Write-Log "Modified: $(Split-Path $FilePath -Leaf)" "SUCCESS"
            } else {
                Write-Log "Would modify: $(Split-Path $FilePath -Leaf)" "WARN"
            }
            $stats.FilesModified++
        }
        
    } catch {
        Write-Log "Error processing $FilePath`: $_" "ERROR"
    }
}

# 主程序
Write-Log "========================================" "INFO"
Write-Log "代码格式统一修复脚本 v2" "INFO"
Write-Log "========================================" "INFO"
Write-Log "目标目录: $TargetDir" "INFO"
Write-Log "模式: $(if ($DryRun) { 'Dry Run (预览)' } else { '实际修改' })" "INFO"
Write-Log "========================================" "INFO"

# 获取所有Markdown文件
$mdFiles = Get-ChildItem -Path $TargetDir -Filter "*.md" -Recurse -File

Write-Log "找到 $($mdFiles.Count) 个Markdown文件" "INFO"
Write-Log "" "INFO"

# 处理每个文件
foreach ($file in $mdFiles) {
    Process-MarkdownFile -FilePath $file.FullName
}

# 输出统计报告
Write-Log "" "INFO"
Write-Log "========================================" "SUCCESS"
Write-Log "修复完成！统计报告：" "SUCCESS"
Write-Log "========================================" "INFO"
Write-Log "扫描文件数: $($stats.FilesScanned)" "INFO"
Write-Log "修改文件数: $($stats.FilesModified)" "INFO"
Write-Log "----------------------------------------" "INFO"
Write-Log "修复问题统计：" "INFO"
Write-Log "代码块尾部空行: $($stats.Issues.TrailingEmptyLinesInCodeBlocks)" "INFO"
Write-Log "连续空行修复: $($stats.Issues.ExcessiveNewlines)" "INFO"
Write-Log "行尾空格移除: $($stats.Issues.TrailingSpaces)" "INFO"
Write-Log "========================================" "INFO"

if ($DryRun) {
    Write-Log "" "WARN"
    Write-Log "这是Dry Run模式，未实际修改文件" "WARN"
    Write-Log "移除 -DryRun 参数以执行实际修改" "WARN"
}


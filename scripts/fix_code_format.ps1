# 代码格式统一修复脚本
# 功能：统一代码块格式、Mermaid图表格式、移除多余空行

param(
    [string]$TargetDir = "docs",
    [switch]$DryRun = $false,
    [switch]$Verbose = $false
)

$ErrorActionPreference = "Stop"

# 统计变量
$stats = @{
    FilesScanned = 0
    FilesModified = 0
    CodeBlocksFixed = 0
    MermaidFixed = 0
    TrailingSpacesRemoved = 0
    ExcessiveNewlinesFixed = 0
}

# 问题类型
$issues = @{
    MissingLanguageTag = 0
    IncorrectIndent = 0
    TrailingEmptyLines = 0
    MermaidFormat = 0
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

function Fix-CodeBlock {
    param([string]$Content)
    
    $modified = $false
    $lines = $Content -split '\r?\n'
    $result = New-Object System.Collections.ArrayList
    $inCodeBlock = $false
    $codeBlockLang = ""
    $i = 0
    
    while ($i -lt $lines.Count) {
        $line = $lines[$i]
        
        # 检测代码块开始
        if ($line -match '^```(\w*)(.*)$') {
            $lang = $matches[1]
            $extra = $matches[2].Trim()
            
            # 修复：统一语言标记
            if ($lang -eq "" -and $i + 1 -lt $lines.Count) {
                # 尝试推断语言
                $nextLine = $lines[$i + 1]
                if ($nextLine -match '^(package|func|import|type|var|const)\s') {
                    $lang = "go"
                    $modified = $true
                    $issues.MissingLanguageTag++
                } elseif ($nextLine -match '^(\$|#|cd|ls|git|npm|go\s)') {
                    $lang = "bash"
                    $modified = $true
                    $issues.MissingLanguageTag++
                }
            }
            
            # 统一格式：```language (不要额外空格或内容)
            if ($extra -ne "") {
                [void]$result.Add("```$lang")
                $modified = $true
            } else {
                [void]$result.Add("```$lang")
            }
            
            $inCodeBlock = $true
            $codeBlockLang = $lang
            $i++
            continue
        }
        
        # 检测代码块结束
        if ($line -match '^```\s*$' -and $inCodeBlock) {
            # 移除代码块末尾的空行
            while ($result.Count -gt 0 -and $result[$result.Count - 1] -match '^\s*$') {
                [void]$result.RemoveAt($result.Count - 1)
                $modified = $true
                $issues.TrailingEmptyLines++
            }
            
            [void]$result.Add("```")
            $inCodeBlock = $false
            $codeBlockLang = ""
            
            # 检查代码块后的空行
            if ($i + 1 -lt $lines.Count -and $lines[$i + 1] -notmatch '^\s*$') {
                # 代码块后应该有一个空行
                [void]$result.Add("")
                $modified = $true
            }
            
            $i++
            continue
        }
        
        # 在代码块内
        if ($inCodeBlock) {
            # 保留代码块内的内容，但移除行尾空格
            $trimmed = $line.TrimEnd()
            if ($trimmed -ne $line) {
                $modified = $true
                $stats.TrailingSpacesRemoved++
            }
            [void]$result.Add($trimmed)
        } else {
            # 在代码块外，移除行尾空格
            $trimmed = $line.TrimEnd()
            if ($trimmed -ne $line) {
                $modified = $true
                $stats.TrailingSpacesRemoved++
            }
            [void]$result.Add($trimmed)
        }
        
        $i++
    }
    
    return @{
        Content = ($result -join [Environment]::NewLine)
        Modified = $modified
    }
}

function Fix-MermaidDiagram {
    param([string]$Content)
    
    $modified = $false
    $lines = $Content -split '\r?\n'
    $result = New-Object System.Collections.ArrayList
    $inMermaid = $false
    $i = 0
    
    while ($i -lt $lines.Count) {
        $line = $lines[$i]
        
        # 检测Mermaid开始
        if ($line -match '^```mermaid\s*(.*)$') {
            $extra = $matches[1].Trim()
            
            # 统一格式：```mermaid (不要额外内容)
            if ($extra -ne "") {
                [void]$result.Add("```mermaid")
                $modified = $true
                $issues.MermaidFormat++
            } else {
                [void]$result.Add("```mermaid")
            }
            
            $inMermaid = $true
            $i++
            continue
        }
        
        # 检测Mermaid结束
        if ($line -match '^```\s*$' -and $inMermaid) {
            # 移除Mermaid末尾的空行
            while ($result.Count -gt 0 -and $result[$result.Count - 1] -match '^\s*$') {
                [void]$result.RemoveAt($result.Count - 1)
                $modified = $true
                $issues.TrailingEmptyLines++
            }
            
            [void]$result.Add("```")
            $inMermaid = $false
            
            # Mermaid后应该有一个空行
            if ($i + 1 -lt $lines.Count -and $lines[$i + 1] -notmatch '^\s*$') {
                [void]$result.Add("")
                $modified = $true
            }
            
            $i++
            continue
        }
        
        # 在Mermaid内，移除行尾空格
        if ($inMermaid) {
            $trimmed = $line.TrimEnd()
            if ($trimmed -ne $line) {
                $modified = $true
                $stats.TrailingSpacesRemoved++
            }
            [void]$result.Add($trimmed)
        } else {
            [void]$result.Add($line)
        }
        
        $i++
    }
    
    return @{
        Content = ($result -join [Environment]::NewLine)
        Modified = $modified
    }
}

function Fix-ExcessiveNewlines {
    param([string]$Content)
    
    # 修复：连续3个以上空行 → 2个空行
    $modified = $false
    $pattern = '(\r?\n\s*){4,}'
    $replacement = [Environment]::NewLine + [Environment]::NewLine
    
    if ($Content -match $pattern) {
        $Content = $Content -replace $pattern, $replacement
        $modified = $true
        $stats.ExcessiveNewlinesFixed++
    }
    
    return @{
        Content = $Content
        Modified = $modified
    }
}

function Process-MarkdownFile {
    param([string]$FilePath)
    
    if ($Verbose) {
        Write-Log "Processing: $FilePath" "INFO"
    }
    
    $stats.FilesScanned++
    
    try {
        # 读取文件
        $content = Get-Content -Path $FilePath -Raw -Encoding UTF8
        $originalContent = $content
        $fileModified = $false
        
        # 1. 修复代码块格式
        $result1 = Fix-CodeBlock -Content $content
        $content = $result1.Content
        if ($result1.Modified) {
            $fileModified = $true
            $stats.CodeBlocksFixed++
        }
        
        # 2. 修复Mermaid图表格式
        $result2 = Fix-MermaidDiagram -Content $content
        $content = $result2.Content
        if ($result2.Modified) {
            $fileModified = $true
            $stats.MermaidFixed++
        }
        
        # 3. 修复连续空行
        $result3 = Fix-ExcessiveNewlines -Content $content
        $content = $result3.Content
        if ($result3.Modified) {
            $fileModified = $true
        }
        
        # 4. 确保文件末尾只有一个换行符
        $content = $content.TrimEnd() + [Environment]::NewLine
        
        # 保存文件
        if ($fileModified) {
            if (-not $DryRun) {
                Set-Content -Path $FilePath -Value $content -Encoding UTF8 -NoNewline
                Write-Log "Modified: $FilePath" "SUCCESS"
            } else {
                Write-Log "Would modify: $FilePath" "WARN"
            }
            $stats.FilesModified++
        }
        
    } catch {
        Write-Log "Error processing $FilePath`: $_" "ERROR"
    }
}

# 主程序
Write-Log "========================================" "INFO"
Write-Log "代码格式统一修复脚本" "INFO"
Write-Log "========================================" "INFO"
Write-Log "目标目录: $TargetDir" "INFO"
Write-Log "模式: $(if ($DryRun) { 'Dry Run (预览)' } else { '实际修改' })" "INFO"
Write-Log "========================================" "INFO"

# 获取所有Markdown文件
$mdFiles = Get-ChildItem -Path $TargetDir -Filter "*.md" -Recurse -File

Write-Log "找到 $($mdFiles.Count) 个Markdown文件" "INFO"

# 处理每个文件
foreach ($file in $mdFiles) {
    Process-MarkdownFile -FilePath $file.FullName
}

# 输出统计报告
Write-Log "`n========================================" "INFO"
Write-Log "修复完成！统计报告：" "SUCCESS"
Write-Log "========================================" "INFO"
Write-Log "扫描文件数: $($stats.FilesScanned)" "INFO"
Write-Log "修改文件数: $($stats.FilesModified)" "INFO"
Write-Log "代码块修复: $($stats.CodeBlocksFixed)" "INFO"
Write-Log "Mermaid修复: $($stats.MermaidFixed)" "INFO"
Write-Log "行尾空格移除: $($stats.TrailingSpacesRemoved)" "INFO"
Write-Log "连续空行修复: $($stats.ExcessiveNewlinesFixed)" "INFO"
Write-Log "----------------------------------------" "INFO"
Write-Log "问题类型统计：" "INFO"
Write-Log "缺失语言标记: $($issues.MissingLanguageTag)" "INFO"
Write-Log "代码块尾部空行: $($issues.TrailingEmptyLines)" "INFO"
Write-Log "Mermaid格式: $($issues.MermaidFormat)" "INFO"
Write-Log "========================================" "INFO"

if ($DryRun) {
    Write-Log "这是Dry Run模式，未实际修改文件" "WARN"
    Write-Log "移除 -DryRun 参数以执行实际修改" "WARN"
}


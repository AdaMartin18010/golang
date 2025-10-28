# 自动生成完整的多层级目录结构
param(
    [Parameter(Mandatory=$false)]
    [string]$TargetFile = "",
    [Parameter(Mandatory=$false)]
    [switch]$Batch = $false
)

function Generate-TOC {
    param([string]$FilePath)
    
    $content = Get-Content $FilePath -Raw -Encoding UTF8
    
    # 提取所有标题（排除目录标题本身）
    $headings = [regex]::Matches($content, '^(#{2,6})\s+(.+)$', [System.Text.RegularExpressions.RegexOptions]::Multiline) |
        Where-Object { $_.Groups[2].Value -notmatch '📋\s*目录' }
    
    if ($headings.Count -eq 0) {
        return $null
    }
    
    # 生成目录
    $toc = @()
    $toc += ""
    
    foreach ($heading in $headings) {
        $level = $heading.Groups[1].Value.Length - 1  # ##=1, ###=2, ####=3...
        $title = $heading.Groups[2].Value.Trim()
        
        # 生成锚点（GitHub风格）
        $anchor = $title -replace '\s+', '-'
        $anchor = $anchor -replace '[^\p{L}\p{N}\-_]', ''
        $anchor = $anchor.ToLower()
        
        # 生成缩进
        $indent = "  " * ($level - 1)
        
        # 生成目录项
        $tocItem = "$indent- [$title](#$anchor)"
        $toc += $tocItem
    }
    
    $toc += ""
    
    return ($toc -join "`n")
}

function Replace-TOC {
    param(
        [string]$FilePath,
        [string]$NewTOC
    )
    
    $content = Get-Content $FilePath -Raw -Encoding UTF8
    
    # 查找并替换目录部分
    if ($content -match '(?s)(##\s*📋\s*目录\s*\n)(.*?)(\n##\s)') {
        $before = $matches[1]
        $after = $matches[3]
        
        # 替换
        $newContent = $content -replace '(?s)(##\s*📋\s*目录\s*\n)(.*?)(\n##\s)', "`$1$NewTOC`$3"
        
        # 写回文件
        [System.IO.File]::WriteAllText($FilePath, $newContent, [System.Text.UTF8Encoding]::new($false))
        
        return $true
    }
    
    return $false
}

# 主逻辑
if ($TargetFile -ne "") {
    # 单文件模式
    Write-Output "处理文件: $TargetFile"
    
    $newTOC = Generate-TOC -FilePath $TargetFile
    
    if ($null -eq $newTOC) {
        Write-Output "  ✗ 无法生成目录（没有找到标题）"
        exit 1
    }
    
    $replaced = Replace-TOC -FilePath $TargetFile -NewTOC $newTOC
    
    if ($replaced) {
        Write-Output "  ✓ 目录已更新"
    } else {
        Write-Output "  ✗ 无法替换目录"
    }
}
elseif ($Batch) {
    # 批量模式
    Write-Output "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    Write-Output "📝 批量生成完整目录"
    Write-Output "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    Write-Output ""
    
    # 读取问题文件列表
    if (-not (Test-Path "toc-issues.json")) {
        Write-Output "❌ 找不到 toc-issues.json，请先运行扫描脚本"
        exit 1
    }
    
    $issues = Get-Content "toc-issues.json" -Raw | ConvertFrom-Json
    
    $filesToFix = @()
    $filesToFix += $issues.simplified_toc | ForEach-Object { $_.File }
    $filesToFix += $issues.incomplete_toc | ForEach-Object { $_.File }
    
    Write-Output "需要修复: $($filesToFix.Count) 个文件"
    Write-Output ""
    
    $success = 0
    $failed = 0
    
    foreach ($file in $filesToFix) {
        Write-Output "[$($success + $failed + 1)/$($filesToFix.Count)] $file"
        
        try {
            $newTOC = Generate-TOC -FilePath $file
            
            if ($null -eq $newTOC) {
                Write-Output "  ✗ 无法生成目录"
                $failed++
                continue
            }
            
            $replaced = Replace-TOC -FilePath $file -NewTOC $newTOC
            
            if ($replaced) {
                Write-Output "  ✓ 已更新"
                $success++
            } else {
                Write-Output "  ✗ 无法替换"
                $failed++
            }
        }
        catch {
            Write-Output "  ✗ 错误: $_"
            $failed++
        }
    }
    
    Write-Output ""
    Write-Output "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    Write-Output "完成！"
    Write-Output "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    Write-Output "成功: $success 个"
    Write-Output "失败: $failed 个"
    Write-Output "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
}
else {
    Write-Output "用法:"
    Write-Output "  单文件: .\generate-complete-toc.ps1 -TargetFile <文件路径>"
    Write-Output "  批量:   .\generate-complete-toc.ps1 -Batch"
}


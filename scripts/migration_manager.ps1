# PowerShell Script: 迁移管理器
# 版本: v1.0
# 日期: 2025-10-22
# 功能: 完整的迁移管理、日志记录、错误追踪

param(
    [string]$ConfigFile = "migration-mapping.json",
    [string]$SourceDir = "docs",
    [string]$TargetDir = "docs-new",
    [string]$LogDir = "logs",
    [ValidateSet("migrate", "validate", "rollback", "status")]
    [string]$Action = "migrate",
    [switch]$DryRun,
    [switch]$Verbose,
    [switch]$Force
)

# 创建日志目录
$timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
$logFile = Join-Path $LogDir "migration-$timestamp.log"
New-Item -ItemType Directory -Path $LogDir -Force | Out-Null

# 日志函数
function Write-Log {
    param(
        [string]$Message,
        [ValidateSet("INFO", "WARN", "ERROR", "SUCCESS", "DEBUG")]
        [string]$Level = "INFO"
    )
    
    $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    $logMessage = "[$timestamp] [$Level] $Message"
    
    # 写入文件
    Add-Content -Path $logFile -Value $logMessage
    
    # 控制台输出
    switch ($Level) {
        "INFO"    { Write-Host $Message -ForegroundColor White }
        "WARN"    { Write-Host "⚠️  $Message" -ForegroundColor Yellow }
        "ERROR"   { Write-Host "❌ $Message" -ForegroundColor Red }
        "SUCCESS" { Write-Host "✅ $Message" -ForegroundColor Green }
        "DEBUG"   { if ($Verbose) { Write-Host "🔍 $Message" -ForegroundColor Gray } }
    }
}

# Banner
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  迁移管理器 v1.0" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

Write-Log "迁移管理器启动" "INFO"
Write-Log "操作模式: $Action" "INFO"
Write-Log "配置文件: $ConfigFile" "INFO"
Write-Log "源目录: $SourceDir" "INFO"
Write-Log "目标目录: $TargetDir" "INFO"
Write-Log "日志文件: $logFile" "INFO"

if ($DryRun) {
    Write-Log "DryRun模式启用" "WARN"
}

# 加载配置
Write-Log "加载配置文件..." "INFO"

if (!(Test-Path $ConfigFile)) {
    Write-Log "配置文件不存在: $ConfigFile" "ERROR"
    exit 1
}

try {
    $config = Get-Content $ConfigFile -Raw | ConvertFrom-Json
    Write-Log "配置加载成功 (版本: $($config.version))" "SUCCESS"
} catch {
    Write-Log "配置文件解析失败: $_" "ERROR"
    exit 1
}

# 迁移状态追踪
$migrationState = @{
    TotalFiles = 0
    Migrated = 0
    Skipped = 0
    Failed = 0
    Errors = @()
    Warnings = @()
    StartTime = Get-Date
}

# 函数: 迁移单个文件
function Migrate-SingleFile {
    param(
        [string]$SourceFile,
        [string]$TargetFile,
        [string]$Module
    )
    
    Write-Log "处理: $SourceFile -> $TargetFile" "DEBUG"
    
    try {
        $migrationState.TotalFiles++
        
        # 检查源文件
        if (!(Test-Path $SourceFile)) {
            Write-Log "源文件不存在: $SourceFile" "WARN"
            $migrationState.Skipped++
            $migrationState.Warnings += "源文件不存在: $SourceFile"
            return $false
        }
        
        # DryRun模式
        if ($DryRun) {
            Write-Log "[DryRun] 将迁移: $SourceFile" "DEBUG"
            $migrationState.Migrated++
            return $true
        }
        
        # 创建目标目录
        $targetDir = Split-Path $TargetFile -Parent
        if (!(Test-Path $targetDir)) {
            New-Item -ItemType Directory -Path $targetDir -Force | Out-Null
            Write-Log "创建目录: $targetDir" "DEBUG"
        }
        
        # 检查目标文件是否已存在
        if ((Test-Path $TargetFile) -and !$Force) {
            Write-Log "目标文件已存在，跳过: $TargetFile" "WARN"
            $migrationState.Skipped++
            return $true
        }
        
        # 读取源文件内容
        $content = Get-Content $SourceFile -Raw -Encoding UTF8
        
        # 应用链接替换
        $modified = $false
        foreach ($replacement in $config.linkReplacements) {
            if ($content -match [regex]::Escape($replacement.from)) {
                $content = $content -replace [regex]::Escape($replacement.from), $replacement.to
                $modified = $true
                Write-Log "应用链接替换: $($replacement.description)" "DEBUG"
            }
        }
        
        # 写入目标文件
        $content | Out-File -FilePath $TargetFile -Encoding UTF8 -NoNewline
        
        $migrationState.Migrated++
        
        if ($modified) {
            Write-Log "迁移完成(已修改链接): $TargetFile" "SUCCESS"
        } else {
            Write-Log "迁移完成: $TargetFile" "DEBUG"
        }
        
        return $true
        
    } catch {
        Write-Log "迁移失败: $SourceFile - $_" "ERROR"
        $migrationState.Failed++
        $migrationState.Errors += @{
            File = $SourceFile
            Error = $_.Exception.Message
            Time = Get-Date
        }
        return $false
    }
}

# 函数: 迁移模块
function Migrate-Module {
    param(
        [string]$ModuleName,
        [object]$ModuleConfig
    )
    
    Write-Host ""
    Write-Log "========================================" "INFO"
    Write-Log "迁移模块: $ModuleName" "INFO"
    Write-Log "目标: $($ModuleConfig.target)" "INFO"
    Write-Log "操作: $($ModuleConfig.action)" "INFO"
    Write-Log "优先级: $($ModuleConfig.priority)" "INFO"
    Write-Log "========================================" "INFO"
    
    $sourcePath = Join-Path $SourceDir $ModuleName
    $targetPath = Join-Path $TargetDir $ModuleConfig.target
    
    # 检查源目录
    if (!(Test-Path $sourcePath)) {
        Write-Log "源目录不存在，跳过: $sourcePath" "WARN"
        return
    }
    
    # 获取所有.md文件
    $mdFiles = Get-ChildItem -Path $sourcePath -Filter "*.md" -Recurse
    Write-Log "找到 $($mdFiles.Count) 个文件" "INFO"
    
    # 迁移文件
    foreach ($file in $mdFiles) {
        $relativePath = $file.FullName.Replace($sourcePath, "").TrimStart('\')
        $targetFile = Join-Path $targetPath $relativePath
        
        Migrate-SingleFile -SourceFile $file.FullName -TargetFile $targetFile -Module $ModuleName
    }
    
    Write-Log "模块迁移完成: $ModuleName ($($mdFiles.Count) 文件)" "SUCCESS"
}

# 函数: 验证迁移
function Validate-Migration {
    Write-Log "开始验证迁移..." "INFO"
    
    $issues = @()
    
    # 检查目标目录结构
    foreach ($module in $config.mappings.PSObject.Properties) {
        $moduleName = $module.Name
        $targetModule = $module.Value.target
        $targetPath = Join-Path $TargetDir $targetModule
        
        if (!(Test-Path $targetPath)) {
            $issues += "目标模块目录不存在: $targetPath"
            Write-Log "验证失败: 目录不存在 - $targetPath" "ERROR"
        }
    }
    
    # 检查README
    $readmeFiles = Get-ChildItem -Path $TargetDir -Filter "README.md" -Recurse
    Write-Log "找到 $($readmeFiles.Count) 个README文件" "INFO"
    
    # 检查链接
    Write-Log "检查链接有效性..." "INFO"
    $brokenLinks = 0
    
    $allMdFiles = Get-ChildItem -Path $TargetDir -Filter "*.md" -Recurse
    foreach ($file in $allMdFiles) {
        $content = Get-Content $file.FullName -Raw -Encoding UTF8
        $links = [regex]::Matches($content, '\[([^\]]+)\]\(([^\)]+)\)')
        
        foreach ($link in $links) {
            $linkUrl = $link.Groups[2].Value
            if ($linkUrl -match "^\.\.?/" -and $linkUrl -notmatch "^#") {
                $targetUrl = $linkUrl -replace '#.*$', ''
                $targetPath = Join-Path (Split-Path $file.FullName) $targetUrl
                $targetPath = [System.IO.Path]::GetFullPath($targetPath)
                
                if (!(Test-Path $targetPath)) {
                    $brokenLinks++
                    Write-Log "失效链接: $($file.Name) -> $linkUrl" "WARN"
                }
            }
        }
    }
    
    Write-Log "链接检查完成，发现 $brokenLinks 个失效链接" "INFO"
    
    if ($issues.Count -eq 0 -and $brokenLinks -eq 0) {
        Write-Log "验证通过！" "SUCCESS"
        return $true
    } else {
        Write-Log "验证发现 $($issues.Count) 个问题和 $brokenLinks 个失效链接" "WARN"
        return $false
    }
}

# 函数: 生成状态报告
function Generate-StatusReport {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "  迁移状态报告" -ForegroundColor Cyan
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host ""
    
    $duration = (Get-Date) - $migrationState.StartTime
    
    Write-Host "执行时间: $($duration.ToString('hh\:mm\:ss'))" -ForegroundColor White
    Write-Host ""
    Write-Host "文件统计:" -ForegroundColor Yellow
    Write-Host "  总文件数: $($migrationState.TotalFiles)" -ForegroundColor White
    Write-Host "  已迁移:   $($migrationState.Migrated)" -ForegroundColor Green
    Write-Host "  跳过:     $($migrationState.Skipped)" -ForegroundColor Yellow
    Write-Host "  失败:     $($migrationState.Failed)" -ForegroundColor Red
    Write-Host ""
    
    if ($migrationState.Failed -gt 0) {
        Write-Host "错误列表:" -ForegroundColor Red
        foreach ($error in $migrationState.Errors) {
            Write-Host "  ❌ $($error.File)" -ForegroundColor Red
            Write-Host "     $($error.Error)" -ForegroundColor Gray
        }
        Write-Host ""
    }
    
    if ($migrationState.Warnings.Count -gt 0) {
        Write-Host "警告列表 (前10条):" -ForegroundColor Yellow
        $migrationState.Warnings | Select-Object -First 10 | ForEach-Object {
            Write-Host "  ⚠️  $_" -ForegroundColor Yellow
        }
        Write-Host ""
    }
    
    $successRate = if ($migrationState.TotalFiles -gt 0) {
        [math]::Round(($migrationState.Migrated / $migrationState.TotalFiles) * 100, 2)
    } else { 0 }
    
    Write-Host "成功率: $successRate%" -ForegroundColor $(if ($successRate -ge 95) { "Green" } elseif ($successRate -ge 80) { "Yellow" } else { "Red" })
    Write-Host ""
    Write-Host "详细日志: $logFile" -ForegroundColor Cyan
    Write-Host "========================================" -ForegroundColor Cyan
}

# 主执行逻辑
try {
    switch ($Action) {
        "migrate" {
            Write-Log "开始迁移操作..." "INFO"
            
            # 按优先级排序模块
            $sortedModules = $config.mappings.PSObject.Properties | 
                Sort-Object { $_.Value.priority } -Descending
            
            foreach ($module in $sortedModules) {
                Migrate-Module -ModuleName $module.Name -ModuleConfig $module.Value
            }
            
            Write-Log "所有模块迁移完成" "SUCCESS"
            
            # 自动验证
            if (!$DryRun) {
                Write-Host ""
                Validate-Migration
            }
        }
        
        "validate" {
            Validate-Migration
        }
        
        "rollback" {
            Write-Log "回滚功能开发中..." "WARN"
            # TODO: 实现回滚功能
        }
        
        "status" {
            Write-Log "显示迁移状态..." "INFO"
            # TODO: 读取上次迁移状态
        }
    }
    
    # 生成报告
    Generate-StatusReport
    
    Write-Log "迁移管理器完成" "SUCCESS"
    
} catch {
    Write-Log "迁移过程发生严重错误: $_" "ERROR"
    Write-Log "堆栈跟踪: $($_.ScriptStackTrace)" "ERROR"
    exit 1
}


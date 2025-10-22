# 安全扫描脚本
# 用于自动化运行安全检查

param(
    [Parameter(HelpMessage="输出格式: text, json, sarif")]
    [ValidateSet("text", "json", "sarif")]
    [string]$Format = "text",
    
    [Parameter(HelpMessage="是否修复可自动修复的问题")]
    [switch]$AutoFix,
    
    [Parameter(HelpMessage="扫描的模块路径")]
    [string]$Path = "./..."
)

Write-Host ""
Write-Host "╔════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║                                                            ║"
Write-Host "║        🔒 安全扫描工具 - Security Scanner 🔒              ║" -ForegroundColor Yellow
Write-Host "║                                                            ║"
Write-Host "╚════════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

# 检查工具是否安装
function Test-CommandExists {
    param($Command)
    $oldPreference = $ErrorActionPreference
    $ErrorActionPreference = 'stop'
    try {
        if (Get-Command $Command) {
            return $true
        }
    } catch {
        return $false
    } finally {
        $ErrorActionPreference = $oldPreference
    }
}

# 安装工具
function Install-SecurityTools {
    Write-Host "📦 检查安全工具..." -ForegroundColor Cyan
    
    if (-not (Test-CommandExists "govulncheck")) {
        Write-Host "  Installing govulncheck..." -ForegroundColor Yellow
        go install golang.org/x/vuln/cmd/govulncheck@latest
    } else {
        Write-Host "  ✅ govulncheck 已安装" -ForegroundColor Green
    }
    
    if (-not (Test-CommandExists "gosec")) {
        Write-Host "  Installing gosec..." -ForegroundColor Yellow
        go install github.com/securego/gosec/v2/cmd/gosec@latest
    } else {
        Write-Host "  ✅ gosec 已安装" -ForegroundColor Green
    }
    
    Write-Host ""
}

# 运行govulncheck
function Invoke-VulnCheck {
    Write-Host "🔍 Step 1: 运行漏洞扫描 (govulncheck)..." -ForegroundColor Cyan
    Write-Host ""
    
    $result = govulncheck $Path
    $exitCode = $LASTEXITCODE
    
    if ($exitCode -eq 0) {
        Write-Host "  ✅ 未发现CVE漏洞" -ForegroundColor Green
    } else {
        Write-Host "  ⚠️  发现安全漏洞！" -ForegroundColor Red
    }
    
    Write-Host ""
    return $exitCode
}

# 运行gosec
function Invoke-GosecScan {
    Write-Host "🔍 Step 2: 运行代码安全分析 (gosec)..." -ForegroundColor Cyan
    Write-Host ""
    
    $outputFile = "security-report-$(Get-Date -Format 'yyyyMMdd-HHmmss').$Format"
    
    $args = @("-fmt", $Format)
    
    if ($Format -ne "text") {
        $args += @("-out", $outputFile)
    }
    
    $args += $Path
    
    try {
        & gosec @args
        $exitCode = $LASTEXITCODE
        
        if ($exitCode -eq 0) {
            Write-Host "  ✅ 未发现安全问题" -ForegroundColor Green
        } else {
            Write-Host "  ⚠️  发现安全问题，详见报告" -ForegroundColor Yellow
            if ($Format -ne "text") {
                Write-Host "  📄 报告已保存到: $outputFile" -ForegroundColor Cyan
            }
        }
    } catch {
        Write-Host "  ❌ gosec 扫描失败: $_" -ForegroundColor Red
        $exitCode = 1
    }
    
    Write-Host ""
    return $exitCode
}

# 扫描各个模块
function Invoke-ModuleScan {
    Write-Host "🔍 Step 3: 分模块扫描..." -ForegroundColor Cyan
    Write-Host ""
    
    $modules = @(
        "pkg/agent",
        "pkg/concurrency",
        "pkg/http3",
        "pkg/memory",
        "pkg/observability"
    )
    
    $results = @{}
    
    foreach ($module in $modules) {
        if (Test-Path $module) {
            Write-Host "  扫描 $module..." -ForegroundColor White
            
            $output = & gosec -fmt json -quiet "./$module/..." 2>&1 | ConvertFrom-Json
            $issueCount = $output.Stats.found
            
            $results[$module] = $issueCount
            
            if ($issueCount -eq 0) {
                Write-Host "    ✅ 安全" -ForegroundColor Green
            } else {
                Write-Host "    ⚠️  $issueCount 个问题" -ForegroundColor Yellow
            }
        }
    }
    
    Write-Host ""
    
    # 总结
    Write-Host "📊 模块扫描总结:" -ForegroundColor Cyan
    Write-Host ""
    
    $totalIssues = 0
    foreach ($module in $results.Keys | Sort-Object) {
        $count = $results[$module]
        $totalIssues += $count
        
        $status = if ($count -eq 0) { "✅ 安全" } else { "⚠️  $count 个问题" }
        $color = if ($count -eq 0) { "Green" } else { "Yellow" }
        
        Write-Host "  $module : $status" -ForegroundColor $color
    }
    
    Write-Host ""
    Write-Host "  总计: $totalIssues 个问题" -ForegroundColor $(if ($totalIssues -eq 0) { "Green" } else { "Yellow" })
    Write-Host ""
}

# 生成安全评分
function Get-SecurityScore {
    param($VulnCheckResult, $GosecResult)
    
    Write-Host "🏆 安全评分:" -ForegroundColor Cyan
    Write-Host ""
    
    # 计算评分
    $score = 100
    
    # CVE漏洞扣分
    if ($VulnCheckResult -ne 0) {
        $score -= 30
        Write-Host "  ❌ CVE漏洞: -30分" -ForegroundColor Red
    } else {
        Write-Host "  ✅ CVE漏洞: 满分 (100/100)" -ForegroundColor Green
    }
    
    # gosec问题扣分
    if ($GosecResult -ne 0) {
        $score -= 15
        Write-Host "  ⚠️  代码安全: -15分" -ForegroundColor Yellow
    } else {
        Write-Host "  ✅ 代码安全: 满分 (100/100)" -ForegroundColor Green
    }
    
    Write-Host ""
    
    $grade = switch ($score) {
        { $_ -ge 95 } { "A+" }
        { $_ -ge 90 } { "A" }
        { $_ -ge 85 } { "B+" }
        { $_ -ge 80 } { "B" }
        { $_ -ge 75 } { "C+" }
        { $_ -ge 70 } { "C" }
        default { "D" }
    }
    
    $color = switch ($grade) {
        "A+" { "Green" }
        "A" { "Green" }
        "B+" { "Cyan" }
        "B" { "Cyan" }
        default { "Yellow" }
    }
    
    Write-Host "  综合评分: $score/100 (等级: $grade)" -ForegroundColor $color
    Write-Host ""
}

# 主函数
function Main {
    $startTime = Get-Date
    
    # 安装工具
    Install-SecurityTools
    
    # 运行扫描
    $vulnResult = Invoke-VulnCheck
    $gosecResult = Invoke-GosecScan
    
    # 分模块扫描
    Invoke-ModuleScan
    
    # 生成评分
    Get-SecurityScore -VulnCheckResult $vulnResult -GosecResult $gosecResult
    
    $duration = (Get-Date) - $startTime
    
    Write-Host "╔════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
    Write-Host "║                                                            ║"
    Write-Host "║        ✅ 安全扫描完成！                                   ║" -ForegroundColor Green
    Write-Host "║                                                            ║"
    Write-Host "║  耗时: $([math]::Round($duration.TotalSeconds, 2)) 秒                                              ║" -ForegroundColor White
    Write-Host "║                                                            ║"
    Write-Host "╚════════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
    Write-Host ""
    
    # 返回退出码
    if ($vulnResult -ne 0 -or $gosecResult -ne 0) {
        exit 1
    }
    exit 0
}

# 执行
Main


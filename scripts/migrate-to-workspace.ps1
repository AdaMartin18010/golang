# Go 1.25.3 Workspace 迁移脚本
# 用途：将现有项目重构为 workspace 模式，分离代码和文档

param(
    [switch]$DryRun,  # 仅预览，不实际执行
    [switch]$Force    # 强制执行，覆盖已存在的文件
)

$ErrorActionPreference = "Stop"

Write-Host "================================================" -ForegroundColor Cyan
Write-Host "🚀 Go 1.25.3 Workspace 迁移工具" -ForegroundColor Cyan
Write-Host "================================================" -ForegroundColor Cyan
Write-Host ""

# 检查 Go 版本
function Check-GoVersion {
    Write-Host "🔍 检查 Go 版本..." -ForegroundColor Yellow
    
    try {
        $goVersion = go version
        Write-Host "   当前版本: $goVersion" -ForegroundColor Green
        
        if ($goVersion -notmatch "go1\.2[5-9]|go1\.[3-9][0-9]") {
            Write-Host "   ⚠️  推荐使用 Go 1.25.3 或更高版本" -ForegroundColor Yellow
            $continue = Read-Host "   是否继续? (y/n)"
            if ($continue -ne "y") {
                exit 0
            }
        }
    } catch {
        Write-Host "   ❌ 未找到 Go 工具链" -ForegroundColor Red
        exit 1
    }
}

# 创建新的目录结构
function Create-NewStructure {
    Write-Host ""
    Write-Host "📁 创建新的目录结构..." -ForegroundColor Yellow
    
    $dirs = @(
        "cmd",
        "pkg/agent/core",
        "pkg/agent/coordination",
        "pkg/concurrency/pipeline",
        "pkg/concurrency/workerpool",
        "pkg/http3/server",
        "pkg/memory/arena",
        "pkg/memory/weakptr",
        "pkg/observability/metrics",
        "pkg/observability/tracing",
        "pkg/observability/logging",
        "internal/config",
        "internal/utils",
        "internal/testutil",
        "tests/integration",
        "tests/e2e",
        "tests/benchmarks",
        "deployments/docker",
        "deployments/kubernetes",
        "reports/phase-reports",
        "reports/code-quality",
        "reports/archive"
    )
    
    foreach ($dir in $dirs) {
        $path = Join-Path $PSScriptRoot "..\$dir"
        if (!(Test-Path $path)) {
            if (!$DryRun) {
                New-Item -ItemType Directory -Path $path -Force | Out-Null
                Write-Host "   ✅ 创建: $dir" -ForegroundColor Green
            } else {
                Write-Host "   [预览] 将创建: $dir" -ForegroundColor Gray
            }
        } else {
            Write-Host "   ⏭️  已存在: $dir" -ForegroundColor Gray
        }
    }
}

# 初始化各模块的 go.mod
function Initialize-Modules {
    Write-Host ""
    Write-Host "📦 初始化 Go 模块..." -ForegroundColor Yellow
    
    $modules = @{
        "pkg/agent" = "github.com/yourusername/agent"
        "pkg/concurrency" = "github.com/yourusername/concurrency"
        "pkg/http3" = "github.com/yourusername/http3"
        "pkg/memory" = "github.com/yourusername/memory"
        "pkg/observability" = "github.com/yourusername/observability"
    }
    
    foreach ($module in $modules.GetEnumerator()) {
        $dir = Join-Path $PSScriptRoot "..\$($module.Key)"
        $modPath = Join-Path $dir "go.mod"
        
        if (!(Test-Path $modPath) -or $Force) {
            if (!$DryRun) {
                Push-Location $dir
                go mod init $module.Value 2>&1 | Out-Null
                go mod edit -go=1.25.3
                Pop-Location
                Write-Host "   ✅ 初始化: $($module.Key)" -ForegroundColor Green
            } else {
                Write-Host "   [预览] 将初始化: $($module.Key)" -ForegroundColor Gray
            }
        } else {
            Write-Host "   ⏭️  已存在: $($module.Key)/go.mod" -ForegroundColor Gray
        }
    }
}

# 迁移 AI Agent 代码
function Migrate-AIAgent {
    Write-Host ""
    Write-Host "🤖 迁移 AI Agent 代码..." -ForegroundColor Yellow
    
    $sourceBase = Join-Path $PSScriptRoot "..\examples\advanced\ai-agent"
    $targetPkg = Join-Path $PSScriptRoot "..\pkg\agent"
    $targetCmd = Join-Path $PSScriptRoot "..\cmd\ai-agent"
    
    if (Test-Path $sourceBase) {
        # 迁移核心代码
        $coreFiles = @("core", "coordination")
        foreach ($item in $coreFiles) {
            $src = Join-Path $sourceBase $item
            $dst = Join-Path $targetPkg $item
            
            if (Test-Path $src) {
                if (!$DryRun) {
                    Copy-Item -Path $src -Destination $dst -Recurse -Force
                    Write-Host "   ✅ 迁移: $item -> pkg/agent/$item" -ForegroundColor Green
                } else {
                    Write-Host "   [预览] 将迁移: $item -> pkg/agent/$item" -ForegroundColor Gray
                }
            }
        }
        
        # 迁移 main.go
        $mainFile = Join-Path $sourceBase "main.go"
        if (Test-Path $mainFile) {
            if (!$DryRun) {
                New-Item -ItemType Directory -Path $targetCmd -Force | Out-Null
                Copy-Item -Path $mainFile -Destination (Join-Path $targetCmd "main.go") -Force
                Write-Host "   ✅ 迁移: main.go -> cmd/ai-agent/" -ForegroundColor Green
            } else {
                Write-Host "   [预览] 将迁移: main.go -> cmd/ai-agent/" -ForegroundColor Gray
            }
        }
    } else {
        Write-Host "   ⚠️  未找到 AI Agent 源代码" -ForegroundColor Yellow
    }
}

# 迁移并发代码
function Migrate-Concurrency {
    Write-Host ""
    Write-Host "⚡ 迁移并发代码..." -ForegroundColor Yellow
    
    $sourceBase = Join-Path $PSScriptRoot "..\examples\concurrency"
    $targetPkg = Join-Path $PSScriptRoot "..\pkg\concurrency"
    
    if (Test-Path $sourceBase) {
        # 迁移 pipeline
        $pipelineTest = Join-Path $sourceBase "pipeline_test.go"
        if (Test-Path $pipelineTest) {
            if (!$DryRun) {
                $dst = Join-Path $targetPkg "pipeline"
                New-Item -ItemType Directory -Path $dst -Force | Out-Null
                Copy-Item -Path $pipelineTest -Destination (Join-Path $dst "pipeline_test.go") -Force
                Write-Host "   ✅ 迁移: pipeline_test.go -> pkg/concurrency/pipeline/" -ForegroundColor Green
            } else {
                Write-Host "   [预览] 将迁移: pipeline_test.go" -ForegroundColor Gray
            }
        }
        
        # 迁移 worker pool
        $workerTest = Join-Path $sourceBase "worker_pool_test.go"
        if (Test-Path $workerTest) {
            if (!$DryRun) {
                $dst = Join-Path $targetPkg "workerpool"
                New-Item -ItemType Directory -Path $dst -Force | Out-Null
                Copy-Item -Path $workerTest -Destination (Join-Path $dst "workerpool_test.go") -Force
                Write-Host "   ✅ 迁移: worker_pool_test.go -> pkg/concurrency/workerpool/" -ForegroundColor Green
            } else {
                Write-Host "   [预览] 将迁移: worker_pool_test.go" -ForegroundColor Gray
            }
        }
    } else {
        Write-Host "   ⚠️  未找到并发代码源文件" -ForegroundColor Yellow
    }
}

# 整理文档
function Organize-Docs {
    Write-Host ""
    Write-Host "📚 整理文档..." -ForegroundColor Yellow
    
    # 移动报告文件
    $reportFiles = Get-ChildItem -Path (Join-Path $PSScriptRoot "..") -Filter "*报告*.md"
    $reportFiles += Get-ChildItem -Path (Join-Path $PSScriptRoot "..") -Filter "Phase*.md"
    
    foreach ($file in $reportFiles) {
        $dst = Join-Path $PSScriptRoot "..\reports\phase-reports\$($file.Name)"
        if (!$DryRun) {
            Move-Item -Path $file.FullName -Destination $dst -Force -ErrorAction SilentlyContinue
            Write-Host "   ✅ 移动: $($file.Name) -> reports/phase-reports/" -ForegroundColor Green
        } else {
            Write-Host "   [预览] 将移动: $($file.Name)" -ForegroundColor Gray
        }
    }
    
    # 提示合并 docs 和 docs-new
    Write-Host ""
    Write-Host "   ⚠️  注意: 需要手动合并 docs/ 和 docs-new/ 目录" -ForegroundColor Yellow
    Write-Host "   建议: 选择结构更好的一个作为主目录，删除另一个" -ForegroundColor Yellow
}

# 更新 examples/go.mod
function Update-ExamplesModule {
    Write-Host ""
    Write-Host "📝 更新 examples/go.mod..." -ForegroundColor Yellow
    
    $examplesDir = Join-Path $PSScriptRoot "..\examples"
    $modFile = Join-Path $examplesDir "go.mod"
    
    if (Test-Path $modFile) {
        if (!$DryRun) {
            Push-Location $examplesDir
            go mod edit -go=1.25.3
            go mod tidy
            Pop-Location
            Write-Host "   ✅ 更新 examples/go.mod 到 Go 1.25.3" -ForegroundColor Green
        } else {
            Write-Host "   [预览] 将更新 examples/go.mod" -ForegroundColor Gray
        }
    }
}

# 验证工作区
function Verify-Workspace {
    Write-Host ""
    Write-Host "🔍 验证工作区配置..." -ForegroundColor Yellow
    
    $workFile = Join-Path $PSScriptRoot "..\go.work"
    
    if (Test-Path $workFile) {
        if (!$DryRun) {
            Push-Location (Join-Path $PSScriptRoot "..")
            
            Write-Host "   检查工作区同步..." -ForegroundColor Gray
            go work sync
            
            Write-Host "   运行测试..." -ForegroundColor Gray
            $testResult = go work test ./examples/... 2>&1
            
            if ($LASTEXITCODE -eq 0) {
                Write-Host "   ✅ 工作区验证通过" -ForegroundColor Green
            } else {
                Write-Host "   ⚠️  部分测试失败，请检查" -ForegroundColor Yellow
            }
            
            Pop-Location
        } else {
            Write-Host "   [预览] 将验证工作区" -ForegroundColor Gray
        }
    } else {
        Write-Host "   ⚠️  未找到 go.work 文件" -ForegroundColor Yellow
        Write-Host "   请先创建 go.work 文件" -ForegroundColor Yellow
    }
}

# 生成迁移报告
function Generate-Report {
    Write-Host ""
    Write-Host "📊 生成迁移报告..." -ForegroundColor Yellow
    
    $reportPath = Join-Path $PSScriptRoot "..\reports\MIGRATION_REPORT_$(Get-Date -Format 'yyyy-MM-dd').md"
    
    $report = @"
# Go 1.25.3 Workspace 迁移报告

**日期**: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')

## 迁移概述

本次迁移将项目重构为 Go 1.25.3 workspace 模式，实现代码与文档的清晰分离。

## 执行的操作

### ✅ 已完成

1. 创建新的目录结构
2. 初始化 Go 模块
3. 迁移 AI Agent 代码
4. 迁移并发代码
5. 整理文档和报告
6. 更新 examples/go.mod

### ⏳ 待完成

1. 手动合并 docs/ 和 docs-new/ 目录
2. 更新所有文档中的代码路径引用
3. 更新 CI/CD 配置
4. 更新 README.md

## 新的项目结构

\`\`\`text
golang/
├── go.work              # Workspace 配置
├── cmd/                 # 可执行程序
├── pkg/                 # 可复用库
├── examples/            # 示例代码
├── internal/            # 内部包
├── docs/                # 文档
└── reports/             # 项目报告
\`\`\`

## 下一步

1. 运行 \`go work sync\` 同步依赖
2. 运行 \`go work test ./...\` 验证所有测试
3. 更新 README.md 和导航文档
4. 提交更改到版本控制

## 参考

- [Go Workspace 文档](https://go.dev/doc/tutorial/workspaces)
- [项目重构方案](../RESTRUCTURE_PROPOSAL_GO1.25.3.md)
"@

    if (!$DryRun) {
        $report | Out-File -FilePath $reportPath -Encoding UTF8
        Write-Host "   ✅ 报告已生成: $reportPath" -ForegroundColor Green
    } else {
        Write-Host "   [预览] 将生成报告" -ForegroundColor Gray
    }
}

# 主流程
function Main {
    Write-Host "执行模式: $(if ($DryRun) { '🔍 预览模式（不会实际修改文件）' } else { '⚙️  执行模式' })" -ForegroundColor Cyan
    Write-Host ""
    
    if (!$DryRun) {
        $confirm = Read-Host "确认开始迁移? 建议先运行 -DryRun 预览 (y/n)"
        if ($confirm -ne "y") {
            Write-Host "已取消" -ForegroundColor Yellow
            exit 0
        }
    }
    
    Check-GoVersion
    Create-NewStructure
    Initialize-Modules
    Migrate-AIAgent
    Migrate-Concurrency
    Organize-Docs
    Update-ExamplesModule
    
    if (!$DryRun) {
        Verify-Workspace
        Generate-Report
    }
    
    Write-Host ""
    Write-Host "================================================" -ForegroundColor Cyan
    Write-Host "🎉 迁移完成！" -ForegroundColor Green
    Write-Host "================================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "下一步:" -ForegroundColor Yellow
    Write-Host "  1. 运行: go work sync" -ForegroundColor White
    Write-Host "  2. 测试: go work test ./..." -ForegroundColor White
    Write-Host "  3. 查看: reports/MIGRATION_REPORT_*.md" -ForegroundColor White
    Write-Host ""
}

Main


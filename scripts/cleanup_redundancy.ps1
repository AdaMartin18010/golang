# ============================================
# 项目冗余文件清理脚本
# ============================================
# 作者: AI 助手
# 日期: 2025-10-25
# 用途: 清理项目中的冗余文件
# 警告: 使用前请先备份！
# ============================================

param(
    [switch]$DryRun = $false,  # 试运行模式，不实际删除
    [switch]$Stage1 = $false,  # 只执行阶段1
    [switch]$Stage2 = $false,  # 只执行阶段2
    [switch]$Stage3 = $false,  # 只执行阶段3
    [switch]$All = $false      # 执行所有阶段
)

$ErrorActionPreference = "Stop"

# 颜色输出函数
function Write-ColorOutput {
    param([string]$Message, [string]$Color = "White")
    Write-Host $Message -ForegroundColor $Color
}

# 确认函数
function Confirm-Action {
    param([string]$Message)
    $response = Read-Host "$Message (y/N)"
    return $response -eq 'y' -or $response -eq 'Y'
}

# 删除文件/目录函数
function Remove-SafeItem {
    param(
        [string]$Path,
        [string]$Description
    )
    
    if (Test-Path $Path) {
        if ($DryRun) {
            Write-ColorOutput "  [DRY RUN] 将删除: $Path" "Yellow"
        } else {
            Write-ColorOutput "  删除: $Description" "Gray"
            Remove-Item -Path $Path -Recurse -Force
            Write-ColorOutput "  ✓ 已删除" "Green"
        }
    } else {
        Write-ColorOutput "  跳过（不存在）: $Path" "DarkGray"
    }
}

# 主函数
function Main {
    Write-ColorOutput "================================================" "Cyan"
    Write-ColorOutput "  Go项目冗余文件清理脚本" "Cyan"
    Write-ColorOutput "================================================" "Cyan"
    Write-ColorOutput ""
    
    if ($DryRun) {
        Write-ColorOutput "⚠️  试运行模式 - 不会实际删除文件" "Yellow"
        Write-ColorOutput ""
    }
    
    # 检查是否在项目根目录
    if (-not (Test-Path "go.work") -and -not (Test-Path "README.md")) {
        Write-ColorOutput "❌ 错误：请在项目根目录运行此脚本！" "Red"
        exit 1
    }
    
    # 如果没有指定阶段，询问用户
    if (-not $Stage1 -and -not $Stage2 -and -not $Stage3 -and -not $All) {
        Write-ColorOutput "请选择要执行的清理阶段：" "Cyan"
        Write-ColorOutput "  1. 阶段1 - 删除 docs_old 和旧状态文件（节省最多空间）" "White"
        Write-ColorOutput "  2. 阶段2 - 合并重复报告和文档" "White"
        Write-ColorOutput "  3. 阶段3 - 全部清理" "White"
        Write-ColorOutput ""
        $choice = Read-Host "请输入选项 (1/2/3)"
        
        switch ($choice) {
            "1" { $Stage1 = $true }
            "2" { $Stage2 = $true }
            "3" { $All = $true }
            default {
                Write-ColorOutput "❌ 无效选项，退出" "Red"
                exit 1
            }
        }
    }
    
    # ============================================
    # 阶段1：删除明确的冗余内容
    # ============================================
    if ($Stage1 -or $All) {
        Write-ColorOutput ""
        Write-ColorOutput "=== 阶段1：删除明确的冗余内容 ===" "Cyan"
        Write-ColorOutput ""
        
        if (-not $DryRun) {
            if (-not (Confirm-Action "⚠️  即将删除大量文件！是否继续？")) {
                Write-ColorOutput "已取消" "Yellow"
                exit 0
            }
        }
        
        # 1.1 备份并删除 docs_old
        Write-ColorOutput "1. 处理 docs_old/ 目录..." "Yellow"
        if (Test-Path "docs_old") {
            if (-not $DryRun) {
                Write-ColorOutput "  创建备份..." "Gray"
                $backupFile = "archive/docs_old_backup_2025-10-25.zip"
                if (-not (Test-Path "archive")) {
                    New-Item -ItemType Directory -Path "archive" -Force | Out-Null
                }
                if (-not (Test-Path $backupFile)) {
                    Compress-Archive -Path "docs_old" -DestinationPath $backupFile -CompressionLevel Fastest
                    Write-ColorOutput "  ✓ 备份已创建: $backupFile" "Green"
                }
            }
            Remove-SafeItem "docs_old" "docs_old/ 目录 (1,428个文件)"
        }
        
        # 1.2 清理根目录状态文件（保留最新）
        Write-ColorOutput ""
        Write-ColorOutput "2. 清理状态文件..." "Yellow"
        $statusFiles = @(
            "📍-当前状态-2025-10-22.md",
            "📍-Go形式化理论体系-最新状态-2025-10-23.md",
            "📍-Go形式化理论体系-最新状态-2025-10-23-Update2.md",
            "📍-Go形式化理论体系-最新状态-2025-10-23-Update3.md",
            "📍-Go形式化理论体系-最新状态-2025-10-23-Update4.md",
            "📍-Go形式化理论体系-最新状态-2025-10-23-Update5.md",
            "📍-Go形式化理论体系-最新状态-2025-10-23-Update6.md",
            "📍-Go形式化理论体系-最新状态-2025-10-23-Update7.md",
            "📍-Go形式化理论体系-最新状态-2025-10-23-Update8.md",
            "📍-Go形式化理论体系-最终状态-2025-10-23.md",
            "📍-Go形式化理论体系-最新状态-2025-10-25-Update.md",
            "📍-Phase3-Week3-当前状态-2025-10-25.md",
            "📍-Phase4启动状态-2025-10-23.md",
            "📍-项目地图-2025-10-23.md"
        )
        
        foreach ($file in $statusFiles) {
            Remove-SafeItem $file "状态文件: $file"
        }
        Write-ColorOutput "  ℹ️  保留: 📍-Go形式化理论体系-最新状态-2025-10-25.md" "Cyan"
        
        # 1.3 删除迁移相关文件
        Write-ColorOutput ""
        Write-ColorOutput "3. 删除迁移文件（迁移已完成）..." "Yellow"
        $migrationFiles = @(
            "MIGRATION_GUIDE.md",
            "MIGRATION_GUIDE_v2.md",
            "MIGRATION_COMPARISON.md",
            "MIGRATION_CHECKLIST.md",
            "WORKSPACE_MIGRATION_PLAN.md",
            "WORKSPACE_MIGRATION_INDEX.md",
            "README_WORKSPACE_MIGRATION.md",
            "快速参考-Workspace迁移.md",
            "新旧结构对照速查.txt",
            "QUICK_START_WORKSPACE.md"
        )
        
        foreach ($file in $migrationFiles) {
            Remove-SafeItem $file "迁移文件: $file"
        }
        
        # 1.4 删除文档优化报告
        Write-ColorOutput ""
        Write-ColorOutput "4. 删除文档优化报告..." "Yellow"
        $optimizationFiles = @(
            "文档结构深度优化方案.md",
            "文档结构优化第二轮完成报告.md",
            "文档优化三轮完成总报告.md",
            "00-开始阅读-重构指南.md"
        )
        
        foreach ($file in $optimizationFiles) {
            Remove-SafeItem $file "优化报告: $file"
        }
        
        Write-ColorOutput ""
        Write-ColorOutput "✅ 阶段1完成！" "Green"
    }
    
    # ============================================
    # 阶段2：合并重复文档
    # ============================================
    if ($Stage2 -or $All) {
        Write-ColorOutput ""
        Write-ColorOutput "=== 阶段2：合并重复文档 ===" "Cyan"
        Write-ColorOutput ""
        
        # 2.1 清理完成报告（docs目录下）
        Write-ColorOutput "1. 清理完成报告..." "Yellow"
        $completionReports = @(
            "docs/🎊-持续推进完成报告-2025-10-23.md",
            "docs/🎊-Golang架构知识库深度优化总结报告-2025-10-24.md",
            "docs/🎊-2025年10月24日完成总结-知识梳理项目终章-2025-10-24.md",
            "docs/🎊-2025年10月24日Phase4推进总结-2025-10-24.md",
            "docs/🎊-2025年10月文档更新计划-100%完成-2025-10-24.md",
            "docs/🎊-2025年10月25日完成总结-Phase3-Week3推进-2025-10-25.md",
            "docs/🎊-2025年10月知识梳理项目-完整交付报告-2025-10-24.md",
            "docs/🎊-docs目录持续推进完成-2025-10-25.md",
            "docs/🎊-Go-1.25.3项目结构梳理完成-2025-10-25.md",
            "docs/🎊-第4轮持续推进圆满完成-2025-10-25.md",
            "docs/🎊-第6轮持续推进圆满完成-2025-10-25.md",
            "docs/🎊-第7轮持续推进圆满完成-2025-10-25.md",
            "docs/🎊-第8轮持续推进圆满完成-2025-10-25.md",
            "docs/🎊-第9轮持续推进圆满完成-2025-10-25.md",
            "docs/🎊-第10轮持续推进圆满完成-2025-10-25.md",
            "docs/🎊-第11轮持续推进圆满完成-2025-10-25.md",
            "docs/🎊-第12轮持续推进圆满完成-2025-10-25.md"
        )
        
        foreach ($file in $completionReports) {
            Remove-SafeItem $file "完成报告: $file"
        }
        
        # 2.2 清理执行计划
        Write-ColorOutput ""
        Write-ColorOutput "2. 清理历史执行计划..." "Yellow"
        $phaseFiles = @(
            "🚀-Phase2执行计划-2025-10-22.md",
            "🚀-Phase3执行计划.md",
            "🚀-持续推进Phase3启动报告-2025-10-23.md",
            "🚀-Phase3-Week2启动报告-2025-10-23.md",
            "🚀-Phase4-3-工具增强计划-2025-10-23.md"
        )
        
        foreach ($file in $phaseFiles) {
            Remove-SafeItem $file "执行计划: $file"
        }
        Write-ColorOutput "  ℹ️  保留: 🚀-Phase4执行计划.md 和 🚀-立即开始-3分钟上手.md" "Cyan"
        
        # 2.3 清理冗余 README
        Write-ColorOutput ""
        Write-ColorOutput "3. 清理冗余 README..." "Yellow"
        $readmeFiles = @(
            "README-项目现状-2025-10-25.md",
            "README-PROJECT-COMPLETE.md",
            "README-WORKSPACE-READY.md",
            "README-重构说明.md",
            "📖-README-项目导航.md"
        )
        
        foreach ($file in $readmeFiles) {
            Remove-SafeItem $file "README: $file"
        }
        Write-ColorOutput "  ℹ️  保留: README.md 和 README_EN.md" "Cyan"
        
        # 2.4 清理项目报告
        Write-ColorOutput ""
        Write-ColorOutput "4. 清理项目报告..." "Yellow"
        $projectReports = @(
            "PROJECT_COMPLETION_REPORT.md",
            "PROJECT_DELIVERY_CHECKLIST.md",
            "📚-项目最终完成报告-2025-10-23.md",
            "📌-项目状态总览.md",
            "📝-持续推进总结-2025-10-22.md"
        )
        
        foreach ($file in $projectReports) {
            Remove-SafeItem $file "项目报告: $file"
        }
        Write-ColorOutput "  ℹ️  保留: PROJECT_PHASES_SUMMARY.md" "Cyan"
        
        # 2.5 清理其他冗余文件
        Write-ColorOutput ""
        Write-ColorOutput "5. 清理其他冗余文件..." "Yellow"
        $otherFiles = @(
            "📚-Workspace文档索引.md",
            "📖-完整学习地图-2025-10-23.md",
            "🤝-贡献指南-CONTRIBUTING-2025.md",
            "🚀-Phase4-持续发展规划-2025-10-23.md",
            "CHART_ENHANCEMENT_SUMMARY.md"
        )
        
        foreach ($file in $otherFiles) {
            Remove-SafeItem $file "冗余文件: $file"
        }
        
        Write-ColorOutput ""
        Write-ColorOutput "✅ 阶段2完成！" "Green"
    }
    
    # ============================================
    # 阶段3：可选的深度清理
    # ============================================
    if ($All) {
        Write-ColorOutput ""
        Write-ColorOutput "=== 阶段3：深度清理（可选） ===" "Cyan"
        Write-ColorOutput ""
        
        Write-ColorOutput "⚠️  以下操作可能需要更多评估，建议手动执行：" "Yellow"
        Write-ColorOutput "  1. 评估 archive/model/ 目录（920个文件）" "Gray"
        Write-ColorOutput "  2. 清理 archive/model/Programming_Language/rust/ （非Go内容）" "Gray"
        Write-ColorOutput "  3. 压缩 archive/ 为 .tar.gz" "Gray"
        Write-ColorOutput "  4. 检查并修复文档内部链接" "Gray"
        Write-ColorOutput ""
    }
    
    # 总结
    Write-ColorOutput ""
    Write-ColorOutput "================================================" "Cyan"
    Write-ColorOutput "  清理完成总结" "Cyan"
    Write-ColorOutput "================================================" "Cyan"
    Write-ColorOutput ""
    
    if (-not $DryRun) {
        Write-ColorOutput "✅ 清理已完成！" "Green"
        Write-ColorOutput ""
        Write-ColorOutput "建议的后续步骤：" "Yellow"
        Write-ColorOutput "  1. 检查 git status" "White"
        Write-ColorOutput "  2. 验证项目功能正常" "White"
        Write-ColorOutput "  3. 更新文档链接" "White"
        Write-ColorOutput "  4. 提交更改" "White"
        Write-ColorOutput ""
        Write-ColorOutput "运行以下命令查看更改：" "Cyan"
        Write-ColorOutput "  git status" "White"
        Write-ColorOutput "  git diff --stat" "White"
    } else {
        Write-ColorOutput "✅ 试运行完成！" "Green"
        Write-ColorOutput ""
        Write-ColorOutput "要实际执行清理，请运行：" "Yellow"
        Write-ColorOutput "  .\scripts\cleanup_redundancy.ps1 -Stage1" "White"
        Write-ColorOutput "  或" "Gray"
        Write-ColorOutput "  .\scripts\cleanup_redundancy.ps1 -All" "White"
    }
    
    Write-ColorOutput ""
}

# 运行主函数
try {
    Main
} catch {
    Write-ColorOutput ""
    Write-ColorOutput "❌ 错误: $_" "Red"
    Write-ColorOutput ""
    Write-ColorOutput "清理已中止，请检查错误信息" "Yellow"
    exit 1
}


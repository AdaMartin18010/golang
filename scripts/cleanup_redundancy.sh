#!/bin/bash

# ============================================
# 项目冗余文件清理脚本 (Linux/macOS)
# ============================================
# 作者: AI 助手
# 日期: 2025-10-25
# 用途: 清理项目中的冗余文件
# 警告: 使用前请先备份！
# ============================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
GRAY='\033[0;90m'
NC='\033[0m' # No Color

# 参数
DRY_RUN=false
STAGE1=false
STAGE2=false
STAGE3=false
ALL=false

# 解析参数
while [[ $# -gt 0 ]]; do
    case $1 in
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --stage1)
            STAGE1=true
            shift
            ;;
        --stage2)
            STAGE2=true
            shift
            ;;
        --stage3)
            STAGE3=true
            shift
            ;;
        --all)
            ALL=true
            shift
            ;;
        *)
            echo -e "${RED}未知参数: $1${NC}"
            echo "用法: $0 [--dry-run] [--stage1|--stage2|--stage3|--all]"
            exit 1
            ;;
    esac
done

# 彩色输出函数
print_color() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# 确认函数
confirm_action() {
    local message=$1
    read -p "$(echo -e ${YELLOW}${message} '(y/N): '${NC})" response
    [[ "$response" == "y" || "$response" == "Y" ]]
}

# 删除函数
remove_safe_item() {
    local path=$1
    local description=$2
    
    if [ -e "$path" ]; then
        if [ "$DRY_RUN" = true ]; then
            print_color "$YELLOW" "  [DRY RUN] 将删除: $path"
        else
            print_color "$GRAY" "  删除: $description"
            rm -rf "$path"
            print_color "$GREEN" "  ✓ 已删除"
        fi
    else
        print_color "$GRAY" "  跳过（不存在）: $path"
    fi
}

# 主函数
main() {
    print_color "$CYAN" "================================================"
    print_color "$CYAN" "  Go项目冗余文件清理脚本"
    print_color "$CYAN" "================================================"
    echo ""
    
    if [ "$DRY_RUN" = true ]; then
        print_color "$YELLOW" "⚠️  试运行模式 - 不会实际删除文件"
        echo ""
    fi
    
    # 检查是否在项目根目录
    if [ ! -f "go.work" ] && [ ! -f "README.md" ]; then
        print_color "$RED" "❌ 错误：请在项目根目录运行此脚本！"
        exit 1
    fi
    
    # 如果没有指定阶段，询问用户
    if [ "$STAGE1" = false ] && [ "$STAGE2" = false ] && [ "$STAGE3" = false ] && [ "$ALL" = false ]; then
        print_color "$CYAN" "请选择要执行的清理阶段："
        echo "  1. 阶段1 - 删除 docs_old 和旧状态文件（节省最多空间）"
        echo "  2. 阶段2 - 合并重复报告和文档"
        echo "  3. 阶段3 - 全部清理"
        echo ""
        read -p "请输入选项 (1/2/3): " choice
        
        case $choice in
            1) STAGE1=true ;;
            2) STAGE2=true ;;
            3) ALL=true ;;
            *)
                print_color "$RED" "❌ 无效选项，退出"
                exit 1
                ;;
        esac
    fi
    
    # ============================================
    # 阶段1：删除明确的冗余内容
    # ============================================
    if [ "$STAGE1" = true ] || [ "$ALL" = true ]; then
        echo ""
        print_color "$CYAN" "=== 阶段1：删除明确的冗余内容 ==="
        echo ""
        
        if [ "$DRY_RUN" = false ]; then
            if ! confirm_action "⚠️  即将删除大量文件！是否继续？"; then
                print_color "$YELLOW" "已取消"
                exit 0
            fi
        fi
        
        # 1.1 备份并删除 docs_old
        print_color "$YELLOW" "1. 处理 docs_old/ 目录..."
        if [ -d "docs_old" ]; then
            if [ "$DRY_RUN" = false ]; then
                print_color "$GRAY" "  创建备份..."
                mkdir -p archive
                backup_file="archive/docs_old_backup_2025-10-25.tar.gz"
                if [ ! -f "$backup_file" ]; then
                    tar -czf "$backup_file" docs_old
                    print_color "$GREEN" "  ✓ 备份已创建: $backup_file"
                fi
            fi
            remove_safe_item "docs_old" "docs_old/ 目录 (1,428个文件)"
        fi
        
        # 1.2 清理根目录状态文件
        echo ""
        print_color "$YELLOW" "2. 清理状态文件..."
        
        status_files=(
            "📍-当前状态-2025-10-22.md"
            "📍-Go形式化理论体系-最新状态-2025-10-23.md"
            "📍-Go形式化理论体系-最新状态-2025-10-23-Update2.md"
            "📍-Go形式化理论体系-最新状态-2025-10-23-Update3.md"
            "📍-Go形式化理论体系-最新状态-2025-10-23-Update4.md"
            "📍-Go形式化理论体系-最新状态-2025-10-23-Update5.md"
            "📍-Go形式化理论体系-最新状态-2025-10-23-Update6.md"
            "📍-Go形式化理论体系-最新状态-2025-10-23-Update7.md"
            "📍-Go形式化理论体系-最新状态-2025-10-23-Update8.md"
            "📍-Go形式化理论体系-最终状态-2025-10-23.md"
            "📍-Go形式化理论体系-最新状态-2025-10-25-Update.md"
            "📍-Phase3-Week3-当前状态-2025-10-25.md"
            "📍-Phase4启动状态-2025-10-23.md"
            "📍-项目地图-2025-10-23.md"
        )
        
        for file in "${status_files[@]}"; do
            remove_safe_item "$file" "状态文件: $file"
        done
        print_color "$CYAN" "  ℹ️  保留: 📍-Go形式化理论体系-最新状态-2025-10-25.md"
        
        # 1.3 删除迁移相关文件
        echo ""
        print_color "$YELLOW" "3. 删除迁移文件（迁移已完成）..."
        
        migration_files=(
            "MIGRATION_GUIDE.md"
            "MIGRATION_GUIDE_v2.md"
            "MIGRATION_COMPARISON.md"
            "MIGRATION_CHECKLIST.md"
            "WORKSPACE_MIGRATION_PLAN.md"
            "WORKSPACE_MIGRATION_INDEX.md"
            "README_WORKSPACE_MIGRATION.md"
            "快速参考-Workspace迁移.md"
            "新旧结构对照速查.txt"
            "QUICK_START_WORKSPACE.md"
        )
        
        for file in "${migration_files[@]}"; do
            remove_safe_item "$file" "迁移文件: $file"
        done
        
        # 1.4 删除文档优化报告
        echo ""
        print_color "$YELLOW" "4. 删除文档优化报告..."
        
        optimization_files=(
            "文档结构深度优化方案.md"
            "文档结构优化第二轮完成报告.md"
            "文档优化三轮完成总报告.md"
            "00-开始阅读-重构指南.md"
        )
        
        for file in "${optimization_files[@]}"; do
            remove_safe_item "$file" "优化报告: $file"
        done
        
        echo ""
        print_color "$GREEN" "✅ 阶段1完成！"
    fi
    
    # ============================================
    # 阶段2：合并重复文档
    # ============================================
    if [ "$STAGE2" = true ] || [ "$ALL" = true ]; then
        echo ""
        print_color "$CYAN" "=== 阶段2：合并重复文档 ==="
        echo ""
        
        # 2.1 清理完成报告
        print_color "$YELLOW" "1. 清理完成报告..."
        
        completion_reports=(
            "docs/🎊-持续推进完成报告-2025-10-23.md"
            "docs/🎊-Golang架构知识库深度优化总结报告-2025-10-24.md"
            "docs/🎊-2025年10月24日完成总结-知识梳理项目终章-2025-10-24.md"
            "docs/🎊-2025年10月24日Phase4推进总结-2025-10-24.md"
            "docs/🎊-2025年10月文档更新计划-100%完成-2025-10-24.md"
            "docs/🎊-2025年10月25日完成总结-Phase3-Week3推进-2025-10-25.md"
            "docs/🎊-2025年10月知识梳理项目-完整交付报告-2025-10-24.md"
            "docs/🎊-docs目录持续推进完成-2025-10-25.md"
            "docs/🎊-Go-1.25.3项目结构梳理完成-2025-10-25.md"
            "docs/🎊-第4轮持续推进圆满完成-2025-10-25.md"
            "docs/🎊-第6轮持续推进圆满完成-2025-10-25.md"
            "docs/🎊-第7轮持续推进圆满完成-2025-10-25.md"
            "docs/🎊-第8轮持续推进圆满完成-2025-10-25.md"
            "docs/🎊-第9轮持续推进圆满完成-2025-10-25.md"
            "docs/🎊-第10轮持续推进圆满完成-2025-10-25.md"
            "docs/🎊-第11轮持续推进圆满完成-2025-10-25.md"
            "docs/🎊-第12轮持续推进圆满完成-2025-10-25.md"
        )
        
        for file in "${completion_reports[@]}"; do
            remove_safe_item "$file" "完成报告: $file"
        done
        
        # 2.2 清理执行计划
        echo ""
        print_color "$YELLOW" "2. 清理历史执行计划..."
        
        phase_files=(
            "🚀-Phase2执行计划-2025-10-22.md"
            "🚀-Phase3执行计划.md"
            "🚀-持续推进Phase3启动报告-2025-10-23.md"
            "🚀-Phase3-Week2启动报告-2025-10-23.md"
            "🚀-Phase4-3-工具增强计划-2025-10-23.md"
        )
        
        for file in "${phase_files[@]}"; do
            remove_safe_item "$file" "执行计划: $file"
        done
        print_color "$CYAN" "  ℹ️  保留: 🚀-Phase4执行计划.md 和 🚀-立即开始-3分钟上手.md"
        
        # 2.3 清理冗余 README
        echo ""
        print_color "$YELLOW" "3. 清理冗余 README..."
        
        readme_files=(
            "README-项目现状-2025-10-25.md"
            "README-PROJECT-COMPLETE.md"
            "README-WORKSPACE-READY.md"
            "README-重构说明.md"
            "📖-README-项目导航.md"
        )
        
        for file in "${readme_files[@]}"; do
            remove_safe_item "$file" "README: $file"
        done
        print_color "$CYAN" "  ℹ️  保留: README.md 和 README_EN.md"
        
        echo ""
        print_color "$GREEN" "✅ 阶段2完成！"
    fi
    
    # 总结
    echo ""
    print_color "$CYAN" "================================================"
    print_color "$CYAN" "  清理完成总结"
    print_color "$CYAN" "================================================"
    echo ""
    
    if [ "$DRY_RUN" = false ]; then
        print_color "$GREEN" "✅ 清理已完成！"
        echo ""
        print_color "$YELLOW" "建议的后续步骤："
        echo "  1. 检查 git status"
        echo "  2. 验证项目功能正常"
        echo "  3. 更新文档链接"
        echo "  4. 提交更改"
        echo ""
        print_color "$CYAN" "运行以下命令查看更改："
        echo "  git status"
        echo "  git diff --stat"
    else
        print_color "$GREEN" "✅ 试运行完成！"
        echo ""
        print_color "$YELLOW" "要实际执行清理，请运行："
        echo "  ./scripts/cleanup_redundancy.sh --stage1"
        echo "  或"
        echo "  ./scripts/cleanup_redundancy.sh --all"
    fi
    
    echo ""
}

# 运行主函数
main


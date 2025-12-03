#!/bin/bash

# 归档旧的报告文档
# 日期: 2025-12-03

set -e

ARCHIVE_DIR="archive/docs-reports-2025-12"
DOCS_DIR="docs"

echo "📦 开始归档旧报告文档..."

# 创建归档目录
mkdir -p "$ARCHIVE_DIR"

# 归档技术栈实施相关报告
echo "📄 归档技术栈实施报告..."
mv "$DOCS_DIR"/00-技术栈实施*.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到技术栈实施报告"

# 归档项目评价相关报告
echo "📄 归档项目评价报告..."
mv "$DOCS_DIR"/00-项目评价*.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到项目评价报告"

# 归档改进任务相关报告
echo "📄 归档改进任务报告..."
mv "$DOCS_DIR"/00-改进任务*.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到改进任务报告"

# 归档文件归档相关报告
echo "📄 归档文件归档报告..."
mv "$DOCS_DIR"/00-文件归档*.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到文件归档报告"

# 归档项目状态相关报告
echo "📄 归档项目状态报告..."
mv "$DOCS_DIR"/00-项目状态*.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到项目状态报告"
mv "$DOCS_DIR"/00-项目完整*.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到项目完整报告"
mv "$DOCS_DIR"/00-项目最终*.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到项目最终报告"
mv "$DOCS_DIR"/00-项目结构*.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到项目结构报告"
mv "$DOCS_DIR"/00-项目重新*.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到项目重新报告"

# 归档各种完成报告
echo "📄 归档完成报告..."
mv "$DOCS_DIR"/*COMPLETE*.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到COMPLETE报告"
mv "$DOCS_DIR"/*FINAL*.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到FINAL报告"
mv "$DOCS_DIR"/*ULTIMATE*.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到ULTIMATE报告"

# 归档文档补充相关报告（今天刚生成的）
echo "📄 归档文档补充报告..."
mv "$DOCS_DIR"/00-文档完善*.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到文档完善报告"
mv "$DOCS_DIR"/00-链接修复*.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到链接修复报告"
mv "$DOCS_DIR"/00-工作完成*.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到工作完成报告"
mv "$DOCS_DIR"/00-工作汇总*.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到工作汇总报告"
mv "$DOCS_DIR"/00-最终完整*.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到最终完整报告"
mv "$DOCS_DIR"/00-下一阶段*.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到下一阶段报告"
mv "$DOCS_DIR"/fundamentals/00-完成声明*.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到完成声明"

# 归档其他总结性报告
echo "📄 归档其他总结报告..."
mv "$DOCS_DIR"/completion-summary.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到completion-summary"
mv "$DOCS_DIR"/features-summary.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到features-summary"
mv "$DOCS_DIR"/final-implementation-summary.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到final-implementation-summary"
mv "$DOCS_DIR"/implementation-status.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到implementation-status"
mv "$DOCS_DIR"/system-monitoring-implementation.md "$ARCHIVE_DIR/" 2>/dev/null || echo "  ℹ️  没有找到system-monitoring-implementation"

# 创建归档说明
cat > "$ARCHIVE_DIR/README.md" << 'EOF'
# 归档的报告文档

**归档日期**: 2025-12-03  
**原因**: 重复的进度报告和完成总结，不再需要

## 归档内容

本目录包含了项目历史上生成的各类报告文档：

- 技术栈实施报告
- 项目评价报告  
- 改进任务报告
- 文件归档报告
- 各种完成/总结报告
- 文档补充工作报告

## 保留的核心文档

以下文档仍保留在 `docs/` 目录：

- `README.md` - 文档总入口
- `architecture/` - 架构设计文档
- `00-项目改进计划总览.md` - 当前改进计划
- `IMPROVEMENT-TASK-BOARD.md` - 任务看板
- `00-架构代码检查与改进计划-2025-12-03.md` - 最新检查计划

## 查看历史

如需查看这些历史报告，请查看本目录中的文件。
EOF

# 统计归档文件数量
ARCHIVED_COUNT=$(find "$ARCHIVE_DIR" -type f -name "*.md" ! -name "README.md" | wc -l)

echo ""
echo "✅ 归档完成！"
echo "📊 归档了 $ARCHIVED_COUNT 个文档"
echo "📁 归档位置: $ARCHIVE_DIR"
echo ""
echo "保留的核心文档："
echo "  - docs/README.md"
echo "  - docs/architecture/"
echo "  - docs/00-项目改进计划总览.md"
echo "  - docs/IMPROVEMENT-TASK-BOARD.md"
echo "  - docs/00-架构代码检查与改进计划-2025-12-03.md"
echo ""


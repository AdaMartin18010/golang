#!/bin/bash
# 验证项目结构是否符合重组规范

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
GRAY='\033[0;90m'
NC='\033[0m' # No Color

# 统计变量
ERROR_COUNT=0
WARNING_COUNT=0
PASS_COUNT=0

echo -e "${CYAN}🔍 开始验证项目结构...${NC}"
echo ""

# 测试函数
test_rule() {
    local name="$1"
    local test_cmd="$2"
    local error_msg="$3"
    local pass_msg="$4"
    
    echo -n "➤ $name... "
    
    if eval "$test_cmd"; then
        echo -e "${GREEN}✅${NC}"
        if [ -n "$pass_msg" ]; then
            echo -e "${GRAY}  └─ $pass_msg${NC}"
        fi
        ((PASS_COUNT++))
    else
        echo -e "${RED}❌${NC}"
        echo -e "${YELLOW}  └─ $error_msg${NC}"
        ((ERROR_COUNT++))
    fi
}

echo "============================================================"
echo -e "${BLUE}📋 规则1: 文档代码分离${NC}"
echo "============================================================"

# 检查 docs/ 目录是否有代码文件
test_rule "docs/ 目录无 .go 文件" \
    "[ \$(find docs -name '*.go' 2>/dev/null | wc -l) -eq 0 ]" \
    "发现 .go 文件，应该移至 examples/" \
    "docs/ 目录纯文档 ✓"

test_rule "docs/ 目录无 go.mod 文件" \
    "[ \$(find docs -name 'go.mod' 2>/dev/null | wc -l) -eq 0 ]" \
    "发现 go.mod 文件" \
    "无 go.mod 文件 ✓"

test_rule "docs/ 目录无可执行文件" \
    "[ \$(find docs -type f -executable 2>/dev/null | wc -l) -eq 0 ]" \
    "发现可执行文件" \
    "无可执行文件 ✓"

echo ""
echo "============================================================"
echo -e "${BLUE}📋 规则2: 根目录清洁${NC}"
echo "============================================================"

test_rule "根目录无 Phase 报告" \
    "[ \$(ls Phase-*.md 2>/dev/null | wc -l) -eq 0 ]" \
    "发现 Phase 报告，应移至 reports/phase-reports/" \
    "无 Phase 报告 ✓"

test_rule "根目录文档数量合理" \
    "count=\$(ls *.md 2>/dev/null | wc -l); [ \$count -ge 8 ] && [ \$count -le 20 ]" \
    "根目录 .md 文件数量不合理" \
    "文件数量合理 ✓"

echo ""
echo "============================================================"
echo -e "${BLUE}📋 规则3: 目录职责${NC}"
echo "============================================================"

test_rule "存在 docs/ 目录" "[ -d docs ]" "缺少 docs/ 目录" "存在 ✓"
test_rule "存在 examples/ 目录" "[ -d examples ]" "缺少 examples/ 目录" "存在 ✓"
test_rule "存在 reports/ 目录" "[ -d reports ]" "缺少 reports/ 目录" "存在 ✓"
test_rule "存在 archive/ 目录" "[ -d archive ]" "缺少 archive/ 目录" "存在 ✓"
test_rule "存在 scripts/ 目录" "[ -d scripts ]" "缺少 scripts/ 目录" "存在 ✓"

echo ""
echo "============================================================"
echo -e "${BLUE}📋 规则4: 关键文件${NC}"
echo "============================================================"

key_files=("README.md" "RESTRUCTURE.md" "MIGRATION_GUIDE.md" "CONTRIBUTING.md" "FAQ.md" "LICENSE")

for file in "${key_files[@]}"; do
    test_rule "存在 $file" "[ -f $file ]" "缺少 $file" "存在 ✓"
done

echo ""
echo "============================================================"
echo -e "${BLUE}📋 规则5: examples/ 结构${NC}"
echo "============================================================"

example_dirs=("advanced" "concurrency" "go125" "modern-features" "testing-framework")

for dir in "${example_dirs[@]}"; do
    test_rule "存在 examples/$dir/" "[ -d examples/$dir ]" "缺少 examples/$dir/ 目录" "存在 ✓"
done

echo ""
echo "============================================================"
echo -e "${BLUE}📋 规则6: 代码质量${NC}"
echo "============================================================"

test_rule "examples/ 中代码可编译" \
    "cd examples && go build ./... 2>/dev/null" \
    "代码编译失败" \
    "编译通过 ✓"

echo ""
echo "============================================================"
echo -e "${MAGENTA}📊 验证结果统计${NC}"
echo "============================================================"

TOTAL=$((PASS_COUNT + ERROR_COUNT + WARNING_COUNT))

echo ""
echo -e "${GREEN}通过: $PASS_COUNT / $TOTAL${NC}"
echo -e "${RED}失败: $ERROR_COUNT / $TOTAL${NC}"
echo -e "${YELLOW}警告: $WARNING_COUNT / $TOTAL${NC}"
echo ""

if [ $ERROR_COUNT -eq 0 ]; then
    echo -e "${GREEN}✅ 项目结构验证通过！${NC}"
    echo -e "${GRAY}项目结构符合重组规范。${NC}"
    exit 0
else
    echo -e "${RED}❌ 项目结构验证失败！${NC}"
    echo -e "${YELLOW}发现 $ERROR_COUNT 个问题需要修复。${NC}"
    echo ""
    echo -e "${GRAY}请参考 RESTRUCTURE.md 了解项目结构规范。${NC}"
    exit 1
fi


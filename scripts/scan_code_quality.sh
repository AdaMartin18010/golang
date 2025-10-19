#!/bin/bash
# Code Quality and Runability Scanner
# 代码质量和可运行性扫描器

set -e

echo "==================================="
echo "Go Project Quality Scanner"
echo "==================================="
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 统计变量
total_go_files=0
total_test_files=0
compilable_modules=0
non_compilable_modules=0
total_modules=0

# 创建报告文件
REPORT_FILE="code_quality_report_$(date +%Y%m%d_%H%M%S).md"

echo "# Code Quality Scan Report" > "$REPORT_FILE"
echo "" >> "$REPORT_FILE"
echo "**Date**: $(date)" >> "$REPORT_FILE"
echo "**Go Version**: $(go version)" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"

# 1. 统计Go文件
echo -e "${YELLOW}[1/6] Counting Go files...${NC}"
total_go_files=$(find . -name "*.go" -not -path "*/vendor/*" -not -path "*/.git/*" | wc -l)
total_test_files=$(find . -name "*_test.go" -not -path "*/vendor/*" | wc -l)
total_code_files=$((total_go_files - total_test_files))

echo "  - Total Go files: $total_go_files"
echo "  - Code files: $total_code_files"
echo "  - Test files: $total_test_files"

# 2. 检查模块编译
echo -e "\n${YELLOW}[2/6] Checking module compilation...${NC}"
echo "" >> "$REPORT_FILE"
echo "## Compilation Status" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"

failed_modules=()

while IFS= read -r mod_file; do
    mod_dir=$(dirname "$mod_file")
    total_modules=$((total_modules + 1))
    
    echo -n "  Checking $mod_dir... "
    
    if (cd "$mod_dir" && go build ./... 2>&1 > /dev/null); then
        echo -e "${GREEN}✓${NC}"
        compilable_modules=$((compilable_modules + 1))
    else
        echo -e "${RED}✗${NC}"
        non_compilable_modules=$((non_compilable_modules + 1))
        failed_modules+=("$mod_dir")
        echo "- ❌ $mod_dir" >> "$REPORT_FILE"
    fi
done < <(find . -name "go.mod" -not -path "*/vendor/*")

if [ ${#failed_modules[@]} -eq 0 ]; then
    echo "" >> "$REPORT_FILE"
    echo "✅ **All modules compile successfully!**" >> "$REPORT_FILE"
fi

# 3. 运行测试
echo -e "\n${YELLOW}[3/6] Running tests...${NC}"
echo "" >> "$REPORT_FILE"
echo "## Test Results" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"

if go test ./... -cover -coverprofile=coverage.out 2>&1 | tee test_output.txt; then
    echo -e "${GREEN}Tests passed!${NC}"
    
    # 计算覆盖率
    if [ -f coverage.out ]; then
        coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
        echo "- **Coverage**: $coverage" >> "$REPORT_FILE"
        echo "  Coverage: $coverage"
    fi
else
    echo -e "${RED}Some tests failed!${NC}"
    echo "- ⚠️ **Some tests failed - see test_output.txt for details**" >> "$REPORT_FILE"
fi

# 4. 代码质量检查
echo -e "\n${YELLOW}[4/6] Running code quality checks...${NC}"
echo "" >> "$REPORT_FILE"
echo "## Code Quality" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"

# go vet
echo -n "  Running go vet... "
if go vet ./... 2>&1 > vet_output.txt; then
    echo -e "${GREEN}✓${NC}"
    echo "- ✅ **go vet**: Passed" >> "$REPORT_FILE"
else
    echo -e "${RED}✗${NC}"
    echo "- ⚠️ **go vet**: Issues found (see vet_output.txt)" >> "$REPORT_FILE"
fi

# gofmt
echo -n "  Checking formatting... "
unformatted=$(gofmt -l -s . 2>&1 | grep -v vendor | grep '.go$' || true)
if [ -z "$unformatted" ]; then
    echo -e "${GREEN}✓${NC}"
    echo "- ✅ **gofmt**: All files properly formatted" >> "$REPORT_FILE"
else
    echo -e "${RED}✗${NC}"
    echo "- ⚠️ **gofmt**: Some files need formatting" >> "$REPORT_FILE"
    echo "$unformatted" | while read -r file; do
        echo "  - $file" >> "$REPORT_FILE"
    done
fi

# 5. 依赖检查
echo -e "\n${YELLOW}[5/6] Checking dependencies...${NC}"
echo "" >> "$REPORT_FILE"
echo "## Dependencies" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"

module_count=0
while IFS= read -r mod_file; do
    module_count=$((module_count + 1))
done < <(find . -name "go.mod" -not -path "*/vendor/*")

echo "- **Total modules**: $module_count" >> "$REPORT_FILE"

# 6. 生成总结
echo -e "\n${YELLOW}[6/6] Generating summary...${NC}"
echo "" >> "$REPORT_FILE"
echo "## Summary" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"
echo "| Metric | Value |" >> "$REPORT_FILE"
echo "|--------|-------|" >> "$REPORT_FILE"
echo "| Total Go files | $total_go_files |" >> "$REPORT_FILE"
echo "| Code files | $total_code_files |" >> "$REPORT_FILE"
echo "| Test files | $total_test_files |" >> "$REPORT_FILE"
echo "| Test coverage | ${coverage:-N/A} |" >> "$REPORT_FILE"
echo "| Total modules | $total_modules |" >> "$REPORT_FILE"
echo "| Compilable modules | $compilable_modules |" >> "$REPORT_FILE"
echo "| Non-compilable modules | $non_compilable_modules |" >> "$REPORT_FILE"

test_coverage_percent=0
if [ "$total_test_files" -gt 0 ] && [ "$total_code_files" -gt 0 ]; then
    test_coverage_percent=$((total_test_files * 100 / total_code_files))
fi

echo "| Test file coverage | ${test_coverage_percent}% |" >> "$REPORT_FILE"

# 计算综合评分
compilation_score=$((compilable_modules * 100 / total_modules))
echo "" >> "$REPORT_FILE"
echo "## Quality Score" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"
echo "- **Compilation Success Rate**: ${compilation_score}%" >> "$REPORT_FILE"
echo "- **Test File Ratio**: ${test_coverage_percent}%" >> "$REPORT_FILE"

# 显示结果
echo ""
echo "==================================="
echo -e "${GREEN}Scan Complete!${NC}"
echo "==================================="
echo ""
echo "Summary:"
echo "  - Total Go files: $total_go_files"
echo "  - Test files: $total_test_files (${test_coverage_percent}%)"
echo "  - Compilable modules: $compilable_modules/$total_modules (${compilation_score}%)"
echo ""
echo "Report saved to: $REPORT_FILE"
echo ""

# 如果有失败的模块，列出来
if [ ${#failed_modules[@]} -gt 0 ]; then
    echo -e "${RED}Failed modules:${NC}"
    for mod in "${failed_modules[@]}"; do
        echo "  - $mod"
    done
    echo ""
    exit 1
fi

exit 0


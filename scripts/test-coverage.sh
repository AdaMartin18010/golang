#!/bin/bash
# 测试覆盖率脚本
# 功能：运行测试并生成覆盖率报告

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}开始运行测试覆盖率检查...${NC}"

# 覆盖率阈值
COVERAGE_THRESHOLD=70
PACKAGE_COVERAGE_THRESHOLD=60

# 覆盖率文件
COVERAGE_FILE="coverage.out"
COVERAGE_HTML="coverage.html"

# 清理旧的覆盖率文件
rm -f ${COVERAGE_FILE} ${COVERAGE_HTML}

# 运行测试并生成覆盖率
echo -e "${YELLOW}运行测试...${NC}"
go test -v -coverprofile=${COVERAGE_FILE} -covermode=atomic ./...

# 检查覆盖率文件是否存在
if [ ! -f ${COVERAGE_FILE} ]; then
    echo -e "${RED}错误: 覆盖率文件未生成${NC}"
    exit 1
fi

# 生成 HTML 报告
echo -e "${YELLOW}生成 HTML 覆盖率报告...${NC}"
go tool cover -html=${COVERAGE_FILE} -o ${COVERAGE_HTML}

# 获取总体覆盖率
TOTAL_COVERAGE=$(go tool cover -func=${COVERAGE_FILE} | grep total | awk '{print $3}' | sed 's/%//')

echo -e "${GREEN}总体覆盖率: ${TOTAL_COVERAGE}%${NC}"

# 检查覆盖率阈值
if (( $(echo "$TOTAL_COVERAGE < $COVERAGE_THRESHOLD" | bc -l) )); then
    echo -e "${RED}警告: 总体覆盖率 ${TOTAL_COVERAGE}% 低于阈值 ${COVERAGE_THRESHOLD}%${NC}"
    exit 1
fi

# 显示各包的覆盖率
echo -e "${YELLOW}各包覆盖率:${NC}"
go tool cover -func=${COVERAGE_FILE} | grep -E "github.com/yourusername/golang" | while read line; do
    PACKAGE=$(echo $line | awk '{print $1}')
    COVERAGE=$(echo $line | awk '{print $3}' | sed 's/%//')

    if (( $(echo "$COVERAGE < $PACKAGE_COVERAGE_THRESHOLD" | bc -l) )); then
        echo -e "${RED}  ${PACKAGE}: ${COVERAGE}% (低于阈值)${NC}"
    else
        echo -e "${GREEN}  ${PACKAGE}: ${COVERAGE}%${NC}"
    fi
done

echo -e "${GREEN}覆盖率报告已生成: ${COVERAGE_HTML}${NC}"
echo -e "${GREEN}测试覆盖率检查完成!${NC}"

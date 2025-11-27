#!/bin/bash
# 开发环境检查脚本
# 功能：检查开发环境是否配置正确

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  开发环境检查${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

ERRORS=0
WARNINGS=0

# 检查函数
check_command() {
    local cmd=$1
    local name=$2
    local required=$3

    if command -v $cmd &> /dev/null; then
        version=$($cmd --version 2>/dev/null || $cmd version 2>/dev/null || echo "已安装")
        echo -e "${GREEN}✅ $name: $version${NC}"
        return 0
    else
        if [ "$required" = "required" ]; then
            echo -e "${RED}❌ $name: 未安装（必需）${NC}"
            ERRORS=$((ERRORS + 1))
            return 1
        else
            echo -e "${YELLOW}⚠️  $name: 未安装（可选）${NC}"
            WARNINGS=$((WARNINGS + 1))
            return 0
        fi
    fi
}

# 必需工具
echo -e "${YELLOW}必需工具:${NC}"
check_command "go" "Go" "required"
check_command "git" "Git" "required"
echo ""

# 开发工具
echo -e "${YELLOW}开发工具:${NC}"
check_command "air" "Air (热重载)" "optional"
check_command "wire" "Wire (依赖注入)" "optional"
check_command "ent" "Ent (ORM)" "optional"
check_command "oapi-codegen" "oapi-codegen (OpenAPI)" "optional"
check_command "golangci-lint" "golangci-lint (代码检查)" "optional"
echo ""

# 可选工具
echo -e "${YELLOW}可选工具:${NC}"
check_command "docker" "Docker" "optional"
check_command "docker-compose" "Docker Compose" "optional"
echo ""

# 检查 Go 版本
echo -e "${YELLOW}Go 版本检查:${NC}"
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.25.3"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" = "$REQUIRED_VERSION" ]; then
    echo -e "${GREEN}✅ Go 版本符合要求: $GO_VERSION${NC}"
else
    echo -e "${RED}❌ Go 版本不符合要求: 需要 $REQUIRED_VERSION+，当前: $GO_VERSION${NC}"
    ERRORS=$((ERRORS + 1))
fi
echo ""

# 检查项目文件
echo -e "${YELLOW}项目文件检查:${NC}"
files=("go.mod" "go.sum" "Makefile" ".air.toml" "api/openapi/openapi.yaml" "api/asyncapi/asyncapi.yaml")
for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        echo -e "${GREEN}✅ $file${NC}"
    else
        echo -e "${YELLOW}⚠️  $file: 不存在${NC}"
        WARNINGS=$((WARNINGS + 1))
    fi
done
echo ""

# 检查依赖
echo -e "${YELLOW}依赖检查:${NC}"
if [ -f "go.mod" ]; then
    if go mod verify &> /dev/null; then
        echo -e "${GREEN}✅ Go 模块依赖验证通过${NC}"
    else
        echo -e "${RED}❌ Go 模块依赖验证失败${NC}"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo -e "${RED}❌ go.mod 文件不存在${NC}"
    ERRORS=$((ERRORS + 1))
fi
echo ""

# 总结
echo -e "${BLUE}========================================${NC}"
if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✅ 环境检查通过！${NC}"
    exit 0
elif [ $ERRORS -eq 0 ]; then
    echo -e "${YELLOW}⚠️  环境检查完成，有 $WARNINGS 个警告${NC}"
    exit 0
else
    echo -e "${RED}❌ 环境检查失败，有 $ERRORS 个错误，$WARNINGS 个警告${NC}"
    echo -e "${YELLOW}运行 'make install-tools' 或 'bash scripts/dev/setup.sh' 安装缺失的工具${NC}"
    exit 1
fi

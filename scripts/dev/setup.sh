#!/bin/bash
# 开发环境设置脚本
# 功能：一键设置开发环境，安装所有必需的工具和依赖

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  框架开发环境设置${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 检查 Go 版本
echo -e "${YELLOW}检查 Go 版本...${NC}"
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.25.3"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo -e "${RED}错误: 需要 Go $REQUIRED_VERSION 或更高版本，当前版本: $GO_VERSION${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Go 版本检查通过: $GO_VERSION${NC}"
echo ""

# 安装依赖
echo -e "${YELLOW}安装 Go 依赖...${NC}"
go mod download
go mod tidy
echo -e "${GREEN}✅ 依赖安装完成${NC}"
echo ""

# 安装开发工具
echo -e "${YELLOW}安装开发工具...${NC}"

install_tool() {
    local tool=$1
    local package=$2
    local name=$3

    if command -v $tool &> /dev/null; then
        echo -e "${GREEN}  ✅ $name 已安装${NC}"
    else
        echo -e "${YELLOW}  安装 $name...${NC}"
        go install $package@latest || echo -e "${RED}  ❌ $name 安装失败${NC}"
    fi
}

install_tool "air" "github.com/cosmtrek/air" "Air (热重载)"
install_tool "wire" "github.com/google/wire/cmd/wire" "Wire (依赖注入)"
install_tool "ent" "entgo.io/ent/cmd/ent" "Ent (ORM)"
install_tool "oapi-codegen" "github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen" "oapi-codegen (OpenAPI)"
install_tool "golangci-lint" "github.com/golangci/golangci-lint/cmd/golangci-lint" "golangci-lint (代码检查)"

echo ""

# 检查 Docker（可选）
echo -e "${YELLOW}检查 Docker...${NC}"
if command -v docker &> /dev/null; then
    echo -e "${GREEN}✅ Docker 已安装${NC}"
    if docker ps &> /dev/null; then
        echo -e "${GREEN}✅ Docker 服务运行中${NC}"
    else
        echo -e "${YELLOW}⚠️  Docker 服务未运行（可选）${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  Docker 未安装（可选，用于 API 规范验证和文档生成）${NC}"
fi
echo ""

# 生成代码
echo -e "${YELLOW}生成代码...${NC}"
make generate || echo -e "${YELLOW}⚠️  代码生成部分失败（可能需要先配置）${NC}"
echo ""

# 运行测试
echo -e "${YELLOW}运行测试...${NC}"
go test ./... || echo -e "${YELLOW}⚠️  部分测试失败${NC}"
echo ""

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  开发环境设置完成！${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${BLUE}可用命令:${NC}"
echo -e "  ${YELLOW}make run-dev${NC}      - 开发模式运行（热重载）"
echo -e "  ${YELLOW}make test${NC}         - 运行测试"
echo -e "  ${YELLOW}make lint${NC}         - 代码检查"
echo -e "  ${YELLOW}make generate${NC}     - 生成代码"
echo -e "  ${YELLOW}make help${NC}         - 查看所有命令"
echo ""

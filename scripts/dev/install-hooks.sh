#!/bin/bash
# 安装 Git hooks 脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  安装 Git Hooks${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 检查是否在 Git 仓库中
if [ ! -d ".git" ]; then
    echo -e "${RED}❌ 错误: 当前目录不是 Git 仓库${NC}"
    exit 1
fi

# 创建 .git/hooks 目录（如果不存在）
mkdir -p .git/hooks

# 复制 pre-commit hook
if [ -f ".githooks/pre-commit" ]; then
    cp .githooks/pre-commit .git/hooks/pre-commit
    chmod +x .git/hooks/pre-commit
    echo -e "${GREEN}✅ Pre-commit hook 已安装${NC}"
else
    echo -e "${YELLOW}⚠️  .githooks/pre-commit 不存在${NC}"
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Git Hooks 安装完成！${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${BLUE}说明:${NC}"
echo -e "  Pre-commit hook 会在每次提交前自动运行："
echo -e "  - 代码格式检查 (gofmt)"
echo -e "  - go vet 检查"
echo -e "  - golangci-lint 检查（如果已安装）"
echo -e "  - 相关测试"
echo ""

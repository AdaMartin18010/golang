#!/bin/bash
# 安全扫描脚本
# 功能：运行多种安全扫描工具

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}开始安全扫描...${NC}"

# 检查工具是否安装
check_tool() {
    if ! command -v $1 &> /dev/null; then
        echo -e "${YELLOW}警告: $1 未安装，跳过${NC}"
        return 1
    fi
    return 0
}

# 1. Gosec 安全扫描
if check_tool gosec; then
    echo -e "${YELLOW}运行 Gosec 安全扫描...${NC}"
    gosec -fmt json -out gosec-report.json ./... || true
    gosec ./...
fi

# 2. Trivy 漏洞扫描
if check_tool trivy; then
    echo -e "${YELLOW}运行 Trivy 漏洞扫描...${NC}"
    trivy fs --format json --output trivy-report.json . || true
    trivy fs .
fi

# 3. 依赖漏洞扫描
if check_tool govulncheck; then
    echo -e "${YELLOW}运行 Go 漏洞检查...${NC}"
    govulncheck ./... || true
fi

# 4. 检查硬编码密钥
echo -e "${YELLOW}检查硬编码密钥...${NC}"
if grep -r "password.*=" --include="*.go" --exclude-dir=vendor . | grep -v "test" | grep -v "example"; then
    echo -e "${RED}警告: 发现可能的硬编码密码${NC}"
fi

if grep -r "secret.*=" --include="*.go" --exclude-dir=vendor . | grep -v "test" | grep -v "example"; then
    echo -e "${RED}警告: 发现可能的硬编码密钥${NC}"
fi

# 5. 检查敏感信息
echo -e "${YELLOW}检查敏感信息...${NC}"
if grep -r "api[_-]key" --include="*.go" --exclude-dir=vendor . -i; then
    echo -e "${RED}警告: 发现可能的 API 密钥${NC}"
fi

if grep -r "private[_-]key" --include="*.go" --exclude-dir=vendor . -i; then
    echo -e "${RED}警告: 发现可能的私钥${NC}"
fi

echo -e "${GREEN}安全扫描完成!${NC}"

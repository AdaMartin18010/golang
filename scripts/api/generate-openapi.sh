#!/bin/bash
# OpenAPI 代码生成脚本
# 功能：从 OpenAPI 规范生成 Go 代码（类型、服务器、客户端）

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

OPENAPI_SPEC="api/openapi/openapi.yaml"
OUTPUT_DIR="internal/interfaces/http/openapi"
CLIENT_DIR="pkg/api/client"

echo -e "${GREEN}开始生成 OpenAPI 代码...${NC}"

# 检查 oapi-codegen 是否安装
if ! command -v oapi-codegen &> /dev/null; then
    echo -e "${YELLOW}oapi-codegen 未安装，正在安装...${NC}"
    go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
fi

# 检查 OpenAPI 规范文件是否存在
if [ ! -f "$OPENAPI_SPEC" ]; then
    echo -e "${RED}错误: OpenAPI 规范文件不存在: $OPENAPI_SPEC${NC}"
    exit 1
fi

# 创建输出目录
mkdir -p "$OUTPUT_DIR"
mkdir -p "$CLIENT_DIR"

# 生成服务器代码（类型、服务器接口、Chi 集成）
echo -e "${YELLOW}生成服务器代码...${NC}"
oapi-codegen \
  -generate types,server,chi-server,spec \
  -package openapi \
  -o "$OUTPUT_DIR/server.gen.go" \
  "$OPENAPI_SPEC"

# 生成客户端代码
echo -e "${YELLOW}生成客户端代码...${NC}"
oapi-codegen \
  -generate types,client \
  -package client \
  -o "$CLIENT_DIR/client.gen.go" \
  "$OPENAPI_SPEC"

echo -e "${GREEN}OpenAPI 代码生成完成!${NC}"
echo -e "${GREEN}服务器代码: $OUTPUT_DIR/server.gen.go${NC}"
echo -e "${GREEN}客户端代码: $CLIENT_DIR/client.gen.go${NC}"

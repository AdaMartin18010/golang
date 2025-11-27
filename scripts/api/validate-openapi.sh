#!/bin/bash
# OpenAPI 规范验证脚本
# 功能：验证 OpenAPI 规范文件的正确性

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

OPENAPI_SPEC="api/openapi/openapi.yaml"

echo -e "${GREEN}开始验证 OpenAPI 规范...${NC}"

# 检查规范文件是否存在
if [ ! -f "$OPENAPI_SPEC" ]; then
    echo -e "${RED}错误: OpenAPI 规范文件不存在: $OPENAPI_SPEC${NC}"
    exit 1
fi

# 使用 Docker 运行 OpenAPI Generator 验证
echo -e "${YELLOW}验证 OpenAPI 规范...${NC}"
docker run --rm \
  -v "${PWD}:/local" \
  openapitools/openapi-generator-cli:latest \
  validate -i "/local/$OPENAPI_SPEC"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}OpenAPI 规范验证通过!${NC}"
else
    echo -e "${RED}OpenAPI 规范验证失败!${NC}"
    exit 1
fi

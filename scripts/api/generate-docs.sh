#!/bin/bash
# API 文档生成脚本
# 功能：生成 OpenAPI 和 AsyncAPI 的 HTML 文档

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

OPENAPI_SPEC="api/openapi/openapi.yaml"
ASYNCAPI_SPEC="api/asyncapi/asyncapi.yaml"
DOCS_DIR="docs/api"

echo -e "${GREEN}开始生成 API 文档...${NC}"

# 创建文档目录
mkdir -p "$DOCS_DIR/openapi"
mkdir -p "$DOCS_DIR/asyncapi"

# 生成 OpenAPI 文档
if [ -f "$OPENAPI_SPEC" ]; then
    echo -e "${YELLOW}生成 OpenAPI 文档...${NC}"
    docker run --rm \
      -v "${PWD}:/local" \
      openapitools/openapi-generator-cli:latest \
      generate -i "/local/$OPENAPI_SPEC" \
      -g html \
      -o "/local/$DOCS_DIR/openapi"

    if [ $? -eq 0 ]; then
        echo -e "${GREEN}OpenAPI 文档已生成: $DOCS_DIR/openapi/index.html${NC}"
    fi
else
    echo -e "${YELLOW}警告: OpenAPI 规范文件不存在: $OPENAPI_SPEC${NC}"
fi

# 生成 AsyncAPI 文档
if [ -f "$ASYNCAPI_SPEC" ]; then
    echo -e "${YELLOW}生成 AsyncAPI 文档...${NC}"
    docker run --rm \
      -v "${PWD}:/local" \
      asyncapi/generator-cli:latest \
      generate -g html \
      -i "/local/$ASYNCAPI_SPEC" \
      -o "/local/$DOCS_DIR/asyncapi"

    if [ $? -eq 0 ]; then
        echo -e "${GREEN}AsyncAPI 文档已生成: $DOCS_DIR/asyncapi/index.html${NC}"
    fi
else
    echo -e "${YELLOW}警告: AsyncAPI 规范文件不存在: $ASYNCAPI_SPEC${NC}"
fi

echo -e "${GREEN}API 文档生成完成!${NC}"

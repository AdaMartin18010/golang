#!/bin/bash
# AsyncAPI 代码生成脚本
# 功能：从 AsyncAPI 规范生成 Go 代码和文档

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

ASYNCAPI_SPEC="api/asyncapi/asyncapi.yaml"
OUTPUT_DIR="pkg/api/async"

echo -e "${GREEN}开始生成 AsyncAPI 代码...${NC}"

# 检查 AsyncAPI 规范文件是否存在
if [ ! -f "$ASYNCAPI_SPEC" ]; then
    echo -e "${RED}错误: AsyncAPI 规范文件不存在: $ASYNCAPI_SPEC${NC}"
    exit 1
fi

# 创建输出目录
mkdir -p "$OUTPUT_DIR"

# 使用 Docker 运行 AsyncAPI Generator
echo -e "${YELLOW}使用 AsyncAPI Generator 生成代码...${NC}"
docker run --rm \
  -v "${PWD}:/local" \
  asyncapi/generator-cli:latest \
  generate -g go \
  -i "/local/$ASYNCAPI_SPEC" \
  -o "/local/$OUTPUT_DIR" \
  -p packageName=async

if [ $? -eq 0 ]; then
    echo -e "${GREEN}AsyncAPI 代码生成完成!${NC}"
    echo -e "${GREEN}输出目录: $OUTPUT_DIR${NC}"
else
    echo -e "${YELLOW}注意: AsyncAPI Generator 可能未安装或 Docker 未运行${NC}"
    echo -e "${YELLOW}可以手动运行: docker pull asyncapi/generator-cli${NC}"
fi

#!/bin/bash

# Docker 构建脚本
# 用于构建应用 Docker 镜像

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
IMAGE_NAME="${IMAGE_NAME:-app}"
IMAGE_TAG="${IMAGE_TAG:-latest}"
DOCKERFILE="${DOCKERFILE:-deployments/docker/Dockerfile}"
BUILD_CONTEXT="${BUILD_CONTEXT:-.}"

echo -e "${GREEN}开始构建 Docker 镜像...${NC}"
echo -e "镜像名称: ${YELLOW}${IMAGE_NAME}:${IMAGE_TAG}${NC}"
echo -e "Dockerfile: ${YELLOW}${DOCKERFILE}${NC}"
echo -e "构建上下文: ${YELLOW}${BUILD_CONTEXT}${NC}"
echo ""

# 构建镜像
docker build \
  -f "${DOCKERFILE}" \
  -t "${IMAGE_NAME}:${IMAGE_TAG}" \
  "${BUILD_CONTEXT}"

echo ""
echo -e "${GREEN}✅ Docker 镜像构建成功！${NC}"
echo -e "镜像: ${YELLOW}${IMAGE_NAME}:${IMAGE_TAG}${NC}"

# 显示镜像信息
echo ""
echo -e "${GREEN}镜像信息:${NC}"
docker images "${IMAGE_NAME}:${IMAGE_TAG}"

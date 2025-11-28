#!/bin/bash

# Docker 推送脚本
# 用于推送 Docker 镜像到镜像仓库

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
IMAGE_NAME="${IMAGE_NAME:-app}"
IMAGE_TAG="${IMAGE_TAG:-latest}"
REGISTRY="${REGISTRY:-docker.io}"
REGISTRY_USER="${REGISTRY_USER:-}"

if [ -z "${REGISTRY_USER}" ]; then
    echo -e "${RED}错误: 请设置 REGISTRY_USER 环境变量${NC}"
    exit 1
fi

FULL_IMAGE_NAME="${REGISTRY}/${REGISTRY_USER}/${IMAGE_NAME}:${IMAGE_TAG}"

echo -e "${GREEN}开始推送 Docker 镜像...${NC}"
echo -e "镜像: ${YELLOW}${FULL_IMAGE_NAME}${NC}"
echo ""

# 标记镜像
docker tag "${IMAGE_NAME}:${IMAGE_TAG}" "${FULL_IMAGE_NAME}"

# 推送镜像
docker push "${FULL_IMAGE_NAME}"

echo ""
echo -e "${GREEN}✅ Docker 镜像推送成功！${NC}"
echo -e "镜像: ${YELLOW}${FULL_IMAGE_NAME}${NC}"

#!/bin/bash

# Kubernetes 删除脚本
# 用于删除 Kubernetes 部署

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
NAMESPACE="${NAMESPACE:-default}"
K8S_DIR="${K8S_DIR:-deployments/kubernetes}"

echo -e "${YELLOW}警告: 即将删除 Kubernetes 部署${NC}"
echo -e "命名空间: ${YELLOW}${NAMESPACE}${NC}"
echo -e "配置目录: ${YELLOW}${K8S_DIR}${NC}"
echo ""
read -p "确认删除? (y/N): " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${RED}取消删除${NC}"
    exit 1
fi

# 删除资源
echo -e "${GREEN}删除 Kubernetes 资源...${NC}"

if [ -f "${K8S_DIR}/hpa.yaml" ]; then
    kubectl delete -f "${K8S_DIR}/hpa.yaml" -n "${NAMESPACE}" --ignore-not-found=true
fi

if [ -f "${K8S_DIR}/service.yaml" ]; then
    kubectl delete -f "${K8S_DIR}/service.yaml" -n "${NAMESPACE}" --ignore-not-found=true
fi

if [ -f "${K8S_DIR}/deployment.yaml" ]; then
    kubectl delete -f "${K8S_DIR}/deployment.yaml" -n "${NAMESPACE}" --ignore-not-found=true
fi

if [ -f "${K8S_DIR}/configmap.yaml" ]; then
    kubectl delete -f "${K8S_DIR}/configmap.yaml" -n "${NAMESPACE}" --ignore-not-found=true
fi

if [ -f "${K8S_DIR}/secret.yaml" ]; then
    kubectl delete -f "${K8S_DIR}/secret.yaml" -n "${NAMESPACE}" --ignore-not-found=true
fi

echo ""
echo -e "${GREEN}✅ 删除完成！${NC}"

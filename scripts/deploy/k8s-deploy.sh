#!/bin/bash

# Kubernetes 部署脚本
# 用于部署应用到 Kubernetes 集群

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
NAMESPACE="${NAMESPACE:-default}"
K8S_DIR="${K8S_DIR:-deployments/kubernetes}"

echo -e "${GREEN}开始部署到 Kubernetes...${NC}"
echo -e "命名空间: ${YELLOW}${NAMESPACE}${NC}"
echo -e "配置目录: ${YELLOW}${K8S_DIR}${NC}"
echo ""

# 检查 kubectl
if ! command -v kubectl &> /dev/null; then
    echo -e "${RED}错误: kubectl 未安装${NC}"
    exit 1
fi

# 检查集群连接
if ! kubectl cluster-info &> /dev/null; then
    echo -e "${RED}错误: 无法连接到 Kubernetes 集群${NC}"
    exit 1
fi

# 创建命名空间（如果不存在）
kubectl create namespace "${NAMESPACE}" --dry-run=client -o yaml | kubectl apply -f -

# 部署 ConfigMap
if [ -f "${K8S_DIR}/configmap.yaml" ]; then
    echo -e "${GREEN}部署 ConfigMap...${NC}"
    kubectl apply -f "${K8S_DIR}/configmap.yaml" -n "${NAMESPACE}"
fi

# 部署 Secret（需要手动创建）
if [ -f "${K8S_DIR}/secret.yaml" ]; then
    echo -e "${GREEN}部署 Secret...${NC}"
    kubectl apply -f "${K8S_DIR}/secret.yaml" -n "${NAMESPACE}"
else
    echo -e "${YELLOW}警告: secret.yaml 不存在，请手动创建 Secret${NC}"
fi

# 部署 Deployment
if [ -f "${K8S_DIR}/deployment.yaml" ]; then
    echo -e "${GREEN}部署 Deployment...${NC}"
    kubectl apply -f "${K8S_DIR}/deployment.yaml" -n "${NAMESPACE}"
fi

# 部署 Service
if [ -f "${K8S_DIR}/service.yaml" ]; then
    echo -e "${GREEN}部署 Service...${NC}"
    kubectl apply -f "${K8S_DIR}/service.yaml" -n "${NAMESPACE}"
fi

# 部署 HPA（可选）
if [ -f "${K8S_DIR}/hpa.yaml" ]; then
    echo -e "${GREEN}部署 HPA...${NC}"
    kubectl apply -f "${K8S_DIR}/hpa.yaml" -n "${NAMESPACE}"
fi

echo ""
echo -e "${GREEN}✅ 部署完成！${NC}"
echo ""
echo -e "${GREEN}检查部署状态:${NC}"
kubectl get pods -n "${NAMESPACE}" -l app=app
kubectl get svc -n "${NAMESPACE}" -l app=app

#!/bin/bash
# 检查所有服务的健康状态

set -e

echo "=== Docker Compose 服务状态检查 ==="
echo ""

# 检查 Docker Compose 服务
docker-compose ps

echo ""
echo "=== 服务健康检查 ==="
echo ""

# 检查应用服务
echo "1. 应用服务 (http://localhost:8080/health):"
curl -s http://localhost:8080/health || echo "应用服务不可用"

echo ""
echo "2. PostgreSQL 主节点:"
docker-compose exec -T db pg_isready -U user -d mydb || echo "主节点不可用"

echo ""
echo "3. PostgreSQL 备节点:"
docker-compose exec -T db-replica pg_isready -U user -d mydb || echo "备节点不可用"

echo ""
echo "4. Redis:"
docker-compose exec -T redis redis-cli -a redispassword ping || echo "Redis 不可用"

echo ""
echo "5. Temporal:"
curl -s http://localhost:8088 || echo "Temporal UI 不可用"

echo ""
echo "6. Prometheus:"
curl -s http://localhost:9090/-/healthy || echo "Prometheus 不可用"

echo ""
echo "7. Grafana:"
curl -s http://localhost:3000/api/health || echo "Grafana 不可用"

echo ""
echo "=== 检查完成 ==="

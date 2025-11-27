#!/bin/bash
# 备节点初始化脚本（在数据目录为空时执行）

set -e

# 这个脚本在 docker-entrypoint-initdb.d 中执行
# 但只在数据目录为空时执行
# 实际的复制设置由 docker-compose.yml 中的 command 处理

echo "Replica setup script loaded"

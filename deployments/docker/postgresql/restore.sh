#!/bin/bash
# PostgreSQL 数据库恢复脚本

set -e

PRIMARY_HOST="${DATABASE_HOST:-db}"
PRIMARY_PORT="${DATABASE_PORT:-5432}"
PRIMARY_USER="${DATABASE_USER:-user}"
PRIMARY_PASSWORD="${DATABASE_PASSWORD:-password}"
PRIMARY_DB="${DATABASE_DBNAME:-mydb}"

BACKUP_FILE="$1"

if [ -z "$BACKUP_FILE" ]; then
    echo "用法: $0 <备份文件路径>"
    echo "示例: $0 ./backups/mydb_20240101_120000.sql.gz"
    exit 1
fi

if [ ! -f "$BACKUP_FILE" ]; then
    echo "错误: 备份文件不存在: $BACKUP_FILE"
    exit 1
fi

export PGPASSWORD="$PRIMARY_PASSWORD"

echo "=== PostgreSQL 数据库恢复 ==="
echo "数据库: $PRIMARY_DB"
echo "主机: $PRIMARY_HOST:$PRIMARY_PORT"
echo "备份文件: $BACKUP_FILE"
echo ""
read -p "警告: 这将覆盖现有数据库。是否继续? (yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo "取消恢复"
    exit 0
fi

echo "开始恢复..."

# 判断文件类型并恢复
if [[ "$BACKUP_FILE" == *.gz ]]; then
    gunzip -c "$BACKUP_FILE" | psql -h "$PRIMARY_HOST" -p "$PRIMARY_PORT" -U "$PRIMARY_USER" -d "$PRIMARY_DB"
else
    psql -h "$PRIMARY_HOST" -p "$PRIMARY_PORT" -U "$PRIMARY_USER" -d "$PRIMARY_DB" < "$BACKUP_FILE"
fi

if [ $? -eq 0 ]; then
    echo "恢复成功"
else
    echo "恢复失败"
    exit 1
fi

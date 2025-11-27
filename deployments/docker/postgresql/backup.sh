#!/bin/bash
# PostgreSQL 数据库备份脚本

set -e

PRIMARY_HOST="${DATABASE_HOST:-db}"
PRIMARY_PORT="${DATABASE_PORT:-5432}"
PRIMARY_USER="${DATABASE_USER:-user}"
PRIMARY_PASSWORD="${DATABASE_PASSWORD:-password}"
PRIMARY_DB="${DATABASE_DBNAME:-mydb}"

BACKUP_DIR="${BACKUP_DIR:-./backups}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/${PRIMARY_DB}_${TIMESTAMP}.sql"
BACKUP_FILE_COMPRESSED="${BACKUP_FILE}.gz"

export PGPASSWORD="$PRIMARY_PASSWORD"

# 创建备份目录
mkdir -p "$BACKUP_DIR"

echo "=== PostgreSQL 数据库备份 ==="
echo "数据库: $PRIMARY_DB"
echo "主机: $PRIMARY_HOST:$PRIMARY_PORT"
echo "备份文件: $BACKUP_FILE_COMPRESSED"
echo ""

# 执行备份
echo "开始备份..."
pg_dump -h "$PRIMARY_HOST" -p "$PRIMARY_PORT" -U "$PRIMARY_USER" -d "$PRIMARY_DB" \
    --format=plain \
    --no-owner \
    --no-privileges \
    --verbose \
    | gzip > "$BACKUP_FILE_COMPRESSED"

if [ $? -eq 0 ]; then
    echo "备份成功: $BACKUP_FILE_COMPRESSED"
    echo "文件大小: $(du -h "$BACKUP_FILE_COMPRESSED" | cut -f1)"

    # 保留最近 7 天的备份
    echo "清理旧备份（保留最近 7 天）..."
    find "$BACKUP_DIR" -name "${PRIMARY_DB}_*.sql.gz" -mtime +7 -delete
    echo "备份完成"
else
    echo "备份失败"
    exit 1
fi

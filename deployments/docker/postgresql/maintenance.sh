#!/bin/bash
# PostgreSQL 数据库维护脚本

set -e

PRIMARY_HOST="${DATABASE_HOST:-db}"
PRIMARY_PORT="${DATABASE_PORT:-5432}"
PRIMARY_USER="${DATABASE_USER:-user}"
PRIMARY_PASSWORD="${DATABASE_PASSWORD:-password}"
PRIMARY_DB="${DATABASE_DBNAME:-mydb}"

export PGPASSWORD="$PRIMARY_PASSWORD"

echo "=== PostgreSQL 数据库维护 ==="
echo ""

# 1. 更新统计信息
echo "1. 更新表统计信息..."
psql -h "$PRIMARY_HOST" -p "$PRIMARY_PORT" -U "$PRIMARY_USER" -d "$PRIMARY_DB" -c "ANALYZE;"
echo "完成"

# 2. 清理过期数据
echo ""
echo "2. 清理过期数据..."
psql -h "$PRIMARY_HOST" -p "$PRIMARY_PORT" -U "$PRIMARY_USER" -d "$PRIMARY_DB" -c "VACUUM ANALYZE;"
echo "完成"

# 3. 检查数据库大小
echo ""
echo "3. 数据库大小:"
psql -h "$PRIMARY_HOST" -p "$PRIMARY_PORT" -U "$PRIMARY_USER" -d "$PRIMARY_DB" -c "
SELECT
    pg_database.datname,
    pg_size_pretty(pg_database_size(pg_database.datname)) AS size
FROM pg_database
WHERE datname = '$PRIMARY_DB';
"

# 4. 检查表大小
echo ""
echo "4. 表大小（前 10 个）:"
psql -h "$PRIMARY_HOST" -p "$PRIMARY_PORT" -U "$PRIMARY_USER" -d "$PRIMARY_DB" -c "
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size,
    pg_size_pretty(pg_relation_size(schemaname||'.'||tablename)) AS table_size,
    pg_size_pretty(pg_indexes_size(schemaname||'.'||tablename)) AS indexes_size
FROM pg_tables
WHERE schemaname NOT IN ('pg_catalog', 'information_schema')
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC
LIMIT 10;
"

# 5. 检查连接数
echo ""
echo "5. 当前连接数:"
psql -h "$PRIMARY_HOST" -p "$PRIMARY_PORT" -U "$PRIMARY_USER" -d "$PRIMARY_DB" -c "
SELECT
    count(*) AS total_connections,
    count(*) FILTER (WHERE state = 'active') AS active_connections,
    count(*) FILTER (WHERE state = 'idle') AS idle_connections,
    count(*) FILTER (WHERE state = 'idle in transaction') AS idle_in_transaction
FROM pg_stat_activity
WHERE datname = '$PRIMARY_DB';
"

# 6. 检查长查询
echo ""
echo "6. 运行时间超过 1 分钟的查询:"
psql -h "$PRIMARY_HOST" -p "$PRIMARY_PORT" -U "$PRIMARY_USER" -d "$PRIMARY_DB" -c "
SELECT
    pid,
    now() - pg_stat_activity.query_start AS duration,
    query,
    state
FROM pg_stat_activity
WHERE (now() - pg_stat_activity.query_start) > interval '1 minute'
    AND state != 'idle'
    AND datname = '$PRIMARY_DB'
ORDER BY duration DESC;
"

echo ""
echo "=== 维护完成 ==="

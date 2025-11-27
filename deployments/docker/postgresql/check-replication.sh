#!/bin/bash
# PostgreSQL 复制状态检查脚本

set -e

PRIMARY_HOST="${DATABASE_HOST:-db}"
PRIMARY_PORT="${DATABASE_PORT:-5432}"
PRIMARY_USER="${DATABASE_USER:-user}"
PRIMARY_PASSWORD="${DATABASE_PASSWORD:-password}"
PRIMARY_DB="${DATABASE_DBNAME:-mydb}"

REPLICA_HOST="${DATABASE_REPLICA_HOST:-db-replica}"
REPLICA_PORT="${DATABASE_REPLICA_PORT:-5432}"

export PGPASSWORD="$PRIMARY_PASSWORD"

echo "=== PostgreSQL 复制状态检查 ==="
echo ""

# 检查主节点状态
echo "1. 主节点状态 ($PRIMARY_HOST:$PRIMARY_PORT):"
psql -h "$PRIMARY_HOST" -p "$PRIMARY_PORT" -U "$PRIMARY_USER" -d "$PRIMARY_DB" -c "
SELECT
    application_name,
    client_addr,
    state,
    sync_state,
    sync_priority,
    pg_wal_lsn_diff(pg_current_wal_lsn(), sent_lsn) AS sent_lag_bytes,
    pg_wal_lsn_diff(sent_lsn, write_lsn) AS write_lag_bytes,
    pg_wal_lsn_diff(write_lsn, flush_lsn) AS flush_lag_bytes,
    pg_wal_lsn_diff(flush_lsn, replay_lsn) AS replay_lag_bytes,
    pg_wal_lsn_diff(pg_current_wal_lsn(), replay_lsn) AS total_lag_bytes
FROM pg_stat_replication;
" || echo "无法连接到主节点"

echo ""
echo "2. 主节点 WAL 信息:"
psql -h "$PRIMARY_HOST" -p "$PRIMARY_PORT" -U "$PRIMARY_USER" -d "$PRIMARY_DB" -c "
SELECT
    pg_current_wal_lsn() AS current_wal_lsn,
    pg_wal_lsn_diff(pg_current_wal_lsn(), '0/0') AS current_wal_bytes;
"

echo ""
export PGPASSWORD="$PRIMARY_PASSWORD"
echo "3. 备节点状态 ($REPLICA_HOST:$REPLICA_PORT):"
psql -h "$REPLICA_HOST" -p "$REPLICA_PORT" -U "$PRIMARY_USER" -d "$PRIMARY_DB" -c "
SELECT
    pg_is_in_recovery() AS is_replica,
    pg_last_wal_receive_lsn() AS received_lsn,
    pg_last_wal_replay_lsn() AS replayed_lsn,
    pg_wal_lsn_diff(pg_last_wal_receive_lsn(), pg_last_wal_replay_lsn()) AS replay_lag_bytes,
    EXTRACT(EPOCH FROM (now() - pg_last_xact_replay_timestamp())) AS replay_lag_seconds;
" || echo "无法连接到备节点"

echo ""
echo "=== 检查完成 ==="

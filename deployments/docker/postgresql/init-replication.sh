#!/bin/bash
set -e

# 创建复制用户
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER $POSTGRES_REPLICATION_USER REPLICATION LOGIN PASSWORD '$POSTGRES_REPLICATION_PASSWORD';
EOSQL

echo "Replication user created successfully"


#!/bin/bash
set -e

# 等待 PostgreSQL 启动
until pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB"; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 2
done

# 创建复制用户（如果不存在）
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    DO \$\$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = '$POSTGRES_REPLICATION_USER') THEN
            CREATE USER $POSTGRES_REPLICATION_USER REPLICATION LOGIN PASSWORD '$POSTGRES_REPLICATION_PASSWORD';
            RAISE NOTICE 'Replication user created';
        ELSE
            RAISE NOTICE 'Replication user already exists';
        END IF;
    END
    \$\$;
EOSQL

echo "Replication user setup completed"

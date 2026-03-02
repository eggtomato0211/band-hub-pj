#!/bin/bash
set -e

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-bandhub}"
DB_NAME="${DB_NAME:-bandhub_dev}"

echo "=== マイグレーション開始 ==="

for file in /migrations/*.sql; do
    echo "実行中: $(basename "$file")"
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$file"
done

echo "=== マイグレーション完了 ==="
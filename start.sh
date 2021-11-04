#!/bin/sh

#脚本将立即退出，if return value ！= 0
set -e

echo "run db migration"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start app"
exec "$@"
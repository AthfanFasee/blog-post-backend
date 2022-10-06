#!/bin/sh

set -e

echo "run db migration"
/app/migrate -path /app/migration -database "$DB_DSN" -verbose up

echo "start the app"
exec "$@"
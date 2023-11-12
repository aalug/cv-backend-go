#!/bin/sh

set -e

source /app/app.env

echo "/run db migrations"
/app/migrate -path /app/migrations -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"
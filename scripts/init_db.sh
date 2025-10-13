#!/bin/bash
set -e

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"

if [ -f "$ROOT_DIR/.env" ]; then
  set -a
  source "$ROOT_DIR/.env"
  set +a
  echo ".env variables loaded from $ROOT_DIR/.env"
else
  echo ".env file NOT found in $ROOT_DIR"
fi

DB_NAME=${DB_NAME:-social_db}
DB_USER=${DB_USER:-social_user}
DB_PASS=${DB_PASS:-social_p@ss}
PG_SUPERUSER=${PG_SUPERUSER:-postgres}

USER_EXISTS=$(psql -U "$PG_SUPERUSER" -tAc "SELECT 1 FROM pg_roles WHERE rolname='$DB_USER'")

if [ "$USER_EXISTS" != "1" ]; then
    echo "Creating new db user $DB_USER..."
    psql -U "$PG_SUPERUSER" -c "CREATE USER $DB_USER WITH PASSWORD '$DB_PASS';"
else
    echo "DB user already exists skipping..."
fi

DB_EXISTS=$(psql -U "$PG_SUPERUSER" -tAc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'")

if [ "$DB_EXISTS" != "1" ]; then
    echo "Creating DB $DB_NAME..."
    psql -U "$PG_SUPERUSER" -c "CREATE DATABASE $DB_NAME OWNER $DB_USER;"
else
    echo "DB ealready exists skipping..."
fi

echo "DB sсript finished successfully!"
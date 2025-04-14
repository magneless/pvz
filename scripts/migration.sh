#!/bin/bash

set -e

if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
else
    echo ".env not found"
    exit 1
fi

if ! command -v yq &> /dev/null; then
    echo "util yq not found"
    exit 1
fi

if ! command -v migrate &> /dev/null; then
    echo "util migrate not found"
    exit 1
fi

CONFIG_FILE="$CONFIG_PATH"
HOST=$(yq '.storage.host' "$CONFIG_FILE")
PORT=$(yq '.storage.port' "$CONFIG_FILE")
USER=$(yq '.storage.username' "$CONFIG_FILE")
DB=$(yq '.storage.dbname' "$CONFIG_FILE")
SSLMODE=$(yq '.storage.sslmode' "$CONFIG_FILE")

DB_URL="postgres://${USER}:${DB_PASSWORD}@${HOST}:${PORT}/${DB}?sslmode=${SSLMODE}"

ACTION=${1:-up}

if [[ "$ACTION" != "up" && "$ACTION" != "down" ]]; then
    echo "wrong argument"
    exit 1
fi

migrate -database "$DB_URL" -path ./migrations  "$ACTION"

echo "migration '$ACTION' done"

#!/usr/bin/env bash
set -euo pipefail

ENV_NAME="${1:-testing}"
ACTION="${2:-up}"
ENV_FILE=".env.${ENV_NAME}"

if [ "$ENV_NAME" = "production" ]; then
    ENV_FILE=".env"
fi

if [ ! -f "$ENV_FILE" ]; then
    echo "error: $ENV_FILE tidak ditemukan"
    exit 1
fi

URL=$(grep -E '^DATABASE_URL=' "$ENV_FILE" | cut -d'=' -f2- | tr -d '\r')

if [ -z "$URL" ]; then
    echo "error: DATABASE_URL tidak ditemukan di $ENV_FILE"
    exit 1
fi

echo "Menjalankan migrate $ACTION ke branch $ENV_NAME..."
migrate -path migrations -database "$URL" "$ACTION"

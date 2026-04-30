#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
SERVICE_DIR="$PROJECT_DIR/services/weclawbotnotify-api"

cd "$SERVICE_DIR"

echo "=== 编译 weclawbotnotify-api ==="
go build -o weclawbotnotify-api ./weclawbotnotify.go

echo "=== 启动服务 ==="
./weclawbotnotify-api -f etc/weclawbotnotify-api.yaml

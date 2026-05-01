#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
API_DIR="$PROJECT_DIR/services/weclawbotnotify-api"
ILINK_DIR="$PROJECT_DIR/services/ilink"

echo "=== 编译 ilink-rpc ==="
cd "$ILINK_DIR"
go build -o ilink-rpc ./ilink.go

echo "=== 编译 weclawbotnotify-api ==="
cd "$API_DIR"
go build -o weclawbotnotify-api ./weclawbotnotify.go

echo "=== 启动 ilink-rpc ==="
cd "$ILINK_DIR"
./ilink-rpc -f etc/ilink.yaml &
ILINK_PID=$!

sleep 2

echo "=== 启动 weclawbotnotify-api ==="
cd "$API_DIR"
./weclawbotnotify-api -f etc/weclawbotnotify-api.yaml &
API_PID=$!

echo "=== 服务已启动 ==="
echo "ilink-rpc PID: $ILINK_PID"
echo "weclawbotnotify-api PID: $API_PID"

trap "kill $ILINK_PID $API_PID 2>/dev/null; exit" INT TERM EXIT

wait

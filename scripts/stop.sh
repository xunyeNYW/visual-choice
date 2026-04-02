#!/bin/bash
# Visual Choice - 停止服务器脚本

set -e

SESSION_DIR="$HOME/.visual-choice/session"

# 检查 PID 文件
if [ ! -f "$SESSION_DIR/server.pid" ]; then
    echo "服务器未运行"
    exit 0
fi

PID=$(cat "$SESSION_DIR/server.pid")

# 检查进程是否存在
if ! ps -p $PID > /dev/null 2>&1; then
    echo "服务器进程不存在 (清理 PID 文件)"
    rm -f "$SESSION_DIR/server.pid"
    exit 0
fi

# 发送 SIGTERM 信号
echo "正在停止服务器 (PID: $PID)..."
kill -TERM $PID

# 等待进程退出
for i in {1..10}; do
    if ! ps -p $PID > /dev/null 2>&1; then
        echo "✅ 服务器已停止"
        rm -f "$SESSION_DIR/server.pid"
        exit 0
    fi
    sleep 0.5
done

# 超时强制杀死
echo "超时，强制停止服务器..."
kill -KILL $PID
rm -f "$SESSION_DIR/server.pid"
echo "✅ 服务器已强制停止"

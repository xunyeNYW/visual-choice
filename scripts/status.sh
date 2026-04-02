#!/bin/bash
# Visual Choice - 查看服务器状态脚本

set -e

SESSION_DIR="$HOME/.visual-choice/session"
PID_FILE="$SESSION_DIR/server.pid"

# 检查 PID 文件
if [ ! -f "$PID_FILE" ]; then
    echo "服务器未运行"
    exit 0
fi

PID=$(cat "$PID_FILE")

# 检查进程是否存在
if ! ps -p $PID > /dev/null 2>&1; then
    echo "服务器未运行 (进程不存在)"
    rm -f "$PID_FILE"
    exit 0
fi

# 显示状态
echo "服务器正在运行"
echo ""
echo "PID: $PID"

# 读取服务器信息
if [ -f "$SESSION_DIR/state/server-info.json" ]; then
    echo ""
    echo "连接信息:"
    cat "$SESSION_DIR/state/server-info.json" | grep -E '"url"|"port"' | sed 's/"/ /g' | sed 's/,//g'
fi

echo ""
echo "目录:"
echo "  Screen: $SESSION_DIR/screens"
echo "  State:  $SESSION_DIR/state"

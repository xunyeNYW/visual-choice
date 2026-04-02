#!/bin/bash
# Visual Choice - 启动服务器脚本

set -e

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SKILL_DIR="$(dirname "$SCRIPT_DIR")"

# Binary 文件路径
VISUAL_CHOICE_BIN="$SKILL_DIR/bin/visual-choice"
SESSION_DIR="$HOME/.visual-choice/session"

# 检查二进制文件是否存在
if [ ! -f "$VISUAL_CHOICE_BIN" ]; then
    echo "错误：找不到 visual-choice 二进制文件"
    echo "请先编译：cd /media/vdc/10226591/github/superpowers/visual-choice && go build"
    exit 1
fi

# 创建会话目录
mkdir -p "$SESSION_DIR"

# 检查服务器是否已在运行
if [ -f "$SESSION_DIR/server.pid" ]; then
    PID=$(cat "$SESSION_DIR/server.pid")
    if ps -p $PID > /dev/null 2>&1; then
        echo "服务器已在运行"
        echo "PID: $PID"
        if [ -f "$SESSION_DIR/state/server-info.json" ]; then
            URL=$(grep -o '"url":"[^"]*"' "$SESSION_DIR/state/server-info.json" | cut -d'"' -f4)
            echo "URL: $URL"
        fi
        exit 0
    else
        echo "检测到残留的 PID 文件，清理中..."
        rm -f "$SESSION_DIR/server.pid"
    fi
fi

# 启动服务器
echo "正在启动 visual-choice 服务器..."
cd "$SESSION_DIR"
"$VISUAL_CHOICE_BIN" start --port 5234 --dir "$SESSION_DIR" &

# 等待服务器启动
sleep 2

# 显示连接信息
if [ -f "$SESSION_DIR/state/server-info.json" ]; then
    echo ""
    echo "✅ 服务器已启动"
    echo ""
    cat "$SESSION_DIR/state/server-info.json" | grep -E '"url"|"screen_dir"|"state_dir"' | sed 's/"/ /g' | sed 's/,//g'
    echo ""
    echo "在浏览器访问上方 URL 开始使用"
    echo "运行 '~/.cursor/skills/visual-choice/scripts/status.sh' 查看状态"
    echo "运行 '~/.cursor/skills/visual-choice/scripts/stop.sh' 停止服务器"
else
    echo "⚠️  服务器启动失败，请检查输出"
    exit 1
fi

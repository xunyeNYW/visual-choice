#!/bin/bash
# Visual Choice - 停止服务器脚本

set -e

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SKILL_DIR="$(dirname "$SCRIPT_DIR")"

# ========================================
# 自动检测安装位置和 Session 目录
# ========================================

if [ -f "$SKILL_DIR/.install-config" ]; then
    source "$SKILL_DIR/.install-config"
else
    case "$SKILL_DIR" in
        $HOME/.cursor/skills/visual-choice|$HOME/.flow/skills/visual-choice|$HOME/.claude/skills/visual-choice)
            INSTALL_MODE="user"
            SESSION_DIR="$HOME/.visual-choice/session"
            ;;
        $HOME/.config/opencode/skills/visual-choice)
            INSTALL_MODE="user"
            SESSION_DIR="$HOME/.visual-choice/session"
            ;;
        */.cursor/skills/visual-choice|*/.flow/skills/visual-choice|*/.claude/skills/visual-choice)
            INSTALL_MODE="project"
            PROJECT_ROOT="$(dirname "$(dirname "$(dirname "$SKILL_DIR")")")"
            SESSION_DIR="$PROJECT_ROOT/.visual-choice/session"
            ;;
        */.config/opencode/skills/visual-choice)
            INSTALL_MODE="project"
            PROJECT_ROOT="$(dirname "$(dirname "$(dirname "$(dirname "$SKILL_DIR")")")")"
            SESSION_DIR="$PROJECT_ROOT/.visual-choice/session"
            ;;
        *)
            INSTALL_MODE="dev"
            SESSION_DIR="$HOME/.visual-choice/session"
            ;;
    esac
fi

SESSION_DIR="${VISUAL_CHOICE_SESSION:-$SESSION_DIR}"

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

#!/bin/bash
# Visual Choice - 查看服务器状态脚本

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

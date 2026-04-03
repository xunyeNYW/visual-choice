#!/bin/bash
# Visual Choice - 查看事件记录脚本

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

EVENTS_FILE="$SESSION_DIR/state/events.jsonl"

# 检查事件文件
if [ ! -f "$EVENTS_FILE" ]; then
    echo "暂无事件记录"
    exit 0
fi

# 检查文件是否为空
if [ ! -s "$EVENTS_FILE" ]; then
    echo "暂无事件记录"
    exit 0
fi

# 解析并显示事件
echo "事件记录:"
echo "---------"

counter=0
while IFS= read -r line; do
    if [ -z "$line" ]; then
        continue
    fi
    
    counter=$((counter + 1))
    
    # 使用 jq 解析 JSON (如果可用)
    if command -v jq &> /dev/null; then
        timestamp=$(echo "$line" | jq -r '.timestamp // 0' | xargs -I {} date -d @{} +"%H:%M:%S" 2>/dev/null || echo "00:00:00")
        event_type=$(echo "$line" | jq -r '.type // "click"')
        choice=$(echo "$line" | jq -r '.choice // ""')
        text=$(echo "$line" | jq -r '.text // ""')
        
        echo "[$counter] $timestamp $event_type - 选择：$choice - $text"
    else
        # 无 jq 时的简单输出
        echo "[$counter] $line"
    fi
done < "$EVENTS_FILE"

echo ""
echo "总事件数：$counter"

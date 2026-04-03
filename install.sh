#!/bin/bash
# Visual Choice - 安装脚本
# 支持用户级和工程级两种安装方式

set -e

INSTALL_MODE="${1:-user}"  # user 或 project
PLATFORM="${2:-auto}"      # cursor, claude, opencode, flow, 或 auto
SKILL_SOURCE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 显示帮助信息
show_help() {
    echo "Visual Choice 安装脚本"
    echo ""
    echo "用法: $0 [安装模式] [平台]"
    echo ""
    echo "安装模式:"
    echo "  user    - 用户级安装 (默认)"
    echo "            安装到用户级 skills 目录"
    echo "            Session 数据: ~/.visual-choice/session/"
    echo ""
    echo "  project - 工程级安装"
    echo "            安装到项目级 skills 目录"
    echo "            Session 数据: 项目/.visual-choice/session/"
    echo ""
    echo "平台及路径:"
    echo "  cursor   - Cursor IDE"
    echo "             用户级: ~/.cursor/skills/"
    echo "             工程级: 项目/.cursor/skills/"
    echo ""
    echo "  flow     - Flow IDE"
    echo "             用户级: ~/.flow/skills/"
    echo "             工程级: 项目/.flow/skills/"
    echo ""
    echo "  claude   - Claude Code CLI"
    echo "             用户级: ~/.claude/skills/"
    echo "             工程级: 项目/.claude/skills/"
    echo ""
    echo "  opencode - OpenCode"
    echo "             用户级: ~/.config/opencode/skills/"
    echo "             工程级: 项目/.config/opencode/skills/"
    echo ""
    echo "  auto     - 自动检测 (默认)"
    echo ""
    echo "示例:"
    echo "  $0                    # 用户级安装，自动检测平台"
    echo "  $0 user cursor        # 用户级安装，指定 Cursor"
    echo "  $0 user flow          # 用户级安装，指定 Flow IDE"
    echo "  $0 user opencode      # 用户级安装，指定 OpenCode"
    echo "  $0 project            # 工程级安装，自动检测平台"
    echo "  $0 project claude     # 工程级安装，指定 Claude Code"
    exit 0
}

# 获取平台安装路径
get_platform_path() {
    local platform="$1"
    local mode="$2"  # user 或 project

    case "$platform" in
        opencode)
            # OpenCode 使用 ~/.config/opencode/ 目录
            if [ "$mode" = "user" ]; then
                echo "$HOME/.config/opencode/skills"
            else
                PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
                echo "$PROJECT_ROOT/.config/opencode/skills"
            fi
            ;;
        cursor|claude|flow)
            # Cursor、Claude、Flow 使用 ~/.xxx/ 目录
            if [ "$mode" = "user" ]; then
                echo "$HOME/.$platform/skills"
            else
                PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
                echo "$PROJECT_ROOT/.$platform/skills"
            fi
            ;;
    esac
}

# 检测平台
detect_platform() {
    local context="$1"  # user 或 project

    if [ "$context" = "project" ]; then
        # 工程级：检查项目目录
        PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
        if [ -d "$PROJECT_ROOT/.cursor" ]; then echo "cursor"
        elif [ -d "$PROJECT_ROOT/.flow" ]; then echo "flow"
        elif [ -d "$PROJECT_ROOT/.claude" ]; then echo "claude"
        elif [ -d "$PROJECT_ROOT/.config/opencode" ]; then echo "opencode"
        else echo "cursor"  # 默认
        fi
    else
        # 用户级：检查用户目录
        if [ -d "$HOME/.cursor" ]; then echo "cursor"
        elif [ -d "$HOME/.flow" ]; then echo "flow"
        elif [ -d "$HOME/.claude" ]; then echo "claude"
        elif [ -d "$HOME/.config/opencode" ]; then echo "opencode"
        else echo "cursor"  # 默认
        fi
    fi
}

# 处理帮助请求
if [ "$INSTALL_MODE" = "-h" ] || [ "$INSTALL_MODE" = "--help" ]; then
    show_help
fi

# 验证安装模式
if [ "$INSTALL_MODE" != "user" ] && [ "$INSTALL_MODE" != "project" ]; then
    echo "错误：安装模式必须是 'user' 或 'project'"
    echo "运行 '$0 --help' 查看帮助"
    exit 1
fi

# 自动检测平台
if [ "$PLATFORM" = "auto" ]; then
    PLATFORM=$(detect_platform "$INSTALL_MODE")
fi

# 验证平台
if [ "$PLATFORM" != "cursor" ] && [ "$PLATFORM" != "flow" ] && [ "$PLATFORM" != "claude" ] && [ "$PLATFORM" != "opencode" ]; then
    echo "错误：平台必须是 'cursor', 'flow', 'claude' 或 'opencode'"
    exit 1
fi

# 计算安装路径
SKILLS_PATH=$(get_platform_path "$PLATFORM" "$INSTALL_MODE")
INSTALL_DIR="$SKILLS_PATH/visual-choice"
SESSION_DIR="$HOME/.visual-choice/session"

if [ "$INSTALL_MODE" = "project" ]; then
    PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
    SESSION_DIR="$PROJECT_ROOT/.visual-choice/session"
fi

# 显示安装信息
echo "========================================"
echo "Visual Choice 安装"
echo "========================================"
echo ""
echo "安装模式: $INSTALL_MODE"
echo "平台:     $PLATFORM"
echo "源目录:   $SKILL_SOURCE"
echo "目标目录: $INSTALL_DIR"
echo "Session:  $SESSION_DIR"
echo ""

# 检查源文件
if [ ! -f "$SKILL_SOURCE/bin/visual-choice" ]; then
    echo "错误：找不到二进制文件 $SKILL_SOURCE/bin/visual-choice"
    echo "请先编译：cd $SKILL_SOURCE && make build"
    exit 1
fi

if [ ! -f "$SKILL_SOURCE/SKILL.md" ]; then
    echo "错误：找不到 SKILL.md 文件"
    exit 1
fi

# 创建目标目录
echo "创建目录..."
mkdir -p "$INSTALL_DIR"
mkdir -p "$INSTALL_DIR/bin"
mkdir -p "$INSTALL_DIR/scripts"

# 复制文件
echo "复制文件..."
cp "$SKILL_SOURCE/bin/visual-choice" "$INSTALL_DIR/bin/"
cp -r "$SKILL_SOURCE/scripts/"*.sh "$INSTALL_DIR/scripts/"
cp "$SKILL_SOURCE/SKILL.md" "$INSTALL_DIR/"

# 设置执行权限
chmod +x "$INSTALL_DIR/bin/visual-choice"
chmod +x "$INSTALL_DIR/scripts/"*.sh

# 记录安装配置
cat > "$INSTALL_DIR/.install-config" << EOF
# Visual Choice 安装配置
# 此文件由 install.sh 自动生成，请勿手动修改

INSTALL_MODE=$INSTALL_MODE
PLATFORM=$PLATFORM
INSTALL_DIR=$INSTALL_DIR
SESSION_DIR=$SESSION_DIR
INSTALLED_AT=$(date +%Y-%m-%d_%H:%M:%S)
SKILL_SOURCE=$SKILL_SOURCE
EOF

# 创建 Session 目录（如果不存在）
mkdir -p "$SESSION_DIR/screens"
mkdir -p "$SESSION_DIR/state"

# 工程级安装：添加 .gitignore 建议
if [ "$INSTALL_MODE" = "project" ]; then
    GITIGNORE_FILE="$PROJECT_ROOT/.gitignore"

    # 检查是否已有 visual-choice 相关规则
    if [ -f "$GITIGNORE_FILE" ] && ! grep -q ".visual-choice" "$GITIGNORE_FILE"; then
        echo ""
        echo "建议添加以下规则到 .gitignore:"
        echo "----------------------------------------"
        echo "# Visual Choice session 数据"
        echo ".visual-choice/session/state/events.jsonl"
        echo ".visual-choice/session/server.pid"
        echo ".visual-choice/session/state/server-info.json"
        echo "----------------------------------------"
    fi
fi

# 验证安装
echo ""
echo "验证安装..."
if [ -f "$INSTALL_DIR/bin/visual-choice" ]; then
    echo "✅ 二进制文件已安装"
else
    echo "❌ 二进制文件安装失败"
    exit 1
fi

if [ -f "$INSTALL_DIR/SKILL.md" ]; then
    echo "✅ SKILL.md 已安装"
else
    echo "❌ SKILL.md 安装失败"
    exit 1
fi

if [ -f "$INSTALL_DIR/.install-config" ]; then
    echo "✅ 安装配置已生成"
else
    echo "❌ 安装配置生成失败"
    exit 1
fi

# 完成
echo ""
echo "========================================"
echo "✅ 安装完成!"
echo "========================================"
echo ""
echo "使用方法:"
echo "  启动服务器: $INSTALL_DIR/scripts/start.sh"
echo "  查看状态:   $INSTALL_DIR/scripts/status.sh"
echo "  查看事件:   $INSTALL_DIR/scripts/events.sh"
echo "  停止服务器: $INSTALL_DIR/scripts/stop.sh"
echo ""
echo "在 AI Agent 中使用:"
echo "  /visual-choice"
echo ""
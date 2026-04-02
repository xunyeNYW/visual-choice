#!/bin/bash
# Visual Choice - 跨平台部署脚本
# 一键部署到所有支持的 AI Agent 平台

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 源目录
SOURCE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 部署模式
DEPLOY_MODE="${1:-symlink}"  # symlink | copy

echo -e "${BLUE}╔════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   Visual Choice - 跨平台 Skill 部署               ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════╝${NC}"
echo ""
echo "源目录：$SOURCE_DIR"
echo "部署模式：$DEPLOY_MODE"
echo ""

# 检查源目录
if [ ! -f "$SOURCE_DIR/SKILL.md" ]; then
    echo -e "${RED}❌ 错误：找不到 SKILL.md${NC}"
    echo "请确保从 skill 目录运行此脚本"
    exit 1
fi

# 平台列表
declare -A PLATFORMS=(
    ["~/.claude/skills"]="Claude Code"
    ["~/.config/opencode/skills"]="OpenCode"
)

# 部署计数器
SUCCESS=0
SKIPPED=0
FAILED=0

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

for dir in "${!PLATFORMS[@]}"; do
    platform="${PLATFORMS[$dir]}"
    target="$dir/visual-choice"
    
    echo -e "${YELLOW}📦 部署到：${platform}${NC}"
    echo "   目标目录：$target"
    
    # 展开 ~
    target="${target/#\~/$HOME}"
    
    # 创建父目录
    mkdir -p "$(dirname "$target")"
    
    if [ "$DEPLOY_MODE" = "symlink" ]; then
        # 符号链接模式
        if [ -L "$target" ]; then
            # 已存在符号链接
            existing=$(readlink "$target")
            if [ "$existing" = "$SOURCE_DIR" ]; then
                echo -e "   ${GREEN}✅ 已存在符号链接，指向正确${NC}"
            else
                echo -e "   ${YELLOW}⚠️  符号链接指向其他位置：$existing${NC}"
                echo -e "   ${YELLOW}   跳过（手动检查冲突）${NC}"
                SKIPPED=$((SKIPPED + 1))
            fi
        elif [ -d "$target" ]; then
            # 已存在目录
            echo -e "   ${YELLOW}⚠️  目录已存在${NC}"
            
            # 检查是否是符号链接目录
            if [ -L "$target/SKILL.md" ]; then
                echo -e "   ${GREEN}✓ 混合部署（目录 + 符号链接文件）${NC}"
            else
                echo -e "   ${YELLOW}   跳过（避免覆盖现有数据）${NC}"
                SKIPPED=$((SKIPPED + 1))
            fi
        else
            # 创建新符号链接
            ln -s "$SOURCE_DIR" "$target"
            echo -e "   ${GREEN}✅ 符号链接已创建${NC}"
            SUCCESS=$((SUCCESS + 1))
        fi
    else
        # 复制模式
        if [ -d "$target" ]; then
            # 更新现有副本
            echo -e "   ${BLUE}🔄 更新现有副本...${NC}"
            cp -r "$SOURCE_DIR"/* "$target/"
            echo -e "   ${GREEN}✅ 副本已更新${NC}"
            SUCCESS=$((SUCCESS + 1))
        else
            # 创建新副本
            echo -e "   ${BLUE}📋 创建新副本...${NC}"
            cp -r "$SOURCE_DIR" "$target"
            echo -e "   ${GREEN}✅ 副本已创建${NC}"
            SUCCESS=$((SUCCESS + 1))
        fi
    fi
    
    echo ""
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# 验证部署
echo -e "${BLUE}🔍 验证部署...${NC}"
echo ""

verify_deployment() {
    local platform_dir="$1"
    local platform_name="$2"
    
    platform_dir="${platform_dir/#\~/$HOME}"
    target="$platform_dir/visual-choice"
    
    if [ -d "$target" ] || [ -L "$target" ]; then
        # 检查关键文件
        if [ -f "$target/SKILL.md" ] && [ -f "$target/bin/visual-choice" ]; then
            echo -e "   ${GREEN}✅ $platform_name: 验证通过${NC}"
            return 0
        else
            echo -e "   ${RED}❌ $platform_name: 文件不完整${NC}"
            return 1
        fi
    else
        echo -e "   ${YELLOW}⚠️  $platform_name: 未部署${NC}"
        return 1
    fi
}

verify_deployment "~/.claude/skills" "Claude Code"
verify_deployment "~/.config/opencode/skills" "OpenCode"
verify_deployment "~/.cursor/skills" "Cursor"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# 输出摘要
echo -e "${BLUE}📊 部署摘要:${NC}"
echo ""
echo -e "   成功：${GREEN}$SUCCESS${NC}"
echo -e "   跳过：${YELLOW}$SKIPPED${NC}"
echo -e "   失败：${RED}$FAILED${NC}"
echo ""

# 输出使用方式
if [ $SUCCESS -gt 0 ] || [ $SKIPPED -gt 0 ]; then
    echo -e "${BLUE}🚀 使用方式:${NC}"
    echo ""
    echo "   Cursor:      ${GREEN}/visual-choice${NC}"
    echo "   Claude Code: ${GREEN}/visual-choice${NC}"
    echo "   OpenCode:    ${GREEN}/visual-choice${NC}"
    echo ""
    
    # 测试 binary
    if [ -f "$SOURCE_DIR/bin/visual-choice" ]; then
        echo -e "${BLUE}🧪 测试 Binary:${NC}"
        echo ""
        "$SOURCE_DIR/bin/visual-choice" --help 2>&1 | head -5
        echo ""
    fi
fi

# 输出下一步
echo -e "${BLUE}📝 下一步:${NC}"
echo ""
echo "   1. 在 AI Agent 中测试：${GREEN}/visual-choice${NC}"
echo "   2. 启动服务器：${GREEN}./scripts/start.sh${NC}"
echo "   3. 查看文档：${GREEN}cat CROSS-PLATFORM.md${NC}"
echo ""

# 退出码
if [ $FAILED -gt 0 ]; then
    exit 1
else
    exit 0
fi

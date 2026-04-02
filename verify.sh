#!/bin/bash
# Visual Choice Skill 验证脚本

set -e

SKILL_DIR="$HOME/.cursor/skills/visual-choice"
ERRORS=0
WARNINGS=0

echo "======================================"
echo "Visual Choice Skill 验证"
echo "======================================"
echo ""

# 1. 检查目录结构
echo "1. 检查目录结构..."
required_files=("SKILL.md" "README.md" "reference.md" "examples.md" "MAINTENANCE.md" "ASSETS.md" "CROSS-PLATFORM.md" "deploy-cross-platform.sh" "scripts/start.sh" "scripts/stop.sh" "scripts/events.sh" "scripts/status.sh" "bin/visual-choice" "src/go.mod" "src/README.md" "src/Makefile" "src/BUILD.md")

for file in "${required_files[@]}"; do
    if [ -f "$SKILL_DIR/$file" ]; then
        echo "   ✅ $file"
    else
        echo "   ❌ $file (缺失)"
        ERRORS=$((ERRORS + 1))
    fi
done

# 检查新的源代码结构
echo ""
echo "   检查源代码结构..."
src_files=("src/cmd/visual-choice/main.go" "src/internal/server/server.go" "src/internal/server/template.go" "src/internal/events/events.go" "src/internal/models/models.go")

for file in "${src_files[@]}"; do
    if [ -f "$SKILL_DIR/$file" ]; then
        echo "   ✅ $file"
    else
        echo "   ❌ $file (缺失)"
        ERRORS=$((ERRORS + 1))
    fi
done

# 检查测试文件
echo ""
echo "   检查测试文件..."
test_files=("src/internal/events/events_test.go" "src/internal/models/models_test.go")

for file in "${test_files[@]}"; do
    if [ -f "$SKILL_DIR/$file" ]; then
        echo "   ✅ $file"
    else
        echo "   ⚠️  $file (可选)"
    fi
done

# 2. 检查 SKILL.md 质量
echo ""
echo "2. 检查 SKILL.md 质量..."

# 检查行数
lines=$(wc -l < "$SKILL_DIR/SKILL.md")
if [ $lines -le 500 ]; then
    echo "   ✅ 行数：$lines (< 500)"
else
    echo "   ⚠️  行数：$lines (> 500，建议拆分)"
    WARNINGS=$((WARNINGS + 1))
fi

# 检查 frontmatter
if grep -q "^---$" "$SKILL_DIR/SKILL.md" && grep -q "^name:" "$SKILL_DIR/SKILL.md" && grep -q "^description:" "$SKILL_DIR/SKILL.md"; then
    echo "   ✅ YAML frontmatter 完整"
else
    echo "   ❌ YAML frontmatter 缺失"
    ERRORS=$((ERRORS + 1))
fi

# 检查 name 字段
name=$(grep "^name:" "$SKILL_DIR/SKILL.md" | cut -d':' -f2 | tr -d ' ')
if [ ${#name} -le 64 ] && [[ "$name" =~ ^[a-z0-9-]+$ ]]; then
    echo "   ✅ name 格式正确：$name"
else
    echo "   ❌ name 格式错误 (应小写，短横线分隔，<64 字符): $name"
    ERRORS=$((ERRORS + 1))
fi

# 检查 description 字段
description=$(grep "^description:" "$SKILL_DIR/SKILL.md" | cut -d':' -f2-)
if [ -n "$description" ]; then
    echo "   ✅ description 非空"
    
    # 检查是否包含触发词
    if [[ "$description" =~ (设计 | 原型 | 评审 | 投票 | 选择 | 浏览器) ]]; then
        echo "   ✅ description 包含触发词"
    else
        echo "   ⚠️  description 缺少触发词"
        WARNINGS=$((WARNINGS + 1))
    fi
    
    # 检查是否第三人称
    if [[ "$description" =~ (我 | 你 | 您) ]]; then
        echo "   ⚠️  description 包含第二人称 (应第三人称)"
        WARNINGS=$((WARNINGS + 1))
    else
        echo "   ✅ description 第三人称"
    fi
else
    echo "   ❌ description 为空"
    ERRORS=$((ERRORS + 1))
fi

# 3. 检查脚本可执行性
echo ""
echo "3. 检查脚本可执行性..."

for script in "$SKILL_DIR"/scripts/*.sh; do
    if [ -x "$script" ]; then
        echo "   ✅ $(basename $script) 可执行"
    else
        echo "   ❌ $(basename $script) 不可执行"
        ERRORS=$((ERRORS + 1))
    fi
done

# 4. 检查文件引用
echo ""
echo "4. 检查文件引用..."

# 检查 SKILL.md 中的引用是否存在
refs=$(grep -oE '\[[^]]+\]\([^)]+\)' "$SKILL_DIR/SKILL.md" | grep -oE '\([^)]+\)' | tr -d '()')
for ref in $refs; do
    # 跳过外部链接 (http 开头)
    if [[ "$ref" =~ ^http ]]; then
        continue
    fi
    
    if [ -f "$SKILL_DIR/$ref" ]; then
        echo "   ✅ 引用存在：$ref"
    else
        echo "   ❌ 引用缺失：$ref"
        ERRORS=$((ERRORS + 1))
    fi
done

# 5. 检查二进制文件
echo ""
echo "5. 检查 binary 文件..."

VISUAL_CHOICE_BIN="$SKILL_DIR/bin/visual-choice"
if [ -f "$VISUAL_CHOICE_BIN" ]; then
    echo "   ✅ Binary 文件存在：$VISUAL_CHOICE_BIN"
    
    if [ -x "$VISUAL_CHOICE_BIN" ]; then
        echo "   ✅ Binary 文件可执行"
    else
        echo "   ⚠️  Binary 文件不可执行"
        WARNINGS=$((WARNINGS + 1))
    fi
else
    echo "   ❌ Binary 文件不存在"
    echo "      请确保 bin/visual-choice 已编译并复制到 skill 目录"
    ERRORS=$((ERRORS + 1))
fi

# 6. 检查术语一致性
echo ""
echo "6. 检查术语一致性..."

# 检查是否有 Windows 风格路径
if grep -q '\\\\' "$SKILL_DIR/SKILL.md"; then
    echo "   ⚠️  发现 Windows 风格路径 (应使用 /)"
    WARNINGS=$((WARNINGS + 1))
else
    echo "   ✅ 路径格式正确 (Unix 风格)"
fi

# 7. 检查源代码完整性
echo ""
echo "7. 检查源代码完整性..."

# 检查新的源代码结构
if [ -f "$SKILL_DIR/src/cmd/visual-choice/main.go" ] && \
   [ -f "$SKILL_DIR/src/internal/server/server.go" ] && \
   [ -f "$SKILL_DIR/src/internal/events/events.go" ] && \
   [ -f "$SKILL_DIR/src/internal/models/models.go" ]; then
    echo "   ✅ 源代码文件完整 (新结构)"
    
    # 检查是否可以编译
    cd "$SKILL_DIR/src"
    if go build -o /tmp/visual-choice-test ./cmd/visual-choice 2>/dev/null; then
        echo "   ✅ 源代码可编译"
        rm -f /tmp/visual-choice-test
    else
        echo "   ⚠️  源代码编译失败 (可能需要更新依赖)"
        WARNINGS=$((WARNINGS + 1))
    fi
else
    echo "   ❌ 源代码文件缺失 (新结构)"
    ERRORS=$((ERRORS + 1))
fi

# 总结
echo ""
echo "======================================"
echo "验证结果"
echo "======================================"
echo "错误：$ERRORS"
echo "警告：$WARNINGS"
echo ""

if [ $ERRORS -eq 0 ]; then
    echo "✅ Skill 验证通过！"
    
    if [ $WARNINGS -gt 0 ]; then
        echo "⚠️  有 $WARNINGS 个警告，建议修复"
    fi
    
    echo ""
    echo "使用方法:"
    echo "  1. 在 Cursor 中调用 /visual-choice"
    echo "  2. 或手动运行：~/.cursor/skills/visual-choice/scripts/start.sh"
    exit 0
else
    echo "❌ Skill 验证失败，请修复 $ERRORS 个错误"
    exit 1
fi

# Visual Choice - 跨平台 Skill 部署指南

## 支持的 AI Agent 平台

Visual Choice 已设计为跨平台通用 skill，支持以下 AI Agent：

| 平台 | 状态 | Skill 目录 | 说明 |
|------|------|-----------|------|
| **Cursor** | ✅ 已部署 | `~/.cursor/skills/visual-choice/` | 当前安装位置 |
| **Claude Code** | ✅ 兼容 | `~/.claude/skills/visual-choice/` | 使用相同 SKILL.md 格式 |
| **OpenCode** | ✅ 兼容 | `~/.config/opencode/skills/visual-choice/` | 使用相同 SKILL.md 格式 |
| **其他平台** | 🟡 待验证 | - | 支持 SKILL.md 格式的平台均可 |

## 跨平台设计原理

### 1. 通用 SKILL.md 格式

所有平台都使用相同的 YAML frontmatter + Markdown body 格式：

```markdown
---
name: visual-choice
description: 浏览器可视化选择工具...
---

# 技能说明

## 快速开始
...
```

### 2. 平台无关的实现

- **Binary 文件**：Go 编译的独立可执行文件，无运行时依赖
- **Shell 脚本**：使用 POSIX shell，跨平台兼容
- **文档**：纯 Markdown，所有平台通用

### 3. 目录结构标准化

```
visual-choice/
├── SKILL.md           # 通用技能定义（所有平台相同）
├── reference.md       # 技术参考
├── examples.md        # 使用示例
├── MAINTENANCE.md     # 维护指南
├── ASSETS.md          # 资产清单
├── scripts/           # 辅助脚本（POSIX shell）
│   ├── start.sh
│   ├── stop.sh
│   ├── status.sh
│   └── events.sh
├── bin/              # Binary 目录
│   └── visual-choice # Go 编译产物
└── src/              # 源代码（可选，用于参考）
    ├── main.go
    ├── server.go
    └── events.go
```

## 部署到不同平台

### 方法 1: 符号链接（推荐，单一来源）

使用符号链接让所有平台共享同一份 skill：

```bash
# 1. 源代码保持在当前目录
SOURCE_DIR="$HOME/.cursor/skills/visual-choice"

# 2. 创建 Claude Code 符号链接
mkdir -p ~/.claude/skills
ln -s "$SOURCE_DIR" ~/.claude/skills/visual-choice

# 3. 创建 OpenCode 符号链接
mkdir -p ~/.config/opencode/skills
ln -s "$SOURCE_DIR" ~/.config/opencode/skills/visual-choice

# 4. 验证
ls -la ~/.claude/skills/visual-choice
ls -la ~/.config/opencode/skills/visual-choice
```

**优点**：
- ✅ 单一来源，更新一次所有平台生效
- ✅ 不占用额外磁盘空间
- ✅ 易于维护

**缺点**：
- ⚠️ 某些平台可能不支持符号链接（罕见）

### 方法 2: 完整复制（独立部署）

为每个平台复制独立副本：

```bash
# 1. 部署到 Claude Code
cp -r ~/.cursor/skills/visual-choice ~/.claude/skills/visual-choice

# 2. 部署到 OpenCode
cp -r ~/.cursor/skills/visual-choice ~/.config/opencode/skills/visual-choice

# 3. 验证部署
~/.claude/skills/visual-choice/verify.sh
~/.config/opencode/skills/visual-choice/verify.sh
```

**优点**：
- ✅ 完全独立，互不影响
- ✅ 每个平台可以自定义配置

**缺点**：
- ❌ 占用更多磁盘空间（~8.6MB × 平台数）
- ❌ 更新需要同步所有副本

### 方法 3: 混合部署（推荐用于生产）

核心文件 symlink，平台特定文件独立：

```bash
# 1. 创建基础目录
mkdir -p ~/.claude/skills/visual-choice

# 2.  symlink 通用文件
ln -s ~/.cursor/skills/visual-choice/SKILL.md ~/.claude/skills/visual-choice/SKILL.md
ln -s ~/.cursor/skills/visual-choice/reference.md ~/.claude/skills/visual-choice/reference.md
ln -s ~/.cursor/skills/visual-choice/examples.md ~/.claude/skills/visual-choice/examples.md
ln -s ~/.cursor/skills/visual-choice/scripts ~/.claude/skills/visual-choice/scripts
ln -s ~/.cursor/skills/visual-choice/bin ~/.claude/skills/visual-choice/bin

# 3. 复制平台特定配置（如有）
cp ~/.cursor/skills/visual-choice/MAINTENANCE.md ~/.claude/skills/visual-choice/MAINTENANCE.md
```

## 平台特定配置

### Cursor

**位置**: `~/.cursor/skills/visual-choice/`

**特点**:
- 已部署，可直接使用
- 支持 `/visual-choice` 命令
- 自动发现技能描述

### Claude Code

**位置**: `~/.claude/skills/visual-choice/`

**部署命令**:
```bash
# 使用符号链接
ln -s ~/.cursor/skills/visual-choice ~/.claude/skills/visual-choice

# 或复制
cp -r ~/.cursor/skills/visual-choice ~/.claude/skills/visual-choice
```

**使用方式**:
```
/visual-choice
```

**特点**:
- 使用完全相同的 SKILL.md 格式
- 自动发现机制相同
- 支持 progressive disclosure

### OpenCode

**位置**: `~/.config/opencode/skills/visual-choice/`

**部署命令**:
```bash
# 使用符号链接
ln -s ~/.cursor/skills/visual-choice ~/.config/opencode/skills/visual-choice

# 或复制
cp -r ~/.cursor/skills/visual-choice ~/.config/opencode/skills/visual-choice
```

**使用方式**:
```
/visual-choice
```

**特点**:
- 支持 SKILL.md 格式
- 自动从 git worktree root 向上搜索
- 也兼容 `.claude/skills/` 目录

## 跨平台二进制编译

### macOS (Apple Silicon + Intel)

```bash
cd ~/.cursor/skills/visual-choice/src

# Apple Silicon (M1/M2/M3)
GOOS=darwin GOARCH=arm64 go build -o ../bin/visual-choice-darwin-arm64

# Intel Mac
GOOS=darwin GOARCH=amd64 go build -o ../bin/visual-choice-darwin-amd64

# 通用二进制（推荐）
lipo -create -output ../bin/visual-choice \
  ../bin/visual-choice-darwin-arm64 \
  ../bin/visual-choice-darwin-amd64
```

### Linux (x86_64 + ARM)

```bash
cd ~/.cursor/skills/visual-choice/src

# x86_64
GOOS=linux GOARCH=amd64 go build -o ../bin/visual-choice-linux-amd64

# ARM64
GOOS=linux GOARCH=arm64 go build -o ../bin/visual-choice-linux-arm64
```

### Windows

```bash
cd ~/.cursor/skills/visual-choice/src
GOOS=windows GOARCH=amd64 go build -o ../bin/visual-choice.exe
```

### 多平台分发

创建 releases 目录：

```bash
mkdir -p ~/.cursor/skills/visual-choice/releases

# 复制所有平台版本
cp ../bin/visual-choice-darwin-arm64 releases/
cp ../bin/visual-choice-darwin-amd64 releases/
cp ../bin/visual-choice-linux-amd64 releases/
cp ../bin/visual-choice-linux-arm64 releases/
cp ../bin/visual-choice.exe releases/

# 创建下载说明
cat > releases/README.md << 'EOF'
# Visual Choice Binaries

## macOS
- `visual-choice-darwin-arm64` - Apple Silicon (M1/M2/M3)
- `visual-choice-darwin-amd64` - Intel Mac

## Linux
- `visual-choice-linux-amd64` - x86_64
- `visual-choice-linux-arm64` - ARM64

## Windows
- `visual-choice.exe` - x86_64

## 验证
```bash
shasum -a 256 visual-choice-*
```
EOF
```

## 自动化部署脚本

创建跨平台部署脚本：

```bash
#!/bin/bash
# deploy-cross-platform.sh
# 一键部署到所有支持的 AI Agent 平台

set -e

SOURCE_DIR="$HOME/.cursor/skills/visual-choice"
DEPLOY_MODE="${1:-symlink}"  # symlink | copy

echo "🚀 部署 Visual Choice Skill"
echo "模式：$DEPLOY_MODE"
echo ""

# 平台列表
declare -A PLATFORMS=(
    ["~/.claude/skills"]="Claude Code"
    ["~/.config/opencode/skills"]="OpenCode"
)

for dir in "${!PLATFORMS[@]}"; do
    platform="${PLATFORMS[$dir]}"
    target="$dir/visual-choice"
    
    echo "📦 部署到 $platform ($dir)"
    
    # 创建父目录
    mkdir -p "$dir"
    
    if [ "$DEPLOY_MODE" = "symlink" ]; then
        # 符号链接模式
        if [ -L "$target" ]; then
            echo "   ✅ 已存在符号链接"
        elif [ -d "$target" ]; then
            echo "   ⚠️  目录已存在，跳过"
        else
            ln -s "$SOURCE_DIR" "$target"
            echo "   ✅ 符号链接已创建"
        fi
    else
        # 复制模式
        if [ -d "$target" ]; then
            echo "   🔄 更新现有副本"
            cp -r "$SOURCE_DIR"/* "$target/"
        else
            echo "   📋 创建新副本"
            cp -r "$SOURCE_DIR" "$target"
        fi
    fi
    
    echo ""
done

echo "✅ 部署完成！"
echo ""
echo "使用方式:"
echo "  Cursor:     /visual-choice"
echo "  Claude Code: /visual-choice"
echo "  OpenCode:   /visual-choice"
```

## 验证跨平台兼容性

### 检查清单

- [ ] SKILL.md frontmatter 符合所有平台规范
- [ ] 名称格式：`^[a-z0-9]+(-[a-z0-9]+)*$`
- [ ] 描述长度：< 1024 字符
- [ ] 路径格式：Unix 风格（`/`）
- [ ] Shell 脚本：POSIX 兼容
- [ ] Binary：目标平台架构匹配

### 测试命令

```bash
# 1. 验证 SKILL.md 格式
cd ~/.cursor/skills/visual-choice
./verify.sh

# 2. 测试符号链接
ls -la ~/.claude/skills/visual-choice
ls -la ~/.config/opencode/skills/visual-choice

# 3. 测试 binary
~/.cursor/skills/visual-choice/bin/visual-choice --help

# 4. 测试脚本
~/.cursor/skills/visual-choice/scripts/start.sh
```

## 更新流程

### 单一来源更新（符号链接模式）

```bash
# 1. 修改源代码
cd ~/.cursor/skills/visual-choice/src
vi main.go

# 2. 重新编译
go build -o ../bin/visual-choice

# 3. 验证
../verify.sh

# ✅ 所有平台自动获得更新！
```

### 多副本同步（复制模式）

```bash
# 1. 更新主副本
cd ~/.cursor/skills/visual-choice
# ... 修改和编译 ...

# 2. 同步到其他平台
cp -r ~/.cursor/skills/visual-choice ~/.claude/skills/visual-choice
cp -r ~/.cursor/skills/visual-choice ~/.config/opencode/skills/visual-choice

# 3. 验证所有平台
~/.claude/skills/visual-choice/verify.sh
~/.config/opencode/skills/visual-choice/verify.sh
```

## 故障排查

### 问题：技能未被发现

**解决**：
1. 检查 SKILL.md frontmatter 格式
2. 验证 name 字段格式
3. 确认目录名称与 name 匹配
4. 重启 AI Agent

### 问题：Binary 无法执行

**解决**：
```bash
# 检查架构
file ~/.cursor/skills/visual-choice/bin/visual-choice

# 重新编译正确架构
cd ~/.cursor/skills/visual-choice/src
GOOS=darwin GOARCH=arm64 go build -o ../bin/visual-choice  # macOS M1
```

### 问题：脚本权限错误

**解决**：
```bash
chmod +x ~/.cursor/skills/visual-choice/scripts/*.sh
chmod +x ~/.cursor/skills/visual-choice/bin/visual-choice
```

## 最佳实践

1. **使用符号链接** - 单一来源，易于维护
2. **编译多平台 binary** - 支持不同操作系统的用户
3. **保持 SKILL.md 简洁** - < 500 行，progressive disclosure
4. **使用 POSIX shell** - 避免 bash 特定语法
5. **文档通用化** - 避免平台特定术语
6. **定期验证** - 在所有目标平台测试

## 参考资料

- [Claude Code Skills](https://code.claude.com/docs/en/skills.md)
- [OpenCode Skills](https://open-code.ai/en/docs/skills)
- [Go Cross Compilation](https://go.dev/doc/install/source#environment)

---

**下一步**: 运行自动化部署脚本

```bash
# 符号链接模式（推荐）
./deploy-cross-platform.sh symlink

# 复制模式
./deploy-cross-platform.sh copy
```

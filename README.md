# Visual Choice - 跨平台可视化选择工具

[![Platform](https://img.shields.io/badge/platform-Cursor%20%7C%20Claude%20Code%20%7C%20OpenCode-blue)]()
[![Skill Format](https://img.shields.io/badge/format-SKILL.md-green)]()
[![Build](https://img.shields.io/badge/build-Go-yellow)]()

## 📖 概述

**Visual Choice** 是一个跨平台的浏览器可视化选择工具，用于在浏览器中展示设计选项、原型、架构方案等，并通过点击交互收集用户反馈。

### 核心特性

- ✅ **跨平台支持** - Cursor、Claude Code、OpenCode 通用
- ✅ **实时文件监听** - 写入 HTML 即生效，无需重启
- ✅ **自动模板包装** - 只需写内容片段，自动包装完整页面
- ✅ **交互记录** - 用户点击自动保存到 JSONL 文件
- ✅ **多种布局模板** - options/cards/split/pros-cons 等
- ✅ **进程管理** - PID 文件管理，优雅启停

### 使用场景

| 场景 | 说明 | 示例 |
|------|------|------|
| 🎨 设计评审 | Logo、配色、布局投票 | 3 个 Logo 方案团队投票 |
| 📱 原型演示 | 产品原型展示和测试 | Dashboard 原型用户反馈 |
| 👥 用户调研 | 一对一访谈，偏好收集 | 功能优先级访谈 |
| 📊 产品规划 | 功能优先级排序 | Sprint 规划投票 |
| 🏗️ 技术决策 | 架构方案对比 | 微服务 vs 单体评估 |
| 🗳️ 团队投票 | 团队决策，意见收集 | 技术选型投票 |

## 🚀 快速开始

### 1. 部署到 AI Agent

```bash
# 进入 skill 目录
cd ~/.cursor/skills/visual-choice

# 一键部署到所有平台（符号链接模式）
./deploy-cross-platform.sh symlink
```

### 2. 在 AI Agent 中使用

```
/visual-choice
```

### 3. 启动服务器

```bash
# 使用启动脚本
~/.cursor/skills/visual-choice/scripts/start.sh

# 输出:
# ✅ 服务器已启动
# URL: http://localhost:5234
```

### 4. 写入 HTML 内容

```bash
cat > ~/.visual-choice/session/screens/design.html << 'EOF'
<h2>喜欢哪种设计风格？</h2>
<p class="subtitle">点击卡片选择</p>

<div class="cards">
  <div class="card" data-choice="minimal" onclick="toggleSelect(this)">
    <div class="card-image" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
      <span style="color: white; font-size: 2rem;">极简</span>
    </div>
    <div class="card-body">
      <h3>极简风格</h3>
      <p>大量留白，简洁排版</p>
    </div>
  </div>
  <div class="card" data-choice="bold" onclick="toggleSelect(this)">
    <div class="card-image" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);">
      <span style="color: white; font-size: 2rem;">大胆</span>
    </div>
    <div class="card-body">
      <h3>大胆风格</h3>
      <p>鲜艳色彩，视觉冲击</p>
    </div>
  </div>
</div>
EOF
```

### 5. 在浏览器查看

打开 http://localhost:5234，点击选项进行选择。

### 6. 查看结果

```bash
~/.cursor/skills/visual-choice/scripts/events.sh

# 输出:
# 事件记录:
# ---------
# [1] 17:28:56 click - 选择：minimal - 极简风格
# [2] 17:29:03 click - 选择：bold - 大胆风格
```

## 📦 目录结构

```
visual-choice/
├── SKILL.md              # Skill 定义（跨平台通用）
├── README.md             # 本文件
├── CROSS-PLATFORM.md     # 跨平台部署指南
├── ASSETS.md             # 资产清单
├── reference.md          # 技术参考
├── examples.md           # 使用示例
├── MAINTENANCE.md        # 维护指南
├── deploy-cross-platform.sh  # 跨平台部署脚本
├── verify.sh             # 验证脚本
├── scripts/              # 辅助脚本
│   ├── start.sh         # 启动服务器
│   ├── stop.sh          # 停止服务器
│   ├── status.sh        # 查看状态
│   └── events.sh        # 查看事件
├── bin/                 # Binary 目录
│   └── visual-choice    # Go 编译产物 (~8.4MB)
└── src/                 # 源代码目录
    ├── main.go          # CLI 入口
    ├── server.go        # HTTP 服务器
    ├── events.go        # 事件处理
    ├── go.mod           # Go 模块
    └── README.md        # 源码说明
```

## 🌐 支持的平台

| 平台 | Skill 目录 | 状态 | 使用方式 |
|------|-----------|------|---------|
| **Cursor** | `~/.cursor/skills/visual-choice/` | ✅ 已部署 | `/visual-choice` |
| **Claude Code** | `~/.claude/skills/visual-choice/` | ✅ 已部署 | `/visual-choice` |
| **OpenCode** | `~/.config/opencode/skills/visual-choice/` | ✅ 已部署 | `/visual-choice` |

### 部署模式

#### 符号链接模式（推荐）

```bash
./deploy-cross-platform.sh symlink
```

- ✅ 单一来源，更新一次所有平台生效
- ✅ 不占用额外磁盘空间
- ✅ 易于维护

#### 复制模式

```bash
./deploy-cross-platform.sh copy
```

- ✅ 完全独立，互不影响
- ❌ 占用更多磁盘空间
- ❌ 更新需要同步

## 🛠️ 技术栈

| 组件 | 技术 | 说明 |
|------|------|------|
| 语言 | Go 1.21+ | 编译为单一二进制 |
| HTTP | net/http | 标准库 |
| 文件监听 | fsnotify | 跨平台文件事件 |
| 并发 | sync.Mutex + Goroutine | 轻量级并发 |
| 数据格式 | JSON/JSONL | 事件存储 |
| 前端 | 原生 HTML/CSS/JS | 零依赖 |

## 📊 性能指标

| 指标 | 数值 | 测试条件 |
|------|------|---------|
| Binary 大小 | 8.4 MB | macOS ARM64 |
| 启动时间 | < 100ms | 冷启动 |
| HTTP 响应 | < 10ms | 本地请求 |
| 文件检测 | < 50ms | 文件写入到检测 |
| 内存占用 | ~10 MB | 运行中 |

## 📚 文档导航

| 文档 | 说明 |
|------|------|
| [SKILL.md](SKILL.md) | Skill 主文档，包含快速开始和使用方式 |
| [CROSS-PLATFORM.md](CROSS-PLATFORM.md) | 跨平台部署指南，支持 Cursor/Claude Code/OpenCode |
| [reference.md](reference.md) | 技术参考，包含 CLI、API、CSS 类参考 |
| [examples.md](examples.md) | 5 个完整使用场景示例 |
| [MAINTENANCE.md](MAINTENANCE.md) | 维护指南，包含编译、更新、故障排查 |
| [ASSETS.md](ASSETS.md) | 资产清单，包含完整统计 |
| [src/README.md](src/README.md) | 源代码说明文档 |

## 🔧 常用命令

### 启动/停止

```bash
# 启动服务器
./scripts/start.sh

# 查看状态
./scripts/status.sh

# 停止服务器
./scripts/stop.sh
```

### 查看事件

```bash
# 查看事件记录
./scripts/events.sh

# 导出为 CSV
cat ~/.visual-choice/session/state/events.jsonl | jq -r '[.timestamp, .choice, .text] | @csv' > events.csv
```

### 验证

```bash
# 验证 skill 完整性
./verify.sh
```

### 重新编译

```bash
cd src
go build -o ../bin/visual-choice
../verify.sh
```

## 🎯 使用示例

### 设计评审

```bash
cat > ~/.visual-choice/session/screens/logo.html << 'EOF'
<h2>哪个 Logo 方案更合适？</h2>
<div class="cards">
  <div class="card" data-choice="a" onclick="toggleSelect(this)">
    <div class="card-image"><!-- Logo A --></div>
    <div class="card-body">
      <h3>方案 A - 极简</h3>
    </div>
  </div>
  <div class="card" data-choice="b" onclick="toggleSelect(this)">
    <div class="card-image"><!-- Logo B --></div>
    <div class="card-body">
      <h3>方案 B - 文字</h3>
    </div>
  </div>
</div>
EOF
```

### 功能优先级

```bash
cat > ~/.visual-choice/session/screens/features.html << 'EOF'
<h2>哪些功能应该优先开发？</h2>
<p class="subtitle">可选择多个（多选）</p>

<div class="options" data-multiselect>
  <div class="option" data-choice="auth" onclick="toggleSelect(this)">
    <div class="letter">A</div>
    <div class="content">
      <h3>用户认证</h3>
    </div>
  </div>
  <div class="option" data-choice="dashboard" onclick="toggleSelect(this)">
    <div class="letter">B</div>
    <div class="content">
      <h3>数据仪表板</h3>
    </div>
  </div>
</div>
EOF
```

### 架构决策

```bash
cat > ~/.visual-choice/session/screens/architecture.html << 'EOF'
<h2>技术选型评估</h2>
<div class="pros-cons">
  <div class="pros">
    <h4>✓ 优点</h4>
    <ul>
      <li>独立部署，快速迭代</li>
    </ul>
  </div>
  <div class="cons">
    <h4>✗ 缺点</h4>
    <ul>
      <li>系统复杂度增加</li>
    </ul>
  </div>
</div>
EOF
```

## ✅ 验证清单

部署完成后，验证以下内容：

- [ ] 所有平台 skill 目录存在
- [ ] SKILL.md 文件存在且格式正确
- [ ] Binary 文件可执行
- [ ] 脚本有执行权限
- [ ] 可以启动服务器
- [ ] 可以写入 HTML 内容
- [ ] 可以查看事件记录

```bash
# 运行自动验证
./verify.sh
```

## 🐛 故障排查

### Skill 未被发现

1. 检查 SKILL.md frontmatter 格式
2. 验证 name 字段：`visual-choice`
3. 确认目录名称匹配 name
4. 重启 AI Agent

### Binary 无法执行

```bash
# 检查架构
file bin/visual-choice

# 重新编译
cd src && go build -o ../bin/visual-choice
```

### 端口被占用

```bash
# 停止旧服务器
./scripts/stop.sh

# 或使用不同端口
bin/visual-choice start --port 5235
```

## 📝 维护

### 更新 Binary

```bash
cd src
go build -o ../bin/visual-choice
../verify.sh
```

### 更新文档

直接修改对应的 `.md` 文件，所有平台自动同步（符号链接模式）。

### 添加新功能

1. 修改 `src/` 目录下的源代码
2. 重新编译：`go build -o ../bin/visual-choice`
3. 验证：`../verify.sh`
4. 更新文档

## 📄 许可证

- 源代码：遵循原始项目许可证
- 文档：MIT
- Binary：包含 fsnotify 依赖（Apache 2.0）

## 🔗 参考资料

- [Claude Code Skills](https://code.claude.com/docs/en/skills.md)
- [OpenCode Skills](https://open-code.ai/en/docs/skills)
- [Go Cross Compilation](https://go.dev/doc/install/source#environment)

## 📞 支持

遇到问题或有改进建议？

1. 查看 [MAINTENANCE.md](MAINTENANCE.md) 故障排查部分
2. 查看 [CROSS-PLATFORM.md](CROSS-PLATFORM.md) 跨平台部署指南
3. 运行 `./verify.sh` 自动验证

---

**快速测试**

```bash
# 1. 部署
./deploy-cross-platform.sh symlink

# 2. 启动
./scripts/start.sh

# 3. 在 AI Agent 中输入
/visual-choice
```

🎉 开始使用吧！

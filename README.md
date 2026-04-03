# Visual Choice - 可视化决策工具

[![Platform](https://img.shields.io/badge/platform-Cursor%20%7C%20Flow%20IDE%20%7C%20Claude%20Code%20%7C%20OpenCode-blue)]()
[![Go Version](https://img.shields.io/badge/go-1.21+-yellow)]()
[![License](https://img.shields.io/badge/license-MIT-green)]()

> **在浏览器中展示可视化选项，记录用户点击反馈，用于设计决策、原型测试、用户调研等场景。**

## 特性亮点

- 浏览器展示设计稿、原型、架构方案
- 自动记录用户点击行为
- 支持 Cursor、Claude Code、OpenCode
- 本地运行，数据完全私有
- Go 编译为单一二进制，无需依赖
- 文件变化实时生效
- 多种模板：options/cards/split/pros-cons

## 安装方式

支持两种安装方式：

### 用户级安装（推荐个人使用）

```bash
# 自动检测平台并安装
./install.sh user

# 或指定平台
./install.sh user cursor    # 安装到 ~/.cursor/skills/
./install.sh user flow      # 安装到 ~/.flow/skills/
./install.sh user claude    # 安装到 ~/.claude/skills/
./install.sh user opencode  # 安装到 ~/.config/opencode/skills/
```

Session 数据：`~/.visual-choice/session/`（所有项目共享）

### 工程级安装（推荐团队项目）

```bash
# 在项目根目录执行
./install.sh project

# 或指定平台
./install.sh project cursor
./install.sh project flow
./install.sh project claude
./install.sh project opencode
```

Session 数据：`项目/.visual-choice/session/`（项目独立，可纳入 Git）

### 查看帮助

```bash
./install.sh --help
```

## 快速开始

### 1. 启动服务器

```bash
./scripts/start.sh
```

### 2. 写入 HTML 内容

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

### 3. 在浏览器查看

打开 http://localhost:5234，点击选项进行选择。

### 4. 查看反馈结果

```bash
./scripts/events.sh
```

## 使用场景

| 场景 | 说明 |
|------|------|
| 设计评审 | Logo 方案、配色方案、界面风格确认 |
| 原型演示 | 产品原型展示、用户测试反馈、A/B 测试 |
| 用户调研 | 一对一访谈、用户偏好收集、可用性测试 |
| 产品规划 | 功能需求投票、用户故事优先级、Sprint 规划 |
| 技术决策 | 技术选型评估、架构方案评审、方案优缺点对比 |
| 团队投票 | 团队决策、意见收集、工具选型投票 |

## 技术架构

### 技术栈

| 组件 | 技术 |
|------|------|
| 语言 | Go 1.21+ |
| HTTP | net/http 标准库 |
| 文件监听 | fsnotify |
| 数据格式 | JSON/JSONL |
| 前端 | 原生 HTML/CSS/JS |

### 项目结构

```
visual-choice/
├── install.sh           # 安装脚本
├── SKILL.md             # Skill 主文档
├── scripts/             # 辅助脚本
│   ├── start.sh
│   ├── stop.sh
│   ├── status.sh
│   └── events.sh
├── bin/                 # 编译产物
└── src/                 # 源代码
```

### 安全特性

- 路径遍历保护
- HTTP 超时配置 (15s/60s)
- 请求体限制 (最大 10MB)
- 文件权限收紧 (0600)
- 端口号验证 (1-65535)

### 性能指标

| 指标 | 数值 |
|------|------|
| Binary 大小 | ~6 MB |
| 启动时间 | < 100ms |
| HTTP 响应 | < 10ms |
| 内存占用 | ~10 MB |

## 开发指南

### 构建项目

```bash
make build          # 标准构建
make build-fast     # 快速构建（无优化）
make build-cross    # 交叉编译多平台
```

### 交叉编译

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o bin/visual-choice-linux

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o bin/visual-choice-darwin-arm64

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o bin/visual-choice-darwin-amd64

# Windows
GOOS=windows GOARCH=amd64 go build -o bin/visual-choice-windows.exe
```

### 运行测试

```bash
make test           # 运行测试
make test-coverage  # 生成覆盖率报告
```

## 故障排查

### Binary 无法执行

```bash
# 检查架构
file bin/visual-choice

# 重新编译
cd src && make build
```

### 端口被占用

```bash
./scripts/stop.sh

# 或使用不同端口
./bin/visual-choice start --port 5235
```

### PID 文件残留

```bash
rm -f ~/.visual-choice/session/server.pid
```

### Skill 未被发现

1. 检查 SKILL.md frontmatter 格式
2. 验证 name 字段：`visual-choice`
3. 重启 AI Agent

## 文档

| 文档 | 说明 |
|------|------|
| [SKILL.md](SKILL.md) | Skill 主文档 |
| [examples.md](examples.md) | 使用示例 |
| [reference.md](reference.md) | 技术参考 |
| [src/README.md](src/README.md) | 源代码说明 |

## 许可证

- 源代码：MIT License
- 依赖：fsnotify (Apache 2.0)

## 参考资料

- [Go 官方文档](https://golang.org/doc/)
- [fsnotify GitHub](https://github.com/fsnotify/fsnotify)
- [Claude Code Skills](https://code.claude.com/docs/en/skills.md)
- [OpenCode Skills](https://open-code.ai/en/docs/skills)
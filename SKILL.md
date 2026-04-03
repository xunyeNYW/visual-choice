---
name: visual-choice
version: 1.1.0
description: 浏览器可视化选择工具 - 展示 mockup、设计选项、架构方案供用户点击选择。使用场景：设计评审、原型演示、用户调研、产品优先级投票、架构决策可视化。跨平台支持：Cursor、Flow IDE、Claude Code、OpenCode。支持用户级和工程级两种安装方式。
---

# Visual Choice - 可视化选择工具

> **跨平台支持**: Cursor · Flow IDE · Claude Code · OpenCode
>
> 在浏览器中展示可视化选项，记录用户点击反馈，用于设计决策、原型测试、用户调研等场景。

## 安装方式

支持两种安装方式，根据使用场景选择：

### 用户级安装（推荐个人使用）

所有项目共享同一套工具和数据：

```bash
# 自动检测平台并安装
./install.sh user

# 或指定平台
./install.sh user cursor    # 安装到 ~/.cursor/skills/
./install.sh user flow      # 安装到 ~/.flow/skills/
./install.sh user claude    # 安装到 ~/.claude/skills/
./install.sh user opencode  # 安装到 ~/.config/opencode/skills/
```

安装位置：
- Cursor: `~/.cursor/skills/visual-choice/`
- Flow IDE: `~/.flow/skills/visual-choice/`
- Claude Code: `~/.claude/skills/visual-choice/`
- OpenCode: `~/.config/opencode/skills/visual-choice/`
- Session 数据: `~/.visual-choice/session/`（所有项目共享）

### 工程级安装（推荐团队项目）

每个项目独立工具和数据，可纳入版本控制：

```bash
# 在项目根目录执行
./install.sh project

# 或指定平台
./install.sh project cursor    # 安装到 项目/.cursor/skills/
./install.sh project flow      # 安装到 项目/.flow/skills/
./install.sh project claude    # 安装到 项目/.claude/skills/
./install.sh project opencode  # 安装到 项目/.config/opencode/skills/
```

安装位置：
- Cursor: `项目/.cursor/skills/visual-choice/`
- Flow IDE: `项目/.flow/skills/visual-choice/`
- Claude Code: `项目/.claude/skills/visual-choice/`
- OpenCode: `项目/.config/opencode/skills/visual-choice/`
- Session 数据: `项目/.visual-choice/session/`（项目独立）

### 环境变量覆盖

临时指定 Session 目录：

```bash
VISUAL_CHOICE_SESSION=/tmp/test-session ./scripts/start.sh
```

### 查看帮助

```bash
./install.sh --help
```

---

## 快速开始

当需要用户进行视觉化决策时（设计选择、原型确认、架构评审），使用此工具在浏览器中展示选项并记录点击反馈。

### 1. 启动服务器

```bash
# 使用启动脚本（推荐）
~/.cursor/skills/visual-choice/scripts/start.sh

# 或手动启动
~/.cursor/skills/visual-choice/bin/visual-choice start --port 5234 --dir ~/.visual-choice/session
```

输出：
```
服务器已启动
URL: http://localhost:5234
Screen 目录：~/.visual-choice/session/screens
State 目录：~/.visual-choice/session/state
```

### 2. 写入 HTML 内容

将 HTML 文件写入 `~/.visual-choice/session/screens/` 目录：

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

### 3. 用户交互

1. 告诉用户："请在浏览器打开 http://localhost:5234 查看选项"
2. 用户在浏览器点击选择
3. 点击事件自动记录到 `~/.visual-choice/session/state/events.jsonl`

### 4. 查看结果

```bash
# 使用查看脚本
~/.cursor/skills/visual-choice/scripts/events.sh

# 或手动查看
cat ~/.visual-choice/session/state/events.jsonl
```

输出示例：
```
事件记录:
---------
[1] 17:28:56 click - 选择：minimal - 极简风格
[2] 17:29:03 click - 选择：bold - 大胆风格
[3] 17:29:10 click - 选择：minimal - 极简风格
```

---

## 核心工作流程

### 设计决策评审

```bash
# 1. 启动服务器
~/.cursor/skills/visual-choice/scripts/start.sh

# 2. 推送设计选项
cat > ~/.visual-choice/session/screens/logo.html << 'EOF'
<h2>哪个 Logo 方案更合适？</h2>
<p class="subtitle">考虑品牌识别度和可扩展性</p>

<div class="cards">
  <div class="card" data-choice="a" onclick="toggleSelect(this)">
    <div class="card-image"><!-- 放入 Logo A 图片 --></div>
    <div class="card-body">
      <h3>方案 A - 极简</h3>
      <p>简洁几何图形</p>
    </div>
  </div>
  <div class="card" data-choice="b" onclick="toggleSelect(this)">
    <div class="card-image"><!-- 放入 Logo B 图片 --></div>
    <div class="card-body">
      <h3>方案 B - 文字</h3>
      <p>品牌名称设计</p>
    </div>
  </div>
</div>
EOF

# 3. 收集团队投票
~/.cursor/skills/visual-choice/scripts/events.sh
```

### 产品功能优先级

```bash
cat > ~/.visual-choice/session/screens/features.html << 'EOF'
<h2>哪些功能应该优先开发？</h2>
<p class="subtitle">可选择多个（多选）</p>

<div class="options" data-multiselect>
  <div class="option" data-choice="auth" onclick="toggleSelect(this)">
    <div class="letter">A</div>
    <div class="content">
      <h3>用户认证</h3>
      <p>登录注册，权限管理</p>
    </div>
  </div>
  <div class="option" data-choice="dashboard" onclick="toggleSelect(this)">
    <div class="letter">B</div>
    <div class="content">
      <h3>数据仪表板</h3>
      <p>可视化数据展示</p>
    </div>
  </div>
  <div class="option" data-choice="report" onclick="toggleSelect(this)">
    <div class="letter">C</div>
    <div class="content">
      <h3>报表导出</h3>
      <p>Excel/PDF 导出</p>
    </div>
  </div>
</div>
EOF
```

### 架构方案对比

```bash
cat > ~/.visual-choice/session/screens/architecture.html << 'EOF'
<h2>技术选型评估</h2>
<p class="subtitle">是否使用微服务架构？</p>

<div class="pros-cons">
  <div class="pros">
    <h4>✓ 优点</h4>
    <ul>
      <li>独立部署，快速迭代</li>
      <li>技术栈灵活，按需选择</li>
      <li>故障隔离，提升稳定性</li>
    </ul>
  </div>
  <div class="cons">
    <h4>✗ 缺点</h4>
    <ul>
      <li>系统复杂度增加</li>
      <li>运维成本上升</li>
      <li>数据一致性问题</li>
    </ul>
  </div>
</div>

<div class="options">
  <div class="option" data-choice="yes" onclick="toggleSelect(this)">
    <div class="letter">A</div>
    <div class="content">
      <h3>采用微服务</h3>
    </div>
  </div>
  <div class="option" data-choice="no" onclick="toggleSelect(this)">
    <div class="letter">B</div>
    <div class="content">
      <h3>保持单体</h3>
    </div>
  </div>
</div>
EOF
```

---

## 可用模板

### 选项卡片（A/B/C 选择）

```html
<div class="options">
  <div class="option" data-choice="a" onclick="toggleSelect(this)">
    <div class="letter">A</div>
    <div class="content">
      <h3>标题</h3>
      <p>描述</p>
    </div>
  </div>
</div>
```

**多选模式：**
```html
<div class="options" data-multiselect>
  <!-- 用户可以多选 -->
</div>
```

### 卡片样式（带图片）

```html
<div class="cards">
  <div class="card" data-choice="design1" onclick="toggleSelect(this)">
    <div class="card-image"><!-- 图片/视频/渲染内容 --></div>
    <div class="card-body">
      <h3>设计名称</h3>
      <p>描述</p>
    </div>
  </div>
</div>
```

### 并排对比

```html
<div class="split">
  <div class="mockup"><!-- 左侧内容 --></div>
  <div class="mockup"><!-- 右侧内容 --></div>
</div>
```

### 优缺点对比

```html
<div class="pros-cons">
  <div class="pros">
    <h4>✓ 优点</h4>
    <ul>
      <li>优势 1</li>
      <li>优势 2</li>
    </ul>
  </div>
  <div class="cons">
    <h4>✗ 缺点</h4>
    <ul>
      <li>劣势 1</li>
      <li>劣势 2</li>
    </ul>
  </div>
</div>
```

---

## 脚本工具

| 脚本 | 功能 |
|------|------|
| `install.sh` | 安装工具（用户级或工程级） |
| `scripts/start.sh` | 启动服务器 |
| `scripts/stop.sh` | 停止服务器 |
| `scripts/events.sh` | 查看事件记录 |
| `scripts/status.sh` | 查看服务器状态 |

---

## 技术原理

- **文件监听**：自动检测 `screens/` 目录的 HTML 文件变化
- **模板包装**：用户只需写内容片段，自动包装完整 HTML
- **事件记录**：用户点击记录到 JSONL 文件（每行一个 JSON 对象）
- **进程管理**：PID 文件管理服务器生命周期

详细技术文档见 [reference.md](reference.md)
完整使用示例见 [examples.md](examples.md)

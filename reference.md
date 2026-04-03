# Visual Choice 技术参考

## 项目位置

Skill 目录：
```
~/.cursor/skills/visual-choice/
```

## 命令行接口

```bash
visual-choice <command> [选项]

命令:
  start     启动服务器
  status    查看服务器状态
  stop      停止服务器
  events    查看事件记录

启动选项:
  --port    端口号 (默认：5234)
  --dir     会话目录 (默认：./session)
```

## 完整示例

### 启动服务器

```bash
# 使用 skill 中的 binary
~/.cursor/skills/visual-choice/bin/visual-choice start --port 5234 --dir ~/.visual-choice/session

# 或使用启动脚本
~/.cursor/skills/visual-choice/scripts/start.sh
```

### 写入内容

```bash
# 内容片段模式（推荐）
cat > ~/.visual-choice/session/screens/platform.html << 'EOF'
<h2>哪个平台更重要？</h2>
<p class="subtitle">选择你的首要目标平台</p>

<div class="options">
  <div class="option" data-choice="web" onclick="toggleSelect(this)">
    <div class="letter">A</div>
    <div class="content">
      <h3>Web 应用</h3>
      <p>跨平台，易访问</p>
    </div>
  </div>
  <div class="option" data-choice="mobile" onclick="toggleSelect(this)">
    <div class="letter">B</div>
    <div class="content">
      <h3>移动应用</h3>
      <p>iOS + Android</p>
    </div>
  </div>
</div>
EOF

# 完整 HTML 文档模式
cat > ~/.visual-choice/session/screens/custom.html << 'EOF'
<!DOCTYPE html>
<html>
<head>
  <title>自定义页面</title>
  <style>
    /* 自定义样式 */
  </style>
</head>
<body>
  <h1>完整自定义内容</h1>
  <!-- 服务器会自动注入交互脚本 -->
</body>
</html>
EOF
```

### 查看状态

```bash
./visual-choice status --dir ~/.visual-choice/session

# 输出:
# 服务器正在运行
# PID: 12345
# URL: http://localhost:5234
# 端口：5234
```

### 停止服务器

```bash
./visual-choice stop --dir ~/.visual-choice/session
```

### 查看事件

```bash
./visual-choice events --dir ~/.visual-choice/session

# 输出:
# 事件记录:
# ---------
# [1] 17:28:56 click - 选择：web - Web 应用
# [2] 17:29:03 click - 选择：mobile - 移动应用
```

## 目录结构

```
session/
├── screens/           # HTML 文件目录
│   ├── platform.html
│   └── layout.html
├── state/             # 状态文件目录
│   ├── events.jsonl   # 用户交互事件
│   └── server-info.json  # 服务器信息
└── server.pid         # 服务器进程 ID
```

## 事件格式

```json
{
  "type": "click",
  "choice": "web",
  "text": "Web 应用",
  "timestamp": 1706000101,
  "server_time": "2024-01-23T10:15:01Z"
}
```

## CSS 类参考

### 选项卡片

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

### 卡片样式

```html
<div class="cards">
  <div class="card" data-choice="design1" onclick="toggleSelect(this)">
    <div class="card-image"><!-- 内容 --></div>
    <div class="card-body">
      <h3>名称</h3>
      <p>描述</p>
    </div>
  </div>
</div>
```

### Mockup 容器

```html
<div class="mockup">
  <div class="mockup-header">预览标题</div>
  <div class="mockup-body"><!-- 内容 --></div>
</div>
```

### 分割视图

```html
<div class="split">
  <div class="mockup"><!-- 左侧 --></div>
  <div class="mockup"><!-- 右侧 --></div>
</div>
```

### 优缺点对比

```html
<div class="pros-cons">
  <div class="pros">
    <h4>✓ 优点</h4>
    <ul>
      <li>优势</li>
    </ul>
  </div>
  <div class="cons">
    <h4>✗ 缺点</h4>
    <ul>
      <li>劣势</li>
    </ul>
  </div>
</div>
```

### Mock 元素

```html
<div class="mock-nav">Logo | 首页 | 关于</div>
<div class="mock-sidebar">导航</div>
<div class="mock-content">主要内容</div>
<button class="mock-button">按钮</button>
<input class="mock-input" placeholder="输入框">
<div class="placeholder">占位区域</div>
```

## 技术栈

- **Go 1.21+**：编译为单一二进制
- **net/http**：HTTP 服务器
- **fsnotify**：文件监听
- **sync.Mutex**：并发控制
- **JSONL**：事件存储格式

## 核心代码位置

| 文件 | 功能 | 行数 |
|------|------|------|
| `main.go` | CLI 入口 + 命令处理 | 390 |
| `server.go` | HTTP 服务器 + 文件监听 + HTML 模板 | 947 |
| `events.go` | 事件记录工具 | 108 |

## 高级用法

### 多会话管理

```bash
# 会话 1 - 设计评审
./visual-choice start --port 5234 --dir ~/.visual-choice/design-review

# 会话 2 - 用户调研
./visual-choice start --port 5235 --dir ~/.visual-choice/user-research

# 查看不同会话的事件
./visual-choice events --dir ~/.visual-choice/design-review
./visual-choice events --dir ~/.visual-choice/user-research
```

### 持久化配置

会话数据存储在 `~/.visual-choice/` 目录，重启服务器后仍然保留。

### API 集成

```bash
# 获取最新屏幕信息
curl http://localhost:5234/latest

# 输出:
# {"file":".../screens/platform.html","filename":"platform.html","timestamp":"..."}
```

### 事件数据分析

```bash
# 导出为 CSV
cat ~/.visual-choice/session/state/events.jsonl | \
  jq -r '[.timestamp, .choice, .text] | @csv' > events.csv

# 统计选择次数
cat ~/.visual-choice/session/state/events.jsonl | \
  jq -r '.choice' | sort | uniq -c | sort -rn
```

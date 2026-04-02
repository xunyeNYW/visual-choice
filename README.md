# Visual Choice - 可视化决策工具 🎨

[![Platform](https://img.shields.io/badge/platform-Cursor%20%7C%20Claude%20Code%20%7C%20OpenCode-blue)]()
[![Go Version](https://img.shields.io/badge/go-1.21+-yellow)]()
[![License](https://img.shields.io/badge/license-MIT-green)]()

> **在浏览器中展示可视化选项，记录用户点击反馈，用于设计决策、原型测试、用户调研等场景。**

## ✨ 特性亮点

- 🌐 **浏览器展示**: 在本地浏览器中展示设计稿、原型、架构方案
- 👆 **点击反馈**: 自动记录用户点击行为，生成反馈报告
- 📊 **决策支持**: 适用于设计评审、原型测试、用户调研、产品优先级投票
- 🔧 **跨平台**: 支持 Cursor、Claude Code、OpenCode 等 AI 助手
- 🛡️ **安全可靠**: 本地运行，数据完全私有
- 📦 **单一二进制**: Go 编译，无需安装依赖
- ⚡ **实时监听**: 文件变化即时生效，无需重启服务
- 🎨 **多种模板**: options/cards/split/pros-cons 等丰富布局

## 🚀 快速开始

### 1. 启动服务器

```bash
# 使用启动脚本
./scripts/start.sh

# 或手动启动
./bin/visual-choice start --port 5234
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

# 输出示例:
# 事件记录:
# ---------
# [1] 17:28:56 click - 选择：minimal - 极简风格
# [2] 17:29:03 click - 选择：bold - 大胆风格
```

## 💡 使用场景

### 🎨 设计评审
- Logo 方案选择
- 配色方案投票
- 界面风格确认
- 设计系统评审

### 📱 原型演示
- 产品原型展示
- 用户测试反馈
- 交互流程验证
- A/B 测试原型

### 👥 用户调研
- 一对一访谈
- 用户偏好收集
- 可用性测试
- 功能优先级访谈

### 📊 产品规划
- 功能需求投票
- 用户故事优先级
- 产品路线图确认
- Sprint 规划投票

### 🏗️ 技术决策
- 技术选型评估
- 架构方案评审
- 方案优缺点对比
- 微服务 vs 单体评估

### 🗳️ 团队投票
- 团队决策
- 意见收集
- 技术栈选择
- 工具选型投票

## 🏗️ 技术架构

### 技术栈

| 组件 | 技术 | 说明 |
|------|------|------|
| **语言** | Go 1.21+ | 编译为单一二进制 |
| **HTTP** | net/http | 标准库 |
| **文件监听** | fsnotify | 跨平台文件事件 |
| **并发** | Goroutine + sync.Mutex | 轻量级并发 |
| **数据格式** | JSON/JSONL | 事件存储 |
| **前端** | 原生 HTML/CSS/JS | 零依赖 |

### 项目结构

```
visual-choice/
├── cmd/visual-choice/     # CLI 入口
├── internal/              # 内部包
│   ├── server/           # HTTP 服务器
│   ├── events/           # 事件处理
│   └── models/           # 数据模型
├── scripts/              # 辅助脚本
├── bin/                  # 编译产物
├── Makefile              # 构建脚本
└── src/                  # 源代码
```

### 安全特性

- ✅ **路径遍历保护**: 验证文件路径在允许目录内
- ✅ **HTTP 超时配置**: 防止资源耗尽攻击 (15s/60s)
- ✅ **请求体限制**: 最大 10MB 请求体
- ✅ **文件权限收紧**: 敏感文件使用 0600 权限
- ✅ **事件字段验证**: 类型检查和长度限制
- ✅ **端口号验证**: 限制在 1-65535 范围
- ✅ **Goroutine 正确退出**: 避免资源泄漏

## 📊 性能指标

| 指标 | 数值 | 测试条件 |
|------|------|---------|
| **Binary 大小** | 6.1 MB | 优化编译后 |
| **启动时间** | < 100ms | 冷启动 |
| **HTTP 响应** | < 10ms | 本地请求 |
| **文件检测** | < 50ms | 文件写入到检测 |
| **内存占用** | ~10 MB | 运行中 |
| **测试覆盖率** | 85%+ | 单元测试 |

## 🔨 开发指南

### 构建项目

```bash
# 标准构建（优化）
make build

# 快速构建（无优化）
make build-fast

# 交叉编译多平台
make build-cross
```

### 运行测试

```bash
# 运行所有测试
make test

# 生成覆盖率报告
make test-coverage

# 显示覆盖率
make test-cover
```

### 代码质量

```bash
# 格式化代码
make fmt

# 运行检查
make lint

# 完整 CI 流程
make ci
```

### 运行服务器

```bash
# 开发模式
make run

# 查看状态
make status

# 停止服务器
make stop
```

## 📚 文档导航

| 文档 | 说明 |
|------|------|
| [SKILL.md](SKILL.md) | Skill 主文档，包含快速开始和使用方式 |
| [CROSS-PLATFORM.md](CROSS-PLATFORM.md) | 跨平台部署指南 |
| [examples.md](examples.md) | 5 个完整使用场景示例 |
| [MAINTENANCE.md](MAINTENANCE.md) | 维护指南，故障排查 |
| [src/README.md](src/README.md) | 源代码说明和 API 文档 |
| [src/BUILD.md](src/BUILD.md) | 详细构建指南 |

## 🎯 使用示例

### 示例 1: Logo 方案投票

```bash
cat > ~/.visual-choice/session/screens/logo.html << 'EOF'
<h2>哪个 Logo 方案更合适？</h2>
<p class="subtitle">考虑品牌识别度和可扩展性</p>

<div class="cards">
  <div class="card" data-choice="a" onclick="toggleSelect(this)">
    <div class="card-image" style="background: #667eea;">
      <span style="color: white; font-size: 3rem;">A</span>
    </div>
    <div class="card-body">
      <h3>方案 A - 极简</h3>
      <p>简洁几何图形</p>
    </div>
  </div>
  <div class="card" data-choice="b" onclick="toggleSelect(this)">
    <div class="card-image" style="background: #f093fb;">
      <span style="color: white; font-size: 3rem;">B</span>
    </div>
    <div class="card-body">
      <h3>方案 B - 文字</h3>
      <p>品牌名称设计</p>
    </div>
  </div>
</div>
EOF
```

### 示例 2: 功能优先级（多选）

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

### 示例 3: 架构方案对比

```bash
cat > ~/.visual-choice/session/screens/architecture.html << 'EOF'
<h2>技术选型评估</h2>
<p class="subtitle">是否采用微服务架构？</p>

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

## ✅ 验证清单

部署完成后，验证以下内容:

```bash
# 运行自动验证
./verify.sh
```

- [ ] ✅ 所有平台 skill 目录存在
- [ ] ✅ SKILL.md 文件存在且格式正确
- [ ] ✅ Binary 文件可执行
- [ ] ✅ 脚本有执行权限
- [ ] ✅ 可以启动服务器
- [ ] ✅ 可以写入 HTML 内容
- [ ] ✅ 可以查看事件记录

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
cd src && make build
```

### 端口被占用

```bash
# 停止旧服务器
./scripts/stop.sh

# 或使用不同端口
./bin/visual-choice start --port 5235
```

## 📝 维护

### 更新 Binary

```bash
cd src
make build
../verify.sh
```

### 添加新功能

1. 修改 `src/` 目录下的源代码
2. 重新编译：`make build`
3. 运行测试：`make test`
4. 验证：`../verify.sh`
5. 更新文档

## 📄 许可证

- **源代码**: MIT License
- **文档**: MIT License
- **依赖**: fsnotify (Apache 2.0)

## 🔗 参考资料

- [Go 官方文档](https://golang.org/doc/)
- [fsnotify GitHub](https://github.com/fsnotify/fsnotify)
- [Claude Code Skills](https://code.claude.com/docs/en/skills.md)
- [OpenCode Skills](https://open-code.ai/en/docs/skills)

## 📞 支持

遇到问题或有改进建议？

1. 查看 [MAINTENANCE.md](MAINTENANCE.md) 故障排查部分
2. 查看 [CROSS-PLATFORM.md](CROSS-PLATFORM.md) 跨平台部署指南
3. 运行 `./verify.sh` 自动验证
4. 提交 Issue 或 Pull Request

---

## 🎉 开始使用

```bash
# 1. 启动服务器
./scripts/start.sh

# 2. 在 AI Agent 中输入
/visual-choice
```

**享受可视化决策的乐趣！** 🚀

# Visual Choice 使用示例

## 示例 1: 设计评审会议

### 场景
需要团队对 3 个 Logo 设计方案进行投票。

### 步骤

```bash
# 1. 启动服务器
~/.cursor/skills/visual-choice/scripts/start.sh

# 2. 准备设计方案
cat > ~/.visual-choice/session/screens/logo.html << 'EOF'
<h2>哪个 Logo 方案更合适？</h2>
<p class="subtitle">考虑品牌识别度、可扩展性、记忆点</p>

<div class="cards">
  <div class="card" data-choice="minimal" onclick="toggleSelect(this)">
    <div class="card-image" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); display: flex; align-items: center; justify-content: center;">
      <svg width="100" height="100" viewBox="0 0 100 100">
        <circle cx="50" cy="50" r="40" fill="white" opacity="0.9"/>
        <rect x="35" y="35" width="30" height="30" fill="#667eea"/>
      </svg>
    </div>
    <div class="card-body">
      <h3>方案 A - 极简几何</h3>
      <p>圆形 + 方形组合</p>
    </div>
  </div>
  <div class="card" data-choice="text" onclick="toggleSelect(this)">
    <div class="card-image" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); display: flex; align-items: center; justify-content: center;">
      <span style="color: white; font-size: 1.5rem; font-weight: bold;">BrandName</span>
    </div>
    <div class="card-body">
      <h3>方案 B - 文字设计</h3>
      <p>品牌名称艺术字</p>
    </div>
  </div>
  <div class="card" data-choice="icon" onclick="toggleSelect(this)">
    <div class="card-image" style="background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); display: flex; align-items: center; justify-content: center;">
      <svg width="80" height="80" viewBox="0 0 80 80">
        <path d="M40 10 L70 70 L10 70 Z" fill="white" opacity="0.9"/>
      </svg>
    </div>
    <div class="card-body">
      <h3>方案 C - 图标符号</h3>
      <p>抽象三角形</p>
    </div>
  </div>
</div>
EOF

# 3. 在会议中展示
echo "请大家在浏览器打开 http://localhost:5234 进行投票"

# 4. 实时查看结果
~/.cursor/skills/visual-choice/scripts/events.sh

# 5. 导出结果
cat ~/.visual-choice/session/state/events.jsonl | jq -r '.choice' | sort | uniq -c
```

---

## 示例 2: 用户调研访谈

### 场景
进行一对一用户访谈，了解用户对功能优先级的看法。

### 步骤

```bash
# 1. 启动专用会话
./visual-choice start --port 5234 --dir ~/.visual-choice/interview-001

# 2. 准备功能优先级问题
cat > ~/.visual-choice/interview-001/screens/features.html << 'EOF'
<h2>哪些功能对你最重要？</h2>
<p class="subtitle">请选择最多 3 个</p>

<div class="options" data-multiselect>
  <div class="option" data-choice="auth" onclick="toggleSelect(this)">
    <div class="letter">A</div>
    <div class="content">
      <h3>单点登录 (SSO)</h3>
      <p>企业账号统一登录</p>
    </div>
  </div>
  <div class="option" data-choice="report" onclick="toggleSelect(this)">
    <div class="letter">B</div>
    <div class="content">
      <h3>数据报表</h3>
      <p>可视化数据分析</p>
    </div>
  </div>
  <div class="option" data-choice="api" onclick="toggleSelect(this)">
    <div class="letter">C</div>
    <div class="content">
      <h3>开放 API</h3>
      <p>与现有系统集成</p>
    </div>
  </div>
  <div class="option" data-choice="mobile" onclick="toggleSelect(this)">
    <div class="letter">D</div>
    <div class="content">
      <h3>移动应用</h3>
      <p>iOS/Android App</p>
    </div>
  </div>
  <div class="option" data-choice="export" onclick="toggleSelect(this)">
    <div class="letter">E</div>
    <div class="content">
      <h3>数据导出</h3>
      <p>Excel/CSV 导出</p>
    </div>
  </div>
</div>
EOF

# 3. 访谈过程中观察用户选择
# 4. 访谈后保存记录
cp ~/.visual-choice/interview-001/state/events.jsonl ./research/interview-001-results.jsonl

# 5. 清理会话
./visual-choice stop --dir ~/.visual-choice/interview-001
```

---

## 示例 3: Sprint 规划会议

### 场景
团队对 Sprint 要开发的故事进行优先级排序和估算。

### 步骤

```bash
# 1. 启动会议会话
~/.cursor/skills/visual-choice/scripts/start.sh

# 2. 故事优先级投票
cat > ~/.visual-choice/session/screens/story-priority.html << 'EOF'
<h2>Sprint 故事优先级</h2>
<p class="subtitle">选择本 Sprint 必须完成的故事（多选）</p>

<div class="options" data-multiselect>
  <div class="option" data-choice="story-1" onclick="toggleSelect(this)">
    <div class="letter">1</div>
    <div class="content">
      <h3>用户登录流程优化</h3>
      <p>简化注册步骤，支持社交登录</p>
    </div>
  </div>
  <div class="option" data-choice="story-2" onclick="toggleSelect(this)">
    <div class="letter">2</div>
    <div class="content">
      <h3>数据仪表板改版</h3>
      <p>新增图表类型，支持自定义布局</p>
    </div>
  </div>
  <div class="option" data-choice="story-3" onclick="toggleSelect(this)">
    <div class="letter">3</div>
    <div class="content">
      <h3>性能优化</h3>
      <p>页面加载时间 < 2 秒</p>
    </div>
  </div>
</div>
EOF

# 3. 故事点估算
cat > ~/.visual-choice/session/screens/story-points.html << 'EOF'
<h2>故事点估算</h2>
<p class="subtitle">"用户登录流程优化"需要多少工作量？</p>

<div class="options">
  <div class="option" data-choice="1" onclick="toggleSelect(this)">
    <div class="letter">1</div>
    <div class="content">
      <h3>1 点</h3>
      <p>很小，几小时完成</p>
    </div>
  </div>
  <div class="option" data-choice="2" onclick="toggleSelect(this)">
    <div class="letter">2</div>
    <div class="content">
      <h3>2 点</h3>
      <p>小，1 天内完成</p>
    </div>
  </div>
  <div class="option" data-choice="3" onclick="toggleSelect(this)">
    <div class="letter">3</div>
    <div class="content">
      <h3>3 点</h3>
      <p>中等，2-3 天</p>
    </div>
  </div>
  <div class="option" data-choice="5" onclick="toggleSelect(this)">
    <div class="letter">5</div>
    <div class="content">
      <h3>5 点</h3>
      <p>较大，1 周+</p>
    </div>
  </div>
</div>
EOF
```

---

## 示例 4: 架构决策记录

### 场景
记录技术选型决策过程，供后续参考。

### 步骤

```bash
# 1. 启动会话
./visual-choice start --port 5234 --dir ~/.visual-choice/architecture-decision

# 2. 数据库选型对比
cat > ~/.visual-choice/architecture-decision/screens/database.html << 'EOF'
<h2>数据库选型</h2>
<p class="subtitle">新项目应该使用哪种数据库？</p>

<div class="pros-cons">
  <div class="pros">
    <h4>✓ PostgreSQL</h4>
    <ul>
      <li>ACID 事务保证</li>
      <li>复杂查询能力强</li>
      <li>JSON 支持良好</li>
      <li>扩展性优秀</li>
    </ul>
  </div>
  <div class="cons">
    <h4>✗ PostgreSQL</h4>
    <ul>
      <li>运维复杂度较高</li>
      <li>水平扩展成本高</li>
    </ul>
  </div>
</div>

<div class="options">
  <div class="option" data-choice="postgres" onclick="toggleSelect(this)">
    <div class="letter">A</div>
    <div class="content">
      <h3>PostgreSQL</h3>
      <p>关系型数据库首选</p>
    </div>
  </div>
  <div class="option" data-choice="mongodb" onclick="toggleSelect(this)">
    <div class="letter">B</div>
    <div class="content">
      <h3>MongoDB</h3>
      <p>文档型数据库</p>
    </div>
  </div>
</div>
EOF

# 3. 记录决策结果
~/.cursor/skills/visual-choice/scripts/events.sh > architecture-decision-$(date +%Y%m%d).md

# 4. 将结果添加到项目文档
echo "## 数据库选型决策 ($(date +%Y-%m-%d))" >> docs/architecture-decisions.md
cat architecture-decision-$(date +%Y%m%d).md >> docs/architecture-decisions.md
```

---

## 示例 5: 产品原型测试

### 场景
向潜在用户展示产品原型，收集反馈。

### 步骤

```bash
# 1. 准备原型展示
cat > ~/.visual-choice/session/screens/prototype.html << 'EOF'
<h2>产品原型预览</h2>
<p class="subtitle">这是我们的新产品概念</p>

<div class="mockup">
  <div class="mockup-header">Dashboard - 数据概览</div>
  <div class="mockup-body">
    <div class="mock-nav">
      <span style="font-weight: bold;">📊 Analytics</span>
      <span>首页</span>
      <span>数据</span>
      <span>报表</span>
      <span>设置</span>
    </div>
    <div style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 1rem; margin-top: 1rem;">
      <div style="background: #667eea; color: white; padding: 2rem; border-radius: 0.5rem; text-align: center;">
        <div style="font-size: 2rem; font-weight: bold;">1,234</div>
        <div style="font-size: 0.875rem; opacity: 0.9;">总用户</div>
      </div>
      <div style="background: #f093fb; color: white; padding: 2rem; border-radius: 0.5rem; text-align: center;">
        <div style="font-size: 2rem; font-weight: bold;">567</div>
        <div style="font-size: 0.875rem; opacity: 0.9;">活跃用户</div>
      </div>
      <div style="background: #4facfe; color: white; padding: 2rem; border-radius: 0.5rem; text-align: center;">
        <div style="font-size: 2rem; font-weight: bold;">89%</div>
        <div style="font-size: 0.875rem; opacity: 0.9;">留存率</div>
      </div>
    </div>
  </div>
</div>

<div class="section" style="margin-top: 2rem;">
  <h3>第一印象反馈</h3>
  <p class="subtitle">这个设计给你的感觉是？</p>
  <div class="options">
    <div class="option" data-choice="professional" onclick="toggleSelect(this)">
      <div class="letter">A</div>
      <div class="content">
        <h3>专业可靠</h3>
        <p>适合企业场景</p>
      </div>
    </div>
    <div class="option" data-choice="modern" onclick="toggleSelect(this)">
      <div class="letter">B</div>
      <div class="content">
        <h3>现代时尚</h3>
        <p>符合当下审美</p>
      </div>
    </div>
    <div class="option" data-choice="complex" onclick="toggleSelect(this)">
      <div class="letter">C</div>
      <div class="content">
        <h3>过于复杂</h3>
        <p>信息密度太高</p>
      </div>
    </div>
  </div>
</div>
EOF

# 2. 在用户测试会议中展示
echo "请查看原型并点击选择你的第一印象"

# 3. 收集反馈
~/.cursor/skills/visual-choice/scripts/events.sh
```

---

## 提示与技巧

### 1. 使用真实内容

```html
<!-- 好：使用真实图片 -->
<div class="card-image">
  <img src="https://images.unsplash.com/photo-xxx" alt="设计预览" style="width: 100%; height: 100%; object-fit: cover;">
</div>

<!-- 差：使用占位符 -->
<div class="card-image">图片</div>
```

### 2. 控制选项数量

- **最佳**: 2-4 个选项
- **可接受**: 5-6 个选项
- **避免**: 超过 7 个选项（认知负荷过高）

### 3. 清晰的决策问题

```html
<!-- 好：问题具体 -->
<h2>哪个布局更适合移动端？</h2>
<p class="subtitle">考虑拇指操作区域和可读性</p>

<!-- 差：问题模糊 -->
<h2>你喜欢哪个？</h2>
```

### 4. 迭代式测试

```bash
# 版本 1 - 测试整体方向
cat > ~/.visual-choice/session/screens/design-v1.html

# 根据反馈调整
cat > ~/.visual-choice/session/screens/design-v2.html

# 最终确认
cat > ~/.visual-choice/session/screens/design-final.html
```

### 5. 结果导出

```bash
# 导出为 JSON
cat ~/.visual-choice/session/state/events.jsonl | jq '.' > results.json

# 导出为 CSV
cat ~/.visual-choice/session/state/events.jsonl | \
  jq -r '[.server_time, .choice, .text] | @csv' > results.csv

# 统计分析
cat ~/.visual-choice/session/state/events.jsonl | \
  jq -r '.choice' | sort | uniq -c | sort -rn
```

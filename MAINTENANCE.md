# Visual Choice Skill 维护指南

## 目录结构

```
~/.cursor/skills/visual-choice/
├── SKILL.md           # 主文档
├── reference.md       # 技术参考
├── examples.md        # 使用示例
├── MAINTENANCE.md     # 维护指南
├── scripts/           # 辅助脚本
│   ├── start.sh      # 启动服务器
│   ├── stop.sh       # 停止服务器
│   ├── status.sh     # 查看状态
│   └── events.sh     # 查看事件
├── bin/              # Binary 文件
│   └── visual-choice # Go 编译的二进制
├── src/              # 源代码目录（资产）
│   ├── main.go       # CLI 入口
│   ├── server.go     # HTTP 服务器
│   ├── events.go     # 事件处理
│   ├── go.mod        # 依赖定义
│   ├── README.md     # 源码说明
│   └── test.sh       # 测试脚本
└── verify.sh         # 验证脚本
```

## 更新 Binary

当源代码更新后，需要重新编译并复制 binary：

### 方法 1: 从源代码编译（推荐）

```bash
# 1. 进入源代码目录
cd ~/.cursor/skills/visual-choice/src

# 2. 编译
go build -o ../bin/visual-choice

# 3. 验证
cd ..
./verify.sh
```

### 方法 2: 从外部源代码更新

```bash
# 1. 从源代码仓库复制
cp /path/to/visual-choice/*.go ~/.cursor/skills/visual-choice/src/
cp /path/to/visual-choice/go.* ~/.cursor/skills/visual-choice/src/

# 2. 重新编译
cd ~/.cursor/skills/visual-choice/src
go build -o ../bin/visual-choice

# 3. 验证
cd ..
./verify.sh
```

## 验证安装

```bash
# 运行验证脚本
~/.cursor/skills/visual-choice/verify.sh

# 预期输出:
# ✅ Skill 验证通过！
```

## 测试流程

```bash
# 1. 启动服务器
~/.cursor/skills/visual-choice/scripts/start.sh

# 2. 查看状态
~/.cursor/skills/visual-choice/scripts/status.sh

# 3. 写入测试内容
cat > ~/.visual-choice/session/screens/test.html << 'EOF'
<h2>测试页面</h2>
<div class="options">
  <div class="option" data-choice="a" onclick="toggleSelect(this)">
    <div class="letter">A</div>
    <div class="content">
      <h3>选项 A</h3>
    </div>
  </div>
</div>
EOF

# 4. 在浏览器访问 http://localhost:5234

# 5. 查看事件
~/.cursor/skills/visual-choice/scripts/events.sh

# 6. 停止服务器
~/.cursor/skills/visual-choice/scripts/stop.sh
```

## 故障排查

### Binary 文件不存在

```bash
# 错误：找不到 visual-choice 二进制文件
# 解决：编译源代码
cd /media/vdc/10226591/github/superpowers/visual-choice
go build -o visual-choice
cp visual-choice ~/.cursor/skills/visual-choice/bin/
```

### 端口被占用

```bash
# 错误：address already in use
# 解决 1: 停止旧服务器
~/.cursor/skills/visual-choice/scripts/stop.sh

# 解决 2: 使用不同端口
~/.cursor/skills/visual-choice/bin/visual-choice start --port 5235 --dir ~/.visual-choice/session
```

### PID 文件残留

```bash
# 错误：服务器已在运行（但实际没有）
# 解决：清理 PID 文件
rm -f ~/.visual-choice/session/server.pid
```

## 版本管理

如果需要管理多个版本：

```bash
# 创建版本目录
~/.cursor/skills/visual-choice/bin/
├── visual-choice        # 当前版本（symlink 或 copy）
├── visual-choice-v1.0   # 版本 1.0
└── visual-choice-v1.1   # 版本 1.1

# 切换版本
ln -sf visual-choice-v1.1 visual-choice
```

## 跨平台注意事项

Binary 文件是平台相关的：

- **macOS (Apple Silicon)**: `GOOS=darwin GOARCH=arm64 go build`
- **macOS (Intel)**: `GOOS=darwin GOARCH=amd64 go build`
- **Linux (x86_64)**: `GOOS=linux GOARCH=amd64 go build`
- **Windows**: `GOOS=windows GOARCH=amd64 go build`

当前安装的 binary 是为当前平台编译的。

# Visual Choice 源代码

本目录包含 visual-choice 的完整源代码，用于参考、修改和重新编译。

## 文件结构

```
src/
├── main.go           # CLI 入口 + 命令处理 (390 行)
├── server.go         # HTTP 服务器 + 文件监听 + HTML 模板 (947 行)
├── events.go         # 事件记录工具 (108 行)
├── go.mod            # Go 模块定义
├── go.sum            # 依赖校验
├── test.sh          # 自动化测试脚本
└── example.html     # 测试页面示例
```

## 核心组件

### main.go (390 行)

**功能**：CLI 命令入口和进程管理

**主要函数**：
- `main()` - 命令分发 (start/status/stop/events)
- `handleStart()` - 启动服务器，创建目录结构，写入 PID
- `handleStatus()` - 检查服务器运行状态
- `handleStop()` - 停止服务器（SIGTERM → SIGKILL）
- `handleEvents()` - 查看事件记录

**关键数据结构**：
```go
type ServerInfo struct {
    Type      string `json:"type"`
    Port      int    `json:"port"`
    URL       string `json:"url"`
    ScreenDir string `json:"screen_dir"`
    StateDir  string `json:"state_dir"`
}
```

### server.go (947 行)

**功能**：HTTP 服务器、文件监听、HTML 模板

**主要组件**：

1. **Server 结构**
```go
type Server struct {
    port       int
    screenDir  string
    stateDir   string
    httpServer *http.Server
    watcher    *fsnotify.Watcher
    wg         sync.WaitGroup
    mu         sync.RWMutex
    latestFile string
}
```

2. **HTML 框架模板** (FrameTemplate, 488 行)
   - 完整 CSS 样式系统
   - 响应式布局
   - 选项卡片样式
   - 点击动画
   - 选择指示器
   - 交互脚本

3. **HTTP 端点**
   - `GET /` - 提供最新 HTML
   - `POST /event` - 记录用户点击
   - `GET /latest` - 获取屏幕信息

4. **文件监听** (`watchFiles()`)
   - 使用 fsnotify 监听 `screens/` 目录
   - 自动更新 `latestFile` 指针
   - 新文件时清空事件

### events.go (108 行)

**功能**：事件记录工具函数

**主要函数**：
- `appendEvent()` - 追加事件到 JSONL 文件
- `clearEvents()` - 清空事件文件

## 编译说明

### 本地编译

```bash
cd ~/.cursor/skills/visual-choice/src
go build -o ../bin/visual-choice
```

### 跨平台编译

```bash
# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o ../bin/visual-choice

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o ../bin/visual-choice

# Linux (x86_64)
GOOS=linux GOARCH=amd64 go build -o ../bin/visual-choice

# Windows
GOOS=windows GOARCH=amd64 go build -o ../bin/visual-choice.exe
```

### 优化编译

```bash
# 减小二进制大小
go build -ldflags="-s -w" -o ../bin/visual-choice
```

## 运行测试

```bash
cd ~/.cursor/skills/visual-choice/src
./test.sh
```

测试内容：
1. 服务器启动/停止
2. HTTP 内容服务
3. 文件监听
4. 事件记录 API
5. 命令行工具

## 修改流程

### 1. 修改源代码

```bash
cd ~/.cursor/skills/visual-choice/src
# 编辑 main.go, server.go, events.go
```

### 2. 重新编译

```bash
go build -o ../bin/visual-choice
```

### 3. 验证

```bash
cd ~/.cursor/skills/visual-choice
./verify.sh
```

### 4. 测试

```bash
# 启动服务器
./scripts/start.sh

# 手动测试功能
# ...

# 停止服务器
./scripts/stop.sh
```

## 技术栈

| 组件 | 技术 | 说明 |
|------|------|------|
| 语言 | Go 1.21+ | 编译为单一二进制 |
| HTTP | net/http | 标准库 |
| 文件监听 | fsnotify | 跨平台文件事件 |
| 并发 | sync.Mutex + Goroutine | 轻量级并发 |
| 数据格式 | JSON/JSONL | 事件存储 |

## 核心代码片段

### 文件监听循环

```go
func (s *Server) watchFiles() {
    for {
        select {
        case event, ok := <-s.watcher.Events:
            if !ok { return }
            
            if !strings.HasSuffix(event.Name, ".html") {
                continue
            }
            
            if event.Op&fsnotify.Create == fsnotify.Create ||
               event.Op&fsnotify.Write == fsnotify.Write {
                s.mu.Lock()
                s.latestFile = event.Name
                s.mu.Unlock()
                
                if event.Op&fsnotify.Create == fsnotify.Create {
                    s.clearEvents()
                }
            }
        }
    }
}
```

### HTML 模板包装

```go
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
    latestFile, err := s.getLatestHTML()
    if err != nil {
        s.serveWelcome(w)
        return
    }
    
    content, err := os.ReadFile(latestFile)
    contentStr := string(content)
    
    isFullDoc := strings.HasPrefix(strings.TrimSpace(contentStr), "<!DOCTYPE") ||
                 strings.HasPrefix(strings.TrimSpace(contentStr), "<html")
    
    var html string
    if isFullDoc {
        html = s.injectScript(contentStr)
    } else {
        html = strings.Replace(FrameTemplate, "{{CONTENT}}", contentStr, 1)
    }
    
    w.Write([]byte(html))
}
```

### 事件记录

```go
func (s *Server) appendEvent(event map[string]interface{}) error {
    data, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    f, err := os.OpenFile(eventsFile, 
        os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    defer f.Close()
    
    _, err = f.Write(append(data, '\n'))
    return err
}
```

## 依赖项

### 外部依赖

```go
github.com/fsnotify/fsnotify  // 文件监听
```

### 标准库

```go
encoding/json    // JSON 编解码
flag            // 命令行参数解析
fmt             // 格式化输出
net/http        // HTTP 服务器
os              // 操作系统接口
path/filepath   // 路径处理
sort            // 排序
strings         // 字符串处理
sync            // 同步原语
syscall         // 系统调用
time            // 时间处理
```

## 性能特征

| 指标 | 数值 |
|------|------|
| 二进制大小 | ~8.4 MB |
| 启动时间 | < 100ms |
| HTTP 响应 | < 10ms |
| 文件检测 | < 50ms |
| 内存占用 | ~10 MB |

## 扩展方向

### 添加新端点

在 `server.go` 的 `NewServer()` 中添加：

```go
mux.HandleFunc("/api/custom", s.handleCustom)
```

### 添加新命令

在 `main.go` 的 `main()` 中添加：

```go
case "custom":
    handleCustom(os.Args[2:])
```

### 修改 HTML 模板

编辑 `server.go` 中的 `FrameTemplate` 常量。

## 参考资料

- [Go 官方文档](https://golang.org/doc/)
- [fsnotify GitHub](https://github.com/fsnotify/fsnotify)
- [net/http 包文档](https://pkg.go.dev/net/http)

## 版本历史

| 版本 | 日期 | 说明 |
|------|------|------|
| 1.0 | 2026-04-01 | 初始版本，核心功能完成 |

## 维护清单

### 定期检查

- [ ] 更新 Go 版本
- [ ] 检查依赖安全更新
- [ ] 测试跨平台编译
- [ ] 验证二进制文件大小

### 升级流程

1. 拉取源代码仓库最新代码
2. 复制源代码到 `src/` 目录
3. 重新编译所有平台版本
4. 运行完整测试
5. 更新版本号

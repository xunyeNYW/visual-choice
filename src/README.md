# Visual Choice 源代码

本目录包含 visual-choice 的完整源代码，用于参考、修改和重新编译。

## 快速开始

### 构建项目
```bash
make build
```

### 运行测试
```bash
make test
```

### 运行服务器
```bash
make run
```

## 文件结构

```
src/
├── cmd/
│   └── visual-choice/
│       └── main.go              # CLI 入口 + 命令处理
├── internal/
│   ├── events/
│   │   ├── events.go            # 事件存储和处理
│   │   └── events_test.go       # 事件测试 (85.2% 覆盖率)
│   ├── models/
│   │   ├── models.go            # 数据模型定义
│   │   └── models_test.go       # 模型测试 (100% 覆盖率)
│   └── server/
│       ├── server.go            # HTTP 服务器 + 文件监听
│       └── template.go          # HTML 模板定义
├── Makefile                     # 构建脚本
├── go.mod                       # Go 模块定义
├── go.sum                       # 依赖校验
├── test.sh                      # 自动化测试脚本
└── example.html                 # 测试页面示例
```

## 核心组件

### cmd/visual-choice/main.go

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

### internal/server/server.go

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

2. **HTML 框架模板** (FrameTemplate)
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

### internal/events/events.go

**功能**：事件存储和处理

**主要函数**：
- `Store.Append()` - 追加事件到 JSONL 文件
- `Store.ReadEvents()` - 读取事件文件
- `Store.Clear()` - 清空事件文件
- `FormatEvents()` - 格式化显示事件
- `GetLastChoice()` - 获取最后一次选择
- `GetAllChoices()` - 获取所有选择（去重）

### internal/models/models.go

**功能**：数据模型定义

**主要结构**：
- `ServerInfo` - 服务器启动信息
- `Event` - 用户交互事件
- `ServerConfig` - 服务器配置

## 编译说明

### 使用 Makefile（推荐）

```bash
# 标准构建（优化）
make build

# 快速构建（无优化）
make build-fast

# 交叉编译多平台
make build-cross

# 完整 CI 流程
make ci
```

### 手动编译

```bash
cd ~/.cursor/skills/visual-choice/src
go build -o ../bin/visual-choice ./cmd/visual-choice
```

### 跨平台编译

```bash
# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o ../bin/visual-choice ./cmd/visual-choice

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o ../bin/visual-choice ./cmd/visual-choice

# Linux (x86_64)
GOOS=linux GOARCH=amd64 go build -o ../bin/visual-choice ./cmd/visual-choice

# Windows
GOOS=windows GOARCH=amd64 go build -o ../bin/visual-choice ./cmd/visual-choice.exe
```

### 优化编译

```bash
# 减小二进制大小
go build -ldflags="-s -w" -o ../bin/visual-choice ./cmd/visual-choice
```

## 运行测试

### 使用 Makefile

```bash
# 运行所有测试
make test

# 生成覆盖率报告
make test-coverage

# 显示覆盖率
make test-cover
```

### 手动运行

```bash
cd ~/.cursor/skills/visual-choice/src
go test -v ./...
```

测试覆盖率：
- `internal/events`: 85.2%
- `internal/models`: 100.0%

### 测试内容

1. 事件追加和读取
2. 事件清空
3. 事件格式化
4. 获取最后选择
5. 获取所有选择（去重）
6. 空文件处理
7. 无效 JSON 处理
8. 默认配置验证

## 修改流程

### 1. 修改源代码

```bash
cd ~/.cursor/skills/visual-choice/src
# 编辑对应包的文件
```

### 2. 重新编译

```bash
# 使用 Makefile
make build

# 或手动编译
go build -o ../bin/visual-choice ./cmd/visual-choice
```

### 3. 运行测试

```bash
# 使用 Makefile
make test

# 或手动运行
go test -v ./...
```

### 4. 验证

```bash
cd ~/.cursor/skills/visual-choice
./verify.sh
```

### 5. 测试功能

```bash
# 启动服务器
make run

# 或手动启动
./../bin/visual-choice start --port 5234 --dir ./session

# 手动测试功能
# ...

# 停止服务器
make stop
```

## 技术栈

| 组件 | 技术 | 说明 |
|------|------|------|
| 语言 | Go 1.21+ | 编译为单一二进制 |
| HTTP | net/http | 标准库 |
| 文件监听 | fsnotify | 跨平台文件事件 |
| 并发 | sync.Mutex + Goroutine | 轻量级并发 |
| 数据格式 | JSON/JSONL | 事件存储 |
| 测试 | testing | 标准库单元测试 |

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
    
    // 路径遍历保护
    cleanPath := filepath.Clean(latestFile)
    cleanScreenDir := filepath.Clean(s.screenDir)
    if !strings.HasPrefix(cleanPath, cleanScreenDir) {
        http.Error(w, "非法的文件路径", http.StatusForbidden)
        return
    }
    
    content, err := os.ReadFile(cleanPath)
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

### 事件记录（带安全加固）

```go
func (s *Server) appendEvent(event map[string]interface{}) error {
    data, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    // 使用 0600 权限
    f, err := os.OpenFile(eventsFile, 
        os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
    defer f.Close()
    
    _, err = f.Write(append(data, '\n'))
    return err
}
```

## 依赖项

### 外部依赖

```go
github.com/fsnotify/fsnotify v1.7.0  // 文件监听
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
context         // 上下文管理
```

## 性能特征

| 指标 | 数值 |
|------|------|
| 二进制大小 | ~8.4 MB |
| 启动时间 | < 100ms |
| HTTP 响应 | < 10ms |
| 文件检测 | < 50ms |
| 内存占用 | ~10 MB |

## 安全特性

### 已实现的安全加固

1. **路径遍历保护** - 验证文件路径在允许目录内
2. **HTTP 超时配置** - 防止资源耗尽攻击
3. **请求体大小限制** - 最大 10MB
4. **文件权限收紧** - 敏感文件使用 0600 权限
5. **事件字段验证** - 类型检查和长度限制
6. **端口号验证** - 限制在 1-65535 范围
7. **Goroutine 正确退出** - 避免资源泄漏

## 扩展方向

### 添加新端点

在 `internal/server/server.go` 的 `Start()` 中添加：

```go
mux.HandleFunc("/api/custom", s.handleCustom)
```

### 添加新命令

在 `cmd/visual-choice/main.go` 的 `main()` 中添加：

```go
case "custom":
    handleCustom(os.Args[2:])
```

### 添加新测试

在对应包的 `_test.go` 文件中添加测试函数：

```go
func TestYourFeature(t *testing.T) {
    // 测试代码
}
```

## Makefile 命令

### 构建命令
- `make build` - 标准构建（优化）
- `make build-fast` - 快速构建
- `make build-cross` - 交叉编译

### 测试命令
- `make test` - 运行所有测试
- `make test-coverage` - 生成覆盖率报告
- `make test-cover` - 显示覆盖率

### 代码质量
- `make fmt` - 格式化代码
- `make fmt-check` - 检查代码格式
- `make vet` - 运行 go vet
- `make lint` - 完整检查

### 依赖管理
- `make tidy` - 整理依赖
- `make deps-update` - 更新依赖

### 清理命令
- `make clean` - 清理构建产物
- `make distclean` - 深度清理

### 运行命令
- `make run` - 运行服务器
- `make status` - 查看状态
- `make stop` - 停止服务器
- `make events` - 查看事件

详见 [BUILD.md](BUILD.md)

## 参考资料

- [Go 官方文档](https://golang.org/doc/)
- [fsnotify GitHub](https://github.com/fsnotify/fsnotify)
- [net/http 包文档](https://pkg.go.dev/net/http)
- [Go 测试指南](https://golang.org/doc/tutorial/add-a-test)
- [Go 项目布局标准](https://github.com/golang-standards/project-layout)

## 版本历史

| 版本 | 日期 | 说明 |
|------|------|------|
| 1.0 | 2026-04-01 | 初始版本，核心功能完成 |
| 1.1 | 2026-04-02 | 重构为标准 Go 项目结构，添加 Makefile 和测试 |

## 维护清单

### 定期检查

- [ ] 更新 Go 版本
- [ ] 检查依赖安全更新
- [ ] 测试跨平台编译
- [ ] 验证二进制文件大小
- [ ] 运行完整测试套件
- [ ] 检查测试覆盖率

### 升级流程

1. 拉取源代码仓库最新代码
2. 复制源代码到 `src/` 目录
3. 运行 `make tidy` 整理依赖
4. 运行 `make ci` 执行完整 CI 流程
5. 重新编译所有平台版本
6. 更新版本号

# Visual Choice 构建指南

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

## Makefile 命令详解

### 构建命令

#### `make build` - 标准构建
构建优化后的生产版本，包含版本信息。
```bash
make build
```
输出：`../bin/visual-choice`

#### `make build-fast` - 快速构建
快速构建开发版本，无优化。
```bash
make build-fast
```

#### `make build-cross` - 交叉编译
编译多平台版本（Linux、macOS、Windows）。
```bash
make build-cross
```
输出：
- `../bin/visual-choice-linux-amd64`
- `../bin/visual-choice-darwin-amd64`
- `../bin/visual-choice-darwin-arm64`
- `../bin/visual-choice-windows-amd64.exe`

### 测试命令

#### `make test` - 运行所有测试
```bash
make test
```

#### `make test-coverage` - 生成覆盖率报告
生成 HTML 格式的覆盖率报告。
```bash
make test-coverage
```
输出：`coverage.html`（可用浏览器打开查看）

#### `make test-cover` - 显示覆盖率
在终端显示测试覆盖率。
```bash
make test-cover
```

#### `make test-unit` - 运行单元测试
只运行以 `Test` 开头的单元测试。
```bash
make test-unit
```

### 代码质量

#### `make fmt` - 格式化代码
自动格式化所有 Go 文件。
```bash
make fmt
```

#### `make fmt-check` - 检查代码格式
检查代码格式，不自动修复。
```bash
make fmt-check
```

#### `make vet` - 运行 go vet
检查代码中的可疑构造。
```bash
make vet
```

#### `make lint` - 完整代码检查
运行所有代码质量检查（fmt-check + vet）。
```bash
make lint
```

### 依赖管理

#### `make tidy` - 整理依赖
更新 `go.mod` 和 `go.sum`。
```bash
make tidy
```

#### `make deps-update` - 更新依赖
更新所有依赖到最新版本。
```bash
make deps-update
```

#### `make deps-clean` - 清理依赖缓存
清理 Go 模块缓存。
```bash
make deps-clean
```

### 清理命令

#### `make clean` - 清理构建产物
删除编译生成的二进制文件和测试报告。
```bash
make clean
```

#### `make distclean` - 深度清理
清理构建产物和依赖缓存。
```bash
make distclean
```

### 运行命令

#### `make run` - 运行服务器
构建并启动开发服务器。
```bash
make run
```
服务器将在 `http://localhost:5234` 启动。

#### `make status` - 查看服务器状态
```bash
make status
```

#### `make stop` - 停止服务器
```bash
make stop
```

#### `make events` - 查看事件记录
```bash
make events
```

### 安装命令

#### `make install` - 安装到 GOPATH
```bash
make install
```
安装位置：`$(GOPATH)/bin/visual-choice`

#### `make uninstall` - 从 GOPATH 卸载
```bash
make uninstall
```

### CI/CD 集成

#### `make ci` - 完整 CI 流程
运行完整的持续集成流程（格式化检查 + vet + 测试 + 构建）。
```bash
make ci
```

#### `make release` - 发布构建
运行完整检查并编译所有平台版本。
```bash
make release
```

### 其他命令

#### `make help` - 显示帮助
```bash
make help
```

#### `make version` - 显示版本信息
```bash
make version
```
输出示例：
```
visual-choice v1.0.0
Go 版本：1.21
构建平台：linux/amd64
```

## 开发工作流

### 日常开发
```bash
# 1. 拉取代码后首先整理依赖
make tidy

# 2. 运行测试确保代码正常
make test

# 3. 快速构建并测试
make build-fast
./../bin/visual-choice start --port 5234 --dir ./session

# 4. 提交前运行完整检查
make lint
```

### 发布流程
```bash
# 1. 运行完整 CI 流程
make ci

# 2. 生成覆盖率报告（可选）
make test-coverage

# 3. 编译多平台版本
make release

# 4. 清理构建产物
make clean
```

## 测试覆盖率

当前测试覆盖率：
- `internal/events`: 85.2%
- `internal/models`: 100.0%

查看覆盖率报告：
```bash
make test-coverage
open coverage.html  # macOS
xdg-open coverage.html  # Linux
start coverage.html  # Windows
```

## 故障排查

### 构建失败
```bash
# 清理后重新构建
make clean
make build
```

### 依赖问题
```bash
# 清理依赖缓存
make deps-clean

# 重新整理依赖
make tidy
```

### 测试失败
```bash
# 查看详细测试输出
make test

# 查看覆盖率
make test-cover
```

### 代码格式问题
```bash
# 自动格式化
make fmt

# 检查格式
make fmt-check
```

## 性能优化

### 快速迭代
使用 `build-fast` 进行快速开发迭代：
```bash
make build-fast
```

### 并行测试
Go 默认并行运行测试，无需额外配置。

### 增量构建
Go 会自动缓存编译结果，第二次构建会更快。

## 环境变量

可以通过环境变量自定义构建：

```bash
# 自定义版本号
VERSION=2.0.0 make build

# 自定义 Go 版本
GO_VERSION=1.22 make build

# 自定义输出路径
BIN_DIR=/usr/local/bin make build
```

## 参考文档

- [Go 项目布局标准](https://github.com/golang-standards/project-layout)
- [Go Makefile 最佳实践](https://www.ardan-labs.com/blog/2020/02/modules-06-makefiles/)
- [Go 测试覆盖率](https://go.dev/blog/cover)

# Visual Choice Makefile
# Go 项目标准构建脚本

# 项目信息
PROJECT_NAME := visual-choice
VERSION := 1.0.0
GO_VERSION := 1.21

# 目录配置
SRC_DIR := .
CMD_DIR := $(SRC_DIR)/cmd/visual-choice
BIN_DIR := $(SRC_DIR)/../bin
INTERNAL_DIR := $(SRC_DIR)/internal

# 输出文件
BINARY_NAME := visual-choice
BINARY_PATH := $(BIN_DIR)/$(BINARY_NAME)

# Go 工具
GO := go
GOFMT := gofmt
GOVET := go vet
GOTEST := go test
GOBUILD := go build
GOMOD := go mod

# Go 参数
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
LDFLAGS := -ldflags "-s -w -X main.Version=$(VERSION)"

# 默认目标
.PHONY: all
all: build

# ==================== 构建命令 ====================

# 构建项目
.PHONY: build
build: tidy
	@echo "==> 构建 $(PROJECT_NAME)..."
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH) $(CMD_DIR)
	@echo "✅ 构建完成：$(BINARY_PATH)"

# 快速构建（不优化）
.PHONY: build-fast
build-fast:
	@echo "==> 快速构建..."
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) -o $(BINARY_PATH) $(CMD_DIR)
	@echo "✅ 快速构建完成"

# 交叉编译
.PHONY: build-cross
build-cross:
	@echo "==> 交叉编译..."
	@mkdir -p $(BIN_DIR)
	
	# Linux AMD64
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)
	
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)
	
	# macOS ARM64
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_DIR)
	
	# Windows AMD64
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)
	
	@echo "✅ 交叉编译完成"

# ==================== 测试命令 ====================

# 运行所有测试
.PHONY: test
test:
	@echo "==> 运行测试..."
	$(GOTEST) -v ./...

# 运行测试并生成覆盖率报告
.PHONY: test-coverage
test-coverage:
	@echo "==> 运行测试并生成覆盖率报告..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "✅ 覆盖率报告已生成：coverage.html"
	@echo "📊 查看报告：open coverage.html"

# 运行测试并显示覆盖率
.PHONY: test-cover
test-cover:
	@echo "==> 运行测试并显示覆盖率..."
	$(GOTEST) -cover ./...

# 只运行单元测试
.PHONY: test-unit
test-unit:
	@echo "==> 运行单元测试..."
	$(GOTEST) -run "Test" ./...

# ==================== 代码质量 ====================

# 代码格式化
.PHONY: fmt
fmt:
	@echo "==> 格式化代码..."
	$(GOFMT) -l -w .
	@echo "✅ 格式化完成"

# 检查格式化
.PHONY: fmt-check
fmt-check:
	@echo "==> 检查代码格式..."
	@output=$$($(GOFMT) -l .); \
	if [ -n "$$output" ]; then \
		echo "❌ 以下文件需要格式化:"; \
		echo "$$output"; \
		exit 1; \
	else \
		echo "✅ 代码格式良好"; \
	fi

# 运行 go vet
.PHONY: vet
vet:
	@echo "==> 运行 go vet..."
	$(GOVET) ./...
	@echo "✅ go vet 检查通过"

# 运行所有检查
.PHONY: lint
lint: fmt-check vet
	@echo "✅ 所有代码检查通过"

# ==================== 依赖管理 ====================

# 整理依赖
.PHONY: tidy
tidy:
	@echo "==> 整理依赖..."
	$(GOMOD) tidy
	@echo "✅ 依赖整理完成"

# 更新依赖
.PHONY: deps-update
deps-update:
	@echo "==> 更新依赖..."
	$(GOMOD) get -u ./...
	$(GOMOD) tidy
	@echo "✅ 依赖更新完成"

# 清理依赖缓存
.PHONY: deps-clean
deps-clean:
	@echo "==> 清理依赖缓存..."
	$(GOMOD) clean
	@echo "✅ 依赖缓存已清理"

# ==================== 清理命令 ====================

# 清理构建产物
.PHONY: clean
clean:
	@echo "==> 清理构建产物..."
	rm -rf $(BIN_DIR)
	rm -f $(SRC_DIR)/coverage.out
	rm -f $(SRC_DIR)/coverage.html
	rm -f $(SRC_DIR)/*.test
	@echo "✅ 清理完成"

# 深度清理（包括缓存）
.PHONY: distclean
distclean: clean
	@echo "==> 深度清理..."
	$(GOMOD) clean
	rm -rf $(SRC_DIR)/vendor
	@echo "✅ 深度清理完成"

# ==================== 运行命令 ====================

# 运行服务器（开发模式）
.PHONY: run
run: build-fast
	@echo "==> 启动服务器（开发模式）..."
	$(BINARY_PATH) start --port 5234 --dir ./session

# 查看服务器状态
.PHONY: status
status:
	@echo "==> 查看服务器状态..."
	$(BINARY_PATH) status --dir ./session

# 停止服务器
.PHONY: stop
stop:
	@echo "==> 停止服务器..."
	$(BINARY_PATH) stop --dir ./session

# 查看事件
.PHONY: events
events:
	@echo "==> 查看事件记录..."
	$(BINARY_PATH) events --dir ./session

# ==================== 安装命令 ====================

# 安装到 GOPATH
.PHONY: install
install:
	@echo "==> 安装到 GOPATH..."
	$(GOBUILD) -o $(GOPATH)/bin/$(BINARY_NAME) $(CMD_DIR)
	@echo "✅ 安装完成：$(GOPATH)/bin/$(BINARY_NAME)"

# 从 GOPATH 卸载
.PHONY: uninstall
uninstall:
	@echo "==> 卸载..."
	rm -f $(GOPATH)/bin/$(BINARY_NAME)
	@echo "✅ 卸载完成"

# ==================== 帮助信息 ====================

.PHONY: help
help:
	@echo "$(PROJECT_NAME) v$(VERSION) - Makefile 帮助"
	@echo ""
	@echo "构建命令:"
	@echo "  make build          构建项目（优化）"
	@echo "  make build-fast     快速构建（无优化）"
	@echo "  make build-cross    交叉编译多平台"
	@echo ""
	@echo "测试命令:"
	@echo "  make test           运行所有测试"
	@echo "  make test-coverage  运行测试并生成覆盖率报告"
	@echo "  make test-cover     运行测试并显示覆盖率"
	@echo "  make test-unit      运行单元测试"
	@echo ""
	@echo "代码质量:"
	@echo "  make fmt            格式化代码"
	@echo "  make fmt-check      检查代码格式"
	@echo "  make vet            运行 go vet"
	@echo "  make lint           运行所有检查"
	@echo ""
	@echo "依赖管理:"
	@echo "  make tidy           整理依赖"
	@echo "  make deps-update    更新依赖"
	@echo "  make deps-clean     清理依赖缓存"
	@echo ""
	@echo "清理命令:"
	@echo "  make clean          清理构建产物"
	@echo "  make distclean      深度清理"
	@echo ""
	@echo "运行命令:"
	@echo "  make run            运行服务器（开发模式）"
	@echo "  make status         查看服务器状态"
	@echo "  make stop           停止服务器"
	@echo "  make events         查看事件记录"
	@echo ""
	@echo "安装命令:"
	@echo "  make install        安装到 GOPATH"
	@echo "  make uninstall      从 GOPATH 卸载"
	@echo ""
	@echo "其他:"
	@echo "  make help           显示此帮助信息"
	@echo "  make version        显示版本信息"
	@echo ""

.PHONY: version
version:
	@echo "$(PROJECT_NAME) v$(VERSION)"
	@echo "Go 版本：$(GO_VERSION)"
	@echo "构建平台：$(GOOS)/$(GOARCH)"

# ==================== CI/CD 集成 ====================

# CI 流程
.PHONY: ci
ci: fmt-check vet test build
	@echo "✅ CI 流程完成"

# 发布构建
.PHONY: release
release: lint test build-cross
	@echo "✅ 发布构建完成"
	@echo "📦 构建产物在：$(BIN_DIR)/"
	@ls -lh $(BIN_DIR)/

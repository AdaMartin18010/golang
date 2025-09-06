# Go语言现代化 - 完整测试体系 Makefile

.PHONY: help build run test clean deps lint format

# 默认目标
.DEFAULT_GOAL := help

# 项目信息
PROJECT_NAME := go-modernization-testing-system
VERSION := 1.0.0
BUILD_TIME := $(shell date +%Y-%m-%d_%H:%M:%S)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go相关变量
GO := go
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
CGO_ENABLED := 0

# 构建目录
BUILD_DIR := build
BINARY_NAME := testing-system
MAIN_FILE := main.go

# 测试相关
TEST_DIR := test-results
PERF_DIR := performance-results
PROFILE_DIR := profiles

# 帮助信息
help: ## 显示帮助信息
	@echo "Go语言现代化 - 完整测试体系"
	@echo "================================"
	@echo "可用命令:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# 依赖管理
deps: ## 下载和整理依赖
	@echo "下载依赖..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "依赖管理完成"

# 代码格式化
format: ## 格式化代码
	@echo "格式化代码..."
	$(GO) fmt ./...
	@echo "代码格式化完成"

# 代码检查
lint: ## 运行代码检查
	@echo "运行代码检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint 未安装，跳过代码检查"; \
	fi

# 构建
build: deps ## 构建可执行文件
	@echo "构建 $(PROJECT_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build \
		-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)" \
		-o $(BUILD_DIR)/$(BINARY_NAME) \
		$(MAIN_FILE)
	@echo "构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 运行
run: build ## 运行测试体系演示
	@echo "运行测试体系演示..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# 测试
test: ## 运行单元测试
	@echo "运行单元测试..."
	@mkdir -p $(TEST_DIR)
	$(GO) test -v -coverprofile=$(TEST_DIR)/coverage.out ./...
	$(GO) tool cover -html=$(TEST_DIR)/coverage.out -o $(TEST_DIR)/coverage.html
	@echo "测试完成，覆盖率报告: $(TEST_DIR)/coverage.html"

# 基准测试
bench: ## 运行基准测试
	@echo "运行基准测试..."
	@mkdir -p $(PERF_DIR)
	$(GO) test -bench=. -benchmem -outputdir=$(PERF_DIR) ./...
	@echo "基准测试完成"

# 性能分析
profile: ## 运行性能分析
	@echo "运行性能分析..."
	@mkdir -p $(PROFILE_DIR)
	$(GO) test -cpuprofile=$(PROFILE_DIR)/cpu.prof -memprofile=$(PROFILE_DIR)/mem.prof -bench=. ./...
	@echo "性能分析完成，结果保存在: $(PROFILE_DIR)/"

# 集成测试
integration: ## 运行集成测试
	@echo "运行集成测试..."
	@mkdir -p $(TEST_DIR)
	$(GO) test -tags=integration -v -outputdir=$(TEST_DIR) ./...
	@echo "集成测试完成"

# 完整测试套件
test-all: test bench integration ## 运行完整测试套件
	@echo "完整测试套件执行完成"

# 清理
clean: ## 清理构建文件
	@echo "清理构建文件..."
	rm -rf $(BUILD_DIR)
	rm -rf $(TEST_DIR)
	rm -rf $(PERF_DIR)
	rm -rf $(PROFILE_DIR)
	@echo "清理完成"

# 安装
install: build ## 安装到系统
	@echo "安装 $(PROJECT_NAME)..."
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "安装完成"

# 卸载
uninstall: ## 从系统卸载
	@echo "卸载 $(PROJECT_NAME)..."
	rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "卸载完成"

# 交叉编译
build-all: ## 交叉编译所有平台
	@echo "交叉编译所有平台..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	GOOS=linux GOARCH=arm64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_FILE)
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)
	GOOS=darwin GOARCH=arm64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE)
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)
	@echo "交叉编译完成"

# Docker构建
docker-build: ## 构建Docker镜像
	@echo "构建Docker镜像..."
	docker build -t $(PROJECT_NAME):$(VERSION) .
	docker tag $(PROJECT_NAME):$(VERSION) $(PROJECT_NAME):latest
	@echo "Docker镜像构建完成"

# Docker运行
docker-run: ## 运行Docker容器
	@echo "运行Docker容器..."
	docker run --rm -p 8080:8080 $(PROJECT_NAME):latest

# 开发模式
dev: ## 开发模式运行
	@echo "开发模式运行..."
	$(GO) run $(MAIN_FILE)

# 调试模式
debug: ## 调试模式运行
	@echo "调试模式运行..."
	$(GO) run -race $(MAIN_FILE)

# 生成文档
docs: ## 生成文档
	@echo "生成文档..."
	@mkdir -p docs
	$(GO) doc -all ./... > docs/api.md
	@echo "文档生成完成: docs/api.md"

# 检查代码质量
quality: lint test ## 检查代码质量
	@echo "代码质量检查完成"

# 发布准备
release: clean build-all test-all ## 发布准备
	@echo "发布准备完成"

# 显示版本信息
version: ## 显示版本信息
	@echo "项目: $(PROJECT_NAME)"
	@echo "版本: $(VERSION)"
	@echo "构建时间: $(BUILD_TIME)"
	@echo "Git提交: $(GIT_COMMIT)"
	@echo "Go版本: $(shell go version)"
	@echo "操作系统: $(GOOS)/$(GOARCH)"

# 监控模式
monitor: ## 启动监控模式
	@echo "启动监控模式..."
	$(GO) run $(MAIN_FILE) --monitor

# 性能测试模式
perf: ## 启动性能测试模式
	@echo "启动性能测试模式..."
	$(GO) run $(MAIN_FILE) --performance

# 完整演示
demo: ## 运行完整演示
	@echo "运行完整演示..."
	$(GO) run $(MAIN_FILE) --demo

# 快速开始
quickstart: deps build run ## 快速开始（下载依赖、构建、运行）
	@echo "快速开始完成"

# Go语言AI-Agent架构构建脚本

.PHONY: build run test clean demo help

# 默认目标
all: build

# 构建项目
build:
	@echo "构建AI-Agent架构演示程序..."
	go build -o ai-agent-demo agent.go learning.go coordinator.go router.go specialized_agents.go main.go
	@echo "构建完成: ai-agent-demo"

# 运行演示
run: build
	@echo "运行AI-Agent架构演示..."
	./ai-agent-demo

# 直接运行（不构建）
demo:
	@echo "直接运行演示程序..."
	go run agent.go learning.go coordinator.go router.go specialized_agents.go main.go

# 运行测试
test:
	@echo "运行测试..."
	go test -v ./...

# 运行基准测试
bench:
	@echo "运行基准测试..."
	go test -bench=. -benchmem ./...

# 代码格式化
fmt:
	@echo "格式化代码..."
	go fmt ./...

# 代码检查
lint:
	@echo "运行代码检查..."
	golangci-lint run

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -f ai-agent-demo
	go clean

# 安装依赖
deps:
	@echo "安装依赖..."
	go mod tidy
	go mod download

# 显示帮助
help:
	@echo "可用的命令:"
	@echo "  build   - 构建项目"
	@echo "  run     - 构建并运行演示"
	@echo "  demo    - 直接运行演示（不构建）"
	@echo "  test    - 运行测试"
	@echo "  bench   - 运行基准测试"
	@echo "  fmt     - 格式化代码"
	@echo "  lint    - 运行代码检查"
	@echo "  clean   - 清理构建文件"
	@echo "  deps    - 安装依赖"
	@echo "  help    - 显示此帮助信息"

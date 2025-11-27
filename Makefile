.PHONY: help build run test clean generate deps setup check-env install-tools fmt fmt-check vet mod-tidy mod-download mod-graph mod-why bench bench-cpu bench-mem check pre-commit clean-all

help: ## 显示帮助信息
	@echo "=========================================="
	@echo "  Go 框架开发命令"
	@echo "=========================================="
	@echo ""
	@echo "开发:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep -E '(run|dev|setup|check)' | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""
	@echo "构建:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep -E '(build|generate|install)' | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""
	@echo "测试:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep -E '(test|bench|coverage)' | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""
	@echo "代码质量:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep -E '(lint|fmt|vet|check)' | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""
	@echo "API:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep -E '(openapi|asyncapi|api|validate)' | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""
	@echo "Docker:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep -E '(docker)' | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""
	@echo "其他:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep -vE '(run|dev|setup|check|build|generate|install|test|bench|coverage|lint|fmt|vet|check|openapi|asyncapi|api|validate|docker)' | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""

deps: ## 安装依赖
	go mod download
	go mod tidy

setup: ## 设置开发环境（运行设置脚本）
	@bash scripts/dev/setup.sh

check-env: ## 检查开发环境配置
	@bash scripts/dev/check-env.sh

install-hooks: ## 安装 Git hooks
	@bash scripts/dev/install-hooks.sh

generate: ## 生成代码（Ent, Wire, OpenAPI等）
	@echo "生成 Ent 代码..."
	cd internal/infrastructure/database/ent && go generate ./...
	@echo "生成 Wire 代码..."
	cd scripts/wire && go generate ./...
	@echo "生成 OpenAPI 代码..."
	@bash scripts/api/generate-openapi.sh || echo "OpenAPI 代码生成跳过（需要安装 oapi-codegen）"

generate-ent: ## 生成 Ent 代码
	cd internal/infrastructure/database/ent && go generate ./...

generate-wire: ## 生成 Wire 代码
	cd scripts/wire && go generate ./...

generate-proto: ## 生成 gRPC 代码
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       internal/interfaces/grpc/proto/user.proto

generate-openapi: ## 生成 OpenAPI 代码
	@bash scripts/api/generate-openapi.sh

generate-asyncapi: ## 生成 AsyncAPI 代码
	@bash scripts/api/generate-asyncapi.sh

generate-api-docs: ## 生成 API 文档（HTML）
	@bash scripts/api/generate-docs.sh

validate-openapi: ## 验证 OpenAPI 规范
	@bash scripts/api/validate-openapi.sh

validate-asyncapi: ## 验证 AsyncAPI 规范
	@bash scripts/api/validate-asyncapi.sh

validate-api: validate-openapi validate-asyncapi ## 验证所有 API 规范

generate-ent: ## 生成 Ent 代码
	go generate ./internal/infrastructure/database/ent

build: ## 构建应用
	go build -o bin/server ./cmd/server
	go build -o bin/temporal-worker ./cmd/temporal-worker

run: ## 运行应用
	go run ./cmd/server

run-dev: ## 开发模式运行（热重载）
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air 未安装，使用 go run 代替"; \
		echo "安装 Air: go install github.com/cosmtrek/air@latest"; \
		go run ./cmd/server; \
	fi

run-worker: ## 运行 Temporal Worker
	go run ./cmd/temporal-worker

test: ## 运行测试
	go test -v ./...

test-coverage: ## 运行测试并生成覆盖率报告
	@bash scripts/test-coverage.sh

test-coverage-simple: ## 运行测试并生成覆盖率报告（简单模式）
	go test -v -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

clean: ## 清理构建文件
	rm -rf bin/
	rm -f coverage.out coverage.html

docker-build: ## 构建 Docker 镜像
	docker build -f deployments/docker/Dockerfile -t golang-app:latest .

docker-run: ## 运行 Docker 容器
	docker-compose -f deployments/docker/docker-compose.yml up -d

docker-stop: ## 停止 Docker 容器
	docker-compose -f deployments/docker/docker-compose.yml down

lint: ## 运行代码检查
	golangci-lint run

fmt: ## 格式化代码
	go fmt ./...

fmt-check: ## 检查代码格式（不修改）
	@if [ "$(shell gofmt -s -l . | wc -l)" -gt 0 ]; then \
		echo "❌ 代码格式不符合规范，请运行 'make fmt'"; \
		gofmt -s -d .; \
		exit 1; \
	fi
	@echo "✅ 代码格式检查通过"

vet: ## 运行 go vet
	go vet ./...

mod-tidy: ## 整理 go.mod 依赖
	go mod tidy
	go mod verify

mod-download: ## 下载依赖
	go mod download

mod-graph: ## 显示依赖关系图
	go mod graph

mod-why: ## 解释为什么需要依赖
	@read -p "请输入要查询的包路径: " pkg; \
	go mod why $$pkg

bench: ## 运行基准测试
	go test -bench=. -benchmem ./...

bench-cpu: ## 运行 CPU 性能分析
	go test -bench=. -cpuprofile=cpu.prof ./...
	@echo "CPU 性能分析文件已生成: cpu.prof"
	@echo "查看: go tool pprof cpu.prof"

bench-mem: ## 运行内存性能分析
	go test -bench=. -memprofile=mem.prof ./...
	@echo "内存性能分析文件已生成: mem.prof"
	@echo "查看: go tool pprof mem.prof"

install-tools: ## 安装开发工具
	@echo "安装开发工具..."
	@go install github.com/cosmtrek/air@latest || echo "Air 安装失败"
	@go install github.com/google/wire/cmd/wire@latest || echo "Wire 安装失败"
	@go install entgo.io/ent/cmd/ent@latest || echo "Ent 安装失败"
	@go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest || echo "oapi-codegen 安装失败"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest || echo "golangci-lint 安装失败"
	@echo "✅ 开发工具安装完成"

check: fmt-check vet lint ## 运行所有代码检查

pre-commit: check test ## 提交前检查（格式、检查、测试）

clean-all: clean ## 清理所有生成的文件
	rm -rf tmp/
	rm -rf build/
	rm -f *.prof
	rm -f *.out
	rm -f *.html
	@echo "✅ 清理完成"

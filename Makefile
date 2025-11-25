.PHONY: help build run test clean generate deps

help: ## 显示帮助信息
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

deps: ## 安装依赖
	go mod download
	go mod tidy

generate: ## 生成代码（Ent, Wire等）
	@echo "生成 Ent 代码..."
	cd internal/infrastructure/database/ent && go generate ./...
	@echo "生成 Wire 代码..."
	cd scripts/wire && go generate ./...

generate-ent: ## 生成 Ent 代码
	cd internal/infrastructure/database/ent && go generate ./...

generate-wire: ## 生成 Wire 代码
	cd scripts/wire && go generate ./...

generate-proto: ## 生成 gRPC 代码
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       internal/interfaces/grpc/proto/user.proto

build: ## 构建应用
	go build -o bin/server ./cmd/server

run: ## 运行应用
	go run ./cmd/server

test: ## 运行测试
	go test -v ./...

test-coverage: ## 运行测试并生成覆盖率报告
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

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

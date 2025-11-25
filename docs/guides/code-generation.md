# 代码生成指南

## 概述

本项目使用代码生成工具来减少样板代码，提高开发效率。

## 生成工具

### 1. Ent - ORM 代码生成

Ent 用于生成类型安全的数据库访问代码。

#### 安装 Ent CLI

```bash
go install entgo.io/ent/cmd/ent@latest
```

#### 生成代码

```bash
# 使用 Makefile
make generate-ent

# 或直接运行
cd internal/infrastructure/database/ent
go generate ./...
```

#### 生成的文件

- `ent/gen/` - 生成的 Ent 客户端代码
- `ent/user.go` - User 实体的生成代码

### 2. Wire - 依赖注入代码生成

Wire 用于生成依赖注入代码。

#### 安装 Wire

```bash
go install github.com/google/wire/cmd/wire@latest
```

#### 生成代码

```bash
# 使用 Makefile
make generate-wire

# 或直接运行
cd scripts/wire
go generate ./...
```

#### 生成的文件

- `scripts/wire/wire_gen.go` - Wire 生成的依赖注入代码

### 3. Protocol Buffers - gRPC 代码生成

Protocol Buffers 用于生成 gRPC 服务代码。

#### 安装 protoc

- **macOS**: `brew install protobuf`
- **Linux**: `apt-get install protobuf-compiler`
- **Windows**: 从 [protobuf releases](https://github.com/protocolbuffers/protobuf/releases) 下载

#### 安装 Go 插件

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

#### 生成代码

```bash
# 使用 Makefile
make generate-proto

# 或直接运行
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       internal/interfaces/grpc/proto/user.proto
```

#### 生成的文件

- `internal/interfaces/grpc/proto/user/user.pb.go` - 生成的 protobuf 代码
- `internal/interfaces/grpc/proto/user/user_grpc.pb.go` - 生成的 gRPC 代码

## 一键生成所有代码

### 使用 Makefile

```bash
make generate
```

### 使用脚本

**Linux/macOS:**
```bash
./scripts/generate/all.sh
```

**Windows:**
```powershell
.\scripts\generate\all.ps1
```

## 开发工作流

1. **修改 Schema/Proto 定义**
   - 修改 `internal/infrastructure/database/ent/schema/` 中的 Ent Schema
   - 修改 `internal/interfaces/grpc/proto/` 中的 .proto 文件

2. **生成代码**
   ```bash
   make generate
   ```

3. **使用生成的代码**
   - 在 Repository 中使用生成的 Ent 客户端
   - 在 gRPC handlers 中使用生成的 protobuf 代码

## 注意事项

1. **不要手动编辑生成的文件**
   - 所有生成的文件都应该被 `.gitignore` 忽略
   - 或者提交到版本控制，但不要手动修改

2. **版本控制**
   - 建议将生成的文件提交到版本控制
   - 这样可以确保团队使用相同的生成代码

3. **CI/CD**
   - 在 CI/CD 流程中运行代码生成
   - 验证生成的代码是否与源代码同步

## 故障排除

### Ent 生成失败

- 确保已安装 Ent CLI: `go install entgo.io/ent/cmd/ent@latest`
- 检查 Schema 定义是否正确
- 查看错误信息，修复 Schema 问题

### Wire 生成失败

- 确保已安装 Wire: `go install github.com/google/wire/cmd/wire@latest`
- 检查 `wire.go` 中的 providers 定义
- 确保所有依赖都已正确导入

### protoc 生成失败

- 确保已安装 protoc
- 确保已安装 Go 插件
- 检查 .proto 文件语法
- 确保 PATH 中包含 protoc 和插件


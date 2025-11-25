#!/bin/bash

# 生成所有代码
echo "Generating all code..."

# 生成 Ent 代码
echo "1. Generating Ent code..."
cd internal/infrastructure/database/ent
go generate ./...
cd ../../../../

# 生成 gRPC 代码
echo "2. Generating gRPC code..."
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       internal/interfaces/grpc/proto/user.proto

# 生成 Wire 代码
echo "3. Generating Wire code..."
cd scripts/wire
go generate ./...
cd ../../

echo "All code generation completed!"


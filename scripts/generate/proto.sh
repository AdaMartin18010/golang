#!/bin/bash

# 生成 gRPC 代码
echo "Generating gRPC code..."

protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       internal/interfaces/grpc/proto/user.proto

echo "gRPC code generation completed!"


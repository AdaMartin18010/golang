#!/bin/bash

set -e

PROTO_DIR="internal/interfaces/grpc/proto"
OUTPUT_DIR="internal/interfaces/grpc/proto"

# 确保输出目录存在
mkdir -p "$OUTPUT_DIR"

# 生成 Go 代码
protoc \
  --go_out="$OUTPUT_DIR" \
  --go_opt=paths=source_relative \
  --go-grpc_out="$OUTPUT_DIR" \
  --go-grpc_opt=paths=source_relative \
  -I="$PROTO_DIR" \
  "$PROTO_DIR"/*.proto

echo "gRPC code generation completed"

#!/bin/bash

# 生成 Ent 代码
echo "Generating Ent code..."

cd internal/infrastructure/database/ent
go generate ./...

echo "Ent code generation completed!"


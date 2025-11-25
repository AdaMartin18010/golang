# PowerShell script to generate all code

Write-Host "Generating all code..." -ForegroundColor Green

# 生成 Ent 代码
Write-Host "1. Generating Ent code..." -ForegroundColor Yellow
Set-Location internal/infrastructure/database/ent
go generate ./...
Set-Location ../../../../

# 生成 gRPC 代码
Write-Host "2. Generating gRPC code..." -ForegroundColor Yellow
protoc --go_out=. --go_opt=paths=source_relative `
       --go-grpc_out=. --go-grpc_opt=paths=source_relative `
       internal/interfaces/grpc/proto/user.proto

# 生成 Wire 代码
Write-Host "3. Generating Wire code..." -ForegroundColor Yellow
Set-Location scripts/wire
go generate ./...
Set-Location ../../

Write-Host "All code generation completed!" -ForegroundColor Green


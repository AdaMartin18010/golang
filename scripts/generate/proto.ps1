# PowerShell script to generate gRPC code

Write-Host "Generating gRPC code..." -ForegroundColor Green

protoc --go_out=. --go_opt=paths=source_relative `
       --go-grpc_out=. --go-grpc_opt=paths=source_relative `
       internal/interfaces/grpc/proto/user.proto

Write-Host "gRPC code generation completed!" -ForegroundColor Green


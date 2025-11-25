# PowerShell script to generate Ent code

Write-Host "Generating Ent code..." -ForegroundColor Green

Set-Location internal/infrastructure/database/ent
go generate ./...

Write-Host "Ent code generation completed!" -ForegroundColor Green


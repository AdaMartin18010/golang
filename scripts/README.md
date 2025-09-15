# scripts 使用说明

- 初始化：
  ```bash
  cd scripts && go mod tidy
  ```
- 生成变更日志：
  ```bash
  echo "Added PGO example" | VERSION=v2025.09-P1 go run ./gen_changelog.go
  ```

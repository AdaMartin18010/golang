# 测试指南

## 测试结构

```text
test/
├── unit/              # 单元测试
├── integration/       # 集成测试
├── e2e/              # 端到端测试
└── fixtures/         # 测试数据
```

## 运行测试

### 运行所有测试

```bash
go test ./...
```

### 运行单元测试

```bash
go test ./test/unit/...
```

### 运行集成测试

```bash
go test ./test/integration/...
```

### 运行 E2E 测试

```bash
go test ./test/e2e/...
```

## 测试覆盖率

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 测试最佳实践

1. **单元测试** - 测试单个函数或方法
2. **集成测试** - 测试组件之间的交互
3. **E2E 测试** - 测试完整的用户流程
4. **Mock** - 使用 mock 隔离依赖
5. **Fixtures** - 使用测试数据文件

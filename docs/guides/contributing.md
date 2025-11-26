# 贡献指南

## 开发流程

### 1. 创建功能分支

```bash
git checkout -b feature/your-feature-name
```

### 2. 开发新功能

遵循 Clean Architecture 原则：

- **Domain Layer**: 添加实体、接口、领域服务
- **Application Layer**: 添加用例、DTO
- **Infrastructure Layer**: 实现技术细节
- **Interfaces Layer**: 添加外部接口

### 3. 编写测试

```bash
# 运行测试
go test ./...

# 查看覆盖率
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 4. 代码检查

```bash
# 格式化代码
go fmt ./...

# 运行 linter
golangci-lint run
```

### 5. 提交代码

```bash
git add .
git commit -m "feat: add new feature"
git push origin feature/your-feature-name
```

## 代码规范

### 命名规范

- **包名**: 小写，简短
- **函数名**: 驼峰命名，公开函数首字母大写
- **变量名**: 驼峰命名
- **常量**: 全大写，下划线分隔

### 注释规范

- 所有公开函数、类型、变量都应该有注释
- 注释应该以被注释的内容开头
- 使用 `//` 进行单行注释
- 使用 `/* */` 进行多行注释

### 错误处理

- 总是检查错误
- 使用 `fmt.Errorf` 包装错误
- 使用 `errors.Is` 和 `errors.As` 检查错误

### 测试规范

- 测试文件以 `_test.go` 结尾
- 测试函数以 `Test` 开头
- 使用表驱动测试
- 使用 `testify` 进行断言

## 提交信息规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建/工具相关

示例：

```
feat: add user authentication
fix: resolve database connection issue
docs: update API documentation
```

## Pull Request

### PR 检查清单

- [ ] 代码通过所有测试
- [ ] 代码通过 linter 检查
- [ ] 添加了必要的测试
- [ ] 更新了相关文档
- [ ] 提交信息符合规范

### PR 描述

PR 描述应该包括：

1. **变更说明**: 简要描述变更内容
2. **变更原因**: 为什么需要这个变更
3. **测试方法**: 如何测试这个变更
4. **相关 Issue**: 关联的 Issue 编号

## 代码审查

### 审查要点

1. **架构**: 是否符合 Clean Architecture
2. **代码质量**: 是否遵循代码规范
3. **测试**: 是否有足够的测试覆盖
4. **性能**: 是否有性能问题
5. **安全性**: 是否有安全隐患

## 问题报告

### Bug 报告

请包含以下信息：

1. **环境信息**: Go 版本、操作系统
2. **复现步骤**: 详细的复现步骤
3. **预期行为**: 应该发生什么
4. **实际行为**: 实际发生了什么
5. **错误信息**: 完整的错误堆栈

### 功能请求

请包含以下信息：

1. **功能描述**: 详细描述功能需求
2. **使用场景**: 在什么场景下需要这个功能
3. **可能的实现**: 如果有想法，可以描述实现方案

## 联系方式

- **Issues**: [GitHub Issues](https://github.com/yourusername/golang/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/golang/discussions)

# GitHub Actions 工作流

本目录包含项目的CI/CD自动化工作流配置。

## 📋 工作流列表

### 1. test.yml - 测试和覆盖率

**触发条件**:
- Push到main/develop分支
- Pull Request到main/develop分支

**功能**:
- ✅ 跨平台测试 (Ubuntu, Windows, macOS)
- ✅ 多Go版本测试 (1.23.x, 1.26.x)
- ✅ 竞态检测 (-race)
- ✅ 覆盖率报告生成
- ✅ Codecov上传
- ✅ 构建验证

**运行时间**: ~5-10分钟

---

### 2. lint.yml - 代码质量检查

**触发条件**:
- Push到main/develop分支
- Pull Request到main/develop分支

**功能**:
- ✅ golangci-lint检查
- ✅ gofmt格式检查
- ✅ go vet静态分析

**运行时间**: ~2-3分钟

---

### 3. security.yml - 安全扫描

**触发条件**:
- Push到main分支
- Pull Request到main分支
- 定时执行（每周日）

**功能**:
- ✅ govulncheck漏洞检测
- ✅ gosec安全扫描
- ✅ SARIF报告上传

**运行时间**: ~3-5分钟

---

## 🚀 使用指南

### 本地验证

在提交代码前，可以本地运行检查：

```bash
# 运行测试
go test -v -race -cover ./...

# 代码格式检查
go fmt ./...

# 静态分析
go vet ./...

# Lint检查（需要安装golangci-lint）
golangci-lint run

# 安全扫描（需要安装gosec）
gosec ./...
```

### 使用gox工具

```bash
# 运行测试
gox test

# 质量检查
gox quality

# 覆盖率报告
gox coverage
```

---

## 📊 CI/CD状态徽章

将以下徽章添加到README.md：

```markdown
[![Tests](https://github.com/your-org/your-repo/workflows/Test%20and%20Coverage/badge.svg)](https://github.com/your-org/your-repo/actions)
[![Lint](https://github.com/your-org/your-repo/workflows/Lint/badge.svg)](https://github.com/your-org/your-repo/actions)
[![Security](https://github.com/your-org/your-repo/workflows/Security/badge.svg)](https://github.com/your-org/your-repo/actions)
[![codecov](https://codecov.io/gh/your-org/your-repo/branch/main/graph/badge.svg)](https://codecov.io/gh/your-org/your-repo)
```

---

## 🔧 配置说明

### 环境变量

工作流使用以下环境变量（在GitHub Secrets中配置）：

- `CODECOV_TOKEN`: Codecov上传令牌（可选，公开仓库不需要）

### 缓存策略

工作流使用GitHub Actions缓存来加速构建：
- Go模块缓存
- Go构建缓存

---

## 🐛 故障排查

### 测试失败

1. 检查测试日志
2. 本地运行相同的测试命令
3. 确认依赖版本一致

### Lint失败

1. 运行 `go fmt ./...`
2. 运行 `go vet ./...`
3. 运行 `golangci-lint run --fix`

### 安全扫描警告

1. 查看具体的漏洞报告
2. 更新依赖版本
3. 评估风险并采取措施

---

## 📝 最佳实践

1. **提交前测试**: 总是在本地运行测试
2. **小步提交**: 频繁提交小的更改
3. **描述性消息**: 使用清晰的commit消息
4. **监控CI**: 关注CI运行结果
5. **及时修复**: 快速修复失败的构建

---

## 🔄 更新工作流

修改工作流后：
1. 测试变更是否有效
2. 更新此README
3. 通知团队成员

---

**最后更新**: 2025-10-22  
**维护者**: Project Team


# 贡献指南

感谢你对本项目的兴趣！以下是参与贡献的指南。

---

## 🚀 快速开始

1. **Fork** 本仓库
2. **Clone** 你的 Fork
3. **创建分支**: `git checkout -b feature/your-feature`
4. **提交更改**: `git commit -m "feat: add feature"`
5. **推送分支**: `git push origin feature/your-feature`
6. **创建 Pull Request**

---

## 📋 开发环境

详见 [开发环境搭建指南](docs/development/setup.md)

```bash
# 快速开始
go mod download
docker-compose -f docker-compose.dev.yml up -d
cp .env.example .env
make generate
make dev
```

---

## 📝 提交规范

使用 [Conventional Commits](https://www.conventionalcommits.org/)：

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**类型**:
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建/工具

**示例**:
```
feat(user): add email validation

fix(cache): resolve redis connection timeout

docs(api): update HTTP endpoint documentation
```

---

## 🧪 测试要求

- 新功能必须包含测试
- Bug 修复必须包含回归测试
- 覆盖率不得低于 80%

```bash
make test
make coverage
```

---

## 📊 代码审查

所有 PR 需要：
- 至少 1 个审查者批准
- 所有 CI 检查通过
- 覆盖率不下降

---

## 📚 资源

- [开发指南](docs/development/setup.md)
- [架构文档](docs/architecture/clean-architecture.md)
- [API 文档](docs/api/README.md)

---

## ❓ 问题?

- 创建 [Issue](https://github.com/yourusername/golang/issues)
- 查看现有 [Issues](https://github.com/yourusername/golang/issues)

---

感谢你的贡献！

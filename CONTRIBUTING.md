# 贡献指南

感谢您对 Go语言现代化项目 的关注！我们欢迎所有形式的贡献。

## 🤝 如何贡献

### 1. 报告问题

- 使用 [GitHub Issues](https://github.com/your-repo/issues) 报告bug
- 提供详细的错误描述和复现步骤
- 包含系统环境信息

### 2. 功能建议

- 在 Issues 中提出新功能建议
- 描述功能的使用场景和价值
- 讨论技术实现方案

### 3. 代码贡献

- Fork 项目仓库
- 创建功能分支
- 提交 Pull Request

## 🚀 开发环境设置

### 前置要求

- Go 1.24+
- Git
- 代码编辑器 (推荐 VS Code)

### 本地开发

```bash
# 克隆仓库
git clone https://github.com/your-username/golang-modernization.git
cd golang-modernization

# 安装依赖
go mod download

# 运行测试
go test ./...

# 运行基准测试
go test -bench=. ./...
```

## 📝 代码规范

### Go 代码风格

- 遵循 [Effective Go](https://golang.org/doc/effective_go.html)
- 使用 `gofmt` 格式化代码
- 运行 `golangci-lint` 检查代码质量

### 提交信息规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 格式：

```text
feat: 添加新功能
fix: 修复bug
docs: 更新文档
style: 代码格式调整
refactor: 代码重构
test: 添加测试
chore: 构建过程或辅助工具的变动
```

### 分支命名规范

- `feature/功能名称`: 新功能开发
- `fix/问题描述`: 问题修复
- `docs/文档更新`: 文档更新
- `refactor/重构描述`: 代码重构

## 🧪 测试要求

### 单元测试

- 所有新功能必须有测试覆盖
- 测试覆盖率不低于 85%
- 使用表驱动测试模式

### 基准测试

- 性能相关功能必须有基准测试
- 提供性能对比数据
- 确保性能回归检测

### 集成测试

- 复杂功能需要集成测试
- 测试真实使用场景
- 验证系统整体功能

## 📚 文档要求

### 代码注释

- 所有导出的函数必须有注释
- 复杂逻辑需要详细说明
- 使用 Go 标准注释格式

### README 文档

- 每个模块都有清晰的 README
- 包含安装和使用说明
- 提供实际使用示例

### API 文档

- 使用 godoc 生成文档
- 包含参数和返回值说明
- 提供使用示例

## 🔄 Pull Request 流程

### 1. 准备工作

- 确保代码通过所有测试
- 更新相关文档
- 添加必要的测试用例

### 2. 提交 PR

- 使用清晰的标题描述变更
- 在描述中详细说明变更内容
- 关联相关的 Issues

### 3. 代码审查

- 至少需要一名维护者审查
- 根据反馈进行修改
- 确保所有检查通过

### 4. 合并

- 维护者确认后合并
- 删除功能分支
- 更新版本和发布说明

## 🏷️ 版本发布

### 版本号规范

使用 [Semantic Versioning](https://semver.org/) 格式：`MAJOR.MINOR.PATCH`

- `MAJOR`: 不兼容的API变更
- `MINOR`: 向后兼容的功能新增
- `PATCH`: 向后兼容的问题修复

### 发布流程

1. 更新版本号
2. 更新 CHANGELOG.md
3. 创建 Git 标签
4. 发布到 GitHub Releases

## 📞 获取帮助

### 社区资源

- [GitHub Discussions](https://github.com/your-repo/discussions)
- [项目 Wiki](https://github.com/your-repo/wiki)
- [技术博客](https://your-blog.com)

### 联系方式

- 邮箱: <your-email@example.com>
- 微信: your-wechat
- 技术交流群: your-group

## 🙏 致谢

感谢所有为项目做出贡献的开发者！您的贡献让这个项目变得更好。

---

**注意**: 参与贡献即表示您同意遵守项目的 [行为准则](CODE_OF_CONDUCT.md) 和 [许可证](LICENSE)。

# 贡献指南

> **欢迎加入 Go 1.25 学习项目！**  
> 感谢您对项目的关注！我们欢迎所有形式的贡献，无论是文档改进、代码示例、问题反馈还是功能建议。

---

## 🎯 贡献方式

我们欢迎以下类型的贡献：

### 📝 文档贡献

- **修正错误**: 修正文档中的错别字、语法错误
- **补充内容**: 添加缺失的说明、示例或最佳实践
- **改进结构**: 优化文档组织和可读性
- **翻译工作**: 帮助翻译成其他语言

### 💻 代码贡献

- **代码示例**: 添加实用的代码示例和用例
- **性能测试**: 提供基准测试和性能数据
- **工具开发**: 开发辅助工具和脚本
- **问题修复**: 修复代码中的 bug

### 🐛 问题反馈

- **Bug 报告**: 报告文档或代码中的问题
- **功能建议**: 提出新功能或改进建议
- **使用反馈**: 分享使用经验和心得

### 🌟 社区贡献

- **项目推广**: 在社交媒体分享项目
- **问题解答**: 帮助其他用户解决问题
- **案例分享**: 分享实际使用案例

## 🚀 开发环境设置

### 前置要求

- **Go 1.25+** (必需)
- **Git** (必需)
- 代码编辑器 (推荐 VS Code 或 GoLand)
- Markdown 编辑器 (可选)

### 本地开发

```bash
# 1. Fork 并克隆仓库
git clone https://github.com/your-username/golang.git
cd golang

# 2. 创建功能分支
git checkout -b feature/your-feature-name

# 3. 运行示例代码
cd docs/02-Go语言现代化/12-Go-1.25运行时优化/examples/gc_optimization
go test -bench=. -benchmem

# 4. 运行项目统计工具
cd scripts
go run project_stats.go ..

# 5. 提交更改
git add .
git commit -m "feat: 添加新功能"
git push origin feature/your-feature-name
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

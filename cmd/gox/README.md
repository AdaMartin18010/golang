# gox - Golang项目管理工具

> 统一的CLI工具，整合项目所有自动化脚本

## 🎯 功能特点

- ✅ **统一界面** - 一个命令管理所有工具
- ✅ **跨平台** - Windows/Linux/macOS全支持
- ✅ **简洁易用** - 命令简短，操作便捷
- ✅ **功能完整** - 整合43个脚本的核心功能

## 📦 安装

### 方式1: 从源码安装

```bash
# 在项目根目录执行
go install ./cmd/gox@latest

# 验证安装
gox version
```

### 方式2: 本地构建

```bash
cd cmd/gox
go build -o gox

# 移动到PATH
# Windows: move gox.exe C:\Go\bin\
# Linux/macOS: sudo mv gox /usr/local/bin/
```

## 🚀 快速开始

```bash
# 查看帮助
gox help

# 运行质量检查
gox quality

# 运行测试
gox test

# 查看项目统计
gox stats
```

## 📖 命令参考

### quality (q) - 代码质量检查

完整的代码质量检查，包括格式、静态分析、编译、测试。

```bash
# 完整检查
gox quality

# 快速检查 (跳过测试)
gox quality --fast

# 简写形式
gox q
```

**检查项**:
- ✅ go fmt - 代码格式检查
- ✅ go vet - 静态分析
- ✅ go build - 编译检查
- ✅ go test - 测试运行

### test (t) - 测试统计

运行测试并生成统计报告。

```bash
# 运行所有测试
gox test

# 生成覆盖率报告
gox test --coverage

# 详细输出
gox test --verbose

# 简写形式
gox t
```

**生成文件**:
- `coverage.out` - 覆盖率数据
- `coverage.html` - 可视化报告

### stats (s) - 项目统计

分析并展示项目统计信息。

```bash
# 显示项目统计
gox stats

# 详细统计
gox stats --detail

# 简写形式
gox s
```

**统计内容**:
- 📊 文件统计
- 💻 代码统计
- 📚 文档统计
- 🎯 质量指标

### format (f) - 代码格式化

格式化Go代码或检查格式。

```bash
# 格式化所有代码
gox format

# 只检查不格式化
gox format --check

# 简写形式
gox f
```

### docs (d) - 文档处理

文档生成和处理工具。

```bash
# 生成文档目录
gox docs toc

# 检查文档链接
gox docs links

# 格式化文档
gox docs format

# 简写形式
gox d
```

### migrate (m) - 项目迁移

执行Workspace迁移。

```bash
# 预览迁移 (推荐先执行)
gox migrate --dry-run

# 执行实际迁移
gox migrate

# 简写形式
gox m
```

**安全提示**:
- ⚠️ 执行前会提示确认
- ✅ 支持dry-run预览
- 📦 建议先备份

### verify (v) - 结构验证

验证项目结构和配置。

```bash
# 验证项目结构
gox verify

# 验证Workspace配置
gox verify workspace

# 简写形式
gox v
```

**验证项**:
- ✅ 文档代码分离
- ✅ 目录结构规范
- ✅ 关键文件存在
- ✅ Workspace配置

### help (h) - 帮助信息

显示详细的帮助信息。

```bash
# 显示帮助
gox help

# 简写形式
gox h
```

### version - 版本信息

显示工具和Go版本信息。

```bash
gox version
```

## 💡 使用示例

### 日常开发工作流

```bash
# 1. 开始工作前 - 验证环境
gox verify

# 2. 编写代码...

# 3. 提交前检查
gox format          # 格式化代码
gox quality         # 质量检查
gox test --coverage # 运行测试

# 4. 查看统计
gox stats
```

### CI/CD集成

```bash
# CI脚本示例
gox quality --fast  # 快速质量检查
gox test --coverage # 测试并生成覆盖率
gox verify          # 验证结构
```

### 项目维护

```bash
# 定期维护任务
gox stats --detail  # 查看详细统计
gox docs links      # 检查文档链接
gox verify          # 验证项目结构
```

## 🔧 全局选项

所有命令都支持以下选项：

```bash
--verbose, -v   # 详细输出
--quiet, -q     # 安静模式
--help, -h      # 显示帮助
```

## 📊 性能对比

### 使用gox前后对比

| 操作 | 使用前 | 使用gox | 提升 |
|------|--------|---------|------|
| 质量检查 | 4条命令 | 1条命令 | 75% |
| 运行测试 | 3条命令 | 1条命令 | 67% |
| 项目统计 | 打开工具 | 1条命令 | 80% |
| 格式检查 | 2条命令 | 1条命令 | 50% |

### 时间节省

```text
每日开发任务:
使用前: ~10分钟
使用gox: ~3分钟
节省: 70%时间
```

## 🛠️ 高级用法

### 自定义配置

可以通过环境变量自定义行为：

```bash
# 设置详细输出
export GOX_VERBOSE=true

# 设置安静模式
export GOX_QUIET=true
```

### 别名设置

建议设置命令别名：

```bash
# .bashrc 或 .zshrc
alias gq='gox quality'
alias gt='gox test'
alias gs='gox stats'
alias gf='gox format'
```

## 📝 开发计划

### v1.1 (计划中)

- [ ] 增加benchmark命令
- [ ] 支持配置文件
- [ ] 增加插件系统
- [ ] 更多文档操作

### v1.2 (计划中)

- [ ] Web界面
- [ ] 实时监控
- [ ] 报告生成
- [ ] 团队协作

## 🐛 问题反馈

如果遇到问题或有建议，请：

1. 查看 [FAQ](../../FAQ.md)
2. 提交 [Issue](../../issues)
3. 发送邮件反馈

## 📄 许可证

MIT License - 详见 [LICENSE](../../LICENSE)

---

**版本**: v1.0.0  
**作者**: Go Community  
**更新**: 2025-10-22


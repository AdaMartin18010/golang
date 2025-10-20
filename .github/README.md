# GitHub Actions 工作流说明

本目录包含Go语言文档库的自动化工作流配置。

---

## 📋 工作流列表

### 1. 文档质量检查 (`docs-check.yml`)

**触发条件**:
- Push到main分支（文档文件变更）
- Pull Request到main分支（文档文件变更）
- 手动触发

**检查项目**:

#### 链接验证
- ✅ 验证所有Markdown文件中的链接
- ✅ 支持相对链接和绝对链接
- ✅ 自动重试失败的链接
- ⚠️ 忽略localhost和GitHub Issues链接

#### 格式检查
- ✅ 检查H1标题格式（不应包含编号）
- ✅ 检查简介完整性
- ✅ 检查元数据完整性

#### 拼写检查
- ✅ 使用typos进行拼写检查
- ✅ 支持自定义词典
- ✅ 忽略Git哈希和版本号

**配置文件**:
- `.markdown-link-check.json` (自动生成)
- `.typos.toml` (自动生成)

---

### 2. 文档部署 (`docs-deploy.yml`)

**触发条件**:
- Push到main分支（文档变更）
- 手动触发

**部署流程**:
1. 构建文档站点
2. 生成首页（index.html）
3. 上传到GitHub Pages
4. 自动部署

**访问地址**:
- https://AdaMartin18010.github.io/golang/

**特性**:
- ✅ 自动生成美观的首页
- ✅ 展示项目统计数据
- ✅ 模块化导航
- ✅ 响应式设计

---

## 🚀 使用方法

### 查看检查结果

1. 访问仓库的**Actions**标签页
2. 选择相应的工作流运行
3. 查看详细日志和报告

### 手动触发

#### 方式1: GitHub网页

1. 访问仓库的**Actions**标签页
2. 选择工作流
3. 点击**Run workflow**按钮

#### 方式2: GitHub CLI

```bash
# 触发文档检查
gh workflow run docs-check.yml

# 触发文档部署
gh workflow run docs-deploy.yml
```

---

## 🔧 配置说明

### 链接检查配置

在`docs-check.yml`中自动生成的`.markdown-link-check.json`：

```json
{
  "ignorePatterns": [
    {"pattern": "^http://localhost"},
    {"pattern": "^https://github.com/.*/issues/new"},
    {"pattern": "^#"}
  ],
  "timeout": "10s",
  "retryOn429": true,
  "retryCount": 3
}
```

### 拼写检查配置

在`docs-check.yml`中自动生成的`.typos.toml`：

```toml
[default]
extend-ignore-re = [
  "[0-9a-f]{7,40}",  # Git哈希
  "v[0-9]+\\.[0-9]+\\.[0-9]+",  # 版本号
]

[files]
extend-exclude = [
  "*.sum",
  "*.mod",
  "*.json",
]
```

---

## 📊 检查报告

每次运行后会生成摘要报告，包含：

- ✅/❌ 链接验证结果
- ✅/❌ 格式检查结果
- ✅/⚠️ 拼写检查结果
- 📝 详细日志链接

---

## 🛠️ 本地测试

### 链接检查

```bash
# 安装工具
npm install -g markdown-link-check

# 检查单个文件
markdown-link-check docs/README.md

# 检查所有文档
find docs -name "*.md" -not -path "*/archive/*" -exec markdown-link-check {} \;
```

### 格式检查

```bash
# 检查H1标题
grep -rn "^# [0-9]" docs/*.md

# 检查简介
grep -L "📚 \*\*简介\*\*" docs/**/*.md

# 检查元数据
grep -L "文档维护者" docs/**/*.md
```

### 拼写检查

```bash
# 安装typos
cargo install typos-cli

# 运行检查
typos docs/
```

---

## 🔍 故障排查

### 链接检查失败

**可能原因**:
1. 链接确实失效
2. 网络超时
3. 目标服务器限流

**解决方案**:
1. 修复或移除失效链接
2. 增加超时时间
3. 添加到忽略列表

### 格式检查失败

**可能原因**:
1. H1标题包含编号
2. 缺少简介或元数据

**解决方案**:
1. 按照文档规范修正格式
2. 使用格式化工具批量修复

### 拼写检查警告

**可能原因**:
1. 真实的拼写错误
2. 专业术语未在词典中

**解决方案**:
1. 修正拼写错误
2. 添加到`.typos.toml`词典

---

## 🎯 最佳实践

### PR检查

1. **提交前本地测试**
   ```bash
   # 运行本地检查
   ./scripts/validate_links.ps1
   ./scripts/format_align_docs_v2.ps1 -DryRun
   ```

2. **PR描述中说明**
   - 文档变更内容
   - 检查结果截图

3. **响应CI反馈**
   - 及时修复CI发现的问题
   - 不要忽略警告

### 持续改进

1. **定期更新工具**
   - 更新GitHub Actions版本
   - 更新依赖工具

2. **优化检查规则**
   - 根据实际情况调整
   - 添加新的检查项

3. **收集反馈**
   - 团队反馈
   - 用户反馈

---

## 📚 相关资源

### GitHub Actions

- [Actions文档](https://docs.github.com/en/actions)
- [Workflow语法](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions)

### 检查工具

- [markdown-link-check](https://github.com/tcort/markdown-link-check)
- [typos](https://github.com/crate-ci/typos)
- [GitHub Pages](https://pages.github.com/)

### 项目文档

- [文档规范](../docs/DOCUMENT_STANDARD.md)
- [维护指南](../CONTRIBUTING.md)

---

**维护团队**: Go Documentation Team  
**最后更新**: 2025年10月20日


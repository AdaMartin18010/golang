# Markdown 格式检查和修复工具使用指南

本项目包含一套完整的 Markdown 格式检查和自动修复工具，帮助保持文档质量和一致性。

## 📋 目录

- [工具概览](#工具概览)
- [配置文件](#配置文件)
- [使用方法](#使用方法)
- [规则说明](#规则说明)
- [常见问题](#常见问题)

## 🛠️ 工具概览

### 1. 基础检查工具（无需外部依赖）

**脚本**: `scripts/check-markdown-basic.ps1`

**功能**:

- 检查多个连续空行 (MD012)
- 检查行尾空格 (MD009)
- 检查代码块缺少语言标识 (MD040)
- 检查空链接 (MD042)
- 检查文件末尾换行符 (MD047)
- 检查重复的版本信息块

**使用**:

```powershell
# 仅检查
.\scripts\check-markdown-basic.ps1 -Path docs

# 检查并自动修复
.\scripts\check-markdown-basic.ps1 -Path docs -Fix
```

### 2. 完整检查工具（需要 markdownlint-cli）

**脚本**: `scripts/check-markdown-format.ps1`

**功能**:

- 使用 markdownlint-cli 进行全面检查
- 支持自动修复大部分格式问题
- 使用项目配置文件 `.markdownlint.json`

**安装依赖**:

```bash
# 使用 npm
npm install -g markdownlint-cli

# 或使用 pnpm
pnpm add -g markdownlint-cli

# 或使用 yarn
yarn global add markdownlint-cli
```

**使用**:

```powershell
# 仅检查
.\scripts\check-markdown-format.ps1 -Path docs

# 检查并自动修复
.\scripts\check-markdown-format.ps1 -Path docs -Fix

# 显示详细信息
.\scripts\check-markdown-format.ps1 -Path docs -Verbose
```

### 3. 目录链接修复工具

**脚本**: `scripts/fix-toc-links.ps1`

**功能**:

- 自动修复 Markdown 文档中的目录链接
- 根据标题生成正确的锚点
- 处理各种特殊格式的标题

**使用**:

```powershell
.\scripts\fix-toc-links.ps1 -Path docs
```

### 4. 重复版本信息修复工具

**脚本**: `scripts/fix-all-duplicates.ps1`

**功能**:

- 检测并删除重复的版本信息块
- 清理多余的分隔符

**使用**:

```powershell
.\scripts\fix-all-duplicates.ps1 -Path docs
```

### 5. 全面修复工具（一键修复）

**脚本**: `scripts/fix-markdown-all.ps1`

**功能**:

- 依次执行所有修复操作
- 包含最终验证步骤

**使用**:

```powershell
.\scripts\fix-markdown-all.ps1 -Path docs
```

## ⚙️ 配置文件

### .markdownlint.json

主配置文件，定义了所有 Markdown 格式规则：

```json
{
  "default": true,
  "MD001": true,          // 标题层级递增
  "MD003": { "style": "atx" },  // 使用 ATX 风格标题
  "MD004": { "style": "dash" }, // 使用 - 作为无序列表标记
  "MD007": { "indent": 2 },     // 列表缩进 2 个空格
  "MD009": true,          // 禁止行尾空格
  "MD010": true,          // 禁止使用硬制表符
  "MD012": true,          // 禁止多个连续空行
  "MD013": false,         // 不限制行长度
  "MD022": true,          // 标题周围需要空行
  "MD023": true,          // 标题不能缩进
  "MD024": { "siblings_only": true }, // 允许不同层级有相同标题
  "MD025": true,          // 一个文档只能有一个 H1
  "MD027": true,          // 列表标记后的空格
  "MD031": true,          // 代码块周围需要空行
  "MD032": true,          // 列表周围需要空行
  "MD033": false,         // 允许内联 HTML
  "MD034": false,         // 允许裸链接
  "MD037": true,          // 强调标记周围不应有空格
  "MD040": true,          // 代码块必须指定语言
  "MD042": true,          // 禁止空链接
  "MD047": true,          // 文件必须以换行符结尾
  "MD048": { "style": "backtick" }, // 使用反引号代码块
  "MD049": { "style": "underscore" }, // 使用下划线斜体
  "MD050": { "style": "asterisk" }    // 使用星号加粗
}
```

### .vscode/settings.json

VS Code 集成配置，实现：

- 保存时自动格式化 Markdown 文件
- 实时 linting
- 自动移除行尾空格
- 文件末尾添加换行符

### .vscode/extensions.json

推荐的 VS Code 扩展：

- `davidanson.vscode-markdownlint` - Markdown linting
- `yzhang.markdown-all-in-one` - Markdown 增强功能
- `bierner.markdown-preview-github-styles` - GitHub 风格预览

## 📖 使用方法

### 日常使用

1. **在 VS Code 中编辑**:
   - 安装推荐的扩展
   - 保存时自动格式化
   - 实时看到 linting 错误

2. **定期检查**:

   ```powershell
   # 快速检查
   .\scripts\check-markdown-basic.ps1 -Path docs

   # 完整检查（需要 markdownlint-cli）
   .\scripts\check-markdown-format.ps1 -Path docs
   ```

3. **批量修复**:

   ```powershell
   # 一键修复所有问题
   .\scripts\fix-markdown-all.ps1 -Path docs
   ```

### CI/CD 集成

可以将检查脚本集成到 CI/CD 流程中：

```yaml
# GitHub Actions 示例
- name: Check Markdown
  run: |
    npm install -g markdownlint-cli
    powershell -File scripts/check-markdown-format.ps1 -Path docs
```

## 📚 规则说明

### 启用的主要规则

| 规则 | 说明 | 可自动修复 |
|------|------|-----------|
| MD001 | 标题层级递增 | ❌ |
| MD009 | 禁止行尾空格 | ✅ |
| MD010 | 禁止硬制表符 | ✅ |
| MD012 | 禁止多个连续空行 | ✅ |
| MD022 | 标题周围需要空行 | ✅ |
| MD031 | 代码块周围需要空行 | ✅ |
| MD032 | 列表周围需要空行 | ✅ |
| MD040 | 代码块必须指定语言 | ❌ |
| MD042 | 禁止空链接 | ❌ |
| MD047 | 文件末尾需要换行符 | ✅ |

### 禁用的规则

| 规则 | 说明 | 原因 |
|------|------|------|
| MD013 | 行长度限制 | Go 文档常有长代码行 |
| MD033 | 内联 HTML | 文档需要高级格式化 |
| MD041 | 首行必须是 H1 | 文档有版本信息块 |
| MD051 | 链接片段存在 | 某些链接是外部的 |

## ❓ 常见问题

### 1. markdownlint-cli 未安装怎么办？

可以使用基础检查工具：

```powershell
.\scripts\check-markdown-basic.ps1 -Path docs -Fix
```

### 2. 如何临时禁用某个规则？

在 Markdown 文件中使用注释：

```markdown
<!-- markdownlint-disable MD033 -->
<div>HTML 内容</div>
<!-- markdownlint-enable MD033 -->
```

### 3. 如何检查单个文件？

```powershell
# 基础检查
$file = "docs/README.md"
Get-Content $file | Select-String "  $"  # 检查行尾空格

# 完整检查（需要 markdownlint-cli）
markdownlint $file
```

### 4. 自动修复会改变什么？

自动修复会：

- ✅ 删除行尾空格
- ✅ 删除多余的空行
- ✅ 在代码块周围添加空行
- ✅ 在文件末尾添加换行符
- ✅ 删除重复的版本信息块
- ✅ 修复目录链接

自动修复不会：

- ❌ 修改文档内容和语义
- ❌ 更改代码示例
- ❌ 修改标题文本

### 5. 如何处理误报？

1. **临时禁用**: 使用 `<!-- markdownlint-disable -->`
2. **更新配置**: 修改 `.markdownlint.json`
3. **提交 Issue**: 如果是工具 bug，提交到对应仓库

## 🎯 最佳实践

1. **提交前检查**: 养成提交前运行检查的习惯
2. **保存即修复**: 在 VS Code 中启用保存时自动格式化
3. **定期全面检查**: 每周运行一次全面检查
4. **团队统一**: 确保团队成员使用相同的配置
5. **CI/CD 集成**: 在 CI 流程中添加自动检查

## 📞 支持

如有问题，请：
1. 查看本文档的常见问题部分
2. 查看工具的详细输出信息
3. 提交 Issue 到项目仓库

---

**最后更新**: 2025-10-30
**维护者**: AI Assistant

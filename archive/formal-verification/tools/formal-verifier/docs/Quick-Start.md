# Quick Start Guide

**Go Formal Verifier (FV)** 快速入门指南

---

## 📋 目录

1. [安装](#安装)
2. [第一次运行](#第一次运行)
3. [基础命令](#基础命令)
4. [配置文件](#配置文件)
5. [常见场景](#常见场景)
6. [下一步](#下一步)

---

## 安装

### 从源码构建

```bash
# 克隆仓库
git clone https://github.com/your-org/formal-verifier.git
cd formal-verifier

# 构建工具
go build -o fv ./cmd/fv

# 安装到系统路径（可选）
sudo mv fv /usr/local/bin/

# 验证安装
fv version
```

### 使用 Go Install

```bash
go install github.com/your-org/formal-verifier/cmd/fv@latest
```

---

## 第一次运行

### 1. 快速分析

最简单的方式是直接在项目根目录运行：

```bash
cd your-go-project
fv analyze
```

这将：

- 递归扫描当前目录的所有 Go 文件
- 检查并发问题、类型安全、复杂度等
- 在终端输出文本格式报告

### 2. 交互式模式

如果你更喜欢图形化菜单：

```bash
fv interactive
```

这将启动交互式菜单，引导你完成分析过程。

### 3. 生成HTML报告

```bash
fv analyze --format=html --output=report.html
```

然后在浏览器中打开 `report.html` 查看可视化报告。

---

## 基础命令

### 项目分析

```bash
# 基本分析
fv analyze

# 指定目录
fv analyze --dir=./mypackage

# 生成HTML报告
fv analyze --format=html --output=report.html

# 生成JSON报告（用于自动化）
fv analyze --format=json --output=report.json

# 包含测试文件
fv analyze --include-tests

# 排除特定目录
fv analyze --exclude="vendor/*,testdata/*"
```

### 配置管理

```bash
# 生成默认配置文件
fv init-config

# 生成严格模式配置（适合CI/CD）
fv init-config --output=.fv-strict.yaml --strict

# 使用配置文件
fv analyze --config=.fv.yaml
```

### 交互式模式

```bash
# 启动交互式界面
fv interactive

# 使用配置文件启动
fv interactive --config=.fv.yaml
```

### 查看帮助

```bash
# 查看所有命令
fv help

# 查看版本
fv version
```

---

## 配置文件

### 创建配置文件

```bash
fv init-config --output=.fv.yaml
```

这将生成一个包含所有选项的配置文件模板。

### 基础配置示例

`.fv.yaml`:

```yaml
project:
  root_dir: .
  recursive: true
  include_tests: false
  exclude_patterns:
    - vendor
    - testdata
    - .git

report:
  format: html
  output_path: fv-report.html

rules:
  complexity:
    cyclomatic_threshold: 10
    max_function_lines: 50

output:
  fail_on_error: false
  min_quality_score: 0
```

### 使用配置

```bash
fv analyze --config=.fv.yaml
```

---

## 常见场景

### 场景 1: 本地开发

**目标**: 快速检查代码质量

```bash
# 方式1: 直接分析，输出到终端
fv analyze

# 方式2: 生成HTML报告查看详情
fv analyze --format=html --output=report.html
open report.html  # macOS
xdg-open report.html  # Linux
start report.html  # Windows
```

### 场景 2: Pull Request 检查

**目标**: 在PR中查看代码质量变化

```bash
# 生成Markdown报告
fv analyze --format=markdown --output=pr-report.md

# 在PR描述中包含报告
cat pr-report.md
```

### 场景 3: CI/CD 集成

**目标**: 自动化质量检查

```bash
# 1. 生成严格模式配置
fv init-config --output=.fv-ci.yaml --strict

# 2. 在CI中运行
fv analyze \
  --config=.fv-ci.yaml \
  --no-color \
  --fail-on-error
```

### 场景 4: 遗留代码评估

**目标**: 评估现有项目的代码质量

```bash
# 1. 生成完整报告
fv analyze --format=html --output=assessment.html

# 2. 生成JSON用于进一步分析
fv analyze --format=json --output=assessment.json

# 3. 使用jq提取质量分数
jq -r '.stats.quality_score' assessment.json
```

### 场景 5: 团队标准化

**目标**: 统一团队的代码质量标准

```bash
# 1. 创建团队配置
fv init-config --output=.fv-team.yaml

# 2. 调整配置（例如更严格的复杂度要求）
# 编辑 .fv-team.yaml:
#   complexity:
#     cyclomatic_threshold: 5
#     max_function_lines: 30

# 3. 提交配置到代码仓库
git add .fv-team.yaml
git commit -m "Add team FV configuration"

# 4. 团队成员使用统一配置
fv analyze --config=.fv-team.yaml
```

---

## 理解报告

### 文本报告示例

```text
========================================
📊 分析报告
========================================

项目: ./myproject
文件数: 45
总行数: 12,543
问题数: 23
质量评分: 87/100

----------------------------------------
问题统计:
  ❌ 错误: 3
  ⚠️  警告: 15
  ℹ️  提示: 5

按类别:
  并发: 5
  类型: 8
  数据流: 3
  复杂度: 7
----------------------------------------

❌ 错误:
  [concurrency] main.go:45:10
    Potential goroutine leak detected
    💡 建议: Add proper goroutine cleanup

⚠️  警告:
  [complexity] handler.go:123:1
    Function processRequest has cyclomatic complexity 15 (threshold: 10)
    💡 建议: Consider breaking down into smaller functions
```

### HTML报告预览

HTML报告包含：

- 📊 **概览仪表板**: 质量分数、问题统计
- 📁 **文件列表**: 按问题数排序
- 🔍 **问题详情**: 每个问题的位置和建议
- 📈 **趋势图表**: 可视化质量指标

### JSON报告结构

```json
{
  "project_info": {
    "root_dir": "./myproject",
    "total_files": 45
  },
  "stats": {
    "total_issues": 23,
    "error_count": 3,
    "warning_count": 15,
    "info_count": 5,
    "quality_score": 87
  },
  "issues": [
    {
      "severity": "error",
      "category": "concurrency",
      "file": "main.go",
      "line": 45,
      "column": 10,
      "message": "Potential goroutine leak detected",
      "suggestion": "Add proper goroutine cleanup"
    }
  ]
}
```

---

## 最佳实践

### 1. 渐进式采用

不要一开始就要求完美：

```bash
# 第一周: 只看错误
fv analyze | grep "❌"

# 第二周: 处理高复杂度
fv analyze --config=.fv.yaml  # cyclomatic_threshold: 15

# 第三周: 收紧阈值
# 调整 .fv.yaml: cyclomatic_threshold: 10

# 最终: 启用所有检查
fv analyze --config=.fv-strict.yaml
```

### 2. 定期分析

添加到开发流程：

```bash
# Git hook (pre-commit)
#!/bin/bash
fv analyze --format=text --fail-on-error || exit 1

# 或使用 pre-commit 框架
# .pre-commit-config.yaml
repos:
  - repo: local
    hooks:
      - id: fv-analysis
        name: FV Analysis
        entry: fv analyze --fail-on-error
        language: system
```

### 3. 配置版本控制

将配置文件提交到仓库：

```bash
git add .fv.yaml
git commit -m "Add FV configuration"
```

### 4. 团队共享

在 README 中添加：

```markdown
## 代码质量

我们使用 FV 工具进行代码质量检查：

\`\`\`bash
# 安装 FV
go install github.com/your-org/formal-verifier/cmd/fv@latest

# 运行检查
fv analyze --config=.fv.yaml
\`\`\`

当前质量分数: ![Quality](https://img.shields.io/badge/FV%20Quality-87%25-green)
```

---

## 常见问题

### Q: 分析很慢怎么办？

A: 调整并发设置：

```yaml
# .fv.yaml
analysis:
  workers: 8  # 增加worker数量
  max_file_size: 512  # 跳过大文件
```

### Q: 如何忽略某些文件？

A: 使用排除模式：

```bash
fv analyze --exclude="vendor/*,generated/*,*.pb.go"
```

或在配置文件中：

```yaml
project:
  exclude_patterns:
    - vendor
    - "*_gen.go"
    - "testdata"
```

### Q: 如何在CI中使用？

A: 参考 [CI/CD Integration Guide](CI-CD-Integration.md)

### Q: 报告显示乱码？

A: 禁用颜色输出：

```bash
fv analyze --no-color
```

或设置环境变量：

```bash
export NO_COLOR=1
fv analyze
```

---

## 下一步

现在你已经掌握了基础用法，可以：

1. 📚 阅读 [详细教程](Tutorial.md) 了解高级功能
2. 🔧 查看 [CI/CD集成指南](CI-CD-Integration.md) 进行自动化
3. ⚙️  探索 [配置参考](Configuration-Reference.md) 了解所有选项
4. 💡 查看 [最佳实践](Best-Practices.md) 学习高效使用技巧

---

## 获取帮助

- 📖 文档: [https://github.com/your-org/formal-verifier/docs](https://github.com/your-org/formal-verifier/docs)
- 🐛 问题反馈: [https://github.com/your-org/formal-verifier/issues](https://github.com/your-org/formal-verifier/issues)
- 💬 讨论: [https://github.com/your-org/formal-verifier/discussions](https://github.com/your-org/formal-verifier/discussions)

---

**开始使用 FV 提升代码质量！** 🚀

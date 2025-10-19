# Scripts 工具集

本目录包含项目维护和自动化工具。

---

## 📋 工具列表

### 1. 项目结构验证 (verify_structure) ⭐ 新增

**功能**: 验证项目结构是否符合重组规范

**使用方法**:

```bash
# Windows
powershell -ExecutionPolicy Bypass -File scripts/verify_structure.ps1

# Linux/macOS
bash scripts/verify_structure.sh
```

**检查项目**:

- ✅ 文档代码分离（docs/ 无代码文件）
- ✅ 根目录清洁（无临时文件）
- ✅ 目录职责明确
- ✅ 关键文件存在
- ✅ examples/ 结构正确
- ✅ 代码可编译

**示例输出**:

```text
🔍 开始验证项目结构...

📋 规则1: 文档代码分离
➤ docs/ 目录无 .go 文件... ✅
➤ docs/ 目录无 go.mod 文件... ✅

📋 规则2: 根目录清洁
➤ 根目录无 Phase 报告... ✅

✅ 项目结构验证通过！
通过: 22 / 23
```

---

### 2. 代码质量扫描 (scan_code_quality)

**功能**: 完整的代码质量检查（格式化、Vet、编译）

**使用方法**:

```bash
# Windows
powershell -ExecutionPolicy Bypass -File scripts/scan_code_quality.ps1

# Linux/macOS
bash scripts/scan_code_quality.sh
```

**检查项目**:

- 代码格式化检查
- go vet 静态分析
- 编译检查
- 测试运行

---

### 3. 测试统计 (test_summary)

**功能**: 运行所有测试并生成统计报告

**使用方法**:

```bash
powershell -ExecutionPolicy Bypass -File scripts/test_summary.ps1
```

---

### 4. 项目统计工具 (project_stats)

**功能**: 自动分析和统计项目信息

**使用方法**:

```bash
# 在项目根目录运行
cd scripts/project_stats
go run main.go

# 或指定目录
go run main.go /path/to/project
```

**输出内容**:

- 📊 总体统计 (文件数、行数、字数)
- 💻 代码统计 (示例数、测试数)
- 📚 文档分类
- 📁 文件类型分布
- 🎯 项目指标
- ✨ 质量评估

---

### 5. 变更日志生成器 (gen_changelog)

**功能**: 自动生成项目变更日志

**使用方法**:

```bash
cd scripts/gen_changelog
go mod tidy

# 生成变更日志
echo "Added PGO example" | VERSION=v2025.09-P1 go run main.go
```

---

### 6. 报告组织工具 (organize_reports)

**功能**: 自动整理和归档项目报告

**使用方法**:

```bash
powershell -ExecutionPolicy Bypass -File scripts/organize_reports.ps1
```

---

## 🚀 快速开始

### 初始化

```bash
cd scripts/project_stats
go mod tidy
```

### 运行项目统计

```bash
cd scripts/project_stats
go run main.go
```

示例输出:

```text
🔍 Go 1.23+ 项目统计分析
======================================================================

📊 总体统计
----------------------------------------------------------------------
  总文件数:       1500
  Markdown 文件:  250
  Go 代码文件:    100
  README 文件:    15
  总行数:         45000
  总字数:         120000

💻 代码统计
----------------------------------------------------------------------
  代码示例:       80
  基准测试:       30

🏆 项目总分: 95/100
📈 项目评级: 卓越 (Excellent) 🏆🏆🏆
```

---

## 📝 开发新工具

### 工具开发指南

1. 在 `scripts/` 目录创建新的子目录 (如 `my_tool/`)
2. 在子目录中创建 `main.go` 和 `go.mod` 文件
3. 实现工具逻辑
4. 更新本 README
5. 提交 PR

### 工具规范

- ✅ 使用 Go 1.23++ 编写
- ✅ 提供清晰的使用说明
- ✅ 处理错误情况
- ✅ 输出格式友好

---

## 🤝 贡献

欢迎贡献新工具！请确保:

- 工具有明确的用途
- 代码质量高
- 文档完善

---

## 🎯 推荐工作流

### 日常开发

```bash
# 1. 验证项目结构
powershell -ExecutionPolicy Bypass -File scripts/verify_structure.ps1

# 2. 代码质量检查
powershell -ExecutionPolicy Bypass -File scripts/scan_code_quality.ps1

# 3. 运行测试
powershell -ExecutionPolicy Bypass -File scripts/test_summary.ps1
```

### 提交前检查

```bash
# 完整检查流程
powershell -ExecutionPolicy Bypass -File scripts/verify_structure.ps1
powershell -ExecutionPolicy Bypass -File scripts/scan_code_quality.ps1
```

---

**维护者**: AI Assistant  
**最后更新**: 2025年10月19日

# go doc -http 本地文档服务器（Go 1.23+）

> **版本要求**: Go 1.23++  
> **平台支持**: 所有平台  
> **实验性**: 否（正式特性）  
>

---

## 📚 目录

- [概述](#概述)
- [为什么需要 go doc -http](#为什么需要-go-doc--http)
- [基本使用](#基本使用)
- [功能特性](#功能特性)
- [配置选项](#配置选项)
- [实践场景](#实践场景)
- [与 pkg.go.dev 对比](#与-pkggodev-对比)
- [最佳实践](#最佳实践)
- [常见问题](#常见问题)
- [参考资料](#参考资料)

---

## 概述

Go 1.23+ 为 `go doc` 命令添加了 `-http` 选项,允许开发者在本地启动一个文档服务器,类似于 pkg.go.dev,但完全离线且自动包含项目代码。

### 什么是 go doc -http?

`go doc -http` 启动一个本地 Web 服务器,提供:

- ✅ **标准库文档**: 所有 Go 标准库文档
- ✅ **项目代码文档**: 自动生成当前项目的文档
- ✅ **依赖包文档**: 项目依赖的第三方包文档
- ✅ **交互式浏览**: 网页界面,支持搜索和跳转
- ✅ **实时刷新**: 代码变化自动更新文档

### 核心优势

- ✅ **离线可用**: 无需网络连接
- ✅ **一键启动**: 无需额外安装
- ✅ **自动发现**: 自动包含项目代码
- ✅ **快速预览**: 本地访问,速度快
- ✅ **团队共享**: 可以在局域网内共享

---

## 为什么需要 go doc -http?

### 传统方式的局限

**方式 1: 在线查看 pkg.go.dev**-

- ❌ 需要网络连接
- ❌ 私有项目看不到文档
- ❌ 不包含未发布的代码
- ❌ 可能有延迟

**方式 2: 命令行 go doc**-

```bash
$ go doc fmt.Println
func Println(a ...any) (n int, err error)
```

- ❌ 只能查看单个符号
- ❌ 没有交互式浏览
- ❌ 不能搜索
- ❌ 格式不友好

**方式 3: 老旧的 godoc 工具**-

```bash
# 需要单独安装
$ go install golang.org/x/tools/cmd/godoc@latest
$ godoc -http=:6060
```

- ❌ 需要额外安装
- ❌ 不再维护
- ❌ 功能有限

### Go 1.23+ 的解决方案

```bash
# 一条命令启动文档服务器
$ go doc -http :6060

# 自动打开浏览器
# 包含标准库 + 项目代码 + 依赖包
# 实时更新,离线可用
```

---

## 基本使用

### 1. 启动文档服务器

```bash
# 基本用法: 启动在 localhost:6060
$ go doc -http :6060

# 输出:
# Serving Go documentation at http://localhost:6060
# Press Ctrl+C to stop
```

**效果**:

1. 启动 HTTP 服务器监听 6060 端口
2. 自动打开浏览器访问 `http://localhost:6060`
3. 显示文档首页

---

### 2. 指定端口

```bash
# 使用不同端口
$ go doc -http :8080

# 使用随机端口
$ go doc -http :0
# 输出: Serving at http://localhost:37281 (随机分配)
```

---

### 3. 指定监听地址

```bash
# 监听所有网络接口 (允许局域网访问)
$ go doc -http 0.0.0.0:6060

# 只监听 localhost (默认)
$ go doc -http localhost:6060
$ go doc -http 127.0.0.1:6060
```

---

### 4. 查看特定包

```bash
# 查看标准库包
$ go doc -http :6060 encoding/json

# 查看项目包
$ go doc -http :6060 ./pkg/utils

# 查看第三方包
$ go doc -http :6060 github.com/gin-gonic/gin
```

---

### 5. 后台运行

```bash
# Unix/Linux/macOS
$ go doc -http :6060 &

# 或使用 nohup
$ nohup go doc -http :6060 > /dev/null 2>&1 &

# Windows PowerShell
$ Start-Process go -ArgumentList "doc", "-http", ":6060" -WindowStyle Hidden
```

---

## 功能特性

### 1. 自动打开浏览器

启动时自动打开默认浏览器:

```bash
$ go doc -http :6060
# 自动打开 http://localhost:6060
```

**禁用自动打开**:

```bash
# 使用环境变量
$ BROWSER=none go doc -http :6060
```

---

### 2. 实时代码跳转

文档中的链接可以直接跳转到源码:

- 点击函数名 → 跳转到函数定义
- 点击类型名 → 跳转到类型定义
- 点击包名 → 跳转到包文档

---

### 3. 源码浏览

可以直接查看源代码:

- 语法高亮
- 行号显示
- 跳转到定义
- 跳转到引用

---

### 4. 搜索功能

支持全文搜索:

- 搜索包名
- 搜索函数名
- 搜索类型名
- 搜索常量和变量

---

### 5. 示例代码展示

自动提取并展示 Example 函数:

```go
// Example函数会自动显示在文档中
func ExamplePrintln() {
    fmt.Println("Hello, World!")
    // Output: Hello, World!
}
```

---

## 配置选项

### 命令行选项

```bash
# 完整语法
go doc [-http=addr] [package|symbol]

# 选项说明
-http string
    启动 HTTP 服务器的地址和端口
    格式: [host]:port
    示例: :6060, localhost:8080, 0.0.0.0:6060
```

---

### 环境变量

```bash
# 禁用自动打开浏览器
export BROWSER=none
go doc -http :6060

# 设置端口 (通过别名)
alias godoc='go doc -http :6060'
```

---

### 配置文件 (可选)

创建 shell 别名简化使用:

```bash
# ~/.bashrc 或 ~/.zshrc
alias godoc='go doc -http :6060'
alias godoc-share='go doc -http 0.0.0.0:6060'  # 局域网共享
```

---

## 实践场景

### 场景 1: 本地开发查阅文档

**需求**: 快速查看标准库和项目 API

```bash
# 启动文档服务器
$ go doc -http :6060

# 在浏览器中访问
http://localhost:6060/pkg/encoding/json/
http://localhost:6060/pkg/github.com/your/project/
```

**优势**:

- ✅ 比 pkg.go.dev 更快
- ✅ 包含未发布的代码
- ✅ 离线可用

---

### 场景 2: API 文档预览

**需求**: 编写完 API 后预览文档效果

```bash
# 在项目目录启动
$ cd /path/to/project
$ go doc -http :6060

# 访问项目文档
http://localhost:6060/pkg/github.com/your/project/pkg/api/
```

**检查内容**:

- ✅ 文档注释是否完整
- ✅ Example 函数是否正确
- ✅ 包结构是否清晰
- ✅ 导出符号是否合理

---

### 场景 3: 团队文档共享

**需求**: 在局域网内共享项目文档

```bash
# 监听所有网络接口
$ go doc -http 0.0.0.0:6060

# 团队成员访问 (假设服务器 IP 是 192.168.1.100)
http://192.168.1.100:6060
```

**应用场景**:

- ✅ 内部培训
- ✅ Code Review
- ✅ API 演示
- ✅ 新人入职

---

### 场景 4: CI/CD 文档验证

**需求**: 在 CI 中验证文档完整性

```bash
# .github/workflows/doc-check.yml
- name: Check documentation
  run: |
    # 启动文档服务器
    go doc -http :6060 &
    DOC_PID=$!
    
    # 等待服务启动
    sleep 2
    
    # 检查文档是否可访问
    curl -f http://localhost:6060/pkg/myproject/ || exit 1
    
    # 停止服务器
    kill $DOC_PID
```

---

### 场景 5: 离线开发

**需求**: 飞机上或无网络环境开发

```bash
# 提前启动文档服务器
$ go doc -http :6060

# 所有标准库和项目文档都可访问
# 无需网络连接
```

---

## 与 pkg.go.dev 对比

| 特性 | `go doc -http` | pkg.go.dev |
|------|----------------|------------|
| **网络要求** | ❌ 离线可用 | ✅ 需要网络 |
| **私有代码** | ✅ 完全支持 | ❌ 只有公开代码 |
| **启动速度** | ⚡ 即时 | - 在线服务 |
| **实时更新** | ✅ 代码变化立即生效 | ⚠️ 有延迟 |
| **标准库** | ✅ 完整 | ✅ 完整 |
| **第三方包** | ✅ 项目依赖的包 | ✅ 所有公开包 |
| **搜索功能** | ✅ 本地搜索 | ✅ 全局搜索 |
| **示例代码** | ✅ 支持 | ✅ 支持 |
| **版本历史** | ❌ 当前版本 | ✅ 所有版本 |
| **使用统计** | ❌ | ✅ |

**选择建议**:

- **本地开发**: 使用 `go doc -http` (更快、更方便)
- **查看公开包**: 使用 pkg.go.dev (更全面)
- **私有项目**: 必须使用 `go doc -http`
- **无网络环境**: 必须使用 `go doc -http`

---

## 最佳实践

### 1. 创建 Shell 别名

```bash
# ~/.bashrc 或 ~/.zshrc
alias godoc='go doc -http :6060'
alias godoc-bg='nohup go doc -http :6060 > /dev/null 2>&1 &'
alias godoc-share='go doc -http 0.0.0.0:6060'

# 使用
$ godoc              # 前台运行
$ godoc-bg           # 后台运行
$ godoc-share        # 局域网共享
```

---

### 2. 在 IDE 中集成

#### VS Code

创建任务 (`.vscode/tasks.json`):

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Start Go Documentation Server",
      "type": "shell",
      "command": "go",
      "args": ["doc", "-http", ":6060"],
      "isBackground": true,
      "problemMatcher": []
    }
  ]
}
```

使用: `Ctrl+Shift+P` → `Tasks: Run Task` → `Start Go Documentation Server`

---

### 3. 开发流程集成

```bash
# Makefile
.PHONY: doc
doc:
    @echo "Starting documentation server..."
    @go doc -http :6060

.PHONY: doc-bg
doc-bg:
    @echo "Starting documentation server in background..."
    @nohup go doc -http :6060 > /dev/null 2>&1 &
    @echo "Documentation available at http://localhost:6060"

# 使用
$ make doc        # 前台运行
$ make doc-bg     # 后台运行
```

---

### 4. Docker 容器化

```dockerfile
# Dockerfile
FROM golang:1.25-alpine

WORKDIR /app
COPY . .

# 安装依赖
RUN go mod download

# 暴露端口
EXPOSE 6060

# 启动文档服务器
CMD ["go", "doc", "-http", "0.0.0.0:6060"]
```

```bash
# 构建和运行
$ docker build -t go-docs .
$ docker run -p 6060:6060 go-docs

# 访问
http://localhost:6060
```

---

### 5. 文档质量检查

```bash
#!/bin/bash
# check-docs.sh

# 启动文档服务器
go doc -http :6060 &
DOC_PID=$!

# 等待启动
sleep 2

# 检查关键包的文档
packages=(
    "myproject/pkg/api"
    "myproject/pkg/service"
    "myproject/pkg/repository"
)

failed=0
for pkg in "${packages[@]}"; do
    if ! curl -s "http://localhost:6060/pkg/$pkg/" | grep -q "Overview"; then
        echo "❌ Missing documentation for $pkg"
        failed=1
    else
        echo "✅ Documentation OK for $pkg"
    fi
done

# 清理
kill $DOC_PID

exit $failed
```

---

## 常见问题

### Q1: 端口被占用怎么办?

**A**: 使用不同端口

```bash
# 检查端口占用
$ lsof -i :6060  # Unix/Linux/macOS
$ netstat -ano | findstr :6060  # Windows

# 使用其他端口
$ go doc -http :8080
$ go doc -http :7070

# 使用随机端口
$ go doc -http :0
```

---

### Q2: 如何在后台运行?

**A**: 多种方式

```bash
# 方式 1: 使用 &
$ go doc -http :6060 &

# 方式 2: 使用 nohup
$ nohup go doc -http :6060 > /dev/null 2>&1 &

# 方式 3: 使用 screen/tmux
$ screen -dmS godoc go doc -http :6060

# 查看后台进程
$ ps aux | grep "go doc"

# 停止后台进程
$ killall go
```

---

### Q3: 如何禁用自动打开浏览器?

**A**: 设置 BROWSER 环境变量

```bash
# 临时禁用
$ BROWSER=none go doc -http :6060

# 永久禁用
$ export BROWSER=none
$ go doc -http :6060

# 或使用别名
$ alias godoc='BROWSER=none go doc -http :6060'
```

---

### Q4: 文档没有更新怎么办?

**A**: 刷新浏览器

`go doc -http` 会自动检测代码变化,但浏览器可能缓存了旧页面:

```bash
# 强制刷新浏览器
Ctrl+F5 (Windows/Linux)
Cmd+Shift+R (macOS)

# 或重启服务器
Ctrl+C  # 停止
go doc -http :6060  # 重新启动
```

---

### Q5: 可以自定义文档样式吗?

**A**: ❌ 不支持

`go doc -http` 使用固定的样式,不支持自定义。如果需要自定义文档,考虑使用:

- [gomarkdoc](https://github.com/princjef/gomarkdoc) - 生成 Markdown
- [godoc2md](https://github.com/davecheney/godoc2md) - godoc 转 Markdown

---

### Q6: 支持多个项目吗?

**A**: ✅ 支持 (通过 Go workspace)

```bash
# 使用 Go 1.18+ workspace
$ go work init ./project1 ./project2

# 启动文档服务器
$ go doc -http :6060

# 所有项目的文档都会显示
```

---

## 参考资料

### 官方文档

- 📘 [Go 1.23+ Release Notes](https://go.dev/doc/go1.23#godoc)
- 📘 [go doc command](https://pkg.go.dev/cmd/go#hdr-Show_documentation_for_package_or_symbol)
- 📘 [How to Write Go Code](https://go.dev/doc/code)

### 相关章节

- 🔗 [Go 1.23+ 工具链增强](./README.md)
- 🔗 [文档编写最佳实践](../../最佳实践/文档编写.md)

---

## 更新日志

| 日期 | 版本 | 更新内容 |
|------|------|----------|
| 2025-10-18 | v1.0 | 初始版本,完整的 go doc -http 指南 |

---

**编写者**: AI Assistant  
**审核者**: [待审核]  

---

<p align="center">
  <b>📚 使用 go doc -http 让文档查阅更便捷! 🚀</b>
</p>

---

**文档维护者**: Go Documentation Team  
**最后更新**: 2025年10月20日  
**文档状态**: 完成  
**适用版本**: Go 1.25.3+

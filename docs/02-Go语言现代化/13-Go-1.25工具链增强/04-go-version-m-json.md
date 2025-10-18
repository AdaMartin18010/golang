# go version -m -json 构建信息（Go 1.25）

> **版本要求**: Go 1.25+  
> **平台支持**: 所有平台  
> **实验性**: 否（正式特性）  
> **最后更新**: 2025年10月18日

---

## 📚 目录

- [概述](#概述)
- [为什么需要 JSON 输出](#为什么需要-json-输出)
- [基本使用](#基本使用)
- [输出格式](#输出格式)
- [应用场景](#应用场景)
- [自动化脚本](#自动化脚本)
- [CI/CD 集成](#cicd-集成)
- [实践案例](#实践案例)
- [最佳实践](#最佳实践)
- [常见问题](#常见问题)
- [参考资料](#参考资料)

---

## 概述

Go 1.25 为 `go version -m` 命令添加了 `-json` 选项,允许以 JSON 格式输出二进制文件的构建信息,极大简化了自动化处理和构建审计。

### 什么是 go version -m -json?

`go version -m -json` 提取并以 JSON 格式输出Go二进制文件中嵌入的构建信息:

- ✅ **模块路径和版本**: 主模块和所有依赖
- ✅ **构建设置**: CGO、GOARCH、GOOS 等
- ✅ **Git 信息**: commit hash、是否有修改
- ✅ **编译器版本**: Go 版本
- ✅ **VCS 信息**: 版本控制系统信息

### 核心优势

- ✅ **机器可读**: JSON 格式易于解析
- ✅ **自动化友好**: 适合脚本处理
- ✅ **审计能力**: 追踪依赖版本
- ✅ **SBOM 生成**: Software Bill of Materials
- ✅ **安全扫描**: 识别已知漏洞

---

## 为什么需要 JSON 输出?

### 传统方式的局限

**go version -m 文本输出**:

```bash
$ go version -m ./myapp
./myapp: go1.25.0
 path   example.com/myapp
 mod    example.com/myapp v1.0.0
 dep    github.com/gin-gonic/gin v1.9.1 h1:abc123...
 build -buildmode=exe
 build CGO_ENABLED=1
 build GOARCH=amd64
 build GOOS=linux
```

**问题**:

- ❌ **难以解析**: 文本格式不规范
- ❌ **脚本复杂**: 需要正则表达式或复杂解析
- ❌ **易出错**: 格式变化导致解析失败
- ❌ **批量处理困难**: 处理多个二进制文件复杂

### Go 1.25 的解决方案

```bash
$ go version -m -json ./myapp
{
  "Path": "example.com/myapp",
  "Main": {
    "Path": "example.com/myapp",
    "Version": "v1.0.0",
    "Sum": "h1:abc123..."
  },
  "Deps": [
    {
      "Path": "github.com/gin-gonic/gin",
      "Version": "v1.9.1",
      "Sum": "h1:def456..."
    }
  ],
  "Settings": [
    {"Key": "CGO_ENABLED", "Value": "1"},
    {"Key": "GOARCH", "Value": "amd64"},
    {"Key": "GOOS", "Value": "linux"}
  ]
}
```

**优势**:

- ✅ 标准 JSON 格式
- ✅ 易于解析 (`jq`, 编程语言)
- ✅ 结构化数据
- ✅ 支持批量处理

---

## 基本使用

### 1. 查看单个二进制文件

```bash
# 查看构建信息 (JSON 格式)
$ go version -m -json ./myapp

# 保存到文件
$ go version -m -json ./myapp > build-info.json
```

---

### 2. 批量处理多个二进制文件

```bash
# 处理目录中所有二进制文件
$ go version -m -json ./bin/* > all-build-info.json

# 或分别输出
$ for bin in ./bin/*; do
    go version -m -json "$bin" > "$(basename $bin).json"
done
```

---

### 3. 管道处理

```bash
# 使用 jq 处理
$ go version -m -json ./myapp | jq '.Main.Version'
"v1.0.0"

# 提取依赖列表
$ go version -m -json ./myapp | jq -r '.Deps[].Path'
github.com/gin-gonic/gin
github.com/stretchr/testify
...
```

---

### 4. 查看 Go 版本

```bash
# 提取 Go 版本
$ go version -m -json ./myapp | jq -r '.GoVersion'
go1.25.0
```

---

## 输出格式

### 完整 JSON 结构

```json
{
  "Path": "example.com/myapp",
  "GoVersion": "go1.25.0",
  "Main": {
    "Path": "example.com/myapp",
    "Version": "v1.0.0",
    "Sum": "h1:abc123...",
    "Replace": null
  },
  "Deps": [
    {
      "Path": "github.com/gin-gonic/gin",
      "Version": "v1.9.1",
      "Sum": "h1:def456...",
      "Replace": null
    },
    {
      "Path": "github.com/stretchr/testify",
      "Version": "v1.8.4",
      "Sum": "h1:ghi789...",
      "Replace": null
    }
  ],
  "Settings": [
    {"Key": "-buildmode", "Value": "exe"},
    {"Key": "-compiler", "Value": "gc"},
    {"Key": "CGO_ENABLED", "Value": "1"},
    {"Key": "CGO_CFLAGS", "Value": ""},
    {"Key": "CGO_CPPFLAGS", "Value": ""},
    {"Key": "CGO_CXXFLAGS", "Value": ""},
    {"Key": "CGO_LDFLAGS", "Value": ""},
    {"Key": "GOARCH", "Value": "amd64"},
    {"Key": "GOOS", "Value": "linux"},
    {"Key": "GOAMD64", "Value": "v1"},
    {"Key": "vcs", "Value": "git"},
    {"Key": "vcs.revision", "Value": "a1b2c3d4..."},
    {"Key": "vcs.time", "Value": "2025-10-18T12:00:00Z"},
    {"Key": "vcs.modified", "Value": "false"}
  ]
}
```

### 字段说明

#### 顶层字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `Path` | string | 主模块路径 |
| `GoVersion` | string | Go 编译器版本 |
| `Main` | object | 主模块信息 |
| `Deps` | array | 依赖模块列表 |
| `Settings` | array | 构建设置 |

#### Main/Deps 对象

| 字段 | 类型 | 说明 |
|------|------|------|
| `Path` | string | 模块路径 |
| `Version` | string | 模块版本 |
| `Sum` | string | 模块校验和 |
| `Replace` | object | 替换信息 (如有) |

#### Settings 数组

常见的构建设置:

| Key | 说明 | 示例值 |
|-----|------|--------|
| `-buildmode` | 构建模式 | exe, pie, plugin |
| `-compiler` | 编译器 | gc, gccgo |
| `CGO_ENABLED` | CGO 是否启用 | 0, 1 |
| `GOARCH` | 目标架构 | amd64, arm64, 386 |
| `GOOS` | 目标操作系统 | linux, darwin, windows |
| `GOAMD64` | AMD64 版本 | v1, v2, v3, v4 |
| `vcs` | 版本控制系统 | git, hg, svn |
| `vcs.revision` | Commit Hash | a1b2c3d4... |
| `vcs.time` | Commit 时间 | RFC3339 格式 |
| `vcs.modified` | 是否有未提交修改 | true, false |

---

## 应用场景

### 场景 1: 依赖版本审计

**需求**: 追踪生产环境二进制使用的依赖版本

```bash
# 提取所有依赖
$ go version -m -json ./prod-app | jq '.Deps[] | {path: .Path, version: .Version}'

# 输出:
{
  "path": "github.com/gin-gonic/gin",
  "version": "v1.9.1"
}
{
  "path": "github.com/go-sql-driver/mysql",
  "version": "v1.7.1"
}
...
```

**用途**:

- ✅ 安全漏洞扫描
- ✅ 依赖合规检查
- ✅ 版本一致性验证

---

### 场景 2: SBOM 生成

**需求**: 生成软件物料清单 (Software Bill of Materials)

```bash
# 生成 SBOM
$ go version -m -json ./myapp | jq '{
    name: .Path,
    version: .Main.Version,
    goVersion: .GoVersion,
    dependencies: [.Deps[] | {
        name: .Path,
        version: .Version,
        checksum: .Sum
    }]
}' > sbom.json
```

**输出**:

```json
{
  "name": "example.com/myapp",
  "version": "v1.0.0",
  "goVersion": "go1.25.0",
  "dependencies": [
    {
      "name": "github.com/gin-gonic/gin",
      "version": "v1.9.1",
      "checksum": "h1:def456..."
    }
  ]
}
```

---

### 场景 3: 构建信息数据库

**需求**: 记录所有构建的二进制信息

```bash
#!/bin/bash
# record-build-info.sh

BINARY="$1"
BUILD_ID=$(date +%Y%m%d-%H%M%S)

# 提取构建信息
go version -m -json "$BINARY" | jq --arg id "$BUILD_ID" --arg bin "$BINARY" '{
    build_id: $id,
    binary: $bin,
    module: .Path,
    version: .Main.Version,
    go_version: .GoVersion,
    build_time: (now | strftime("%Y-%m-%dT%H:%M:%SZ")),
    dependencies: [.Deps[] | {path: .Path, version: .Version}],
    settings: .Settings | map({(.Key): .Value}) | add
}' >> build-history.jsonl
```

---

### 场景 4: 版本合规检查

**需求**: 确保所有服务使用允许的依赖版本

```bash
#!/bin/bash
# check-compliance.sh

BINARY="$1"

# 不允许的依赖版本
BLACKLIST='github.com/gin-gonic/gin@v1.8.0'

# 检查依赖
VIOLATIONS=$(go version -m -json "$BINARY" | jq -r --arg bl "$BLACKLIST" '
    .Deps[] | 
    select(.Path + "@" + .Version == $bl) |
    .Path + "@" + .Version
')

if [ -n "$VIOLATIONS" ]; then
    echo "❌ Compliance violation found: $VIOLATIONS"
    exit 1
else
    echo "✅ Compliance check passed"
fi
```

---

## 自动化脚本

### 脚本 1: 依赖版本报告

```bash
#!/bin/bash
# dependency-report.sh

echo "# Dependency Report"
echo "Generated: $(date)"
echo ""

for binary in ./bin/*; do
    echo "## $(basename $binary)"
    echo ""
    go version -m -json "$binary" | jq -r '
        "**Go Version:** \(.GoVersion)\n",
        "**Module:** \(.Path)@\(.Main.Version)\n",
        "**Dependencies:**\n",
        (.Deps[] | "- \(.Path)@\(.Version)")
    '
    echo ""
done
```

---

### 脚本 2: 安全漏洞扫描

```bash
#!/bin/bash
# security-scan.sh

BINARY="$1"
VULN_DB="vuln-database.json"

# 提取依赖
DEPS=$(go version -m -json "$BINARY" | jq -r '.Deps[] | "\(.Path)@\(.Version)"')

# 检查已知漏洞
echo "Scanning for known vulnerabilities..."
while read -r dep; do
    if grep -q "$dep" "$VULN_DB"; then
        echo "⚠️  Found vulnerability in $dep"
        grep "$dep" "$VULN_DB" | jq '.'
    fi
done <<< "$DEPS"

echo "Scan complete"
```

---

### 脚本 3: 构建信息对比

```bash
#!/bin/bash
# compare-builds.sh

BIN1="$1"
BIN2="$2"

echo "Comparing $BIN1 and $BIN2..."

# 提取依赖差异
diff <(go version -m -json "$BIN1" | jq -S '.Deps') \
     <(go version -m -json "$BIN2" | jq -S '.Deps')

if [ $? -eq 0 ]; then
    echo "✅ Dependencies are identical"
else
    echo "❌ Dependencies differ"
fi
```

---

## CI/CD 集成

### GitHub Actions

```yaml
# .github/workflows/build-audit.yml
name: Build Audit

on:
  push:
    branches: [main]
  release:
    types: [created]

jobs:
  audit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.25'
      
      - name: Build
        run: go build -o myapp ./cmd/app
      
      - name: Extract build info
        run: |
          go version -m -json ./myapp > build-info.json
          cat build-info.json
      
      - name: Generate SBOM
        run: |
          go version -m -json ./myapp | jq '{
            name: .Path,
            version: .Main.Version,
            dependencies: [.Deps[] | {name: .Path, version: .Version}]
          }' > sbom.json
      
      - name: Upload artifacts
        uses: actions/upload-artifact@v3
        with:
          name: build-metadata
          path: |
            build-info.json
            sbom.json
      
      - name: Check for known vulnerabilities
        run: |
          # 这里集成漏洞扫描工具
          # 例如: govulncheck
```

---

### GitLab CI

```yaml
# .gitlab-ci.yml
build_audit:
  stage: build
  script:
    - go build -o myapp ./cmd/app
    - go version -m -json ./myapp > build-info.json
    
    # 生成 SBOM
    - |
      go version -m -json ./myapp | jq '{
        name: .Path,
        version: .Main.Version,
        dependencies: [.Deps[] | {name: .Path, version: .Version}]
      }' > sbom.json
    
    # 上传到构建服务器
    - curl -X POST -F "file=@build-info.json" https://audit.example.com/api/builds
  
  artifacts:
    paths:
      - build-info.json
      - sbom.json
    expire_in: 1 year
```

---

## 实践案例

### 案例 1: 生产环境审计

**场景**: 审计生产环境所有服务的依赖版本

```bash
#!/bin/bash
# prod-audit.sh

echo "=== Production Environment Audit ==="
echo "Date: $(date)"
echo ""

# 从所有服务器收集二进制文件
servers=(prod-1 prod-2 prod-3)

for server in "${servers[@]}"; do
    echo "## Server: $server"
    
    # SSH 到服务器,获取二进制信息
    ssh "$server" "go version -m -json /usr/local/bin/myapp" | jq '{
        server: "'$server'",
        module: .Path,
        version: .Main.Version,
        go_version: .GoVersion,
        dependencies: [.Deps[] | {path: .Path, version: .Version}]
    }' >> prod-audit.jsonl
    
    echo "✅ Audited"
done

echo ""
echo "Audit complete. Results saved to prod-audit.jsonl"
```

---

### 案例 2: 依赖更新影响分析

**场景**: 分析依赖更新对二进制大小和依赖树的影响

```bash
#!/bin/bash
# dependency-impact.sh

# 构建当前版本
go build -o app-before ./cmd/app
SIZE_BEFORE=$(stat -f%z app-before)  # macOS
# SIZE_BEFORE=$(stat -c%s app-before)  # Linux

# 保存当前依赖
go version -m -json app-before > deps-before.json

# 更新依赖
go get -u github.com/gin-gonic/gin

# 构建更新后版本
go build -o app-after ./cmd/app
SIZE_AFTER=$(stat -f%z app-after)

# 保存更新后依赖
go version -m -json app-after > deps-after.json

# 分析影响
echo "=== Dependency Update Impact ==="
echo "Binary size: $SIZE_BEFORE → $SIZE_AFTER bytes"
echo "Change: $((SIZE_AFTER - SIZE_BEFORE)) bytes"
echo ""
echo "Dependency changes:"
diff <(jq -S '.Deps' deps-before.json) <(jq -S '.Deps' deps-after.json)
```

---

## 最佳实践

### 1. 自动化构建信息记录

```yaml
# 在 Makefile 中集成
build:
    go build -o $(BINARY) ./cmd/app
    go version -m -json $(BINARY) > $(BINARY).build-info.json
    @echo "Build info saved to $(BINARY).build-info.json"
```

---

### 2. 版本标签和构建信息关联

```bash
#!/bin/bash
# release.sh

VERSION="$1"

# 构建
go build -ldflags "-X main.version=$VERSION" -o myapp-$VERSION ./cmd/app

# 提取构建信息
go version -m -json myapp-$VERSION | jq --arg v "$VERSION" '. + {release_tag: $v}' > myapp-$VERSION.json

# 创建 Git 标签
git tag -a "$VERSION" -m "Release $VERSION"

echo "✅ Release $VERSION created"
```

---

### 3. 构建信息归档

```bash
# 归档到 S3
aws s3 cp build-info.json s3://builds/$(date +%Y/%m/%d)/build-info-$(git rev-parse --short HEAD).json

# 归档到本地数据库
sqlite3 builds.db "INSERT INTO builds (date, commit, info) VALUES (datetime('now'), '$(git rev-parse HEAD)', '$(cat build-info.json)');"
```

---

## 常见问题

### Q1: 旧版本 Go 编译的二进制可以用吗?

**A**: ✅ 可以,但信息可能不完整

- Go 1.18+ 嵌入完整构建信息
- Go 1.13-1.17 嵌入部分信息
- Go <1.13 可能没有构建信息

---

### Q2: 如何处理批量二进制文件?

**A**: 使用 `jq -s` (slurp mode)

```bash
# 合并多个 JSON
for bin in ./bin/*; do
    go version -m -json "$bin"
done | jq -s '.' > all-builds.json
```

---

### Q3: 输出太大怎么办?

**A**: 选择性提取字段

```bash
# 只提取关键信息
go version -m -json ./myapp | jq '{
    module: .Path,
    version: .Main.Version,
    go_version: .GoVersion,
    dependencies: [.Deps[] | .Path]
}'
```

---

## 参考资料

### 官方文档

- 📘 [Go 1.25 Release Notes](https://go.dev/doc/go1.25)
- 📘 [go version command](https://pkg.go.dev/cmd/go#hdr-Print_Go_version)
- 📘 [Build Info](https://pkg.go.dev/runtime/debug#BuildInfo)

### 相关章节

- 🔗 [Go 1.25 工具链增强](./README.md)
- 🔗 [CI/CD 最佳实践](../../最佳实践/CI-CD.md)

---

**编写者**: AI Assistant  
**审核者**: [待审核]  
**最后更新**: 2025年10月18日

---

<p align="center">
  <b>📊 使用 go version -m -json 实现构建信息自动化! 🤖</b>
</p>

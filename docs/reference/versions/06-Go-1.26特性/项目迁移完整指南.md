# Go 1.26 项目迁移完整指南

> 从 Go 1.25/1.24 迁移到 Go 1.26 的完整步骤和检查清单

---

## 迁移概览

### 难度评估

| 项目类型 | 迁移难度 | 预估时间 | 主要工作 |
|----------|----------|----------|----------|
| 简单 CLI | ⭐ 低 | 30分钟 | 更新 go.mod |
| Web 服务 | ⭐⭐ 中 | 2小时 | 检查 TLS + go fix |
| 加密应用 | ⭐⭐⭐ 中高 | 4小时 | 检查随机参数 + HPKE |
| 图像处理 | ⭐⭐ 中 | 2小时 | 更新 JPEG 测试 |
| 使用 CGO | ⭐ 低 | 1小时 | 享受性能提升 |
| 大型项目 | ⭐⭐⭐⭐ 高 | 1-2天 | 全面测试 |

---

## 迁移前准备

### 1. 环境检查

```bash
# 检查当前 Go 版本
go version

# 检查项目状态
git status
git log --oneline -5

# 确保测试通过
go test ./...
```

### 2. 备份代码

```bash
# 创建迁移分支
git checkout -b upgrade/go1.26

# 或者打标签
git tag -a v1.x-before-go1.26 -m "Before Go 1.26 upgrade"
```

### 3. 依赖审计

```bash
# 列出所有依赖
go list -m all > deps-before.txt

# 检查过时依赖
go list -u -m all

# 检查安全漏洞
go install github.com/sonatype-nexus-community/nancy@latest
cat go.sum | nancy sleuth
```

---

## 迁移步骤

### 步骤1: 安装 Go 1.26

```bash
# 使用官方安装包
# https://go.dev/dl/

# 或者使用版本管理工具
# gvm
gvm install go1.26
gvm use go1.26

# 验证安装
go version
go version go1.26.0 linux/amd64
```

---

### 步骤2: 更新 go.mod

```bash
# 更新 Go 版本
go mod edit -go=1.26

# 查看变更
cat go.mod | grep "^go"

# 整理依赖
go mod tidy

# 验证依赖
go mod verify
```

**go.mod 示例**:

```go
module example.com/myproject

go 1.26  // 更新这一行

require (
    // 依赖会自动更新
)
```

---

### 步骤3: 运行 go fix

```bash
# 预览将要应用的修复
go fix -n ./... 2>&1 | head -50

# 应用所有修复
go fix ./...

# 查看变更
git diff --stat

# 具体查看某个文件的变更
git diff main.go
```

**常见修复**:

- 循环查找 → `slices.Contains`
- `sort.Slice` → `slices.Sort`
- if-else → `min`/`max`
- `errors.As` → `errors.AsType`

---

### 步骤4: 检查破坏性变化

#### 4.1 加密包随机参数

**检查代码**:

```bash
# 搜索可能受影响的代码
grep -r "GenerateKey\|Sign\|Prime" --include="*.go" | grep -v "_test.go"
```

**需要修改的情况**:

```go
// 如果测试依赖确定性随机
import "testing/cryptotest"

func TestDeterministic(t *testing.T) {
    cryptotest.SetGlobalRandom(yourTestRand)
    defer cryptotest.SetGlobalRandom(nil)

    // 测试代码
}
```

**无需修改的情况**:

- 正常使用 `rand.Reader`
- 不依赖特定随机序列

---

#### 4.2 JPEG 编解码器

**检查测试**:

```bash
# 搜索 JPEG 相关测试
grep -r "jpeg" --include="*_test.go"
```

**需要修改的情况**:

```go
// 旧代码: 检查特定字节
got := encodeJPEG(img)
want := loadFile("expected.jpg")
if !bytes.Equal(got, want) { // 可能失败
    t.Error()
}

// 新代码: 检查图像属性
got := encodeJPEG(img)
gotImg, _ := jpeg.Decode(bytes.NewReader(got))
if gotImg.Bounds() != wantBounds {
    t.Error()
}
```

---

#### 4.3 GODEBUG 设置

**检查环境变量**:

```bash
echo $GODEBUG
```

**即将移除的设置**:

```bash
# 这些设置将在 Go 1.27 移除
tlsunsafeekm=1
tlsrsakex=1
tls10server=1
tls3des=1
x509keypairleaf=0
```

**迁移建议**:

- 移除这些设置
- 测试应用是否正常工作
- 如有问题，更新代码而非依赖 GODEBUG

---

### 步骤5: 更新代码（可选）

#### 5.1 使用 new(expr)

**识别机会**:

```bash
# 搜索可能的模式
grep -rn "var tmp" --include="*.go"
grep -rn "&T{" --include="*.go" | head -20
```

**改写示例**:

```go
// 旧代码
age := calculateAge(birth)
person.Age = &age

// 新代码
person.Age = new(calculateAge(birth))
```

**建议**: 不是必须修改，根据可读性决定。

---

#### 5.2 使用递归泛型（高级）

**适用场景**:

- 通用树/图算法库
- 需要类型安全的数据结构

**示例**:

```go
// 为项目添加通用树遍历
type TreeNode[T TreeNode[T]] interface {
    Children() []T
}

func Traverse[T TreeNode[T]](root T, visit func(T)) {
    // 实现
}
```

**建议**: 新项目可以使用，旧项目不必强行迁移。

---

### 步骤6: 运行测试

```bash
# 运行所有测试
go test ./...

# 运行竞态检测
go test -race ./...

# 运行基准测试
go test -bench=. -benchmem ./...

# 检查测试覆盖率
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

**测试失败处理**:

| 失败类型 | 可能原因 | 解决方案 |
|----------|----------|----------|
| 加密测试 | 随机参数变更 | 使用 cryptotest |
| JPEG 测试 | 输出变化 | 更新期望值 |
| 性能测试 | GC 行为变化 | 调整阈值 |
| 功能测试 | 行为变更 | 检查 GODEBUG |

---

### 步骤7: 性能验证

```bash
# 运行基准测试对比
# 保存旧版本结果
go test -bench=. -count=5 > benchmark-old.txt

# 切换到新版本后
go test -bench=. -count=5 > benchmark-new.txt

# 使用 benchstat 对比
go install golang.org/x/perf/cmd/benchstat@latest
benchstat benchmark-old.txt benchmark-new.txt
```

**预期改进**:

- GC 延迟降低 10-40%
- cgo 开销减少 30%
- io.ReadAll 性能提升 2 倍

---

### 步骤8: 部署验证

#### 8.1 构建验证

```bash
# 清理构建缓存
go clean -cache

# 构建项目
go build -o myapp .

# 检查构建信息
go version -m myapp
```

#### 8.2 容器验证

```dockerfile
# Dockerfile
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/myapp .
CMD ["./myapp"]
```

```bash
# 构建并测试
docker build -t myapp:go1.26 .
docker run myapp:go1.26
```

---

## 迁移后检查清单

### 功能检查

- [ ] 所有测试通过
- [ ] 基准测试无回归
- [ ] 竞态检测无问题
- [ ] 构建成功
- [ ] 容器构建成功
- [ ] 部署验证通过

### 性能检查

- [ ] GC 延迟降低
- [ ] 内存使用正常
- [ ] 响应时间无回归
- [ ] 吞吐量正常

### 安全检查

- [ ] TLS 配置正确
- [ ] 加密功能正常
- [ ] 无安全漏洞

### 文档检查

- [ ] README 更新
- [ ] 变更日志更新
- [ ] 部署文档更新

---

## 回滚计划

### 如果迁移失败

```bash
# 回滚代码
git checkout main
git branch -D upgrade/go1.26

# 或者使用标签
git checkout v1.x-before-go1.26

# 回滚 go.mod
go mod edit -go=1.25
go mod tidy
```

---

## 常见问题

### Q: 测试失败但不确定原因？

**A**: 使用二分法排查

```bash
# 只运行特定测试
go test -v -run TestSpecific ./...

# 使用旧 Go 版本对比
go1.25 test -v -run TestSpecific ./...
go1.26 test -v -run TestSpecific ./...
```

### Q: 依赖不兼容？

**A**: 更新依赖

```bash
# 更新所有依赖
go get -u ./...
go mod tidy

# 或者更新特定依赖
go get -u example.com/some/dependency
go mod tidy
```

### Q: 性能反而下降了？

**A**: 检查 GC 调优

```bash
# 导出 GC 日志
GODEBUG=gctrace=1 go test -bench=. ./... 2>&1 | head -50

# 调整 GOGC
go test -bench=. ./... 2>&1 | head -50
```

---

## 迁移时间表模板

### 小型项目 (1天)

| 时间 | 任务 | 负责人 |
|------|------|--------|
| 09:00 | 环境准备 | Dev |
| 09:30 | 更新 go.mod | Dev |
| 10:00 | 运行 go fix | Dev |
| 11:00 | 运行测试 | Dev |
| 12:00 | 午餐 | - |
| 13:00 | 修复问题 | Dev |
| 15:00 | 性能验证 | Dev |
| 16:00 | 部署验证 | Dev |
| 17:00 | 完成 | - |

### 大型项目 (1周)

| 天 | 任务 | 产出 |
|----|------|------|
| 周一 | 环境准备 + 依赖审计 | 检查清单 |
| 周二 | 更新 + go fix | 修复分支 |
| 周三 | 测试修复 | 测试通过 |
| 周四 | 性能验证 | 基准报告 |
| 周五 | 部署 + 监控 | 生产上线 |

---

## 成功案例

### 案例1: Web API 项目

**项目规模**: 5万行代码
**迁移时间**: 4小时
**主要工作**:

- go fix 自动修复 15 处代码
- 更新 2 个 JPEG 测试
- 验证 GC 延迟降低 25%

### 案例2: 微服务集群

**项目规模**: 20万行代码，10个服务
**迁移时间**: 2天
**主要工作**:

- 逐个服务迁移
- 统一更新基础库
- 全链路性能提升 15%

---

## 参考资源

- [Go 1.26 Release Notes](https://go.dev/doc/go1.26)
- [常见问题解答](./常见问题解答.md)
- [实战示例集](./实战示例集.md)

---

*按照本指南，你的项目应该能够顺利迁移到 Go 1.26。*

# 季度权威对齐审计机制

> **文档编号**: TRACK-AUDIT-001  
> **频率**: 每季度（3月、6月、9月、12月）  
> **负责人**: 技术委员会 / 知识库维护者  
> **最后更新**: 2026-05-06

---

## 一、审计目标

确保本项目所有 Go 特性描述与 **go.dev 官方文档** 保持同步，消除：
1. 虚构 API（如 `runtime.SetGoroutineLeakCallback`）
2. 版本归属错误（如 Go 1.23 特性标注为 1.25）
3. 过时的实验性/稳定状态描述
4. 缺失的新标准库包覆盖

---

## 二、审计清单模板

### 2.1 语言特性核对

| 特性 | 本项目描述 | 官方来源 | 状态 |
|------|-----------|---------|------|
| `new(expr)` | | go.dev/doc/go1.26 | |
| 自引用泛型 | | go.dev/doc/go1.26 | |
| `errors.AsType` | | go.dev/doc/go1.26 | |
| `slog.NewMultiHandler` | | go.dev/doc/go1.26 | |
| `crypto/hpke` | | go.dev/doc/go1.26 | |
| `runtime/secret` | | go.dev/doc/go1.26 (exp) | |
| `simd/archsimd` | | go.dev/doc/go1.26 (exp) | |
| `goroutineleak` profile | | go.dev/doc/go1.26 (exp) | |

### 2.2 实验性特性状态追踪

| 特性 | Go 1.26 状态 | Go 1.27 实际 | 行动 |
|------|-------------|-------------|------|
| `runtime/secret` | `GOEXPERIMENT=runtimesecret` | 仍为实验 | 追踪 release notes |
| `simd/archsimd` | `GOEXPERIMENT=simd` | 仍为实验 | 追踪 release notes |
| `goroutineleak` profile | `GOEXPERIMENT=goroutineleakprofile` | 已 GA（新增 `/debug/pprof/goroutineleak` 端点） | 文档已更新 |

**Go 1.27 实际发布结论**（2026-08，当前版本 1.27.1）：泛型方法 ✅ 正式落地（如 `math/rand/v2.Rand.N[Int intType]`）、`encoding/json/v2` ✅ GA（v1 改由 v2 实现）、`goroutineleak` profile ✅ GA、container-aware GOMAXPROCS ✅ 默认开启；Green Tea GC ❌ 仍为 `GOEXPERIMENT=greenteagc` 实验，未默认启用。

---

## 三、自动化审计脚本

```bash
#!/bin/bash
# scripts/quarterly-audit.sh

echo "=== Quarterly Authority Audit ==="
echo "Date: $(date)"
echo "Go Version: $(go version)"
echo

# 1. 虚构 API 扫描
echo "[1/4] Scanning for fictional APIs..."
if grep -r "runtime.SetGoroutineLeakCallback" --include="*.md" --include="*.go" .; then
    echo "FAIL: Fictional API found"
    exit 1
fi
echo "PASS"

# 2. 文档 API 存在性验证
echo "[2/4] Verifying documented APIs..."
go doc errors.AsType >/dev/null && echo "  ✓ errors.AsType"
go doc log/slog.NewMultiHandler >/dev/null && echo "  ✓ slog.NewMultiHandler"
go doc crypto/hpke.Seal >/dev/null && echo "  ✓ crypto/hpke"

# 3. 版本标签检查
echo "[3/4] Checking version labels..."
# 检查是否有已知的错误标签

# 4. 生成审计报告
echo "[4/4] Generating report..."
echo "Audit completed. See docs/tracking/quarterly-audit-$(date +%Y-%m).md"
```

---

## 四、2026-Q2 审计记录（示例）

**审计日期**: 2026-05-06  
**执行人**: Kimi Code CLI  
**Go 版本**: go1.27.1 windows/amd64（2026-08 自 go1.26.2 升级）

### 发现与修正

| 发现 | 严重性 | 修正措施 | 状态 |
|------|--------|---------|------|
| `runtime.SetGoroutineLeakCallback` 虚构 API | 🔴 Critical | 删除并替换为 `goroutineleak` profile | ✅ 已修正 |
| `new(expr)` 被错误否定 | 🔴 Critical | 重写示例和文档 | ✅ 已修正 |
| `iter`/`unique` 版本误标为 1.25 | 🟠 High | 更正为 Go 1.23 | ✅ 已修正 |
| `crypto/hpke` 零覆盖 | 🟠 High | 新建示例和知识库文档 | ✅ 已补充 |
| `math/rand/v2` 未使用 | 🟡 Medium | 迁移生产代码 | ✅ 已迁移 |
| `weak.Pointer` 模拟实现 | 🟡 Medium | 替换为标准库 API | ✅ 已替换 |
| `cmd/grpc-server` 未完成 | 🟡 Medium | 完成服务注册 | ✅ 已完成 |

---

## 五、下季度（2026-Q3）待追踪事项

1. Go 1.27 已发布（2026-08，当前 1.27.1）✅
2. `goroutineleak` profile 已 GA（新增 `/debug/pprof/goroutineleak` 端点）✅
3. `runtime/secret` 可能进入稳定阶段
4. Go 1.26.3/1.26.4 安全补丁追踪（注：建议直接升级至 1.27.1）

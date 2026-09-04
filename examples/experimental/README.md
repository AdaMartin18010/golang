# 实验性特性沙箱 (Experimental Features Sandbox)

> **分类**: examples/experimental/  
> **Go 版本**: 1.26+  
> **最后更新**: 2026-05-06

---

## 使用说明

本目录统一收纳需要 `GOEXPERIMENT` 开启的 Go 特性。

与 `examples/go126-features/` 中的**稳定特性**不同，实验性特性：
- 需要特定 `GOEXPERIMENT` 标志编译
- API 可能在后续版本中变化
- 不建议在生产环境直接使用

---

## 目录结构

```
examples/experimental/
├── README.md              # 本文件
├── simd-archsimd/         # SIMD 向量运算 (GOEXPERIMENT=simd)
├── goroutineleak-profile/ # Goroutine 泄漏检测 (GOEXPERIMENT=goroutineleakprofile)
└── runtime-secret/        # 敏感数据内存擦除 (GOEXPERIMENT=runtimesecret)
```

---

## 实验性特性清单

| 特性 | 实验标志 | Go 版本 | 预计稳定 |
|------|---------|---------|---------|
| `simd/archsimd` | `GOEXPERIMENT=simd` | 1.26 | 待定 |
| `goroutineleak` profile | `GOEXPERIMENT=goroutineleakprofile` | 1.26 | **Go 1.27 已 GA（/debug/pprof/goroutineleak）** |
| `runtime/secret` | `GOEXPERIMENT=runtimesecret` | 1.26 | 待定 |

---

## 运行示例

```bash
# SIMD 示例
GOEXPERIMENT=simd go run ./experimental/simd-archsimd/

# Goroutine leak profile
GOEXPERIMENT=goroutineleakprofile go run ./experimental/goroutineleak-profile/

# Runtime secret
GOEXPERIMENT=runtimesecret go run ./experimental/runtime-secret/
```

---

## 追踪状态

每季度检查实验性特性的稳定化进度，见 `docs/tracking/quarterly-authority-audit.md`。

# Go 1.26.2 全面现代化后续计划

> 状态：核心现代化 100% 完成，以下任务为可选增强项
> 日期：2026-05-06
> 当前 Go 版本：go1.27.1 windows/amd64（2026-08 自 go1.26.2 升级，Go 1.27 升级任务 ✅ 已完成）

---

## 已完成的基线（无需再投入）

- [x] 虚构 API 清理（`runtime.SetGoroutineLeakCallback`）
- [x] 版本归属校正（`iter`/`unique` → Go 1.23，`slog.NewMultiHandler` → Go 1.26 等）
- [x] `math/rand` → `math/rand/v2` 迁移
- [x] `weak.Make[T]` 正确实现
- [x] `pkg/registry` 死锁修复 + 测试修复
- [x] `pkg/utils/crypto` AES key size 修复
- [x] `pkg/utils/id` `SequentialID` 竞态修复
- [x] `pkg/logger` sample rate 逻辑修复
- [x] `cmd/grpc-server` 服务注册完成
- [x] 所有 `go.mod` 统一为 `go 1.26`
- [x] 所有生产 Dockerfile 更新为 `golang:1.26.2-alpine`
- [x] 所有 README 版本徽章更新
- [x] `gofmt` 全面格式化
- [x] CI workflow 版本矩阵精简为 `1.26.x`
- [x] `scripts/quality-check.sh` 与季度审计机制建立

---

## 候选后续任务（按优先级分组）

### P1 — CI/CD 与质量门禁（建议优先）

| # | 任务 | 影响 | 预估工作量 |
|---|------|------|----------|
| 1.1 | **CI 新增 `gofmt --check`** | 防止未来未格式化代码入库 | 1 个 workflow 文件 |
| 1.2 | **CI 新增 `go vet` gate** | 静态分析强制通过 | 1 个 workflow 文件 |
| 1.3 | **CI 集成 `scripts/quality-check.sh`** | 自动化虚构 API 与版本扫描 | 1 个 workflow 文件 |
| 1.4 | **CI 新增 `govulncheck ./...`** | 安全漏洞扫描 | 1 个 job，可能需缓存 |

### P2 — Go 1.26 新特性深度落地（可选）

| # | 任务 | 影响 | 风险 | 预估工作量 |
|---|------|------|------|----------|
| 2.1 | **`omitzero` struct tag 评估** | JSON 序列化更精确（避免 `omitempty` 误删 0/false） | ent 生成代码不可改；手写的 DTO/响应结构可改 | 审计 + 局部替换 |
| 2.2 | **`new(expr)` 语法识别与应用** | 简化 `*T` 创建（如 `new(time.Now().Year())`） | 低：纯语法糖，不改变语义 | 扫描后局部重构 |
| 2.3 | **`errors.AsType` 评估** | 类型断言更简洁 | 需 Go 1.26+，低向后兼容风险 | 扫描后局部重构 |
| 2.4 | **`runtime/secret` 实验示例扩展** | 安全敏感数据内存保护演示 | 需 `GOEXPERIMENT=secret` | 仅示例/文档 |

### P3 — 性能与可观测性（中长期）

| # | 任务 | 影响 | 预估工作量 |
|---|------|------|----------|
| 3.1 | **PGO（Profile-Guided Optimization）集成** | 5-10% 性能提升 | 添加收集 profile 的 Makefile 目标 + CI 支持 |
| 3.2 | **`testing.B.Loop` 基准测试现代化** | 更准确的 benchmark（Go 1.24+） | 重写关键 benchmark（pkg/loadbalancer, pkg/memory 等） |
| 3.3 | **Green Tea GC 调优文档** | 利用 Go 1.26 默认 GC 优化 | 更新 ops 文档 |

### P4 — 知识库与文档（低优先）

| # | 任务 | 说明 | 预估工作量 |
|---|------|------|----------|
| 4.1 | **历史模板文档 Dockerfile 更新** | `go-knowledge-base/` 中 200+ 处 `golang:1.21-alpine` 示例属于历史模板/报告，如作为未来项目模板复用则应更新 | 批量替换，但需区分历史记录与模板 |
| 4.2 | **CHANGELOG 正式化** | 将多轮现代化改动写入项目根目录 CHANGELOG.md | 中等 |
| 4.3 | **安全公告 CVE-2026-27143/27144 内部培训** | 基于已编写的 `docs/go126-security-advisory-2026-05.md` 制作团队通知 | 低 |

---

## 推荐执行顺序

```
Phase A（本周）: P1 全部 → CI 质量门禁上线
Phase B（下周）: P2.1 + P2.2 → Go 1.26 语法糖落地
Phase C（下月）: P3.1 + P3.2 → 性能优化
Phase D（按需）: P4.1 → 知识库模板更新
```

---

## 决策点（需您确认）

1. **是否执行 P1（CI 门禁）？**  建议执行，防止退化。
2. **P2 中是否优先尝试 `omitzero`？** 需先审计 handwritten struct（非 ent 生成）。
3. **P4.1 是否批量更新知识库历史模板？** 277 处引用，但大部分是历史完成报告，修改价值有限。
4. **是否现在提交当前所有改动并打 tag（如 `v1.26.2-modernized`）？** 便于回滚与追踪。

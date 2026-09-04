# 可复用 Prompt：Rust patch release 响应

> **用途**：针对 Rust patch release（如 1.98.1、1.97.1）新建或更新 `concept/07_future/00_version_tracking/rust_1_XX_Y.md`。
> **配套文件**：[`.kimi/templates/kimi_project_requirements.md`](../templates/kimi_project_requirements.md) §9。

---

## Prompt 模板

```text
请在 E:/_src/rust-lang 仓库中响应 Rust {{patch_version}} patch release。

**已知信息**：
- 前一个稳定版本：{{previous_version}}
- 发布日期：{{release_date}}
- 修复/变更摘要：{{summary}}
- 官方 Release Notes：{{release_notes_url}}
- 相关 issue/PR：{{issue_url}}

要求：
1. 若 `concept/07_future/00_version_tracking/rust_1_XX_Y.md` 已存在，则补全；否则新建。
2. 按 `.kimi/templates/kimi_project_requirements.md` §9 版本补丁页格式书写头部（EN、Summary、Rust 版本、Bloom 层级、权威来源、前置/后置概念）。
3. 正文必须包含：
   - 补丁要点（发布日期、修复类型、影响范围、是否建议升级）
   - 对称差分析（A∩B、B\A、A\B），明确说明 1.98.1 与 1.98.0 在语法/语义/API/平台支持上的差异
   - 技术细节与触发路径
   - 迁移建议（rustup update stable、MSRV 是否变更、CI/安全关键项目注意事项）
   - 最小风险代码模式示例
   - 反命题与边界分析表（至少 6 条）
   - Mermaid mindmap
   - P0/P1/P2 References
4. MSRV 保持 {{previous_msrv}}；`Cargo.toml` 的 `rust-version` 不提升。
5. 在受影响的 concept/ 权威页注入补丁提示并建立双向链接（典型页：Trait、Unsafe、Memory Model、Toolchain、Async、Pin）。
6. 运行：
   python scripts/check_msrv_consistency.py --strict
   python scripts/kb_auditor.py --link-check
   python scripts/check_version_semantic_injection.py --strict
7. 不要声明“已完成”，除非上述命令 exit 0。
```

---

## 示例填充

```text
**patch_version**: 1.98.1
**previous_version**: 1.98.0
**release_date**: 2026-09-03
**summary**: rustc vtable 生成错误编译修复：复杂 trait bound 下方法槽位可能被置零，导致 dyn Trait 调用时 segfault。
**release_notes_url**: https://doc.rust-lang.org/stable/releases.html
**issue_url**: https://github.com/rust-lang/rust/issues/161441
**previous_msrv**: 1.98.0
```

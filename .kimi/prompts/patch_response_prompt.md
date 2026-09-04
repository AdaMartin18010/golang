# 可复用 Prompt：Go patch release 响应

> **用途**：针对 Go patch release（如 1.27.1、1.27.2）更新版本跟踪页与受影响权威页。
> **本项目版本页位置**：`go-knowledge-base/02-Language-Design/`（LD-026/027/029/030/035/036/037）、`go-knowledge-base/02-Language-Design/03-Evolution/`、`view/formal/Go/`、`docs/go1xx-*`。
> **配套文件**：[`.kimi/templates/kimi_project_requirements.md`](../templates/kimi_project_requirements.md)「红线」节。

---

## Prompt 模板

```text
请在 E:/_src/golang 仓库中响应 Go {{patch_version}} patch release。

**已知信息**：
- 前一个稳定版本：{{previous_version}}
- 发布日期：{{release_date}}
- 修复/变更摘要：{{summary}}
- 官方 Release Notes：{{release_notes_url}}
- 相关 issue/CL：{{issue_url}}

要求：
1. 定位或新建版本跟踪页（优先复用已有 03-Evolution 或 LD-02x/03x 页面，不新建重复权威页）。
2. 按 AGENTS.md §3 头部模板补全元数据（维度/级别/标签/Go 版本/Bloom/前置/后置/定理链）。
3. 正文必须包含：
   - 补丁要点（发布日期、修复类型、影响范围、是否建议升级）
   - 对称差分析（A∩B、B\A、A\B），明确 {{patch_version}} 与 {{previous_version}} 在语法/语义/标准库/平台支持/工具链上的差异
   - 技术细节与触发路径
   - 迁移建议（go.mod 的 go 指令是否需变更、CI 镜像、安全关键项目注意事项）
   - 最小风险代码模式示例
   - 反命题与边界分析表（至少 6 条）
   - Mermaid mindmap
   - P0/P1/P2 References
4. go.mod / go.work 的 go 指令保持 {{previous_go_directive}}；不提升语言版本（除非 release 明确变更）。
5. 在受影响的权威页注入补丁提示并建立双向链接（典型页：内存模型、GC、调度器、编译器、标准库对应包页）。
6. 验证：
   python scripts/tmp/rescan_deadlinks.py   # 死链 0
   cd examples/go127-features && GOWORK=off go vet ./...   # 示例模块未受影响也应通过
7. 不要声明"已完成"，除非上述检查通过。
```

---

## 示例填充

```text
**patch_version**: 1.27.1
**previous_version**: 1.27.0
**release_date**: 2026-09-0X
**summary**: 编译器在泛型方法实例化 + 内联组合场景下错误折叠边界检查（示意）。
**release_notes_url**: https://go.dev/doc/devel/release#go1.27.1
**issue_url**: https://github.com/golang/go/issues/XXXXX
**previous_go_directive**: go 1.27
```

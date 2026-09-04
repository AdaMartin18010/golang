# 可复用 Prompt：形式化与定理链补全

> **用途**：为 FT/LD 维度权威页（L4-L5）补充或修正定理链、形式化语义、不变式与学术来源对齐。
> **配套文件**：[`.kimi/templates/kimi_project_requirements.md`](../templates/kimi_project_requirements.md)「定理链与推理」节。

---

## Prompt 模板

```text
请在 E:/_src/golang 仓库中补全/修正以下权威页的形式化与定理链。

**目标页面**：{{target_path}}
**主题**：{{topic}}
**当前 Bloom 层级**：{{bloom_level}}
**相关形式化来源**：{{formal_sources}}  <!-- 如 Go 内存模型、Featherweight Go 论文、CSP 原始论文、Raft/Paxos 论文 -->

要求：
1. 先读取目标页与所有前置/后置概念页，确保定理链不自相矛盾、不循环。
2. 在头部添加或修正 `定理链` 字段，格式如：
   前提 → 操作 → 结论 / 不变式
3. 在正文中为每个推理步骤写出至少一句话解释，并用 `⟹` / `⟸` 标记推理方向。
4. 补充以下至少一种形式化内容：
   - 小步/大步操作语义规则（ inference rules ）
   - 不变式表格（happens-before 关系、类型保持/进展等）
   - TLA+ 规约或伪代码
   - 竞态/未定义行为触发条件表
5. 反命题与边界分析节必须引用至少一个推理步骤编号或不变式。
6. 所有形式化声明必须对齐 P0/P1 权威来源（Go Spec、Go 内存模型、原始论文）。
7. 不要引入无法验证的"伪定理"；引用论文给出链接。
8. 运行：
   python scripts/tmp/rescan_deadlinks.py
   cd <涉及模块> && GOWORK=off go vet ./...   # 若含代码块
9. 不要声明"已完成"，除非命令 exit 0。
```

---

## 示例填充

```text
**target_path**: go-knowledge-base/02-Language-Design/LD-001-Go-Memory-Model-Formal.md
**topic**: happens-before 的传递性与 channel 同步
**bloom_level**: L4
**formal_sources**: Go Memory Model (go.dev/ref/mem), Hoare CSP (1978), Featherweight Go (Griesemer et al. 2020)
```

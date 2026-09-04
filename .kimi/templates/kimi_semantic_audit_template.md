# 语义审计报告：{审计主题}

> **EN**: Semantic Audit Report — {English Topic}
> **Summary**: Quarterly/monthly audit of semantic consistency, drift, and coverage for `{审计主题}` against canonical sources.
> **Scope**: `E:/_src/rust-lang`
> **Audit date**: YYYY-MM-DD
> **Auditor**: Kimi / Human maintainer

---

## 1. 审计范围与方法

### 1.1 范围

- 受审 `concept/` 页：
- 受审 `docs/` / `content/` 页：
- 对照权威来源：The Rust Reference / The Rustonomicon / TRPL / RFCs / RustBelt / Tree Borrows / ...

### 1.2 方法

1. 抽取受审页中的**定义句**、**不变式**、**定理链**。
2. 与权威来源逐条比对，标记：
   - ✅ 一致
   - ⚠️ 表述差异但语义等价
   - ❌ 语义漂移或事实错误
   - ❓ 待补充证据
3. 运行相关质量门：

   ```bash
   python scripts/concept_consistency_auditor.py --strict
   python scripts/check_kg_relation_precision.py --strict
   python scripts/check_cross_domain_coverage.py --strict
   python scripts/semantic_health.py --strict
   ```

---

## 2. 发现清单

| # | 位置 | 发现类型 | 具体问题 | 建议动作 | 优先级 |
|---:|---|---|---|---|---|
| 1 | `concept/.../xx.md` §3.2 | 语义漂移 | 定义与 Reference 第 X 章不一致 | 按权威来源修正或加注 scope | P0 |
| 2 | `concept/.../yy.md` §4 | 缺失链接 | 未链接到相关交叉域页 | 添加双向链接 | P1 |
| 3 | `concept/.../zz.md` 元数据 | 元数据错误 | Bloom 层级与正文不匹配 | 修正为元数据或正文 | P2 |

---

## 3. 跨文件一致性

| 术语 | 权威页定义 | 其他出现位置 | 状态 |
|---|---|---|---|
| `{term}` | `concept/.../canonical.md` "..." | `concept/.../other.md` "..." | ✅/⚠️/❌ |

---

## 4. KG / 拓扑检查

- 核心 50 实体周边 `generic_ratio`：____%
- 新增关系是否使用具体谓词：是 / 否
- Atlas 符号 ⟹/⊣/⟺ 是否映射到 KG 谓词：是 / 否
- 决策树节点是否关联 rustc error code / KG 实体：是 / 否

---

## 5. 修复与后续动作

| # | 动作 | 责任人 | 验收标准 | 截止日期 |
|---:|---|---|---|---|
| 1 | 修正 `xx.md` 定义漂移 | Kimi | `concept_consistency_auditor.py` exit 0 | YYYY-MM-DD |
| 2 | 补充交叉域双向链接 | Kimi | `kb_auditor.py --link-check` 0 死链 | YYYY-MM-DD |
| 3 | 重新生成 KG | Kimi | `check_kg_relation_precision.py --strict` exit 0 | YYYY-MM-DD |

---

## 6. 证据与基线

- `semantic_health.py` 输出：
- `check_cross_domain_coverage.py` 输出：
- 相关 CI 链接：

---

> **完成声明纪律**：本报告的所有“已修复”结论必须附质量门 exit 0 截图或命令输出；否则只能标记为“待验证”。

# 语义审计报告：{审计主题}

> **EN**: Semantic Audit Report — {English Topic}
> **Summary**: Quarterly/monthly audit of semantic consistency, drift, and coverage for `{审计主题}` against canonical sources.
> **Scope**: `E:/_src/golang`
> **Audit date**: YYYY-MM-DD
> **Auditor**: Kimi / Human maintainer

---

## 1. 审计范围与方法

### 1.1 范围

- 受审五维权威页（`go-knowledge-base/0X-维度/{FT|LD|EC|TS|AD}-NNN-*.md`）：
- 受审 `docs/` / `view/` 页（须回链五维权威页）：
- 对照权威来源：Go 语言规范（go.dev/ref/spec）/ Effective Go / Go Memory Model / Go Blog / Proposals（go.googlesource.com/proposal）/ pkg.go.dev / ...

### 1.2 方法

1. 抽取受审页中的**定义句**、**不变式**、**定理链**。
2. 与权威来源逐条比对，标记：
   - ✅ 一致
   - ⚠️ 表述差异但语义等价
   - ❌ 语义漂移或事实错误
   - ❓ 待补充证据
3. 运行相关质量门：

   ```bash
   # 六件套覆盖与结构检查
   python scripts/tmp/verify_sixpiece.py

   # canonical 唯一性 / 编号唯一性 / 双向链接
   python scripts/tmp/dup_canon_scan.py
   python scripts/tmp/dup_id_scan.py
   python scripts/tmp/check_bidir.py

   # 死链复扫（跳过 archive/、代码围栏，排除 http(s) 外链）
   python scripts/tmp/rescan_deadlinks.py

   # 代码块实测：提取权威页 go 块并编译/测试
   python scripts/tmp/extract_go_blocks.py
   ```

---

## 2. 发现清单

| # | 位置 | 发现类型 | 具体问题 | 建议动作 | 优先级 |
|---:|---|---|---|---|---|
| 1 | `go-knowledge-base/0X-维度/xx.md` §3.2 | 语义漂移 | 定义与 Go Spec 第 X 节不一致 | 按权威来源修正或加注 scope | P0 |
| 2 | `go-knowledge-base/0X-维度/yy.md` §4 | 缺失链接 | 未链接到相关交叉域页 | 添加双向链接 | P1 |
| 3 | `go-knowledge-base/0X-维度/zz.md` 元数据 | 元数据错误 | Bloom 层级与正文不匹配 | 修正元数据或正文 | P2 |

---

## 3. 跨文件一致性

| 术语 | 权威页定义 | 其他出现位置 | 状态 |
|---|---|---|---|
| `{term}` | `go-knowledge-base/0X-维度/canonical.md` "..." | `docs/...other.md` "..." | ✅/⚠️/❌ |

---

## 4. KG / 拓扑检查

- KG 关系是否全部使用具体语义谓词（`dependsOn`/`entails`/`mutexWith`/`refines`/`equivalentTo`/`counterExample`）：是 / 否
- 定理链符号 ⟹ / ⟸ 是否映射到 KG 谓词：是 / 否
- 决策树节点（如启用）是否关联「Go 编译失败类别」/ KG 实体：是 / 否
- 权威页前置/后置概念链接有效性：是 / 否

---

## 5. 修复与后续动作

| # | 动作 | 责任人 | 验收标准 | 截止日期 |
|---:|---|---|---|---|
| 1 | 修正 `xx.md` 定义漂移 | Kimi | 与 Go Spec 逐条比对一致，复跑 verify_sixpiece.py exit 0 | YYYY-MM-DD |
| 2 | 补充交叉域双向链接 | Kimi | `python scripts/tmp/rescan_deadlinks.py` 0 死链 | YYYY-MM-DD |
| 3 | 重新登记/校验 indices | Kimi | complete-map 登记数 == 实际篇数 | YYYY-MM-DD |
| 4 | 修正受影响代码块 | Kimi | `cd <模块> && GOWORK=off go vet ./... && GOWORK=off go test ./...` exit 0 | YYYY-MM-DD |

---

## 6. 证据与基线

- `verify_sixpiece.py` 输出：
- `rescan_deadlinks.py` 输出：
- 相关 `GOWORK=off go vet` / `go test` 输出：
- 相关 CI 链接：

---

> **完成声明纪律**：本报告的所有"已修复"结论必须附质量门 exit 0 截图或命令输出；否则只能标记为"待验证"。

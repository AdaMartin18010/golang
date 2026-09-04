# Kimi 质量门 Checklist

> **EN**: Kimi Quality Gate Checklist
> **Summary**: Scenario-based command checklist for Kimi (and human maintainers) to verify content changes before committing.
> **Scope**: `E:/_src/golang`

---

## 使用说明

根据你本次修改的类型，勾选并执行对应命令。所有命令在仓库根目录执行；涉及 go 命令时**必须 `GOWORK=off`**（根 `go.work` 不覆盖全部 42 个模块），或先 `cd` 进入具体模块目录。**只有本地质量门全部 exit 0，才可声明完成；push 由用户决定。**

---

## 场景 A：新建/修改 `go-knowledge-base/` 权威页

### A1 内容自检（生成后立即）

- [ ] 双语标题与元数据齐全：维度 / 级别（S>16KB、A>8KB、B）/ 标签 / Go 版本 / Bloom 层级 / 前置·后置概念 / 定理链
- [ ] Bloom 层级、FT/LD/EC/TS/AD 编号、维度归属一致
- [ ] 前置/后置概念链接有效且含低层链接（高层 L6/L7 必须引用至少一个低层页）
- [ ] 包含至少 1 个可运行 ```go 块和 1 个首行注释 `// 编译失败: <原因>` 的编译失败反例
- [ ] 包含 Mermaid mindmap（覆盖定义、机制、边界、实践）
- [ ] References 覆盖 P0/P1/P2（go.dev/ref/spec、pkg.go.dev、go.googlesource.com/proposal、github.com/golang ＋ 学术 ＋ go.dev/blog / 知名 Go 项目）
- [ ] 形成至少一对双向链接，并已登记 `go-knowledge-base/indices/`（by-date / by-topic / complete-index）

### A2 必跑命令

```bash
# 1. 死链检查（活跃区 0 死链；跳过 archive/ 与代码围栏）
python scripts/tmp/rescan_deadlinks.py

# 2. 代码块实测（可运行 go 块必须 vet/build 通过；编译失败反例必须确实失败）
python scripts/tmp/extract_go_blocks.py
python scripts/tmp/test_go_blocks.py

# 3. 六件套机械校验（权威定义/核心机制/工程实践/反命题与边界/mindmap/References）
python scripts/tmp/verify_sixpiece.py

# 4. 权威页唯一性（五维内部同主题撞车扫描，发现重复按 canonical 规则合并或 stub 化）
python scripts/tmp/dup_canon_scan.py

# 5. 头部元数据一致性（维度/级别/标签/Go 版本/Bloom；可用标注脚本辅助，最终人工核对 AGENTS.md §3 模板）
python scripts/tmp/annotate_bloom.py

# 6. 双向链接核查（前置/后置链接的目标页须回链）
python scripts/tmp/check_bidir.py

# 7. 全部门（提交前）
bash .githooks/pre-commit
```

---

## 场景 B：版本补丁响应（如 1.27.1）

### B1 内容自检

- [ ] 补丁页头部包含双语标题、摘要、发布日期、Bloom 层级
- [ ] 正文包含对称差分析（A∩B、B\A、A\B）
- [ ] 反命题表 ≥6 条
- [ ] go 指令未提升：全部 `go.mod` 仍为 `go 1.27`（Go 的 patch release 语义：patch 不提升 go 指令，只修工具链与标准库）
- [ ] 受影响的权威页已注入补丁提示并双向链接（版本特性 → 概念权威页，如 → `LD-037` + `examples/go127-features/`）

### B2 必跑命令

```bash
# 1. go 指令一致性（42 个模块统一 go 1.27；archive/ 只读除外）
find . -name go.mod -not -path "./archive/*" -exec grep -H "^go " {} \; | grep -v "go 1.27" || echo "go directives OK"

# 2. 死链检查
python scripts/tmp/rescan_deadlinks.py

# 3. 版本语义注入双向链接（人工核对 + 死链复扫兜底）
python scripts/tmp/check_bidir.py

# 4. 全部门（提交前）
bash .githooks/pre-commit
```

---

## 场景 C：新建/修改 `docs/` / `view/` / `examples/` 工程页

### C1 内容自检

- [ ] 头部声明 canonical 来源链接（概念解释回链五维权威页，禁止在非权威位置复制概念推导）
- [ ] `examples/` README 只链接权威页，不复制通用概念解释
- [ ] 保留应用场景、决策树、操作步骤、链接

### C2 必跑命令

```bash
# 1. 死链检查（覆盖 go-knowledge-base/ + docs/ + view/）
python scripts/tmp/rescan_deadlinks.py

# 2. Markdown 格式规范
powershell -File scripts/check-markdown-format.ps1

# 3. 通用质量检查
powershell -File scripts/check_quality.ps1

# 4. 全部门（提交前）
bash .githooks/pre-commit
```

---

## 场景 D：只改 Kimi 配置/模板/Prompt（如 `.kimi/`）

### D1 内容自检

- [ ] 新文件路径符合 `.kimi/` 命名规范（kebab-case，模板在 `templates/`、提示词在 `prompts/`）
- [ ] `AGENTS.md` 已同步引用新资产
- [ ] 未在 `scripts/tmp/`、构建产物目录中新增正式内容

### D2 必跑命令

```bash
# 1. 死链检查（确保新增 markdown 内部链接有效）
python scripts/tmp/rescan_deadlinks.py

# 2. Markdown 格式规范
powershell -File scripts/check-markdown-format.ps1

# 3. 全部门（提交前）
bash .githooks/pre-commit
```

---

## 场景 E：新增/修改 `pkg/` / `internal/` / `cmd/` / `examples/` Go 代码

### E1 内容自检

- [ ] 所在模块 `go.mod` 的 go 指令为 `go 1.27`
- [ ] 新增示例在对应模块内可编译、可测试，不假装自包含
- [ ] 未在仓库根裸跑 go 命令（一律 `GOWORK=off` 或进入模块目录）

### E2 必跑命令

```bash
# 1. 构建 + vet + 测试（进入具体模块执行）
cd <module-dir> && GOWORK=off go build ./... && GOWORK=off go vet ./... && GOWORK=off go test ./...

# 2. Lint
golangci-lint run

# 3. 示例代码校验（如适用）
powershell -File scripts/validate_code_samples.ps1

# 4. 全部门（提交前）
bash .githooks/pre-commit
```

---

## 场景 F：新增/修改 Quiz（如启用）

### F1 内容自检

- [ ] quiz 文件位于五维权威页配套 quiz 目录（结构以 `.kimi/prompts/quiz_generation_prompt.md` 定义为准）
- [ ] 题型 ≥3 种
- [ ] 每题标注难度与 Bloom 层级
- [ ] quiz→权威页与权威页→quiz 双向链接

### F2 必跑命令

```bash
# 1. 死链检查
python scripts/tmp/rescan_deadlinks.py

# 2. 双向链接核查
python scripts/tmp/check_bidir.py

# 3. 全部门（提交前）
bash .githooks/pre-commit
```

---

## 通用红线

- 禁止在 `scripts/tmp/`、构建产物目录中直接写正式内容。
- 禁止复制已有权威页正文到非权威位置（重复内容按 canonical 规则合并或 stub 化，stub ≤ 25 行 / 2000 字节）。
- 禁止未经验证声明"已完成"或"全部通过"——以构建/测试/链接检查结果为准。
- 新增权威页必须同步登记 `go-knowledge-base/indices/`。
- KG 关系必须用语义谓词（`dependsOn`/`entails`/`mutexWith`/`refines`/`equivalentTo`/`counterExample`），禁止通用 `RelationAnnotation`。
- 提交信息惯例为 `update`；push 由用户决定。

---

## 快速命令速查

```bash
# 一键全部门（预提交钩子：gofmt → go vet → golangci-lint → 相关包测试）
bash .githooks/pre-commit

# 单项快速检查
python scripts/tmp/rescan_deadlinks.py     # 死链扫描
python scripts/tmp/extract_go_blocks.py    # 提取权威页 go 代码块
python scripts/tmp/test_go_blocks.py       # go 块实测（可运行必过、反例必失败）
python scripts/tmp/verify_sixpiece.py      # 六件套机械校验
python scripts/tmp/dup_canon_scan.py       # 权威页唯一性/同主题撞车
python scripts/tmp/check_bidir.py          # 前置/后置双向链接核查
python scripts/tmp/annotate_bloom.py       # Bloom/概念链标注辅助

# 文档格式与通用质量
powershell -File scripts/check-markdown-format.ps1
powershell -File scripts/check_quality.ps1
powershell -File scripts/check-unfixed-links.ps1
```

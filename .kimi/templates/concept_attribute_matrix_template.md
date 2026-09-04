# 概念属性矩阵模板

> **配套要求**：[`.kimi/kimi_thinking_representation_requirements.md`](../kimi_thinking_representation_requirements.md) §4

将以下矩阵嵌入五维权威页 `go-knowledge-base/0X-维度/{FT|LD|EC|TS|AD}-NNN-{Kebab-Title}.md` 的「权威定义」节后。

---

## 五元组表格

```markdown
| 元素 | 内容 |
|---|---|
| **定义** | 一句话精确语义定义。 |
| **属性** | 1. 属性 A：... <br> 2. 属性 B：... <br> 3. 属性 C：... |
| **关系** | - `dependsOn` [概念 B](../../xx/xx.md) <br> - `mutexWith` [概念 C](../../yy/yy.md) <br> - `equivalentTo` [概念 D](../../zz/zz.md) |
| **示例** | ```go\npackage main\n\nfunc main() {\n\t// 正确用法\n}\n``` |
| **反例** | ```go\n// 编译失败: declared and not used: x\npackage main\n\nfunc main() {\n\tx := 1\n}\n``` |
```

---

## 关系谓词速查

| 谓词 | 使用场景 |
| --- | --- |
| `dependsOn` | A 依赖 B 的前提/定义 |
| `entails` | A 语义上蕴含 B |
| `mutexWith` | A 与 B 不能同时成立 |
| `refines` | A 是 B 的细化/特例 |
| `equivalentTo` | A 与 B 在特定上下文等价 |
| `counterExample` | A 是 B 的反例 |

---

## 检查清单

- [ ] 定义可转化为判定过程
- [ ] 属性为不可撤销的不变式（Go 语境：语言规范保证 / 运行时 / GC / happens-before 约束）
- [ ] 关系使用具体谓词，非 generic `related to`
- [ ] 示例为可运行 `go` 块（自包含，`GOWORK=off go vet` 可通过）
- [ ] 反例为 `go` 块，首行注释 `// 编译失败: <确定性编译期错误原因>`（Go 编译器无公开错误码体系，写错误类别文本，如 `declared and not used`、`cannot use … (variable of type X) as Y`、`invalid operation`），并已在真实模块中 `GOWORK=off go build` 复现

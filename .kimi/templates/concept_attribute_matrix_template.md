# 概念属性矩阵模板

> **配套要求**：[`.kimi/kimi_thinking_representation_requirements.md`](../kimi_thinking_representation_requirements.md) §4

将以下矩阵嵌入 `concept/` 权威页的「权威定义」节后。

---

## 五元组表格

```markdown
| 元素 | 内容 |
|---|---|
| **定义** | 一句话精确语义定义。 |
| **属性** | 1. 属性 A：... <br> 2. 属性 B：... <br> 3. 属性 C：... |
| **关系** | - `dependsOn` [概念 B](../../xx/xx/xx.md) <br> - `mutexWith` [概念 C](../../yy/yy/yy.md) <br> - `equivalentTo` [概念 D](../../zz/zz/zz.md) |
| **示例** | ```rust\nfn main() {\n    // 正确用法\n}\n``` |
| **反例** | ```rust,compile_fail\nfn main() {\n    // 错误用法，期望 E0xxx\n}\n``` |
```

---

## 关系谓词速查

| 谓词 | 使用场景 |
|---|---|
| `dependsOn` | A 依赖 B 的前提/定义 |
| `entails` | A 语义上蕴含 B |
| `mutexWith` | A 与 B 不能同时成立 |
| `refines` | A 是 B 的细化/特例 |
| `equivalentTo` | A 与 B 在特定上下文等价 |
| `counterExample` | A 是 B 的反例 |

---

## 检查清单

- [ ] 定义可转化为判定过程
- [ ] 属性为不可撤销的不变式
- [ ] 关系使用具体谓词，非 generic `related to`
- [ ] 示例为可运行 `rust` 块
- [ ] 反例为 `rust,compile_fail` 并标注错误码

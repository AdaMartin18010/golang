# R2: 定理索引

> **层级**: 参考层 (Reference)
> **地位**: 定理快速索引

---

## 定理清单

| 编号 | 名称 | 陈述 | 位置 | 依赖 |
|------|------|------|------|------|
| **Th1.1** | new语义等价 | ∀T,v. new(v) ≡ &v' | [D3](../derivation/D3-theorems.md#Th1.1) | A1-A3, A6 |
| **Th1.2** | 递归泛型终止 | wellformed(C) → terminates(unfold(C)) | [D3](../derivation/D3-theorems.md#Th1.2) | A5, A7 |
| **Th2.1** | GC低延迟 | P(pause < 1ms) ≥ 0.99 | [D3](../derivation/D3-theorems.md#Th2.1) | A8 |

## 公理清单

| 编号 | 名称 | 位置 |
|------|------|------|
| A1 | 内存分配公理 | [F2](../foundation/F2-axioms.md#A1) |
| A2 | 值存储公理 | [F2](../foundation/F2-axioms.md#A2) |
| A3 | 指针语义公理 | [F2](../foundation/F2-axioms.md#A3) |
| A4 | 类型等价公理 | [F2](../foundation/F2-axioms.md#A4) |
| A5 | 泛型实例化公理 | [F2](../foundation/F2-axioms.md#A5) |
| A6 | new表达式公理 | [F2](../foundation/F2-axioms.md#A6) |
| A7 | 递归约束公理 | [F2](../foundation/F2-axioms.md#A7) |
| A8 | 并发GC公理 | [F2](../foundation/F2-axioms.md#A8) |

## 概念清单

| 编号 | 名称 | 位置 |
|------|------|------|
| C1 | new表达式 | [F3](../foundation/F3-core-concepts.md#C1) |
| C2 | 递归泛型约束 | [F3](../foundation/F3-core-concepts.md#C2) |
| C3 | GreenTeaGC | [F3](../foundation/F3-core-concepts.md#C3) |
| C4 | HPKE | [F3](../foundation/F3-core-concepts.md#C4) |

---

**快速参考**: [R3-快速参考](R3-quick-reference.md)

# {中文标题}

> **EN**: {English Title}
> **Summary**: One-sentence abstract covering the design space, current syntax proposal, and stability status of the {effect name} effect.
>
> **Rust 版本**: 1.99 beta / nightly / 设计提案
> **Bloom 层级**: L7
> **权威来源**: 本文件为 `concept/` 权威页（Rust 1.99+ 预览特性跟踪页）。
> **受众**: 研究者 / 前沿特性关注者
> **内容分级**: 综述级
> **A/S/P 标记**: **S+A** — Structure + Application
> **双维定位**: C×Res   <!-- Concept × Research -->
> **前置概念**: [Effects and Purity](../../01_foundation/21_effects_and_purity.md) · [Async](../../03_advanced/01_async/01_async.md)
> **后置概念**: [Rust 1.99 前沿特性预览](../../07_future/rust_1_99_preview.md)
> **定理链**: 设计动机 → 语法提案 → 组合规则 → 稳定性边界

---

## 1. 权威定义

一句话精确定义该效果是什么、它解决什么问题、当前稳定状态如何。

| 属性 | 说明 |
|---|---|
| **效果名称** | ... |
| **稳定状态** | nightly / beta / 设计提案 / 已废弃方向 |
| **核心问题** | ... |
| **最小公倍数** | 该效果如何消除 2^N 组合爆炸 |

---

## 2. 学术谱系与设计动机

### 2.1 关键文章 / RFC / 论文时间线

| 时间 | 来源 | 核心贡献 | 本文件对应节 |
|---|---|---|---|
| YYYY-MM | Author / RFC / Paper | ... | §X.Y |

### 2.2 为什么需要这个效果？

- 当前 `Result`/`Option`/`async` 等机制的局限性。
- 该效果提供的统一抽象。

---

## 3. 设计空间分类矩阵

| 维度 | 选项 A | 选项 B | Rust 当前方向 |
|---|---|---|---|
| 开放 / 封闭 | 开放效果系统 | 封闭效果系统 | ... |
| 显式 / 隐式 | 显式 `throws` | 隐式异常 | ... |
| 静态 / 动态 | 编译期检查 | 运行时检查 | ... |

---

## 4. 当前语法提案

> ⚠️ 以下语法为**设计提案**，非稳定决策。在稳定 Rust 上可能无法编译。

```rust,ignore
// 示例：函数签名中的效果声明
fn may_fail() -> i32 throws MyError {
    // ...
}
```

### 4.1 语法变体对比

| 变体 | 写法 | 适用场景 | 风险 |
|---|---|---|---|
| 变体 1 | ... | ... | ... |
| 变体 2 | ... | ... | ... |

### 4.2 效果传播

```rust,ignore
// .do / with-clauses / ? 在效果上下文中的传播示例
```

---

## 5. 效果代数与组合规则

| 操作 | 语义 | 示例 |
|---|---|---|
| 并集 | 同时具有两种效果 | `eff A + B` |
| 排除 | 明确不具有某种效果 | `eff A - B` |
| 互斥 | 两种效果不能共存 | `eff A mutexWith B` |
| 别名 | 常用效果组合命名 | `eff Pure = panic + diverge` |

---

## 6. 与现有概念的交叉分析

### 6.1 Effect × Async

- `Future::poll` 的 effect 语义。
- async 状态机如何携带效果信息。

### 6.2 Effect × Pin

- `Pin` 作为 async 效果的“附属类型系统”。
- 统一效果系统中 `Pin` 的角色演变。

### 6.3 Effect × Const

- const fn 与效果泛型的边界。
- `~const` 等历史语法的废弃方向说明。

---

## 7. 反命题与边界分析

| 命题 | 真假 | 说明 |
|---|---|---|
| 该效果已在稳定 Rust 中可用 | ❌ | 当前仅为提案 / nightly |
| 该效果可完全替代现有 `Result` 机制 | ❌/⚠️ | 取决于最终设计 |
| 效果系统选择封闭而非开放 | ✅ | 见 §3 矩阵 |

### 7.1 当前不可编译示例

```rust,ignore
// 说明：以下写法在稳定 Rust 下会报错，因为该效果尚未稳定。
fn stable_fn() -> i32 throws Error { 0 }
```

---

## 8. 版本语义注入与双向链接

- 受影响的 `concept/` 权威页：
  - [](../xx/xx.md)
  - [](../yy/yy.md)
- 确保上述页在“版本兼容性”或“相关概念”节反向链接回本页。

---

## 9. 思维导图

```mermaid
mindmap
  root(({中文标题}))
    学术谱系
      关键文章
      设计动机
    设计空间
      开放/封闭
      显式/隐式
      静态/动态
    语法提案
      函数签名
      传播规则
      效果别名
    交叉概念
      Async
      Pin
      Const
    边界
      未稳定
      与 Result 的对比
      迁移风险
```

---

## 10. 参考来源 / References

- **P0 官方**: [Rust Project Goals](https://github.com/rust-lang/rust-project-goals) · [RFC 索引](https://rust-lang.github.io/rfcs/)
- **P1 学术/设计文章**: [Author blog / paper](...) · [相关论文 DOI](...)
- **P2 生态/社区**: [相关 crate docs](https://docs.rs/...) · [讨论帖](...)

---

## 11. 相关概念（双向链接）

- [Effects and Purity](../../01_foundation/21_effects_and_purity.md)
- [Async](../../03_advanced/01_async/01_async.md)
- [Pin / Unpin](../../03_advanced/01_async/08_pin_unpin.md)
- [Rust 1.99 前沿特性预览](../../07_future/rust_1_99_preview.md)

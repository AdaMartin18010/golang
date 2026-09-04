# 可复用 Prompt：代码示例与反例生成

> **用途**：为 `concept/` 权威页生成符合 10 桶规则的 `rust` 可运行示例与 `rust,compile_fail` 反例。
> **配套文件**：[`.kimi/templates/kimi_project_requirements.md`](../templates/kimi_project_requirements.md) §4、§12。

---

## Prompt 模板

```text
请在 E:/_src/rust-lang 仓库中为以下概念页生成或补全代码示例与反例。

**目标页面**：{{target_path}}
**概念**：{{concept}}
**需要的示例类型**：{{example_types}}  <!-- 可运行示例 / 反例 / 依赖 crate 示例 / 测试 / 反模式 -->
**是否依赖外部 crate**：{{needs_deps}}  <!-- 是 / 否；若是，列出 crate -->

要求：
1. 每个 `rust` 可运行示例必须自包含，包含 `fn main()` 或 `#[test]`，且不依赖未声明的变量/函数。
2. 每个 `rust,compile_fail` 反例必须：
   - 确实在当前稳定 Rust（{{rust_version}}）下编译失败
   - 标注期望的错误码（如 `E0382`）或错误类别
   - 在正文解释失败原因与正确写法
3. 若需要外部 crate，使用 `rust,ignore` 并在对应 crate 的 `examples/` 中提供可运行版本；或先确认 `--with-deps` 能否编译。
4. 反例必须放在「反命题与边界分析」节，并与定理链/命题表对应。
5. 禁止在代码块中使用未引入的 crate 或假定存在外部上下文。
6. 生成后运行：
   python scripts/check_concept_code_blocks.py --strict
   # 若含依赖块：
   cargo build --workspace
   python scripts/check_concept_code_blocks.py --strict --with-deps
7. 不要声明“已完成”，除非命令 exit 0。
```

---

## 示例填充

```text
**target_path**: concept/02_intermediate/00_traits/01_traits.md
**concept**: Trait object upcasting 与 vtable
**example_types**: 可运行示例 + 反例
**needs_deps**: 否
**rust_version**: 1.98.0
```

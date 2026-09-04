# 可复用 Prompt：代码示例与反例生成

> **用途**：为五维权威页（`go-knowledge-base/01..05-*`）生成可运行的 ```go 示例与 `// 编译失败:` 反例。
> **配套文件**：[`.kimi/templates/kimi_project_requirements.md`](../templates/kimi_project_requirements.md)「内容质量要求」节。

---

## Prompt 模板

```text
请在 E:/_src/golang 仓库中为以下权威页生成或补全代码示例与反例。

**目标页面**：{{target_path}}
**概念**：{{concept}}
**需要的示例类型**：{{example_types}}  <!-- 可运行示例 / 反例 / 依赖外部模块示例 / 测试 / 反模式 -->
**是否依赖外部模块**：{{needs_deps}}  <!-- 是 / 否；若是，列出 module -->

要求：
1. 每个可运行 ```go 示例必须自包含（package main + func main()，或独立 _test.go），不依赖未声明的标识符。
2. 每个反例必须：
   - 首行注释 `// 编译失败: <一句话原因>`
   - 确实在当前稳定 Go（{{go_version}}）下编译失败
   - 标注期望的编译器错误（如 `cannot use ... as ...`、`invalid operation`）
   - 在正文解释失败原因与正确写法
3. 若需要外部依赖，把示例放入对应 examples/<module>/（含 go.mod），正文内只做引用不假装自包含。
4. 反例必须放在「反命题与边界分析」节，并与定理链/命题表对应。
5. 禁止在代码块中使用未 import 的包或假定存在外部上下文。
6. 生成后在临时模块中验证（禁止在仓库根裸跑 go 命令）：
   cd <临时模块> && GOWORK=off go vet ./... && GOWORK=off go test ./...
7. 不要声明"已完成"，除非命令 exit 0。
```

---

## 示例填充

```text
**target_path**: go-knowledge-base/02-Language-Design/LD-037-Go-1.27-Generic-Methods.md
**concept**: 泛型方法的接收者类型参数推断
**example_types**: 可运行示例 + 反例
**needs_deps**: 否
**go_version**: 1.27.1
```

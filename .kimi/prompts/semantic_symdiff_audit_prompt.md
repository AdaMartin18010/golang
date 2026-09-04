# 可复用 Prompt：语义对称差审计

> **用途**：比较两个 Go 版本、两个概念或两种写法之间的语义/语法/API/行为差异，输出集合论形式的对称差分析。
> **配套文件**：[`.kimi/templates/kimi_project_requirements.md`](../templates/kimi_project_requirements.md)「红线」节（版本特性必须映射回概念权威页）。

---

## Prompt 模板

```text
请在 E:/_src/golang 仓库中对以下主题进行语义对称差审计。

**对比对象 A**：{{object_a}}
**对比对象 B**：{{object_b}}
**审计维度**：{{dimensions}}  <!-- 例如：语法 / 语义 / 标准库 API / 运行时行为 / 平台支持 / 工具链 -->
**目标位置**：{{target_path}}  <!-- 已有权威页或版本跟踪页 -->

要求：
1. 先读取 A、B 在五维中的权威页（若存在），避免重复定义。
2. 用集合记号明确写出：
   - A ∩ B：两者完全相同的部分
   - B \ A：仅 B 有的部分
   - A \ B：仅 A 有的部分
3. 每个维度用表格呈现，避免大段文字堆砌。
4. 对 B \ A 和 A \ B 中的每一项给出：
   - 具体变更/差异描述
   - 触发条件或代码模式
   - 迁移/修复建议
   - 相关编译器/vet 错误文本（如适用）
5. 至少包含一个 ```go 反例（`// 编译失败:` 注释）展示"在 A 上可行、在 B 上失败"或反之，并在临时模块中实测：
   cd <临时模块> && GOWORK=off go build ./...   # 反例应非零退出
6. 生成 Mermaid mindmap 总结对称差。
7. 建立与相关权威页的双向链接；版本差异必须映射回概念权威页。
8. 验证：
   python scripts/tmp/rescan_deadlinks.py
9. 不要声明"已完成"，除非检查通过。
```

---

## 示例填充

```text
**object_a**: Go 1.26
**object_b**: Go 1.27
**dimensions**: 语法（泛型方法）、运行时（Green Tea GC）、标准库（json/v2 实验包）、工具链（go.mod go 指令）
**target_path**: go-knowledge-base/02-Language-Design/03-Evolution/07-Go126-to-Go127.md
```

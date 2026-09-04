# Prompt: 生成 Rust 交叉 / 边界语义域权威页

## 角色

你是 Rust 语义深度文档作者。你的任务是为两个或多个核心概念交叉形成的边界语义域新建 `concept/` 权威页。

## 输入

- 交叉域名称：{name}
- 参与概念：{domain A}, {domain B}, {optional domain C}
- 涉及的 rustc error codes：{E0xxx, E0yyy}
- 相关 `concept/` 页路径：{paths}

## 任务

1. 使用 `.kimi/templates/kimi_cross_domain_concept_template.md` 作为骨架。
2. 页面头部必须包含：
   - EN 标题与 Summary
   - Rust 版本
   - Bloom 层级 L4
   - 权威来源声明
   - 前置/后置概念（必须包含所有参与概念）
   - 定理链
3. 正文必须包含：
   - 一句话权威定义与组合不变式
   - 涉及概念矩阵（≥4 个维度）
   - 核心机制与组合规则（使用 `⟹` / `⟸`）
   - 形式化契约或不变式表格
   - 反命题与边界分析（≥6 条命题）
   - 至少一个 `rust,compile_fail` 反例，标注错误码
   - 决策树 / 判定表，并链接到 `decision_trees.yaml`
   - Mermaid mindmap
   - References（P0/P1/P2）
4. 确保与所有参与概念页形成双向链接。
5. 若涉及 KG 实体，提示应生成 `dependsOn`/`mutexWith`/`counterExample` 等具体谓词关系。

## 输出格式

- 输出完整 markdown 文件内容（可直接写入 `concept/03_advanced/...` 等）。
- 使用中文，EN 标题与 Summary 必须英文。
- 代码块自包含；std-only 示例必须可直接 `rustc --edition 2024` 编译。

## 约束

- 禁止将通用概念推导复制到本页；只保留组合后的新语义。
- 每个反例必须解释失败原因与正确写法。
- 必须引用至少一个 P1 形式化/学术来源。

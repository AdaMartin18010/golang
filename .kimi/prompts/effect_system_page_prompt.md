# Prompt: 生成 Rust 效应系统 / 预览特性权威页

## 角色

你是 Rust 前沿特性文档作者。你的任务是为一个 Rust 预览特性（尤其是效应系统相关主题）新建或更新 `concept/` 权威页。

## 输入

- 特性名称：{feature name}
- 当前稳定状态：{nightly / beta / 设计提案 / 已废弃方向}
- 关键权威来源：{blog posts / RFCs / papers / project goals}
- 相关 `concept/` 页：{paths}

## 任务

1. 使用 `.kimi/templates/kimi_effect_system_page_template.md` 作为骨架。
2. 在页面头部明确标注：
   - Rust 版本（nightly / beta / 设计提案）
   - Bloom 层级 L7
   - 权威来源声明
   - 前置/后置概念
3. 正文必须包含：
   - 权威定义与稳定状态
   - 学术谱系 / 设计动机时间线
   - 设计空间分类矩阵（≥3 个维度）
   - 当前语法提案（使用 `rust,ignore`，明确标注“提案”）
   - 效果代数 / 组合规则
   - 与现有概念（async / Pin / const 等）的交叉分析
   - 反命题与边界分析（≥6 条命题）
   - 版本语义注入与双向链接
   - Mermaid mindmap
   - References（P0/P1/P2）
4. 确保所有代码块使用 `rust,ignore` 或明确说明不可编译。
5. 不要复制已有权威页正文；对已有概念只给出链接和一句话摘要。

## 输出格式

- 输出完整 markdown 文件内容（可直接写入 `concept/07_future/...`）。
- 使用中文，EN 标题与 Summary 必须英文。
- 代码块自包含，避免依赖未声明变量。

## 约束

- 禁止声明该特性“已稳定”除非有官方 release notes 证据。
- 语法提案必须标注为设计提案。
- 必须链接到受影响的 `concept/` 权威页，并提示需要在这些页中添加反向链接。

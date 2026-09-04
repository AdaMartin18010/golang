# Prompt: 生成 Go 前沿提案跟踪页（原效应系统 / 预览特性页）

> **状态**: 可选扩展 — 本仓库 AGENTS.md 体系尚未启用该机制，从 rust-lang 项目适配保留以便未来复用；不进入质量门。

## 角色

你是 Go 语言前沿提案文档作者。你的任务是为一个尚未进入稳定工具链的 Go 前沿主题（语言设计提案、实验特性、runtime/GC/调度器演进方向）新建或更新跟踪页。

## 输入

- 主题名称：{proposal / feature name}
- 当前状态：{设计草案 / 提案评审中（go.googlesource.com/proposal）/ 实验实现（GOEXPERIMENT）/ 已废弃方向}
- 关键权威来源：{design documents / proposal 讨论 / proposal review meetings / Go Blog}
- 相关五维权威页：{paths}

## 任务

1. 使用 `.kimi/templates/kimi_effect_system_page_template.md` 作为骨架。
2. 在页面头部明确标注：
   - Go 版本（实验：`GOEXPERIMENT=<name>` 于 1.xx；或"未合入任何版本，提案阶段"）
   - Bloom 层级 L7
   - 权威来源声明
   - 前置/后置概念
3. 正文必须包含：
   - 权威定义与稳定状态
   - 设计动机时间线与学术谱系
   - 设计空间分类矩阵（≥3 个维度）
   - 当前提案形态（```go + 首行 `// 不可编译: 提案语法，当前工具链不支持`，明确标注"设计提案"）
   - 与现有概念（interface 方法集、泛型实现方式、channel happens-before、cgo/unsafe 边界、GC 语义等）的交叉分析
   - 反命题与边界分析（≥6 条命题）
   - 版本语义注入与双向链接
   - Mermaid mindmap
   - References（P0/P1/P2）
4. 确保所有代码块使用 ```go + 首行 `// 不可编译: <原因>` 注释，或明确说明不可编译。
5. 不要复制已有权威页正文；对已有概念只给出链接和一句话摘要。

## 输出格式

- 输出完整 markdown 文件内容（可直接写入 `go-knowledge-base/0X-维度/` 或 `docs/tracking/`）。
- 使用中文，标题双语（`# 中文标题 (English Title)`）。
- 代码块自包含，避免依赖未声明变量。

## 约束

- 禁止声明该特性"已稳定"除非有官方 release notes 证据。
- 提案语法必须标注为设计提案，不得伪装成可运行代码。
- 必须链接到受影响的五维权威页，并提示需要在这些页中添加反向链接。

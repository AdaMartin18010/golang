# .kimi — 协作治理模板

本目录存放 Kimi 协作的治理模板，已从 rust-lang 项目（E:/_src/rust-lang）的要求全面适配为本 Go 项目（Go 1.27.1）版本。

| 模板 | 用途 |
| --- | --- |
| [`templates/kimi_project_requirements.md`](templates/kimi_project_requirements.md) | 项目协作总要求：架构、命名、元数据、内容质量、权威来源分级、构建、质量门、红线 |
| [`templates/monthly_semantic_review.md`](templates/monthly_semantic_review.md) | 月度语义审查清单：定义漂移、边界语义、stub 纯净度、KG 关系、版本语义注入 |
| [`templates/quarterly_international_source_audit.md`](templates/quarterly_international_source_audit.md) | 季度权威源审计：对照 Go Spec / Effective Go / Go Blog / pkg.go.dev / Proposals |
| [`templates/concept_page_template.md`](templates/concept_page_template.md) | 权威概念页骨架：AGENTS.md §3 头部 + 六件套章节 |
| [`prompts/`](prompts/) | 内容生成提示词：概念页 / 代码示例 / 形式化定理补全 / 语义对称差审计 / patch 响应 / quiz（可选扩展，本项目未启用） |

执行摘要已同步到仓库根 [`AGENTS.md`](../AGENTS.md)，日常协作以 AGENTS.md 为准；本目录模板用于周期性审查与复用到其他 Go 项目。

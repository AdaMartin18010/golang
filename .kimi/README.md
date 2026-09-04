# .kimi — 协作治理模板

本目录存放 Kimi 协作的治理模板，已从 rust-lang 项目（E:/_src/rust-lang）的要求全面适配为本 Go 项目（Go 1.27.1）版本。

| 类别 | 文件 | 状态 |
| --- | --- | --- |
| 项目总要求 | [`templates/kimi_project_requirements.md`](templates/kimi_project_requirements.md) | 核心启用 |
| 内容生成要求 | [`kimi_content_requirements.md`](kimi_content_requirements.md) | 核心启用 |
| 语义要求 | [`kimi_semantic_requirements.md`](kimi_semantic_requirements.md) | 核心启用 |
| 质量门清单 | [`kimi_quality_gate_checklist.md`](kimi_quality_gate_checklist.md) | 核心启用 |
| KG 拓扑要求 | [`kimi_kg_topology_requirements.md`](kimi_kg_topology_requirements.md) | 核心启用 |
| 表征空间分析 | [`kimi_semantic_space_analysis_requirements.md`](kimi_semantic_space_analysis_requirements.md) | 可选扩展 |
| 思维表征 | [`kimi_thinking_representation_requirements.md`](kimi_thinking_representation_requirements.md) | 可选扩展 |
| 月度审查 | [`templates/monthly_semantic_review.md`](templates/monthly_semantic_review.md) | 核心启用 |
| 季度权威源审计 | [`templates/quarterly_international_source_audit.md`](templates/quarterly_international_source_audit.md) | 核心启用 |
| 概念页骨架 | [`templates/concept_page_template.md`](templates/concept_page_template.md) | 核心启用 |
| 属性矩阵 | [`templates/concept_attribute_matrix_template.md`](templates/concept_attribute_matrix_template.md) | 核心启用 |
| 交叉域概念页 | [`templates/kimi_cross_domain_concept_template.md`](templates/kimi_cross_domain_concept_template.md) | 核心启用 |
| 语义审计记录 | [`templates/kimi_semantic_audit_template.md`](templates/kimi_semantic_audit_template.md) | 核心启用 |
| mindmap | [`templates/mindmap_template.md`](templates/mindmap_template.md) | 核心启用 |
| 决策树 | [`templates/decision_tree_template.md`](templates/decision_tree_template.md) | 可选扩展（Go 无错误码体系，映射为编译失败类别） |
| 前沿提案跟踪页 | [`templates/kimi_effect_system_page_template.md`](templates/kimi_effect_system_page_template.md) | 可选扩展 |
| 表征空间页 | [`templates/kimi_semantic_space_page_template.md`](templates/kimi_semantic_space_page_template.md) | 可选扩展 |
| 提示词 | [`prompts/`](prompts/)（14 个：概念生成 / 代码示例 / 定理补全 / 语义审计×3 / 漂移评审 / 交叉域 / mindmap / KG 谓词 / patch / quiz（未启用）/ 对称差审计 / 决策树 / 提案页） | 核心启用为主，决策树与提案页为可选扩展 |

**核心启用** = 本仓库 AGENTS.md 体系正在使用的治理机制；**可选扩展** = 从 rust-lang 项目适配保留、尚未启用、不进入质量门的机制（文件头均有降级标注）。

执行摘要已同步到仓库根 [`AGENTS.md`](../AGENTS.md)，日常协作以 AGENTS.md 为准；本目录模板用于周期性审查与复用到其他 Go 项目。

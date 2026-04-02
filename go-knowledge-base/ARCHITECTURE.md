# 知识库架构 (Knowledge Base Architecture)

## 统一编号体系

```
[维度代码]-[三位序号] 文档标题
```

| 维度代码 | 维度名称 | 范围 | 文档数 |
|---------|---------|------|--------|
| FT | Formal Theory (形式理论) | FT-001 ~ FT-099 | 26 |
| LD | Language Design (语言设计) | LD-001 ~ LD-099 | 40 |
| EC | Engineering CloudNative (工程与云原生) | EC-001 ~ EC-199 | 148 |
| TS | Technology Stack (技术栈) | TS-001 ~ TS-099 | 56 |
| AD | Application Domains (应用领域) | AD-001 ~ AD-099 | 43 |

## 内容质量标准

| 级别 | 代码 | 最小大小 | 要求 |
|------|------|---------|------|
| 基础 | B | 5 KB | 概念清晰，可运行代码片段 |
| 进阶 | A | 10 KB | 完整实现，生产级错误处理 |
| 深入 | S | 15 KB | 源码分析，性能基准，形式化定义 |

## 文档元数据格式

```markdown
---
id: EC-042
dimension: Engineering CloudNative
title: 任务调度器完整实现
level: S
size: 28KB
tags: [scheduler, distributed-systems, production]
related: [EC-038, EC-039, FT-015]
last_updated: 2026-04-02
---
```

## 目录结构

```
go-knowledge-base/
├── ARCHITECTURE.md          # 本文件：架构说明
├── INDEX.md                 # 主索引
├── FT-001~026/              # 形式理论
├── LD-001~040/              # 语言设计
├── EC-001~148/              # 工程与云原生
│   └── scheduled-tasks/     # 子分类：计划任务
├── TS-001~056/              # 技术栈
├── AD-001~043/              # 应用领域
└── cross-references.md      # 跨维度关联
```

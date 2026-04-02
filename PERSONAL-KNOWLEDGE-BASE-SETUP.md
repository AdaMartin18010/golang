# 个人 Go 技术知识库 - 配置指南

> **类型**: 个人知识管理项目
> **目标**: 建立系统化的 Go 技术知识体系
> **更新策略**: 持续跟踪 + 深度研究
> **日期**: 2026-04-02

---

## 一、知识库结构

```
pknowledge/                    # 个人知识库 (Personal Knowledge)
├── 00-inbox/                  # 收件箱 - 待处理信息
│   ├── articles/              # 待读文章
│   ├── papers/                # 待读论文
│   ├── code-snippets/         # 代码片段
│   └── ideas/                 # 想法记录
│
├── 01-reading/                # 阅读中
│   ├── go-source/             # Go 源码阅读
│   ├── papers/                # 论文精读
│   └── books/                 # 书籍笔记
│
├── 02-zettelkasten/           # 卡片盒笔记
│   ├── concepts/              # 概念卡片
│   ├── patterns/              # 模式卡片
│   ├── theories/              # 理论卡片
│   └── implementations/       # 实现卡片
│
├── 03-structures/             # 结构化知识
│   ├── language/              # 语言特性
│   ├── concurrency/           # 并发模型
│   ├── memory/                # 内存模型
│   ├── types/                 # 类型系统
│   ├── compiler/              # 编译器
│   ├── runtime/               # 运行时
│   └── ecosystem/             # 生态系统
│
├── 04-projects/               # 项目实践
│   ├── experiments/           # 实验代码
│   ├── benchmarks/            # 性能测试
│   ├── proofs/                # 形式化证明
│   └── implementations/       # 实现练习
│
├── 05-references/             # 参考资料
│   ├── papers/                # 学术论文
│   ├── articles/              # 技术文章
│   ├── videos/                # 视频课程
│   └── talks/                 # 演讲记录
│
└── 99-archive/                # 归档
    ├── outdated/              # 过时内容
    └── completed/             # 已完成
```

---

## 二、个人工作流

### 2.1 信息收集流程

```
发现信息 → 00-inbox/ → 每日整理 → 分类处理
                ↓
        ┌───────┼───────┐
        ↓       ↓       ↓
   深度阅读  快速参考   归档
        ↓       ↓       ↓
   02-zettelkasten  03-structures  99-archive
```

### 2.2 知识处理流程

#### 第一步：收集 (Inbox)

```bash
# 发现有趣的内容，先放入收件箱
pknowledge/00-inbox/articles/2026-04-02-green-tea-gc.md
```

**每日清理**:

- 阅读并决定去向
- 不超过 10 个未处理项目

#### 第二步：精炼 (Zettelkasten)

创建原子笔记卡片：

```markdown
---
id: 2026040201
title: Green Tea GC Span-based Scanning
date: 2026-04-02
tags: [gc, memory, go126]
references: [go-blog-2026-02]
---

## 核心概念

Green Tea GC 使用 span-based scanning 替代 object-based scanning。

## 关键特性

- 8 KiB 连续内存 span
- 每对象 2 bits 元数据
- SIMD 向量优化

## 性能数据

| 指标 | 改善 |
|------|------|
| GC 开销 | -10~40% |
| Cache Miss | -50% |

## 关联

- [[DRF-SC Memory Model]]
- [[Go Memory Allocator]]

## 问题

- RSS 增加 8-15% 是否可接受？
```

#### 第三步：结构化

整合卡片为专题文档：

```markdown
# Go 垃圾收集器演进

## 版本对比

| 版本 | 算法 | 特点 |
|------|------|------|
| Go 1.0-1.4 | STW | 完全暂停 |
| Go 1.5 | 并发标记 | 写屏障 |
| Go 1.8 | 亚毫秒 | 优化 STW |
| Go 1.26 | Green Tea | Span-based |

## 深入分析

基于 [[Green Tea GC Span-based Scanning]] 卡片...
```

### 2.3 日常维护

#### 每日 (15分钟)

- [ ] 清理 inbox
- [ ] 回顾昨日笔记
- [ ] 记录新想法

#### 每周 (2小时)

- [ ] 整理 Zettelkasten
- [ ] 更新结构化知识
- [ ] 运行跟踪脚本

#### 每月 (1天)

- [ ] 深度研究一个主题
- [ ] 撰写技术文章
- [ ] 复习归档知识

---

## 三、技术深度专题

### 3.1 形式化方法专题

```
pknowledge/03-structures/formal-methods/
├── README.md                  # 专题总览
├── featherweight-go/          # Featherweight Go 研究
│   ├── syntax.md
│   ├── typing-rules.md
│   ├── operational-semantics.md
│   └── proofs.md
├── csp/                       # CSP 模型
│   ├── hoare-csp.md
│   ├── go-adaptation.md
│   └── proof-obligations.md
├── memory-model/              # 内存模型
│   ├── happens-before.md
│   ├── drf-sc-proof.md
│   └── examples.md
└── pi-calculus/               # π-演算
    ├── basic-syntax.md
    ├── mobility.md
    └── go-encoding.md
```

### 3.2 源码阅读专题

```
pknowledge/01-reading/go-source/
├── runtime/
│   ├── malloc.go              # 内存分配
│   ├── mgc.go                 # 垃圾收集
│   ├── proc.go                # 调度器
│   └── chan.go                # 通道实现
├── compiler/
│   ├── ssagen/                # SSA 生成
│   └── walk/                  # AST 遍历
└── reflect/
    └── type.go                # 反射类型系统
```

### 3.3 实验项目

```
pknowledge/04-projects/experiments/
├── green-tea-gc-benchmark/    # GC 性能测试
├── simd-vector-ops/           # SIMD 实验
├── generics-patterns/         # 泛型模式
└── formal-verification/       # 形式化验证练习
```

---

## 四、个人工具链

### 4.1 笔记工具

| 工具 | 用途 | 配置 |
|------|------|------|
| **Obsidian** | 主笔记工具 | Zettelkasten 插件 |
| **VS Code** | 代码 + Markdown | Go 插件 |
| **Git** | 版本控制 | 私有仓库 |

### 4.2 阅读工具

| 工具 | 用途 |
|------|------|
| **Zotero** | 论文管理 |
| **PDF Expert** | PDF 批注 |
| **Readwise** | 文章高亮同步 |

### 4.3 实验环境

```bash
# 个人实验目录
~/experiments/go/
├── go1.26-custom/            # 自定义 Go 构建
├── benchmark-suite/          # 性能测试套件
└── proof-of-concepts/        # 概念验证
```

---

## 五、跟踪清单

### 5.1 官方跟踪

- [ ] Go Blog 每周阅读
- [ ] Go Release 每月检查
- [ ] Go Proposal 每周审查
- [ ] CL 每周概览

### 5.2 学术跟踪

- [ ] arXiv 每周搜索
- [ ] POPL/PLDI 每届论文
- [ ] 相关期刊每月检查

### 5.3 社区跟踪

- [ ] Reddit r/golang 每日浏览
- [ ] Hacker News 每日检查
- [ ] Go Time 每期收听

---

## 六、输出目标

### 6.1 个人产出

| 产出类型 | 频率 | 目标 |
|----------|------|------|
| 原子笔记 | 每日 | 3-5 张 |
| 深度文章 | 每月 | 2-3 篇 |
| 实验代码 | 每周 | 1-2 个 |
| 源码阅读 | 每周 | 1 个包 |

### 6.2 质量指标

- 笔记引用网络密度 > 50%
- 深度阅读完成率 > 80%
- 实验可复现率 100%

---

## 七、快速开始

### 7.1 初始化命令

```bash
# 创建知识库结构
mkdir -p pknowledge/{00-inbox/{articles,papers,code-snippets,ideas},01-reading/{go-source,papers,books},02-zettelkasten/{concepts,patterns,theories,implementations},03-structures/{language,concurrency,memory,types,compiler,runtime,ecosystem},04-projects/{experiments,benchmarks,proofs,implementations},05-references/{papers,articles,videos,talks},99-archive/{outdated,completed}}

# 初始化 Git
cd pknowledge
git init
echo "# Personal Go Knowledge Base" > README.md
git add .
git commit -m "Initial knowledge base setup"
```

### 7.2 第一条笔记

```bash
# 创建第一张卡片
cat > 02-zettelkasten/concepts/2026040201-green-tea-gc.md << 'EOF'
---
id: 2026040201
title: Green Tea GC
date: 2026-04-02
tags: [gc, memory, go126]
---

Go 1.26 引入的新垃圾收集器。

主要改进：
1. Span-based scanning
2. SIMD 优化
3. 更低的 GC 延迟

需要深入研究内存布局。
EOF
```

---

## 八、持续改进

### 8.1 每月回顾

- 哪些笔记最有价值？
- 哪些链接缺失？
- 系统如何优化？

### 8.2 季度调整

- 更新知识结构
- 归档过时内容
- 调整专题重点

---

*个人知识库配置完成*
*开始系统化学习之旅*

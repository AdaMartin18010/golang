# AI原生可观测性开源项目分析

> 构建面向AI的开源项目知识图谱与可调用的运行时分析系统

---

## 项目愿景

让AI能够像资深工程师一样理解和使用开源可观测性项目：

- 🔍 **深度理解** - 掌握项目架构、API语义和运行时行为
- 💬 **自然语言交互** - 用日常语言查询项目知识
- 🛠️ **智能生成** - 自动生成集成代码和配置
- 📊 **运行时诊断** - 基于实时数据提供优化建议

## 核心特性

| 特性 | 描述 | 状态 |
|------|------|------|
| **知识图谱** | 以图结构存储项目组件、接口、配置、指标关系 | 🚧 开发中 |
| **MCP服务** | 通过Model Context Protocol提供AI调用接口 | 🚧 开发中 |
| **自然语言查询** | 支持自然语言到Cypher查询转换 | 📋 计划中 |
| **代码生成** | 基于项目知识生成集成代码 | 📋 计划中 |
| **运行时诊断** | 结合实时指标提供优化建议 | 📋 计划中 |

## 项目文档

| 文档 | 内容 | 路径 |
|------|------|------|
| **主任务计划** | 整体规划、目标架构、阶段安排 | [TASK-PLAN-MASTER.md](./TASK-PLAN-MASTER.md) |
| **OTel分析** | OpenTelemetry Collector深度分析 | [01-OTEL-COLLECTOR-ANALYSIS.md](./01-OTEL-COLLECTOR-ANALYSIS.md) |
| **知识图谱Schema** | 实体定义、关系类型、查询规范 | [02-KNOWLEDGE-GRAPH-SCHEMA.md](./02-KNOWLEDGE-GRAPH-SCHEMA.md) |
| **实现路线图** | 详细实施计划、里程碑 | [03-IMPLEMENTATION-ROADMAP.md](./03-IMPLEMENTATION-ROADMAP.md) |

## 快速开始

### 环境准备

```bash
# 1. 启动知识图谱数据库
docker-compose up -d neo4j weaviate

# 2. 克隆目标项目
git clone https://github.com/open-telemetry/opentelemetry-collector.git

# 3. 运行分析工具
go run ./cmd/analyze --project=./opentelemetry-collector --output=./data

# 4. 导入知识图谱
go run ./cmd/import --input=./data

# 5. 启动MCP服务
go run ./cmd/mcp-server
```

### 使用MCP客户端查询

```python
import mcp

client = mcp.Client("http://localhost:8080/mcp")

# 查询项目信息
result = client.call("query_project", {"name": "OpenTelemetry Collector"})
print(result)

# 搜索组件
result = client.call("search_components", {
    "capability": "batch",
    "signal": "traces"
})
print(result)

# 生成配置
result = client.call("generate_config", {
    "use_case": "high-throughput-traces",
    "constraints": {"memory_limit": "2GB"}
})
print(result)
```

## 项目分析范围

### Tier 1 - 核心基础设施

- [ ] OpenTelemetry Collector
- [ ] Go Runtime 1.26
- [ ] Prometheus
- [ ] Cilium
- [ ] Istio

### Tier 2 - 应用框架

- [ ] gRPC-Go
- [ ] Gin
- [ ] Ent
- [ ] Temporal

### Tier 3 - 运维平台

- [ ] Kubernetes
- [ ] Flagd
- [ ] Keptn

## 架构概览

```
┌─────────────────────────────────────────────────────────────────┐
│                         AI 交互层                                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐  │
│  │ 自然语言查询 │  │ 代码生成   │  │ 智能诊断与优化建议       │  │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    MCP Server (Model Context Protocol)          │
│              query_project │ search │ generate │ diagnose       │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    知识图谱服务层                                │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  Neo4j (图存储)        │  Weaviate (向量检索)           │    │
│  │  • 组件依赖关系        │  • 语义相似度搜索              │    │
│  │  • 配置-指标关联       │  • 自然语言匹配                │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    数据采集层                                    │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐  │
│  │ 静态分析   │  │ 运行时采集 │  │ LLM语义增强            │  │
│  │ • AST解析 │  │ • OTel指标 │  │ • 自动生成描述         │  │
│  │ • 依赖图  │  │ • pprof    │  │ • 使用场景推断         │  │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

## 技术栈

- **知识存储**: Neo4j (图) + Weaviate (向量)
- **静态分析**: Go AST + go/analysis
- **AI集成**: LangChain-Go + OpenAI/Anthropic
- **API协议**: Model Context Protocol (MCP)
- **可观测性**: OpenTelemetry (自举)

## 贡献指南

### 分析新项目

1. 在 `projects/` 下创建项目配置
2. 实现项目特定的分析器
3. 运行验证测试
4. 提交PR

### 扩展Schema

1. 在 `schema/` 下修改本体定义
2. 编写迁移脚本
3. 更新文档
4. 提交PR

## 路线图

| 阶段 | 时间 | 目标 |
|------|------|------|
| Phase 1 | Week 1-2 | 技术验证、基础架构 |
| Phase 2 | Week 3-8 | 核心项目分析 |
| Phase 3 | Week 9-12 | AI能力集成 |
| Phase 4 | 持续 | 生态扩展 |

详见 [实现路线图](./03-IMPLEMENTATION-ROADMAP.md)

## 社区

- 💬 讨论: GitHub Discussions
- 🐛 问题: GitHub Issues
- 📧 联系: [邮箱]

## License

MIT License

---

*让AI真正理解开源项目，解放工程师的生产力。*

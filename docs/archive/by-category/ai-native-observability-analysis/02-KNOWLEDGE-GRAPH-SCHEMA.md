# 知识图谱Schema详细设计

> 面向AI友好的开源项目知识表示规范

---

## 一、设计原则

### 1.1 核心原则

| 原则 | 说明 | 实现方式 |
|------|------|----------|
| **语义丰富** | 每个实体都有AI可理解的描述 | 每个节点包含semantic_description |
| **关系明确** | 关系类型标准化，有明确语义 | 预定义关系类型枚举 |
| **可验证** | 知识可通过代码/运行时验证 | 关联源码位置、测试用例 |
| **可扩展** | 支持新项目和自定义属性 | 灵活的属性Schema |

### 1.2 命名规范

```
项目命名: {owner}-{repo}-{component}
示例: open-telemetry-opentelemetry-collector-processor-batch

接口命名: {project}:{package}.{Interface}.{Method}
示例: otelcol:processor.BatchProcessor.ConsumeTraces

配置命名: {component}:{config_path}
示例: processor-batch:timeout

指标命名: {component}:metric.{name}
示例: processor-batch:metric.batch_send_size
```

---

## 二、实体类型详解

### 2.1 Project（项目）

```cypher
// Neo4j Cypher 创建语句示例
CREATE (p:Project {
  id: "open-telemetry-opentelemetry-collector",
  name: "OpenTelemetry Collector",
  short_name: "otelcol",

  // 基本信息
  description: "Vendor-agnostic way to receive, process and export telemetry data",
  description_zh: "与供应商无关的遥测数据接收、处理和导出方式",

  // 仓库信息
  repository: "https://github.com/open-telemetry/opentelemetry-collector",
  license: "Apache-2.0",
  primary_language: "Go",

  // 版本信息
  current_version: "v0.96.0",
  versioning_scheme: "semver",

  // 分类标签
  categories: ["observability", "telemetry", "collector"],
  capabilities: ["traces", "metrics", "logs"],

  // AI语义
  ai_summary: "可观测性数据管道基础设施，支持接收OTLP数据、处理后导出到多种后端",
  typical_use_cases: [
    "Kubernetes集群遥测收集",
    "多租户数据路由",
    "数据格式转换"
  ],

  // 元数据
  stars: 4500,
  created_at: "2019-01-01T00:00:00Z",
  last_updated: "2024-03-01T00:00:00Z"
})
```

**查询示例**：

```cypher
// 查找支持metrics和traces的Go项目
MATCH (p:Project)
WHERE "Go" IN p.languages
  AND "metrics" IN p.capabilities
  AND "traces" IN p.capabilities
RETURN p.name, p.ai_summary
```

### 2.2 Component（组件）

```cypher
CREATE (c:Component {
  id: "otelcol-processor-batch",
  name: "batch processor",
  type: "processor",  // receiver|processor|exporter|extension|connector

  // 功能描述
  purpose: "批量处理遥测数据以提高传输效率",
  when_to_use: "当需要减少网络调用次数、提高吞吐量时",
  when_not_to_use: "当需要低延迟、实时处理时",

  // 信号类型支持
  supported_signals: ["traces", "metrics", "logs"],

  // AI语义
  semantic_description: "累积多条记录后批量发送，平衡延迟与吞吐量",
  algorithm: "时间窗口 + 大小阈值触发",

  // 性能特征
  performance_characteristics: {
    throughput: "high",
    latency: "medium",
    memory_usage: "medium",
    cpu_usage: "low"
  },

  // 源码位置
  source_location: "processor/batchprocessor/batch_processor.go",

  // 状态
  stability: "stable",  // deprecated|experimental|stable
  distributions: ["core", "contrib"]
})
```

**关系建立**：

```cypher
MATCH (p:Project {id: "open-telemetry-opentelemetry-collector"})
MATCH (c:Component {id: "otelcol-processor-batch"})
CREATE (c)-[:BELONGS_TO]->(p)
```

### 2.3 Interface（接口）

```cypher
CREATE (i:Interface {
  id: "otelcol:processor.TracesProcessor.ConsumeTraces",
  name: "ConsumeTraces",
  type: "method",

  // 签名
  signature: {
    receiver: "TracesProcessor",
    method: "ConsumeTraces",
    parameters: [
      {name: "ctx", type: "context.Context", description: "上下文，用于取消和传播trace"},
      {name: "td", type: "ptrace.Traces", description: "要处理的trace数据"}
    ],
    returns: [
      {type: "error", description: "处理错误，nil表示成功"}
    ]
  },

  // AI语义
  semantic_description: "接收trace数据并进行处理，是Processor的核心入口",
  contract: [
    "不应修改输入数据的ownership",
    "错误应包装上下文信息",
    "应尊重context的取消信号"
  ],

  // 示例代码
  example_usage: '''
func (bp *batchProcessor) ConsumeTraces(ctx context.Context, td ptrace.Traces) error {
    bp.batch.Add(td)
    if bp.batch.ShouldExport() {
        return bp.batch.Export(ctx)
    }
    return nil
}
  ''',

  // 性能特征
  time_complexity: "O(n)",
  space_complexity: "O(n)",
  is_blocking: true,
  concurrency_safe: true
})
```

### 2.4 Configuration（配置）

```cypher
CREATE (cfg:Configuration {
  id: "otelcol:processor-batch:send_batch_size",
  key: "send_batch_size",
  component: "otelcol-processor-batch",

  // 类型信息
  value_type: "int",
  go_type: "int",
  schema: {
    minimum: 1,
    maximum: null,
    default: 8192
  },

  // 文档
  description: "在触发导出前批量中的最小span/metric/log记录数",
  description_zh: "触发导出前等待的最小记录数",

  // AI语义
  semantic_description: "控制批处理的粒度，较大值提高吞吐量但增加延迟",
  tuning_guidance: '''
调优建议：
- 高吞吐量场景：设为8192-16384
- 低延迟场景：设为1024-4096
- 内存受限：设为较小的值
  ''',

  // 影响分析
  affects: [
    {
      metric: "processor_batch_batch_send_size",
      relationship: "determines_maximum"
    },
    {
      metric: "memory_usage",
      relationship: "positive_correlation"
    },
    {
      metric: "latency",
      relationship: "negative_correlation"
    }
  ],

  // 验证规则
  validation_rules: [
    "必须为正整数",
    "应小于send_batch_max_size"
  ],

  // 运行时可调性
  hot_reloadable: false,  // 当前版本不支持
  requires_restart: true
})
```

### 2.5 RuntimeBehavior（运行时行为）

```cypher
CREATE (rb:RuntimeBehavior {
  id: "otelcol:processor-batch:metric.batch_send_size",
  name: "batch_send_size",
  full_name: "processor_batch_batch_send_size",

  // 指标类型
  metric_type: "histogram",  // counter|gauge|histogram|summary
  unit: "1",  // 无单位

  // 描述
  description: "实际发送的批次大小分布",
  semantic_description: "反映批处理效率，用于判断send_batch_size配置是否合理",

  // 标签维度
  labels: [
    {name: "processor", description: "处理器ID"},
    {name: "signal", description: "信号类型：traces/metrics/logs"}
  ],

  // 解释性指导
  interpretation: {
    normal_range: "send_batch_size的50%-100%",
    warning_signs: [
      "经常小于50%：timeout可能过短或流量不足",
      "经常等于100%：可能需要增加send_batch_size"
    ],
    action_items: [
      "监控p99值判断尾部延迟",
      "结合timeout_trigger_send比率分析"
    ]
  },

  // 与配置的关系
  configuration_drivers: [
    "send_batch_size",
    "timeout",
    "send_batch_max_size"
  ]
})
```

### 2.6 UseCase（使用场景）

```cypher
CREATE (uc:UseCase {
  id: "otelcol:high-throughput-trace-ingestion",
  name: "高吞吐量Trace数据摄入",
  domain: "observability",

  // 场景描述
  scenario: "Kubernetes集群产生大量trace数据，需要高效收集并导出到Jaeger",
  requirements: [
    "支持每秒10万+ span",
    "内存占用可控",
    "数据不丢失"
  ],

  // 推荐的组件组合
  recommended_pipeline: {
    receivers: ["otlp"],
    processors: ["memory_limiter", "batch"],
    exporters: ["otlp"]
  },

  // 关键配置
  key_configurations: {
    "batch.send_batch_size": 16384,
    "batch.timeout": "500ms",
    "otlp.sending_queue.num_consumers": 20
  },

  // 成功指标
  success_criteria: {
    throughput: ">100k spans/sec",
    memory: "<2GB",
    drop_rate: "<0.1%"
  }
})
```

---

## 三、关系类型详解

### 3.1 核心关系

| 关系 | 方向 | 语义 | 示例 |
|------|------|------|------|
| **BELONGS_TO** | Component -> Project | 组件属于项目 | batch processor -> OTel Collector |
| **IMPLEMENTS** | Component -> Interface | 实现接口 | batch processor -> ConsumeTraces |
| **DEPENDS_ON** | Component -> Component | 依赖关系 | exporter -> processor |
| **CONFIGURED_BY** | Component -> Configuration | 由配置控制 | batch -> send_batch_size |
| **EMITS** | Component -> RuntimeBehavior | 产生指标 | batch -> batch_send_size |
| **AFFECTS** | Configuration -> RuntimeBehavior | 配置影响指标 | send_batch_size -> batch_send_size |
| **USES_FOR** | Interface -> UseCase | 用于场景 | ConsumeTraces -> 高吞吐量摄入 |
| **SIMILAR_TO** | Component -> Component | 相似组件 | batch -> groupbyattrs |

### 3.2 关系创建示例

```cypher
// 配置影响指标
MATCH (cfg:Configuration {id: "otelcol:processor-batch:send_batch_size"})
MATCH (rb:RuntimeBehavior {id: "otelcol:processor-batch:metric.batch_send_size"})
CREATE (cfg)-[:AFFECTS {
  relationship_type: "determines_maximum",
  correlation: "direct",
  strength: "strong"
}]->(rb)

// 相似组件推荐
MATCH (c1:Component {id: "otelcol-processor-batch"})
MATCH (c2:Component {id: "otelcol-processor-groupbyattrs"})
CREATE (c1)-[:SIMILAR_TO {
  reason: "都用于数据聚合",
  difference: "batch按数量聚合，groupbyattrs按属性聚合"
}]->(c2)
```

---

## 四、向量索引设计

### 4.1 语义搜索支持

```cypher
// 为实体创建向量嵌入
CREATE (p:Project {
  id: "otelcol",
  name: "OpenTelemetry Collector"
})

// 向量属性存储（通过Weaviate或Neo4j GDS）
// embedding 由LLM生成
```

### 4.2 搜索场景

```python
# 语义搜索示例
async def search_by_natural_language(query: str):
    """
    将自然语言查询转换为向量搜索

    示例查询：
    - "如何批量处理trace数据？"
    - " Collector内存满了怎么办？"
    - "推荐一个支持动态配置的处理器"
    """

    # 1. 生成查询向量
    query_embedding = await llm.embed(query)

    # 2. 向量检索
    results = await weaviate.search(
        class_name="Component",
        vector=query_embedding,
        limit=5
    )

    # 3. 知识图谱增强
    for result in results:
        # 获取相关配置和指标
        related = await graph.query("""
            MATCH (c:Component {id: $id})-[:CONFIGURED_BY]->(cfg:Configuration)
            MATCH (c)-[:EMITS]->(rb:RuntimeBehavior)
            RETURN cfg, rb
        """, {"id": result.id})

        result.enrich(related)

    return results
```

---

## 五、Schema演进管理

### 5.1 版本控制

```yaml
# schema_version.yaml
version: "1.0.0"
last_updated: "2024-03-15"

changes:
  - version: "1.0.0"
    date: "2024-03-15"
    changes:
      - "初始Schema定义"
      - "支持Project/Component/Interface/Configuration/RuntimeBehavior"

  - version: "1.1.0"  # 计划中
    planned_date: "2024-04-01"
    changes:
      - "添加PerformanceProfile实体"
      - "添加VersionCompatibility关系"
```

### 5.2 迁移脚本

```cypher
// 示例：添加新属性的迁移
MATCH (c:Component)
WHERE c.stability IS NULL
SET c.stability = "unknown"
```

---

*本Schema将作为所有项目知识图谱构建的规范依据。*

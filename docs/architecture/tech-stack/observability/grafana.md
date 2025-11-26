# 1. 📊 Grafana 深度解析

> **简介**: 本文档详细阐述了 Grafana 的核心特性、选型论证、实际应用和最佳实践。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [1. 📊 Grafana 深度解析](#1--grafana-深度解析)
  - [📋 目录](#-目录)
  - [1.1 核心特性](#11-核心特性)
  - [1.2 选型论证](#12-选型论证)
  - [1.3 实际应用](#13-实际应用)
    - [1.3.1 数据源配置](#131-数据源配置)
    - [1.3.2 仪表板创建](#132-仪表板创建)
    - [1.3.3 查询编写](#133-查询编写)
    - [1.3.4 告警配置](#134-告警配置)
  - [1.4 最佳实践](#14-最佳实践)
    - [1.4.1 仪表板设计最佳实践](#141-仪表板设计最佳实践)
  - [📚 扩展阅读](#-扩展阅读)

---

## 1.1 核心特性

**Grafana 是什么？**

Grafana 是一个开源的可视化平台，用于监控和数据分析。

**核心特性**:

- ✅ **可视化**: 丰富的可视化图表
- ✅ **数据源**: 支持多种数据源
- ✅ **仪表板**: 灵活的仪表板配置
- ✅ **告警**: 支持告警通知

---

## 1.2 选型论证

**为什么选择 Grafana？**

**论证矩阵**:

| 评估维度 | 权重 | Grafana | Kibana | DataDog | New Relic | 说明 |
|---------|------|---------|--------|---------|-----------|------|
| **可视化能力** | 30% | 10 | 8 | 9 | 9 | Grafana 可视化最强大 |
| **数据源支持** | 25% | 10 | 7 | 8 | 7 | Grafana 支持最多数据源 |
| **Prometheus 集成** | 20% | 10 | 6 | 8 | 7 | Grafana 与 Prometheus 集成最好 |
| **社区生态** | 15% | 10 | 8 | 6 | 6 | Grafana 社区最活跃 |
| **成本** | 10% | 10 | 9 | 3 | 3 | Grafana 开源免费 |
| **加权总分** | - | **9.80** | 7.50 | 7.20 | 7.00 | Grafana 得分最高 |

**核心优势**:

1. **可视化能力（权重 30%）**:
   - 丰富的图表类型，支持自定义面板
   - 灵活的仪表板配置，支持变量和模板
   - 实时数据刷新，支持流式数据

2. **数据源支持（权重 25%）**:
   - 支持 100+ 数据源，包括 Prometheus、InfluxDB、MySQL 等
   - 统一的数据源接口，易于扩展
   - 支持多数据源联合查询

3. **Prometheus 集成（权重 20%）**:
   - 原生支持 Prometheus，查询性能优秀
   - 支持 PromQL，查询语法一致
   - 与 Prometheus 生态完美集成

**为什么不选择其他可视化工具？**

1. **Kibana**:
   - ✅ 与 Elasticsearch 集成好
   - ❌ 主要面向日志分析，监控功能有限
   - ❌ 与 Prometheus 集成不如 Grafana
   - ❌ 可视化能力不如 Grafana

2. **DataDog**:
   - ✅ 功能完善，SaaS 服务
   - ❌ 成本高，不适合中小型项目
   - ❌ 数据存储在第三方
   - ❌ 依赖外部服务

3. **New Relic**:
   - ✅ APM 功能强大
   - ❌ 成本高
   - ❌ 数据存储在第三方
   - ❌ 与 Prometheus 集成不如 Grafana

---

## 1.3 实际应用

### 1.3.1 数据源配置

**Prometheus 数据源配置**:

```json
{
  "name": "Prometheus",
  "type": "prometheus",
  "url": "http://prometheus:9090",
  "access": "proxy",
  "isDefault": true,
  "jsonData": {
    "timeInterval": "15s"
  }
}
```

### 1.3.2 仪表板创建

**创建仪表板示例**:

```json
{
  "dashboard": {
    "title": "Golang Service Dashboard",
    "panels": [
      {
        "title": "HTTP Requests",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{method}} {{path}}"
          }
        ]
      },
      {
        "title": "Request Duration",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "95th percentile"
          }
        ]
      }
    ]
  }
}
```

### 1.3.3 查询编写

**PromQL 查询示例**:

```promql
# 请求速率
rate(http_requests_total[5m])

# 错误率
rate(http_requests_total{status="error"}[5m]) / rate(http_requests_total[5m])

# 95 分位延迟
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# 活跃连接数
active_connections
```

### 1.3.4 告警配置

**告警配置示例**:

```json
{
  "alert": {
    "name": "High Error Rate",
    "message": "Error rate is above threshold",
    "conditions": [
      {
        "evaluator": {
          "params": [0.05],
          "type": "gt"
        },
        "query": {
          "params": ["A", "5m", "now"]
        },
        "reducer": {
          "type": "avg"
        }
      }
    ],
    "notifications": [
      {
        "uid": "slack"
      }
    ]
  }
}
```

---

## 1.4 最佳实践

### 1.4.1 仪表板设计最佳实践

**为什么需要良好的仪表板设计？**

良好的仪表板设计可以提高监控效率，便于快速发现问题。

**仪表板设计原则**:

1. **信息层次**: 按照重要性组织信息
2. **可视化类型**: 根据数据类型选择合适的可视化类型
3. **时间范围**: 设置合理的时间范围
4. **告警集成**: 集成告警信息

**实际应用示例**:

```json
{
  "dashboard": {
    "title": "Service Overview",
    "panels": [
      {
        "title": "Request Rate",
        "type": "graph",
        "gridPos": {"x": 0, "y": 0, "w": 12, "h": 8},
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{method}} {{path}}"
          }
        ]
      },
      {
        "title": "Error Rate",
        "type": "graph",
        "gridPos": {"x": 12, "y": 0, "w": 12, "h": 8},
        "targets": [
          {
            "expr": "rate(http_requests_total{status=\"error\"}[5m]) / rate(http_requests_total[5m])",
            "legendFormat": "Error Rate"
          }
        ]
      }
    ]
  }
}
```

**最佳实践要点**:

1. **信息层次**: 将最重要的指标放在顶部
2. **可视化类型**: 根据数据类型选择合适的可视化类型（折线图、柱状图、饼图等）
3. **时间范围**: 设置合理的时间范围，便于查看趋势
4. **告警集成**: 在仪表板上显示告警信息，便于快速发现问题

---

## 📚 扩展阅读

- [Grafana 官方文档](https://grafana.com/docs/)
- [技术栈概览](../00-技术栈概览.md)
- [技术栈集成](../01-技术栈集成.md)
- [技术栈选型决策树](../02-技术栈选型决策树.md)

---

> 📚 **简介**
> 本文档提供了 Grafana 的完整解析，包括核心特性、选型论证、实际应用和最佳实践。

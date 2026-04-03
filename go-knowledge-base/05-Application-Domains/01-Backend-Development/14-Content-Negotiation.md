# 内容协商 (Content Negotiation)

> **分类**: 成熟应用领域  
> **标签**: #content-negotiation #api #rest

---

## Accept Header 解析

```go
func ParseAccept(header string) []MediaType {
    var types []MediaType
    
    for _, part := range strings.Split(header, ",") {
        part = strings.TrimSpace(part)
        
        mediaType, q := part, 1.0
        if idx := strings.Index(part, ";"); idx != -1 {
            mediaType = strings.TrimSpace(part[:idx])
            params := part[idx+1:]
            
            // 解析 q 值
            if strings.Contains(params, "q=") {
                qStr := strings.TrimPrefix(params[strings.Index(params, "q="):], "q=")
                q, _ = strconv.ParseFloat(strings.TrimSpace(qStr), 64)
            }
        }
        
        types = append(types, MediaType{
            Type:    mediaType,
            Quality: q,
        })
    }
    
    // 按 q 值排序
    sort.Slice(types, func(i, j int) bool {
        return types[i].Quality > types[j].Quality
    })
    
    return types
}
```

---

## 响应格式协商

```go
func ContentNegotiation() gin.HandlerFunc {
    return func(c *gin.Context) {
        accept := c.GetHeader("Accept")
        if accept == "" {
            accept = "application/json"  // 默认
        }
        
        mediaTypes := ParseAccept(accept)
        
        // 查找最佳匹配
        supported := []string{"application/json", "application/xml", "text/csv"}
        
        for _, mt := range mediaTypes {
            for _, s := range supported {
                if match(mt.Type, s) {
                    c.Set("response_format", s)
                    c.Next()
                    return
                }
            }
        }
        
        c.AbortWithStatus(406)  // Not Acceptable
    }
}

func match(accept, supported string) bool {
    // application/json 匹配 application/*
    if strings.HasSuffix(accept, "/*") {
        prefix := strings.TrimSuffix(accept, "/*")
        return strings.HasPrefix(supported, prefix)
    }
    return accept == supported
}
```

---

## 多格式响应

```go
func RespondWithFormat(c *gin.Context, data interface{}) {
    format := c.GetString("response_format")
    
    switch format {
    case "application/xml":
        c.XML(200, data)
    case "text/csv":
        respondCSV(c, data)
    default:
        c.JSON(200, data)
    }
}

func respondCSV(c *gin.Context, data interface{}) {
    c.Header("Content-Type", "text/csv")
    c.Header("Content-Disposition", "attachment; filename=data.csv")
    
    writer := csv.NewWriter(c.Writer)
    defer writer.Flush()
    
    // 写入数据...
}
```

---

## 语言协商

```go
func LanguageNegotiation(supported []string) gin.HandlerFunc {
    return func(c *gin.Context) {
        acceptLang := c.GetHeader("Accept-Language")
        if acceptLang == "" {
            c.Set("lang", supported[0])
            c.Next()
            return
        }
        
        // 解析语言偏好
        langs := parseLanguages(acceptLang)
        
        // 查找匹配
        for _, l := range langs {
            for _, s := range supported {
                if matchLanguage(l.Tag, s) {
                    c.Set("lang", s)
                    c.Next()
                    return
                }
            }
        }
        
        c.Set("lang", supported[0])
        c.Next()
    }
}

func matchLanguage(accepted, supported string) bool {
    // en-US 匹配 en
    if strings.HasPrefix(accepted, supported) {
        return true
    }
    return accepted == supported
}
```

---

## 架构决策记录

### 决策矩阵

| 方案 | 优点 | 缺点 | 适用场景 |
|------|------|------|----------|
| A | 高性能 | 复杂 | 大规模 |
| B | 简单 | 扩展性差 | 小规模 |

### 风险评估

**风险 R.1**: 性能瓶颈
- 概率: 中
- 影响: 高
- 缓解: 缓存、分片

**风险 R.2**: 单点故障
- 概率: 低
- 影响: 极高
- 缓解: 冗余、故障转移

### 实施路线图

`
Phase 1: 基础设施 (Week 1-2)
Phase 2: 核心功能 (Week 3-6)
Phase 3: 优化加固 (Week 7-8)
`

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 架构决策记录 (ADR)

### 上下文

业务需求和技术约束分析。

### 决策

选择方案A作为主要架构方向。

### 后果

正面：
- 可扩展性提升
- 维护成本降低

负面：
- 初期开发复杂度增加
- 团队学习成本

### 实施指南

`
Week 1-2: 基础设施搭建
Week 3-4: 核心功能开发
Week 5-6: 集成测试
Week 7-8: 性能优化
`

### 风险评估

| 风险 | 概率 | 影响 | 缓解措施 |
|------|------|------|----------|
| 性能不足 | 中 | 高 | 缓存、分片 |
| 兼容性 | 低 | 中 | 接口适配层 |

### 监控指标

- 系统吞吐量
- 响应延迟
- 错误率
- 资源利用率

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 系统设计

### 需求分析

功能需求和非功能需求的完整梳理。

### 架构视图

`
┌─────────────────────────────────────┐
│           API Gateway               │
└─────────────┬───────────────────────┘
              │
    ┌─────────┴─────────┐
    ▼                   ▼
┌─────────┐       ┌─────────┐
│ Service │       │ Service │
│   A     │       │   B     │
└────┬────┘       └────┬────┘
     │                 │
     └────────┬────────┘
              ▼
        ┌─────────┐
        │  Data   │
        │  Store  │
        └─────────┘
`

### 技术选型

| 组件 | 技术 | 理由 |
|------|------|------|
| API | gRPC | 性能 |
| DB | PostgreSQL | 可靠 |
| Cache | Redis | 速度 |
| Queue | Kafka | 吞吐 |

### 性能指标

- QPS: 10K+
- P99 Latency: <100ms
- Availability: 99.99%

### 运维手册

- 部署流程
- 监控配置
- 应急预案
- 容量规划

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 扩展分析

### 理论基础

深入探讨相关理论概念和数学基础。

### 实现细节

完整的代码实现和配置示例。

### 最佳实践

- 设计原则
- 编码规范
- 测试策略
- 部署流程

### 性能优化

| 技术 | 效果 | 复杂度 |
|------|------|--------|
| 缓存 | 10x | 低 |
| 批处理 | 5x | 中 |
| 异步 | 3x | 中 |

### 常见问题

Q: 如何处理高并发？
A: 使用连接池、限流、熔断等模式。

### 相关资源

- 官方文档
- 学术论文
- 开源项目

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
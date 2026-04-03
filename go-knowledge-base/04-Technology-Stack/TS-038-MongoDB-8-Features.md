# TS-038-MongoDB-8-Features

> **Dimension**: 04-Technology-Stack
> **Status**: S-Level
> **Created**: 2026-04-03
> **Version**: MongoDB 8.0 (2025)
> **Size**: >20KB

---

## 1. MongoDB 8.0 概览

### 1.1 发布信息

- **发布日期**: 2025年5月
- **版本类型**: Major Release
- **主要焦点**: 性能、可扩展性、AI/ML

### 1.2 性能提升概览

| 操作类型 | 性能提升 |
|---------|---------|
| 读操作 | 20-36% |
| 写操作 | 35-58% |
| 批量写入 | 56% |
| 时间序列查询 | 200%+ |
| 分片数据分布 | 50倍更快 |

---

## 2. 核心性能优化

### 2.1 升级版TCMalloc

**改进**:

- 每CPU缓存(替代每线程缓存)
- 减少内存碎片
- 更好的高压力工作负载适应性

**内存释放**:

```javascript
// 后台线程每秒释放内存
// 默认释放速率: 10MB/s

db.adminCommand({
    setParameter: 1,
    tcmallocReleaseRate: 10485760  // 10MB/s
});
```

**监控指标**:

```javascript
// 查看TCMalloc统计
db.serverStatus().tcmalloc

// 关键指标:
// - usingPerCPUCaches
// - peak_memory_usage
// - total_free_bytes
```

### 2.2 批量写入优化

**bulkWrite命令**:

```javascript
// MongoDB 8.0新命令
// 单请求多集合操作
db.adminCommand({
    bulkWrite: 1,
    operations: [
        { insert: 0, document: { _id: 1, name: "A" } },
        { update: 1, filter: { _id: 2 }, updateMods: { $set: { name: "B" } } },
        { delete: 2, filter: { _id: 3 } }
    ],
    nsInfo: [
        { ns: "db.coll1" },
        { ns: "db.coll2" },
        { ns: "db.coll3" }
    ]
});

// 比7.0快56%
```

### 2.3 时间序列块处理

**块处理(Chunk Processing)**:

```javascript
// 使用时间序列集合
// 自动启用块处理

db.createCollection("metrics", {
    timeseries: {
        timeField: "timestamp",
        metaField: "metadata",
        granularity: "minutes"
    }
});

// 查询使用块处理
// 吞吐量提升200%+
```

**查看执行计划**:

```javascript
db.metrics.explain("executionStats").aggregate([
    { $group: { _id: "$metadata.sensor", avg: { $avg: "$value" } } }
]);

// 检查: queryPlanner.winningPlan.slotBasedPlan.stages
// 查看是否使用块处理
```

---

## 3. Express查询阶段

### 3.1 概述

针对简单查询优化，跳过常规查询计划。

**适用场景**:

- 单_id索引查询
- 简单点查

### 3.2 使用示例

```javascript
// 简单_id查询
db.customer.find({_id: ObjectId('670ec6b005b98857588f5b6a')}).explain()

// 执行计划显示: EXPRESS_IXSCAN
// 性能提升: 17%
```

**Express阶段类型**:

- `EXPRESS_CLUSTERED_IXSCAN`
- `EXPRESS_DELETE`
- `EXPRESS_IXSCAN`
- `EXPRESS_UPDATE`

---

## 4. 可查询加密增强

### 4.1 范围查询支持

```javascript
// 加密字段范围查询
db.users.find({
    "ssn": {
        $gt: "100-00-0000",
        $lt: "999-99-9999"
    }
});

// 支持操作符: $lt, $lte, $gt, $gte
```

### 4.2 Java集成示例

```java
// 可查询加密范围查询
Bson filter = Filters.and(
    Filters.gt("age", 18),
    Filters.lt("age", 65)
);

FindIterable<Document> results = collection.find(filter);
```

---

## 5. 分片增强

### 5.1 分片性能提升

- 数据分布速度: **50倍更快**
- 成本降低: **50%**

### 5.2 Reshard性能优化

```javascript
// 重新分片命令优化
// 更快完成数据重分布
db.adminCommand({
    reshardCollection: "db.largeCollection",
    key: { newShardKey: 1 }
});
```

### 5.3 配置分片(Config Shard)

```javascript
// 配置服务器存储应用数据
// 节省资源

// 在配置集群上启用
db.adminCommand({
    setClusterParameter: {
        configShard: { enabled: true }
    }
});
```

---

## 6. 查询设置(Query Settings)

### 6.1 概述

替代索引过滤器，提供更精细的查询控制。

### 6.2 设置查询配置

```javascript
// 添加查询设置
db.adminCommand({
    setQuerySettings: {
        queryShape: {
            cmdNs: { db: "test", coll: "coll" },
            command: "find",
            filter: { category: "electronics" }
        }
    },
    settings: {
        indexHints: [{ category: 1 }],
        reject: false
    }
});

// 限流: 拒绝特定查询形状
/*
db.adminCommand({
    setQuerySettings: { ... },
    settings: { reject: true }
});
*/
```

### 6.3 查看查询设置

```javascript
// 查看当前设置
db.aggregate([{ $querySettings: {} }]);

// 删除设置
db.adminCommand({
    removeQuerySettings: { queryShape: { ... } }
});
```

---

## 7. 索引构建改进

### 7.1 快速错误检测

```
行为变化:
- 8.0: 收集扫描阶段发现错误立即返回
- <8.0: 错误在提交阶段返回

优势: 快速诊断索引错误
```

### 7.2 弹性部署

- 辅助成员可请求主成员停止索引构建
- 避免辅助成员崩溃

### 7.3 磁盘空间管理

```javascript
// 设置最小可用磁盘空间
db.adminCommand({
    setParameter: 1,
    indexBuildMinAvailableDiskSpaceMB: 1000
});

// 低于阈值自动停止索引构建
```

### 7.4 后台压缩

```javascript
// 自动后台压缩
db.adminCommand({
    autoCompact: 1,
    freeSpaceTargetMB: 100
});

// 保持空闲空间在指定值以下
```

---

## 8. 变更数据捕获(CDC)

### 8.1 减少重复消息

**改进**:

- 重启和重试期间显著减少重复消息
- 大表或不均匀处理速度场景受益

### 8.2 丰富变更流格式

```javascript
// 新的丰富格式 (类似Debezium)
db.collection.watch([], {
    fullDocument: "updateLookup",
    enrich: true  // 启用丰富格式
});

// 包含: 源信息、模式、操作类型、时间戳
```

### 8.3 变更流改进

```javascript
// 批量插入性能提升
// 非事务批量插入生成统一oplog
// 所有文档在Change Stream中具有相同clusterTime
```

---

## 9. 日志和监控

### 9.1 WorkingMillis字段

```javascript
// 慢日志新增字段
db.system.profile.find().pretty()

// workingMillis: 实际执行时间(不含锁等待)
// durationMillis: 总延迟(含等待)
```

### 9.2 查询统计

```javascript
// $queryStats聚合阶段
db.aggregate([{ $queryStats: {} }]);

// 返回已记录查询的统计数据
// 优化Change Stream中的跟踪
```

---

## 10. 聚合增强

### 10.1 BinData转换

```javascript
// 字符串与BinData互转
db.collection.aggregate([
    {
        $project: {
            // 字符串转BinData
            binary: { $convert: { input: "$hexString", to: "binData" } },
            // BinData转字符串
            string: { $convert: { input: "$binaryData", to: "string" } }
        }
    }
]);

// $toUUID: 简化字符串转UUID
```

### 10.2 updateOne增强

```javascript
// 支持sort选项
db.collection.updateOne(
    { status: "pending" },
    { $set: { status: "processing" } },
    { sort: { priority: -1 } }  // 先更新高优先级
);
```

---

## 11. 升级指南

### 11.1 兼容性注意事项

**undefined类型**:

```javascript
// 8.0前: null查询返回undefined值
// 8.0: null查询不返回undefined

db.people.find({name: null})
// 8.0前返回: [{name: null}, {name: undefined}]
// 8.0返回: [{name: null}]
```

**迁移建议**:

```javascript
// 查找undefined数据
db.people.find({name: {$type: "undefined"}})

// 迁移到null
db.people.updateMany(
    {name: {$type: "undefined"}},
    {$set: {name: null}}
);
```

### 11.2 索引过滤器弃用

```javascript
// 旧方式(已弃用)
db.collection.planCacheSetFilter(
    { query: { status: "active" } },
    { indexes: [{ status: 1, created: -1 }] }
);

// 新方式: setQuerySettings
// (见第6节)
```

---

## 12. 云数据库MongoDB版特性

### 12.1 阿里云MongoDB 8.0

**2025年发布功能**:

| 功能 | 描述 | 发布时间 |
|------|------|---------|
| MongoDB 8.0支持 | 读20-36%↑, 写35-58%↑ | 2025-05 |
| 全球多活数据库 | 跨地域多活架构 | 2025-08 |
| 单节点架构 | 4.2+版本支持 | 2025-07 |
| SQL限流 | 语句并发度限制 | 2025-12 |
| 索引优化推荐 | 基于慢日志的优化建议 | 2025-12 |

### 12.2 全球多活数据库

```
架构:
- 基于MongoDB高可用架构
- 与DTS无缝集成
- 异地灾备
- 就近访问
```

---

## 13. 参考文献

1. MongoDB 8.0 Release Notes
2. MongoDB Performance Best Practices
3. MongoDB Security Guide
4. MongoDB Time Series Collections
5. MongoDB Queryable Encryption

---

*Last Updated: 2026-04-03*

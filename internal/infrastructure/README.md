# Infrastructure Layer (基础设施层)

Clean Architecture 的基础设施层，包含技术实现。

## 结构

```
infrastructure/
├── database/      # 数据库实现
│   ├── postgres/  # PostgreSQL
│   └── ent/       # Ent ORM
├── messaging/     # 消息队列
│   ├── kafka/     # Kafka
│   └── mqtt/      # MQTT
├── cache/         # 缓存
└── observability/ # 可观测性
    ├── otlp/      # OpenTelemetry
    └── ebpf/      # eBPF
```

## 规则

- ✅ 实现 domain 层定义的接口
- ✅ 包含技术实现细节
- ✅ 可以导入外部库

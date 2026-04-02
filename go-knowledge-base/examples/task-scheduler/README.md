# 分布式任务调度器示例

完整可运行的分布式任务调度器实现。

## 快速开始

```bash
docker-compose up -d
make build
./bin/scheduler --config configs/scheduler.yaml
./bin/worker --id worker-1
./bin/cli submit --type email --payload '{"to":"user@example.com"}'
```

## 架构

- Go 1.26
- etcd (协调)
- PostgreSQL (状态)
- Redis (队列)
- gRPC/HTTP API

## 特性

- 领导者选举
- 任务调度 (轮询/最少任务/资源匹配)
- 任务执行与重试
- 死信队列
- 健康检查
- 优雅关闭

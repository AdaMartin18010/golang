# Observability Examples (OTel)

## 目录

- docker-compose.yaml：一键启动 Collector/Tempo/Prometheus/Loki/Grafana
- otelcol.yaml：Collector 最小管道配置
- app/main.go：最小 Go 服务，暴露 /hello 并发送 Trace，日志含 trace_id

## 运行

1. 启动后端组件

```bash
docker compose up -d
```

2. 运行应用

```bash
go run ./examples/observability/app
```

3. 访问接口

```bash
curl http://localhost:8080/hello
```

随后在 Grafana/Tempo 中查看 Trace，在 Prometheus/Loki 中查看指标与日志。



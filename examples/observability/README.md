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

   或在 Windows/PowerShell 使用脚本：

   ```powershell
   ./up.ps1
   ```

2. 运行应用

    ```bash
    go run ./examples/observability/app
    ```

   在 Windows/PowerShell 下，推荐使用脚本（会自动设置缺省 OTEL 环境变量并后台启动）：

   ```powershell
   ./run-app.ps1   # 若未设置，则默认：
                   #   OTEL_SERVICE_NAME=example-observability-app
                   #   OTEL_ENV=dev
                   #   OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
                   #   OTEL_EXPORTER_OTLP_INSECURE=true
   ```

3. 访问接口

    ```bash
    curl http://localhost:8080/hello
    curl -i http://localhost:8080/error  # 返回 500，用于错误率验证
    ```

随后在 Grafana/Tempo 中查看 Trace，在 Prometheus/Loki 中查看指标与日志。
Grafana 预置数据源（Prometheus/Tempo/Loki），RED 看板可查看请求率、错误率与 P99。

### Windows/PowerShell 提示

- 后台启动应用（不使用 `&`/`&&`）：

  ```powershell
  powershell -NoProfile -Command "Start-Process -FilePath 'go' -ArgumentList 'run ./app' -WorkingDirectory 'G:\_src\golang\golang\examples\observability' -WindowStyle Hidden"
  ```

- 或使用脚本：

  ```powershell
  ./run-app.ps1   # 启动
  ./stop-app.ps1  # 停止
  ```

- 如果未运行 Docker Desktop，应用日志可能出现 4317 连接失败（OTLP 导出）。这不影响本地接口验证，可先忽略；要完整链路与指标，请先启动 `docker compose up -d`。

提示：`run-app.ps1` 会在未设置端点时默认将导出器指向本机 Collector（`localhost:4317`，insecure）。

### 常见排障

- 端口占用：确保本地 `:8080` 未被占用。
- Collector 未启动：应用会报 `dial tcp 127.0.0.1:4317 refused`，属预期；启动 compose 后恢复。
- Grafana 登录：默认 `http://localhost:3000`，首登创建账号；已预置 Tempo/Prometheus/Loki 数据源与 RED 看板。

K8s（可选）：

```bash
kubectl apply -f k8s/otel-collector.yaml
kubectl apply -f k8s/tempo.yaml
```

说明：

- 已启用 Tail Sampling（错误与>200ms 慢请求更易被采样保留）
- 应用上报直方图 `http_server_request_duration_seconds`，在 Grafana RED 面板可查看 P99；Trace→Metrics 支持 exemplars 跳转

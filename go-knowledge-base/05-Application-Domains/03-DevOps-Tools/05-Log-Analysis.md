# 日志分析工具

> **分类**: 成熟应用领域

---

## 结构化日志

```go
import "go.uber.org/zap"

logger, _ := zap.NewProduction()
logger.Info("request",
    zap.String("method", "GET"),
    zap.Int("status", 200),
    zap.Duration("latency", time.Millisecond*45),
)
```

---

## 日志收集

### Fluent Bit

```go
// 发送日志到 Fluent Bit
conn, _ := net.Dial("tcp", "localhost:24224")
msg := fmt.Sprintf("[\"tag\", %d, {\"log\":\"%s\"}]\n", time.Now().Unix(), logData)
conn.Write([]byte(msg))
```

---

## 日志分析

```go
// 解析日志
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
    // 解析 JSON 日志
    var log LogEntry
    json.Unmarshal([]byte(line), &log)
}
```

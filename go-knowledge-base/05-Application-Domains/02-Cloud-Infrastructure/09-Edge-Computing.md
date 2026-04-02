# 边缘计算 (Edge Computing)

> **分类**: 成熟应用领域
> **标签**: #edge #iot #wasm

---

## WebAssembly (WASM)

### 编译到 WASM

```bash
# 编译为 WASM
GOOS=js GOARCH=wasm go build -o main.wasm

# 复制 JS 支持文件
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
```

### WASM 运行时

```go
import "github.com/tetratelabs/wazero"

ctx := context.Background()

// 创建运行时
r := wazero.NewRuntime(ctx)
defer r.Close(ctx)

// 加载 WASM 模块
mod, err := r.InstantiateFromPath(ctx, "plugin.wasm")
if err != nil {
    log.Fatal(err)
}

// 调用函数
add := mod.ExportedFunction("add")
result, err := add.Call(ctx, 1, 2)
```

---

## 边缘函数

### Cloudflare Workers

```go
package main

import (
    "github.com/syumai/workers"
    "github.com/syumai/workers/cloudflare"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        env := cloudflare.FromContext(ctx)

        // 访问 KV
        kv := env.KV("MY_KV")
        value, _ := kv.GetString(ctx, "key")

        w.Write([]byte(value))
    })

    workers.Serve(nil)
}
```

---

## IoT 设备通信

### MQTT 客户端

```go
import mqtt "github.com/eclipse/paho.mqtt.golang"

opts := mqtt.NewClientOptions()
opts.AddBroker("tcp://broker.hivemq.com:1883")
opts.SetClientID("go-client")

client := mqtt.NewClient(opts)
if token := client.Connect(); token.Wait() && token.Error() != nil {
    log.Fatal(token.Error())
}

// 订阅
token := client.Subscribe("sensors/temperature", 0, func(client mqtt.Client, msg mqtt.Message) {
    fmt.Printf("Received: %s from %s\n", msg.Payload(), msg.Topic())
})
token.Wait()

// 发布
token = client.Publish("sensors/temperature", 0, false, "25.5")
token.Wait()
```

---

## 边缘数据处理

```go
type EdgeProcessor struct {
    window time.Duration
    buffer []SensorData
    mu     sync.Mutex
}

func (p *EdgeProcessor) Process(data SensorData) {
    p.mu.Lock()
    defer p.mu.Unlock()

    p.buffer = append(p.buffer, data)

    // 窗口满，聚合发送
    if len(p.buffer) >= 100 {
        p.flush()
    }
}

func (p *EdgeProcessor) flush() {
    if len(p.buffer) == 0 {
        return
    }

    // 本地聚合
    avg := calculateAverage(p.buffer)
    max := findMax(p.buffer)

    // 只发送聚合结果到云端
    cloud.Send(AggregatedData{
        Average: avg,
        Max:     max,
        Count:   len(p.buffer),
    })

    p.buffer = p.buffer[:0]
}
```

---

## 离线优先

```go
type OfflineQueue struct {
    db     *sql.DB
    client *http.Client
}

func (q *OfflineQueue) Enqueue(data []byte) error {
    _, err := q.db.Exec("INSERT INTO queue (data, created_at) VALUES (?, ?)",
        data, time.Now())
    return err
}

func (q *OfflineQueue) Sync() error {
    rows, err := q.db.Query("SELECT id, data FROM queue ORDER BY created_at LIMIT 100")
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var id int
        var data []byte
        rows.Scan(&id, &data)

        resp, err := q.client.Post("https://api.example.com/data",
            "application/json", bytes.NewReader(data))

        if err == nil && resp.StatusCode == 200 {
            // 成功，删除记录
            q.db.Exec("DELETE FROM queue WHERE id = ?", id)
        }
    }

    return nil
}
```

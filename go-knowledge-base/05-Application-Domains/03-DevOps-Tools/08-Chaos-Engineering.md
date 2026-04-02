# 混沌工程 (Chaos Engineering)

> **分类**: 成熟应用领域
> **标签**: #chaos #reliability #testing

---

## 故障注入

### 网络延迟

```go
// 使用 toxiproxy
import "github.com/Shopify/toxiproxy/v2/client"

cli := toxiproxy.NewClient("localhost:8474")

// 创建代理
proxy, err := cli.CreateProxy("mysql", "localhost:3306", "mysql:3306")
if err != nil {
    log.Fatal(err)
}

// 添加延迟
_, err = proxy.AddToxic("latency_down", "latency", "downstream", 1.0, toxiproxy.Attributes{
    "latency": 1000,  // 1000ms 延迟
    "jitter":  100,
})
```

### HTTP 故障

```go
// 故障注入中间件
func ChaosMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 10% 概率返回错误
        if rand.Float32() < 0.1 {
            w.WriteHeader(http.StatusServiceUnavailable)
            return
        }

        // 5% 概率延迟
        if rand.Float32() < 0.05 {
            time.Sleep(5 * time.Second)
        }

        next.ServeHTTP(w, r)
    })
}
```

---

## 资源压力测试

### CPU 压力

```go
func CPULoad(ctx context.Context, cores int) {
    for i := 0; i < cores; i++ {
        go func() {
            for {
                select {
                case <-ctx.Done():
                    return
                default:
                    // 空循环消耗 CPU
                    for j := 0; j < 1000000; j++ {
                        _ = j * j
                    }
                }
            }
        }()
    }
}
```

### 内存压力

```go
func MemoryLoad(ctx context.Context, sizeMB int) {
    data := make([][]byte, 0)

    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            chunk := make([]byte, sizeMB*1024*1024)
            data = append(data, chunk)
        }
    }
}
```

---

## 自动化混沌测试

```go
// 定义实验
type ChaosExperiment struct {
    Name        string
    Duration    time.Duration
    Faults      []Fault
    AbortOnError bool
}

type Fault interface {
    Inject(ctx context.Context) error
    Recover() error
}

func RunExperiment(exp ChaosExperiment) error {
    ctx, cancel := context.WithTimeout(context.Background(), exp.Duration)
    defer cancel()

    for _, fault := range exp.Faults {
        if err := fault.Inject(ctx); err != nil {
            if exp.AbortOnError {
                return err
            }
            log.Printf("fault injection failed: %v", err)
        }
    }

    <-ctx.Done()

    // 清理
    for _, fault := range exp.Faults {
        fault.Recover()
    }

    return nil
}
```

# 性能剖析 (Profiling)

> **分类**: 工程与云原生

---

## CPU 分析

```go
import _ "net/http/pprof"

go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

```bash
# 采集 30 秒 CPU 数据
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

---

## 内存分析

```bash
# 堆内存
go tool pprof http://localhost:6060/debug/pprof/heap

# 分配
go tool pprof http://localhost:6060/debug/pprof/allocs
```

---

## 阻塞分析

```bash
go tool pprof http://localhost:6060/debug/pprof/block
```

---

## 可视化

```bash
# 生成火焰图
go tool pprof -http=:8080 profile.out

# 交互式命令
(pprof) top
(pprof) list functionName
(pprof) web
```

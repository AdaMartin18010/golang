# etcd 客户端

> **分类**: 开源技术堆栈

---

## 连接

```go
import clientv3 "go.etcd.io/etcd/client/v3"

cli, err := clientv3.New(clientv3.Config{
    Endpoints:   []string{"localhost:2379"},
    DialTimeout: 5 * time.Second,
})
```

---

## KV 操作

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Put
cli.Put(ctx, "foo", "bar")

// Get
resp, _ := cli.Get(ctx, "foo")
for _, ev := range resp.Kvs {
    fmt.Printf("%s : %s\n", ev.Key, ev.Value)
}

// Delete
cli.Delete(ctx, "foo")
```

---

## Watch

```go
rch := cli.Watch(ctx, "foo")
for wresp := range rch {
    for _, ev := range wresp.Events {
        fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
    }
}
```

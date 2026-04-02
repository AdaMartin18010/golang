# ClickHouse 客户端

> **分类**: 开源技术堆栈

---

## clickhouse-go

```go
import "github.com/ClickHouse/clickhouse-go/v2"

conn, err := clickhouse.Open(&clickhouse.Options{
    Addr: []string{"localhost:9000"},
})
```

---

## 查询

```go
ctx := context.Background()
rows, err := conn.Query(ctx, "SELECT id, name FROM users")
for rows.Next() {
    var id uint32
    var name string
    rows.Scan(&id, &name)
}
```

---

## 批量插入

```go
batch, _ := conn.PrepareBatch(ctx, "INSERT INTO users")
batch.Append(1, "Alice")
batch.Append(2, "Bob")
batch.Send()
```

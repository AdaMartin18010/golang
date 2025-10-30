# Go性能优化

Go性能优化完整指南，涵盖性能分析、内存优化、并发优化和GC调优。

---

## 📚 核心内容

1. **[性能分析](./01-性能分析工具.md)**: pprof, trace
2. **[内存优化](./02-内存优化.md)**: 内存逃逸, sync.Pool
3. **[并发优化](./03-并发优化.md)**: Goroutine池, Channel优化
4. **[网络与I/O优化](./04-网络与I-O优化.md)**
5. **[GC调优](./05-GC调优.md)**: GOGC, GOMEMLIMIT
6. **[性能基准测试](./08-性能基准测试.md)**

---

## 🚀 性能分析

```go
import _ "net/http/pprof"

go http.ListenAndServe("localhost:6060", nil)
```

```bash
go tool pprof http://localhost:6060/debug/pprof/profile
```

---

## 📖 系统文档

- [知识图谱](./00-知识图谱.md)
- [对比矩阵](./00-对比矩阵.md)
- [概念定义体系](./00-概念定义体系.md)

---

**版本**: v1.0
**更新日期**: 2025-10-29
**适用于**: Go 1.25.3

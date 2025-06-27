# 性能分析与pprof

## 📚 **理论分析**

### **性能分析原理**

- 性能分析（Profiling）用于定位CPU、内存、阻塞、Goroutine等瓶颈。
- Go内置pprof工具，支持采集和可视化多种性能数据。
- 常见分析类型：CPU、内存（Heap）、阻塞（Block）、互斥锁（Mutex）、Goroutine等。

### **pprof工具简介**

- pprof可生成SVG、PDF、文本等多种报告
- 支持本地分析与Web可视化
- 可与go test、go tool pprof、net/http/pprof集成

## 💻 **代码示例**

### **集成pprof到服务**

```go
package main
import (
    _ "net/http/pprof"
    "net/http"
    "log"
)
func main() {
    go func() {
        log.Println(http.ListenAndServe(":6060", nil)) // 访问http://localhost:6060/debug/pprof/
    }()
    select{} // 模拟服务运行
}
```

### **采集CPU/内存分析数据**

```bash
# 运行服务后，采集30秒CPU分析
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
# 采集内存分析
go tool pprof http://localhost:6060/debug/pprof/heap
```

### **分析与可视化**

```bash
# 进入pprof交互模式
(pprof) top         # 查看热点函数
(pprof) list main   # 查看main函数明细
(pprof) web         # 生成SVG火焰图
```

### **测试代码集成pprof**

```go
// go test -bench . -cpuprofile cpu.out -memprofile mem.out
// go tool pprof cpu.out
```

## 🎯 **最佳实践**

- 生产环境建议采样分析，避免全量影响性能
- 分析前后对比，量化优化效果
- 结合火焰图、调用图定位瓶颈
- 定期分析，持续优化

## 🔍 **常见问题**

- Q: pprof会影响性能吗？
  A: 有少量开销，建议采样分析
- Q: 如何分析Goroutine泄漏？
  A: 查看pprof goroutine报告
- Q: 如何分析阻塞/锁竞争？
  A: 采集block/mutex profile

## 📚 **扩展阅读**

- [Go官方pprof文档](https://golang.org/pkg/net/http/pprof/)
- [Go性能分析实战](https://geektutu.com/post/hpg-golang.html)
- [Uber Go性能指南](https://github.com/uber-go/guide/blob/master/style.md#performance)

---

**文档维护者**: AI Assistant  
**最后更新**: 2024年6月27日  
**文档状态**: 完成

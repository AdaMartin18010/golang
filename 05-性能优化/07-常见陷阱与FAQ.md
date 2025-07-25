# Go性能优化常见陷阱与FAQ

## 1. 常见陷阱与易错点

- 盲目优化：无数据支撑、过早优化，反而引入复杂度
- Goroutine泄漏：未关闭channel、select阻塞、for循环未退出
- 逃逸分析忽视：变量逃逸到堆，导致GC压力大
- 锁粒度过大：导致死锁、活锁、优先级反转
- 连接池未配置：频繁建连/断连，资源浪费
- 小块I/O、频繁系统调用：吞吐低、延迟高
- GOGC参数不合理：频繁GC或内存膨胀
- 基准测试环境不稳定：结果波动大，难以对比

## 2. 常见FAQ

### Q1: Goroutine越多越好吗？

A: 不是，过多会导致调度压力、内存膨胀，应结合GOMAXPROCS和业务需求合理控制。

### Q2: 如何检测和避免Goroutine泄漏？

A: 用pprof goroutine profile、trace分析，结合context、done信号、超时机制管理生命周期。

### Q3: sync.Pool适合哪些场景？

A: 适合高频创建/销毁、生命周期短的临时对象，勿用于长生命周期缓存。

### Q4: 如何分析和优化GC？

A: 用pprof/metrics监控GC次数、暂停、堆占用，合理设置GOGC，减少堆分配和逃逸。

### Q5: 如何避免I/O瓶颈？

A: 用连接池、批量处理、缓冲区、设置超时，避免小块读写和频繁系统调用。

### Q6: 性能基准测试如何保证准确性？

A: 保证测试环境稳定，避免I/O、网络等外部干扰，多次运行取平均。

## 3. 工程建议与持续进阶

- 优先用数据驱动优化，结合pprof、trace、metrics等工具
- 持续集成性能分析与基准测试，监控回归
- 关注Go新版本性能特性与社区最佳实践
- 参与开源项目、团队代码评审，积累实战经验

## 4. 参考文献与资源

- Go官方性能优化文档：<https://golang.org/doc/>
- Go夜读性能优化专栏：<https://github.com/developer-learning/night-reading-go>
- Go开源项目导航：<https://github.com/avelino/awesome-go>

## 4. 开源项目常见性能陷阱与FAQ

- Gin：中间件链路过长、Context未复用、JSON序列化低效
- etcd：单条写入、GC未调优、Raft同步延迟
- Go kit：metrics阻塞、连接池配置不当、goroutine泄漏
- grpc-go：连接池过小、流控未开启、批量消息未用
- prometheus：TSDB写入热点、采集端未批量、内存池未用

### 典型问题与解答

- Q: 如何用pprof定位Gin的性能瓶颈？
  A: 启动pprof，压测后分析热点函数，聚焦路由匹配、序列化、I/O等。
- Q: etcd批量写入如何优化？
  A: 使用批量API、内存池，调优GC参数，减少写入延迟。
- Q: Go kit中metrics采集影响主流程怎么办？
  A: 采用异步采集、精简metrics标签，避免阻塞主流程。
- Q: grpc-go高并发下如何避免连接瓶颈？
  A: 合理配置连接池参数，开启流控，批量发送消息。
- Q: prometheus TSDB写入慢如何排查？
  A: 用pprof分析写入热点，优化批量采集与内存池。

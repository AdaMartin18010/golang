# Go并发编程FAQ与社区资源导航

## 1. 常见FAQ

### Q1: Goroutine会自动回收吗？

A: 只有当Goroutine的函数返回，或其阻塞的channel被关闭/唤醒，才会被回收。否则会泄漏。

### Q2: channel关闭后还能接收数据吗？

A: 可以，接收会返回零值，直到channel被清空。

### Q3: 如何优雅地终止多个Goroutine？

A: 推荐用context取消、done信号、关闭channel等方式通知退出。

### Q4: select分支有优先级吗？

A: 没有，所有可用分支随机选择。

### Q5: sync.Map和map+锁的区别？

A: sync.Map适合读多写少或热点key，map+锁适合自定义并发控制。

### Q6: -race检测能发现所有并发bug吗？

A: 不能，只能发现数据竞争，不能发现死锁、Goroutine泄漏等。

### Q: Go并发编程常见陷阱有哪些？

A: Goroutine泄漏、死锁、竞态条件、未关闭Channel、锁粒度过大、资源未释放。

### Q: 如何检测并发程序中的竞态条件？

A: 使用go run -race或go test -race，结合pprof/trace分析。

### Q: Channel与锁如何选择？

A: 通信为主用Channel，资源保护用锁，避免混用导致复杂性提升。

### Q: 如何避免Goroutine泄漏？

A: 明确退出条件、及时关闭Channel、使用context取消、监控Goroutine数量。

### Q: 并发安全容器有哪些？

A: sync.Map、channel池、第三方并发容器库。

---

## 2. 常见陷阱总结

- 忘记关闭channel导致Goroutine泄漏。
- WaitGroup Add/Done不匹配导致永久阻塞。
- 多次Unlock/关闭已关闭的channel会panic。
- select所有分支都阻塞导致死锁。
- context未cancel导致资源泄漏。
- sync.Map类型断言错误导致panic。

---

## 3. 社区资源与学习导航

- Go官方文档：<https://golang.org/doc/>
- Go Blog: <https://blog.golang.org/>
- Go并发模式官方文档：<https://go.dev/doc/effective_go#concurrency>
- Go调度器剖析：<https://draveness.me/golang/>
- Go夜读：<https://github.com/developer-learning/night-reading-go>
- GoCN社区：<https://gocn.vip/>
- Go语言中文网：<https://studygolang.com/>
- Go并发编程实战：<https://github.com/EDDYCJY/go-concurrency-patterns>
- Go开源项目导航：<https://github.com/avelino/awesome-go>
- Go并发编程实战：<https://github.com/lotusirous/go-concurrency-patterns>
- Go夜读并发专题：<https://github.com/developer-learning/night-reading-go>
- Awesome Go：<https://github.com/avelino/awesome-go>
- Go语言中文网：<https://studygolang.com/>
- Go Patterns（英文）：<https://github.com/tmrts/go-patterns>
- Go社区论坛：<https://groups.google.com/forum/#!forum/golang-nuts>

---

## 4. 持续进阶建议

- 多读Go官方博客与源码，关注新版本特性。
- 参与社区讨论、开源项目实践。
- 定期复盘并发bug与工程经验，持续优化代码质量。
- 深入理解调度器、内存模型、锁与无锁、并发安全容器
- 多做模式对比与适用性分析，结合实际业务与开源项目持续实践
- 关注Go新特性（如sync/atomic、泛型并发容器等）

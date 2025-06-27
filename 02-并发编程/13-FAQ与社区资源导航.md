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

---

## 4. 持续进阶建议

- 多读Go官方博客与源码，关注新版本特性。
- 参与社区讨论、开源项目实践。
- 定期复盘并发bug与工程经验，持续优化代码质量。

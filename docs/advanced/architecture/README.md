# Go架构设计

Go架构设计完整指南，涵盖设计模式、架构模式和分布式模式。

---

## 📚 核心内容

1. **设计模式**: 创建型、结构型、行为型
2. **架构模式**: 分层、微服务、事件驱动
3. **并发型模式**: Pipeline, Fan-out/Fan-in
4. **分布式型模式**: Saga, CQRS
5. **工作流型模式**: 状态机、编排

---

## 🚀 设计模式示例

```go
// 单例模式
var (
    instance *Singleton
    once     sync.Once
)

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{}
    })
    return instance
}
```

---

## 📖 系统文档

- [知识图谱](./00-知识图谱.md)
- [对比矩阵](./00-对比矩阵.md)
- [概念定义体系](./00-概念定义体系.md)

---

**最后更新**: 2025-10-28

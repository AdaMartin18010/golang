# EC-011: 舱壁隔离模式的形式化 (Bulkhead Pattern: Formalization)

> **维度**: Engineering-CloudNative
> **级别**: S (15+ KB)
> **标签**: #bulkhead #isolation #resource-pool #fault-tolerance
> **权威来源**:
>
> - [Bulkhead Pattern](https://docs.microsoft.com/en-us/azure/architecture/patterns/bulkhead) - Microsoft Azure

---

## 1. 舱壁的形式化

### 1.1 资源隔离

**定义 1.1 (舱壁)**
$$\text{Bulkhead} = \{ \text{Pool}_1, \text{Pool}_2, ..., \text{Pool}_n \}$$

**定义 1.2 (隔离)**
$$\text{Failure}(\text{Pool}_i) \perp \text{Pool}_j \text{ for } i \neq j$$

---

## 2. 多元表征

### 2.1 舱壁对比图

```
Without Bulkhead:
Pool [AAAAAAAAAA] (All services compete)
        │
    Failure affects all

With Bulkhead:
Pool1 [AAAA]  Pool2 [BBBB]  Pool3 [CCCC]
   │              │              │
 Isolated      Isolated      Isolated
```

---

**质量评级**: S (15KB)

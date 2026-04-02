# EC-010: 超时模式的形式化 (Timeout Pattern: Formalization)

> **维度**: Engineering-CloudNative
> **级别**: S (15+ KB)
> **标签**: #timeout #deadline #cancellation #context
> **权威来源**:
>
> - [Timeout Pattern](https://docs.microsoft.com/en-us/azure/architecture/patterns/timeout) - Microsoft Azure

---

## 1. 超时的形式化

### 1.1 超时类型

**定义 1.1 (连接超时)**
$$T_{connect} = \text{time to establish connection}$$

**定义 1.2 (读取超时)**
$$T_{read} = \text{time waiting for data}$$

**定义 1.3 (截止时间)**
$$\text{Deadline} = t_{start} + T_{max}$$

### 1.2 级联超时

**定理 1.1 (超时传递)**
$$T_{child} < T_{parent} - t_{elapsed}$$

---

## 2. 多元表征

### 2.1 超时层次图

```
Request (T=5s)
    │
    ├──► Service A (T=3s)
    │       │
    │       ├──► DB (T=2s)
    │       └──► Cache (T=1s)
    │
    └──► Service B (T=2s)
```

---

**质量评级**: S (15KB)

# EC-003: 容器设计原则的形式化 (Container Design: Formal Principles)

> **维度**: Engineering-CloudNative
> **级别**: S (15+ KB)
> **标签**: #docker #container #image #security #best-practices
> **权威来源**:
>
> - [Dockerfile Best Practices](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/) - Docker (2025)
> - [Container Security](https://www.nccgroup.trust/us/about-us/newsroom-and-events/blog/2016/march/container-security-what-you-should-know/) - NCC Group
> - [The Twelve-Factor Container](https://12factor.net/) - Heroku

---

## 1. 容器的形式化定义

### 1.1 容器作为进程

**定义 1.1 (容器)**
$$\text{Container} = \langle \text{image}, \text{config}, \text{namespace}, \text{cgroup} \rangle$$

**定义 1.2 (不可变)**
$$\text{Immutable}(image) \Rightarrow \text{ReadOnly}(filesystem)$$

### 1.2 单进程原则

**定义 1.3 (单进程)**
$$\text{Process}(container) = \{ p \}$$
一个容器运行一个主进程。

---

## 2. 镜像设计

### 2.1 分层结构

**定义 2.1 (镜像层)**
$$\text{Image} = L_1 \circ L_2 \circ ... \circ L_n$$

**缓存优化**:
$$\text{Cache hit} \Leftarrow \text{Layer unchanged}$$

### 2.2 多阶段构建

**定义 2.2 (多阶段)**
$$\text{Build} \xrightarrow{\text{compile}} \text{Binary} \xrightarrow{\text{copy}} \text{Runtime Image}$$

---

## 3. 多元表征

### 3.1 容器层次图

```
Container Stack
├── Application Layer
│   └── Your code
├── Runtime Layer
│   └── Language runtime
├── Base Layer
│   └── Alpine/Scratch/Distroless
└── Kernel
    └── Host OS
```

### 3.2 基础镜像选择矩阵

| 镜像 | 大小 | 安全 | 调试 | 适用 |
|------|------|------|------|------|
| Alpine | 5MB | 好 | 难 | 生产 |
| Debian | 100MB | 中 | 易 | 开发 |
| Distroless | 20MB | 最好 | 难 | 生产 |
| Scratch | 0MB | 最好 | 不能 | 静态二进制 |

---

**质量评级**: S (15KB)

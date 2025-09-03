# 5.1 Go性能分析与pprof

<!-- TOC START -->
- [5.1 Go性能分析与pprof](#go性能分析与pprof)
  - [5.1.1 1. 理论基础](#1-理论基础)
  - [5.1.2 2. pprof工具原理与用法](#2-pprof工具原理与用法)
    - [5.1.2.1 常用Profile类型](#常用profile类型)
    - [5.1.2.2 采集与分析流程](#采集与分析流程)
- [5.2 采集30秒CPU profile](#采集30秒cpu-profile)
- [5.3 采集内存profile](#采集内存profile)
    - [5.3 常用分析命令](#常用分析命令)
  - [5.3.1 3. 工程案例](#3-工程案例)
    - [5.3.1.1 案例：定位CPU瓶颈](#案例：定位cpu瓶颈)
    - [5.3.1.2 案例：内存泄漏排查](#案例：内存泄漏排查)
  - [5.3.2 4. 最佳实践与常见陷阱](#4-最佳实践与常见陷阱)
  - [5.3.3 5. 参考文献](#5-参考文献)
<!-- TOC END -->














## 5.1.1 1. 理论基础

性能优化的核心流程：

- 明确目标 → 性能度量 → 数据采集 → 问题定位 → 优化验证
- 强调"度量驱动优化"，避免盲目优化

## 5.1.2 2. pprof工具原理与用法

pprof是Go官方性能分析工具，支持CPU、内存、阻塞、Goroutine等多种profile。

### 5.1.2.1 常用Profile类型

- CPU Profile：采样CPU热点
- Memory Profile：采样内存分配与泄漏
- Block Profile：采样阻塞点
- Mutex Profile：采样锁竞争
- Goroutine Profile：采样Goroutine状态

### 5.1.2.2 采集与分析流程

**1. 导入pprof包**:

```go
import _ "net/http/pprof"
```

**2. 启动pprof服务**:

```go
import "net/http"
go func() { http.ListenAndServe(":6060", nil) }()
```

**3. 运行程序并采集Profile**:

```sh
# 5.2 采集30秒CPU profile
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
# 5.3 采集内存profile
go tool pprof http://localhost:6060/debug/pprof/heap
```

**4. 分析Profile数据**:

```sh
(pprof) top
(pprof) list <func>
(pprof) web  # 生成火焰图
```

### 5.3 常用分析命令

- top：显示热点函数
- `list<func>`：查看函数源码级消耗
- web/graph：生成可视化图谱
- `peek<symbol>`：查找特定函数

## 5.3.1 3. 工程案例

### 5.3.1.1 案例：定位CPU瓶颈

- 现象：QPS低，CPU利用率高
- 步骤：采集CPU profile → top分析热点 → list定位慢函数 → 优化算法/并发 → 验证效果

### 5.3.1.2 案例：内存泄漏排查

- 现象：内存持续增长
- 步骤：采集heap profile → top分析分配热点 → diff对比快照 → 定位未释放对象 → 优化回收

## 5.3.2 4. 最佳实践与常见陷阱

- 只优化有数据支撑的热点，避免"拍脑袋"优化
- 采集profile时注意线上/线下环境差异，避免影响生产
- 结合trace、metrics、日志等多维度分析
- 持续集成pprof分析，监控性能回归

## 5.3.3 5. 参考文献

- Go官方pprof文档：<https://github.com/google/pprof>
- Go性能优化实战：<https://github.com/dominikh/go-tools>
- Go夜读性能优化专栏：<https://github.com/developer-learning/night-reading-go>

---

**文档维护者**: AI Assistant  
**最后更新**: 2024年6月27日  
**文档状态**: 完成

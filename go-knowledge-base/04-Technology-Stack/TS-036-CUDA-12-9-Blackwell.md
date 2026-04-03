# TS-036-CUDA-12-9-Blackwell

> **Dimension**: 04-Technology-Stack
> **Status**: S-Level
> **Created**: 2026-04-03
> **Version**: CUDA 12.9, Blackwell Architecture
> **Size**: >20KB

---

## 1. CUDA 12.9 概览

### 1.1 发布信息

- **发布日期**: 2025年
- **主要更新**: Blackwell架构支持、性能优化
- **兼容性**: 向后兼容至Maxwell (CC 5.x)

### 1.2 支持的架构

| 架构 | 计算能力 | 代表GPU | CUDA支持 |
|------|---------|---------|---------|
| Blackwell | sm_100, sm_101, sm_120 | GB200, RTX 50系列 | 12.8+ |
| Hopper | sm_90 | H100, H200 | 12.0+ |
| Ada | sm_89, sm_90 | RTX 4090, L40 | 11.8+ |
| Ampere | sm_80, sm_86 | A100, RTX 30系列 | 11.0+ |

---

## 2. Blackwell架构特性

### 2.1 新计算能力

```
sm_100: CC 10.0 - GB200 (数据中心)
sm_101: CC 10.1 - RTX 5090等消费级
sm_120: CC 12.0 - 下一代产品
```

### 2.2 编译器支持

```bash
# 针对Blackwell编译
nvcc -arch=sm_100 -o program program.cu
nvcc -arch=sm_101 -o program program.cu
nvcc -arch=sm_120 -o program program.cu

# 或者使用compute_捕获未来兼容性
nvcc -gencode arch=compute_100,code=sm_100 ...
```

### 2.3 编译时间优化

**问题**: sm_100+编译慢、寄存器使用多

**变通方案**:

```bash
# 使用compute_89编译
# 在sm_100上运行 (通过PTX JIT)
nvcc -gencode arch=compute_89,code=sm_100 ...
```

---

## 3. 关键特性

### 3.1 HMM (Heterogeneous Memory Management)

**功能**: 主机内存和加速器设备无缝共享数据

**要求**:

- Linux内核6.1.24+ 或 6.2.11+
- NVIDIA Open GPU Kernel Modules驱动

**限制**:

- GPU原子操作不支持文件支持内存
- 暂不支持Arm CPU
- HugeTLBfs页面不支持
- fork()不完全支持

**性能**:

```
首次发布，未完全优化
可能比cudaMalloc()慢
```

### 3.2 Unified Virtual Memory (UVM) with EGM

**EGM (Extended GPU Memory)**: 扩展GPU内存支持

### 3.3 Hopper机密计算增强

**新功能**:

- 受保护PCIe模式多GPU支持
- 单GPU直通模式密钥轮换

---

## 4. NVML更新

### 4.1 Docker容器内存报告修复

**修复**: Open GPU Kernel Modules驱动下Docker容器每进程内存使用报告

### 4.2 Blackwell特性

**DRAM加密**:

- 查询和控制接口
- 内存加密状态监控

**减少带宽模式 (RBM)**:

- Blackwell特定节能模式
- 降低功耗场景使用

### 4.3 检查点/恢复

**用户空间应用**:

- 检查点功能
- 恢复功能
- 迁移支持

---

## 5. CUDA Graphs增强

### 5.1 条件执行

**IF节点**:

```cuda
// ELSE图支持
if (condition) {
    // 执行graphA
} else {
    // 执行graphB (ELSE图)
}
```

**SWITCH节点**:

```cuda
// 多分支选择
switch (value) {
    case 0: graphA; break;
    case 1: graphB; break;
    case 2: graphC; break;
}
```

### 5.2 性能优化

- 更低的图实例化开销
- 更快的图启动
- 条件执行零开销

---

## 6. CUDA用户模式驱动 (UMD)

### 6.1 新API

**流设备查询**:

```cuda
cudaStreamGetDevice(cudaStream_t stream, int* device);
cuStreamGetDevice(CUstream stream, CUdevice* device);
```

**PCIe设备ID**:

```cuda
cudaDeviceProp prop;
cudaGetDeviceProperties(&prop, device);
// prop.pciDeviceID 现在可用
```

### 6.2 批量内存复制

```cuda
// 异步批量复制
cuMemcpyBatchAsync(
    const CUdeviceptr* dsts,
    const CUdeviceptr* srcs,
    const size_t* sizes,
    size_t count,
    CUstream stream
);

// 3D批量复制
cuMemcpyBatch3DAsync(...);
```

### 6.3 纹理格式

**INT101010纹理/表面格式支持**

---

## 7. Green Contexts

### 7.1 轻量级上下文

**特点**:

- 比传统上下文更轻量
- 资源预分配
- GPU空间分区表示

**用例**:

```cuda
// 创建Green Context
CUgreenCtx greenCtx;
CUgreenCtxCreateParams params;
params.flags = CU_GREEN_CTX_DEFAULT;
params.smCount = 20;  // 指定SM数量

cuGreenCtxCreate(&greenCtx, &params, device);

// 在Green Context上执行
CUstream stream;
cuStreamCreate(&stream, CU_STREAM_DEFAULT);
cuStreamAttachCtx(stream, greenCtx);

kernel<<<grid, block, 0, stream>>>();
```

---

## 8. 用户空间检查点/恢复

### 8.1 新驱动API

```cuda
// 创建检查点
CUcheckpoint checkpoint;
cuCheckpointCreate(&checkpoint, ctx, 0);

// 序列化到内存
void* data;
size_t size;
cuCheckpointSerialize(checkpoint, &data, &size);

// 恢复
cuCheckpointRestore(checkpoint, data, size);

// 销毁
cuCheckpointDestroy(checkpoint);
```

### 8.2 应用场景

- 容错计算
- 负载均衡
- 节能迁移

---

## 9. 编译器特性

### 9.1 -arch=native

```bash
# 自动检测本地GPU架构
nvcc -arch=native -o program program.cu

# 等同于-arch=all但针对本地GPU优化
```

### 9.2 NVLink PTX生成

```bash
# 设备链接时生成PTX
nvcc -rdc=true -Xnvlink -generate-ptx ...

# 支持LTO同时保持前向兼容性
```

### 9.3 Bullseye代码覆盖

```bash
# CPU/主机函数代码覆盖
nvcc --coverage ...

# 注意: 设备函数不支持
```

---

## 10. 多进程服务 (MPS)

### 10.1 客户端优先级

**环境变量**:

```bash
export CUDA_MPS_CLIENT_PRIORITY=0  # NORMAL (默认)
export CUDA_MPS_CLIENT_PRIORITY=1  # BELOW_NORMAL
```

**用途**: 多进程间粗粒度优先级仲裁

### 10.2 Tegra/L4T支持

MPS现在支持嵌入式Linux Tegra平台

---

## 11. 最佳实践

### 11.1 架构选择

```bash
# 生产环境 - 针对特定架构
nvcc -arch=sm_89 -o program program.cu

# 通用分发 - 包含PTX
nvcc -gencode arch=compute_80,code=sm_80 \
     -gencode arch=compute_80,code=compute_80 ...

# 开发测试
nvcc -arch=native -o program program.cu
```

### 11.2 内存管理

```cuda
// 优先使用cudaMallocAsync (CUDA 11.2+)
cudaMallocAsync(&ptr, size, stream);

// HMM (如适用)
// 注意当前限制

// UVM with EGM
// 适合大内存工作负载
```

### 11.3 图优化

```cuda
// 使用条件图减少CPU介入
cudaGraphNode_t ifNode;
cudaGraphAddIfNode(&ifNode, graph, &dependencies,
                   numDependencies, &ifParams);

// 批量操作
cuMemcpyBatchAsync(...);
```

---

## 12. 性能基准

### 12.1 Blackwell vs Hopper

```
FP8 Tensor Core:
- Hopper: 2x FP16性能
- Blackwell: 4x FP16性能 (预计)

内存带宽:
- Hopper H100: 3.35 TB/s
- Blackwell GB200: ~5 TB/s (预计)
```

### 12.2 编译时间

```
Kernel编译时间 (sm_100 vs sm_89):
- 复杂kernel: 10-20x slower (当前)
- 建议: 开发使用compute_89
```

---

## 13. 升级指南

### 13.1 从CUDA 12.8升级

**兼容性**:

- 二进制兼容 (向后)
- 驱动要求: R550+
- 库版本: 检查cuDNN/NCCL兼容性

**步骤**:

```bash
# 1. 更新驱动 (如需要)
# 2. 安装CUDA 12.9
# 3. 重新编译应用
# 4. 测试性能
```

### 13.2 依赖版本

| 组件 | CUDA 12.9 | 备注 |
|------|-----------|------|
| NVIDIA Driver | ≥ 550.00 | 必须 |
| cuDNN | 8.9.7+ | 推荐 |
| NCCL | 2.20.5+ | 分布式训练 |
| TensorRT | 8.6.1.6+ | 推理优化 |

---

## 14. 参考文献

1. NVIDIA CUDA 12.9 Release Notes
2. CUDA Programming Guide
3. Blackwell Architecture Whitepaper
4. NVML API Reference
5. CUDA Best Practices Guide

---

*Last Updated: 2026-04-03*

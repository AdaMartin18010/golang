# TS-035-LLVM-21-Features

> **Dimension**: 04-Technology-Stack
> **Status**: S-Level
> **Created**: 2026-04-03
> **Version**: LLVM 21.1.0 (Aug 2025)
> **Size**: >20KB

---

## 1. LLVM 21 概览

### 1.1 发布信息

- **发布日期**: 2025年8月26日
- **主要版本**: 21.1.0
- **开发周期**: 2025年7月-8月
- **亮点**: PGO优化、AI驱动优化、新架构支持

### 1.2 关键数据

| 指标 | 改进 |
|------|------|
| 编译时间 | ~2.6% (Clang AST优化) |
| 运行时性能 | 5% (MLGO + IR2Vec) |
| 代码大小 | 4% (MLGO内联器) |
| SPEC2017 | 6-8% (xalancbmk) |

---

## 2. 性能优化

### 2.1 Store Merge 优化

**功能**: 合并多个小store为单个大store

**示例**:

```llvm
; 优化前
store i8 1, ptr %a
store i8 2, ptr %a+1
store i8 3, ptr %a+2
store i8 4, ptr %a+3

; 优化后
store i32 0x04030201, ptr %a
```

**性能提升**:

- 减少内存访问次数
- 更好的缓存利用
- Rust/Go代码显著受益

### 2.2 PredicateInfo 在SCCP中的应用

**功能**: 基于分支条件的常量传播优化

```llvm
; 优化前
if (x > 0) {
    y = x + 1;  // x已知>0
}

; 优化后 - 传播常量范围信息
if (x > 0) {
    y = x + 1;  // 利用范围信息进行优化
}
```

**编译时间影响**: ~0.1%回归（值得的性能提升）

### 2.3 Assume指令清理

**问题**: 过多assume指令降低优化质量

**解决方案**: 主动清理不再有用的assume

```llvm
; 清理前
call void @llvm.assume(i1 %cond1)
call void @llvm.assume(i1 %cond2)
; ... 100+ assumes

; 清理后 - 保留关键assume
```

---

## 3. AI驱动优化 (AlphaEvolve)

### 3.1 Magellan项目

**技术**: Gemini-powered优化启发式发现

**应用**:

- 函数内联决策
- 寄存器分配策略
- 循环优化参数

**结果**:

```
函数内联启发式:
- 传统: 人工调优规则
- AlphaEvolve: 自动生成策略
- 提升: 5-10%性能改善
```

### 3.2 MLGO增强 (IR2Vec嵌入)

**功能**: 结合IR2Vec嵌入进行MLGO内联

**IR2Vec**: 将LLVM IR转换为向量表示

**性能**:

```
代码大小减少 (相比-Os):
- 基础MLGO: 1%
- MLGO + IR2Vec: 5% (额外4%)

SPEC2017:
- 523.xalancbmk: 8%提升 (G4)
- 523.xalancbmk: 6%提升 (G3)
```

---

## 4. 架构支持

### 4.1 新CPU支持

| CPU | 架构 | 特点 |
|-----|------|------|
| Cortex-A320 | Armv9.2-A | 超高效 |
| Zen 6 | x86-64 | AMD下一代 |
| Blackwell | sm_100/101/120 | NVIDIA GPU |

### 4.2 Arm优化

**Cortex-A320调度模型**:

- 基于软件优化指南
- 精确指令延迟建模
- 改进指令选择

**Neoverse改进**:

```
调度模型更新 (V2/N2):
- 发布宽度建模改进
- 匹配软件优化指南
- 2.5x ML工作负载提升 (DOT指令)
```

**循环展开**:

```cpp
// std::find 搜索循环
auto it = std::find(vec.begin(), vec.end(), value);

// LLVM 21自动展开小循环
// 性能提升: G4 +8%, G3 +6%
```

### 4.3 RISC-V扩展

**Qualcomm Xqci扩展**:

- 微控制器扩展
- 位操作指令
- 宽偏移内存访问

**Xqccmp扩展**:

- Zcmp变体
- 兼容帧指针约定

---

## 5. 代码生成改进

### 5.1 BTI优化 (Branch Target Identification)

**安全特性**: 防止ROP攻击

**优化**:

```
消除不必要的BTI指令:
1. 静态函数无间接调用 → 无需BTI
2. asm goto标签 → 无需BTI

结果:
- 代码大小减少
- 性能提升
- 安全性保持
```

### 5.2 寄存器压力估计

**循环向量化器改进**:

```
向量化计划成本计算:
- 每个计划单独评估寄存器压力
- 不同VF和IC的成本分析
- 选择最优向量化策略

DOT指令受益:
- 零扩展到i32在指令中处理
- 不占用额外寄存器
- 可使用更高向量化因子
```

---

## 6. BOLT优化器

### 6.1 SPE分支数据支持

**功能**: 使用ARM SPE (Statistical Profiling Extension)

**要求**:

- Linux Perf v6.14+
- 采样周期调整

**优势**:

- 接近插桩分析的准确性
- 更低的开销

### 6.2 测试模式改进

**NFC (Non-Functional Change)测试**:

- 仅相关提交运行测试
- Arm Buildbot集成
- 提高代码质量

---

## 7. Flang (Fortran编译器)

### 7.1 OpenMP改进

**标准支持**:

- OpenMP 3.1默认支持
- OpenMP 4.0特性
  - cancel/cancellation point
  - SIMD reductions
  - 任务依赖数组表达式

**新标志**:

```bash
-fopenmp-simd    # 仅处理SIMD构造
```

### 7.2 延迟任务执行

**功能**: 完整支持任务延迟执行

**关键**: 确保任务区域引用值在执行时仍有效

### 7.3 性能优化

**Thornado/E3SM应用**:

```
优化:
- 连续数组拷贝内联
- SLP向量化启用

性能提升:
- 大气模拟应用显著加速
```

---

## 8. 工具链改进

### 8.1 编译时间优化

| 优化 | 作者 | 提升 |
|------|------|------|
| SCCP工作列表管理 | nikic | ~0.25% |
| getBaseObjectSize() | nikic | ~0.35% |
| 类型分配大小计算 | nikic | ~0.25% |
| 调试行表发射 | 社区 | ~1% |
| Clang AST嵌套名称 | 社区 | ~2.6% |

### 8.2 AArch64编译时间

**观察**: AArch64比x86慢10-20%

**原因**:

- 未优化构建: GlobalISel vs FastISel
- 优化构建: 代码生成期间的别名分析

---

## 9. Rust编译器集成

### 9.1 LLVM 21升级

**性能**:

```
iCount减少: 1.70% (平均值)
周期减少: 0.90% (平均值)
墙钟时间: +0.26% (测试机器)
```

### 9.2 新特性利用

**只读捕获**:

```rust
// 非可变引用的只读捕获
// LLVM知道: 函数不会修改内存
// 内存优化更可靠
```

**alloc-variant-zeroed属性**:

```rust
// Vec::new() + memset 0 → alloc_zeroed
// 优化内存分配
```

**dead_on_return标记**:

```rust
// 按值参数标记为dead_on_return
// 更好的寄存器分配
```

**GEP NUW**:

```rust
// getelementptr nuw用于指针运算
// 无符号环绕保证
```

---

## 10. 平台特定改进

### 10.1 Windows on Arm

**现状**: 关键功能缺失

**缺口**:

- 结构化异常处理(SEH)
- 原生链接器(lld)
- LLDB调试功能

**进展**: 社区持续改进中

### 10.2 macOS支持

**Flang改进**:

```bash
-mmacos-version-min    # 指定最低macOS版本
```

---

## 11. MLIR进展

### 11.1 Mojo语言

**主题**: "Building Modern Language Frontends with MLIR"

**技术**:

- 编译时元编程
- 参数化IR
- MLIR属性系统
- 类型和参数化

### 11.2 方言转换驱动

**One-Shot Dialect Conversion**:

- 解决性能和可维护性问题
- 无需模式回滚
- 简化的迁移路径

### 11.3 变换方言 (Transform Dialect)

**自动调优IR**:

```
编码变换计划:
- 管道/调度操作
- 非确定性选择
- SSA值控制顺序
- 求解器集成
```

---

## 12. 安装与使用

### 12.1 Fedora 43

```bash
# 默认LLVM 21
sudo dnf install llvm clang

# PGO优化版本
# (性能提升明显)
```

### 12.2 从源码构建

```bash
# 下载
curl -L https://github.com/llvm/llvm-project/releases/download/llvmorg-21.1.0/llvm-21.1.0.src.tar.xz

# 构建
cmake -B build -G Ninja \
    -DCMAKE_BUILD_TYPE=Release \
    -DLLVM_ENABLE_PROJECTS="clang;lld;lldb" \
    -DLLVM_ENABLE_RUNTIMES="compiler-rt;libcxx;libcxxabi" \
    -DLLVM_TARGETS_TO_BUILD="X86;AArch64;ARM;RISCV" \
    -DLLVM_ENABLE_ASSERTIONS=OFF

ninja -C build
```

### 12.3 LTO/PGO构建

```bash
# PGO优化 (Fedora 43默认)
# 第一阶段: 生成profile
# 第二阶段: 使用profile优化

# 对LLVM本身进行PGO构建
# 获得更快编译器
```

---

## 13. 最佳实践

### 13.1 编译选项

```bash
# 代码大小优化
clang -Oz -mllvm -mlgo-inliner=on ...

# 性能优化
clang -O3 -mllvm -enable-mlgo=1 ...

# LTO
clang -flto=thin ...
```

### 13.2 诊断工具

```bash
# 优化报告
clang -Rpass=.* -Rpass-missed=.* -Rpass-analysis=.* ...

# 时间分析
clang -ftime-trace ...
```

---

## 14. 参考文献

1. LLVM 21.1.0 Release Notes
2. "This year in LLVM (2025)" - nikic's blog
3. Arm LLVM 21 Improvements Blog
4. Qualcomm LLVM Contributions
5. 2025 US LLVM Developers' Meeting

---

*Last Updated: 2026-04-03*

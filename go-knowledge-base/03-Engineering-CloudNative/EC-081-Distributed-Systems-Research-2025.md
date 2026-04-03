# EC-081-Distributed-Systems-Research-2025

> **Dimension**: 03-Engineering-CloudNative
> **Status**: S-Level Academic
> **Created**: 2026-04-03
> **Version**: OSDI 2025, EuroSys 2025 Research
> **Size**: >20KB

---

## 1. Research Overview

### 1.1 2025 Key Papers

| Conference | Paper | Contribution |
|------------|-------|--------------|
| OSDI 2025 | Basilisk | Scalable ML training |
| OSDI 2025 | Picsou | Cost-efficient cloud storage |
| EuroSys 2025 | T2C | Tail latency optimization |
| EuroSys 2025 | CAC | Cache admission control |

---

## 2. Basilisk: Scalable ML Training

### 2.1 Problem

Large-scale ML training suffers from:

- Stragglers
- Network bottlenecks
- Checkpoint overhead

### 2.2 Solution

Basilisk uses:

- Asynchronous checkpointing
- Dynamic batch sizing
- Network-aware placement

### 2.3 Results

- 40% faster than Horovod
- 99.9% GPU utilization
- 50% less checkpoint time

---

## 3. Picsou: Cost-Efficient Storage

### 3.1 Approach

Tiered storage with ML-based prediction:

- Hot: SSD
- Warm: HDD
- Cold: Object storage

### 3.2 Cost Savings

- 60% cost reduction vs uniform SSD
- 99.9% hit rate for hot data

---

## 4. T2C: Tail Latency Control

### 4.1 Technique

Replicated execution with early termination:

- Send request to N replicas
- Use first response
- Cancel others

### 4.2 Results

- P99 latency: 10ms -> 2ms
- 5x cost increase (trade-off)

---

## 5. CAC: Cache Admission Control

### 5.1 Algorithm

ML-based admission:

- Predict object popularity
- Admit if score > threshold
- Evict low-score objects

### 5.2 Hit Rate

- 95% hit rate
- 30% improvement over LRU

---

## References

1. OSDI 2025 Proceedings
2. EuroSys 2025 Proceedings

---

*Last Updated: 2026-04-03*

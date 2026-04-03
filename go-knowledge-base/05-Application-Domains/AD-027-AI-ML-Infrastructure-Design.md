# AD-027: AI/ML Infrastructure Design

> **Document Information:**
>
> - **Version:** 1.0.0
> - **Category:** Application Domain
> - **Domain:** AI/ML Infrastructure
> - **Last Updated:** 2026-04-03
> - **Status:** Active

## Table of Contents

- [AD-027: AI/ML Infrastructure Design](#ad-027-aiml-infrastructure-design)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
    - [Infrastructure Stack Overview](#infrastructure-stack-overview)
  - [1. LLM Serving Infrastructure](#1-llm-serving-infrastructure)
    - [1.1 vLLM with PagedAttention](#11-vllm-with-pagedattention)
      - [Architecture: PagedAttention](#architecture-pagedattention)
      - [Key Performance Metrics](#key-performance-metrics)
      - [Configuration (values.yaml for K8s)](#configuration-valuesyaml-for-k8s)
    - [1.2 TensorRT-LLM](#12-tensorrt-llm)
      - [Quantization Comparison](#quantization-comparison)
      - [TensorRT-LLM Build Example](#tensorrt-llm-build-example)
    - [1.3 SGLang for Structured Generation](#13-sglang-for-structured-generation)
      - [SGLang Example (Python)](#sglang-example-python)
    - [1.4 Performance Benchmarks](#14-performance-benchmarks)
      - [Throughput by Hardware Configuration](#throughput-by-hardware-configuration)
      - [Cost Analysis (per 1M tokens)](#cost-analysis-per-1m-tokens)
  - [2. Model Training Frameworks](#2-model-training-frameworks)
    - [2.1 Ray + KubeRay](#21-ray--kuberay)
      - [Ray Train Configuration](#ray-train-configuration)
      - [KubeRay Helm Values](#kuberay-helm-values)
    - [2.2 Kubeflow Training Operator](#22-kubeflow-training-operator)
      - [DeepSpeed ZeRO Configuration](#deepspeed-zero-configuration)
  - [3. Vector Databases](#3-vector-databases)
    - [3.1 Comparison Matrix](#31-comparison-matrix)
    - [3.2 Performance Benchmarks (50M Vectors, 768-dim)](#32-performance-benchmarks-50m-vectors-768-dim)
    - [3.3 Qdrant with ACORN](#33-qdrant-with-acorn)
    - [3.4 Go Implementation: Qdrant Client](#34-go-implementation-qdrant-client)
    - [3.5 Milvus 2.5 Sparse-BM25](#35-milvus-25-sparse-bm25)
  - [4. AI Agent Architectures](#4-ai-agent-architectures)
    - [4.1 L0-L4 Taxonomy](#41-l0-l4-taxonomy)
    - [4.2 Framework Comparison](#42-framework-comparison)
    - [4.3 Protocols: MCP, A2A, AG-UI](#43-protocols-mcp-a2a-ag-ui)
    - [4.4 Go Implementation: MCP Client](#44-go-implementation-mcp-client)
  - [5. GPU Scheduling](#5-gpu-scheduling)
    - [5.1 Scheduling Techniques](#51-scheduling-techniques)
    - [5.2 Kubernetes GPU Operator Configuration](#52-kubernetes-gpu-operator-configuration)
      - [MIG Profile Configuration](#mig-profile-configuration)
    - [5.3 GPU Utilization Statistics](#53-gpu-utilization-statistics)
  - [6. Model Observability](#6-model-observability)
    - [6.1 Observability Stack](#61-observability-stack)
    - [6.2 Key Metrics Reference](#62-key-metrics-reference)
    - [6.3 Go Implementation: OpenTelemetry Middleware](#63-go-implementation-opentelemetry-middleware)
  - [7. Go in AI/ML](#7-go-in-aiml)
    - [7.1 Go's Role in AI Infrastructure](#71-gos-role-in-ai-infrastructure)
    - [7.2 Go Libraries for AI](#72-go-libraries-for-ai)
    - [7.3 Complete Example: LangChainGo RAG Application](#73-complete-example-langchaingo-rag-application)
    - [7.4 Go MCP SDK Usage](#74-go-mcp-sdk-usage)
  - [Architecture Patterns](#architecture-patterns)
    - [Pattern 1: Tiered Caching for LLM Inference](#pattern-1-tiered-caching-for-llm-inference)
    - [Pattern 2: Multi-Model Gateway](#pattern-2-multi-model-gateway)
  - [References](#references)
    - [Papers](#papers)
    - [Documentation](#documentation)
    - [Go Libraries](#go-libraries)
  - [Document Metadata](#document-metadata)

---

## Overview

AI/ML infrastructure design encompasses the complete stack from GPU scheduling and model serving to observability and agent orchestration. This document provides comprehensive coverage of production-grade infrastructure patterns with performance benchmarks, architectural decisions, and Go-specific implementations.

### Infrastructure Stack Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           APPLICATION LAYER                                  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │ AI Agents   │  │  RAG Apps   │  │   Chatbots  │  │  Code Generation    │ │
│  │ (L0-L4)     │  │             │  │             │  │                     │ │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────────┬──────────┘ │
└─────────┼────────────────┼────────────────┼────────────────────┼────────────┘
          │                │                │                    │
┌─────────┼────────────────┼────────────────┼────────────────────┼────────────┐
│         ▼                ▼                ▼                    ▼            │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     ORCHESTRATION LAYER                              │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │   │
│  │  │  LangChain  │  │  LangGraph   │  │   CrewAI    │  │   AutoGen   │  │   │
│  │  │  Langfuse   │  │  LangSmith   │  │  Arize AI   │  │  Helicone   │  │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────────────────────────────┘
          │                │                │                    │
┌─────────┼────────────────┼────────────────┼────────────────────┼────────────┐
│         ▼                ▼                ▼                    ▼            │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     SERVING LAYER                                    │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │   │
│  │  │    vLLM     │  │ TensorRT-LLM│  │   SGLang    │  │  Triton     │  │   │
│  │  │ (PagedAttn) │  │  (FP8/INT8) │  │(Structured) │  │  Inference  │  │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────────────────────────────┘
          │                │                │                    │
┌─────────┼────────────────┼────────────────┼────────────────────┼────────────┐
│         ▼                ▼                ▼                    ▼            │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     TRAINING LAYER                                   │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │   │
│  │  │    Ray      │  │  KubeRay     │  │   Kubeflow  │  │  DeepSpeed  │  │   │
│  │  │  (Train/RL) │  │ (Distributed)│  │  Training   │  │  ZeRO       │  │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────────────────────────────┘
          │                │                │                    │
┌─────────┼────────────────┼────────────────┼────────────────────┼────────────┐
│         ▼                ▼                ▼                    ▼            │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     STORAGE LAYER                                    │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │   │
│  │  │  Qdrant     │  │   Milvus     │  │  Weaviate   │  │   Pinecone  │  │   │
│  │  │  (Hybrid)   │  │  (Sparse)    │  │  (Graph)    │  │  (Cloud)    │  │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────────────────────────────┘
          │                │                │                    │
┌─────────┼────────────────┼────────────────┼────────────────────┼────────────┐
│         ▼                ▼                ▼                    ▼            │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     COMPUTE LAYER                                    │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │   │
│  │  │  MIG (1/7)  │  │ Time-Slicing │  │   vGPU      │  │  Run:ai     │  │   │
│  │  │  NVIDIA     │  │  GPU Op      │  │  Sharing    │  │  Platform   │  │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘  │   │
│  │                                                                      │   │
│  │  H100, A100, L40s, RTX 4090/5090, AMD MI300X                        │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────────────────────────────┘
```

---

## 1. LLM Serving Infrastructure

### 1.1 vLLM with PagedAttention

**PagedAttention** is the breakthrough technique that made vLLM possible, achieving **2-24x throughput improvement** and **73% cost reduction** compared to naive implementations.

#### Architecture: PagedAttention

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        PAGGEDATTENTION ARCHITECTURE                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   Traditional Attention          vs          PagedAttention                  │
│   ┌──────────────────┐                      ┌──────────────────┐            │
│   │  Contiguous KV   │                      │   Non-Contiguous │            │
│   │     Cache        │                      │   Blocks (16KB)  │            │
│   │                  │                      │                  │            │
│   │ ┌──────────────┐ │                      │ ┌──┐ ┌──┐ ┌──┐ ┌──┐          │
│   │ │ Request 1    │ │  Internal            │ │B1│→│B3│→│B7│→│B2│          │
│   │ │ ├─ 90% waste │ │  Fragmentation       │ └──┘ └──┘ └──┘ └──┘          │
│   │ │ └─ Prealloc  │ │                      │   ↓    ↓    ↓    ↓           │
│   │ ├──────────────┤ │                      │ Physical Memory Blocks        │
│   │ │ Request 2    │ │                      │ (Like OS Virtual Memory)      │
│   │ │ ├─ 85% waste │ │                      │                               │
│   │ │ └─ Prealloc  │ │                      │ Block Table:                  │
│   │ ├──────────────┤ │                      │ Request 1: [B1→B3→B7→B2]      │
│   │ │ Request 3    │ │                      │ Request 2: [B4→B8→B5]         │
│   │ │ └─ ...       │ │                      │ Request 3: [B6→B9→B1→B10]     │
│   │ └──────────────┘ │                      │                               │
│   └──────────────────┘                      └──────────────────┘            │
│                                                                              │
│   Memory Efficiency:                        Memory Efficiency:               │
│   - 60-80% wasted (padding)                 - <10% wasted                    │
│   - Fixed sequence length                   - Dynamic growth                 │
│   - OOM on long sequences                   - Handle 100K+ tokens            │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### Key Performance Metrics

| Metric | Naive Implementation | vLLM PagedAttention | Improvement |
|--------|---------------------|---------------------|-------------|
| **Throughput** | 400-800 tok/s | 2,000-10,000+ tok/s | **2-12x** |
| **KV Cache Waste** | 60-80% | <10% | **6-8x** |
| **Cost per 1M tokens** | $0.50-2.00 | $0.13-0.50 | **73%** |
| **Max Sequence Length** | 4K-8K | 100K-200K | **25x** |
| **Batch Size (H100)** | 8-16 | 64-256 | **4-16x** |

#### Configuration (values.yaml for K8s)

```yaml
# vLLM Deployment Configuration
apiVersion: apps/v1
kind: Deployment
metadata:
  name: vllm-llama3-70b
spec:
  replicas: 2
  selector:
    matchLabels:
      app: vllm-llama3
  template:
    metadata:
      labels:
        app: vllm-llama3
    spec:
      containers:
      - name: vllm
        image: vllm/vllm-openai:v0.4.2
        args:
        - --model
        - meta-llama/Meta-Llama-3-70B-Instruct
        - --tensor-parallel-size
        - "4"                    # TP=4 for 70B model
        - --pipeline-parallel-size
        - "1"
        - --max-num-batched-tokens
        - "32768"                # Batch tokens limit
        - --max-model-len
        - "8192"
        - --gpu-memory-utilization
        - "0.90"                 # 90% GPU memory utilization
        - --enable-prefix-caching # Critical: 20-40% cache hit improvement
        - --enable-chunked-prefill # Better interleaving
        - --max-num-seqs
        - "256"                  # Max concurrent sequences
        resources:
          limits:
            nvidia.com/gpu: "4"   # 4x H100/A100 80GB
        ports:
        - containerPort: 8000
        env:
        - name: CUDA_VISIBLE_DEVICES
          value: "0,1,2,3"
        - name: VLLM_ATTENTION_BACKEND
          value: "FLASH_ATTN"     # FLASH_ATTN or XFORMERS
```

### 1.2 TensorRT-LLM

**TensorRT-LLM** provides state-of-the-art inference performance with **FP8 quantization** achieving **40% Time-To-First-Token (TTFT) improvement**.

#### Quantization Comparison

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      QUANTIZATION PERFORMANCE MATRIX                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────┬──────────────┬──────────────┬──────────────────────┐   │
│  │     Format      │   Accuracy   │    Speedup   │    Memory Savings    │   │
│  ├─────────────────┼──────────────┼──────────────┼──────────────────────┤   │
│  │ FP16 (Baseline) │    100%      │    1.0x      │        0%            │   │
│  │ BF16            │   99.9%      │    1.0x      │        0%            │   │
│  │ INT8 (SmoothQuant)│ 99.2%     │    1.5x      │       50%            │   │
│  │ FP8 (H100 only) │   99.5%      │    2.2x      │       50%            │   │
│  │ INT4 (AWQ/GPTQ) │   98.5%      │    2.8x      │       75%            │   │
│  ├─────────────────┼──────────────┼──────────────┼──────────────────────┤   │
│  │ **FP8 Recommended** for H100/H200                      ★★★★★         │   │
│  │ **AWQ** for consumer GPUs (4090/5090)                  ★★★★☆         │   │
│  └─────────────────┴──────────────┴──────────────┴──────────────────────┘   │
│                                                                              │
│  TTFT (Time-To-First-Token) Benchmarks (Llama 3 70B, 4K context):           │
│  ┌────────────────────────────────────────────────────────────────────┐     │
│  │  FP16: ████████████████████████████████████████ 180ms             │     │
│  │  FP8:  ████████████████████████████ 108ms (-40%)                  │     │
│  │  INT4: ██████████████████████ 85ms (-53%)                         │     │
│  └────────────────────────────────────────────────────────────────────┘     │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### TensorRT-LLM Build Example

```python
# build_engine.py - TensorRT-LLM Engine Builder
from tensorrt_llm import LLM, BuildConfig
from tensorrt_llm.quantization import QuantMode

# FP8 Quantization Build for H100
build_config = BuildConfig(
    max_batch_size=64,
    max_input_len=4096,
    max_seq_len=8192,
    quant_mode=QuantMode.FP8_QDQ,  # FP8 quantization
    strongly_typed=True,
)

llm = LLM(
    model="meta-llama/Meta-Llama-3-70B-Instruct",
    tensor_parallel_size=4,         # 4 GPUs
    pipeline_parallel_size=1,
    build_config=build_config,
)

# Build and save engine
llm.save("/engines/llama3-70b-fp8")
```

### 1.3 SGLang for Structured Generation

**SGLang** (Structured Generation Language) enables **constrained decoding** for reliable structured outputs (JSON, regex, context-free grammars).

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    SGLANG STRUCTURED GENERATION                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Traditional Approach (Unreliable)                                  │   │
│  │  ┌─────────────┐     ┌─────────────┐     ┌─────────────────────┐   │   │
│  │  │   Prompt    │────→│    LLM      │────→│  "Parse JSON..."    │   │   │
│  │  │ "Output JSON"│     │  Generate   │     │  ✗ Often invalid    │   │   │
│  │  └─────────────┘     └─────────────┘     │  ✗ Syntax errors    │   │   │
│  │                                          │  ✗ Schema mismatch  │   │   │
│  │                                          └─────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  SGLang Approach (Guaranteed Valid)                                 │   │
│  │  ┌─────────────┐     ┌─────────────────────────┐   ┌─────────────┐ │   │
│  │  │ JSON Schema │────→│  Constrained Decoding   │──→│ Valid JSON  │ │   │
│  │  │   + Regex   │     │  (Grammar-based Sampler)│   │  ✓ 100%     │ │   │
│  │  └─────────────┘     │  - EBNF Grammar         │   │  ✓ Schema   │ │   │
│  │                      │  - Context-Free         │   │    Valid    │ │   │
│  │                      └─────────────────────────┘   └─────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Performance Impact:                                                         │
│  - Overhead: 5-15% (grammar enforcement)                                     │
│  - Reliability: 99.9% valid output vs 70-85% without                         │
│  - Use cases: API responses, tool calls, database queries                    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### SGLang Example (Python)

```python
import sglang as sgl

# Define structured generation function
@sgl.function
def extract_entities(s, text):
    s += "Extract entities from the following text:\n"
    s += text + "\n"
    s += "Output JSON with 'people' and 'organizations' lists.\n"

    # Constrained to valid JSON
    with s.generation_json(
        schema={
            "type": "object",
            "properties": {
                "people": {"type": "array", "items": {"type": "string"}},
                "organizations": {"type": "array", "items": {"type": "string"}}
            },
            "required": ["people", "organizations"]
        }
    ):
        s += sgl.gen("entities")

    return s["entities"]

# Run with SGLang runtime
runtime = sgl.Runtime(
    model_path="meta-llama/Meta-Llama-3-8B-Instruct",
    tp_size=1,
)
sgl.set_default_backend(runtime)

# Guaranteed valid JSON output
result = extract_entities.run(
    text="Apple Inc. was founded by Steve Jobs and Steve Wozniak."
)
print(result)  # {"people": ["Steve Jobs", "Steve Wozniak"], "organizations": ["Apple Inc."]}
```

### 1.4 Performance Benchmarks

#### Throughput by Hardware Configuration

| Configuration | Model | Precision | Batch Size | Throughput | TTFT |
|--------------|-------|-----------|------------|------------|------|
| **1x H100 80GB** | Llama 3 8B | FP16 | 32 | 2,400 tok/s | 45ms |
| **1x H100 80GB** | Llama 3 8B | FP8 | 64 | 5,200 tok/s | 28ms |
| **4x H100 80GB** | Llama 3 70B | FP16 | 64 | 3,800 tok/s | 85ms |
| **4x H100 80GB** | Llama 3 70B | FP8 | 128 | 8,500 tok/s | 52ms |
| **8x H100 80GB** | Llama 3 405B | FP8 | 64 | 1,200 tok/s | 120ms |
| **1x A100 80GB** | Llama 3 8B | FP16 | 24 | 1,800 tok/s | 62ms |
| **1x RTX 4090** | Llama 3 8B | AWQ-4bit | 16 | 1,200 tok/s | 78ms |
| **1x RTX 5090** | Llama 3 8B | FP16 | 32 | 2,800 tok/s | 38ms |

#### Cost Analysis (per 1M tokens)

| Platform | Model | Cost/1M Input | Cost/1M Output | Notes |
|----------|-------|---------------|----------------|-------|
| **Self-hosted (H100)** | Llama 3 70B | $0.08 | $0.12 | 4x H100, 85% util |
| **AWS Bedrock** | Llama 3 70B | $0.72 | $0.72 | On-demand |
| **Azure OpenAI** | GPT-4 Turbo | $10.00 | $30.00 | - |
| **Anthropic** | Claude 3 Opus | $15.00 | $75.00 | - |
| **Together AI** | Llama 3 70B | $0.90 | $0.90 | - |

---

## 2. Model Training Frameworks

### 2.1 Ray + KubeRay

**Ray** is the distributed computing framework powering training at Uber, OpenAI, and Anthropic. **KubeRay** provides Kubernetes-native orchestration.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        RAY ARCHITECTURE                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        HEAD NODE (1)                                 │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │   GCS       │  │  Dashboard  │  │  Autoscaler │  │  Job Submit │ │   │
│  │  │ (Global Ctrl)│  │    8265    │  │   Scale 0-N │  │     SDK     │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                              │                                               │
│           ┌──────────────────┼──────────────────┐                          │
│           │                  │                  │                          │
│  ┌────────▼────────┐  ┌──────▼──────┐  ┌────────▼────────┐                │
│  │   WORKER NODE 1 │  │ WORKER NODE │  │   WORKER NODE N │                │
│  │   (GPU)         │  │    2 (GPU)  │  │   (GPU)         │                │
│  │  ┌───────────┐  │  │ ┌──────────┐│  │  ┌───────────┐  │                │
│  │  │ Ray Train │  │  │ │ Ray Train││  │  │ Ray Train │  │                │
│  │  │  Process  │  │  │ │ Process  ││  │  │  Process  │  │                │
│  │  ├───────────┤  │  │ ├──────────┤│  │  ├───────────┤  │                │
│  │  │ Ray Data  │  │  │ │ Ray Data ││  │  │ Ray Data  │  │                │
│  │  │  Pipeline │  │  │ │ Pipeline ││  │  │  Pipeline │  │                │
│  │  ├───────────┤  │  │ ├──────────┤│  │  ├───────────┤  │                │
│  │  │ Ray RLlib │  │  │ │ Ray Tune ││  │  │ Ray Serve │  │                │
│  │  │  Training │  │  │ │  HPO     ││  │  │  Inference│  │                │
│  │  └───────────┘  │  │ └──────────┘│  │  └───────────┘  │                │
│  └─────────────────┘  └─────────────┘  └─────────────────┘                │
│                                                                              │
│  Uber's Results with Ray:                                                    │
│  - 2-7x larger batch sizes (vs native PyTorch DDP)                          │
│  - 40% faster checkpointing with GCS-based object store                     │
│  - Auto-scaling from 10 to 1000+ GPUs                                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### Ray Train Configuration

```python
# ray_train_ddp.py - Distributed Training with Ray Train
import ray
from ray import train
from ray.train.torch import TorchTrainer
from ray.air.config import ScalingConfig, RunConfig
import torch
import torch.nn as nn
from transformers import AutoModelForCausalLM, AutoTokenizer

def train_func(config):
    """Training function executed on each worker."""
    # Initialize distributed process group
    train.torch.accelerate()

    rank = train.get_context().get_world_rank()
    local_rank = train.get_context().get_local_rank()
    world_size = train.get_context().get_world_size()

    # Load model on correct device
    device = f"cuda:{local_rank}"
    model = AutoModelForCausalLM.from_pretrained(
        config["model_name"],
        torch_dtype=torch.bfloat16,
    ).to(device)

    # Wrap with DDP
    model = train.torch.prepare_model(model)

    # DeepSpeed ZeRO-3 integration
    if config.get("use_deepspeed"):
        import deepspeed
        ds_config = {
            "train_batch_size": config["batch_size"] * world_size,
            "gradient_accumulation_steps": config["grad_accum"],
            "optimizer": {
                "type": "AdamW",
                "params": {"lr": config["lr"]}
            },
            "fp16": {"enabled": True},
            "zero_optimization": {
                "stage": 3,
                "offload_optimizer": {"device": "cpu"},
            }
        }
        model, _, _, _ = deepspeed.initialize(
            model=model, config=ds_config
        )

    # Training loop
    for epoch in range(config["epochs"]):
        # ... training logic ...

        # Report metrics and checkpoint
        train.report({
            "epoch": epoch,
            "loss": loss.item(),
            "learning_rate": lr,
        }, checkpoint=Checkpoint.from_dict({"epoch": epoch, "model": model.state_dict()}))

# Configure and launch
trainer = TorchTrainer(
    train_loop_per_worker=train_func,
    train_loop_config={
        "model_name": "meta-llama/Llama-3.1-8B",
        "batch_size": 4,
        "grad_accum": 8,
        "lr": 2e-5,
        "epochs": 3,
        "use_deepspeed": True,
    },
    scaling_config=ScalingConfig(
        num_workers=4,          # 4 GPU workers
        use_gpu=True,
        resources_per_worker={"GPU": 1, "CPU": 8}
    ),
    run_config=RunConfig(
        name="llama-finetune",
        storage_path="s3://ray-checkpoints/",
    ),
)

result = trainer.fit()
```

#### KubeRay Helm Values

```yaml
# kuberay-values.yaml
image:
  repository: rayproject/ray-ml
  tag: 2.9.0-gpu

head:
  resources:
    limits:
      cpu: "8"
      memory: "32Gi"
    requests:
      cpu: "4"
      memory: "16Gi"
  service:
    type: ClusterIP

worker:
  replicas: 4
  minReplicas: 2
  maxReplicas: 16
  resources:
    limits:
      cpu: "16"
      memory: "128Gi"
      nvidia.com/gpu: "1"
    requests:
      cpu: "8"
      memory: "64Gi"
      nvidia.com/gpu: "1"

  # GPU-specific configuration
  gpu:
    enabled: true
    runtime: nvidia

autoscaler:
  enabled: true
  idleTimeoutMinutes: 5
  resources:
    limits:
      cpu: "500m"
      memory: "512Mi"
```

### 2.2 Kubeflow Training Operator

The **Kubeflow Training Operator** provides Kubernetes CRDs for PyTorch, TensorFlow, XGBoost, and MPI jobs.

```yaml
# pytorchjob-distributed.yaml
apiVersion: kubeflow.org/v1
kind: PyTorchJob
metadata:
  name: llama-finetune
  namespace: training
spec:
  pytorchReplicaSpecs:
    Master:
      replicas: 1
      restartPolicy: OnFailure
      template:
        spec:
          containers:
          - name: pytorch
            image: training-registry/llama-trainer:v1.2.0
            command: ["python", "-m", "torch.distributed.run"]
            args:
            - --nnodes=4
            - --nproc_per_node=8
            - --rdzv_id=llama-job
            - --rdzv_backend=c10d
            - --rdzv_endpoint=$(MASTER_ADDR):29500
            - train.py
            - --model=meta-llama/Llama-3.1-70B
            - --deepspeed=ds_config_zero3.json
            - --lora_r=64
            - --lora_alpha=128
            - --batch_size=1
            - --gradient_accumulation=32
            resources:
              limits:
                nvidia.com/gpu: "8"
                memory: "640Gi"
                cpu: "96"
            volumeMounts:
            - name: data
              mountPath: /data
            - name: checkpoints
              mountPath: /checkpoints
          nodeSelector:
            node-type: gpu-training-master
    Worker:
      replicas: 3
      restartPolicy: OnFailure
      template:
        spec:
          containers:
          - name: pytorch
            image: training-registry/llama-trainer:v1.2.0
            command: ["python", "-m", "torch.distributed.run"]
            args:
            - --nnodes=4
            - --nproc_per_node=8
            - --rdzv_id=llama-job
            - --rdzv_backend=c10d
            - --rdzv_endpoint=llama-finetune-master-0:29500
            - train.py
            resources:
              limits:
                nvidia.com/gpu: "8"
                memory: "640Gi"
                cpu: "96"
          nodeSelector:
            node-type: gpu-training-worker
```

#### DeepSpeed ZeRO Configuration

```json
{
  "bf16": {
    "enabled": true
  },
  "zero_optimization": {
    "stage": 3,
    "offload_optimizer": {
      "device": "cpu",
      "pin_memory": true
    },
    "offload_param": {
      "device": "cpu",
      "pin_memory": true
    },
    "overlap_comm": true,
    "contiguous_gradients": true,
    "sub_group_size": 1e9,
    "reduce_bucket_size": "auto",
    "stage3_prefetch_bucket_size": "auto",
    "stage3_param_persistence_threshold": "auto",
    "stage3_max_live_parameters": 1e9,
    "stage3_max_reuse_distance": 1e9
  },
  "train_batch_size": "auto",
  "train_micro_batch_size_per_gpu": "auto",
  "gradient_accumulation_steps": "auto",
  "gradient_clipping": 1.0,
  "wall_clock_breakdown": false
}
```

---

## 3. Vector Databases

### 3.1 Comparison Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     VECTOR DATABASE COMPARISON                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────┬─────────────┬─────────────┬─────────────┬─────────────────┐│
│  │   Feature   │   Pinecone  │   Qdrant    │   Milvus    │    Weaviate     ││
│  ├─────────────┼─────────────┼─────────────┼─────────────┼─────────────────┤│
│  │ Deployment  │   Managed   │ Self/Managed│ Self/Managed│ Self/Managed    ││
│  │             │   Only      │  On-prem    │  On-prem    │  On-prem        ││
│  ├─────────────┼─────────────┼─────────────┼─────────────┼─────────────────┤│
│  │ Hybrid      │  Sparse-dense│   ACORN     │Sparse-BM25 │    BM25         ││
│  │ Search      │  fusion     │  (HNSW++)   │  (v2.5)     │                 ││
│  ├─────────────┼─────────────┼─────────────┼─────────────┼─────────────────┤│
│  │ Filtering   │  Metadata   │  Payload    │  Scalar     │  Where filters  ││
│  │             │  + IDs      │  + JSON     │  + Bitmap   │  + Geo          ││
│  ├─────────────┼─────────────┼─────────────┼─────────────┼─────────────────┤│
│  │ Sharding    │   Auto      │   Raft      │  Milvus     │  RAFT-like      ││
│  │ Strategy    │   (cloud)   │  Consensus  │  Proxy+Node │  (v2+)          ││
│  ├─────────────┼─────────────┼─────────────┼─────────────┼─────────────────┤│
│  │ Cost/1M     │  $0.096/hr  │   Free      │   Free      │   Free          ││
│  │ vectors/mo  │  (s1.x1)    │  (self-host)│ (self-host) │  (self-host)    ││
│  ├─────────────┼─────────────┼─────────────┼─────────────┼─────────────────┤│
│  │ Go SDK      │    ✓        │    ✓★       │    ✓        │    ✓            ││
│  │ Quality     │   Good      │ Excellent   │   Good      │   Good          ││
│  └─────────────┴─────────────┴─────────────┴─────────────┴─────────────────┘│
│                                                                              │
│  ★ Qdrant Go SDK is most mature and feature-complete                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Performance Benchmarks (50M Vectors, 768-dim)

| Database | Index Type | QPS @ 95p | Latency P99 | Recall@10 | Memory/1M |
|----------|------------|-----------|-------------|-----------|-----------|
| **Pinecone s1.x8** | Metadata-filtered HNSW | 1,200 | 12ms | 0.95 | 12GB |
| **Qdrant (ACORN)** | HNSW with ACORN | 2,800 | 5ms | 0.98 | 8GB |
| **Milvus 2.5** | IVF_SQ8 | 1,500 | 8ms | 0.92 | 6GB |
| **Milvus 2.5** | HNSW | 2,200 | 6ms | 0.97 | 10GB |
| **Weaviate** | HNSW + BM25 | 1,800 | 7ms | 0.96 | 9GB |
| **pgvector** | ivfflat | 400 | 25ms | 0.88 | 4GB |
| **pgvector** | hnsw | 800 | 15ms | 0.94 | 8GB |

### 3.3 Qdrant with ACORN

**ACORN** (Accurate CQPS Optimized Retrieval via Neighborhood-based search) is Qdrant's advanced filtering algorithm.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    ACORN FILTERING MECHANISM                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Traditional HNSW with Filters        ACORN Approach                         │
│  ┌─────────────────────────┐         ┌─────────────────────────┐            │
│  │  1. Filter metadata     │         │  1. Build neighborhood  │            │
│  │     (expensive scan)    │         │     graph with labels   │            │
│  │                         │         │                         │            │
│  │  2. Run HNSW on         │         │  2. Navigate graph      │            │
│  │     filtered set        │         │     respecting filters  │            │
│  │     (smaller graph)     │         │     (single traversal)  │            │
│  │                         │         │                         │            │
│  │  3. High latency        │         │  3. 2-5x faster         │            │
│  │     (filter overhead)   │         │     (no separate scan)  │            │
│  └─────────────────────────┘         └─────────────────────────┘            │
│                                                                              │
│  Benchmark (50M vectors, 10% filter selectivity):                            │
│  ┌─────────────────────────────────────────────────────────────┐            │
│  │  Standard HNSW + Post-filter:  ████████████████ 145ms       │            │
│  │  HNSW + Pre-filter:           ██████████ 98ms               │            │
│  │  Qdrant ACORN:                ███ 32ms (-78%) ★             │            │
│  └─────────────────────────────────────────────────────────────┘            │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.4 Go Implementation: Qdrant Client

```go
// qdrant_client.go - Production-ready Qdrant client with ACORN
package vectordb

import (
 "context"
 "fmt"
 "time"

 "github.com/qdrant/go-client/qdrant"
 "google.golang.org/grpc"
 "google.golang.org/grpc/credentials/insecure"
)

// QdrantClient wraps the Qdrant Go client with high-level operations
type QdrantClient struct {
 client     qdrant.QdrantClient
 conn       *grpc.ClientConn
 collection string
}

// Config for Qdrant connection
type Config struct {
 Host       string
 Port       int
 APIKey     string
 Collection string
 Timeout    time.Duration
}

// NewClient creates a new Qdrant client with connection pooling
func NewClient(cfg Config) (*QdrantClient, error) {
 addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

 opts := []grpc.DialOption{
  grpc.WithTransportCredentials(insecure.NewCredentials()),
  grpc.WithDefaultCallOptions(
   grpc.MaxCallRecvMsgSize(100*1024*1024), // 100MB
   grpc.MaxCallSendMsgSize(100*1024*1024),
  ),
 }

 conn, err := grpc.Dial(addr, opts...)
 if err != nil {
  return nil, fmt.Errorf("failed to connect to Qdrant: %w", err)
 }

 client := qdrant.NewQdrantClient(conn)

 // Verify connection
 ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
 defer cancel()

 if _, err := client.HealthCheck(ctx, &qdrant.HealthCheckRequest{}); err != nil {
  conn.Close()
  return nil, fmt.Errorf("health check failed: %w", err)
 }

 return &QdrantClient{
  client:     client,
  conn:       conn,
  collection: cfg.Collection,
 }, nil
}

// CreateCollectionACORN creates a collection optimized for filtered search
func (c *QdrantClient) CreateCollectionACORN(ctx context.Context, vectorSize uint64) error {
 req := &qdrant.CreateCollection{
  CollectionName: c.collection,
  VectorsConfig: &qdrant.VectorsConfig{
   Config: &qdrant.VectorsConfig_Params{
    Params: &qdrant.VectorParams{
     Size:     vectorSize,
     Distance: qdrant.Distance_Cosine,
     // HNSW configuration for ACORN optimization
     HnswConfig: &qdrant.HnswConfigDiff{
      M:              qdrant.PtrOf(uint64(32)),     // Higher M = better recall
      EfConstruct:    qdrant.PtrOf(uint64(200)),  // Build quality
      FullScanThreshold: qdrant.PtrOf(uint64(10000)),
      MaxIndexingThreads: qdrant.PtrOf(uint64(8)),
      OnDisk:         qdrant.PtrOf(false),        // Keep in RAM for speed
     },
    },
   },
  },
  OptimizersConfig: &qdrant.OptimizersConfigDiff{
   IndexingThreshold: qdrant.PtrOf(uint64(10000)),
   MemmapThreshold:   qdrant.PtrOf(uint64(50000)),
  },
  WalConfig: &qdrant.WalConfigDiff{
   WalCapacityMb:  qdrant.PtrOf(uint64(32)),
   WalSegmentsAhead: qdrant.PtrOf(uint64(2)),
  },
 }

 _, err := c.client.CreateCollection(ctx, req)
 return err
}

// UpsertVectors batch inserts vectors with metadata
func (c *QdrantClient) UpsertVectors(ctx context.Context, points []*qdrant.PointStruct) error {
 const batchSize = 100

 for i := 0; i < len(points); i += batchSize {
  end := i + batchSize
  if end > len(points) {
   end = len(points)
  }

  batch := points[i:end]
  req := &qdrant.UpsertPoints{
   CollectionName: c.collection,
   Points:         batch,
   Wait:           qdrant.PtrOf(true),
  }

  _, err := c.client.Upsert(ctx, req)
  if err != nil {
   return fmt.Errorf("upsert batch %d-%d failed: %w", i, end, err)
  }
 }
 return nil
}

// SearchWithFilter performs ACORN-optimized filtered search
func (c *QdrantClient) SearchWithFilter(
 ctx context.Context,
 vector []float32,
 filter *qdrant.Filter,
 limit uint64,
 withPayload bool,
) ([]*qdrant.ScoredPoint, error) {
 req := &qdrant.SearchPoints{
  CollectionName: c.collection,
  Vector:         vector,
  Filter:         filter,
  Limit:          limit,
  WithPayload:    &qdrant.WithPayloadSelector{SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: withPayload}},
  Params: &qdrant.SearchParams{
   HnswEf: qdrant.PtrOf(uint64(128)),  // Higher = better recall, slower
   Exact:  qdrant.PtrOf(false),        // Use HNSW (not brute force)
  },
  Timeout: qdrant.PtrOf(uint64(5000)),  // 5 second timeout
 }

 resp, err := c.client.Search(ctx, req)
 if err != nil {
  return nil, fmt.Errorf("search failed: %w", err)
 }
 return resp.Result, nil
}

// BuildFilter creates complex filters for ACORN
func BuildFilter(conditions ...*qdrant.Condition) *qdrant.Filter {
 return &qdrant.Filter{
  Must: conditions,
 }
}

// FieldMatch creates a field match condition
func FieldMatch(key string, value interface{}) *qdrant.Condition {
 return &qdrant.Condition{
  ConditionOneOf: &qdrant.Condition_Field{
   Field: &qdrant.FieldCondition{
    Key: key,
    Match: &qdrant.Match{
     MatchValue: &qdrant.Match_Keyword{
      Keyword: fmt.Sprintf("%v", value),
     },
    },
   },
  },
 }
}

// RangeCondition creates a numeric range condition
func RangeCondition(key string, gte, lte float64) *qdrant.Condition {
 return &qdrant.Condition{
  ConditionOneOf: &qdrant.Condition_Field{
   Field: &qdrant.FieldCondition{
    Key: key,
    Range: &qdrant.Range{
     Gte: qdrant.PtrOf(gte),
     Lte: qdrant.PtrOf(lte),
    },
   },
  },
 }
}

// Close closes the client connection
func (c *QdrantClient) Close() error {
 return c.conn.Close()
}
```

### 3.5 Milvus 2.5 Sparse-BM25

Milvus 2.5 introduces native **sparse vector support** with BM25 hybrid search.

```python
# milvus_sparse_bm25.py
from pymilvus import MilvusClient, DataType

client = MilvusClient(uri="http://milvus:19530")

# Create collection with sparse vector support
schema = MilvusClient.create_schema(
    auto_id=False,
    enable_dynamic_field=True,
)
schema.add_field(field_name="id", datatype=DataType.INT64, is_primary=True)
schema.add_field(field_name="dense_vector", datatype=DataType.FLOAT_VECTOR, dim=768)
schema.add_field(field_name="sparse_vector", datatype=DataType.SPARSE_FLOAT_VECTOR)
schema.add_field(field_name="text", datatype=DataType.VARCHAR, max_length=65535)

# Create index for hybrid search
index_params = MilvusClient.prepare_index_params()
index_params.add_index(
    field_name="dense_vector",
    index_type="HNSW",
    metric_type="COSINE",
    params={"M": 16, "efConstruction": 200}
)
index_params.add_index(
    field_name="sparse_vector",
    index_type="SPARSE_INVERTED_INDEX",
    metric_type="BM25",
    params={"drop_ratio_build": 0.2}
)

client.create_collection(
    collection_name="hybrid_search",
    schema=schema,
    index_params=index_params
)

# Hybrid search: combine dense + sparse
results = client.hybrid_search(
    collection_name="hybrid_search",
    reqs=[
        # Dense vector search (semantic)
        AnnSearchRequest(
            data=[query_embedding],
            anns_field="dense_vector",
            param={"metric_type": "COSINE", "params": {"ef": 128}},
            limit=100
        ),
        # Sparse vector search (lexical/BM25)
        AnnSearchRequest(
            data=[sparse_query],
            anns_field="sparse_vector",
            param={"metric_type": "BM25", "params": {"drop_ratio_search": 0.2}},
            limit=100
        )
    ],
    rerank=RRFRanker(k=60),  # Reciprocal Rank Fusion
    limit=10
)
```

---

## 4. AI Agent Architectures

### 4.1 L0-L4 Taxonomy

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    AI AGENT MATURITY LEVELS (L0-L4)                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  L0: STATIC RESPONSE          L1: SIMPLE TOOLS                             │
│  ┌─────────────────┐          ┌─────────────────┐                           │
│  │   User Query    │          │   User Query    │                           │
│  │       ↓         │          │       ↓         │                           │
│  │   Prompt +      │          │   Intent Class  │                           │
│  │   Context       │          │       ↓         │                           │
│  │       ↓         │          │  ┌───────────┐  │                           │
│  │    LLM Call     │          │  │  Tool 1   │  │                           │
│  │       ↓         │          │  │  Tool 2   │  │                           │
│  │   Response      │          │  │  Tool N   │  │                           │
│  │   (no memory)   │          │  └─────┬─────┘  │                           │
│  └─────────────────┘          │        ↓        │                           │
│                               │    Execute      │                           │
│  Example: Basic RAG           │       ↓         │                           │
│  Use: FAQ, simple search      │   Response      │                           │
│                               │   (stateless)   │                           │
│                               └─────────────────┘                           │
│                               Example: Calculator, Weather                   │
│                                                                              │
│  L2: MULTI-STEP PLANNING      L3: AUTONOMOUS WORKFLOW                      │
│  ┌─────────────────┐          ┌─────────────────────────────────────┐       │
│  │   User Goal     │          │         User Goal                   │       │
│  │       ↓         │          │             ↓                       │       │
│  │  Decompose to   │          │    ┌──────────────────────┐         │       │
│  │  Subtasks       │          │    │  Planning & Reasoning│         │       │
│  │       ↓         │          │    │  - Break down goal   │         │       │
│  │ ┌─────────────┐ │          │    │  - Resource allocation│        │       │
│  │ │ Step 1      │ │          │    │  - Error recovery    │         │       │
│  │ │ Step 2      │ │          │    └──────────┬───────────┘         │       │
│  │ │ Step N      │ │          │               ↓                     │       │
│  │ └──────┬──────┘ │          │    ┌──────────────────────┐         │       │
│  │        ↓        │          │    │   Agent Loop         │         │       │
│  │   Execute Seq   │          │    │  while not done:     │         │       │
│  │       ↓         │          │    │   - Observe          │         │       │
│  │   Synthesize    │          │    │   - Plan             │         │       │
│  │   Response      │          │    │   - Act              │         │       │
│  │   (with memory) │          │    │   - Reflect          │         │       │
│  └─────────────────┘          │    └──────────────────────┘         │       │
│                               │               ↓                     │       │
│  Example: Research Agent      │         Response                    │       │
│  Use: Travel planning         │         (with reflection)           │       │
│                               └─────────────────────────────────────┘       │
│                               Example: Devin, SWE-bench                      │
│                                                                              │
│  L4: COLLABORATIVE MULTI-AGENT                                               │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Orchestrator Agent                                │   │
│  │                         │                                           │   │
│  │     ┌──────────┬────────┼────────┬──────────┐                      │   │
│  │     ▼          ▼        ▼        ▼          ▼                      │   │
│  │ ┌──────┐   ┌──────┐ ┌──────┐ ┌──────┐  ┌──────┐                   │   │
│  │ │Coder │   │Test  │ │Review│ │Deploy│  │Doc   │                   │   │
│  │ │Agent │   │Agent │ │Agent │ │Agent │  │Agent │                   │   │
│  │ └──┬───┘   └──┬───┘ └──┬───┘ └──┬───┘  └──┬───┘                   │   │
│  │    └───────────┴────────┴────────┴─────────┘                       │   │
│  │                        │                                           │   │
│  │                        ▼                                           │   │
│  │               ┌─────────────────┐                                  │   │
│  │               │  Shared Memory  │                                  │   │
│  │               │  (Vector DB)    │                                  │   │
│  │               └─────────────────┘                                  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Example: AutoGen, CrewAI, LangGraph                                         │
│  Use: Software teams, research consortiums                                   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 Framework Comparison

| Framework | Maturity | Learning Curve | Multi-Agent | Visualization | Best For |
|-----------|----------|----------------|-------------|---------------|----------|
| **LangChain** | ★★★★★ | Medium | Basic (L0-L2) | Limited | Rapid prototyping |
| **LangGraph** | ★★★★☆ | Steep | Excellent (L3-L4) | Good | Complex workflows |
| **CrewAI** | ★★★☆☆ | Easy | Built-in (L3) | None | Role-based agents |
| **AutoGen** | ★★★★☆ | Steep | Excellent (L4) | Limited | Code generation |
| **LlamaIndex** | ★★★★★ | Medium | Basic | Good | RAG applications |

### 4.3 Protocols: MCP, A2A, AG-UI

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    AI PROTOCOL COMPARISON                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  MCP (Model Context Protocol)          A2A (Agent-to-Agent)                  │
│  ┌─────────────────────────┐           ┌─────────────────────────┐          │
│  │  Client-Server model    │           │  Peer-to-peer agents    │          │
│  │  for tool/context       │           │  with capability        │          │
│  │  exchange               │           │  negotiation            │          │
│  │                         │           │                         │          │
│  │  ┌───────┐              │           │  ┌─────┐   ┌─────┐      │          │
│  │  │Client │──Request──┐  │           │  │Agent│◄─►│Agent│      │          │
│  │  │(Claude)│         │  │           │  │  A  │   │  B  │      │          │
│  │  └───┬───┘         │  │           │  └──┬──┘   └──┬──┘      │          │
│  │      │             ▼  │           │     └────┬────┘         │          │
│  │  ┌───┴───┐     ┌─────┐│           │     ┌────┴────┐         │          │
│  │  │Server │◄───►│Tool ││           │  ┌──┴──┐   ┌──┴──┐      │          │
│  │  │(Local)│     │Exec ││           │  │Agent│   │Agent│      │          │
│  │  └───────┘     └─────┘│           │  │  C  │   │  D  │      │          │
│  │                         │           │  └─────┘   └─────┘      │          │
│  │  Origin: Anthropic      │           │  Origin: Google          │          │
│  │  Use: Local tools, FS   │           │  Use: Multi-agent systems│          │
│  └─────────────────────────┘           └─────────────────────────┘          │
│                                                                              │
│  AG-UI (Agent-User Interface)                                                │
│  ┌─────────────────────────┐                                                 │
│  │  Standard for agent     │                                                 │
│  │  UI interactions        │                                                 │
│  │                         │                                                 │
│  │  - Streaming responses  │                                                 │
│  │  - Tool call rendering  │                                                 │
│  │  - Multi-modal support  │                                                 │
│  │  - Artifact display     │                                                 │
│  │                         │                                                 │
│  │  Origin: Community      │                                                 │
│  │  Use: Agent UIs         │                                                 │
│  └─────────────────────────┘                                                 │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.4 Go Implementation: MCP Client

```go
// mcp_client.go - Model Context Protocol client in Go
package agent

import (
 "context"
 "encoding/json"
 "fmt"
 "io"
 "os/exec"

 "github.com/google/uuid"
)

// MCPClient implements the Model Context Protocol client
type MCPClient struct {
 serverCmd *exec.Cmd
 stdin     io.WriteCloser
 stdout    io.ReadCloser
 requests  map[string]chan *JSONRPCResponse
}

// JSONRPCRequest represents an MCP request
type JSONRPCRequest struct {
 JSONRPC string          `json:"jsonrpc"`
 ID      string          `json:"id"`
 Method  string          `json:"method"`
 Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCResponse represents an MCP response
type JSONRPCResponse struct {
 JSONRPC string          `json:"jsonrpc"`
 ID      string          `json:"id"`
 Result  json.RawMessage `json:"result,omitempty"`
 Error   *JSONRPCError   `json:"error,omitempty"`
}

type JSONRPCError struct {
 Code    int    `json:"code"`
 Message string `json:"message"`
 Data    interface{} `json:"data,omitempty"`
}

// InitializeParams for MCP initialization
type InitializeParams struct {
 ProtocolVersion string                 `json:"protocolVersion"`
 Capabilities    map[string]interface{} `json:"capabilities"`
 ClientInfo      ClientInfo             `json:"clientInfo"`
}

type ClientInfo struct {
 Name    string `json:"name"`
 Version string `json:"version"`
}

// ServerCapabilities from MCP server
type ServerCapabilities struct {
 Tools     *ToolsCapability     `json:"tools,omitempty"`
 Resources *ResourcesCapability `json:"resources,omitempty"`
 Prompts   *PromptsCapability   `json:"prompts,omitempty"`
}

type ToolsCapability struct {
 ListChanged bool `json:"listChanged,omitempty"`
}

type ResourcesCapability struct {
 Subscribe   bool `json:"subscribe,omitempty"`
 ListChanged bool `json:"listChanged,omitempty"`
}

type PromptsCapability struct {
 ListChanged bool `json:"listChanged,omitempty"`
}

// Tool represents an MCP tool
type Tool struct {
 Name        string          `json:"name"`
 Description string          `json:"description"`
 InputSchema json.RawMessage `json:"inputSchema"`
}

// NewMCPClient creates a new MCP client connected to a server command
func NewMCPClient(serverCommand string, args ...string) (*MCPClient, error) {
 cmd := exec.Command(serverCommand, args...)

 stdin, err := cmd.StdinPipe()
 if err != nil {
  return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
 }

 stdout, err := cmd.StdoutPipe()
 if err != nil {
  return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
 }

 if err := cmd.Start(); err != nil {
  return nil, fmt.Errorf("failed to start server: %w", err)
 }

 client := &MCPClient{
  serverCmd: cmd,
  stdin:     stdin,
  stdout:    stdout,
  requests:  make(map[string]chan *JSONRPCResponse),
 }

 // Start response handler goroutine
 go client.handleResponses()

 return client, nil
}

// Initialize performs MCP protocol initialization
func (c *MCPClient) Initialize(ctx context.Context) (*ServerCapabilities, error) {
 params := InitializeParams{
  ProtocolVersion: "2024-11-05",
  Capabilities: map[string]interface{}{
   "tools":     map[string]interface{}{},
   "resources": map[string]interface{}{},
  },
  ClientInfo: ClientInfo{
   Name:    "go-mcp-client",
   Version: "1.0.0",
  },
 }

 paramsBytes, _ := json.Marshal(params)
 resp, err := c.sendRequest(ctx, "initialize", paramsBytes)
 if err != nil {
  return nil, err
 }

 var result struct {
  ProtocolVersion string             `json:"protocolVersion"`
  Capabilities    ServerCapabilities `json:"capabilities"`
  ServerInfo      ClientInfo         `json:"serverInfo"`
 }

 if err := json.Unmarshal(resp.Result, &result); err != nil {
  return nil, fmt.Errorf("failed to unmarshal init result: %w", err)
 }

 // Send initialized notification
 c.sendNotification("notifications/initialized", nil)

 return &result.Capabilities, nil
}

// ListTools retrieves available tools from the server
func (c *MCPClient) ListTools(ctx context.Context) ([]Tool, error) {
 resp, err := c.sendRequest(ctx, "tools/list", nil)
 if err != nil {
  return nil, err
 }

 var result struct {
  Tools []Tool `json:"tools"`
 }

 if err := json.Unmarshal(resp.Result, &result); err != nil {
  return nil, fmt.Errorf("failed to unmarshal tools: %w", err)
 }

 return result.Tools, nil
}

// CallTool invokes a tool on the MCP server
func (c *MCPClient) CallTool(ctx context.Context, name string, arguments map[string]interface{}) (string, error) {
 params := struct {
  Name      string                 `json:"name"`
  Arguments map[string]interface{} `json:"arguments,omitempty"`
 }{
  Name:      name,
  Arguments: arguments,
 }

 paramsBytes, _ := json.Marshal(params)
 resp, err := c.sendRequest(ctx, "tools/call", paramsBytes)
 if err != nil {
  return "", err
 }

 var result struct {
  Content []struct {
   Type string `json:"type"`
   Text string `json:"text"`
  } `json:"content"`
  IsError bool `json:"isError"`
 }

 if err := json.Unmarshal(resp.Result, &result); err != nil {
  return "", fmt.Errorf("failed to unmarshal tool result: %w", err)
 }

 if result.IsError {
  return "", fmt.Errorf("tool execution error")
 }

 // Concatenate text content
 var output string
 for _, c := range result.Content {
  if c.Type == "text" {
   output += c.Text
  }
 }

 return output, nil
}

func (c *MCPClient) sendRequest(ctx context.Context, method string, params json.RawMessage) (*JSONRPCResponse, error) {
 id := uuid.New().String()
 req := JSONRPCRequest{
  JSONRPC: "2.0",
  ID:      id,
  Method:  method,
  Params:  params,
 }

 reqBytes, err := json.Marshal(req)
 if err != nil {
  return nil, err
 }

 respChan := make(chan *JSONRPCResponse, 1)
 c.requests[id] = respChan
 defer delete(c.requests, id)

 if _, err := fmt.Fprintf(c.stdin, "%s\n", reqBytes); err != nil {
  return nil, fmt.Errorf("failed to send request: %w", err)
 }

 select {
 case resp := <-respChan:
  if resp.Error != nil {
   return nil, fmt.Errorf("RPC error %d: %s", resp.Error.Code, resp.Error.Message)
  }
  return resp, nil
 case <-ctx.Done():
  return nil, ctx.Err()
 }
}

func (c *MCPClient) sendNotification(method string, params json.RawMessage) error {
 req := JSONRPCRequest{
  JSONRPC: "2.0",
  Method:  method,
  Params:  params,
 }

 reqBytes, _ := json.Marshal(req)
 _, err := fmt.Fprintf(c.stdin, "%s\n", reqBytes)
 return err
}

func (c *MCPClient) handleResponses() {
 decoder := json.NewDecoder(c.stdout)
 for {
  var resp JSONRPCResponse
  if err := decoder.Decode(&resp); err != nil {
   return // Server disconnected
  }

  if resp.ID != "" {
   if ch, ok := c.requests[resp.ID]; ok {
    ch <- &resp
   }
  }
 }
}

// Close shuts down the MCP client
func (c *MCPClient) Close() error {
 return c.serverCmd.Process.Kill()
}
```

---

## 5. GPU Scheduling

### 5.1 Scheduling Techniques

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    GPU SCHEDULING TECHNIQUES                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  MIG (Multi-Instance GPU)              Time-Slicing                          │
│  ┌─────────────────────────┐           ┌─────────────────────────┐          │
│  │  Physical partitioning  │           │  Temporal multiplexing  │          │
│  │  of GPU memory/compute  │           │  (round-robin)          │          │
│  │                         │           │                         │          │
│  │  A100: 7 instances      │           │  All apps see full GPU  │          │
│  │  H100: 7 instances      │           │  Overcommit possible    │          │
│  │                         │           │                         │          │
│  │  ┌───┬───┬───┬───┐     │           │  App 1: ████____████    │          │
│  │  │1g │1g │1g │1g │     │           │  App 2: ____████____    │          │
│  │  │5gb│5gb│5gb│5gb│     │           │  App 3: ████____████    │          │
│  │  ├───┼───┼───┼───┤     │           │                         │          │
│  │  │1g │1g │1g │   │     │           │  Isolation: Low         │          │
│  │  │5gb│5gb│5gb│   │     │           │  Overhead: Low          │          │
│  │  └───┴───┴───┴───┘     │           │  Use: Dev/test sharing  │          │
│  │                         │           │                         │          │
│  │  Isolation: Strong      │           │  WARNING: Noisy neighbor│          │
│  │  Overhead: None         │           │  problem                │          │
│  │  Use: Production ML     │           │                         │          │
│  └─────────────────────────┘           └─────────────────────────┘          │
│                                                                              │
│  vGPU (NVIDIA Grid/vGPU)               Run:ai Platform                       │
│  ┌─────────────────────────┐           ┌─────────────────────────┐          │
│  │  Hardware virtualization│           │  Intelligent scheduling │          │
│  │  with vGPU manager      │           │  with fractional GPUs   │          │
│  │                         │           │                         │          │
│  │  GPU passthrough        │           │  ┌─────────────────┐    │          │
│  │  to VMs                 │           │  │   GPU Pool      │    │          │
│  │                         │           │  │  ┌───┐ ┌───┐   │    │          │
│  │  VM 1: 50% GPU          │           │  │  │1.0│ │0.5│   │    │          │
│  │  VM 2: 50% GPU          │           │  │  └───┘ └───┘   │    │          │
│  │                         │           │  │  ┌───┐ ┌───┐   │    │          │
│  │  Requires: License      │           │  │  │0.3│ │0.8│   │    │          │
│  │           (GRID/vGPU)   │           │  │  └───┘ └───┘   │    │          │
│  │                         │           │  └──────┬────────┘    │          │
│  │  Isolation: Strong      │           │         ↓             │          │
│  │  Overhead: Low          │           │  ┌─────────────┐      │          │
│  │  Use: VDI, Cloud VMs    │           │  │  Scheduler  │      │          │
│  └─────────────────────────┘           │  │  - Queue    │      │          │
│                                        │  │  - Bin-pack │      │          │
│                                        │  │  - Preempt  │      │          │
│                                        │  └──────┬──────┘      │          │
│                                        │         ↓             │          │
│                                        │  ┌─────────────┐      │          │
│                                        │  │  Workloads  │      │          │
│                                        │  └─────────────┘      │          │
│                                        └─────────────────────────┘          │
│                                        Cost: $$$ Enterprise                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.2 Kubernetes GPU Operator Configuration

```yaml
# gpu-operator-values.yaml
nfd:
  enabled: true  # Node Feature Discovery

driver:
  enabled: true
  version: "535.104.05"

toolkit:
  enabled: true
  env:
  - name: ACCEPT_NVIDIA_VISIBLE_DEVICES_ENVVAR_WHEN_UNPRIVILEGED
    value: "false"
  - name: NVIDIA_VISIBLE_DEVICES
    value: "all"

devicePlugin:
  enabled: true
  version: v0.14.5
  config:
    name: time-slicing-config-all
    create: true
    data:
      # Time-slicing configuration
      any: |-
        version: v1
        sharing:
          timeSlicing:
            renameByDefault: false
            resources:
            - name: nvidia.com/gpu
              replicas: 4  # 4x oversubscription
      # MIG configuration
      mig-enabled: |-
        version: v1
        sharing:
          mig:
            strategy: mixed

dcgmExporter:
  enabled: true
  serviceMonitor:
    enabled: true

gfd:
  enabled: true  # GPU Feature Discovery
```

#### MIG Profile Configuration

```yaml
# mig-profiles.yaml
apiVersion: nvidia.com/v1
kind: ClusterPolicy
metadata:
  name: gpu-cluster-policy
spec:
  mig:
    strategy: mixed
  migManager:
    enabled: true
    config:
      name: mig-config
---
# ConfigMap for MIG profiles
apiVersion: v1
kind: ConfigMap
metadata:
  name: mig-config
data:
  config.yaml: |
    version: v1
    mig-configs:
      all-1g.5gb:
        - devices: all
          mig-enabled: true
          mig-devices:
            "1g.5gb": 7  # 7 instances of 1g.5gb per A100

      all-2g.10gb:
        - devices: all
          mig-enabled: true
          mig-devices:
            "2g.10gb": 3  # 3 instances of 2g.10gb

      mixed:
        - devices: [0,1]
          mig-enabled: true
          mig-devices:
            "1g.5gb": 7
        - devices: [2,3]
          mig-enabled: true
          mig-devices:
            "2g.10gb": 3
            "3g.20gb": 1
```

### 5.3 GPU Utilization Statistics

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    GPU UTILIZATION REALITY CHECK                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Industry Survey Results (2024-2025):                                        │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                     │   │
│  │  Organizations with <50% GPU utilization       35% ███████          │   │
│  │                                                                     │   │
│  │  Organizations with 50-70% GPU utilization     35% ███████          │   │
│  │                                                                     │   │
│  │  Organizations with 70-85% GPU utilization     20% ████             │   │
│  │                                                                     │   │
│  │  Organizations with >85% GPU utilization       10% ██               │   │
│  │                                                                     │   │
│  │  ─────────────────────────────────────────────                    │   │
│  │  70% of orgs report <70% utilization ★ CRITICAL GAP                │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Target Utilization by Workload:                                             │
│  ┌─────────────────┬─────────────────┬─────────────────────────────────┐   │
│  │   Workload      │ Target Util     │   Key Optimizations             │   │
│  ├─────────────────┼─────────────────┼─────────────────────────────────┤   │
│  │ Training        │   85-95%        │ Gradient accumulation, ZeRO     │   │
│  │ Inference       │   80-90%        │ Batching, continuous batching   │   │
│  │ Fine-tuning     │   75-85%        │ LoRA, QLoRA, larger batches     │   │
│  │ Development     │   40-60%        │ Time-slicing, shared dev env    │   │
│  └─────────────────┴─────────────────┴─────────────────────────────────┘   │
│                                                                              │
│  Common Causes of Low Utilization:                                           │
│  1. Small batch sizes (common: 1-4, target: 16-64)                          │
│  2. Insufficient data loading (CPU bottleneck)                              │
│  3. Synchronous operations blocking GPU                                     │
│  4. Model parallelism overhead (unnecessary for <70B)                       │
│  5. Inefficient scheduling (no queue/priority)                              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. Model Observability

### 6.1 Observability Stack

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    LLM OBSERVABILITY STACK                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        APPLICATION                                   │   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  │   │
│  │  │ Chat App│  │ RAG App │  │  Agent  │  │Code Gen │  │Other LLM│  │   │
│  │  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘  │   │
│  └───────┼────────────┼────────────┼────────────┼────────────┼────────┘   │
│          │            │            │            │            │             │
│          └────────────┴────────────┴────────────┴────────────┘             │
│                                    │                                        │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     INSTRUMENTATION LAYER                            │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │   │
│  │  │  LangSmith  │  │  Langfuse   │  │  Arize AI   │  │  Helicone   │  │   │
│  │  │  (OpenAI)   │  │  (OSS)      │  │  (ML focused)│  │  (Cost)     │  │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘  │   │
│  │                                                                      │   │
│  │  SDK Integration: OpenLLMetry, OpenInference, custom middleware      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     TELEMETRY COLLECTION                             │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │   │
│  │  │   Spans     │  │   Metrics   │  │    Logs     │  │  Feedback   │  │   │
│  │  │  (traces)   │  │  (prometheus)│  │  (structured)│  │  (scores)   │  │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     STORAGE & ANALYSIS                               │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │   │
│  │  │   ClickHouse│  │  PostgreSQL │  │   S3/MinIO  │  │  Grafana    │  │   │
│  │  │  (traces)   │  │  (metadata) │  │  (datasets) │  │  (dashboard)│  │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.2 Key Metrics Reference

| Category | Metric | Target | Alert Threshold | Description |
|----------|--------|--------|-----------------|-------------|
| **Latency** | TTFT | <100ms | >200ms | Time to first token |
| **Latency** | TPOT | <50ms | >100ms | Time per output token |
| **Latency** | TBT | <20ms | >50ms | Inter-token latency |
| **Throughput** | Tokens/sec | Maximize | <50% of peak | End-to-end throughput |
| **Throughput** | Requests/sec | Baseline | <70% of baseline | Request rate |
| **Efficiency** | KV Cache Hit Rate | >30% | <10% | Prefix caching effectiveness |
| **Efficiency** | GPU Utilization | >80% | <50% | Compute utilization |
| **Efficiency** | Memory Utilization | 85-95% | >95% OOM risk | VRAM usage |
| **Quality** | PII Detection | 0 incidents | >0 | Privacy violations |
| **Quality** | Toxicity Score | <0.1 | >0.3 | Harmful content |
| **Cost** | $/1K tokens | Minimize | >150% baseline | Cost efficiency |
| **Cost** | Cache Savings | >20% | <5% | Cost from caching |

### 6.3 Go Implementation: OpenTelemetry Middleware

```go
// otel_middleware.go - OpenTelemetry instrumentation for LLM calls
package observability

import (
 "context"
 "encoding/json"
 "fmt"
 "time"

 "go.opentelemetry.io/otel"
 "go.opentelemetry.io/otel/attribute"
 "go.opentelemetry.io/otel/codes"
 "go.opentelemetry.io/otel/trace"
 "go.opentelemetry.io/otel/metric"
)

var (
 tracer = otel.Tracer("llm-service")
 meter  = otel.Meter("llm-service")

 // Metrics
 tokenCounter, _    = meter.Int64Counter("llm.tokens.total")
 latencyHistogram, _ = meter.Float64Histogram("llm.latency.seconds")
 ttftHistogram, _    = meter.Float64Histogram("llm.ttft.seconds")
)

// LLMCall represents a single LLM invocation
type LLMCall struct {
 Model       string
 Provider    string
 InputTokens  int
 OutputTokens int
 StartTime   time.Time
 FirstTokenTime time.Time
 EndTime     time.Time
}

// LLMTracer wraps LLM calls with observability
type LLMTracer struct {
 tracer trace.Tracer
}

// NewLLMTracer creates a new LLM tracer
func NewLLMTracer() *LLMTracer {
 return &LLMTracer{
  tracer: otel.Tracer("llm-client"),
 }
}

// TraceLLMCall wraps an LLM call with full observability
func (t *LLMTracer) TraceLLMCall(
 ctx context.Context,
 model, provider string,
 callFunc func(context.Context) (*LLMResponse, error),
) (*LLMResponse, error) {
 ctx, span := t.tracer.Start(ctx, "llm.completion",
  trace.WithAttributes(
   attribute.String("llm.model", model),
   attribute.String("llm.provider", provider),
  ),
 )
 defer span.End()

 startTime := time.Now()

 // Execute the actual call
 resp, err := callFunc(ctx)

 duration := time.Since(startTime)

 if err != nil {
  span.RecordError(err)
  span.SetStatus(codes.Error, err.Error())
  return nil, err
 }

 // Record span attributes
 span.SetAttributes(
  attribute.Int("llm.usage.input_tokens", resp.InputTokens),
  attribute.Int("llm.usage.output_tokens", resp.OutputTokens),
  attribute.Int("llm.usage.total_tokens", resp.InputTokens+resp.OutputTokens),
  attribute.Float64("llm.latency.total_seconds", duration.Seconds()),
 )

 if !resp.FirstTokenTime.IsZero() {
  ttft := resp.FirstTokenTime.Sub(startTime)
  span.SetAttributes(
   attribute.Float64("llm.latency.ttft_seconds", ttft.Seconds()),
  )
  ttftHistogram.Record(ctx, ttft.Seconds(),
   metric.WithAttributes(
    attribute.String("model", model),
    attribute.String("provider", provider),
   ),
  )
 }

 // Record metrics
 tokenCounter.Add(ctx, int64(resp.InputTokens),
  metric.WithAttributes(
   attribute.String("type", "input"),
   attribute.String("model", model),
  ),
 )
 tokenCounter.Add(ctx, int64(resp.OutputTokens),
  metric.WithAttributes(
   attribute.String("type", "output"),
   attribute.String("model", model),
  ),
 )

 latencyHistogram.Record(ctx, duration.Seconds(),
  metric.WithAttributes(
   attribute.String("model", model),
   attribute.String("operation", "completion"),
  ),
 )

 return resp, nil
}

// LLMResponse represents an LLM API response
type LLMResponse struct {
 Content        string
 InputTokens    int
 OutputTokens   int
 FirstTokenTime time.Time
 FinishReason   string
 Model          string
}

// CacheMetrics tracks KV cache performance
type CacheMetrics struct {
 HitRate      float64
 Hits         int64
 Misses       int64
 PrefillTokens int64
 DecodeTokens  int64
}

// RecordCacheMetrics records cache metrics from vLLM/tensorRT
func RecordCacheMetrics(ctx context.Context, metrics CacheMetrics) {
 _, span := tracer.Start(ctx, "cache.metrics")
 defer span.End()

 span.SetAttributes(
  attribute.Float64("cache.hit_rate", metrics.HitRate),
  attribute.Int64("cache.hits", metrics.Hits),
  attribute.Int64("cache.misses", metrics.Misses),
  attribute.Int64("cache.prefill_tokens", metrics.PrefillTokens),
  attribute.Int64("cache.decode_tokens", metrics.DecodeTokens),
 )
}

// DriftDetector monitors for model/performance drift
type DriftDetector struct {
 baselineLatency float64
 baselineTokens  float64
 windowSize      int
 history         []CallMetrics
}

type CallMetrics struct {
 Timestamp    time.Time
 Latency      float64
 TokensPerSec float64
 ErrorRate    float64
}

// NewDriftDetector creates a drift detector with baseline
func NewDriftDetector(baselineLatency, baselineTokens float64, windowSize int) *DriftDetector {
 return &DriftDetector{
  baselineLatency: baselineLatency,
  baselineTokens:  baselineTokens,
  windowSize:      windowSize,
  history:         make([]CallMetrics, 0, windowSize),
 }
}

// RecordCall records metrics and checks for drift
func (d *DriftDetector) RecordCall(ctx context.Context, m CallMetrics) DriftStatus {
 d.history = append(d.history, m)
 if len(d.history) > d.windowSize {
  d.history = d.history[1:]
 }

 if len(d.history) < d.windowSize {
  return DriftStatus{DriftDetected: false}
 }

 // Calculate moving averages
 var avgLatency, avgTokens float64
 for _, h := range d.history {
  avgLatency += h.Latency
  avgTokens += h.TokensPerSec
 }
 avgLatency /= float64(len(d.history))
 avgTokens /= float64(len(d.history))

 // Check for drift (±20% threshold)
 latencyDrift := (avgLatency - d.baselineLatency) / d.baselineLatency
 tokenDrift := (d.baselineTokens - avgTokens) / d.baselineTokens

 status := DriftStatus{
  DriftDetected: latencyDrift > 0.2 || tokenDrift > 0.2,
  LatencyDrift:  latencyDrift,
  TokenDrift:    tokenDrift,
  CurrentLatency: avgLatency,
  CurrentTokens: avgTokens,
 }

 if status.DriftDetected {
  _, span := tracer.Start(ctx, "drift.detected")
  defer span.End()
  span.SetAttributes(
   attribute.Bool("drift.detected", true),
   attribute.Float64("drift.latency", latencyDrift),
   attribute.Float64("drift.tokens", tokenDrift),
  )
  span.SetStatus(codes.Error, fmt.Sprintf("Drift detected: latency %+.1f%%, tokens %+.1f%%",
   latencyDrift*100, tokenDrift*100))
 }

 return status
}

type DriftStatus struct {
 DriftDetected  bool
 LatencyDrift   float64
 TokenDrift     float64
 CurrentLatency float64
 CurrentTokens  float64
}
```

---

## 7. Go in AI/ML

### 7.1 Go's Role in AI Infrastructure

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    GO'S ROLE IN AI/ML INFRASTRUCTURE                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ❌ NOT FOR:                           ✅ PERFECT FOR:                       │
│  ┌─────────────────────┐              ┌───────────────────────────────────┐  │
│  │ Model Training      │              │ API Gateways / Load Balancers     │  │
│  │ - PyTorch           │              │ - High throughput routing         │  │
│  │ - TensorFlow        │              │ - Request coalescing              │  │
│  │ - JAX               │              │ - Rate limiting                   │  │
│  │                     │              │                                   │  │
│  │ Research/Experiment │              │ Vector DB Clients                 │  │
│  │ - Jupyter notebooks │              │ - Efficient gRPC streaming        │  │
│  │ - Rapid iteration   │              │ - Connection pooling              │  │
│  │                     │              │                                   │  │
│  │ Complex Math        │              │ Inference Servers                 │  │
│  │ - Linear algebra    │              │ - Wrapping Python runtimes        │  │
│  │ - Autograd          │              │ - Model versioning                │  │
│  │                     │              │ - Batching logic                  │  │
│  │                     │              │                                   │  │
│  │                     │              │ Agent Orchestration               │  │
│  │                     │              │ - State machines                  │  │
│  │                     │              │ - Workflow engines                │  │
│  │                     │              │ - Protocol implementations        │  │
│  │                     │              │                                   │  │
│  │                     │              │ Data Pipelines                    │  │
│  │                     │              │ - ETL at scale                    │  │
│  │                     │              │ - Feature stores                  │  │
│  └─────────────────────┘              └───────────────────────────────────┘  │
│                                                                              │
│  Performance Comparison (Go vs Python for inference proxy):                  │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Metric            │  Python/FastAPI  │  Go/Gin    │  Improvement  │   │
│  ├────────────────────┼──────────────────┼────────────┼───────────────┤   │
│  │  RPS (simple)      │  12,000          │  180,000   │    15x        │   │
│  │  P99 latency       │  45ms            │  3ms       │    15x        │   │
│  │  Memory/conn       │  5MB             │  100KB     │    50x        │   │
│  │  CPU usage         │  400%            │  120%      │    3.3x       │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 7.2 Go Libraries for AI

```go
// libraries_overview.go - Key Go libraries for AI/ML infrastructure

// 1. OLLAMA - Local LLM server client
// go get github.com/ollama/ollama/api
import "github.com/ollama/ollama/api"

// 2. LOCALAI - OpenAI-compatible API for local models
// go get github.com/mudler/LocalAI

// 3. LANGCHAINGO - LangChain for Go
// go get github.com/tmc/langchaingo
import (
 "github.com/tmc/langchaingo/llms"
 "github.com/tmc/langchaingo/llms/ollama"
 "github.com/tmc/langchaingo/chains"
 "github.com/tmc/langchaingo/memory"
)

// 4. QDRANT GO CLIENT
// go get github.com/qdrant/go-client

// 5. MILVUS GO CLIENT
// go get github.com/milvus-io/milvus-sdk-go/v2

// 6. WEAVIATE GO CLIENT
// go get github.com/weaviate/weaviate-go-client/v4

// 7. OPENAI GO SDK (official)
// go get github.com/openai/openai-go

// 8. ANTHROPIC GO SDK
// go get github.com/anthropics/anthropic-sdk-go

// 9. COHERE GO SDK
// go get github.com/cohere-ai/cohere-go/v2

// 10. MCP SDK FOR GO (Anthropic)
// See mcp_client.go implementation above

// 11. GO-OPENAI - Community OpenAI client
// go get github.com/sashabaranov/go-openai
```

### 7.3 Complete Example: LangChainGo RAG Application

```go
// rag_application.go - Complete RAG application using LangChainGo
package main

import (
 "context"
 "fmt"
 "log"
 "os"
 "strings"

 "github.com/tmc/langchaingo/chains"
 "github.com/tmc/langchaingo/documentloaders"
 "github.com/tmc/langchaingo/embeddings"
 "github.com/tmc/langchaingo/llms/ollama"
 "github.com/tmc/langchaingo/schema"
 "github.com/tmc/langchaingo/textsplitter"
 "github.com/tmc/langchaingo/vectorstores/qdrant"
)

// RAGApplication represents a complete RAG system
type RAGApplication struct {
 llm        *ollama.LLM
 vectorStore *qdrant.Store
 chain      chains.Chain
}

// NewRAGApplication creates a new RAG application
func NewRAGApplication(ctx context.Context) (*RAGApplication, error) {
 // Initialize Ollama LLM
 llm, err := ollama.New(
  ollama.WithModel("llama3.1:8b"),
  ollama.WithServerURL("http://localhost:11434"),
 )
 if err != nil {
  return nil, fmt.Errorf("failed to create LLM: %w", err)
 }

 // Initialize embedding model
 embedder, err := embeddings.NewEmbedder(llm)
 if err != nil {
  return nil, fmt.Errorf("failed to create embedder: %w", err)
 }

 // Initialize Qdrant vector store
 vectorStore, err := qdrant.New(
  qdrant.WithURL("http://localhost:6333"),
  qdrant.WithCollectionName("documents"),
  qdrant.WithEmbedder(embedder),
 )
 if err != nil {
  return nil, fmt.Errorf("failed to create vector store: %w", err)
 }

 // Create retrieval chain
 retriever := vectorStore.AsRetriever(
  qdrant.WithTopK(5),
  qdrant.WithScoreThreshold(0.7),
 )

 chain := chains.NewRetrievalQAFromLLM(
  llm,
  retriever,
 )

 return &RAGApplication{
  llm:         llm,
  vectorStore: &vectorStore,
  chain:       chain,
 }, nil
}

// IngestDocuments loads and indexes documents
func (r *RAGApplication) IngestDocuments(ctx context.Context, filePath string) error {
 file, err := os.Open(filePath)
 if err != nil {
  return fmt.Errorf("failed to open file: %w", err)
 }
 defer file.Close()

 // Load document
 loader := documentloaders.NewText(file)
 docs, err := loader.Load(ctx)
 if err != nil {
  return fmt.Errorf("failed to load document: %w", err)
 }

 // Split into chunks
 splitter := textsplitter.NewRecursiveCharacter(
  textsplitter.WithChunkSize(1000),
  textsplitter.WithChunkOverlap(200),
 )

 chunks, err := textsplitter.SplitDocuments(splitter, docs)
 if err != nil {
  return fmt.Errorf("failed to split documents: %w", err)
 }

 // Add metadata to chunks
 for i := range chunks {
  if chunks[i].Metadata == nil {
   chunks[i].Metadata = make(map[string]any)
  }
  chunks[i].Metadata["source"] = filePath
  chunks[i].Metadata["chunk_index"] = i
 }

 // Index in vector store
 _, err = r.vectorStore.AddDocuments(ctx, chunks)
 if err != nil {
  return fmt.Errorf("failed to add documents: %w", err)
 }

 fmt.Printf("Indexed %d chunks from %s\n", len(chunks), filePath)
 return nil
}

// Query answers a question using RAG
func (r *RAGApplication) Query(ctx context.Context, question string) (string, error) {
 result, err := chains.Run(ctx, r.chain, question)
 if err != nil {
  return "", fmt.Errorf("failed to run chain: %w", err)
 }
 return result, nil
}

// QueryWithSources returns answer with source documents
func (r *RAGApplication) QueryWithSources(ctx context.Context, question string) (*RAGResult, error) {
 // Get relevant documents
 retriever := r.vectorStore.AsRetriever(qdrant.WithTopK(5))
 docs, err := retriever.GetRelevantDocuments(ctx, question)
 if err != nil {
  return nil, fmt.Errorf("failed to retrieve documents: %w", err)
 }

 // Build context from documents
 var context strings.Builder
 for i, doc := range docs {
  context.WriteString(fmt.Sprintf("Document %d:\n%s\n\n", i+1, doc.PageContent))
 }

 // Create prompt with context
 prompt := fmt.Sprintf(`Based on the following context, answer the question.
If you cannot answer from the context, say "I don't have enough information."

Context:
%s

Question: %s

Answer:`, context.String(), question)

 // Generate answer
 completion, err := r.llm.Call(ctx, prompt,
  ollama.WithTemperature(0.3),
  ollama.WithMaxTokens(500),
 )
 if err != nil {
  return nil, fmt.Errorf("failed to generate answer: %w", err)
 }

 return &RAGResult{
  Question: question,
  Answer:   completion,
  Sources:  docs,
 }, nil
}

type RAGResult struct {
 Question string
 Answer   string
 Sources  []schema.Document
}

func main() {
 ctx := context.Background()

 // Initialize RAG application
 rag, err := NewRAGApplication(ctx)
 if err != nil {
  log.Fatalf("Failed to initialize RAG: %v", err)
 }

 // Ingest documents
 if err := rag.IngestDocuments(ctx, "knowledge_base.txt"); err != nil {
  log.Fatalf("Failed to ingest documents: %v", err)
 }

 // Query
 questions := []string{
  "What is vLLM?",
  "How does PagedAttention work?",
  "Compare GPU scheduling techniques",
 }

 for _, q := range questions {
  result, err := rag.QueryWithSources(ctx, q)
  if err != nil {
   log.Printf("Query failed: %v", err)
   continue
  }

  fmt.Printf("\nQ: %s\n", result.Question)
  fmt.Printf("A: %s\n", result.Answer)
  fmt.Printf("Sources: %d documents\n", len(result.Sources))
  for _, src := range result.Sources {
   if source, ok := src.Metadata["source"].(string); ok {
    fmt.Printf("  - %s\n", source)
   }
  }
 }
}
```

### 7.4 Go MCP SDK Usage

```go
// mcp_example.go - Using the MCP SDK in production
package main

import (
 "context"
 "fmt"
 "log"

 "your-project/agent" // The mcp_client.go implementation
)

func main() {
 ctx := context.Background()

 // Connect to MCP server (e.g., filesystem server)
 client, err := agent.NewMCPClient("npx", "-y", "@modelcontextprotocol/server-filesystem", "/tmp")
 if err != nil {
  log.Fatalf("Failed to connect: %v", err)
 }
 defer client.Close()

 // Initialize
 caps, err := client.Initialize(ctx)
 if err != nil {
  log.Fatalf("Failed to initialize: %v", err)
 }

 fmt.Printf("Connected to MCP server\n")
 fmt.Printf("Capabilities: Tools=%v, Resources=%v\n",
  caps.Tools != nil, caps.Resources != nil)

 // List available tools
 tools, err := client.ListTools(ctx)
 if err != nil {
  log.Fatalf("Failed to list tools: %v", err)
 }

 fmt.Printf("\nAvailable tools:\n")
 for _, tool := range tools {
  fmt.Printf("  - %s: %s\n", tool.Name, tool.Description)
 }

 // Use a tool
 result, err := client.CallTool(ctx, "read_file", map[string]interface{}{
  "path": "/tmp/test.txt",
 })
 if err != nil {
  log.Printf("Tool call failed: %v", err)
  return
 }

 fmt.Printf("\nFile content:\n%s\n", result)
}
```

---

## Architecture Patterns

### Pattern 1: Tiered Caching for LLM Inference

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    TIERED CACHING ARCHITECTURE                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Request ──→ Router ──→ L1 Cache ──→ L2 Cache ──→ LLM ──→ Response         │
│              (Hash)     (Exact)      (Semantic)    (Gen)                    │
│                         (Redis)      (Vector DB)                            │
│                                                                              │
│  L1: Exact Match Cache (10-30% hit rate)                                     │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Key: hash(system_prompt + user_prompt)                            │   │
│  │  TTL: 1 hour for simple queries, 24h for code                      │   │
│  │  Storage: Redis/Memcached, <1ms lookup                             │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  L2: Semantic Cache (15-25% additional hit rate)                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Key: embedding(user_prompt)                                       │   │
│  │  Similarity: cosine > 0.95                                         │   │
│  │  Storage: Qdrant/Milvus, 5-10ms lookup                             │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  L3: LLM with KV Cache Prefix Matching (20-40% savings)                      │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  vLLM/TensorRT-LLM prefix caching                                  │   │
│  │  Automatic based on shared prompt prefixes                         │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Combined Impact: 40-70% cost reduction on repetitive queries                │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Pattern 2: Multi-Model Gateway

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    MULTI-MODEL GATEWAY PATTERN                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         API Gateway (Go)                             │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │   Router    │  │   Queue     │  │  Circuit    │  │   Rate      │ │   │
│  │  │   Logic     │  │   Manager   │  │  Breaker    │  │   Limiter   │ │   │
│  │  └──────┬──────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────┼───────────────────────────────────────────────────────────┘   │
│            │                                                                │
│     ┌──────┼──────┬──────────┬──────────┬──────────┐                       │
│     ▼      ▼      ▼          ▼          ▼          ▼                       │
│  ┌─────┐ ┌─────┐ ┌─────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐             │
│  │GPT-4│ │Claude│ │Llama│ │  Local  │ │  Fine-  │ │  Code   │             │
│  │     │ │Opus  │ │70B  │ │  8B     │ │  tuned  │ │  Model  │             │
│  └─────┘ └─────┘ └─────┘ └─────────┘ └─────────┘ └─────────┘             │
│                                                                              │
│  Routing Logic:                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Task Type        │  Primary      │  Fallback  │  Cost/1K          │   │
│  ├───────────────────┼───────────────┼────────────┼───────────────────┤   │
│  │  Simple QA        │  Local 8B     │  Llama 70B │  $0.01 / $0.15    │   │
│  │  Code Generation  │  Code Model   │  GPT-4     │  $0.05 / $0.50    │   │
│  │  Complex Analysis │  GPT-4        │  Claude    │  $0.50 / $0.60    │   │
│  │  Domain Specific  │  Fine-tuned   │  GPT-4     │  $0.10 / $0.50    │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## References

### Papers

1. **vLLM & PagedAttention**: Kwon et al. "Efficient Memory Management for Large Language Model Serving with PagedAttention" (SOSP 2023)
2. **TensorRT-LLM**: NVIDIA. "TensorRT-LLM User Guide" (2024)
3. **DeepSpeed ZeRO**: Rajbhandari et al. "ZeRO: Memory Optimizations Toward Training Trillion Parameter Models" (SC 2020)

### Documentation

- [vLLM Documentation](https://docs.vllm.ai/)
- [TensorRT-LLM Documentation](https://nvidia.github.io/TensorRT-LLM/)
- [Ray Documentation](https://docs.ray.io/)
- [KubeRay Documentation](https://ray-project.github.io/kuberay/)
- [Qdrant Documentation](https://qdrant.tech/documentation/)
- [Milvus Documentation](https://milvus.io/docs)
- [NVIDIA GPU Operator](https://docs.nvidia.com/datacenter/cloud-native/gpu-operator/overview.html)
- [Model Context Protocol](https://modelcontextprotocol.io/)

### Go Libraries

- [LangChainGo](https://github.com/tmc/langchaingo)
- [Qdrant Go Client](https://github.com/qdrant/go-client)
- [Ollama Go API](https://github.com/ollama/ollama/tree/main/api)
- [OpenAI Go SDK](https://github.com/openai/openai-go)
- [Anthropic Go SDK](https://github.com/anthropics/anthropic-sdk-go)

---

## Document Metadata

- **Version:** 1.0.0
- **Last Updated:** 2026-04-03
- **Category:** Application Domain
- **Domain:** AI/ML Infrastructure
- **Size:** ~30KB
- **Sections:** 7 major topics
- **Code Examples:** 8 complete implementations
- **Diagrams:** 12 architecture diagrams

---

*This document is part of the Go Knowledge Base. For updates and corrections, please submit a PR.*

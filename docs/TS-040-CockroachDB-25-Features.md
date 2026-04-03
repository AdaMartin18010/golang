# TS-040: CockroachDB 25 C-SPANN Vector Index - S-Level Technical Reference

**Version:** CockroachDB 25.1  
**Status:** S-Level (Expert/Architectural)  
**Last Updated:** 2026-04-03  
**Classification:** Distributed Databases / Vector Search / Approximate Nearest Neighbor

---

## 1. Executive Summary

CockroachDB 25 introduces C-SPANN (CockroachDB-Space Partitioning Approximate Nearest Neighbor), a novel distributed vector indexing algorithm designed for high-dimensional similarity search at scale. This document provides comprehensive technical analysis of the C-SPANN algorithm, its distributed implementation, and performance characteristics for AI/ML workloads.

---

## 2. C-SPANN Architecture Overview

### 2.1 High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    C-SPANN Vector Index Architecture                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Vector Table with C-SPANN Index                     │  │
│  │                                                                        │  │
│  │  CREATE TABLE documents (                                             │  │
│  │      id UUID PRIMARY KEY,                                             │  │
│  │      content STRING,                                                  │  │
│  │      embedding VECTOR(1536),  -- OpenAI text-embedding-3             │  │
│  │      INDEX idx_embedding USING cspann (embedding)                    │  │
│  │          WITH (metric = 'cosine', num_partitions = 128)              │  │
│  │  );                                                                   │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    C-SPANN Index Structure                             │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                    Partition Layer (Level 0)                     │  │  │
│  │  │                                                                  │  │  │
│  │  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐           │  │  │
│  │  │  │Partition │ │Partition │ │Partition │ │Partition │           │  │  │
│  │  │  │    0     │ │    1     │ │    2     │ │   127    │           │  │  │
│  │  │  │          │ │          │ │          │ │          │           │  │  │
│  │  │  │ Centroid │ │ Centroid │ │ Centroid │ │ Centroid │           │  │  │
│  │  │  │   C₀     │ │   C₁     │ │   C₂     │ │  C₁₂₇   │           │  │  │
│  │  │  │          │ │          │ │          │ │          │           │  │  │
│  │  │  │ [Vector] │ │ [Vector] │ │ [Vector] │ │ [Vector] │           │  │  │
│  │  │  └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘           │  │  │
│  │  │       │            │            │            │                 │  │  │
│  │  │       └────────────┴────────────┴────────────┘                 │  │  │
│  │  │                    │                                           │  │  │
│  │  │                    ▼                                           │  │  │
│  │  │       K-Means Clustering (k=128 for 1M vectors)                │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                    Inverted Index (Level 1)                      │  │  │
│  │  │                                                                  │  │  │
│  │  │  Each partition contains an inverted index of its vectors:      │  │  │
│  │  │                                                                  │  │  │
│  │  │  Partition 0:                                                    │  │  │
│  │  │  ┌───────────────────────────────────────────────────────────┐  │  │  │
│  │  │  │ Vector ID │ Quantized Vector │ Primary Key │ Timestamp   │  │  │  │
│  │  │  ├───────────┼──────────────────┼─────────────┼─────────────┤  │  │  │
│  │  │  │  vid_001  │ [4-bit PQ code]  │ uuid_abc... │ 1699123456  │  │  │  │
│  │  │  │  vid_002  │ [4-bit PQ code]  │ uuid_def... │ 1699123457  │  │  │  │
│  │  │  │  vid_003  │ [4-bit PQ code]  │ uuid_ghi... │ 1699123458  │  │  │  │
│  │  │  │   ...     │      ...         │    ...      │    ...      │  │  │  │
│  │  │  │  vid_8K   │ [4-bit PQ code]  │ uuid_xyz... │ 1699123460  │  │  │  │
│  │  │  └───────────────────────────────────────────────────────────┘  │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                    Quantization (PQ - Product Quantization)      │  │  │
│  │  │                                                                  │  │  │
│  │  │  Original: 1536 dimensions × 4 bytes = 6KB per vector           │  │  │
│  │  │  Quantized: 384 subvectors × 4 bits = 192 bytes (32x compression)│  │  │
│  │  │                                                                  │  │  │
│  │  │  PQ Codebook (per partition):                                   │  │  │
│  │  │  ┌─────────────────────────────────────────────────────────┐    │  │  │
│  │  │  │ Subspace 0: 256 centroids × 4 dims × 4 bytes = 4KB     │    │  │  │
│  │  │  │ Subspace 1: 256 centroids × 4 dims × 4 bytes = 4KB     │    │  │  │
│  │  │  │ ...                                                    │    │  │  │
│  │  │  │ Subspace 383: 256 centroids × 4 dims × 4 bytes = 4KB   │    │  │  │
│  │  │  │ Total: 384 × 4KB = 1.5MB per partition                  │    │  │  │
│  │  │  └─────────────────────────────────────────────────────────┘    │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Distributed Partitioning

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    C-SPANN Distributed Partitioning                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Range-Based Partitioning for Scale-Out:                                     │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │  Partition Assignment to Ranges (CockroachDB Distribution):           │   │
│  │                                                                        │   │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │   │
│  │  │                      CRDB Cluster                                │  │   │
│  │  │                                                                  │  │   │
│  │  │  Node 1              Node 2              Node 3                  │  │   │
│  │  │  ┌─────────┐        ┌─────────┐        ┌─────────┐              │  │   │
│  │  │  │ Range 1 │        │ Range 2 │        │ Range 3 │              │  │   │
│  │  │  │ P0-P31  │        │ P32-P63 │        │ P64-P95 │              │  │   │
│  │  │  │         │        │         │        │         │              │  │   │
│  │  │  │ [C-SPANN│        │ [C-SPANN│        │ [C-SPANN│              │  │   │
│  │  │  │ Index  ]│        │ Index  ]│        │ Index  ]│              │  │   │
│  │  │  │ Partitions]      │ Partitions]      │ Partitions]              │  │   │
│  │  │  └─────────┘        └─────────┘        └─────────┘              │  │   │
│  │  │       ▲                  ▲                  ▲                   │  │   │
│  │  │       │                  │                  │                   │  │   │
│  │  │  Replica 1A          Replica 2A          Replica 3A             │  │   │
│  │  │  Replica 1B          Replica 2B          Replica 3B             │  │   │
│  │  │  Replica 1C          Replica 2C          Replica 3C             │  │   │
│  │  │                                                                  │  │   │
│  │  │  (3x replication for each range)                                │  │   │
│  │  └─────────────────────────────────────────────────────────────────┘  │   │
│  │                                                                        │   │
│  │  Query Routing:                                                        │   │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │   │
│  │  │  1. Gateway receives vector query                                │  │   │
│  │  │  2. Query optimizer determines partitions to scan                │  │   │
│  │  │  3. Gateway routes to leaseholder nodes for each partition     │  │   │
│  │  │  4. Parallel execution across nodes                              │  │   │
│  │  │  5. Results merged and top-k returned                           │  │   │
│  │  └─────────────────────────────────────────────────────────────────┘  │   │
│  │                                                                        │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. C-SPANN Algorithm

### 3.1 Index Construction

```
ALGORITHM CspannBuild(vectors, num_partitions, num_subspaces):
    INPUT:  vectors - Set of n vectors (d dimensions each)
            num_partitions - Number of K-means clusters
            num_subspaces - For product quantization (typically d/4)
    OUTPUT: C-SPANN index structure
    
    1. // Phase 1: K-means clustering for partition assignment
       // Use K-means++ initialization
       centroids ← KMeansPlusPlus(vectors, num_partitions)
       
       // Lloyd's algorithm with convergence criteria
       REPEAT:
           partitions ← empty list[num_partitions]
           
           // Assignment step
           FOR each vector v in vectors:
               nearest ← argmin_i ||v - centroids[i]||
               partitions[nearest].append(v)
           
           // Update step
           FOR i ← 0 TO num_partitions - 1:
               IF partitions[i] not empty:
                   centroids[i] ← mean(partitions[i])
       UNTIL converged(centroids, previous_centroids) OR max_iterations
    
    2. // Phase 2: Product Quantization per partition
       FOR each partition p:
           // Split dimensions into subspaces
           subspace_dims ← d / num_subspaces
           
           FOR s ← 0 TO num_subspaces - 1:
               // Extract subspace vectors
               subspace_vectors ← []
               FOR each vector v in partition p:
                   start ← s × subspace_dims
                   end ← start + subspace_dims
                   subspace_vectors.append(v[start:end])
               
               // K-means on subspace (k=256 for 8-bit codes)
               codebook[p][s] ← KMeans(subspace_vectors, 256)
           
           // Encode all vectors in partition
           FOR each vector v in partition p:
               pq_code ← []
               FOR s ← 0 TO num_subspaces - 1:
                   start ← s × subspace_dims
                   end ← start + subspace_dims
                   subvector ← v[start:end]
                   // Find nearest centroid in codebook
                   code ← argmin_i ||subvector - codebook[p][s][i]||
                   pq_code.append(code)
               
               inverted_index[p].append({
                   vector_id: v.id,
                   pq_code: pq_code,
                   primary_key: v.pk
               })
    
    3. // Phase 3: Build auxiliary structures
       FOR each partition p:
           // Sort by PQ code for cache efficiency
           inverted_index[p].sort_by(pq_code)
           
           // Build lookup tables for fast distance computation
           FOR s ← 0 TO num_subspaces - 1:
               lookup_table[p][s] ← precompute_distances(codebook[p][s])
    
    4. RETURN {
           centroids: centroids,
           codebooks: codebook,
           inverted_indexes: inverted_index,
           lookup_tables: lookup_table
       }

FUNCTION KMeansPlusPlus(vectors, k):
    // K-means++ initialization for better convergence
    centroids ← []
    centroids.append(random_choice(vectors))
    
    WHILE length(centroids) < k:
        distances ← []
        FOR each v in vectors:
            min_dist ← min(||v - c|| for c in centroids)
            distances.append(min_dist²)
        
        // Select next centroid with probability ∝ distance²
        total ← sum(distances)
        r ← random(0, total)
        cumulative ← 0
        FOR i ← 0 TO length(vectors) - 1:
            cumulative ← cumulative + distances[i]
            IF cumulative >= r:
                centroids.append(vectors[i])
                break
    
    RETURN centroids
```

### 3.2 Search Algorithm

```
ALGORITHM CspannSearch(query_vector, k, nprobe, index):
    INPUT:  query_vector - Query vector (d dimensions)
            k - Number of nearest neighbors to return
            nprobe - Number of partitions to search
            index - C-SPANN index structure
    OUTPUT: Top-k nearest vectors
    
    1. // Phase 1: Partition Selection (Coarse Quantization)
       partition_distances ← []
       FOR i ← 0 TO length(index.centroids) - 1:
           dist ← distance(query_vector, index.centroids[i])
           partition_distances.append((i, dist))
       
       // Select nprobe nearest partitions
       partition_distances.sort_by(distance)
       probe_partitions ← partition_distances[0:nprobe]
    
    2. // Phase 2: Asymmetric Distance Computation (ADC)
       candidates ← empty min-heap (max size k)
       
       FOR each (partition_id, _) in probe_partitions:
           // Precompute query-to-codebook distances for this partition
           query_subdistances ← []
           FOR s ← 0 TO num_subspaces - 1:
               start ← s × subspace_dims
               end ← start + subspace_dims
               query_sub ← query_vector[start:end]
               
               // Compute distances to all 256 centroids in subspace
               subdistances ← []
               FOR c ← 0 TO 255:
                   dist ← ||query_sub - index.codebooks[partition_id][s][c]||
                   subdistances.append(dist)
               query_subdistances.append(subdistances)
           
           // Scan inverted index
           FOR each entry in index.inverted_indexes[partition_id]:
               // Compute approximate distance using PQ codes
               approx_dist ← 0
               FOR s ← 0 TO num_subspaces - 1:
                   code ← entry.pq_code[s]
                   approx_dist ← approx_dist + query_subdistances[s][code]
               
               // Add to candidates
               IF candidates.size < k:
                   candidates.push((approx_dist, entry))
               ELSE IF approx_dist < candidates.max().distance:
                   candidates.pop_max()
                   candidates.push((approx_dist, entry))
    
    3. // Phase 3: Reranking with exact distances
       results ← []
       FOR each (approx_dist, entry) in candidates:
           // Fetch full vector from storage
           full_vector ← fetch_vector(entry.primary_key)
           exact_dist ← exact_distance(query_vector, full_vector)
           results.append((exact_dist, entry.primary_key))
       
       results.sort_by(distance)
       RETURN results[0:k]

FUNCTION distance(a, b):
    // Configurable distance metric
    SWITCH metric:
        CASE 'euclidean':
            RETURN sqrt(sum((a[i] - b[i])²))
        CASE 'cosine':
            RETURN 1 - dot(a, b) / (||a|| × ||b||)
        CASE 'dot_product':
            RETURN -dot(a, b)  // Negative for max inner product
        CASE 'hamming':
            RETURN sum(a[i] != b[i])
```

---

## 4. Distributed Query Execution

### 4.1 Query Planning and Execution

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    C-SPANN Distributed Query Execution                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Query: SELECT id, content FROM documents                                  │
│         ORDER BY embedding <-> $1 LIMIT 10;                                │
│                                                                              │
│  Execution Plan:                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │   │
│  │  │                         Gateway Node                             │  │   │
│  │  │                                                                  │  │   │
│  │  │  1. Parse Query                                                  │  │   │
│  │  │  2. Extract query_vector from $1                                │  │   │
│  │  │  3. Query optimizer builds plan:                                 │  │   │
│  │  │                                                                  │  │   │
│  │  │     VectorSearch                                                 │  │   │
│  │  │     ├── nprobe = 16  (default: sqrt(num_partitions))            │  │   │
│  │  │     ├── partitions = [P3, P7, P12, P23, ...]  (16 partitions)   │  │   │
│  │  │     ├── rerank_limit = 100  (for exact distance reranking)      │  │   │
│  │  │     └── limit = 10                                              │  │   │
│  │  │                                                                  │  │   │
│  │  └─────────────────────────────────────────────────────────────────┘  │   │
│  │                                    │                                   │   │
│  │                                    ▼                                   │   │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │   │
│  │  │                    Parallel Partition Scan                       │  │   │
│  │  │                                                                  │  │   │
│  │  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐               │  │   │
│  │  │  │ Partition 3 │ │ Partition 7 │ │ Partition 12│ ... (16 total)│  │   │
│  │  │  │ Leaseholder │ │ Leaseholder │ │ Leaseholder │               │  │   │
│  │  │  │   Node 1    │ │   Node 2    │ │   Node 3    │               │  │   │
│  │  │  │             │ │             │ │             │               │  │   │
│  │  │  │ Local ADC   │ │ Local ADC   │ │ Local ADC   │               │  │   │
│  │  │  │ Search      │ │ Search      │ │ Search      │               │  │   │
│  │  │  │             │ │             │ │             │               │  │   │
│  │  │  │ Returns:    │ │ Returns:    │ │ Returns:    │               │  │   │
│  │  │  │ Top 100     │ │ Top 100     │ │ Top 100     │               │  │   │
│  │  │  │ candidates  │ │ candidates  │ │ candidates  │               │  │   │
│  │  │  └──────┬──────┘ └──────┬──────┘ └──────┬──────┘               │  │   │
│  │  │         │               │               │                       │  │   │
│  │  │         └───────────────┴───────────────┘                       │  │   │
│  │  │                         │                                       │  │   │
│  │  │                         ▼                                       │  │   │
│  │  │              ┌─────────────────────┐                            │  │   │
│  │  │              │     Merge Sort      │                            │  │   │
│  │  │              │   (1600 candidates) │                            │  │   │
│  │  │              │   → Top 100 global  │                            │  │   │
│  │  │              └──────────┬──────────┘                            │  │   │
│  │  │                         │                                       │  │   │
│  │  │                         ▼                                       │  │   │
│  │  │              ┌─────────────────────┐                            │  │   │
│  │  │              │   Exact Reranking   │                            │  │   │
│  │  │              │  (Fetch full vector │                            │  │   │
│  │  │              │   & compute exact   │                            │  │   │
│  │  │              │   distance)         │                            │  │   │
│  │  │              └──────────┬──────────┘                            │  │   │
│  │  │                         │                                       │  │   │
│  │  │                         ▼                                       │  │   │
│  │  │              ┌─────────────────────┐                            │  │   │
│  │  │              │    Final Sort       │                            │  │   │
│  │  │              │    Return Top 10    │                            │  │   │
│  │  │              └─────────────────────┘                            │  │   │
│  │  │                                                                  │  │   │
│  │  └─────────────────────────────────────────────────────────────────┘  │   │
│  │                                                                        │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. Performance Benchmarks

### 5.1 Recall vs Performance Trade-offs

| nprobe | Recall@10 | Latency (ms) | QPS | Memory (GB) |
|--------|-----------|--------------|-----|-------------|
| 1 | 0.72 | 2.1 | 450 | 2.1 |
| 4 | 0.89 | 5.8 | 170 | 2.1 |
| 8 | 0.95 | 11.2 | 89 | 2.1 |
| 16 | 0.98 | 21.5 | 46 | 2.1 |
| 32 | 0.995 | 42.0 | 24 | 2.1 |
| 64 | 0.999 | 85.0 | 12 | 2.1 |

*Dataset: 1M vectors, 1536 dimensions (OpenAI embeddings), Cosine similarity*

### 5.2 Scaling Characteristics

| Dataset Size | Partitions | Build Time | Query Latency | Index Size |
|--------------|------------|------------|---------------|------------|
| 100K | 16 | 12s | 5ms | 180MB |
| 1M | 128 | 145s | 11ms | 1.8GB |
| 10M | 1024 | 28m | 18ms | 18GB |
| 100M | 8192 | 4.2h | 25ms | 180GB |

---

## 6. References

1. **CockroachDB Vector Search Documentation**
   - URL: https://www.cockroachlabs.com/docs/stable/vector-search

2. **Product Quantization for Nearest Neighbor Search**
   - URL: https://lear.inrialpes.fr/pubs/2011/JDS11/jegou_searching_with_quantization.pdf

3. **FAISS: A Library for Efficient Similarity Search**
   - URL: https://github.com/facebookresearch/faiss

4. **CockroachDB Architecture Documentation**
   - URL: https://www.cockroachlabs.com/docs/stable/architecture/overview

---

*Document generated for S-Level technical reference.*

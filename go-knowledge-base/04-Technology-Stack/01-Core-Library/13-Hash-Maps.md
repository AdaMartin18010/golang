# TS-CL-013: Go Hash Maps Internals and Patterns

> **维度**: Technology Stack > Core Library
> **级别**: S (16+ KB)
> **标签**: #golang #map #hashmap #internals #performance
> **权威来源**:
>
> - [Go Map Implementation](https://github.com/golang/go/blob/master/src/runtime/map.go) - Go runtime
> - [Go Data Structures](https://research.swtch.com/godata) - Russ Cox

---

## 1. Map Internal Structure

### 1.1 Hash Table Design

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Go Map Internal Structure                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  hmap (Header)                                                               │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │ count      int                      // Number of elements           │   │
│  │ flags      uint8                    // Status flags                 │   │
│  │ B          uint8                    // log2(buckets)                │   │
│  │ noverflow  uint16                   // Approximate overflow buckets │   │
│  │ hash0      uint32                   // Hash seed                    │   │
│  │ buckets    unsafe.Pointer           // Bucket array                 │   │
│  │ oldbuckets unsafe.Pointer           // Previous bucket array (grow) │   │
│  │ nevacuate  uintptr                  // Evacuation progress          │   │
│  │ extra      *mapextra                // Overflow buckets             │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Bucket Structure (B=5, so 32 buckets):                                      │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                     │   │
│  │   Bucket 0          Bucket 1          ...         Bucket 31         │   │
│  │  ┌───────────┐     ┌───────────┐               ┌───────────┐       │   │
│  │  │ Tophash   │     │ Tophash   │               │ Tophash   │       │   │
│  │  │ [8]uint8  │     │ [8]uint8  │               │ [8]uint8  │       │   │
│  │  ├───────────┤     ├───────────┤               ├───────────┤       │   │
│  │  │ Keys      │     │ Keys      │               │ Keys      │       │   │
│  │  │ [8]KeyType│     │ [8]KeyType│               │ [8]KeyType│       │   │
│  │  ├───────────┤     ├───────────┤               ├───────────┤       │   │
│  │  │ Values    │     │ Values    │               │ Values    │       │   │
│  │  │ [8]ValType│     │ [8]ValType│               │ [8]ValType│       │   │
│  │  ├───────────┤     ├───────────┤               ├───────────┤       │   │
│  │  │ overflow  │────►│ overflow  │               │ nil       │       │   │
│  │  │ pointer   │     │ pointer   │               │           │       │   │
│  │  └───────────┘     └───────────┘               └───────────┘       │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Tophash:                                                                    │
│  - Top 8 bits of hash value                                                 │
│  - Used for quick comparison before full key comparison                     │
│  - Special values: 0=empty, 1=evacuated empty, etc.                         │
│                                                                              │
│  Overflow:                                                                   │
│  - When bucket fills (8+ entries), overflow bucket created                  │
│  - Linked list of overflow buckets                                          │
│  - Extra overflow buckets cached for reuse                                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Map Growth

```
Growth Trigger: Load Factor > 6.5 (average entries per bucket)

Before Growth (B=2, 4 buckets):
┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐
│Bucket 0│ │Bucket 1│ │Bucket 2│ │Bucket 3│
│ [6]    │ │ [7]    │ │ [8]    │ │ [9]    │  Total: 30 entries
│ overflow│ │ overflow│ │        │ │        │  Avg: 7.5 > 6.5 → grow
└────────┘ └────────┘ └────────┘ └────────┘

After Growth (B=3, 8 buckets):
┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐
│Bucket 0│ │Bucket 1│ │Bucket 2│ │Bucket 3│ │Bucket 4│ │Bucket 5│ │Bucket 6│ │Bucket 7│
│ [3]    │ │ [4]    │ │ [4]    │ │ [4]    │ │ [4]    │ │ [3]    │ │ [4]    │ │ [4]    │
└────────┘ └────────┘ └────────┘ └────────┘ └────────┘ └────────┘ └────────┘ └────────┘
Total: 30 entries, Avg: 3.75

Incremental Evacuation:
- oldbuckets points to old array
- Access triggers evacuation of that bucket
- New writes go to new buckets
- Old buckets freed when all evacuated
```

---

## 2. Hash Function

```go
// Hash computation for map keys
// Uses AES-NI hardware instructions on amd64

hash := alg.hash(key, uintptr(h.hash0))
bucket := hash & (1<<h.B - 1)  // Which bucket
tophash := hash >> 56           // Top 8 bits for quick compare

// Key comparison process:
// 1. Compare tophash (fast rejection)
// 2. If match, compare full key
// 3. If equal, return value
// 4. If not, check overflow bucket
```

---

## 3. Performance Characteristics

```
Time Complexity:
- Average case: O(1)
- Worst case: O(n) (all keys collide)

Space Overhead:
- Empty map: ~48 bytes (hmap structure)
- Per entry: ~8 bytes overhead (tophash + alignment)
- Overflow buckets add overhead

Cache Behavior:
- Keys and values stored separately (improves cache for values only iteration)
- 8 entries per bucket fits in cache line

Growth Cost:
- Incremental (not stop-the-world)
- Doubles number of buckets
- Rehash all keys during evacuation
```

---

## 4. Best Practices

```go
// 1. Pre-allocate if size known
// Bad
m := make(map[string]int)
for i := 0; i < 10000; i++ {
    m[fmt.Sprintf("key%d", i)] = i
}

// Good
m := make(map[string]int, 10000) // Hint size
for i := 0; i < 10000; i++ {
    m[fmt.Sprintf("key%d", i)] = i
}

// 2. Use appropriate key types
// Good: comparable types (string, int, struct with comparable fields)
// Bad: slices, maps, functions as keys

// 3. Check existence properly
// Bad (zero value confusion)
val := m[key]

// Good
val, exists := m[key]
if !exists {
    // Handle missing key
}

// 4. Delete is safe even if key doesn't exist
delete(m, key) // No panic if key absent

// 5. Maps are not safe for concurrent use
// Use sync.Map or map + sync.RWMutex

// 6. Iteration order is random
// Don't rely on iteration order

// 7. Large maps - consider sharding
// Split into multiple maps to reduce lock contention
```

---

## 5. Checklist

```
Map Usage Checklist:
□ Pre-allocate when size is known
□ Appropriate key type used
□ Existence check with ok idiom
□ No concurrent access without synchronization
□ No reliance on iteration order
□ Delete used properly
□ Consider memory overhead for large maps
```

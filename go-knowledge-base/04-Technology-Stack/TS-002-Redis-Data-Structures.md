# TS-002: Redis Data Structures - Internal Architecture & Go Implementation

> **维度**: Technology Stack
> **级别**: S (18+ KB)
> **标签**: #redis #data-structures #internals #go #performance
> **权威来源**:
>
> - [Redis Documentation](https://redis.io/docs/) - Redis Ltd.
> - [Redis Internals](https://redis.io/docs/reference/internals/) - Redis Source Code Analysis
> - [Redis Data Types](https://redis.io/docs/data-types/) - Official Reference
> - [Go-Redis Client](https://github.com/redis/go-redis) - Official Go Client

---

## 1. Redis Data Structures Internal Architecture

### 1.1 String (SDS - Simple Dynamic String)

**Internal Structure**:

```c
// sds.h - Redis 7.0+ implementation
struct __attribute__ ((__packed__)) sdshdr8 {
    uint8_t len;        // 已使用长度
    uint8_t alloc;      // 分配总长度
    unsigned char flags; // 类型标记
    char buf[];         // 柔性数组
};

struct __attribute__ ((__packed__)) sdshdr16 {
    uint16_t len;
    uint16_t alloc;
    unsigned char flags;
    char buf[];
};

// 64-bit systems use sdshdr64 for large strings
struct __attribute__ ((__packed__)) sdshdr64 {
    uint64_t len;
    uint64_t alloc;
    unsigned char flags;
    char buf[];
};
```

**Design Rationale**:

- **O(1) 长度获取**: `len` 字段直接存储，无需遍历
- **预分配策略**: 减少内存重分配次数
- **二进制安全**: 支持任意字节序列，不仅限于文本
- **兼容 C 字符串**: 末尾隐式 `\0`，可直接使用 C 字符串函数

**内存布局**:

```
┌─────────────────────────────────────────────────────────┐
│                    SDS Memory Layout                     │
├─────────────────────────────────────────────────────────┤
│  Header  │  Content Buffer  │  Free Space  │  Implicit  │
│  (1-8B)  │     (len bytes)  │(alloc-len B) │    \0      │
└─────────────────────────────────────────────────────────┘
         ▲
         └── 用户可见指针 (buf[0])
```

### 1.2 List (Quicklist + Listpack)

**Redis 7.0+ Architecture**:

```c
// quicklist.h
struct quicklist {
    quicklistNode *head;
    quicklistNode *tail;
    unsigned long count;        // 总元素数
    unsigned long len;          // 节点数
    int fill : QL_FILL_BITS;    // 每个节点填充因子
    unsigned int compress : QL_COMP_BITS; // 压缩深度
};

struct quicklistNode {
    quicklistNode *prev;
    quicklistNode *next;
    unsigned char *entry;       // Listpack 数据
    size_t sz;                  // 数据大小
    unsigned int count : 16;    // 元素数量
    unsigned int encoding : 2;  // RAW=1, LZF=2
    unsigned int container : 2; // PLAIN=1, PACKED=2
    unsigned int recompress : 1;
    unsigned int attempted_compress : 1;
};
```

**Quicklist + Listpack 演进**:

- **Redis 3.2**: 引入 quicklist（双向链表 + ziplist）
- **Redis 7.0**: ziplist 升级为 listpack，解决级联更新问题

**Listpack Structure**:

```
┌────────────────────────────────────────────────────────────┐
│                     Listpack Structure                      │
├────────────────────────────────────────────────────────────┤
│  <total-bytes> <num-elements> <element-1> ... <element-n>  │
│      4B             2B                                         │
├────────────────────────────────────────────────────────────┤
│  Element Format:                                            │
│  ┌─────────────┬──────────────┬───────────────┬───────────┐ │
│  │ encoding    │ content-len  │ actual-content│ end-len   │ │
│  │ (1-5B)      │ (optional)   │               │ (1B)      │ │
│  └─────────────┴──────────────┴───────────────┴───────────┘ │
└────────────────────────────────────────────────────────────┘
```

**Encoding Types**:

- `1111xxxx`: 小整数 (0-12)，encoding 本身存储值
- `0xxxxxxx`: 7-bit 无符号整数
- `10xxxxxx`: 6-bit 字符串，长度存储在 encoding
- `110xxxxx yyyyyyyy`: 13-bit 有符号整数
- `1110xxxx`: 12-bit 无符号整数
- `11110001`: 16-bit 有符号整数 (3 bytes total)
- `11110010`: 24-bit 有符号整数 (4 bytes total)
- `11110011`: 32-bit 有符号整数 (5 bytes total)

### 1.3 Hash (Listpack / Hash Table)

**Small Hash (< 512 entries, value < 64 bytes)**:

- 使用 listpack 编码，节省内存
- 字段-值对连续存储

**Large Hash**:

```c
// dict.h - Hash Table Implementation
struct dict {
    dictType *type;
    dictEntry **ht_table[2];    // 两个哈希表，用于渐进式 rehash
    unsigned long ht_used[2];
    long rehashidx;             // -1 表示未在进行 rehash
    int16_t pauserehash;        // 安全迭代器计数
};

struct dictEntry {
    void *key;
    union {
        void *val;
        uint64_t u64;
        int64_t s64;
        double d;
    } v;
    struct dictEntry *next;     // 链地址法处理冲突
};
```

**渐进式 Rehash**:

```
┌─────────────────────────────────────────────────────────────┐
│              Progressive Rehash Process                      │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Step 1: 分配 ht[1]，大小为 2^n >= ht[0].used * 2            │
│                                                              │
│  Step 2: rehashidx = 0，开始渐进迁移                         │
│                                                              │
│  Step 3: 每次增删查改时，迁移 rehashidx 索引的桶             │
│          ┌─────────┐         ┌─────────┐                    │
│          │ Bucket 0│────────►│ Bucket 0│ (新表)              │
│          ├─────────┤         ├─────────┤                    │
│          │ Bucket 1│         │         │                    │
│          ├─────────┤         │         │                    │
│          │    ...  │         │    ...  │                    │
│          ├─────────┤         ├─────────┤                    │
│          │rehashidx│────────►│rehashidx│ 每次处理一个桶      │
│          ├─────────┤         ├─────────┤                    │
│          │    ...  │         │         │                    │
│          └─────────┘         └─────────┘                    │
│                                                              │
│  Step 4: 查询时先查 ht[0]，再查 ht[1]                        │
│                                                              │
│  Step 5: 全部迁移完成后，ht[0] = ht[1]，rehashidx = -1       │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 1.4 Set (Intset / Hash Table)

**Intset (整数集合)**:

```c
// intset.h
struct intset {
    uint32_t encoding;  // INT16_ENC=2, INT32_ENC=4, INT64_ENC=8
    uint32_t length;    // 元素数量
    int8_t contents[];  // 柔性数组，实际类型由 encoding 决定
};
```

**升级策略**:

- 所有元素初始以 int16_t 存储
- 插入 int32_t 时，整体升级为 int32_t
- 升级不可逆，但保证 O(1) 查找（有序数组 + 二分）

**Set Operation Complexity**:

| Operation | Intset | Hash Table |
|-----------|--------|------------|
| SADD | O(n) 升级可能 | O(1) |
| SISMEMBER | O(log n) | O(1) |
| SRANDMEMBER | O(1) | O(1) |
| SPOP | O(1) | O(1) |
| SUNION | O(n) | O(n) |
| SINTER | O(n*m) | O(min(n,m)) |

### 1.5 Sorted Set (Skiplist + Hash Table)

**Dual Structure Design**:

```c
// server.h
struct zset {
    dict *dict;           // member -> score 映射，O(1) 查找
    zskiplist *zsl;       // 按 score 排序，范围查询
};

struct zskiplist {
    zskiplistNode *header, *tail;
    unsigned long length;
    int level;            // 当前最大层数
};

struct zskiplistNode {
    sds ele;              // 成员
    double score;
    struct zskiplistNode *backward;
    struct zskiplistLevel {
        zskiplistNode *forward;
        unsigned int span;    // 到 forward 的跨度
    } level[];            // 柔性数组，层数随机 1-32
};
```

**Skiplist Level Generation**:

```
层数计算 (幂次定律):
Level 1: 概率 1/2    (50% 节点)
Level 2: 概率 1/4    (25% 节点)
Level 3: 概率 1/8    (12.5% 节点)
...
Level 32: 概率 1/2^32

期望值: E[level] = 1/(1-p) = 2 (当 p=0.5)
```

**Visual Skiplist Structure**:

```
┌─────────────────────────────────────────────────────────────────┐
│                     Sorted Set Skiplist                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Level 3:  Head ───────────────────────────────►  NULL          │
│                    │                                             │
│  Level 2:  Head ───┼──────────────►  [30] ──────►  NULL          │
│                    │                 │                           │
│  Level 1:  Head ───┼──►  [10] ─────►│ [30] ────►  [50] ──► NULL │
│                    │      │          │    │        │             │
│  Level 0:  Head ───┴──►  [10] ────► [20]─►[30]──► [40]─►[50]─►  │
│                                                                  │
│  Backward Links:                                                 │
│  [10] ←── [20] ←── [30] ←── [40] ←── [50]                       │
│                                                                  │
│  Span Calculation:                                               │
│  Head.level[2].span = 2 (跳过了 [10], [20])                      │
│  Head.level[1].span = 1                                          │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 1.6 Bitmap / Bitfield

**Bitmap Storage**:

```
┌─────────────────────────────────────────────────────────────┐
│                    Bitmap Memory Layout                      │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Bit 0-7    Bit 8-15   Bit 16-23  Bit 24-31                 │
│  ┌─────┐   ┌─────┐   ┌─────┐   ┌─────┐                     │
│  │Byte0│   │Byte1│   │Byte2│   │Byte3│  ...                 │
│  └─────┘   └─────┘   └─────┘   └─────┘                     │
│                                                              │
│  SETBIT key 7 1  →  Byte0 = 0b0000_0001                     │
│  SETBIT key 9 1  →  Byte1 = 0b0000_0010                     │
│                                                              │
│  底层使用 SDS 存储，支持动态扩展                              │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

**Bitfield 编码**:

- `i8/i16/i32/i64`: 有符号整数
- `u8/u16/u32/u64`: 无符号整数
- 支持溢出控制 (WRAP/SAT/FAIL)

### 1.7 HyperLogLog

**HLL Structure**:

```
┌─────────────────────────────────────────────────────────────┐
│                 HyperLogLog Structure                        │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Sparse Representation (小基数):                            │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  XZERO (6B run) │ XZERO │ VAL (4-bit len + 4-bit val)│   │
│  └─────────────────────────────────────────────────────┘   │
│                                                              │
│  Dense Representation (大基数):                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Magic(4B) │ HLLD(4B) │ Register[0] ... Register[2^14-1]│   │
│  └─────────────────────────────────────────────────────┘   │
│                                                              │
│  Register: 6 bits each, stores max run of leading zeros     │
│  Standard Error: 1.04 / sqrt(2^14) ≈ 0.81%                  │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 1.8 Stream (Redis 5.0+)

**Stream Structure**:

```c
// stream.h
struct stream {
    rax *rax;                 // Radix tree 存储消息
    uint64_t length;          // 条目数
    streamID last_id;         // 最后生成的 ID
    streamID first_id;        // 第一个有效 ID
    streamID max_deleted_entry_id;
    uint64_t entries_added;
    rax *cgroups;             // 消费者组
};

struct streamID {
    uint64_t ms;              // 毫秒时间戳
    uint64_t seq;             // 序列号
};
```

**Radix Tree for Stream**:

```
┌─────────────────────────────────────────────────────────────────┐
│              Radix Tree Stream Storage                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Key Format: 1712345678900-0 (timestamp-seq)                    │
│                                                                  │
│                    [root]                                       │
│                     /  \                                        │
│                   17    18                                      │
│                  /        \                                     │
│              12...        ...                                   │
│             /                                                   │
│        34...                                                    │
│       /    \                                                    │
│    56...   57...                                                │
│    /                                                            │
│  [ID] ──► {field-value pairs}                                   │
│                                                                  │
│  Consumer Group:                                                 │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  Group Name │ Last Delivered ID │ Consumers (PEL)       │    │
│  │  mygroup    │ 1712345678900-5   │ consumer1, consumer2  │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 2. Redis Configuration Best Practices

### 2.1 Memory Management

```conf
# redis.conf - Memory Optimization

# 最大内存限制 (必须设置)
maxmemory 4gb

# 淘汰策略选择
# allkeys-lru: 所有键按 LRU 淘汰 (推荐缓存场景)
# volatile-lru: 仅淘汰有过期时间的键
# allkeys-lfu: 按访问频率淘汰 (适合幂律分布)
# volatile-ttl: 淘汰即将过期的键
# noeviction: 不淘汰，直接返回错误
maxmemory-policy allkeys-lru

# 采样数量，越大越精确但消耗 CPU
maxmemory-samples 5

# 开启内存碎片整理
activedefrag yes

# 碎片整理配置
active-defrag-ignore-bytes 100mb
active-defrag-threshold-lower 10
active-defrag-threshold-upper 100
active-defrag-cycle-min 5
active-defrag-cycle-max 75
```

### 2.2 Persistence Configuration

**RDB Configuration**:

```conf
# 自动保存策略 (频率根据数据重要性调整)
save 900 1      # 900秒内有1次修改
save 300 10     # 300秒内有10次修改
save 60 10000   # 60秒内有10000次修改

# RDB 压缩
rdbcompression yes
rdbchecksum yes
dbfilename dump.rdb

# 后台保存失败时停止写入 (数据安全)
stop-writes-on-bgsave-error yes
```

**AOF Configuration**:

```conf
# 开启 AOF
appendonly yes
appendfilename "appendonly.aof"

# 同步策略
# always: 每个命令都 fsync (最安全，最慢)
# everysec: 每秒 fsync (推荐，平衡安全和性能)
# no: 由操作系统决定 (最快，风险最高)
appendfsync everysec

# AOF 重写配置
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb

# 加载截断 AOF
aof-load-truncated yes

# RDB-AOF 混合持久化 (Redis 7.0+ 推荐)
aof-use-rdb-preamble yes
```

### 2.3 Connection & Network

```conf
# 网络配置
bind 0.0.0.0
port 6379
tcp-backlog 511
timeout 300
tcp-keepalive 300

# 连接限制
maxclients 10000

# 慢查询日志
slowlog-log-slower-than 10000    # 10ms
slowlog-max-len 128

# 大 key 检测
latency-monitor-threshold 100
```

### 2.4 Security Hardening

```conf
# 认证
requirepass your-strong-password

# 命令重命名/禁用 (危险命令)
rename-command FLUSHDB ""
rename-command FLUSHALL ""
rename-command CONFIG "CONFIG_9f3b2a1e"
rename-command DEBUG ""

# 受保护模式
protected-mode yes

# ACL 配置 (Redis 6.0+)
user default on >password ~* &* +@all
user app-read on >app-pass ~app:* &* +@read
user app-write on >write-pass ~app:* &* +@write -@dangerous
```

---

## 3. Performance Tuning Guidelines

### 3.1 Key Design Patterns

**Hot Key Sharding**:

```
问题: 单个热门键 (如计数器) 导致单分片热点

解决方案:
┌─────────────────────────────────────────────────────────────┐
│                    Hot Key Sharding                          │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  原始设计:                                                   │
│  INCR user:12345:view_count  (热点)                         │
│                                                              │
│  分片设计:                                                   │
│  INCR user:12345:view_count:0  ──┐                          │
│  INCR user:12345:view_count:1  ──┼── 100 个分片             │
│  INCR user:12345:view_count:2  ──┘                          │
│  ...                                                         │
│                                                              │
│  获取总值:                                                   │
│  SUM = GET user:12345:view_count:0 + ... + :99              │
│                                                              │
│  或使用 Redis Lua 原子计算:                                  │
│  local sum = 0                                               │
│  for i=0,99 do                                               │
│    sum = sum + redis.call('GET', KEYS[1]..':'..i)            │
│  end                                                         │
│  return sum                                                  │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

**Big Key Prevention**:

```go
// 检测大 key 的 Go 代码
package main

import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
)

func scanBigKeys(ctx context.Context, rdb *redis.Client, threshold int64) {
    var cursor uint64
    var bigKeys []string

    for {
        keys, nextCursor, err := rdb.Scan(ctx, cursor, "*", 100).Result()
        if err != nil {
            panic(err)
        }

        pipe := rdb.Pipeline()
        cmds := make([]*redis.Cmd, len(keys))

        for i, key := range keys {
            cmds[i] = pipe.Do(ctx, "MEMORY", "USAGE", key)
        }

        _, err = pipe.Exec(ctx)
        if err != nil {
            continue
        }

        for i, cmd := range cmds {
            size, err := cmd.Int64()
            if err == nil && size > threshold {
                bigKeys = append(bigKeys, fmt.Sprintf("%s: %d bytes", keys[i], size))
            }
        }

        cursor = nextCursor
        if cursor == 0 {
            break
        }
    }

    for _, k := range bigKeys {
        fmt.Println(k)
    }
}
```

### 3.2 Pipeline & Transaction Optimization

```go
package main

import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
    "sync"
    "time"
)

// PipelineExample 管道批处理示例
func PipelineExample() {
    rdb := redis.NewClient(&redis.Options{
        Addr:         "localhost:6379",
        PoolSize:     100,
        MinIdleConns: 10,
    })
    defer rdb.Close()

    ctx := context.Background()

    // 非管道方式 (慢)
    start := time.Now()
    for i := 0; i < 1000; i++ {
        rdb.Set(ctx, fmt.Sprintf("key:%d", i), i, 0)
    }
    fmt.Printf("Non-pipelined: %v\n", time.Since(start))

    // 管道方式 (快 ~100x)
    start = time.Now()
    pipe := rdb.Pipeline()
    for i := 0; i < 1000; i++ {
        pipe.Set(ctx, fmt.Sprintf("key:%d", i), i, 0)
    }
    _, err := pipe.Exec(ctx)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Pipelined: %v\n", time.Since(start))
}

// TxPipelineExample 事务管道
func TxPipelineExample() {
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    defer rdb.Close()

    ctx := context.Background()

    // 乐观锁事务
    fn := func(tx *redis.Tx) error {
        // 获取当前值
        val, err := tx.Get(ctx, "counter").Int()
        if err != nil && err != redis.Nil {
            return err
        }

        // 业务逻辑
        newVal := val + 1

        // 提交事务
        _, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
            pipe.Set(ctx, "counter", newVal, 0)
            return nil
        })
        return err
    }

    // 重试机制
    maxRetries := 100
    for i := 0; i < maxRetries; i++ {
        err := rdb.Watch(ctx, fn, "counter")
        if err == nil {
            break
        }
        if err == redis.TxFailedErr {
            continue // 冲突，重试
        }
        panic(err)
    }
}

// BatchGet 批量获取优化
func BatchGet(ctx context.Context, rdb *redis.Client, keys []string) []interface{} {
    const batchSize = 100
    var wg sync.WaitGroup
    results := make([]interface{}, len(keys))

    for i := 0; i < len(keys); i += batchSize {
        end := i + batchSize
        if end > len(keys) {
            end = len(keys)
        }

        wg.Add(1)
        go func(start, end int) {
            defer wg.Done()

            pipe := rdb.Pipeline()
            cmds := make([]*redis.StringCmd, end-start)

            for j := start; j < end; j++ {
                cmds[j-start] = pipe.Get(ctx, keys[j])
            }

            _, err := pipe.Exec(ctx)
            if err != nil && err != redis.Nil {
                return
            }

            for j, cmd := range cmds {
                val, err := cmd.Result()
                if err == nil {
                    results[start+j] = val
                }
            }
        }(i, end)
    }

    wg.Wait()
    return results
}
```

### 3.3 Connection Pool Tuning

```go
// ConnectionPoolOptimization 连接池优化配置
func ConnectionPoolOptimization() *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,

        // 连接池大小 = max(10, CPU cores * 4)
        PoolSize:     100,
        MinIdleConns: 20,

        // 连接超时配置
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,

        // 连接生命周期
        PoolTimeout:  4 * time.Second,  // 从池获取连接等待时间
        IdleTimeout:  5 * time.Minute,  // 空闲连接超时
        MaxConnAge:   30 * time.Minute, // 连接最大存活时间
    })
}

// ClusterClientOptimization 集群客户端优化
func ClusterClientOptimization() *redis.ClusterClient {
    return redis.NewClusterClient(&redis.ClusterOptions{
        Addrs: []string{
            "localhost:7000",
            "localhost:7001",
            "localhost:7002",
        },

        // 连接池
        PoolSize:     50,
        MinIdleConns: 10,

        // 路由配置
        ReadOnly:       false, // 允许从副本读
        RouteByLatency: true,  // 按延迟路由
        RouteRandomly:  false,

        // 重试配置
        MaxRetries:      3,
        MinRetryBackoff: 8 * time.Millisecond,
        MaxRetryBackoff: 512 * time.Millisecond,
    })
}
```

---

## 4. Go Client Implementation Patterns

### 4.1 Production-Ready Client

```go
package redisclient

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

// Client Redis 客户端封装
type Client struct {
    rdb    *redis.Client
    config *Config
}

// Config Redis 配置
type Config struct {
    Addr         string
    Password     string
    DB           int
    PoolSize     int
    MinIdleConns int
    MaxRetries   int
}

// NewClient 创建客户端
func NewClient(cfg *Config) *Client {
    rdb := redis.NewClient(&redis.Options{
        Addr:         cfg.Addr,
        Password:     cfg.Password,
        DB:           cfg.DB,
        PoolSize:     cfg.PoolSize,
        MinIdleConns: cfg.MinIdleConns,
        MaxRetries:   cfg.MaxRetries,
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
    })

    return &Client{rdb: rdb, config: cfg}
}

// HealthCheck 健康检查
func (c *Client) HealthCheck(ctx context.Context) error {
    return c.rdb.Ping(ctx).Err()
}

// Close 关闭连接
func (c *Client) Close() error {
    return c.rdb.Close()
}

// Redis Client methods...
```

### 4.2 Distributed Lock Implementation

```go
package redisclient

import (
    "context"
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

// DistributedLock Redis 分布式锁
type DistributedLock struct {
    client *Client
    key    string
    token  string
    ttl    time.Duration
}

// NewDistributedLock 创建分布式锁
func (c *Client) NewDistributedLock(key string, ttl time.Duration) *DistributedLock {
    token := generateToken()
    return &DistributedLock{
        client: c,
        key:    fmt.Sprintf("lock:%s", key),
        token:  token,
        ttl:    ttl,
    }
}

// Acquire 获取锁 (阻塞模式)
func (l *DistributedLock) Acquire(ctx context.Context) error {
    for {
        ok, err := l.TryAcquire(ctx)
        if err != nil {
            return err
        }
        if ok {
            return nil
        }

        // 等待后重试
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(100 * time.Millisecond):
        }
    }
}

// TryAcquire 尝试获取锁 (非阻塞)
func (l *DistributedLock) TryAcquire(ctx context.Context) (bool, error) {
    // SET key token NX EX ttl
    ok, err := l.client.rdb.SetNX(ctx, l.key, l.token, l.ttl).Result()
    if err != nil {
        return false, err
    }
    return ok, nil
}

// Release 释放锁 (Lua 脚本保证原子性)
func (l *DistributedLock) Release(ctx context.Context) error {
    script := `
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        else
            return 0
        end
    `

    result, err := l.client.rdb.Eval(ctx, script, []string{l.key}, l.token).Result()
    if err != nil {
        return err
    }

    if result.(int64) == 0 {
        return fmt.Errorf("lock not held or expired")
    }
    return nil
}

// Extend 延长锁有效期
func (l *DistributedLock) Extend(ctx context.Context, additionalTTL time.Duration) error {
    script := `
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("pexpire", KEYS[1], ARGV[2])
        else
            return 0
        end
    `

    result, err := l.client.rdb.Eval(ctx, script, []string{l.key},
        l.token,
        int64((l.ttl + additionalTTL).Milliseconds())).Result()
    if err != nil {
        return err
    }

    if result.(int64) == 0 {
        return fmt.Errorf("lock not held or expired")
    }
    l.ttl += additionalTTL
    return nil
}

// AutoRefresh 自动续期
func (l *DistributedLock) AutoRefresh(ctx context.Context, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            // 续期半段时间
            if err := l.Extend(ctx, l.ttl/2); err != nil {
                return
            }
        }
    }
}

func generateToken() string {
    b := make([]byte, 16)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)
}
```

### 4.3 Cache-Aside Pattern with Redis

```go
package redisclient

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
)

// Cache Cache-Aside 模式实现
type Cache struct {
    client     *Client
    prefix     string
    defaultTTL time.Duration
}

// NewCache 创建缓存
func (c *Client) NewCache(prefix string, defaultTTL time.Duration) *Cache {
    return &Cache{
        client:     c,
        prefix:     prefix,
        defaultTTL: defaultTTL,
    }
}

// Get 获取缓存
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
    fullKey := c.fullKey(key)

    data, err := c.client.rdb.Get(ctx, fullKey).Result()
    if err != nil {
        if err == redis.Nil {
            return false, nil
        }
        return false, err
    }

    if err := json.Unmarshal([]byte(data), dest); err != nil {
        return false, err
    }
    return true, nil
}

// Set 设置缓存
func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error {
    fullKey := c.fullKey(key)

    data, err := json.Marshal(value)
    if err != nil {
        return err
    }

    duration := c.defaultTTL
    if len(ttl) > 0 {
        duration = ttl[0]
    }

    return c.client.rdb.Set(ctx, fullKey, data, duration).Err()
}

// Delete 删除缓存
func (c *Cache) Delete(ctx context.Context, keys ...string) error {
    if len(keys) == 0 {
        return nil
    }

    fullKeys := make([]string, len(keys))
    for i, k := range keys {
        fullKeys[i] = c.fullKey(k)
    }

    return c.client.rdb.Del(ctx, fullKeys...).Err()
}

// GetOrSet 缓存不存在时计算并设置 (带防击穿)
func (c *Cache) GetOrSet(
    ctx context.Context,
    key string,
    dest interface{},
    fn func() (interface{}, error),
) error {
    // 1. 尝试从缓存获取
    found, err := c.Get(ctx, key, dest)
    if err != nil {
        return err
    }
    if found {
        return nil
    }

    // 2. 获取分布式锁，防止缓存击穿
    lock := c.client.NewDistributedLock(fmt.Sprintf("cache:%s", key), 10*time.Second)
    if err := lock.Acquire(ctx); err != nil {
        // 获取锁失败，尝试再读一次
        found, _ = c.Get(ctx, key, dest)
        if found {
            return nil
        }
        return err
    }
    defer lock.Release(ctx)

    // 3. 双重检查
    found, _ = c.Get(ctx, key, dest)
    if found {
        return nil
    }

    // 4. 执行计算
    val, err := fn()
    if err != nil {
        return err
    }

    // 5. 写入缓存
    if err := c.Set(ctx, key, val); err != nil {
        return err
    }

    // 6. 返回结果
    data, _ := json.Marshal(val)
    return json.Unmarshal(data, dest)
}

// Invalidate 缓存失效模式
func (c *Cache) Invalidate(ctx context.Context, pattern string) error {
    iter := c.client.rdb.Scan(ctx, 0, c.fullKey(pattern), 100).Iterator()

    var keys []string
    for iter.Next(ctx) {
        keys = append(keys, iter.Val())
        if len(keys) >= 100 {
            c.client.rdb.Del(ctx, keys...)
            keys = keys[:0]
        }
    }

    if len(keys) > 0 {
        return c.client.rdb.Del(ctx, keys...).Err()
    }
    return iter.Err()
}

func (c *Cache) fullKey(key string) string {
    return fmt.Sprintf("%s:%s", c.prefix, key)
}
```

### 4.4 Redis Stream Consumer Group

```go
package redisclient

import (
    "context"
    "fmt"
    "log"
    "time"
)

// StreamConsumer Stream 消费者
type StreamConsumer struct {
    client    *Client
    stream    string
    group     string
    consumer  string
    batchSize int64
}

// NewStreamConsumer 创建消费者
func (c *Client) NewStreamConsumer(stream, group, consumer string) *StreamConsumer {
    return &StreamConsumer{
        client:    c,
        stream:    stream,
        group:     group,
        consumer:  consumer,
        batchSize: 10,
    }
}

// CreateGroup 创建消费者组
func (s *StreamConsumer) CreateGroup(ctx context.Context, startID string) error {
    // MKSTREAM 参数: 流不存在时自动创建
    return s.client.rdb.XGroupCreateMkStream(ctx, s.stream, s.group, startID).Err()
}

// Consume 消费消息
func (s *StreamConsumer) Consume(ctx context.Context, handler func(map[string]interface{}) error) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        // 读取消息
        streams, err := s.client.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
            Group:    s.group,
            Consumer: s.consumer,
            Streams:  []string{s.stream, ">"}, // ">" 表示未传递给消费者的消息
            Count:    s.batchSize,
            Block:    5 * time.Second,
        }).Result()

        if err != nil {
            if err == redis.Nil {
                continue // 无消息，继续
            }
            return err
        }

        // 处理消息
        for _, stream := range streams {
            for _, msg := range stream.Messages {
                if err := handler(msg.Values); err != nil {
                    log.Printf("Message handler error: %v", err)
                    // 处理失败，不确认，消息会重新投递
                    continue
                }

                // 确认消息
                if err := s.client.rdb.XAck(ctx, s.stream, s.group, msg.ID).Err(); err != nil {
                    log.Printf("XAck error: %v", err)
                }
            }
        }
    }
}

// ClaimPending 认领挂起消息 (故障转移)
func (s *StreamConsumer) ClaimPending(ctx context.Context, minIdle time.Duration) error {
    // 获取挂起的消息
    pending, err := s.client.rdb.XPendingExt(ctx, &redis.XPendingExtArgs{
        Stream: s.stream,
        Group:  s.group,
        Start:  "-",
        End:    "+",
        Count:  10,
    }).Result()

    if err != nil {
        return err
    }

    var ids []string
    for _, p := range pending {
        if p.Idle >= minIdle {
            ids = append(ids, p.ID)
        }
    }

    if len(ids) == 0 {
        return nil
    }

    // 认领消息
    _, err = s.client.rdb.XClaim(ctx, &redis.XClaimArgs{
        Stream:   s.stream,
        Group:    s.group,
        Consumer: s.consumer,
        MinIdle:  minIdle,
        Messages: ids,
    }).Result()

    return err
}

// Produce 生产消息
func (s *StreamConsumer) Produce(ctx context.Context, values map[string]interface{}) (string, error) {
    return s.client.rdb.XAdd(ctx, &redis.XAddArgs{
        Stream: s.stream,
        Values: values,
    }).Result()
}
```

---

## 5. Visual Representations

### 5.1 Redis Data Structure Hierarchy

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     Redis Data Structure Architecture                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                         Redis Object System                          │    │
│  │                    (redisObject / robj)                              │    │
│  ├─────────────────────────────────────────────────────────────────────┤    │
│  │  type: STRING | LIST | HASH | SET | ZSET | STREAM | ...             │    │
│  │  encoding: RAW | INT | EMBSTR | HT | ZIPLIST | INTSET | ...         │    │
│  │  lru: 24-bit access time                                             │    │
│  │  refcount: reference count                                           │    │
│  │  ptr: ──────┐                                                        │    │
│  └─────────────┼───────────────────────────────────────────────────────┘    │
│                │                                                             │
│                ▼                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    Encoding-Specific Implementations                 │    │
│  ├─────────────────────────────────────────────────────────────────────┤    │
│  │                                                                      │    │
│  │  STRING ──► SDS (Simple Dynamic String)                             │    │
│  │      ├─► EMBSTR (< 44 bytes, embedded object)                       │    │
│  │      ├─► RAW (> 44 bytes, separate allocation)                      │    │
│  │      └─► INT (integer, 8-byte long)                                 │    │
│  │                                                                      │    │
│  │  LIST ────► QUICKLIST ──► Listpack Nodes                            │    │
│  │      ├─► Node 1: [elem1, elem2, ...]                                │    │
│  │      ├─► Node 2: [elemN, ...]                                       │    │
│  │      └─► (Each node is a listpack for memory efficiency)            │    │
│  │                                                                      │    │
│  │  HASH ────► LISTPACK (< 512 entries, < 64B values)                  │    │
│  │      └─► HASHTABLE (dict) - O(1) operations                         │    │
│  │                                                                      │    │
│  │  SET ─────► INTSET (integer-only, sorted array)                     │    │
│  │      └─► HASHTABLE (general purpose)                                │    │
│  │                                                                      │    │
│  │  ZSET ────► SKIPLIST + HASHTABLE (dual structure)                   │    │
│  │      ├─► Dict: member -> score (O(1) lookup)                        │    │
│  │      └─► Skiplist: sorted by score (range queries)                  │    │
│  │                                                                      │    │
│  │  STREAM ──► RADIX TREE ──► Listpack Chunks                          │    │
│  │      └─► Consumer Groups (PEL tracking)                             │    │
│  │                                                                      │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.2 Redis Cluster Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Redis Cluster Architecture                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│                         ┌──────────────┐                                    │
│                         │   Client     │                                    │
│                         └──────┬───────┘                                    │
│                                │                                            │
│              ┌─────────────────┼─────────────────┐                          │
│              │                 │                 │                          │
│              ▼                 ▼                 ▼                          │
│       ┌────────────┐   ┌────────────┐   ┌────────────┐                     │
│       │  Master A  │◄──┤  Master B  │◄──┤  Master C  │                     │
│       │ (0-5460)   │   │ (5461-10922)│   │ (10923-16383)│                   │
│       └─────┬──────┘   └─────┬──────┘   └─────┬──────┘                     │
│             │                │                │                             │
│             │ Replication    │ Replication    │ Replication                 │
│             │                │                │                             │
│             ▼                ▼                ▼                             │
│       ┌────────────┐   ┌────────────┐   ┌────────────┐                     │
│       │  Replica A │   │  Replica B │   │  Replica C │                     │
│       │ (Read Only)│   │ (Read Only)│   │ (Read Only)│                     │
│       └────────────┘   └────────────┘   └────────────┘                     │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    Slot Distribution (16384 slots)                   │    │
│  ├─────────────────────────────────────────────────────────────────────┤    │
│  │                                                                      │    │
│  │  Slot 0 ────────────────────────────► Slot 5460     [Master A]      │    │
│  │  Slot 5461 ─────────────────────────► Slot 10922    [Master B]      │    │
│  │  Slot 10923 ────────────────────────► Slot 16383    [Master C]      │    │
│  │                                                                      │    │
│  │  Key → Slot: HASH_SLOT = CRC16(key) mod 16384                       │    │
│  │                                                                      │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                      Gossip Protocol                                 │    │
│  ├─────────────────────────────────────────────────────────────────────┤    │
│  │                                                                      │    │
│  │   ┌─────────┐      Ping/Pong      ┌─────────┐                       │    │
│  │   │ Node A  │◄───────────────────►│ Node B  │                       │    │
│  │   └────┬────┘                     └────┬────┘                       │    │
│  │        │         ╲            ╱        │                            │    │
│  │        │          ╲          ╱         │                            │    │
│  │        │           ╲        ╱          │                            │    │
│  │        ▼            ▼      ▼           ▼                            │    │
│  │   ┌─────────┐    ┌─────────┐    ┌─────────┐                        │    │
│  │   │ Node C  │◄──►│ Node D  │◄──►│ Node E  │                        │    │
│  │   └─────────┘    └─────────┘    └─────────┘                        │    │
│  │                                                                      │    │
│  │  Gossip Content:                                                     │    │
│  │  - Node address and flags                                            │    │
│  │  - Slot ownership                                                    │    │
│  │  - Last Ping timestamp                                               │    │
│  │  - Configuration epoch                                               │    │
│  │                                                                      │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.3 Redis Persistence & Replication Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                Redis Persistence & Replication Architecture                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                       Memory (Dataset)                                 │  │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐          │  │
│  │  │ String  │ │  List   │ │  Hash   │ │  Set    │ │ Stream  │  ...     │  │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘          │  │
│  └──────────┬────────────────────────────────────────────────────────────┘  │
│             │                                                                │
│  ┌──────────┼────────────────────────────────────────────────────────────┐  │
│  │          │                    Persistence Layer                        │  │
│  │          ▼                                                             │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐   │  │
│  │  │                      AOF Rewrite Process                         │   │  │
│  │  │  ┌───────────┐    fork()    ┌───────────┐    ┌───────────┐     │   │  │
│  │  │  │  Parent   │──────────────►│  Child    │───►│ New AOF   │     │   │  │
│  │  │  │ Process   │               │ Process   │    │ (Compact) │     │   │  │
│  │  │  └─────┬─────┘               └───────────┘    └───────────┘     │   │  │
│  │  │        │                                                     │   │  │
│  │  │        └──► Write to AOF Rewrite Buffer (Incremental)         │   │  │
│  │  └─────────────────────────────────────────────────────────────────┘   │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐   │  │
│  │  │                      RDB Snapshot Process                        │   │  │
│  │  │  ┌───────────┐    fork()    ┌───────────┐    ┌───────────┐     │   │  │
│  │  │  │  Parent   │──────────────►│  Child    │───►│ RDB File  │     │   │  │
│  │  │  │ Process   │               │ Process   │    │ (Binary)  │     │   │  │
│  │  │  └─────┬─────┘               └───────────┘    └───────────┘     │   │  │
│  │  │        │                                                     │   │  │
│  │  │        └──► Copy-on-Write (COW) Pages Modified During Save    │   │  │
│  │  └─────────────────────────────────────────────────────────────────┘   │  │
│  └────────────────────────────────────────────────────────────────────────┘  │
│                                                                               │
│  ┌────────────────────────────────────────────────────────────────────────┐  │
│  │                      Replication Flow                                   │  │
│  ├────────────────────────────────────────────────────────────────────────┤  │
│  │                                                                         │  │
│  │   Master                              Replica                           │  │
│  │  ┌──────────┐                        ┌──────────┐                       │  │
│  │  │ 1. BGSAVE│────RDB File───────────►│  Load    │                       │  │
│  │  │          │                        │  RDB     │                       │  │
│  │  │ 2. Buffer│                        └────┬─────┘                       │  │
│  │  │ Commands │◄─────────────────────────────┘                            │  │
│  │  │ (Offset) │  Replication Offset Sync                                  │  │
│  │  └────┬─────┘                                                           │  │
│  │       │                                                                 │  │
│  │       │ Stream Commands to Replicas                                     │  │
│  │       ├──► SET foo bar                                                  │  │
│  │       ├──► DEL mykey                                                    │  │
│  │       ├──► HSET user:1 name bob                                         │  │
│  │       └──► XADD stream * field value                                    │  │
│  │                                                                         │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐    │  │
│  │  │              Partial Resynchronization (PSYNC)                   │    │  │
│  │  │  Replica: PSYNC <master_replid> <replica_offset>                 │    │  │
│  │  │  Master:  +CONTINUE (if backlog contains offset)                 │    │  │
│  │  │           +FULLRESYNC (otherwise, send new RDB)                  │    │  │
│  │  └─────────────────────────────────────────────────────────────────┘    │  │
│  │                                                                         │  │
│  └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                               │
└───────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. Performance Checklist

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Redis Performance Checklist                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  内存优化:                                                                   │
│  □ 使用合适的数据结构 (ziplist/listpack vs hashtable)                       │
│  □ 设置合理的 maxmemory 和淘汰策略                                          │
│  □ 启用内存碎片整理 (activedefrag)                                          │
│  □ 监控内存使用: INFO memory                                                │
│  □ 避免大 key (MEMORY USAGE 检测)                                           │
│                                                                              │
│  命令优化:                                                                   │
│  □ 使用 O(1) 命令，避免 O(N) 操作                                           │
│  □ 批量操作使用 Pipeline                                                    │
│  □ 范围查询控制返回数量                                                     │
│  □ 使用 UNLINK 替代 DEL (异步删除)                                          │
│                                                                              │
│  连接管理:                                                                   │
│  □ 使用连接池 (合理设置 PoolSize)                                           │
│  □ 控制连接超时时间                                                         │
│  □ 监控连接数: CLIENT LIST | INFO clients                                   │
│                                                                              │
│  持久化优化:                                                                 │
│  □ 低峰期执行 BGSAVE                                                        │
│  □ AOF 使用 everysec 模式                                                   │
│  □ 开启 AOF 重写压缩                                                        │
│  □ 监控 fork 耗时: INFO stats latest_fork_usec                              │
│                                                                              │
│  集群优化:                                                                   │
│  □ 避免跨 slot 事务 (MULTI)                                                 │
│  □ 使用 hashtag 确保相关 key 在同一 slot                                    │
│  □ 监控槽位均衡                                                             │
│                                                                              │
│  Go 客户端最佳实践:                                                          │
│  □ 使用 go-redis/v9 官方客户端                                              │
│  □ 配置合理的连接池参数                                                     │
│  □ 大批量操作分批处理                                                       │
│  □ 实现优雅关闭                                                             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 7. References

1. **Redis Documentation** (2024). Data Types. Redis Ltd.
2. **Redis Documentation** (2024). Memory Optimization. Redis Ltd.
3. **Redis Source Code** (v7.2). github.com/redis/redis
4. **Go-Redis Documentation** (v9). github.com/redis/go-redis
5. **Zhang, T.** (2023). Redis 5设计与源码分析. 机械工业出版社.
6. **Josiah Carlson** (2023). Redis in Action, 2nd Edition. Manning Publications.

---

*Document Version: 1.0 | Last Updated: 2024*

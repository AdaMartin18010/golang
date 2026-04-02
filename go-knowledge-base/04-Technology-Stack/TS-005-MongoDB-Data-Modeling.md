# TS-005: MongoDB Data Modeling - Schema Design & Go Implementation

> **维度**: Technology Stack
> **级别**: S (18+ KB)
> **标签**: #mongodb #nosql #data-modeling #document #go
> **权威来源**:
> - [MongoDB Documentation](https://docs.mongodb.com/) - MongoDB Inc.
> - [MongoDB: The Definitive Guide](https://www.oreilly.com/library/view/mongodb-the-definitive/) - O'Reilly Media
> - [Designing Data-Intensive Applications](https://dataintensive.net/) - Martin Kleppmann

---

## 1. MongoDB Storage Architecture

### 1.1 WiredTiger Storage Engine

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    WiredTiger Storage Engine Architecture                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Memory Layer (Cache)                                │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  WiredTiger Cache (50% RAM - 1GB by default)                     │  │  │
│  │  │                                                                  │  │  │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │  │  │
│  │  │  │ Collection  │  │ Collection  │  │ Index B-tree│             │  │  │
│  │  │  │    A        │  │    B        │  │    Data     │             │  │  │
│  │  │  │  (Pages)    │  │  (Pages)    │  │  (Pages)    │             │  │  │
│  │  │  └─────────────┘  └─────────────┘  └─────────────┘             │  │  │
│  │  │                                                                  │  │  │
│  │  │  Page Structure:                                                 │  │  │
│  │  │  ┌────────────────────────────────────────────────────────────┐ │  │  │
│  │  │  │ Header │ Key/Value Pairs │ Trailer (checksum)              │ │  │  │
│  │  │  └────────────────────────────────────────────────────────────┘ │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Disk Layer                                          │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  /data/db/                                                             │  │
│  │  ├── WiredTiger                        (存储引擎元数据)                 │  │
│  │  ├── WiredTiger.lock                   (锁文件)                        │  │
│  │  ├── WiredTiger.turtle                 (检查点元数据)                   │  │
│  │  ├── WiredTiger.wt                     (根表)                          │  │
│  │  ├── collection-*.wt                   (集合数据文件)                   │  │
│  │  ├── index-*.wt                        (索引数据文件)                   │  │
│  │  ├── journal/                          (WAL 日志)                       │  │
│  │  │   ├── WiredTigerLog.0000000001                                    │  │
│  │  │   ├── WiredTigerLog.0000000002                                    │  │
│  │  │   └── ...                                                         │  │
│  │  └── diagnostic.data/                  (诊断数据)                       │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    B-tree Structure                                    │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │                    [Root Page]                                         │  │
│  │                   /    |    \                                          │  │
│  │                  /     |     \                                         │  │
│  │            [Internal] [Internal] [Internal]                           │  │
│  │             /   |   \   ...                                            │  │
│  │            /    |    \                                                 │  │
│  │       [Leaf] [Leaf] [Leaf] [Leaf] [Leaf]                              │  │
│  │                                                                        │  │
│  │  Leaf Page Content:                                                    │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Key: _id    │  Value: BSON Document  │  Addresses              │  │  │
│  │  │  ────────────┼─────────────────────────┼────────────────         │  │  │
│  │  │  ObjectId(1) │  {name:"A", ...}        │  Disk addr              │  │  │
│  │  │  ObjectId(2) │  {name:"B", ...}        │  Disk addr              │  │  │
│  │  │  ...         │  ...                    │  ...                    │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  Key Features:                                                         │  │
│  │  - Prefix compression on keys                                          │  │
│  │  - Dictionary compression on values                                    │  │
│  │  - Block compression (snappy/zstd/zlib)                                │  │
│  │  - Multi-version concurrency control (MVCC)                            │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Document Storage Format (BSON)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    BSON Document Structure                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  BSON = Binary JSON (BSON Spec v1.1)                                         │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    BSON Document Layout                                │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  int32 document_size  │  总字节数 (包含自身)                      │  │  │
│  │  │  ─────────────────────┼─────────────────────────────────────     │  │  │
│  │  │  byte element_1       │  元素1: type(1B) + cstring(name) + value  │  │  │
│  │  │  byte element_2       │  元素2 ...                                │  │  │
│  │  │  ...                  │                                           │  │  │
│  │  │  byte element_n       │  元素n                                     │  │  │
│  │  │  ─────────────────────┼─────────────────────────────────────     │  │  │
│  │  │  0x00                 │  文档结束标记                              │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  Element Types:                                                        │  │
│  │  ┌──────┬─────────────────┬────────────────────────────────────────┐  │  │
│  │  │ Type │ Name            │ Size / Description                     │  │  │
│  │  ├──────┼─────────────────┼────────────────────────────────────────┤  │  │
│  │  │ 0x01 │ Double          │ 8 bytes (IEEE 754)                     │  │  │
│  │  │ 0x02 │ String          │ int32 length + cstring + 0x00          │  │  │
│  │  │ 0x03 │ Document        │ Embedded document                      │  │  │
│  │  │ 0x04 │ Array           │ Embedded document (keys are indexes)   │  │  │
│  │  │ 0x05 │ Binary          │ int32 len + subtype(1) + bytes         │  │  │
│  │  │ 0x07 │ ObjectId        │ 12 bytes                               │  │  │
│  │  │ 0x08 │ Boolean         │ 1 byte (0x00/0x01)                     │  │  │
│  │  │ 0x09 │ UTC datetime    │ 8 bytes (int64 millis since epoch)     │  │  │
│  │  │ 0x0A │ Null            │ 0 bytes                                │  │  │
│  │  │ 0x10 │ Int32           │ 4 bytes                                │  │  │
│  │  │ 0x12 │ Int64           │ 8 bytes                                │  │  │
│  │  │ 0x13 │ Decimal128      │ 16 bytes                               │  │  │
│  │  └──────┴─────────────────┴────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Example Document Encoding                           │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  JSON:                                                                 │  │
│  │  {                                                                     │  │
│  │    "_id": ObjectId("507f1f77bcf86cd799439011"),                       │  │
│  │    "name": "John",                                                     │  │
│  │    "age": 30,                                                          │  │
│  │    "active": true                                                      │  │
│  │  }                                                                     │  │
│  │                                                                        │  │
│  │  BSON (hex):                                                           │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │ 4B 00 00 00  │  Document size: 75 bytes                         │  │  │
│  │  │                                                                      │  │
│  │  │ 07 5F 69 64 00  │  Type 0x07 (ObjectId), "_id\0"                  │  │  │
│  │  │ 50 7F 1F 77 BC F8 6C D7 99 43 90 11  │  12-byte ObjectId        │  │  │
│  │  │                                                                      │  │
│  │  │ 02 6E 61 6D 65 00  │  Type 0x02 (String), "name\0"                │  │  │
│  │  │ 05 00 00 00 4A 6F 68 6E 00  │  "John\0" (len=5)                 │  │  │
│  │  │                                                                      │  │
│  │  │ 10 61 67 65 00  │  Type 0x10 (Int32), "age\0"                    │  │  │
│  │  │ 1E 00 00 00  │  30 (little-endian)                               │  │  │
│  │  │                                                                      │  │
│  │  │ 08 61 63 74 69 76 65 00  │  Type 0x08 (Bool), "active\0"           │  │  │
│  │  │ 01  │  true (0x01)                                              │  │  │
│  │  │                                                                      │  │
│  │  │ 00  │  Document terminator                                       │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.3 Replica Set Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    MongoDB Replica Set Architecture                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Three-Node Replica Set                              │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │      ┌─────────────────────────────────────────────────────────┐      │  │
│  │      │                      Heartbeats                          │      │  │
│  │      │  ┌──────────┐◄──────────────►┌──────────┐               │      │  │
│  │      │  │ Primary  │◄──────────────►│ Secondary│               │      │  │
│  │      │  │  (P)     │                │   (S1)   │               │      │  │
│  │      │  └────┬─────┘◄──────────────►└────┬─────┘               │      │  │
│  │      │       │                           │                      │      │  │
│  │      │       │     ┌──────────┐          │                      │      │  │
│  │      │       └────►│ Secondary│◄─────────┘                      │      │  │
│  │      │             │   (S2)   │                                 │      │  │
│  │      │             └──────────┘                                 │      │  │
│  │      └─────────────────────────────────────────────────────────┘      │  │
│  │                                                                        │  │
│  │  Replication Flow:                                                     │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                                                                  │  │  │
│  │  │  1. Client writes to Primary                                     │  │  │
│  │  │     db.users.insertOne({name: "John"})                           │  │  │
│  │  │                                                                  │  │  │
│  │  │  2. Primary writes to local oplog (capped collection)            │  │  │
│  │  │     local.oplog.rs: {ts: Timestamp, op: "i", ns: "db.users", ...}│  │  │
│  │  │                                                                  │  │  │
│  │  │  3. Secondaries tail oplog                                       │  │  │
│  │  │     - Read oplog entries                                         │  │  │
│  │  │     - Apply operations                                           │  │  │
│  │  │     - Update local.lastAppliedOpTime                             │  │  │
│  │  │                                                                  │  │  │
│  │  │  4. Write Concern Acknowledgment                                 │  │  │
│  │  │     w: "majority" ──► wait for >50% nodes acknowledge            │  │  │
│  │  │     w: 2 ──► wait for 2 nodes                                    │  │  │
│  │  │     w: 1 ──► primary only                                        │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  Failover Process:                                                     │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  1. Primary becomes unreachable                                  │  │  │
│  │  │                                                                  │  │  │
│  │  │  2. Secondaries initiate election (after electionTimeoutMillis)  │  │  │
│  │  │                                                                  │  │  │
│  │  │  3. Node with highest priority + up-to-date data becomes Primary │  │  │
│  │  │                                                                  │  │  │
│  │  │  4. Client connections redirected (via mongos or driver)         │  │  │
│  │  │                                                                  │  │  │
│  │  │  5. Old primary rejoins as secondary (rolls back if needed)      │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Read Preferences & Write Concerns                   │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Read Preferences:                                                     │  │
│  │  ┌────────────────────┬────────────────────────────────────────────┐  │  │
│  │  │ Preference         │ Description                                │  │  │
│  │  ├────────────────────┼────────────────────────────────────────────┤  │  │
│  │  │ primary            │ Default, read from primary only            │  │  │
│  │  │ primaryPreferred   │ Prefer primary, fallback to secondary      │  │  │
│  │  │ secondary          │ Read from secondary only                   │  │  │
│  │  │ secondaryPreferred │ Prefer secondary, fallback to primary      │  │  │
│  │  │ nearest            │ Read from nearest node (by latency)        │  │  │
│  │  └────────────────────┴────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  Write Concerns:                                                       │  │
│  │  ┌────────────────────┬────────────────────────────────────────────┐  │  │
│  │  │ Concern            │ Durability Guarantee                       │  │  │
│  │  ├────────────────────┼────────────────────────────────────────────┤  │  │
│  │  │ {w: 0}             │ Fire-and-forget, no acknowledgment         │  │  │
│  │  │ {w: 1}             │ Acknowledged by primary                    │  │  │
│  │  │ {w: "majority"}    │ Committed to majority nodes                │  │  │
│  │  │ {w: "majority",    │ Committed to majority + journaled          │  │  │
│  │  │  j: true}          │                                            │  │  │
│  │  │ {w: 3, wtimeout:   │ Timeout if not acknowledged in time        │  │  │
│  │  │  5000}             │                                            │  │  │
│  │  └────────────────────┴────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Data Modeling Patterns

### 2.1 Schema Design Decision Tree

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    MongoDB Schema Design Decision Tree                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Start Schema Design                                                         │
│       │                                                                      │
│       ▼                                                                      │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  What is the cardinality of the relationship?                          │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│       │                                                                      │
│       ├──► One-to-Few (1:1 - 1:10)                                         │
│       │       │                                                              │
│       │       ├──► Access pattern: Often read together?                      │
│       │       │       ├──► YES ──► Embed (Single Document)                   │
│       │       │       │       ┌──────────────────────────────────────┐       │
│       │       │       │       │  {                                   │       │
│       │       │       │       │    user: "john",                     │       │
│       │       │       │       │    address: {                        │       │
│       │       │       │       │      street: "123 Main",             │       │
│       │       │       │       │      city: "NYC"                     │       │
│       │       │       │       │    }                                 │       │
│       │       │       │       │  }                                   │       │
│       │       │       │       └──────────────────────────────────────┘       │
│       │       │       │                                                      │
│       │       │       └──► NO ──► Reference (Separate Collection)            │
│       │       │               ┌──────────────────────────────────────┐       │
│       │       │               │  // users collection               │       │
│       │       │               │  { _id: 1, name: "john" }          │       │
│       │       │               │                                    │       │
│       │       │               │  // addresses collection           │       │
│       │       │               │  { user_id: 1, street: "123 Main" }│       │
│       │       │               └──────────────────────────────────────┘       │
│       │       │                                                              │
│       │       └──► Data volatility? Frequently updated?                      │
│       │               ├──► YES ──► Consider reference to reduce write lock   │
│       │               └──► NO  ──► Embedding is fine                         │
│       │                                                                      │
│       ├──► One-to-Many (1:10 - 1:1000)                                     │
│       │       │                                                              │
│       │       ├──► Subset pattern: Frequently accessed subset?               │
│       │       │       ├──► YES ──► Partial embedding + reference             │
│       │       │       │       ┌──────────────────────────────────────┐       │
│       │       │       │       │  {                                   │       │
│       │       │       │       │    product: "Widget",                │       │
│       │       │       │       │    recent_reviews: [                 │       │
│       │       │       │       │      {user: "A", rating: 5},         │       │
│       │       │       │       │      {user: "B", rating: 4}          │       │
│       │       │       │       │    ],                                │       │
│       │       │       │       │    review_count: 1000                │       │
│       │       │       │       │  }                                   │       │
│       │       │       │       │  // Full reviews in separate collection      │
│       │       │       │       └──────────────────────────────────────┘       │
│       │       │       │                                                      │
│       │       │       └──► NO ──► Use reference with proper indexing         │
│       │       │                                                              │
│       │       └──► Document size approaching 16MB?                           │
│       │               ├──► YES ──► MUST use reference                        │
│       │               └──► NO  ──► Can consider embedding                    │
│       │                                                                      │
│       └──► One-to-Many/Many (1:1000+ or N:M)                               │
│               │                                                              │
│               ├──► Always use Reference pattern                              │
│               │   ┌──────────────────────────────────────────────────────┐   │
│               │   │  // authors collection                               │   │
│               │   │  { _id: 1, name: "Orwell" }                          │   │
│               │   │                                                      │   │
│               │   │  // books collection                                 │   │
│               │   │  { _id: 101, title: "1984", author_ids: [1, 2] }    │   │
│               │   │                                                      │   │
│               │   │  // Many-to-Many requires array of references        │   │
│               │   └──────────────────────────────────────────────────────┘   │
│               │                                                              │
│               └──► Consider intermediate collection for relationship data    │
│                   ┌──────────────────────────────────────────────────────┐   │
│                   │  // book_authors junction                              │   │
│                   │  { book_id: 101, author_id: 1, role: "primary" }       │   │
│                   └──────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Design Patterns

```go
package mongodb

import (
    "context"
    "fmt"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// Client MongoDB 客户端封装
type Client struct {
    client   *mongo.Client
    database *mongo.Database
}

// NewClient 创建客户端
func NewClient(uri, dbName string) (*Client, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOptions := options.Client().ApplyURI(uri)
    clientOptions.SetMaxPoolSize(100)
    clientOptions.SetMinPoolSize(10)
    clientOptions.SetMaxConnIdleTime(30 * time.Second)

    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return nil, err
    }

    // 验证连接
    if err := client.Ping(ctx, nil); err != nil {
        return nil, err
    }

    return &Client{
        client:   client,
        database: client.Database(dbName),
    }, nil
}

// Close 关闭连接
func (c *Client) Close() error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    return c.client.Disconnect(ctx)
}

// ==================== Pattern 1: Polymorphic Pattern ====================

// BaseEvent 基础事件结构
type BaseEvent struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    Type      string             `bson:"type"`
    Timestamp time.Time          `bson:"timestamp"`
    UserID    string             `bson:"user_id"`
}

// PageViewEvent 页面浏览事件
type PageViewEvent struct {
    BaseEvent `bson:",inline"`
    URL       string `bson:"url"`
    Referrer  string `bson:"referrer"`
}

// PurchaseEvent 购买事件
type PurchaseEvent struct {
    BaseEvent `bson:",inline"`
    ProductID string  `bson:"product_id"`
    Amount    float64 `bson:"amount"`
    Currency  string  `bson:"currency"`
}

// InsertEvent 插入事件 (使用多态模式)
func (c *Client) InsertEvent(ctx context.Context, event interface{}) error {
    collection := c.database.Collection("events")
    _, err := collection.InsertOne(ctx, event)
    return err
}

// QueryEventsByType 按类型查询事件
func (c *Client) QueryEventsByType(ctx context.Context, eventType string, from, to time.Time) ([]bson.M, error) {
    collection := c.database.Collection("events")

    filter := bson.M{
        "type":      eventType,
        "timestamp": bson.M{"$gte": from, "$lte": to},
    }

    cursor, err := collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var results []bson.M
    if err := cursor.All(ctx, &results); err != nil {
        return nil, err
    }

    return results, nil
}

// ==================== Pattern 2: Bucket Pattern ====================

// SensorReading 传感器读数
type SensorReading struct {
    Timestamp time.Time `bson:"timestamp"`
    Value     float64   `bson:"value"`
    Unit      string    `bson:"unit"`
}

// SensorBucket 传感器数据桶
type SensorBucket struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    SensorID  string             `bson:"sensor_id"`
    Date      time.Time          `bson:"date"`       // 桶的日期
    Readings  []SensorReading    `bson:"readings"`   // 该天的所有读数
    Count     int                `bson:"count"`
    AvgValue  float64            `bson:"avg_value"`
    MinValue  float64            `bson:"min_value"`
    MaxValue  float64            `bson:"max_value"`
}

// AddSensorReading 添加传感器读数 (使用桶模式)
func (c *Client) AddSensorReading(ctx context.Context, sensorID string, reading SensorReading) error {
    collection := c.database.Collection("sensor_data")

    // 按天分桶
    date := reading.Timestamp.Truncate(24 * time.Hour)
    bucketID := fmt.Sprintf("%s_%s", sensorID, date.Format("2006-01-02"))

    filter := bson.M{
        "sensor_id": sensorID,
        "date":      date,
    }

    update := bson.M{
        "$push": bson.M{"readings": reading},
        "$inc":  bson.M{"count": 1},
        "$min":  bson.M{"min_value": reading.Value},
        "$max":  bson.M{"max_value": reading.Value},
    }

    // 计算新平均值的复杂更新
    // 实际生产中可能需要聚合管道

    opts := options.Update().SetUpsert(true)
    _, err := collection.UpdateOne(ctx, filter, update, opts)
    return err
}

// GetSensorStats 获取传感器统计
func (c *Client) GetSensorStats(ctx context.Context, sensorID string, from, to time.Time) ([]SensorBucket, error) {
    collection := c.database.Collection("sensor_data")

    filter := bson.M{
        "sensor_id": sensorID,
        "date":      bson.M{"$gte": from, "$lte": to},
    }

    cursor, err := collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var buckets []SensorBucket
    if err := cursor.All(ctx, &buckets); err != nil {
        return nil, err
    }

    return buckets, nil
}

// ==================== Pattern 3: Outlier Pattern ====================

// ProductReview 产品评论
type ProductReview struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    ProductID   string             `bson:"product_id"`
    UserID      string             `bson:"user_id"`
    Rating      int                `bson:"rating"`
    Comment     string             `bson:"comment,omitempty"`
    IsOutlier   bool               `bson:"is_outlier,omitempty"`
}

// ProductWithReviews 带评论的产品 (使用异常值分离模式)
type ProductWithReviews struct {
    ID             primitive.ObjectID `bson:"_id,omitempty"`
    Name           string             `bson:"name"`
    RecentReviews  []ProductReview    `bson:"recent_reviews"`   // 最常访问的最近评论
    ReviewCount    int                `bson:"review_count"`
    AvgRating      float64            `bson:"avg_rating"`
}

// AddReview 添加评论
func (c *Client) AddReview(ctx context.Context, review ProductReview) error {
    // 保存到评论集合
    reviewsColl := c.database.Collection("reviews")
    _, err := reviewsColl.InsertOne(ctx, review)
    if err != nil {
        return err
    }

    // 更新产品文档中的最近评论
    productsColl := c.database.Collection("products")

    // 获取当前评论列表
    var product ProductWithReviews
    err = productsColl.FindOne(ctx, bson.M{"_id": review.ProductID}).Decode(&product)
    if err != nil {
        return err
    }

    // 保持最近10条评论
    recentReviews := append([]ProductReview{review}, product.RecentReviews...)
    if len(recentReviews) > 10 {
        recentReviews = recentReviews[:10]
    }

    // 更新产品
    update := bson.M{
        "$set": bson.M{"recent_reviews": recentReviews},
        "$inc": bson.M{"review_count": 1},
    }

    _, err = productsColl.UpdateOne(ctx, bson.M{"_id": review.ProductID}, update)
    return err
}

// GetFullReviews 获取所有评论 (分页)
func (c *Client) GetFullReviews(ctx context.Context, productID string, page, limit int64) ([]ProductReview, error) {
    collection := c.database.Collection("reviews")

    opts := options.Find().
        SetSkip(page * limit).
        SetLimit(limit).
        SetSort(bson.M{"_id": -1})

    cursor, err := collection.Find(ctx, bson.M{"product_id": productID}, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var reviews []ProductReview
    if err := cursor.All(ctx, &reviews); err != nil {
        return nil, err
    }

    return reviews, nil
}

// ==================== Pattern 4: Computed Pattern ====================

// Order 订单 (带计算字段)
type Order struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    UserID      string             `bson:"user_id"`
    Items       []OrderItem        `bson:"items"`
    SubTotal    float64            `bson:"sub_total"`     // 计算字段
    Tax         float64            `bson:"tax"`           // 计算字段
    Total       float64            `bson:"total"`         // 计算字段
    Status      string             `bson:"status"`
    CreatedAt   time.Time          `bson:"created_at"`
}

// OrderItem 订单项
type OrderItem struct {
    ProductID string  `bson:"product_id"`
    Quantity  int     `bson:"quantity"`
    UnitPrice float64 `bson:"unit_price"`
    Total     float64 `bson:"total"`  // quantity * unit_price
}

// CreateOrder 创建订单 (预计算)
func (c *Client) CreateOrder(ctx context.Context, userID string, items []OrderItem) (*Order, error) {
    // 计算所有派生字段
    var subTotal float64
    for i := range items {
        items[i].Total = float64(items[i].Quantity) * items[i].UnitPrice
        subTotal += items[i].Total
    }

    taxRate := 0.08
    tax := subTotal * taxRate
    total := subTotal + tax

    order := &Order{
        ID:        primitive.NewObjectID(),
        UserID:    userID,
        Items:     items,
        SubTotal:  subTotal,
        Tax:       tax,
        Total:     total,
        Status:    "pending",
        CreatedAt: time.Now(),
    }

    collection := c.database.Collection("orders")
    _, err := collection.InsertOne(ctx, order)
    if err != nil {
        return nil, err
    }

    return order, nil
}

// ==================== Pattern 5: Schema Versioning ====================

// UserV1 用户版本1
type UserV1 struct {
    ID       primitive.ObjectID `bson:"_id,omitempty"`
    Name     string             `bson:"name"`
    Email    string             `bson:"email"`
    SchemaV  int                `bson:"schema_version"`  // 模式版本
}

// UserV2 用户版本2 (新增字段)
type UserV2 struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    Name        string             `bson:"name"`
    Email       string             `bson:"email"`
    Phone       string             `bson:"phone"`           // 新增
    Preferences UserPreferences    `bson:"preferences"`     // 新增
    SchemaV     int                `bson:"schema_version"`
}

type UserPreferences struct {
    Newsletter bool   `bson:"newsletter"`
    Theme      string `bson:"theme"`
}

// UpsertUser 保存用户 (向后兼容)
func (c *Client) UpsertUser(ctx context.Context, user UserV2) error {
    collection := c.database.Collection("users")

    user.SchemaV = 2  // 设置模式版本

    filter := bson.M{"_id": user.ID}
    update := bson.M{"$set": user}
    opts := options.Update().SetUpsert(true)

    _, err := collection.UpdateOne(ctx, filter, update, opts)
    return err
}

// GetUser 获取用户 (处理多版本)
func (c *Client) GetUser(ctx context.Context, userID string) (*UserV2, error) {
    collection := c.database.Collection("users")

    objID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return nil, err
    }

    var raw bson.M
    err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&raw)
    if err != nil {
        return nil, err
    }

    // 根据 schema_version 进行迁移
    schemaVersion, _ := raw["schema_version"].(int32)

    var user UserV2

    switch schemaVersion {
    case 1:
        // 从 V1 迁移到 V2
        user.ID = raw["_id"].(primitive.ObjectID)
        user.Name = raw["name"].(string)
        user.Email = raw["email"].(string)
        user.SchemaV = 2
        // Phone 和 Preferences 使用默认值
    case 2, 0:  // 0 表示没有版本字段 (新文档)
        // 直接解码
        if err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user); err != nil {
            return nil, err
        }
    }

    return &user, nil
}

// ==================== Transactions ====================

// TransferMoney 转账事务
func (c *Client) TransferMoney(ctx context.Context, fromAccount, toAccount string, amount float64) error {
    accountsColl := c.database.Collection("accounts")

    session, err := c.client.StartSession()
    if err != nil {
        return err
    }
    defer session.EndSession(ctx)

    _, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
        // 扣款
        _, err := accountsColl.UpdateOne(sessCtx,
            bson.M{"account_id": fromAccount, "balance": bson.M{"$gte": amount}},
            bson.M{"$inc": bson.M{"balance": -amount}},
        )
        if err != nil {
            return nil, err
        }

        // 收款
        _, err = accountsColl.UpdateOne(sessCtx,
            bson.M{"account_id": toAccount},
            bson.M{"$inc": bson.M{"balance": amount}},
        )
        if err != nil {
            return nil, err
        }

        return nil, nil
    })

    return err
}

// ==================== Aggregation Pipeline ====================

// SalesSummary 销售汇总
type SalesSummary struct {
    ProductID    string  `bson:"_id"`
    TotalRevenue float64 `bson:"total_revenue"`
    TotalUnits   int     `bson:"total_units"`
    AvgPrice     float64 `bson:"avg_price"`
    OrderCount   int     `bson:"order_count"`
}

// GetSalesSummary 获取销售汇总
func (c *Client) GetSalesSummary(ctx context.Context, from, to time.Time) ([]SalesSummary, error) {
    collection := c.database.Collection("orders")

    pipeline := mongo.Pipeline{
        // 阶段1: 匹配时间范围
        {{
            Key: "$match", Value: bson.M{
                "created_at": bson.M{"$gte": from, "$lte": to},
                "status":     "completed",
            },
        }},
        // 阶段2: 展开订单项
        {{
            Key: "$unwind", Value: "$items",
        }},
        // 阶段3: 按产品分组汇总
        {{
            Key: "$group", Value: bson.M{
                "_id": "$items.product_id",
                "total_revenue": bson.M{"$sum": "$items.total"},
                "total_units":   bson.M{"$sum": "$items.quantity"},
                "avg_price":     bson.M{"$avg": "$items.unit_price"},
                "order_count":   bson.M{"$sum": 1},
            },
        }},
        // 阶段4: 排序
        {{
            Key: "$sort", Value: bson.M{"total_revenue": -1},
        }},
    }

    cursor, err := collection.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var results []SalesSummary
    if err := cursor.All(ctx, &results); err != nil {
        return nil, err
    }

    return results, nil
}

// ==================== Indexes ====================

// CreateIndexes 创建常用索引
func (c *Client) CreateIndexes(ctx context.Context) error {
    // 产品集合索引
    productIndexes := []mongo.IndexModel{
        {
            Keys:    bson.D{{Key: "name", Value: "text"}, {Key: "description", Value: "text"}},
            Options: options.Index().SetName("text_search"),
        },
        {
            Keys:    bson.D{{Key: "category", Value: 1}, {Key: "price", Value: -1}},
            Options: options.Index().SetName("category_price"),
        },
        {
            Keys:    bson.D{{Key: "tags", Value: 1}},
            Options: options.Index().SetName("tags"),
        },
        {
            Keys:    bson.D{{Key: "location", Value: "2dsphere"}},
            Options: options.Index().SetName("geo"),
        },
    }

    _, err := c.database.Collection("products").Indexes().CreateMany(ctx, productIndexes)
    if err != nil {
        return err
    }

    // 订单集合索引
    orderIndexes := []mongo.IndexModel{
        {
            Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "created_at", Value: -1}},
            Options: options.Index().SetName("user_orders"),
        },
        {
            Keys:    bson.D{{Key: "status", Value: 1}},
            Options: options.Index().SetPartialFilterExpression(bson.M{"status": "pending"}),
        },
        {
            Keys:    bson.D{{Key: "created_at", Value: 1}},
            Options: options.Index().SetExpireAfterSeconds(30 * 24 * 60 * 60), // TTL 索引
        },
    }

    _, err = c.database.Collection("orders").Indexes().CreateMany(ctx, orderIndexes)
    return err
}
```

---

## 3. Configuration Best Practices

```yaml
# MongoDB 配置文件 mongod.conf

# ===== 存储配置 =====
storage:
  dbPath: /var/lib/mongodb
  journal:
    enabled: true
  engine: wiredTiger
  wiredTiger:
    engineConfig:
      cacheSizeGB: 8           # RAM 的 50%，最大不超过 32GB
      journalCompressor: zlib
      directoryForIndexes: false
    collectionConfig:
      blockCompressor: zstd    # none/snappy/zlib/zstd
    indexConfig:
      prefixCompression: true

# ===== 系统日志 =====
systemLog:
  destination: file
  logAppend: true
  path: /var/log/mongodb/mongod.log
  logRotate: reopen

# ===== 网络配置 =====
net:
  port: 27017
  bindIp: 127.0.0.1,10.0.0.1
  maxIncomingConnections: 65536
  wireObjectCheck: true

# ===== 进程管理 =====
processManagement:
  fork: true
  pidFilePath: /var/run/mongodb/mongod.pid
  timeZoneInfo: /usr/share/zoneinfo

# ===== 复制集配置 =====
replication:
  replSetName: rs0
  enableMajorityReadConcern: true

# ===== 分片配置 =====
sharding:
  clusterRole: shardsvr  # 或 configsvr

# ===== 安全配置 =====
security:
  authorization: enabled
  keyFile: /etc/mongodb/keyfile
  # enableEncryption: true
  # encryptionCipherMode: AES256-GCM

# ===== 性能分析 =====
operationProfiling:
  slowOpThresholdMs: 100
  mode: slowOp  # off/slowOp/all
```

---

## 4. Performance Tuning Guidelines

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    MongoDB Performance Tuning                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Index Optimization                                  │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  1. 复合索引设计                                                         │  │
│  │     - ESR 规则: Equality, Sort, Range                                  │  │
│  │     - {status: 1, created_at: -1} 优于 {created_at: -1, status: 1}     │  │
│  │                                                                        │  │
│  │  2. 覆盖查询 (Covered Query)                                            │  │
│  │     - 查询和投影都在索引中                                               │  │
│  │     - db.orders.find({user_id: "A"}, {_id: 0, status: 1})              │  │
│  │     - 索引: {user_id: 1, status: 1}                                   │  │
│  │                                                                        │  │
│  │  3. 索引交集                                                             │  │
│  │     - MongoDB 可以使用多个索引的交集                                     │  │
│  │     - 但复合索引通常更好                                                 │  │
│  │                                                                        │  │
│  │  4. 部分索引 (Partial Index)                                            │  │
│  │     - 只为符合条件的文档创建索引                                         │  │
│  │     - db.orders.createIndex(                                            │  │
│  │         {user_id: 1},                                                   │  │
│  │         {partialFilterExpression: {status: "pending"}}                  │  │
│  │       )                                                                 │  │
│  │                                                                        │  │
│  │  5. 稀疏索引 (Sparse Index)                                             │  │
│  │     - 只包含有该字段的文档                                               │  │
│  │     - 适用于可选的唯一字段                                               │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Query Optimization                                  │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  1. 使用投影减少返回字段                                                 │  │
│  │     db.users.find({}, {name: 1, email: 1})                             │  │
│  │                                                                        │  │
│  │  2. 使用 limit 和 skip (分页优化)                                        │  │
│  │     - 避免大 offset: skip(100000) 性能差                               │  │
│  │     - 使用范围查询: find({_id: {$gt: lastId}}).limit(20)               │  │
│  │                                                                        │  │
│  │  3. 避免 $where 和 JavaScript                                           │  │
│  │     - 使用聚合管道或标准操作符                                           │  │
│  │                                                                        │  │
│  │  4. $in 查询优化                                                         │  │
│  │     - 保持数组有序 (与索引顺序一致)                                      │  │
│  │     - 限制 $in 数组大小 (< 1000)                                        │  │
│  │                                                                        │  │
│  │  5. 聚合管道优化                                                         │  │
│  │     - 尽早使用 $match 减少数据量                                         │  │
│  │     - 使用 $project 减少字段传递                                         │  │
│  │     - 避免内存限制溢出 (allowDiskUse: true)                              │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Storage Tuning                                      │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  1. WiredTiger Cache                                                   │  │
│  │     - 默认: 50% RAM - 1GB                                              │  │
│  │     - 不要设置超过 32GB (压缩指针失效)                                   │  │
│  │     - 监控 cache usage 和 dirty percentage                             │  │
│  │                                                                        │  │
│  │  2. 文档大小                                                             │  │
│  │     - 保持文档 < 1MB (更新性能)                                          │  │
│  │     - 避免频繁增长的数组                                                 │  │
│  │                                                                        │  │
│  │  3. 预分配和填充                                                         │  │
│  │     - 使用 collation 进行排序规则优化                                    │  │
│  │     - 考虑使用 power-of-2 sizes (已废弃，现在默认)                       │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Go Driver Optimization                              │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  1. 连接池配置                                                           │  │
│  │     clientOptions.SetMaxPoolSize(100)                                  │  │
│  │     clientOptions.SetMinPoolSize(10)                                   │  │
│  │                                                                        │  │
│  │  2. 批量操作                                                             │  │
│  │     - InsertMany 替代多次 InsertOne                                     │  │
│  │     - BulkWrite 用于混合操作                                             │  │
│  │                                                                        │  │
│  │  3. 使用 context 控制超时                                                │  │
│  │     ctx, cancel := context.WithTimeout(context.Background(), 5s)        │  │
│  │                                                                        │  │
│  │  4. 重用 bson.M 对象 (谨慎)                                              │  │
│  │     - 避免在热路径频繁分配                                               │  │
│  │                                                                        │  │
│  │  5. 使用合适的解码方式                                                   │  │
│  │     - Decode(&struct) 比 bson.M 更快                                    │  │
│  │     - 预定义结构体避免反射                                               │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. Visual Representations

### 5.1 Document Growth and Padding

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Document Growth & Storage                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Document Update & Growth                            │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Initial Document:                                                     │  │
│  │  ┌──────────────────────────────────────────────────────────────┐     │  │
│  │  │  { name: "Product A", tags: ["electronics"] } │ Size: 100B  │     │  │
│  │  └──────────────────────────────────────────────────────────────┘     │  │
│  │                                                                        │  │
│  │  After Update (add tags):                                              │  │
│  │  ┌──────────────────────────────────────────────────────────────┐     │  │
│  │  │  { name: "Product A", tags: ["electronics", "gadget", ...] } │     │  │
│  │  └──────────────────────────────────────────────────────────────┘     │  │
│  │                                                                        │  │
│  │  Storage Behavior:                                                     │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                                                                 │  │  │
│  │  │  Record 1 (Original, now stale)                                 │  │  │
│  │  │  ┌─────────────────────────────────────────────────────────┐    │  │  │
│  │  │  │ { name: "Product A", tags: ["electronics"] } │ 100B     │    │  │  │
│  │  │  └─────────────────────────────────────────────────────────┘    │  │  │
│  │  │           │                                                     │  │  │
│  │  │           │ Update ──► New location if doesn't fit              │  │  │
│  │  │           ▼                                                     │  │  │
│  │  │  Record 2 (Updated)                                             │  │  │
│  │  │  ┌─────────────────────────────────────────────────────────┐    │  │  │
│  │  │  │ { name: "Product A", tags: ["electronics", ...] } │ 200B │    │  │  │
│  │  │  └─────────────────────────────────────────────────────────┘    │  │  │
│  │  │                                                                 │  │  │
│  │  │  Result: Document moved, space left in original location        │  │  │
│  │  │          (storage fragmentation)                                │  │  │
│  │  │                                                                 │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  Solutions:                                                            │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  1. Pre-allocation: Create document with expected final size    │  │  │
│  │  │  2. Bucketing: Separate frequently growing arrays               │  │  │
│  │  │  3. Compaction: use compact command to defragment               │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.2 Sharding Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    MongoDB Sharding Architecture                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                         Query Flow                                     │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Application                                                           │  │
│  │     │                                                                  │  │
│  │     │  db.users.find({user_id: "user_12345"})                          │  │
│  │     ▼                                                                  │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                        mongos (Router)                           │  │  │
│  │  │  1. Check cache: user_id hashed ──► shard?                      │  │  │
│  │  │  2. Cache miss ──► Query config servers                         │  │  │
│  │  │  3. Route to Shard 2                                            │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │     │                                                                  │  │
│  │     │  Route                                                            │  │
│  │     ▼                                                                  │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                    Config Servers (Replica Set)                  │  │  │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │  │  │
│  │  │  │ Config 1    │  │ Config 2    │  │ Config 3    │             │  │  │
│  │  │  │ (Primary)   │  │ (Secondary) │  │ (Secondary) │             │  │  │
│  │  │  └─────────────┘  └─────────────┘  └─────────────┘             │  │  │
│  │  │                                                                  │  │  │
│  │  │  Stores:                                                         │  │  │
│  │  │  - Chunk ranges (shard key ranges)                              │  │  │
│  │  │  - Shard locations                                              │  │  │
│  │  │  - Balancer state                                               │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │     ▲                                                                  │  │
│  │     │  Metadata                                                        │  │
│  │     │                                                                  │  │
│  │  ┌──┴───────────────────────────────────────────────────────────────┐  │  │
│  │  │                       Shard Cluster                               │  │  │
│  │  │  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐   │  │  │
│  │  │  │   Shard 0       │  │   Shard 1       │  │   Shard 2       │   │  │  │
│  │  │  │ ┌─────────────┐ │  │ ┌─────────────┐ │  │ ┌─────────────┐ │   │  │  │
│  │  │  │ │ Chunk       │ │  │ │ Chunk       │ │  │ │ Chunk       │ │   │  │  │
│  │  │  │ │ user_id:    │ │  │ │ user_id:    │ │  │ │ user_id:    │ │   │  │  │
│  │  │  │ │ [min, 100)  │ │  │ │ [100, 200)  │ │  │ │ [200, max)  │ │   │  │  │
│  │  │  │ └─────────────┘ │  │ └─────────────┘ │  │ └─────────────┘ │   │  │  │
│  │  │  │   Replica Set   │  │   Replica Set   │  │   Replica Set   │   │  │  │
│  │  │  └─────────────────┘  └─────────────────┘  └─────────────────┘   │  │  │
│  │  │                                                                    │  │  │
│  │  │  Chunk Split & Migration (Balancer):                               │  │  │
│  │  │  ┌────────────────────────────────────────────────────────────┐   │  │  │
│  │  │  │  Shard 2 chunk grows too large                             │   │  │  │
│  │  │  │  ─────────────────────────────────────────                 │   │  │  │
│  │  │  │  Split: [200, max) ──► [200, 300) + [300, max)             │   │  │  │
│  │  │  │                                                            │   │  │  │
│  │  │  │  Migrate [200, 300) to Shard 0 (to balance load)           │   │  │  │
│  │  │  └────────────────────────────────────────────────────────────┘   │  │  │
│  │  └───────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Shard Key Selection                                 │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Good Shard Key:                                                       │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  1. High Cardinality                                             │  │  │
│  │  │     - user_id (millions of unique values)                        │  │  │
│  │  │     - Avoid: country (only ~200 values)                          │  │  │
│  │  │                                                                  │  │  │
│  │  │  2. Even Distribution                                              │  │  │
│  │  │     - Hash sharding for monotonic _id                              │  │  │
│  │  │     - { _id: "hashed" }                                          │  │  │
│  │  │                                                                  │  │  │
│  │  │  3. Query Isolation                                                  │  │  │
│  │  │     - Targeted queries should include shard key                      │  │  │
│  │  │     - Find({user_id: "x"}) ──► Single shard                          │  │  │
│  │  │     - Find({}) ──► Scatter-gather (all shards)                       │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.3 Replica Set Oplog

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    MongoDB Oplog (Operations Log)                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Oplog Structure                                     │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  local.oplog.rs (Capped Collection, default 5% disk or 990MB)         │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  ts:   Timestamp    │  Operation timestamp                     │  │  │
│  │  │  t:    Long         │  Term (for Raft consensus)               │  │  │
│  │  │  h:    Long         │  Hash of the op (unique identifier)      │  │  │
│  │  │  v:    Int          │  Version of oplog entry format           │  │  │
│  │  │  op:   String       │  Operation type (i/u/d/c/n)              │  │  │
│  │  │  ns:   String       │  Namespace (db.collection)               │  │  │
│  │  │  ui:   UUID         │  Collection UUID                         │  │  │
│  │  │  o:    Document     │  Operation payload (insert/update/delete)│  │  │
│  │  │  o2:   Document     │  Query criteria (for updates)            │  │  │
│  │  │  wall: Date         │  Wall clock time                         │  │  │
│  │  │  prevOpTime: Object │  Previous operation timestamp            │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  Operation Types:                                                      │  │
│  │  ┌────────┬─────────────────────────────────────────────────────────┐  │  │
│  │  │  "i"   │ Insert document ──► o: {full document}                │  │  │
│  │  │  "u"   │ Update document ──► o: {$set: {...}}, o2: {_id: ...}  │  │  │
│  │  │  "d"   │ Delete document ──► o: {_id: ...}                     │  │  │
│  │  │  "c"   │ Command (DDL)     ──► o: {create, drop, ...}          │  │  │
│  │  │  "n"   │ No-op             ──► Used for heartbeats             │  │  │
│  │  └────────┴─────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Replication Flow                                    │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Primary                            Secondary                          │  │
│  │  ┌─────────────┐                    ┌─────────────┐                   │  │
│  │  │ Client Write│                    │ Oplog Query │                   │  │
│  │  │     │       │                    │     │       │                   │  │
│  │  │     ▼       │                    │     ▼       │                   │  │
│  │  │ ┌─────────┐ │                    │ ┌─────────┐ │                   │  │
│  │  │ │ Apply   │ │────oplog entry────►│ │ Fetch   │ │                   │  │
│  │  │ │ to Data │ │                    │ │ Entry   │ │                   │  │
│  │  │ └────┬────┘ │                    │ └────┬────┘ │                   │  │
│  │  │      │      │                    │      │      │                   │  │
│  │  │ ┌────┴────┐ │                    │ ┌────┴────┐ │                   │  │
│  │  │ │ Oplog   │ │◄────sync from──────┤ │ Apply   │ │                   │  │
│  │  │ │ (local) │ │   (timestamp)      │ │ to Data │ │                   │  │
│  │  │ └─────────┘ │                    │ └─────────┘ │                   │  │
│  │  └─────────────┘                    └─────────────┘                   │  │
│  │                                                                        │  │
│  │  Sync Process:                                                         │  │
│  │  1. Secondary queries oplog with last synced timestamp                 │  │
│  │  2. Primary returns entries since that timestamp                       │  │
│  │  3. Secondary applies entries in order                                 │  │
│  │  4. Update lastAppliedOpTime                                           │  │
│  │  5. Repeat (tailable cursor)                                           │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. References

1. **Chodorow, K.** (2013). MongoDB: The Definitive Guide, 2nd Edition. O'Reilly Media.
2. **Banker, K., Bakkum, P., et al.** (2016). MongoDB in Action, 2nd Edition. Manning Publications.
3. **MongoDB Documentation** (2024). docs.mongodb.com
4. **Kleppmann, M.** (2017). Designing Data-Intensive Applications. O'Reilly Media.
5. **MongoDB University** (2024). M320: Data Modeling Course.

---

*Document Version: 1.0 | Last Updated: 2024*

# TS-DB-005: MongoDB Architecture and Go Integration

> **维度**: Technology Stack > Database
> **级别**: S (16+ KB)
> **标签**: #mongodb #nosql #document #replica-set #sharding #go-mongo
> **权威来源**:
>
> - [MongoDB Documentation](https://docs.mongodb.com/) - MongoDB Inc.
> - [MongoDB WiredTiger](https://docs.mongodb.com/manual/core/wiredtiger/) - Storage Engine
> - [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver) - Official driver

---

## 1. MongoDB Architecture

### 1.1 Document Model

```
┌─────────────────────────────────────────────────────────────────┐
│                    MongoDB Document Structure                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  BSON Document (Binary JSON):                                   │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ Document Size (4 bytes)                                 │   │
│  ├─────────────────────────────────────────────────────────┤   │
│  │ Element 1:                                              │   │
│  │   Field Name: "_id"                                     │   │
│  │   Type: 0x07 (ObjectId)                                 │   │
│  │   Value: 12-byte ObjectId                               │   │
│  ├─────────────────────────────────────────────────────────┤   │
│  │ Element 2:                                              │   │
│  │   Field Name: "name"                                    │   │
│  │   Type: 0x02 (String)                                   │   │
│  │   Value: Length + UTF-8 string + null                   │   │
│  ├─────────────────────────────────────────────────────────┤   │
│  │ ...                                                     │   │
│  ├─────────────────────────────────────────────────────────┤   │
│  │ Null terminator (0x00)                                  │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                  │
│  BSON Types:                                                    │
│  - Double (0x01), String (0x02), Document (0x03)               │
│  - Array (0x04), Binary (0x05), Undefined (0x06, deprecated)   │
│  - ObjectId (0x07), Boolean (0x08), DateTime (0x09)            │
│  - Null (0x0A), Regex (0x0B), DBPointer (0x0C, deprecated)     │
│  - JavaScript (0x0D), Symbol (0x0E), Int32 (0x10)              │
│  - Timestamp (0x11), Int64 (0x12), Decimal128 (0x13)           │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 1.2 Storage Architecture (WiredTiger)

```
┌─────────────────────────────────────────────────────────────────┐
│                    WiredTiger Storage Engine                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                     MongoDB Layer                          │  │
│  │  (Query Planner, Index Selection, Document Validation)     │  │
│  └───────────────────────┬───────────────────────────────────┘  │
│                          │                                       │
│  ┌───────────────────────▼───────────────────────────────────┐  │
│  │                    WiredTiger API                          │  │
│  │  (Sessions, Transactions, Cursors, Checkpoints)            │  │
│  └───────────────────────┬───────────────────────────────────┘  │
│                          │                                       │
│  ┌───────────────────────▼───────────────────────────────────┐  │
│  │                   B-Tree Layer                              │  │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐                    │  │
│  │  │ Root    │──►│ Internal│──►│ Leaf    │                    │  │
│  │  │ Page    │  │ Pages   │  │ Pages   │                    │  │
│  │  └─────────┘  └─────────┘  └─────────┘                    │  │
│  │  Page size: 32KB default                                   │  │
│  └───────────────────────┬───────────────────────────────────┘  │
│                          │                                       │
│  ┌───────────────────────▼───────────────────────────────────┐  │
│  │              Cache & Eviction Layer                        │  │
│  │  ┌───────────────────────────────────────────────┐        │  │
│  │  │  Dirty Cache (modified pages)                 │        │  │
│  │  ├───────────────────────────────────────────────┤        │  │
│  │  │  Clean Cache (read pages)                     │        │  │
│  │  └───────────────────────────────────────────────┘        │  │
│  │  Cache size: 50% RAM - 1GB (default)                      │  │
│  └───────────────────────┬───────────────────────────────────┘  │
│                          │                                       │
│  ┌───────────────────────▼───────────────────────────────────┐  │
│  │                  Logging & Recovery                        │  │
│  │  ┌─────────────────────────────────────────────────────┐  │  │
│  │  │  Write-Ahead Log (WAL/Journal)                      │  │  │
│  │  │  - Compresses at 128KB boundaries                   │  │  │
│  │  │  - 100ms fsync interval (default)                   │  │  │
│  │  └─────────────────────────────────────────────────────┘  │  │
│  │  Checkpoints: Every 60s or 2GB of journal                 │  │
│  └───────────────────────────────────────────────────────────┘  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 2. Replication Architecture

### 2.1 Replica Set

```
┌─────────────────────────────────────────────────────────────────┐
│                      Replica Set Architecture                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│                    ┌─────────────┐                              │
│                    │   Primary   │ ◄──── Write operations       │
│                    │  (Active)   │                              │
│                    └──────┬──────┘                              │
│                           │                                      │
│         ┌─────────────────┼─────────────────┐                   │
│         │                 │                 │                   │
│         ▼                 ▼                 ▼                   │
│  ┌─────────────┐   ┌─────────────┐   ┌─────────────┐           │
│  │  Secondary  │   │  Secondary  │   │   Arbiter   │           │
│  │  (Active)   │   │  (Active)   │   │  (Voting)   │           │
│  └─────────────┘   └─────────────┘   └─────────────┘           │
│                                                                  │
│  Replication:                                                    │
│  - Asynchronous oplog replication                               │
│  - Oplog (operations log): capped collection on local           │
│  - Default size: 5% disk space (min 990MB)                      │
│                                                                  │
│  Election Process:                                               │
│  - Priority-based voting                                        │
│  - Majority needed for primary election                         │
│  - Arbiter participates only in elections                       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 Oplog Structure

```go
// Oplog entry structure
type OplogEntry struct {
    Timestamp  primitive.Timestamp `bson:"ts"`       // Operation time
    Hash       int64               `bson:"h"`        // Document hash
    Version    int                 `bson:"v"`        // Version
    Operation  string              `bson:"op"`       // i=insert, u=update, d=delete
    Namespace  string              `bson:"ns"`       // db.collection
    Object     bson.D              `bson:"o"`        // Full document (insert)
    Query      bson.D              `bson:"o2"`       // Query (update/delete)
    Update     bson.D              `bson:"o"`        // Update operators
}
```

---

## 3. Sharding Architecture

### 3.1 Shard Cluster

```
┌─────────────────────────────────────────────────────────────────┐
│                     Sharded Cluster Architecture                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                      Application                           │  │
│  └───────────────────────┬───────────────────────────────────┘  │
│                          │                                       │
│  ┌───────────────────────▼───────────────────────────────────┐  │
│  │                   mongos (Router)                          │  │
│  │  - Routes queries to appropriate shards                    │  │
│  │  - Caches cluster metadata from config servers             │  │
│  │  - No persistent state                                     │  │
│  └───────────────────────┬───────────────────────────────────┘  │
│                          │                                       │
│         ┌────────────────┼────────────────┐                     │
│         │                │                │                     │
│  ┌──────▼──────┐  ┌─────▼──────┐  ┌──────▼──────┐             │
│  │   Shard 1   │  │   Shard 2  │  │   Shard N   │             │
│  │ ┌─────────┐ │  │ ┌─────────┐│  │ ┌─────────┐ │             │
│  │ │Primary  │ │  │ │Primary  ││  │ │Primary  │ │             │
│  │ │Secondary│ │  │ │Secondary││  │ │Secondary│ │             │
│  │ └─────────┘ │  │ └─────────┘│  │ └─────────┘ │             │
│  │  Chunk 1    │  │  Chunk 2   │  │  Chunk N    │             │
│  └─────────────┘  └────────────┘  └─────────────┘             │
│                                                                  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                Config Servers (Replica Set)                │  │
│  │  - Store cluster metadata                                  │  │
│  │  - Chunk distribution, shard ranges                        │  │
│  │  - Must be deployed as replica set (3 nodes)               │  │
│  └───────────────────────────────────────────────────────────┘  │
│                                                                  │
│  Chunk: Contiguous range of shard key values (64MB default)     │
│  Balancer: Automatic chunk migration for even distribution      │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 4. Go Integration

### 4.1 Connection and Basic Operations

```go
package main

import (
    "context"
    "fmt"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connection configuration
func connect() (*mongo.Client, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().
        ApplyURI("mongodb://localhost:27017").
        SetMaxPoolSize(100).
        SetMinPoolSize(10).
        SetMaxConnIdleTime(30 * time.Second).
        SetServerSelectionTimeout(5 * time.Second).
        SetRetryWrites(true).
        SetRetryReads(true))

    if err != nil {
        return nil, err
    }

    // Verify connection
    if err := client.Ping(ctx, readpref.Primary()); err != nil {
        return nil, err
    }

    return client, nil
}

// Document structure
type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name      string             `bson:"name" json:"name"`
    Email     string             `bson:"email" json:"email"`
    Age       int                `bson:"age" json:"age"`
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
    Tags      []string           `bson:"tags" json:"tags"`
    Address   Address            `bson:"address" json:"address"`
}

type Address struct {
    Street  string `bson:"street" json:"street"`
    City    string `bson:"city" json:"city"`
    Country string `bson:"country" json:"country"`
}

// CRUD Operations
func crudOperations(client *mongo.Client) error {
    ctx := context.Background()
    collection := client.Database("app").Collection("users")

    // CREATE
    user := User{
        Name:      "John Doe",
        Email:     "john@example.com",
        Age:       30,
        CreatedAt: time.Now(),
        Tags:      []string{"developer", "golang"},
        Address: Address{
            Street:  "123 Main St",
            City:    "San Francisco",
            Country: "USA",
        },
    }

    result, err := collection.InsertOne(ctx, user)
    if err != nil {
        return err
    }
    fmt.Printf("Inserted ID: %v\n", result.InsertedID)

    // READ with filter
    var found User
    err = collection.FindOne(ctx, bson.M{"email": "john@example.com"}).Decode(&found)
    if err == mongo.ErrNoDocuments {
        fmt.Println("User not found")
    } else if err != nil {
        return err
    }

    // UPDATE
    update := bson.M{
        "$set": bson.M{
            "age": 31,
        },
        "$push": bson.M{
            "tags": "senior",
        },
    }
    _, err = collection.UpdateOne(ctx, bson.M{"email": "john@example.com"}, update)

    // DELETE
    _, err = collection.DeleteOne(ctx, bson.M{"email": "john@example.com"})

    return err
}
```

### 4.2 Query Operations

```go
// Advanced queries
func advancedQueries(collection *mongo.Collection) error {
    ctx := context.Background()

    // Find with options
    opts := options.Find().
        SetLimit(10).
        SetSkip(20).
        SetSort(bson.D{{"created_at", -1}}).
        SetProjection(bson.M{"password": 0}) // Exclude password field

    cursor, err := collection.Find(ctx, bson.M{
        "age": bson.M{"$gte": 18, "$lte": 65},
        "tags": bson.M{"$in": []string{"developer"}},
    }, opts)
    if err != nil {
        return err
    }
    defer cursor.Close(ctx)

    var users []User
    if err = cursor.All(ctx, &users); err != nil {
        return err
    }

    // Aggregation pipeline
    pipeline := mongo.Pipeline{
        {{"$match", bson.M{"status": "active"}}},
        {{"$group", bson.M{
            "_id": "$department",
            "avg_age": bson.M{"$avg": "$age"},
            "count": bson.M{"$sum": 1},
        }}},
        {{"$sort", bson.M{"count": -1}}},
    }

    aggCursor, err := collection.Aggregate(ctx, pipeline)
    if err != nil {
        return err
    }
    defer aggCursor.Close(ctx)

    return nil
}
```

### 4.3 Transactions

```go
// Multi-document transaction
func transferFunds(client *mongo.Client, fromID, toID string, amount float64) error {
    ctx := context.Background()

    session, err := client.StartSession()
    if err != nil {
        return err
    }
    defer session.EndSession(ctx)

    callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
        accounts := client.Database("bank").Collection("accounts")

        // Deduct from sender
        _, err := accounts.UpdateOne(sessCtx,
            bson.M{"_id": fromID, "balance": bson.M{"$gte": amount}},
            bson.M{"$inc": bson.M{"balance": -amount}})
        if err != nil {
            return nil, err
        }

        // Add to receiver
        _, err = accounts.UpdateOne(sessCtx,
            bson.M{"_id": toID},
            bson.M{"$inc": bson.M{"balance": amount}})

        return nil, err
    }

    _, err = session.WithTransaction(ctx, callback)
    return err
}
```

### 4.4 Connection to Replica Set

```go
func connectReplicaSet() (*mongo.Client, error) {
    uri := "mongodb://user:pass@host1:27017,host2:27017,host3:27017/dbname?replicaSet=rs0"

    clientOpts := options.Client().ApplyURI(uri).
        SetReadPreference(readpref.SecondaryPreferred()).
        SetRetryWrites(true)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, clientOpts)
    if err != nil {
        return nil, err
    }

    return client, client.Ping(ctx, readpref.SecondaryPreferred())
}

// Read preferences
func readPreferences() {
    // Primary - default, read from primary only
    readpref.Primary()

    // PrimaryPreferred - read from primary, fallback to secondary
    readpref.PrimaryPreferred()

    // Secondary - read from secondary only
    readpref.Secondary()

    // SecondaryPreferred - read from secondary, fallback to primary
    readpref.SecondaryPreferred()

    // Nearest - read from nearest member by latency
    readpref.Nearest()
}
```

---

## 5. Indexing Strategies

```go
// Index creation
func createIndexes(collection *mongo.Collection) error {
    ctx := context.Background()

    // Single field index
    _, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
        Keys:    bson.D{{"email", 1}}, // 1 = ascending
        Options: options.Index().SetUnique(true),
    })

    // Compound index
    _, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
        Keys: bson.D{
            {"department", 1},
            {"created_at", -1},
        },
    })

    // Text index for search
    _, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
        Keys: bson.D{{"$text", bson.D{
            {"title", "text"},
            {"content", "text"},
        }}},
    })

    // TTL index (auto-expire documents)
    _, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
        Keys:    bson.D{{"created_at", 1}},
        Options: options.Index().SetExpireAfterSeconds(3600), // 1 hour
    })

    return err
}
```

---

## 6. Monitoring and Diagnostics

```go
// Get server status
func getServerStatus(client *mongo.Client) (bson.M, error) {
    ctx := context.Background()
    var result bson.M
    err := client.Database("admin").RunCommand(ctx, bson.D{{"serverStatus", 1}}).Decode(&result)
    return result, err
}

// Current operations
func getCurrentOps(client *mongo.Client) error {
    ctx := context.Background()

    cursor, err := client.Database("admin").Collection("$cmd.sys.inprog").
        Find(ctx, bson.M{"$all": true})
    if err != nil {
        return err
    }
    defer cursor.Close(ctx)

    var ops []bson.M
    return cursor.All(ctx, &ops)
}
```

---

## 7. Checklist

```
MongoDB Production Checklist:
□ Deploy as replica set (minimum 3 nodes)
□ Enable authentication
□ Configure proper firewall rules
□ Set up monitoring (MongoDB Atlas or self-hosted)
□ Configure backup strategy
□ Size oplog appropriately
□ Create necessary indexes
□ Configure working set to fit in RAM
□ Set up alerts for slow queries

Go Application Checklist:
□ Use connection pooling
□ Always close cursors
□ Handle context cancellation
□ Use transactions for multi-document operations
□ Implement retry logic for transient errors
□ Create indexes before going to production
□ Monitor connection pool statistics
```

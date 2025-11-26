# MongoDB编程 - Go语言实战指南

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [MongoDB编程 - Go语言实战指南](#mongodb编程---go语言实战指南)
  - [📋 目录](#-目录)
  - [MongoDB概述](#mongodb概述)
    - [特点](#特点)
    - [适用场景](#适用场景)
  - [安装与配置](#安装与配置)
    - [安装驱动](#安装驱动)
    - [项目结构](#项目结构)
  - [连接管理](#连接管理)
    - [基本连接](#基本连接)
    - [连接字符串配置](#连接字符串配置)
  - [CRUD操作](#crud操作)
    - [数据模型定义](#数据模型定义)
    - [插入操作](#插入操作)
    - [查询操作](#查询操作)
    - [更新操作](#更新操作)
    - [删除操作](#删除操作)
  - [聚合查询](#聚合查询)
    - [聚合管道](#聚合管道)
  - [索引管理](#索引管理)
    - [创建索引](#创建索引)
  - [事务处理](#事务处理)
    - [ACID事务](#acid事务)
  - [Change Streams](#change-streams)
    - [实时监听数据变化](#实时监听数据变化)
  - [GridFS文件存储](#gridfs文件存储)
    - [大文件存储](#大文件存储)
  - [性能优化](#性能优化)
    - [连接池优化](#连接池优化)
    - [批量操作](#批量操作)
    - [投影查询](#投影查询)
  - [最佳实践](#最佳实践)
    - [1. 错误处理](#1-错误处理)
    - [2. Context超时控制](#2-context超时控制)
    - [3. 连接复用](#3-连接复用)
    - [4. 索引策略](#4-索引策略)
    - [5. 数据建模](#5-数据建模)
  - [总结](#总结)

---

## MongoDB概述

### 特点

MongoDB是一个基于文档的NoSQL数据库：

- **文档型存储**: 使用BSON格式（类似JSON）
- **Schema灵活**: 无需预定义表结构
- **高性能**: 内存映射、索引支持
- **水平扩展**: 分片（Sharding）支持
- **高可用**: 副本集（Replica Set）

### 适用场景

```text
✅ 内容管理系统
✅ 实时分析
✅ 物联网数据存储
✅ 移动应用后端
✅ 目录服务
✅ 社交网络
```

---

## 安装与配置

### 安装驱动

```bash
go get go.mongodb.org/mongo-driver/mongo
go get go.mongodb.org/mongo-driver/bson
```

### 项目结构

```text
project/
├── main.go
├── config/
│   └── database.go      # 数据库配置
├── models/
│   ├── user.go          # 数据模型
│   └── product.go
├── repository/
│   ├── user_repo.go     # 数据访问层
│   └── product_repo.go
└── go.mod
```

---

## 连接管理

### 基本连接

```go
package database

import (
    "context"
    "fmt"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDB 数据库实例
type MongoDB struct {
    Client   *mongo.Client
    Database *mongo.Database
}

// NewMongoDB 创建MongoDB连接
func NewMongoDB(uri, dbName string) (*MongoDB, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // 设置客户端选项
    clientOptions := options.Client().
        ApplyURI(uri).
        SetMaxPoolSize(100).                    // 最大连接池大小
        SetMinPoolSize(10).                     // 最小连接池大小
        SetMaxConnIdleTime(30 * time.Second).   // 最大空闲时间
        SetConnectTimeout(10 * time.Second).    // 连接超时
        SetServerSelectionTimeout(5 * time.Second) // 服务器选择超时

    // 连接MongoDB
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return nil, fmt.Errorf("failed to connect: %w", err)
    }

    // 测试连接
    if err := client.Ping(ctx, readpref.Primary()); err != nil {
        return nil, fmt.Errorf("ping failed: %w", err)
    }

    fmt.Println("✅ MongoDB连接成功")

    return &MongoDB{
        Client:   client,
        Database: client.Database(dbName),
    }, nil
}

// Close 关闭连接
func (m *MongoDB) Close() error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    return m.Client.Disconnect(ctx)
}

// GetCollection 获取集合
func (m *MongoDB) GetCollection(name string) *mongo.Collection {
    return m.Database.Collection(name)
}
```

### 连接字符串配置

```go
const (
    // 本地开发
    LocalURI = "mongodb://localhost:27017"

    // 带认证
    AuthURI = "mongodb://username:password@localhost:27017/?authSource=admin"

    // 副本集
    ReplicaSetURI = "mongodb://host1:27017,host2:27017,host3:27017/?replicaSet=myReplSet"

    // MongoDB Atlas（云服务）
    AtlasURI = "mongodb+srv://username:password@cluster0.mongodb.net/?retryWrites=true&w=majority"
)
```

---

## CRUD操作

### 数据模型定义

```go
package models

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

// User 用户模型
type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Username  string             `bson:"username" json:"username"`
    Email     string             `bson:"email" json:"email"`
    Age       int                `bson:"age" json:"age"`
    Tags      []string           `bson:"tags,omitempty" json:"tags,omitempty"`
    Profile   Profile            `bson:"profile" json:"profile"`
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// Profile 用户资料（嵌套文档）
type Profile struct {
    Bio       string `bson:"bio,omitempty" json:"bio,omitempty"`
    Avatar    string `bson:"avatar,omitempty" json:"avatar,omitempty"`
    Location  string `bson:"location,omitempty" json:"location,omitempty"`
}
```

### 插入操作

```go
package repository

import (
    "context"
    "fmt"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

// UserRepository 用户仓储
type UserRepository struct {
    collection *mongo.Collection
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *mongo.Database) *UserRepository {
    return &UserRepository{
        collection: db.Collection("users"),
    }
}

// Insert 插入单个用户
func (r *UserRepository) Insert(ctx Context.Context, user *models.User) error {
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()

    result, err := r.collection.InsertOne(ctx, user)
    if err != nil {
        return fmt.Errorf("insert failed: %w", err)
    }

    user.ID = result.InsertedID.(primitive.ObjectID)
    fmt.Printf("✅ 插入成功，ID: %s\n", user.ID.Hex())

    return nil
}

// InsertMany 批量插入
func (r *UserRepository) InsertMany(ctx Context.Context, users []*models.User) error {
    docs := make([]interface{}, len(users))
    for i, user := range users {
        user.CreatedAt = time.Now()
        user.UpdatedAt = time.Now()
        docs[i] = user
    }

    result, err := r.collection.InsertMany(ctx, docs)
    if err != nil {
        return fmt.Errorf("insert many failed: %w", err)
    }

    fmt.Printf("✅ 批量插入成功，插入了 %d 条记录\n", len(result.InsertedIDs))

    return nil
}
```

### 查询操作

```go
// FindByID 根据ID查询
func (r *UserRepository) FindByID(ctx Context.Context, id primitive.ObjectID) (*models.User, error) {
    var user models.User

    filter := bson.M{"_id": id}
    err := r.collection.FindOne(ctx, filter).Decode(&user)
    if err == mongo.ErrNoDocuments {
        return nil, fmt.Errorf("user not found")
    }
    if err != nil {
        return nil, err
    }

    return &user, nil
}

// FindByUsername 根据用户名查询
func (r *UserRepository) FindByUsername(ctx Context.Context, username string) (*models.User, error) {
    var user models.User

    filter := bson.M{"username": username}
    err := r.collection.FindOne(ctx, filter).Decode(&user)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

// FindAll 查询所有用户（带分页）
func (r *UserRepository) FindAll(ctx Context.Context, page, pageSize int64) ([]*models.User, int64, error) {
    // 计算跳过数量
    skip := (page - 1) * pageSize

    // 设置查询选项
    opts := options.Find().
        SetLimit(pageSize).
        SetSkip(skip).
        SetSort(bson.D{{Key: "created_at", Value: -1}}) // 按创建时间降序

    // 执行查询
    cursor, err := r.collection.Find(ctx, bson.M{}, opts)
    if err != nil {
        return nil, 0, err
    }
    defer cursor.Close(ctx)

    // 解码结果
    var users []*models.User
    if err := cursor.All(ctx, &users); err != nil {
        return nil, 0, err
    }

    // 获取总数
    total, err := r.collection.CountDocuments(ctx, bson.M{})
    if err != nil {
        return nil, 0, err
    }

    return users, total, nil
}

// FindByAge 根据年龄范围查询
func (r *UserRepository) FindByAge(ctx Context.Context, minAge, maxAge int) ([]*models.User, error) {
    filter := bson.M{
        "age": bson.M{
            "$gte": minAge,
            "$lte": maxAge,
        },
    }

    cursor, err := r.collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var users []*models.User
    if err := cursor.All(ctx, &users); err != nil {
        return nil, err
    }

    return users, nil
}

// FindByTags 根据标签查询（数组查询）
func (r *UserRepository) FindByTags(ctx Context.Context, tags []string) ([]*models.User, error) {
    filter := bson.M{
        "tags": bson.M{
            "$in": tags, // 包含任意一个标签
        },
    }

    cursor, err := r.collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var users []*models.User
    if err := cursor.All(ctx, &users); err != nil {
        return nil, err
    }

    return users, nil
}

// Search 全文搜索（需要创建文本索引）
func (r *UserRepository) Search(ctx Context.Context, keyword string) ([]*models.User, error) {
    filter := bson.M{
        "$text": bson.M{
            "$search": keyword,
        },
    }

    // 按相关性排序
    opts := options.Find().SetProjection(bson.M{
        "score": bson.M{"$meta": "textScore"},
    }).SetSort(bson.D{{Key: "score", Value: bson.M{"$meta": "textScore"}}})

    cursor, err := r.collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var users []*models.User
    if err := cursor.All(ctx, &users); err != nil {
        return nil, err
    }

    return users, nil
}
```

### 更新操作

```go
// Update 更新用户
func (r *UserRepository) Update(ctx Context.Context, id primitive.ObjectID, update bson.M) error {
    filter := bson.M{"_id": id}

    // 添加更新时间
    update["updated_at"] = time.Now()

    updateDoc := bson.M{"$set": update}

    result, err := r.collection.UpdateOne(ctx, filter, updateDoc)
    if err != nil {
        return err
    }

    if result.MatchedCount == 0 {
        return fmt.Errorf("user not found")
    }

    fmt.Printf("✅ 更新成功，修改了 %d 条记录\n", result.ModifiedCount)

    return nil
}

// UpdateUsername 更新用户名
func (r *UserRepository) UpdateUsername(ctx Context.Context, id primitive.ObjectID, newUsername string) error {
    return r.Update(ctx, id, bson.M{"username": newUsername})
}

// IncrementAge 增加年龄（原子操作）
func (r *UserRepository) IncrementAge(ctx Context.Context, id primitive.ObjectID, increment int) error {
    filter := bson.M{"_id": id}
    update := bson.M{
        "$inc": bson.M{"age": increment},
        "$set": bson.M{"updated_at": time.Now()},
    }

    result, err := r.collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }

    if result.MatchedCount == 0 {
        return fmt.Errorf("user not found")
    }

    return nil
}

// AddTag 添加标签
func (r *UserRepository) AddTag(ctx Context.Context, id primitive.ObjectID, tag string) error {
    filter := bson.M{"_id": id}
    update := bson.M{
        "$addToSet": bson.M{"tags": tag}, // 避免重复
        "$set":      bson.M{"updated_at": time.Now()},
    }

    _, err := r.collection.UpdateOne(ctx, filter, update)
    return err
}

// RemoveTag 移除标签
func (r *UserRepository) RemoveTag(ctx Context.Context, id primitive.ObjectID, tag string) error {
    filter := bson.M{"_id": id}
    update := bson.M{
        "$pull": bson.M{"tags": tag},
        "$set":  bson.M{"updated_at": time.Now()},
    }

    _, err := r.collection.UpdateOne(ctx, filter, update)
    return err
}

// UpdateMany 批量更新
func (r *UserRepository) UpdateMany(ctx Context.Context, filter, update bson.M) (int64, error) {
    update["updated_at"] = time.Now()
    updateDoc := bson.M{"$set": update}

    result, err := r.collection.UpdateMany(ctx, filter, updateDoc)
    if err != nil {
        return 0, err
    }

    return result.ModifiedCount, nil
}
```

### 删除操作

```go
// Delete 删除用户
func (r *UserRepository) Delete(ctx Context.Context, id primitive.ObjectID) error {
    filter := bson.M{"_id": id}

    result, err := r.collection.DeleteOne(ctx, filter)
    if err != nil {
        return err
    }

    if result.DeletedCount == 0 {
        return fmt.Errorf("user not found")
    }

    fmt.Println("✅ 删除成功")

    return nil
}

// DeleteByUsername 根据用户名删除
func (r *UserRepository) DeleteByUsername(ctx Context.Context, username string) error {
    filter := bson.M{"username": username}

    _, err := r.collection.DeleteOne(ctx, filter)
    return err
}

// DeleteMany 批量删除
func (r *UserRepository) DeleteMany(ctx Context.Context, filter bson.M) (int64, error) {
    result, err := r.collection.DeleteMany(ctx, filter)
    if err != nil {
        return 0, err
    }

    return result.DeletedCount, nil
}

// DeleteOldUsers 删除旧用户（示例：删除超过1年未登录的用户）
func (r *UserRepository) DeleteOldUsers(ctx Context.Context, daysAgo int) (int64, error) {
    cutoff := time.Now().AddDate(0, 0, -daysAgo)

    filter := bson.M{
        "last_login_at": bson.M{
            "$lt": cutoff,
        },
    }

    result, err := r.collection.DeleteMany(ctx, filter)
    if err != nil {
        return 0, err
    }

    return result.DeletedCount, nil
}
```

---

## 聚合查询

### 聚合管道

```go
package repository

import (
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

// AggregateByAge 按年龄分组统计
func (r *UserRepository) AggregateByAge(ctx Context.Context) ([]bson.M, error) {
    pipeline := mongo.Pipeline{
        // 阶段1: 分组
        {{Key: "$group", Value: bson.D{
            {Key: "_id", Value: "$age"},
            {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
            {Key: "usernames", Value: bson.D{{Key: "$push", Value: "$username"}}},
        }}},
        // 阶段2: 排序
        {{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
    }

    cursor, err := r.collection.Aggregate(ctx, pipeline)
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

// GetStatistics 获取用户统计信息
func (r *UserRepository) GetStatistics(ctx Context.Context) (*UserStatistics, error) {
    pipeline := mongo.Pipeline{
        {{Key: "$group", Value: bson.D{
            {Key: "_id", Value: nil},
            {Key: "total", Value: bson.D{{Key: "$sum", Value: 1}}},
            {Key: "avgAge", Value: bson.D{{Key: "$avg", Value: "$age"}}},
            {Key: "minAge", Value: bson.D{{Key: "$min", Value: "$age"}}},
            {Key: "$maxAge", Value: bson.D{{Key: "$max", Value: "$age"}}},
        }}},
    }

    cursor, err := r.collection.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var results []UserStatistics
    if err := cursor.All(ctx, &results); err != nil {
        return nil, err
    }

    if len(results) == 0 {
        return nil, fmt.Errorf("no data")
    }

    return &results[0], nil
}

// UserStatistics 用户统计
type UserStatistics struct {
    Total  int     `bson:"total"`
    AvgAge float64 `bson:"avgAge"`
    MinAge int     `bson:"minAge"`
    MaxAge int     `bson:"maxAge"`
}

// AggregateTopTags 获取最热门的标签
func (r *UserRepository) AggregateTopTags(ctx Context.Context, limit int) ([]TagCount, error) {
    pipeline := mongo.Pipeline{
        // 阶段1: 展开数组
        {{Key: "$unwind", Value: "$tags"}},
        // 阶段2: 分组计数
        {{Key: "$group", Value: bson.D{
            {Key: "_id", Value: "$tags"},
            {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
        }}},
        // 阶段3: 排序
        {{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
        // 阶段4: 限制结果数量
        {{Key: "$limit", Value: limit}},
        // 阶段5: 重命名字段
        {{Key: "$project", Value: bson.D{
            {Key: "tag", Value: "$_id"},
            {Key: "count", Value: 1},
            {Key: "_id", Value: 0},
        }}},
    }

    cursor, err := r.collection.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var results []TagCount
    if err := cursor.All(ctx, &results); err != nil {
        return nil, err
    }

    return results, nil
}

// TagCount 标签计数
type TagCount struct {
    Tag   string `bson:"tag"`
    Count int    `bson:"count"`
}

// SearchWithFacets 带分面的搜索
func (r *UserRepository) SearchWithFacets(ctx Context.Context, keyword string) (*SearchResult, error) {
    pipeline := mongo.Pipeline{
        // 阶段1: 文本搜索
        {{Key: "$match", Value: bson.M{
            "$text": bson.M{"$search": keyword},
        }}},
        // 阶段2: 分面聚合
        {{Key: "$facet", Value: bson.D{
            // 分面1: 按年龄段
            {Key: "byAge", Value: []bson.D{
                {{Key: "$bucket", Value: bson.D{
                    {Key: "groupBy", Value: "$age"},
                    {Key: "boundaries", Value: []int{0, 18, 30, 50, 100}},
                    {Key: "default", Value: "Other"},
                    {Key: "output", Value: bson.D{
                        {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
                    }},
                }}},
            }},
            // 分面2: 按位置
            {Key: "byLocation", Value: []bson.D{
                {{Key: "$group", Value: bson.D{
                    {Key: "_id", Value: "$profile.location"},
                    {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
                }}},
                {{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
                {{Key: "$limit", Value: 10}},
            }},
            // 分面3: 结果列表
            {Key: "results", Value: []bson.D{
                {{Key: "$limit", Value: 20}},
            }},
        }}},
    }

    cursor, err := r.collection.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var results []SearchResult
    if err := cursor.All(ctx, &results); err != nil {
        return nil, err
    }

    if len(results) == 0 {
        return nil, fmt.Errorf("no results")
    }

    return &results[0], nil
}

// SearchResult 搜索结果
type SearchResult struct {
    ByAge      []bson.M        `bson:"byAge"`
    ByLocation []bson.M        `bson:"byLocation"`
    Results    []*models.User  `bson:"results"`
}
```

---

## 索引管理

### 创建索引

```go
package repository

import (
    "context"
    "fmt"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// CreateIndexes 创建索引
func (r *UserRepository) CreateIndexes(ctx Context.Context) error {
    indexes := []mongo.IndexModel{
        // 1. 单字段索引（用户名唯一）
        {
            Keys:    bson.D{{Key: "username", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
        // 2. 单字段索引（邮箱唯一）
        {
            Keys:    bson.D{{Key: "email", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
        // 3. 复合索引（年龄+创建时间）
        {
            Keys: bson.D{
                {Key: "age", Value: 1},
                {Key: "created_at", Value: -1},
            },
        },
        // 4. 文本索引（用于全文搜索）
        {
            Keys: bson.D{
                {Key: "username", Value: "text"},
                {Key: "profile.bio", Value: "text"},
            },
            Options: options.Index().SetWeights(bson.M{
                "username":    10, // 用户名权重更高
                "profile.bio": 5,
            }),
        },
        // 5. TTL索引（自动删除过期文档）
        {
            Keys: bson.D{{Key: "created_at", Value: 1}},
            Options: options.Index().SetExpireAfterSeconds(86400 * 30), // 30天后过期
        },
    }

    // 批量创建索引
    names, err := r.collection.Indexes().CreateMany(ctx, indexes)
    if err != nil {
        return fmt.Errorf("create indexes failed: %w", err)
    }

    fmt.Printf("✅ 创建索引成功: %v\n", names)

    return nil
}

// ListIndexes 列出所有索引
func (r *UserRepository) ListIndexes(ctx Context.Context) error {
    cursor, err := r.collection.Indexes().List(ctx)
    if err != nil {
        return err
    }
    defer cursor.Close(ctx)

    var indexes []bson.M
    if err := cursor.All(ctx, &indexes); err != nil {
        return err
    }

    fmt.Println("=== 索引列表 ===")
    for _, index := range indexes {
        fmt.Printf("名称: %v, 键: %v\n", index["name"], index["key"])
    }

    return nil
}

// DropIndex 删除索引
func (r *UserRepository) DropIndex(ctx Context.Context, name string) error {
    _, err := r.collection.Indexes().DropOne(ctx, name)
    if err != nil {
        return err
    }

    fmt.Printf("✅ 删除索引成功: %s\n", name)

    return nil
}
```

---

## 事务处理

### ACID事务

```go
package repository

import (
    "context"
    "fmt"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

// TransferPoints 积分转账（使用事务）
func TransferPoints(ctx Context.Context, client *mongo.Client, fromUserID, toUserID primitive.ObjectID, points int) error {
    // 开启会话
    session, err := client.StartSession()
    if err != nil {
        return err
    }
    defer session.EndSession(ctx)

    // 定义事务函数
    callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
        usersCollection := client.Database("mydb").Collection("users")

        // 1. 扣减发送方积分
        filter := bson.M{
            "_id": fromUserID,
            "points": bson.M{"$gte": points}, // 确保积分足够
        }
        update := bson.M{
            "$inc": bson.M{"points": -points},
        }

        result, err := usersCollection.UpdateOne(sessCtx, filter, update)
        if err != nil {
            return nil, err
        }

        if result.MatchedCount == 0 {
            return nil, fmt.Errorf("insufficient points or user not found")
        }

        // 2. 增加接收方积分
        filter = bson.M{"_id": toUserID}
        update = bson.M{
            "$inc": bson.M{"points": points},
        }

        result, err = usersCollection.UpdateOne(sessCtx, filter, update)
        if err != nil {
            return nil, err
        }

        if result.MatchedCount == 0 {
            return nil, fmt.Errorf("recipient user not found")
        }

        // 3. 记录转账历史
        transactionsCollection := client.Database("mydb").Collection("transactions")
        transaction := bson.M{
            "from_user_id": fromUserID,
            "to_user_id":   toUserID,
            "points":       points,
            "created_at":   time.Now(),
        }

        _, err = transactionsCollection.InsertOne(sessCtx, transaction)
        if err != nil {
            return nil, err
        }

        return nil, nil
    }

    // 执行事务
    _, err = session.WithTransaction(ctx, callback)
    if err != nil {
        return fmt.Errorf("transaction failed: %w", err)
    }

    fmt.Println("✅ 积分转账成功")

    return nil
}
```

---

## Change Streams

### 实时监听数据变化

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// WatchUsers 监听用户集合的变化
func WatchUsers(ctx Context.Context, collection *mongo.Collection) {
    // 创建Change Stream
    pipeline := mongo.Pipeline{}

    // 可选：过滤特定操作
    // pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{
    //     {Key: "operationType", Value: bson.M{"$in": []string{"insert", "update"}}},
    // }}})

    changeStream, err := collection.Watch(ctx, pipeline)
    if err != nil {
        log.Fatal(err)
    }
    defer changeStream.Close(ctx)

    fmt.Println("🔔 开始监听用户数据变化...")

    // 持续监听
    for changeStream.Next(ctx) {
        var changeEvent bson.M
        if err := changeStream.Decode(&changeEvent); err != nil {
            log.Printf("解码失败: %v", err)
            continue
        }

        operationType := changeEvent["operationType"]
        fmt.Printf("\n=== 检测到变化 ===\n")
        fmt.Printf("操作类型: %v\n", operationType)

        switch operationType {
        case "insert":
            doc := changeEvent["fullDocument"]
            fmt.Printf("新增文档: %v\n", doc)

        case "update":
            docID := changeEvent["documentKey"]
            updatedFields := changeEvent["updateDescription"].(bson.M)["updatedFields"]
            fmt.Printf("更新文档 ID: %v\n", docID)
            fmt.Printf("更新字段: %v\n", updatedFields)

        case "delete":
            docID := changeEvent["documentKey"]
            fmt.Printf("删除文档 ID: %v\n", docID)
        }
    }

    if err := changeStream.Err(); err != nil {
        log.Fatal(err)
    }
}

// 使用示例
func main() {
    client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(context.Background())

    collection := client.Database("mydb").Collection("users")

    // 在goroutine中监听
    go WatchUsers(context.Background(), collection)

    // 主程序继续运行
    time.Sleep(10 * time.Minute)
}
```

---

## GridFS文件存储

### 大文件存储

```go
package storage

import (
    "context"
    "fmt"
    "io"
    "os"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/gridfs"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// FileStorage GridFS文件存储
type FileStorage struct {
    bucket *gridfs.Bucket
}

// NewFileStorage 创建文件存储
func NewFileStorage(db *mongo.Database) (*FileStorage, error) {
    bucket, err := gridfs.NewBucket(db)
    if err != nil {
        return nil, err
    }

    return &FileStorage{bucket: bucket}, nil
}

// Upload 上传文件
func (fs *FileStorage) Upload(ctx Context.Context, filename string, file io.Reader) (primitive.ObjectID, error) {
    opts := options.GridFSUpload().SetMetadata(bson.M{
        "uploaded_at": time.Now(),
        "content_type": getContentType(filename),
    })

    fileID, err := fs.bucket.UploadFromStream(filename, file, opts)
    if err != nil {
        return primitive.NilObjectID, fmt.Errorf("upload failed: %w", err)
    }

    fmt.Printf("✅ 文件上传成功，ID: %s\n", fileID.Hex())

    return fileID.(primitive.ObjectID), nil
}

// Download 下载文件
func (fs *FileStorage) Download(ctx Context.Context, fileID primitive.ObjectID, dest io.Writer) error {
    _, err := fs.bucket.DownloadToStream(fileID, dest)
    if err != nil {
        return fmt.Errorf("download failed: %w", err)
    }

    fmt.Println("✅ 文件下载成功")

    return nil
}

// DownloadByName 根据文件名下载
func (fs *FileStorage) DownloadByName(ctx Context.Context, filename string, dest io.Writer) error {
    _, err := fs.bucket.DownloadToStreamByName(filename, dest)
    if err != nil {
        return fmt.Errorf("download failed: %w", err)
    }

    return nil
}

// Delete 删除文件
func (fs *FileStorage) Delete(ctx Context.Context, fileID primitive.ObjectID) error {
    err := fs.bucket.Delete(fileID)
    if err != nil {
        return fmt.Errorf("delete failed: %w", err)
    }

    fmt.Println("✅ 文件删除成功")

    return nil
}

// ListFiles 列出所有文件
func (fs *FileStorage) ListFiles(ctx Context.Context) ([]FileInfo, error) {
    cursor, err := fs.bucket.Find(bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var files []FileInfo
    for cursor.Next(ctx) {
        var file gridfs.File
        if err := cursor.Decode(&file); err != nil {
            return nil, err
        }

        files = append(files, FileInfo{
            ID:         file.ID.(primitive.ObjectID),
            Filename:   file.Name,
            Length:     file.Length,
            UploadDate: file.UploadDate,
        })
    }

    return files, nil
}

// FileInfo 文件信息
type FileInfo struct {
    ID         primitive.ObjectID
    Filename   string
    Length     int64
    UploadDate time.Time
}

// getContentType 获取内容类型
func getContentType(filename string) string {
    // 简化实现
    ext := filepath.Ext(filename)
    switch ext {
    case ".jpg", ".jpeg":
        return "image/jpeg"
    case ".png":
        return "image/png"
    case ".pdf":
        return "application/pdf"
    default:
        return "application/octet-stream"
    }
}

// 使用示例
func ExampleFileStorage() {
    client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
    db := client.Database("mydb")

    storage, _ := NewFileStorage(db)

    // 上传文件
    file, _ := os.Open("large-file.pdf")
    defer file.Close()

    fileID, _ := storage.Upload(context.Background(), "large-file.pdf", file)

    // 下载文件
    outFile, _ := os.Create("downloaded.pdf")
    defer outFile.Close()

    storage.Download(context.Background(), fileID, outFile)
}
```

---

## 性能优化

### 连接池优化

```go
clientOptions := options.Client().
    ApplyURI("mongodb://localhost:27017").
    SetMaxPoolSize(100).                  // 最大连接数
    SetMinPoolSize(10).                   // 最小连接数
    SetMaxConnIdleTime(30 * time.Second). // 最大空闲时间
    SetServerSelectionTimeout(5 * time.Second)
```

### 批量操作

```go
// BulkWrite 批量写入
func (r *UserRepository) BulkWrite(ctx Context.Context, operations []mongo.WriteModel) error {
    opts := options.BulkWrite().SetOrdered(false) // 无序执行，更快

    result, err := r.collection.BulkWrite(ctx, operations, opts)
    if err != nil {
        return err
    }

    fmt.Printf("✅ 批量操作成功: 插入%d, 更新%d, 删除%d\n",
        result.InsertedCount, result.ModifiedCount, result.DeletedCount)

    return nil
}

// 示例：批量更新和插入
func ExampleBulkWrite(repo *UserRepository) {
    operations := []mongo.WriteModel{
        mongo.NewInsertOneModel().SetDocument(bson.M{"username": "alice", "age": 25}),
        mongo.NewInsertOneModel().SetDocument(bson.M{"username": "bob", "age": 30}),
        mongo.NewUpdateOneModel().
            SetFilter(bson.M{"username": "charlie"}).
            SetUpdate(bson.M{"$set": bson.M{"age": 35}}),
        mongo.NewDeleteOneModel().SetFilter(bson.M{"username": "dave"}),
    }

    repo.BulkWrite(context.Background(), operations)
}
```

### 投影查询

```go
// FindUsernamesOnly 只查询用户名（减少网络传输）
func (r *UserRepository) FindUsernamesOnly(ctx Context.Context) ([]string, error) {
    opts := options.Find().SetProjection(bson.M{
        "username": 1,
        "_id":      0,
    })

    cursor, err := r.collection.Find(ctx, bson.M{}, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var results []struct {
        Username string `bson:"username"`
    }

    if err := cursor.All(ctx, &results); err != nil {
        return nil, err
    }

    usernames := make([]string, len(results))
    for i, result := range results {
        usernames[i] = result.Username
    }

    return usernames, nil
}
```

---

## 最佳实践

### 1. 错误处理

```go
if err == mongo.ErrNoDocuments {
    return nil, fmt.Errorf("user not found")
}
if mongo.IsDuplicateKeyError(err) {
    return fmt.Errorf("username already exists")
}
if mongo.IsTimeout(err) {
    return fmt.Errorf("operation timeout")
}
```

### 2. Context超时控制

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

user, err := repo.FindByID(ctx, userID)
```

### 3. 连接复用

```go
// ✅ 全局单例
var mongoClient *mongo.Client

func init() {
    client, _ := mongo.Connect(context.Background(), clientOptions)
    mongoClient = client
}

// ❌ 每次都创建新连接
func BadExample() {
    client, _ := mongo.Connect(context.Background(), clientOptions)
    defer client.Disconnect(context.Background())
    // ...
}
```

### 4. 索引策略

- 为常用查询字段创建索引
- 使用复合索引覆盖多字段查询
- 避免创建过多索引（影响写入性能）
- 定期分析慢查询日志

### 5. 数据建模

- 嵌入 vs 引用：相关数据嵌入，独立数据引用
- 避免文档过大（16MB限制）
- 使用数组存储一对多关系

---

## 总结

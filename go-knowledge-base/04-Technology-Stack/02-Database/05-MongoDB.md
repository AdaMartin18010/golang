# MongoDB 驱动

> **分类**: 开源技术堆栈

---

## mongo-go-driver

```go
import "go.mongodb.org/mongo-driver/mongo"
import "go.mongodb.org/mongo-driver/mongo/options"

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
defer client.Disconnect(ctx)

collection := client.Database("test").Collection("users")
```

---

## CRUD

### 插入

```go
user := bson.D{{"name", "Alice"}, {"age", 30}}
result, _ := collection.InsertOne(ctx, user)
fmt.Println(result.InsertedID)
```

### 查询

```go
var result bson.M
err := collection.FindOne(ctx, bson.D{{"name", "Alice"}}).Decode(&result)

cursor, _ := collection.Find(ctx, bson.D{{"age", bson.D{{"$gte", 18}}}})
var results []bson.M
cursor.All(ctx, &results)
```

### 更新

```go
filter := bson.D{{"name", "Alice"}}
update := bson.D{{"$set", bson.D{{"age", 31}}}}
collection.UpdateOne(ctx, filter, update)
```

---

## 索引

```go
indexModel := mongo.IndexModel{
    Keys:    bson.D{{"email", 1}},
    Options: options.Index().SetUnique(true),
}
collection.Indexes().CreateOne(ctx, indexModel)
```

# Ent Schema 定义示例

本目录包含 Ent Schema 定义的示例代码，展示如何定义数据模型。

## ⚠️ 重要说明

这些是**示例代码**，用于展示 Ent Schema 的定义方式。用户应该在自己的项目中定义自己的业务 Schema。

## 目录结构

```
ent-schema/
├── basic/         # 基础 Schema 示例
│   └── user.go    # User 实体示例
└── advanced/      # 高级 Schema 示例（待添加）
```

## 使用方式

### 1. 在自己的项目中定义 Schema

```go
// 用户项目中的 Ent Schema
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
)

type Product struct {
    ent.Schema
}

func (Product) Fields() []ent.Field {
    return []ent.Field{
        field.String("id").Unique().Immutable(),
        field.String("name").NotEmpty(),
        field.Float("price").Positive(),
        // ...
    }
}
```

### 2. 生成 Ent 代码

```bash
go generate ./ent
```

### 3. 使用生成的客户端

```go
client, err := ent.Open("postgres", "connection_string")
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// 使用生成的客户端
user, err := client.User.
    Create().
    SetEmail("user@example.com").
    SetName("John").
    Save(ctx)
```

## 相关资源

- [Ent 官方文档](https://entgo.io/)
- [Ent Schema 定义指南](https://entgo.io/docs/schema-def/)

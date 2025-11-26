# GraphQL

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [GraphQL](#graphql)
  - [📋 目录](#-目录)
  - [1. 📖 gqlgen入门](#1--gqlgen入门)
  - [🔍 DataLoader](#-dataloader)
  - [📊 分页](#-分页)
  - [💡 订阅 (Subscriptions)](#-订阅-subscriptions)
  - [📚 相关资源](#-相关资源)

---

## 1. 📖 gqlgen入门

```go
// schema.graphql
type Query {
    users: [User!]!
    user(id: ID!): User
}

type Mutation {
    createUser(input: NewUser!): User!
}

type User {
    id: ID!
    name: String!
    email: String!
}

input NewUser {
    name: String!
    email: String!
}
```

```go
// resolver.go
package graph

import (
    "context"
    "github.com/99designs/gqlgen/graphql/handler"
)

type Resolver struct {
    users map[string]*model.User
}

func (r *queryResolver) Users(ctx Context.Context) ([]*model.User, error) {
    users := make([]*model.User, 0, len(r.users))
    for _, user := range r.users {
        users = append(users, user)
    }
    return users, nil
}

func (r *queryResolver) User(ctx Context.Context, id string) (*model.User, error) {
    user, ok := r.users[id]
    if !ok {
        return nil, fmt.Errorf("user not found")
    }
    return user, nil
}

func (r *mutationResolver) CreateUser(ctx Context.Context, input model.NewUser) (*model.User, error) {
    user := &model.User{
        ID:    uuid.New().String(),
        Name:  input.Name,
        Email: input.Email,
    }

    r.users[user.ID] = user
    return user, nil
}

// 启动服务器
func main() {
    srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
        Resolvers: &Resolver{
            users: make(map[string]*model.User),
        },
    }))

    http.Handle("/graphql", srv)
    http.ListenAndServe(":8080", nil)
}
```

---

## 🔍 DataLoader

```go
import "github.com/graph-gophers/dataloader"

type UserLoader struct {
    db *sql.DB
}

func (u *UserLoader) BatchGetUsers(ctx Context.Context, keys dataloader.Keys) []*dataloader.Result {
    userIDs := make([]string, len(keys))
    for i, key := range keys {
        userIDs[i] = key.String()
    }

    // 批量查询
    users, err := u.db.GetUsersByIDs(userIDs)

    results := make([]*dataloader.Result, len(keys))
    for i, key := range keys {
        if err != nil {
            results[i] = &dataloader.Result{Error: err}
            continue
        }

        user := findUserByID(users, key.String())
        results[i] = &dataloader.Result{Data: user}
    }

    return results
}

// 在resolver中使用
func (r *queryResolver) Posts(ctx Context.Context, userID string) ([]*model.Post, error) {
    loader := ctx.Value("userLoader").(*dataloader.Loader)

    thunk := loader.Load(ctx, dataloader.StringKey(userID))
    result, err := thunk()
    if err != nil {
        return nil, err
    }

    user := result.(*model.User)
    return user.Posts, nil
}
```

---

## 📊 分页

```go
type Connection struct {
    Edges    []*Edge
    PageInfo *PageInfo
}

type Edge struct {
    Node   *model.User
    Cursor string
}

type PageInfo struct {
    HasNextPage     bool
    HasPreviousPage bool
    StartCursor     string
    EndCursor       string
}

func (r *queryResolver) UsersConnection(
    ctx Context.Context,
    first *int,
    after *string,
) (*Connection, error) {
    limit := 10
    if first != nil {
        limit = *first
    }

    var offset int
    if after != nil {
        offset = decodeCursor(*after)
    }

    users, err := r.db.GetUsers(limit+1, offset)
    if err != nil {
        return nil, err
    }

    hasNextPage := len(users) > limit
    if hasNextPage {
        users = users[:limit]
    }

    edges := make([]*Edge, len(users))
    for i, user := range users {
        edges[i] = &Edge{
            Node:   user,
            Cursor: encodeCursor(offset + i),
        }
    }

    return &Connection{
        Edges: edges,
        PageInfo: &PageInfo{
            HasNextPage: hasNextPage,
            StartCursor: edges[0].Cursor,
            EndCursor:   edges[len(edges)-1].Cursor,
        },
    }, nil
}
```

---

## 💡 订阅 (Subscriptions)

```go
type Subscription {
    messageAdded(roomID: ID!): Message!
}

func (r *subscriptionResolver) MessageAdded(
    ctx Context.Context,
    roomID string,
) (<-Channel *model.Message, error) {
    messages := make(Channel *model.Message, 1)

    // 订阅消息
    r.messageBus.Subscribe(roomID, func(msg *model.Message) {
        select {
        case messages <- msg:
        case <-ctx.Done():
        }
    })

    go func() {
        <-ctx.Done()
        r.messageBus.Unsubscribe(roomID)
        close(messages)
    }()

    return messages, nil
}
```

---

## 📚 相关资源

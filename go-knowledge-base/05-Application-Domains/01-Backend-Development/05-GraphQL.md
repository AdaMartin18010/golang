# GraphQL

> **分类**: 成熟应用领域

---

## gqlgen

```bash
go get github.com/99designs/gqlgen
```

---

## Schema

```graphql
type Query {
    users: [User!]!
    user(id: ID!): User
}

type User {
    id: ID!
    name: String!
    email: String!
}
```

---

## Resolver

```go
type Resolver struct{}

func (r *Resolver) Query() QueryResolver {
    return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Users(ctx context.Context) ([]*User, error) {
    return db.GetUsers()
}
```

---

## 生成代码

```bash
go run github.com/99designs/gqlgen generate
```

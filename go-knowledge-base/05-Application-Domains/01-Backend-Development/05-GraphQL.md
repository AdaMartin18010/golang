# GraphQL API Development in Go

> **Dimension**: Application Domains
> **Level**: S (18+ KB)
> **Tags**: #graphql #api #gqlgen #schema #resolver

---

## 1. Domain Requirements Analysis

### 1.1 GraphQL vs REST Comparison

| Aspect | GraphQL | REST |
|--------|---------|------|
| Data Fetching | Precise, single request | Multiple endpoints, over/under-fetching |
| Schema | Strong typing, introspection | OpenAPI/Swagger documentation |
| Versioning | Schema evolution | URL versioning (v1, v2) |
| Real-time | Native subscriptions | WebSockets, SSE separately |
| Caching | Application-level | HTTP caching, CDN-friendly |
| Complexity | Steeper learning curve | Simpler, widely understood |

### 1.2 When to Use GraphQL

**Ideal Scenarios:**

- Mobile applications with limited bandwidth
- Complex data relationships (nested resources)
- Rapidly evolving frontend requirements
- Aggregating multiple data sources
- Real-time data requirements

**Avoid When:**

- Simple CRUD operations suffice
- Heavy caching requirements at edge
- Team unfamiliar with GraphQL concepts
- File upload/download is primary use case

---

## 2. Architecture Formalization

### 2.1 GraphQL System Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         GraphQL API Architecture                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────┐     ┌──────────────┐     ┌─────────────────────────────┐  │
│  │   Client    │────►│   GraphQL    │────►│      Schema Layer           │  │
│  │             │     │   Server     │     │  ┌─────────────────────┐    │  │
│  │  Query      │     │  (gqlgen)    │     │  │ Type Definitions    │    │  │
│  │  Mutation   │     │              │     │  │ Directives          │    │  │
│  │  Subscription│    │              │     │  │ Extensions          │    │  │
│  └─────────────┘     └──────────────┘     │  └─────────────────────┘    │  │
│                              │            └─────────────────────────────┘  │
│                              ▼                              │               │
│                     ┌─────────────────┐                     │               │
│                     │   Middleware    │                     ▼               │
│                     │   Pipeline      │          ┌─────────────────────┐   │
│                     │                 │          │    Resolver Layer   │   │
│                     │  - Auth         │          │                     │   │
│                     │  - Logging      │          │  ┌───────────────┐  │   │
│                     │  - Metrics      │          │  │ Query Resolver│  │   │
│                     │  - Validation   │          │  ├───────────────┤  │   │
│                     └─────────────────┘          │  │ Mutation Resolver│ │   │
│                              │                   │  ├───────────────┤  │   │
│                              ▼                   │  │ Subscription  │  │   │
│                     ┌─────────────────┐          │  └───────────────┘  │   │
│                     │   Data Loaders  │          │                     │   │
│                     │   (N+1 Prevention)│        └─────────────────────┘   │
│                     └─────────────────┘                     │               │
│                              │                              ▼               │
│                              ▼                   ┌─────────────────────┐   │
│                     ┌─────────────────┐          │  Data Access Layer  │   │
│                     │   Cache Layer   │          │                     │   │
│                     │   (Redis)       │          │  - Repository       │   │
│                     └─────────────────┘          │  - Service          │   │
│                                                  │  - External APIs    │   │
│                                                  └─────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Schema-First Design

```graphql
# schema.graphql
scalar Time
scalar Upload

# Domain Types
type User {
    id: ID!
    email: String!
    name: String!
    avatar: String
    createdAt: Time!
    updatedAt: Time!

    # Relationships
    orders: [Order!]!
    posts: [Post!]!
}

type Order {
    id: ID!
    userId: ID!
    status: OrderStatus!
    total: Float!
    items: [OrderItem!]!
    createdAt: Time!

    # Computed field
    user: User!
}

type OrderItem {
    productId: ID!
    product: Product!
    quantity: Int!
    unitPrice: Float!
}

type Product {
    id: ID!
    name: String!
    description: String
    price: Float!
    stock: Int!
    category: Category!
}

type Category {
    id: ID!
    name: String!
    products: [Product!]!
}

type Post {
    id: ID!
    title: String!
    content: String!
    author: User!
    published: Boolean!
    createdAt: Time!
}

# Enums
enum OrderStatus {
    PENDING
    CONFIRMED
    SHIPPED
    DELIVERED
    CANCELLED
}

# Input Types
input CreateUserInput {
    email: String!
    name: String!
    password: String!
}

input UpdateUserInput {
    name: String
    avatar: String
}

input CreateOrderInput {
    items: [OrderItemInput!]!
    shippingAddress: AddressInput!
}

input OrderItemInput {
    productId: ID!
    quantity: Int!
}

input AddressInput {
    street: String!
    city: String!
    country: String!
    zipCode: String!
}

# Query Type
type Query {
    # User queries
    me: User!
    user(id: ID!): User
    users(pagination: PaginationInput): UserConnection!

    # Product queries
    product(id: ID!): Product
    products(
        filter: ProductFilter
        sort: ProductSort
        pagination: PaginationInput
    ): ProductConnection!

    # Order queries
    order(id: ID!): Order
    myOrders(pagination: PaginationInput): OrderConnection!
}

# Mutation Type
type Mutation {
    # User mutations
    createUser(input: CreateUserInput!): User!
    updateUser(id: ID!, input: UpdateUserInput!): User!
    deleteUser(id: ID!): Boolean!

    # Order mutations
    createOrder(input: CreateOrderInput!): Order!
    cancelOrder(id: ID!): Order!

    # File upload
    uploadAvatar(file: Upload!): String!
}

# Subscription Type
type Subscription {
    orderUpdated(userId: ID): Order!
    productStockUpdated(productId: ID): Product!
}

# Pagination
type PageInfo {
    hasNextPage: Boolean!
    hasPreviousPage: Boolean!
    startCursor: String
    endCursor: String
    totalCount: Int!
}

input PaginationInput {
    first: Int
    after: String
    last: Int
    before: String
}

type UserConnection {
    edges: [UserEdge!]!
    pageInfo: PageInfo!
}

type UserEdge {
    node: User!
    cursor: String!
}

type ProductConnection {
    edges: [ProductEdge!]!
    pageInfo: PageInfo!
}

type ProductEdge {
    node: Product!
    cursor: String!
}

type OrderConnection {
    edges: [OrderEdge!]!
    pageInfo: PageInfo!
}

type OrderEdge {
    node: Order!
    cursor: String!
}

# Filter inputs
input ProductFilter {
    categoryId: ID
    minPrice: Float
    maxPrice: Float
    inStock: Boolean
    search: String
}

enum ProductSortField {
    NAME
    PRICE
    CREATED_AT
}

input ProductSort {
    field: ProductSortField!
    direction: SortDirection!
}

enum SortDirection {
    ASC
    DESC
}
```

---

## 3. Scalability and Performance Considerations

### 3.1 Query Complexity Analysis

```go
package graphql

import (
    "context"
    "fmt"

    "github.com/99designs/gqlgen/graphql"
)

// ComplexityCalculator defines field complexity
type ComplexityCalculator struct {
    defaultComplexity int
    maxDepth          int
}

func NewComplexityCalculator() *ComplexityCalculator {
    return &ComplexityCalculator{
        defaultComplexity: 1,
        maxDepth:          10,
    }
}

// ComplexityConfig for gqlgen
func ComplexityConfig() graphql.ComplexityRoot {
    var c graphql.ComplexityRoot

    // User queries
    c.Query.User = func(childComplexity int, id string) int {
        return childComplexity + 1
    }

    c.Query.Users = func(childComplexity int, pagination *model.PaginationInput) int {
        // Limit based on requested page size
        limit := 20
        if pagination != nil && pagination.First != nil {
            limit = *pagination.First
        }
        return childComplexity * limit
    }

    // User relationships
    c.User.Orders = func(childComplexity int) int {
        return childComplexity * 10 // Assume avg 10 orders per user
    }

    c.User.Posts = func(childComplexity int) int {
        return childComplexity * 20 // Assume avg 20 posts per user
    }

    // Order relationships
    c.Order.Items = func(childComplexity int) int {
        return childComplexity * 5 // Assume avg 5 items per order
    }

    c.Order.User = func(childComplexity int) int {
        return childComplexity + 1
    }

    return c
}

// Max complexity limit
const MaxQueryComplexity = 1000
const MaxQueryDepth = 10
```

### 3.2 DataLoader Pattern (N+1 Prevention)

```go
package dataloader

import (
    "context"
    "net/http"
    "time"

    "github.com/graph-gophers/dataloader"
)

// Loaders holds all dataloaders
type Loaders struct {
    UserLoader     *dataloader.Loader
    OrderLoader    *dataloader.Loader
    ProductLoader  *dataloader.Loader
    CategoryLoader *dataloader.Loader
}

// Context key
type ctxKey string
const loadersKey = ctxKey("dataloaders")

// Middleware injects dataloaders into context
func Middleware(userRepo UserRepository, orderRepo OrderRepository) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            loaders := &Loaders{
                UserLoader:     newUserLoader(userRepo),
                OrderLoader:    newOrderLoader(orderRepo),
                ProductLoader:  newProductLoader(productRepo),
                CategoryLoader: newCategoryLoader(categoryRepo),
            }

            r = r.WithContext(context.WithValue(r.Context(), loadersKey, loaders))
            next.ServeHTTP(w, r)
        })
    }
}

// For returns loaders from context
func For(ctx context.Context) *Loaders {
    return ctx.Value(loadersKey).(*Loaders)
}

// User batch loader
func newUserLoader(repo UserRepository) *dataloader.Loader {
    batchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
        var results []*dataloader.Result

        // Fetch all users in one query
        userIDs := make([]string, len(keys))
        for i, key := range keys {
            userIDs[i] = key.String()
        }

        users, err := repo.FindByIDs(ctx, userIDs)
        if err != nil {
            // Return error for all
            for range keys {
                results = append(results, &dataloader.Result{Error: err})
            }
            return results
        }

        // Map users by ID
        userMap := make(map[string]*model.User)
        for _, user := range users {
            userMap[user.ID] = user
        }

        // Return results in same order as keys
        for _, key := range keys {
            if user, ok := userMap[key.String()]; ok {
                results = append(results, &dataloader.Result{Data: user})
            } else {
                results = append(results, &dataloader.Result{Data: nil})
            }
        }

        return results
    }

    return dataloader.NewBatchedLoader(batchFn,
        dataloader.WithWait(5*time.Millisecond),
        dataloader.WithBatchCapacity(100),
    )
}

// Order batch loader for user relationship
func newOrderLoader(repo OrderRepository) *dataloader.Loader {
    batchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
        var results []*dataloader.Result

        userIDs := make([]string, len(keys))
        for i, key := range keys {
            userIDs[i] = key.String()
        }

        // Fetch orders for all users
        ordersMap, err := repo.FindByUserIDs(ctx, userIDs)
        if err != nil {
            for range keys {
                results = append(results, &dataloader.Result{Error: err})
            }
            return results
        }

        for _, key := range keys {
            if orders, ok := ordersMap[key.String()]; ok {
                results = append(results, &dataloader.Result{Data: orders})
            } else {
                results = append(results, &dataloader.Result{Data: []*model.Order{}})
            }
        }

        return results
    }

    return dataloader.NewBatchedLoader(batchFn,
        dataloader.WithWait(5*time.Millisecond),
        dataloader.WithBatchCapacity(100),
    )
}
```

### 3.3 Query Response Caching

```go
package graphql

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

// CacheMiddleware implements response caching
type CacheMiddleware struct {
    redis      *redis.Client
    defaultTTL time.Duration
}

func NewCacheMiddleware(redis *redis.Client, ttl time.Duration) *CacheMiddleware {
    return &CacheMiddleware{
        redis:      redis,
        defaultTTL: ttl,
    }
}

// AroundResponses caches query responses
func (c *CacheMiddleware) AroundResponses(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
    // Skip caching for mutations
    if !isQuery(ctx) {
        return next(ctx)
    }

    // Generate cache key from query
    cacheKey := generateCacheKey(ctx)

    // Try to get from cache
    cached, err := c.redis.Get(ctx, cacheKey).Result()
    if err == nil {
        var response graphql.Response
        if json.Unmarshal([]byte(cached), &response) == nil {
            return &response
        }
    }

    // Execute query
    response := next(ctx)

    // Cache successful responses
    if len(response.Errors) == 0 {
        data, _ := json.Marshal(response)
        c.redis.Set(ctx, cacheKey, data, c.getTTL(ctx))
    }

    return response
}

func generateCacheKey(ctx context.Context) string {
    // Get operation context
    oc := graphql.GetOperationContext(ctx)

    // Hash query + variables
    data := oc.RawQuery
    if vars, err := json.Marshal(oc.Variables); err == nil {
        data += string(vars)
    }

    hash := sha256.Sum256([]byte(data))
    return fmt.Sprintf("graphql:%s", hex.EncodeToString(hash[:]))
}

func isQuery(ctx context.Context) bool {
    oc := graphql.GetOperationContext(ctx)
    if oc == nil || oc.Operation == nil {
        return false
    }
    return oc.Operation.Operation == ast.Query
}

func (c *CacheMiddleware) getTTL(ctx context.Context) time.Duration {
    // Check for @cacheControl directive
    // Return custom TTL or default
    return c.defaultTTL
}
```

---

## 4. Technology Stack Recommendations

### 4.1 Recommended Libraries

| Component | Library | Purpose |
|-----------|---------|---------|
| Schema Generation | gqlgen | Type-safe GraphQL server |
| DataLoader | graph-gophers/dataloader | Batching and caching |
| Subscriptions | gqlgen (built-in) | WebSocket subscriptions |
| Validation | go-playground/validator | Input validation |
| Auth | casbin | Authorization |
| Caching | redis/go-redis | Response caching |
| Tracing | OpenTelemetry | Distributed tracing |

### 4.2 gqlgen Configuration

```yaml
# gqlgen.yml
schema:
  - graph/*.graphql

exec:
  filename: graph/generated.go
  package: graph

model:
  filename: graph/model/models_gen.go
  package: model

resolver:
  layout: follow-schema
  dir: graph
  package: graph
  filename_template: "{name}.resolvers.go"

autobind:
  - "github.com/example/app/graph/model"

models:
  ID:
    model:
      - github.com/99designs/gqlgen/graphql.ID
      - github.com/99designs/gqlgen/graphql.Int
      - github.com/99designs/gqlgen/graphql.Int64
      - github.com/99designs/gqlgen/graphql.Int32
  Int:
    model:
      - github.com/99designs/gqlgen/graphql.Int
      - github.com/99designs/gqlgen/graphql.Int64
      - github.com/99designs/gqlgen/graphql.Int32
  Time:
    model: github.com/99designs/gqlgen/graphql.Time
  Upload:
    model: github.com/99designs/gqlgen/graphql.Upload
```

---

## 5. Case Studies

### 5.1 GitHub GraphQL API

**Scale:** 40M+ developers, billions of queries/month

**Architecture Decisions:**

- Rate limiting based on query complexity (node count)
- Query timeout: 10 seconds
- Max depth: 20 levels
- Complexity score per query tracked

**Lessons Learned:**

- Complexity calculation prevents abuse
- Cursor-based pagination essential
- Deprecation strategy for schema changes

### 5.2 Shopify Storefront API

**Scale:** 1M+ merchants, peak 100K requests/second

**Key Features:**

- GraphQL over HTTP and WebSocket
- Persisted queries for mobile apps
- CDN caching for public data

**Performance Optimizations:**

- Query result caching at edge
- DataLoader batching
- Field-level instrumentation

---

## 6. Go Implementation Examples

### 6.1 Complete Resolver Implementation

```go
package graph

import (
    "context"
    "fmt"

    "github.com/example/app/graph/model"
    "github.com/example/app/internal/dataloader"
    "github.com/example/app/internal/repository"
)

// Resolver is the root resolver
type Resolver struct {
    UserRepo    repository.UserRepository
    OrderRepo   repository.OrderRepository
    ProductRepo repository.ProductRepository
}

// Query resolver
func (r *Resolver) Query() QueryResolver {
    return &queryResolver{r}
}

// Mutation resolver
func (r *Resolver) Mutation() MutationResolver {
    return &mutationResolver{r}
}

// Subscription resolver
func (r *Resolver) Subscription() SubscriptionResolver {
    return &subscriptionResolver{r}
}

type queryResolver struct{ *Resolver }

// Me returns current authenticated user
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
    userID, err := auth.UserIDFromContext(ctx)
    if err != nil {
        return nil, fmt.Errorf("unauthenticated: %w", err)
    }

    return r.UserRepo.FindByID(ctx, userID)
}

// User returns a user by ID
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
    // Use dataloader to batch/cache
    return dataloader.For(ctx).UserLoader.Load(ctx, dataloader.StringKey(id))()
}

// Users returns paginated users
func (r *queryResolver) Users(ctx context.Context, pagination *model.PaginationInput) (*model.UserConnection, error) {
    limit := 20
    if pagination != nil && pagination.First != nil {
        limit = *pagination.First
    }

    users, hasMore, err := r.UserRepo.FindAll(ctx, limit, pagination.After)
    if err != nil {
        return nil, err
    }

    edges := make([]*model.UserEdge, len(users))
    for i, user := range users {
        edges[i] = &model.UserEdge{
            Node:   user,
            Cursor: encodeCursor(user.ID),
        }
    }

    return &model.UserConnection{
        Edges: edges,
        PageInfo: &model.PageInfo{
            HasNextPage: hasMore,
            EndCursor:   encodeCursor(users[len(users)-1].ID),
            TotalCount:  r.UserRepo.Count(ctx),
        },
    }, nil
}

// Products with filtering and sorting
func (r *queryResolver) Products(ctx context.Context, filter *model.ProductFilter, sort *model.ProductSort, pagination *model.PaginationInput) (*model.ProductConnection, error) {
    return r.ProductRepo.FindWithFilter(ctx, filter, sort, pagination)
}

type mutationResolver struct{ *Resolver }

// CreateUser creates a new user
func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
    // Validate input
    if err := validate.Struct(input); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    user := &model.User{
        ID:        generateID(),
        Email:     input.Email,
        Name:      input.Name,
        Password:  string(hashedPassword),
        CreatedAt: time.Now(),
    }

    if err := r.UserRepo.Create(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}

// CreateOrder creates a new order with inventory check
func (r *mutationResolver) CreateOrder(ctx context.Context, input model.CreateOrderInput) (*model.Order, error) {
    userID, err := auth.UserIDFromContext(ctx)
    if err != nil {
        return nil, err
    }

    // Validate products and check inventory
    var total float64
    items := make([]*model.OrderItem, len(input.Items))

    for i, itemInput := range input.Items {
        product, err := r.ProductRepo.FindByID(ctx, itemInput.ProductID)
        if err != nil {
            return nil, fmt.Errorf("product %s not found", itemInput.ProductID)
        }

        if product.Stock < itemInput.Quantity {
            return nil, fmt.Errorf("insufficient stock for product %s", product.Name)
        }

        items[i] = &model.OrderItem{
            ProductID: itemInput.ProductID,
            Product:   product,
            Quantity:  itemInput.Quantity,
            UnitPrice: product.Price,
        }
        total += product.Price * float64(itemInput.Quantity)
    }

    order := &model.Order{
        ID:     generateID(),
        UserID: userID,
        Status: model.OrderStatusPending,
        Total:  total,
        Items:  items,
    }

    // Transaction: create order and update inventory
    if err := r.OrderRepo.CreateWithInventoryUpdate(ctx, order); err != nil {
        return nil, err
    }

    // Publish event for subscriptions
    r.publishOrderEvent(ctx, order)

    return order, nil
}

type subscriptionResolver struct{ *Resolver }

// OrderUpdated subscribes to order updates
func (r *subscriptionResolver) OrderUpdated(ctx context.Context, userID *string) (<-chan *model.Order, error) {
    ch := make(chan *model.Order, 1)

    go func() {
        defer close(ch)

        // Subscribe to events
        events := r.EventBus.Subscribe(ctx, "order.updated")

        for event := range events {
            order := event.Payload.(*model.Order)

            // Filter by user if specified
            if userID != nil && order.UserID != *userID {
                continue
            }

            select {
            case ch <- order:
            case <-ctx.Done():
                return
            }
        }
    }()

    return ch, nil
}
```

### 6.2 Custom Directives

```go
package graph

import (
    "context"

    "github.com/99designs/gqlgen/graphql"
    "github.com/vektah/gqlparser/v2/ast"
)

// Auth directive - requires authentication
func Auth(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
    if _, err := auth.UserIDFromContext(ctx); err != nil {
        return nil, fmt.Errorf("unauthorized")
    }

    return next(ctx)
}

// HasRole directive - requires specific role
func HasRole(ctx context.Context, obj interface{}, next graphql.Resolver, roles []model.Role) (interface{}, error) {
    userRole, err := auth.RoleFromContext(ctx)
    if err != nil {
        return nil, err
    }

    for _, role := range roles {
        if role == userRole {
            return next(ctx)
        }
    }

    return nil, fmt.Errorf("forbidden: insufficient permissions")
}

// CacheControl directive - sets cache headers
func CacheControl(ctx context.Context, obj interface{}, next graphql.Resolver, maxAge int) (interface{}, error) {
    response := next(ctx)

    // Set cache control header
    if fc := graphql.GetFieldContext(ctx); fc != nil {
        fc.ResultContext.Stats.SetCacheControl(maxAge)
    }

    return response, nil
}

// Constraint directive - custom validation
func Constraint(ctx context.Context, obj interface{}, next graphql.Resolver, format string) (interface{}, error) {
    val, err := next(ctx)
    if err != nil {
        return nil, err
    }

    str, ok := val.(string)
    if !ok {
        return val, nil
    }

    switch format {
    case "email":
        if !isValidEmail(str) {
            return nil, fmt.Errorf("invalid email format")
        }
    case "uuid":
        if !isValidUUID(str) {
            return nil, fmt.Errorf("invalid UUID format")
        }
    }

    return val, nil
}
```

---

## 7. Visual Representations

### 7.1 GraphQL Query Execution Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      GraphQL Query Execution Flow                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Client Request                                                              │
│       │                                                                      │
│       ▼                                                                      │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     GraphQL Server                                  │   │
│  │  ┌───────────┐   ┌───────────┐   ┌───────────┐   ┌───────────┐    │   │
│  │  │  Parse    │──►│ Validate  │──►│  Execute  │──►│  Format   │    │   │
│  │  │           │   │           │   │           │   │           │    │   │
│  │  │ - Lexical │   │ - Schema  │   │ - Resolve │   │ - JSON    │    │   │
│  │  │ - Syntax  │   │ - Types   │   │ - Batch   │   │ - Errors  │    │   │
│  │  │ - AST     │   │ - Depth   │   │ - Cache   │   │           │    │   │
│  │  └───────────┘   │ - Complexity│  └───────────┘   └───────────┘    │   │
│  │                  └───────────┘                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│       │                                                                      │
│       ▼                                                                      │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Data Layer                                       │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐│   │
│  │  │ DataLoader  │  │   Cache     │  │  Database   │  │  External   ││   │
│  │  │ (Batch)     │  │   (Redis)   │  │             │  │    APIs     ││   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘│   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│       │                                                                      │
│       ▼                                                                      │
│  Response                                                                    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 7.2 Subscription Architecture with WebSockets

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    GraphQL Subscription Architecture                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Clients                           Server                          Events    │
│                                                                              │
│  ┌─────────┐                    ┌──────────────────┐              ┌────────┐ │
│  │ Web App │◄────WebSocket─────►│ Connection       │              │        │ │
│  └─────────┘      (GraphQL WS)  │ Manager          │              │ Redis  │ │
│                                  │                  │              │ Pub/Sub│ │
│  ┌─────────┐                    │ ┌──────────────┐ │◄────────────►│        │ │
│  │ Mobile  │◄────WebSocket─────►│ │ Subscription │ │   Subscribe  └────────┘ │
│  └─────────┘                    │ │ Handler      │ │                            │
│                                  │ └──────────────┘ │              ┌────────┐ │
│  ┌─────────┐                    │                  │              │        │ │
│  │ IoT     │◄────WebSocket─────►│ ┌──────────────┐ │              │ Kafka  │ │
│  └─────────┘                    │ │ Event Bus    │ │◄────────────►│        │ │
│                                  │ └──────────────┘ │   Publish    └────────┘ │
│                                  └──────────────────┘                            │
│                                                                              │
│  Flow:                                                                       │
│  1. Client subscribes via WebSocket                                          │
│  2. Server registers subscription                                            │
│  3. Events published to bus                                                  │
│  4. Matching subscriptions receive events                                    │
│  5. Server pushes to client                                                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 7.3 Schema Stitching and Federation

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    GraphQL Federation Architecture                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│                              Gateway                                         │
│                    ┌───────────────────────┐                                 │
│                    │    Apollo/GraphQL     │                                 │
│                    │       Gateway         │                                 │
│                    │                       │                                 │
│                    │  Query Planning       │                                 │
│                    │  Entity Resolution    │                                 │
│                    └───────────┬───────────┘                                 │
│                                │                                             │
│        ┌───────────────────────┼───────────────────────┐                     │
│        │                       │                       │                     │
│        ▼                       ▼                       ▼                     │
│  ┌─────────────┐       ┌─────────────┐       ┌─────────────┐                │
│  │   Users     │       │   Orders    │       │  Products   │                │
│  │   Service   │       │   Service   │       │   Service   │                │
│  │             │       │             │       │             │                │
│  │ type User { │       │ type Order {│       │ type Product│                │
│  │   id: ID!   │◄──────│   user: User│       │   id: ID!   │                │
│  │   name:     │       │   products: │◄──────│   name:     │                │
│  │ }           │       │   Product[] │       │ }           │                │
│  │             │       │             │       │             │                │
│  └─────────────┘       └─────────────┘       └─────────────┘                │
│                                                                              │
│  Key Concepts:                                                               │
│  - @key directive marks entity fields                                        │
│  - @extends for cross-service references                                     │
│  - @external for fields from other services                                  │
│  - Query planner optimizes cross-service calls                               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 8. Security Requirements

### 8.1 Security Checklist

| Category | Requirement | Implementation |
|----------|-------------|----------------|
| Query | Depth limiting | Max 10 levels |
| Query | Complexity limiting | Max 1000 points |
| Query | Timeout | 10 seconds max |
| Query | Persisted queries | For production |
| Auth | Authentication | JWT, OAuth2 |
| Auth | Authorization | RBAC, ABAC |
| Input | Validation | Schema + custom |
| Transport | HTTPS | TLS 1.3 |
| Transport | CORS | Restrict origins |
| Rate | Limiting | Per client/IP |

### 8.2 Security Middleware

```go
package security

import (
    "context"
    "time"

    "github.com/99designs/gqlgen/graphql"
)

// SecurityConfig holds security settings
type SecurityConfig struct {
    MaxDepth      int
    MaxComplexity int
    QueryTimeout  time.Duration
    Introspection bool
}

// DefaultSecurityConfig returns secure defaults
func DefaultSecurityConfig() *SecurityConfig {
    return &SecurityConfig{
        MaxDepth:      10,
        MaxComplexity: 1000,
        QueryTimeout:  10 * time.Second,
        Introspection: false,
    }
}

// SecurityMiddleware implements security checks
func SecurityMiddleware(config *SecurityConfig) graphql.HandlerExtension {
    return &securityExtension{config: config}
}

type securityExtension struct {
    config *SecurityConfig
}

func (s *securityExtension) ExtensionName() string {
    return "SecurityMiddleware"
}

func (s *securityExtension) Validate(schema graphql.ExecutableSchema) error {
    return nil
}

func (s *securityExtension) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
    // Check introspection in production
    if !s.config.Introspection {
        op := graphql.GetOperationContext(ctx)
        if op.Operation.Operation == ast.SchemaDefinition {
            return func(ctx context.Context) *graphql.Response {
                return graphql.ErrorResponse(ctx, "introspection disabled")
            }
        }
    }

    // Apply timeout
    ctx, cancel := context.WithTimeout(ctx, s.config.QueryTimeout)
    defer cancel()

    return next(ctx)
}

func (s *securityExtension) InterceptField(ctx context.Context, next graphql.Resolver) (interface{}, error) {
    // Check depth
    fc := graphql.GetFieldContext(ctx)
    if fc != nil && fc.Field.Depth > s.config.MaxDepth {
        return nil, fmt.Errorf("query exceeds maximum depth of %d", s.config.MaxDepth)
    }

    return next(ctx)
}
```

### 8.3 Persisted Queries

```go
package persistedqueries

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"

    "github.com/redis/go-redis/v9"
)

// PersistedQueryStore manages query persistence
type PersistedQueryStore struct {
    redis *redis.Client
}

// Save persists a query
func (s *PersistedQueryStore) Save(ctx context.Context, query string) (string, error) {
    hash := sha256.Sum256([]byte(query))
    key := hex.EncodeToString(hash[:])

    err := s.redis.Set(ctx, "pq:"+key, query, 0).Err()
    return key, err
}

// Get retrieves a persisted query
func (s *PersistedQueryStore) Get(ctx context.Context, hash string) (string, error) {
    return s.redis.Get(ctx, "pq:"+hash).Result()
}

// ApolloPersistedQueryHandler handles automatic persisted queries
func ApolloPersistedQueryHandler(store *PersistedQueryStore) graphql.HandlerExtension {
    return &apqExtension{store: store}
}

type apqExtension struct {
    store *PersistedQueryStore
}

func (a *apqExtension) ExtensionName() string {
    return "ApolloPersistedQueries"
}

func (a *apqExtension) Validate(schema graphql.ExecutableSchema) error {
    return nil
}

func (a *apqExtension) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
    // Check for persisted query extensions
    // Implementation details...
    return next(ctx)
}
```

---

**Quality Rating**: S (18+ KB)
**Last Updated**: 2026-04-02

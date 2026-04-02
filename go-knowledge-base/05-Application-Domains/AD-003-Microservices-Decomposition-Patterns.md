# AD-003: Microservices Decomposition Patterns

> **Dimension**: Application Domains
> **Level**: S (18+ KB)
> **Tags**: #microservices #decomposition #architecture

---

## 1. Formal Definition

### 1.1 Microservices Decomposition Problem

**Definition**: Given monolithic application M = <Components, Dependencies, Functions>, find partition P = {S1, S2, ..., Sn} minimizing inter-service coupling while maximizing intra-service cohesion.

### 1.2 Cohesion and Coupling Metrics

**Cohesion** = Internal Dependencies / Total Dependencies
**Coupling** = Cross-Service Calls / Total Calls

Optimal decomposition: Cohesion >= 0.7, Coupling <= 0.3

---

## 2. Decomposition Strategies

### 2.1 Decompose by Business Capability

Partition services based on organizational business capabilities.

**Example - E-commerce**:

- Product Catalog Service
- Order Management Service
- Payment Processing Service
- Inventory Management Service
- Shipping Service
- Customer Management Service

**Go Implementation**:

```go
package product

type Service struct {
    repo     Repository
    cache    Cache
    eventBus EventBus
}

type Product struct {
    ID          string
    SKU         string
    Name        string
    Price       Money
    Category    Category
}

func (s *Service) CreateProduct(ctx context.Context, cmd CreateProductCommand) (*Product, error) {
    product := &Product{
        ID:   generateID(),
        SKU:  cmd.SKU,
        Name: cmd.Name,
        Price: cmd.Price,
    }

    if err := s.repo.Save(ctx, product); err != nil {
        return nil, err
    }

    s.eventBus.Publish(ProductCreatedEvent{
        ProductID: product.ID,
        Timestamp: time.Now(),
    })

    return product, nil
}
```

### 2.2 Decompose by Subdomain (DDD)

Based on Domain-Driven Design subdomains:

- **Core Domain**: Order Service (complex business logic)
- **Supporting Subdomain**: Catalog Service (simplified)
- **Generic Subdomain**: Auth Service (buy/outsource)

```go
// Core Domain: Order Service
package order

type OrderService struct {
    pricingEngine   PricingEngine
    inventoryClient InventoryClient
    paymentClient   PaymentClient
}

func (s *OrderService) CreateOrder(ctx context.Context, cmd CreateOrderCommand) (*Order, error) {
    // Validate inventory
    if err := s.inventoryClient.CheckAvailability(ctx, cmd.Items); err != nil {
        return nil, err
    }

    // Calculate pricing
    pricing, err := s.pricingEngine.Calculate(ctx, cmd.Items)
    if err != nil {
        return nil, err
    }

    // Create order
    order := NewOrder(cmd.CustomerID, cmd.Items, pricing)

    return order, nil
}
```

### 2.3 Decompose by Transaction

Partition based on transaction consistency boundaries using Saga pattern.

```go
package saga

type SagaOrchestrator struct {
    steps []SagaStep
}

type SagaStep struct {
    Name       string
    Action     func(ctx context.Context) error
    Compensate func(ctx context.Context) error
}

func (o *SagaOrchestrator) Execute(ctx context.Context) error {
    completed := []int{}

    for i, step := range o.steps {
        if err := step.Action(ctx); err != nil {
            // Execute compensation
            for j := len(completed) - 1; j >= 0; j-- {
                o.steps[completed[j]].Compensate(ctx)
            }
            return err
        }
        completed = append(completed, i)
    }
    return nil
}
```

---

## 3. Implementation Strategies

### 3.1 Strangler Fig Pattern

Incremental migration from monolith to microservices.

```go
type Router struct {
    routes         map[string]http.Handler
    defaultHandler http.Handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    path := req.URL.Path

    for prefix, handler := range r.routes {
        if strings.HasPrefix(path, prefix) {
            handler.ServeHTTP(w, req)
            return
        }
    }

    r.defaultHandler.ServeHTTP(w, req)
}
```

### 3.2 Branch by Abstraction

```go
type PricingService interface {
    CalculatePrice(ctx context.Context, items []Item) (*Price, error)
}

// Legacy implementation
type LegacyPricingService struct{}

// New microservice implementation
type RemotePricingService struct {
    client *grpc.Client
}
```

---

## 4. Case Studies

### Netflix Evolution

- 1998-2008: DVD rental monolith
- 2009-2012: Split by business capability, AWS migration
- 2013-2018: 500+ services, service mesh
- 2019+: Merge small services, optimize

### Monzo Bank

- 1500+ microservices from day one
- Kubernetes-native
- One database per service
- Chaos engineering

---

## 5. Decision Framework

### When to Use Microservices

- Multiple teams need independent delivery
- Different scaling requirements
- Technology diversity needed
- Clear domain boundaries

### When NOT to Use

- Team size < 10
- No clear boundaries
- Infrastructure not ready
- Rapidly changing business

### Service Size Guidelines

- Lines of Code: 1K-10K
- Team Size: 5-9 people
- Business Capabilities: 1-3 per service

---

## 6. Anti-patterns

| Anti-pattern | Solution |
|--------------|----------|
| Distributed Monolith | Truly decouple interfaces |
| Over-splitting | Merge related services |
| Shared Database | Each service owns its data |
| Circular Dependency | Refactor or merge |
| Anemic Service | Encapsulate business logic |

---

## References

1. Building Microservices - Sam Newman
2. Microservices Patterns - Chris Richardson
3. Domain-Driven Design - Eric Evans
4. The Art of Scalability - Abbott & Fisher

---

**Quality Rating**: S (18+ KB)
**Last Updated**: 2026-04-02

# EC-024: Scatter-Gather Pattern

## Problem Formalization

### The Distributed Query Challenge

When a client request requires data or processing from multiple services, we need an efficient pattern to parallelize these requests and aggregate results without cascading latency.

#### Problem Statement

Given:

- Request R requiring data from services S = {s₁, s₂, ..., sₙ}
- Each service sᵢ has latency L(sᵢ) and availability P(sᵢ)
- Client timeout constraint T

Find an execution strategy that:

```
Minimize: TotalLatency = max(L(sᵢ) for successful sᵢ)
Maximize: Availability = 1 - Π(1 - P(sᵢ))
Subject to:
    - TotalLatency ≤ T
    - Partial results are acceptable if specified
    - Failed services don't cascade to overall failure
```

### Parallel Request Patterns

```
Sequential (Anti-pattern):
┌─────────────────────────────────────────────────────────────────────────┐
│                                                                         │
│  Client ──► S1 (100ms) ──► S2 (100ms) ──► S3 (100ms) ──► Response     │
│                                                                         │
│  Total: 300ms + overhead                                                │
│  Failure: Any service fails = total failure                             │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘

Scatter-Gather:
┌─────────────────────────────────────────────────────────────────────────┐
│                                                                         │
│                    ┌──► S1 (100ms) ──┐                                  │
│                    │                  │                                  │
│  Client ──Scatter─►├──► S2 (100ms) ──┼──Gather──► Response             │
│                    │                  │   (120ms total)                 │
│                    └──► S3 (100ms) ──┘                                  │
│                                                                         │
│  Total: ~100ms (parallel)                                               │
│  Failure: Can tolerate partial failures                                 │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Solution Architecture

### Core Components

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Scatter-Gather Architecture                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Scatter Phase                                                   │   │
│  │  ┌─────────────────────────────────────────────────────────────┐ │   │
│  │  │  1. Request Decomposition                                     │ │   │
│  │  │     • Split request into sub-requests per service           │ │   │
│  │  │     • Add correlation IDs for tracking                      │ │   │
│  │  │     • Set individual timeouts per service                   │ │   │
│  │  └─────────────────────────────────────────────────────────────┘ │   │
│  │                              │                                    │   │
│  │                              ▼                                    │   │
│  │  ┌─────────────────────────────────────────────────────────────┐ │   │
│  │  │  2. Concurrent Dispatch                                       │ │   │
│  │  │     • Spawn goroutines/workers for each service             │ │   │
│  │  │     • Apply backpressure if needed                          │ │   │
│  │  │     • Track in-flight requests                              │ │   │
│  │  └─────────────────────────────────────────────────────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                              │                                          │
│                              ▼                                          │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Gather Phase                                                    │   │
│  │  ┌─────────────────────────────────────────────────────────────┐ │   │
│  │  │  3. Result Collection                                         │ │   │
│  │  │     • Aggregate successful responses                        │ │   │
│  │  │     • Handle partial failures                               │ │   │
│  │  │     • Apply timeout/cancellation                            │ │   │
│  │  └─────────────────────────────────────────────────────────────┘ │   │
│  │                              │                                    │   │
│  │                              ▼                                    │   │
│  │  ┌─────────────────────────────────────────────────────────────┐ │   │
│  │  │  4. Response Aggregation                                      │ │   │
│  │  │     • Merge data according to strategy                      │ │   │
│  │  │     • Handle conflicts/duplicates                           │ │   │
│  │  │     • Format final response                                 │ │   │
│  │  └─────────────────────────────────────────────────────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### Aggregation Strategies

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Aggregation Strategies                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  1. Wait for All (Default)                                              │
│     Wait for all services to respond or timeout                         │
│     • Pros: Complete results                                            │
│     • Cons: Slowest service determines latency                          │
│     • Use: When completeness is critical                                │
│                                                                         │
│  2. First N (Quorum)                                                    │
│     Return after receiving N successful responses                       │
│     • Pros: Predictable latency                                         │
│     • Cons: May miss some data                                          │
│     • Use: When partial results are acceptable                          │
│                                                                         │
│  3. Timeout-Based                                                       │
│     Return all results collected within time window                     │
│     • Pros: Bounded latency                                             │
│     • Cons: Non-deterministic completeness                              │
│     • Use: Real-time systems with SLA requirements                      │
│                                                                         │
│  4. Priority-Based                                                      │
│     Critical services must respond, others are optional                 │
│     • Pros: Guaranteed core functionality                               │
│     • Cons: Complex priority management                                 │
│     • Use: Tiered service architectures                                 │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Production-Ready Go Implementation

### Core Scatter-Gather Implementation

```go
// pkg/scattergather/scattergather.go
package scattergather

import (
    "context"
    "fmt"
    "sync"
    "time"

    "golang.org/x/sync/errgroup"
)

// Request represents a request to a service
type Request struct {
    ID      string
    Service string
    Payload interface{}
}

// Response represents a response from a service
type Response struct {
    RequestID string
    Service   string
    Result    interface{}
    Error     error
    Latency   time.Duration
}

// ServiceCaller defines how to call a service
type ServiceCaller interface {
    Call(ctx context.Context, req Request) (interface{}, error)
}

// Strategy defines how to aggregate results
type Strategy int

const (
    // WaitForAll waits for all services to respond
    WaitForAll Strategy = iota
    // FirstN waits for N successful responses
    FirstN
    // Timeout returns results after timeout
    Timeout
    // Priority requires critical services, optional for others
    Priority
)

// Config configures scatter-gather behavior
type Config struct {
    Strategy        Strategy
    FirstN          int           // For FirstN strategy
    Timeout         time.Duration // Overall timeout
    ServiceTimeout  time.Duration // Per-service timeout
    ContinueOnError bool          // Whether to continue if one service fails
}

// ScatterGather orchestrates parallel service calls
type ScatterGather struct {
    config Config

    // Metrics
    metrics *Metrics
}

func New(config Config) *ScatterGather {
    return &ScatterGather{
        config:  config,
        metrics: NewMetrics(),
    }
}

// Execute performs scatter-gather across services
func (sg *ScatterGather) Execute(
    ctx context.Context,
    requests []Request,
    callers map[string]ServiceCaller,
) ([]Response, error) {
    start := time.Now()
    defer func() {
        sg.metrics.RecordTotalLatency(time.Since(start))
    }()

    // Create context with overall timeout
    ctx, cancel := context.WithTimeout(ctx, sg.config.Timeout)
    defer cancel()

    // Response collection
    responses := make(chan Response, len(requests))

    // Use errgroup for structured concurrency
    g, ctx := errgroup.WithContext(ctx)

    // Scatter phase
    for _, req := range requests {
        req := req // capture for goroutine
        caller, ok := callers[req.Service]
        if !ok {
            responses <- Response{
                RequestID: req.ID,
                Service:   req.Service,
                Error:     fmt.Errorf("no caller for service: %s", req.Service),
            }
            continue
        }

        g.Go(func() error {
            resp := sg.callService(ctx, req, caller)
            responses <- resp

            // Return error only if we shouldn't continue
            if resp.Error != nil && !sg.config.ContinueOnError {
                return resp.Error
            }
            return nil
        })
    }

    // Wait for completion in separate goroutine
    go func() {
        g.Wait()
        close(responses)
    }()

    // Gather phase with strategy
    return sg.gather(ctx, responses, len(requests))
}

func (sg *ScatterGather) callService(
    ctx context.Context,
    req Request,
    caller ServiceCaller,
) Response {
    start := time.Now()

    // Create per-service timeout context
    ctx, cancel := context.WithTimeout(ctx, sg.config.ServiceTimeout)
    defer cancel()

    result, err := caller.Call(ctx, req)

    return Response{
        RequestID: req.ID,
        Service:   req.Service,
        Result:    result,
        Error:     err,
        Latency:   time.Since(start),
    }
}

func (sg *ScatterGather) gather(
    ctx context.Context,
    responses <-chan Response,
    total int,
) ([]Response, error) {
    var results []Response

    switch sg.config.Strategy {
    case WaitForAll:
        results = sg.gatherAll(ctx, responses, total)
    case FirstN:
        results = sg.gatherFirstN(ctx, responses, sg.config.FirstN)
    case Timeout:
        results = sg.gatherTimeout(ctx, responses)
    case Priority:
        results = sg.gatherPriority(ctx, responses)
    }

    return results, nil
}

func (sg *ScatterGather) gatherAll(
    ctx context.Context,
    responses <-chan Response,
    total int,
) []Response {
    results := make([]Response, 0, total)
    received := 0

    for resp := range responses {
        results = append(results, resp)
        received++

        if resp.Error != nil {
            sg.metrics.RecordServiceError(resp.Service)
        } else {
            sg.metrics.RecordServiceLatency(resp.Service, resp.Latency)
        }

        if received >= total {
            break
        }
    }

    return results
}

func (sg *ScatterGather) gatherFirstN(
    ctx context.Context,
    responses <-chan Response,
    n int,
) []Response {
    results := make([]Response, 0, n)
    successful := 0

    for resp := range responses {
        if resp.Error == nil {
            results = append(results, resp)
            successful++

            if successful >= n {
                break
            }
        } else {
            sg.metrics.RecordServiceError(resp.Service)
        }
    }

    return results
}

func (sg *ScatterGather) gatherTimeout(
    ctx context.Context,
    responses <-chan Response,
) []Response {
    results := []Response{}

    for {
        select {
        case resp, ok := <-responses:
            if !ok {
                return results
            }
            results = append(results, resp)

        case <-ctx.Done():
            return results
        }
    }
}

// Aggregator merges responses into final result
type Aggregator interface {
    Aggregate(responses []Response) (interface{}, error)
}

// DefaultAggregator provides basic aggregation
type DefaultAggregator struct{}

func (a *DefaultAggregator) Aggregate(responses []Response) (interface{}, error) {
    successes := make([]Response, 0)
    failures := make([]Response, 0)

    for _, resp := range responses {
        if resp.Error == nil {
            successes = append(successes, resp)
        } else {
            failures = append(failures, resp)
        }
    }

    return &AggregationResult{
        Successes: successes,
        Failures:  failures,
        Success:   len(failures) == 0,
    }, nil
}

type AggregationResult struct {
    Successes []Response
    Failures  []Response
    Success   bool
}
```

### Product Search Example

```go
// internal/search/service.go
package search

import (
    "context"
    "fmt"
    "sort"
    "time"

    "github.com/company/project/pkg/scattergather"
)

// ProductSearchService searches across multiple product sources
type ProductSearchService struct {
    sg          *scattergather.ScatterGather
    aggregators map[string]Aggregator

    // Service callers
    catalogCaller   scattergather.ServiceCaller
    inventoryCaller scattergather.ServiceCaller
    pricingCaller   scattergather.ServiceCaller
    reviewsCaller   scattergather.ServiceCaller
}

// SearchRequest contains search parameters
type SearchRequest struct {
    Query      string
    Category   string
    MinPrice   float64
    MaxPrice   float64
    InStock    bool
    SortBy     string
    Page       int
    PageSize   int
}

// SearchResult contains aggregated results
type SearchResult struct {
    Products      []Product
    TotalCount    int
    Categories    []CategoryFacet
    PriceRange    PriceRange
    ResponseTime  time.Duration
    PartialResult bool
    Errors        []string
}

func (s *ProductSearchService) Search(ctx context.Context, req SearchRequest) (*SearchResult, error) {
    start := time.Now()

    // Build scatter-gather requests
    requests := []scattergather.Request{
        {
            ID:      "catalog-" + generateID(),
            Service: "catalog",
            Payload: CatalogSearchRequest{
                Query:    req.Query,
                Category: req.Category,
                Limit:    1000, // Get more for filtering
            },
        },
        {
            ID:      "inventory-" + generateID(),
            Service: "inventory",
            Payload: InventoryRequest{
                // Will be filled after catalog results
            },
        },
        {
            ID:      "pricing-" + generateID(),
            Service: "pricing",
            Payload: PricingRequest{
                // Will be filled after catalog results
            },
        },
        {
            ID:      "reviews-" + generateID(),
            Service: "reviews",
            Payload: ReviewsRequest{
                // Will be filled after catalog results
            },
        },
    }

    // Phase 1: Get catalog results first (prerequisite for other calls)
    catalogOnly := requests[:1]
    catalogCallers := map[string]scattergather.ServiceCaller{
        "catalog": s.catalogCaller,
    }

    catalogConfig := scattergather.Config{
        Strategy:        scattergather.WaitForAll,
        Timeout:         2 * time.Second,
        ServiceTimeout:  1 * time.Second,
        ContinueOnError: false,
    }

    sg := scattergather.New(catalogConfig)
    catalogResponses, err := sg.Execute(ctx, catalogOnly, catalogCallers)
    if err != nil {
        return nil, fmt.Errorf("catalog search failed: %w", err)
    }

    catalogResult := catalogResponses[0].Result.(*CatalogSearchResponse)

    // Phase 2: Parallel enrichment calls
    if len(catalogResult.Products) > 0 {
        productIDs := extractProductIDs(catalogResult.Products)

        enrichmentRequests := []scattergather.Request{
            {
                ID:      "inventory-" + generateID(),
                Service: "inventory",
                Payload: InventoryRequest{ProductIDs: productIDs},
            },
            {
                ID:      "pricing-" + generateID(),
                Service: "pricing",
                Payload: PricingRequest{ProductIDs: productIDs},
            },
            {
                ID:      "reviews-" + generateID(),
                Service: "reviews",
                Payload: ReviewsRequest{ProductIDs: productIDs},
            },
        }

        enrichmentCallers := map[string]scattergather.ServiceCaller{
            "inventory": s.inventoryCaller,
            "pricing":   s.pricingCaller,
            "reviews":   s.reviewsCaller,
        }

        enrichmentConfig := scattergather.Config{
            Strategy:        scattergather.FirstN,
            FirstN:          2, // Can proceed with partial data
            Timeout:         1 * time.Second,
            ServiceTimeout:  500 * time.Millisecond,
            ContinueOnError: true,
        }

        sg := scattergather.New(enrichmentConfig)
        enrichmentResponses, _ := sg.Execute(ctx, enrichmentRequests, enrichmentCallers)

        // Merge results
        return s.aggregateResults(catalogResult, enrichmentResponses, req, time.Since(start))
    }

    return &SearchResult{
        Products:     []Product{},
        TotalCount:   0,
        ResponseTime: time.Since(start),
    }, nil
}

func (s *ProductSearchService) aggregateResults(
    catalog *CatalogSearchResponse,
    enrichments []scattergather.Response,
    req SearchRequest,
    elapsed time.Duration,
) (*SearchResult, error) {

    // Build lookup maps from enrichment responses
    inventoryMap := make(map[string]InventoryInfo)
    pricingMap := make(map[string]PricingInfo)
    reviewsMap := make(map[string]ReviewsInfo)

    partialResult := false
    errors := []string{}

    for _, resp := range enrichments {
        if resp.Error != nil {
            partialResult = true
            errors = append(errors, fmt.Sprintf("%s: %v", resp.Service, resp.Error))
            continue
        }

        switch resp.Service {
        case "inventory":
            if inv, ok := resp.Result.(*InventoryResponse); ok {
                for _, item := range inv.Items {
                    inventoryMap[item.ProductID] = item
                }
            }
        case "pricing":
            if price, ok := resp.Result.(*PricingResponse); ok {
                for _, item := range price.Items {
                    pricingMap[item.ProductID] = item
                }
            }
        case "reviews":
            if rev, ok := resp.Result.(*ReviewsResponse); ok {
                for _, item := range rev.Items {
                    reviewsMap[item.ProductID] = item
                }
            }
        }
    }

    // Build final products with enrichment data
    products := make([]Product, 0, len(catalog.Products))

    for _, baseProduct := range catalog.Products {
        product := Product{
            ID:          baseProduct.ID,
            Name:        baseProduct.Name,
            Description: baseProduct.Description,
            Category:    baseProduct.Category,
        }

        // Add inventory data if available
        if inv, ok := inventoryMap[baseProduct.ID]; ok {
            product.InStock = inv.Quantity > 0
            product.Quantity = inv.Quantity
            product.Warehouse = inv.Warehouse
        }

        // Add pricing data if available
        if price, ok := pricingMap[baseProduct.ID]; ok {
            product.Price = price.Price
            product.Currency = price.Currency
            product.Discount = price.Discount
        }

        // Add reviews data if available
        if rev, ok := reviewsMap[baseProduct.ID]; ok {
            product.Rating = rev.AverageRating
            product.ReviewCount = rev.TotalReviews
        }

        // Apply filters
        if req.InStock && !product.InStock {
            continue
        }
        if product.Price > 0 {
            if req.MinPrice > 0 && product.Price < req.MinPrice {
                continue
            }
            if req.MaxPrice > 0 && product.Price > req.MaxPrice {
                continue
            }
        }

        products = append(products, product)
    }

    // Sort results
    products = sortProducts(products, req.SortBy)

    // Pagination
    totalCount := len(products)
    start := (req.Page - 1) * req.PageSize
    end := start + req.PageSize
    if start > totalCount {
        start = totalCount
    }
    if end > totalCount {
        end = totalCount
    }
    paginated := products[start:end]

    return &SearchResult{
        Products:      paginated,
        TotalCount:    totalCount,
        ResponseTime:  elapsed,
        PartialResult: partialResult,
        Errors:        errors,
    }, nil
}

func sortProducts(products []Product, sortBy string) []Product {
    switch sortBy {
    case "price_asc":
        sort.Slice(products, func(i, j int) bool {
            return products[i].Price < products[j].Price
        })
    case "price_desc":
        sort.Slice(products, func(i, j int) bool {
            return products[i].Price > products[j].Price
        })
    case "rating":
        sort.Slice(products, func(i, j int) bool {
            return products[i].Rating > products[j].Rating
        })
    case "newest":
        sort.Slice(products, func(i, j int) bool {
            return products[i].CreatedAt.After(products[j].CreatedAt)
        })
    default:
        // Default relevance (already sorted by catalog service)
    }
    return products
}
```

## Trade-off Analysis

### Strategy Comparison

| Strategy | Latency | Completeness | Complexity | Best For |
|----------|---------|--------------|------------|----------|
| WaitForAll | Highest | 100% | Low | Critical transactions |
| FirstN | Configurable | N/M | Medium | Quorum-based systems |
| Timeout | Bounded | Variable | Medium | Real-time systems |
| Priority | Configurable | Core+ | High | Tier-1/2/3 services |

### Resource Usage

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Resource Usage by Service Count                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Goroutines: O(N) where N = number of services                          │
│  Memory: O(N * avg_response_size)                                       │
│  Network: O(N * request_size + N * response_size)                       │
│                                                                         │
│  Guidelines:                                                            │
│  • Limit concurrent scatter operations per client                       │
│  • Use connection pooling to backend services                           │
│  • Implement request coalescing for identical queries                   │
│  • Consider circuit breakers for failing services                       │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Testing Strategies

### Scatter-Gather Testing

```go
// test/scattergather/scattergather_test.go
package scattergather

import (
    "context"
    "errors"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestScatterGatherWaitForAll(t *testing.T) {
    // Setup mock callers
    caller1 := new(MockServiceCaller)
    caller2 := new(MockServiceCaller)

    caller1.On("Call", mock.Anything, mock.Anything).Return("result1", nil)
    caller2.On("Call", mock.Anything, mock.Anything).Return("result2", nil)

    callers := map[string]ServiceCaller{
        "service1": caller1,
        "service2": caller2,
    }

    requests := []Request{
        {ID: "1", Service: "service1"},
        {ID: "2", Service: "service2"},
    }

    sg := New(Config{
        Strategy:        WaitForAll,
        Timeout:         5 * time.Second,
        ServiceTimeout:  1 * time.Second,
        ContinueOnError: false,
    })

    results, err := sg.Execute(context.Background(), requests, callers)

    assert.NoError(t, err)
    assert.Len(t, results, 2)
    assert.Equal(t, "result1", results[0].Result)
    assert.Equal(t, "result2", results[1].Result)
}

func TestScatterGatherFirstN(t *testing.T) {
    caller1 := new(MockServiceCaller)
    caller2 := new(MockServiceCaller)
    caller3 := new(MockServiceCaller)

    // Two fast, one slow
    caller1.On("Call", mock.Anything, mock.Anything).Return("result1", nil)
    caller2.On("Call", mock.Anything, mock.Anything).Return("result2", nil)
    caller3.On("Call", mock.Anything, mock.Anything).After(10*time.Second).Return("result3", nil)

    callers := map[string]ServiceCaller{
        "service1": caller1,
        "service2": caller2,
        "service3": caller3,
    }

    requests := []Request{
        {ID: "1", Service: "service1"},
        {ID: "2", Service: "service2"},
        {ID: "3", Service: "service3"},
    }

    sg := New(Config{
        Strategy:        FirstN,
        FirstN:          2,
        Timeout:         5 * time.Second,
        ServiceTimeout:  1 * time.Second,
    })

    start := time.Now()
    results, _ := sg.Execute(context.Background(), requests, callers)
    elapsed := time.Since(start)

    // Should complete quickly, not waiting for slow service
    assert.Less(t, elapsed, 2*time.Second)
    assert.Len(t, results, 2)
}

func TestScatterGatherWithFailures(t *testing.T) {
    caller1 := new(MockServiceCaller)
    caller2 := new(MockServiceCaller)

    caller1.On("Call", mock.Anything, mock.Anything).Return("result1", nil)
    caller2.On("Call", mock.Anything, mock.Anything).Return(nil, errors.New("service down"))

    callers := map[string]ServiceCaller{
        "service1": caller1,
        "service2": caller2,
    }

    requests := []Request{
        {ID: "1", Service: "service1"},
        {ID: "2", Service: "service2"},
    }

    // With ContinueOnError = true
    sg := New(Config{
        Strategy:        WaitForAll,
        Timeout:         5 * time.Second,
        ContinueOnError: true,
    })

    results, err := sg.Execute(context.Background(), requests, callers)

    assert.NoError(t, err) // No error because we continue on error
    assert.Len(t, results, 2)
    assert.NoError(t, results[0].Error)
    assert.Error(t, results[1].Error)

    // With ContinueOnError = false
    sg2 := New(Config{
        Strategy:        WaitForAll,
        Timeout:         5 * time.Second,
        ContinueOnError: false,
    })

    _, err = sg2.Execute(context.Background(), requests, callers)
    assert.Error(t, err) // Should return error
}

func TestScatterGatherTimeout(t *testing.T) {
    caller := new(MockServiceCaller)
    caller.On("Call", mock.Anything, mock.Anything).After(5 * time.Second).Return("result", nil)

    callers := map[string]ServiceCaller{
        "slow": caller,
    }

    requests := []Request{
        {ID: "1", Service: "slow"},
    }

    sg := New(Config{
        Strategy:       WaitForAll,
        Timeout:        500 * time.Millisecond,
        ServiceTimeout: 1 * time.Second,
    })

    start := time.Now()
    results, _ := sg.Execute(context.Background(), requests, callers)
    elapsed := time.Since(start)

    assert.Less(t, elapsed, 1*time.Second)
    assert.Len(t, results, 1)
    assert.Error(t, results[0].Error) // Should have timeout error
    assert.Contains(t, results[0].Error.Error(), "timeout")
}
```

## Summary

The Scatter-Gather Pattern provides:

1. **Parallel Execution**: Reduce total latency by executing requests concurrently
2. **Resilience**: Handle partial failures gracefully
3. **Flexibility**: Multiple aggregation strategies for different use cases
4. **Observability**: Track per-service latency and errors
5. **Scalability**: Add more services without increasing latency proportionally

Key considerations:

- Timeout management is critical
- Consider circuit breakers for failing services
- Design for partial results
- Monitor goroutine usage to prevent resource exhaustion

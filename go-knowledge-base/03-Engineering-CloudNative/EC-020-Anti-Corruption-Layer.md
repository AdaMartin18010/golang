# EC-020: Anti-Corruption Layer Pattern

## Problem Formalization

### The Legacy Integration Problem

When integrating with legacy systems or external services, modern applications face a fundamental challenge: how to consume functionality from systems with incompatible models, protocols, or data formats without contaminating the modern application's domain model.

#### Problem Statement

Given:

- Modern domain model M with bounded context boundaries
- Legacy/external system L with model Mₗ
- Integration requirements I = {i₁, i₂, ..., iₙ}

Find a translation layer T such that:

```
∀i ∈ I:
    - T.translate_to_legacy(i) produces valid L input
    - T.translate_from_legacy(L.output) produces valid M input
    - T maintains semantic equivalence: Meaning(M) = Meaning(L)

Constraints:
    - No dependencies from M to L's domain model
    - T isolates M from L's changes
    - T handles all protocol/format conversions
```

### Domain Model Contamination Risks

Without an ACL, developers may:

1. **Leaky Abstractions**: Use legacy data structures directly in modern code
2. **Semantic Coupling**: Depend on legacy system behavior and quirks
3. **Language Pollution**: Mix legacy terminology with modern ubiquitous language
4. **Technology Lock-in**: Expose legacy protocols to modern consumers

```go
// WITHOUT ACL - Contaminated domain model
// Modern service directly using legacy structures

type OrderService struct {
    legacyClient *legacy.SOAPClient  // Direct dependency!
}

func (s *OrderService) CreateOrder(ctx context.Context, req OrderRequest) (*Order, error) {
    // Converting modern request to legacy format inline
    legacyReq := &legacy.LegacyOrderRequest{
        CustomerId: strconv.Itoa(req.CustomerID), // Type conversion scattered
        ProductCode: req.ProductSKU,              // Different naming
        Qty: req.Quantity,                        // Abbreviation pollution
        // Missing fields handled inconsistently
    }

    // Direct call to legacy system
    legacyResp, err := s.legacyClient.CreateOrder(ctx, legacyReq)
    if err != nil {
        return nil, err
    }

    // Ad-hoc translation logic everywhere
    return &Order{
        ID: fmt.Sprintf("ORD-%d", legacyResp.OrderNumber),
        Status: mapLegacyStatus(legacyResp.Status), // Scattered mapping
    }, nil
}
```

## Solution Architecture

### Anti-Corruption Layer Structure

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Anti-Corruption Layer Architecture                   │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Modern Bounded Context (Clean Domain)                          │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │   │
│  │  │   Order      │  │   Customer   │  │   Payment            │  │   │
│  │  │   Aggregate  │  │   Aggregate  │  │   Aggregate          │  │   │
│  │  └──────────────┘  └──────────────┘  └──────────────────────┘  │   │
│  │                                                                  │   │
│  │  Domain Language: "Place Order", "Customer", "Payment Method"   │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                    │                                    │
│                                    │ Uses clean interface               │
│                                    ▼                                    │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Anti-Corruption Layer (ACL)                                    │   │
│  │  ┌─────────────────────────────────────────────────────────┐   │   │
│  │  │  Translation Layer                                      │   │   │
│  │  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │   │   │
│  │  │  │   Model      │  │   Error      │  │   Event      │  │   │   │
│  │  │  │   Translator │  │   Translator │  │   Translator │  │   │   │
│  │  │  └──────────────┘  └──────────────┘  └──────────────┘  │   │   │
│  │  └─────────────────────────────────────────────────────────┘   │   │
│  │                                                                  │   │
│  │  ┌─────────────────────────────────────────────────────────┐   │   │
│  │  │  Adapter Layer                                          │   │   │
│  │  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │   │   │
│  │  │  │   Service    │  │   Repository │  │   Client     │  │   │   │
│  │  │  │   Facade     │  │   Adapter    │  │   Adapter    │  │   │   │
│  │  │  └──────────────┘  └──────────────┘  └──────────────┘  │   │   │
│  │  └─────────────────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                    │                                    │
│                                    │ Protocol translation               │
│                                    ▼                                    │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Legacy/External System                                         │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │   │
│  │  │  Mainframe   │  │  SOAP API    │  │  External Partner    │  │   │
│  │  │  (CICS)      │  │  (WSDL)      │  │  (REST/GraphQL)      │  │   │
│  │  └──────────────┘  └──────────────┘  └──────────────────────┘  │   │
│  │                                                                  │   │
│  │  Legacy Language: "CUST-REC", "ORD-TRANS", "PYMT-TYPE"          │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### ACL Component Breakdown

```
┌─────────────────────────────────────────────────────────────────────────┐
│                     ACL Component Responsibilities                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Facade (Simplified Interface)                                  │   │
│  │  ┌─────────────────────────────────────────────────────────┐   │   │
│  │  │  Interface exposed to modern domain                      │   │   │
│  │  │  - Hides complexity of legacy system                     │   │   │
│  │  │  - Provides coarse-grained operations                    │   │   │
│  │  │  - Handles transactions if needed                        │   │   │
│  │  │                                                          │   │   │
│  │  │  Example:                                                │   │   │
│  │  │  RegisterCustomer() instead of:                          │   │   │
│  │  │  - CreateCustomerRecord()                                │   │   │
│  │  │  - UpdateCreditScore()                                   │   │   │
│  │  │  - LinkAccountToBranch()                                 │   │   │
│  │  └─────────────────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Translator (Model Conversion)                                  │   │
│  │  ┌─────────────────────────────────────────────────────────┐   │   │
│  │  │  Bidirectional mapping between domain models             │   │   │
│  │  │                                                          │   │   │
│  │  │  Modern Model              Legacy Model                  │   │   │
│  │  │  ───────────               ───────────                   │   │   │
│  │  │  Customer.Email    ────►   CUST-EMAIL (40 CHAR)          │   │   │
│  │  │  Order.Items       ────►   ORD-LINES (OCCURS 100)        │   │   │
│  │  │  Money.Amount      ────►   DECIMAL(15,2)                 │   │   │
│  │  │  UUID              ────►   COMP-3 PIC 9(18)               │   │   │
│  │  └─────────────────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Adapter (Protocol/Transport)                                   │   │
│  │  ┌─────────────────────────────────────────────────────────┐   │   │
│  │  │  Handles communication mechanics                         │   │   │
│  │  │                                                          │   │   │
│  │  │  - SOAP envelope construction                            │   │   │
│  │  │  - Flat file parsing                                     │   │   │
│  │  │  - Message queue protocols                               │   │   │
│  │  │  - Batch file processing                                 │   │   │
│  │  │  - Connection pooling                                    │   │   │
│  │  │  - Retry/circuit breaker logic                           │   │   │
│  │  └─────────────────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Production-Ready Go Implementation

### Complete ACL Implementation

```go
// internal/acl/legacyinventory/facade.go
package legacyinventory

import (
    "context"
    "fmt"
    "time"

    "github.com/company/project/internal/domain/inventory"
    "go.uber.org/zap"
)

// Facade is the only interface exposed to the domain layer
type Facade interface {
    // GetStockLevel returns current inventory for a SKU
    GetStockLevel(ctx context.Context, sku string) (*inventory.StockLevel, error)

    // ReserveStock attempts to reserve inventory for an order
    ReserveStock(ctx context.Context, req ReservationRequest) (*Reservation, error)

    // ReleaseReservation releases previously reserved stock
    ReleaseReservation(ctx context.Context, reservationID string) error

    // UpdateStock handles inbound shipments and adjustments
    UpdateStock(ctx context.Context, adjustment StockAdjustment) error

    // GetWarehouses returns all available warehouses
    GetWarehouses(ctx context.Context) ([]inventory.Warehouse, error)
}

// facade implements the Facade interface
type facade struct {
    client      LegacyClient
    translator  *ModelTranslator
    validator   *RequestValidator
    cache       *StockCache
    logger      *zap.Logger
    metrics     *Metrics
}

// ReservationRequest is the domain's request to reserve stock
type ReservationRequest struct {
    OrderID   string
    Items     []ReservedItem
    TTL       time.Duration
    Priority  ReservationPriority
}

type ReservedItem struct {
    SKU       string
    Quantity  int
    Warehouse *string // Optional: specific warehouse preference
}

type ReservationPriority int

const (
    PriorityLow ReservationPriority = iota
    PriorityNormal
    PriorityHigh
    PriorityCritical
)

// Reservation represents a successful stock reservation
type Reservation struct {
    ID              string
    OrderID         string
    Status          ReservationStatus
    ReservedItems   []ReservedItem
    ExpiresAt       time.Time
    LegacyReference string // Internal tracking
}

type ReservationStatus int

const (
    ReservationPending ReservationStatus = iota
    ReservationConfirmed
    ReservationPartial
    ReservationFailed
)

// StockAdjustment represents stock level changes
type StockAdjustment struct {
    Type        AdjustmentType
    SKU         string
    WarehouseID string
    Quantity    int
    Reason      string
    Reference   string // PO number, ticket ID, etc.
}

type AdjustmentType int

const (
    AdjustmentInbound AdjustmentType = iota
    AdjustmentOutbound
    AdjustmentCorrection
    AdjustmentDamage
)

// NewFacade creates a new ACL facade
func NewFacade(cfg Config, logger *zap.Logger) (Facade, error) {
    client, err := NewLegacyClient(cfg.LegacyEndpoint, cfg.Timeout)
    if err != nil {
        return nil, fmt.Errorf("creating legacy client: %w", err)
    }

    return &facade{
        client:     client,
        translator: NewModelTranslator(),
        validator:  NewRequestValidator(),
        cache:      NewStockCache(cfg.CacheTTL),
        logger:     logger,
        metrics:    NewMetrics(),
    }, nil
}

// GetStockLevel retrieves current stock with caching
func (f *facade) GetStockLevel(ctx context.Context, sku string) (*inventory.StockLevel, error) {
    start := time.Now()
    defer func() {
        f.metrics.RecordLatency("get_stock_level", time.Since(start))
    }()

    // Check cache first
    if cached := f.cache.Get(sku); cached != nil {
        f.metrics.RecordCacheHit("stock_level")
        return cached, nil
    }

    f.metrics.RecordCacheMiss("stock_level")

    // Build legacy request
    legacyReq := f.translator.ToLegacyStockRequest(sku)

    // Call legacy system
    legacyResp, err := f.client.InquiryStock(ctx, legacyReq)
    if err != nil {
        f.metrics.RecordError("get_stock_level")
        f.logger.Error("legacy stock inquiry failed",
            zap.Error(err),
            zap.String("sku", sku),
        )
        return nil, f.translator.TranslateError(err)
    }

    // Translate response
    stockLevel, err := f.translator.FromLegacyStockResponse(legacyResp)
    if err != nil {
        return nil, fmt.Errorf("translating response: %w", err)
    }

    // Cache result
    f.cache.Set(sku, stockLevel)

    return stockLevel, nil
}

// ReserveStock implements complex reservation logic with fallback
func (f *facade) ReserveStock(ctx context.Context, req ReservationRequest) (*Reservation, error) {
    // Validate request
    if err := f.validator.ValidateReservation(req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // Try primary reservation method
    legacyReq := f.translator.ToLegacyReservation(req)
    legacyResp, err := f.client.ReserveStock(ctx, legacyReq)

    if err != nil {
        // Check if error is retryable
        if f.translator.IsRetryableError(err) && req.Priority >= PriorityHigh {
            // Attempt manual reservation for high priority
            return f.manualReserve(ctx, req)
        }
        return nil, f.translator.TranslateError(err)
    }

    // Handle partial reservations
    if legacyResp.Status == LegacyStatusPartial {
        return f.handlePartialReservation(ctx, req, legacyResp)
    }

    // Translate and return
    reservation := f.translator.FromLegacyReservation(legacyResp)

    // Invalidate cache for affected SKUs
    for _, item := range req.Items {
        f.cache.Invalidate(item.SKU)
    }

    return reservation, nil
}

// manualReserve attempts reservation through alternative channels
func (f *facade) manualReserve(ctx context.Context, req ReservationRequest) (*Reservation, error) {
    // Some legacy systems require manual entry for priority orders
    // This method uses a different API or even file-based submission

    f.logger.Warn("attempting manual reservation",
        zap.String("order_id", req.OrderID),
        zap.Int("item_count", len(req.Items)),
    )

    manualReq := f.translator.ToManualReservationRequest(req)

    // Use async file-based submission for manual processing
    ticketID, err := f.client.SubmitManualReservation(ctx, manualReq)
    if err != nil {
        return nil, fmt.Errorf("manual reservation failed: %w", err)
    }

    // Return pending reservation
    return &Reservation{
        ID:              generateReservationID(),
        OrderID:         req.OrderID,
        Status:          ReservationPending,
        ReservedItems:   req.Items,
        ExpiresAt:       time.Now().Add(req.TTL),
        LegacyReference: ticketID,
    }, nil
}

// handlePartialReservation deals with partial stock availability
func (f *facade) handlePartialReservation(ctx context.Context,
    req ReservationRequest, legacyResp *LegacyReservationResponse) (*Reservation, error) {

    // Analyze what was reserved vs requested
    reserved := f.translator.ParseReservedItems(legacyResp)

    // Check if partial fulfillment is acceptable
    if req.Priority == PriorityCritical {
        // For critical orders, try to source from alternative warehouses
        alternative, err := f.sourceFromAlternative(ctx, req, reserved)
        if err == nil {
            // Merge reservations
            reserved = mergeReservations(reserved, alternative)
        }
    }

    return &Reservation{
        ID:            generateReservationID(),
        OrderID:       req.OrderID,
        Status:        ReservationPartial,
        ReservedItems: reserved,
        ExpiresAt:     time.Now().Add(req.TTL),
    }, nil
}

// sourceFromAlternative attempts to fulfill remaining from other warehouses
func (f *facade) sourceFromAlternative(ctx context.Context,
    req ReservationRequest, alreadyReserved []ReservedItem) ([]ReservedItem, error) {

    // Build map of fulfilled quantities
    fulfilled := make(map[string]int)
    for _, item := range alreadyReserved {
        fulfilled[item.SKU] += item.Quantity
    }

    // Find remaining requirements
    var alternatives []AlternativeSource
    for _, item := range req.Items {
        needed := item.Quantity - fulfilled[item.SKU]
        if needed <= 0 {
            continue
        }

        // Query alternative warehouses
        warehouses, err := f.findAlternativeWarehouses(ctx, item.SKU, needed)
        if err != nil {
            continue
        }

        alternatives = append(alternatives, AlternativeSource{
            SKU:         item.SKU,
            Needed:      needed,
            Warehouses:  warehouses,
        })
    }

    // Attempt to reserve from alternatives
    return f.reserveFromAlternatives(ctx, alternatives)
}
```

### Model Translator

```go
// internal/acl/legacyinventory/translator.go
package legacyinventory

import (
    "errors"
    "fmt"
    "strconv"
    "strings"
    "time"

    "github.com/company/project/internal/domain/inventory"
)

// ModelTranslator handles bidirectional model conversion
type ModelTranslator struct {
    skuMapper      *SKUMapper
    warehouseCodes map[string]string // modern ID -> legacy code
    errorMappings  map[string]error  // legacy error -> domain error
}

func NewModelTranslator() *ModelTranslator {
    return &ModelTranslator{
        skuMapper:      NewSKUMapper(),
        warehouseCodes: loadWarehouseCodeMappings(),
        errorMappings:  initErrorMappings(),
    }
}

// ToLegacyStockRequest converts domain query to legacy format
func (t *ModelTranslator) ToLegacyStockRequest(sku string) *LegacyStockRequest {
    // Legacy system uses different SKU format
    legacySKU := t.skuMapper.ToLegacy(sku)

    // Legacy uses packed decimal for SKU
    packedSKU := packDecimal(legacySKU)

    return &LegacyStockRequest{
        ItemNumber: packedSKU,
        // Legacy uses specific date format
        EffectiveDate: time.Now().Format("20060102"),
        // Always request all warehouse levels
        IncludeAllLocations: true,
    }
}

// FromLegacyStockResponse converts legacy response to domain model
func (t *ModelTranslator) FromLegacyStockResponse(resp *LegacyStockResponse) (*inventory.StockLevel, error) {
    // Unpack legacy SKU
    modernSKU := t.skuMapper.FromLegacy(unpackDecimal(resp.ItemNumber))

    // Parse warehouse levels
    warehouses := make([]inventory.WarehouseStock, 0, len(resp.LocationLevels))
    for _, loc := range resp.LocationLevels {
        whCode, ok := t.warehouseCodes[loc.LocationCode]
        if !ok {
            // Unknown warehouse - log but don't fail
            continue
        }

        warehouses = append(warehouses, inventory.WarehouseStock{
            WarehouseID: whCode,
            Available:   int(loc.AvailableQty),
            Reserved:    int(loc.ReservedQty),
            OnOrder:     int(loc.OnOrderQty),
        })
    }

    return &inventory.StockLevel{
        SKU:         modernSKU,
        Warehouses:  warehouses,
        LastUpdated: parseLegacyTimestamp(resp.LastUpdateTime),
        // Legacy system doesn't track at this level
        Version:     0,
    }, nil
}

// ToLegacyReservation converts domain reservation to legacy format
func (t *ModelTranslator) ToLegacyReservation(req ReservationRequest) *LegacyReservationRequest {
    legacyItems := make([]LegacyReservationItem, len(req.Items))

    for i, item := range req.Items {
        legacyItems[i] = LegacyReservationItem{
            ItemNumber:   packDecimal(t.skuMapper.ToLegacy(item.SKU)),
            Quantity:     item.Quantity,
            LocationCode: t.resolveWarehouseCode(item.Warehouse),
        }
    }

    // Legacy uses different order ID format
    legacyOrderID := formatLegacyOrderID(req.OrderID)

    // Map priority to legacy reservation type
    resType := t.mapPriorityToType(req.Priority)

    return &LegacyReservationRequest{
        OrderNumber:       legacyOrderID,
        ReservationItems:  legacyItems,
        ReservationType:   resType,
        // Convert duration to legacy format (hours)
        DurationHours:     int(req.TTL.Hours()),
    }
}

// FromLegacyReservation converts legacy reservation response
func (t *ModelTranslator) FromLegacyReservation(resp *LegacyReservationResponse) *Reservation {
    items := make([]ReservedItem, len(resp.ReservedItems))

    for i, item := range resp.ReservedItems {
        items[i] = ReservedItem{
            SKU:      t.skuMapper.FromLegacy(unpackDecimal(item.ItemNumber)),
            Quantity: item.QuantityReserved,
            Warehouse: t.reverseWarehouseCode(item.LocationCode),
        }
    }

    return &Reservation{
        ID:              resp.ConfirmationNumber,
        OrderID:         parseLegacyOrderID(resp.OrderNumber),
        Status:          t.mapLegacyStatus(resp.Status),
        ReservedItems:   items,
        ExpiresAt:       parseLegacyExpiry(resp.ExpirationDate),
        LegacyReference: resp.BatchReference,
    }
}

// TranslateError converts legacy errors to domain errors
func (t *ModelTranslator) TranslateError(err error) error {
    if err == nil {
        return nil
    }

    // Extract error code from legacy error message
    // Legacy format: "ERR-XXX: Description"
    code := extractErrorCode(err.Error())

    if domainErr, ok := t.errorMappings[code]; ok {
        return domainErr
    }

    // Unknown error - wrap with context
    return fmt.Errorf("legacy system error: %w", err)
}

// IsRetryableError determines if error warrants retry
func (t *ModelTranslator) IsRetryableError(err error) bool {
    if err == nil {
        return false
    }

    code := extractErrorCode(err.Error())

    // These legacy error codes indicate temporary issues
    retryableCodes := []string{"ERR-503", "ERR-504", "ERR-001"}

    for _, rc := range retryableCodes {
        if code == rc {
            return true
        }
    }

    return false
}

// Error mappings
type ErrorDefinitions struct {
    ErrInsufficientStock = errors.New("insufficient stock")
    ErrInvalidSKU        = errors.New("invalid SKU")
    ErrWarehouseOffline  = errors.New("warehouse offline")
    ErrReservationFailed = errors.New("reservation failed")
    ErrTimeout           = errors.New("legacy system timeout")
)

func initErrorMappings() map[string]error {
    return map[string]error{
        "ERR-100": ErrInsufficientStock,
        "ERR-200": ErrInvalidSKU,
        "ERR-300": ErrWarehouseOffline,
        "ERR-400": ErrReservationFailed,
        "ERR-504": ErrTimeout,
    }
}

// Helper functions

func packDecimal(s string) []byte {
    // Pack string into COBOL COMP-3 format
    // Each digit takes 4 bits, sign in last nibble
    // Implementation details...
    return []byte(s) // Simplified
}

func unpackDecimal(data []byte) string {
    // Unpack COMP-3 to string
    return string(data) // Simplified
}

func extractErrorCode(msg string) string {
    // Extract ERR-XXX pattern from message
    parts := strings.Split(msg, ":")
    if len(parts) > 0 {
        return strings.TrimSpace(parts[0])
    }
    return "UNKNOWN"
}

func parseLegacyTimestamp(ts string) time.Time {
    // Parse various legacy timestamp formats
    formats := []string{
        "20060102150405",
        "2006-01-02-15.04.05.000000",
    }

    for _, f := range formats {
        if t, err := time.Parse(f, ts); err == nil {
            return t
        }
    }

    return time.Time{}
}
```

### Legacy Client Adapter

```go
// internal/acl/legacyinventory/client.go
package legacyinventory

import (
    "bytes"
    "context"
    "encoding/xml"
    "fmt"
    "io"
    "net/http"
    "time"

    "github.com/company/project/pkg/retry"
)

// LegacyClient handles low-level communication with legacy system
type LegacyClient interface {
    InquiryStock(ctx context.Context, req *LegacyStockRequest) (*LegacyStockResponse, error)
    ReserveStock(ctx context.Context, req *LegacyReservationRequest) (*LegacyReservationResponse, error)
    ReleaseReservation(ctx context.Context, reservationID string) error
    SubmitManualReservation(ctx context.Context, req *ManualReservationRequest) (string, error)
}

// SOAPClient implements LegacyClient using SOAP protocol
type SOAPClient struct {
    endpoint    string
    httpClient  *http.Client
    retryConfig retry.Config
    auth        *LegacyAuth
}

type LegacyAuth struct {
    Username   string
    Password   string
    SystemID   string
    TerminalID string
}

func NewLegacyClient(endpoint string, timeout time.Duration) (LegacyClient, error) {
    return &SOAPClient{
        endpoint: endpoint,
        httpClient: &http.Client{
            Timeout: timeout,
            Transport: &http.Transport{
                MaxIdleConns:        100,
                MaxIdleConnsPerHost: 100,
                IdleConnTimeout:     90 * time.Second,
            },
        },
        retryConfig: retry.Config{
            MaxRetries:  3,
            BaseDelay:   100 * time.Millisecond,
            MaxDelay:    5 * time.Second,
            Multiplier:  2.0,
            RetryableErrors: []error{ErrTimeout, ErrConnectionReset},
        },
        auth: loadLegacyCredentials(),
    }, nil
}

// InquiryStock performs stock level inquiry
func (c *SOAPClient) InquiryStock(ctx context.Context, req *LegacyStockRequest) (*LegacyStockResponse, error) {
    // Build SOAP envelope
    envelope := c.buildInquiryEnvelope(req)

    // Execute with retry
    var resp LegacyStockResponse
    err := retry.Do(ctx, c.retryConfig, func() error {
        httpResp, err := c.executeSOAP(ctx, "InquiryStock", envelope)
        if err != nil {
            return err
        }
        defer httpResp.Body.Close()

        return c.parseInquiryResponse(httpResp.Body, &resp)
    })

    if err != nil {
        return nil, err
    }

    return &resp, nil
}

// ReserveStock performs stock reservation
func (c *SOAPClient) ReserveStock(ctx context.Context, req *LegacyReservationRequest) (*LegacyReservationResponse, error) {
    envelope := c.buildReservationEnvelope(req)

    var resp LegacyReservationResponse
    err := retry.Do(ctx, c.retryConfig, func() error {
        httpResp, err := c.executeSOAP(ctx, "ReserveStock", envelope)
        if err != nil {
            return err
        }
        defer httpResp.Body.Close()

        return c.parseReservationResponse(httpResp.Body, &resp)
    })

    if err != nil {
        return nil, err
    }

    return &resp, nil
}

func (c *SOAPClient) executeSOAP(ctx context.Context, operation string, envelope []byte) (*http.Response, error) {
    req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewReader(envelope))
    if err != nil {
        return nil, err
    }

    // SOAP headers
    req.Header.Set("Content-Type", "text/xml; charset=utf-8")
    req.Header.Set("SOAPAction", fmt.Sprintf("\"%s\"", operation))

    // Legacy system requires specific headers
    req.Header.Set("X-System-ID", c.auth.SystemID)
    req.Header.Set("X-Terminal-ID", c.auth.TerminalID)

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }

    // Check for SOAP faults
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        resp.Body.Close()
        return nil, fmt.Errorf("SOAP error: %s - %s", resp.Status, string(body))
    }

    return resp, nil
}

func (c *SOAPClient) buildInquiryEnvelope(req *LegacyStockRequest) []byte {
    // Build SOAP envelope with WS-Security
    envelope := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"
               xmlns:inv="http://legacy.example.com/inventory">
    <soap:Header>
        <wsse:Security xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd">
            <wsse:UsernameToken>
                <wsse:Username>%s</wsse:Username>
                <wsse:Password>%s</wsse:Password>
            </wsse:UsernameToken>
        </wsse:Security>
    </soap:Header>
    <soap:Body>
        <inv:InquiryStockRequest>
            <inv:ItemNumber>%s</inv:ItemNumber>
            <inv:EffectiveDate>%s</inv:EffectiveDate>
            <inv:IncludeAllLocations>%t</inv:IncludeAllLocations>
        </inv:InquiryStockRequest>
    </soap:Body>
</soap:Envelope>`, c.auth.Username, c.auth.Password,
        req.ItemNumber, req.EffectiveDate, req.IncludeAllLocations)

    return []byte(envelope)
}
```

## Trade-off Analysis

### ACL vs Direct Integration

| Aspect | With ACL | Direct Integration | Notes |
|--------|----------|-------------------|-------|
| **Domain Purity** | High | Low | ACL isolates legacy concepts |
| **Maintenance Cost** | Higher | Lower | Additional translation layer |
| **Change Flexibility** | High | Low | Easy to swap legacy systems |
| **Performance** | Lower | Higher | Translation overhead |
| **Testing** | Easier | Harder | Can mock legacy at ACL boundary |
| **Team Boundaries** | Clear | Blurred | ACL as team interface |

### Translation Granularity

```
┌─────────────────────────────────────────────────────────────────────────┐
│  Translation Granularity Options                                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  1. Field-level Translation                                             │
│     Pros: Maximum flexibility, fine-grained control                     │
│     Cons: High maintenance, complex mappers                             │
│                                                                         │
│  2. Entity-level Translation (Recommended)                              │
│     Pros: Balanced complexity, natural boundaries                       │
│     Cons: May need sub-translation for complex entities                 │
│                                                                         │
│  3. Use-case-level Translation                                          │
│     Pros: Optimized for specific flows, high performance                │
│     Cons: Duplication across use cases, harder to maintain              │
│                                                                         │
│  4. Aggregate-level Translation                                         │
│     Pros: Transactional consistency preserved                           │
│     Cons: Large translation surface, harder to test                     │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Testing Strategies

### ACL Testing Approach

```go
// test/acl/acl_test.go
package acl

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// MockLegacyClient for testing
type MockLegacyClient struct {
    mock.Mock
}

func (m *MockLegacyClient) InquiryStock(ctx context.Context, req *LegacyStockRequest) (*LegacyStockResponse, error) {
    args := m.Called(ctx, req)
    return args.Get(0).(*LegacyStockResponse), args.Error(1)
}

// Test translator bidirectional conversion
func TestTranslatorRoundTrip(t *testing.T) {
    translator := NewModelTranslator()

    original := &inventory.StockLevel{
        SKU: "PROD-12345",
        Warehouses: []inventory.WarehouseStock{
            {WarehouseID: "WH-EAST", Available: 100},
            {WarehouseID: "WH-WEST", Available: 50},
        },
    }

    // Convert to legacy and back
    legacyReq := translator.ToLegacyStockRequest(original.SKU)
    legacyResp := &LegacyStockResponse{
        ItemNumber: packDecimal("12345"),
        LocationLevels: []LocationLevel{
            {LocationCode: "01", AvailableQty: 100},
            {LocationCode: "02", AvailableQty: 50},
        },
    }

    converted, err := translator.FromLegacyStockResponse(legacyResp)
    assert.NoError(t, err)

    // Verify round-trip preserves semantics
    assert.Equal(t, original.SKU, converted.SKU)
    assert.Equal(t, len(original.Warehouses), len(converted.Warehouses))
}

// Test facade with mocked legacy
func TestFacadeWithMockedLegacy(t *testing.T) {
    mockClient := new(MockLegacyClient)
    facade := &facade{client: mockClient}

    // Setup expectations
    mockClient.On("InquiryStock", mock.Anything, mock.MatchedBy(func(req *LegacyStockRequest) bool {
        return req.ItemNumber != nil
    })).Return(&LegacyStockResponse{
        ItemNumber: packDecimal("12345"),
        AvailableQty: 100,
    }, nil)

    // Execute
    stock, err := facade.GetStockLevel(context.Background(), "PROD-12345")

    // Verify
    assert.NoError(t, err)
    assert.Equal(t, "PROD-12345", stock.SKU)
    mockClient.AssertExpectations(t)
}

// Test error translation
func TestErrorTranslation(t *testing.T) {
    translator := NewModelTranslator()

    testCases := []struct {
        legacyError    string
        expectedDomain error
    }{
        {"ERR-100: Insufficient stock", ErrInsufficientStock},
        {"ERR-200: Invalid item", ErrInvalidSKU},
        {"ERR-999: Unknown error", nil}, // Falls through
    }

    for _, tc := range testCases {
        t.Run(tc.legacyError, func(t *testing.T) {
            err := errors.New(tc.legacyError)
            translated := translator.TranslateError(err)

            if tc.expectedDomain != nil {
                assert.ErrorIs(t, translated, tc.expectedDomain)
            } else {
                assert.Contains(t, translated.Error(), "legacy system error")
            }
        })
    }
}
```

## Summary

The Anti-Corruption Layer pattern:

1. **Protects Domain Model**: Legacy concepts don't leak into modern system
2. **Enables Evolution**: Change legacy systems without affecting domain
3. **Clarifies Boundaries**: Clear interface between contexts
4. **Improves Testability**: Can test domain without legacy dependencies

Key implementation considerations:

- One ACL per bounded context/legacy system pair
- Keep translation logic in ACL, not domain
- Handle all legacy quirks internally
- Cache and optimize where possible
- Document legacy mappings thoroughly

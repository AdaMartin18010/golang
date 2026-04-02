# EC-048: Compensating Transaction Pattern

> **Dimension**: Engineering-CloudNative
> **Level**: S (>15KB)
> **Tags**: #compensating-transaction #saga-pattern #distributed-transactions #eventual-consistency #microservices
> **Authoritative Sources**:
>
> - [Principles of Transaction Processing](https://www.morganclaypool.com/doi/abs/10.2200/S00193ED1V01Y200810DTM002) - Bernstein & Newcomer (2009)
> - [Enterprise Integration Patterns](https://www.enterpriseintegrationpatterns.com/) - Hohpe & Woolf (2004)
> - [Designing Data-Intensive Applications](https://dataintensive.net/) - Kleppmann (2017)
> - [AWS Saga Pattern](https://docs.aws.amazon.com/prescriptive-guidance/latest/cloud-design-patterns/saga.html)

---

## 1. Problem Formalization

### 1.1 System Context and Constraints

**Definition 1.1 (Distributed Transaction Domain)**
Let $\mathcal{T} = \{T_1, T_2, ..., T_n\}$ be a distributed long-running transaction where:

- Each $T_i$ operates on local database $D_i$
- $T_i$ produces effects $E_i = effects(T_i)$
- No global atomic commit protocol exists (no 2PC)

**System Constraints:**

| Constraint | Formal Definition | Impact |
|------------|-------------------|--------|
| **No Global Lock** | $\nexists L_{global}: \forall T_i: acquire(L_{global})$ | Cannot enforce serializability globally |
| **Network Partition** | $\exists t: partition(D_i, D_j)$ | Sub-transactions may commit independently |
| **Partial Failure** | $\exists i: failure(T_i) \land \forall j \neq i: commit(T_j)$ | Need semantic undo capability |
| **Long Duration** | $\Delta(T) \gg \Delta(local\_transaction)$ | Locks cannot be held for duration |
| **Heterogeneous Data** | $\forall i, j: schema(D_i) \neq schema(D_j)$ | Compensations must be domain-aware |

### 1.2 Problem Statement

**Problem 1.1 (Compensating Transaction Problem)**
Given a saga $S = \{T_1, T_2, ..., T_m\}$ where each $T_i$ has committed, and a failure at step $T_k$, find a compensation sequence $C = \{C_{k-1}, C_{k-2}, ..., C_1\}$ such that:

$$\forall E_i \in \{E_1, ..., E_{k-1}\}: C_i(E_i) \Rightarrow \neg E_i \land consistent(D_{global})$$

**Semantic Requirements:**

1. **Compensatability**: $\forall T_i: \exists C_i: C_i \circ T_i = id$
2. **Reversibility**: $\forall E_i: compensatable(E_i) \Rightarrow deterministic(C_i)$
3. **Consistency**: $\forall D_i: apply(C) \Rightarrow \phi(D_i)$ where $\phi$ is the consistency predicate
4. **Idempotency**: $\forall C_i: C_i(C_i(x)) = C_i(x)$

### 1.3 Formal Requirements Specification

**Requirement 1.1 (Semantic Compensation)**
$$\forall T_i: compensatable(T_i) \Rightarrow semantically\_undoable(effects(T_i))$$

**Requirement 1.2 (Ordering Invariant)**
$$\forall i < j: commit(T_i) \land commit(T_j) \Rightarrow C_j \prec C_i \text{ (compensation order)}$$

**Requirement 1.3 (Durability)**
$$\forall C_i: issued(C_i) \Rightarrow \Diamond (committed(C_i) \lor failed(C_i))$$

---

## 2. Solution Architecture

### 2.1 Formal Compensating Transaction Model

**Definition 2.1 (Compensating Transaction)**
A Compensating Transaction $CT$ is a 5-tuple $\langle T, C, \prec, S, R \rangle$:

- $T$: Original transaction that produced effects
- $C: Effects \to Operations$: Compensation function mapping effects to undo operations
- $\prec$: Partial order on compensation execution
- $S$: State recorded at $T$ commit for compensation use
- $R$: Recovery information for incomplete compensations

**Compensation Function Types:**

| Type | Formula | Example |
|------|---------|---------|
| **Exact Inverse** | $C(T(x)) = T^{-1}(x)$ | Credit $\leftrightarrow$ Debit |
| **Semantic Undo** | $C(T(x)) = x'$ where $x' \approx x$ | Cancel order, issue refund |
| **Retrospective** | $C(T(x)) = notify \land record$ | Alert admins, log for manual handling |
| **Pivotal** | $\nexists C(T_{pivot})$: pivot marks point of no return | Cannot compensate past pivot |

### 2.2 Compensation State Machine

```
States:
- COMMITTED: T has committed, C not yet triggered
- COMPENSATING: C is executing
- COMPENSATED: C has committed successfully
- PARTIAL: C partially executed (some steps failed)
- IRREVERSIBLE: T beyond pivot, cannot compensate

Transitions:
COMMITTED ──trigger──► COMPENSATING ──success──► COMPENSATED
                         │
                         └──failure──► PARTIAL ──retry──► COMPENSATING
```

---

## 3. Visual Representations

### 3.1 Compensating Transaction Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    COMPENSATING TRANSACTION SYSTEM                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     COMPENSATION COORDINATOR                         │   │
│  │                                                                      │   │
│  │  ┌─────────────┐   ┌─────────────┐   ┌─────────────┐   ┌─────────┐  │   │
│  │  │   Saga      │──►│  Compensate │──►│   Execute   │──►│  Verify │  │   │
│  │  │   Monitor   │   │   Planner   │   │   Engine    │   │  State  │  │   │
│  │  └─────────────┘   └─────────────┘   └─────────────┘   └────┬────┘  │   │
│  │         │                      │              │              │       │   │
│  │         ▼                      ▼              ▼              ▼       │   │
│  │  ┌─────────────────────────────────────────────────────────────────┐ │   │
│  │  │                    STATE REPOSITORY                              │ │   │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │ │   │
│  │  │  │ Saga State  │  │ Compensate  │  │   Events    │              │ │   │
│  │  │  │   Store     │  │    Log      │  │    Store    │              │ │   │
│  │  │  └─────────────┘  └─────────────┘  └─────────────┘              │ │   │
│  │  └─────────────────────────────────────────────────────────────────┘ │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    │ Compensate commands                     │
│              ┌─────────────────────┼─────────────────────┐                  │
│              │                     │                     │                  │
│              ▼                     ▼                     ▼                  │
│  ┌───────────────────┐  ┌───────────────────┐  ┌───────────────────┐       │
│  │   Service A       │  │   Service B       │  │   Service C       │       │
│  │                   │  │                   │  │                   │       │
│  │ ┌───────────────┐ │  │ ┌───────────────┐ │  │ ┌───────────────┐ │       │
│  │ │ Compensatable │ │  │ │ Compensatable │ │  │ │ Compensatable │ │       │
│  │ │  Operation    │ │  │ │  Operation    │ │  │ │  Operation    │ │       │
│  │ │               │ │  │ │               │ │  │ │               │ │       │
│  │ │ • Record      │ │  │ │ • Record      │ │  │ │ • Record      │ │       │
│  │ │   state       │ │  │ │   state       │ │  │ │   state       │ │       │
│  │ │ • Implement   │ │  │ │ • Implement   │ │  │ │ • Implement   │ │       │
│  │ │   compensate()│ │  │ │   compensate()│ │  │ │   compensate()│ │       │
│  │ │ • Idempotent  │ │  │ │ • Idempotent  │ │  │ │ • Idempotent  │ │       │
│  │ └───────────────┘ │  │ └───────────────┘ │  │ └───────────────┘ │       │
│  │                   │  │                   │  │                   │       │
│  │  ┌─────────────┐  │  │  ┌─────────────┐  │  │  ┌─────────────┐  │       │
│  │  │  Database   │  │  │  │  Database   │  │  │  │  Database   │  │       │
│  │  │   (Local    │  │  │  │   (Local    │  │  │  │   (Local    │  │       │
│  │  │    State)   │  │  │  │    State)   │  │  │  │    State)   │  │       │
│  │  └─────────────┘  │  │  └─────────────┘  │  │  └─────────────┘  │       │
│  └───────────────────┘  └───────────────────┘  └───────────────────┘       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Compensation Execution Flow

```
Time ─────────────────────────────────────────────────────────────────────────►

Step 1: Reserve Inventory     Step 2: Process Payment      Step 3: Ship Order
  │                               │                              │
  │  T1: RESERVE_INV              │  T2: CHARGE_PAYMENT          │  T3: SHIP
  ├──────────────────────────────►├─────────────────────────────►│
  │                               │                              │
  │  C1: State=RESERVED           │  C2: State=CHARGED           │  C3: State=SHIPPED
  │      amount=100               │      amount=100              │      tracking=XYZ
  │                               │                              │
  │                               │                              │ X (FAILURE)
  │                               │                              │
  ▼                               ▼                              ▼
Compensation Phase:

         ┌─────────────────────────────────────────────────────────────┐
         │              COMPENSATION IN REVERSE ORDER                   │
         └─────────────────────────────────────────────────────────────┘

  C3: CANCEL_SHIP                                          ┌─────────────┐
  ────────────────────────────────────────────────────────►│  Service C  │
  • Void tracking number                                   │  (Success)  │
  • Return to warehouse                                    └──────┬──────┘
       │                                                          │
       │  C3 Success                                              │
       ◄──────────────────────────────────────────────────────────┘
       │
       ▼
  C2: REFUND_PAYMENT                                       ┌─────────────┐
  ────────────────────────────────────────────────────────►│  Service B  │
  • Reverse charge $100                                    │  (Success)  │
  • Record refund transaction                              └──────┬──────┘
       │                                                          │
       │  C2 Success                                              │
       ◄──────────────────────────────────────────────────────────┘
       │
       ▼
  C1: RELEASE_INVENTORY                                    ┌─────────────┐
  ────────────────────────────────────────────────────────►│  Service A  │
  • Remove reservation                                     │  (Success)  │
  • Make inventory available                               └──────┬──────┘
       │                                                          │
       │  C1 Success                                              │
       ◄──────────────────────────────────────────────────────────┘
       │
       ▼
  ┌─────────────┐
  │   SAGA      │
  │  ABORTED &  │
  │ CONSISTENT  │
  └─────────────┘
```

### 3.3 State Recovery After Failure

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    COMPENSATION STATE RECOVERY                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Scenario: Compensation partially executed before coordinator crash         │
│                                                                             │
│  Step 1: RESERVE     Step 2: CHARGE      Step 3: SHIP      (FAILED)        │
│     │                    │                   │                              │
│     │ C1                 │ C2                │ C3                           │
│     ▼                    ▼                   ▼                              │
│  ┌────────┐          ┌────────┐         ┌────────┐                         │
│  │RESERVED│          │CHARGED │         │SHIPPED │                         │
│  └────┬───┘          └────┬───┘         └────┬───┘                         │
│       │                   │                  │                              │
│       │ COMPENSATE        │ COMPENSATE      │ COMPENSATE                    │
│       ▼                   ▼                  ▼                              │
│  ┌──────────┐        ┌──────────┐      ┌──────────┐                        │
│  │C1 Success│        │C2 Success│      │C3 Failed │◄─── Coordinator crashes│
│  └──────────┘        └──────────┘      └──────────┘                        │
│       │                   │                  │                              │
│       │                   │                  │                              │
│       ▼                   ▼                  ▼                              │
│  [RELEASED]           [REFUNDED]         [SHIPPED] ◄── Inconsistent!        │
│                                                                             │
│  ═══════════════════════════════════════════════════════════════════════   │
│                                                                             │
│  Recovery Process:                                                          │
│                                                                             │
│  1. Coordinator restarts, reads compensation log                            │
│  2. Identifies pending C3 for saga S-12345                                  │
│  3. Retries C3 (idempotent):                                               │
│                                                                             │
│     ┌──────────┐                                                            │
│     │  RETRY   │────────────────────────────────────────► Service C          │
│     │   C3     │                                         • Check state      │
│     └──────────┘                                         • If SHIPPED:      │
│       │                                                    - Void tracking  │
│       │  Success                                           - Return stock   │
│       ◄────────────────────────────────────────────────────────────────    │
│       │                                                                     │
│       ▼                                                                     │
│  ┌─────────────┐                                                            │
│  │   SAGA      │                                                            │
│  │  ABORTED &  │                                                            │
│  │ CONSISTENT  │                                                            │
│  └─────────────┘                                                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Production Go Implementation

### 4.1 Core Compensation Framework

```go
package compensation

import (
 "context"
 "encoding/json"
 "errors"
 "fmt"
 "time"

 "go.opentelemetry.io/otel/attribute"
 "go.opentelemetry.io/otel/codes"
 "go.opentelemetry.io/otel/metric"
 "go.opentelemetry.io/otel/trace"
 "go.uber.org/zap"
)

// Compensatable defines the interface for compensatable operations
type Compensatable interface {
 // Execute performs the forward operation
 Execute(ctx context.Context) (*CompensationState, error)

 // Compensate undoes the operation using recorded state
 Compensate(ctx context.Context, state *CompensationState) error

 // IsCompensatable returns true if this operation can be compensated
 IsCompensatable() bool

 // GetCompensationType returns the type of compensation supported
 GetCompensationType() CompensationType
}

// CompensationType indicates how an operation can be compensated
type CompensationType int

const (
 CompensationTypeNone CompensationType = iota
 CompensationTypeExactInverse      // Mathematical inverse exists
 CompensationTypeSemanticUndo      // Business-level undo
 CompensationTypeRetrospective     // Record and notify only
 CompensationTypePivotal          // Point of no return
)

func (t CompensationType) String() string {
 switch t {
 case CompensationTypeNone:
  return "NONE"
 case CompensationTypeExactInverse:
  return "EXACT_INVERSE"
 case CompensationTypeSemanticUndo:
  return "SEMANTIC_UNDO"
 case CompensationTypeRetrospective:
  return "RETROSPECTIVE"
 case CompensationTypePivotal:
  return "PIVOTAL"
 default:
  return "UNKNOWN"
 }
}

// CompensationState records information needed for compensation
type CompensationState struct {
 OperationID   string                 `json:"operation_id"`
 OperationType string                 `json:"operation_type"`
 Input         map[string]interface{} `json:"input"`
 Output        map[string]interface{} `json:"output"`
 Timestamp     time.Time              `json:"timestamp"`
 Context       map[string]string      `json:"context"`

 // Compensation-specific data
 CompensationData json.RawMessage `json:"compensation_data,omitempty"`
}

// CompensatedOperation represents an operation with its compensation
type CompensatedOperation struct {
 ID          string
 Operation   Compensatable
 State       *CompensationState
 Status      OperationStatus
 ExecutedAt  time.Time
 CompensatedAt *time.Time
 Attempts    int
}

type OperationStatus int

const (
 OperationStatusPending OperationStatus = iota
 OperationStatusExecuting
 OperationStatusSucceeded
 OperationStatusFailed
 OperationStatusCompensating
 OperationStatusCompensated
 OperationStatusCompensationFailed
)

func (s OperationStatus) String() string {
 names := []string{
  "PENDING", "EXECUTING", "SUCCEEDED", "FAILED",
  "COMPENSATING", "COMPENSATED", "COMPENSATION_FAILED",
 }
 if int(s) < len(names) {
  return names[s]
 }
 return "UNKNOWN"
}

// CompensationCoordinator manages the compensation lifecycle
type CompensationCoordinator struct {
 store       CompensationStore
 logger      *zap.Logger
 tracer      trace.Tracer
 meter       metric.Meter

 // Metrics
 compensationsTotal   metric.Int64Counter
 compensationsFailed  metric.Int64Counter
 compensationDuration metric.Float64Histogram

 // Configuration
 maxRetries    int
 retryBackoff  time.Duration
}

// CompensationStore persists compensation state
type CompensationStore interface {
 SaveOperation(ctx context.Context, op *CompensatedOperation) error
 GetOperation(ctx context.Context, id string) (*CompensatedOperation, error)
 UpdateOperation(ctx context.Context, op *CompensatedOperation) error
 ListPendingCompensations(ctx context.Context) ([]*CompensatedOperation, error)
 SaveCompensationLog(ctx context.Context, log *CompensationLog) error
}

// CompensationLog records compensation attempts for audit
type CompensationLog struct {
 ID            string    `json:"id"`
 OperationID   string    `json:"operation_id"`
 Attempt       int       `json:"attempt"`
 Status        string    `json:"status"`
 Error         string    `json:"error,omitempty"`
 Timestamp     time.Time `json:"timestamp"`
 DurationMs    int64     `json:"duration_ms"`
}

// NewCompensationCoordinator creates a new coordinator
func NewCompensationCoordinator(
 store CompensationStore,
 logger *zap.Logger,
 tracer trace.Tracer,
 meter metric.Meter,
) (*CompensationCoordinator, error) {
 cc := &CompensationCoordinator{
  store:        store,
  logger:       logger,
  tracer:       tracer,
  meter:        meter,
  maxRetries:   3,
  retryBackoff: time.Second,
 }

 if meter != nil {
  var err error
  cc.compensationsTotal, err = meter.Int64Counter(
   "compensations_total",
   metric.WithDescription("Total number of compensation attempts"),
  )
  if err != nil {
   return nil, err
  }

  cc.compensationsFailed, err = meter.Int64Counter(
   "compensations_failed_total",
   metric.WithDescription("Total number of failed compensations"),
  )
  if err != nil {
   return nil, err
  }

  cc.compensationDuration, err = meter.Float64Histogram(
   "compensation_duration_seconds",
   metric.WithDescription("Duration of compensation execution"),
  )
  if err != nil {
   return nil, err
  }
 }

 return cc, nil
}

// ExecuteOperation runs an operation and records it for potential compensation
func (cc *CompensationCoordinator) ExecuteOperation(
 ctx context.Context,
 op Compensatable,
) (*CompensationState, error) {
 ctx, span := cc.tracer.Start(ctx, "compensation.ExecuteOperation",
  trace.WithAttributes(
   attribute.String("operation.type", op.GetCompensationType().String()),
  ))
 defer span.End()

 if !op.IsCompensatable() {
  span.SetAttributes(attribute.Bool("operation.compensatable", false))
  // Execute without compensation tracking
  return op.Execute(ctx)
 }

 // Create operation record
 compOp := &CompensatedOperation{
  ID:         cc.generateOperationID(),
  Operation:  op,
  Status:     OperationStatusExecuting,
  ExecutedAt: time.Now().UTC(),
 }

 // Execute the operation
 state, err := op.Execute(ctx)
 if err != nil {
  compOp.Status = OperationStatusFailed
  if saveErr := cc.store.SaveOperation(ctx, compOp); saveErr != nil {
   cc.logger.Error("Failed to save failed operation", zap.Error(saveErr))
  }
  span.RecordError(err)
  span.SetStatus(codes.Error, "operation execution failed")
  return nil, fmt.Errorf("operation execution failed: %w", err)
 }

 // Record successful execution
 compOp.State = state
 compOp.Status = OperationStatusSucceeded
 if err := cc.store.SaveOperation(ctx, compOp); err != nil {
  cc.logger.Error("Failed to save operation state", zap.Error(err))
  // Continue - compensation might still be possible
 }

 span.SetAttributes(
  attribute.String("operation.id", compOp.ID),
  attribute.String("operation.status", compOp.Status.String()),
 )

 return state, nil
}

// CompensateOperation executes compensation for a completed operation
func (cc *CompensationCoordinator) CompensateOperation(
 ctx context.Context,
 operationID string,
) error {
 ctx, span := cc.tracer.Start(ctx, "compensation.CompensateOperation",
  trace.WithAttributes(
   attribute.String("operation.id", operationID),
  ))
 defer span.End()

 start := time.Now()

 // Retrieve operation
 compOp, err := cc.store.GetOperation(ctx, operationID)
 if err != nil {
  span.RecordError(err)
  return fmt.Errorf("failed to retrieve operation: %w", err)
 }

 if compOp.Status != OperationStatusSucceeded {
  return fmt.Errorf("operation %s cannot be compensated: status is %s",
   operationID, compOp.Status.String())
 }

 // Update status
 compOp.Status = OperationStatusCompensating
 if err := cc.store.UpdateOperation(ctx, compOp); err != nil {
  cc.logger.Error("Failed to update operation status", zap.Error(err))
 }

 // Execute compensation with retries
 var lastErr error
 for attempt := 0; attempt <= cc.maxRetries; attempt++ {
  if attempt > 0 {
   time.Sleep(cc.retryBackoff * time.Duration(attempt))
  }

  err = compOp.Operation.Compensate(ctx, compOp.State)

  // Log attempt
  log := &CompensationLog{
   ID:          cc.generateLogID(),
   OperationID: operationID,
   Attempt:     attempt + 1,
   Timestamp:   time.Now().UTC(),
   DurationMs:  time.Since(start).Milliseconds(),
  }

  if err == nil {
   // Success
   now := time.Now().UTC()
   compOp.Status = OperationStatusCompensated
   compOp.CompensatedAt = &now
   compOp.Attempts = attempt + 1

   if updateErr := cc.store.UpdateOperation(ctx, compOp); updateErr != nil {
    cc.logger.Error("Failed to update operation after compensation",
     zap.Error(updateErr))
   }

   log.Status = "SUCCESS"
   if saveErr := cc.store.SaveCompensationLog(ctx, log); saveErr != nil {
    cc.logger.Error("Failed to save compensation log", zap.Error(saveErr))
   }

   duration := time.Since(start).Seconds()
   if cc.compensationsTotal != nil {
    cc.compensationsTotal.Add(ctx, 1, metric.WithAttributes(
     attribute.String("result", "success"),
    ))
   }
   if cc.compensationDuration != nil {
    cc.compensationDuration.Record(ctx, duration)
   }

   cc.logger.Info("Compensation succeeded",
    zap.String("operation_id", operationID),
    zap.Int("attempts", attempt+1),
    zap.Float64("duration_seconds", duration))

   return nil
  }

  lastErr = err
  log.Status = "FAILED"
  log.Error = err.Error()
  if saveErr := cc.store.SaveCompensationLog(ctx, log); saveErr != nil {
   cc.logger.Error("Failed to save compensation log", zap.Error(saveErr))
  }

  cc.logger.Warn("Compensation attempt failed",
   zap.String("operation_id", operationID),
   zap.Int("attempt", attempt+1),
   zap.Error(err))
 }

 // All retries exhausted
 compOp.Status = OperationStatusCompensationFailed
 compOp.Attempts = cc.maxRetries + 1
 if err := cc.store.UpdateOperation(ctx, compOp); err != nil {
  cc.logger.Error("Failed to update operation after compensation failure",
   zap.Error(err))
 }

 if cc.compensationsFailed != nil {
  cc.compensationsFailed.Add(ctx, 1, metric.WithAttributes(
   attribute.String("reason", "max_retries_exceeded"),
  ))
 }

 span.RecordError(lastErr)
 span.SetStatus(codes.Error, "compensation failed")

 return fmt.Errorf("compensation failed after %d attempts: %w", cc.maxRetries+1, lastErr)
}

// generateOperationID generates a unique operation ID
func (cc *CompensationCoordinator) generateOperationID() string {
 return fmt.Sprintf("op-%d-%d", time.Now().UnixNano(), time.Now().Nanosecond())
}

// generateLogID generates a unique log ID
func (cc *CompensationCoordinator) generateLogID() string {
 return fmt.Sprintf("log-%d-%d", time.Now().UnixNano(), time.Now().Nanosecond())
}
```

### 4.2 Example Compensatable Operations

```go
package compensation

import (
 "context"
 "encoding/json"
 "fmt"
)

// InventoryReservation implements compensatable inventory reservation
type InventoryReservation struct {
 ProductID string
 Quantity  int
 Warehouse string
}

func (ir *InventoryReservation) Execute(ctx context.Context) (*CompensationState, error) {
 // Call inventory service to reserve
 reservationID, err := ir.reserveInventory(ctx)
 if err != nil {
  return nil, err
 }

 // Record state for compensation
 state := &CompensationState{
  OperationType: "INVENTORY_RESERVE",
  Input: map[string]interface{}{
   "product_id": ir.ProductID,
   "quantity":   ir.Quantity,
   "warehouse":  ir.Warehouse,
  },
  Output: map[string]interface{}{
   "reservation_id": reservationID,
  },
  Timestamp: time.Now().UTC(),
 }

 // Store compensation-specific data
 compData, _ := json.Marshal(map[string]string{
  "reservation_id": reservationID,
 })
 state.CompensationData = compData

 return state, nil
}

func (ir *InventoryReservation) Compensate(ctx context.Context, state *CompensationState) error {
 // Parse compensation data
 var compData map[string]string
 if err := json.Unmarshal(state.CompensationData, &compData); err != nil {
  return fmt.Errorf("failed to parse compensation data: %w", err)
 }

 reservationID := compData["reservation_id"]
 if reservationID == "" {
  return fmt.Errorf("reservation_id not found in compensation data")
 }

 // Release the reservation
 return ir.releaseReservation(ctx, reservationID)
}

func (ir *InventoryReservation) IsCompensatable() bool {
 return true
}

func (ir *InventoryReservation) GetCompensationType() CompensationType {
 return CompensationTypeExactInverse
}

func (ir *InventoryReservation) reserveInventory(ctx context.Context) (string, error) {
 // Implementation: call inventory service
 return fmt.Sprintf("RES-%d", time.Now().Unix()), nil
}

func (ir *InventoryReservation) releaseReservation(ctx context.Context, reservationID string) error {
 // Implementation: call inventory service to release
 return nil
}

// PaymentCharge implements compensatable payment processing
type PaymentCharge struct {
 OrderID       string
 Amount        float64
 Currency      string
 PaymentMethod string
 CustomerID    string
}

func (pc *PaymentCharge) Execute(ctx context.Context) (*CompensationState, error) {
 // Process payment
 transactionID, err := pc.chargePayment(ctx)
 if err != nil {
  return nil, err
 }

 state := &CompensationState{
  OperationType: "PAYMENT_CHARGE",
  Input: map[string]interface{}{
   "order_id":       pc.OrderID,
   "amount":         pc.Amount,
   "currency":       pc.Currency,
   "payment_method": pc.PaymentMethod,
   "customer_id":    pc.CustomerID,
  },
  Output: map[string]interface{}{
   "transaction_id": transactionID,
  },
  Timestamp: time.Now().UTC(),
 }

 compData, _ := json.Marshal(map[string]interface{}{
  "transaction_id": transactionID,
  "amount":         pc.Amount,
  "currency":       pc.Currency,
 })
 state.CompensationData = compData

 return state, nil
}

func (pc *PaymentCharge) Compensate(ctx context.Context, state *CompensationState) error {
 var compData map[string]interface{}
 if err := json.Unmarshal(state.CompensationData, &compData); err != nil {
  return err
 }

 transactionID, ok := compData["transaction_id"].(string)
 if !ok {
  return fmt.Errorf("transaction_id not found")
 }

 // Issue refund
 return pc.refundPayment(ctx, transactionID)
}

func (pc *PaymentCharge) IsCompensatable() bool {
 return true
}

func (pc *PaymentCharge) GetCompensationType() CompensationType {
 return CompensationTypeSemanticUndo // Refund is semantic, not exact inverse
}

func (pc *PaymentCharge) chargePayment(ctx context.Context) (string, error) {
 return fmt.Sprintf("TXN-%d", time.Now().Unix()), nil
}

func (pc *PaymentCharge) refundPayment(ctx context.Context, transactionID string) error {
 return nil
}

// NotificationSend implements non-compensatable notification
type NotificationSend struct {
 Recipient string
 Message   string
 Channel   string
}

func (ns *NotificationSend) Execute(ctx context.Context) (*CompensationState, error) {
 // Send notification
 err := ns.sendNotification(ctx)
 if err != nil {
  return nil, err
 }

 return &CompensationState{
  OperationType: "NOTIFICATION_SEND",
  Input: map[string]interface{}{
   "recipient": ns.Recipient,
   "channel":   ns.Channel,
  },
  Timestamp: time.Now().UTC(),
 }, nil
}

func (ns *NotificationSend) Compensate(ctx context.Context, state *CompensationState) error {
 // Cannot unsend notification - no-op
 return nil
}

func (ns *NotificationSend) IsCompensatable() bool {
 return false
}

func (ns *NotificationSend) GetCompensationType() CompensationType {
 return CompensationTypeNone
}

func (ns *NotificationSend) sendNotification(ctx context.Context) error {
 return nil
}
```

---

## 5. Failure Scenarios and Mitigations

### 5.1 Failure Taxonomy

| Scenario | Impact | Detection | Mitigation |
|----------|--------|-----------|------------|
| **Compensation Service Unavailable** | Cannot undo | Connection timeout | Queue for retry + Alert |
| **Partial Compensation** | Inconsistent state | State verification | Idempotent retry + Reconciliation |
| **Irreversible Operation** | Cannot compensate | Type check before execute | Design-time validation |
| **State Corruption** | Wrong compensation | Checksum validation | Event sourcing + Replay |
| **Race Condition** | Double compensation | Unique constraint | Idempotency keys |
| **Timeout** | Unknown outcome | Context deadline | Idempotent query + Retry |

### 5.2 Recovery Mechanisms

```go
// RecoveryManager handles compensation failures
type RecoveryManager struct {
 coordinator *CompensationCoordinator
 store       CompensationStore
 logger      *zap.Logger
 alerter     Alerter
}

// RecoverIncompleteCompensations retries failed compensations
func (rm *RecoveryManager) RecoverIncompleteCompensations(ctx context.Context) error {
 pending, err := rm.store.ListPendingCompensations(ctx)
 if err != nil {
  return fmt.Errorf("failed to list pending compensations: %w", err)
 }

 for _, op := range pending {
  if op.Status == OperationStatusCompensationFailed {
   rm.logger.Warn("Recovering failed compensation",
    zap.String("operation_id", op.ID))

   if err := rm.coordinator.CompensateOperation(ctx, op.ID); err != nil {
    // Alert for manual intervention
    rm.alerter.Alert(ctx, Alert{
     Level:   AlertLevelCritical,
     Message: fmt.Sprintf("Compensation recovery failed for %s", op.ID),
     Error:   err.Error(),
    })
   }
  }
 }

 return nil
}
```

---

## 6. Semantic Trade-off Analysis

### 6.1 Compensation Strategy Comparison

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    COMPENSATION STRATEGY COMPARISON                          │
├─────────────────────┬─────────────────┬─────────────────┬───────────────────┤
│     Dimension       │  Exact Inverse  │ Semantic Undo   │  Retrospective    │
├─────────────────────┼─────────────────┼─────────────────┼───────────────────┤
│ Consistency Level   │ Strong          │ Eventual        │ Eventual + Manual │
│ Implementation      │ Automatic       │ Domain-specific │ Audit + Alert     │
│ Complexity          │ Low             │ Medium          │ High              │
│ Applicability       │ Limited         │ Broad           │ Last resort       │
│ Example             │ Financial       │ Cancel order    │ Physical shipment │
│                     │ transactions    │ Send apology    │ already delivered │
└─────────────────────┴─────────────────┴─────────────────┴───────────────────┘
```

### 6.2 Saga vs 2PC Trade-offs

| Aspect | Saga + Compensation | Two-Phase Commit |
|--------|---------------------|------------------|
| **Availability** | High (local commits) | Lower (blocking) |
| **Consistency** | Eventual | Strong |
| **Complexity** | Business logic | Infrastructure |
| **Performance** | High throughput | Lower throughput |
| **Recovery** | Compensation logic | Coordinator recovery |

---

## 7. References

1. Bernstein, P. A., & Newcomer, E. (2009). *Principles of Transaction Processing*. Morgan Kaufmann.
2. García-Molina, H., & Salem, K. (1987). Sagas. *ACM SIGMOD Record*.
3. Kleppmann, M. (2017). *Designing Data-Intensive Applications*. O'Reilly Media.
4. Hohpe, G., & Woolf, B. (2004). *Enterprise Integration Patterns*. Addison-Wesley.
5. Fowler, M. (2005). Patterns of Enterprise Application Architecture. Addison-Wesley.

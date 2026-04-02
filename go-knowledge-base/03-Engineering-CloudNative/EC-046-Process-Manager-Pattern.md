# EC-046: Process Manager Pattern (Saga Orchestrator)

> **Dimension**: Engineering-CloudNative
> **Level**: S (>15KB)
> **Tags**: #saga-pattern #process-manager #distributed-transactions #orchestration #event-driven #microservices
> **Authoritative Sources**:
>
> - [Enterprise Integration Patterns](https://www.enterpriseintegrationpatterns.com/) - Hohpe & Woolf (2004)
> - [Microservices Patterns](https://microservices.io/patterns/data/saga.html) - Richardson (2018)
> - [Saga Pattern for Microservices](https://docs.aws.amazon.com/prescriptive-guidance/latest/cloud-design-patterns/saga.html) - AWS (2024)
> - [Designing Data-Intensive Applications](https://dataintensive.net/) - Kleppmann (2017)

---

## 1. Problem Formalization

### 1.1 System Context and Constraints

**Definition 1.1 (Distributed Transaction Domain)**
Let $\mathcal{S} = \{S_1, S_2, ..., S_n\}$ be a set of distributed services, each maintaining local state $s_i \in \Sigma_i$. A distributed transaction $T$ spans multiple services requiring atomic consistency:

$$T = \{(S_i, op_i) \mid i \in [1, n], op_i \in \{invoke, compensate\}\}$$

**System Constraints:**

| Constraint | Formal Definition | Impact |
|------------|-------------------|--------|
| **Network Partition** | $\exists t: S_i \nleftrightarrow S_j$ | Services may become unreachable during execution |
| **Partial Failure** | $\exists S_k \in T: failure(S_k) \land \forall i \neq k: success(S_i)$ | Need compensation mechanism |
| **Eventual Consistency** | $\forall s_i: \Diamond \square (s_i \models \phi)$ | Temporary inconsistency is acceptable |
| **No Global Clock** | $\nexists C: \forall S_i, S_j: |clock_i - clock_j| < \epsilon$ | Ordering requires vector clocks |
| **Idempotency Required** | $\forall op: op(x) = op(op(x))$ | Retry safety for at-least-once delivery |

### 1.2 Problem Statement

**Problem 1.1 (Saga Coordination Problem)**
Given a long-running business process $P = \{step_1, step_2, ..., step_m\}$ where each $step_i$ involves service invocation with potential compensation $comp_i$, design a coordinator $PM$ such that:

$$\forall P: PM(P) \Rightarrow (\square success(P) \lor \Diamond consistent\_abort(P))$$

Where:

- $\square success(P)$: All steps complete successfully
- $\Diamond consistent\_abort(P)$: Completed steps are compensated to maintain consistency

**Key Challenges:**

1. **State Machine Complexity**: Managing $3^m$ possible state combinations for $m$ steps
2. **Failure Detection**: Distinguishing between service failure and communication failure
3. **Compensation Ordering**: Ensuring compensations execute in reverse order with proper dependencies
4. **Observability**: Tracking saga execution across distributed boundaries
5. **Scalability**: Handling thousands of concurrent saga instances

### 1.3 Formal Requirements Specification

**Requirement 1.1 (Atomicity)**
$$\forall P: (\forall i \in [1,m]: success(step_i)) \lor (\exists j: failure(step_j) \Rightarrow \forall i < j: executed(comp_i))$$

**Requirement 1.2 (Durability)**
$$\forall state_i \in PM: committed(state_i) \Rightarrow \square \neg lost(state_i)$$

**Requirement 1.3 (Liveness)**
$$\forall P: \neg deadlock(P) \land \neg livelock(P)$$

**Requirement 1.4 (Observability)**
$$\forall event_e \in saga: traced(e) \land recoverable(e)$$

---

## 2. Solution Architecture

### 2.1 Formal Process Manager Definition

**Definition 2.1 (Process Manager)**
A Process Manager $PM$ is a 7-tuple $\langle Q, \Sigma, \delta, q_0, F, \Gamma, \eta \rangle$:

- $Q = \{start, running_i, compensating_j, completed, failed\}$: State set
- $\Sigma = \{invoke, success, failure, compensate, timeout\}$: Input alphabet
- $\delta: Q \times \Sigma \to Q$: Transition function
- $q_0 = start$: Initial state
- $F = \{completed, failed\}$: Final states
- $\Gamma$: Saga definition (step sequence with compensations)
- $\eta: Q \to Action$: Output function mapping states to actions

### 2.2 State Transition Functions

**Primary Transitions:**

$$
\delta(running_i, success_i) = \begin{cases}
running_{i+1} & i < m \\
completed & i = m
\end{cases}
$$

$$
\delta(running_i, failure_i) = compensating_i
$$

$$
\delta(compensating_i, comp\_success_i) = \begin{cases}
compensating_{i-1} & i > 1 \\
failed & i = 1
\end{cases}
$$

**Output Function:**

$$
\eta(q) = \begin{cases}
invoke(step_{i+1}) & q = running_i, i < m \\
compensate(step_i) & q = compensating_i \\
nop & q \in \{completed, failed\}
\end{cases}
$$

### 2.3 Saga Definition Grammar

```
Saga        ::= Definition Steps Compensations
Definition  ::= "saga" Identifier "{" Meta "}"
Steps       ::= "steps" "{" Step+ "}"
Step        ::= Identifier "{"
                service: ServiceRef,
                action: ActionRef,
                input:  InputSchema,
                output: OutputSchema,
                timeout: Duration,
                retry: RetryPolicy
              "}"
Compensations ::= "compensations" "{" Comp+ "}"
Comp        ::= Identifier "compensates" Identifier
```

---

## 3. Visual Representations

### 3.1 Process Manager Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Client Request Gateway                               │
│                         (API Gateway / Load Balancer)                        │
└───────────────────────────────────┬─────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Process Manager (Orchestrator)                       │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                        Saga State Machine                              │  │
│  │                                                                         │  │
│  │   ┌─────────┐    invoke     ┌──────────┐    success    ┌──────────┐   │  │
│  │   │  START  │──────────────►│ RUNNING  │──────────────►│ RUNNING  │   │  │
│  │   │         │               │  step_1  │               │  step_2  │   │  │
│  │   └────┬────┘               └────┬─────┘               └────┬─────┘   │  │
│  │        │                         │    failure               │         │  │
│  │        │                         ▼                          │         │  │
│  │        │                    ┌──────────┐                   │         │  │
│  │        │                    │COMPENSATE│◄──────────────────┘         │  │
│  │        │                    │  step_1  │                             │  │
│  │        │                    └────┬─────┘                             │  │
│  │        │                         │ comp_done                         │  │
│  │        │                         ▼                                  │  │
│  │        │                    ┌──────────┐                            │  │
│  │        └───────────────────►│  FAILED  │                            │  │
│  │                             │          │                            │  │
│  │                             └──────────┘                            │  │
│  │                                                                       │  │
│  │   ═══════════════════════════════════════════════════════════════    │  │
│  │                         Persistence Layer                              │  │
│  │   • Saga Instance State                                                │  │
│  │   • Event Store (Event Sourcing)                                       │  │
│  │   • Compensation Log                                                   │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────────────┐  │
│  │   Command Bus    │  │   Event Bus      │  │   Observability Stack     │  │
│  │  (NATS/Kafka)    │  │  (Event Store)   │  │  (Traces/Metrics/Logs)   │  │
│  └────────┬─────────┘  └────────┬─────────┘  └───────────┬──────────────┘  │
└───────────┼─────────────────────┼────────────────────────┼─────────────────┘
            │                     │                        │
            ▼                     ▼                        ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Service Mesh / Message Bus                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │   Service A  │  │   Service B  │  │   Service C  │  │   Service D  │    │
│  │  (Inventory) │  │  (Payment)   │  │  (Shipping)  │  │ (Notification)│   │
│  │              │  │              │  │              │  │              │    │
│  │ ┌──────────┐ │  │ ┌──────────┐ │  │ ┌──────────┐ │  │ ┌──────────┐ │    │
│  │ │ Command  │ │  │ │ Command  │ │  │ │ Command  │ │  │ │ Command  │ │    │
│  │ │ Handler  │ │  │ │ Handler  │ │  │ │ Handler  │ │  │ │ Handler  │ │    │
│  │ └──────────┘ │  │ └──────────┘ │  │ └──────────┘ │  │ └──────────┘ │    │
│  │              │  │              │  │              │  │              │    │
│  │ ┌──────────┐ │  │ ┌──────────┐ │  │ ┌──────────┐ │  │ ┌──────────┐ │    │
│  │ │Compensate│ │  │ │Compensate│ │  │ │Compensate│ │  │ │Compensate│ │    │
│  │ │ Handler  │ │  │ │ Handler  │ │  │ │ Handler  │ │  │ │ Handler  │ │    │
│  │ └──────────┘ │  │ └──────────┘ │  │ └──────────┘ │  │ └──────────┘ │    │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘    │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Saga Execution Flow

```
Time ─────────────────────────────────────────────────────────────────────────►

Client          Process Manager          Service A          Service B          Service C
  │                   │                      │                  │                  │
  │  Start Saga       │                      │                  │                  │
  ├──────────────────►│                      │                  │                  │
  │                   │  Execute Step 1      │                  │                  │
  │                   ├─────────────────────►│                  │                  │
  │                   │                      │                  │                  │
  │                   │  Step 1 Success      │                  │                  │
  │                   │◄─────────────────────┤                  │                  │
  │                   │                      │                  │                  │
  │                   │  Execute Step 2                         │                  │
  │                   ├────────────────────────────────────────►│                  │
  │                   │                                         │                  │
  │                   │  Step 2 Success                         │                  │
  │                   │◄────────────────────────────────────────┤                  │
  │                   │                                         │                  │
  │                   │  Execute Step 3                                            │
  │                   ├───────────────────────────────────────────────────────────►│
  │                   │                                                            │
  │                   │  Step 3 FAILURE ◄─────────────────────────────────────────┤
  │                   │                                                            │
  │                   │════════════ COMPENSATION PHASE ════════════════════════════│
  │                   │                                                            │
  │                   │  Compensate Step 2                                         │
  │                   ├────────────────────────────────────────►│                  │
  │                   │                                         │                  │
  │                   │  Compensate Step 1                      │                  │
  │                   ├─────────────────────►│                  │                  │
  │                   │                      │                  │                  │
  │  Saga Failed      │                      │                  │                  │
  │◄──────────────────┤                      │                  │                  │
  │                   │                      │                  │                  │
```

### 3.3 State Machine Visualization

```
                         ╔═══════════════════════════════════════════════╗
                         ║           SAGA STATE MACHINE                   ║
                         ╚═══════════════════════════════════════════════╝

    ┌──────────┐
    │          │
    │  START   │
    │  (init)  │
    │          │
    └────┬─────┘
         │ begin(saga)
         ▼
    ┌──────────┐         success          ┌──────────┐
    │          │──────────────────────────│          │
    │ RUNNING  │                          │ RUNNING  │◄────┐
    │ step_1   │                          │ step_n   │     │
    │          │◄─────────────────────────│          │     │ success
    └────┬─────┘         failure          └────┬─────┘     │
         │                                      │           │
         │ failure                                │ success   │
         ▼                                      │           │
    ┌──────────┐                                 │           │
    │COMPENSATE│                                 │           │
    │ step_1   │                                 │           │
    │          │                                 │           │
    └────┬─────┘                                 │           │
         │ comp_done                             │           │
         ▼                                      │           │
    ┌──────────┐                                 │           │
    │  FAILED  │◄────────────────────────────────┘           │
    │ (final)  │◄────────────────────────────────────────────┘
    │          │              failure (all steps compensated)
    └──────────┘
         ▲
         │
         │ success (all steps complete)
         │
    ┌────┴─────┐
    │ COMPLETED│
    │ (final)  │
    │          │
    └──────────┘


┌─────────────────────────────────────────────────────────────────────────────┐
│                              STATE ATTRIBUTES                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │    START    │  │   RUNNING   │  │ COMPENSATE  │  │   FINAL     │        │
│  │             │  │             │  │             │  │             │        │
│  │ saga_id     │  │ saga_id     │  │ saga_id     │  │ saga_id     │        │
│  │ definition  │  │ current_step│  │ comp_step   │  │ status      │        │
│  │ input_data  │  │ step_data[] │  │ comp_data[] │  │ result      │        │
│  │ started_at  │  │ started_at  │  │ failed_step │  │ completed_at│        │
│  │             │  │ attempts    │  │             │  │ duration    │        │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Production Go Implementation

### 4.1 Core Process Manager Types

```go
package processmanager

import (
 "context"
 "encoding/json"
 "errors"
 "fmt"
 "sync"
 "time"

 "go.opentelemetry.io/otel/attribute"
 "go.opentelemetry.io/otel/codes"
 "go.opentelemetry.io/otel/metric"
 "go.opentelemetry.io/otel/trace"
 "go.uber.org/zap"
)

// SagaStatus represents the current status of a saga instance
type SagaStatus int

const (
 SagaStatusPending SagaStatus = iota
 SagaStatusRunning
 SagaStatusCompleted
 SagaStatusCompensating
 SagaStatusFailed
 SagaStatusCancelled
)

func (s SagaStatus) String() string {
 switch s {
 case SagaStatusPending:
  return "PENDING"
 case SagaStatusRunning:
  return "RUNNING"
 case SagaStatusCompleted:
  return "COMPLETED"
 case SagaStatusCompensating:
  return "COMPENSATING"
 case SagaStatusFailed:
  return "FAILED"
 case SagaStatusCancelled:
  return "CANCELLED"
 default:
  return "UNKNOWN"
 }
}

// StepStatus represents the status of an individual step
type StepStatus int

const (
 StepStatusPending StepStatus = iota
 StepStatusExecuting
 StepStatusSucceeded
 StepStatusFailed
 StepStatusCompensating
 StepStatusCompensated
 StepStatusCompensationFailed
)

// StepDefinition defines a single step in the saga
type StepDefinition struct {
 Name            string                 `json:"name"`
 Service         string                 `json:"service"`
 Action          string                 `json:"action"`
 Compensable     bool                   `json:"compensable"`
 CompensationAction string              `json:"compensation_action,omitempty"`
 Timeout         time.Duration          `json:"timeout"`
 MaxRetries      int                    `json:"max_retries"`
 RetryDelay      time.Duration          `json:"retry_delay"`
 InputSchema     map[string]interface{} `json:"input_schema,omitempty"`
 OutputSchema    map[string]interface{} `json:"output_schema,omitempty"`
}

// SagaDefinition defines the complete saga structure
type SagaDefinition struct {
 ID          string           `json:"id"`
 Name        string           `json:"name"`
 Version     string           `json:"version"`
 Steps       []StepDefinition `json:"steps"`
 Description string           `json:"description,omitempty"`
 CreatedAt   time.Time        `json:"created_at"`
}

// StepInstance represents a running instance of a step
type StepInstance struct {
 Definition   StepDefinition         `json:"definition"`
 Status       StepStatus             `json:"status"`
 Input        map[string]interface{} `json:"input"`
 Output       map[string]interface{} `json:"output,omitempty"`
 Error        string                 `json:"error,omitempty"`
 StartedAt    *time.Time             `json:"started_at,omitempty"`
 CompletedAt  *time.Time             `json:"completed_at,omitempty"`
 Attempts     int                    `json:"attempts"`
 TraceID      string                 `json:"trace_id,omitempty"`
}

// SagaInstance represents a running saga instance
type SagaInstance struct {
 ID           string                 `json:"id"`
 DefinitionID string                 `json:"definition_id"`
 Status       SagaStatus             `json:"status"`
 Steps        []*StepInstance        `json:"steps"`
 CurrentStep  int                    `json:"current_step"`
 Input        map[string]interface{} `json:"input"`
 Output       map[string]interface{} `json:"output,omitempty"`
 Error        string                 `json:"error,omitempty"`
 CreatedAt    time.Time              `json:"created_at"`
 UpdatedAt    time.Time              `json:"updated_at"`
 CompletedAt  *time.Time             `json:"completed_at,omitempty"`
 TraceID      string                 `json:"trace_id"`
 SpanID       string                 `json:"span_id"`

 mu sync.RWMutex `json:"-"`
}

// ProcessManager is the main saga orchestrator
type ProcessManager struct {
 definitions map[string]*SagaDefinition
 instances   map[string]*SagaInstance
 store       SagaStore
 executor    StepExecutor
 logger      *zap.Logger
 tracer      trace.Tracer
 meter       metric.Meter

 // Metrics
 sagaCreated      metric.Int64Counter
 sagaCompleted    metric.Int64Counter
 sagaFailed       metric.Int64Counter
 sagaCompensated  metric.Int64Counter
 stepExecuted     metric.Int64Counter
 stepDuration     metric.Float64Histogram

 mu sync.RWMutex
}

// SagaStore defines the persistence interface
type SagaStore interface {
 SaveSaga(ctx context.Context, instance *SagaInstance) error
 GetSaga(ctx context.Context, id string) (*SagaInstance, error)
 UpdateSaga(ctx context.Context, instance *SagaInstance) error
 ListActiveSagas(ctx context.Context) ([]*SagaInstance, error)
 SaveEvent(ctx context.Context, sagaID string, event *SagaEvent) error
 GetEvents(ctx context.Context, sagaID string) ([]*SagaEvent, error)
}

// StepExecutor executes steps on target services
type StepExecutor interface {
 Execute(ctx context.Context, step *StepInstance, saga *SagaInstance) error
 Compensate(ctx context.Context, step *StepInstance, saga *SagaInstance) error
}

// SagaEvent represents a saga lifecycle event
type SagaEvent struct {
 ID        string          `json:"id"`
 SagaID    string          `json:"saga_id"`
 Type      SagaEventType   `json:"type"`
 Payload   json.RawMessage `json:"payload"`
 Timestamp time.Time       `json:"timestamp"`
 TraceID   string          `json:"trace_id"`
}

type SagaEventType string

const (
 EventSagaStarted       SagaEventType = "SAGA_STARTED"
 EventStepStarted       SagaEventType = "STEP_STARTED"
 EventStepCompleted     SagaEventType = "STEP_COMPLETED"
 EventStepFailed        SagaEventType = "STEP_FAILED"
 EventSagaCompleted     SagaEventType = "SAGA_COMPLETED"
 EventCompensationStarted SagaEventType = "COMPENSATION_STARTED"
 EventCompensationCompleted SagaEventType = "COMPENSATION_COMPLETED"
 EventSagaFailed        SagaEventType = "SAGA_FAILED"
)

// Config holds the process manager configuration
type Config struct {
 MaxConcurrentSagas int
 DefaultTimeout     time.Duration
 CompensationOrder  CompensationOrder
 EnableMetrics      bool
 EnableTracing      bool
}

type CompensationOrder int

const (
 CompensateReverseOrder CompensationOrder = iota
 CompensateParallel
)

// New creates a new ProcessManager
func New(cfg Config, store SagaStore, executor StepExecutor, logger *zap.Logger, tracer trace.Tracer, meter metric.Meter) (*ProcessManager, error) {
 pm := &ProcessManager{
  definitions: make(map[string]*SagaDefinition),
  instances:   make(map[string]*SagaInstance),
  store:       store,
  executor:    executor,
  logger:      logger,
  tracer:      tracer,
  meter:       meter,
 }

 // Initialize metrics
 if cfg.EnableMetrics && meter != nil {
  var err error
  pm.sagaCreated, err = meter.Int64Counter(
   "saga_created_total",
   metric.WithDescription("Total number of sagas created"),
  )
  if err != nil {
   return nil, fmt.Errorf("failed to create saga_created counter: %w", err)
  }

  pm.sagaCompleted, err = meter.Int64Counter(
   "saga_completed_total",
   metric.WithDescription("Total number of sagas completed successfully"),
  )
  if err != nil {
   return nil, fmt.Errorf("failed to create saga_completed counter: %w", err)
  }

  pm.sagaFailed, err = meter.Int64Counter(
   "saga_failed_total",
   metric.WithDescription("Total number of sagas that failed"),
  )
  if err != nil {
   return nil, fmt.Errorf("failed to create saga_failed counter: %w", err)
  }

  pm.sagaCompensated, err = meter.Int64Counter(
   "saga_compensated_total",
   metric.WithDescription("Total number of sagas that required compensation"),
  )
  if err != nil {
   return nil, fmt.Errorf("failed to create saga_compensated counter: %w", err)
  }

  pm.stepExecuted, err = meter.Int64Counter(
   "saga_step_executed_total",
   metric.WithDescription("Total number of steps executed"),
  )
  if err != nil {
   return nil, fmt.Errorf("failed to create step_executed counter: %w", err)
  }

  pm.stepDuration, err = meter.Float64Histogram(
   "saga_step_duration_seconds",
   metric.WithDescription("Duration of step execution in seconds"),
  )
  if err != nil {
   return nil, fmt.Errorf("failed to create step_duration histogram: %w", err)
  }
 }

 return pm, nil
}

// RegisterDefinition registers a saga definition
func (pm *ProcessManager) RegisterDefinition(def *SagaDefinition) error {
 pm.mu.Lock()
 defer pm.mu.Unlock()

 if _, exists := pm.definitions[def.ID]; exists {
  return fmt.Errorf("saga definition %s already registered", def.ID)
 }

 // Validate definition
 if err := pm.validateDefinition(def); err != nil {
  return fmt.Errorf("invalid saga definition: %w", err)
 }

 pm.definitions[def.ID] = def
 pm.logger.Info("Registered saga definition",
  zap.String("id", def.ID),
  zap.String("name", def.Name),
  zap.Int("steps", len(def.Steps)))

 return nil
}

func (pm *ProcessManager) validateDefinition(def *SagaDefinition) error {
 if def.ID == "" {
  return errors.New("saga ID is required")
 }
 if len(def.Steps) == 0 {
  return errors.New("saga must have at least one step")
 }

 stepNames := make(map[string]bool)
 for i, step := range def.Steps {
  if step.Name == "" {
   return fmt.Errorf("step %d: name is required", i)
  }
  if stepNames[step.Name] {
   return fmt.Errorf("duplicate step name: %s", step.Name)
  }
  stepNames[step.Name] = true
  if step.Service == "" {
   return fmt.Errorf("step %s: service is required", step.Name)
  }
  if step.Action == "" {
   return fmt.Errorf("step %s: action is required", step.Name)
  }
 }

 return nil
}
```

### 4.2 Saga Execution Engine

```go
// StartSaga initiates a new saga instance
func (pm *ProcessManager) StartSaga(ctx context.Context, definitionID string, input map[string]interface{}) (*SagaInstance, error) {
 ctx, span := pm.tracer.Start(ctx, "processmanager.StartSaga",
  trace.WithAttributes(
   attribute.String("saga.definition_id", definitionID),
  ))
 defer span.End()

 pm.mu.RLock()
 def, exists := pm.definitions[definitionID]
 pm.mu.RUnlock()

 if !exists {
  span.SetStatus(codes.Error, "definition not found")
  return nil, fmt.Errorf("saga definition %s not found", definitionID)
 }

 // Create saga instance
 now := time.Now().UTC()
 instance := &SagaInstance{
  ID:           pm.generateSagaID(),
  DefinitionID: definitionID,
  Status:       SagaStatusPending,
  Input:        input,
  Steps:        make([]*StepInstance, len(def.Steps)),
  CurrentStep:  -1,
  CreatedAt:    now,
  UpdatedAt:    now,
  TraceID:      span.SpanContext().TraceID().String(),
  SpanID:       span.SpanContext().SpanID().String(),
 }

 // Initialize steps
 for i, stepDef := range def.Steps {
  instance.Steps[i] = &StepInstance{
   Definition: stepDef,
   Status:     StepStatusPending,
  }
 }

 // Persist initial state
 if err := pm.store.SaveSaga(ctx, instance); err != nil {
  span.RecordError(err)
  return nil, fmt.Errorf("failed to save saga: %w", err)
 }

 // Record event
 event := &SagaEvent{
  ID:        pm.generateEventID(),
  SagaID:    instance.ID,
  Type:      EventSagaStarted,
  Timestamp: now,
  TraceID:   instance.TraceID,
 }
 if err := pm.store.SaveEvent(ctx, instance.ID, event); err != nil {
  pm.logger.Error("Failed to save saga started event", zap.Error(err))
 }

 // Update metrics
 if pm.sagaCreated != nil {
  pm.sagaCreated.Add(ctx, 1, metric.WithAttributes(
   attribute.String("saga_definition", definitionID),
  ))
 }

 pm.logger.Info("Saga started",
  zap.String("saga_id", instance.ID),
  zap.String("definition", definitionID),
  zap.String("trace_id", instance.TraceID))

 // Start execution
 go pm.executeSaga(context.Background(), instance)

 return instance, nil
}

// executeSaga runs the saga state machine
func (pm *ProcessManager) executeSaga(ctx context.Context, instance *SagaInstance) {
 ctx, span := pm.tracer.Start(ctx, "processmanager.executeSaga",
  trace.WithAttributes(
   attribute.String("saga.id", instance.ID),
   attribute.String("saga.definition", instance.DefinitionID),
  ),
  trace.WithLinks(trace.Link{SpanContext: span.SpanContext()}),
 )
 defer span.End()

 instance.mu.Lock()
 instance.Status = SagaStatusRunning
 instance.UpdatedAt = time.Now().UTC()
 instance.mu.Unlock()

 if err := pm.store.UpdateSaga(ctx, instance); err != nil {
  pm.logger.Error("Failed to update saga status", zap.Error(err))
 }

 // Execute steps sequentially
 for instance.CurrentStep < len(instance.Steps)-1 {
  nextStep := instance.CurrentStep + 1
  step := instance.Steps[nextStep]

  pm.logger.Debug("Executing step",
   zap.String("saga_id", instance.ID),
   zap.String("step", step.Definition.Name),
   zap.Int("step_index", nextStep))

  if err := pm.executeStep(ctx, instance, step); err != nil {
   pm.logger.Error("Step execution failed",
    zap.String("saga_id", instance.ID),
    zap.String("step", step.Definition.Name),
    zap.Error(err))

   // Initiate compensation
   if err := pm.compensateSaga(ctx, instance, nextStep); err != nil {
    pm.logger.Error("Compensation failed",
     zap.String("saga_id", instance.ID),
     zap.Error(err))
   }
   return
  }

  instance.mu.Lock()
  instance.CurrentStep = nextStep
  instance.UpdatedAt = time.Now().UTC()
  instance.mu.Unlock()

  if err := pm.store.UpdateSaga(ctx, instance); err != nil {
   pm.logger.Error("Failed to update saga progress", zap.Error(err))
  }
 }

 // Saga completed successfully
 now := time.Now().UTC()
 instance.mu.Lock()
 instance.Status = SagaStatusCompleted
 instance.CompletedAt = &now
 instance.UpdatedAt = now
 instance.mu.Unlock()

 if err := pm.store.UpdateSaga(ctx, instance); err != nil {
  pm.logger.Error("Failed to update saga completion", zap.Error(err))
 }

 // Record completion event
 event := &SagaEvent{
  ID:        pm.generateEventID(),
  SagaID:    instance.ID,
  Type:      EventSagaCompleted,
  Timestamp: now,
  TraceID:   instance.TraceID,
 }
 if err := pm.store.SaveEvent(ctx, instance.ID, event); err != nil {
  pm.logger.Error("Failed to save completion event", zap.Error(err))
 }

 if pm.sagaCompleted != nil {
  pm.sagaCompleted.Add(ctx, 1, metric.WithAttributes(
   attribute.String("saga_definition", instance.DefinitionID),
  ))
 }

 pm.logger.Info("Saga completed successfully",
  zap.String("saga_id", instance.ID),
  zap.Int("steps_completed", len(instance.Steps)))
}

// executeStep executes a single step with retry logic
func (pm *ProcessManager) executeStep(ctx context.Context, saga *SagaInstance, step *StepInstance) error {
 stepCtx, span := pm.tracer.Start(ctx, "processmanager.executeStep",
  trace.WithAttributes(
   attribute.String("saga.id", saga.ID),
   attribute.String("step.name", step.Definition.Name),
   attribute.String("step.service", step.Definition.Service),
  ))
 defer span.End()

 start := time.Now()

 step.mu.Lock()
 step.Status = StepStatusExecuting
 step.StartedAt = &start
 step.TraceID = span.SpanContext().TraceID().String()
 step.mu.Unlock()

 // Record step started event
 event := &SagaEvent{
  ID:        pm.generateEventID(),
  SagaID:    saga.ID,
  Type:      EventStepStarted,
  Timestamp: start,
  TraceID:   saga.TraceID,
 }
 if err := pm.store.SaveEvent(stepCtx, saga.ID, event); err != nil {
  pm.logger.Error("Failed to save step started event", zap.Error(err))
 }

 // Execute with retry logic
 var lastErr error
 maxRetries := step.Definition.MaxRetries
 if maxRetries == 0 {
  maxRetries = 3 // Default retries
 }

 for attempt := 0; attempt <= maxRetries; attempt++ {
  if attempt > 0 {
   delay := step.Definition.RetryDelay
   if delay == 0 {
    delay = time.Duration(attempt) * time.Second // Exponential backoff
   }
   pm.logger.Debug("Retrying step",
    zap.String("saga_id", saga.ID),
    zap.String("step", step.Definition.Name),
    zap.Int("attempt", attempt))
   time.Sleep(delay)
  }

  step.mu.Lock()
  step.Attempts = attempt + 1
  step.mu.Unlock()

  err := pm.executor.Execute(stepCtx, step, saga)
  if err == nil {
   // Success
   completedAt := time.Now().UTC()
   step.mu.Lock()
   step.Status = StepStatusSucceeded
   step.CompletedAt = &completedAt
   step.mu.Unlock()

   duration := time.Since(start).Seconds()
   if pm.stepDuration != nil {
    pm.stepDuration.Record(stepCtx, duration, metric.WithAttributes(
     attribute.String("step_name", step.Definition.Name),
     attribute.String("result", "success"),
    ))
   }

   if pm.stepExecuted != nil {
    pm.stepExecuted.Add(stepCtx, 1, metric.WithAttributes(
     attribute.String("step_name", step.Definition.Name),
     attribute.String("result", "success"),
    ))
   }

   // Record step completed event
   event := &SagaEvent{
    ID:        pm.generateEventID(),
    SagaID:    saga.ID,
    Type:      EventStepCompleted,
    Timestamp: completedAt,
    TraceID:   saga.TraceID,
   }
   if err := pm.store.SaveEvent(stepCtx, saga.ID, event); err != nil {
    pm.logger.Error("Failed to save step completed event", zap.Error(err))
   }

   return nil
  }

  lastErr = err
  pm.logger.Warn("Step execution attempt failed",
   zap.String("saga_id", saga.ID),
   zap.String("step", step.Definition.Name),
   zap.Int("attempt", attempt+1),
   zap.Error(err))
 }

 // All retries exhausted
 step.mu.Lock()
 step.Status = StepStatusFailed
 step.Error = lastErr.Error()
 step.mu.Unlock()

 span.RecordError(lastErr)
 span.SetStatus(codes.Error, "step execution failed")

 if pm.stepExecuted != nil {
  pm.stepExecuted.Add(stepCtx, 1, metric.WithAttributes(
   attribute.String("step_name", step.Definition.Name),
   attribute.String("result", "failure"),
  ))
 }

 // Record step failed event
 event = &SagaEvent{
  ID:        pm.generateEventID(),
  SagaID:    saga.ID,
  Type:      EventStepFailed,
  Timestamp: time.Now().UTC(),
  TraceID:   saga.TraceID,
 }
 if err := pm.store.SaveEvent(stepCtx, saga.ID, event); err != nil {
  pm.logger.Error("Failed to save step failed event", zap.Error(err))
 }

 return fmt.Errorf("step %s failed after %d attempts: %w", step.Definition.Name, maxRetries+1, lastErr)
}

// compensateSaga initiates compensation for completed steps
func (pm *ProcessManager) compensateSaga(ctx context.Context, saga *SagaInstance, failedStep int) error {
 ctx, span := pm.tracer.Start(ctx, "processmanager.compensateSaga",
  trace.WithAttributes(
   attribute.String("saga.id", saga.ID),
   attribute.Int("failed_step", failedStep),
  ))
 defer span.End()

 pm.logger.Info("Starting saga compensation",
  zap.String("saga_id", saga.ID),
  zap.Int("completed_steps", failedStep))

 saga.mu.Lock()
 saga.Status = SagaStatusCompensating
 saga.UpdatedAt = time.Now().UTC()
 saga.mu.Unlock()

 if err := pm.store.UpdateSaga(ctx, saga); err != nil {
  pm.logger.Error("Failed to update saga status", zap.Error(err))
 }

 // Record compensation started event
 event := &SagaEvent{
  ID:        pm.generateEventID(),
  SagaID:    saga.ID,
  Type:      EventCompensationStarted,
  Timestamp: time.Now().UTC(),
  TraceID:   saga.TraceID,
 }
 if err := pm.store.SaveEvent(ctx, saga.ID, event); err != nil {
  pm.logger.Error("Failed to save compensation started event", zap.Error(err))
 }

 if pm.sagaCompensated != nil {
  pm.sagaCompensated.Add(ctx, 1, metric.WithAttributes(
   attribute.String("saga_definition", saga.DefinitionID),
  ))
 }

 // Compensate in reverse order
 for i := failedStep - 1; i >= 0; i-- {
  step := saga.Steps[i]
  if !step.Definition.Compensable {
   pm.logger.Warn("Step is not compensable, skipping",
    zap.String("saga_id", saga.ID),
    zap.String("step", step.Definition.Name))
   continue
  }

  if step.Status != StepStatusSucceeded {
   // Only compensate successfully completed steps
   continue
  }

  if err := pm.compensateStep(ctx, saga, step); err != nil {
   pm.logger.Error("Compensation failed",
    zap.String("saga_id", saga.ID),
    zap.String("step", step.Definition.Name),
    zap.Error(err))
   // Continue with other compensations even if one fails
   // This is a design decision - could also abort and alert
  }
 }

 // Mark saga as failed
 now := time.Now().UTC()
 saga.mu.Lock()
 saga.Status = SagaStatusFailed
 saga.CompletedAt = &now
 saga.UpdatedAt = now
 saga.mu.Unlock()

 if err := pm.store.UpdateSaga(ctx, saga); err != nil {
  pm.logger.Error("Failed to update saga failure status", zap.Error(err))
 }

 // Record compensation completed event
 event = &SagaEvent{
  ID:        pm.generateEventID(),
  SagaID:    saga.ID,
  Type:      EventCompensationCompleted,
  Timestamp: now,
  TraceID:   saga.TraceID,
 }
 if err := pm.store.SaveEvent(ctx, saga.ID, event); err != nil {
  pm.logger.Error("Failed to save compensation completed event", zap.Error(err))
 }

 if pm.sagaFailed != nil {
  pm.sagaFailed.Add(ctx, 1, metric.WithAttributes(
   attribute.String("saga_definition", saga.DefinitionID),
  ))
 }

 pm.logger.Info("Saga compensation completed",
  zap.String("saga_id", saga.ID),
  zap.Int("steps_compensated", failedStep))

 return nil
}

// compensateStep executes compensation for a single step
func (pm *ProcessManager) compensateStep(ctx context.Context, saga *SagaInstance, step *StepInstance) error {
 stepCtx, span := pm.tracer.Start(ctx, "processmanager.compensateStep",
  trace.WithAttributes(
   attribute.String("saga.id", saga.ID),
   attribute.String("step.name", step.Definition.Name),
  ))
 defer span.End()

 step.mu.Lock()
 step.Status = StepStatusCompensating
 step.mu.Unlock()

 pm.logger.Info("Compensating step",
  zap.String("saga_id", saga.ID),
  zap.String("step", step.Definition.Name))

 if err := pm.executor.Compensate(stepCtx, step, saga); err != nil {
  step.mu.Lock()
  step.Status = StepStatusCompensationFailed
  step.Error = err.Error()
  step.mu.Unlock()

  span.RecordError(err)
  return fmt.Errorf("compensation failed for step %s: %w", step.Definition.Name, err)
 }

 step.mu.Lock()
 step.Status = StepStatusCompensated
 step.mu.Unlock()

 pm.logger.Info("Step compensated successfully",
  zap.String("saga_id", saga.ID),
  zap.String("step", step.Definition.Name))

 return nil
}

// generateSagaID generates a unique saga ID
func (pm *ProcessManager) generateSagaID() string {
 return fmt.Sprintf("saga-%d-%d", time.Now().UnixNano(), time.Now().Nanosecond())
}

// generateEventID generates a unique event ID
func (pm *ProcessManager) generateEventID() string {
 return fmt.Sprintf("evt-%d-%d", time.Now().UnixNano(), time.Now().Nanosecond())
}
```

### 4.3 Step Executor Implementation

```go
package processmanager

import (
 "bytes"
 "context"
 "encoding/json"
 "fmt"
 "net/http"
 "time"

 "go.opentelemetry.io/otel/codes"
 "go.opentelemetry.io/otel/trace"
)

// HTTPStepExecutor executes steps via HTTP calls
type HTTPStepExecutor struct {
 client     *http.Client
 baseURLs   map[string]string // service name -> base URL
 tracer     trace.Tracer
}

// NewHTTPStepExecutor creates a new HTTP step executor
func NewHTTPStepExecutor(timeout time.Duration, tracer trace.Tracer) *HTTPStepExecutor {
 return &HTTPStepExecutor{
  client: &http.Client{
   Timeout: timeout,
  },
  baseURLs: make(map[string]string),
  tracer:   tracer,
 }
}

// RegisterService registers a service base URL
func (e *HTTPStepExecutor) RegisterService(name, baseURL string) {
 e.baseURLs[name] = baseURL
}

// Execute executes a step by calling the target service
func (e *HTTPStepExecutor) Execute(ctx context.Context, step *StepInstance, saga *SagaInstance) error {
 ctx, span := e.tracer.Start(ctx, "http_executor.Execute",
  trace.WithAttributes(
   trace.String("step.name", step.Definition.Name),
   trace.String("step.service", step.Definition.Service),
  ),
 )
 defer span.End()

 baseURL, exists := e.baseURLs[step.Definition.Service]
 if !exists {
  err := fmt.Errorf("service %s not registered", step.Definition.Service)
  span.RecordError(err)
  span.SetStatus(codes.Error, err.Error())
  return err
 }

 // Prepare request payload
 payload := map[string]interface{}{
  "saga_id":     saga.ID,
  "step_name":   step.Definition.Name,
  "input":       step.Input,
  "saga_input":  saga.Input,
  "trace_id":    saga.TraceID,
 }

 body, err := json.Marshal(payload)
 if err != nil {
  span.RecordError(err)
  span.SetStatus(codes.Error, "failed to marshal request")
  return fmt.Errorf("failed to marshal request: %w", err)
 }

 // Build URL
 url := fmt.Sprintf("%s/%s", baseURL, step.Definition.Action)

 // Create request
 req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
 if err != nil {
  span.RecordError(err)
  span.SetStatus(codes.Error, "failed to create request")
  return fmt.Errorf("failed to create request: %w", err)
 }

 req.Header.Set("Content-Type", "application/json")
 req.Header.Set("X-Saga-ID", saga.ID)
 req.Header.Set("X-Trace-ID", saga.TraceID)

 // Execute request
 resp, err := e.client.Do(req)
 if err != nil {
  span.RecordError(err)
  span.SetStatus(codes.Error, "request failed")
  return fmt.Errorf("request failed: %w", err)
 }
 defer resp.Body.Close()

 if resp.StatusCode >= 200 && resp.StatusCode < 300 {
  // Parse response
  var result map[string]interface{}
  if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
   step.mu.Lock()
   step.Output = result
   step.mu.Unlock()
  }
  return nil
 }

 return fmt.Errorf("step execution failed with status %d", resp.StatusCode)
}

// Compensate executes compensation for a step
func (e *HTTPStepExecutor) Compensate(ctx context.Context, step *StepInstance, saga *SagaInstance) error {
 ctx, span := e.tracer.Start(ctx, "http_executor.Compensate",
  trace.WithAttributes(
   trace.String("step.name", step.Definition.Name),
   trace.String("step.service", step.Definition.Service),
  ),
 )
 defer span.End()

 baseURL, exists := e.baseURLs[step.Definition.Service]
 if !exists {
  err := fmt.Errorf("service %s not registered", step.Definition.Service)
  span.RecordError(err)
  return err
 }

 compAction := step.Definition.CompensationAction
 if compAction == "" {
  compAction = fmt.Sprintf("compensate/%s", step.Definition.Action)
 }

 payload := map[string]interface{}{
  "saga_id":    saga.ID,
  "step_name":  step.Definition.Name,
  "step_input": step.Input,
  "step_output": step.Output,
  "trace_id":   saga.TraceID,
 }

 body, err := json.Marshal(payload)
 if err != nil {
  span.RecordError(err)
  return fmt.Errorf("failed to marshal compensation request: %w", err)
 }

 url := fmt.Sprintf("%s/%s", baseURL, compAction)
 req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
 if err != nil {
  span.RecordError(err)
  return fmt.Errorf("failed to create compensation request: %w", err)
 }

 req.Header.Set("Content-Type", "application/json")
 req.Header.Set("X-Saga-ID", saga.ID)
 req.Header.Set("X-Compensation", "true")
 req.Header.Set("X-Trace-ID", saga.TraceID)

 resp, err := e.client.Do(req)
 if err != nil {
  span.RecordError(err)
  return fmt.Errorf("compensation request failed: %w", err)
 }
 defer resp.Body.Close()

 if resp.StatusCode >= 200 && resp.StatusCode < 300 {
  return nil
 }

 return fmt.Errorf("compensation failed with status %d", resp.StatusCode)
}
```

---

## 5. Failure Scenarios and Mitigations

### 5.1 Failure Taxonomy

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        FAILURE SCENARIOS MATRIX                              │
├───────────────────────────────┬───────────────────┬─────────────────────────┤
│         Scenario              │     Detection     │      Mitigation         │
├───────────────────────────────┼───────────────────┼─────────────────────────┤
│ Service Timeout               │ P99 latency spike │ Circuit breaker + Retry │
│ Service Unavailable           │ Connection error  │ Mark step failed →      │
│                               │                   │ Trigger compensation    │
├───────────────────────────────┼───────────────────┼─────────────────────────┤
│ Process Manager Crash         │ Heartbeat timeout │ Recovery from event     │
│                               │                   │ store on restart        │
├───────────────────────────────┼───────────────────┼─────────────────────────┤
│ Compensation Failure          │ Non-2xx response  │ Alert + Manual          │
│                               │                   │ intervention queue      │
├───────────────────────────────┼───────────────────┼─────────────────────────┤
│ Duplicate Saga Execution      │ Idempotency key   │ Idempotency check in    │
│                               │ collision         │ event store             │
├───────────────────────────────┼───────────────────┼─────────────────────────┤
│ Partial Step Execution        │ Ambiguous status  │ Idempotent operations + │
│                               │                   │ At-least-once delivery  │
├───────────────────────────────┼───────────────────┼─────────────────────────┤
│ Event Store Unavailable       │ Write timeout     │ Local WAL buffer +      │
│                               │                   │ Replay on recovery      │
├───────────────────────────────┼───────────────────┼─────────────────────────┤
│ Cascading Compensation        │ Resource spike    │ Rate limiting +         │
│                               │                   │ Exponential backoff     │
└───────────────────────────────┴───────────────────┴─────────────────────────┘
```

### 5.2 Recovery Mechanisms

```go
// RecoveryService handles PM crash recovery
type RecoveryService struct {
 store  SagaStore
 pm     *ProcessManager
 logger *zap.Logger
}

// Recover recovers incomplete sagas after crash
func (r *RecoveryService) Recover(ctx context.Context) error {
 // Get all active sagas
 sagas, err := r.store.ListActiveSagas(ctx)
 if err != nil {
  return fmt.Errorf("failed to list active sagas: %w", err)
 }

 r.logger.Info("Recovering incomplete sagas", zap.Int("count", len(sagas)))

 for _, saga := range sagas {
  switch saga.Status {
  case SagaStatusRunning:
   // Resume from current step
   r.logger.Info("Resuming saga",
    zap.String("saga_id", saga.ID),
    zap.Int("current_step", saga.CurrentStep))
   go r.pm.executeSaga(ctx, saga)

  case SagaStatusCompensating:
   // Resume compensation
   r.logger.Info("Resuming compensation",
    zap.String("saga_id", saga.ID))
   go r.pm.compensateSaga(ctx, saga, saga.CurrentStep+1)

  case SagaStatusPending:
   // Start execution
   r.logger.Info("Starting pending saga",
    zap.String("saga_id", saga.ID))
   go r.pm.executeSaga(ctx, saga)
  }
 }

 return nil
}
```

### 5.3 Production Hardening

| Aspect | Strategy | Implementation |
|--------|----------|----------------|
| **Idempotency** | Request deduplication | Idempotency key in headers, check before execution |
| **Observability** | Distributed tracing | OpenTelemetry spans for all operations |
| **Metrics** | Key performance indicators | Counter/Histogram for success rate, duration |
| **Alerting** | Threshold-based alerts | P95 latency, error rate, compensation rate |
| **Rate Limiting** | Prevent overload | Token bucket for step execution |
| **Circuit Breaking** | Fail fast on errors | Per-service circuit breakers |
| **Dead Letter Queue** | Handle persistent failures | DLQ for failed compensations |

---

## 6. Semantic Trade-off Analysis

### 6.1 Orchestration vs Choreography

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    ORCHESTRATION vs CHOREOGRAPHY                             │
├─────────────────────┬──────────────────────────┬────────────────────────────┤
│     Dimension       │      Orchestration       │       Choreography         │
├─────────────────────┼──────────────────────────┼────────────────────────────┤
│ Central Control     │ ✅ Explicit state machine │ ❌ Distributed logic       │
│ Coupling            │ ❌ Tight to PM           │ ✅ Loosely coupled         │
│ Observability       │ ✅ Single source of truth │ ❌ Distributed tracing     │
│ Complexity          │ ⚠️  Centralized complexity │ ⚠️  Distributed complexity │
│ Recovery            │ ✅ Coordinated recovery   │ ❌ Requires saga discovery │
│ Scalability         │ ⚠️  PM bottleneck         │ ✅ Natural horizontal      │
│ Transaction Isolation│ ✅ ACID at PM            │ ❌ Eventual consistency    │
└─────────────────────┴──────────────────────────┴────────────────────────────┘

Recommendation: Use Orchestration for complex business processes requiring
strong consistency guarantees; use Choreography for simple, independent flows.
```

### 6.2 Synchronous vs Asynchronous Compensation

| Approach | Latency | Reliability | Complexity | Use Case |
|----------|---------|-------------|------------|----------|
| **Synchronous** | Higher | Easier to track | Lower | Financial transactions |
| **Asynchronous** | Lower | Requires queue | Higher | High-throughput flows |
| **Parallel** | Lowest | Harder to debug | Medium | Independent compensations |

### 6.3 Storage Backend Trade-offs

| Backend | Consistency | Availability | Performance | Cost |
|---------|-------------|--------------|-------------|------|
| PostgreSQL | Strong | Good | Medium | Low |
| MongoDB | Eventual | Excellent | High | Medium |
| Event Store DB | Strong | Good | Medium | High |
| Redis + WAL | Eventual | Excellent | Very High | Low |

---

## 7. References

1. Hohpe, G., & Woolf, B. (2004). *Enterprise Integration Patterns*. Addison-Wesley.
2. Richardson, C. (2018). *Microservices Patterns*. Manning Publications.
3. Kleppmann, M. (2017). *Designing Data-Intensive Applications*. O'Reilly Media.
4. AWS Architecture Center. (2024). Saga Pattern Implementation.
5. Microsoft Azure. (2024). Saga Pattern for Microservices.
6. Fowler, M. (2005). Event Sourcing. martinfowler.com.

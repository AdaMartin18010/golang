# FT-023: SAGA Pattern - Formal Specification

## Overview

The SAGA pattern is a design pattern for managing distributed transactions across multiple services. Unlike ACID transactions, SAGAs split a long-running transaction into a sequence of local transactions, each with a corresponding compensating action to undo its effects if needed.

## Theoretical Foundations

### 1.1 SAGA Model

**SAGA Definition**:

```
A SAGA S = ⟨T, C, ≺, σ⟩ where:
- T = {t₁, t₂, ..., tₙ}: set of transaction steps
- C = {c₁, c₂, ..., cₙ}: set of compensating actions (cᵢ compensates tᵢ)
- ≺: partial order defining execution sequence
- σ: state machine governing SAGA execution
```

**SAGA Execution States**:

```
State(S) ∈ {
  STARTED,       // SAGA initiated
  RUNNING(i),    // Executing step i
  COMPENSATING(j), // Compensating from step j backwards
  COMPLETED,     // All steps executed successfully
  ABORTED        // SAGA aborted, compensations complete
}
```

### 1.2 SAGA Types

**Choreography-Based SAGA**:

```
Each service:
  1. Completes local transaction
  2. Publishes event
  3. Next service listens and executes
  4. On failure, publishes compensation event

No central coordinator; services communicate via events.
```

**Orchestration-Based SAGA**:

```
Central SAGA Orchestrator:
  1. Sends command to Service 1
  2. On success, sends command to Service 2
  3. On failure, sends compensation to Service 1
  4. Continues or compensates as needed

Explicit coordinator manages the flow.
```

### 1.3 SAGA Properties

**Compensating Transaction**:

```
For each transaction step tᵢ with effect E(tᵢ),
compensating action cᵢ satisfies:
  E(cᵢ) ∘ E(tᵢ) = identity

In practice, cᵢ restores the system to a semantically
consistent state (not necessarily bitwise identical).
```

**SAGA Guarantee**:

```
Theorem: A SAGA ensures eventual consistency.

Proof:
Case 1: All steps succeed
  - SAGA executes t₁, t₂, ..., tₙ sequentially
  - No compensations needed
  - Result: COMPLETED

Case 2: Step tₖ fails
  - Steps t₁, ..., tₖ₋₁ have committed
  - Execute compensations cₖ₋₁, ..., c₁ in reverse order
  - Each cᵢ undoes the effect of tᵢ
  - Result: Semantically equivalent to never having started

Therefore, system is always left in a consistent state. ∎
```

**Compensation Properties**:

```
1. Order Preservation: Compensations execute in reverse order
   cᵢ executes before cⱼ if i > j

2. Partial Compensation: If compensation cᵢ fails, retry or escalate

3. Idempotence: Each compensation should be idempotent
   applying cᵢ twice has same effect as once

4. Durability: Compensated state persists across failures
```

### 1.4 Comparison with 2PC

| Aspect | 2PC/3PC | SAGA |
|--------|---------|------|
| **Isolation** | Strong (locks) | Weak (compensations) |
| **Consistency** | Immediate | Eventual |
| **Availability** | Lower (blocking) | Higher (non-blocking) |
| **Complexity** | Protocol complexity | Business logic complexity |
| **Latency** | Synchronous | Asynchronous |
| **Use Case** | Short transactions | Long-running processes |

## TLA+ Specification

```tla
----------------------------- MODULE SAGA -----------------------------
EXTENDS Integers, Sequences, FiniteSets, TLC

CONSTANTS Steps,           \* Set of SAGA steps
          Services,        \* Set of services
          MaxRetries       \* Maximum retry attempts

VARIABLES sagaState,       \* Current SAGA state
          stepStates,      \* State of each step
          compensations,   \* Compensation tracking
          executionOrder,  \* Order of step execution
          completed        \* Completed steps

\* Step state
StepState == {"pending", "executing", "completed", "failed",
              "compensating", "compensated"}

\* SAGA states
SagaState == {"started", "running", "compensating", "completed", "aborted"}

\* Step definition
Step == [id: Steps,
         service: Services,
         compensation: Steps,
         order: Nat]

\* Initial state
Init ==
  ∧ sagaState = "started"
  ∧ stepStates = [s ∈ Steps ↦ "pending"]
  ∧ compensations = ⟨⟩
  ∧ executionOrder = ⟨⟩
  ∧ completed = {}

\* Execute a step
ExecuteStep(s) ==
  ∧ sagaState ∈ {"started", "running"}
  ∧ stepStates[s] = "pending"
  ∧ ∀s2 ∈ Steps: StepOrder(s2) < StepOrder(s) ⇒ stepStates[s2] = "completed"
  ∧ stepStates' = [stepStates EXCEPT ![s] = "executing"]
  ∧ sagaState' = "running"
  ∧ UNCHANGED ⟨compensations, executionOrder, completed⟩

\* Step completes successfully
CompleteStep(s) ==
  ∧ stepStates[s] = "executing"
  ∧ stepStates' = [stepStates EXCEPT ![s] = "completed"]
  ∧ executionOrder' = Append(executionOrder, s)
  ∧ completed' = completed ∪ {s}
  ∧ IF completed' = Steps
     THEN sagaState' = "completed"
     ELSE UNCHANGED sagaState
  ∧ UNCHANGED compensations

\* Step fails
FailStep(s) ==
  ∧ stepStates[s] = "executing"
  ∧ stepStates' = [stepStates EXCEPT ![s] = "failed"]
  ∧ sagaState' = "compensating"
  ∧ compensations' = Reverse(executionOrder)  \* Compensate in reverse
  ∧ UNCHANGED ⟨executionOrder, completed⟩

\* Compensate a completed step
CompensateStep(s) ==
  ∧ sagaState = "compensating"
  ∧ stepStates[s] = "completed"
  ∧ Head(compensations) = s
  ∧ stepStates' = [stepStates EXCEPT ![s] = "compensating"]
  ∧ UNCHANGED ⟨sagaState, compensations, executionOrder, completed⟩

\* Compensation completes
CompleteCompensation(s) ==
  ∧ stepStates[s] = "compensating"
  ∧ stepStates' = [stepStates EXCEPT ![s] = "compensated"]
  ∧ compensations' = Tail(compensations)
  ∧ IF compensations' = ⟨⟩
     THEN sagaState' = "aborted"
     ELSE UNCHANGED sagaState
  ∧ UNCHANGED ⟨executionOrder, completed⟩

\* Helper: Determine step order
StepOrder(s) == CHOOSE n ∈ Nat : TRUE  \* Abstract ordering

\* Helper: Reverse a sequence
Reverse(seq) ==
  IF seq = ⟨⟩
  THEN ⟨⟩
  ELSE Append(Reverse(Tail(seq)), Head(seq))

\* Next state
Next ==
  ∨ ∃s ∈ Steps: ExecuteStep(s)
  ∨ ∃s ∈ Steps: CompleteStep(s)
  ∨ ∃s ∈ Steps: FailStep(s)
  ∨ ∃s ∈ Steps: CompensateStep(s)
  ∨ ∃s ∈ Steps: CompleteCompensation(s)

\* Invariants

\* At most one step executing at a time
SingleExecution ==
  Cardinality({s ∈ Steps : stepStates[s] = "executing"}) ≤ 1

\* Compensation order is reverse of execution
CompensationOrder ==
  ∀i, j ∈ 1..Len(compensations):
    i < j ⇒ Position(executionOrder, compensations[i]) >
            Position(executionOrder, compensations[j])

\* Final states are consistent
FinalConsistency ==
  (sagaState = "completed") ⇒ ∀s ∈ Steps: stepStates[s] = "completed"

\* Compensated state consistency
CompensationConsistency ==
  (sagaState = "aborted") ⇒
    ∀s ∈ completed: stepStates[s] = "compensated"

=============================================================================
```

## Algorithm Pseudocode

### Orchestration-Based SAGA

```
Algorithm: Orchestration-Based SAGA

Types:
  StepID: unique identifier for each step
  StepResult: SUCCESS | FAILURE | TIMEOUT
  SagaState: STARTED | RUNNING | COMPENSATING | COMPLETED | ABORTED

SagaDefinition:
  steps: List<Step> where Step = {
    id: StepID,
    service: Service,
    action: () -> Result,
    compensation: () -> Result,
    retryPolicy: RetryPolicy
  }

Orchestrator:
  State:
    sagaId: unique identifier
    definition: SagaDefinition
    currentStep: int = 0
    state: SagaState = STARTED
    executedSteps: Stack<StepID>
    results: Map<StepID, Result>

  ExecuteSaga():
    state = RUNNING
    log(SAGA_STARTED, sagaId)

    for i from 0 to definition.steps.length - 1:
      currentStep = i
      step = definition.steps[i]

      log(STEP_STARTED, sagaId, step.id)
      result = ExecuteStep(step)

      if result == SUCCESS:
        log(STEP_COMPLETED, sagaId, step.id)
        executedSteps.push(step.id)
        results[step.id] = result
      else:
        log(STEP_FAILED, sagaId, step.id, result)
        return HandleFailure(step, result)

    state = COMPLETED
    log(SAGA_COMPLETED, sagaId)
    return SUCCESS

  ExecuteStep(step):
    attempts = 0

    while attempts < step.retryPolicy.maxAttempts:
      try:
        result = step.action()
        if result.success:
          return SUCCESS
      catch Exception e:
        log(ERROR, sagaId, step.id, e)

      attempts++
      if attempts < step.retryPolicy.maxAttempts:
        sleep(step.retryPolicy.backoff(attempts))

    return FAILURE

  HandleFailure(failedStep, error):
    state = COMPENSATING
    log(COMPENSATION_STARTED, sagaId, failedStep.id)

    while not executedSteps.empty():
      stepId = executedSteps.pop()
      step = definition.getStep(stepId)

      log(COMPENSATING_STEP, sagaId, stepId)
      compResult = ExecuteCompensation(step)

      if compResult != SUCCESS:
        // Compensation failed - critical error
        log(COMPENSATION_FAILED, sagaId, stepId)
        return HandleCompensationFailure(step, compResult)

    state = ABORTED
    log(SAGA_ABORTED, sagaId)
    return FAILURE

  ExecuteCompensation(step):
    attempts = 0

    while attempts < step.retryPolicy.maxCompensationAttempts:
      try:
        result = step.compensation()
        if result.success:
          log(COMPENSATION_SUCCESS, sagaId, step.id)
          return SUCCESS
      catch Exception e:
        log(COMPENSATION_ERROR, sagaId, step.id, e)

      attempts++
      if attempts < step.retryPolicy.maxCompensationAttempts:
        sleep(step.retryPolicy.compensationBackoff(attempts))

    return FAILURE

  HandleCompensationFailure(step, error):
    // Escalate - may require human intervention
    log(CRITICAL_ERROR, sagaId, step.id, error)
    AlertOperationsTeam(sagaId, step, error)
    // Leave SAGA in COMPENSATING state for manual resolution
    return CRITICAL_FAILURE
```

### Choreography-Based SAGA

```
Algorithm: Choreography-Based SAGA

Types:
  Event: SAGA events exchanged between services
  EventType: STEP_COMPLETED | STEP_FAILED | COMPENSATION_REQUIRED

Service A:
  On StartSagaCommand:
    result = ExecuteLocalTransaction()
    if result.success:
      PublishEvent(STEP_COMPLETED, {
        sagaId: sagaId,
        stepId: "A",
        nextStep: "B",
        payload: result.data
      })
    else:
      PublishEvent(STEP_FAILED, {
        sagaId: sagaId,
        stepId: "A",
        error: result.error
      })

  On CompensationRequiredEvent(event):
    if event.stepId == "A":
      ExecuteCompensation()
      PublishEvent(COMPENSATION_COMPLETED, {
        sagaId: sagaId,
        stepId: "A"
      })

Service B:
  On StepCompletedEvent(event):
    if event.nextStep == "B":
      result = ExecuteLocalTransaction(event.payload)
      if result.success:
        PublishEvent(STEP_COMPLETED, {
          sagaId: sagaId,
          stepId: "B",
          nextStep: "C",
          payload: result.data
        })
      else:
        PublishEvent(STEP_FAILED, {
          sagaId: sagaId,
          stepId: "B",
          error: result.error
        })

  On StepFailedEvent(event):
    // If previous step failed, may need to compensate
    if event.stepId != "B":
      // Check if we need to compensate
      if HasExecuted(sagaId, "B"):
        ExecuteCompensation()
        PublishEvent(COMPENSATION_COMPLETED, {
          sagaId: sagaId,
          stepId: "B"
        })

Event Processor:
  State:
    sagaStates: Map<SagaID, SagaState>
    stepHistory: Map<SagaID, List<StepEvent>>

  On Event(event):
    switch event.type:
      case STEP_COMPLETED:
        RecordStep(event.sagaId, event)
        if IsFinalStep(event):
          MarkSagaCompleted(event.sagaId)

      case STEP_FAILED:
        RecordStep(event.sagaId, event)
        InitiateCompensation(event.sagaId, event.stepId)

      case COMPENSATION_COMPLETED:
        RecordCompensation(event.sagaId, event)
        if AllCompensationsComplete(event.sagaId):
          MarkSagaAborted(event.sagaId)

  InitiateCompensation(sagaId, failedStep):
    // Determine which steps need compensation
    stepsToCompensate = GetExecutedStepsBefore(sagaId, failedStep)

    for step in reverse(stepsToCompensate):
      PublishEvent(COMPENSATION_REQUIRED, {
        sagaId: sagaId,
        stepId: step.id
      })
```

### Parallel SAGA Execution

```
Algorithm: Parallel SAGA with Dependency Management

SAGA Structure:
  steps: List<Step>
  dependencies: Map<StepID, Set<StepID>>
    // Step cannot start until all dependencies complete

Execution Graph:
  Build dependency graph from SAGA definition
  Identify parallelizable steps (no dependencies between them)

Orchestrator.ExecuteParallel():
  state = RUNNING
  completedSteps = Set()
  executingSteps = Set()
  failedStep = null

  while completedSteps ≠ allSteps:
    // Find ready steps (all dependencies satisfied)
    readySteps = {s ∈ steps :
      dependencies[s] ⊆ completedSteps ∧ s ∉ completedSteps ∧ s ∉ executingSteps}

    if readySteps is empty and executingSteps is empty:
      // Deadlock or all done
      break

    // Execute ready steps in parallel
    for step in readySteps:
      executingSteps.add(step)
      go ExecuteStepAsync(step)

    // Wait for any step to complete
    result = WaitForAnyCompletion()

    if result.success:
      executingSteps.remove(result.step)
      completedSteps.add(result.step)
      log(STEP_COMPLETED, result.step)
    else:
      // Failure - stop new executions, compensate completed
      failedStep = result.step
      break

  if failedStep ≠ null:
    // Cancel pending executions
    for step in executingSteps:
      CancelExecution(step)

    // Compensate completed steps
    return CompensateParallel(completedSteps)

  state = COMPLETED
  return SUCCESS

CompensateParallel(completedSteps):
  state = COMPENSATING

  // Build compensation dependency graph (reverse of execution)
  compDeps = BuildCompensationDependencies(completedSteps)

  compensatedSteps = Set()
  compensatingSteps = Set()

  while compensatedSteps ≠ completedSteps:
    // Find steps ready for compensation
    readyCompensations = {s ∈ completedSteps :
      s ∉ compensatedSteps ∧
      s ∉ compensatingSteps ∧
      compDeps[s] ⊆ compensatedSteps}

    // Execute compensations in parallel
    for step in readyCompensations:
      compensatingSteps.add(step)
      go CompensateStepAsync(step)

    // Wait for completions
    result = WaitForAnyCompensation()

    if result.success:
      compensatingSteps.remove(result.step)
      compensatedSteps.add(result.step)
    else:
      // Compensation failure - escalate
      return HandleCompensationFailure(result.step)

  state = ABORTED
  return FAILURE
```

## Go Implementation

```go
// Package saga implements the SAGA pattern for distributed transactions
package saga

import (
 "context"
 "fmt"
 "sync"
 "time"
)

// StepState represents the state of a SAGA step
type StepState int

const (
 StepPending StepState = iota
 StepExecuting
 StepCompleted
 StepFailed
 StepCompensating
 StepCompensated
)

func (s StepState) String() string {
 switch s {
 case StepPending:
  return "PENDING"
 case StepExecuting:
  return "EXECUTING"
 case StepCompleted:
  return "COMPLETED"
 case StepFailed:
  return "FAILED"
 case StepCompensating:
  return "COMPENSATING"
 case StepCompensated:
  return "COMPENSATED"
 default:
  return "UNKNOWN"
 }
}

// SagaState represents the overall SAGA state
type SagaState int

const (
 SagaStarted SagaState = iota
 SagaRunning
 SagaCompensating
 SagaCompleted
 SagaAborted
)

// StepResult represents the result of executing a step
type StepResult struct {
 Success bool
 Data    interface{}
 Error   error
}

// Step defines a SAGA step
type Step struct {
 ID            string
 Name          string
 Service       string
 Action        func(ctx context.Context, input interface{}) (interface{}, error)
 Compensation  func(ctx context.Context, input interface{}) error
 RetryPolicy   *RetryPolicy
 Dependencies  []string // IDs of steps that must complete before this one
}

// RetryPolicy defines retry behavior
type RetryPolicy struct {
 MaxAttempts            int
 InitialBackoff         time.Duration
 MaxBackoff             time.Duration
 BackoffMultiplier      float64
 MaxCompensationAttempts int
}

// Saga represents a SAGA orchestration
type Saga struct {
 mu           sync.RWMutex
 ID           string
 Name         string
 Steps        []*Step
 state        SagaState
 stepStates   map[string]StepState
 stepResults  map[string]*StepResult
 stepData     map[string]interface{}
 executionOrder []string
 log          SagaLog
 onStepComplete func(stepID string, result *StepResult)
 onStepFail     func(stepID string, err error)
}

// SagaLog interface for persistence
type SagaLog interface {
 Log(event SagaEvent) error
 GetSagaState(sagaID string) (*SagaStateSnapshot, error)
}

// SagaEvent represents a SAGA log event
type SagaEvent struct {
 SagaID    string
 StepID    string
 EventType string
 Timestamp time.Time
 Data      interface{}
}

// SagaStateSnapshot captures SAGA state for recovery
type SagaStateSnapshot struct {
 SagaID         string
 State          SagaState
 StepStates     map[string]StepState
 ExecutionOrder []string
 StepData       map[string]interface{}
}

// NewSaga creates a new SAGA
func NewSaga(id, name string, log SagaLog) *Saga {
 return &Saga{
  ID:             id,
  Name:           name,
  Steps:          make([]*Step, 0),
  state:          SagaStarted,
  stepStates:     make(map[string]StepState),
  stepResults:    make(map[string]*StepResult),
  stepData:       make(map[string]interface{}),
  executionOrder: make([]string, 0),
  log:            log,
 }
}

// AddStep adds a step to the SAGA
func (s *Saga) AddStep(step *Step) {
 s.mu.Lock()
 defer s.mu.Unlock()
 s.Steps = append(s.Steps, step)
 s.stepStates[step.ID] = StepPending
}

// Execute runs the SAGA
func (s *Saga) Execute(ctx context.Context) error {
 s.mu.Lock()
 s.state = SagaRunning
 s.mu.Unlock()

 // Log start
 s.log.Log(SagaEvent{
  SagaID:    s.ID,
  EventType: "SAGA_STARTED",
  Timestamp: time.Now(),
 })

 // Build execution plan respecting dependencies
 executionPlan := s.buildExecutionPlan()

 // Execute steps
 for _, stepID := range executionPlan {
  step := s.getStep(stepID)
  if step == nil {
   return fmt.Errorf("step %s not found", stepID)
  }

  // Check dependencies
  if !s.dependenciesSatisfied(step) {
   return fmt.Errorf("dependencies not satisfied for step %s", stepID)
  }

  // Execute step
  result := s.executeStep(ctx, step)

  if !result.Success {
   // Log failure
   s.log.Log(SagaEvent{
    SagaID:    s.ID,
    StepID:    step.ID,
    EventType: "STEP_FAILED",
    Timestamp: time.Now(),
    Data:      result.Error,
   })

   if s.onStepFail != nil {
    s.onStepFail(step.ID, result.Error)
   }

   // Compensate
   return s.compensate(ctx)
  }

  // Success
  s.mu.Lock()
  s.stepResults[step.ID] = result
  s.executionOrder = append(s.executionOrder, step.ID)
  s.mu.Unlock()

  if s.onStepComplete != nil {
   s.onStepComplete(step.ID, result)
  }
 }

 // All steps completed
 s.mu.Lock()
 s.state = SagaCompleted
 s.mu.Unlock()

 s.log.Log(SagaEvent{
  SagaID:    s.ID,
  EventType: "SAGA_COMPLETED",
  Timestamp: time.Now(),
 })

 return nil
}

func (s *Saga) buildExecutionPlan() []string {
 s.mu.RLock()
 defer s.mu.RUnlock()

 // Simple topological sort
 plan := make([]string, 0, len(s.Steps))
 completed := make(map[string]bool)

 for len(plan) < len(s.Steps) {
  progress := false
  for _, step := range s.Steps {
   if completed[step.ID] {
    continue
   }

   // Check dependencies
   depsSatisfied := true
   for _, dep := range step.Dependencies {
    if !completed[dep] {
     depsSatisfied = false
     break
    }
   }

   if depsSatisfied {
    plan = append(plan, step.ID)
    completed[step.ID] = true
    progress = true
   }
  }

  if !progress {
   // Cycle detected or missing dependencies
   break
  }
 }

 return plan
}

func (s *Saga) dependenciesSatisfied(step *Step) bool {
 s.mu.RLock()
 defer s.mu.RUnlock()

 for _, dep := range step.Dependencies {
  if s.stepStates[dep] != StepCompleted {
   return false
  }
 }
 return true
}

func (s *Saga) getStep(id string) *Step {
 s.mu.RLock()
 defer s.mu.RUnlock()

 for _, step := range s.Steps {
  if step.ID == id {
   return step
  }
 }
 return nil
}

func (s *Saga) executeStep(ctx context.Context, step *Step) *StepResult {
 s.mu.Lock()
 s.stepStates[step.ID] = StepExecuting
 s.mu.Unlock()

 // Log start
 s.log.Log(SagaEvent{
  SagaID:    s.ID,
  StepID:    step.ID,
  EventType: "STEP_STARTED",
  Timestamp: time.Now(),
 })

 retryPolicy := step.RetryPolicy
 if retryPolicy == nil {
  retryPolicy = &RetryPolicy{MaxAttempts: 1}
 }

 var lastErr error
 backoff := retryPolicy.InitialBackoff

 for attempt := 0; attempt < retryPolicy.MaxAttempts; attempt++ {
  if attempt > 0 {
   time.Sleep(backoff)
   if backoff < retryPolicy.MaxBackoff {
    backoff = time.Duration(float64(backoff) * retryPolicy.BackoffMultiplier)
    if backoff > retryPolicy.MaxBackoff {
     backoff = retryPolicy.MaxBackoff
    }
   }
  }

  // Get input from previous steps
  input := s.getStepInput(step)

  result, err := step.Action(ctx, input)
  if err == nil {
   s.mu.Lock()
   s.stepStates[step.ID] = StepCompleted
   s.stepData[step.ID] = result
   s.mu.Unlock()

   // Log success
   s.log.Log(SagaEvent{
    SagaID:    s.ID,
    StepID:    step.ID,
    EventType: "STEP_COMPLETED",
    Timestamp: time.Now(),
   })

   return &StepResult{Success: true, Data: result}
  }

  lastErr = err
 }

 s.mu.Lock()
 s.stepStates[step.ID] = StepFailed
 s.mu.Unlock()

 return &StepResult{Success: false, Error: lastErr}
}

func (s *Saga) getStepInput(step *Step) interface{} {
 s.mu.RLock()
 defer s.mu.RUnlock()

 // Return data from last completed step
 if len(s.executionOrder) > 0 {
  lastStep := s.executionOrder[len(s.executionOrder)-1]
  return s.stepData[lastStep]
 }
 return nil
}

// compensate runs compensations in reverse order
func (s *Saga) compensate(ctx context.Context) error {
 s.mu.Lock()
 s.state = SagaCompensating
 order := make([]string, len(s.executionOrder))
 copy(order, s.executionOrder)
 s.mu.Unlock()

 s.log.Log(SagaEvent{
  SagaID:    s.ID,
  EventType: "COMPENSATION_STARTED",
  Timestamp: time.Now(),
 })

 // Compensate in reverse order
 for i := len(order) - 1; i >= 0; i-- {
  stepID := order[i]
  step := s.getStep(stepID)
  if step == nil || step.Compensation == nil {
   continue
  }

  // Update state
  s.mu.Lock()
  s.stepStates[stepID] = StepCompensating
  s.mu.Unlock()

  // Execute compensation with retries
  retryPolicy := step.RetryPolicy
  if retryPolicy == nil {
   retryPolicy = &RetryPolicy{MaxCompensationAttempts: 3}
  }

  success := false
  backoff := retryPolicy.InitialBackoff

  for attempt := 0; attempt < retryPolicy.MaxCompensationAttempts; attempt++ {
   if attempt > 0 {
    time.Sleep(backoff)
   }

   err := step.Compensation(ctx, s.stepData[stepID])
   if err == nil {
    success = true
    break
   }

   backoff = time.Duration(float64(backoff) * retryPolicy.BackoffMultiplier)
  }

  if !success {
   // Compensation failed - critical
   s.log.Log(SagaEvent{
    SagaID:    s.ID,
    StepID:    stepID,
    EventType: "COMPENSATION_FAILED",
    Timestamp: time.Now(),
   })
   return fmt.Errorf("compensation failed for step %s", stepID)
  }

  s.mu.Lock()
  s.stepStates[stepID] = StepCompensated
  s.mu.Unlock()

  s.log.Log(SagaEvent{
   SagaID:    s.ID,
   StepID:    stepID,
   EventType: "COMPENSATION_COMPLETED",
   Timestamp: time.Now(),
  })
 }

 s.mu.Lock()
 s.state = SagaAborted
 s.mu.Unlock()

 s.log.Log(SagaEvent{
  SagaID:    s.ID,
  EventType: "SAGA_ABORTED",
  Timestamp: time.Now(),
 })

 return fmt.Errorf("saga aborted")
}

// Orchestrator manages multiple SAGAs
type Orchestrator struct {
 mu      sync.RWMutex
 sagas   map[string]*Saga
 log     SagaLog
}

// NewOrchestrator creates a new orchestrator
func NewOrchestrator(log SagaLog) *Orchestrator {
 return &Orchestrator{
  sagas: make(map[string]*Saga),
  log:   log,
 }
}

// RegisterSaga registers a SAGA definition
func (o *Orchestrator) RegisterSaga(saga *Saga) {
 o.mu.Lock()
 defer o.mu.Unlock()
 o.sagas[saga.ID] = saga
}

// ExecuteSaga executes a SAGA by ID
func (o *Orchestrator) ExecuteSaga(ctx context.Context, sagaID string) error {
 o.mu.RLock()
 saga, exists := o.sagas[sagaID]
 o.mu.RUnlock()

 if !exists {
  return fmt.Errorf("saga %s not found", sagaID)
 }

 return saga.Execute(ctx)
}

// GetSagaState returns the current state of a SAGA
func (o *Orchestrator) GetSagaState(sagaID string) (SagaState, error) {
 o.mu.RLock()
 saga, exists := o.sagas[sagaID]
 o.mu.RUnlock()

 if !exists {
  return SagaStarted, fmt.Errorf("saga %s not found", sagaID)
 }

 saga.mu.RLock()
 defer saga.mu.RUnlock()
 return saga.state, nil
}

// RecoverSaga recovers a SAGA from persistent state
func (o *Orchestrator) RecoverSaga(ctx context.Context, sagaID string) error {
 snapshot, err := o.log.GetSagaState(sagaID)
 if err != nil {
  return err
 }

 saga := &Saga{
  ID:             snapshot.SagaID,
  state:          snapshot.State,
  stepStates:     snapshot.StepStates,
  executionOrder: snapshot.ExecutionOrder,
  stepData:       snapshot.StepData,
  log:            o.log,
 }

 o.RegisterSaga(saga)

 // Resume execution
 if snapshot.State == SagaRunning {
  return saga.Execute(ctx)
 }

 return nil
}

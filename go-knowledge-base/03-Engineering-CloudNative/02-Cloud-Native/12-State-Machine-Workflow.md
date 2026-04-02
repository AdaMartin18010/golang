# 状态机工作流 (State Machine Workflow)

> **分类**: 工程与云原生
> **标签**: #state-machine #workflow #saga

---

## 有限状态机 (FSM)

```go
type State string
const (
    StatePending    State = "pending"
    StateProcessing State = "processing"
    StateCompleted  State = "completed"
    StateFailed     State = "failed"
)

type Event string
const (
    EventStart   Event = "start"
    EventSuccess Event = "success"
    EventFail    Event = "fail"
    EventRetry   Event = "retry"
)

type Transition struct {
    From  State
    Event Event
    To    State
    Action func(ctx context.Context, data interface{}) error
}

type StateMachine struct {
    current     State
    transitions map[State]map[Event]Transition
    mu          sync.RWMutex
}

func NewStateMachine(initial State) *StateMachine {
    return &StateMachine{
        current:     initial,
        transitions: make(map[State]map[Event]Transition),
    }
}

func (sm *StateMachine) AddTransition(t Transition) {
    if sm.transitions[t.From] == nil {
        sm.transitions[t.From] = make(map[Event]Transition)
    }
    sm.transitions[t.From][t.Event] = t
}

func (sm *StateMachine) Trigger(ctx context.Context, event Event, data interface{}) error {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    transitions, ok := sm.transitions[sm.current]
    if !ok {
        return fmt.Errorf("no transitions from state %s", sm.current)
    }

    transition, ok := transitions[event]
    if !ok {
        return fmt.Errorf("event %s not valid from state %s", event, sm.current)
    }

    // 执行动作
    if transition.Action != nil {
        if err := transition.Action(ctx, data); err != nil {
            return fmt.Errorf("action failed: %w", err)
        }
    }

    sm.current = transition.To
    return nil
}

func (sm *StateMachine) Current() State {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    return sm.current
}
```

---

## 订单状态机示例

```go
func NewOrderStateMachine() *StateMachine {
    sm := NewStateMachine(StatePending)

    sm.AddTransition(Transition{
        From:  StatePending,
        Event: EventStart,
        To:    StateProcessing,
        Action: func(ctx context.Context, data interface{}) error {
            order := data.(*Order)
            return reserveInventory(ctx, order)
        },
    })

    sm.AddTransition(Transition{
        From:  StateProcessing,
        Event: EventSuccess,
        To:    StateCompleted,
        Action: func(ctx context.Context, data interface{}) error {
            order := data.(*Order)
            return chargePayment(ctx, order)
        },
    })

    sm.AddTransition(Transition{
        From:  StateProcessing,
        Event: EventFail,
        To:    StateFailed,
        Action: func(ctx context.Context, data interface{}) error {
            order := data.(*Order)
            return releaseInventory(ctx, order)
        },
    })

    sm.AddTransition(Transition{
        From:  StateFailed,
        Event: EventRetry,
        To:    StateProcessing,
        Action: func(ctx context.Context, data interface{}) error {
            // 重试逻辑
            return nil
        },
    })

    return sm
}
```

---

## 持久化状态机

```go
type PersistentStateMachine struct {
    *StateMachine
    store StateStore
    id    string
}

func (psm *PersistentStateMachine) Trigger(ctx context.Context, event Event, data interface{}) error {
    // 先持久化意图
    if err := psm.store.SaveEvent(ctx, psm.id, event, data); err != nil {
        return err
    }

    // 执行状态转换
    if err := psm.StateMachine.Trigger(ctx, event, data); err != nil {
        return err
    }

    // 保存新状态
    return psm.store.SaveState(ctx, psm.id, psm.Current())
}

func (psm *PersistentStateMachine) Recover(ctx context.Context) error {
    state, err := psm.store.GetState(ctx, psm.id)
    if err != nil {
        return err
    }

    psm.StateMachine.current = state

    // 重放未处理的事件
    events, _ := psm.store.GetPendingEvents(ctx, psm.id)
    for _, event := range events {
        psm.StateMachine.Trigger(ctx, event.Type, event.Data)
    }

    return nil
}
```

---

## 并行状态

```go
type ParallelStateMachine struct {
    regions []*StateMachine
}

func (psm *ParallelStateMachine) Trigger(ctx context.Context, event Event, data interface{}) error {
    var wg sync.WaitGroup
    errChan := make(chan error, len(psm.regions))

    for _, region := range psm.regions {
        wg.Add(1)
        go func(sm *StateMachine) {
            defer wg.Done()
            if err := sm.Trigger(ctx, event, data); err != nil {
                errChan <- err
            }
        }(region)
    }

    wg.Wait()
    close(errChan)

    for err := range errChan {
        return err
    }
    return nil
}

func (psm *ParallelStateMachine) AllInState(state State) bool {
    for _, region := range psm.regions {
        if region.Current() != state {
            return false
        }
    }
    return true
}
```

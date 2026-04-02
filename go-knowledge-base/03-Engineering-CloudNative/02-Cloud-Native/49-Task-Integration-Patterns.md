# 任务系统集成模式 (Task System Integration Patterns)

> **分类**: 工程与云原生
> **标签**: #integration #patterns #external-systems

---

## 外部系统集成

```go
// Webhook 集成模式
type WebhookIntegration struct {
    client    *http.Client
    endpoints map[string]WebhookEndpoint
    retries   RetryConfig
}

type WebhookEndpoint struct {
    URL       string
    Secret    string
    Events    []string
    Timeout   time.Duration
    Retries   int
}

func (wi *WebhookIntegration) Notify(ctx context.Context, event TaskEvent) error {
    endpoint := wi.endpoints[event.Type]

    payload, _ := json.Marshal(event)

    req, _ := http.NewRequestWithContext(ctx, "POST", endpoint.URL, bytes.NewReader(payload))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-Webhook-Secret", wi.generateSignature(payload, endpoint.Secret))

    return wi.sendWithRetry(req, endpoint.Retries)
}

func (wi *WebhookIntegration) sendWithRetry(req *http.Request, maxRetries int) error {
    backoff := time.Second

    for i := 0; i <= maxRetries; i++ {
        resp, err := wi.client.Do(req)
        if err != nil {
            if i == maxRetries {
                return err
            }
            time.Sleep(backoff)
            backoff *= 2
            continue
        }
        defer resp.Body.Close()

        if resp.StatusCode >= 200 && resp.StatusCode < 300 {
            return nil
        }

        if resp.StatusCode >= 500 || resp.StatusCode == 429 {
            if i == maxRetries {
                return fmt.Errorf("webhook failed with status: %d", resp.StatusCode)
            }
            time.Sleep(backoff)
            backoff *= 2
            continue
        }

        return fmt.Errorf("webhook failed with status: %d", resp.StatusCode)
    }

    return nil
}
```

---

## 适配器模式

```go
// 适配不同任务系统
type TaskSystemAdapter interface {
    Submit(ctx context.Context, task ExternalTask) error
    GetStatus(ctx context.Context, taskID string) (TaskStatus, error)
    Cancel(ctx context.Context, taskID string) error
}

// Celery 适配器
type CeleryAdapter struct {
    broker CeleryBroker
}

func (ca *CeleryAdapter) Submit(ctx context.Context, task ExternalTask) error {
    celeryTask := CeleryTask{
        Task: task.Name,
        Args: task.Payload,
        ID:   task.ID,
    }

    return ca.broker.SendTask(ctx, celeryTask)
}

// AWS Step Functions 适配器
type StepFunctionsAdapter struct {
    client *sfn.Client
}

func (sfa *StepFunctionsAdapter) Submit(ctx context.Context, task ExternalTask) error {
    input, _ := json.Marshal(task.Payload)

    _, err := sfa.client.StartExecution(ctx, &sfn.StartExecutionInput{
        StateMachineArn: aws.String(sfa.getStateMachineARN(task.Type)),
        Name:            aws.String(task.ID),
        Input:           aws.String(string(input)),
    })

    return err
}

// 统一接口
func SubmitToExternalSystem(ctx context.Context, system string, task ExternalTask) error {
    adapters := map[string]TaskSystemAdapter{
        "celery":         &CeleryAdapter{},
        "stepfunctions":  &StepFunctionsAdapter{},
        "temporal":       &TemporalAdapter{},
        "argo":           &ArgoAdapter{},
    }

    adapter, ok := adapters[system]
    if !ok {
        return fmt.Errorf("unknown external system: %s", system)
    }

    return adapter.Submit(ctx, task)
}
```

---

## Saga 模式

```go
// 分布式事务 Saga 实现
type Saga struct {
    steps []SagaStep
    state SagaState
}

type SagaStep struct {
    Name      string
    Action    func(ctx context.Context) error
    Compensate func(ctx context.Context) error
}

func (s *Saga) Execute(ctx context.Context) error {
    for i, step := range s.steps {
        if err := step.Action(ctx); err != nil {
            // 执行补偿
            s.compensate(ctx, i)
            return fmt.Errorf("saga failed at step %s: %w", step.Name, err)
        }
        s.state.CompletedSteps = append(s.state.CompletedSteps, step.Name)
    }
    return nil
}

func (s *Saga) compensate(ctx context.Context, failedIndex int) {
    // 逆序执行补偿
    for i := failedIndex - 1; i >= 0; i-- {
        step := s.steps[i]
        if step.Compensate != nil {
            if err := step.Compensate(ctx); err != nil {
                // 补偿失败，需要人工介入
                log.Printf("Compensation failed for step %s: %v", step.Name, err)
            }
        }
    }
}

// 订单创建 Saga
func CreateOrderSaga(order Order) *Saga {
    return &Saga{
        steps: []SagaStep{
            {
                Name: "reserve_inventory",
                Action: func(ctx context.Context) error {
                    return inventoryService.Reserve(ctx, order.Items)
                },
                Compensate: func(ctx context.Context) error {
                    return inventoryService.Release(ctx, order.Items)
                },
            },
            {
                Name: "process_payment",
                Action: func(ctx context.Context) error {
                    return paymentService.Charge(ctx, order.Payment)
                },
                Compensate: func(ctx context.Context) error {
                    return paymentService.Refund(ctx, order.Payment)
                },
            },
            {
                Name: "create_shipment",
                Action: func(ctx context.Context) error {
                    return shippingService.Create(ctx, order)
                },
                Compensate: func(ctx context.Context) error {
                    return shippingService.Cancel(ctx, order)
                },
            },
        },
    }
}
```

---

## 事件驱动集成

```go
// 事件总线集成
type EventDrivenIntegration struct {
    eventBus EventBus
    handlers map[string][]EventHandler
}

func (edi *EventDrivenIntegration) Subscribe(eventType string, handler EventHandler) {
    edi.handlers[eventType] = append(edi.handlers[eventType], handler)
}

func (edi *EventDrivenIntegration) Publish(ctx context.Context, event DomainEvent) error {
    return edi.eventBus.Publish(ctx, EventEnvelope{
        Type:      event.Type(),
        Payload:   event,
        Timestamp: time.Now(),
    })
}

func (edi *EventDrivenIntegration) Start() {
    for eventType, handlers := range edi.handlers {
        go edi.consume(eventType, handlers)
    }
}

func (edi *EventDrivenIntegration) consume(eventType string, handlers []EventHandler) {
    messages := edi.eventBus.Subscribe(eventType)

    for msg := range messages {
        for _, handler := range handlers {
            go func(h EventHandler) {
                if err := h.Handle(context.Background(), msg); err != nil {
                    log.Printf("Handler error: %v", err)
                }
            }(handler)
        }
    }
}
```

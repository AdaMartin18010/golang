# 任务批量处理 (Task Batch Processing)

> **分类**: 工程与云原生
> **标签**: #batch-processing #bulk-operations #performance

---

## 批量执行器

```go
type BatchExecutor struct {
    batchSize     int
    flushInterval time.Duration
    buffer        []Task
    mu            sync.Mutex
    processor     BatchProcessor
    ticker        *time.Ticker
}

type BatchProcessor interface {
    ProcessBatch(ctx context.Context, tasks []Task) []Result
}

func NewBatchExecutor(size int, interval time.Duration, processor BatchProcessor) *BatchExecutor {
    be := &BatchExecutor{
        batchSize:     size,
        flushInterval: interval,
        buffer:        make([]Task, 0, size),
        processor:     processor,
        ticker:        time.NewTicker(interval),
    }

    go be.flushLoop()
    return be
}

func (be *BatchExecutor) Submit(task Task) {
    be.mu.Lock()
    be.buffer = append(be.buffer, task)
    shouldFlush := len(be.buffer) >= be.batchSize
    be.mu.Unlock()

    if shouldFlush {
        be.Flush()
    }
}

func (be *BatchExecutor) flushLoop() {
    for range be.ticker.C {
        be.Flush()
    }
}

func (be *BatchExecutor) Flush() {
    be.mu.Lock()
    if len(be.buffer) == 0 {
        be.mu.Unlock()
        return
    }

    batch := make([]Task, len(be.buffer))
    copy(batch, be.buffer)
    be.buffer = be.buffer[:0]
    be.mu.Unlock()

    // 批量处理
    ctx := context.Background()
    results := be.processor.ProcessBatch(ctx, batch)

    // 回调结果
    for i, result := range results {
        be.notifyResult(batch[i], result)
    }
}

func (be *BatchExecutor) Stop() {
    be.ticker.Stop()
    be.Flush()  // 刷新剩余
}
```

---

## 微批量处理

```go
type MicroBatcher struct {
    maxDelay    time.Duration
    maxSize     int
    buffer      chan Task
    processor   func([]Task) error
}

func (mb *MicroBatcher) Start() {
    go mb.processLoop()
}

func (mb *MicroBatcher) processLoop() {
    var batch []Task
    timer := time.NewTimer(mb.maxDelay)

    for {
        select {
        case task := <-mb.buffer:
            batch = append(batch, task)
            if len(batch) >= mb.maxSize {
                mb.process(batch)
                batch = nil
                timer.Reset(mb.maxDelay)
            }

        case <-timer.C:
            if len(batch) > 0 {
                mb.process(batch)
                batch = nil
            }
            timer.Reset(mb.maxDelay)
        }
    }
}

func (mb *MicroBatcher) process(batch []Task) {
    if err := mb.processor(batch); err != nil {
        // 失败处理：逐个重试
        for _, task := range batch {
            go mb.retryTask(task)
        }
    }
}
```

---

## 批量数据库操作

```go
func BatchInsertUsers(ctx context.Context, db *sql.DB, users []User) error {
    if len(users) == 0 {
        return nil
    }

    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    stmt, err := tx.PrepareContext(ctx,
        "INSERT INTO users (name, email) VALUES (?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, user := range users {
        if _, err := stmt.ExecContext(ctx, user.Name, user.Email); err != nil {
            return err
        }
    }

    return tx.Commit()
}

// 使用 COPY (PostgreSQL)
func BatchCopyUsers(ctx context.Context, conn *pgx.Conn, users []User) error {
    copyCount, err := conn.CopyFrom(ctx,
        pgx.Identifier{"users"},
        []string{"name", "email"},
        pgx.CopyFromSlice(len(users), func(i int) ([]interface{}, error) {
            return []interface{}{users[i].Name, users[i].Email}, nil
        }),
    )

    if err != nil {
        return err
    }

    if int(copyCount) != len(users) {
        return fmt.Errorf("expected %d rows, got %d", len(users), copyCount)
    }

    return nil
}
```

---

## 批量 API 调用

```go
type BatchAPIClient struct {
    client      *http.Client
    endpoint    string
    maxBatchSize int
}

func (bac *BatchAPIClient) BatchRequest(ctx context.Context, requests []APIRequest) ([]APIResponse, error) {
    // 分批发送
    var allResponses []APIResponse

    for i := 0; i < len(requests); i += bac.maxBatchSize {
        end := i + bac.maxBatchSize
        if end > len(requests) {
            end = len(requests)
        }

        batch := requests[i:end]
        responses, err := bac.sendBatch(ctx, batch)
        if err != nil {
            return nil, err
        }

        allResponses = append(allResponses, responses...)
    }

    return allResponses, nil
}

func (bac *BatchAPIClient) sendBatch(ctx context.Context, requests []APIRequest) ([]APIResponse, error) {
    payload := BatchRequest{Requests: requests}

    data, _ := json.Marshal(payload)
    req, _ := http.NewRequestWithContext(ctx, "POST", bac.endpoint, bytes.NewReader(data))
    req.Header.Set("Content-Type", "application/json")

    resp, err := bac.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var batchResp BatchResponse
    if err := json.NewDecoder(resp.Body).Decode(&batchResp); err != nil {
        return nil, err
    }

    return batchResp.Responses, nil
}
```

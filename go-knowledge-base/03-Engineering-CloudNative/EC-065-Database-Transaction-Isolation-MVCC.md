# 数据库事务隔离与 MVCC (Database Transaction Isolation & MVCC)

> **分类**: 工程与云原生
> **标签**: #database #transaction #mvcc #isolation-level #acid
> **参考**: PostgreSQL MVCC, MySQL InnoDB, ACID Theory

---

## ACID 属性

```
┌─────────────────────────────────────────────────────────────────┐
│                        ACID Properties                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  A - Atomicity (原子性)                                          │
│     事务是不可分割的工作单位，要么全部完成，要么全部不完成           │
│                                                                 │
│     BEGIN;                                                      │
│       UPDATE account SET balance = balance - 100 WHERE id = 1;  │
│       UPDATE account SET balance = balance + 100 WHERE id = 2;  │
│     COMMIT;  -- 要么都成功                                       │
│     ROLLBACK; -- 要么都失败                                      │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│  C - Consistency (一致性)                                        │
│     事务执行前后，数据库必须处于一致状态                            │
│                                                                  │
│     约束：balance >= 0                                           │
│     事务必须维护这个约束，不能产生负余额                            │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│  I - Isolation (隔离性)                                          │
│     并发事务之间相互隔离，互不干扰                                 │
│                                                                 │
│     通过 MVCC 或锁机制实现                                        │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│  D - Durability (持久性)                                         │
│     一旦提交，数据永久保存，即使系统故障                           │
│                                                                 │
│     通过 WAL (Write-Ahead Logging) 实现                          │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 隔离级别与异常现象

```
┌─────────────────────────────────────────────────────────────────┐
│              Isolation Levels vs Read Phenomena                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Phenomenon          │ Description                              │
├──────────────────────┼──────────────────────────────────────────┤
│  Dirty Read          │ 读取到其他事务未提交的数据                 │
│  Non-Repeatable Read │ 同一事务内，两次读取同一行数据结果不同      │
│  Phantom Read        │ 同一事务内，两次查询结果集的行数不同        │
│  Lost Update         │ 两个事务同时更新同一行，一个被覆盖          │
│  Write Skew          │ 两个事务基于快照读取，同时写入违反约束      │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│  Isolation Level     │ Dirty │ Non-Repeat │ Phantom │ Lost Up  │
├──────────────────────┼───────┼────────────┼─────────┼──────────┤
│  Read Uncommitted    │  ✓    │     ✗      │    ✗    │    ✗    │
│  Read Committed      │  ✗    │     ✓      │    ✓    │    ✗    │
│  Repeatable Read     │  ✗    │     ✗      │    ✓*   │    ✗    │
│  Serializable        │  ✗    │     ✗      │    ✗    │    ✗    │
│                                                                 │
│  * PostgreSQL 的 Repeatable Read 也防止 Phantom Read             │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## MVCC 实现原理

### PostgreSQL MVCC

```go
// PostgreSQL 行格式（简化）
type HeapTupleHeader struct {
    // 事务 ID
    t_xmin TransactionId  // 插入事务 ID
    t_xmax TransactionId  // 删除事务 ID（0 表示未删除）

    // 命令 ID（同一事务内的版本）
    t_cid CommandId

    // 指向旧版本的指针（TOAST 或 HOT）
    t_ctid ItemPointerData  // 当前版本或指向新版本

    // 实际数据
    data []byte
}

// 可见性判断逻辑
func IsVisible(tuple *HeapTupleHeader, snapshot *Snapshot) bool {
    // 规则 1：插入事务未提交，不可见
    if !TransactionIdIsCommitted(tuple.t_xmin) {
        // 如果是自己插入的，可见
        if tuple.t_xmin != snapshot.xmin {
            return false
        }
    }

    // 规则 2：删除事务已提交，不可见
    if TransactionIdIsCommitted(tuple.t_xmax) {
        return false
    }

    // 规则 3：插入事务在快照开始前提交
    if !TransactionIdPrecedes(tuple.t_xmin, snapshot.xmax) {
        return false
    }

    return true
}
```

### 快照机制

```go
// PostgreSQL Snapshot
type Snapshot struct {
    xmin        TransactionId  // 最老的活跃事务 ID
    xmax        TransactionId  // 最新已分配事务 ID + 1
    xip         []TransactionId // 活跃事务列表（snapshot 时正在运行的）
    snapshotcsn CommitSeqNumber // 提交序列号（用于可重复读）
}

// 创建快照
func GetSnapshot() *Snapshot {
    snap := &Snapshot{}

    // 1. 获取当前事务 ID
    snap.xmin = GetOldestActiveTransactionId()
    snap.xmax = GetNextTransactionId()

    // 2. 复制活跃事务列表
    snap.xip = GetActiveTransactionIds()

    return snap
}

// 可见性判断（简化）
func (s *Snapshot) IsVisible(xmin, xmax TransactionId) bool {
    // 插入者已提交
    if TransactionIdIsCommitted(xmin) {
        // 且不是在快照开始后提交的
        if TransactionIdPrecedes(xmin, s.xmax) {
            // 且不是在快照创建时活跃的
            if !contains(s.xip, xmin) {
                // 检查是否被删除
                if xmax == 0 || !TransactionIdIsCommitted(xmax) {
                    return true
                }
            }
        }
    }
    return false
}
```

---

## MySQL InnoDB MVCC

### 行结构

```go
// InnoDB 行格式（简化）
type InnoDBRow struct {
    // 系统列
    DB_TRX_ID  uint64  // 最后修改事务 ID（6 字节）
    DB_ROLL_PTR uint64 // 回滚指针（7 字节）

    // 实际数据
    columns []Column
}

// Undo Log 记录
type UndoLogRecord struct {
    TRX_ID      uint64  // 事务 ID
    ROLL_PTR    uint64  // 指向上一个版本
    OLD_VALUES  []byte  // 旧值（用于回滚和构建旧版本）
}
```

### Read View

```go
// InnoDB Read View
type ReadView struct {
    m_ids           []uint64  // 创建时活跃事务 ID 列表
    m_up_limit_id   uint64    // 活跃事务最小 ID（m_ids 最小值）
    m_low_limit_id  uint64    // 下一个分配的事务 ID
    m_creator_trx_id uint64   // 创建者事务 ID
}

// 可见性判断
func (rv *ReadView) IsVisible(trxID uint64) bool {
    // 1. 如果 trx_id 等于 creator_trx_id，可见（自己修改的）
    if trxID == rv.m_creator_trx_id {
        return true
    }

    // 2. 如果 trx_id < up_limit_id，可见（已提交）
    if trxID < rv.m_up_limit_id {
        return true
    }

    // 3. 如果 trx_id >= low_limit_id，不可见（快照后启动的）
    if trxID >= rv.m_low_limit_id {
        return false
    }

    // 4. 检查是否在 m_ids 中（活跃事务）
    if contains(rv.m_ids, trxID) {
        return false  // 未提交，不可见
    }

    return true  // 已提交，可见
}
```

---

## Go 实现简单 MVCC 存储

```go
package mvcc

import (
    "sync"
    "sync/atomic"
    "time"
)

// 版本化值
type VersionedValue struct {
    Value      []byte
    Version    uint64
    TxID       uint64
    Timestamp  time.Time
    Deleted    bool
}

// MVCC 存储
type MVCCStore struct {
    mu      sync.RWMutex
    data    map[string][]*VersionedValue  // key -> versions
    txCounter uint64

    // 活跃事务
    activeTx map[uint64]*Transaction
}

type Transaction struct {
    ID        uint64
    StartTime time.Time
    ReadTS    uint64
    WriteSet  map[string]*VersionedValue
    mu        sync.Mutex
}

func NewMVCCStore() *MVCCStore {
    return &MVCCStore{
        data:     make(map[string][]*VersionedValue),
        activeTx: make(map[uint64]*Transaction),
    }
}

// 开始事务
func (s *MVCCStore) Begin() *Transaction {
    txID := atomic.AddUint64(&s.txCounter, 1)
    tx := &Transaction{
        ID:       txID,
        ReadTS:   txID,  // 简化：使用事务 ID 作为时间戳
        WriteSet: make(map[string]*VersionedValue),
    }

    s.mu.Lock()
    s.activeTx[txID] = tx
    s.mu.Unlock()

    return tx
}

// 读取（快照读）
func (s *MVCCStore) Get(tx *Transaction, key string) ([]byte, bool) {
    s.mu.RLock()
    versions, exists := s.data[key]
    s.mu.RUnlock()

    if !exists {
        return nil, false
    }

    // 查找对当前事务可见的最新版本
    for i := len(versions) - 1; i >= 0; i-- {
        v := versions[i]
        // 版本在事务开始前创建，且未删除
        if v.TxID <= tx.ReadTS && !v.Deleted {
            return v.Value, true
        }
    }

    return nil, false
}

// 写入
func (s *MVCCStore) Put(tx *Transaction, key string, value []byte) {
    tx.mu.Lock()
    defer tx.mu.Unlock()

    newVersion := &VersionedValue{
        Value:     value,
        Version:   tx.ReadTS,
        TxID:      tx.ID,
        Timestamp: time.Now(),
    }

    tx.WriteSet[key] = newVersion
}

// 删除
func (s *MVCCStore) Delete(tx *Transaction, key string) {
    tx.mu.Lock()
    defer tx.mu.Unlock()

    tx.WriteSet[key] = &VersionedValue{
        TxID:    tx.ID,
        Deleted: true,
    }
}

// 提交事务
func (s *MVCCStore) Commit(tx *Transaction) error {
    tx.mu.Lock()
    defer tx.mu.Unlock()

    s.mu.Lock()
    defer s.mu.Unlock()

    // 检查写写冲突（简化：检查最后修改时间）
    for key, newVer := range tx.WriteSet {
        versions := s.data[key]
        if len(versions) > 0 {
            lastVer := versions[len(versions)-1]
            // 如果在我们读取后，有其他事务修改了该键
            if lastVer.TxID > tx.ReadTS && lastVer.TxID != tx.ID {
                return ErrWriteConflict
            }
        }

        // 追加新版本
        s.data[key] = append(s.data[key], newVer)
    }

    delete(s.activeTx, tx.ID)
    return nil
}

// 回滚事务
func (s *MVCCStore) Rollback(tx *Transaction) {
    s.mu.Lock()
    delete(s.activeTx, tx.ID)
    s.mu.Unlock()
}

var ErrWriteConflict = errors.New("write conflict")
```

---

## 两阶段锁 (2PL) 实现

```go
package twophase

import (
    "sync"
)

// 锁类型
type LockType int

const (
    LockNone LockType = iota
    LockShared       // 读锁
    LockExclusive    // 写锁
)

// 锁管理器
type LockManager struct {
    mu     sync.Mutex
    locks  map[string]*Lock  // 资源 -> 锁
    graph  map[uint64]map[uint64]bool  // 等待图（用于死锁检测）
}

type Lock struct {
    Resource    string
    LockType    LockType
    Holder      uint64           // 持有者
    WaitQueue   []LockRequest    // 等待队列
}

type LockRequest struct {
    TxID     uint64
    Type     LockType
    Granted  chan bool
}

func NewLockManager() *LockManager {
    return &LockManager{
        locks: make(map[string]*Lock),
        graph: make(map[uint64]map[uint64]bool),
    }
}

// 获取锁
func (lm *LockManager) Acquire(txID uint64, resource string, lockType LockType) bool {
    lm.mu.Lock()

    lock, exists := lm.locks[resource]
    if !exists {
        // 无锁，直接授予
        lm.locks[resource] = &Lock{
            Resource: resource,
            LockType: lockType,
            Holder:   txID,
        }
        lm.mu.Unlock()
        return true
    }

    // 检查兼容性
    if lm.isCompatible(lock, txID, lockType) {
        // 兼容，授予锁
        if lockType == LockExclusive {
            lock.LockType = LockExclusive
        }
        lm.mu.Unlock()
        return true
    }

    // 不兼容，加入等待队列
    granted := make(chan bool, 1)
    req := LockRequest{
        TxID:    txID,
        Type:    lockType,
        Granted: granted,
    }
    lock.WaitQueue = append(lock.WaitQueue, req)

    // 更新等待图
    lm.addEdge(txID, lock.Holder)

    // 检查死锁
    if lm.detectDeadlock(txID) {
        // 死锁，移除请求
        lm.removeFromQueue(lock, txID)
        lm.removeEdge(txID, lock.Holder)
        lm.mu.Unlock()
        return false
    }

    lm.mu.Unlock()

    // 等待锁授予
    return <-granted
}

// 检查锁兼容性
func (lm *LockManager) isCompatible(lock *Lock, txID uint64, reqType LockType) bool {
    // 同一事务，总是兼容（锁升级）
    if lock.Holder == txID {
        return true
    }

    // 检查现有锁和请求锁的兼容性
    if lock.LockType == LockShared && reqType == LockShared {
        return true  // 多个读锁兼容
    }

    return false  // 其他情况不兼容
}

// 释放锁
func (lm *LockManager) Release(txID uint64, resource string) {
    lm.mu.Lock()
    defer lm.mu.Unlock()

    lock, exists := lm.locks[resource]
    if !exists || lock.Holder != txID {
        return
    }

    // 处理等待队列
    for len(lock.WaitQueue) > 0 {
        req := lock.WaitQueue[0]
        lock.WaitQueue = lock.WaitQueue[1:]

        if lm.isCompatible(lock, req.TxID, req.Type) {
            // 授予锁
            lock.Holder = req.TxID
            lock.LockType = req.Type
            lm.removeEdge(req.TxID, txID)
            req.Granted <- true
            return
        }
        // 不兼容，继续等待
        lock.WaitQueue = append(lock.WaitQueue, req)
    }

    // 无等待者，删除锁
    delete(lm.locks, resource)
}

// 死锁检测（深度优先搜索）
func (lm *LockManager) detectDeadlock(start uint64) bool {
    visited := make(map[uint64]bool)
    recStack := make(map[uint64]bool)

    var dfs func(uint64) bool
    dfs = func(node uint64) bool {
        visited[node] = true
        recStack[node] = true

        for neighbor := range lm.graph[node] {
            if !visited[neighbor] {
                if dfs(neighbor) {
                    return true
                }
            } else if recStack[neighbor] {
                return true  // 发现回边，有环
            }
        }

        recStack[node] = false
        return false
    }

    return dfs(start)
}

func (lm *LockManager) addEdge(from, to uint64) {
    if lm.graph[from] == nil {
        lm.graph[from] = make(map[uint64]bool)
    }
    lm.graph[from][to] = true
}

func (lm *LockManager) removeEdge(from, to uint64) {
    if lm.graph[from] != nil {
        delete(lm.graph[from], to)
    }
}
```

---

## 事务管理器

```go
package txmanager

import (
    "context"
    "sync"
)

// 事务状态
type TxState int

const (
    TxActive TxState = iota
    TxCommitted
    TxAborted
)

// 事务管理器
type TransactionManager struct {
    mu          sync.RWMutex
    activeTx    map[uint64]*Tx
    committedTx map[uint64]*Tx

    storage     *MVCCStore
    lm          *LockManager

    txCounter   uint64
}

type Tx struct {
    ID       uint64
    State    TxState
    StartTS  uint64
    CommitTS uint64

    ReadSet  map[string]uint64   // key -> version
    WriteSet map[string][]byte

    mu       sync.Mutex
}

func (tm *TransactionManager) Begin() *Tx {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    tm.txCounter++
    tx := &Tx{
        ID:       tm.txCounter,
        State:    TxActive,
        StartTS:  tm.txCounter,
        ReadSet:  make(map[string]uint64),
        WriteSet: make(map[string][]byte),
    }

    tm.activeTx[tx.ID] = tx
    return tx
}

// 读取（带锁）
func (tm *TransactionManager) Read(tx *Tx, key string) ([]byte, error) {
    tx.mu.Lock()
    defer tx.mu.Unlock()

    if tx.State != TxActive {
        return nil, ErrTxNotActive
    }

    // 获取读锁
    if !tm.lm.Acquire(tx.ID, key, LockShared) {
        return nil, ErrLockAcquisitionFailed
    }

    // 读取值
    value, exists := tm.storage.Get(&Transaction{ID: tx.ID, ReadTS: tx.StartTS}, key)
    if exists {
        tx.ReadSet[key] = tx.StartTS  // 记录读取版本
    }

    return value, nil
}

// 写入（带锁）
func (tm *TransactionManager) Write(tx *Tx, key string, value []byte) error {
    tx.mu.Lock()
    defer tx.mu.Unlock()

    if tx.State != TxActive {
        return ErrTxNotActive
    }

    // 获取写锁
    if !tm.lm.Acquire(tx.ID, key, LockExclusive) {
        return ErrLockAcquisitionFailed
    }

    tx.WriteSet[key] = value
    return nil
}

// 提交（2PC 简化版）
func (tm *TransactionManager) Commit(tx *Tx) error {
    tx.mu.Lock()
    defer tx.mu.Unlock()

    if tx.State != TxActive {
        return ErrTxNotActive
    }

    // Phase 1: 验证（乐观并发控制）
    for key := range tx.WriteSet {
        // 检查是否有冲突
        // 简化实现
    }

    // Phase 2: 提交
    tx.State = TxCommitted
    tx.CommitTS = tm.txCounter + 1

    // 写入存储
    for key, value := range tx.WriteSet {
        tm.storage.Put(&Transaction{ID: tx.ID, ReadTS: tx.StartTS}, key, value)
    }

    // 释放锁
    for key := range tx.ReadSet {
        tm.lm.Release(tx.ID, key)
    }
    for key := range tx.WriteSet {
        tm.lm.Release(tx.ID, key)
    }

    // 从事务列表移除
    tm.mu.Lock()
    delete(tm.activeTx, tx.ID)
    tm.committedTx[tx.ID] = tx
    tm.mu.Unlock()

    return nil
}

// 回滚
func (tm *TransactionManager) Abort(tx *Tx) {
    tx.mu.Lock()
    defer tx.mu.Unlock()

    if tx.State != TxActive {
        return
    }

    tx.State = TxAborted

    // 释放所有锁
    for key := range tx.ReadSet {
        tm.lm.Release(tx.ID, key)
    }
    for key := range tx.WriteSet {
        tm.lm.Release(tx.ID, key)
    }

    tm.mu.Lock()
    delete(tm.activeTx, tx.ID)
    tm.mu.Unlock()
}

var (
    ErrTxNotActive          = errors.New("transaction not active")
    ErrLockAcquisitionFailed = errors.New("lock acquisition failed")
)
```

---

## 性能优化策略

```
┌─────────────────────────────────────────────────────────────────┐
│                   MVCC Performance Tuning                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1. 版本清理（Vacuum/GC）                                         │
│     - PostgreSQL: AUTOVACUUM, 定期清理死元组                      │
│     - MySQL: Purge thread, 清理 undo log                         │
│     - 手动触发：VACUUM ANALYZE                                   │
│                                                                  │
│  2. 快照管理                                                      │
│     - 避免长时间事务（长事务阻止版本清理）                         │
│     - 使用合理的隔离级别（RC 比 RR 清理更及时）                    │
│                                                                  │
│  3. 索引优化                                                      │
│     - 索引也遵循 MVCC 规则                                       │
│     - 定期 REINDEX 减少膨胀                                       │
│                                                                  │
│  4. 存储优化                                                      │
│     - 使用表分区减少单表数据量                                    │
│     - 大字段使用 TOAST/外部存储                                   │
│                                                                  │
│  5. 并发控制                                                      │
│     - 短事务优先                                                  │
│     - 热点数据使用应用层缓存                                      │
│     - 批量操作减少事务数量                                        │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

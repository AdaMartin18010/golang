# etcd 分布式协调实现

> **分类**: 工程与云原生
> **标签**: #etcd #distributed-systems #coordination #consensus
> **参考**: etcd v3 API, Raft Consensus Algorithm

---

## etcd 架构概览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      etcd Distributed Key-Value Store                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                         Client Layer                                │    │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │    │
│  │  │  gRPC API   │  │  HTTP API   │  │   Watch     │  │  Lease API  │ │    │
│  │  │             │  │  (Legacy)   │  │   Stream    │  │             │ │    │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                    │                                        │
│  ┌─────────────────────────────────▼─────────────────────────────────────┐  │
│  │                    etcdserver (Raft Node)                             │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │   Raft      │  │   WAL       │  │  Snapshot   │  │   Backend   │ │   │
│  │  │  Consensus  │  │  (Log)      │  │  (Periodic) │  │  (BoltDB)   │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                         │
│  ┌─────────────────────────────────▼─────────────────────────────────────┐   │
│  │                    Storage Layer                                     │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │   BoltDB    │  │   bbolt     │  │   Index     │  │   Key       │ │   │
│  │  │  (MVCC)     │  │  (Backend)  │  │  (B-tree)   │  │  Revision   │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 租约机制 (Lease)

```go
package etcd

import (
    "context"
    "sync"
    "time"

    clientv3 "go.etcd.io/etcd/client/v3"
)

// LeaseManager 租约管理器
type LeaseManager struct {
    client *clientv3.Client

    // 活跃租约
    leases map[clientv3.LeaseID]*ManagedLease
    mu     sync.RWMutex

    // 续约间隔
    renewInterval time.Duration
}

type ManagedLease struct {
    ID       clientv3.LeaseID
    TTL      int64
    Keys     []string // 关联的 key
    RenewAt  time.Time
    Cancel   context.CancelFunc
}

func NewLeaseManager(client *clientv3.Client) *LeaseManager {
    return &LeaseManager{
        client:        client,
        leases:        make(map[clientv3.LeaseID]*ManagedLease),
        renewInterval: 5 * time.Second,
    }
}

// Grant 创建租约
func (m *LeaseManager) Grant(ctx context.Context, ttl int64) (*ManagedLease, error) {
    resp, err := m.client.Grant(ctx, ttl)
    if err != nil {
        return nil, err
    }

    lease := &ManagedLease{
        ID:      resp.ID,
        TTL:     ttl,
        RenewAt: time.Now().Add(time.Duration(ttl) * time.Second / 3),
    }

    m.mu.Lock()
    m.leases[resp.ID] = lease
    m.mu.Unlock()

    // 启动自动续约
    ctx, cancel := context.WithCancel(context.Background())
    lease.Cancel = cancel
    go m.keepAlive(ctx, lease)

    return lease, nil
}

// keepAlive 自动续约
func (m *LeaseManager) keepAlive(ctx context.Context, lease *ManagedLease) {
    ticker := time.NewTicker(m.renewInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            _, err := m.client.KeepAliveOnce(ctx, lease.ID)
            if err != nil {
                // 续约失败，标记租约过期
                m.revoke(ctx, lease.ID)
                return
            }
            lease.RenewAt = time.Now()
        }
    }
}

// AttachKey 将 key 关联到租约
func (m *ManagedLease) AttachKey(key string) {
    m.Keys = append(m.Keys, key)
}

// Revoke 撤销租约
func (m *LeaseManager) Revoke(ctx context.Context, leaseID clientv3.LeaseID) error {
    m.mu.Lock()
    lease, ok := m.leases[leaseID]
    if ok {
        delete(m.leases, leaseID)
    }
    m.mu.Unlock()

    if lease != nil && lease.Cancel != nil {
        lease.Cancel()
    }

    _, err := m.client.Revoke(ctx, leaseID)
    return err
}

func (m *LeaseManager) revoke(ctx context.Context, leaseID clientv3.LeaseID) {
    m.mu.Lock()
    delete(m.leases, leaseID)
    m.mu.Unlock()
}

// GetLeaseInfo 获取租约信息
func (m *LeaseManager) GetLeaseInfo(ctx context.Context, leaseID clientv3.LeaseID) (*clientv3.LeaseTimeToLiveResponse, error) {
    resp, err := m.client.TimeToLive(ctx, leaseID, clientv3.WithAttachedKeys())
    if err != nil {
        return nil, err
    }
    return resp, nil
}

// LeaseWithKeys 创建租约并自动关联 key
func (m *LeaseManager) LeaseWithKeys(ctx context.Context, ttl int64, keys []string, values []string) (*ManagedLease, error) {
    if len(keys) != len(values) {
        return nil, fmt.Errorf("keys and values length mismatch")
    }

    // 创建租约
    lease, err := m.Grant(ctx, ttl)
    if err != nil {
        return nil, err
    }

    // 批量写入带租约的 key
    ops := make([]clientv3.Op, len(keys))
    for i, key := range keys {
        ops[i] = clientv3.OpPut(key, values[i], clientv3.WithLease(lease.ID))
    }

    _, err = m.client.Txn(ctx).Then(ops...).Commit()
    if err != nil {
        m.Revoke(ctx, lease.ID)
        return nil, err
    }

    // 记录关联的 key
    for _, key := range keys {
        lease.AttachKey(key)
    }

    return lease, nil
}
```

---

## 分布式锁实现

```go
package etcd

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "sync"
    "time"

    clientv3 "go.etcd.io/etcd/client/v3"
    "go.etcd.io/etcd/client/v3/concurrency"
)

// DistributedLock 基于 etcd 的分布式锁
type DistributedLock struct {
    client *clientv3.Client
    key    string
    value  string
    lease  clientv3.LeaseID
    ttl    int64

    mu     sync.Mutex
    locked bool
}

// NewDistributedLock 创建分布式锁
func NewDistributedLock(client *clientv3.Client, key string, ttl int64) (*DistributedLock, error) {
    // 生成唯一值
    b := make([]byte, 8)
    rand.Read(b)
    value := hex.EncodeToString(b)

    return &DistributedLock{
        client: client,
        key:    key,
        value:  value,
        ttl:    ttl,
        locked: false,
    }, nil
}

// Lock 获取锁（阻塞）
func (l *DistributedLock) Lock(ctx context.Context) error {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.locked {
        return fmt.Errorf("already locked")
    }

    // 创建租约
    lease, err := l.client.Grant(ctx, l.ttl)
    if err != nil {
        return err
    }
    l.lease = lease.ID

    // 尝试获取锁：使用事务确保原子性
    for {
        txn := l.client.Txn(ctx).
            If(clientv3.Compare(clientv3.CreateRevision(l.key), "=", 0)).
            Then(clientv3.OpPut(l.key, l.value, clientv3.WithLease(l.lease))).
            Else(clientv3.OpGet(l.key))

        resp, err := txn.Commit()
        if err != nil {
            l.client.Revoke(ctx, l.lease)
            return err
        }

        if resp.Succeeded {
            // 获取成功，启动续约
            l.locked = true
            go l.keepAlive(context.Background())
            return nil
        }

        // 锁已被占用，监听删除事件
        watchResp := l.client.Watch(ctx, l.key)
        select {
        case <-ctx.Done():
            l.client.Revoke(ctx, l.lease)
            return ctx.Err()
        case <-watchResp:
            // 锁被释放，重试
            continue
        case <-time.After(time.Duration(l.ttl) * time.Second):
            // 超时，重试
            continue
        }
    }
}

// TryLock 尝试获取锁（非阻塞）
func (l *DistributedLock) TryLock(ctx context.Context) (bool, error) {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.locked {
        return false, fmt.Errorf("already locked")
    }

    // 创建租约
    lease, err := l.client.Grant(ctx, l.ttl)
    if err != nil {
        return false, err
    }
    l.lease = lease.ID

    // 事务尝试获取
    txn := l.client.Txn(ctx).
        If(clientv3.Compare(clientv3.CreateRevision(l.key), "=", 0)).
        Then(clientv3.OpPut(l.key, l.value, clientv3.WithLease(l.lease))).
        Else(clientv3.OpGet(l.key))

    resp, err := txn.Commit()
    if err != nil {
        l.client.Revoke(ctx, l.lease)
        return false, err
    }

    if resp.Succeeded {
        l.locked = true
        go l.keepAlive(context.Background())
        return true, nil
    }

    // 获取失败，撤销租约
    l.client.Revoke(ctx, l.lease)
    return false, nil
}

// Unlock 释放锁
func (l *DistributedLock) Unlock(ctx context.Context) error {
    l.mu.Lock()
    defer l.mu.Unlock()

    if !l.locked {
        return fmt.Errorf("not locked")
    }

    // 删除 key（只能删除自己创建的）
    _, err := l.client.Txn(ctx).
        If(clientv3.Compare(clientv3.Value(l.key), "=", l.value)).
        Then(clientv3.OpDelete(l.key)).
        Commit()

    // 撤销租约
    if l.lease != 0 {
        l.client.Revoke(ctx, l.lease)
    }

    l.locked = false
    return err
}

// keepAlive 自动续约
func (l *DistributedLock) keepAlive(ctx context.Context) {
    ticker := time.NewTicker(time.Duration(l.ttl/3) * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            l.mu.Lock()
            if !l.locked {
                l.mu.Unlock()
                return
            }
            l.mu.Unlock()

            _, err := l.client.KeepAliveOnce(ctx, l.lease)
            if err != nil {
                return
            }
        }
    }
}

// IsLocked 检查锁状态
func (l *DistributedLock) IsLocked() bool {
    l.mu.Lock()
    defer l.mu.Unlock()
    return l.locked
}
```

---

## Session 封装（使用 concurrency 包）

```go
package etcd

import (
    "context"
    "fmt"

    clientv3 "go.etcd.io/etcd/client/v3"
    "go.etcd.io/etcd/client/v3/concurrency"
)

// Session 封装 etcd 会话
type Session struct {
    *concurrency.Session
    client *clientv3.Client
}

// NewSession 创建会话
func NewSession(client *clientv3.Client, ttl int) (*Session, error) {
    s, err := concurrency.NewSession(client, concurrency.WithTTL(ttl))
    if err != nil {
        return nil, err
    }

    return &Session{
        Session: s,
        client:  client,
    }, nil
}

// NewMutex 创建互斥锁
func (s *Session) NewMutex(key string) *concurrency.Mutex {
    return concurrency.NewMutex(s.Session, key)
}

// MutexLock 带超时的锁
func MutexLock(ctx context.Context, mutex *concurrency.Mutex, timeout time.Duration) error {
    ctx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()
    return mutex.Lock(ctx)
}

// Election 领导者选举
type Election struct {
    *concurrency.Election
    session *Session
    key     string
}

// NewElection 创建选举
func (s *Session) NewElection(key string) *Election {
    return &Election{
        Election: concurrency.NewElection(s.Session, key),
        session:  s,
        key:      key,
    }
}

// Campaign 竞选领导者
func (e *Election) Campaign(ctx context.Context, val string) error {
    return e.Election.Campaign(ctx, val)
}

// Proclaim 宣布领导权（已有领导者时）
func (e *Election) Proclaim(ctx context.Context, val string) error {
    return e.Election.Proclaim(ctx, val)
}

// Resign 放弃领导权
func (e *Election) Resign(ctx context.Context) error {
    return e.Election.Resign(ctx)
}

// Leader 获取领导者
func (e *Election) Leader(ctx context.Context) (string, error) {
    resp, err := e.Election.Leader(ctx)
    if err != nil {
        return "", err
    }
    return string(resp.Kvs[0].Value), nil
}

// Observe 观察领导者变化
func (e *Election) Observe(ctx context.Context) <-chan string {
    ch := make(chan string)

    go func() {
        defer close(ch)

        for resp := range e.Election.Observe(ctx) {
            if len(resp.Kvs) > 0 {
                select {
                case ch <- string(resp.Kvs[0].Value):
                case <-ctx.Done():
                    return
                }
            }
        }
    }()

    return ch
}

// LeaderElectionConfig 领导者选举配置
type LeaderElectionConfig struct {
    Name           string
    Key            string
    LeaseDuration  time.Duration
    RenewDeadline  time.Duration
    RetryPeriod    time.Duration
    Callbacks      LeaderCallbacks
}

type LeaderCallbacks struct {
    OnStartedLeading func(ctx context.Context)
    OnStoppedLeading func()
    OnNewLeader      func(identity string)
}

// RunLeaderElection 运行领导者选举
func RunLeaderElection(ctx context.Context, client *clientv3.Client, cfg LeaderElectionConfig) error {
    // 创建会话
    ttl := int(cfg.LeaseDuration.Seconds())
    session, err := NewSession(client, ttl)
    if err != nil {
        return err
    }
    defer session.Close()

    election := session.NewElection(cfg.Key)

    // 启动观察者
    leaderCh := election.Observe(ctx)
    go func() {
        for leader := range leaderCh {
            if cfg.Callbacks.OnNewLeader != nil {
                cfg.Callbacks.OnNewLeader(leader)
            }
        }
    }()

    // 竞选
    err = election.Campaign(ctx, cfg.Name)
    if err != nil {
        return err
    }

    // 成为领导者
    leaderCtx, cancel := context.WithCancel(ctx)
    defer cancel()

    go func() {
        <-leaderCtx.Done()
        election.Resign(context.Background())
        if cfg.Callbacks.OnStoppedLeading != nil {
            cfg.Callbacks.OnStoppedLeading()
        }
    }()

    if cfg.Callbacks.OnStartedLeading != nil {
        cfg.Callbacks.OnStartedLeading(leaderCtx)
    }

    // 保持领导地位
    <-leaderCtx.Done()
    return nil
}
```

---

## Watch 机制

```go
package etcd

import (
    "context"
    "fmt"

    clientv3 "go.etcd.io/etcd/client/v3"
)

// WatchConfig Watch 配置
type WatchConfig struct {
    Prefix      bool              // 是否监听前缀
    FilterPut   bool              // 只监听 PUT 事件
    FilterDelete bool             // 只监听 DELETE 事件
    PrevKV      bool              // 包含修改前的值
    Revision    int64             // 从指定版本开始监听
}

// WatchHandler Watch 事件处理器
type WatchHandler func(ev *clientv3.Event) error

// Watcher 包装 etcd Watch
type Watcher struct {
    client  *clientv3.Client
    watchers map[string]context.CancelFunc
}

func NewWatcher(client *clientv3.Client) *Watcher {
    return &Watcher{
        client:   client,
        watchers: make(map[string]context.CancelFunc),
    }
}

// Watch 监听 key
func (w *Watcher) Watch(ctx context.Context, key string, cfg WatchConfig, handler WatchHandler) error {
    opts := []clientv3.OpOption{}

    if cfg.Prefix {
        opts = append(opts, clientv3.WithPrefix())
    }
    if cfg.PrevKV {
        opts = append(opts, clientv3.WithPrevKV())
    }
    if cfg.Revision > 0 {
        opts = append(opts, clientv3.WithRev(cfg.Revision))
    }

    // 创建过滤
    filters := []clientv3.OpOption{}
    if cfg.FilterPut {
        filters = append(filters, clientv3.WithFilterPut())
    }
    if cfg.FilterDelete {
        filters = append(filters, clientv3.WithFilterDelete())
    }

    watchCtx, cancel := context.WithCancel(ctx)
    w.watchers[key] = cancel
    defer delete(w.watchers, key)

    watchChan := w.client.Watch(watchCtx, key, opts...)

    for watchResp := range watchChan {
        if watchResp.Err() != nil {
            return watchResp.Err()
        }

        for _, ev := range watchResp.Events {
            if err := handler(ev); err != nil {
                return err
            }
        }
    }

    return nil
}

// WatchOnce 监听一次变化
func (w *Watcher) WatchOnce(ctx context.Context, key string) (*clientv3.Event, error) {
    watchChan := w.client.Watch(ctx, key)

    for watchResp := range watchChan {
        if watchResp.Err() != nil {
            return nil, watchResp.Err()
        }

        for _, ev := range watchResp.Events {
            return ev, nil
        }
    }

    return nil, fmt.Errorf("watch closed")
}

// Stop 停止监听
func (w *Watcher) Stop(key string) {
    if cancel, ok := w.watchers[key]; ok {
        cancel()
    }
}

// WatchWithProgressNotify 带进度通知的监听
func (w *Watcher) WatchWithProgressNotify(ctx context.Context, key string) {
    opts := []clientv3.OpOption{
        clientv3.WithProgressNotify(),
    }

    watchChan := w.client.Watch(ctx, key, opts...)

    for watchResp := range watchChan {
        if watchResp.Canceled {
            fmt.Printf("Watch canceled: %v\n", watchResp.Err())
            return
        }

        // 处理进度通知
        if watchResp.Header.Revision > 0 && len(watchResp.Events) == 0 {
            fmt.Printf("Progress notification: revision=%d\n", watchResp.Header.Revision)
            continue
        }

        for _, ev := range watchResp.Events {
            fmt.Printf("Event: %s %s -> %s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
        }
    }
}

// DistributedNotifier 分布式通知器
type DistributedNotifier struct {
    client *clientv3.Client
    prefix string
}

func NewDistributedNotifier(client *clientv3.Client, prefix string) *DistributedNotifier {
    return &DistributedNotifier{
        client: client,
        prefix: prefix,
    }
}

// Notify 发送通知
func (n *DistributedNotifier) Notify(ctx context.Context, channel string, message []byte) error {
    key := fmt.Sprintf("%s/%s/%d", n.prefix, channel, time.Now().UnixNano())
    lease, err := n.client.Grant(ctx, 60) // 1分钟 TTL
    if err != nil {
        return err
    }
    _, err = n.client.Put(ctx, key, string(message), clientv3.WithLease(lease.ID))
    return err
}

// Subscribe 订阅通知
func (n *DistributedNotifier) Subscribe(ctx context.Context, channel string, handler func([]byte)) error {
    key := fmt.Sprintf("%s/%s", n.prefix, channel)

    watchChan := n.client.Watch(ctx, key, clientv3.WithPrefix())

    for watchResp := range watchChan {
        for _, ev := range watchResp.Events {
            if ev.Type == clientv3.EventTypePut {
                handler(ev.Kv.Value)
            }
        }
    }

    return nil
}
```

---

## 事务操作

```go
package etcd

import (
    "context"

    clientv3 "go.etcd.io/etcd/client/v3"
)

// TxnBuilder 事务构建器
type TxnBuilder struct {
    client  *clientv3.Client
    cmps    []clientv3.Cmp
    thenOps []clientv3.Op
    elseOps []clientv3.Op
}

func NewTxn(client *clientv3.Client) *TxnBuilder {
    return &TxnBuilder{
        client: client,
    }
}

// If 添加比较条件
func (b *TxnBuilder) If(cmps ...clientv3.Cmp) *TxnBuilder {
    b.cmps = append(b.cmps, cmps...)
    return b
}

// Then 添加成功操作
func (b *TxnBuilder) Then(ops ...clientv3.Op) *TxnBuilder {
    b.thenOps = append(b.thenOps, ops...)
    return b
}

// Else 添加失败操作
func (b *TxnBuilder) Else(ops ...clientv3.Op) *TxnBuilder {
    b.elseOps = append(b.elseOps, ops...)
    return b
}

// Commit 提交事务
func (b *TxnBuilder) Commit(ctx context.Context) (*clientv3.TxnResponse, error) {
    return b.client.Txn(ctx).
        If(b.cmps...).
        Then(b.thenOps...).
        Else(b.elseOps...).
        Commit()
}

// Compare Helpers

// KeyExists 检查 key 存在
func KeyExists(key string) clientv3.Cmp {
    return clientv3.Compare(clientv3.CreateRevision(key), ">", 0)
}

// KeyNotExists 检查 key 不存在
func KeyNotExists(key string) clientv3.Cmp {
    return clientv3.Compare(clientv3.CreateRevision(key), "=", 0)
}

// KeyEquals 检查 key 值相等
func KeyEquals(key string, value string) clientv3.Cmp {
    return clientv3.Compare(clientv3.Value(key), "=", value)
}

// KeyVersionEquals 检查 key 版本
func KeyVersionEquals(key string, version int64) clientv3.Cmp {
    return clientv3.Compare(clientv3.Version(key), "=", version)
}

// KeyModRevisionEquals 检查修改版本
func KeyModRevisionEquals(key string, modRevision int64) clientv3.Cmp {
    return clientv3.Compare(clientv3.ModRevision(key), "=", modRevision)
}

// 常用事务模式

// CAS (Compare-And-Swap) 比较并交换
func CAS(ctx context.Context, client *clientv3.Client, key, oldVal, newVal string) (bool, error) {
    resp, err := client.Txn(ctx).
        If(clientv3.Compare(clientv3.Value(key), "=", oldVal)).
        Then(clientv3.OpPut(key, newVal)).
        Else(clientv3.OpGet(key)).
        Commit()

    if err != nil {
        return false, err
    }

    return resp.Succeeded, nil
}

// Increment 原子递增
func Increment(ctx context.Context, client *clientv3.Client, key string) (int64, error) {
    for {
        // 获取当前值
        getResp, err := client.Get(ctx, key)
        if err != nil {
            return 0, err
        }

        var current int64 = 0
        var modRev int64 = 0

        if len(getResp.Kvs) > 0 {
            // 解析当前值
            fmt.Sscanf(string(getResp.Kvs[0].Value), "%d", &current)
            modRev = getResp.Kvs[0].ModRevision
        }

        newVal := current + 1

        // 尝试更新
        txn := client.Txn(ctx).
            If(clientv3.Compare(clientv3.ModRevision(key), "=", modRev)).
            Then(clientv3.OpPut(key, fmt.Sprintf("%d", newVal)))

        if len(getResp.Kvs) == 0 {
            txn = client.Txn(ctx).
                If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
                Then(clientv3.OpPut(key, fmt.Sprintf("%d", newVal)))
        }

        resp, err := txn.Commit()
        if err != nil {
            return 0, err
        }

        if resp.Succeeded {
            return newVal, nil
        }

        // 失败重试
    }
}

// DeletePrefix 删除前缀
func DeletePrefix(ctx context.Context, client *clientv3.Client, prefix string) (int64, error) {
    resp, err := client.Delete(ctx, prefix, clientv3.WithPrefix())
    if err != nil {
        return 0, err
    }
    return resp.Deleted, nil
}
```

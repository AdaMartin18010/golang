# 分布式任务分片 (Distributed Task Sharding)

> **分类**: 工程与云原生
> **标签**: #sharding #distributed-tasks #consistent-hashing
> **参考**: Elasticsearch Sharding, Kafka Partitioning

---

## 分片策略架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Distributed Task Sharding Architecture                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. Hash-Based Sharding (一致性哈希)                                         │
│                                                                              │
│     ┌─────────────────────────────────────────────────────────────────┐     │
│     │                     Consistent Hash Ring                         │     │
│     │                                                                  │     │
│     │   0°        60°       120°      180°      240°      300°       │     │
│     │   ┌─────────┬─────────┬─────────┬─────────┬─────────┐           │     │
│     │   │ Node A  │ Node B  │ Node C  │ Node A  │ Node B  │           │     │
│     │   │ (0-60)  │(60-120) │(120-180)│(180-240)│(240-300)│           │     │
│     │   └─────────┴─────────┴─────────┴─────────┴─────────┘           │     │
│     │                                                                  │     │
│     │   Task ID Hash ──► Position on Ring ──► Responsible Node       │     │
│     │                                                                  │     │
│     └─────────────────────────────────────────────────────────────────┘     │
│                                                                              │
│  2. Range-Based Sharding (范围分片)                                          │
│                                                                              │
│     Shard 1: UserID 0 - 1000000                                             │
│     Shard 2: UserID 1000001 - 2000000                                       │
│     Shard 3: UserID 2000001 - 3000000                                       │
│                                                                              │
│  3. Round-Robin Sharding (轮询分片)                                          │
│                                                                              │
│     Task 1 ──► Node A                                                       │
│     Task 2 ──► Node B                                                       │
│     Task 3 ──► Node C                                                       │
│     Task 4 ──► Node A                                                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 一致性哈希实现

```go
package sharding

import (
    "context"
    "encoding/binary"
    "errors"
    "fmt"
    "hash/crc32"
    "sort"
    "sync"
)

// Shard 分片
type Shard struct {
    ID       string
    NodeID   string
    Start    uint32 // 哈希环起始位置
    End      uint32 // 哈希环结束位置
    Weight   int    // 权重
}

// ConsistentHashRing 一致性哈希环
type ConsistentHashRing struct {
    replicas int               // 虚拟节点数
    ring     map[uint32]string // hash -> node
    nodes    map[string]bool   // 实际节点
    sortedHashes []uint32      // 排序后的哈希值
    mu       sync.RWMutex
}

// NewConsistentHashRing 创建一致性哈希环
func NewConsistentHashRing(replicas int) *ConsistentHashRing {
    if replicas <= 0 {
        replicas = 150 // 默认虚拟节点数
    }

    return &ConsistentHashRing{
        replicas: replicas,
        ring:     make(map[uint32]string),
        nodes:    make(map[string]bool),
    }
}

// AddNode 添加节点
func (ch *ConsistentHashRing) AddNode(nodeID string) {
    ch.mu.Lock()
    defer ch.mu.Unlock()

    if _, ok := ch.nodes[nodeID]; ok {
        return
    }

    ch.nodes[nodeID] = true

    // 添加虚拟节点
    for i := 0; i < ch.replicas; i++ {
        virtualKey := fmt.Sprintf("%s#%d", nodeID, i)
        hash := ch.hash(virtualKey)
        ch.ring[hash] = nodeID
    }

    ch.updateSortedHashes()
}

// RemoveNode 移除节点
func (ch *ConsistentHashRing) RemoveNode(nodeID string) {
    ch.mu.Lock()
    defer ch.mu.Unlock()

    if _, ok := ch.nodes[nodeID]; !ok {
        return
    }

    delete(ch.nodes, nodeID)

    // 移除虚拟节点
    for i := 0; i < ch.replicas; i++ {
        virtualKey := fmt.Sprintf("%s#%d", nodeID, i)
        hash := ch.hash(virtualKey)
        delete(ch.ring, hash)
    }

    ch.updateSortedHashes()
}

// GetNode 获取负责 key 的节点
func (ch *ConsistentHashRing) GetNode(key string) (string, error) {
    ch.mu.RLock()
    defer ch.mu.RUnlock()

    if len(ch.ring) == 0 {
        return "", errors.New("no nodes available")
    }

    hash := ch.hash(key)

    // 顺时针找到第一个节点
    idx := sort.Search(len(ch.sortedHashes), func(i int) bool {
        return ch.sortedHashes[i] >= hash
    })

    if idx == len(ch.sortedHashes) {
        idx = 0
    }

    return ch.ring[ch.sortedHashes[idx]], nil
}

// GetNodes 获取 key 的 N 个节点（用于副本）
func (ch *ConsistentHashRing) GetNodes(key string, n int) ([]string, error) {
    ch.mu.RLock()
    defer ch.mu.RUnlock()

    if len(ch.ring) == 0 {
        return nil, errors.New("no nodes available")
    }

    if n > len(ch.nodes) {
        n = len(ch.nodes)
    }

    hash := ch.hash(key)

    idx := sort.Search(len(ch.sortedHashes), func(i int) bool {
        return ch.sortedHashes[i] >= hash
    })

    nodes := make([]string, 0, n)
    seen := make(map[string]bool)

    for len(nodes) < n {
        if idx >= len(ch.sortedHashes) {
            idx = 0
        }

        node := ch.ring[ch.sortedHashes[idx]]
        if !seen[node] {
            seen[node] = true
            nodes = append(nodes, node)
        }
        idx++
    }

    return nodes, nil
}

func (ch *ConsistentHashRing) updateSortedHashes() {
    hashes := make([]uint32, 0, len(ch.ring))
    for h := range ch.ring {
        hashes = append(hashes, h)
    }
    sort.Slice(hashes, func(i, j int) bool { return hashes[i] < hashes[j] })
    ch.sortedHashes = hashes
}

func (ch *ConsistentHashRing) hash(key string) uint32 {
    return crc32.ChecksumIEEE([]byte(key))
}

// GetAllNodes 获取所有节点
func (ch *ConsistentHashRing) GetAllNodes() []string {
    ch.mu.RLock()
    defer ch.mu.RUnlock()

    nodes := make([]string, 0, len(ch.nodes))
    for node := range ch.nodes {
        nodes = append(nodes, node)
    }
    return nodes
}
```

---

## 分片管理器实现

```go
package sharding

import (
    "context"
    "fmt"
    "sync"
)

// ShardManager 分片管理器
type ShardManager struct {
    ring       *ConsistentHashRing
    shardCount int

    // 分片分配
    assignments map[string][]string // shardID -> nodeIDs

    // 节点管理
    mu          sync.RWMutex
    nodes       map[string]*Node

    // 重平衡控制
    rebalancer  Rebalancer
}

// Node 节点信息
type Node struct {
    ID       string
    Address  string
    Capacity int // 处理能力
    Load     int // 当前负载
    Healthy  bool
}

// Rebalancer 重平衡器接口
type Rebalancer interface {
    Rebalance(ctx context.Context, current map[string][]string, changes NodeChanges) (map[string][]string, error)
}

// NodeChanges 节点变更
type NodeChanges struct {
    Added   []string
    Removed []string
}

// NewShardManager 创建分片管理器
func NewShardManager(shardCount int, replicas int) *ShardManager {
    return &ShardManager{
        ring:        NewConsistentHashRing(replicas),
        shardCount:  shardCount,
        assignments: make(map[string][]string),
        nodes:       make(map[string]*Node),
    }
}

// RegisterNode 注册节点
func (sm *ShardManager) RegisterNode(node *Node) error {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    if _, exists := sm.nodes[node.ID]; exists {
        return fmt.Errorf("node %s already registered", node.ID)
    }

    sm.nodes[node.ID] = node
    sm.ring.AddNode(node.ID)

    // 触发重平衡
    changes := NodeChanges{Added: []string{node.ID}}
    sm.triggerRebalance(context.Background(), changes)

    return nil
}

// DeregisterNode 注销节点
func (sm *ShardManager) DeregisterNode(nodeID string) error {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    if _, exists := sm.nodes[nodeID]; !exists {
        return fmt.Errorf("node %s not found", node.ID)
    }

    delete(sm.nodes, nodeID)
    sm.ring.RemoveNode(nodeID)

    // 触发重平衡
    changes := NodeChanges{Removed: []string{nodeID}}
    sm.triggerRebalance(context.Background(), changes)

    return nil
}

// GetShardForKey 获取 key 对应的分片
func (sm *ShardManager) GetShardForKey(key string) (string, error) {
    nodeID, err := sm.ring.GetNode(key)
    if err != nil {
        return "", err
    }

    // 根据节点获取分片
    sm.mu.RLock()
    defer sm.mu.RUnlock()

    for shardID, nodes := range sm.assignments {
        for _, nid := range nodes {
            if nid == nodeID {
                return shardID, nil
            }
        }
    }

    return "", fmt.Errorf("no shard found for key %s", key)
}

// GetNodesForShard 获取分片的所有节点（主节点+副本）
func (sm *ShardManager) GetNodesForShard(shardID string) ([]string, error) {
    sm.mu.RLock()
    defer sm.mu.RUnlock()

    nodes, ok := sm.assignments[shardID]
    if !ok {
        return nil, fmt.Errorf("shard %s not found", shardID)
    }

    return nodes, nil
}

// GetShardAssignment 获取分片分配信息
func (sm *ShardManager) GetShardAssignment() map[string][]string {
    sm.mu.RLock()
    defer sm.mu.RUnlock()

    result := make(map[string][]string)
    for k, v := range sm.assignments {
        nodes := make([]string, len(v))
        copy(nodes, v)
        result[k] = nodes
    }
    return result
}

// triggerRebalance 触发重平衡
func (sm *ShardManager) triggerRebalance(ctx context.Context, changes NodeChanges) {
    if sm.rebalancer == nil {
        return
    }

    newAssignments, err := sm.rebalancer.Rebalance(ctx, sm.assignments, changes)
    if err != nil {
        // 记录错误
        return
    }

    sm.assignments = newAssignments
}

// SetRebalancer 设置重平衡器
func (sm *ShardManager) SetRebalancer(r Rebalancer) {
    sm.rebalancer = r
}

// GetNodeLoad 获取节点负载
func (sm *ShardManager) GetNodeLoad(nodeID string) (int, error) {
    sm.mu.RLock()
    defer sm.mu.RUnlock()

    node, ok := sm.nodes[nodeID]
    if !ok {
        return 0, fmt.Errorf("node %s not found", nodeID)
    }

    return node.Load, nil
}
```

---

## 任务分片执行器

```go
package sharding

import (
    "context"
    "fmt"
    "sync"
)

// ShardedTaskExecutor 分片任务执行器
type ShardedTaskExecutor struct {
    shardManager *ShardManager
    localNodeID  string

    // 执行器
    executors    map[string]ShardExecutor
    mu           sync.RWMutex
}

// ShardExecutor 分片执行器接口
type ShardExecutor interface {
    Execute(ctx context.Context, shardID string, tasks []Task) ([]Result, error)
    GetLoad() int
}

// Task 任务
type Task struct {
    ID       string
    Key      string
    Payload  interface{}
}

// Result 执行结果
type Result struct {
    TaskID string
    Data   interface{}
    Error  error
}

// NewShardedTaskExecutor 创建分片任务执行器
func NewShardedTaskExecutor(sm *ShardManager, localNodeID string) *ShardedTaskExecutor {
    return &ShardedTaskExecutor{
        shardManager: sm,
        localNodeID:  localNodeID,
        executors:    make(map[string]ShardExecutor),
    }
}

// RegisterExecutor 注册分片执行器
func (ste *ShardedTaskExecutor) RegisterExecutor(shardID string, executor ShardExecutor) {
    ste.mu.Lock()
    defer ste.mu.Unlock()
    ste.executors[shardID] = executor
}

// ExecuteTask 执行任务（自动路由到对应分片）
func (ste *ShardedTaskExecutor) ExecuteTask(ctx context.Context, task Task) (*Result, error) {
    shardID, err := ste.shardManager.GetShardForKey(task.Key)
    if err != nil {
        return nil, err
    }

    nodes, err := ste.shardManager.GetNodesForShard(shardID)
    if err != nil {
        return nil, err
    }

    if len(nodes) == 0 {
        return nil, fmt.Errorf("no nodes available for shard %s", shardID)
    }

    // 主节点执行
    primaryNode := nodes[0]

    if primaryNode == ste.localNodeID {
        // 本地执行
        return ste.executeLocal(ctx, shardID, task)
    }

    // 远程调用（简化实现）
    return ste.executeRemote(ctx, primaryNode, shardID, task)
}

// ExecuteTasksBatch 批量执行任务
func (ste *ShardedTaskExecutor) ExecuteTasksBatch(ctx context.Context, tasks []Task) ([]Result, error) {
    // 按分片分组
    shardGroups := make(map[string][]Task)
    for _, task := range tasks {
        shardID, err := ste.shardManager.GetShardForKey(task.Key)
        if err != nil {
            continue
        }
        shardGroups[shardID] = append(shardGroups[shardID], task)
    }

    var wg sync.WaitGroup
    resultChan := make(chan []Result, len(shardGroups))
    errChan := make(chan error, len(shardGroups))

    for shardID, shardTasks := range shardGroups {
        wg.Add(1)
        go func(sid string, st []Task) {
            defer wg.Done()

            results, err := ste.executeShardTasks(ctx, sid, st)
            if err != nil {
                errChan <- err
                return
            }
            resultChan <- results
        }(shardID, shardTasks)
    }

    wg.Wait()
    close(resultChan)
    close(errChan)

    // 收集结果
    var allResults []Result
    for results := range resultChan {
        allResults = append(allResults, results...)
    }

    // 检查错误
    for err := range errChan {
        if err != nil {
            return nil, err
        }
    }

    return allResults, nil
}

func (ste *ShardedTaskExecutor) executeLocal(ctx context.Context, shardID string, task Task) (*Result, error) {
    ste.mu.RLock()
    executor, ok := ste.executors[shardID]
    ste.mu.RUnlock()

    if !ok {
        return nil, fmt.Errorf("no executor for shard %s", shardID)
    }

    results, err := executor.Execute(ctx, shardID, []Task{task})
    if err != nil {
        return nil, err
    }

    if len(results) == 0 {
        return nil, fmt.Errorf("no result returned")
    }

    return &results[0], nil
}

func (ste *ShardedTaskExecutor) executeRemote(ctx context.Context, nodeID, shardID string, task Task) (*Result, error) {
    // 实际实现中通过 RPC 调用远程节点
    // 这里简化处理
    return nil, fmt.Errorf("remote execution not implemented")
}

func (ste *ShardedTaskExecutor) executeShardTasks(ctx context.Context, shardID string, tasks []Task) ([]Result, error) {
    nodes, err := ste.shardManager.GetNodesForShard(shardID)
    if err != nil {
        return nil, err
    }

    if len(nodes) == 0 {
        return nil, fmt.Errorf("no nodes for shard %s", shardID)
    }

    primaryNode := nodes[0]

    if primaryNode == ste.localNodeID {
        return ste.executors[shardID].Execute(ctx, shardID, tasks)
    }

    // 远程执行
    return nil, fmt.Errorf("remote execution not implemented")
}

// RebalanceShards 重平衡分片
func (ste *ShardedTaskExecutor) RebalanceShards(ctx context.Context) error {
    assignment := ste.shardManager.GetShardAssignment()

    for shardID, nodes := range assignment {
        primaryNode := nodes[0]

        if primaryNode != ste.localNodeID {
            // 如果当前节点不再是主节点，需要迁移数据
            if err := ste.migrateShard(ctx, shardID, primaryNode); err != nil {
                return err
            }
        }
    }

    return nil
}

func (ste *ShardedTaskExecutor) migrateShard(ctx context.Context, shardID string, targetNode string) error {
    // 实现分片数据迁移
    return nil
}
```

---

## 使用示例

```go
package main

import (
    "context"
    "fmt"

    "sharding"
)

func main() {
    // 创建分片管理器
    sm := sharding.NewShardManager(32, 150)

    // 注册节点
    for i := 1; i <= 3; i++ {
        node := &sharding.Node{
            ID:       fmt.Sprintf("node-%d", i),
            Address:  fmt.Sprintf("192.168.1.%d", i),
            Capacity: 100,
        }
        sm.RegisterNode(node)
    }

    // 创建分片执行器
    executor := sharding.NewShardedTaskExecutor(sm, "node-1")

    // 注册执行器
    for i := 0; i < 32; i++ {
        shardID := fmt.Sprintf("shard-%d", i)
        executor.RegisterExecutor(shardID, &MyShardExecutor{})
    }

    // 执行任务
    task := sharding.Task{
        ID:      "task-1",
        Key:     "user-12345",
        Payload: "data",
    }

    result, err := executor.ExecuteTask(context.Background(), task)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Result: %+v\n", result)
}

type MyShardExecutor struct{}

func (e *MyShardExecutor) Execute(ctx context.Context, shardID string, tasks []sharding.Task) ([]sharding.Result, error) {
    var results []sharding.Result
    for _, task := range tasks {
        results = append(results, sharding.Result{
            TaskID: task.ID,
            Data:   "processed",
        })
    }
    return results, nil
}

func (e *MyShardExecutor) GetLoad() int {
    return 0
}
```

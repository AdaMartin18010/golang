# FT-014: 两阶段提交形式化分析 (Two-Phase Commit Formalization)

> **维度**: Formal Theory
> **级别**: S (16+ KB)
> **标签**: #2pc #distributed-transactions #atomic-commit #consensus
> **权威来源**: [Consensus on Transaction Commit](https://dl.acm.org/doi/10.1145/114539.114544) - Jim Gray

---

## 2PC 协议概述

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Two-Phase Commit Protocol                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Phase 1: Voting (准备阶段)                                                   │
│                                                                              │
│  Coordinator                Participants                                     │
│       │                          │                                           │
│       │──────PREPARE────────────►│                                           │
│       │                          │                                           │
│       │◄─────YES/NO─────────────│                                           │
│       │                          │                                           │
│  (收集所有投票)                                                               │
│                                                                              │
│  Phase 2: Decision (决策阶段)                                                 │
│                                                                              │
│       │                          │                                           │
│       │──────COMMIT─────────────►│  (如果所有 YES)                            │
│       │      or                  │                                           │
│       │──────ABORT──────────────►│  (如果任何 NO)                             │
│       │                          │                                           │
│       │◄────ACK─────────────────│                                           │
│       │                          │                                           │
│                                                                              │
│  状态机:                                                                      │
│  Coordinator: INIT → PREPARING → (PREPARED) → COMMITTING/ABORTING → DONE     │
│  Participant: INIT → PREPARING → (READY/ABORTED) → COMMITTED/ABORTED         │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 形式化规范

### 时序逻辑规约

```
安全性质 (Safety):
1. 一致性: ¬(∃i: committed(i) ∧ ∃j: aborted(j))
   不可能同时存在已提交和已中止的参与者

2. 有效性:
   - 如果任何参与者提交，则协调者必须发送 COMMIT
   - 如果任何参与者中止，则协调者必须发送 ABORT

活性性质 (Liveness):
3. 终止性: ◇(all_committed ∨ all_aborted)
   最终所有参与者都会达到终止状态

4. 原子性:
   - 如果协调者决定 COMMIT，则所有已就绪的参与者最终提交
   - 如果协调者决定 ABORT，则所有参与者最终中止
```

### TLA+ 风格规约

```tla
(* 2PC 规约 *)

VARIABLES
    coordinator_state,
    participant_states,
    coordinator_vote_count,
    messages

Init ==
    ∧ coordinator_state = "INIT"
    ∧ participant_states = [p ∈ Participants ↦ "INIT"]
    ∧ coordinator_vote_count = 0
    ∧ messages = {}

(* Phase 1: Coordinator sends PREPARE *)
SendPrepare ==
    ∧ coordinator_state = "INIT"
    ∧ coordinator_state' = "PREPARING"
    ∧ messages' = messages ∪ {[type ↦ "PREPARE", to ↦ p] : p ∈ Participants}
    ∧ UNCHANGED ⟨participant_states, coordinator_vote_count⟩

(* Phase 1: Participant votes *)
ParticipantVote(p) ==
    ∧ participant_states[p] = "INIT"
    ∧ [type ↦ "PREPARE", to ↦ p] ∈ messages
    ∧ participant_states' = [participant_states EXCEPT ![p] = "READY"]
    ∧ messages' = messages ∪ {[type ↦ "YES", from ↦ p]}
    ∧ UNCHANGED ⟨coordinator_state, coordinator_vote_count⟩

(* Phase 1: Coordinator collects votes *)
CollectVote(p) ==
    ∧ coordinator_state = "PREPARING"
    ∧ [type ↦ "YES", from ↦ p] ∈ messages
    ∧ coordinator_vote_count' = coordinator_vote_count + 1
    ∧ IF coordinator_vote_count' = Cardinality(Participants)
       THEN coordinator_state' = "PREPARED"
       ELSE UNCHANGED coordinator_state
    ∧ UNCHANGED ⟨participant_states, messages⟩

(* Phase 2: Coordinator decides and sends COMMIT *)
DecideCommit ==
    ∧ coordinator_state = "PREPARED"
    ∧ coordinator_state' = "COMMITTING"
    ∧ messages' = messages ∪ {[type ↦ "COMMIT", to ↦ p] : p ∈ Participants}
    ∧ UNCHANGED ⟨participant_states, coordinator_vote_count⟩

(* Phase 2: Participant commits *)
ParticipantCommit(p) ==
    ∧ participant_states[p] = "READY"
    ∧ [type ↦ "COMMIT", to ↦ p] ∈ messages
    ∧ participant_states' = [participant_states EXCEPT ![p] = "COMMITTED"]
    ∧ UNCHANGED ⟨coordinator_state, coordinator_vote_count, messages⟩

Next ==
    ∨ SendPrepare
    ∨ ∃p ∈ Participants : ParticipantVote(p)
    ∨ ∃p ∈ Participants : CollectVote(p)
    ∨ DecideCommit
    ∨ ∃p ∈ Participants : ParticipantCommit(p)

(* 不变式: 不可能同时存在已提交和已中止 *)
Consistency ==
    ¬(∃p, q ∈ Participants :
        participant_states[p] = "COMMITTED" ∧
        participant_states[q] = "ABORTED")
```

---

## Go 实现

```go
package twopc

import (
    "context"
    "errors"
    "fmt"
    "sync"
    "time"
)

// Coordinator 协调者
type Coordinator struct {
    participants []Participant
    timeout      time.Duration
}

// Participant 参与者接口
type Participant interface {
    ID() string
    Prepare(ctx context.Context) error
    Commit(ctx context.Context) error
    Abort(ctx context.Context) error
}

// Transaction 事务
type Transaction struct {
    ID   string
    Data interface{}
}

// Execute 执行 2PC
func (c *Coordinator) Execute(ctx context.Context, txn *Transaction) error {
    // Phase 1: Prepare
    votes := make(map[string]bool)
    var mu sync.Mutex
    var wg sync.WaitGroup

    ctx, cancel := context.WithTimeout(ctx, c.timeout)
    defer cancel()

    // 并行发送 PREPARE
    for _, p := range c.participants {
        wg.Add(1)
        go func(part Participant) {
            defer wg.Done()

            err := part.Prepare(ctx)
            mu.Lock()
            votes[part.ID()] = (err == nil)
            mu.Unlock()
        }(p)
    }

    wg.Wait()

    // 检查所有投票
    allYes := true
    for _, vote := range votes {
        if !vote {
            allYes = false
            break
        }
    }

    // Phase 2: Decision
    if allYes {
        // 发送 COMMIT
        return c.broadcastCommit(ctx)
    } else {
        // 发送 ABORT
        c.broadcastAbort(ctx)
        return errors.New("transaction aborted")
    }
}

func (c *Coordinator) broadcastCommit(ctx context.Context) error {
    var wg sync.WaitGroup
    errChan := make(chan error, len(c.participants))

    for _, p := range c.participants {
        wg.Add(1)
        go func(part Participant) {
            defer wg.Done()
            if err := part.Commit(ctx); err != nil {
                errChan <- fmt.Errorf("commit failed on %s: %w", part.ID(), err)
            }
        }(p)
    }

    wg.Wait()
    close(errChan)

    for err := range errChan {
        return err
    }
    return nil
}

func (c *Coordinator) broadcastAbort(ctx context.Context) {
    var wg sync.WaitGroup

    for _, p := range c.participants {
        wg.Add(1)
        go func(part Participant) {
            defer wg.Done()
            part.Abort(ctx)
        }(p)
    }

    wg.Wait()
}

// LocalParticipant 本地参与者实现
type LocalParticipant struct {
    id      string
    data    map[string]interface{}
    mu      sync.RWMutex
    prepared bool
    backup  map[string]interface{}
}

func (p *LocalParticipant) ID() string {
    return p.id
}

func (p *LocalParticipant) Prepare(ctx context.Context) error {
    p.mu.Lock()
    defer p.mu.Unlock()

    // 创建备份
    p.backup = make(map[string]interface{})
    for k, v := range p.data {
        p.backup[k] = v
    }

    p.prepared = true
    return nil
}

func (p *LocalParticipant) Commit(ctx context.Context) error {
    p.mu.Lock()
    defer p.mu.Unlock()

    if !p.prepared {
        return errors.New("not prepared")
    }

    // 清理备份，正式提交
    p.backup = nil
    p.prepared = false
    return nil
}

func (p *LocalParticipant) Abort(ctx context.Context) error {
    p.mu.Lock()
    defer p.mu.Unlock()

    if !p.prepared {
        return nil
    }

    // 回滚
    p.data = p.backup
    p.backup = nil
    p.prepared = false
    return nil
}
```

---

## 故障分析

```
故障场景与恢复:

1. 协调者故障 (Phase 1)
   - 参与者等待 PREPARE 超时
   - 单方面中止事务
   - 无需恢复

2. 协调者故障 (Phase 2)
   - 参与者已投票 YES，等待 DECISION
   - 阻塞直到协调者恢复
   - 解决方案: 3PC 或协调者高可用

3. 参与者故障 (Phase 1)
   - 协调者等待投票超时
   - 决定 ABORT
   - 正常处理

4. 参与者故障 (Phase 2)
   - 协调者已决定 COMMIT
   - 参与者恢复后需要重放 COMMIT
   - 需要持久化日志

阻塞问题:
- 2PC 在协调者故障时可能阻塞
- 最坏情况: 参与者持有锁，等待协调者
- 解决: 3PC (增加预提交阶段) 或 Paxos/ Raft 替换 2PC
```

---

## 参考文献

1. [Consensus on Transaction Commit](https://dl.acm.org/doi/10.1145/114539.114544) - Jim Gray
2. [Notes on Data Base Operating Systems](https://dl.acm.org/doi/book/10.5555/914517)
3. [Atomic Commit](https://en.wikipedia.org/wiki/Atomic_commit)

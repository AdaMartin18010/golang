# 任务调度策略 (Task Scheduling Strategies)

> **分类**: 工程与云原生
> **标签**: #scheduling-strategy #load-balancing #affinity

---

## 负载均衡策略

```go
type SchedulingStrategy interface {
    SelectWorker(workers []Worker, task *Task) (Worker, error)
}

// 轮询
type RoundRobinStrategy struct {
    current uint64
}

func (rr *RoundRobinStrategy) SelectWorker(workers []Worker, task *Task) (Worker, error) {
    if len(workers) == 0 {
        return Worker{}, errors.New("no workers available")
    }

    idx := atomic.AddUint64(&rr.current, 1) % uint64(len(workers))
    return workers[idx], nil
}

// 最少连接
type LeastConnectionsStrategy struct{}

func (lc *LeastConnectionsStrategy) SelectWorker(workers []Worker, task *Task) (Worker, error) {
    var best Worker
    minConn := int(^uint(0) >> 1)  // MaxInt

    for _, w := range workers {
        if w.ActiveTasks < minConn {
            minConn = w.ActiveTasks
            best = w
        }
    }

    return best, nil
}

// 加权随机
type WeightedRandomStrategy struct{}

func (wr *WeightedRandomStrategy) SelectWorker(workers []Worker, task *Task) (Worker, error) {
    totalWeight := 0
    for _, w := range workers {
        totalWeight += w.Weight
    }

    r := rand.Intn(totalWeight)

    for _, w := range workers {
        r -= w.Weight
        if r < 0 {
            return w, nil
        }
    }

    return workers[0], nil
}

// 资源感知
type ResourceAwareStrategy struct{}

func (ra *ResourceAwareStrategy) SelectWorker(workers []Worker, task *Task) (Worker, error) {
    var candidates []Worker

    for _, w := range workers {
        // 检查资源是否满足
        if w.CPUUsage < 70 && w.MemoryUsage < 80 {
            // 计算资源匹配度分数
            score := ra.calculateScore(w, task)
            w.Score = score
            candidates = append(candidates, w)
        }
    }

    if len(candidates) == 0 {
        return Worker{}, errors.New("no worker with sufficient resources")
    }

    // 选择分数最高的
    sort.Slice(candidates, func(i, j int) bool {
        return candidates[i].Score > candidates[j].Score
    })

    return candidates[0], nil
}

func (ra *ResourceAwareStrategy) calculateScore(w Worker, task *Task) float64 {
    score := 100.0

    // CPU 余量
    score += (100 - w.CPUUsage) * 0.5

    // 内存余量
    score += (100 - w.MemoryUsage) * 0.3

    // 任务亲和性
    if w.LastTaskType == task.Type {
        score += 10  // 缓存友好
    }

    return score
}
```

---

## 亲和性与反亲和性

```go
type AffinityRule struct {
    Type      string  // affinity, anti-affinity
    Key       string  // 标签键
    Values    []string
    Weight    int     // 权重
}

func calculateAffinityScore(worker Worker, task *Task, rules []AffinityRule) int {
    score := 0

    for _, rule := range rules {
        matches := 0
        for k, v := range worker.Labels {
            if k == rule.Key {
                for _, val := range rule.Values {
                    if v == val {
                        matches++
                    }
                }
            }
        }

        if rule.Type == "affinity" {
            // 亲和性：匹配越多分越高
            score += matches * rule.Weight
        } else {
            // 反亲和性：匹配越多分越低
            score -= matches * rule.Weight
        }
    }

    return score
}

// 使用示例
task := &Task{
    AffinityRules: []AffinityRule{
        {
            Type:   "affinity",
            Key:    "zone",
            Values: []string{"zone-a"},
            Weight: 100,
        },
        {
            Type:   "anti-affinity",
            Key:    "task-type",
            Values: []string{"heavy"},
            Weight: 50,
        },
    },
}
```

---

## 抢占式调度

```go
type PreemptiveScheduler struct {
    workers []Worker
    tasks   map[string]*RunningTask
}

func (ps *PreemptiveScheduler) Schedule(task *Task) error {
    // 查找最佳 worker
    worker, err := ps.findWorker(task)
    if err != nil {
        // 资源不足，尝试抢占
        return ps.preempt(task)
    }

    return ps.assign(worker, task)
}

func (ps *PreemptiveScheduler) preempt(highPriTask *Task) error {
    // 查找可以被抢占的低优先级任务
    for id, runningTask := range ps.tasks {
        if runningTask.Priority < highPriTask.Priority {
            // 抢占
            ps.evict(id)
            return ps.assign(runningTask.Worker, highPriTask)
        }
    }

    return errors.New("no preemptable tasks found")
}

func (ps *PreemptiveScheduler) evict(taskID string) {
    task := ps.tasks[taskID]

    // 发送抢占信号
    task.Cancel()

    // 保存状态（检查点）
    ps.saveCheckpoint(task)

    // 移除
    delete(ps.tasks, taskID)
}
```

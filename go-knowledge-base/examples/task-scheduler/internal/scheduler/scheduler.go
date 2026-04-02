package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// Config 调度器配置
type Config struct {
	Strategy      string        // 调度策略: round-robin, least-tasks, priority
	BatchSize     int           // 批处理大小
	CheckInterval time.Duration // 检查间隔
}

// Scheduler 任务调度器
type Scheduler struct {
	config    *Config
	isLeader  bool
	leaderMu  sync.RWMutex
	tasks     chan Task
	workers   map[string]*Worker
	workersMu sync.RWMutex
	storage   Storage
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

// Task 任务定义
type Task struct {
	ID        string
	Type      string
	Payload   []byte
	Priority  int
	Scheduled time.Time
}

// Worker 工作节点
type Worker struct {
	ID       string
	Addr     string
	Capacity int
	Tasks    int
	LastSeen time.Time
}

// Storage 存储接口
type Storage interface {
	AcquireLeadership(ctx context.Context, nodeID string, ttl time.Duration) (bool, error)
	ReleaseLeadership(ctx context.Context, nodeID string) error
	GetTasks(ctx context.Context, status string, limit int) ([]Task, error)
	UpdateTaskStatus(ctx context.Context, taskID string, status string) error
	RegisterWorker(ctx context.Context, worker *Worker) error
	GetWorkers(ctx context.Context) ([]Worker, error)
}

// New 创建调度器
func New(config *Config, storage Storage) *Scheduler {
	return &Scheduler{
		config:  config,
		tasks:   make(chan Task, 1000),
		workers: make(map[string]*Worker),
		storage: storage,
	}
}

// Start 启动调度器
func (s *Scheduler) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)
	
	// 1. 领导选举
	s.wg.Add(1)
	go s.leaderElection()
	
	// 2. 任务调度
	s.wg.Add(1)
	go s.taskScheduler()
	
	// 3. 工作节点监控
	s.wg.Add(1)
	go s.workerMonitor()
	
	// 4. 指标收集
	s.wg.Add(1)
	go s.metricsCollector()
	
	s.wg.Wait()
	return nil
}

// Shutdown 优雅关闭
func (s *Scheduler) Shutdown(ctx context.Context) error {
	s.cancel()
	
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("shutdown timeout")
	}
}

// IsLeader 是否为主节点
func (s *Scheduler) IsLeader() bool {
	s.leaderMu.RLock()
	defer s.leaderMu.RUnlock()
	return s.isLeader
}

// leaderElection 领导选举
func (s *Scheduler) leaderElection() {
	defer s.wg.Done()
	
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	nodeID := fmt.Sprintf("scheduler-%d", time.Now().UnixNano())
	
	for {
		select {
		case <-s.ctx.Done():
			if s.isLeader {
				s.storage.ReleaseLeadership(s.ctx, nodeID)
			}
			return
		case <-ticker.C:
			isLeader, err := s.storage.AcquireLeadership(s.ctx, nodeID, 10*time.Second)
			if err != nil {
				log.Printf("Leadership acquisition failed: %v", err)
				continue
			}
			
			s.leaderMu.Lock()
			s.isLeader = isLeader
			s.leaderMu.Unlock()
			
			if isLeader {
				log.Println("Became leader")
			}
		}
	}
}

// taskScheduler 任务调度
func (s *Scheduler) taskScheduler() {
	defer s.wg.Done()
	
	ticker := time.NewTicker(s.config.CheckInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			if !s.IsLeader() {
				continue
			}
			s.distributeTasks()
		}
	}
}

// distributeTasks 分发任务
func (s *Scheduler) distributeTasks() {
	tasks, err := s.storage.GetTasks(s.ctx, "pending", s.config.BatchSize)
	if err != nil {
		log.Printf("Failed to get tasks: %v", err)
		return
	}
	
	workers, err := s.storage.GetWorkers(s.ctx)
	if err != nil {
		log.Printf("Failed to get workers: %v", err)
		return
	}
	
	for _, task := range tasks {
		worker := s.selectWorker(workers, task)
		if worker == nil {
			continue
		}
		
		// 分发任务到工作节点
		if err := s.dispatchToWorker(task, worker); err != nil {
			log.Printf("Failed to dispatch task %s: %v", task.ID, err)
			continue
		}
		
		// 更新任务状态
		if err := s.storage.UpdateTaskStatus(s.ctx, task.ID, "assigned"); err != nil {
			log.Printf("Failed to update task status: %v", err)
		}
	}
}

// selectWorker 选择工作节点
func (s *Scheduler) selectWorker(workers []Worker, task Task) *Worker {
	if len(workers) == 0 {
		return nil
	}
	
	switch s.config.Strategy {
	case "least-tasks":
		// 最少任务优先
		var selected *Worker
		minTasks := int(^uint(0) >> 1)
		for i := range workers {
			if workers[i].Tasks < minTasks && workers[i].Capacity > workers[i].Tasks {
				minTasks = workers[i].Tasks
				selected = &workers[i]
			}
		}
		return selected
		
	case "priority":
		// 高优先级任务分配到最强节点
		var selected *Worker
		maxCapacity := 0
		for i := range workers {
			if workers[i].Capacity > maxCapacity {
				maxCapacity = workers[i].Capacity
				selected = &workers[i]
			}
		}
		return selected
		
	default: // round-robin
		// 轮询
		for i := range workers {
			if workers[i].Capacity > workers[i].Tasks {
				return &workers[i]
			}
		}
		return nil
	}
}

// dispatchToWorker 分发任务到工作节点
func (s *Scheduler) dispatchToWorker(task Task, worker *Worker) error {
	// 实际实现: HTTP/gRPC 调用工作节点
	log.Printf("Dispatching task %s to worker %s", task.ID, worker.ID)
	return nil
}

// workerMonitor 工作节点监控
func (s *Scheduler) workerMonitor() {
	defer s.wg.Done()
	
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			if !s.IsLeader() {
				continue
			}
			s.checkWorkerHealth()
		}
	}
}

// checkWorkerHealth 检查工作节点健康状态
func (s *Scheduler) checkWorkerHealth() {
	workers, err := s.storage.GetWorkers(s.ctx)
	if err != nil {
		log.Printf("Failed to get workers: %v", err)
		return
	}
	
	for _, worker := range workers {
		if time.Since(worker.LastSeen) > 30*time.Second {
			log.Printf("Worker %s is unhealthy, last seen %v ago", worker.ID, time.Since(worker.LastSeen))
			// 重新分配该节点的任务
			s.reassignWorkerTasks(worker.ID)
		}
	}
}

// reassignWorkerTasks 重新分配工作节点的任务
func (s *Scheduler) reassignWorkerTasks(workerID string) {
	log.Printf("Reassigning tasks from worker %s", workerID)
	// 实现: 查询该节点的任务并重新调度
}

// metricsCollector 指标收集
func (s *Scheduler) metricsCollector() {
	defer s.wg.Done()
	
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.collectMetrics()
		}
	}
}

// collectMetrics 收集指标
func (s *Scheduler) collectMetrics() {
	// 实现: Prometheus 指标上报
	// - 任务队列长度
	// - 调度延迟
	// - 工作节点数量
	// - 任务成功率
}

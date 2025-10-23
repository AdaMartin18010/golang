// Package patterns - 同步模式简化版本
package patterns

import "fmt"

// GenerateMutexSimple 生成简化的Mutex模式
func GenerateMutexSimple(pkg string) string {
	return fmt.Sprintf("package %s\n\nimport \"sync\"\n\ntype SafeCounter struct {\n\tmu sync.Mutex\n\tvalue int\n}\n\nfunc (c *SafeCounter) Inc() {\n\tc.mu.Lock()\n\tdefer c.mu.Unlock()\n\tc.value++\n}\n", pkg)
}

// GenerateRWMutexSimple 生成简化的RWMutex模式
func GenerateRWMutexSimple(pkg string) string {
	return fmt.Sprintf("package %s\n\nimport \"sync\"\n\ntype Cache struct {\n\tmu sync.RWMutex\n\tdata map[string]interface{}\n}\n\nfunc (c *Cache) Get(key string) interface{} {\n\tc.mu.RLock()\n\tdefer c.mu.RUnlock()\n\treturn c.data[key]\n}\n", pkg)
}

// GenerateWaitGroupSimple 生成简化的WaitGroup模式
func GenerateWaitGroupSimple(pkg string) string {
	return fmt.Sprintf("package %s\n\nimport \"sync\"\n\nfunc ParallelTasks(tasks []func()) {\n\tvar wg sync.WaitGroup\n\tfor _, task := range tasks {\n\t\twg.Add(1)\n\t\tgo func(t func()) {\n\t\t\tdefer wg.Done()\n\t\t\tt()\n\t\t}(task)\n\t}\n\twg.Wait()\n}\n", pkg)
}

// GenerateOnceSimple 生成简化的Once模式
func GenerateOnceSimple(pkg string) string {
	return fmt.Sprintf("package %s\n\nimport \"sync\"\n\nvar (\n\tinstance *Singleton\n\tonce sync.Once\n)\n\ntype Singleton struct{}\n\nfunc GetInstance() *Singleton {\n\tonce.Do(func() {\n\t\tinstance = &Singleton{}\n\t})\n\treturn instance\n}\n", pkg)
}

// GenerateSemaphoreSimple 生成简化的Semaphore模式
func GenerateSemaphoreSimple(pkg string) string {
	return fmt.Sprintf("package %s\n\ntype Semaphore struct {\n\tslots chan struct{}\n}\n\nfunc NewSemaphore(n int) *Semaphore {\n\treturn &Semaphore{slots: make(chan struct{}, n)}\n}\n\nfunc (s *Semaphore) Acquire() {\n\ts.slots <- struct{}{}\n}\n\nfunc (s *Semaphore) Release() {\n\t<-s.slots\n}\n", pkg)
}

// GenerateBarrierSimple 生成简化的Barrier模式
func GenerateBarrierSimple(pkg string) string {
	return fmt.Sprintf("package %s\n\nimport \"sync\"\n\ntype Barrier struct {\n\tcount int\n\ttarget int\n\tmu sync.Mutex\n\tcond *sync.Cond\n}\n\nfunc NewBarrier(n int) *Barrier {\n\tb := &Barrier{target: n}\n\tb.cond = sync.NewCond(&b.mu)\n\treturn b\n}\n", pkg)
}

// GenerateCountDownLatchSimple 生成简化的CountDownLatch模式
func GenerateCountDownLatchSimple(pkg string) string {
	return fmt.Sprintf("package %s\n\nimport \"sync\"\n\ntype CountDownLatch struct {\n\tcount int\n\tmu sync.Mutex\n\tcond *sync.Cond\n}\n\nfunc NewCountDownLatch(n int) *CountDownLatch {\n\tl := &CountDownLatch{count: n}\n\tl.cond = sync.NewCond(&l.mu)\n\treturn l\n}\n", pkg)
}

// GenerateCondSimple 生成简化的Cond模式
func GenerateCondSimple(pkg string) string {
	return fmt.Sprintf("package %s\n\nimport \"sync\"\n\ntype BoundedQueue struct {\n\tmu sync.Mutex\n\tcond *sync.Cond\n\titems []interface{}\n\tmax int\n}\n\nfunc NewBoundedQueue(n int) *BoundedQueue {\n\tq := &BoundedQueue{max: n}\n\tq.cond = sync.NewCond(&q.mu)\n\treturn q\n}\n", pkg)
}

# 系统优化框架分析

<!-- TOC START -->
- [1 1 1 1 1 1 1 系统优化框架分析](#1-1-1-1-1-1-1-系统优化框架分析)
  - [1.1 1. 概述](#1-概述)
    - [1.1.1 正式定义](#正式定义)
    - [1.1.2 优化目标函数](#优化目标函数)
  - [1.2 2. 优化算法 (Optimization Algorithms)](#2-优化算法-optimization-algorithms)
    - [1.2.1 遗传算法 (Genetic Algorithm)](#遗传算法-genetic-algorithm)
      - [1.2.1.1 正式定义](#正式定义)
      - [1.2.1.2 Golang实现](#golang实现)
    - [1.2.2 粒子群优化 (Particle Swarm Optimization)](#粒子群优化-particle-swarm-optimization)
      - [1.2.2.1 正式定义](#正式定义)
      - [1.2.2.2 Golang实现](#golang实现)
  - [1.3 3. 资源分配策略 (Resource Allocation)](#3-资源分配策略-resource-allocation)
    - [1.3.1 动态资源分配](#动态资源分配)
      - [1.3.1.1 正式定义](#正式定义)
      - [1.3.1.2 Golang实现](#golang实现)
  - [1.4 4. 缓存机制 (Caching Mechanisms)](#4-缓存机制-caching-mechanisms)
    - [1.4.1 LRU缓存 (Least Recently Used)](#lru缓存-least-recently-used)
      - [1.4.1.1 正式定义](#正式定义)
      - [1.4.1.2 Golang实现](#golang实现)
  - [1.5 5. 负载均衡算法 (Load Balancing)](#5-负载均衡算法-load-balancing)
    - [1.5.1 一致性哈希 (Consistent Hashing)](#一致性哈希-consistent-hashing)
      - [1.5.1.1 正式定义](#正式定义)
      - [1.5.1.2 Golang实现](#golang实现)
  - [1.6 6. 容错模式 (Fault Tolerance)](#6-容错模式-fault-tolerance)
    - [1.6.1 重试模式 (Retry Pattern)](#重试模式-retry-pattern)
      - [1.6.1.1 正式定义](#正式定义)
      - [1.6.1.2 Golang实现](#golang实现)
  - [1.7 7. 监控系统 (Monitoring Systems)](#7-监控系统-monitoring-systems)
    - [1.7.1 指标收集 (Metrics Collection)](#指标收集-metrics-collection)
      - [1.7.1.1 正式定义](#正式定义)
      - [1.7.1.2 Golang实现](#golang实现)
  - [1.8 8. 性能分析](#8-性能分析)
    - [1.8.1 优化效果评估](#优化效果评估)
  - [1.9 9. 总结](#9-总结)
    - [1.9.1 关键成果](#关键成果)
    - [1.9.2 应用价值](#应用价值)
  - [1.10 参考文献](#参考文献)
<!-- TOC END -->

## 1.1 1. 概述

### 1.1.1 正式定义

系统优化框架是一个形式化的性能优化知识体系，定义为：

$$\mathcal{O} = (\mathcal{A}, \mathcal{R}, \mathcal{C}, \mathcal{L}, \mathcal{F}, \mathcal{M})$$

其中：

- $\mathcal{A}$: 优化算法集合 (Optimization Algorithms)
- $\mathcal{R}$: 资源分配策略集合 (Resource Allocation)
- $\mathcal{C}$: 缓存机制集合 (Caching Mechanisms)
- $\mathcal{L}$: 负载均衡算法集合 (Load Balancing)
- $\mathcal{F}$: 容错模式集合 (Fault Tolerance)
- $\mathcal{M}$: 监控系统集合 (Monitoring Systems)

### 1.1.2 优化目标函数

$$\text{Optimize}(S) = \arg\min_{s \in S} \left( \alpha \cdot \text{Latency}(s) + \beta \cdot \text{Resource}(s) + \gamma \cdot \text{Cost}(s) \right)$$

其中 $\alpha, \beta, \gamma$ 是权重系数，满足 $\alpha + \beta + \gamma = 1$。

## 1.2 2. 优化算法 (Optimization Algorithms)

### 1.2.1 遗传算法 (Genetic Algorithm)

#### 1.2.1.1 正式定义

遗传算法模拟自然选择过程：

$$\text{GA}(P, G) = \text{Selection}(P) \circ \text{Crossover}(P) \circ \text{Mutation}(P)$$

其中 $P$ 是种群，$G$ 是代数。

#### 1.2.1.2 Golang实现

```go
package genetic

import (
    "math/rand"
    "sort"
    "time"
)

// Individual represents a solution
type Individual struct {
    Genes   []float64
    Fitness float64
}

// GeneticAlgorithm implementation
type GeneticAlgorithm struct {
    PopulationSize int
    MutationRate   float64
    CrossoverRate  float64
    Generations    int
    Population     []Individual
}

func NewGeneticAlgorithm(popSize int, mutationRate, crossoverRate float64, generations int) *GeneticAlgorithm {
    return &GeneticAlgorithm{
        PopulationSize: popSize,
        MutationRate:   mutationRate,
        CrossoverRate:  crossoverRate,
        Generations:    generations,
        Population:     make([]Individual, popSize),
    }
}

func (ga *GeneticAlgorithm) Initialize() {
    for i := 0; i < ga.PopulationSize; i++ {
        ga.Population[i] = Individual{
            Genes:   make([]float64, 10), // 10-dimensional problem
            Fitness: 0,
        }
        
        // Random initialization
        for j := range ga.Population[i].Genes {
            ga.Population[i].Genes[j] = rand.Float64() * 100
        }
    }
}

func (ga *GeneticAlgorithm) Evaluate() {
    for i := range ga.Population {
        ga.Population[i].Fitness = ga.fitnessFunction(ga.Population[i].Genes)
    }
}

func (ga *GeneticAlgorithm) fitnessFunction(genes []float64) float64 {
    // Example: minimize sum of squares
    sum := 0.0
    for _, gene := range genes {
        sum += gene * gene
    }
    return sum
}

func (ga *GeneticAlgorithm) Selection() []Individual {
    // Tournament selection
    selected := make([]Individual, ga.PopulationSize)
    
    for i := 0; i < ga.PopulationSize; i++ {
        // Select 3 random individuals
        tournament := make([]Individual, 3)
        for j := 0; j < 3; j++ {
            idx := rand.Intn(ga.PopulationSize)
            tournament[j] = ga.Population[idx]
        }
        
        // Select the best from tournament
        best := tournament[0]
        for _, individual := range tournament[1:] {
            if individual.Fitness < best.Fitness {
                best = individual
            }
        }
        selected[i] = best
    }
    
    return selected
}

func (ga *GeneticAlgorithm) Crossover(parents []Individual) []Individual {
    offspring := make([]Individual, ga.PopulationSize)
    
    for i := 0; i < ga.PopulationSize; i += 2 {
        if i+1 < ga.PopulationSize && rand.Float64() < ga.CrossoverRate {
            // Single-point crossover
            crossoverPoint := rand.Intn(len(parents[i].Genes))
            
            offspring[i] = Individual{
                Genes:   make([]float64, len(parents[i].Genes)),
                Fitness: 0,
            }
            offspring[i+1] = Individual{
                Genes:   make([]float64, len(parents[i].Genes)),
                Fitness: 0,
            }
            
            // Copy genes before crossover point
            copy(offspring[i].Genes[:crossoverPoint], parents[i].Genes[:crossoverPoint])
            copy(offspring[i+1].Genes[:crossoverPoint], parents[i+1].Genes[:crossoverPoint])
            
            // Copy genes after crossover point
            copy(offspring[i].Genes[crossoverPoint:], parents[i+1].Genes[crossoverPoint:])
            copy(offspring[i+1].Genes[crossoverPoint:], parents[i].Genes[crossoverPoint:])
        } else {
            offspring[i] = parents[i]
            if i+1 < ga.PopulationSize {
                offspring[i+1] = parents[i+1]
            }
        }
    }
    
    return offspring
}

func (ga *GeneticAlgorithm) Mutation(individuals []Individual) {
    for i := range individuals {
        for j := range individuals[i].Genes {
            if rand.Float64() < ga.MutationRate {
                // Gaussian mutation
                individuals[i].Genes[j] += rand.NormFloat64() * 10
            }
        }
    }
}

func (ga *GeneticAlgorithm) Run() Individual {
    rand.Seed(time.Now().UnixNano())
    ga.Initialize()
    
    bestIndividual := Individual{Fitness: float64(^uint(0) >> 1)}
    
    for generation := 0; generation < ga.Generations; generation++ {
        ga.Evaluate()
        
        // Find best individual
        for _, individual := range ga.Population {
            if individual.Fitness < bestIndividual.Fitness {
                bestIndividual = individual
            }
        }
        
        // Selection
        parents := ga.Selection()
        
        // Crossover
        offspring := ga.Crossover(parents)
        
        // Mutation
        ga.Mutation(offspring)
        
        // Update population
        ga.Population = offspring
    }
    
    return bestIndividual
}

// Usage example
func Example() {
    ga := NewGeneticAlgorithm(100, 0.01, 0.8, 100)
    best := ga.Run()
    
    fmt.Printf("Best fitness: %f\n", best.Fitness)
    fmt.Printf("Best genes: %v\n", best.Genes)
}

```

### 1.2.2 粒子群优化 (Particle Swarm Optimization)

#### 1.2.2.1 正式定义

粒子群优化算法模拟群体行为：

$$v_i^{t+1} = w \cdot v_i^t + c_1 \cdot r_1 \cdot (p_i - x_i^t) + c_2 \cdot r_2 \cdot (g - x_i^t)$$

$$x_i^{t+1} = x_i^t + v_i^{t+1}$$

其中 $v_i$ 是速度，$x_i$ 是位置，$p_i$ 是个人最佳，$g$ 是全局最佳。

#### 1.2.2.2 Golang实现

```go
package pso

import (
    "math"
    "math/rand"
    "time"
)

// Particle represents a solution
type Particle struct {
    Position     []float64
    Velocity     []float64
    BestPosition []float64
    BestFitness  float64
}

// PSO implementation
type PSO struct {
    Particles    []Particle
    GlobalBest   []float64
    GlobalFitness float64
    Dimensions   int
    ParticleCount int
    MaxIterations int
    W, C1, C2    float64
}

func NewPSO(dimensions, particleCount, maxIterations int, w, c1, c2 float64) *PSO {
    pso := &PSO{
        Dimensions:     dimensions,
        ParticleCount:  particleCount,
        MaxIterations:  maxIterations,
        W:             w,
        C1:            c1,
        C2:            c2,
        Particles:     make([]Particle, particleCount),
        GlobalBest:    make([]float64, dimensions),
        GlobalFitness: math.Inf(1),
    }
    
    pso.initialize()
    return pso
}

func (pso *PSO) initialize() {
    for i := 0; i < pso.ParticleCount; i++ {
        pso.Particles[i] = Particle{
            Position:     make([]float64, pso.Dimensions),
            Velocity:     make([]float64, pso.Dimensions),
            BestPosition: make([]float64, pso.Dimensions),
            BestFitness:  math.Inf(1),
        }
        
        // Random initialization
        for j := 0; j < pso.Dimensions; j++ {
            pso.Particles[i].Position[j] = rand.Float64()*200 - 100
            pso.Particles[i].Velocity[j] = rand.Float64()*2 - 1
        }
    }
}

func (pso *PSO) fitnessFunction(position []float64) float64 {
    // Example: Sphere function
    sum := 0.0
    for _, x := range position {
        sum += x * x
    }
    return sum
}

func (pso *PSO) updateParticle(particle *Particle) {
    for i := 0; i < pso.Dimensions; i++ {
        r1 := rand.Float64()
        r2 := rand.Float64()
        
        // Update velocity
        particle.Velocity[i] = pso.W*particle.Velocity[i] +
            pso.C1*r1*(particle.BestPosition[i]-particle.Position[i]) +
            pso.C2*r2*(pso.GlobalBest[i]-particle.Position[i])
        
        // Update position
        particle.Position[i] += particle.Velocity[i]
    }
}

func (pso *PSO) Run() ([]float64, float64) {
    rand.Seed(time.Now().UnixNano())
    
    for iteration := 0; iteration < pso.MaxIterations; iteration++ {
        // Evaluate all particles
        for i := range pso.Particles {
            fitness := pso.fitnessFunction(pso.Particles[i].Position)
            
            // Update personal best
            if fitness < pso.Particles[i].BestFitness {
                pso.Particles[i].BestFitness = fitness
                copy(pso.Particles[i].BestPosition, pso.Particles[i].Position)
                
                // Update global best
                if fitness < pso.GlobalFitness {
                    pso.GlobalFitness = fitness
                    copy(pso.GlobalBest, pso.Particles[i].Position)
                }
            }
        }
        
        // Update particles
        for i := range pso.Particles {
            pso.updateParticle(&pso.Particles[i])
        }
    }
    
    return pso.GlobalBest, pso.GlobalFitness
}

// Usage example
func Example() {
    pso := NewPSO(10, 50, 100, 0.7, 1.5, 1.5)
    bestPosition, bestFitness := pso.Run()
    
    fmt.Printf("Best position: %v\n", bestPosition)
    fmt.Printf("Best fitness: %f\n", bestFitness)
}

```

## 1.3 3. 资源分配策略 (Resource Allocation)

### 1.3.1 动态资源分配

#### 1.3.1.1 正式定义

动态资源分配基于负载变化调整资源：

$$\text{Allocate}(R, L) = \arg\min_{r \in R} \left( \text{Cost}(r) + \lambda \cdot \text{Load}(L) \right)$$

其中 $R$ 是资源集合，$L$ 是负载，$\lambda$ 是权重。

#### 1.3.1.2 Golang实现

```go
package resourceallocation

import (
    "fmt"
    "sync"
    "time"
)

// Resource represents a computing resource
type Resource struct {
    ID       string
    CPU      float64
    Memory   float64
    Network  float64
    Cost     float64
    Load     float64
    mu       sync.RWMutex
}

// ResourcePool manages resource allocation
type ResourcePool struct {
    Resources map[string]*Resource
    mu        sync.RWMutex
}

func NewResourcePool() *ResourcePool {
    return &ResourcePool{
        Resources: make(map[string]*Resource),
    }
}

func (rp *ResourcePool) AddResource(resource *Resource) {
    rp.mu.Lock()
    defer rp.mu.Unlock()
    rp.Resources[resource.ID] = resource
}

func (rp *ResourcePool) AllocateResource(requiredCPU, requiredMemory float64) (*Resource, error) {
    rp.mu.RLock()
    defer rp.mu.RUnlock()
    
    var bestResource *Resource
    bestScore := float64(^uint(0) >> 1)
    
    for _, resource := range rp.Resources {
        resource.mu.RLock()
        
        if resource.CPU >= requiredCPU && resource.Memory >= requiredMemory {
            // Calculate allocation score (lower is better)
            score := resource.Cost + resource.Load*10
            
            if score < bestScore {
                bestScore = score
                bestResource = resource
            }
        }
        
        resource.mu.RUnlock()
    }
    
    if bestResource == nil {
        return nil, fmt.Errorf("no suitable resource found")
    }
    
    // Update resource load
    bestResource.mu.Lock()
    bestResource.Load += (requiredCPU + requiredMemory) / (bestResource.CPU + bestResource.Memory)
    bestResource.mu.Unlock()
    
    return bestResource, nil
}

func (rp *ResourcePool) ReleaseResource(resourceID string, cpu, memory float64) error {
    rp.mu.RLock()
    resource, exists := rp.Resources[resourceID]
    rp.mu.RUnlock()
    
    if !exists {
        return fmt.Errorf("resource %s not found", resourceID)
    }
    
    resource.mu.Lock()
    defer resource.mu.Unlock()
    
    loadReduction := (cpu + memory) / (resource.CPU + resource.Memory)
    resource.Load = math.Max(0, resource.Load-loadReduction)
    
    return nil
}

// LoadBalancer for resource allocation
type LoadBalancer struct {
    pool *ResourcePool
}

func NewLoadBalancer(pool *ResourcePool) *LoadBalancer {
    return &LoadBalancer{pool: pool}
}

func (lb *LoadBalancer) RoundRobin(requests []ResourceRequest) []*Resource {
    resources := make([]*Resource, len(requests))
    resourceList := make([]*Resource, 0, len(lb.pool.Resources))
    
    lb.pool.mu.RLock()
    for _, resource := range lb.pool.Resources {
        resourceList = append(resourceList, resource)
    }
    lb.pool.mu.RUnlock()
    
    for i, request := range requests {
        resource, err := lb.pool.AllocateResource(request.CPU, request.Memory)
        if err != nil {
            // Fallback to round-robin
            resource = resourceList[i%len(resourceList)]
        }
        resources[i] = resource
    }
    
    return resources
}

type ResourceRequest struct {
    CPU    float64
    Memory float64
}

// Usage example
func Example() {
    pool := NewResourcePool()
    
    // Add resources
    resources := []*Resource{
        {ID: "r1", CPU: 4, Memory: 8, Cost: 10},
        {ID: "r2", CPU: 8, Memory: 16, Cost: 20},
        {ID: "r3", CPU: 2, Memory: 4, Cost: 5},
    }
    
    for _, resource := range resources {
        pool.AddResource(resource)
    }
    
    // Create load balancer
    lb := NewLoadBalancer(pool)
    
    // Allocate resources
    requests := []ResourceRequest{
        {CPU: 1, Memory: 2},
        {CPU: 2, Memory: 4},
        {CPU: 1, Memory: 1},
    }
    
    allocated := lb.RoundRobin(requests)
    
    for i, resource := range allocated {
        fmt.Printf("Request %d allocated to resource %s\n", i, resource.ID)
    }
}

```

## 1.4 4. 缓存机制 (Caching Mechanisms)

### 1.4.1 LRU缓存 (Least Recently Used)

#### 1.4.1.1 正式定义

LRU缓存基于最近使用时间进行淘汰：

$$\text{LRU}(C, K) = \arg\min_{k \in K} \text{LastAccess}(k)$$

其中 $C$ 是缓存，$K$ 是键集合。

#### 1.4.1.2 Golang实现

```go
package cache

import (
    "container/list"
    "sync"
    "time"
)

// LRUCache implementation
type LRUCache struct {
    capacity int
    cache    map[string]*list.Element
    list     *list.List
    mu       sync.RWMutex
}

type CacheEntry struct {
    Key       string
    Value     interface{}
    Timestamp time.Time
}

func NewLRUCache(capacity int) *LRUCache {
    return &LRUCache{
        capacity: capacity,
        cache:    make(map[string]*list.Element),
        list:     list.New(),
    }
}

func (lru *LRUCache) Get(key string) (interface{}, bool) {
    lru.mu.Lock()
    defer lru.mu.Unlock()
    
    if element, exists := lru.cache[key]; exists {
        // Move to front (most recently used)
        lru.list.MoveToFront(element)
        entry := element.Value.(*CacheEntry)
        entry.Timestamp = time.Now()
        return entry.Value, true
    }
    
    return nil, false
}

func (lru *LRUCache) Put(key string, value interface{}) {
    lru.mu.Lock()
    defer lru.mu.Unlock()
    
    if element, exists := lru.cache[key]; exists {
        // Update existing entry
        lru.list.MoveToFront(element)
        entry := element.Value.(*CacheEntry)
        entry.Value = value
        entry.Timestamp = time.Now()
    } else {
        // Add new entry
        entry := &CacheEntry{
            Key:       key,
            Value:     value,
            Timestamp: time.Now(),
        }
        
        element := lru.list.PushFront(entry)
        lru.cache[key] = element
        
        // Remove oldest if capacity exceeded
        if lru.list.Len() > lru.capacity {
            oldest := lru.list.Back()
            lru.list.Remove(oldest)
            delete(lru.cache, oldest.Value.(*CacheEntry).Key)
        }
    }
}

func (lru *LRUCache) Remove(key string) bool {
    lru.mu.Lock()
    defer lru.mu.Unlock()
    
    if element, exists := lru.cache[key]; exists {
        lru.list.Remove(element)
        delete(lru.cache, key)
        return true
    }
    
    return false
}

func (lru *LRUCache) Size() int {
    lru.mu.RLock()
    defer lru.mu.RUnlock()
    return lru.list.Len()
}

// TTL Cache with expiration
type TTLCache struct {
    cache map[string]*TTLEntry
    mu    sync.RWMutex
}

type TTLEntry struct {
    Value      interface{}
    Expiration time.Time
}

func NewTTLCache() *TTLCache {
    cache := &TTLCache{
        cache: make(map[string]*TTLEntry),
    }
    
    // Start cleanup goroutine
    go cache.cleanup()
    
    return cache
}

func (ttl *TTLCache) Put(key string, value interface{}, ttl time.Duration) {
    ttl.mu.Lock()
    defer ttl.mu.Unlock()
    
    ttl.cache[key] = &TTLEntry{
        Value:      value,
        Expiration: time.Now().Add(ttl),
    }
}

func (ttl *TTLCache) Get(key string) (interface{}, bool) {
    ttl.mu.Lock()
    defer ttl.mu.Unlock()
    
    if entry, exists := ttl.cache[key]; exists {
        if time.Now().Before(entry.Expiration) {
            return entry.Value, true
        } else {
            // Expired, remove it
            delete(ttl.cache, key)
        }
    }
    
    return nil, false
}

func (ttl *TTLCache) cleanup() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        ttl.mu.Lock()
        now := time.Now()
        for key, entry := range ttl.cache {
            if now.After(entry.Expiration) {
                delete(ttl.cache, key)
            }
        }
        ttl.mu.Unlock()
    }
}

// Usage example
func Example() {
    // LRU Cache
    lru := NewLRUCache(3)
    
    lru.Put("a", 1)
    lru.Put("b", 2)
    lru.Put("c", 3)
    
    fmt.Printf("LRU size: %d\n", lru.Size())
    
    if value, exists := lru.Get("a"); exists {
        fmt.Printf("LRU get 'a': %v\n", value)
    }
    
    // TTL Cache
    ttl := NewTTLCache()
    ttl.Put("key", "value", 5*time.Second)
    
    if value, exists := ttl.Get("key"); exists {
        fmt.Printf("TTL get 'key': %v\n", value)
    }
    
    time.Sleep(6 * time.Second)
    
    if value, exists := ttl.Get("key"); exists {
        fmt.Printf("TTL get 'key' after expiration: %v\n", value)
    } else {
        fmt.Println("TTL key expired")
    }
}

```

## 1.5 5. 负载均衡算法 (Load Balancing)

### 1.5.1 一致性哈希 (Consistent Hashing)

#### 1.5.1.1 正式定义

一致性哈希将节点映射到哈希环上：

$$\text{Hash}(key) \rightarrow [0, 2^{32})$$

$$\text{Node}(key) = \arg\min_{n \in N} \text{Hash}(n) \geq \text{Hash}(key)$$

其中 $N$ 是节点集合。

#### 1.5.1.2 Golang实现

```go
package loadbalancer

import (
    "crypto/md5"
    "fmt"
    "sort"
    "sync"
)

// ConsistentHash implementation
type ConsistentHash struct {
    nodes    map[uint32]string
    sorted   []uint32
    replicas int
    mu       sync.RWMutex
}

func NewConsistentHash(replicas int) *ConsistentHash {
    return &ConsistentHash{
        nodes:    make(map[uint32]string),
        replicas: replicas,
    }
}

func (ch *ConsistentHash) hash(key string) uint32 {
    hash := md5.Sum([]byte(key))
    return uint32(hash[0])<<24 | uint32(hash[1])<<16 | uint32(hash[2])<<8 | uint32(hash[3])
}

func (ch *ConsistentHash) AddNode(node string) {
    ch.mu.Lock()
    defer ch.mu.Unlock()
    
    for i := 0; i < ch.replicas; i++ {
        virtualNode := fmt.Sprintf("%s-%d", node, i)
        hash := ch.hash(virtualNode)
        ch.nodes[hash] = node
    }
    
    // Update sorted list
    ch.updateSorted()
}

func (ch *ConsistentHash) RemoveNode(node string) {
    ch.mu.Lock()
    defer ch.mu.Unlock()
    
    for i := 0; i < ch.replicas; i++ {
        virtualNode := fmt.Sprintf("%s-%d", node, i)
        hash := ch.hash(virtualNode)
        delete(ch.nodes, hash)
    }
    
    // Update sorted list
    ch.updateSorted()
}

func (ch *ConsistentHash) GetNode(key string) string {
    ch.mu.RLock()
    defer ch.mu.RUnlock()
    
    if len(ch.nodes) == 0 {
        return ""
    }
    
    hash := ch.hash(key)
    
    // Find the first node with hash >= key hash
    idx := sort.Search(len(ch.sorted), func(i int) bool {
        return ch.sorted[i] >= hash
    })
    
    // Wrap around if necessary
    if idx == len(ch.sorted) {
        idx = 0
    }
    
    return ch.nodes[ch.sorted[idx]]
}

func (ch *ConsistentHash) updateSorted() {
    ch.sorted = make([]uint32, 0, len(ch.nodes))
    for hash := range ch.nodes {
        ch.sorted = append(ch.sorted, hash)
    }
    sort.Slice(ch.sorted, func(i, j int) bool {
        return ch.sorted[i] < ch.sorted[j]
    })
}

// Weighted Round Robin
type WeightedRoundRobin struct {
    nodes    []*WeightedNode
    current  int
    mu       sync.Mutex
}

type WeightedNode struct {
    Name     string
    Weight   int
    Current  int
}

func NewWeightedRoundRobin() *WeightedRoundRobin {
    return &WeightedRoundRobin{
        nodes: make([]*WeightedNode, 0),
    }
}

func (wrr *WeightedRoundRobin) AddNode(name string, weight int) {
    wrr.mu.Lock()
    defer wrr.mu.Unlock()
    
    wrr.nodes = append(wrr.nodes, &WeightedNode{
        Name:   name,
        Weight: weight,
        Current: 0,
    })
}

func (wrr *WeightedRoundRobin) GetNode() string {
    wrr.mu.Lock()
    defer wrr.mu.Unlock()
    
    if len(wrr.nodes) == 0 {
        return ""
    }
    
    // Find node with highest current weight
    maxWeight := -1
    selectedNode := 0
    
    for i, node := range wrr.nodes {
        if node.Current > maxWeight {
            maxWeight = node.Current
            selectedNode = i
        }
    }
    
    // Update weights
    for i, node := range wrr.nodes {
        if i == selectedNode {
            node.Current -= wrr.getTotalWeight()
        }
        node.Current += node.Weight
    }
    
    return wrr.nodes[selectedNode].Name
}

func (wrr *WeightedRoundRobin) getTotalWeight() int {
    total := 0
    for _, node := range wrr.nodes {
        total += node.Weight
    }
    return total
}

// Usage example
func Example() {
    // Consistent Hashing
    ch := NewConsistentHash(3)
    
    ch.AddNode("node1")
    ch.AddNode("node2")
    ch.AddNode("node3")
    
    keys := []string{"key1", "key2", "key3", "key4", "key5"}
    
    fmt.Println("Consistent Hashing:")
    for _, key := range keys {
        node := ch.GetNode(key)
        fmt.Printf("Key %s -> Node %s\n", key, node)
    }
    
    // Weighted Round Robin
    wrr := NewWeightedRoundRobin()
    wrr.AddNode("server1", 3)
    wrr.AddNode("server2", 2)
    wrr.AddNode("server3", 1)
    
    fmt.Println("\nWeighted Round Robin:")
    for i := 0; i < 10; i++ {
        node := wrr.GetNode()
        fmt.Printf("Request %d -> %s\n", i+1, node)
    }
}

```

## 1.6 6. 容错模式 (Fault Tolerance)

### 1.6.1 重试模式 (Retry Pattern)

#### 1.6.1.1 正式定义

重试模式在失败时重复执行操作：

$$\text{Retry}(f, n) = \begin{cases}
f() & \text{if } f() \text{ succeeds} \\
\text{Retry}(f, n-1) & \text{if } n > 0 \text{ and } f() \text{ fails}
\end{cases}$$

#### 1.6.1.2 Golang实现

```go
package faulttolerance

import (
    "context"
    "fmt"
    "math"
    "time"
)

// RetryPolicy defines retry behavior
type RetryPolicy struct {
    MaxAttempts     int
    InitialDelay    time.Duration
    MaxDelay        time.Duration
    BackoffMultiplier float64
    RetryableErrors  []error
}

// RetryableOperation represents an operation that can be retried
type RetryableOperation func() error

// RetryWithBackoff implements exponential backoff retry
func RetryWithBackoff(ctx context.Context, operation RetryableOperation, policy RetryPolicy) error {
    var lastErr error
    delay := policy.InitialDelay

    for attempt := 0; attempt < policy.MaxAttempts; attempt++ {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        if err := operation(); err == nil {
            return nil
        } else {
            lastErr = err

            // Check if error is retryable
            if !isRetryableError(err, policy.RetryableErrors) {
                return err
            }

            // Wait before next attempt
            if attempt < policy.MaxAttempts-1 {
                select {
                case <-time.After(delay):
                case <-ctx.Done():
                    return ctx.Err()
                }

                // Calculate next delay with exponential backoff
                delay = time.Duration(float64(delay) * policy.BackoffMultiplier)
                if delay > policy.MaxDelay {
                    delay = policy.MaxDelay
                }
            }
        }
    }

    return fmt.Errorf("operation failed after %d attempts: %v", policy.MaxAttempts, lastErr)
}

func isRetryableError(err error, retryableErrors []error) bool {
    if len(retryableErrors) == 0 {
        return true // Default to retryable
    }

    for _, retryableErr := range retryableErrors {
        if err == retryableErr {
            return true
        }
    }
    return false
}

// Circuit Breaker with retry
type CircuitBreakerWithRetry struct {
    breaker *CircuitBreaker
    policy  RetryPolicy
}

func NewCircuitBreakerWithRetry(failureThreshold int, timeout time.Duration, policy RetryPolicy) *CircuitBreakerWithRetry {
    return &CircuitBreakerWithRetry{
        breaker: NewCircuitBreaker(failureThreshold, timeout),
        policy:  policy,
    }
}

func (cb *CircuitBreakerWithRetry) Execute(ctx context.Context, operation RetryableOperation) error {
    return cb.breaker.Execute(ctx, func() error {
        return RetryWithBackoff(ctx, operation, cb.policy)
    })
}

// Bulkhead pattern for resource isolation
type Bulkhead struct {
    maxConcurrency int
    semaphore      chan struct{}
}

func NewBulkhead(maxConcurrency int) *Bulkhead {
    return &Bulkhead{
        maxConcurrency: maxConcurrency,
        semaphore:      make(chan struct{}, maxConcurrency),
    }
}

func (b *Bulkhead) Execute(ctx context.Context, operation func() error) error {
    select {
    case b.semaphore <- struct{}{}:
        defer func() { <-b.semaphore }()
        return operation()
    case <-ctx.Done():
        return ctx.Err()
    }
}

// Usage example
func Example() {
    // Retry policy
    policy := RetryPolicy{
        MaxAttempts:      3,
        InitialDelay:     100 * time.Millisecond,
        MaxDelay:         1 * time.Second,
        BackoffMultiplier: 2.0,
    }

    // Simulate flaky operation
    operation := func() error {
        if time.Now().UnixNano()%3 == 0 {
            return fmt.Errorf("temporary error")
        }
        return nil
    }

    ctx := context.Background()

    // Retry with backoff
    err := RetryWithBackoff(ctx, operation, policy)
    if err != nil {
        fmt.Printf("Retry failed: %v\n", err)
    } else {
        fmt.Println("Retry succeeded")
    }

    // Circuit breaker with retry
    cb := NewCircuitBreakerWithRetry(3, 5*time.Second, policy)

    err = cb.Execute(ctx, operation)
    if err != nil {
        fmt.Printf("Circuit breaker with retry failed: %v\n", err)
    } else {
        fmt.Println("Circuit breaker with retry succeeded")
    }

    // Bulkhead
    bulkhead := NewBulkhead(2)

    for i := 0; i < 5; i++ {
        go func(id int) {
            err := bulkhead.Execute(ctx, func() error {
                time.Sleep(100 * time.Millisecond)
                return nil
            })
            if err != nil {
                fmt.Printf("Bulkhead operation %d failed: %v\n", id, err)
            } else {
                fmt.Printf("Bulkhead operation %d succeeded\n", id)
            }
        }(i)
    }

    time.Sleep(1 * time.Second)
}

```

## 1.7 7. 监控系统 (Monitoring Systems)

### 1.7.1 指标收集 (Metrics Collection)

#### 1.7.1.1 正式定义

指标收集系统聚合性能数据：

$$\text{Metrics}(T) = \sum_{t \in T} \text{Collect}(t)$$

其中 $T$ 是时间序列。

#### 1.7.1.2 Golang实现

```go
package monitoring

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Metric types
type MetricType int

const (
    Counter MetricType = iota
    Gauge
    Histogram
    Summary
)

// Metric represents a single metric
type Metric struct {
    Name   string
    Type   MetricType
    Value  float64
    Labels map[string]string
    Time   time.Time
}

// MetricsCollector collects and stores metrics
type MetricsCollector struct {
    metrics map[string][]Metric
    mu      sync.RWMutex
}

func NewMetricsCollector() *MetricsCollector {
    return &MetricsCollector{
        metrics: make(map[string][]Metric),
    }
}

func (mc *MetricsCollector) Record(name string, metricType MetricType, value float64, labels map[string]string) {
    mc.mu.Lock()
    defer mc.mu.Unlock()

    metric := Metric{
        Name:   name,
        Type:   metricType,
        Value:  value,
        Labels: labels,
        Time:   time.Now(),
    }

    mc.metrics[name] = append(mc.metrics[name], metric)
}

func (mc *MetricsCollector) GetMetrics(name string) []Metric {
    mc.mu.RLock()
    defer mc.mu.RUnlock()

    if metrics, exists := mc.metrics[name]; exists {
        return append([]Metric{}, metrics...)
    }
    return nil
}

func (mc *MetricsCollector) GetLatestMetric(name string) (Metric, bool) {
    mc.mu.RLock()
    defer mc.mu.RUnlock()

    if metrics, exists := mc.metrics[name]; exists && len(metrics) > 0 {
        return metrics[len(metrics)-1], true
    }
    return Metric{}, false
}

// Prometheus-style metrics
type PrometheusMetrics struct {
    counters   map[string]float64
    gauges     map[string]float64
    histograms map[string][]float64
    mu         sync.RWMutex
}

func NewPrometheusMetrics() *PrometheusMetrics {
    return &PrometheusMetrics{
        counters:   make(map[string]float64),
        gauges:     make(map[string]float64),
        histograms: make(map[string][]float64),
    }
}

func (pm *PrometheusMetrics) IncrementCounter(name string) {
    pm.mu.Lock()
    defer pm.mu.Unlock()
    pm.counters[name]++
}

func (pm *PrometheusMetrics) SetGauge(name string, value float64) {
    pm.mu.Lock()
    defer pm.mu.Unlock()
    pm.gauges[name] = value
}

func (pm *PrometheusMetrics) ObserveHistogram(name string, value float64) {
    pm.mu.Lock()
    defer pm.mu.Unlock()
    pm.histograms[name] = append(pm.histograms[name], value)
}

func (pm *PrometheusMetrics) GetCounter(name string) float64 {
    pm.mu.RLock()
    defer pm.mu.RUnlock()
    return pm.counters[name]
}

func (pm *PrometheusMetrics) GetGauge(name string) float64 {
    pm.mu.RLock()
    defer pm.mu.RUnlock()
    return pm.gauges[name]
}

func (pm *PrometheusMetrics) GetHistogram(name string) []float64 {
    pm.mu.RLock()
    defer pm.mu.RUnlock()
    if values, exists := pm.histograms[name]; exists {
        return append([]float64{}, values...)
    }
    return nil
}

// Health check system
type HealthChecker struct {
    checks map[string]HealthCheck
    mu     sync.RWMutex
}

type HealthCheck func() error

func NewHealthChecker() *HealthChecker {
    return &HealthChecker{
        checks: make(map[string]HealthCheck),
    }
}

func (hc *HealthChecker) AddCheck(name string, check HealthCheck) {
    hc.mu.Lock()
    defer hc.mu.Unlock()
    hc.checks[name] = check
}

func (hc *HealthChecker) CheckHealth() map[string]error {
    hc.mu.RLock()
    defer hc.mu.RUnlock()

    results := make(map[string]error)

    for name, check := range hc.checks {
        results[name] = check()
    }

    return results
}

func (hc *HealthChecker) IsHealthy() bool {
    results := hc.CheckHealth()

    for _, err := range results {
        if err != nil {
            return false
        }
    }

    return true
}

// Usage example
func Example() {
    // Metrics collector
    collector := NewMetricsCollector()

    // Record some metrics
    collector.Record("http_requests_total", Counter, 1, map[string]string{"method": "GET", "status": "200"})
    collector.Record("http_request_duration", Histogram, 0.15, map[string]string{"method": "GET"})
    collector.Record("memory_usage", Gauge, 1024.5, map[string]string{"component": "api"})

    // Get metrics
    requests := collector.GetMetrics("http_requests_total")
    fmt.Printf("HTTP requests: %+v\n", requests)

    // Prometheus metrics
    prom := NewPrometheusMetrics()

    prom.IncrementCounter("requests_total")
    prom.SetGauge("memory_usage", 1024.5)
    prom.ObserveHistogram("request_duration", 0.15)

    fmt.Printf("Counter: %f\n", prom.GetCounter("requests_total"))
    fmt.Printf("Gauge: %f\n", prom.GetGauge("memory_usage"))
    fmt.Printf("Histogram: %v\n", prom.GetHistogram("request_duration"))

    // Health checker
    checker := NewHealthChecker()

    checker.AddCheck("database", func() error {
        // Simulate database check
        return nil
    })

    checker.AddCheck("cache", func() error {
        // Simulate cache check
        return nil
    })

    if checker.IsHealthy() {
        fmt.Println("System is healthy")
    } else {
        fmt.Println("System has health issues")
    }
}

```

## 1.8 8. 性能分析

### 1.8.1 优化效果评估

```go
// Performance metrics
type PerformanceMetrics struct {
    Latency    time.Duration
    Throughput float64
    ResourceUsage float64
    Cost       float64
}

// Optimization comparison
type OptimizationComparison struct {
    Before PerformanceMetrics
    After  PerformanceMetrics
    Improvement float64
}

func CalculateImprovement(before, after PerformanceMetrics) OptimizationComparison {
    latencyImprovement := float64(before.Latency-after.Latency) / float64(before.Latency) * 100
    throughputImprovement := (after.Throughput - before.Throughput) / before.Throughput * 100
    resourceImprovement := (before.ResourceUsage - after.ResourceUsage) / before.ResourceUsage * 100
    costImprovement := (before.Cost - after.Cost) / before.Cost * 100

    overallImprovement := (latencyImprovement + throughputImprovement + resourceImprovement + costImprovement) / 4

    return OptimizationComparison{
        Before:      before,
        After:       after,
        Improvement: overallImprovement,
    }
}

```

## 1.9 9. 总结

系统优化框架提供了一个完整的性能优化知识体系，涵盖了从算法优化到系统监控的各个方面。通过形式化的数学定义、完整的Golang实现和详细的性能分析，为系统优化提供了坚实的理论基础和实践指导。

### 1.9.1 关键成果

1. **优化算法**: 遗传算法、粒子群优化等高级优化技术
2. **资源管理**: 动态资源分配、负载均衡策略
3. **缓存机制**: LRU缓存、TTL缓存等缓存策略
4. **容错模式**: 重试模式、熔断器、隔离模式
5. **监控系统**: 指标收集、健康检查、性能监控

### 1.9.2 应用价值

- **性能提升**: 通过算法优化和资源管理提升系统性能
- **可靠性增强**: 通过容错模式提高系统可靠性
- **成本优化**: 通过资源优化降低运营成本
- **可观测性**: 通过监控系统提供系统可观测性

## 1.10 参考文献

1. Goldberg, D. E. (1989). Genetic Algorithms in Search, Optimization and Machine Learning
2. Kennedy, J., & Eberhart, R. (1995). Particle Swarm Optimization
3. Karger, D., et al. (1997). Consistent Hashing and Random Trees
4. Go Performance: <https://golang.org/doc/effective_go.html#performance>
5. Go Profiling: <https://golang.org/pkg/runtime/pprof/>

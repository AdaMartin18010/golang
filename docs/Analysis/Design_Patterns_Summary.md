# Design Pattern System Summary

## Core Pattern Categories

### 1. Creational Patterns

**Factory Pattern**

```go
type Product interface {
    Operation() string
}

type Creator interface {
    CreateProduct() Product
}

func NewCreator(creatorType string) Creator {
    switch creatorType {
    case "A":
        return &ConcreteCreatorA{}
    case "B":
        return &ConcreteCreatorB{}
    default:
        return &ConcreteCreatorA{}
    }
}
```

**Singleton Pattern**

```go
var (
    instance *Singleton
    once     sync.Once
)

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{data: "Singleton instance"}
    })
    return instance
}
```

**Builder Pattern**

```go
type ComputerBuilder interface {
    SetCPU(cpu string) ComputerBuilder
    SetMemory(memory string) ComputerBuilder
    SetStorage(storage string) ComputerBuilder
    Build() *Computer
}

func (cb *ConcreteComputerBuilder) SetCPU(cpu string) ComputerBuilder {
    cb.computer.CPU = cpu
    return cb
}
```

### 2. Structural Patterns

**Adapter Pattern**

```go
type Target interface {
    Request() string
}

type Adapter struct {
    adaptee *Adaptee
}

func (a *Adapter) Request() string {
    return "Adapter: " + a.adaptee.SpecificRequest()
}
```

**Decorator Pattern**

```go
type Component interface {
    Operation() string
}

type Decorator struct {
    component Component
}

func (d *Decorator) Operation() string {
    return d.component.Operation()
}

func (cda *ConcreteDecoratorA) Operation() string {
    return "ConcreteDecoratorA(" + cda.component.Operation() + ")"
}
```

**Proxy Pattern**

```go
type Proxy struct {
    realSubject *RealSubject
    accessCount int
    mutex       sync.Mutex
}

func (p *Proxy) Request() string {
    p.mutex.Lock()
    defer p.mutex.Unlock()
    
    if p.realSubject == nil {
        p.realSubject = &RealSubject{}
    }
    
    p.accessCount++
    return fmt.Sprintf("Proxy (access #%d): %s", p.accessCount, p.realSubject.Request())
}
```

### 3. Behavioral Patterns

**Observer Pattern**

```go
type Observer interface {
    Update(data interface{})
}

type Subject interface {
    Attach(observer Observer)
    Detach(observer Observer)
    Notify()
}

func (cs *ConcreteSubject) Notify() {
    cs.mutex.RLock()
    observers := make([]Observer, len(cs.observers))
    copy(observers, cs.observers)
    cs.mutex.RUnlock()
    
    for _, observer := range observers {
        observer.Update(cs.data)
    }
}
```

**Strategy Pattern**

```go
type Strategy interface {
    Algorithm(data []int) []int
}

type Context struct {
    strategy Strategy
}

func (c *Context) ExecuteStrategy(data []int) []int {
    return c.strategy.Algorithm(data)
}
```

**Command Pattern**

```go
type Command interface {
    Execute()
    Undo()
}

type Invoker struct {
    commands []Command
    history  []Command
    mutex    sync.Mutex
}

func (i *Invoker) ExecuteCommands() {
    i.mutex.Lock()
    defer i.mutex.Unlock()
    
    for _, command := range i.commands {
        command.Execute()
        i.history = append(i.history, command)
    }
    i.commands = i.commands[:0]
}
```

### 4. Concurrency Patterns

**Worker Pool Pattern**

```go
type WorkerPool struct {
    workers    int
    taskQueue  chan Task
    resultChan chan error
    wg         sync.WaitGroup
    ctx        context.Context
    cancel     context.CancelFunc
}

func (wp *WorkerPool) worker(id int) {
    defer wp.wg.Done()
    
    for {
        select {
        case task := <-wp.taskQueue:
            if task == nil {
                return
            }
            err := task.Execute()
            wp.resultChan <- err
        case <-wp.ctx.Done():
            return
        }
    }
}
```

**Pipeline Pattern**

```go
type Pipeline struct {
    stages []Stage
}

func (p *Pipeline) Execute(data interface{}) (interface{}, error) {
    result := data
    
    for _, stage := range p.stages {
        var err error
        result, err = stage.Process(result)
        if err != nil {
            return nil, fmt.Errorf("pipeline stage failed: %w", err)
        }
        
        if result == nil {
            return nil, nil
        }
    }
    
    return result, nil
}
```

### 5. Enterprise Patterns

**Repository Pattern**

```go
type Repository[T any] interface {
    FindByID(id string) (T, error)
    FindAll() ([]T, error)
    Save(entity T) error
    Delete(id string) error
    Update(entity T) error
}

type InMemoryRepository[T any] struct {
    data map[string]T
    mutex sync.RWMutex
}

func (imr *InMemoryRepository[T]) FindByID(id string) (T, error) {
    imr.mutex.RLock()
    defer imr.mutex.RUnlock()
    
    entity, exists := imr.data[id]
    if !exists {
        var zero T
        return zero, fmt.Errorf("entity with id %s not found", id)
    }
    
    return entity, nil
}
```

**Unit of Work Pattern**

```go
type UnitOfWork interface {
    RegisterNew(entity interface{})
    RegisterDirty(entity interface{})
    RegisterDeleted(entity interface{})
    Commit() error
    Rollback() error
}

func (uow *ConcreteUnitOfWork) Commit() error {
    uow.mutex.Lock()
    defer uow.mutex.Unlock()
    
    // Process new entities
    for entityID, tracker := range uow.newEntities {
        repo := uow.getRepository(tracker.entity)
        if err := repo.Save(tracker.entity); err != nil {
            return fmt.Errorf("failed to save new entity %s: %w", entityID, err)
        }
    }
    
    // Process dirty and deleted entities...
    uow.clearTracking()
    return nil
}
```

## Performance Characteristics

| Pattern | Creation | Operation | Memory | Thread Safety |
|---------|----------|-----------|---------|---------------|
| Factory | O(1) | O(1) | O(1) | ✅ |
| Singleton | O(1) | O(1) | O(1) | ✅ |
| Builder | O(n) | O(1) | O(n) | ✅ |
| Adapter | O(1) | O(1) | O(1) | ✅ |
| Decorator | O(1) | O(n) | O(n) | ✅ |
| Proxy | O(1) | O(1) | O(1) | ✅ |
| Observer | O(1) | O(n) | O(n) | ✅ |
| Strategy | O(1) | O(1) | O(1) | ✅ |
| Command | O(1) | O(1) | O(n) | ✅ |
| Worker Pool | O(w) | O(1) | O(w+t) | ✅ |
| Pipeline | O(s) | O(s) | O(s) | ✅ |

## Anti-Patterns and Solutions

### Goroutine Leak

```go
// Anti-pattern
func processItems(items []string) {
    for _, item := range items {
        go processItem(item) // Unbounded goroutines
    }
}

// Solution: Worker Pool
func processItemsWithPool(items []string) {
    pool := NewWorkerPool(10, 100)
    pool.Start()
    
    for _, item := range items {
        task := NewConcreteTask(item, item)
        pool.Submit(task)
    }
    
    pool.Stop()
}
```

### Mutex Misuse

```go
// Anti-pattern
func (s *Service) ProcessData(data []int) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    result := expensiveOperation(data) // Holding lock too long
    s.data = result
}

// Solution: Minimize critical section
func (s *Service) ProcessData(data []int) {
    result := expensiveOperation(data) // Outside lock
    
    s.mutex.Lock()
    s.data = result
    s.mutex.Unlock()
}
```

## Best Practices

1. **Thread Safety**: Use sync.Once for singletons, mutexes for shared state
2. **Memory Management**: Implement object pooling for high-frequency objects
3. **Error Handling**: Proper error propagation and recovery
4. **Performance**: Choose patterns based on complexity requirements
5. **Testing**: Unit tests for each pattern implementation
6. **Documentation**: Clear interfaces and usage examples

## Key Principles

- **Composition over Inheritance**: Use interfaces and embedding
- **Dependency Inversion**: Depend on abstractions, not concretions
- **Single Responsibility**: Each pattern has one clear purpose
- **Open/Closed**: Open for extension, closed for modification
- **Interface Segregation**: Small, focused interfaces
- **Dependency Injection**: Inject dependencies rather than creating them

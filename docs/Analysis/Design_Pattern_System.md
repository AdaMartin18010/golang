# Design Pattern System

## Executive Summary

This document provides a comprehensive framework for design patterns in Golang, with formal definitions, mathematical models, and production-ready implementations.

## 1. Creational Patterns

### 1.1 Factory Pattern

**Definition 1.1.1 (Factory Pattern)**
The Factory pattern provides an interface for creating objects without specifying their exact classes.

**Mathematical Model:**

```
Factory = (Creator, Product, CreationMethod)
where:
- Creator: Interface defining creation method
- Product: Interface for created objects
- CreationMethod: Creator → Product
```

**Golang Implementation:**

```go
// Product interface
type Product interface {
    Operation() string
}

// Concrete products
type ConcreteProductA struct{}

func (cpa *ConcreteProductA) Operation() string {
    return "ConcreteProductA operation"
}

type ConcreteProductB struct{}

func (cpb *ConcreteProductB) Operation() string {
    return "ConcreteProductB operation"
}

// Creator interface
type Creator interface {
    CreateProduct() Product
}

// Concrete creators
type ConcreteCreatorA struct{}

func (cca *ConcreteCreatorA) CreateProduct() Product {
    return &ConcreteProductA{}
}

type ConcreteCreatorB struct{}

func (ccb *ConcreteCreatorB) CreateProduct() Product {
    return &ConcreteProductB{}
}

// Factory function
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

### 1.2 Singleton Pattern

**Definition 1.2.1 (Singleton Pattern)**
The Singleton pattern ensures a class has only one instance and provides a global point of access to it.

**Mathematical Model:**

```
Singleton = (Instance, GetInstance)
where:
- Instance: Single object instance
- GetInstance: () → Instance (always returns same instance)
```

**Golang Implementation:**

```go
// Thread-safe singleton
type Singleton struct {
    data string
}

var (
    instance *Singleton
    once     sync.Once
)

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{
            data: "Singleton instance",
        }
    })
    return instance
}

// Alternative implementation with mutex
type SingletonMutex struct {
    data string
}

var (
    instanceMutex *SingletonMutex
    mutex         sync.Mutex
)

func GetInstanceMutex() *SingletonMutex {
    if instanceMutex == nil {
        mutex.Lock()
        defer mutex.Unlock()
        
        if instanceMutex == nil {
            instanceMutex = &SingletonMutex{
                data: "Singleton with mutex",
            }
        }
    }
    return instanceMutex
}
```

### 1.3 Builder Pattern

**Definition 1.3.1 (Builder Pattern)**
The Builder pattern constructs complex objects step by step, allowing the same construction process to create different representations.

**Mathematical Model:**

```
Builder = (Builder, Director, Product)
where:
- Builder: Interface with build methods
- Director: Uses builder to construct product
- Product: Complex object being built
```

**Golang Implementation:**

```go
// Product
type Computer struct {
    CPU       string
    Memory    string
    Storage   string
    Graphics  string
}

// Builder interface
type ComputerBuilder interface {
    SetCPU(cpu string) ComputerBuilder
    SetMemory(memory string) ComputerBuilder
    SetStorage(storage string) ComputerBuilder
    SetGraphics(graphics string) ComputerBuilder
    Build() *Computer
}

// Concrete builder
type ConcreteComputerBuilder struct {
    computer *Computer
}

func NewComputerBuilder() ComputerBuilder {
    return &ConcreteComputerBuilder{
        computer: &Computer{},
    }
}

func (cb *ConcreteComputerBuilder) SetCPU(cpu string) ComputerBuilder {
    cb.computer.CPU = cpu
    return cb
}

func (cb *ConcreteComputerBuilder) SetMemory(memory string) ComputerBuilder {
    cb.computer.Memory = memory
    return cb
}

func (cb *ConcreteComputerBuilder) SetStorage(storage string) ComputerBuilder {
    cb.computer.Storage = storage
    return cb
}

func (cb *ConcreteComputerBuilder) SetGraphics(graphics string) ComputerBuilder {
    cb.computer.Graphics = graphics
    return cb
}

func (cb *ConcreteComputerBuilder) Build() *Computer {
    return cb.computer
}

// Director
type ComputerDirector struct {
    builder ComputerBuilder
}

func NewComputerDirector(builder ComputerBuilder) *ComputerDirector {
    return &ComputerDirector{builder: builder}
}

func (cd *ComputerDirector) ConstructGamingComputer() *Computer {
    return cd.builder.
        SetCPU("Intel i9").
        SetMemory("32GB DDR4").
        SetStorage("1TB NVMe SSD").
        SetGraphics("RTX 4080").
        Build()
}

func (cd *ComputerDirector) ConstructOfficeComputer() *Computer {
    return cd.builder.
        SetCPU("Intel i5").
        SetMemory("8GB DDR4").
        SetStorage("256GB SSD").
        SetGraphics("Integrated").
        Build()
}
```

## 2. Structural Patterns

### 2.1 Adapter Pattern

**Definition 2.1.1 (Adapter Pattern)**
The Adapter pattern allows incompatible interfaces to work together by wrapping an existing class with a new interface.

**Mathematical Model:**

```
Adapter = (Target, Adaptee, Adapter)
where:
- Target: Expected interface
- Adaptee: Existing incompatible interface
- Adapter: Target → Adaptee (translation)
```

**Golang Implementation:**

```go
// Target interface
type Target interface {
    Request() string
}

// Adaptee (incompatible interface)
type Adaptee struct{}

func (a *Adaptee) SpecificRequest() string {
    return "Specific request from adaptee"
}

// Adapter
type Adapter struct {
    adaptee *Adaptee
}

func NewAdapter(adaptee *Adaptee) *Adapter {
    return &Adapter{adaptee: adaptee}
}

func (a *Adapter) Request() string {
    // Translate target request to adaptee request
    return "Adapter: " + a.adaptee.SpecificRequest()
}

// Usage
func UseTarget(target Target) {
    fmt.Println(target.Request())
}
```

### 2.2 Decorator Pattern

**Definition 2.2.1 (Decorator Pattern)**
The Decorator pattern attaches additional responsibilities to an object dynamically, providing a flexible alternative to subclassing.

**Mathematical Model:**

```
Decorator = (Component, ConcreteComponent, Decorator)
where:
- Component: Interface for objects that can have responsibilities added
- ConcreteComponent: Basic implementation
- Decorator: Wraps component and adds behavior
```

**Golang Implementation:**

```go
// Component interface
type Component interface {
    Operation() string
}

// Concrete component
type ConcreteComponent struct{}

func (cc *ConcreteComponent) Operation() string {
    return "ConcreteComponent"
}

// Base decorator
type Decorator struct {
    component Component
}

func NewDecorator(component Component) *Decorator {
    return &Decorator{component: component}
}

func (d *Decorator) Operation() string {
    return d.component.Operation()
}

// Concrete decorators
type ConcreteDecoratorA struct {
    Decorator
}

func NewConcreteDecoratorA(component Component) *ConcreteDecoratorA {
    return &ConcreteDecoratorA{
        Decorator: *NewDecorator(component),
    }
}

func (cda *ConcreteDecoratorA) Operation() string {
    return "ConcreteDecoratorA(" + cda.component.Operation() + ")"
}

type ConcreteDecoratorB struct {
    Decorator
}

func NewConcreteDecoratorB(component Component) *ConcreteDecoratorB {
    return &ConcreteDecoratorB{
        Decorator: *NewDecorator(component),
    }
}

func (cdb *ConcreteDecoratorB) Operation() string {
    return "ConcreteDecoratorB(" + cdb.component.Operation() + ")"
}
```

### 2.3 Proxy Pattern

**Definition 2.3.1 (Proxy Pattern)**
The Proxy pattern provides a surrogate or placeholder for another object to control access to it.

**Mathematical Model:**

```
Proxy = (Subject, RealSubject, Proxy)
where:
- Subject: Interface for real object and proxy
- RealSubject: Actual object
- Proxy: Controls access to RealSubject
```

**Golang Implementation:**

```go
// Subject interface
type Subject interface {
    Request() string
}

// Real subject
type RealSubject struct{}

func (rs *RealSubject) Request() string {
    return "RealSubject request"
}

// Proxy
type Proxy struct {
    realSubject *RealSubject
    accessCount int
    mutex       sync.Mutex
}

func NewProxy() *Proxy {
    return &Proxy{}
}

func (p *Proxy) Request() string {
    p.mutex.Lock()
    defer p.mutex.Unlock()
    
    // Lazy initialization
    if p.realSubject == nil {
        p.realSubject = &RealSubject{}
    }
    
    p.accessCount++
    
    // Add proxy behavior
    return fmt.Sprintf("Proxy (access #%d): %s", p.accessCount, p.realSubject.Request())
}

// Virtual proxy for expensive operations
type VirtualProxy struct {
    realSubject *RealSubject
    mutex       sync.Mutex
}

func NewVirtualProxy() *VirtualProxy {
    return &VirtualProxy{}
}

func (vp *VirtualProxy) Request() string {
    vp.mutex.Lock()
    defer vp.mutex.Unlock()
    
    if vp.realSubject == nil {
        // Simulate expensive initialization
        time.Sleep(100 * time.Millisecond)
        vp.realSubject = &RealSubject{}
    }
    
    return vp.realSubject.Request()
}
```

## 3. Behavioral Patterns

### 3.1 Observer Pattern

**Definition 3.1.1 (Observer Pattern)**
The Observer pattern defines a one-to-many dependency between objects so that when one object changes state, all its dependents are notified and updated automatically.

**Mathematical Model:**

```
Observer = (Subject, Observer, Notify)
where:
- Subject: Object being observed
- Observer: Objects to be notified
- Notify: Subject → Observer[] (notification mechanism)
```

**Golang Implementation:**

```go
// Observer interface
type Observer interface {
    Update(data interface{})
}

// Subject interface
type Subject interface {
    Attach(observer Observer)
    Detach(observer Observer)
    Notify()
}

// Concrete subject
type ConcreteSubject struct {
    observers []Observer
    data      interface{}
    mutex     sync.RWMutex
}

func NewConcreteSubject() *ConcreteSubject {
    return &ConcreteSubject{
        observers: make([]Observer, 0),
    }
}

func (cs *ConcreteSubject) Attach(observer Observer) {
    cs.mutex.Lock()
    defer cs.mutex.Unlock()
    
    cs.observers = append(cs.observers, observer)
}

func (cs *ConcreteSubject) Detach(observer Observer) {
    cs.mutex.Lock()
    defer cs.mutex.Unlock()
    
    for i, obs := range cs.observers {
        if obs == observer {
            cs.observers = append(cs.observers[:i], cs.observers[i+1:]...)
            break
        }
    }
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

func (cs *ConcreteSubject) SetData(data interface{}) {
    cs.mutex.Lock()
    cs.data = data
    cs.mutex.Unlock()
    
    cs.Notify()
}

// Concrete observers
type ConcreteObserverA struct {
    id string
}

func (coa *ConcreteObserverA) Update(data interface{}) {
    fmt.Printf("ObserverA %s received: %v\n", coa.id, data)
}

type ConcreteObserverB struct {
    id string
}

func (cob *ConcreteObserverB) Update(data interface{}) {
    fmt.Printf("ObserverB %s received: %v\n", cob.id, data)
}
```

### 3.2 Strategy Pattern

**Definition 3.2.1 (Strategy Pattern)**
The Strategy pattern defines a family of algorithms, encapsulates each one, and makes them interchangeable.

**Mathematical Model:**

```
Strategy = (Context, Strategy, Algorithm)
where:
- Context: Uses strategy
- Strategy: Interface for algorithms
- Algorithm: Concrete strategy implementations
```

**Golang Implementation:**

```go
// Strategy interface
type Strategy interface {
    Algorithm(data []int) []int
}

// Concrete strategies
type BubbleSortStrategy struct{}

func (bss *BubbleSortStrategy) Algorithm(data []int) []int {
    result := make([]int, len(data))
    copy(result, data)
    
    n := len(result)
    for i := 0; i < n-1; i++ {
        for j := 0; j < n-i-1; j++ {
            if result[j] > result[j+1] {
                result[j], result[j+1] = result[j+1], result[j]
            }
        }
    }
    
    return result
}

type QuickSortStrategy struct{}

func (qss *QuickSortStrategy) Algorithm(data []int) []int {
    result := make([]int, len(data))
    copy(result, data)
    
    if len(result) <= 1 {
        return result
    }
    
    pivot := result[0]
    var left, right []int
    
    for i := 1; i < len(result); i++ {
        if result[i] < pivot {
            left = append(left, result[i])
        } else {
            right = append(right, result[i])
        }
    }
    
    left = qss.Algorithm(left)
    right = qss.Algorithm(right)
    
    return append(append(left, pivot), right...)
}

// Context
type Context struct {
    strategy Strategy
}

func NewContext(strategy Strategy) *Context {
    return &Context{strategy: strategy}
}

func (c *Context) SetStrategy(strategy Strategy) {
    c.strategy = strategy
}

func (c *Context) ExecuteStrategy(data []int) []int {
    return c.strategy.Algorithm(data)
}
```

### 3.3 Command Pattern

**Definition 3.3.1 (Command Pattern)**
The Command pattern encapsulates a request as an object, allowing parameterization of clients with different requests, queuing of requests, and logging of operations.

**Mathematical Model:**

```
Command = (Command, Receiver, Invoker)
where:
- Command: Interface for executing operations
- Receiver: Object that performs the operation
- Invoker: Sends command to receiver
```

**Golang Implementation:**

```go
// Command interface
type Command interface {
    Execute()
    Undo()
}

// Receiver
type Receiver struct {
    data string
}

func (r *Receiver) Action1() {
    fmt.Printf("Receiver performing Action1 on: %s\n", r.data)
}

func (r *Receiver) Action2() {
    fmt.Printf("Receiver performing Action2 on: %s\n", r.data)
}

// Concrete commands
type ConcreteCommand1 struct {
    receiver *Receiver
}

func NewConcreteCommand1(receiver *Receiver) *ConcreteCommand1 {
    return &ConcreteCommand1{receiver: receiver}
}

func (cc1 *ConcreteCommand1) Execute() {
    cc1.receiver.Action1()
}

func (cc1 *ConcreteCommand1) Undo() {
    fmt.Println("Undoing Action1")
}

type ConcreteCommand2 struct {
    receiver *Receiver
}

func NewConcreteCommand2(receiver *Receiver) *ConcreteCommand2 {
    return &ConcreteCommand2{receiver: receiver}
}

func (cc2 *ConcreteCommand2) Execute() {
    cc2.receiver.Action2()
}

func (cc2 *ConcreteCommand2) Undo() {
    fmt.Println("Undoing Action2")
}

// Invoker
type Invoker struct {
    commands []Command
    history  []Command
    mutex    sync.Mutex
}

func NewInvoker() *Invoker {
    return &Invoker{
        commands: make([]Command, 0),
        history:  make([]Command, 0),
    }
}

func (i *Invoker) AddCommand(command Command) {
    i.mutex.Lock()
    defer i.mutex.Unlock()
    
    i.commands = append(i.commands, command)
}

func (i *Invoker) ExecuteCommands() {
    i.mutex.Lock()
    defer i.mutex.Unlock()
    
    for _, command := range i.commands {
        command.Execute()
        i.history = append(i.history, command)
    }
    i.commands = i.commands[:0] // Clear commands
}

func (i *Invoker) UndoLast() {
    i.mutex.Lock()
    defer i.mutex.Unlock()
    
    if len(i.history) > 0 {
        lastCommand := i.history[len(i.history)-1]
        lastCommand.Undo()
        i.history = i.history[:len(i.history)-1]
    }
}
```

## 4. Concurrency Patterns

### 4.1 Worker Pool Pattern

**Definition 4.1.1 (Worker Pool Pattern)**
The Worker Pool pattern maintains a pool of worker goroutines to process tasks from a queue, providing controlled concurrency and resource management.

**Mathematical Model:**

```
WorkerPool = (Workers, Tasks, Queue)
where:
- Workers: Pool of goroutines
- Tasks: Work to be performed
- Queue: Task distribution mechanism
```

**Golang Implementation:**

```go
// Task interface
type Task interface {
    Execute() error
    ID() string
}

// Concrete task
type ConcreteTask struct {
    id   string
    data interface{}
}

func NewConcreteTask(id string, data interface{}) *ConcreteTask {
    return &ConcreteTask{
        id:   id,
        data: data,
    }
}

func (ct *ConcreteTask) Execute() error {
    // Simulate work
    time.Sleep(100 * time.Millisecond)
    fmt.Printf("Task %s executed with data: %v\n", ct.id, ct.data)
    return nil
}

func (ct *ConcreteTask) ID() string {
    return ct.id
}

// Worker pool
type WorkerPool struct {
    workers    int
    taskQueue  chan Task
    resultChan chan error
    wg         sync.WaitGroup
    ctx        context.Context
    cancel     context.CancelFunc
}

func NewWorkerPool(workers int, queueSize int) *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &WorkerPool{
        workers:    workers,
        taskQueue:  make(chan Task, queueSize),
        resultChan: make(chan error, queueSize),
        ctx:        ctx,
        cancel:     cancel,
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.worker(i)
    }
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

func (wp *WorkerPool) Submit(task Task) error {
    select {
    case wp.taskQueue <- task:
        return nil
    case <-wp.ctx.Done():
        return fmt.Errorf("worker pool is stopped")
    default:
        return fmt.Errorf("task queue is full")
    }
}

func (wp *WorkerPool) Stop() {
    wp.cancel()
    close(wp.taskQueue)
    wp.wg.Wait()
    close(wp.resultChan)
}

func (wp *WorkerPool) Results() <-chan error {
    return wp.resultChan
}
```

### 4.2 Pipeline Pattern

**Definition 4.2.1 (Pipeline Pattern)**
The Pipeline pattern processes data through a series of stages, where each stage performs a specific transformation on the data.

**Mathematical Model:**

```
Pipeline = (Stage, Data, Transformation)
where:
- Stage: Processing step
- Data: Information flowing through pipeline
- Transformation: Stage → Data → Data
```

**Golang Implementation:**

```go
// Pipeline stage interface
type Stage interface {
    Process(data interface{}) (interface{}, error)
}

// Concrete stages
type FilterStage struct {
    predicate func(interface{}) bool
}

func NewFilterStage(predicate func(interface{}) bool) *FilterStage {
    return &FilterStage{predicate: predicate}
}

func (fs *FilterStage) Process(data interface{}) (interface{}, error) {
    if fs.predicate(data) {
        return data, nil
    }
    return nil, nil // Filter out
}

type TransformStage struct {
    transformer func(interface{}) interface{}
}

func NewTransformStage(transformer func(interface{}) interface{}) *TransformStage {
    return &TransformStage{transformer: transformer}
}

func (ts *TransformStage) Process(data interface{}) (interface{}, error) {
    return ts.transformer(data), nil
}

// Pipeline
type Pipeline struct {
    stages []Stage
}

func NewPipeline() *Pipeline {
    return &Pipeline{
        stages: make([]Stage, 0),
    }
}

func (p *Pipeline) AddStage(stage Stage) {
    p.stages = append(p.stages, stage)
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
            return nil, nil // Filtered out
        }
    }
    
    return result, nil
}

// Concurrent pipeline
type ConcurrentPipeline struct {
    stages []Stage
}

func NewConcurrentPipeline() *ConcurrentPipeline {
    return &ConcurrentPipeline{
        stages: make([]Stage, 0),
    }
}

func (cp *ConcurrentPipeline) AddStage(stage Stage) {
    cp.stages = append(cp.stages, stage)
}

func (cp *ConcurrentPipeline) Execute(data []interface{}) ([]interface{}, error) {
    if len(cp.stages) == 0 {
        return data, nil
    }
    
    // Create channels for each stage
    channels := make([]chan interface{}, len(cp.stages)+1)
    for i := range channels {
        channels[i] = make(chan interface{}, len(data))
    }
    
    // Start workers for each stage
    var wg sync.WaitGroup
    
    for i, stage := range cp.stages {
        wg.Add(1)
        go func(stageIndex int, stage Stage) {
            defer wg.Done()
            defer close(channels[stageIndex+1])
            
            for item := range channels[stageIndex] {
                result, err := stage.Process(item)
                if err != nil {
                    fmt.Printf("Stage %d error: %v\n", stageIndex, err)
                    continue
                }
                
                if result != nil {
                    channels[stageIndex+1] <- result
                }
            }
        }(i, stage)
    }
    
    // Send input data
    go func() {
        defer close(channels[0])
        for _, item := range data {
            channels[0] <- item
        }
    }()
    
    // Collect results
    var results []interface{}
    go func() {
        for item := range channels[len(channels)-1] {
            results = append(results, item)
        }
    }()
    
    wg.Wait()
    return results, nil
}
```

## 5. Enterprise Patterns

### 5.1 Repository Pattern

**Definition 5.1.1 (Repository Pattern)**
The Repository pattern abstracts the data persistence layer, providing a collection-like interface for accessing domain objects.

**Mathematical Model:**

```
Repository = (Entity, Repository, DataSource)
where:
- Entity: Domain object
- Repository: Interface for data access
- DataSource: Actual data storage
```

**Golang Implementation:**

```go
// Entity
type User struct {
    ID       string
    Name     string
    Email    string
    CreatedAt time.Time
}

// Repository interface
type Repository[T any] interface {
    FindByID(id string) (T, error)
    FindAll() ([]T, error)
    Save(entity T) error
    Delete(id string) error
    Update(entity T) error
}

// In-memory repository
type InMemoryRepository[T any] struct {
    data map[string]T
    mutex sync.RWMutex
}

func NewInMemoryRepository[T any]() *InMemoryRepository[T] {
    return &InMemoryRepository[T]{
        data: make(map[string]T),
    }
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

func (imr *InMemoryRepository[T]) FindAll() ([]T, error) {
    imr.mutex.RLock()
    defer imr.mutex.RUnlock()
    
    entities := make([]T, 0, len(imr.data))
    for _, entity := range imr.data {
        entities = append(entities, entity)
    }
    
    return entities, nil
}

func (imr *InMemoryRepository[T]) Save(entity T) error {
    imr.mutex.Lock()
    defer imr.mutex.Unlock()
    
    // Assume entity has an ID field (simplified)
    // In practice, you'd use reflection or interfaces to get ID
    imr.data["generated-id"] = entity
    return nil
}

func (imr *InMemoryRepository[T]) Delete(id string) error {
    imr.mutex.Lock()
    defer imr.mutex.Unlock()
    
    if _, exists := imr.data[id]; !exists {
        return fmt.Errorf("entity with id %s not found", id)
    }
    
    delete(imr.data, id)
    return nil
}

func (imr *InMemoryRepository[T]) Update(entity T) error {
    imr.mutex.Lock()
    defer imr.mutex.Unlock()
    
    // Assume entity has an ID field (simplified)
    imr.data["generated-id"] = entity
    return nil
}

// User-specific repository
type UserRepository interface {
    Repository[User]
    FindByEmail(email string) (User, error)
}

type InMemoryUserRepository struct {
    *InMemoryRepository[User]
    emailIndex map[string]string // email -> userID
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
    return &InMemoryUserRepository{
        InMemoryRepository: NewInMemoryRepository[User](),
        emailIndex:         make(map[string]string),
    }
}

func (imur *InMemoryUserRepository) FindByEmail(email string) (User, error) {
    imur.mutex.RLock()
    defer imur.mutex.RUnlock()
    
    userID, exists := imur.emailIndex[email]
    if !exists {
        var zero User
        return zero, fmt.Errorf("user with email %s not found", email)
    }
    
    return imur.data[userID], nil
}
```

### 5.2 Unit of Work Pattern

**Definition 5.2.1 (Unit of Work Pattern)**
The Unit of Work pattern maintains a list of objects affected by a business transaction and coordinates the writing out of changes and the resolution of concurrency problems.

**Mathematical Model:**

```
UnitOfWork = (Entities, Changes, Transaction)
where:
- Entities: Objects being tracked
- Changes: Modifications to entities
- Transaction: Atomic operation boundary
```

**Golang Implementation:**

```go
// Unit of Work interface
type UnitOfWork interface {
    RegisterNew(entity interface{})
    RegisterDirty(entity interface{})
    RegisterDeleted(entity interface{})
    Commit() error
    Rollback() error
}

// Entity tracking
type EntityState int

const (
    EntityStateNew EntityState = iota
    EntityStateDirty
    EntityStateDeleted
    EntityStateClean
)

type EntityTracker struct {
    entity interface{}
    state  EntityState
}

// Concrete Unit of Work
type ConcreteUnitOfWork struct {
    newEntities     map[string]*EntityTracker
    dirtyEntities   map[string]*EntityTracker
    deletedEntities map[string]*EntityTracker
    repositories    map[string]Repository[interface{}]
    mutex           sync.Mutex
}

func NewConcreteUnitOfWork() *ConcreteUnitOfWork {
    return &ConcreteUnitOfWork{
        newEntities:     make(map[string]*EntityTracker),
        dirtyEntities:   make(map[string]*EntityTracker),
        deletedEntities: make(map[string]*EntityTracker),
        repositories:    make(map[string]Repository[interface{}]),
    }
}

func (uow *ConcreteUnitOfWork) RegisterNew(entity interface{}) {
    uow.mutex.Lock()
    defer uow.mutex.Unlock()
    
    entityID := uow.getEntityID(entity)
    uow.newEntities[entityID] = &EntityTracker{
        entity: entity,
        state:  EntityStateNew,
    }
}

func (uow *ConcreteUnitOfWork) RegisterDirty(entity interface{}) {
    uow.mutex.Lock()
    defer uow.mutex.Unlock()
    
    entityID := uow.getEntityID(entity)
    
    // Remove from new if already there
    delete(uow.newEntities, entityID)
    
    uow.dirtyEntities[entityID] = &EntityTracker{
        entity: entity,
        state:  EntityStateDirty,
    }
}

func (uow *ConcreteUnitOfWork) RegisterDeleted(entity interface{}) {
    uow.mutex.Lock()
    defer uow.mutex.Unlock()
    
    entityID := uow.getEntityID(entity)
    
    // Remove from new and dirty
    delete(uow.newEntities, entityID)
    delete(uow.dirtyEntities, entityID)
    
    uow.deletedEntities[entityID] = &EntityTracker{
        entity: entity,
        state:  EntityStateDeleted,
    }
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
    
    // Process dirty entities
    for entityID, tracker := range uow.dirtyEntities {
        repo := uow.getRepository(tracker.entity)
        if err := repo.Update(tracker.entity); err != nil {
            return fmt.Errorf("failed to update entity %s: %w", entityID, err)
        }
    }
    
    // Process deleted entities
    for entityID, tracker := range uow.deletedEntities {
        repo := uow.getRepository(tracker.entity)
        if err := repo.Delete(entityID); err != nil {
            return fmt.Errorf("failed to delete entity %s: %w", entityID, err)
        }
    }
    
    // Clear all tracking
    uow.clearTracking()
    
    return nil
}

func (uow *ConcreteUnitOfWork) Rollback() error {
    uow.mutex.Lock()
    defer uow.mutex.Unlock()
    
    uow.clearTracking()
    return nil
}

func (uow *ConcreteUnitOfWork) getEntityID(entity interface{}) string {
    // Simplified - in practice, use reflection or interfaces
    return fmt.Sprintf("%p", entity)
}

func (uow *ConcreteUnitOfWork) getRepository(entity interface{}) Repository[interface{}] {
    // Simplified - in practice, use type assertions or dependency injection
    return nil
}

func (uow *ConcreteUnitOfWork) clearTracking() {
    uow.newEntities = make(map[string]*EntityTracker)
    uow.dirtyEntities = make(map[string]*EntityTracker)
    uow.deletedEntities = make(map[string]*EntityTracker)
}
```

## 6. Anti-Patterns and Refactoring

### 6.1 Common Anti-Patterns

**Definition 6.1.1 (Anti-Pattern)**
An anti-pattern is a common response to a recurring problem that is usually ineffective and risks being highly counterproductive.

**Common Anti-Patterns in Go:**

1. **Goroutine Leak**

```go
// Anti-pattern: Unbounded goroutines
func processItems(items []string) {
    for _, item := range items {
        go processItem(item) // Creates unbounded goroutines
    }
}

// Solution: Worker pool
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

2. **Mutex Misuse**

```go
// Anti-pattern: Holding locks too long
func (s *Service) ProcessData(data []int) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    // Expensive operation while holding lock
    result := expensiveOperation(data)
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

3. **Channel Misuse**

```go
// Anti-pattern: Unbuffered channels for high-throughput
func processHighThroughput(data []int) {
    ch := make(chan int) // Unbuffered
    
    go func() {
        for _, item := range data {
            ch <- item // Blocks on each send
        }
        close(ch)
    }()
    
    for item := range ch {
        process(item)
    }
}

// Solution: Buffered channels
func processHighThroughput(data []int) {
    ch := make(chan int, len(data)) // Buffered
    
    go func() {
        for _, item := range data {
            ch <- item // Non-blocking
        }
        close(ch)
    }()
    
    for item := range ch {
        process(item)
    }
}
```

### 6.2 Refactoring Strategies

**Definition 6.2.1 (Refactoring)**
Refactoring is the process of restructuring existing code without changing its external behavior.

**Refactoring Techniques:**

1. **Extract Method**

```go
// Before
func processUser(user *User) {
    // Validate user
    if user.Name == "" {
        return errors.New("name is required")
    }
    if user.Email == "" {
        return errors.New("email is required")
    }
    if !strings.Contains(user.Email, "@") {
        return errors.New("invalid email format")
    }
    
    // Process user
    // ... complex logic
}

// After
func processUser(user *User) error {
    if err := validateUser(user); err != nil {
        return err
    }
    
    return processUserLogic(user)
}

func validateUser(user *User) error {
    if user.Name == "" {
        return errors.New("name is required")
    }
    if user.Email == "" {
        return errors.New("email is required")
    }
    if !strings.Contains(user.Email, "@") {
        return errors.New("invalid email format")
    }
    return nil
}
```

2. **Replace Conditional with Strategy**

```go
// Before
func calculateDiscount(order *Order, customerType string) float64 {
    switch customerType {
    case "VIP":
        return order.Total * 0.2
    case "Regular":
        return order.Total * 0.1
    case "New":
        return order.Total * 0.05
    default:
        return 0
    }
}

// After
type DiscountStrategy interface {
    CalculateDiscount(order *Order) float64
}

type VIPDiscountStrategy struct{}

func (vds *VIPDiscountStrategy) CalculateDiscount(order *Order) float64 {
    return order.Total * 0.2
}

type RegularDiscountStrategy struct{}

func (rds *RegularDiscountStrategy) CalculateDiscount(order *Order) float64 {
    return order.Total * 0.1
}

func calculateDiscount(order *Order, strategy DiscountStrategy) float64 {
    return strategy.CalculateDiscount(order)
}
```

## 7. Performance Characteristics

### 7.1 Pattern Performance Analysis

**Time Complexity Analysis:**

| Pattern | Creation | Operation | Memory |
|---------|----------|-----------|---------|
| Factory | O(1) | O(1) | O(1) |
| Singleton | O(1) | O(1) | O(1) |
| Builder | O(n) | O(1) | O(n) |
| Adapter | O(1) | O(1) | O(1) |
| Decorator | O(1) | O(n) | O(n) |
| Proxy | O(1) | O(1) | O(1) |
| Observer | O(1) | O(n) | O(n) |
| Strategy | O(1) | O(1) | O(1) |
| Command | O(1) | O(1) | O(n) |
| Worker Pool | O(w) | O(1) | O(w+t) |
| Pipeline | O(s) | O(s) | O(s) |

Where:

- n = number of elements
- w = number of workers
- t = number of tasks
- s = number of stages

### 7.2 Memory Management

**Memory Optimization Strategies:**

1. **Object Pooling**

```go
type ObjectPool[T any] struct {
    pool chan T
    new  func() T
}

func NewObjectPool[T any](size int, newFunc func() T) *ObjectPool[T] {
    pool := &ObjectPool[T]{
        pool: make(chan T, size),
        new:  newFunc,
    }
    
    // Pre-populate pool
    for i := 0; i < size; i++ {
        pool.pool <- newFunc()
    }
    
    return pool
}

func (op *ObjectPool[T]) Get() T {
    select {
    case obj := <-op.pool:
        return obj
    default:
        return op.new()
    }
}

func (op *ObjectPool[T]) Put(obj T) {
    select {
    case op.pool <- obj:
    default:
        // Pool is full, discard object
    }
}
```

2. **Lazy Initialization**

```go
type LazySingleton struct {
    data string
}

var (
    lazyInstance *LazySingleton
    lazyOnce     sync.Once
)

func GetLazyInstance() *LazySingleton {
    lazyOnce.Do(func() {
        lazyInstance = &LazySingleton{
            data: "Lazy initialized",
        }
    })
    return lazyInstance
}
```

## 8. Conclusion

This design pattern system provides comprehensive coverage of:

1. **Creational Patterns**: Factory, Singleton, Builder with thread-safe implementations
2. **Structural Patterns**: Adapter, Decorator, Proxy with flexible composition
3. **Behavioral Patterns**: Observer, Strategy, Command with proper concurrency handling
4. **Concurrency Patterns**: Worker Pool, Pipeline with performance optimization
5. **Enterprise Patterns**: Repository, Unit of Work with transaction management
6. **Anti-Patterns**: Common pitfalls and refactoring strategies
7. **Performance Analysis**: Time complexity and memory management

The framework emphasizes:

- **Thread Safety**: Proper synchronization mechanisms
- **Performance**: Optimized implementations with complexity analysis
- **Flexibility**: Generic implementations for reusability
- **Best Practices**: Go-specific patterns and idioms

## References

1. Gamma, E., Helm, R., Johnson, R., & Vlissides, J. (1994). Design Patterns: Elements of Reusable Object-Oriented Software. Addison-Wesley.
2. Freeman, E., Robson, E., Sierra, K., & Bates, B. (2004). Head First Design Patterns. O'Reilly.
3. Go Documentation: <https://golang.org/doc/>
4. Go Concurrency Patterns: <https://golang.org/doc/effective_go.html#concurrency>

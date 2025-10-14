# 3.3.1 备忘录模式 (Memento Pattern)

## 3.3.1.1 目录

## 3.3.1.2 1. 概述

### 3.3.1.2.1 定义

备忘录模式在不破坏封装的前提下，捕获并外部化一个对象的内部状态，这样以后就可以将该对象恢复到原先保存的状态。

**形式化定义**:
$$Memento = (Originator, Memento, Caretaker, State)$$

其中：

- $Originator$ 是发起人
- $Memento$ 是备忘录
- $Caretaker$ 是管理者
- $State$ 是状态

### 3.3.1.2.2 核心特征

- **状态保存**: 保存对象的内部状态
- **状态恢复**: 将对象恢复到之前的状态
- **封装保护**: 不破坏对象的封装性
- **历史管理**: 管理多个历史状态

## 3.3.1.3 2. 理论基础

### 3.3.1.3.1 数学形式化

**定义 2.1** (备忘录模式): 备忘录模式是一个四元组 $M = (O, Mem, Car, S)$

其中：

- $O$ 是发起人集合
- $Mem$ 是备忘录集合
- $Car$ 是管理者集合
- $S$ 是状态保存函数，$S: O \rightarrow Mem$

**定理 2.1** (状态一致性): 对于任意发起人 $o \in O$ 和备忘录 $m \in Mem$，状态恢复后对象状态与备忘录状态一致。

### 3.3.1.3.2 范畴论视角

在范畴论中，备忘录模式可以表示为：

$$Memento : Originator \rightarrow State$$

## 3.3.1.4 3. Go语言实现

### 3.3.1.4.1 基础备忘录模式

```go
package memento

import (
    "fmt"
    "time"
)

// Memento 备忘录接口
type Memento interface {
    GetTimestamp() time.Time
    GetDescription() string
}

// Originator 发起人接口
type Originator interface {
    CreateMemento() Memento
    RestoreFromMemento(memento Memento)
    GetState() string
}

// ConcreteMemento 具体备忘录
type ConcreteMemento struct {
    state       string
    timestamp   time.Time
    description string
}

func NewConcreteMemento(state, description string) *ConcreteMemento {
    return &ConcreteMemento{
        state:       state,
        timestamp:   time.Now(),
        description: description,
    }
}

func (c *ConcreteMemento) GetState() string {
    return c.state
}

func (c *ConcreteMemento) GetTimestamp() time.Time {
    return c.timestamp
}

func (c *ConcreteMemento) GetDescription() string {
    return c.description
}

// ConcreteOriginator 具体发起人
type ConcreteOriginator struct {
    state string
}

func NewConcreteOriginator(initialState string) *ConcreteOriginator {
    return &ConcreteOriginator{
        state: initialState,
    }
}

func (c *ConcreteOriginator) SetState(state string) {
    c.state = state
    fmt.Printf("Originator state changed to: %s\n", state)
}

func (c *ConcreteOriginator) CreateMemento() Memento {
    fmt.Printf("Creating memento for state: %s\n", c.state)
    return NewConcreteMemento(c.state, fmt.Sprintf("State: %s", c.state))
}

func (c *ConcreteOriginator) RestoreFromMemento(memento Memento) {
    if concreteMemento, ok := memento.(*ConcreteMemento); ok {
        c.state = concreteMemento.GetState()
        fmt.Printf("Originator restored to state: %s\n", c.state)
    }
}

func (c *ConcreteOriginator) GetState() string {
    return c.state
}

// Caretaker 管理者
type Caretaker struct {
    mementos []Memento
}

func NewCaretaker() *Caretaker {
    return &Caretaker{
        mementos: make([]Memento, 0),
    }
}

func (c *Caretaker) AddMemento(memento Memento) {
    c.mementos = append(c.mementos, memento)
    fmt.Printf("Caretaker: Added memento with description: %s\n", 
        memento.GetDescription())
}

func (c *Caretaker) GetMemento(index int) Memento {
    if index >= 0 && index < len(c.mementos) {
        return c.mementos[index]
    }
    return nil
}

func (c *Caretaker) GetMementoCount() int {
    return len(c.mementos)
}

func (c *Caretaker) ListMementos() {
    fmt.Println("Caretaker: List of mementos:")
    for i, memento := range c.mementos {
        fmt.Printf("  %d: %s (at %s)\n", i, 
            memento.GetDescription(), 
            memento.GetTimestamp().Format("15:04:05"))
    }
}

```

### 3.3.1.4.2 文本编辑器备忘录模式

```go
package texteditor

import (
    "fmt"
    "time"
)

// TextMemento 文本备忘录
type TextMemento struct {
    content     string
    cursor      int
    timestamp   time.Time
    description string
}

func NewTextMemento(content string, cursor int, description string) *TextMemento {
    return &TextMemento{
        content:     content,
        cursor:      cursor,
        timestamp:   time.Now(),
        description: description,
    }
}

func (t *TextMemento) GetContent() string {
    return t.content
}

func (t *TextMemento) GetCursor() int {
    return t.cursor
}

func (t *TextMemento) GetTimestamp() time.Time {
    return t.timestamp
}

func (t *TextMemento) GetDescription() string {
    return t.description
}

// TextEditor 文本编辑器
type TextEditor struct {
    content string
    cursor  int
}

func NewTextEditor() *TextEditor {
    return &TextEditor{
        content: "",
        cursor:  0,
    }
}

func (t *TextEditor) Insert(text string) {
    t.content = t.content[:t.cursor] + text + t.content[t.cursor:]
    t.cursor += len(text)
    fmt.Printf("Inserted '%s' at position %d\n", text, t.cursor-len(text))
}

func (t *TextEditor) Delete(length int) {
    if t.cursor+length <= len(t.content) {
        deleted := t.content[t.cursor : t.cursor+length]
        t.content = t.content[:t.cursor] + t.content[t.cursor+length:]
        fmt.Printf("Deleted '%s' at position %d\n", deleted, t.cursor)
    }
}

func (t *TextEditor) MoveCursor(position int) {
    if position >= 0 && position <= len(t.content) {
        t.cursor = position
        fmt.Printf("Cursor moved to position %d\n", t.cursor)
    }
}

func (t *TextEditor) CreateMemento() *TextMemento {
    description := fmt.Sprintf("Content length: %d, Cursor: %d", 
        len(t.content), t.cursor)
    return NewTextMemento(t.content, t.cursor, description)
}

func (t *TextEditor) RestoreFromMemento(memento *TextMemento) {
    t.content = memento.GetContent()
    t.cursor = memento.GetCursor()
    fmt.Printf("Editor restored to: %s\n", memento.GetDescription())
}

func (t *TextEditor) GetContent() string {
    return t.content
}

func (t *TextEditor) GetCursor() int {
    return t.cursor
}

func (t *TextEditor) Display() {
    fmt.Printf("Content: '%s'\n", t.content)
    fmt.Printf("Cursor: %d\n", t.cursor)
}

// UndoManager 撤销管理器
type UndoManager struct {
    mementos []*TextMemento
    current  int
}

func NewUndoManager() *UndoManager {
    return &UndoManager{
        mementos: make([]*TextMemento, 0),
        current:  -1,
    }
}

func (u *UndoManager) SaveMemento(memento *TextMemento) {
    // 移除当前位置之后的所有备忘录
    u.mementos = u.mementos[:u.current+1]
    
    u.mementos = append(u.mementos, memento)
    u.current++
    
    fmt.Printf("UndoManager: Saved memento %d\n", u.current)
}

func (u *UndoManager) Undo() *TextMemento {
    if u.current > 0 {
        u.current--
        memento := u.mementos[u.current]
        fmt.Printf("UndoManager: Undoing to memento %d\n", u.current)
        return memento
    }
    fmt.Println("UndoManager: Nothing to undo")
    return nil
}

func (u *UndoManager) Redo() *TextMemento {
    if u.current < len(u.mementos)-1 {
        u.current++
        memento := u.mementos[u.current]
        fmt.Printf("UndoManager: Redoing to memento %d\n", u.current)
        return memento
    }
    fmt.Println("UndoManager: Nothing to redo")
    return nil
}

func (u *UndoManager) CanUndo() bool {
    return u.current > 0
}

func (u *UndoManager) CanRedo() bool {
    return u.current < len(u.mementos)-1
}

func (u *UndoManager) GetHistory() []*TextMemento {
    return u.mementos[:u.current+1]
}

```

### 3.3.1.4.3 游戏状态备忘录模式

```go
package gamestate

import (
    "fmt"
    "time"
)

// GameState 游戏状态
type GameState struct {
    playerHealth int
    playerLevel  int
    playerScore  int
    inventory    map[string]int
    position     Position
    timestamp    time.Time
}

type Position struct {
    X, Y int
}

func NewGameState() *GameState {
    return &GameState{
        playerHealth: 100,
        playerLevel:  1,
        playerScore:  0,
        inventory:    make(map[string]int),
        position:     Position{0, 0},
        timestamp:    time.Now(),
    }
}

func (g *GameState) String() string {
    return fmt.Sprintf("Health: %d, Level: %d, Score: %d, Position: (%d, %d)", 
        g.playerHealth, g.playerLevel, g.playerScore, g.position.X, g.position.Y)
}

// GameMemento 游戏备忘录
type GameMemento struct {
    state       *GameState
    description string
}

func NewGameMemento(state *GameState, description string) *GameMemento {
    return &GameMemento{
        state:       state,
        description: description,
    }
}

func (g *GameMemento) GetState() *GameState {
    return g.state
}

func (g *GameMemento) GetDescription() string {
    return g.description
}

func (g *GameMemento) GetTimestamp() time.Time {
    return g.state.timestamp
}

// Game 游戏
type Game struct {
    currentState *GameState
    saveSlots    map[string]*GameMemento
}

func NewGame() *Game {
    return &Game{
        currentState: NewGameState(),
        saveSlots:    make(map[string]*GameMemento),
    }
}

func (g *Game) TakeDamage(damage int) {
    g.currentState.playerHealth -= damage
    if g.currentState.playerHealth < 0 {
        g.currentState.playerHealth = 0
    }
    fmt.Printf("Player took %d damage. Health: %d\n", damage, g.currentState.playerHealth)
}

func (g *Game) GainExperience(exp int) {
    g.currentState.playerScore += exp
    newLevel := g.currentState.playerScore/100 + 1
    if newLevel > g.currentState.playerLevel {
        g.currentState.playerLevel = newLevel
        fmt.Printf("Player leveled up to level %d!\n", g.currentState.playerLevel)
    }
    fmt.Printf("Gained %d experience. Score: %d\n", exp, g.currentState.playerScore)
}

func (g *Game) Move(dx, dy int) {
    g.currentState.position.X += dx
    g.currentState.position.Y += dy
    fmt.Printf("Player moved to position (%d, %d)\n", 
        g.currentState.position.X, g.currentState.position.Y)
}

func (g *Game) AddItem(item string, quantity int) {
    g.currentState.inventory[item] += quantity
    fmt.Printf("Added %d %s to inventory\n", quantity, item)
}

func (g *Game) SaveGame(slotName string) {
    // 创建当前状态的副本
    savedState := &GameState{
        playerHealth: g.currentState.playerHealth,
        playerLevel:  g.currentState.playerLevel,
        playerScore:  g.currentState.playerScore,
        inventory:    make(map[string]int),
        position:     g.currentState.position,
        timestamp:    time.Now(),
    }
    
    // 复制库存
    for item, quantity := range g.currentState.inventory {
        savedState.inventory[item] = quantity
    }
    
    description := fmt.Sprintf("Level %d, Score %d, Health %d", 
        savedState.playerLevel, savedState.playerScore, savedState.playerHealth)
    
    memento := NewGameMemento(savedState, description)
    g.saveSlots[slotName] = memento
    
    fmt.Printf("Game saved to slot '%s': %s\n", slotName, description)
}

func (g *Game) LoadGame(slotName string) bool {
    if memento, exists := g.saveSlots[slotName]; exists {
        g.currentState = memento.GetState()
        fmt.Printf("Game loaded from slot '%s': %s\n", slotName, memento.GetDescription())
        return true
    }
    fmt.Printf("Save slot '%s' not found\n", slotName)
    return false
}

func (g *Game) GetCurrentState() *GameState {
    return g.currentState
}

func (g *Game) ListSaveSlots() {
    fmt.Println("Available save slots:")
    for slotName, memento := range g.saveSlots {
        fmt.Printf("  %s: %s (saved at %s)\n", 
            slotName, memento.GetDescription(), 
            memento.GetTimestamp().Format("2006-01-02 15:04:05"))
    }
}

// AutoSaveManager 自动保存管理器
type AutoSaveManager struct {
    game        *Game
    autoSaveSlot string
    interval     time.Duration
    lastSave     time.Time
}

func NewAutoSaveManager(game *Game, interval time.Duration) *AutoSaveManager {
    return &AutoSaveManager{
        game:         game,
        autoSaveSlot: "autosave",
        interval:     interval,
        lastSave:     time.Now(),
    }
}

func (a *AutoSaveManager) CheckAutoSave() {
    if time.Since(a.lastSave) >= a.interval {
        a.game.SaveGame(a.autoSaveSlot)
        a.lastSave = time.Now()
    }
}

func (a *AutoSaveManager) LoadAutoSave() bool {
    return a.game.LoadGame(a.autoSaveSlot)
}

```

## 3.3.1.5 4. 工程案例

### 3.3.1.5.1 数据库事务备忘录模式

```go
package transactionmemento

import (
    "fmt"
    "time"
)

// TransactionState 事务状态
type TransactionState struct {
    operations []Operation
    timestamp  time.Time
    id         string
}

type Operation struct {
    Type      string
    Table     string
    Data      map[string]interface{}
    Timestamp time.Time
}

func NewTransactionState(id string) *TransactionState {
    return &TransactionState{
        operations: make([]Operation, 0),
        timestamp:  time.Now(),
        id:         id,
    }
}

func (t *TransactionState) AddOperation(opType, table string, data map[string]interface{}) {
    operation := Operation{
        Type:      opType,
        Table:     table,
        Data:      data,
        Timestamp: time.Now(),
    }
    t.operations = append(t.operations, operation)
}

func (t *TransactionState) GetOperations() []Operation {
    return t.operations
}

func (t *TransactionState) GetID() string {
    return t.id
}

// TransactionMemento 事务备忘录
type TransactionMemento struct {
    state       *TransactionState
    description string
}

func NewTransactionMemento(state *TransactionState, description string) *TransactionMemento {
    return &TransactionMemento{
        state:       state,
        description: description,
    }
}

func (t *TransactionMemento) GetState() *TransactionState {
    return t.state
}

func (t *TransactionMemento) GetDescription() string {
    return t.description
}

// Database 数据库
type Database struct {
    tables map[string]map[string]interface{}
    currentTransaction *TransactionState
}

func NewDatabase() *Database {
    return &Database{
        tables: make(map[string]map[string]interface{}),
    }
}

func (d *Database) BeginTransaction() *TransactionState {
    transactionID := fmt.Sprintf("txn_%d", time.Now().Unix())
    d.currentTransaction = NewTransactionState(transactionID)
    fmt.Printf("Database: Started transaction %s\n", transactionID)
    return d.currentTransaction
}

func (d *Database) Insert(table, id string, data map[string]interface{}) {
    if d.currentTransaction == nil {
        fmt.Println("Database: No active transaction")
        return
    }
    
    if d.tables[table] == nil {
        d.tables[table] = make(map[string]interface{})
    }
    
    d.tables[table][id] = data
    d.currentTransaction.AddOperation("INSERT", table, data)
    
    fmt.Printf("Database: Inserted into %s with id %s\n", table, id)
}

func (d *Database) Update(table, id string, data map[string]interface{}) {
    if d.currentTransaction == nil {
        fmt.Println("Database: No active transaction")
        return
    }
    
    if d.tables[table] != nil {
        d.tables[table][id] = data
        d.currentTransaction.AddOperation("UPDATE", table, data)
        fmt.Printf("Database: Updated %s with id %s\n", table, id)
    }
}

func (d *Database) Delete(table, id string) {
    if d.currentTransaction == nil {
        fmt.Println("Database: No active transaction")
        return
    }
    
    if d.tables[table] != nil {
        delete(d.tables[table], id)
        d.currentTransaction.AddOperation("DELETE", table, map[string]interface{}{"id": id})
        fmt.Printf("Database: Deleted from %s with id %s\n", table, id)
    }
}

func (d *Database) Commit() {
    if d.currentTransaction == nil {
        fmt.Println("Database: No active transaction to commit")
        return
    }
    
    fmt.Printf("Database: Committed transaction %s\n", d.currentTransaction.GetID())
    d.currentTransaction = nil
}

func (d *Database) Rollback() {
    if d.currentTransaction == nil {
        fmt.Println("Database: No active transaction to rollback")
        return
    }
    
    fmt.Printf("Database: Rolling back transaction %s\n", d.currentTransaction.GetID())
    d.currentTransaction = nil
}

func (d *Database) CreateCheckpoint() *TransactionMemento {
    if d.currentTransaction == nil {
        return nil
    }
    
    description := fmt.Sprintf("Transaction %s with %d operations", 
        d.currentTransaction.GetID(), len(d.currentTransaction.GetOperations()))
    
    return NewTransactionMemento(d.currentTransaction, description)
}

func (d *Database) RestoreFromCheckpoint(memento *TransactionMemento) {
    if memento != nil {
        d.currentTransaction = memento.GetState()
        fmt.Printf("Database: Restored from checkpoint: %s\n", memento.GetDescription())
    }
}

```

## 3.3.1.6 5. 批判性分析

### 3.3.1.6.1 优势

1. **状态保存**: 可以保存对象的完整状态
2. **状态恢复**: 可以恢复到任意历史状态
3. **封装保护**: 不破坏对象的封装性
4. **历史管理**: 支持撤销/重做功能

### 3.3.1.6.2 劣势

1. **内存开销**: 大量备忘录占用内存
2. **性能影响**: 状态序列化/反序列化开销
3. **存储复杂**: 复杂对象状态难以序列化
4. **版本管理**: 备忘录版本管理复杂

### 3.3.1.6.3 行业对比

| 语言 | 实现方式 | 性能 | 复杂度 |
|------|----------|------|--------|
| Go | 结构体 | 高 | 中 |
| Java | 序列化 | 中 | 中 |
| C++ | 深拷贝 | 高 | 高 |
| Python | pickle | 中 | 低 |

### 3.3.1.6.4 最新趋势

1. **增量保存**: 只保存状态变化
2. **压缩存储**: 压缩备忘录数据
3. **分布式存储**: 分布式备忘录存储
4. **版本控制**: 集成版本控制系统

## 3.3.1.7 6. 面试题与考点

### 3.3.1.7.1 基础考点

1. **Q**: 备忘录模式与命令模式的区别？
   **A**: 备忘录关注状态保存，命令关注操作封装

2. **Q**: 什么时候使用备忘录模式？
   **A**: 需要撤销/重做、状态恢复功能时

3. **Q**: 备忘录模式的优缺点？
   **A**: 优点：状态保存、恢复；缺点：内存开销、性能影响

### 3.3.1.7.2 进阶考点

1. **Q**: 如何处理大量备忘录的内存问题？
   **A**: 增量保存、压缩存储、定期清理

2. **Q**: 备忘录模式在分布式系统中的应用？
   **A**: 状态同步、故障恢复、数据备份

3. **Q**: 如何优化备忘录的性能？
   **A**: 延迟保存、批量处理、缓存机制

## 3.3.1.8 7. 术语表

| 术语 | 定义 | 英文 |
|------|------|------|
| 备忘录模式 | 保存对象状态的设计模式 | Memento Pattern |
| 发起人 | 需要保存状态的对象 | Originator |
| 备忘录 | 保存状态的对象 | Memento |
| 管理者 | 管理备忘录的对象 | Caretaker |

## 3.3.1.9 8. 常见陷阱

| 陷阱 | 问题 | 解决方案 |
|------|------|----------|
| 内存泄漏 | 大量备忘录占用内存 | 定期清理、增量保存 |
| 性能问题 | 状态序列化开销大 | 延迟保存、压缩存储 |
| 存储复杂 | 复杂状态难以序列化 | 简化状态、分步保存 |
| 版本冲突 | 备忘录版本管理困难 | 版本控制、时间戳 |

## 3.3.1.10 9. 相关主题

- [观察者模式](./01-Observer-Pattern.md)
- [策略模式](./02-Strategy-Pattern.md)
- [命令模式](./03-Command-Pattern.md)
- [状态模式](./04-State-Pattern.md)
- [模板方法模式](./05-Template-Method-Pattern.md)
- [迭代器模式](./06-Iterator-Pattern.md)
- [访问者模式](./07-Visitor-Pattern.md)
- [中介者模式](./08-Mediator-Pattern.md)

## 3.3.1.11 10. 学习路径

### 3.3.1.11.1 新手路径

1. 理解备忘录模式的基本概念
2. 学习发起人和备忘录的关系
3. 实现简单的备忘录模式
4. 理解状态保存和恢复机制

### 3.3.1.11.2 进阶路径

1. 学习复杂的备忘录实现
2. 理解备忘录的性能优化
3. 掌握备忘录的应用场景
4. 学习备忘录的最佳实践

### 3.3.1.11.3 高阶路径

1. 分析备忘录在大型项目中的应用
2. 理解备忘录与架构设计的关系
3. 掌握备忘录的性能调优
4. 学习备忘录的替代方案

---

**相关文档**: [行为型模式总览](./README.md) | [设计模式总览](../README.md)

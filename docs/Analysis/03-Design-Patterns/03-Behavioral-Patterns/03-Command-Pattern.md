# 命令模式 (Command Pattern)

## 目录

1. [概述](#1-概述)
2. [理论基础](#2-理论基础)
3. [Go语言实现](#3-go语言实现)
4. [工程案例](#4-工程案例)
5. [批判性分析](#5-批判性分析)
6. [面试题与考点](#6-面试题与考点)
7. [术语表](#7-术语表)
8. [常见陷阱](#8-常见陷阱)
9. [相关主题](#9-相关主题)
10. [学习路径](#10-学习路径)

## 1. 概述

### 1.1 定义

命令模式将请求封装成对象，从而可以用不同的请求对客户进行参数化，对请求排队或记录请求日志，以及支持可撤销的操作。

**形式化定义**:
$$Command = (Command, ConcreteCommand, Invoker, Receiver, Client)$$

其中：
- $Command$ 是命令接口
- $ConcreteCommand$ 是具体命令
- $Invoker$ 是调用者
- $Receiver$ 是接收者
- $Client$ 是客户端

### 1.2 核心特征

- **请求封装**: 将请求封装成对象
- **参数化**: 可以用不同请求参数化客户端
- **队列支持**: 支持请求排队
- **撤销操作**: 支持可撤销的操作

## 2. 理论基础

### 2.1 数学形式化

**定义 2.1** (命令模式): 命令模式是一个五元组 $C = (Cmd, Inv, Rec, Exe, Undo)$

其中：
- $Cmd$ 是命令集合
- $Inv$ 是调用者集合
- $Rec$ 是接收者集合
- $Exe$ 是执行函数，$Exe: Cmd \times Rec \rightarrow Result$
- $Undo$ 是撤销函数，$Undo: Cmd \times Rec \rightarrow Result$

**定理 2.1** (命令可撤销性): 对于任意命令 $c \in Cmd$ 和接收者 $r \in Rec$，如果存在撤销函数，则命令可以撤销。

### 2.2 范畴论视角

在范畴论中，命令模式可以表示为：

$$Command : Invoker \times Receiver \rightarrow Result$$

## 3. Go语言实现

### 3.1 基础命令模式

```go
package command

import "fmt"

// Command 命令接口
type Command interface {
    Execute()
    Undo()
    GetName() string
}

// Receiver 接收者
type Receiver struct {
    name string
}

func NewReceiver(name string) *Receiver {
    return &Receiver{name: name}
}

func (r *Receiver) Action1() {
    fmt.Printf("Receiver %s performing Action1\n", r.name)
}

func (r *Receiver) Action2() {
    fmt.Printf("Receiver %s performing Action2\n", r.name)
}

// ConcreteCommand1 具体命令1
type ConcreteCommand1 struct {
    receiver *Receiver
}

func NewConcreteCommand1(receiver *Receiver) *ConcreteCommand1 {
    return &ConcreteCommand1{receiver: receiver}
}

func (c *ConcreteCommand1) Execute() {
    c.receiver.Action1()
}

func (c *ConcreteCommand1) Undo() {
    fmt.Printf("Undoing Action1 for receiver %s\n", c.receiver.name)
}

func (c *ConcreteCommand1) GetName() string {
    return "Action1"
}

// ConcreteCommand2 具体命令2
type ConcreteCommand2 struct {
    receiver *Receiver
}

func NewConcreteCommand2(receiver *Receiver) *ConcreteCommand2 {
    return &ConcreteCommand2{receiver: receiver}
}

func (c *ConcreteCommand2) Execute() {
    c.receiver.Action2()
}

func (c *ConcreteCommand2) Undo() {
    fmt.Printf("Undoing Action2 for receiver %s\n", c.receiver.name)
}

func (c *ConcreteCommand2) GetName() string {
    return "Action2"
}

// Invoker 调用者
type Invoker struct {
    commands []Command
    history  []Command
}

func NewInvoker() *Invoker {
    return &Invoker{
        commands: make([]Command, 0),
        history:  make([]Command, 0),
    }
}

func (i *Invoker) AddCommand(command Command) {
    i.commands = append(i.commands, command)
}

func (i *Invoker) ExecuteCommands() {
    for _, command := range i.commands {
        command.Execute()
        i.history = append(i.history, command)
    }
    i.commands = make([]Command, 0)
}

func (i *Invoker) UndoLast() {
    if len(i.history) > 0 {
        lastCommand := i.history[len(i.history)-1]
        lastCommand.Undo()
        i.history = i.history[:len(i.history)-1]
    }
}
```

### 3.2 文本编辑器命令模式

```go
package editorcommand

import (
    "fmt"
    "strings"
)

// Document 文档
type Document struct {
    content string
    cursor  int
}

func NewDocument() *Document {
    return &Document{
        content: "",
        cursor:  0,
    }
}

func (d *Document) Insert(text string) {
    d.content = d.content[:d.cursor] + text + d.content[d.cursor:]
    d.cursor += len(text)
}

func (d *Document) Delete(length int) string {
    if d.cursor+length > len(d.content) {
        length = len(d.content) - d.cursor
    }
    deleted := d.content[d.cursor : d.cursor+length]
    d.content = d.content[:d.cursor] + d.content[d.cursor+length:]
    return deleted
}

func (d *Document) MoveCursor(position int) {
    if position >= 0 && position <= len(d.content) {
        d.cursor = position
    }
}

func (d *Document) GetContent() string {
    return d.content
}

func (d *Document) GetCursor() int {
    return d.cursor
}

// EditorCommand 编辑器命令接口
type EditorCommand interface {
    Execute()
    Undo()
    GetName() string
}

// InsertCommand 插入命令
type InsertCommand struct {
    document *Document
    text     string
    position int
}

func NewInsertCommand(document *Document, text string) *InsertCommand {
    return &InsertCommand{
        document: document,
        text:     text,
        position: document.GetCursor(),
    }
}

func (i *InsertCommand) Execute() {
    i.document.Insert(i.text)
}

func (i *InsertCommand) Undo() {
    i.document.MoveCursor(i.position)
    i.document.Delete(len(i.text))
}

func (i *InsertCommand) GetName() string {
    return fmt.Sprintf("Insert '%s'", i.text)
}

// DeleteCommand 删除命令
type DeleteCommand struct {
    document *Document
    length   int
    position int
    deleted  string
}

func NewDeleteCommand(document *Document, length int) *DeleteCommand {
    return &DeleteCommand{
        document: document,
        length:   length,
        position: document.GetCursor(),
    }
}

func (d *DeleteCommand) Execute() {
    d.deleted = d.document.Delete(d.length)
}

func (d *DeleteCommand) Undo() {
    d.document.MoveCursor(d.position)
    d.document.Insert(d.deleted)
}

func (d *DeleteCommand) GetName() string {
    return fmt.Sprintf("Delete %d characters", d.length)
}

// TextEditor 文本编辑器
type TextEditor struct {
    document *Document
    history  []EditorCommand
    redoStack []EditorCommand
}

func NewTextEditor() *TextEditor {
    return &TextEditor{
        document:  NewDocument(),
        history:   make([]EditorCommand, 0),
        redoStack: make([]EditorCommand, 0),
    }
}

func (t *TextEditor) Insert(text string) {
    command := NewInsertCommand(t.document, text)
    command.Execute()
    t.history = append(t.history, command)
    t.redoStack = make([]EditorCommand, 0) // 清空重做栈
}

func (t *TextEditor) Delete(length int) {
    command := NewDeleteCommand(t.document, length)
    command.Execute()
    t.history = append(t.history, command)
    t.redoStack = make([]EditorCommand, 0)
}

func (t *TextEditor) Undo() {
    if len(t.history) > 0 {
        lastCommand := t.history[len(t.history)-1]
        lastCommand.Undo()
        t.history = t.history[:len(t.history)-1]
        t.redoStack = append(t.redoStack, lastCommand)
    }
}

func (t *TextEditor) Redo() {
    if len(t.redoStack) > 0 {
        lastCommand := t.redoStack[len(t.redoStack)-1]
        lastCommand.Execute()
        t.redoStack = t.redoStack[:len(t.redoStack)-1]
        t.history = append(t.history, lastCommand)
    }
}

func (t *TextEditor) GetContent() string {
    return t.document.GetContent()
}
```

### 3.3 遥控器命令模式

```go
package remotecommand

import "fmt"

// Device 设备接口
type Device interface {
    TurnOn()
    TurnOff()
    GetName() string
}

// Light 灯
type Light struct {
    name  string
    isOn  bool
}

func NewLight(name string) *Light {
    return &Light{
        name: name,
        isOn: false,
    }
}

func (l *Light) TurnOn() {
    l.isOn = true
    fmt.Printf("%s light is now ON\n", l.name)
}

func (l *Light) TurnOff() {
    l.isOn = false
    fmt.Printf("%s light is now OFF\n", l.name)
}

func (l *Light) GetName() string {
    return l.name
}

func (l *Light) IsOn() bool {
    return l.isOn
}

// TV 电视
type TV struct {
    name  string
    isOn  bool
    channel int
}

func NewTV(name string) *TV {
    return &TV{
        name:    name,
        isOn:    false,
        channel: 1,
    }
}

func (t *TV) TurnOn() {
    t.isOn = true
    fmt.Printf("%s TV is now ON (Channel %d)\n", t.name, t.channel)
}

func (t *TV) TurnOff() {
    t.isOn = false
    fmt.Printf("%s TV is now OFF\n", t.name)
}

func (t *TV) ChangeChannel(channel int) {
    if t.isOn {
        t.channel = channel
        fmt.Printf("%s TV changed to channel %d\n", t.name, channel)
    }
}

func (t *TV) GetName() string {
    return t.name
}

func (t *TV) IsOn() bool {
    return t.isOn
}

// RemoteCommand 遥控器命令接口
type RemoteCommand interface {
    Execute()
    Undo()
    GetName() string
}

// LightOnCommand 开灯命令
type LightOnCommand struct {
    light *Light
}

func NewLightOnCommand(light *Light) *LightOnCommand {
    return &LightOnCommand{light: light}
}

func (l *LightOnCommand) Execute() {
    l.light.TurnOn()
}

func (l *LightOnCommand) Undo() {
    l.light.TurnOff()
}

func (l *LightOnCommand) GetName() string {
    return fmt.Sprintf("Turn on %s light", l.light.GetName())
}

// LightOffCommand 关灯命令
type LightOffCommand struct {
    light *Light
}

func NewLightOffCommand(light *Light) *LightOffCommand {
    return &LightOffCommand{light: light}
}

func (l *LightOffCommand) Execute() {
    l.light.TurnOff()
}

func (l *LightOffCommand) Undo() {
    l.light.TurnOn()
}

func (l *LightOffCommand) GetName() string {
    return fmt.Sprintf("Turn off %s light", l.light.GetName())
}

// TVOnCommand 开电视命令
type TVOnCommand struct {
    tv *TV
}

func NewTVOnCommand(tv *TV) *TVOnCommand {
    return &TVOnCommand{tv: tv}
}

func (t *TVOnCommand) Execute() {
    t.tv.TurnOn()
}

func (t *TVOnCommand) Undo() {
    t.tv.TurnOff()
}

func (t *TVOnCommand) GetName() string {
    return fmt.Sprintf("Turn on %s TV", t.tv.GetName())
}

// TVOffCommand 关电视命令
type TVOffCommand struct {
    tv *TV
}

func NewTVOffCommand(tv *TV) *TVOffCommand {
    return &TVOffCommand{tv: tv}
}

func (t *TVOffCommand) Execute() {
    t.tv.TurnOff()
}

func (t *TVOffCommand) Undo() {
    t.tv.TurnOn()
}

func (t *TVOffCommand) GetName() string {
    return fmt.Sprintf("Turn off %s TV", t.tv.GetName())
}

// RemoteControl 遥控器
type RemoteControl struct {
    onCommands  []RemoteCommand
    offCommands []RemoteCommand
    undoCommand RemoteCommand
}

func NewRemoteControl() *RemoteControl {
    return &RemoteControl{
        onCommands:  make([]RemoteCommand, 7),
        offCommands: make([]RemoteCommand, 7),
    }
}

func (r *RemoteControl) SetCommand(slot int, onCommand, offCommand RemoteCommand) {
    r.onCommands[slot] = onCommand
    r.offCommands[slot] = offCommand
}

func (r *RemoteControl) PressOnButton(slot int) {
    if r.onCommands[slot] != nil {
        r.onCommands[slot].Execute()
        r.undoCommand = r.onCommands[slot]
    }
}

func (r *RemoteControl) PressOffButton(slot int) {
    if r.offCommands[slot] != nil {
        r.offCommands[slot].Execute()
        r.undoCommand = r.offCommands[slot]
    }
}

func (r *RemoteControl) PressUndoButton() {
    if r.undoCommand != nil {
        r.undoCommand.Undo()
    }
}
```

## 4. 工程案例

### 4.1 数据库事务命令模式

```go
package transactioncommand

import "fmt"

// Database 数据库接口
type Database interface {
    Insert(table string, data map[string]interface{}) error
    Update(table string, id int, data map[string]interface{}) error
    Delete(table string, id int) error
    GetName() string
}

// MockDatabase 模拟数据库
type MockDatabase struct {
    name   string
    tables map[string][]map[string]interface{}
    nextID map[string]int
}

func NewMockDatabase(name string) *MockDatabase {
    return &MockDatabase{
        name:   name,
        tables: make(map[string][]map[string]interface{}),
        nextID: make(map[string]int),
    }
}

func (m *MockDatabase) Insert(table string, data map[string]interface{}) error {
    if m.tables[table] == nil {
        m.tables[table] = make([]map[string]interface{}, 0)
        m.nextID[table] = 1
    }
    
    data["id"] = m.nextID[table]
    m.tables[table] = append(m.tables[table], data)
    m.nextID[table]++
    
    fmt.Printf("Inserted into %s: %v\n", table, data)
    return nil
}

func (m *MockDatabase) Update(table string, id int, data map[string]interface{}) error {
    if m.tables[table] == nil {
        return fmt.Errorf("table %s does not exist", table)
    }
    
    for i, row := range m.tables[table] {
        if row["id"] == id {
            for key, value := range data {
                row[key] = value
            }
            fmt.Printf("Updated %s id %d: %v\n", table, id, row)
            return nil
        }
    }
    
    return fmt.Errorf("record with id %d not found in table %s", id, table)
}

func (m *MockDatabase) Delete(table string, id int) error {
    if m.tables[table] == nil {
        return fmt.Errorf("table %s does not exist", table)
    }
    
    for i, row := range m.tables[table] {
        if row["id"] == id {
            m.tables[table] = append(m.tables[table][:i], m.tables[table][i+1:]...)
            fmt.Printf("Deleted from %s id %d\n", table, id)
            return nil
        }
    }
    
    return fmt.Errorf("record with id %d not found in table %s", id, table)
}

func (m *MockDatabase) GetName() string {
    return m.name
}

// TransactionCommand 事务命令接口
type TransactionCommand interface {
    Execute() error
    Rollback() error
    GetName() string
}

// InsertCommand 插入命令
type InsertCommand struct {
    database Database
    table    string
    data     map[string]interface{}
    insertedID int
}

func NewInsertCommand(database Database, table string, data map[string]interface{}) *InsertCommand {
    return &InsertCommand{
        database: database,
        table:    table,
        data:     data,
    }
}

func (i *InsertCommand) Execute() error {
    err := i.database.Insert(i.table, i.data)
    if err == nil {
        i.insertedID = i.data["id"].(int)
    }
    return err
}

func (i *InsertCommand) Rollback() error {
    if i.insertedID > 0 {
        return i.database.Delete(i.table, i.insertedID)
    }
    return nil
}

func (i *InsertCommand) GetName() string {
    return fmt.Sprintf("Insert into %s", i.table)
}

// UpdateCommand 更新命令
type UpdateCommand struct {
    database Database
    table    string
    id       int
    data     map[string]interface{}
    oldData  map[string]interface{}
}

func NewUpdateCommand(database Database, table string, id int, data map[string]interface{}) *UpdateCommand {
    return &UpdateCommand{
        database: database,
        table:    table,
        id:       id,
        data:     data,
    }
}

func (u *UpdateCommand) Execute() error {
    // 在实际实现中，这里需要先获取旧数据
    u.oldData = make(map[string]interface{})
    return u.database.Update(u.table, u.id, u.data)
}

func (u *UpdateCommand) Rollback() error {
    if u.oldData != nil {
        return u.database.Update(u.table, u.id, u.oldData)
    }
    return nil
}

func (u *UpdateCommand) GetName() string {
    return fmt.Sprintf("Update %s id %d", u.table, u.id)
}

// Transaction 事务
type Transaction struct {
    commands []TransactionCommand
    executed []TransactionCommand
}

func NewTransaction() *Transaction {
    return &Transaction{
        commands: make([]TransactionCommand, 0),
        executed: make([]TransactionCommand, 0),
    }
}

func (t *Transaction) AddCommand(command TransactionCommand) {
    t.commands = append(t.commands, command)
}

func (t *Transaction) Commit() error {
    for _, command := range t.commands {
        if err := command.Execute(); err != nil {
            // 回滚已执行的命令
            t.rollback()
            return err
        }
        t.executed = append(t.executed, command)
    }
    
    fmt.Println("Transaction committed successfully")
    return nil
}

func (t *Transaction) rollback() {
    fmt.Println("Rolling back transaction...")
    for i := len(t.executed) - 1; i >= 0; i-- {
        t.executed[i].Rollback()
    }
    t.executed = make([]TransactionCommand, 0)
}
```

## 5. 批判性分析

### 5.1 优势

1. **请求封装**: 将请求封装成对象
2. **参数化**: 可以用不同请求参数化客户端
3. **队列支持**: 支持请求排队
4. **撤销操作**: 支持可撤销的操作

### 5.2 劣势

1. **内存开销**: 每个命令都需要创建对象
2. **复杂度增加**: 增加了系统的复杂度
3. **调试困难**: 命令链可能难以调试
4. **性能影响**: 大量命令可能影响性能

### 5.3 行业对比

| 语言 | 实现方式 | 性能 | 复杂度 |
|------|----------|------|--------|
| Go | 接口 | 高 | 中 |
| Java | 接口 | 中 | 中 |
| C++ | 虚函数 | 高 | 中 |
| Python | 函数对象 | 中 | 低 |

### 5.4 最新趋势

1. **函数式命令**: 使用函数作为命令
2. **命令队列**: 异步命令处理
3. **命令持久化**: 命令的持久化存储
4. **微服务命令**: 分布式命令模式

## 6. 面试题与考点

### 6.1 基础考点

1. **Q**: 命令模式与策略模式的区别？
   **A**: 命令模式关注请求封装，策略模式关注算法选择

2. **Q**: 什么时候使用命令模式？
   **A**: 需要撤销操作、请求排队、日志记录时

3. **Q**: 命令模式的优缺点？
   **A**: 优点：请求封装、撤销支持；缺点：内存开销、复杂度增加

### 6.2 进阶考点

1. **Q**: 如何实现命令的撤销？
   **A**: 保存状态、反向操作、命令历史

2. **Q**: 命令模式在微服务中的应用？
   **A**: 事件溯源、CQRS、分布式事务

3. **Q**: 如何处理命令的性能问题？
   **A**: 命令池、异步处理、批量操作

## 7. 术语表

| 术语 | 定义 | 英文 |
|------|------|------|
| 命令模式 | 封装请求的设计模式 | Command Pattern |
| 命令 | 封装请求的对象 | Command |
| 调用者 | 执行命令的对象 | Invoker |
| 接收者 | 执行具体操作的对象 | Receiver |

## 8. 常见陷阱

| 陷阱 | 问题 | 解决方案 |
|------|------|----------|
| 内存泄漏 | 命令对象无法释放 | 正确管理生命周期 |
| 撤销复杂 | 撤销操作实现复杂 | 状态保存、反向操作 |
| 性能问题 | 大量命令影响性能 | 命令池、异步处理 |
| 调试困难 | 命令链难以调试 | 日志记录、可视化 |

## 9. 相关主题

- [观察者模式](./01-Observer-Pattern.md)
- [策略模式](./02-Strategy-Pattern.md)
- [状态模式](./04-State-Pattern.md)
- [模板方法模式](./05-Template-Method-Pattern.md)
- [迭代器模式](./06-Iterator-Pattern.md)

## 10. 学习路径

### 10.1 新手路径

1. 理解命令模式的基本概念
2. 学习命令的封装和执行
3. 实现简单的命令模式
4. 理解撤销机制

### 10.2 进阶路径

1. 学习复杂的命令实现
2. 理解命令的性能优化
3. 掌握命令的应用场景
4. 学习命令的最佳实践

### 10.3 高阶路径

1. 分析命令在大型项目中的应用
2. 理解命令与架构设计的关系
3. 掌握命令的性能调优
4. 学习命令的替代方案

---

**相关文档**: [行为型模式总览](./README.md) | [设计模式总览](../README.md) 
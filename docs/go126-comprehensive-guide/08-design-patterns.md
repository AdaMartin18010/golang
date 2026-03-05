# 第八章：23种设计模式 Go 实现

> 基于 GoF (Gang of Four) 设计模式在 Go 语言中的实现

---

## 8.1 创建型模式 (Creational Patterns)

### 8.1.1 单例模式 (Singleton)

```go
package singleton

import (
    "sync"
)

// 饿汉式单例（推荐）
type Config struct {
    DatabaseURL string
    Port        int
}

var (
    instance *Config
    once     sync.Once
)

// GetInstance 返回单例实例
func GetInstance() *Config {
    once.Do(func() {
        instance = &Config{
            DatabaseURL: "localhost:5432",
            Port:        8080,
        }
    })
    return instance
}

// 测试
type singletonTest struct{}

func (t singletonTest) Test() {
    c1 := GetInstance()
    c2 := GetInstance()
    // c1 == c2 (同一个实例)
}
```

### 8.1.2 工厂方法 (Factory Method)

```go
package factory

// 产品接口
type Animal interface {
    Speak() string
}

// 具体产品
type Dog struct{}

func (d Dog) Speak() string { return "Woof!" }

type Cat struct{}

func (c Cat) Speak() string { return "Meow!" }

// 工厂接口
type AnimalFactory interface {
    Create() Animal
}

// 具体工厂
type DogFactory struct{}

func (f DogFactory) Create() Animal { return Dog{} }

type CatFactory struct{}

func (f CatFactory) Create() Animal { return Cat{} }

// 使用
func UseFactory() {
    var factory AnimalFactory = DogFactory{}
    animal := factory.Create()
    println(animal.Speak()) // Woof!
}
```

### 8.1.3 抽象工厂 (Abstract Factory)

```go
package abstractfactory

// UI 组件家族
type Button interface {
    Render() string
}

type Checkbox interface {
    Check() string
}

// Windows 风格
type WindowsButton struct{}

func (b WindowsButton) Render() string { return "Windows Button" }

type WindowsCheckbox struct{}

func (c WindowsCheckbox) Check() string { return "Windows Checkbox" }

// Mac 风格
type MacButton struct{}

func (b MacButton) Render() string { return "Mac Button" }

type MacCheckbox struct{}

func (c MacCheckbox) Check() string { return "Mac Checkbox" }

// 抽象工厂
type GUIFactory interface {
    CreateButton() Button
    CreateCheckbox() Checkbox
}

// Windows 工厂
type WindowsFactory struct{}

func (f WindowsFactory) CreateButton() Button     { return WindowsButton{} }
func (f WindowsFactory) CreateCheckbox() Checkbox { return WindowsCheckbox{} }

// Mac 工厂
type MacFactory struct{}

func (f MacFactory) CreateButton() Button     { return MacButton{} }
func (f MacFactory) CreateCheckbox() Checkbox { return MacCheckbox{} }
```

### 8.1.4 建造者模式 (Builder)

```go
package builder

// 产品
type House struct {
    Windows int
    Doors   int
    Garage  bool
    Pool    bool
}

// 建造者
type HouseBuilder struct {
    house House
}

func NewHouseBuilder() *HouseBuilder {
    return &HouseBuilder{house: House{}}
}

func (b *HouseBuilder) SetWindows(n int) *HouseBuilder {
    b.house.Windows = n
    return b
}

func (b *HouseBuilder) SetDoors(n int) *HouseBuilder {
    b.house.Doors = n
    return b
}

func (b *HouseBuilder) AddGarage() *HouseBuilder {
    b.house.Garage = true
    return b
}

func (b *HouseBuilder) AddPool() *HouseBuilder {
    b.house.Pool = true
    return b
}

func (b *HouseBuilder) Build() House {
    return b.house
}

// 使用
func UseBuilder() {
    house := NewHouseBuilder().
        SetWindows(4).
        SetDoors(2).
        AddGarage().
        AddPool().
        Build()

    println(house.Windows) // 4
}
```

### 8.1.5 原型模式 (Prototype)

```go
package prototype

// 克隆接口
type Cloneable interface {
    Clone() Cloneable
}

// 具体原型
type Document struct {
    Title   string
    Content string
    Author  string
}

func (d *Document) Clone() Cloneable {
    // 深拷贝
    return &Document{
        Title:   d.Title,
        Content: d.Content,
        Author:  d.Author,
    }
}

// 原型注册表
type Registry struct {
    prototypes map[string]Cloneable
}

func NewRegistry() *Registry {
    return &Registry{
        prototypes: make(map[string]Cloneable),
    }
}

func (r *Registry) Register(name string, prototype Cloneable) {
    r.prototypes[name] = prototype
}

func (r *Registry) Create(name string) Cloneable {
    if prototype, ok := r.prototypes[name]; ok {
        return prototype.Clone()
    }
    return nil
}
```

---

## 8.2 结构型模式 (Structural Patterns)

### 8.2.1 适配器模式 (Adapter)

```go
package adapter

// 目标接口
type MediaPlayer interface {
    Play(audioType string, fileName string)
}

// 被适配者
type AdvancedMediaPlayer interface {
    PlayVlc(fileName string)
    PlayMp4(fileName string)
}

type VlcPlayer struct{}

func (v VlcPlayer) PlayVlc(fileName string) {
    println("Playing vlc file:", fileName)
}
func (v VlcPlayer) PlayMp4(fileName string) {}

type Mp4Player struct{}

func (m Mp4Player) PlayVlc(fileName string) {}
func (m Mp4Player) PlayMp4(fileName string) {
    println("Playing mp4 file:", fileName)
}

// 适配器
type MediaAdapter struct {
    advancedPlayer AdvancedMediaPlayer
}

func NewMediaAdapter(audioType string) *MediaAdapter {
    switch audioType {
    case "vlc":
        return &MediaAdapter{advancedPlayer: VlcPlayer{}}
    case "mp4":
        return &MediaAdapter{advancedPlayer: Mp4Player{}}
    }
    return nil
}

func (a *MediaAdapter) Play(audioType string, fileName string) {
    switch audioType {
    case "vlc":
        a.advancedPlayer.PlayVlc(fileName)
    case "mp4":
        a.advancedPlayer.PlayMp4(fileName)
    }
}

// 客户端
type AudioPlayer struct {
    mediaAdapter *MediaAdapter
}

func (a *AudioPlayer) Play(audioType string, fileName string) {
    if audioType == "mp3" {
        println("Playing mp3 file:", fileName)
    } else if audioType == "vlc" || audioType == "mp4" {
        a.mediaAdapter = NewMediaAdapter(audioType)
        a.mediaAdapter.Play(audioType, fileName)
    }
}
```

### 8.2.2 桥接模式 (Bridge)

```go
package bridge

// 实现接口
type DrawAPI interface {
    DrawCircle(radius int, x int, y int)
}

// 具体实现 1
type RedCircle struct{}

func (r RedCircle) DrawCircle(radius int, x int, y int) {
    println("Drawing Circle[color: red, radius:", radius, ", x:", x, ", y:", y, "]")
}

// 具体实现 2
type GreenCircle struct{}

func (g GreenCircle) DrawCircle(radius int, x int, y int) {
    println("Drawing Circle[color: green, radius:", radius, ", x:", x, ", y:", y, "]")
}

// 抽象
type Shape struct {
    drawAPI DrawAPI
}

func (s *Shape) Shape(drawAPI DrawAPI) {
    s.drawAPI = drawAPI
}

// 扩展抽象
type Circle struct {
    shape  Shape
    x      int
    y      int
    radius int
}

func NewCircle(x, y, radius int, drawAPI DrawAPI) *Circle {
    c := &Circle{x: x, y: y, radius: radius}
    c.shape.Shape(drawAPI)
    return c
}

func (c *Circle) Draw() {
    c.shape.drawAPI.DrawCircle(c.radius, c.x, c.y)
}
```

### 8.2.3 组合模式 (Composite)

```go
package composite

import "fmt"

// 组件接口
type Employee interface {
    ShowDetails()
}

// 叶子节点
type Developer struct {
    Name string
}

func (d Developer) ShowDetails() {
    fmt.Println("Developer:", d.Name)
}

// 叶子节点
type Manager struct {
    Name string
}

func (m Manager) ShowDetails() {
    fmt.Println("Manager:", m.Name)
}

// 复合节点
type CompanyDirectory struct {
    Employees []Employee
}

func (c *CompanyDirectory) Add(emp Employee) {
    c.Employees = append(c.Employees, emp)
}

func (c *CompanyDirectory) Remove(emp Employee) {
    // 实现删除逻辑
}

func (c *CompanyDirectory) ShowDetails() {
    for _, emp := range c.Employees {
        emp.ShowDetails()
    }
}

// 使用
func UseComposite() {
    dev1 := Developer{Name: "John"}
    dev2 := Developer{Name: "Jane"}

    mgr := Manager{Name: "Bob"}

    team := CompanyDirectory{}
    team.Add(dev1)
    team.Add(dev2)
    team.Add(mgr)

    team.ShowDetails()
}
```

### 8.2.4 装饰器模式 (Decorator)

```go
package decorator

// 组件
type Pizza interface {
    Cost() int
    Description() string
}

// 具体组件
type PlainPizza struct{}

func (p PlainPizza) Cost() int        { return 10 }
func (p PlainPizza) Description() string { return "Plain Pizza" }

// 装饰器基础
type PizzaDecorator struct {
    pizza Pizza
}

func (d PizzaDecorator) Cost() int        { return d.pizza.Cost() }
func (d PizzaDecorator) Description() string { return d.pizza.Description() }

// 具体装饰器
type CheeseDecorator struct {
    PizzaDecorator
}

func NewCheeseDecorator(p Pizza) Pizza {
    return &CheeseDecorator{PizzaDecorator{pizza: p}}
}

func (c *CheeseDecorator) Cost() int {
    return c.pizza.Cost() + 2
}

func (c *CheeseDecorator) Description() string {
    return c.pizza.Description() + ", Cheese"
}

// 另一个装饰器
type PepperoniDecorator struct {
    PizzaDecorator
}

func NewPepperoniDecorator(p Pizza) Pizza {
    return &PepperoniDecorator{PizzaDecorator{pizza: p}}
}

func (p *PepperoniDecorator) Cost() int {
    return p.pizza.Cost() + 3
}

func (p *PepperoniDecorator) Description() string {
    return p.pizza.Description() + ", Pepperoni"
}

// 使用
func UseDecorator() {
    pizza := PlainPizza{}
    pizza = NewCheeseDecorator(pizza)
    pizza = NewPepperoniDecorator(pizza)

    println(pizza.Description()) // Plain Pizza, Cheese, Pepperoni
    println(pizza.Cost())        // 15
}
```

### 8.2.5 外观模式 (Facade)

```go
package facade

import "fmt"

// 子系统
type CPU struct{}

func (c *CPU) Freeze()     { fmt.Println("CPU freeze") }
func (c *CPU) Jump(pos int) { fmt.Println("CPU jump to", pos) }
func (c *CPU) Execute()    { fmt.Println("CPU execute") }

type Memory struct{}

func (m *Memory) Load(pos int, data string) {
    fmt.Println("Memory load", data, "at", pos)
}

type HardDrive struct{}

func (h *HardDrive) Read(lba int, size int) string {
    return fmt.Sprintf("data from %d size %d", lba, size)
}

// 外观
type ComputerFacade struct {
    cpu       CPU
    memory    Memory
    hardDrive HardDrive
}

func NewComputerFacade() *ComputerFacade {
    return &ComputerFacade{
        cpu:       CPU{},
        memory:    Memory{},
        hardDrive: HardDrive{},
    }
}

func (c *ComputerFacade) Start() {
    c.cpu.Freeze()
    data := c.hardDrive.Read(0, 1024)
    c.memory.Load(0, data)
    c.cpu.Jump(0)
    c.cpu.Execute()
}

// 使用
func UseFacade() {
    computer := NewComputerFacade()
    computer.Start() // 简化的启动流程
}
```

### 8.2.6 享元模式 (Flyweight)

```go
package flyweight

import "fmt"

// 享元接口
type Shape interface {
    Draw(x, y int)
}

// 具体享元
type Circle struct {
    Color string  // 内部状态
}

func (c *Circle) Draw(x, y int) {
    fmt.Printf("Draw %s circle at (%d, %d)\n", c.Color, x, y)
}

// 享元工厂
type ShapeFactory struct {
    circles map[string]*Circle
}

func NewShapeFactory() *ShapeFactory {
    return &ShapeFactory{circles: make(map[string]*Circle)}
}

func (f *ShapeFactory) GetCircle(color string) *Circle {
    if circle, ok := f.circles[color]; ok {
        return circle
    }
    circle := &Circle{Color: color}
    f.circles[color] = circle
    fmt.Println("Creating new", color, "circle")
    return circle
}

// 使用
func UseFlyweight() {
    factory := NewShapeFactory()

    // 获取或创建享元对象
    red := factory.GetCircle("red")
    red.Draw(10, 20)

    red2 := factory.GetCircle("red") // 复用已有对象
    red2.Draw(30, 40)
}
```

### 8.2.7 代理模式 (Proxy)

```go
package proxy

import "fmt"

// 主题
type Image interface {
    Display()
}

// 真实主题
type RealImage struct {
    filename string
}

func NewRealImage(filename string) *RealImage {
    img := &RealImage{filename: filename}
    img.loadFromDisk()
    return img
}

func (r *RealImage) loadFromDisk() {
    fmt.Println("Loading", r.filename)
}

func (r *RealImage) Display() {
    fmt.Println("Displaying", r.filename)
}

// 代理
type ProxyImage struct {
    filename  string
    realImage *RealImage
}

func NewProxyImage(filename string) *ProxyImage {
    return &ProxyImage{filename: filename}
}

func (p *ProxyImage) Display() {
    if p.realImage == nil {
        p.realImage = NewRealImage(p.filename)
    }
    p.realImage.Display()
}

// 使用
func UseProxy() {
    image := NewProxyImage("photo.jpg")
    // 此时还未加载

    image.Display() // 第一次，加载并显示
    image.Display() // 第二次，直接显示
}
```

---

## 8.3 行为型模式 (Behavioral Patterns)

### 8.3.1 责任链模式 (Chain of Responsibility)

```go
package chain

import "fmt"

// 处理者接口
type Handler interface {
    SetNext(handler Handler) Handler
    Handle(request int)
}

// 基础处理者
type BaseHandler struct {
    next Handler
}

func (b *BaseHandler) SetNext(handler Handler) Handler {
    b.next = handler
    return handler
}

func (b *BaseHandler) HandleNext(request int) {
    if b.next != nil {
        b.next.Handle(request)
    }
}

// 具体处理者
type ConcreteHandlerA struct {
    BaseHandler
}

func (h *ConcreteHandlerA) Handle(request int) {
    if request < 10 {
        fmt.Println("Handler A processed request", request)
    } else {
        h.HandleNext(request)
    }
}

type ConcreteHandlerB struct {
    BaseHandler
}

func (h *ConcreteHandlerB) Handle(request int) {
    if request >= 10 && request < 20 {
        fmt.Println("Handler B processed request", request)
    } else {
        h.HandleNext(request)
    }
}

// 使用
func UseChain() {
    handlerA := &ConcreteHandlerA{}
    handlerB := &ConcreteHandlerB{}

    handlerA.SetNext(handlerB)

    handlerA.Handle(5)  // Handler A
    handlerA.Handle(15) // Handler B
}
```

### 8.3.2 命令模式 (Command)

```go
package command

import "fmt"

// 命令接口
type Command interface {
    Execute()
    Undo()
}

// 接收者
type Light struct {
    isOn bool
}

func (l *Light) On() {
    l.isOn = true
    fmt.Println("Light is on")
}

func (l *Light) Off() {
    l.isOn = false
    fmt.Println("Light is off")
}

// 具体命令
type LightOnCommand struct {
    light *Light
}

func NewLightOnCommand(light *Light) *LightOnCommand {
    return &LightOnCommand{light: light}
}

func (c *LightOnCommand) Execute() {
    c.light.On()
}

func (c *LightOnCommand) Undo() {
    c.light.Off()
}

// 调用者
type RemoteControl struct {
    command Command
}

func (r *RemoteControl) SetCommand(cmd Command) {
    r.command = cmd
}

func (r *RemoteControl) PressButton() {
    r.command.Execute()
}
```

### 8.3.3 迭代器模式 (Iterator)

```go
package iterator

// 迭代器接口
type Iterator interface {
    HasNext() bool
    Next() interface{}
}

// 集合接口
type Container interface {
    GetIterator() Iterator
}

// 具体集合
type NameRepository struct {
    names []string
}

func NewNameRepository() *NameRepository {
    return &NameRepository{
        names: []string{"Robert", "John", "Julie", "Lora"},
    }
}

func (n *NameRepository) GetIterator() Iterator {
    return &NameIterator{repository: n, index: 0}
}

// 具体迭代器
type NameIterator struct {
    repository *NameRepository
    index      int
}

func (n *NameIterator) HasNext() bool {
    return n.index < len(n.repository.names)
}

func (n *NameIterator) Next() interface{} {
    if n.HasNext() {
        name := n.repository.names[n.index]
        n.index++
        return name
    }
    return nil
}

// Go 风格实现（使用 channel）
func (n *NameRepository) Iterate() <-chan string {
    ch := make(chan string)
    go func() {
        defer close(ch)
        for _, name := range n.names {
            ch <- name
        }
    }()
    return ch
}
```

### 8.3.4 观察者模式 (Observer)

```go
package observer

import "fmt"

// 观察者接口
type Observer interface {
    Update(temperature float64)
}

// 主题接口
type Subject interface {
    RegisterObserver(o Observer)
    RemoveObserver(o Observer)
    NotifyObservers()
}

// 具体主题
type WeatherStation struct {
    observers   []Observer
    temperature float64
}

func (w *WeatherStation) RegisterObserver(o Observer) {
    w.observers = append(w.observers, o)
}

func (w *WeatherStation) RemoveObserver(o Observer) {
    // 实现删除逻辑
}

func (w *WeatherStation) NotifyObservers() {
    for _, observer := range w.observers {
        observer.Update(w.temperature)
    }
}

func (w *WeatherStation) SetTemperature(temp float64) {
    w.temperature = temp
    w.NotifyObservers()
}

// 具体观察者
type PhoneDisplay struct {
    temperature float64
}

func (p *PhoneDisplay) Update(temp float64) {
    p.temperature = temp
    fmt.Println("Phone Display: Temperature updated to", temp)
}

type TVDisplay struct {
    temperature float64
}

func (t *TVDisplay) Update(temp float64) {
    t.temperature = temp
    fmt.Println("TV Display: Temperature updated to", temp)
}
```

### 8.3.5 策略模式 (Strategy)

```go
package strategy

import "fmt"

// 策略接口
type PaymentStrategy interface {
    Pay(amount int)
}

// 具体策略
type CreditCardPayment struct {
    cardNumber string
}

func NewCreditCardPayment(card string) *CreditCardPayment {
    return &CreditCardPayment{cardNumber: card}
}

func (c *CreditCardPayment) Pay(amount int) {
    fmt.Printf("Paid %d using Credit Card %s\n", amount, c.cardNumber)
}

type PayPalPayment struct {
    email string
}

func NewPayPalPayment(email string) *PayPalPayment {
    return &PayPalPayment{email: email}
}

func (p *PayPalPayment) Pay(amount int) {
    fmt.Printf("Paid %d using PayPal account %s\n", amount, p.email)
}

// 上下文
type ShoppingCart struct {
    strategy PaymentStrategy
}

func (s *ShoppingCart) SetPaymentStrategy(strategy PaymentStrategy) {
    s.strategy = strategy
}

func (s *ShoppingCart) Checkout(amount int) {
    s.strategy.Pay(amount)
}

// 使用
func UseStrategy() {
    cart := ShoppingCart{}

    cart.SetPaymentStrategy(NewCreditCardPayment("1234-5678"))
    cart.Checkout(100)

    cart.SetPaymentStrategy(NewPayPalPayment("user@example.com"))
    cart.Checkout(200)
}
```

### 8.3.6 模板方法模式 (Template Method)

```go
package template

import "fmt"

// 抽象类
type Game interface {
    Initialize()
    StartPlay()
    EndPlay()
    Play() // 模板方法
}

// 基础实现
type BaseGame struct{}

func (b *BaseGame) Play(game Game) {
    game.Initialize()
    game.StartPlay()
    game.EndPlay()
}

// 具体实现
type Football struct {
    BaseGame
}

func (f *Football) Initialize() {
    fmt.Println("Football Game Initialized")
}

func (f *Football) StartPlay() {
    fmt.Println("Football Game Started")
}

func (f *Football) EndPlay() {
    fmt.Println("Football Game Finished")
}

type Basketball struct {
    BaseGame
}

func (b *Basketball) Initialize() {
    fmt.Println("Basketball Game Initialized")
}

func (b *Basketball) StartPlay() {
    fmt.Println("Basketball Game Started")
}

func (b *Basketball) EndPlay() {
    fmt.Println("Basketball Game Finished")
}
```

---

## 8.4 Go 特有的模式简化

Go 语言特性可以简化某些设计模式：

```go
// 函数类型简化策略模式
type PaymentFunc func(amount int)

func CreditCardPay(card string) PaymentFunc {
    return func(amount int) {
        fmt.Printf("Pay %d with card %s\n", amount, card)
    }
}

// 使用
var pay PaymentFunc = CreditCardPay("1234")
pay(100)

// 闭包简化装饰器
func LogDecorator(fn func()) func() {
    return func() {
        fmt.Println("Before")
        fn()
        fmt.Println("After")
    }
}

// 接口组合简化适配器
// Go 的隐式接口让适配器更灵活
```

---

*注：其余行为型模式（解释器、中介者、备忘录、状态、访问者）因篇幅限制省略，但遵循类似实现方式。*

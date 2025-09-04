# 3.2.1 组合模式 (Composite Pattern)

<!-- TOC START -->
- [3.2.1 组合模式 (Composite Pattern)](#321-组合模式-composite-pattern)
  - [3.2.1.1 目录](#3211-目录)
  - [3.2.1.2 1. 概述](#3212-1-概述)
    - [3.2.1.2.1 定义](#32121-定义)
    - [3.2.1.2.2 核心特征](#32122-核心特征)
  - [3.2.1.3 2. 理论基础](#3213-2-理论基础)
    - [3.2.1.3.1 数学形式化](#32131-数学形式化)
    - [3.2.1.3.2 范畴论视角](#32132-范畴论视角)
  - [3.2.1.4 3. Go语言实现](#3214-3-go语言实现)
    - [3.2.1.4.1 基础组合模式](#32141-基础组合模式)
    - [3.2.1.4.2 文件系统组合模式](#32142-文件系统组合模式)
    - [3.2.1.4.3 UI组件组合模式](#32143-ui组件组合模式)
  - [3.2.1.5 4. 工程案例](#3215-4-工程案例)
    - [3.2.1.5.1 组织架构组合模式](#32151-组织架构组合模式)
    - [3.2.1.5.2 表达式树组合模式](#32152-表达式树组合模式)
  - [3.2.1.6 5. 批判性分析](#3216-5-批判性分析)
    - [3.2.1.6.1 优势](#32161-优势)
    - [3.2.1.6.2 劣势](#32162-劣势)
    - [3.2.1.6.3 行业对比](#32163-行业对比)
    - [3.2.1.6.4 最新趋势](#32164-最新趋势)
  - [3.2.1.7 6. 面试题与考点](#3217-6-面试题与考点)
    - [3.2.1.7.1 基础考点](#32171-基础考点)
    - [3.2.1.7.2 进阶考点](#32172-进阶考点)
  - [3.2.1.8 7. 术语表](#3218-7-术语表)
  - [3.2.1.9 8. 常见陷阱](#3219-8-常见陷阱)
  - [3.2.1.10 9. 相关主题](#32110-9-相关主题)
  - [3.2.1.11 10. 学习路径](#32111-10-学习路径)
    - [3.2.1.11.1 新手路径](#321111-新手路径)
    - [3.2.1.11.2 进阶路径](#321112-进阶路径)
    - [3.2.1.11.3 高阶路径](#321113-高阶路径)
<!-- TOC END -->

## 3.2.1.1 目录

## 3.2.1.2 1. 概述

### 3.2.1.2.1 定义

组合模式将对象组合成树形结构以表示"部分-整体"的层次结构，使得用户对单个对象和组合对象的使用具有一致性。

**形式化定义**:
$$Composite = (Component, Leaf, Composite, Client, TreeStructure)$$

其中：

- $Component$ 是组件接口
- $Leaf$ 是叶子节点
- $Composite$ 是组合节点
- $Client$ 是客户端
- $TreeStructure$ 是树形结构

### 3.2.1.2.2 核心特征

- **统一接口**: 叶子节点和组合节点使用相同接口
- **递归结构**: 支持递归的树形结构
- **透明性**: 客户端无需区分叶子节点和组合节点
- **扩展性**: 易于添加新的组件类型

## 3.2.1.3 2. 理论基础

### 3.2.1.3.1 数学形式化

**定义 2.1** (组合模式): 组合模式是一个六元组 $C = (V, E, L, N, F, T)$

其中：

- $V$ 是节点集合
- $E$ 是边集合
- $L$ 是叶子节点集合，$L \subseteq V$
- $N$ 是组合节点集合，$N \subseteq V$
- $F$ 是操作函数，$F: V \rightarrow Result$
- $T$ 是树形结构约束

**定理 2.1** (递归性): 对于任意组合节点 $n \in N$，$F(n) = \sum_{c \in children(n)} F(c)$

**证明**: 由组合模式的递归性质保证。

### 3.2.1.3.2 范畴论视角

在范畴论中，组合模式可以表示为：

$$Composite : Component \times Component \rightarrow Component$$

其中 $Component$ 是对象范畴，满足幺半群性质。

## 3.2.1.4 3. Go语言实现

### 3.2.1.4.1 基础组合模式

```go
package composite

import (
    "fmt"
    "strings"
)

// Component 组件接口
type Component interface {
    Operation() string
    Add(component Component)
    Remove(component Component)
    GetChild(index int) Component
    GetChildren() []Component
    IsLeaf() bool
}

// Leaf 叶子节点
type Leaf struct {
    name string
}

func NewLeaf(name string) *Leaf {
    return &Leaf{name: name}
}

func (l *Leaf) Operation() string {
    return fmt.Sprintf("Leaf(%s)", l.name)
}

func (l *Leaf) Add(component Component) {
    // 叶子节点不支持添加子节点
    fmt.Printf("Cannot add to leaf: %s\n", l.name)
}

func (l *Leaf) Remove(component Component) {
    // 叶子节点不支持移除子节点
    fmt.Printf("Cannot remove from leaf: %s\n", l.name)
}

func (l *Leaf) GetChild(index int) Component {
    // 叶子节点没有子节点
    return nil
}

func (l *Leaf) GetChildren() []Component {
    // 叶子节点没有子节点
    return []Component{}
}

func (l *Leaf) IsLeaf() bool {
    return true
}

// Composite 组合节点
type Composite struct {
    name     string
    children []Component
}

func NewComposite(name string) *Composite {
    return &Composite{
        name:     name,
        children: make([]Component, 0),
    }
}

func (c *Composite) Operation() string {
    results := []string{fmt.Sprintf("Composite(%s)", c.name)}
    
    for _, child := range c.children {
        results = append(results, "  "+child.Operation())
    }
    
    return strings.Join(results, "\n")
}

func (c *Composite) Add(component Component) {
    c.children = append(c.children, component)
}

func (c *Composite) Remove(component Component) {
    for i, child := range c.children {
        if child == component {
            c.children = append(c.children[:i], c.children[i+1:]...)
            break
        }
    }
}

func (c *Composite) GetChild(index int) Component {
    if index >= 0 && index < len(c.children) {
        return c.children[index]
    }
    return nil
}

func (c *Composite) GetChildren() []Component {
    return c.children
}

func (c *Composite) IsLeaf() bool {
    return false
}

```

### 3.2.1.4.2 文件系统组合模式

```go
package filesystem

import (
    "fmt"
    "strings"
    "time"
)

// FileSystemItem 文件系统项接口
type FileSystemItem interface {
    GetName() string
    GetSize() int64
    GetModifiedTime() time.Time
    GetPath() string
    IsDirectory() bool
    List() []FileSystemItem
    Add(item FileSystemItem) error
    Remove(item FileSystemItem) error
    Search(pattern string) []FileSystemItem
}

// File 文件
type File struct {
    name         string
    size         int64
    modifiedTime time.Time
    path         string
}

func NewFile(name string, size int64, path string) *File {
    return &File{
        name:         name,
        size:         size,
        modifiedTime: time.Now(),
        path:         path,
    }
}

func (f *File) GetName() string {
    return f.name
}

func (f *File) GetSize() int64 {
    return f.size
}

func (f *File) GetModifiedTime() time.Time {
    return f.modifiedTime
}

func (f *File) GetPath() string {
    return f.path
}

func (f *File) IsDirectory() bool {
    return false
}

func (f *File) List() []FileSystemItem {
    return []FileSystemItem{}
}

func (f *File) Add(item FileSystemItem) error {
    return fmt.Errorf("cannot add item to file: %s", f.name)
}

func (f *File) Remove(item FileSystemItem) error {
    return fmt.Errorf("cannot remove item from file: %s", f.name)
}

func (f *File) Search(pattern string) []FileSystemItem {
    if strings.Contains(strings.ToLower(f.name), strings.ToLower(pattern)) {
        return []FileSystemItem{f}
    }
    return []FileSystemItem{}
}

// Directory 目录
type Directory struct {
    name         string
    path         string
    modifiedTime time.Time
    items        []FileSystemItem
}

func NewDirectory(name, path string) *Directory {
    return &Directory{
        name:         name,
        path:         path,
        modifiedTime: time.Now(),
        items:        make([]FileSystemItem, 0),
    }
}

func (d *Directory) GetName() string {
    return d.name
}

func (d *Directory) GetSize() int64 {
    var totalSize int64
    for _, item := range d.items {
        totalSize += item.GetSize()
    }
    return totalSize
}

func (d *Directory) GetModifiedTime() time.Time {
    return d.modifiedTime
}

func (d *Directory) GetPath() string {
    return d.path
}

func (d *Directory) IsDirectory() bool {
    return true
}

func (d *Directory) List() []FileSystemItem {
    return d.items
}

func (d *Directory) Add(item FileSystemItem) error {
    d.items = append(d.items, item)
    d.modifiedTime = time.Now()
    return nil
}

func (d *Directory) Remove(item FileSystemItem) error {
    for i, existingItem := range d.items {
        if existingItem == item {
            d.items = append(d.items[:i], d.items[i+1:]...)
            d.modifiedTime = time.Now()
            return nil
        }
    }
    return fmt.Errorf("item not found: %s", item.GetName())
}

func (d *Directory) Search(pattern string) []FileSystemItem {
    var results []FileSystemItem
    
    // 搜索当前目录
    if strings.Contains(strings.ToLower(d.name), strings.ToLower(pattern)) {
        results = append(results, d)
    }
    
    // 递归搜索子项
    for _, item := range d.items {
        results = append(results, item.Search(pattern)...)
    }
    
    return results
}

// FileSystem 文件系统
type FileSystem struct {
    root Directory
}

func NewFileSystem() *FileSystem {
    return &FileSystem{
        root: *NewDirectory("root", "/"),
    }
}

func (fs *FileSystem) GetRoot() *Directory {
    return &fs.root
}

func (fs *FileSystem) CreateFile(path, name string, size int64) error {
    // 简化的文件创建逻辑
    file := NewFile(name, size, path+"/"+name)
    return fs.root.Add(file)
}

func (fs *FileSystem) CreateDirectory(path, name string) error {
    // 简化的目录创建逻辑
    dir := NewDirectory(name, path+"/"+name)
    return fs.root.Add(dir)
}

func (fs *FileSystem) Search(pattern string) []FileSystemItem {
    return fs.root.Search(pattern)
}

func (fs *FileSystem) GetTotalSize() int64 {
    return fs.root.GetSize()
}

func (fs *FileSystem) ListAll() string {
    return fs.listRecursive(&fs.root, 0)
}

func (fs *FileSystem) listRecursive(item FileSystemItem, depth int) string {
    indent := strings.Repeat("  ", depth)
    result := indent + item.GetName()
    
    if item.IsDirectory() {
        result += "/"
        for _, child := range item.List() {
            result += "\n" + fs.listRecursive(child, depth+1)
        }
    }
    
    return result
}

```

### 3.2.1.4.3 UI组件组合模式

```go
package uicomposite

import (
    "fmt"
    "strings"
)

// UIComponent UI组件接口
type UIComponent interface {
    GetName() string
    GetPosition() (x, y int)
    GetSize() (width, height int)
    Render() string
    Add(component UIComponent) error
    Remove(component UIComponent) error
    GetChildren() []UIComponent
    IsContainer() bool
    HandleEvent(event string) bool
}

// Widget 基础组件
type Widget struct {
    name     string
    x, y     int
    width    int
    height   int
    visible  bool
    enabled  bool
}

func NewWidget(name string, x, y, width, height int) *Widget {
    return &Widget{
        name:    name,
        x:       x,
        y:       y,
        width:   width,
        height:  height,
        visible: true,
        enabled: true,
    }
}

func (w *Widget) GetName() string {
    return w.name
}

func (w *Widget) GetPosition() (x, y int) {
    return w.x, w.y
}

func (w *Widget) GetSize() (width, height int) {
    return w.width, w.height
}

func (w *Widget) Render() string {
    if !w.visible {
        return ""
    }
    return fmt.Sprintf("Widget(%s) at (%d,%d) size %dx%d", w.name, w.x, w.y, w.width, w.height)
}

func (w *Widget) Add(component UIComponent) error {
    return fmt.Errorf("cannot add to widget: %s", w.name)
}

func (w *Widget) Remove(component UIComponent) error {
    return fmt.Errorf("cannot remove from widget: %s", w.name)
}

func (w *Widget) GetChildren() []UIComponent {
    return []UIComponent{}
}

func (w *Widget) IsContainer() bool {
    return false
}

func (w *Widget) HandleEvent(event string) bool {
    if !w.enabled {
        return false
    }
    fmt.Printf("Widget %s handling event: %s\n", w.name, event)
    return true
}

// Button 按钮
type Button struct {
    Widget
    text    string
    onClick func()
}

func NewButton(name, text string, x, y, width, height int) *Button {
    return &Button{
        Widget: *NewWidget(name, x, y, width, height),
        text:   text,
    }
}

func (b *Button) Render() string {
    if !b.visible {
        return ""
    }
    return fmt.Sprintf("Button(%s) '%s' at (%d,%d) size %dx%d", b.name, b.text, b.x, b.y, b.width, b.height)
}

func (b *Button) SetOnClick(handler func()) {
    b.onClick = handler
}

func (b *Button) HandleEvent(event string) bool {
    if event == "click" && b.onClick != nil {
        b.onClick()
        return true
    }
    return b.Widget.HandleEvent(event)
}

// TextBox 文本框
type TextBox struct {
    Widget
    text     string
    maxLength int
}

func NewTextBox(name string, x, y, width, height int) *TextBox {
    return &TextBox{
        Widget: *NewWidget(name, x, y, width, height),
    }
}

func (t *TextBox) Render() string {
    if !t.visible {
        return ""
    }
    return fmt.Sprintf("TextBox(%s) '%s' at (%d,%d) size %dx%d", t.name, t.text, t.x, t.y, t.width, t.height)
}

func (t *TextBox) SetText(text string) {
    if t.maxLength > 0 && len(text) > t.maxLength {
        t.text = text[:t.maxLength]
    } else {
        t.text = text
    }
}

func (t *TextBox) GetText() string {
    return t.text
}

func (t *TextBox) HandleEvent(event string) bool {
    if event == "input" {
        fmt.Printf("TextBox %s received input: %s\n", t.name, t.text)
        return true
    }
    return t.Widget.HandleEvent(event)
}

// Container 容器
type Container struct {
    Widget
    children []UIComponent
    layout   string
}

func NewContainer(name string, x, y, width, height int, layout string) *Container {
    return &Container{
        Widget: *NewWidget(name, x, y, width, height),
        children: make([]UIComponent, 0),
        layout:   layout,
    }
}

func (c *Container) Render() string {
    if !c.visible {
        return ""
    }
    
    result := fmt.Sprintf("Container(%s) at (%d,%d) size %dx%d layout=%s", 
        c.name, c.x, c.y, c.width, c.height, c.layout)
    
    for _, child := range c.children {
        childRender := child.Render()
        if childRender != "" {
            result += "\n  " + childRender
        }
    }
    
    return result
}

func (c *Container) Add(component UIComponent) error {
    c.children = append(c.children, component)
    return nil
}

func (c *Container) Remove(component UIComponent) error {
    for i, child := range c.children {
        if child == component {
            c.children = append(c.children[:i], c.children[i+1:]...)
            return nil
        }
    }
    return fmt.Errorf("component not found: %s", component.GetName())
}

func (c *Container) GetChildren() []UIComponent {
    return c.children
}

func (c *Container) IsContainer() bool {
    return true
}

func (c *Container) HandleEvent(event string) bool {
    // 容器处理事件，然后传递给子组件
    handled := c.Widget.HandleEvent(event)
    
    for _, child := range c.children {
        if child.HandleEvent(event) {
            handled = true
        }
    }
    
    return handled
}

// Window 窗口
type Window struct {
    Container
    title    string
    resizable bool
}

func NewWindow(name, title string, x, y, width, height int) *Window {
    return &Window{
        Container: *NewContainer(name, x, y, width, height, "absolute"),
        title:     title,
        resizable: true,
    }
}

func (w *Window) Render() string {
    if !w.visible {
        return ""
    }
    
    result := fmt.Sprintf("Window(%s) '%s' at (%d,%d) size %dx%d", 
        w.name, w.title, w.x, w.y, w.width, w.height)
    
    for _, child := range w.children {
        childRender := child.Render()
        if childRender != "" {
            result += "\n  " + childRender
        }
    }
    
    return result
}

func (w *Window) SetTitle(title string) {
    w.title = title
}

func (w *Window) GetTitle() string {
    return w.title
}

func (w *Window) HandleEvent(event string) bool {
    if event == "close" {
        fmt.Printf("Window %s closing\n", w.name)
        w.visible = false
        return true
    }
    return w.Container.HandleEvent(event)
}

```

## 3.2.1.5 4. 工程案例

### 3.2.1.5.1 组织架构组合模式

```go
package organization

import (
    "fmt"
    "strings"
    "time"
)

// Employee 员工接口
type Employee interface {
    GetID() string
    GetName() string
    GetPosition() string
    GetSalary() float64
    GetDepartment() string
    GetSubordinates() []Employee
    AddSubordinate(employee Employee) error
    RemoveSubordinate(employee Employee) error
    IsManager() bool
    GetTotalSubordinates() int
    GetTotalSalary() float64
    GetHierarchy() string
}

// IndividualEmployee 个人员工
type IndividualEmployee struct {
    id         string
    name       string
    position   string
    salary     float64
    department string
    hireDate   time.Time
}

func NewIndividualEmployee(id, name, position, department string, salary float64) *IndividualEmployee {
    return &IndividualEmployee{
        id:         id,
        name:       name,
        position:   position,
        salary:     salary,
        department: department,
        hireDate:   time.Now(),
    }
}

func (e *IndividualEmployee) GetID() string {
    return e.id
}

func (e *IndividualEmployee) GetName() string {
    return e.name
}

func (e *IndividualEmployee) GetPosition() string {
    return e.position
}

func (e *IndividualEmployee) GetSalary() float64 {
    return e.salary
}

func (e *IndividualEmployee) GetDepartment() string {
    return e.department
}

func (e *IndividualEmployee) GetSubordinates() []Employee {
    return []Employee{}
}

func (e *IndividualEmployee) AddSubordinate(employee Employee) error {
    return fmt.Errorf("individual employee cannot have subordinates")
}

func (e *IndividualEmployee) RemoveSubordinate(employee Employee) error {
    return fmt.Errorf("individual employee cannot have subordinates")
}

func (e *IndividualEmployee) IsManager() bool {
    return false
}

func (e *IndividualEmployee) GetTotalSubordinates() int {
    return 0
}

func (e *IndividualEmployee) GetTotalSalary() float64 {
    return e.salary
}

func (e *IndividualEmployee) GetHierarchy() string {
    return fmt.Sprintf("%s (%s)", e.name, e.position)
}

// Manager 经理
type Manager struct {
    IndividualEmployee
    subordinates []Employee
}

func NewManager(id, name, position, department string, salary float64) *Manager {
    return &Manager{
        IndividualEmployee: *NewIndividualEmployee(id, name, position, department, salary),
        subordinates:       make([]Employee, 0),
    }
}

func (m *Manager) GetSubordinates() []Employee {
    return m.subordinates
}

func (m *Manager) AddSubordinate(employee Employee) error {
    m.subordinates = append(m.subordinates, employee)
    return nil
}

func (m *Manager) RemoveSubordinate(employee Employee) error {
    for i, sub := range m.subordinates {
        if sub.GetID() == employee.GetID() {
            m.subordinates = append(m.subordinates[:i], m.subordinates[i+1:]...)
            return nil
        }
    }
    return fmt.Errorf("subordinate not found: %s", employee.GetID())
}

func (m *Manager) IsManager() bool {
    return true
}

func (m *Manager) GetTotalSubordinates() int {
    total := len(m.subordinates)
    for _, sub := range m.subordinates {
        total += sub.GetTotalSubordinates()
    }
    return total
}

func (m *Manager) GetTotalSalary() float64 {
    total := m.salary
    for _, sub := range m.subordinates {
        total += sub.GetTotalSalary()
    }
    return total
}

func (m *Manager) GetHierarchy() string {
    result := fmt.Sprintf("%s (%s)", m.name, m.position)
    
    if len(m.subordinates) > 0 {
        result += "\n"
        for i, sub := range m.subordinates {
            indent := "  "
            subHierarchy := strings.ReplaceAll(sub.GetHierarchy(), "\n", "\n"+indent)
            result += indent + subHierarchy
            if i < len(m.subordinates)-1 {
                result += "\n"
            }
        }
    }
    
    return result
}

// Department 部门
type Department struct {
    name      string
    manager   *Manager
    employees []Employee
}

func NewDepartment(name string, manager *Manager) *Department {
    dept := &Department{
        name:      name,
        manager:   manager,
        employees: make([]Employee, 0),
    }
    
    if manager != nil {
        dept.employees = append(dept.employees, manager)
    }
    
    return dept
}

func (d *Department) GetName() string {
    return d.name
}

func (d *Department) GetManager() *Manager {
    return d.manager
}

func (d *Department) AddEmployee(employee Employee) error {
    d.employees = append(d.employees, employee)
    return nil
}

func (d *Department) RemoveEmployee(employee Employee) error {
    for i, emp := range d.employees {
        if emp.GetID() == employee.GetID() {
            d.employees = append(d.employees[:i], d.employees[i+1:]...)
            return nil
        }
    }
    return fmt.Errorf("employee not found: %s", employee.GetID())
}

func (d *Department) GetEmployees() []Employee {
    return d.employees
}

func (d *Department) GetTotalEmployees() int {
    total := len(d.employees)
    for _, emp := range d.employees {
        total += emp.GetTotalSubordinates()
    }
    return total
}

func (d *Department) GetTotalSalary() float64 {
    total := 0.0
    for _, emp := range d.employees {
        total += emp.GetTotalSalary()
    }
    return total
}

func (d *Department) GetHierarchy() string {
    result := fmt.Sprintf("Department: %s", d.name)
    
    if d.manager != nil {
        result += "\n" + d.manager.GetHierarchy()
    }
    
    return result
}

// Organization 组织
type Organization struct {
    name        string
    departments []*Department
}

func NewOrganization(name string) *Organization {
    return &Organization{
        name:        name,
        departments: make([]*Department, 0),
    }
}

func (o *Organization) GetName() string {
    return o.name
}

func (o *Organization) AddDepartment(department *Department) {
    o.departments = append(o.departments, department)
}

func (o *Organization) RemoveDepartment(department *Department) error {
    for i, dept := range o.departments {
        if dept.GetName() == department.GetName() {
            o.departments = append(o.departments[:i], o.departments[i+1:]...)
            return nil
        }
    }
    return fmt.Errorf("department not found: %s", department.GetName())
}

func (o *Organization) GetDepartments() []*Department {
    return o.departments
}

func (o *Organization) GetTotalEmployees() int {
    total := 0
    for _, dept := range o.departments {
        total += dept.GetTotalEmployees()
    }
    return total
}

func (o *Organization) GetTotalSalary() float64 {
    total := 0.0
    for _, dept := range o.departments {
        total += dept.GetTotalSalary()
    }
    return total
}

func (o *Organization) GetHierarchy() string {
    result := fmt.Sprintf("Organization: %s", o.name)
    
    for _, dept := range o.departments {
        deptHierarchy := strings.ReplaceAll(dept.GetHierarchy(), "\n", "\n  ")
        result += "\n  " + deptHierarchy
    }
    
    return result
}

func (o *Organization) FindEmployee(id string) Employee {
    for _, dept := range o.departments {
        for _, emp := range dept.GetEmployees() {
            if emp.GetID() == id {
                return emp
            }
            // 递归搜索子员工
            if emp.IsManager() {
                for _, sub := range emp.GetSubordinates() {
                    if sub.GetID() == id {
                        return sub
                    }
                }
            }
        }
    }
    return nil
}

```

### 3.2.1.5.2 表达式树组合模式

```go
package expression

import (
    "fmt"
    "strconv"
)

// Expression 表达式接口
type Expression interface {
    Evaluate() float64
    GetType() string
    GetValue() string
    GetChildren() []Expression
    AddChild(child Expression) error
    RemoveChild(child Expression) error
    IsLeaf() bool
    ToString() string
}

// NumberExpression 数字表达式
type NumberExpression struct {
    value float64
}

func NewNumberExpression(value float64) *NumberExpression {
    return &NumberExpression{value: value}
}

func (n *NumberExpression) Evaluate() float64 {
    return n.value
}

func (n *NumberExpression) GetType() string {
    return "number"
}

func (n *NumberExpression) GetValue() string {
    return strconv.FormatFloat(n.value, 'f', -1, 64)
}

func (n *NumberExpression) GetChildren() []Expression {
    return []Expression{}
}

func (n *NumberExpression) AddChild(child Expression) error {
    return fmt.Errorf("cannot add child to number expression")
}

func (n *NumberExpression) RemoveChild(child Expression) error {
    return fmt.Errorf("cannot remove child from number expression")
}

func (n *NumberExpression) IsLeaf() bool {
    return true
}

func (n *NumberExpression) ToString() string {
    return n.GetValue()
}

// VariableExpression 变量表达式
type VariableExpression struct {
    name  string
    value float64
}

func NewVariableExpression(name string, value float64) *VariableExpression {
    return &VariableExpression{
        name:  name,
        value: value,
    }
}

func (v *VariableExpression) Evaluate() float64 {
    return v.value
}

func (v *VariableExpression) GetType() string {
    return "variable"
}

func (v *VariableExpression) GetValue() string {
    return v.name
}

func (v *VariableExpression) GetChildren() []Expression {
    return []Expression{}
}

func (v *VariableExpression) AddChild(child Expression) error {
    return fmt.Errorf("cannot add child to variable expression")
}

func (v *VariableExpression) RemoveChild(child Expression) error {
    return fmt.Errorf("cannot remove child from variable expression")
}

func (v *VariableExpression) IsLeaf() bool {
    return true
}

func (v *VariableExpression) ToString() string {
    return v.name
}

func (v *VariableExpression) SetValue(value float64) {
    v.value = value
}

// BinaryExpression 二元表达式
type BinaryExpression struct {
    operator string
    left     Expression
    right    Expression
}

func NewBinaryExpression(operator string, left, right Expression) *BinaryExpression {
    return &BinaryExpression{
        operator: operator,
        left:     left,
        right:    right,
    }
}

func (b *BinaryExpression) Evaluate() float64 {
    leftVal := b.left.Evaluate()
    rightVal := b.right.Evaluate()
    
    switch b.operator {
    case "+":
        return leftVal + rightVal
    case "-":
        return leftVal - rightVal
    case "*":
        return leftVal * rightVal
    case "/":
        if rightVal == 0 {
            panic("division by zero")
        }
        return leftVal / rightVal
    case "^":
        return pow(leftVal, rightVal)
    default:
        panic(fmt.Sprintf("unknown operator: %s", b.operator))
    }
}

func (b *BinaryExpression) GetType() string {
    return "binary"
}

func (b *BinaryExpression) GetValue() string {
    return b.operator
}

func (b *BinaryExpression) GetChildren() []Expression {
    return []Expression{b.left, b.right}
}

func (b *BinaryExpression) AddChild(child Expression) error {
    if b.left == nil {
        b.left = child
        return nil
    }
    if b.right == nil {
        b.right = child
        return nil
    }
    return fmt.Errorf("binary expression already has two children")
}

func (b *BinaryExpression) RemoveChild(child Expression) error {
    if b.left == child {
        b.left = nil
        return nil
    }
    if b.right == child {
        b.right = nil
        return nil
    }
    return fmt.Errorf("child not found")
}

func (b *BinaryExpression) IsLeaf() bool {
    return false
}

func (b *BinaryExpression) ToString() string {
    if b.left == nil || b.right == nil {
        return "incomplete expression"
    }
    
    leftStr := b.left.ToString()
    rightStr := b.right.ToString()
    
    // 添加括号以确保正确的优先级
    if b.operator == "*" || b.operator == "/" {
        if b.left.GetType() == "binary" && (b.left.GetValue() == "+" || b.left.GetValue() == "-") {
            leftStr = "(" + leftStr + ")"
        }
        if b.right.GetType() == "binary" && (b.right.GetValue() == "+" || b.right.GetValue() == "-") {
            rightStr = "(" + rightStr + ")"
        }
    }
    
    return fmt.Sprintf("%s %s %s", leftStr, b.operator, rightStr)
}

// FunctionExpression 函数表达式
type FunctionExpression struct {
    name     string
    argument Expression
}

func NewFunctionExpression(name string, argument Expression) *FunctionExpression {
    return &FunctionExpression{
        name:     name,
        argument: argument,
    }
}

func (f *FunctionExpression) Evaluate() float64 {
    argVal := f.argument.Evaluate()
    
    switch f.name {
    case "sin":
        return sin(argVal)
    case "cos":
        return cos(argVal)
    case "tan":
        return tan(argVal)
    case "log":
        if argVal <= 0 {
            panic("logarithm of non-positive number")
        }
        return log(argVal)
    case "sqrt":
        if argVal < 0 {
            panic("square root of negative number")
        }
        return sqrt(argVal)
    default:
        panic(fmt.Sprintf("unknown function: %s", f.name))
    }
}

func (f *FunctionExpression) GetType() string {
    return "function"
}

func (f *FunctionExpression) GetValue() string {
    return f.name
}

func (f *FunctionExpression) GetChildren() []Expression {
    return []Expression{f.argument}
}

func (f *FunctionExpression) AddChild(child Expression) error {
    if f.argument == nil {
        f.argument = child
        return nil
    }
    return fmt.Errorf("function expression already has an argument")
}

func (f *FunctionExpression) RemoveChild(child Expression) error {
    if f.argument == child {
        f.argument = nil
        return nil
    }
    return fmt.Errorf("child not found")
}

func (f *FunctionExpression) IsLeaf() bool {
    return false
}

func (f *FunctionExpression) ToString() string {
    if f.argument == nil {
        return f.name + "()"
    }
    return fmt.Sprintf("%s(%s)", f.name, f.argument.ToString())
}

// ExpressionTree 表达式树
type ExpressionTree struct {
    root Expression
}

func NewExpressionTree(root Expression) *ExpressionTree {
    return &ExpressionTree{root: root}
}

func (et *ExpressionTree) Evaluate() float64 {
    if et.root == nil {
        return 0
    }
    return et.root.Evaluate()
}

func (et *ExpressionTree) ToString() string {
    if et.root == nil {
        return "empty"
    }
    return et.root.ToString()
}

func (et *ExpressionTree) GetRoot() Expression {
    return et.root
}

func (et *ExpressionTree) SetRoot(root Expression) {
    et.root = root
}

// 辅助数学函数
func pow(x, y float64) float64 {
    result := 1.0
    for i := 0; i < int(y); i++ {
        result *= x
    }
    return result
}

func sin(x float64) float64 {
    // 简化的sin实现
    return x // 实际应用中应使用math.Sin
}

func cos(x float64) float64 {
    // 简化的cos实现
    return 1.0 // 实际应用中应使用math.Cos
}

func tan(x float64) float64 {
    // 简化的tan实现
    return x // 实际应用中应使用math.Tan
}

func log(x float64) float64 {
    // 简化的log实现
    return x // 实际应用中应使用math.Log
}

func sqrt(x float64) float64 {
    // 简化的sqrt实现
    return x // 实际应用中应使用math.Sqrt
}

```

## 3.2.1.6 5. 批判性分析

### 3.2.1.6.1 优势

1. **统一接口**: 叶子节点和组合节点使用相同接口
2. **递归结构**: 支持递归的树形结构
3. **透明性**: 客户端无需区分叶子节点和组合节点
4. **扩展性**: 易于添加新的组件类型

### 3.2.1.6.2 劣势

1. **类型安全**: 可能违反类型安全原则
2. **性能开销**: 递归操作可能影响性能
3. **复杂性**: 树形结构可能变得复杂
4. **内存使用**: 大量节点可能占用大量内存

### 3.2.1.6.3 行业对比

| 语言 | 实现方式 | 性能 | 复杂度 |
|------|----------|------|--------|
| Go | 接口 + 结构体 | 高 | 中 |
| Java | 抽象类 + 接口 | 中 | 中 |
| C++ | 虚函数 | 中 | 中 |
| Python | 抽象基类 | 高 | 低 |

### 3.2.1.6.4 最新趋势

1. **函数式组合**: 使用函数式编程思想
2. **不可变组合**: 不可变数据结构
3. **并行组合**: 并行处理树形结构
4. **流式组合**: 流式处理组合数据

## 3.2.1.7 6. 面试题与考点

### 3.2.1.7.1 基础考点

1. **Q**: 组合模式与装饰器模式的区别？
   **A**: 组合表示部分-整体关系，装饰器增强功能

2. **Q**: 什么时候使用组合模式？
   **A**: 需要表示树形结构、统一处理叶子节点和组合节点时

3. **Q**: 组合模式的优缺点？
   **A**: 优点：统一接口、递归结构；缺点：类型安全、性能开销

### 3.2.1.7.2 进阶考点

1. **Q**: 如何优化组合模式的性能？
   **A**: 缓存、懒加载、并行处理

2. **Q**: 组合模式在大型项目中的应用？
   **A**: 文件系统、UI组件、组织架构

3. **Q**: 如何处理组合模式的类型安全？
   **A**: 泛型、类型断言、接口设计

## 3.2.1.8 7. 术语表

| 术语 | 定义 | 英文 |
|------|------|------|
| 组合模式 | 表示部分-整体层次结构的设计模式 | Composite Pattern |
| 组件 | 抽象接口 | Component |
| 叶子节点 | 没有子节点的组件 | Leaf |
| 组合节点 | 包含子节点的组件 | Composite |
| 树形结构 | 层次化的数据结构 | Tree Structure |

## 3.2.1.9 8. 常见陷阱

| 陷阱 | 问题 | 解决方案 |
|------|------|----------|
| 类型安全 | 违反类型安全原则 | 使用泛型或类型断言 |
| 性能问题 | 递归操作影响性能 | 缓存、懒加载 |
| 内存泄漏 | 循环引用导致内存泄漏 | 弱引用、清理机制 |
| 过度复杂 | 树形结构过于复杂 | 简化设计、分层处理 |

## 3.2.1.10 9. 相关主题

- [适配器模式](./01-Adapter-Pattern.md)
- [装饰器模式](./02-Decorator-Pattern.md)
- [代理模式](./03-Proxy-Pattern.md)
- [外观模式](./04-Facade-Pattern.md)
- [桥接模式](./05-Bridge-Pattern.md)

## 3.2.1.11 10. 学习路径

### 3.2.1.11.1 新手路径

1. 理解组合模式的基本概念
2. 学习树形结构的表示
3. 实现简单的组合模式
4. 理解递归操作

### 3.2.1.11.2 进阶路径

1. 学习复杂的组合实现
2. 理解组合的性能优化
3. 掌握组合的应用场景
4. 学习组合的最佳实践

### 3.2.1.11.3 高阶路径

1. 分析组合在大型项目中的应用
2. 理解组合与架构设计的关系
3. 掌握组合的性能调优
4. 学习组合的替代方案

---

**相关文档**: [结构型模式总览](./README.md) | [设计模式总览](../README.md)

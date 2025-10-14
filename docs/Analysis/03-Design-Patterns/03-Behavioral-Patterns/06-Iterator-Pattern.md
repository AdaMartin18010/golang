# 3.3.1 迭代器模式 (Iterator Pattern)

## 3.3.1.1 目录

## 3.3.1.2 1. 概述

### 3.3.1.2.1 定义

迭代器模式提供一种方法顺序访问一个聚合对象中的各个元素，而又不暴露其内部的表示。

**形式化定义**:
$$Iterator = (Aggregate, Iterator, ConcreteAggregate, ConcreteIterator)$$

其中：

- $Aggregate$ 是聚合接口
- $Iterator$ 是迭代器接口
- $ConcreteAggregate$ 是具体聚合
- $ConcreteIterator$ 是具体迭代器

### 3.3.1.2.2 核心特征

- **顺序访问**: 提供顺序访问聚合元素的方法
- **封装内部**: 不暴露聚合的内部结构
- **多种遍历**: 支持不同的遍历方式
- **统一接口**: 提供统一的迭代接口

## 3.3.1.3 2. 理论基础

### 3.3.1.3.1 数学形式化

**定义 2.1** (迭代器模式): 迭代器模式是一个四元组 $I = (A, It, N, H)$

其中：

- $A$ 是聚合集合
- $It$ 是迭代器集合
- $N$ 是下一个函数，$N: It \rightarrow Element$
- $H$ 是判断函数，$H: It \rightarrow Boolean$

**定理 2.1** (迭代完整性): 对于任意聚合 $a \in A$，其迭代器 $it \in It$ 能够遍历所有元素。

### 3.3.1.3.2 范畴论视角

在范畴论中，迭代器模式可以表示为：

$$Iterator : Aggregate \rightarrow Element$$

## 3.3.1.4 3. Go语言实现

### 3.3.1.4.1 基础迭代器模式

```go
package iterator

import "fmt"

// Iterator 迭代器接口
type Iterator interface {
    HasNext() bool
    Next() interface{}
    Current() interface{}
    Reset()
}

// Aggregate 聚合接口
type Aggregate interface {
    CreateIterator() Iterator
    Add(item interface{})
    Remove(item interface{})
    GetSize() int
}

// ConcreteIterator 具体迭代器
type ConcreteIterator struct {
    aggregate *ConcreteAggregate
    index     int
}

func NewConcreteIterator(aggregate *ConcreteAggregate) *ConcreteIterator {
    return &ConcreteIterator{
        aggregate: aggregate,
        index:     0,
    }
}

func (c *ConcreteIterator) HasNext() bool {
    return c.index < len(c.aggregate.items)
}

func (c *ConcreteIterator) Next() interface{} {
    if c.HasNext() {
        item := c.aggregate.items[c.index]
        c.index++
        return item
    }
    return nil
}

func (c *ConcreteIterator) Current() interface{} {
    if c.index < len(c.aggregate.items) {
        return c.aggregate.items[c.index]
    }
    return nil
}

func (c *ConcreteIterator) Reset() {
    c.index = 0
}

// ConcreteAggregate 具体聚合
type ConcreteAggregate struct {
    items []interface{}
}

func NewConcreteAggregate() *ConcreteAggregate {
    return &ConcreteAggregate{
        items: make([]interface{}, 0),
    }
}

func (c *ConcreteAggregate) CreateIterator() Iterator {
    return NewConcreteIterator(c)
}

func (c *ConcreteAggregate) Add(item interface{}) {
    c.items = append(c.items, item)
}

func (c *ConcreteAggregate) Remove(item interface{}) {
    for i, existingItem := range c.items {
        if existingItem == item {
            c.items = append(c.items[:i], c.items[i+1:]...)
            break
        }
    }
}

func (c *ConcreteAggregate) GetSize() int {
    return len(c.items)
}

func (c *ConcreteAggregate) GetItem(index int) interface{} {
    if index >= 0 && index < len(c.items) {
        return c.items[index]
    }
    return nil
}

```

### 3.3.1.4.2 图书管理迭代器

```go
package bookiterator

import "fmt"

// Book 图书
type Book struct {
    ID     string
    Title  string
    Author string
    Year   int
}

func NewBook(id, title, author string, year int) *Book {
    return &Book{
        ID:     id,
        Title:  title,
        Author: author,
        Year:   year,
    }
}

func (b *Book) String() string {
    return fmt.Sprintf("Book{ID: %s, Title: %s, Author: %s, Year: %d}", 
        b.ID, b.Title, b.Author, b.Year)
}

// BookIterator 图书迭代器接口
type BookIterator interface {
    HasNext() bool
    Next() *Book
    Current() *Book
    Reset()
}

// BookCollection 图书集合
type BookCollection struct {
    books []*Book
}

func NewBookCollection() *BookCollection {
    return &BookCollection{
        books: make([]*Book, 0),
    }
}

func (b *BookCollection) AddBook(book *Book) {
    b.books = append(b.books, book)
}

func (b *BookCollection) RemoveBook(book *Book) {
    for i, existingBook := range b.books {
        if existingBook.ID == book.ID {
            b.books = append(b.books[:i], b.books[i+1:]...)
            break
        }
    }
}

func (b *BookCollection) GetSize() int {
    return len(b.books)
}

func (b *BookCollection) CreateIterator() BookIterator {
    return NewBookIterator(b)
}

func (b *BookCollection) CreateReverseIterator() BookIterator {
    return NewReverseBookIterator(b)
}

func (b *BookCollection) CreateAuthorIterator(author string) BookIterator {
    return NewAuthorBookIterator(b, author)
}

// BookIterator 正向图书迭代器
type BookIterator struct {
    collection *BookCollection
    index      int
}

func NewBookIterator(collection *BookCollection) *BookIterator {
    return &BookIterator{
        collection: collection,
        index:      0,
    }
}

func (b *BookIterator) HasNext() bool {
    return b.index < len(b.collection.books)
}

func (b *BookIterator) Next() *Book {
    if b.HasNext() {
        book := b.collection.books[b.index]
        b.index++
        return book
    }
    return nil
}

func (b *BookIterator) Current() *Book {
    if b.index < len(b.collection.books) {
        return b.collection.books[b.index]
    }
    return nil
}

func (b *BookIterator) Reset() {
    b.index = 0
}

// ReverseBookIterator 反向图书迭代器
type ReverseBookIterator struct {
    collection *BookCollection
    index      int
}

func NewReverseBookIterator(collection *BookCollection) *ReverseBookIterator {
    return &ReverseBookIterator{
        collection: collection,
        index:      len(collection.books) - 1,
    }
}

func (r *ReverseBookIterator) HasNext() bool {
    return r.index >= 0
}

func (r *ReverseBookIterator) Next() *Book {
    if r.HasNext() {
        book := r.collection.books[r.index]
        r.index--
        return book
    }
    return nil
}

func (r *ReverseBookIterator) Current() *Book {
    if r.index >= 0 && r.index < len(r.collection.books) {
        return r.collection.books[r.index]
    }
    return nil
}

func (r *ReverseBookIterator) Reset() {
    r.index = len(r.collection.books) - 1
}

// AuthorBookIterator 按作者过滤的图书迭代器
type AuthorBookIterator struct {
    collection *BookCollection
    author     string
    index      int
}

func NewAuthorBookIterator(collection *BookCollection, author string) *AuthorBookIterator {
    return &AuthorBookIterator{
        collection: collection,
        author:     author,
        index:      0,
    }
}

func (a *AuthorBookIterator) HasNext() bool {
    for i := a.index; i < len(a.collection.books); i++ {
        if a.collection.books[i].Author == a.author {
            return true
        }
    }
    return false
}

func (a *AuthorBookIterator) Next() *Book {
    for a.index < len(a.collection.books) {
        book := a.collection.books[a.index]
        a.index++
        if book.Author == a.author {
            return book
        }
    }
    return nil
}

func (a *AuthorBookIterator) Current() *Book {
    if a.index > 0 && a.index <= len(a.collection.books) {
        book := a.collection.books[a.index-1]
        if book.Author == a.author {
            return book
        }
    }
    return nil
}

func (a *AuthorBookIterator) Reset() {
    a.index = 0
}

```

### 3.3.1.4.3 文件系统迭代器

```go
package filesystemiterator

import (
    "fmt"
    "os"
    "path/filepath"
)

// FileNode 文件节点
type FileNode struct {
    Name     string
    Path     string
    IsDir    bool
    Size     int64
    ModTime  string
}

func NewFileNode(name, path string, isDir bool, size int64, modTime string) *FileNode {
    return &FileNode{
        Name:    name,
        Path:    path,
        IsDir:   isDir,
        Size:    size,
        ModTime: modTime,
    }
}

func (f *FileNode) String() string {
    nodeType := "File"
    if f.IsDir {
        nodeType = "Directory"
    }
    return fmt.Sprintf("%s{Name: %s, Path: %s, Size: %d, ModTime: %s}", 
        nodeType, f.Name, f.Path, f.Size, f.ModTime)
}

// FileSystemIterator 文件系统迭代器接口
type FileSystemIterator interface {
    HasNext() bool
    Next() *FileNode
    Current() *FileNode
    Reset()
}

// FileSystem 文件系统
type FileSystem struct {
    rootPath string
}

func NewFileSystem(rootPath string) *FileSystem {
    return &FileSystem{
        rootPath: rootPath,
    }
}

func (f *FileSystem) CreateIterator() FileSystemIterator {
    return NewFileSystemIterator(f)
}

func (f *FileSystem) CreateDirectoryIterator() FileSystemIterator {
    return NewDirectoryIterator(f)
}

func (f *FileSystem) CreateFileIterator() FileSystemIterator {
    return NewFileIterator(f)
}

func (f *FileSystem) CreateRecursiveIterator() FileSystemIterator {
    return NewRecursiveFileSystemIterator(f)
}

// FileSystemIterator 基础文件系统迭代器
type FileSystemIterator struct {
    fileSystem *FileSystem
    files      []*FileNode
    index      int
}

func NewFileSystemIterator(fileSystem *FileSystem) *FileSystemIterator {
    iterator := &FileSystemIterator{
        fileSystem: fileSystem,
        index:      0,
    }
    iterator.loadFiles()
    return iterator
}

func (f *FileSystemIterator) loadFiles() {
    f.files = make([]*FileNode, 0)
    
    entries, err := os.ReadDir(f.fileSystem.rootPath)
    if err != nil {
        return
    }
    
    for _, entry := range entries {
        info, err := entry.Info()
        if err != nil {
            continue
        }
        
        fileNode := NewFileNode(
            entry.Name(),
            filepath.Join(f.fileSystem.rootPath, entry.Name()),
            entry.IsDir(),
            info.Size(),
            info.ModTime().Format("2006-01-02 15:04:05"),
        )
        f.files = append(f.files, fileNode)
    }
}

func (f *FileSystemIterator) HasNext() bool {
    return f.index < len(f.files)
}

func (f *FileSystemIterator) Next() *FileNode {
    if f.HasNext() {
        file := f.files[f.index]
        f.index++
        return file
    }
    return nil
}

func (f *FileSystemIterator) Current() *FileNode {
    if f.index > 0 && f.index <= len(f.files) {
        return f.files[f.index-1]
    }
    return nil
}

func (f *FileSystemIterator) Reset() {
    f.index = 0
}

// DirectoryIterator 目录迭代器
type DirectoryIterator struct {
    FileSystemIterator
}

func NewDirectoryIterator(fileSystem *FileSystem) *DirectoryIterator {
    iterator := &DirectoryIterator{
        FileSystemIterator: *NewFileSystemIterator(fileSystem),
    }
    iterator.filterDirectories()
    return iterator
}

func (d *DirectoryIterator) filterDirectories() {
    filtered := make([]*FileNode, 0)
    for _, file := range d.files {
        if file.IsDir {
            filtered = append(filtered, file)
        }
    }
    d.files = filtered
}

// FileIterator 文件迭代器
type FileIterator struct {
    FileSystemIterator
}

func NewFileIterator(fileSystem *FileSystem) *FileIterator {
    iterator := &FileIterator{
        FileSystemIterator: *NewFileSystemIterator(fileSystem),
    }
    iterator.filterFiles()
    return iterator
}

func (f *FileIterator) filterFiles() {
    filtered := make([]*FileNode, 0)
    for _, file := range f.files {
        if !file.IsDir {
            filtered = append(filtered, file)
        }
    }
    f.files = filtered
}

// RecursiveFileSystemIterator 递归文件系统迭代器
type RecursiveFileSystemIterator struct {
    fileSystem *FileSystem
    files      []*FileNode
    index      int
}

func NewRecursiveFileSystemIterator(fileSystem *FileSystem) *RecursiveFileSystemIterator {
    iterator := &RecursiveFileSystemIterator{
        fileSystem: fileSystem,
        index:      0,
    }
    iterator.loadFilesRecursively()
    return iterator
}

func (r *RecursiveFileSystemIterator) loadFilesRecursively() {
    r.files = make([]*FileNode, 0)
    
    err := filepath.Walk(r.fileSystem.rootPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        fileNode := NewFileNode(
            info.Name(),
            path,
            info.IsDir(),
            info.Size(),
            info.ModTime().Format("2006-01-02 15:04:05"),
        )
        r.files = append(r.files, fileNode)
        return nil
    })
    
    if err != nil {
        fmt.Printf("Error walking directory: %v\n", err)
    }
}

func (r *RecursiveFileSystemIterator) HasNext() bool {
    return r.index < len(r.files)
}

func (r *RecursiveFileSystemIterator) Next() *FileNode {
    if r.HasNext() {
        file := r.files[r.index]
        r.index++
        return file
    }
    return nil
}

func (r *RecursiveFileSystemIterator) Current() *FileNode {
    if r.index > 0 && r.index <= len(r.files) {
        return r.files[r.index-1]
    }
    return nil
}

func (r *RecursiveFileSystemIterator) Reset() {
    r.index = 0
}

```

## 3.3.1.5 4. 工程案例

### 3.3.1.5.1 数据库结果集迭代器

```go
package databaseiterator

import (
    "database/sql"
    "fmt"
)

// RowIterator 行迭代器接口
type RowIterator interface {
    HasNext() bool
    Next() map[string]interface{}
    Current() map[string]interface{}
    Reset()
    Close()
}

// DatabaseIterator 数据库迭代器
type DatabaseIterator struct {
    rows *sql.Rows
    cols []string
    current map[string]interface{}
    hasNext bool
}

func NewDatabaseIterator(rows *sql.Rows) (*DatabaseIterator, error) {
    cols, err := rows.Columns()
    if err != nil {
        return nil, err
    }
    
    return &DatabaseIterator{
        rows: rows,
        cols: cols,
    }, nil
}

func (d *DatabaseIterator) HasNext() bool {
    if d.hasNext {
        return true
    }
    
    d.hasNext = d.rows.Next()
    if d.hasNext {
        d.scanCurrentRow()
    }
    return d.hasNext
}

func (d *DatabaseIterator) Next() map[string]interface{} {
    if !d.hasNext {
        return nil
    }
    
    result := d.current
    d.hasNext = false
    return result
}

func (d *DatabaseIterator) Current() map[string]interface{} {
    return d.current
}

func (d *DatabaseIterator) Reset() {
    // 数据库迭代器通常不支持重置
    fmt.Println("Database iterator does not support reset")
}

func (d *DatabaseIterator) Close() {
    if d.rows != nil {
        d.rows.Close()
    }
}

func (d *DatabaseIterator) scanCurrentRow() {
    // 创建值的切片
    values := make([]interface{}, len(d.cols))
    valuePtrs := make([]interface{}, len(d.cols))
    
    for i := range values {
        valuePtrs[i] = &values[i]
    }
    
    // 扫描当前行
    err := d.rows.Scan(valuePtrs...)
    if err != nil {
        fmt.Printf("Error scanning row: %v\n", err)
        return
    }
    
    // 构建结果映射
    d.current = make(map[string]interface{})
    for i, col := range d.cols {
        d.current[col] = values[i]
    }
}

// QueryExecutor 查询执行器
type QueryExecutor struct {
    db *sql.DB
}

func NewQueryExecutor(db *sql.DB) *QueryExecutor {
    return &QueryExecutor{
        db: db,
    }
}

func (q *QueryExecutor) ExecuteQuery(query string, args ...interface{}) (RowIterator, error) {
    rows, err := q.db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    
    return NewDatabaseIterator(rows)
}

```

## 3.3.1.6 5. 批判性分析

### 3.3.1.6.1 优势

1. **封装内部**: 不暴露聚合的内部结构
2. **统一接口**: 提供统一的迭代接口
3. **多种遍历**: 支持不同的遍历方式
4. **并发安全**: 支持并发迭代

### 3.3.1.6.2 劣势

1. **性能开销**: 迭代器对象创建开销
2. **内存使用**: 可能占用额外内存
3. **复杂性**: 增加系统复杂度
4. **调试困难**: 迭代器状态难以调试

### 3.3.1.6.3 行业对比

| 语言 | 实现方式 | 性能 | 复杂度 |
|------|----------|------|--------|
| Go | 接口 | 高 | 中 |
| Java | Iterator接口 | 中 | 中 |
| C++ | 迭代器类 | 高 | 中 |
| Python | 生成器 | 高 | 低 |

### 3.3.1.6.4 最新趋势

1. **函数式迭代**: 使用函数式编程
2. **流式处理**: 流式迭代器
3. **异步迭代**: 异步迭代器
4. **内存优化**: 内存高效的迭代器

## 3.3.1.7 6. 面试题与考点

### 3.3.1.7.1 基础考点

1. **Q**: 迭代器模式与访问者模式的区别？
   **A**: 迭代器关注遍历，访问者关注操作

2. **Q**: 什么时候使用迭代器模式？
   **A**: 需要遍历复杂数据结构、隐藏内部实现时

3. **Q**: 迭代器模式的优缺点？
   **A**: 优点：封装内部、统一接口；缺点：性能开销、复杂度增加

### 3.3.1.7.2 进阶考点

1. **Q**: 如何实现并发安全的迭代器？
   **A**: 使用锁、快照、不可变集合

2. **Q**: 迭代器模式在微服务中的应用？
   **A**: 分页查询、流式处理、事件迭代

3. **Q**: 如何处理迭代器的性能问题？
   **A**: 延迟加载、批量处理、内存优化

## 3.3.1.8 7. 术语表

| 术语 | 定义 | 英文 |
|------|------|------|
| 迭代器模式 | 提供顺序访问的设计模式 | Iterator Pattern |
| 聚合 | 包含元素的集合 | Aggregate |
| 迭代器 | 遍历聚合的对象 | Iterator |
| 遍历 | 访问聚合中的元素 | Traversal |

## 3.3.1.9 8. 常见陷阱

| 陷阱 | 问题 | 解决方案 |
|------|------|----------|
| 并发修改 | 迭代时修改集合 | 使用快照、不可变集合 |
| 内存泄漏 | 迭代器未正确关闭 | 使用defer、资源管理 |
| 性能问题 | 大量迭代器对象 | 对象池、延迟创建 |
| 状态混乱 | 迭代器状态不一致 | 状态验证、重置机制 |

## 3.3.1.10 9. 相关主题

- [观察者模式](./01-Observer-Pattern.md)
- [策略模式](./02-Strategy-Pattern.md)
- [命令模式](./03-Command-Pattern.md)
- [状态模式](./04-State-Pattern.md)
- [模板方法模式](./05-Template-Method-Pattern.md)

## 3.3.1.11 10. 学习路径

### 3.3.1.11.1 新手路径

1. 理解迭代器模式的基本概念
2. 学习聚合和迭代器的关系
3. 实现简单的迭代器模式
4. 理解遍历机制

### 3.3.1.11.2 进阶路径

1. 学习复杂的迭代器实现
2. 理解迭代器的性能优化
3. 掌握迭代器的应用场景
4. 学习迭代器的最佳实践

### 3.3.1.11.3 高阶路径

1. 分析迭代器在大型项目中的应用
2. 理解迭代器与架构设计的关系
3. 掌握迭代器的性能调优
4. 学习迭代器的替代方案

---

**相关文档**: [行为型模式总览](./README.md) | [设计模式总览](../README.md)

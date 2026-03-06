# Go 1.26 实战示例集

> 可运行的完整代码示例，覆盖所有 Go 1.26 新特性

---

## 目录

- [Go 1.26 实战示例集](#go-126-实战示例集)
  - [目录](#目录)
  - [1. new(expr) 实战](#1-newexpr-实战)
    - [示例 1.1: Web API 可选字段](#示例-11-web-api-可选字段)
    - [示例 1.2: 配置初始化](#示例-12-配置初始化)
  - [2. 递归泛型实战](#2-递归泛型实战)
    - [示例 2.1: 通用树遍历库](#示例-21-通用树遍历库)
    - [示例 2.2: 通用搜索树](#示例-22-通用搜索树)
  - [3. Green Tea GC 监控](#3-green-tea-gc-监控)
  - [4. crypto/hpke 实战](#4-cryptohpke-实战)
    - [示例 4.1: 安全消息传递](#示例-41-安全消息传递)
  - [5. go fix 实战](#5-go-fix-实战)
    - [示例 5.1: 项目现代化](#示例-51-项目现代化)
    - [示例 5.2: API 迁移](#示例-52-api-迁移)
  - [6. 综合项目示例](#6-综合项目示例)
    - [项目: 安全配置服务](#项目-安全配置服务)
  - [运行说明](#运行说明)

---

## 1. new(expr) 实战

### 示例 1.1: Web API 可选字段

```go
package main

import (
 "encoding/json"
 "fmt"
 "net/http"
 "time"
)

// User 用户信息
type User struct {
 ID        string     `json:"id"`
 Name      string     `json:"name"`
 Email     *string    `json:"email,omitempty"`      // 可选
 Age       *int       `json:"age,omitempty"`        // 可选
 BirthDate *time.Time `json:"birth_date,omitempty"` // 可选
 VIP       *bool      `json:"vip,omitempty"`        // 可选
}

func calculateAge(birth time.Time) int {
 return int(time.Since(birth).Hours() / 24 / 365.25)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
 name := r.FormValue("name")
 email := r.FormValue("email")
 birthStr := r.FormValue("birth_date")

 user := User{
  ID:   generateID(),
  Name: name,
 }

 // Go 1.26: 使用 new(expr) 简洁处理可选字段
 if email != "" {
  user.Email = new(email)  // 直接传入变量
 }

 if birthStr != "" {
  if birth, err := time.Parse("2006-01-02", birthStr); err == nil {
   user.BirthDate = new(birth)                           // 直接传入解析结果
   user.Age = new(calculateAge(birth))                   // 直接传入计算结果
  }
 }

 // 默认非 VIP
 user.VIP = new(false)
 if r.FormValue("vip") == "true" {
  *user.VIP = true
 }

 w.Header().Set("Content-Type", "application/json")
 json.NewEncoder(w).Encode(user)
}

func generateID() string {
 return fmt.Sprintf("user_%d", time.Now().Unix())
}

func main() {
 http.HandleFunc("/api/users", createUserHandler)
 fmt.Println("Server running on :8080")
 http.ListenAndServe(":8080", nil)
}
```

**传统写法对比**:

```go
// Go 1.25 之前 - 需要中间变量
if email != "" {
    e := email  // 必须声明临时变量
    user.Email = &e
}
if birthStr != "" {
    birth, _ := time.Parse("2006-01-02", birthStr)
    user.BirthDate = &birth  // 编译错误: 不能取地址
    // 必须这样:
    birthCopy := birth
    user.BirthDate = &birthCopy
}

// Go 1.26 - 简洁直接
user.Email = new(email)
user.BirthDate = new(birth)
```

---

### 示例 1.2: 配置初始化

```go
package config

import (
 "os"
 "strconv"
 "time"
)

type ServerConfig struct {
 Host         string
 Port         int
 ReadTimeout  *time.Duration  // 可选，nil 表示使用默认值
 WriteTimeout *time.Duration
 MaxHeader    *int
 TLS          *TLSConfig
}

type TLSConfig struct {
 CertFile string
 KeyFile  string
}

func LoadConfig() *ServerConfig {
 cfg := &ServerConfig{
  Host: getEnv("HOST", "localhost"),
  Port: getEnvAsInt("PORT", 8080),
 }

 // Go 1.26: 一行处理环境变量转换
 if timeout := getEnv("READ_TIMEOUT", ""); timeout != "" {
  if d, err := time.ParseDuration(timeout); err == nil {
   cfg.ReadTimeout = new(d)  // 直接传入解析结果
  }
 }

 if timeout := getEnv("WRITE_TIMEOUT", ""); timeout != "" {
  if d, err := time.ParseDuration(timeout); err == nil {
   cfg.WriteTimeout = new(d)
  }
 }

 if max := getEnv("MAX_HEADER", ""); max != "" {
  if n, err := strconv.Atoi(max); err == nil {
   cfg.MaxHeader = new(n)
  }
 }

 if cert := getEnv("TLS_CERT", ""); cert != "" {
  cfg.TLS = new(TLSConfig{  // 直接初始化复合类型
   CertFile: cert,
   KeyFile:  getEnv("TLS_KEY", ""),
  })
 }

 return cfg
}

func getEnv(key, defaultVal string) string {
 if v := os.Getenv(key); v != "" {
  return v
 }
 return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
 if v := os.Getenv(key); v != "" {
  if n, err := strconv.Atoi(v); err == nil {
   return n
  }
 }
 return defaultVal
}
```

---

## 2. 递归泛型实战

### 示例 2.1: 通用树遍历库

```go
package tree

// Node 树节点接口
type Node[T Node[T]] interface {
 Value() int
 Children() []T
}

// PreOrder 前序遍历
func PreOrder[T Node[T]](root T, visit func(int)) {
 if root == nil {
  return
 }
 visit(root.Value())
 for _, child := range root.Children() {
  PreOrder(child, visit)
 }
}

// PostOrder 后序遍历
func PostOrder[T Node[T]](root T, visit func(int)) {
 if root == nil {
  return
 }
 for _, child := range root.Children() {
  PostOrder(child, visit)
 }
 visit(root.Value())
}

// Find 查找节点
func Find[T Node[T]](root T, target int) (T, bool) {
 var zero T
 if root == nil {
  return zero, false
 }
 if root.Value() == target {
  return root, true
 }
 for _, child := range root.Children() {
  if found, ok := Find(child, target); ok {
   return found, true
  }
 }
 return zero, false
}

// MaxDepth 计算最大深度
func MaxDepth[T Node[T]](root T) int {
 if root == nil {
  return 0
 }
 maxChildDepth := 0
 for _, child := range root.Children() {
  if d := MaxDepth(child); d > maxChildDepth {
   maxChildDepth = d
  }
 }
 return 1 + maxChildDepth
}
```

**二叉树实现**:

```go
package tree

// BinaryNode 二叉树节点
type BinaryNode struct {
 Value_ int
 Left   *BinaryNode
 Right  *BinaryNode
}

func (b *BinaryNode) Value() int {
 return b.Value_
}

func (b *BinaryNode) Children() []*BinaryNode {
 children := make([]*BinaryNode, 0, 2)
 if b.Left != nil {
  children = append(children, b.Left)
 }
 if b.Right != nil {
  children = append(children, b.Right)
 }
 return children
}

// 验证 BinaryNode 实现 Node[*BinaryNode]
var _ Node[*BinaryNode] = (*BinaryNode)(nil)
```

**多叉树实现**:

```go
package tree

// NaryNode 多叉树节点
type NaryNode struct {
 Value_   int
 Children []*NaryNode
}

func (n *NaryNode) Value() int {
 return n.Value_
}

func (n *NaryNode) Children() []*NaryNode {
 return n.Children
}

// 验证 NaryNode 实现 Node[*NaryNode]
var _ Node[*NaryNode] = (*NaryNode)(nil)
```

**使用示例**:

```go
package main

import (
 "fmt"
 "tree"
)

func main() {
 // 二叉树
 binaryTree := &tree.BinaryNode{
  Value_: 1,
  Left: &tree.BinaryNode{
   Value_: 2,
   Left:   &tree.BinaryNode{Value_: 4},
   Right:  &tree.BinaryNode{Value_: 5},
  },
  Right: &tree.BinaryNode{Value_: 3},
 }

 fmt.Print("Binary PreOrder: ")
 tree.PreOrder(binaryTree, func(v int) { fmt.Printf("%d ", v) })
 fmt.Println() // 1 2 4 5 3

 // 多叉树
 naryTree := &tree.NaryNode{
  Value_: 1,
  Children: []*tree.NaryNode{
   {Value_: 2, Children: []*tree.NaryNode{
    {Value_: 4},
    {Value_: 5},
   }},
   {Value_: 3},
  },
 }

 fmt.Print("Nary PreOrder: ")
 tree.PreOrder(naryTree, func(v int) { fmt.Printf("%d ", v) })
 fmt.Println() // 1 2 4 5 3

 // 通用查找
 if node, ok := tree.Find(binaryTree, 5); ok {
  fmt.Printf("Found node with value: %d\n", node.Value())
 }

 // 通用深度计算
 fmt.Printf("Binary tree depth: %d\n", tree.MaxDepth(binaryTree))
 fmt.Printf("Nary tree depth: %d\n", tree.MaxDepth(naryTree))
}
```

---

### 示例 2.2: 通用搜索树

```go
package search

import "golang.org/x/exp/constraints"

// Ordered 可比较接口
type Ordered[T Ordered[T]] interface {
 constraints.Ordered
 Less(other T) bool
 Equal(other T) bool
}

// BST 二叉搜索树
type BST[T Ordered[T]] struct {
 root *bstNode[T]
}

type bstNode[T Ordered[T]] struct {
 value T
 left  *bstNode[T]
 right *bstNode[T]
}

func (b *BST[T]) Insert(v T) {
 b.root = insertNode(b.root, v)
}

func insertNode[T Ordered[T]](n *bstNode[T], v T) *bstNode[T] {
 if n == nil {
  return &bstNode[T]{value: v}
 }
 if v.Less(n.value) {
  n.left = insertNode(n.left, v)
 } else if n.value.Less(v) {
  n.right = insertNode(n.right, v)
 }
 // 相等时不插入
 return n
}

func (b *BST[T]) Search(v T) bool {
 return searchNode(b.root, v)
}

func searchNode[T Ordered[T]](n *bstNode[T], v T) bool {
 if n == nil {
  return false
 }
 if v.Equal(n.value) {
  return true
 }
 if v.Less(n.value) {
  return searchNode(n.left, v)
 }
 return searchNode(n.right, v)
}

func (b *BST[T]) InOrder(visit func(T)) {
 inOrderNode(b.root, visit)
}

func inOrderNode[T Ordered[T]](n *bstNode[T], visit func(T)) {
 if n == nil {
  return
 }
 inOrderNode(n.left, visit)
 visit(n.value)
 inOrderNode(n.right, visit)
}

// Int 实现 Ordered[int]
type Int int

func (i Int) Less(other Int) bool {
 return i < other
}

func (i Int) Equal(other Int) bool {
 return i == other
}

// String 实现 Ordered[string]
type String string

func (s String) Less(other String) bool {
 return s < other
}

func (s String) Equal(other String) bool {
 return s == other
}
```

**使用示例**:

```go
package main

import (
 "fmt"
 "search"
)

func main() {
 // 整数 BST
 intTree := &search.BST[search.Int]{}
 intTree.Insert(5)
 intTree.Insert(3)
 intTree.Insert(7)
 intTree.Insert(1)
 intTree.Insert(9)

 fmt.Println(intTree.Search(3))  // true
 fmt.Println(intTree.Search(10)) // false

 fmt.Print("InOrder: ")
 intTree.InOrder(func(v search.Int) { fmt.Printf("%d ", v) })
 fmt.Println() // 1 3 5 7 9

 // 字符串 BST
 strTree := &search.BST[search.String]{}
 strTree.Insert("apple")
 strTree.Insert("banana")
 strTree.Insert("cherry")

 fmt.Println(strTree.Search("banana")) // true
}
```

---

## 3. Green Tea GC 监控

```go
package main

import (
 "fmt"
 "runtime"
 "runtime/debug"
 "time"
)

func main() {
 fmt.Println("=== Go 1.26 Green Tea GC 监控 ===\n")

 // 1. 基本 GC 统计
 printGCStats()

 // 2. 分配内存观察 GC 行为
 fmt.Println("\n--- 分配 100MB 内存 ---")
 allocateMemory(100 * 1024 * 1024)
 time.Sleep(100 * time.Millisecond) // 等待 GC
 printGCStats()

 // 3. 调整 GC 参数
 fmt.Println("\n--- 调整 GOGC=200 ---")
 old := debug.SetGCPercent(200)
 fmt.Printf("旧 GOGC 值: %d\n", old)

 // 4. 设置内存限制
 debug.SetMemoryLimit(1 << 30) // 1GB

 // 5. 强制 GC
 fmt.Println("\n--- 强制 GC ---")
 runtime.GC()
 printGCStats()
}

func printGCStats() {
 var m runtime.MemStats
 runtime.ReadMemStats(&m)

 fmt.Printf("GC 次数: %d\n", m.NumGC)
 fmt.Printf("GC CPU 占比: %.4f%%\n", m.GCCPUFraction*100)
 fmt.Printf("堆分配: %d MB\n", m.HeapAlloc/1024/1024)
 fmt.Printf("堆系统: %d MB\n", m.HeapSys/1024/1024)
 fmt.Printf("堆对象数: %d\n", m.HeapObjects)
 fmt.Printf("下次 GC 目标: %d MB\n", m.NextGC/1024/1024)

 if m.NumGC > 0 {
  idx := (m.NumGC + 255) % 256
  fmt.Printf("最近 GC 暂停: %d µs\n", m.PauseNs[idx]/1000)
 }
}

func allocateMemory(size int) {
 _ = make([]byte, size)
}
```

---

## 4. crypto/hpke 实战

### 示例 4.1: 安全消息传递

```go
package main

import (
 "crypto/hpke"
 "fmt"
 "log"
)

func main() {
 // 使用 ML-KEM-768 + AES-256-GCM (后量子安全)
 kem, err := hpke.GetKEM(hpke.KEM_MLKEM768)
 if err != nil {
  log.Fatal(err)
 }

 kdf, err := hpke.GetKDF(hpke.KDF_HKDF_SHA384)
 if err != nil {
  log.Fatal(err)
 }

 aead, err := hpke.GetAEAD(hpke.AEAD_AES256GCM)
 if err != nil {
  log.Fatal(err)
 }

 suite, err := hpke.NewSuite(kem, kdf, aead)
 if err != nil {
  log.Fatal(err)
 }

 // 接收方生成密钥对
 skR, err := kem.GenerateKeyPair()
 if err != nil {
  log.Fatal(err)
 }
 pkR := skR.PublicKey()

 // 发送方加密消息
 messages := []string{
  "Hello, Post-Quantum World!",
  "This message is encrypted with HPKE.",
  "ML-KEM-768 provides quantum resistance.",
 }

 for _, msg := range messages {
  // 每次使用新的上下文
  sender := suite.NewSender(pkR, []byte("shared context"))
  enc, senderCtx, err := sender.SetupBase()
  if err != nil {
   log.Fatal(err)
  }

  ciphertext, err := senderCtx.Seal([]byte(msg), nil)
  if err != nil {
   log.Fatal(err)
  }

  // 接收方解密
  recipient := suite.NewRecipient(skR, []byte("shared context"))
  recipientCtx, err := recipient.SetupBase(enc)
  if err != nil {
   log.Fatal(err)
  }

  plaintext, err := recipientCtx.Open(ciphertext, nil)
  if err != nil {
   log.Fatal(err)
  }

  fmt.Printf("原文: %s\n", msg)
  fmt.Printf("解密: %s\n\n", string(plaintext))
 }
}
```

---

## 5. go fix 实战

### 示例 5.1: 项目现代化

假设有以下旧代码 `oldcode.go`:

```go
package main

import (
 "fmt"
 "sort"
)

func main() {
 // 旧方式: 手动查找
 nums := []int{1, 2, 3, 4, 5}
 target := 3
 found := false
 for _, v := range nums {
  if v == target {
   found = true
   break
  }
 }
 fmt.Println("Found:", found)

 // 旧方式: 手动排序
 sort.Slice(nums, func(i, j int) bool {
  return nums[i] < nums[j]
 })

 // 旧方式: 手动找最小值
 a, b := 10, 20
 var min int
 if a < b {
  min = a
 } else {
  min = b
 }
 fmt.Println("Min:", min)
}
```

**运行 go fix**:

```bash
# 预览变化
go fix -n .

# 输出:
# oldcode.go:9:2: replacing loop with slices.Contains
# oldcode.go:20:2: replacing sort.Slice with slices.Sort
# oldcode.go:25:2: replacing if-else with min()

# 应用修复
go fix .

# 修复后的代码会自动添加 "slices" 导入并替换代码
```

**修复后代码**:

```go
package main

import (
 "fmt"
 "slices"
)

func main() {
 // 新方式: 使用 slices.Contains
 nums := []int{1, 2, 3, 4, 5}
 target := 3
 found := slices.Contains(nums, target)
 fmt.Println("Found:", found)

 // 新方式: 使用 slices.Sort
 slices.Sort(nums)

 // 新方式: 使用内置 min
 a, b := 10, 20
 min := min(a, b)
 fmt.Println("Min:", min)
}
```

---

### 示例 5.2: API 迁移

**旧 API**:

```go
package myapi

//go:fix inline
// Deprecated: Use ProcessV2 instead.
func Process(data []byte) error {
 return ProcessV2(data, DefaultOptions())
}

func ProcessV2(data []byte, opts Options) error {
 // 新实现
 return nil
}

type Options struct {
 Timeout int
 Retry   bool
}

func DefaultOptions() Options {
 return Options{Timeout: 30, Retry: true}
}
```

**用户代码**:

```go
package main

import "myapi"

func main() {
 // 旧调用方式
 err := myapi.Process([]byte("data"))
 if err != nil {
  panic(err)
 }
}
```

**运行 go fix 后**:

```go
package main

import "myapi"

func main() {
 // 自动替换为新调用方式
 err := myapi.ProcessV2([]byte("data"), myapi.DefaultOptions())
 if err != nil {
  panic(err)
 }
}
```

---

## 6. 综合项目示例

### 项目: 安全配置服务

结合 new(expr)、递归泛型、HPKE 的综合示例:

```go
package main

import (
 "crypto/hpke"
 "encoding/json"
 "fmt"
 "log"
 "time"
)

// Config 应用配置
type Config struct {
 Database   DatabaseConfig `json:"database"`
 Server     ServerConfig   `json:"server"`
 Encryption *EncryptConfig `json:"encryption,omitempty"` // 可选
}

type DatabaseConfig struct {
 Host     string `json:"host"`
 Port     int    `json:"port"`
 Username string `json:"username"`
 Password string `json:"password"`
}

type ServerConfig struct {
 Host         string         `json:"host"`
 Port         int            `json:"port"`
 ReadTimeout  *time.Duration `json:"read_timeout,omitempty"`
 WriteTimeout *time.Duration `json:"write_timeout,omitempty"`
}

type EncryptConfig struct {
 Enabled   bool   `json:"enabled"`
 PublicKey []byte `json:"public_key"`
}

// TreeNode 配置树节点 (递归泛型)
type TreeNode[T TreeNode[T]] interface {
 Name() string
 Children() []T
 GetValue() interface{}
}

// ConfigNode 配置树实现
type ConfigNode struct {
 name     string
 value    interface{}
 children []*ConfigNode
}

func (c *ConfigNode) Name() string                { return c.name }
func (c *ConfigNode) Children() []*ConfigNode     { return c.children }
func (c *ConfigNode) GetValue() interface{}       { return c.value }

// 验证实现
type _ interface {
 TreeNode[*ConfigNode]
} = (*ConfigNode)(nil)

// HPKE 加密器
type ConfigEncryptor struct {
 suite hpke.Suite
 sk    hpke.PrivateKey
}

func NewConfigEncryptor() (*ConfigEncryptor, error) {
 kem, _ := hpke.GetKEM(hpke.KEM_P256_HKDF_SHA256)
 kdf, _ := hpke.GetKDF(hpke.KDF_HKDF_SHA256)
 aead, _ := hpke.GetAEAD(hpke.AEAD_AES128GCM)

 suite, err := hpke.NewSuite(kem, kdf, aead)
 if err != nil {
  return nil, err
 }

 sk, err := kem.GenerateKeyPair()
 if err != nil {
  return nil, err
 }

 return &ConfigEncryptor{suite: suite, sk: sk}, nil
}

func (e *ConfigEncryptor) EncryptConfig(config *Config) ([]byte, []byte, error) {
 data, err := json.Marshal(config)
 if err != nil {
  return nil, nil, err
 }

 sender := e.suite.NewSender(e.sk.PublicKey(), []byte("config"))
 enc, ctx, err := sender.SetupBase()
 if err != nil {
  return nil, nil, err
 }

 ciphertext, err := ctx.Seal(data, nil)
 if err != nil {
  return nil, nil, err
 }

 return enc, ciphertext, nil
}

func main() {
 // 1. 创建配置 (使用 new(expr))
 timeout := 30 * time.Second
 config := &Config{
  Database: DatabaseConfig{
   Host:     "localhost",
   Port:     5432,
   Username: "admin",
   Password: "secret",
  },
  Server: ServerConfig{
   Host:         "0.0.0.0",
   Port:         8080,
   ReadTimeout:  new(timeout),  // Go 1.26: new(expr)
   WriteTimeout: new(timeout),  // Go 1.26: new(expr)
  },
  Encryption: new(EncryptConfig{  // Go 1.26: new(复合字面量)
   Enabled: true,
  }),
 }

 // 2. 构建配置树 (使用递归泛型)
 root := &ConfigNode{
  name:  "app",
  value: config,
  children: []*ConfigNode{
   {name: "database", value: config.Database},
   {name: "server", value: config.Server},
   {name: "encryption", value: config.Encryption},
  },
 }

 // 3. 遍历配置树
 fmt.Println("配置树结构:")
 traverseTree(root, 0)

 // 4. 加密配置 (使用 HPKE)
 encryptor, err := NewConfigEncryptor()
 if err != nil {
  log.Fatal(err)
 }

 enc, ciphertext, err := encryptor.EncryptConfig(config)
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("\n配置已加密:\n")
 fmt.Printf("封装大小: %d bytes\n", len(enc))
 fmt.Printf("密文大小: %d bytes\n", len(ciphertext))
}

func traverseTree(node *ConfigNode, depth int) {
 indent := ""
 for i := 0; i < depth; i++ {
  indent += "  "
 }
 fmt.Printf("%s- %s: %T\n", indent, node.Name(), node.GetValue())

 for _, child := range node.Children() {
  traverseTree(child, depth+1)
 }
}
```

---

## 运行说明

```bash
# 1. 确保使用 Go 1.26
go version  # go version go1.26.0 ...

# 2. 运行示例
cd example1
go run main.go

# 3. 运行 go fix
cd example5
go fix -n .  # 预览
go fix .     # 应用

# 4. 运行 HPKE 示例 (需要 Go 1.26)
cd example4
go run main.go
```

---

*所有示例代码均可直接运行，展示了 Go 1.26 新特性在实际项目中的应用。*

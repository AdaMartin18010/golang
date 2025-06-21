# 高级数据结构分析

## 目录

- [高级数据结构分析](#高级数据结构分析)
  - [目录](#目录)
  - [概述](#概述)
    - [核心目标](#核心目标)
    - [应用领域](#应用领域)
  - [形式化定义](#形式化定义)
    - [高级数据结构形式化定义](#高级数据结构形式化定义)
    - [高级数据结构分类体系](#高级数据结构分类体系)
    - [高级数据结构评价指标](#高级数据结构评价指标)
  - [树形结构](#树形结构)
    - [3.1 二叉树](#31-二叉树)
      - [3.1.1 形式化定义](#311-形式化定义)
      - [3.1.2 Golang实现](#312-golang实现)
      - [3.1.3 复杂度分析](#313-复杂度分析)
    - [3.2 平衡树](#32-平衡树)
      - [3.2.1 AVL树](#321-avl树)
        - [Golang实现](#golang实现)
        - [复杂度分析](#复杂度分析)
      - [3.2.2 其他平衡树](#322-其他平衡树)
    - [3.3 B树和B+树](#33-b树和b树)
      - [3.3.1 B树](#331-b树)
        - [Golang实现](#golang实现-1)
      - [3.3.2 B+树](#332-b树)
        - [Golang实现](#golang实现-2)
      - [3.3.3 复杂度分析与应用场景](#333-复杂度分析与应用场景)
    - [3.4 红黑树](#34-红黑树)
      - [3.4.1 定义与性质](#341-定义与性质)
      - [3.4.2 Golang实现](#342-golang实现)
      - [3.4.3 复杂度分析](#343-复杂度分析)
      - [3.4.4 应用场景](#344-应用场景)
  - [图结构](#图结构)
    - [4.1 图的表示](#41-图的表示)
      - [4.1.1 形式化定义](#411-形式化定义)
      - [4.1.2 图的表示方法](#412-图的表示方法)
        - [邻接矩阵](#邻接矩阵)
        - [邻接表](#邻接表)
      - [4.1.3 综合比较](#413-综合比较)
    - [4.2 图遍历算法](#42-图遍历算法)
      - [4.2.1 深度优先搜索 (DFS)](#421-深度优先搜索-dfs)
      - [4.2.2 广度优先搜索 (BFS)](#422-广度优先搜索-bfs)
      - [4.2.3 连通分量](#423-连通分量)
      - [4.2.4 强连通分量](#424-强连通分量)
    - [4.3 最短路径算法](#43-最短路径算法)
      - [4.3.1 Dijkstra算法](#431-dijkstra算法)
      - [4.3.2 Bellman-Ford算法](#432-bellman-ford算法)
      - [4.3.3 Floyd-Warshall算法](#433-floyd-warshall算法)
    - [4.4 最小生成树算法](#44-最小生成树算法)
      - [4.4.1 Prim算法](#441-prim算法)
      - [4.4.2 Kruskal算法](#442-kruskal算法)
    - [4.5 拓扑排序](#45-拓扑排序)
  - [高级散列表](#高级散列表)
    - [5.1 开放寻址法](#51-开放寻址法)
    - [5.2 链式散列](#52-链式散列)
    - [5.3 一致性哈希](#53-一致性哈希)
    - [5.4 完美哈希](#54-完美哈希)
    - [5.5 布谷鸟哈希](#55-布谷鸟哈希)
  - [特殊数据结构](#特殊数据结构)
    - [6.1 布隆过滤器](#61-布隆过滤器)
    - [6.2 跳表](#62-跳表)
    - [6.3 字典树变体](#63-字典树变体)
    - [6.4 前缀和与区间树](#64-前缀和与区间树)
    - [6.5 稀疏表与RMQ](#65-稀疏表与rmq)
  - [空间效率数据结构](#空间效率数据结构)
    - [7.1 位图](#71-位图)
    - [7.2 基数树](#72-基数树)
    - [7.3 压缩前缀树](#73-压缩前缀树)
  - [概率数据结构](#概率数据结构)
    - [8.1 Count-Min Sketch](#81-count-min-sketch)
    - [8.2 HyperLogLog](#82-hyperloglog)
    - [8.3 跳表变体](#83-跳表变体)
  - [Golang实现](#golang实现-3)
  - [性能分析与对比](#性能分析与对比)
  - [最佳实践](#最佳实践)
  - [案例研究](#案例研究)
  - [总结](#总结)

## 概述

高级数据结构是解决复杂问题的强大工具，在处理大规模数据、提高算法效率、减少空间消耗等方面具有显著优势。本章节将系统分析各种高级数据结构，包括树形结构、图结构、高级散列表和特殊用途数据结构，并提供它们在Go语言中的实现和优化策略。

### 核心目标

- **形式化定义**: 使用严格的数学定义描述高级数据结构
- **算法分析**: 详细分析各种数据结构的时间和空间复杂度
- **Golang实现**: 提供高效的Go语言实现和使用模式
- **实际应用**: 探讨各种数据结构在实际问题中的应用
- **性能对比**: 比较不同数据结构在各种场景下的性能

### 应用领域

高级数据结构广泛应用于以下领域：

1. **数据库系统**: B树、B+树用于索引结构
2. **网络路由**: 前缀树用于IP路由表查找
3. **搜索引擎**: 倒排索引、布隆过滤器
4. **地理信息系统**: 四叉树、R树
5. **图形学和游戏开发**: 八叉树、BSP树
6. **大数据处理**: 概率数据结构用于流数据分析
7. **分布式系统**: 一致性哈希用于数据分片

## 形式化定义

### 高级数据结构形式化定义

**定义 1.1** (高级数据结构)
一个高级数据结构可以定义为一个七元组：
$$\mathcal{DS} = (T, O, R, I, C, P, A)$$

其中：

- $T$ 是存储的数据类型集合
- $O$ 是支持的操作集合
- $R$ 是内部表示形式
- $I$ 是不变性条件
- $C$ 是复杂度特征
- $P$ 是性能属性
- $A$ 是算法属性

**定义 1.2** (数据结构复杂度模型)
数据结构复杂度模型是一个三元组：
$$\mathcal{CM} = (T, S, M)$$

其中：

- $T: O \times \mathbb{N} \rightarrow \mathbb{R}^+$ 是时间复杂度函数
- $S: \mathbb{N} \rightarrow \mathbb{R}^+$ 是空间复杂度函数
- $M: O \times \mathbb{N} \rightarrow \mathbb{R}^+$ 是访存复杂度函数

**定义 1.3** (数据结构性能模型)
数据结构性能模型是一个四元组：
$$\mathcal{PM} = (L, H, P, D)$$

其中：

- $L$ 是访存局部性函数
- $H$ 是缓存命中率函数
- $P$ 是并行性衡量函数
- $D$ 是数据依赖性函数

### 高级数据结构分类体系

**定义 1.4** (分类体系)
高级数据结构分类体系是一个有向无环图 $G = (V, E)$，其中：

- $V$ 是数据结构类型集合
- $E \subset V \times V$ 是"继承关系"或"特化关系"

主要分类维度包括：

1. **组织方式**: 线性、层次、网络
2. **访问方式**: 随机访问、顺序访问、键值访问
3. **平衡性**: 完全平衡、近似平衡、非平衡
4. **持久性**: 临时性、半持久性、完全持久性
5. **并发性**: 串行、并发安全、无锁

### 高级数据结构评价指标

**定义 1.5** (评价指标)
高级数据结构的评价指标是一个六元组：
$$\mathcal{MI} = (E, F, S, C, M, A)$$

其中：

- $E$ 是效率指标
- $F$ 是功能完备性
- $S$ 是空间利用率
- $C$ 是并发能力
- $M$ 是可维护性
- $A$ 是适应性

## 树形结构

### 3.1 二叉树

#### 3.1.1 形式化定义

**定义 3.1** (二叉树)
二叉树是一个有向无环图 $T = (V, E)$，其中：

- $V$ 是节点集合
- $E \subset V \times V$ 是边集合
- 每个节点最多有两个子节点：左子树和右子树
- 存在唯一的根节点 $r \in V$，没有入边
- 除根节点外，每个节点有且仅有一个父节点

**定义 3.1.1** (完全二叉树)
完全二叉树是指除了最后一层外，其他层都填满节点，且最后一层的节点都集中在左侧。

**定义 3.1.2** (满二叉树)
满二叉树是指除了叶子节点外，每个节点都有两个子节点。

**定理 3.1** (节点数与高度关系)
对于深度为 $h$ 的满二叉树，节点数 $n = 2^{h+1} - 1$。

#### 3.1.2 Golang实现

```go
// 二叉树节点
type BinaryTreeNode struct {
    Value       int
    Left, Right *BinaryTreeNode
}

// 二叉树
type BinaryTree struct {
    Root *BinaryTreeNode
}

// 创建新二叉树
func NewBinaryTree() *BinaryTree {
    return &BinaryTree{}
}

// 插入节点（简单实现，不保证平衡）
func (t *BinaryTree) Insert(value int) {
    newNode := &BinaryTreeNode{Value: value}
    
    if t.Root == nil {
        t.Root = newNode
        return
    }
    
    insertNode(t.Root, newNode)
}

// 递归插入节点
func insertNode(node, newNode *BinaryTreeNode) {
    if newNode.Value < node.Value {
        if node.Left == nil {
            node.Left = newNode
        } else {
            insertNode(node.Left, newNode)
        }
    } else {
        if node.Right == nil {
            node.Right = newNode
        } else {
            insertNode(node.Right, newNode)
        }
    }
}

// 中序遍历
func (t *BinaryTree) InOrderTraversal() []int {
    result := []int{}
    inOrder(t.Root, &result)
    return result
}

func inOrder(node *BinaryTreeNode, result *[]int) {
    if node == nil {
        return
    }
    
    inOrder(node.Left, result)
    *result = append(*result, node.Value)
    inOrder(node.Right, result)
}

// 前序遍历
func (t *BinaryTree) PreOrderTraversal() []int {
    result := []int{}
    preOrder(t.Root, &result)
    return result
}

func preOrder(node *BinaryTreeNode, result *[]int) {
    if node == nil {
        return
    }
    
    *result = append(*result, node.Value)
    preOrder(node.Left, result)
    preOrder(node.Right, result)
}

// 后序遍历
func (t *BinaryTree) PostOrderTraversal() []int {
    result := []int{}
    postOrder(t.Root, &result)
    return result
}

func postOrder(node *BinaryTreeNode, result *[]int) {
    if node == nil {
        return
    }
    
    postOrder(node.Left, result)
    postOrder(node.Right, result)
    *result = append(*result, node.Value)
}

// 层序遍历（广度优先）
func (t *BinaryTree) LevelOrderTraversal() []int {
    if t.Root == nil {
        return nil
    }
    
    result := []int{}
    queue := []*BinaryTreeNode{t.Root}
    
    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        
        result = append(result, node.Value)
        
        if node.Left != nil {
            queue = append(queue, node.Left)
        }
        
        if node.Right != nil {
            queue = append(queue, node.Right)
        }
    }
    
    return result
}
```

#### 3.1.3 复杂度分析

| 操作 | 平均时间复杂度 | 最坏时间复杂度 | 空间复杂度 |
|------|----------------|----------------|------------|
| 访问 | $O(\log n)$    | $O(n)$         | $O(1)$     |
| 搜索 | $O(\log n)$    | $O(n)$         | $O(1)$     |
| 插入 | $O(\log n)$    | $O(n)$         | $O(1)$     |
| 删除 | $O(\log n)$    | $O(n)$         | $O(1)$     |

### 3.2 平衡树

#### 3.2.1 AVL树

**定义 3.2** (AVL树)
AVL树是一种自平衡的二叉搜索树，满足以下性质：

1. 是一个二叉搜索树
2. 对于任意节点，其左右子树高度差不超过1
3. 每个节点都有一个平衡因子(左子树高度减右子树高度)，取值为 $\{-1, 0, 1\}$

**定理 3.2** (AVL树高度)
包含 $n$ 个节点的AVL树高度 $h = O(\log n)$。

##### Golang实现

```go
// AVL树节点
type AVLNode struct {
    Value       int
    Left, Right *AVLNode
    Height      int
}

// AVL树
type AVLTree struct {
    Root *AVLNode
}

// 获取节点高度
func height(node *AVLNode) int {
    if node == nil {
        return 0
    }
    return node.Height
}

// 获取平衡因子
func balanceFactor(node *AVLNode) int {
    if node == nil {
        return 0
    }
    return height(node.Left) - height(node.Right)
}

// 更新节点高度
func updateHeight(node *AVLNode) {
    node.Height = max(height(node.Left), height(node.Right)) + 1
}

// 右旋
func rightRotate(y *AVLNode) *AVLNode {
    x := y.Left
    T2 := x.Right
    
    // 旋转
    x.Right = y
    y.Left = T2
    
    // 更新高度
    updateHeight(y)
    updateHeight(x)
    
    return x
}

// 左旋
func leftRotate(x *AVLNode) *AVLNode {
    y := x.Right
    T2 := y.Left
    
    // 旋转
    y.Left = x
    x.Right = T2
    
    // 更新高度
    updateHeight(x)
    updateHeight(y)
    
    return y
}

// 插入节点
func (t *AVLTree) Insert(value int) {
    t.Root = insertAVL(t.Root, value)
}

func insertAVL(node *AVLNode, value int) *AVLNode {
    // 标准BST插入
    if node == nil {
        return &AVLNode{Value: value, Height: 1}
    }
    
    if value < node.Value {
        node.Left = insertAVL(node.Left, value)
    } else if value > node.Value {
        node.Right = insertAVL(node.Right, value)
    } else {
        // 值已存在，不做任何操作
        return node
    }
    
    // 更新高度
    updateHeight(node)
    
    // 获取平衡因子
    balance := balanceFactor(node)
    
    // 左左情况
    if balance > 1 && value < node.Left.Value {
        return rightRotate(node)
    }
    
    // 右右情况
    if balance < -1 && value > node.Right.Value {
        return leftRotate(node)
    }
    
    // 左右情况
    if balance > 1 && value > node.Left.Value {
        node.Left = leftRotate(node.Left)
        return rightRotate(node)
    }
    
    // 右左情况
    if balance < -1 && value < node.Right.Value {
        node.Right = rightRotate(node.Right)
        return leftRotate(node)
    }
    
    return node
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
```

##### 复杂度分析

| 操作 | 平均时间复杂度 | 最坏时间复杂度 | 空间复杂度 |
|------|----------------|----------------|------------|
| 搜索 | $O(\log n)$    | $O(\log n)$    | $O(1)$     |
| 插入 | $O(\log n)$    | $O(\log n)$    | $O(1)$     |
| 删除 | $O(\log n)$    | $O(\log n)$    | $O(1)$     |

#### 3.2.2 其他平衡树

**伸展树 (Splay Tree)**：通过旋转操作将最近访问的节点移至根部的自平衡二叉搜索树，适用于频繁访问相同元素的场景。

**Treap**：将二叉搜索树与堆相结合的数据结构，每个节点同时具有键值和随机优先级，用于维持树的平衡性。

**替罪羊树 (Scapegoat Tree)**：一种高度平衡的二叉搜索树，通过定期全局重构而不是旋转来保持平衡。

### 3.3 B树和B+树

#### 3.3.1 B树

**定义 3.3** (B树)
B树是一种自平衡的多路搜索树，满足以下条件：

1. 每个节点至多有 $m$ 个子节点（$m$ 阶B树）
2. 除根节点外，每个非叶节点至少有 $\lceil m/2 \rceil$ 个子节点
3. 根节点至少有2个子节点（除非整树只有一个节点）
4. 所有叶节点都在同一层
5. 包含 $n$ 个关键字的节点有 $n+1$ 个子节点

**定理 3.3** (B树高度)
包含 $n$ 个关键字的 $m$ 阶B树的高度为 $h = O(\log_m n)$。

##### Golang实现

```go
// B树节点
type BTreeNode struct {
    Keys     []int
    Children []*BTreeNode
    IsLeaf   bool
    Degree   int // 最小度数，B树的阶为 2*t
}

// B树
type BTree struct {
    Root   *BTreeNode
    Degree int
}

// 创建B树
func NewBTree(degree int) *BTree {
    if degree < 2 {
        panic("B树的最小度数必须至少为2")
    }
    
    return &BTree{
        Root: &BTreeNode{
            Keys:     []int{},
            Children: []*BTreeNode{},
            IsLeaf:   true,
            Degree:   degree,
        },
        Degree: degree,
    }
}

// 搜索关键字
func (t *BTree) Search(key int) bool {
    return searchKey(t.Root, key)
}

func searchKey(node *BTreeNode, key int) bool {
    i := 0
    for i < len(node.Keys) && key > node.Keys[i] {
        i++
    }
    
    if i < len(node.Keys) && key == node.Keys[i] {
        return true
    }
    
    if node.IsLeaf {
        return false
    }
    
    return searchKey(node.Children[i], key)
}

// 插入关键字
func (t *BTree) Insert(key int) {
    root := t.Root
    
    // 如果根节点已满，分裂根节点
    if len(root.Keys) == 2*t.Degree - 1 {
        newRoot := &BTreeNode{
            Keys:     []int{},
            Children: []*BTreeNode{root},
            IsLeaf:   false,
            Degree:   t.Degree,
        }
        
        t.Root = newRoot
        splitChild(newRoot, 0)
        insertNonFull(newRoot, key)
    } else {
        insertNonFull(root, key)
    }
}

// 分裂子节点
func splitChild(node *BTreeNode, index int) {
    t := node.Degree
    y := node.Children[index]
    z := &BTreeNode{
        Keys:     make([]int, t-1),
        IsLeaf:   y.IsLeaf,
        Degree:   y.Degree,
    }
    
    // 将y的后半部分关键字复制到z
    for j := 0; j < t-1; j++ {
        z.Keys[j] = y.Keys[j+t]
    }
    
    // 如果不是叶节点，复制子节点
    if !y.IsLeaf {
        z.Children = make([]*BTreeNode, t)
        for j := 0; j < t; j++ {
            z.Children[j] = y.Children[j+t]
        }
    }
    
    // 截断y的关键字和子节点
    y.Keys = y.Keys[:t-1]
    if !y.IsLeaf {
        y.Children = y.Children[:t]
    }
    
    // 将中间关键字插入父节点
    node.Keys = append(node.Keys, 0)
    copy(node.Keys[index+1:], node.Keys[index:])
    node.Keys[index] = y.Keys[t-1]
    
    // 在父节点中添加z的引用
    node.Children = append(node.Children, nil)
    copy(node.Children[index+2:], node.Children[index+1:])
    node.Children[index+1] = z
}

// 在非满节点中插入
func insertNonFull(node *BTreeNode, key int) {
    i := len(node.Keys) - 1
    
    if node.IsLeaf {
        // 在叶节点中插入
        node.Keys = append(node.Keys, 0)
        for i >= 0 && key < node.Keys[i] {
            node.Keys[i+1] = node.Keys[i]
            i--
        }
        node.Keys[i+1] = key
    } else {
        // 找到插入的子节点
        for i >= 0 && key < node.Keys[i] {
            i--
        }
        i++
        
        // 如果子节点已满，先分裂
        if len(node.Children[i].Keys) == 2*node.Degree - 1 {
            splitChild(node, i)
            if key > node.Keys[i] {
                i++
            }
        }
        
        insertNonFull(node.Children[i], key)
    }
}
```

#### 3.3.2 B+树

**定义 3.4** (B+树)
B+树是B树的变种，满足以下附加条件：

1. 只有叶节点存储数据，内部节点仅用于索引
2. 所有叶节点通过指针链接形成一个有序链表
3. 内部节点的每个键值也存在于叶节点中
4. 所有数据记录都存储在叶节点中

##### Golang实现

```go
// B+树节点
type BPlusTreeNode struct {
    Keys     []int
    Children []*BPlusTreeNode
    IsLeaf   bool
    Next     *BPlusTreeNode // 用于链接叶节点
    Degree   int
}

// B+树
type BPlusTree struct {
    Root   *BPlusTreeNode
    Degree int
}

// 创建B+树
func NewBPlusTree(degree int) *BPlusTree {
    if degree < 2 {
        panic("B+树的最小度数必须至少为2")
    }
    
    return &BPlusTree{
        Root: &BPlusTreeNode{
            Keys:     []int{},
            Children: []*BPlusTreeNode{},
            IsLeaf:   true,
            Degree:   degree,
        },
        Degree: degree,
    }
}

// 搜索关键字
func (t *BPlusTree) Search(key int) bool {
    node := t.Root
    
    for !node.IsLeaf {
        i := 0
        for i < len(node.Keys) && key >= node.Keys[i] {
            i++
        }
        node = node.Children[i]
    }
    
    // 在叶节点中搜索
    for _, k := range node.Keys {
        if k == key {
            return true
        }
    }
    
    return false
}

// 范围查询
func (t *BPlusTree) RangeQuery(start, end int) []int {
    var result []int
    
    // 找到第一个包含start的叶节点
    node := t.Root
    for !node.IsLeaf {
        i := 0
        for i < len(node.Keys) && start >= node.Keys[i] {
            i++
        }
        node = node.Children[i]
    }
    
    // 收集范围内的所有键
    for node != nil {
        for _, key := range node.Keys {
            if key >= start && key <= end {
                result = append(result, key)
            }
            if key > end {
                return result
            }
        }
        node = node.Next
    }
    
    return result
}
```

#### 3.3.3 复杂度分析与应用场景

| 数据结构 | 搜索 | 插入 | 删除 | 范围查询 | 空间 |
|----------|------|------|------|----------|------|
| B树      | $O(\log_m n)$ | $O(\log_m n)$ | $O(\log_m n)$ | $O(\log_m n + k)$ | $O(n)$ |
| B+树     | $O(\log_m n)$ | $O(\log_m n)$ | $O(\log_m n)$ | $O(\log_m n + k)$ | $O(n)$ |

**B树应用场景**:

- 文件系统和数据库索引
- 多级索引

**B+树应用场景**:

- 关系型数据库索引 (如MySQL的InnoDB)
- 需要高效范围查询的场景
- 文件系统

### 3.4 红黑树

#### 3.4.1 定义与性质

**定义 3.5** (红黑树)
红黑树是一种自平衡的二叉搜索树，每个节点都有一个颜色属性（红色或黑色），且满足以下性质：

1. 每个节点要么是红色，要么是黑色
2. 根节点是黑色
3. 所有叶节点（NIL节点）是黑色
4. 如果一个节点是红色，则其两个子节点都是黑色
5. 对于每个节点，从该节点到其所有后代叶节点的简单路径上，包含相同数目的黑色节点

**定理 3.4** (红黑树高度)
包含 $n$ 个节点的红黑树的高度为 $h = O(\log n)$。

#### 3.4.2 Golang实现

```go
// 颜色常量
const (
    RED   = true
    BLACK = false
)

// 红黑树节点
type RBNode struct {
    Value       int
    Left, Right *RBNode
    Color       bool // true为红色，false为黑色
    Parent      *RBNode
}

// 红黑树
type RedBlackTree struct {
    Root *RBNode
    NIL  *RBNode // 哨兵节点
}

// 创建红黑树
func NewRedBlackTree() *RedBlackTree {
    nil := &RBNode{Color: BLACK}
    return &RedBlackTree{
        NIL:  nil,
        Root: nil,
    }
}

// 左旋
func (t *RedBlackTree) leftRotate(x *RBNode) {
    y := x.Right
    x.Right = y.Left
    
    if y.Left != t.NIL {
        y.Left.Parent = x
    }
    
    y.Parent = x.Parent
    
    if x.Parent == t.NIL {
        t.Root = y
    } else if x == x.Parent.Left {
        x.Parent.Left = y
    } else {
        x.Parent.Right = y
    }
    
    y.Left = x
    x.Parent = y
}

// 右旋
func (t *RedBlackTree) rightRotate(y *RBNode) {
    x := y.Left
    y.Left = x.Right
    
    if x.Right != t.NIL {
        x.Right.Parent = y
    }
    
    x.Parent = y.Parent
    
    if y.Parent == t.NIL {
        t.Root = x
    } else if y == y.Parent.Left {
        y.Parent.Left = x
    } else {
        y.Parent.Right = x
    }
    
    x.Right = y
    y.Parent = x
}

// 插入修复
func (t *RedBlackTree) insertFixup(z *RBNode) {
    for z.Parent.Color == RED {
        if z.Parent == z.Parent.Parent.Left {
            y := z.Parent.Parent.Right
            
            if y.Color == RED {
                // 情况1: 叔叔是红色
                z.Parent.Color = BLACK
                y.Color = BLACK
                z.Parent.Parent.Color = RED
                z = z.Parent.Parent
            } else {
                if z == z.Parent.Right {
                    // 情况2: 叔叔是黑色，当前节点是右子节点
                    z = z.Parent
                    t.leftRotate(z)
                }
                
                // 情况3: 叔叔是黑色，当前节点是左子节点
                z.Parent.Color = BLACK
                z.Parent.Parent.Color = RED
                t.rightRotate(z.Parent.Parent)
            }
        } else {
            // 对称情况
            y := z.Parent.Parent.Left
            
            if y.Color == RED {
                z.Parent.Color = BLACK
                y.Color = BLACK
                z.Parent.Parent.Color = RED
                z = z.Parent.Parent
            } else {
                if z == z.Parent.Left {
                    z = z.Parent
                    t.rightRotate(z)
                }
                
                z.Parent.Color = BLACK
                z.Parent.Parent.Color = RED
                t.leftRotate(z.Parent.Parent)
            }
        }
        
        if z == t.Root {
            break
        }
    }
    
    t.Root.Color = BLACK
}

// 插入节点
func (t *RedBlackTree) Insert(value int) {
    z := &RBNode{
        Value:  value,
        Color:  RED,
        Left:   t.NIL,
        Right:  t.NIL,
        Parent: t.NIL,
    }
    
    var y *RBNode = t.NIL
    x := t.Root
    
    // 找到插入位置
    for x != t.NIL {
        y = x
        if z.Value < x.Value {
            x = x.Left
        } else {
            x = x.Right
        }
    }
    
    z.Parent = y
    
    if y == t.NIL {
        t.Root = z // 树为空
    } else if z.Value < y.Value {
        y.Left = z
    } else {
        y.Right = z
    }
    
    // 修复红黑树属性
    t.insertFixup(z)
}
```

#### 3.4.3 复杂度分析

| 操作 | 平均时间复杂度 | 最坏时间复杂度 | 空间复杂度 |
|------|----------------|----------------|------------|
| 搜索 | $O(\log n)$    | $O(\log n)$    | $O(1)$     |
| 插入 | $O(\log n)$    | $O(\log n)$    | $O(1)$     |
| 删除 | $O(\log n)$    | $O(\log n)$    | $O(1)$     |

#### 3.4.4 应用场景

- C++ STL中的`map`和`set`
- Java中的`TreeMap`和`TreeSet`
- 内核数据结构
- 实时系统中的调度

## 图结构

### 4.1 图的表示

#### 4.1.1 形式化定义

**定义 4.1** (图)
图是一个二元组 $G = (V, E)$，其中：

- $V$ 是顶点集合
- $E \subseteq V \times V$ 是边集合

对于有向图，边 $(u, v) \in E$ 表示从 $u$ 到 $v$ 的有向边。
对于无向图，边 $\{u, v\} \in E$ 表示连接 $u$ 和 $v$ 的无向边。

**定义 4.1.1** (有权图)
有权图是一个三元组 $G = (V, E, w)$，其中：

- $V$ 是顶点集合
- $E \subseteq V \times V$ 是边集合
- $w: E \rightarrow \mathbb{R}$ 是权重函数

#### 4.1.2 图的表示方法

##### 邻接矩阵

**定义 4.1.2** (邻接矩阵)
对于 $n$ 个顶点的图 $G = (V, E)$，其邻接矩阵 $A$ 是一个 $n \times n$ 的矩阵，定义为：

$$
A[i][j] = \begin{cases}
1 & \text{若} (i, j) \in E \\
0 & \text{否则}
\end{cases}
$$

对于有权图：
$$
A[i][j] = \begin{cases}
w(i, j) & \text{若} (i, j) \in E \\
\infty & \text{否则}
\end{cases}
$$

```go
// 基于邻接矩阵的图
type AdjMatrixGraph struct {
    V     int       // 顶点数
    E     int       // 边数
    Directed bool   // 是否是有向图
    Matrix [][]int  // 邻接矩阵
}

// 创建图
func NewAdjMatrixGraph(v int, directed bool) *AdjMatrixGraph {
    matrix := make([][]int, v)
    for i := range matrix {
        matrix[i] = make([]int, v)
    }

    return &AdjMatrixGraph{
        V:        v,
        E:        0,
        Directed: directed,
        Matrix:   matrix,
    }
}

// 添加边
func (g *AdjMatrixGraph) AddEdge(v, w int, weight int) {
    if v < 0 || v >= g.V || w < 0 || w >= g.V {
        panic("顶点超出范围")
    }

    // 检查边是否已存在
    if g.Matrix[v][w] == 0 {
        g.E++
    }

    g.Matrix[v][w] = weight

    if !g.Directed {
        g.Matrix[w][v] = weight
    }
}

// 检查边是否存在
func (g *AdjMatrixGraph) HasEdge(v, w int) bool {
    if v < 0 || v >= g.V || w < 0 || w >= g.V {
        return false
    }
    return g.Matrix[v][w] != 0
}

// 获取顶点的所有邻居
func (g *AdjMatrixGraph) Adj(v int) []int {
    if v < 0 || v >= g.V {
        panic("顶点超出范围")
    }

    var adj []int
    for i := 0; i < g.V; i++ {
        if g.Matrix[v][i] != 0 {
            adj = append(adj, i)
        }
    }

    return adj
}
```

**空间复杂度**: $O(V^2)$

##### 邻接表

**定义 4.1.3** (邻接表)
对于图 $G = (V, E)$，其邻接表是一个长度为 $|V|$ 的数组，其中每个元素是一个链表，存储与对应顶点相邻的所有顶点。

```go
// 邻接表中的边
type Edge struct {
    To     int
    Weight int
    Next   *Edge
}

// 基于邻接表的图
type AdjListGraph struct {
    V        int      // 顶点数
    E        int      // 边数
    Directed bool     // 是否是有向图
    Adj      []*Edge  // 邻接表
}

// 创建图
func NewAdjListGraph(v int, directed bool) *AdjListGraph {
    adj := make([]*Edge, v)
    return &AdjListGraph{
        V:        v,
        E:        0,
        Directed: directed,
        Adj:      adj,
    }
}

// 添加边
func (g *AdjListGraph) AddEdge(v, w int, weight int) {
    if v < 0 || v >= g.V || w < 0 || w >= g.V {
        panic("顶点超出范围")
    }

    // 添加 v->w 的边
    newEdge := &Edge{To: w, Weight: weight, Next: g.Adj[v]}
    g.Adj[v] = newEdge
    g.E++

    // 如果是无向图，添加 w->v 的边
    if !g.Directed {
        newEdge = &Edge{To: v, Weight: weight, Next: g.Adj[w]}
        g.Adj[w] = newEdge
    }
}

// 获取顶点的所有邻居
func (g *AdjListGraph) GetAdjVertices(v int) []int {
    if v < 0 || v >= g.V {
        panic("顶点超出范围")
    }

    var vertices []int
    for e := g.Adj[v]; e != nil; e = e.Next {
        vertices = append(vertices, e.To)
    }

    return vertices
}
```

**空间复杂度**: $O(V + E)$

#### 4.1.3 综合比较

| 表示方法 | 空间复杂度 | 查询边 | 枚举邻居 | 适用场景 |
|----------|----------|--------|---------|---------|
| 邻接矩阵 | $O(V^2)$ | $O(1)$ | $O(V)$ | 稠密图，边的查询频繁 |
| 邻接表   | $O(V+E)$ | $O(V)$ | $O(\text{degree}(v))$ | 稀疏图，需要遍历邻居 |

### 4.2 图遍历算法

#### 4.2.1 深度优先搜索 (DFS)

**定义 4.2** (深度优先搜索)
深度优先搜索是一种用于遍历图的算法，从根节点开始，尽可能深地探索一条路径，直到无法继续前进时回溯到上一个节点继续探索。

```go
// 深度优先搜索
func (g *AdjListGraph) DFS(start int) []int {
    visited := make([]bool, g.V)
    var result []int

    // 递归DFS函数
    var dfsRecursive func(v int)
    dfsRecursive = func(v int) {
        visited[v] = true
        result = append(result, v)

        for e := g.Adj[v]; e != nil; e = e.Next {
            if !visited[e.To] {
                dfsRecursive(e.To)
            }
        }
    }

    dfsRecursive(start)
    return result
}

// 非递归实现的DFS
func (g *AdjListGraph) IterativeDFS(start int) []int {
    visited := make([]bool, g.V)
    var result []int

    // 使用栈来模拟递归
    stack := []int{start}

    for len(stack) > 0 {
        // 弹出栈顶元素
        v := stack[len(stack)-1]
        stack = stack[:len(stack)-1]

        if !visited[v] {
            visited[v] = true
            result = append(result, v)

            // 将邻居入栈 (注意：为保持与递归DFS相同的访问顺序，我们需要逆序入栈)
            neighbors := g.GetAdjVertices(v)
            for i := len(neighbors) - 1; i >= 0; i-- {
                if !visited[neighbors[i]] {
                    stack = append(stack, neighbors[i])
                }
            }
        }
    }

    return result
}
```

**时间复杂度**: $O(V + E)$
**空间复杂度**: $O(V)$

#### 4.2.2 广度优先搜索 (BFS)

**定义 4.3** (广度优先搜索)
广度优先搜索是一种用于遍历图的算法，从根节点开始，先访问所有邻近节点，然后再访问下一层节点。

```go
// 广度优先搜索
func (g *AdjListGraph) BFS(start int) []int {
    visited := make([]bool, g.V)
    var result []int

    // 使用队列进行BFS
    queue := []int{start}
    visited[start] = true

    for len(queue) > 0 {
        // 出队
        v := queue[0]
        queue = queue[1:]
        result = append(result, v)

        // 将所有未访问的邻居入队
        for e := g.Adj[v]; e != nil; e = e.Next {
            if !visited[e.To] {
                visited[e.To] = true
                queue = append(queue, e.To)
            }
        }
    }

    return result
}
```

**时间复杂度**: $O(V + E)$
**空间复杂度**: $O(V)$

#### 4.2.3 连通分量

**定义 4.4** (连通分量)
在无向图中，连通分量是极大连通子图，即子图中任意两点都是连通的，且不能再添加任何顶点使其保持连通性。

```go
// 查找连通分量
func (g *AdjListGraph) ConnectedComponents() [][]int {
    visited := make([]bool, g.V)
    var components [][]int

    for i := 0; i < g.V; i++ {
        if !visited[i] {
            // 对每个未访问的顶点开始DFS
            var component []int

            var dfs func(v int)
            dfs = func(v int) {
                visited[v] = true
                component = append(component, v)

                for e := g.Adj[v]; e != nil; e = e.Next {
                    if !visited[e.To] {
                        dfs(e.To)
                    }
                }
            }

            dfs(i)
            components = append(components, component)
        }
    }

    return components
}
```

**时间复杂度**: $O(V + E)$
**空间复杂度**: $O(V)$

#### 4.2.4 强连通分量

**定义 4.5** (强连通分量)
在有向图中，强连通分量是极大强连通子图，即子图中任意两点都是互相可达的，且不能再添加任何顶点使其保持强连通性。

**Kosaraju算法**:

```go
// 使用Kosaraju算法查找强连通分量
func (g *AdjListGraph) StronglyConnectedComponents() [][]int {
    if !g.Directed {
        return g.ConnectedComponents()
    }

    // 第一次DFS，填充结束时间栈
    visited := make([]bool, g.V)
    var finishOrder []int

    var dfs1 func(v int)
    dfs1 = func(v int) {
        visited[v] = true

        for e := g.Adj[v]; e != nil; e = e.Next {
            if !visited[e.To] {
                dfs1(e.To)
            }
        }

        finishOrder = append(finishOrder, v)
    }

    for i := 0; i < g.V; i++ {
        if !visited[i] {
            dfs1(i)
        }
    }

    // 创建图的转置
    transpose := NewAdjListGraph(g.V, true)
    for v := 0; v < g.V; v++ {
        for e := g.Adj[v]; e != nil; e = e.Next {
            transpose.AddEdge(e.To, v, e.Weight)
        }
    }

    // 第二次DFS，按照完成时间的逆序访问顶点
    for i := range visited {
        visited[i] = false
    }

    var components [][]int

    for i := len(finishOrder) - 1; i >= 0; i-- {
        v := finishOrder[i]
        if !visited[v] {
            var component []int

            var dfs2 func(v int)
            dfs2 = func(v int) {
                visited[v] = true
                component = append(component, v)

                for e := transpose.Adj[v]; e != nil; e = e.Next {
                    if !visited[e.To] {
                        dfs2(e.To)
                    }
                }
            }

            dfs2(v)
            components = append(components, component)
        }
    }

    return components
}
```

**时间复杂度**: $O(V + E)$
**空间复杂度**: $O(V)$

### 4.3 最短路径算法

#### 4.3.1 Dijkstra算法

**定义 4.6** (Dijkstra算法)
Dijkstra算法用于计算带非负权重的图中，从源顶点到所有其他顶点的最短路径。

```go
// Dijkstra算法实现
func (g *AdjListGraph) Dijkstra(start int) ([]int, []int) {
    // dist[i]表示从start到i的最短距离
    dist := make([]int, g.V)
    // prev[i]表示从start到i的最短路径上i的前驱节点
    prev := make([]int, g.V)
    // 标记是否已确定最短路径
    visited := make([]bool, g.V)

    // 初始化
    for i := range dist {
        dist[i] = math.MaxInt32
        prev[i] = -1
    }
    dist[start] = 0

    // 使用优先队列优化
    pq := &PriorityQueue{}
    heap.Init(pq)
    heap.Push(pq, &Item{value: start, priority: 0})

    for pq.Len() > 0 {
        // 获取当前最小距离的顶点
        item := heap.Pop(pq).(*Item)
        u := item.value

        if visited[u] {
            continue
        }
        visited[u] = true

        // 更新所有邻居的距离
        for e := g.Adj[u]; e != nil; e = e.Next {
            v := e.To
            if !visited[v] && dist[u] != math.MaxInt32 && dist[u]+e.Weight < dist[v] {
                dist[v] = dist[u] + e.Weight
                prev[v] = u
                heap.Push(pq, &Item{value: v, priority: dist[v]})
            }
        }
    }

    return dist, prev
}

// 优先队列实现
type Item struct {
    value    int
    priority int
    index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
    return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
    pq[i].index = i
    pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
    n := len(*pq)
    item := x.(*Item)
    item.index = n
    *pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
    old := *pq
    n := len(old)
    item := old[n-1]
    old[n-1] = nil
    item.index = -1
    *pq = old[0 : n-1]
    return item
}

// 重建最短路径
func ReconstructPath(start, end int, prev []int) []int {
    if prev[end] == -1 {
        return nil // 路径不存在
    }

    path := []int{end}
    for curr := end; curr != start; curr = prev[curr] {
        path = append([]int{prev[curr]}, path...)
    }

    return path
}
```

**时间复杂度**：使用二叉堆优化的Dijkstra算法时间复杂度为$O((V+E)\log V)$
**空间复杂度**：$O(V)$

#### 4.3.2 Bellman-Ford算法

**定义 4.7** (Bellman-Ford算法)
Bellman-Ford算法用于计算带权重（可以为负）的图中，从源顶点到所有其他顶点的最短路径。它还能检测负权回路。

```go
// Bellman-Ford算法实现
func (g *AdjListGraph) BellmanFord(start int) ([]int, []int, bool) {
    // dist[i]表示从start到i的最短距离
    dist := make([]int, g.V)
    // prev[i]表示从start到i的最短路径上i的前驱节点
    prev := make([]int, g.V)

    // 初始化
    for i := range dist {
        dist[i] = math.MaxInt32
        prev[i] = -1
    }
    dist[start] = 0

    // 松弛操作，最多V-1次
    for i := 0; i < g.V-1; i++ {
        for v := 0; v < g.V; v++ {
            for e := g.Adj[v]; e != nil; e = e.Next {
                u := v
                v := e.To
                w := e.Weight

                if dist[u] != math.MaxInt32 && dist[u]+w < dist[v] {
                    dist[v] = dist[u] + w
                    prev[v] = u
                }
            }
        }
    }

    // 检测负权回路
    for v := 0; v < g.V; v++ {
        for e := g.Adj[v]; e != nil; e = e.Next {
            u := v
            v := e.To
            w := e.Weight

            if dist[u] != math.MaxInt32 && dist[u]+w < dist[v] {
                return nil, nil, true // 存在负权回路
            }
        }
    }

    return dist, prev, false
}
```

**时间复杂度**：$O(V \cdot E)$
**空间复杂度**：$O(V)$

#### 4.3.3 Floyd-Warshall算法

**定义 4.8** (Floyd-Warshall算法)
Floyd-Warshall算法用于计算加权图中所有顶点对之间的最短路径。

```go
// Floyd-Warshall算法实现
func (g *AdjMatrixGraph) FloydWarshall() [][]int {
    // dist[i][j]表示从i到j的最短距离
    dist := make([][]int, g.V)
    for i := range dist {
        dist[i] = make([]int, g.V)
        for j := range dist[i] {
            if i == j {
                dist[i][j] = 0
            } else if g.Matrix[i][j] != 0 {
                dist[i][j] = g.Matrix[i][j]
            } else {
                dist[i][j] = math.MaxInt32
            }
        }
    }

    // 动态规划核心
    for k := 0; k < g.V; k++ {
        for i := 0; i < g.V; i++ {
            for j := 0; j < g.V; j++ {
                if dist[i][k] != math.MaxInt32 && dist[k][j] != math.MaxInt32 &&
                   dist[i][k]+dist[k][j] < dist[i][j] {
                    dist[i][j] = dist[i][k] + dist[k][j]
                }
            }
        }
    }

    return dist
}
```

**时间复杂度**：$O(V^3)$
**空间复杂度**：$O(V^2)$

### 4.4 最小生成树算法

#### 4.4.1 Prim算法

**定义 4.9** (最小生成树)
最小生成树是一张连通加权无向图的一棵权重和最小的生成树。

**定义 4.10** (Prim算法)
Prim算法是一种用于计算连通加权无向图的最小生成树的贪心算法。

```go
// Prim算法实现
func (g *AdjListGraph) Prim(start int) ([]int, int) {
    if g.Directed {
        panic("Prim算法仅适用于无向图")
    }

    // 标记顶点是否已加入MST
    inMST := make([]bool, g.V)
    // 记录各顶点加入MST的边的权重
    key := make([]int, g.V)
    // 记录各顶点的父节点
    parent := make([]int, g.V)

    // 初始化
    for i := range key {
        key[i] = math.MaxInt32
        parent[i] = -1
    }
    key[start] = 0

    // 使用优先队列优化
    pq := &PriorityQueue{}
    heap.Init(pq)
    heap.Push(pq, &Item{value: start, priority: 0})

    for pq.Len() > 0 {
        // 获取当前权重最小的顶点
        item := heap.Pop(pq).(*Item)
        u := item.value

        if inMST[u] {
            continue
        }
        inMST[u] = true

        // 更新所有邻居的权重
        for e := g.Adj[u]; e != nil; e = e.Next {
            v := e.To
            w := e.Weight

            if !inMST[v] && w < key[v] {
                key[v] = w
                parent[v] = u
                heap.Push(pq, &Item{value: v, priority: key[v]})
            }
        }
    }

    // 计算MST的总权重
    mstWeight := 0
    for i := 0; i < g.V; i++ {
        if i != start && parent[i] != -1 {
            mstWeight += key[i]
        }
    }

    return parent, mstWeight
}
```

**时间复杂度**：$O((V+E)\log V)$ (使用二叉堆优化)
**空间复杂度**：$O(V)$

#### 4.4.2 Kruskal算法

**定义 4.11** (Kruskal算法)
Kruskal算法是另一种计算最小生成树的贪心算法，通过按权重递增的顺序添加边来构建MST。

```go
// 并查集数据结构
type DisjointSet struct {
    Parent []int
    Rank   []int
}

// 初始化并查集
func NewDisjointSet(n int) *DisjointSet {
    parent := make([]int, n)
    rank := make([]int, n)

    for i := range parent {
        parent[i] = i
    }

    return &DisjointSet{
        Parent: parent,
        Rank:   rank,
    }
}

// 查找集合代表元素（带路径压缩）
func (ds *DisjointSet) Find(x int) int {
    if ds.Parent[x] != x {
        ds.Parent[x] = ds.Find(ds.Parent[x])
    }
    return ds.Parent[x]
}

// 合并两个集合
func (ds *DisjointSet) Union(x, y int) {
    rootX := ds.Find(x)
    rootY := ds.Find(y)

    if rootX == rootY {
        return
    }

    // 按秩合并
    if ds.Rank[rootX] < ds.Rank[rootY] {
        ds.Parent[rootX] = rootY
    } else if ds.Rank[rootX] > ds.Rank[rootY] {
        ds.Parent[rootY] = rootX
    } else {
        ds.Parent[rootY] = rootX
        ds.Rank[rootX]++
    }
}

// 存储边的结构
type KruskalEdge struct {
    From   int
    To     int
    Weight int
}

// Kruskal算法实现
func (g *AdjListGraph) Kruskal() ([]KruskalEdge, int) {
    if g.Directed {
        panic("Kruskal算法仅适用于无向图")
    }

    // 收集所有边
    var edges []KruskalEdge
    for u := 0; u < g.V; u++ {
        for e := g.Adj[u]; e != nil; e = e.Next {
            v := e.To
            // 对于无向图，只添加一次边(u,v)和(v,u)中的一个
            if u < v {
                edges = append(edges, KruskalEdge{From: u, To: v, Weight: e.Weight})
            }
        }
    }

    // 按权重排序边
    sort.Slice(edges, func(i, j int) bool {
        return edges[i].Weight < edges[j].Weight
    })

    // 初始化并查集
    ds := NewDisjointSet(g.V)

    // 存储MST的边和总权重
    var mstEdges []KruskalEdge
    mstWeight := 0

    // 贪心选择边
    for _, edge := range edges {
        if ds.Find(edge.From) != ds.Find(edge.To) {
            mstEdges = append(mstEdges, edge)
            mstWeight += edge.Weight
            ds.Union(edge.From, edge.To)

            // 当MST有V-1条边时结束
            if len(mstEdges) == g.V-1 {
                break
            }
        }
    }

    return mstEdges, mstWeight
}
```

**时间复杂度**：$O(E \log E)$ (主要来自边的排序)
**空间复杂度**：$O(V + E)$

### 4.5 拓扑排序

**定义 4.12** (拓扑排序)
拓扑排序是有向无环图中所有顶点的线性排序，使得对于每个顶点 $u$ 到 $v$ 的有向边 $(u,v)$，$u$ 在排序中都出现在 $v$ 之前。

```go
// Kahn算法实现拓扑排序
func (g *AdjListGraph) TopologicalSort() ([]int, bool) {
    if !g.Directed {
        panic("拓扑排序仅适用于有向图")
    }

    // 计算每个顶点的入度
    inDegree := make([]int, g.V)
    for u := 0; u < g.V; u++ {
        for e := g.Adj[u]; e != nil; e = e.Next {
            inDegree[e.To]++
        }
    }

    // 将所有入度为0的顶点加入队列
    var queue []int
    for i := 0; i < g.V; i++ {
        if inDegree[i] == 0 {
            queue = append(queue, i)
        }
    }

    // 存储拓扑排序结果
    var result []int

    // BFS遍历
    for len(queue) > 0 {
        u := queue[0]
        queue = queue[1:]
        result = append(result, u)

        // 减少所有邻居的入度
        for e := g.Adj[u]; e != nil; e = e.Next {
            inDegree[e.To]--
            if inDegree[e.To] == 0 {
                queue = append(queue, e.To)
            }
        }
    }

    // 如果结果中的顶点数小于图的顶点数，说明图中有环
    if len(result) != g.V {
        return nil, false
    }

    return result, true
}

// 使用DFS实现拓扑排序
func (g *AdjListGraph) TopologicalSortDFS() ([]int, bool) {
    if !g.Directed {
        panic("拓扑排序仅适用于有向图")
    }

    visited := make([]bool, g.V)
    temp := make([]bool, g.V)  // 用于检测环
    var order []int

    // DFS函数
    var hasCycle bool
    var dfs func(v int)
    dfs = func(v int) {
        if temp[v] {
            hasCycle = true
            return
        }

        if visited[v] || hasCycle {
            return
        }

        temp[v] = true

        for e := g.Adj[v]; e != nil; e = e.Next {
            dfs(e.To)
        }

        temp[v] = false
        visited[v] = true
        order = append([]int{v}, order...) // 前插入
    }

    // 对每个顶点调用DFS
    for i := 0; i < g.V; i++ {
        if !visited[i] {
            dfs(i)
        }
    }

    if hasCycle {
        return nil, false
    }

    return order, true
}
```

**时间复杂度**：$O(V + E)$ (对于Kahn算法和DFS算法)
**空间复杂度**：$O(V)$

## 高级散列表

### 5.1 开放寻址法

### 5.2 链式散列

### 5.3 一致性哈希

### 5.4 完美哈希

### 5.5 布谷鸟哈希

## 特殊数据结构

### 6.1 布隆过滤器

### 6.2 跳表

### 6.3 字典树变体

### 6.4 前缀和与区间树

### 6.5 稀疏表与RMQ

## 空间效率数据结构

### 7.1 位图

### 7.2 基数树

### 7.3 压缩前缀树

## 概率数据结构

### 8.1 Count-Min Sketch

### 8.2 HyperLogLog

### 8.3 跳表变体

## Golang实现

## 性能分析与对比

## 最佳实践

## 案例研究

## 总结

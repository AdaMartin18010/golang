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
      - [定义 5.1](#定义-51)
      - [定义 5.1.1](#定义-511)
      - [定义 5.1.2](#定义-512)
      - [定义 5.1.3](#定义-513)
      - [Golang实现](#golang实现-4)
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
      - [定义 7.1](#定义-71)
      - [Golang实现](#golang实现-5)
    - [7.2 基数树](#72-基数树)
      - [定义 7.2](#定义-72)
      - [Golang实现](#golang实现-6)
    - [7.3 压缩前缀树](#73-压缩前缀树)
      - [定义 7.3](#定义-73)
      - [Golang实现](#golang实现-7)
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
    Value       interface{}
    Left, Right *BinaryTreeNode
}

// 创建新节点
func NewBinaryTreeNode(value interface{}) *BinaryTreeNode {
    return &BinaryTreeNode{
        Value: value,
        Left:  nil,
        Right: nil,
    }
}

// 二叉树
type BinaryTree struct {
    Root *BinaryTreeNode
    Size int
}

// 创建新二叉树
func NewBinaryTree() *BinaryTree {
    return &BinaryTree{
        Root: nil,
        Size: 0,
    }
}

// 插入值（构建二叉搜索树）
func (t *BinaryTree) Insert(value interface{}) {
    // 使用类型断言确保可比较
    var compareValue interface{} = value
    
    newNode := NewBinaryTreeNode(value)
    t.Size++
    
    if t.Root == nil {
        t.Root = newNode
        return
    }
    
    var insertNode func(node *BinaryTreeNode)
    insertNode = func(node *BinaryTreeNode) {
        // 使用类型断言进行比较
        if compare(compareValue, node.Value) < 0 {
            if node.Left == nil {
                node.Left = newNode
            } else {
                insertNode(node.Left)
            }
        } else {
            if node.Right == nil {
                node.Right = newNode
            } else {
                insertNode(node.Right)
            }
        }
    }
    
    insertNode(t.Root)
}

// 比较函数
func compare(a, b interface{}) int {
    switch va := a.(type) {
    case int:
        if vb, ok := b.(int); ok {
            return va - vb
        }
    case string:
        if vb, ok := b.(string); ok {
            if va < vb {
                return -1
            } else if va > vb {
                return 1
            }
            return 0
        }
    case float64:
        if vb, ok := b.(float64); ok {
            if va < vb {
                return -1
            } else if va > vb {
                return 1
            }
            return 0
        }
    }
    
    // 默认情况，比较指针地址
    return 0
}

// 前序遍历
func (t *BinaryTree) PreOrderTraversal() []interface{} {
    result := make([]interface{}, 0, t.Size)
    
    var traverse func(node *BinaryTreeNode)
    traverse = func(node *BinaryTreeNode) {
        if node == nil {
            return
        }
        
        result = append(result, node.Value)
        traverse(node.Left)
        traverse(node.Right)
    }
    
    traverse(t.Root)
    return result
}

// 中序遍历
func (t *BinaryTree) InOrderTraversal() []interface{} {
    result := make([]interface{}, 0, t.Size)
    
    var traverse func(node *BinaryTreeNode)
    traverse = func(node *BinaryTreeNode) {
        if node == nil {
            return
        }
        
        traverse(node.Left)
        result = append(result, node.Value)
        traverse(node.Right)
    }
    
    traverse(t.Root)
    return result
}

// 后序遍历
func (t *BinaryTree) PostOrderTraversal() []interface{} {
    result := make([]interface{}, 0, t.Size)
    
    var traverse func(node *BinaryTreeNode)
    traverse = func(node *BinaryTreeNode) {
        if node == nil {
            return
        }
        
        traverse(node.Left)
        traverse(node.Right)
        result = append(result, node.Value)
    }
    
    traverse(t.Root)
    return result
}

// 层序遍历
func (t *BinaryTree) LevelOrderTraversal() []interface{} {
    if t.Root == nil {
        return []interface{}{}
    }
    
    result := make([]interface{}, 0, t.Size)
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

// 查找值
func (t *BinaryTree) Search(value interface{}) bool {
    var searchNode func(node *BinaryTreeNode) bool
    searchNode = func(node *BinaryTreeNode) bool {
        if node == nil {
            return false
        }
        
        cmp := compare(value, node.Value)
        if cmp == 0 {
            return true
        } else if cmp < 0 {
            return searchNode(node.Left)
        } else {
            return searchNode(node.Right)
        }
    }
    
    return searchNode(t.Root)
}

// 获取高度
func (t *BinaryTree) Height() int {
    var height func(node *BinaryTreeNode) int
    height = func(node *BinaryTreeNode) int {
        if node == nil {
            return 0
        }
        
        leftHeight := height(node.Left)
        rightHeight := height(node.Right)
        
        if leftHeight > rightHeight {
            return leftHeight + 1
        }
        return rightHeight + 1
    }
    
    return height(t.Root)
}

// 判断是否为完全二叉树
func (t *BinaryTree) IsComplete() bool {
    if t.Root == nil {
        return true
    }
    
    queue := []*BinaryTreeNode{t.Root}
    var endFlag bool
    
    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        
        if node == nil {
            endFlag = true
        } else {
            if endFlag {
                return false
            }
            
            queue = append(queue, node.Left)
            queue = append(queue, node.Right)
        }
    }
    
    return true
}

// 判断是否为满二叉树
func (t *BinaryTree) IsFull() bool {
    height := t.Height()
    nodeCount := t.Size
    
    // 满二叉树节点数 = 2^h - 1
    return nodeCount == (1<<height) - 1
}
```

#### 3.1.3 复杂度分析

**时间复杂度**:

| 操作 | 平均情况 | 最坏情况 |
|------|----------|----------|
| 插入 | $O(\log n)$ | $O(n)$ |
| 查找 | $O(\log n)$ | $O(n)$ |
| 删除 | $O(\log n)$ | $O(n)$ |
| 遍历 | $O(n)$ | $O(n)$ |

**空间复杂度**:

- 存储: $O(n)$
- 递归遍历: $O(h)$ 其中 $h$ 是树的高度

**定理 3.2** (二叉搜索树高度)
对于包含 $n$ 个节点的随机构建的二叉搜索树，其期望高度为 $O(\log n)$。

**证明**:
随机插入 $n$ 个不同的键，生成的二叉搜索树的期望高度为 $O(\log n)$。这是因为每个节点被放置在树中的位置是随机的，使得树保持相对平衡。

**定理 3.3** (二叉树遍历复杂度)
任何二叉树的前序、中序、后序和层序遍历的时间复杂度均为 $O(n)$，其中 $n$ 是节点数。

**证明**:
在任何遍历算法中，每个节点都被访问一次且仅一次，因此时间复杂度为 $O(n)$。

**注意**:
- 在最坏情况下（如链状树），二叉搜索树的高度可达 $O(n)$，此时查找、插入和删除操作的时间复杂度退化为 $O(n)$。
- 为了避免这种情况，可以使用自平衡二叉树，如AVL树或红黑树。

### 3.2 平衡树

#### 3.2.1 AVL树

**定义 3.2.1** (AVL树)
AVL树是一种自平衡二叉搜索树，其中每个节点的左右子树高度差不超过1。

**定义 3.2.2** (平衡因子)
节点的平衡因子定义为其左子树高度减去右子树高度。在AVL树中，每个节点的平衡因子必须是 -1、0 或 1。

**定理 3.4** (AVL树高度)
包含 $n$ 个节点的AVL树高度不超过 $1.44 \log_2(n+2) - 1.328$。

**证明**:
设 $N(h)$ 是高度为 $h$ 的AVL树的最少节点数。可以证明 $N(h) = N(h-1) + N(h-2) + 1$，这与斐波那契数列类似。解此递归关系，得到 $N(h) \approx \frac{\phi^h}{\sqrt{5}}$，其中 $\phi = \frac{1+\sqrt{5}}{2}$ 是黄金比例。因此，$h \leq 1.44 \log_2(n+2) - 1.328$。

##### Golang实现

```go
// AVL树节点
type AVLNode struct {
    Value       interface{}
    Left, Right *AVLNode
    Height      int
}

// 创建新节点
func NewAVLNode(value interface{}) *AVLNode {
    return &AVLNode{
        Value:  value,
        Left:   nil,
        Right:  nil,
        Height: 1, // 新节点高度为1
    }
}

// AVL树
type AVLTree struct {
    Root *AVLNode
    Size int
}

// 创建新AVL树
func NewAVLTree() *AVLTree {
    return &AVLTree{
        Root: nil,
        Size: 0,
    }
}

// 获取节点高度
func getHeight(node *AVLNode) int {
    if node == nil {
        return 0
    }
    return node.Height
}

// 获取平衡因子
func getBalanceFactor(node *AVLNode) int {
    if node == nil {
        return 0
    }
    return getHeight(node.Left) - getHeight(node.Right)
}

// 更新节点高度
func updateHeight(node *AVLNode) {
    if node == nil {
        return
    }
    
    leftHeight := getHeight(node.Left)
    rightHeight := getHeight(node.Right)
    
    if leftHeight > rightHeight {
        node.Height = leftHeight + 1
    } else {
        node.Height = rightHeight + 1
    }
}

// 右旋转
func rightRotate(y *AVLNode) *AVLNode {
    x := y.Left
    T2 := x.Right
    
    // 执行旋转
    x.Right = y
    y.Left = T2
    
    // 更新高度
    updateHeight(y)
    updateHeight(x)
    
    return x
}

// 左旋转
func leftRotate(x *AVLNode) *AVLNode {
    y := x.Right
    T2 := y.Left
    
    // 执行旋转
    y.Left = x
    x.Right = T2
    
    // 更新高度
    updateHeight(x)
    updateHeight(y)
    
    return y
}

// 插入节点
func (t *AVLTree) Insert(value interface{}) {
    t.Root = t.insertNode(t.Root, value)
    t.Size++
}

// 递归插入节点
func (t *AVLTree) insertNode(node *AVLNode, value interface{}) *AVLNode {
    // 标准BST插入
    if node == nil {
        return NewAVLNode(value)
    }
    
    cmp := compare(value, node.Value)
    if cmp < 0 {
        node.Left = t.insertNode(node.Left, value)
    } else if cmp > 0 {
        node.Right = t.insertNode(node.Right, value)
    } else {
        // 值已存在，不做改变
        t.Size-- // 抵消外部的Size++
        return node
    }
    
    // 更新当前节点高度
    updateHeight(node)
    
    // 获取平衡因子
    balance := getBalanceFactor(node)
    
    // 左左情况
    if balance > 1 && compare(value, node.Left.Value) < 0 {
        return rightRotate(node)
    }
    
    // 右右情况
    if balance < -1 && compare(value, node.Right.Value) > 0 {
        return leftRotate(node)
    }
    
    // 左右情况
    if balance > 1 && compare(value, node.Left.Value) > 0 {
        node.Left = leftRotate(node.Left)
        return rightRotate(node)
    }
    
    // 右左情况
    if balance < -1 && compare(value, node.Right.Value) < 0 {
        node.Right = rightRotate(node.Right)
        return leftRotate(node)
    }
    
    return node
}

// 查找最小值节点
func findMinNode(node *AVLNode) *AVLNode {
    current := node
    for current != nil && current.Left != nil {
        current = current.Left
    }
    return current
}

// 删除节点
func (t *AVLTree) Delete(value interface{}) {
    t.Root = t.deleteNode(t.Root, value)
}

// 递归删除节点
func (t *AVLTree) deleteNode(node *AVLNode, value interface{}) *AVLNode {
    if node == nil {
        return nil
    }
    
    cmp := compare(value, node.Value)
    if cmp < 0 {
        node.Left = t.deleteNode(node.Left, value)
    } else if cmp > 0 {
        node.Right = t.deleteNode(node.Right, value)
    } else {
        // 找到要删除的节点
        t.Size--
        
        // 节点有0或1个子节点
        if node.Left == nil {
            return node.Right
        } else if node.Right == nil {
            return node.Left
        }
        
        // 节点有2个子节点，找到中序后继
        successor := findMinNode(node.Right)
        
        // 用后继节点的值替换当前节点
        node.Value = successor.Value
        
        // 删除后继节点
        node.Right = t.deleteNode(node.Right, successor.Value)
        t.Size++ // 抵消上面的Size--，因为实际只删除一个节点
    }
    
    // 更新高度
    updateHeight(node)
    
    // 检查平衡因子
    balance := getBalanceFactor(node)
    
    // 左左情况
    if balance > 1 && getBalanceFactor(node.Left) >= 0 {
        return rightRotate(node)
    }
    
    // 左右情况
    if balance > 1 && getBalanceFactor(node.Left) < 0 {
        node.Left = leftRotate(node.Left)
        return rightRotate(node)
    }
    
    // 右右情况
    if balance < -1 && getBalanceFactor(node.Right) <= 0 {
        return leftRotate(node)
    }
    
    // 右左情况
    if balance < -1 && getBalanceFactor(node.Right) > 0 {
        node.Right = rightRotate(node.Right)
        return leftRotate(node)
    }
    
    return node
}

// 搜索
func (t *AVLTree) Search(value interface{}) bool {
    return t.searchNode(t.Root, value)
}

// 递归搜索节点
func (t *AVLTree) searchNode(node *AVLNode, value interface{}) bool {
    if node == nil {
        return false
    }
    
    cmp := compare(value, node.Value)
    if cmp == 0 {
        return true
    } else if cmp < 0 {
        return t.searchNode(node.Left, value)
    } else {
        return t.searchNode(node.Right, value)
    }
}

// 中序遍历
func (t *AVLTree) InOrderTraversal() []interface{} {
    result := make([]interface{}, 0, t.Size)
    
    var traverse func(node *AVLNode)
    traverse = func(node *AVLNode) {
        if node == nil {
            return
        }
        
        traverse(node.Left)
        result = append(result, node.Value)
        traverse(node.Right)
    }
    
    traverse(t.Root)
    return result
}
```

##### 复杂度分析

**时间复杂度**:

| 操作 | 平均情况 | 最坏情况 |
|------|----------|----------|
| 插入 | $O(\log n)$ | $O(\log n)$ |
| 查找 | $O(\log n)$ | $O(\log n)$ |
| 删除 | $O(\log n)$ | $O(\log n)$ |
| 遍历 | $O(n)$ | $O(n)$ |

**空间复杂度**:

- 存储: $O(n)$
- 递归调用栈: $O(\log n)$

**定理 3.5** (AVL树操作复杂度)
在AVL树中，插入、删除和查找操作的时间复杂度均为 $O(\log n)$。

**证明**:
由于AVL树的高度被限制在 $O(\log n)$，所有基于树高的操作复杂度都是 $O(\log n)$。插入和删除操作可能需要旋转来维持平衡，但每次插入或删除最多需要 $O(\log n)$ 次旋转。

**定理 3.6** (AVL树旋转次数)
在AVL树中，每次插入操作最多需要1次单旋转或1次双旋转（2次单旋转）；每次删除操作最多需要 $O(\log n)$ 次旋转。

**证明**:
插入操作只会改变插入路径上节点的平衡因子，而且最多只有一个节点会变得不平衡，因此最多需要1次单旋转或1次双旋转。
删除操作可能会导致从删除点到根节点的路径上的多个节点变得不平衡，因此可能需要 $O(\log n)$ 次旋转。

#### 3.2.2 其他平衡树

除了AVL树，还有其他几种常见的平衡树结构：

1. **红黑树**: 通过着色机制和旋转操作保持平衡，插入和删除操作的旋转次数比AVL树少。
2. **伸展树**: 不严格保持平衡，但通过将访问过的节点移到根部来提高常用节点的访问效率。
3. **树堆(Treap)**: 结合了二叉搜索树和堆的特性，每个节点有一个随机优先级。
4. **替罪羊树**: 允许树暂时不平衡，但当不平衡度超过阈值时进行重建。

这些平衡树各有优缺点，适用于不同的应用场景。例如，红黑树在插入和删除频繁的场景中比AVL树更高效，而AVL树在查询频繁的场景中表现更好。

### 3.3 B树和B+树

#### 3.3.1 B树

**定义 3.3.1** (B树)
B树是一种自平衡的搜索树，具有以下特性：

1. 每个节点最多有m个子节点
2. 每个内部节点（除根节点外）至少有⌈m/2⌉个子节点
3. 根节点至少有2个子节点（除非树只有一个节点）
4. 所有叶节点都在同一层
5. 每个非叶节点包含n个键和n+1个指针，其中⌈m/2⌉-1 ≤ n ≤ m-1

B树广泛应用于数据库和文件系统中，因为它能有效减少磁盘I/O操作。

#### 3.3.2 B+树

**定义 3.3.2** (B+树)
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

**定义 5.1** (开放寻址法)
开放寻址法是一种解决哈希冲突的方法，当发生冲突时，通过某种探测序列查找下一个可用的位置。

**定义 5.1.1** (线性探测)
线性探测是开放寻址法的一种，当位置 $h(k)$ 被占用时，尝试 $h(k)+1$, $h(k)+2$, ... 直到找到空位。

**定义 5.1.2** (二次探测)
二次探测是开放寻址法的一种，当位置 $h(k)$ 被占用时，尝试 $h(k)+1^2$, $h(k)+2^2$, $h(k)-1^2$, $h(k)-2^2$, ... 直到找到空位。

**定义 5.1.3** (双重哈希)
双重哈希是开放寻址法的一种，使用两个哈希函数 $h_1(k)$ 和 $h_2(k)$，探测序列为 $h_1(k)$, $h_1(k)+h_2(k)$, $h_1(k)+2h_2(k)$, ... 直到找到空位。

```go
// 使用开放寻址法的哈希表
type OpenAddressingHashTable struct {
    size     int
    capacity int
    table    []*KeyValue
    deleted  []bool
}

// 键值对
type KeyValue struct {
    Key   interface{}
    Value interface{}
}

// 创建哈希表
func NewOpenAddressingHashTable(capacity int) *OpenAddressingHashTable {
    return &OpenAddressingHashTable{
        size:     0,
        capacity: capacity,
        table:    make([]*KeyValue, capacity),
        deleted:  make([]bool, capacity),
    }
}

// 哈希函数
func (h *OpenAddressingHashTable) hash(key interface{}) int {
    switch v := key.(type) {
    case int:
        return v % h.capacity
    case string:
        hash := 0
        for i := 0; i < len(v); i++ {
            hash = (hash*31 + int(v[i])) % h.capacity
        }
        return hash
    default:
        return 0
    }
}

// 二次探测
func (h *OpenAddressingHashTable) quadraticProbe(key interface{}) int {
    hash := h.hash(key)
    i := 0
    pos := hash
    
    for h.table[pos] != nil && !h.keyEquals(h.table[pos].Key, key) && !h.deleted[pos] {
        i++
        pos = (hash + i*i) % h.capacity
    }
    
    return pos
}

// 键比较
func (h *OpenAddressingHashTable) keyEquals(k1, k2 interface{}) bool {
    return k1 == k2
}

// 插入键值对
func (h *OpenAddressingHashTable) Put(key, value interface{}) bool {
    // 检查负载因子
    if float64(h.size)/float64(h.capacity) > 0.7 {
        h.resize(h.capacity * 2)
    }
    
    pos := h.quadraticProbe(key)
    
    if h.table[pos] != nil && h.keyEquals(h.table[pos].Key, key) {
        h.table[pos].Value = value
        return true
    }
    
    h.table[pos] = &KeyValue{Key: key, Value: value}
    h.deleted[pos] = false
    h.size++
    return true
}

// 获取值
func (h *OpenAddressingHashTable) Get(key interface{}) (interface{}, bool) {
    pos := h.quadraticProbe(key)
    
    if h.table[pos] != nil && h.keyEquals(h.table[pos].Key, key) && !h.deleted[pos] {
        return h.table[pos].Value, true
    }
    
    return nil, false
}

// 删除键值对
func (h *OpenAddressingHashTable) Remove(key interface{}) bool {
    pos := h.quadraticProbe(key)
    
    if h.table[pos] != nil && h.keyEquals(h.table[pos].Key, key) && !h.deleted[pos] {
        h.deleted[pos] = true
        h.size--
        return true
    }
    
    return false
}

// 重新调整大小
func (h *OpenAddressingHashTable) resize(newCapacity int) {
    oldTable := h.table
    oldDeleted := h.deleted
    
    h.table = make([]*KeyValue, newCapacity)
    h.deleted = make([]bool, newCapacity)
    h.capacity = newCapacity
    h.size = 0
    
    for i, kv := range oldTable {
        if kv != nil && !oldDeleted[i] {
            h.Put(kv.Key, kv.Value)
        }
    }
}
```

**时间复杂度**:

| 操作 | 平均情况 | 最坏情况 |
|------|----------|----------|
| 插入 | $O(1)$ | $O(n)$ |
| 查找 | $O(1)$ | $O(n)$ |
| 删除 | $O(1)$ | $O(n)$ |

**空间复杂度**: $O(n)$

**定理 5.1** (开放寻址法的性能)
当负载因子 $\alpha = n/m < 1$ 时，使用线性探测的开放寻址法中，成功查找的平均探测次数约为 $\frac{1}{2}(1+\frac{1}{1-\alpha})$，不成功查找的平均探测次数约为 $\frac{1}{2}(1+(\frac{1}{1-\alpha})^2)$。

### 5.2 链式散列

**定义 5.2** (链式散列)
链式散列是一种解决哈希冲突的方法，将具有相同哈希值的元素存储在一个链表中。

```go
// 链表节点
type ChainNode struct {
    Key   interface{}
    Value interface{}
    Next  *ChainNode
}

// 链式哈希表
type ChainHashTable struct {
    size     int
    capacity int
    table    []*ChainNode
}

// 创建链式哈希表
func NewChainHashTable(capacity int) *ChainHashTable {
    return &ChainHashTable{
        size:     0,
        capacity: capacity,
        table:    make([]*ChainNode, capacity),
    }
}

// 哈希函数
func (h *ChainHashTable) hash(key interface{}) int {
    switch v := key.(type) {
    case int:
        return v % h.capacity
    case string:
        hash := 0
        for i := 0; i < len(v); i++ {
            hash = (hash*31 + int(v[i])) % h.capacity
        }
        return hash
    default:
        return 0
    }
}

// 键比较
func (h *ChainHashTable) keyEquals(k1, k2 interface{}) bool {
    return k1 == k2
}

// 插入键值对
func (h *ChainHashTable) Put(key, value interface{}) {
    // 检查负载因子
    if float64(h.size)/float64(h.capacity) > 0.75 {
        h.resize(h.capacity * 2)
    }
    
    hash := h.hash(key)
    node := h.table[hash]
    
    // 检查键是否已存在
    for node != nil {
        if h.keyEquals(node.Key, key) {
            node.Value = value
            return
        }
        node = node.Next
    }
    
    // 创建新节点并插入链表头部
    newNode := &ChainNode{Key: key, Value: value, Next: h.table[hash]}
    h.table[hash] = newNode
    h.size++
}

// 获取值
func (h *ChainHashTable) Get(key interface{}) (interface{}, bool) {
    hash := h.hash(key)
    node := h.table[hash]
    
    for node != nil {
        if h.keyEquals(node.Key, key) {
            return node.Value, true
        }
        node = node.Next
    }
    
    return nil, false
}

// 删除键值对
func (h *ChainHashTable) Remove(key interface{}) bool {
    hash := h.hash(key)
    
    if h.table[hash] == nil {
        return false
    }
    
    // 如果是头节点
    if h.keyEquals(h.table[hash].Key, key) {
        h.table[hash] = h.table[hash].Next
        h.size--
        return true
    }
    
    // 查找前一个节点
    prev := h.table[hash]
    for prev.Next != nil && !h.keyEquals(prev.Next.Key, key) {
        prev = prev.Next
    }
    
    if prev.Next != nil {
        prev.Next = prev.Next.Next
        h.size--
        return true
    }
    
    return false
}

// 重新调整大小
func (h *ChainHashTable) resize(newCapacity int) {
    oldTable := h.table
    
    h.table = make([]*ChainNode, newCapacity)
    h.capacity = newCapacity
    h.size = 0
    
    // 重新插入所有键值对
    for _, node := range oldTable {
        for node != nil {
            h.Put(node.Key, node.Value)
            node = node.Next
        }
    }
}
```

**时间复杂度**:

| 操作 | 平均情况 | 最坏情况 |
|------|----------|----------|
| 插入 | $O(1)$ | $O(n)$ |
| 查找 | $O(1)$ | $O(n)$ |
| 删除 | $O(1)$ | $O(n)$ |

**空间复杂度**: $O(n)$

**定理 5.2** (链式散列的性能)
当负载因子 $\alpha = n/m$ 时，链式散列中，成功查找的平均时间复杂度为 $O(1+\alpha)$。

### 5.3 一致性哈希

**定义 5.3** (一致性哈希)
一致性哈希是一种特殊的哈希算法，用于分布式系统中，可以在节点增加或删除时，最小化键的重新映射。

**定理 5.3** (一致性哈希的重新映射)
在一致性哈希中，当添加或删除一个节点时，平均只有 $\frac{K}{n}$ 个键需要重新映射，其中 $K$ 是键的总数，$n$ 是节点数。

```go
// 一致性哈希
type ConsistentHash struct {
    replicas  int           // 每个节点的虚拟节点数
    hashRing  []int         // 有序哈希环
    nodeMap   map[int]string // 哈希值到节点的映射
    hashFunc  func(data []byte) uint32 // 哈希函数
}

// 创建一致性哈希
func NewConsistentHash(replicas int, fn func(data []byte) uint32) *ConsistentHash {
    ch := &ConsistentHash{
        replicas: replicas,
        hashRing: []int{},
        nodeMap:  make(map[int]string),
        hashFunc: fn,
    }
    
    if ch.hashFunc == nil {
        ch.hashFunc = crc32.ChecksumIEEE
    }
    
    return ch
}

// 添加节点
func (ch *ConsistentHash) Add(nodes ...string) {
    for _, node := range nodes {
        // 为每个节点创建多个虚拟节点
        for i := 0; i < ch.replicas; i++ {
            hash := int(ch.hashFunc([]byte(fmt.Sprintf("%s-%d", node, i))))
            ch.hashRing = append(ch.hashRing, hash)
            ch.nodeMap[hash] = node
        }
    }
    
    // 排序哈希环
    sort.Ints(ch.hashRing)
}

// 删除节点
func (ch *ConsistentHash) Remove(node string) {
    for i := 0; i < ch.replicas; i++ {
        hash := int(ch.hashFunc([]byte(fmt.Sprintf("%s-%d", node, i))))
        
        // 删除哈希环中的虚拟节点
        idx := sort.SearchInts(ch.hashRing, hash)
        if idx < len(ch.hashRing) && ch.hashRing[idx] == hash {
            ch.hashRing = append(ch.hashRing[:idx], ch.hashRing[idx+1:]...)
        }
        
        delete(ch.nodeMap, hash)
    }
}

// 获取键对应的节点
func (ch *ConsistentHash) Get(key string) string {
    if len(ch.hashRing) == 0 {
        return ""
    }
    
    hash := int(ch.hashFunc([]byte(key)))
    
    // 二分查找，找到第一个大于等于hash的位置
    idx := sort.Search(len(ch.hashRing), func(i int) bool {
        return ch.hashRing[i] >= hash
    })
    
    // 如果没有找到，则使用第一个节点
    if idx == len(ch.hashRing) {
        idx = 0
    }
    
    return ch.nodeMap[ch.hashRing[idx]]
}
```

**时间复杂度**:

| 操作 | 复杂度 |
|------|--------|
| 添加节点 | $O(r \log n)$ |
| 删除节点 | $O(r \log n)$ |
| 查找节点 | $O(\log n)$ |

其中，$r$ 是每个节点的虚拟节点数，$n$ 是节点总数。

### 5.4 完美哈希

**定义 5.4** (完美哈希)
完美哈希是一种特殊的哈希函数，对于一组已知的键集合，它能够无冲突地将每个键映射到不同的位置。

**定义 5.4.1** (最小完美哈希)
最小完美哈希是一种完美哈希，它将 $n$ 个键映射到 $[0, n-1]$ 的范围内，且每个位置恰好有一个键。

```go
// 完美哈希表的第一级哈希函数
type PerfectHashTable struct {
    size     int
    g        []int      // 第一级哈希函数的参数
    secondLevel []*SecondLevelTable // 第二级哈希表
    keys     []string   // 存储的键集合
}

// 第二级哈希表
type SecondLevelTable struct {
    size int
    g    []int    // 第二级哈希函数的参数
    table []string // 存储的键
}

// 创建完美哈希表
func NewPerfectHashTable(keys []string) *PerfectHashTable {
    n := len(keys)
    m := nextPrime(n) // 选择一个合适的素数作为表大小
    
    pht := &PerfectHashTable{
        size:       m,
        g:          make([]int, m),
        secondLevel: make([]*SecondLevelTable, m),
        keys:       keys,
    }
    
    // 构建第一级哈希表
    buckets := make([][]string, m)
    for _, key := range keys {
        h := pht.firstHash(key, 0) % m
        buckets[h] = append(buckets[h], key)
    }
    
    // 为每个桶构建第二级哈希表
    for i, bucket := range buckets {
        if len(bucket) > 0 {
            // 尝试不同的g值，直到找到一个完美哈希
            for g := 1; ; g++ {
                if pht.trySecondLevel(i, bucket, g) {
                    break
                }
            }
        }
    }
    
    return pht
}

// 尝试为桶构建第二级哈希表
func (pht *PerfectHashTable) trySecondLevel(bucketIndex int, bucket []string, g int) bool {
    m := nextPrime(len(bucket) * len(bucket)) // 第二级表大小
    
    occupied := make([]bool, m)
    pht.g[bucketIndex] = g
    
    for _, key := range bucket {
        h := pht.secondHash(key, g) % m
        if occupied[h] {
            return false // 发生冲突，尝试下一个g值
        }
        occupied[h] = true
    }
    
    // 创建第二级表
    secondTable := &SecondLevelTable{
        size:  m,
        g:     []int{g},
        table: make([]string, m),
    }
    
    // 填充第二级表
    for _, key := range bucket {
        h := pht.secondHash(key, g) % m
        secondTable.table[h] = key
    }
    
    pht.secondLevel[bucketIndex] = secondTable
    return true
}

// 第一级哈希函数
func (pht *PerfectHashTable) firstHash(key string, g int) int {
    hash := 0
    for i := 0; i < len(key); i++ {
        hash = (hash*33 + int(key[i])) % pht.size
    }
    return (hash + g) % pht.size
}

// 第二级哈希函数
func (pht *PerfectHashTable) secondHash(key string, g int) int {
    hash := 0
    for i := 0; i < len(key); i++ {
        hash = (hash*37 + int(key[i])) % 1000000007
    }
    return (hash + g) % 1000000007
}

// 查找键
func (pht *PerfectHashTable) Contains(key string) bool {
    h1 := pht.firstHash(key, 0) % pht.size
    
    if pht.secondLevel[h1] == nil {
        return false
    }
    
    secondTable := pht.secondLevel[h1]
    h2 := pht.secondHash(key, pht.g[h1]) % secondTable.size
    
    return secondTable.table[h2] == key
}

// 获取下一个素数
func nextPrime(n int) int {
    if n <= 2 {
        return 2
    }
    
    prime := n
    if prime%2 == 0 {
        prime++
    }
    
    for {
        if isPrime(prime) {
            return prime
        }
        prime += 2
    }
}

// 判断是否为素数
func isPrime(n int) bool {
    if n <= 1 {
        return false
    }
    if n <= 3 {
        return true
    }
    if n%2 == 0 || n%3 == 0 {
        return false
    }
    
    i := 5
    for i*i <= n {
        if n%i == 0 || n%(i+2) == 0 {
            return false
        }
        i += 6
    }
    
    return true
}
```

**时间复杂度**:

| 操作 | 复杂度 |
|------|--------|
| 构建 | $O(n)$ (期望) |
| 查找 | $O(1)$ (最坏情况) |

**空间复杂度**: $O(n)$

**定理 5.4** (完美哈希的空间效率)
对于 $n$ 个键，最小完美哈希函数的理论下界空间复杂度为 $\Omega(n)$ 比特。

### 5.5 布谷鸟哈希

**定义 5.5** (布谷鸟哈希)
布谷鸟哈希是一种解决冲突的方法，使用两个或多个哈希函数，当发生冲突时，将已有元素踢出，并重新插入到另一个位置。

```go
// 布谷鸟哈希表
type CuckooHashTable struct {
    size     int
    table1   []*KeyValue
    table2   []*KeyValue
    hashFunc1 func(interface{}) int
    hashFunc2 func(interface{}) int
    maxLoop  int // 最大重试次数
}

// 创建布谷鸟哈希表
func NewCuckooHashTable(size int) *CuckooHashTable {
    return &CuckooHashTable{
        size:     size,
        table1:   make([]*KeyValue, size),
        table2:   make([]*KeyValue, size),
        hashFunc1: func(key interface{}) int {
            h := 0
            switch v := key.(type) {
            case int:
                h = v
            case string:
                for i := 0; i < len(v); i++ {
                    h = (h*31 + int(v[i])) % size
                }
            }
            return h % size
        },
        hashFunc2: func(key interface{}) int {
            h := 0
            switch v := key.(type) {
            case int:
                h = v*16777619
            case string:
                for i := 0; i < len(v); i++ {
                    h = (h*37 + int(v[i])) % size
                }
            }
            return h % size
        },
        maxLoop:  100,
    }
}

// 插入键值对
func (h *CuckooHashTable) Put(key, value interface{}) bool {
    // 检查键是否已存在
    if kv, found := h.Get(key); found {
        kv.Value = value
        return true
    }
    
    kv := &KeyValue{Key: key, Value: value}
    
    // 尝试插入
    for i := 0; i < h.maxLoop; i++ {
        pos1 := h.hashFunc1(key)
        
        // 尝试插入表1
        if h.table1[pos1] == nil {
            h.table1[pos1] = kv
            return true
        }
        
        // 交换并尝试插入表2
        h.table1[pos1], kv = kv, h.table1[pos1]
        
        pos2 := h.hashFunc2(kv.Key)
        if h.table2[pos2] == nil {
            h.table2[pos2] = kv
            return true
        }
        
        // 交换并继续尝试表1
        h.table2[pos2], kv = kv, h.table2[pos2]
    }
    
    // 达到最大重试次数，需要重新哈希
    h.rehash()
    return h.Put(key, value)
}

// 获取值
func (h *CuckooHashTable) Get(key interface{}) (*KeyValue, bool) {
    pos1 := h.hashFunc1(key)
    if h.table1[pos1] != nil && h.keyEquals(h.table1[pos1].Key, key) {
        return h.table1[pos1], true
    }
    
    pos2 := h.hashFunc2(key)
    if h.table2[pos2] != nil && h.keyEquals(h.table2[pos2].Key, key) {
        return h.table2[pos2], true
    }
    
    return nil, false
}

// 删除键值对
func (h *CuckooHashTable) Remove(key interface{}) bool {
    pos1 := h.hashFunc1(key)
    if h.table1[pos1] != nil && h.keyEquals(h.table1[pos1].Key, key) {
        h.table1[pos1] = nil
        return true
    }
    
    pos2 := h.hashFunc2(key)
    if h.table2[pos2] != nil && h.keyEquals(h.table2[pos2].Key, key) {
        h.table2[pos2] = nil
        return true
    }
    
    return false
}

// 键比较
func (h *CuckooHashTable) keyEquals(k1, k2 interface{}) bool {
    return k1 == k2
}

// 重新哈希
func (h *CuckooHashTable) rehash() {
    oldTable1 := h.table1
    oldTable2 := h.table2
    newSize := h.size * 2
    
    h.size = newSize
    h.table1 = make([]*KeyValue, newSize)
    h.table2 = make([]*KeyValue, newSize)
    
    // 重新插入所有键值对
    for _, kv := range oldTable1 {
        if kv != nil {
            h.Put(kv.Key, kv.Value)
        }
    }
    
    for _, kv := range oldTable2 {
        if kv != nil {
            h.Put(kv.Key, kv.Value)
        }
    }
}
```

**时间复杂度**:

| 操作 | 平均情况 | 最坏情况 |
|------|----------|----------|
| 插入 | $O(1)$ | $O(n)$ (需要重哈希) |
| 查找 | $O(1)$ | $O(1)$ |
| 删除 | $O(1)$ | $O(1)$ |

**空间复杂度**: $O(n)$

**定理 5.5** (布谷鸟哈希的负载因子)
当使用两个哈希函数时，布谷鸟哈希的最大负载因子约为 50%。使用三个或更多哈希函数可以提高负载因子。

## 特殊数据结构

### 6.1 布隆过滤器

**定义 6.1** (布隆过滤器)
布隆过滤器是一种空间效率很高的概率性数据结构，用于判断一个元素是否在集合中。它可能会产生误判（假阳性），但不会漏判（假阴性）。

**定理 6.1** (布隆过滤器的误判率)
对于使用 $k$ 个哈希函数和 $m$ 比特的布隆过滤器，当存储 $n$ 个元素时，其误判率约为 $p = (1-e^{-kn/m})^k$。

```go
// 布隆过滤器
type BloomFilter struct {
    bitset    []bool
    size      int
    hashFuncs []func(string) int
}

// 创建布隆过滤器
func NewBloomFilter(size int, numHashes int) *BloomFilter {
    bf := &BloomFilter{
        bitset:    make([]bool, size),
        size:      size,
        hashFuncs: make([]func(string) int, numHashes),
    }
    
    // 创建不同的哈希函数
    for i := 0; i < numHashes; i++ {
        seed := i + 1
        bf.hashFuncs[i] = func(s string) int {
            hash := 0
            for j := 0; j < len(s); j++ {
                hash = (hash*seed + int(s[j])) % size
            }
            return hash
        }
    }
    
    return bf
}

// 添加元素
func (bf *BloomFilter) Add(s string) {
    for _, hashFunc := range bf.hashFuncs {
        position := hashFunc(s)
        bf.bitset[position] = true
    }
}

// 检查元素是否可能存在
func (bf *BloomFilter) Contains(s string) bool {
    for _, hashFunc := range bf.hashFuncs {
        position := hashFunc(s)
        if !bf.bitset[position] {
            return false // 肯定不存在
        }
    }
    return true // 可能存在
}

// 计算最佳哈希函数数量
func OptimalNumHashes(n, m int) int {
    // n是预期元素数量，m是比特数组大小
    return int(math.Ceil(float64(m) / float64(n) * math.Log(2)))
}

// 计算给定误判率所需的比特数组大小
func OptimalSize(n int, p float64) int {
    // n是预期元素数量，p是期望的误判率
    return int(math.Ceil(-float64(n) * math.Log(p) / math.Pow(math.Log(2), 2)))
}
```

**时间复杂度**:

| 操作 | 复杂度 |
|------|--------|
| 添加 | $O(k)$ |
| 查询 | $O(k)$ |

其中，$k$ 是哈希函数的数量。

**空间复杂度**: $O(m)$，其中 $m$ 是比特数组的大小。

**应用场景**:
- 网络爬虫的URL去重
- 垃圾邮件过滤
- 缓存穿透防止
- 大规模集合的快速查询

### 6.2 跳表

**定义 6.2** (跳表)
跳表是一种随机化的数据结构，基于并联的有序链表，其效率可比拟于平衡树，但实现更为简单。

**定理 6.2** (跳表的高度)
包含 $n$ 个元素的跳表的期望高度为 $O(\log n)$。

```go
// 跳表节点
type SkipListNode struct {
    Value    int
    Forward  []*SkipListNode // 不同层级的前向指针
}

// 跳表
type SkipList struct {
    Head     *SkipListNode
    MaxLevel int
    Level    int
    P        float64 // 层级提升概率
}

// 创建跳表
func NewSkipList(maxLevel int) *SkipList {
    return &SkipList{
        Head:     &SkipListNode{Forward: make([]*SkipListNode, maxLevel)},
        MaxLevel: maxLevel,
        Level:    0,
        P:        0.5,
    }
}

// 随机生成层级
func (sl *SkipList) randomLevel() int {
    level := 0
    for rand.Float64() < sl.P && level < sl.MaxLevel-1 {
        level++
    }
    return level
}

// 搜索元素
func (sl *SkipList) Search(value int) bool {
    current := sl.Head
    
    // 从最高层开始向下搜索
    for i := sl.Level; i >= 0; i-- {
        // 在当前层向前移动，直到找到大于等于目标值的节点
        for current.Forward[i] != nil && current.Forward[i].Value < value {
            current = current.Forward[i]
        }
    }
    
    // 现在current.Forward[0]是第一个大于等于value的节点
    current = current.Forward[0]
    
    return current != nil && current.Value == value
}

// 插入元素
func (sl *SkipList) Insert(value int) {
    // 保存每一层的前驱节点
    update := make([]*SkipListNode, sl.MaxLevel)
    current := sl.Head
    
    // 从最高层开始向下查找插入位置
    for i := sl.Level; i >= 0; i-- {
        for current.Forward[i] != nil && current.Forward[i].Value < value {
            current = current.Forward[i]
        }
        update[i] = current
    }
    
    // 获取随机层级
    level := sl.randomLevel()
    
    // 更新跳表的最大层级
    if level > sl.Level {
        for i := sl.Level + 1; i <= level; i++ {
            update[i] = sl.Head
        }
        sl.Level = level
    }
    
    // 创建新节点
    newNode := &SkipListNode{
        Value:   value,
        Forward: make([]*SkipListNode, level+1),
    }
    
    // 更新指针
    for i := 0; i <= level; i++ {
        newNode.Forward[i] = update[i].Forward[i]
        update[i].Forward[i] = newNode
    }
}

// 删除元素
func (sl *SkipList) Delete(value int) bool {
    // 保存每一层的前驱节点
    update := make([]*SkipListNode, sl.MaxLevel)
    current := sl.Head
    
    // 从最高层开始向下查找删除位置
    for i := sl.Level; i >= 0; i-- {
        for current.Forward[i] != nil && current.Forward[i].Value < value {
            current = current.Forward[i]
        }
        update[i] = current
    }
    
    // 找到要删除的节点
    current = current.Forward[0]
    
    // 如果节点存在且值匹配
    if current != nil && current.Value == value {
        // 从最底层开始向上删除
        for i := 0; i <= sl.Level; i++ {
            if update[i].Forward[i] != current {
                break
            }
            update[i].Forward[i] = current.Forward[i]
        }
        
        // 更新跳表的最大层级
        for sl.Level > 0 && sl.Head.Forward[sl.Level] == nil {
            sl.Level--
        }
        
        return true
    }
    
    return false
}
```

**时间复杂度**:

| 操作 | 平均情况 | 最坏情况 |
|------|----------|----------|
| 搜索 | $O(\log n)$ | $O(n)$ |
| 插入 | $O(\log n)$ | $O(n)$ |
| 删除 | $O(\log n)$ | $O(n)$ |

**空间复杂度**: $O(n)$

**应用场景**:
- Redis中的有序集合
- 高效的范围查询
- 作为平衡树的替代

### 6.3 字典树变体

#### 6.3.1 基本字典树

**定义 6.3** (字典树)
字典树，又称前缀树或Trie树，是一种树形数据结构，用于高效地存储和检索字符串数据集中的键。

```go
// 字典树节点
type TrieNode struct {
    Children map[rune]*TrieNode
    IsEnd    bool
}

// 字典树
type Trie struct {
    Root *TrieNode
}

// 创建字典树
func NewTrie() *Trie {
    return &Trie{
        Root: &TrieNode{
            Children: make(map[rune]*TrieNode),
            IsEnd:    false,
        },
    }
}

// 插入单词
func (t *Trie) Insert(word string) {
    node := t.Root
    
    for _, ch := range word {
        if _, exists := node.Children[ch]; !exists {
            node.Children[ch] = &TrieNode{
                Children: make(map[rune]*TrieNode),
                IsEnd:    false,
            }
        }
        node = node.Children[ch]
    }
    
    node.IsEnd = true
}

// 搜索单词
func (t *Trie) Search(word string) bool {
    node := t.Root
    
    for _, ch := range word {
        if _, exists := node.Children[ch]; !exists {
            return false
        }
        node = node.Children[ch]
    }
    
    return node.IsEnd
}

// 检查前缀
func (t *Trie) StartsWith(prefix string) bool {
    node := t.Root
    
    for _, ch := range prefix {
        if _, exists := node.Children[ch]; !exists {
            return false
        }
        node = node.Children[ch]
    }
    
    return true
}
```

#### 6.3.2 压缩前缀树

**定义 6.3.2** (压缩前缀树)
压缩前缀树是字典树的一种优化变体，通过合并只有一个子节点的节点来减少内存使用。

```go
// 压缩前缀树节点
type CompressedTrieNode struct {
    Prefix   string
    Children map[rune]*CompressedTrieNode
    IsEnd    bool
}

// 压缩前缀树
type CompressedTrie struct {
    Root *CompressedTrieNode
}

// 创建压缩前缀树
func NewCompressedTrie() *CompressedTrie {
    return &CompressedTrie{
        Root: &CompressedTrieNode{
            Prefix:   "",
            Children: make(map[rune]*CompressedTrieNode),
            IsEnd:    false,
        },
    }
}

// 插入单词
func (t *CompressedTrie) Insert(word string) {
    t.insertRecursive(t.Root, word)
}

// 递归插入
func (t *CompressedTrie) insertRecursive(node *CompressedTrieNode, word string) {
    // 如果单词为空，标记当前节点为单词结束
    if len(word) == 0 {
        node.IsEnd = true
        return
    }
    
    // 获取第一个字符
    firstChar := rune(word[0])
    
    // 如果存在以该字符开头的子节点
    if child, exists := node.Children[firstChar]; exists {
        i := 0
        prefixLen := min(len(word), len(child.Prefix))
        
        // 找到共同前缀的长度
        for i < prefixLen && word[i] == child.Prefix[i] {
            i++
        }
        
        // 如果共同前缀长度小于子节点的前缀长度，需要拆分子节点
        if i < len(child.Prefix) {
            // 创建新的中间节点
            newNode := &CompressedTrieNode{
                Prefix:   child.Prefix[:i],
                Children: make(map[rune]*CompressedTrieNode),
                IsEnd:    false,
            }
            
            // 更新原子节点
            child.Prefix = child.Prefix[i:]
            newNode.Children[rune(child.Prefix[0])] = child
            
            // 更新父节点指向新节点
            node.Children[firstChar] = newNode
            
            // 如果还有剩余单词，继续插入
            if i < len(word) {
                t.insertRecursive(newNode, word[i:])
            } else {
                newNode.IsEnd = true
            }
        } else {
            // 共同前缀等于子节点前缀，继续在子节点中插入剩余部分
            if i < len(word) {
                t.insertRecursive(child, word[i:])
            } else {
                child.IsEnd = true
            }
        }
    } else {
        // 创建新节点
        node.Children[firstChar] = &CompressedTrieNode{
            Prefix:   word,
            Children: make(map[rune]*CompressedTrieNode),
            IsEnd:    true,
        }
    }
}

// 搜索单词
func (t *CompressedTrie) Search(word string) bool {
    return t.searchPrefix(word, true)
}

// 检查前缀
func (t *CompressedTrie) StartsWith(prefix string) bool {
    return t.searchPrefix(prefix, false)
}

// 搜索前缀
func (t *CompressedTrie) searchPrefix(prefix string, isWord bool) bool {
    node := t.Root
    
    for len(prefix) > 0 {
        firstChar := rune(prefix[0])
        
        if child, exists := node.Children[firstChar]; exists {
            if strings.HasPrefix(prefix, child.Prefix) {
                // 前缀匹配，继续搜索
                prefix = prefix[len(child.Prefix):]
                node = child
            } else {
                // 前缀不匹配
                return false
            }
        } else {
            // 没有匹配的子节点
            return false
        }
    }
    
    // 如果是搜索单词，检查是否是单词结束
    if isWord {
        return node.IsEnd
    }
    
    return true
}

// 辅助函数
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
```

#### 6.3.3 后缀树

**定义 6.3.3** (后缀树)
后缀树是一种特殊的压缩前缀树，包含给定字符串的所有后缀。

**应用场景**:
- 字符串匹配
- 自动补全
- 拼写检查
- 文本压缩
- 生物信息学中的DNA序列匹配

### 6.4 前缀和与区间树

#### 6.4.1 前缀和数组

**定义 6.4.1** (前缀和)
前缀和是一种预处理技术，用于快速计算数组的区间和。

```go
// 构建前缀和数组
func BuildPrefixSum(arr []int) []int {
    n := len(arr)
    prefixSum := make([]int, n+1)
    
    for i := 0; i < n; i++ {
        prefixSum[i+1] = prefixSum[i] + arr[i]
    }
    
    return prefixSum
}

// 查询区间和 [left, right]
func QueryRange(prefixSum []int, left, right int) int {
    return prefixSum[right+1] - prefixSum[left]
}
```

**时间复杂度**:

| 操作 | 复杂度 |
|------|--------|
| 构建 | $O(n)$ |
| 查询 | $O(1)$ |

#### 6.4.2 线段树

**定义 6.4.2** (线段树)
线段树是一种用于解决区间查询问题的树形数据结构，每个节点代表一个区间。

```go
// 线段树节点
type SegmentTreeNode struct {
    Start, End int
    Sum        int
    Left, Right *SegmentTreeNode
}

// 构建线段树
func BuildSegmentTree(arr []int, start, end int) *SegmentTreeNode {
    if start > end {
        return nil
    }
    
    root := &SegmentTreeNode{
        Start: start,
        End:   end,
    }
    
    if start == end {
        root.Sum = arr[start]
    } else {
        mid := start + (end-start)/2
        root.Left = BuildSegmentTree(arr, start, mid)
        root.Right = BuildSegmentTree(arr, mid+1, end)
        root.Sum = root.Left.Sum + root.Right.Sum
    }
    
    return root
}

// 查询区间和 [start, end]
func (root *SegmentTreeNode) Query(start, end int) int {
    // 区间不重叠
    if root.End < start || root.Start > end {
        return 0
    }
    
    // 当前节点区间被查询区间包含
    if start <= root.Start && end >= root.End {
        return root.Sum
    }
    
    // 查询区间与当前节点区间部分重叠，递归查询
    return root.Left.Query(start, end) + root.Right.Query(start, end)
}

// 更新单个元素
func (root *SegmentTreeNode) Update(index, val int) {
    if root.Start == root.End {
        root.Sum = val
        return
    }
    
    mid := root.Start + (root.End-root.Start)/2
    if index <= mid {
        root.Left.Update(index, val)
    } else {
        root.Right.Update(index, val)
    }
    
    root.Sum = root.Left.Sum + root.Right.Sum
}
```

**时间复杂度**:

| 操作 | 复杂度 |
|------|--------|
| 构建 | $O(n)$ |
| 查询 | $O(\log n)$ |
| 更新 | $O(\log n)$ |

#### 6.4.3 树状数组

**定义 6.4.3** (树状数组)
树状数组，又称Binary Indexed Tree或Fenwick Tree，是一种支持高效更新和查询前缀和的数据结构。

```go
// 树状数组
type BinaryIndexedTree struct {
    Tree []int
    Size int
}

// 创建树状数组
func NewBinaryIndexedTree(size int) *BinaryIndexedTree {
    return &BinaryIndexedTree{
        Tree: make([]int, size+1),
        Size: size,
    }
}

// 更新操作
func (bit *BinaryIndexedTree) Update(index, val int) {
    index++ // 转为1-indexed
    for index <= bit.Size {
        bit.Tree[index] += val
        index += index & -index // 加上最低位的1
    }
}

// 查询前缀和 [1, index]
func (bit *BinaryIndexedTree) Query(index int) int {
    index++ // 转为1-indexed
    sum := 0
    for index > 0 {
        sum += bit.Tree[index]
        index -= index & -index // 减去最低位的1
    }
    return sum
}

// 查询区间和 [left, right]
func (bit *BinaryIndexedTree) QueryRange(left, right int) int {
    return bit.Query(right) - bit.Query(left-1)
}

// 从数组构建
func (bit *BinaryIndexedTree) BuildFromArray(arr []int) {
    for i, val := range arr {
        bit.Update(i, val)
    }
}
```

**时间复杂度**:

| 操作 | 复杂度 |
|------|--------|
| 构建 | $O(n \log n)$ |
| 查询 | $O(\log n)$ |
| 更新 | $O(\log n)$ |

**应用场景**:
- 区间求和
- 区间最大/最小值查询
- 动态更新的区间查询问题

### 6.5 稀疏表与RMQ

**定义 6.5** (稀疏表)
稀疏表是一种用于解决静态区间查询问题的数据结构，特别适用于区间最大值、最小值查询(RMQ)。

```go
// 稀疏表
type SparseTable struct {
    ST    [][]int // ST[i][j]表示从i开始，长度为2^j的区间的最小值
    Log   []int   // Log[i]表示log2(i)的向下取整
    Arr   []int   // 原始数组
}

// 创建稀疏表
func NewSparseTable(arr []int) *SparseTable {
    n := len(arr)
    logN := int(math.Log2(float64(n))) + 1
    
    // 初始化稀疏表
    st := make([][]int, n)
    for i := range st {
        st[i] = make([]int, logN)
    }
    
    // 初始化对数表
    log := make([]int, n+1)
    for i := 2; i <= n; i++ {
        log[i] = log[i/2] + 1
    }
    
    // 填充稀疏表
    for i := 0; i < n; i++ {
        st[i][0] = arr[i]
    }
    
    for j := 1; (1 << j) <= n; j++ {
        for i := 0; i + (1 << j) - 1 < n; i++ {
            st[i][j] = min(st[i][j-1], st[i+(1<<(j-1))][j-1])
        }
    }
    
    return &SparseTable{
        ST:  st,
        Log: log,
        Arr: arr,
    }
}

// 查询区间最小值 [l, r]
func (st *SparseTable) Query(l, r int) int {
    length := r - l + 1
    j := st.Log[length]
    return min(st.ST[l][j], st.ST[r-(1<<j)+1][j])
}

// 最小值
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
```

**时间复杂度**:

| 操作 | 复杂度 |
|------|--------|
| 构建 | $O(n \log n)$ |
| 查询 | $O(1)$ |

**应用场景**:
- 区间最值查询
- 静态数据的快速查询
- 最近公共祖先(LCA)问题

## 空间效率数据结构

### 7.1 位图

**定义 7.1** (位图)
位图是一种使用位来表示集合的数据结构，每个元素使用一个位来表示其是否存在，非常适合表示稠密集合。

```go
// 位图
type BitMap struct {
    bits []uint64
    size int
}

// 创建位图
func NewBitMap(size int) *BitMap {
    return &BitMap{
        bits: make([]uint64, (size+63)/64), // 向上取整到64的倍数
        size: size,
    }
}

// 设置位
func (b *BitMap) Set(pos int) {
    if pos < 0 || pos >= b.size {
        return
    }
    b.bits[pos/64] |= 1 << (pos % 64)
}

// 清除位
func (b *BitMap) Clear(pos int) {
    if pos < 0 || pos >= b.size {
        return
    }
    b.bits[pos/64] &= ^(1 << (pos % 64))
}

// 测试位
func (b *BitMap) Test(pos int) bool {
    if pos < 0 || pos >= b.size {
        return false
    }
    return (b.bits[pos/64] & (1 << (pos % 64))) != 0
}

// 翻转位
func (b *BitMap) Flip(pos int) {
    if pos < 0 || pos >= b.size {
        return
    }
    b.bits[pos/64] ^= 1 << (pos % 64)
}

// 计算设置的位数
func (b *BitMap) Count() int {
    count := 0
    for _, word := range b.bits {
        // 使用位操作计算一个64位整数中1的个数
        v := word
        v = v - ((v >> 1) & 0x5555555555555555)
        v = (v & 0x3333333333333333) + ((v >> 2) & 0x3333333333333333)
        count += int((((v + (v >> 4)) & 0xF0F0F0F0F0F0F0F) * 0x101010101010101) >> 56)
    }
    return count
}
```

**时间复杂度**:

| 操作 | 复杂度 |
|------|--------|
| 设置/清除/测试/翻转 | $O(1)$ |
| 计数 | $O(n/64)$ |

**空间复杂度**: $O(n/64)$ 字节，其中 $n$ 是元素数量。

**应用场景**:
- 表示稠密集合
- 布隆过滤器的底层实现
- 位向量
- 存储大量布尔值

### 7.2 基数树

**定义 7.2** (基数树)
基数树是一种特殊的前缀树，用于存储整数或字符串，其中每个节点代表一个位或字符。

```go
// 基数树节点
type RadixTreeNode struct {
    Children map[byte]*RadixTreeNode
    IsEnd    bool
    Value    interface{}
}

// 基数树
type RadixTree struct {
    Root *RadixTreeNode
}

// 创建基数树
func NewRadixTree() *RadixTree {
    return &RadixTree{
        Root: &RadixTreeNode{
            Children: make(map[byte]*RadixTreeNode),
            IsEnd:    false,
        },
    }
}

// 插入键值对
func (rt *RadixTree) Insert(key string, value interface{}) {
    node := rt.Root
    
    // 查找公共前缀
    i := 0
    for i < len(key) {
        if child, exists := node.Children[key[i]]; exists {
            node = child
            i++
        } else {
            break
        }
    }
    
    // 插入剩余的键
    for i < len(key) {
        newNode := &RadixTreeNode{
            Children: make(map[byte]*RadixTreeNode),
            IsEnd:    false,
        }
        node.Children[key[i]] = newNode
        node = newNode
        i++
    }
    
    node.IsEnd = true
    node.Value = value
}

// 搜索键
func (rt *RadixTree) Search(key string) (interface{}, bool) {
    node := rt.Root
    
    for i := 0; i < len(key); i++ {
        if child, exists := node.Children[key[i]]; exists {
            node = child
        } else {
            return nil, false
        }
    }
    
    if node.IsEnd {
        return node.Value, true
    }
    
    return nil, false
}

// 前缀搜索
func (rt *RadixTree) SearchPrefix(prefix string) []string {
    node := rt.Root
    
    // 查找前缀节点
    for i := 0; i < len(prefix); i++ {
        if child, exists := node.Children[prefix[i]]; exists {
            node = child
        } else {
            return nil
        }
    }
    
    // 收集所有以该前缀开头的键
    var result []string
    rt.collectKeys(node, prefix, &result)
    return result
}

// 收集键
func (rt *RadixTree) collectKeys(node *RadixTreeNode, prefix string, result *[]string) {
    if node.IsEnd {
        *result = append(*result, prefix)
    }
    
    for ch, child := range node.Children {
        rt.collectKeys(child, prefix+string(ch), result)
    }
}
```

**时间复杂度**:

| 操作 | 复杂度 |
|------|--------|
| 插入 | $O(k)$ |
| 查找 | $O(k)$ |
| 前缀搜索 | $O(n)$ |

其中，$k$ 是键的长度，$n$ 是匹配前缀的键的数量。

**空间复杂度**: $O(n \cdot k)$，其中 $n$ 是键的数量，$k$ 是平均键长度。

### 7.3 压缩前缀树

**定义 7.3** (压缩前缀树)
压缩前缀树是一种优化的前缀树，通过合并只有一个子节点的节点来减少空间使用。

```go
// 压缩前缀树节点
type CompactPrefixTreeNode struct {
    Prefix   string
    Children map[byte]*CompactPrefixTreeNode
    IsEnd    bool
    Value    interface{}
}

// 压缩前缀树
type CompactPrefixTree struct {
    Root *CompactPrefixTreeNode
}

// 创建压缩前缀树
func NewCompactPrefixTree() *CompactPrefixTree {
    return &CompactPrefixTree{
        Root: &CompactPrefixTreeNode{
            Prefix:   "",
            Children: make(map[byte]*CompactPrefixTreeNode),
            IsEnd:    false,
        },
    }
}

// 查找公共前缀长度
func findCommonPrefixLength(a, b string) int {
    minLen := min(len(a), len(b))
    for i := 0; i < minLen; i++ {
        if a[i] != b[i] {
            return i
        }
    }
    return minLen
}

// 插入键值对
func (cpt *CompactPrefixTree) Insert(key string, value interface{}) {
    if len(key) == 0 {
        cpt.Root.IsEnd = true
        cpt.Root.Value = value
        return
    }
    
    node := cpt.Root
    
    for {
        // 查找匹配的子节点
        var matchingChild *CompactPrefixTreeNode
        var firstChar byte
        
        if len(key) > 0 {
            firstChar = key[0]
            if child, exists := node.Children[firstChar]; exists {
                matchingChild = child
            }
        }
        
        if matchingChild == nil {
            // 没有匹配的子节点，创建新节点
            newNode := &CompactPrefixTreeNode{
                Prefix:   key,
                Children: make(map[byte]*CompactPrefixTreeNode),
                IsEnd:    true,
                Value:    value,
            }
            node.Children[firstChar] = newNode
            return
        }
        
        // 找到匹配的子节点，计算公共前缀长度
        prefixLen := findCommonPrefixLength(key, matchingChild.Prefix)
        
        if prefixLen == len(matchingChild.Prefix) {
            // 子节点的前缀是键的前缀
            if prefixLen == len(key) {
                // 键与子节点前缀完全匹配
                matchingChild.IsEnd = true
                matchingChild.Value = value
                return
            }
            
            // 继续在子节点中查找
            node = matchingChild
            key = key[prefixLen:]
            continue
        }
        
        // 需要分割子节点
        splitNode := &CompactPrefixTreeNode{
            Prefix:   matchingChild.Prefix[:prefixLen],
            Children: make(map[byte]*CompactPrefixTreeNode),
            IsEnd:    prefixLen == len(key),
        }
        
        if splitNode.IsEnd {
            splitNode.Value = value
        }
        
        // 更新原子节点
        matchingChild.Prefix = matchingChild.Prefix[prefixLen:]
        splitNode.Children[matchingChild.Prefix[0]] = matchingChild
        
        // 添加新键的剩余部分
        if prefixLen < len(key) {
            newNode := &CompactPrefixTreeNode{
                Prefix:   key[prefixLen:],
                Children: make(map[byte]*CompactPrefixTreeNode),
                IsEnd:    true,
                Value:    value,
            }
            splitNode.Children[key[prefixLen]] = newNode
        }
        
        // 更新父节点
        node.Children[firstChar] = splitNode
        return
    }
}

// 搜索键
func (cpt *CompactPrefixTree) Search(key string) (interface{}, bool) {
    node := cpt.Root
    
    if len(key) == 0 {
        return node.Value, node.IsEnd
    }
    
    for len(key) > 0 {
        firstChar := key[0]
        child, exists := node.Children[firstChar]
        
        if !exists {
            return nil, false
        }
        
        if len(key) < len(child.Prefix) {
            return nil, false
        }
        
        if key[:len(child.Prefix)] != child.Prefix {
            return nil, false
        }
        
        key = key[len(child.Prefix):]
        node = child
        
        if len(key) == 0 && node.IsEnd {
            return node.Value, true
        }
    }
    
    return nil, false
}

// 前缀搜索
func (cpt *CompactPrefixTree) SearchPrefix(prefix string) []string {
    node := cpt.Root
    currentPrefix := ""
    
    // 查找前缀节点
    for len(prefix) > 0 {
        firstChar := prefix[0]
        child, exists := node.Children[firstChar]
        
        if !exists {
            return nil
        }
        
        if len(prefix) < len(child.Prefix) {
            if prefix == child.Prefix[:len(prefix)] {
                // 前缀是子节点前缀的前缀
                currentPrefix += child.Prefix
                node = child
                prefix = ""
                break
            }
            return nil
        }
        
        if prefix[:len(child.Prefix)] != child.Prefix {
            return nil
        }
        
        currentPrefix += child.Prefix
        prefix = prefix[len(child.Prefix):]
        node = child
    }
    
    // 收集所有以该前缀开头的键
    var result []string
    if node.IsEnd {
        result = append(result, currentPrefix)
    }
    
    cpt.collectKeys(node, currentPrefix, &result)
    return result
}

// 收集键
func (cpt *CompactPrefixTree) collectKeys(node *CompactPrefixTreeNode, prefix string, result *[]string) {
    for _, child := range node.Children {
        newPrefix := prefix + child.Prefix
        
        if child.IsEnd {
            *result = append(*result, newPrefix)
        }
        
        cpt.collectKeys(child, newPrefix, result)
    }
}
```

**时间复杂度**:

| 操作 | 复杂度 |
|------|--------|
| 插入 | $O(k)$ |
| 查找 | $O(k)$ |
| 前缀搜索 | $O(n)$ |

其中，$k$ 是键的长度，$n$ 是匹配前缀的键的数量。

**空间复杂度**: 最坏情况下为 $O(n \cdot k)$，但通常比标准前缀树小得多。

**应用场景**:
- IP路由表
- 电话号码查找
- 字典实现
- 自动补全系统
- 拼写检查

## 概率数据结构

### 8.1 Count-Min Sketch

**定义 8.1** (Count-Min Sketch)
Count-Min Sketch是一种概率数据结构，用于估计数据流中元素的频率，使用多个哈希函数和计数器数组来减少内存使用。

```go
// Count-Min Sketch
type CountMinSketch struct {
    depth     int       // 哈希函数数量
    width     int       // 每个哈希函数的计数器数量
    counts    [][]int   // 计数器数组
    hashFuncs []func(string) int // 哈希函数
}

// 创建Count-Min Sketch
func NewCountMinSketch(depth, width int) *CountMinSketch {
    cms := &CountMinSketch{
        depth:  depth,
        width:  width,
        counts: make([][]int, depth),
    }
    
    // 初始化计数器数组
    for i := range cms.counts {
        cms.counts[i] = make([]int, width)
    }
    
    // 创建哈希函数
    cms.hashFuncs = make([]func(string) int, depth)
    for i := range cms.hashFuncs {
        seed := i + 1
        cms.hashFuncs[i] = func(s string) int {
            hash := 0
            for j := 0; j < len(s); j++ {
                hash = (hash*seed + int(s[j])) % width
            }
            return hash
        }
    }
    
    return cms
}

// 增加元素计数
func (cms *CountMinSketch) Add(item string, count int) {
    for i := 0; i < cms.depth; i++ {
        position := cms.hashFuncs[i](item)
        cms.counts[i][position] += count
    }
}

// 估计元素频率
func (cms *CountMinSketch) Estimate(item string) int {
    min := math.MaxInt32
    
    for i := 0; i < cms.depth; i++ {
        position := cms.hashFuncs[i](item)
        if cms.counts[i][position] < min {
            min = cms.counts[i][position]
        }
    }
    
    return min
}

// 合并两个Count-Min Sketch
func (cms *CountMinSketch) Merge(other *CountMinSketch) error {
    if cms.depth != other.depth || cms.width != other.width {
        return errors.New("无法合并不同维度的Count-Min Sketch")
    }
    
    for i := 0; i < cms.depth; i++ {
        for j := 0; j < cms.width; j++ {
            cms.counts[i][j] += other.counts[i][j]
        }
    }
    
    return nil
}
```

**时间复杂度**:

| 操作 | 复杂度 |
|------|--------|
| 添加 | $O(d)$ |
| 估计 | $O(d)$ |
| 合并 | $O(d \cdot w)$ |

其中，$d$ 是哈希函数的数量，$w$ 是每个哈希函数的计数器数量。

**空间复杂度**: $O(d \cdot w)$

**定理 8.1** (Count-Min Sketch的误差保证)
使用 $d = \lceil \ln(1/\delta) \rceil$ 个哈希函数和 $w = \lceil e/\epsilon \rceil$ 个计数器，Count-Min Sketch可以以 $1-\delta$ 的概率保证估计误差不超过 $\epsilon \cdot ||f||_1$，其中 $||f||_1$ 是所有计数的总和。

**应用场景**:
- 网络流量监控
- 数据库查询优化
- 频繁项集挖掘
- 异常检测

### 8.2 HyperLogLog

**定义 8.2** (HyperLogLog)
HyperLogLog是一种用于估计集合基数（不同元素数量）的概率算法，使用极少的内存就能估计非常大的集合。

```go
// HyperLogLog
type HyperLogLog struct {
    registers []byte   // 寄存器数组
    m         int      // 寄存器数量
    p         uint8    // 精度参数
    alpha     float64  // 修正因子
}

// 创建HyperLogLog
func NewHyperLogLog(precision uint8) *HyperLogLog {
    if precision < 4 || precision > 16 {
        precision = 10 // 默认精度
    }
    
    m := 1 << precision
    
    // 计算修正因子
    var alpha float64
    switch m {
    case 16:
        alpha = 0.673
    case 32:
        alpha = 0.697
    case 64:
        alpha = 0.709
    default:
        alpha = 0.7213 / (1.0 + 1.079/float64(m))
    }
    
    return &HyperLogLog{
        registers: make([]byte, m),
        m:         m,
        p:         precision,
        alpha:     alpha,
    }
}

// 添加元素
func (hll *HyperLogLog) Add(item string) {
    // 计算哈希值
    hash := murmurHash([]byte(item))
    
    // 使用前p位确定寄存器索引
    idx := hash & ((1 << hll.p) - 1)
    
    // 计算前导零的位置（从p+1位开始）
    hash = hash >> hll.p
    rank := uint8(leadingZeros(hash) + 1)
    
    // 更新寄存器
    if rank > hll.registers[idx] {
        hll.registers[idx] = rank
    }
}

// 估计基数
func (hll *HyperLogLog) Estimate() float64 {
    sum := 0.0
    m := float64(hll.m)
    
    // 计算调和平均数
    for _, val := range hll.registers {
        sum += math.Pow(2.0, -float64(val))
    }
    
    // 应用修正因子
    estimate := hll.alpha * m * m / sum
    
    // 小值修正
    if estimate <= 2.5*m {
        // 计算零寄存器的数量
        zeros := 0
        for _, val := range hll.registers {
            if val == 0 {
                zeros++
            }
        }
        
        if zeros > 0 {
            return m * math.Log(m/float64(zeros))
        }
    }
    
    // 大值修正
    if estimate > math.Pow(2.0, 32)/30.0 {
        return -math.Pow(2.0, 32) * math.Log(1.0-estimate/math.Pow(2.0, 32))
    }
    
    return estimate
}

// 合并两个HyperLogLog
func (hll *HyperLogLog) Merge(other *HyperLogLog) error {
    if hll.p != other.p {
        return errors.New("无法合并不同精度的HyperLogLog")
    }
    
    for i := 0; i < hll.m; i++ {
        if other.registers[i] > hll.registers[i] {
            hll.registers[i] = other.registers[i]
        }
    }
    
    return nil
}

// 计算前导零数量
func leadingZeros(x uint64) int {
    if x == 0 {
        return 64
    }
    
    n := 0
    if x <= 0x00000000FFFFFFFF {
        n += 32
        x <<= 32
    }
    if x <= 0x0000FFFFFFFFFFFF {
        n += 16
        x <<= 16
    }
    if x <= 0x00FFFFFFFFFFFFFF {
        n += 8
        x <<= 8
    }
    if x <= 0x0FFFFFFFFFFFFFFF {
        n += 4
        x <<= 4
    }
    if x <= 0x3FFFFFFFFFFFFFFF {
        n += 2
        x <<= 2
    }
    if x <= 0x7FFFFFFFFFFFFFFF {
        n += 1
    }
    
    return n
}

// MurmurHash实现
func murmurHash(data []byte) uint64 {
    const (
        c1 = uint64(0x87c37b91114253d5)
        c2 = uint64(0x4cf5ad432745937f)
        r1 = 31
        r2 = 27
        m  = uint64(5)
        n  = uint64(0x52dce729)
    )
    
    h1 := uint64(0x9747b28c)
    h2 := uint64(0x5493a0e0)
    
    // 处理完整的块
    for len(data) >= 16 {
        k1 := binary.LittleEndian.Uint64(data[:8])
        k2 := binary.LittleEndian.Uint64(data[8:16])
        
        k1 *= c1
        k1 = (k1 << r1) | (k1 >> (64 - r1))
        k1 *= c2
        h1 ^= k1
        
        h1 = (h1 << r2) | (h1 >> (64 - r2))
        h1 = h1*m + n
        
        k2 *= c2
        k2 = (k2 << r1) | (k2 >> (64 - r1))
        k2 *= c1
        h2 ^= k2
        
        h2 = (h2 << r2) | (h2 >> (64 - r2))
        h2 = h2*m + n
        
        data = data[16:]
    }
    
    // 处理剩余字节
    if len(data) > 0 {
        var k1, k2 uint64
        
        for i := 0; i < len(data) && i < 8; i++ {
            k1 |= uint64(data[i]) << uint(i*8)
        }
        
        for i := 8; i < len(data); i++ {
            k2 |= uint64(data[i]) << uint((i-8)*8)
        }
        
        k1 *= c1
        k1 = (k1 << r1) | (k1 >> (64 - r1))
        k1 *= c2
        h1 ^= k1
        
        k2 *= c2
        k2 = (k2 << r1) | (k2 >> (64 - r1))
        k2 *= c1
        h2 ^= k2
    }
    
    // 最终混合
    h1 ^= uint64(len(data))
    h2 ^= uint64(len(data))
    
    h1 += h2
    h2 += h1
    
    h1 ^= h1 >> 33
    h1 *= 0xff51afd7ed558ccd
    h1 ^= h1 >> 33
    h1 *= 0xc4ceb9fe1a85ec53
    h1 ^= h1 >> 33
    
    h2 ^= h2 >> 33
    h2 *= 0xff51afd7ed558ccd
    h2 ^= h2 >> 33
    h2 *= 0xc4ceb9fe1a85ec53
    h2 ^= h2 >> 33
    
    return h1 ^ h2
}
```

**时间复杂度**:

| 操作 | 复杂度 |
|------|--------|
| 添加 | $O(1)$ |
| 估计 | $O(m)$ |
| 合并 | $O(m)$ |

其中，$m$ 是寄存器的数量。

**空间复杂度**: $O(m)$

**定理 8.2** (HyperLogLog的误差保证)
使用 $m = 2^p$ 个寄存器的HyperLogLog算法，其标准误差约为 $1.04/\sqrt{m}$。

**应用场景**:
- 大数据集的去重统计
- 网络流量分析
- 数据库查询优化
- 用户行为分析

### 8.3 跳表变体

#### 8.3.1 确定性跳表

**定义 8.3.1** (确定性跳表)
确定性跳表是跳表的一种变体，其层级结构不是随机生成的，而是按照确定性规则构建，通常基于元素的位置或值。

```go
// 确定性跳表节点
type DeterministicSkipListNode struct {
    Value    int
    Forward  []*DeterministicSkipListNode
}

// 确定性跳表
type DeterministicSkipList struct {
    Head     *DeterministicSkipListNode
    MaxLevel int
    Size     int
}

// 创建确定性跳表
func NewDeterministicSkipList(maxLevel int) *DeterministicSkipList {
    return &DeterministicSkipList{
        Head:     &DeterministicSkipListNode{Forward: make([]*DeterministicSkipListNode, maxLevel)},
        MaxLevel: maxLevel,
        Size:     0,
    }
}

// 计算元素的层级
func (dsl *DeterministicSkipList) levelOf(index int) int {
    level := 0
    for (index & 1) == 0 && level < dsl.MaxLevel-1 {
        index >>= 1
        level++
    }
    return level
}

// 插入元素
func (dsl *DeterministicSkipList) Insert(value int) {
    // 保存每一层的前驱节点
    update := make([]*DeterministicSkipListNode, dsl.MaxLevel)
    current := dsl.Head
    
    // 从最高层开始向下查找插入位置
    for i := dsl.MaxLevel - 1; i >= 0; i-- {
        for current.Forward[i] != nil && current.Forward[i].Value < value {
            current = current.Forward[i]
        }
        update[i] = current
    }
    
    // 检查是否已存在
    if current.Forward[0] != nil && current.Forward[0].Value == value {
        return
    }
    
    // 创建新节点
    level := dsl.levelOf(dsl.Size)
    newNode := &DeterministicSkipListNode{
        Value:   value,
        Forward: make([]*DeterministicSkipListNode, level+1),
    }
    
    // 更新指针
    for i := 0; i <= level; i++ {
        newNode.Forward[i] = update[i].Forward[i]
        update[i].Forward[i] = newNode
    }
    
    dsl.Size++
    
    // 重建跳表以保持确定性结构
    if dsl.Size > 1 && (dsl.Size & (dsl.Size - 1)) == 0 {
        dsl.rebuild()
    }
}

// 重建跳表
func (dsl *DeterministicSkipList) rebuild() {
    // 收集所有元素
    var elements []int
    current := dsl.Head.Forward[0]
    for current != nil {
        elements = append(elements, current.Value)
        current = current.Forward[0]
    }
    
    // 创建新的头节点
    dsl.Head = &DeterministicSkipListNode{Forward: make([]*DeterministicSkipListNode, dsl.MaxLevel)}
    dsl.Size = 0
    
    // 重新插入所有元素
    for _, value := range elements {
        dsl.Insert(value)
    }
}

// 搜索元素
func (dsl *DeterministicSkipList) Search(value int) bool {
    current := dsl.Head
    
    // 从最高层开始向下搜索
    for i := dsl.MaxLevel - 1; i >= 0; i-- {
        for current.Forward[i] != nil && current.Forward[i].Value < value {
            current = current.Forward[i]
        }
    }
    
    current = current.Forward[0]
    return current != nil && current.Value == value
}
```

**时间复杂度**:

| 操作 | 平均情况 | 最坏情况 |
|------|----------|----------|
| 搜索 | $O(\log n)$ | $O(\log n)$ |
| 插入 | $O(\log n)$ | $O(n)$ (重建时) |
| 删除 | $O(\log n)$ | $O(\log n)$ |

**空间复杂度**: $O(n \log n)$

#### 8.3.2 跳表变种：Biased Skip List

**定义 8.3.2** (偏向跳表)
偏向跳表是跳表的一种变体，其中某些频繁访问的元素被提升到更高的层级，以加速访问。

```go
// 偏向跳表节点
type BiasedSkipListNode struct {
    Value      int
    Forward    []*BiasedSkipListNode
    AccessCount int
}

// 偏向跳表
type BiasedSkipList struct {
    Head       *BiasedSkipListNode
    MaxLevel   int
    Level      int
    Size       int
    Threshold  int  // 访问计数阈值
}

// 创建偏向跳表
func NewBiasedSkipList(maxLevel, threshold int) *BiasedSkipList {
    return &BiasedSkipList{
        Head:      &BiasedSkipListNode{Forward: make([]*BiasedSkipListNode, maxLevel)},
        MaxLevel:  maxLevel,
        Level:     0,
        Size:      0,
        Threshold: threshold,
    }
}

// 随机生成层级
func (bsl *BiasedSkipList) randomLevel() int {
    level := 0
    for rand.Float64() < 0.5 && level < bsl.MaxLevel-1 {
        level++
    }
    return level
}

// 搜索元素
func (bsl *BiasedSkipList) Search(value int) bool {
    current := bsl.Head
    
    // 从最高层开始向下搜索
    for i := bsl.Level; i >= 0; i-- {
        for current.Forward[i] != nil && current.Forward[i].Value < value {
            current = current.Forward[i]
        }
    }
    
    current = current.Forward[0]
    
    if current != nil && current.Value == value {
        // 增加访问计数
        current.AccessCount++
        
        // 如果访问次数超过阈值，提升层级
        if current.AccessCount >= bsl.Threshold {
            bsl.promoteNode(current)
            current.AccessCount = 0
        }
        
        return true
    }
    
    return false
}

// 提升节点层级
func (bsl *BiasedSkipList) promoteNode(node *BiasedSkipListNode) {
    if len(node.Forward) >= bsl.MaxLevel {
        return
    }
    
    // 保存每一层的前驱节点
    update := make([]*BiasedSkipListNode, bsl.MaxLevel)
    current := bsl.Head
    
    // 查找节点的前驱
    for i := bsl.Level; i >= 0; i-- {
        for current.Forward[i] != nil && current.Forward[i].Value < node.Value {
            current = current.Forward[i]
        }
        update[i] = current
    }
    
    // 增加一层
    newLevel := len(node.Forward)
    if newLevel < bsl.MaxLevel {
        node.Forward = append(node.Forward, make([]*BiasedSkipListNode, 1)...)
        
        // 更新最大层级
        if newLevel > bsl.Level {
            bsl.Level = newLevel
        }
        
        // 更新指针
        node.Forward[newLevel] = update[newLevel].Forward[newLevel]
        update[newLevel].Forward[newLevel] = node
    }
}

// 插入元素
func (bsl *BiasedSkipList) Insert(value int) {
    // 保存每一层的前驱节点
    update := make([]*BiasedSkipListNode, bsl.MaxLevel)
    current := bsl.Head
    
    // 从最高层开始向下查找插入位置
    for i := bsl.Level; i >= 0; i-- {
        for current.Forward[i] != nil && current.Forward[i].Value < value {
            current = current.Forward[i]
        }
        update[i] = current
    }
    
    // 检查是否已存在
    if current.Forward[0] != nil && current.Forward[0].Value == value {
        return
    }
    
    // 获取随机层级
    level := bsl.randomLevel()
    
    // 更新跳表的最大层级
    if level > bsl.Level {
        for i := bsl.Level + 1; i <= level; i++ {
            update[i] = bsl.Head
        }
        bsl.Level = level
    }
    
    // 创建新节点
    newNode := &BiasedSkipListNode{
        Value:      value,
        Forward:    make([]*BiasedSkipListNode, level+1),
        AccessCount: 0,
    }
    
    // 更新指针
    for i := 0; i <= level; i++ {
        newNode.Forward[i] = update[i].Forward[i]
        update[i].Forward[i] = newNode
    }
    
    bsl.Size++
}
```

**时间复杂度**:

| 操作 | 平均情况 | 最坏情况 |
|------|----------|----------|
| 搜索 | $O(\log n)$ | $O(n)$ |
| 插入 | $O(\log n)$ | $O(n)$ |
| 删除 | $O(\log n)$ | $O(n)$ |

**空间复杂度**: $O(n \log n)$

**应用场景**:
- 缓存实现
- 频繁访问元素的快速查找
- 自适应数据结构

## Golang实现

## 性能分析与对比

## 最佳实践

## 案例研究

## 总结

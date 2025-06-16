# 算法与数据结构分析

## 目录

1. [基础算法 (Basic Algorithms)](01-Basic-Algorithms/README.md)
2. [数据结构 (Data Structures)](02-Data-Structures/README.md)
3. [并发算法 (Concurrent Algorithms)](03-Concurrent-Algorithms/README.md)
4. [分布式算法 (Distributed Algorithms)](04-Distributed-Algorithms/README.md)
5. [机器学习算法 (Machine Learning Algorithms)](05-Machine-Learning-Algorithms/README.md)
6. [图算法 (Graph Algorithms)](06-Graph-Algorithms/README.md)

## 概述

算法与数据结构是计算机科学的核心基础，本章节基于形式化方法，对常用算法和数据结构进行严格的数学定义、复杂度分析和Golang实现。

### 算法分析基础

#### 定义 1.1 (算法)

算法是一个四元组 $\mathcal{A} = (\mathcal{I}, \mathcal{O}, \mathcal{P}, \mathcal{C})$，其中：

- $\mathcal{I}$ 是输入集合
- $\mathcal{O}$ 是输出集合
- $\mathcal{P}$ 是处理步骤
- $\mathcal{C}$ 是复杂度分析

#### 定义 1.2 (时间复杂度)

算法 $\mathcal{A}$ 的时间复杂度定义为：
$$T(n) = O(f(n))$$

其中 $n$ 是输入规模，$f(n)$ 是增长函数。

#### 定义 1.3 (空间复杂度)

算法 $\mathcal{A}$ 的空间复杂度定义为：
$$S(n) = O(g(n))$$

其中 $g(n)$ 是空间增长函数。

### 基础算法

#### 1.1 排序算法

##### 定义 1.4 (排序问题)

排序问题：给定序列 $A = [a_1, a_2, ..., a_n]$，找到排列 $\pi$ 使得：
$$a_{\pi(1)} \leq a_{\pi(2)} \leq ... \leq a_{\pi(n)}$$

##### 1.1.1 快速排序 (QuickSort)

###### 算法描述

```go
func QuickSort(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    
    pivot := arr[0]
    var left, right []int
    
    for i := 1; i < len(arr); i++ {
        if arr[i] <= pivot {
            left = append(left, arr[i])
        } else {
            right = append(right, arr[i])
        }
    }
    
    left = QuickSort(left)
    right = QuickSort(right)
    
    return append(append(left, pivot), right...)
}
```

###### 复杂度分析

- **时间复杂度**: $O(n \log n)$ 平均情况，$O(n^2)$ 最坏情况
- **空间复杂度**: $O(\log n)$ 平均情况，$O(n)$ 最坏情况

###### 数学证明

**定理 1.1**: 快速排序的平均时间复杂度为 $O(n \log n)$

**证明**:

1. 设 $T(n)$ 为排序 $n$ 个元素的时间复杂度
2. 在平均情况下，pivot将数组分为两个大致相等的部分
3. 递归关系：$T(n) = 2T(n/2) + O(n)$
4. 根据主定理，$T(n) = O(n \log n)$

##### 1.1.2 归并排序 (MergeSort)

###### 算法描述

```go
func MergeSort(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    
    mid := len(arr) / 2
    left := MergeSort(arr[:mid])
    right := MergeSort(arr[mid:])
    
    return merge(left, right)
}

func merge(left, right []int) []int {
    result := make([]int, 0, len(left)+len(right))
    i, j := 0, 0
    
    for i < len(left) && j < len(right) {
        if left[i] <= right[j] {
            result = append(result, left[i])
            i++
        } else {
            result = append(result, right[j])
            j++
        }
    }
    
    result = append(result, left[i:]...)
    result = append(result, right[j:]...)
    
    return result
}
```

###### 复杂度分析

- **时间复杂度**: $O(n \log n)$
- **空间复杂度**: $O(n)$

#### 1.2 搜索算法

##### 定义 1.5 (搜索问题)

搜索问题：在集合 $S$ 中查找元素 $x$，返回位置或不存在。

##### 1.2.1 二分搜索 (Binary Search)

###### 算法描述

```go
func BinarySearch(arr []int, target int) int {
    left, right := 0, len(arr)-1
    
    for left <= right {
        mid := left + (right-left)/2
        if arr[mid] == target {
            return mid
        } else if arr[mid] < target {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    
    return -1
}
```

###### 复杂度分析

- **时间复杂度**: $O(\log n)$
- **空间复杂度**: $O(1)$

###### 数学证明

**定理 1.2**: 二分搜索的时间复杂度为 $O(\log n)$

**证明**:

1. 每次迭代将搜索空间减半
2. 设 $k$ 为迭代次数，则 $n/2^k = 1$
3. 解得 $k = \log_2 n$
4. 因此时间复杂度为 $O(\log n)$

### 数据结构

#### 2.1 线性数据结构

##### 2.1.1 数组 (Array)

###### 定义 2.1 (数组)

数组是一个有序的元素序列：
$$Array = [a_1, a_2, ..., a_n]$$

###### Golang实现

```go
type Array[T any] struct {
    data []T
    size int
}

func NewArray[T any](capacity int) *Array[T] {
    return &Array[T]{
        data: make([]T, capacity),
        size: 0,
    }
}

func (a *Array[T]) Get(index int) (T, error) {
    if index < 0 || index >= a.size {
        var zero T
        return zero, errors.New("index out of bounds")
    }
    return a.data[index], nil
}

func (a *Array[T]) Set(index int, value T) error {
    if index < 0 || index >= a.size {
        return errors.New("index out of bounds")
    }
    a.data[index] = value
    return nil
}
```

##### 2.1.2 链表 (Linked List)

###### 定义 2.2 (链表)

链表是由节点组成的线性结构：
$$LinkedList = Node_1 \rightarrow Node_2 \rightarrow ... \rightarrow Node_n$$

###### Golang实现

```go
type Node[T any] struct {
    data T
    next *Node[T]
}

type LinkedList[T any] struct {
    head *Node[T]
    size int
}

func (ll *LinkedList[T]) Insert(value T) {
    newNode := &Node[T]{data: value, next: ll.head}
    ll.head = newNode
    ll.size++
}

func (ll *LinkedList[T]) Delete(value T) bool {
    if ll.head == nil {
        return false
    }
    
    if ll.head.data == value {
        ll.head = ll.head.next
        ll.size--
        return true
    }
    
    current := ll.head
    for current.next != nil {
        if current.next.data == value {
            current.next = current.next.next
            ll.size--
            return true
        }
        current = current.next
    }
    
    return false
}
```

#### 2.2 树形数据结构

##### 2.2.1 二叉树 (Binary Tree)

###### 定义 2.3 (二叉树)

二叉树是每个节点最多有两个子节点的树：
$$BinaryTree = (Node, Left, Right)$$

###### Golang实现

```go
type TreeNode[T any] struct {
    data  T
    left  *TreeNode[T]
    right *TreeNode[T]
}

type BinaryTree[T any] struct {
    root *TreeNode[T]
}

func (bt *BinaryTree[T]) Insert(value T) {
    bt.root = bt.insertRecursive(bt.root, value)
}

func (bt *BinaryTree[T]) insertRecursive(node *TreeNode[T], value T) *TreeNode[T] {
    if node == nil {
        return &TreeNode[T]{data: value}
    }
    
    if value < node.data {
        node.left = bt.insertRecursive(node.left, value)
    } else {
        node.right = bt.insertRecursive(node.right, value)
    }
    
    return node
}

func (bt *BinaryTree[T]) InorderTraversal() []T {
    var result []T
    bt.inorderRecursive(bt.root, &result)
    return result
}

func (bt *BinaryTree[T]) inorderRecursive(node *TreeNode[T], result *[]T) {
    if node != nil {
        bt.inorderRecursive(node.left, result)
        *result = append(*result, node.data)
        bt.inorderRecursive(node.right, result)
    }
}
```

##### 2.2.2 红黑树 (Red-Black Tree)

###### 定义 2.4 (红黑树)

红黑树是自平衡的二叉搜索树，满足以下性质：

1. 每个节点是红色或黑色
2. 根节点是黑色
3. 红节点的子节点都是黑色
4. 从根到叶子的所有路径包含相同数量的黑节点

###### Golang实现

```go
type Color bool

const (
    Red   Color = false
    Black Color = true
)

type RBNode[T any] struct {
    data  T
    color Color
    left  *RBNode[T]
    right *RBNode[T]
    parent *RBNode[T]
}

type RedBlackTree[T any] struct {
    root *RBNode[T]
    nil  *RBNode[T]
}

func (rbt *RedBlackTree[T]) Insert(value T) {
    node := &RBNode[T]{
        data:  value,
        color: Red,
        left:  rbt.nil,
        right: rbt.nil,
    }
    
    rbt.insertNode(node)
    rbt.fixInsert(node)
}

func (rbt *RedBlackTree[T]) fixInsert(node *RBNode[T]) {
    for node.parent.color == Red {
        if node.parent == node.parent.parent.left {
            uncle := node.parent.parent.right
            if uncle.color == Red {
                node.parent.color = Black
                uncle.color = Black
                node.parent.parent.color = Red
                node = node.parent.parent
            } else {
                if node == node.parent.right {
                    node = node.parent
                    rbt.leftRotate(node)
                }
                node.parent.color = Black
                node.parent.parent.color = Red
                rbt.rightRotate(node.parent.parent)
            }
        } else {
            // 对称情况
        }
    }
    rbt.root.color = Black
}
```

#### 2.3 图数据结构

##### 2.3.1 邻接矩阵 (Adjacency Matrix)

###### 定义 2.5 (邻接矩阵)

图的邻接矩阵是一个 $n \times n$ 的矩阵 $A$，其中：
$$A[i][j] = \begin{cases}
1 & \text{if } (i,j) \in E \\
0 & \text{otherwise}
\end{cases}$$

###### Golang实现
```go
type AdjacencyMatrix struct {
    matrix [][]bool
    size   int
}

func NewAdjacencyMatrix(size int) *AdjacencyMatrix {
    matrix := make([][]bool, size)
    for i := range matrix {
        matrix[i] = make([]bool, size)
    }

    return &AdjacencyMatrix{
        matrix: matrix,
        size:   size,
    }
}

func (am *AdjacencyMatrix) AddEdge(from, to int) {
    if from >= 0 && from < am.size && to >= 0 && to < am.size {
        am.matrix[from][to] = true
    }
}

func (am *AdjacencyMatrix) HasEdge(from, to int) bool {
    if from >= 0 && from < am.size && to >= 0 && to < am.size {
        return am.matrix[from][to]
    }
    return false
}
```

##### 2.3.2 邻接表 (Adjacency List)

###### 定义 2.6 (邻接表)
图的邻接表是一个数组，每个元素是一个链表：
$$AdjacencyList = [List_1, List_2, ..., List_n]$$

###### Golang实现
```go
type AdjacencyList struct {
    lists [][]int
    size  int
}

func NewAdjacencyList(size int) *AdjacencyList {
    lists := make([][]int, size)
    return &AdjacencyList{
        lists: lists,
        size:  size,
    }
}

func (al *AdjacencyList) AddEdge(from, to int) {
    if from >= 0 && from < al.size {
        al.lists[from] = append(al.lists[from], to)
    }
}

func (al *AdjacencyList) GetNeighbors(vertex int) []int {
    if vertex >= 0 && vertex < al.size {
        return al.lists[vertex]
    }
    return nil
}
```

### 并发算法

#### 3.1 并发数据结构

##### 3.1.1 并发队列 (Concurrent Queue)

###### 定义 3.1 (并发队列)
并发队列是线程安全的队列数据结构：
$$ConcurrentQueue = (Queue, Mutex)$$

###### Golang实现
```go
type ConcurrentQueue[T any] struct {
    data []T
    mutex sync.RWMutex
}

func (cq *ConcurrentQueue[T]) Enqueue(item T) {
    cq.mutex.Lock()
    defer cq.mutex.Unlock()
    cq.data = append(cq.data, item)
}

func (cq *ConcurrentQueue[T]) Dequeue() (T, bool) {
    cq.mutex.Lock()
    defer cq.mutex.Unlock()

    if len(cq.data) == 0 {
        var zero T
        return zero, false
    }

    item := cq.data[0]
    cq.data = cq.data[1:]
    return item, true
}

func (cq *ConcurrentQueue[T]) Size() int {
    cq.mutex.RLock()
    defer cq.mutex.RUnlock()
    return len(cq.data)
}
```

##### 3.1.2 无锁栈 (Lock-Free Stack)

###### 定义 3.2 (无锁栈)
无锁栈使用原子操作实现线程安全：
$$LockFreeStack = (Stack, AtomicOperations)$$

###### Golang实现
```go
type LockFreeStack[T any] struct {
    head *atomic.Value
}

type node[T any] struct {
    data T
    next *atomic.Value
}

func NewLockFreeStack[T any]() *LockFreeStack[T] {
    return &LockFreeStack[T]{
        head: &atomic.Value{},
    }
}

func (lfs *LockFreeStack[T]) Push(item T) {
    newNode := &node[T]{
        data: item,
        next: &atomic.Value{},
    }

    for {
        oldHead := lfs.head.Load()
        if oldHead == nil {
            newNode.next.Store(nil)
        } else {
            newNode.next.Store(oldHead)
        }

        if lfs.head.CompareAndSwap(oldHead, newNode) {
            break
        }
    }
}

func (lfs *LockFreeStack[T]) Pop() (T, bool) {
    for {
        oldHead := lfs.head.Load()
        if oldHead == nil {
            var zero T
            return zero, false
        }

        headNode := oldHead.(*node[T])
        newHead := headNode.next.Load()

        if lfs.head.CompareAndSwap(oldHead, newHead) {
            return headNode.data, true
        }
    }
}
```

#### 3.2 并发算法

##### 3.2.1 并行归并排序

###### 算法描述
```go
func ParallelMergeSort(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }

    mid := len(arr) / 2

    var left, right []int
    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        left = ParallelMergeSort(arr[:mid])
    }()

    go func() {
        defer wg.Done()
        right = ParallelMergeSort(arr[mid:])
    }()

    wg.Wait()
    return merge(left, right)
}
```

###### 复杂度分析
- **时间复杂度**: $O(n \log n)$
- **空间复杂度**: $O(n)$
- **并行度**: $O(\log n)$

##### 3.2.2 并行矩阵乘法

###### 算法描述
```go
func ParallelMatrixMultiply(a, b [][]int) [][]int {
    rows := len(a)
    cols := len(b[0])
    result := make([][]int, rows)

    for i := range result {
        result[i] = make([]int, cols)
    }

    var wg sync.WaitGroup
    for i := 0; i < rows; i++ {
        for j := 0; j < cols; j++ {
            wg.Add(1)
            go func(i, j int) {
                defer wg.Done()
                for k := 0; k < len(a[0]); k++ {
                    result[i][j] += a[i][k] * b[k][j]
                }
            }(i, j)
        }
    }

    wg.Wait()
    return result
}
```

### 分布式算法

#### 4.1 共识算法

##### 4.1.1 Raft算法

###### 定义 4.1 (Raft状态)
Raft节点有三种状态：
$$RaftState \in \{Follower, Candidate, Leader\}$$

###### Golang实现
```go
type RaftNode struct {
    id        int
    state     RaftState
    term      int
    votedFor  *int
    log       []LogEntry
    commitIndex int
    lastApplied int
    nextIndex  map[int]int
    matchIndex map[int]int

    electionTimeout time.Duration
    heartbeatInterval time.Duration

    mutex sync.RWMutex
}

type LogEntry struct {
    Term    int
    Command interface{}
}

func (rn *RaftNode) StartElection() {
    rn.mutex.Lock()
    rn.state = Candidate
    rn.term++
    rn.votedFor = &rn.id
    rn.mutex.Unlock()

    votes := 1 // 自己的一票

    for _, peer := range rn.peers {
        go func(p *RaftNode) {
            if rn.requestVote(p) {
                votes++
                if votes > len(rn.peers)/2 {
                    rn.becomeLeader()
                }
            }
        }(peer)
    }
}
```

##### 4.1.2 Paxos算法

###### 定义 4.2 (Paxos阶段)
Paxos算法分为两个阶段：
1. **准备阶段**: 提议者选择提案编号
2. **接受阶段**: 提议者提交提案

###### Golang实现
```go
type PaxosNode struct {
    id           int
    proposers    map[int]*Proposer
    acceptors    map[int]*Acceptor
    learners     map[int]*Learner

    mutex        sync.RWMutex
}

type Proposer struct {
    id        int
    round     int
    value     interface{}
    promises  map[int]Promise
}

type Acceptor struct {
    id           int
    promisedNum  int
    acceptedNum  int
    acceptedValue interface{}
}

func (pn *PaxosNode) Propose(value interface{}) {
    proposer := &Proposer{
        id:    pn.id,
        round: pn.nextRound(),
        value: value,
    }

    // 准备阶段
    promises := pn.prepare(proposer)

    // 接受阶段
    if len(promises) > len(pn.acceptors)/2 {
        pn.accept(proposer, promises)
    }
}
```

#### 4.2 分布式排序

##### 4.2.1 分布式归并排序

###### 算法描述
```go
type DistributedMergeSort struct {
    nodes []*Node
    data  []int
}

func (dms *DistributedMergeSort) Sort() []int {
    // 1. 数据分片
    chunks := dms.partition()

    // 2. 并行排序
    var wg sync.WaitGroup
    for i, chunk := range chunks {
        wg.Add(1)
        go func(i int, chunk []int) {
            defer wg.Done()
            sort.Ints(chunk)
            chunks[i] = chunk
        }(i, chunk)
    }
    wg.Wait()

    // 3. 分布式归并
    return dms.distributedMerge(chunks)
}

func (dms *DistributedMergeSort) distributedMerge(chunks [][]int) []int {
    if len(chunks) == 1 {
        return chunks[0]
    }

    // 两两归并
    var wg sync.WaitGroup
    for i := 0; i < len(chunks)-1; i += 2 {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            chunks[i] = merge(chunks[i], chunks[i+1])
        }(i)
    }
    wg.Wait()

    // 递归归并
    return dms.distributedMerge(chunks)
}
```

### 机器学习算法

#### 5.1 分类算法

##### 5.1.1 K近邻算法 (K-NN)

###### 定义 5.1 (K-NN)
K-NN算法基于距离度量进行分类：
$$Class(x) = \arg\max_{c} \sum_{i \in N_k(x)} I(y_i = c)$$

其中 $N_k(x)$ 是 $x$ 的 $k$ 个最近邻。

###### Golang实现
```go
type KNN struct {
    k      int
    data   [][]float64
    labels []int
}

func (knn *KNN) Predict(features []float64) int {
    distances := make([]Distance, len(knn.data))

    for i, point := range knn.data {
        distances[i] = Distance{
            index:    i,
            distance: euclideanDistance(features, point),
        }
    }

    // 排序获取k个最近邻
    sort.Slice(distances, func(i, j int) bool {
        return distances[i].distance < distances[j].distance
    })

    // 投票
    votes := make(map[int]int)
    for i := 0; i < knn.k; i++ {
        label := knn.labels[distances[i].index]
        votes[label]++
    }

    // 返回得票最多的类别
    maxVotes := 0
    predictedClass := 0
    for class, votes := range votes {
        if votes > maxVotes {
            maxVotes = votes
            predictedClass = class
        }
    }

    return predictedClass
}

func euclideanDistance(a, b []float64) float64 {
    sum := 0.0
    for i := range a {
        diff := a[i] - b[i]
        sum += diff * diff
    }
    return math.Sqrt(sum)
}
```

##### 5.1.2 决策树算法

###### 定义 5.2 (决策树)
决策树是一个树形结构，每个内部节点表示一个特征测试：
$$DecisionTree = (Feature, Threshold, LeftChild, RightChild)$$

###### Golang实现
```go
type DecisionTreeNode struct {
    feature   int
    threshold float64
    isLeaf    bool
    class     int
    left      *DecisionTreeNode
    right     *DecisionTreeNode
}

type DecisionTree struct {
    root *DecisionTreeNode
}

func (dt *DecisionTree) Train(data [][]float64, labels []int) {
    dt.root = dt.buildTree(data, labels, 0)
}

func (dt *DecisionTree) buildTree(data [][]float64, labels []int, depth int) *DecisionTreeNode {
    // 停止条件
    if len(data) == 0 || depth >= 10 {
        return &DecisionTreeNode{isLeaf: true, class: majorityClass(labels)}
    }

    // 找到最佳分割
    bestFeature, bestThreshold := dt.findBestSplit(data, labels)

    if bestFeature == -1 {
        return &DecisionTreeNode{isLeaf: true, class: majorityClass(labels)}
    }

    // 分割数据
    leftData, leftLabels, rightData, rightLabels := dt.splitData(data, labels, bestFeature, bestThreshold)

    return &DecisionTreeNode{
        feature:   bestFeature,
        threshold: bestThreshold,
        left:      dt.buildTree(leftData, leftLabels, depth+1),
        right:     dt.buildTree(rightData, rightLabels, depth+1),
    }
}

func (dt *DecisionTree) Predict(features []float64) int {
    return dt.predictRecursive(dt.root, features)
}

func (dt *DecisionTree) predictRecursive(node *DecisionTreeNode, features []float64) int {
    if node.isLeaf {
        return node.class
    }

    if features[node.feature] <= node.threshold {
        return dt.predictRecursive(node.left, features)
    } else {
        return dt.predictRecursive(node.right, features)
    }
}
```

#### 5.2 聚类算法

##### 5.2.1 K均值算法 (K-Means)

###### 定义 5.3 (K-Means)
K-Means算法最小化簇内平方误差：
$$\arg\min_{S} \sum_{i=1}^{k} \sum_{x \in S_i} \|x - \mu_i\|^2$$

其中 $S_i$ 是第 $i$ 个簇，$\mu_i$ 是簇中心。

###### Golang实现
```go
type KMeans struct {
    k       int
    centers [][]float64
}

func (km *KMeans) Fit(data [][]float64) {
    // 初始化中心点
    km.centers = km.initializeCenters(data)

    for iteration := 0; iteration < 100; iteration++ {
        // 分配点到最近的中心
        clusters := km.assignToClusters(data)

        // 更新中心点
        newCenters := km.updateCenters(data, clusters)

        // 检查收敛
        if km.isConverged(newCenters) {
            break
        }

        km.centers = newCenters
    }
}

func (km *KMeans) assignToClusters(data [][]float64) []int {
    clusters := make([]int, len(data))

    for i, point := range data {
        minDistance := math.Inf(1)
        bestCluster := 0

        for j, center := range km.centers {
            distance := euclideanDistance(point, center)
            if distance < minDistance {
                minDistance = distance
                bestCluster = j
            }
        }

        clusters[i] = bestCluster
    }

    return clusters
}

func (km *KMeans) updateCenters(data [][]float64, clusters []int) [][]float64 {
    newCenters := make([][]float64, km.k)
    counts := make([]int, km.k)

    // 初始化
    for i := range newCenters {
        newCenters[i] = make([]float64, len(data[0]))
    }

    // 累加
    for i, point := range data {
        cluster := clusters[i]
        for j, value := range point {
            newCenters[cluster][j] += value
        }
        counts[cluster]++
    }

    // 平均
    for i := range newCenters {
        if counts[i] > 0 {
            for j := range newCenters[i] {
                newCenters[i][j] /= float64(counts[i])
            }
        }
    }

    return newCenters
}
```

### 图算法

#### 6.1 最短路径算法

##### 6.1.1 Dijkstra算法

###### 定义 6.1 (Dijkstra算法)
Dijkstra算法找到从源点到所有其他点的最短路径：
$$d[v] = \min(d[v], d[u] + w(u,v))$$

###### Golang实现
```go
type Graph struct {
    vertices int
    edges    [][]Edge
}

type Edge struct {
    to     int
    weight int
}

func (g *Graph) Dijkstra(source int) []int {
    distances := make([]int, g.vertices)
    visited := make([]bool, g.vertices)

    // 初始化距离
    for i := range distances {
        distances[i] = math.MaxInt32
    }
    distances[source] = 0

    for i := 0; i < g.vertices; i++ {
        // 找到未访问的最小距离顶点
        u := g.findMinDistance(distances, visited)
        visited[u] = true

        // 更新邻接顶点的距离
        for _, edge := range g.edges[u] {
            v := edge.to
            if !visited[v] && distances[u] != math.MaxInt32 &&
               distances[u]+edge.weight < distances[v] {
                distances[v] = distances[u] + edge.weight
            }
        }
    }

    return distances
}

func (g *Graph) findMinDistance(distances []int, visited []bool) int {
    min := math.MaxInt32
    minIndex := -1

    for i := 0; i < g.vertices; i++ {
        if !visited[i] && distances[i] <= min {
            min = distances[i]
            minIndex = i
        }
    }

    return minIndex
}
```

##### 6.1.2 Floyd-Warshall算法

###### 定义 6.2 (Floyd-Warshall算法)
Floyd-Warshall算法找到所有顶点对之间的最短路径：
$$d[i][j] = \min(d[i][j], d[i][k] + d[k][j])$$

###### Golang实现
```go
func (g *Graph) FloydWarshall() [][]int {
    distances := make([][]int, g.vertices)

    // 初始化距离矩阵
    for i := range distances {
        distances[i] = make([]int, g.vertices)
        for j := range distances[i] {
            if i == j {
                distances[i][j] = 0
            } else {
                distances[i][j] = math.MaxInt32
            }
        }
    }

    // 设置初始边权重
    for i := range g.edges {
        for _, edge := range g.edges[i] {
            distances[i][edge.to] = edge.weight
        }
    }

    // Floyd-Warshall算法
    for k := 0; k < g.vertices; k++ {
        for i := 0; i < g.vertices; i++ {
            for j := 0; j < g.vertices; j++ {
                if distances[i][k] != math.MaxInt32 && distances[k][j] != math.MaxInt32 &&
                   distances[i][k]+distances[k][j] < distances[i][j] {
                    distances[i][j] = distances[i][k] + distances[k][j]
                }
            }
        }
    }

    return distances
}
```

#### 6.2 最小生成树算法

##### 6.2.1 Kruskal算法

###### 定义 6.3 (Kruskal算法)
Kruskal算法按权重递增顺序选择边，避免环：
$$MST = \arg\min_{T} \sum_{e \in T} w(e)$$

###### Golang实现
```go
type Edge struct {
    from   int
    to     int
    weight int
}

func (g *Graph) Kruskal() []Edge {
    var mst []Edge
    edges := g.getAllEdges()

    // 按权重排序
    sort.Slice(edges, func(i, j int) bool {
        return edges[i].weight < edges[j].weight
    })

    // 并查集
    uf := NewUnionFind(g.vertices)

    for _, edge := range edges {
        if uf.Find(edge.from) != uf.Find(edge.to) {
            mst = append(mst, edge)
            uf.Union(edge.from, edge.to)
        }
    }

    return mst
}

type UnionFind struct {
    parent []int
    rank   []int
}

func NewUnionFind(size int) *UnionFind {
    parent := make([]int, size)
    rank := make([]int, size)

    for i := range parent {
        parent[i] = i
    }

    return &UnionFind{parent: parent, rank: rank}
}

func (uf *UnionFind) Find(x int) int {
    if uf.parent[x] != x {
        uf.parent[x] = uf.Find(uf.parent[x])
    }
    return uf.parent[x]
}

func (uf *UnionFind) Union(x, y int) {
    rootX := uf.Find(x)
    rootY := uf.Find(y)

    if rootX == rootY {
        return
    }

    if uf.rank[rootX] < uf.rank[rootY] {
        uf.parent[rootX] = rootY
    } else if uf.rank[rootX] > uf.rank[rootY] {
        uf.parent[rootY] = rootX
    } else {
        uf.parent[rootY] = rootX
        uf.rank[rootX]++
    }
}
```

### 性能优化策略

#### 7.1 算法优化

##### 7.1.1 缓存优化
```go
// 使用缓存优化斐波那契计算
type FibonacciCache struct {
    cache map[int]int
    mutex sync.RWMutex
}

func (fc *FibonacciCache) Calculate(n int) int {
    fc.mutex.RLock()
    if result, exists := fc.cache[n]; exists {
        fc.mutex.RUnlock()
        return result
    }
    fc.mutex.RUnlock()

    fc.mutex.Lock()
    defer fc.mutex.Unlock()

    // 双重检查
    if result, exists := fc.cache[n]; exists {
        return result
    }

    var result int
    if n <= 1 {
        result = n
    } else {
        result = fc.Calculate(n-1) + fc.Calculate(n-2)
    }

    fc.cache[n] = result
    return result
}
```

##### 7.1.2 并行优化
```go
// 并行计算数组和
func ParallelSum(arr []int) int {
    numCPU := runtime.NumCPU()
    chunkSize := len(arr) / numCPU

    var wg sync.WaitGroup
    results := make(chan int, numCPU)

    for i := 0; i < numCPU; i++ {
        wg.Add(1)
        go func(start int) {
            defer wg.Done()
            end := start + chunkSize
            if end > len(arr) {
                end = len(arr)
            }

            sum := 0
            for j := start; j < end; j++ {
                sum += arr[j]
            }
            results <- sum
        }(i * chunkSize)
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    total := 0
    for sum := range results {
        total += sum
    }

    return total
}
```

#### 7.2 数据结构优化

##### 7.2.1 内存池
```go
type ObjectPool[T any] struct {
    pool chan T
    new  func() T
}

func NewObjectPool[T any](size int, newFunc func() T) *ObjectPool[T] {
    pool := make(chan T, size)
    for i := 0; i < size; i++ {
        pool <- newFunc()
    }

    return &ObjectPool[T]{
        pool: pool,
        new:  newFunc,
    }
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
        // 池已满，丢弃对象
    }
}
```

##### 7.2.2 对象复用
```go
type ReusableBuffer struct {
    buffer []byte
    size   int
}

func (rb *ReusableBuffer) Reset() {
    rb.size = 0
}

func (rb *ReusableBuffer) Write(data []byte) {
    if rb.size+len(data) > len(rb.buffer) {
        rb.buffer = make([]byte, rb.size+len(data))
    }
    copy(rb.buffer[rb.size:], data)
    rb.size += len(data)
}

func (rb *ReusableBuffer) Bytes() []byte {
    return rb.buffer[:rb.size]
}
```

### 最佳实践

#### 8.1 算法选择指南

1. **排序算法选择**:
   - 小数据集 (< 50): 插入排序
   - 中等数据集 (50-1000): 快速排序
   - 大数据集 (> 1000): 归并排序

2. **搜索算法选择**:
   - 有序数据: 二分搜索
   - 无序数据: 线性搜索
   - 频繁搜索: 哈希表

3. **图算法选择**:
   - 单源最短路径: Dijkstra
   - 全源最短路径: Floyd-Warshall
   - 最小生成树: Kruskal

#### 8.2 性能调优

1. **时间复杂度优化**:
   - 使用更高效的算法
   - 减少不必要的计算
   - 利用数据结构特性

2. **空间复杂度优化**:
   - 重用内存
   - 使用对象池
   - 及时释放资源

3. **并发优化**:
   - 并行处理独立任务
   - 减少锁竞争
   - 使用无锁数据结构

### 持续更新

本文档将根据算法理论的发展和Golang语言特性的变化持续更新。

---

*最后更新时间: 2024-01-XX*
*版本: 1.0.0*

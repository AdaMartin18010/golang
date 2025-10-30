# Go数据结构

Go数据结构完整指南，涵盖基础数据结构、高级数据结构和算法实现。

---

## 📚 核心内容

### 基础数据结构

- 数组与切片
- Map与Set
- 链表
- 栈与队列

### 高级数据结构

- 树 (二叉树、AVL、红黑树)
- 图
- 堆
- Trie

### 算法

- 排序算法
- 查找算法
- 动态规划
- 图算法

---

## 🚀 快速示例

```go
// 链表
type Node struct {
    Val  int
    Next *Node
}

// 二叉树
type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}

// 图
type Graph map[int][]int
```

---

## 📖 系统文档

- [知识图谱](./00-知识图谱.md)
- [对比矩阵](./00-对比矩阵.md)
- [概念定义体系](./00-概念定义体系.md)

---

**版本**: v1.0
**更新日期**: 2025-10-29
**适用于**: Go 1.25.3

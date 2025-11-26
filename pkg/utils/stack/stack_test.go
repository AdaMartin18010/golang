package stack

import (
	"testing"
)

func TestSimpleStack(t *testing.T) {
	s := NewSimpleStack[int]()
	
	// 测试入栈
	s.Push(1)
	s.Push(2)
	s.Push(3)
	
	// 测试大小
	if s.Size() != 3 {
		t.Errorf("Expected size 3, got %d", s.Size())
	}
	
	// 测试查看
	item, ok := s.Peek()
	if !ok || item != 3 {
		t.Errorf("Expected peek 3, got %d", item)
	}
	
	// 测试出栈
	item, ok = s.Pop()
	if !ok || item != 3 {
		t.Errorf("Expected pop 3, got %d", item)
	}
	
	if s.Size() != 2 {
		t.Errorf("Expected size 2, got %d", s.Size())
	}
	
	// 测试清空
	s.Clear()
	if !s.IsEmpty() {
		t.Error("Expected empty stack")
	}
}

func TestMaxStack(t *testing.T) {
	ms := NewMaxStack[int](func(a, b int) bool {
		return a > b
	})
	
	ms.Push(3)
	ms.Push(1)
	ms.Push(5)
	ms.Push(2)
	
	max, ok := ms.Max()
	if !ok || max != 5 {
		t.Errorf("Expected max 5, got %d", max)
	}
	
	ms.Pop() // 弹出2
	max, ok = ms.Max()
	if !ok || max != 5 {
		t.Errorf("Expected max 5, got %d", max)
	}
	
	ms.Pop() // 弹出5
	max, ok = ms.Max()
	if !ok || max != 3 {
		t.Errorf("Expected max 3, got %d", max)
	}
}

func TestMinStack(t *testing.T) {
	ms := NewMinStack[int](func(a, b int) bool {
		return a < b
	})
	
	ms.Push(3)
	ms.Push(1)
	ms.Push(5)
	ms.Push(2)
	
	min, ok := ms.Min()
	if !ok || min != 1 {
		t.Errorf("Expected min 1, got %d", min)
	}
	
	ms.Pop() // 弹出2
	min, ok = ms.Min()
	if !ok || min != 1 {
		t.Errorf("Expected min 1, got %d", min)
	}
	
	ms.Pop() // 弹出5
	min, ok = ms.Min()
	if !ok || min != 1 {
		t.Errorf("Expected min 1, got %d", min)
	}
	
	ms.Pop() // 弹出1
	min, ok = ms.Min()
	if !ok || min != 3 {
		t.Errorf("Expected min 3, got %d", min)
	}
}


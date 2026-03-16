package queue

import (
	"testing"
)

func TestSimpleQueue(t *testing.T) {
	q := NewSimpleQueue[int]()
	
	// 测试入队
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	
	// 测试大小
	if q.Size() != 3 {
		t.Errorf("Expected size 3, got %d", q.Size())
	}
	
	// 测试查看
	item, ok := q.Peek()
	if !ok || item != 1 {
		t.Errorf("Expected peek 1, got %d", item)
	}
	
	// 测试出队
	item, ok = q.Dequeue()
	if !ok || item != 1 {
		t.Errorf("Expected dequeue 1, got %d", item)
	}
	
	if q.Size() != 2 {
		t.Errorf("Expected size 2, got %d", q.Size())
	}
	
	// 测试清空
	q.Clear()
	if !q.IsEmpty() {
		t.Error("Expected empty queue")
	}
}

func TestPriorityQueue(t *testing.T) {
	pq := NewPriorityQueue[string]()
	
	// 测试入队
	pq.Enqueue("low", 1)
	pq.Enqueue("high", 10)
	pq.Enqueue("medium", 5)
	
	// 测试出队（应该按优先级）
	item, ok := pq.Dequeue()
	if !ok || item != "high" {
		t.Errorf("Expected 'high', got %s", item)
	}
	
	item, ok = pq.Dequeue()
	if !ok || item != "medium" {
		t.Errorf("Expected 'medium', got %s", item)
	}
	
	item, ok = pq.Dequeue()
	if !ok || item != "low" {
		t.Errorf("Expected 'low', got %s", item)
	}
}

func TestCircularQueue(t *testing.T) {
	cq := NewCircularQueue[int](3)
	
	// 测试入队
	if !cq.Enqueue(1) {
		t.Error("Expected enqueue success")
	}
	if !cq.Enqueue(2) {
		t.Error("Expected enqueue success")
	}
	if !cq.Enqueue(3) {
		t.Error("Expected enqueue success")
	}
	
	// 测试队列已满
	if cq.Enqueue(4) {
		t.Error("Expected queue full")
	}
	
	// 测试出队
	item, ok := cq.Dequeue()
	if !ok || item != 1 {
		t.Errorf("Expected dequeue 1, got %d", item)
	}
	
	// 测试可以继续入队
	if !cq.Enqueue(4) {
		t.Error("Expected enqueue success")
	}
	
	// 测试大小
	if cq.Size() != 3 {
		t.Errorf("Expected size 3, got %d", cq.Size())
	}
}


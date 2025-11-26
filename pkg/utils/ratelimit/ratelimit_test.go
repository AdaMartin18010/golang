package ratelimit

import (
	"testing"
	"time"
)

func TestTokenBucket(t *testing.T) {
	tb := NewTokenBucket(10, 5)
	
	// 应该允许前10个请求
	for i := 0; i < 10; i++ {
		if !tb.Allow() {
			t.Errorf("Expected allow at request %d", i)
		}
	}
	
	// 第11个请求应该被拒绝
	if tb.Allow() {
		t.Error("Expected deny after 10 requests")
	}
	
	// 等待一段时间后应该可以继续
	time.Sleep(200 * time.Millisecond)
	if !tb.Allow() {
		t.Error("Expected allow after refill")
	}
}

func TestLeakyBucket(t *testing.T) {
	lb := NewLeakyBucket(10, 5)
	
	// 应该允许前10个请求
	for i := 0; i < 10; i++ {
		if !lb.Allow() {
			t.Errorf("Expected allow at request %d", i)
		}
	}
	
	// 第11个请求应该被拒绝
	if lb.Allow() {
		t.Error("Expected deny after 10 requests")
	}
	
	// 等待一段时间后应该可以继续
	time.Sleep(200 * time.Millisecond)
	if !lb.Allow() {
		t.Error("Expected allow after leak")
	}
}

func TestSlidingWindow(t *testing.T) {
	sw := NewSlidingWindow(1*time.Second, 5)
	
	// 应该允许前5个请求
	for i := 0; i < 5; i++ {
		if !sw.Allow() {
			t.Errorf("Expected allow at request %d", i)
		}
	}
	
	// 第6个请求应该被拒绝
	if sw.Allow() {
		t.Error("Expected deny after 5 requests")
	}
	
	// 等待窗口过期后应该可以继续
	time.Sleep(1100 * time.Millisecond)
	if !sw.Allow() {
		t.Error("Expected allow after window expires")
	}
}

func TestFixedWindow(t *testing.T) {
	fw := NewFixedWindow(1*time.Second, 5)
	
	// 应该允许前5个请求
	for i := 0; i < 5; i++ {
		if !fw.Allow() {
			t.Errorf("Expected allow at request %d", i)
		}
	}
	
	// 第6个请求应该被拒绝
	if fw.Allow() {
		t.Error("Expected deny after 5 requests")
	}
	
	// 等待窗口重置后应该可以继续
	time.Sleep(1100 * time.Millisecond)
	if !fw.Allow() {
		t.Error("Expected allow after window resets")
	}
}


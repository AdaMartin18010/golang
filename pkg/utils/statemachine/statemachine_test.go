package statemachine

import (
	"testing"
)

func TestStateMachine(t *testing.T) {
	sm := NewStateMachine[string, string]("idle")
	
	// 添加状态转换
	sm.AddTransition("idle", "start", "running")
	sm.AddTransition("running", "stop", "idle")
	sm.AddTransition("running", "pause", "paused")
	sm.AddTransition("paused", "resume", "running")
	
	// 检查当前状态
	if sm.Current() != "idle" {
		t.Error("Expected initial state 'idle'")
	}
	
	// 触发事件
	err := sm.Trigger("start")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if sm.Current() != "running" {
		t.Error("Expected state 'running'")
	}
	
	// 检查是否可以触发
	if !sm.CanTrigger("stop") {
		t.Error("Expected can trigger 'stop'")
	}
	
	// 无效转换
	err = sm.Trigger("invalid")
	if err == nil {
		t.Error("Expected error for invalid transition")
	}
}

func TestStateMachineCallbacks(t *testing.T) {
	sm := NewStateMachine[string, string]("idle")
	
	var entered, exited bool
	
	sm.AddTransition("idle", "start", "running")
	sm.OnEnter("running", func() {
		entered = true
	})
	sm.OnExit("idle", func() {
		exited = true
	})
	
	sm.Trigger("start")
	
	if !entered {
		t.Error("Expected onEnter callback")
	}
	if !exited {
		t.Error("Expected onExit callback")
	}
}

func TestSimpleStateMachine(t *testing.T) {
	sm := NewSimpleStateMachine("idle")
	
	sm.AddTransition("idle", "start", "running")
	sm.Trigger("start")
	
	if sm.Current() != "running" {
		t.Error("Expected state 'running'")
	}
}


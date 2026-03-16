package statemachine

import (
	"fmt"
	"sync"
)

// State 状态类型
type State string

// Event 事件类型
type Event string

// Transition 状态转换
type Transition[S comparable, E comparable] struct {
	From  S
	Event E
	To    S
}

// StateMachine 状态机
type StateMachine[S comparable, E comparable] struct {
	current     S
	transitions map[Transition[S, E]]S
	onEnter     map[S]func()
	onExit      map[S]func()
	onTransition map[Transition[S, E]]func()
	mu          sync.RWMutex
}

// NewStateMachine 创建状态机
func NewStateMachine[S comparable, E comparable](initial S) *StateMachine[S, E] {
	return &StateMachine[S, E]{
		current:     initial,
		transitions: make(map[Transition[S, E]]S),
		onEnter:     make(map[S]func()),
		onExit:      make(map[S]func()),
		onTransition: make(map[Transition[S, E]]func()),
	}
}

// AddTransition 添加状态转换
func (sm *StateMachine[S, E]) AddTransition(from S, event E, to S) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	transition := Transition[S, E]{From: from, Event: event, To: to}
	sm.transitions[transition] = to
}

// AddTransitions 批量添加状态转换
func (sm *StateMachine[S, E]) AddTransitions(transitions []Transition[S, E]) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	for _, t := range transitions {
		sm.transitions[t] = t.To
	}
}

// OnEnter 设置进入状态时的回调
func (sm *StateMachine[S, E]) OnEnter(state S, callback func()) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.onEnter[state] = callback
}

// OnExit 设置离开状态时的回调
func (sm *StateMachine[S, E]) OnExit(state S, callback func()) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.onExit[state] = callback
}

// OnTransition 设置状态转换时的回调
func (sm *StateMachine[S, E]) OnTransition(from S, event E, callback func()) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	transition := Transition[S, E]{From: from, Event: event}
	sm.onTransition[transition] = callback
}

// Trigger 触发事件
func (sm *StateMachine[S, E]) Trigger(event E) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	current := sm.current
	transition := Transition[S, E]{From: current, Event: event}
	
	to, ok := sm.transitions[transition]
	if !ok {
		return fmt.Errorf("invalid transition from %v on event %v", current, event)
	}
	
	// 执行退出回调
	if callback, ok := sm.onExit[current]; ok {
		sm.mu.Unlock()
		callback()
		sm.mu.Lock()
	}
	
	// 执行转换回调
	if callback, ok := sm.onTransition[transition]; ok {
		sm.mu.Unlock()
		callback()
		sm.mu.Lock()
	}
	
	// 更新状态
	sm.current = to
	
	// 执行进入回调
	if callback, ok := sm.onEnter[to]; ok {
		sm.mu.Unlock()
		callback()
		sm.mu.Lock()
	}
	
	return nil
}

// Current 获取当前状态
func (sm *StateMachine[S, E]) Current() S {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.current
}

// CanTrigger 检查是否可以触发事件
func (sm *StateMachine[S, E]) CanTrigger(event E) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	current := sm.current
	transition := Transition[S, E]{From: current, Event: event}
	_, ok := sm.transitions[transition]
	return ok
}

// Reset 重置状态机
func (sm *StateMachine[S, E]) Reset(initial S) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.current = initial
}

// GetTransitions 获取所有状态转换
func (sm *StateMachine[S, E]) GetTransitions() []Transition[S, E] {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	transitions := make([]Transition[S, E], 0, len(sm.transitions))
	for t := range sm.transitions {
		transitions = append(transitions, t)
	}
	return transitions
}

// GetAvailableEvents 获取当前状态可用的所有事件
func (sm *StateMachine[S, E]) GetAvailableEvents() []E {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	current := sm.current
	events := make([]E, 0)
	
	for t := range sm.transitions {
		if t.From == current {
			events = append(events, t.Event)
		}
	}
	return events
}

// SimpleStateMachine 简单状态机（使用字符串作为状态和事件）
type SimpleStateMachine struct {
	*StateMachine[string, string]
}

// NewSimpleStateMachine 创建简单状态机
func NewSimpleStateMachine(initial string) *SimpleStateMachine {
	return &SimpleStateMachine{
		StateMachine: NewStateMachine[string, string](initial),
	}
}


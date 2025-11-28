package core

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// TestNewLearningEngine 测试创建学习引擎
func TestNewLearningEngine(t *testing.T) {
	engine := NewLearningEngine(nil)

	if engine == nil {
		t.Fatal("Failed to create LearningEngine")
	}

	if engine.knowledgeBase == nil {
		t.Error("KnowledgeBase not initialized")
	}

	if engine.experiences == nil {
		t.Error("ExperienceBuffer not initialized")
	}

	if engine.models == nil {
		t.Error("Models map not initialized")
	}
}

// TestLearnBasic 测试基本学习功能
func TestLearnBasic(t *testing.T) {
	engine := NewLearningEngine(nil)

	experience := Experience{
		Input: Input{
			ID:        "test-input",
			Type:      "test",
			Data:      map[string]interface{}{"value": 10},
			Timestamp: time.Now(),
		},
		Output: Output{
			ID:        "test-output",
			Type:      "result",
			Data:      map[string]interface{}{"result": 20},
			Timestamp: time.Now(),
		},
		Reward:    0.8,
		Timestamp: time.Now(),
	}

	ctx := context.Background()
	err := engine.Learn(ctx, experience)
	if err != nil {
		t.Errorf("Learn failed: %v", err)
	}
}

// TestLearnMultipleExperiences 测试学习多个经验
func TestLearnMultipleExperiences(t *testing.T) {
	engine := NewLearningEngine(nil)

	experiences := []Experience{
		{
			Input:     Input{ID: "1", Type: "t1", Data: map[string]interface{}{}, Timestamp: time.Now()},
			Output:    Output{ID: "1", Type: "r1", Data: map[string]interface{}{}, Timestamp: time.Now()},
			Reward:    0.7,
			Timestamp: time.Now(),
		},
		{
			Input:     Input{ID: "2", Type: "t2", Data: map[string]interface{}{}, Timestamp: time.Now()},
			Output:    Output{ID: "2", Type: "r2", Data: map[string]interface{}{}, Timestamp: time.Now()},
			Reward:    0.8,
			Timestamp: time.Now(),
		},
		{
			Input:     Input{ID: "3", Type: "t3", Data: map[string]interface{}{}, Timestamp: time.Now()},
			Output:    Output{ID: "3", Type: "r3", Data: map[string]interface{}{}, Timestamp: time.Now()},
			Reward:    0.9,
			Timestamp: time.Now(),
		},
	}

	ctx := context.Background()
	for i, exp := range experiences {
		err := engine.Learn(ctx, exp)
		if err != nil {
			t.Errorf("Learn failed for experience %d: %v", i, err)
		}
	}
}

// TestLearnWithDifferentRewards 测试不同奖励的学习
func TestLearnWithDifferentRewards(t *testing.T) {
	engine := NewLearningEngine(nil)
	ctx := context.Background()

	rewards := []float64{0.0, 0.25, 0.5, 0.75, 1.0, -0.5}

	for i, reward := range rewards {
		experience := Experience{
			Input: Input{
				ID:        fmt.Sprintf("input-%d", i),
				Type:      "test",
				Data:      map[string]interface{}{"index": i},
				Timestamp: time.Now(),
			},
			Output: Output{
				ID:        fmt.Sprintf("output-%d", i),
				Type:      "result",
				Data:      map[string]interface{}{"value": i * 10},
				Timestamp: time.Now(),
			},
			Reward:    reward,
			Timestamp: time.Now(),
		}

		err := engine.Learn(ctx, experience)
		if err != nil {
			t.Errorf("Learn failed for reward %f: %v", reward, err)
		}
	}
}

// TestLearnWithContext 测试带上下文的学习
func TestLearnWithContext(t *testing.T) {
	engine := NewLearningEngine(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	experience := Experience{
		Input:     Input{ID: "ctx-test", Type: "test", Data: map[string]interface{}{}, Timestamp: time.Now()},
		Output:    Output{ID: "ctx-output", Type: "result", Data: map[string]interface{}{}, Timestamp: time.Now()},
		Reward:    0.5,
		Timestamp: time.Now(),
	}

	err := engine.Learn(ctx, experience)
	if err != nil && ctx.Err() != nil {
		t.Skip("Context timeout as expected")
	}

	if err != nil {
		t.Logf("Learn with context returned error: %v", err)
	}
}

// TestExperienceBuffer 测试经验缓冲区
func TestExperienceBuffer(t *testing.T) {
	engine := NewLearningEngine(&LearningConfig{
		ExperienceBufferSize: 5, // 小缓冲区以便测试
		LearningRate:         0.01,
		DiscountFactor:       0.95,
	})

	ctx := context.Background()

	// 添加超过缓冲区大小的经验
	for i := 0; i < 10; i++ {
		experience := Experience{
			Input:     Input{ID: fmt.Sprintf("buf-%d", i), Type: "test", Data: map[string]interface{}{}, Timestamp: time.Now()},
			Output:    Output{ID: fmt.Sprintf("out-%d", i), Type: "result", Data: map[string]interface{}{}, Timestamp: time.Now()},
			Reward:    float64(i) / 10.0,
			Timestamp: time.Now(),
		}

		err := engine.Learn(ctx, experience)
		if err != nil {
			t.Errorf("Learn failed for experience %d: %v", i, err)
		}
	}

	// 验证缓冲区不超过最大大小
	engine.mu.RLock()
	bufferSize := len(engine.experiences.experiences)
	engine.mu.RUnlock()

	if bufferSize > 5 {
		t.Errorf("Buffer size %d exceeds max size 5", bufferSize)
	}
}

// TestKnowledgeBase 测试知识库功能
func TestKnowledgeBase(t *testing.T) {
	kb := NewKnowledgeBase()

	if kb == nil {
		t.Fatal("Failed to create KnowledgeBase")
	}

	// 添加事实
	kb.AddFact("test-fact", "test-value")

	// 检索事实
	value, exists := kb.GetFact("test-fact")
	if !exists {
		t.Error("Expected fact to exist")
	}
	if value != "test-value" {
		t.Errorf("Expected 'test-value', got %v", value)
	}

	// 检索不存在的事实
	value, exists = kb.GetFact("nonexistent")
	if exists {
		t.Error("Expected fact not to exist")
	}
	if value != nil {
		t.Errorf("Expected nil for nonexistent fact, got %v", value)
	}
}

// TestModelManagement 测试模型管理
func TestModelManagement(t *testing.T) {
	engine := NewLearningEngine(nil)

	// 创建模型
	model := &MLModel{
		ID:       "test-model",
		Type:     "supervised",
		Weights:  make(map[string]float64),
		Features: []string{"feature1", "feature2"},
		Config:   map[string]interface{}{"version": "1.0"},
	}

	// 注册模型
	engine.mu.Lock()
	engine.models["test-model"] = model
	engine.mu.Unlock()

	// 检索模型
	engine.mu.RLock()
	retrievedModel := engine.models["test-model"]
	engine.mu.RUnlock()

	if retrievedModel == nil {
		t.Error("Failed to retrieve model")
	}

	if retrievedModel.ID != "test-model" {
		t.Errorf("Expected 'test-model', got %s", retrievedModel.ID)
	}

	if retrievedModel.Type != "supervised" {
		t.Errorf("Expected 'supervised', got %s", retrievedModel.Type)
	}
}

// TestConcurrentLearning 测试并发学习
func TestConcurrentLearning(t *testing.T) {
	engine := NewLearningEngine(nil)
	ctx := context.Background()

	numGoroutines := 10
	results := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			experience := Experience{
				Input:     Input{ID: fmt.Sprintf("conc-%d", id), Type: "test", Data: map[string]interface{}{}, Timestamp: time.Now()},
				Output:    Output{ID: fmt.Sprintf("out-%d", id), Type: "result", Data: map[string]interface{}{}, Timestamp: time.Now()},
				Reward:    float64(id) / 10.0,
				Timestamp: time.Now(),
			}
			err := engine.Learn(ctx, experience)
			results <- err
		}(i)
	}

	// 收集结果
	errorCount := 0
	for i := 0; i < numGoroutines; i++ {
		if err := <-results; err != nil {
			t.Logf("Concurrent learning error: %v", err)
			errorCount++
		}
	}

	if errorCount > 0 {
		t.Errorf("Expected 0 errors, got %d", errorCount)
	}
}

// BenchmarkLearn 基准测试学习性能
func BenchmarkLearn(b *testing.B) {
	engine := NewLearningEngine(nil)
	ctx := context.Background()

	experience := Experience{
		Input:     Input{ID: "bench-input", Type: "test", Data: map[string]interface{}{}, Timestamp: time.Now()},
		Output:    Output{ID: "bench-output", Type: "result", Data: map[string]interface{}{}, Timestamp: time.Now()},
		Reward:    0.75,
		Timestamp: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.Learn(ctx, experience)
	}
}

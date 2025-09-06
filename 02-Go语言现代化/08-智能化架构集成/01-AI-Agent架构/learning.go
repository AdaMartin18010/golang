package main

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

// LearningEngine 学习引擎接口
type LearningEngine interface {
	UpdateModel(experience Experience) error
	Predict(input Input) (Prediction, error)
	GetModelMetrics() ModelMetrics
	Train(data []Experience) error
}

// DecisionEngine 决策引擎接口
type DecisionEngine interface {
	MakeDecision(ctx context.Context, input Input) (Decision, error)
	EvaluateDecision(decision Decision, outcome Outcome) error
	OptimizeStrategy(strategy Strategy) error
}

// Prediction 预测结果
type Prediction struct {
	Confidence float64                `json:"confidence"`
	Result     map[string]interface{} `json:"result"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// Decision 决策结果
type Decision struct {
	Action     string                 `json:"action"`
	Parameters map[string]interface{} `json:"parameters"`
	Confidence float64                `json:"confidence"`
	Reasoning  string                 `json:"reasoning"`
}

// Outcome 决策结果
type Outcome struct {
	Success   bool               `json:"success"`
	Reward    float64            `json:"reward"`
	Metrics   map[string]float64 `json:"metrics"`
	Timestamp time.Time          `json:"timestamp"`
}

// Strategy 策略
type Strategy struct {
	Name       string                 `json:"name"`
	Parameters map[string]interface{} `json:"parameters"`
	Rules      []Rule                 `json:"rules"`
}

// Rule 规则
type Rule struct {
	Condition string  `json:"condition"`
	Action    string  `json:"action"`
	Weight    float64 `json:"weight"`
}

// ModelMetrics 模型指标
type ModelMetrics struct {
	Accuracy    float64   `json:"accuracy"`
	Precision   float64   `json:"precision"`
	Recall      float64   `json:"recall"`
	F1Score     float64   `json:"f1_score"`
	Loss        float64   `json:"loss"`
	LastUpdated time.Time `json:"last_updated"`
}

// SimpleLearningEngine 简单学习引擎实现
type SimpleLearningEngine struct {
	patterns    map[string][]Experience
	weights     map[string]float64
	mu          sync.RWMutex
	lastUpdated time.Time
}

// NewSimpleLearningEngine 创建简单学习引擎
func NewSimpleLearningEngine() *SimpleLearningEngine {
	return &SimpleLearningEngine{
		patterns: make(map[string][]Experience),
		weights:  make(map[string]float64),
	}
}

// UpdateModel 更新模型
func (s *SimpleLearningEngine) UpdateModel(experience Experience) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 根据输入类型分组经验
	patternKey := s.getPatternKey(experience.Input)
	s.patterns[patternKey] = append(s.patterns[patternKey], experience)

	// 保持每个模式的经验数量在合理范围内
	if len(s.patterns[patternKey]) > 100 {
		s.patterns[patternKey] = s.patterns[patternKey][len(s.patterns[patternKey])-100:]
	}

	// 更新权重
	s.updateWeights(patternKey, experience)

	s.lastUpdated = time.Now()
	return nil
}

// Predict 预测
func (s *SimpleLearningEngine) Predict(input Input) (Prediction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	patternKey := s.getPatternKey(input)
	experiences, exists := s.patterns[patternKey]

	if !exists || len(experiences) == 0 {
		return Prediction{
			Confidence: 0.0,
			Result:     map[string]interface{}{"default": "no_data"},
		}, nil
	}

	// 基于历史经验进行简单预测
	successRate := s.calculateSuccessRate(experiences)
	avgFeedback := s.calculateAverageFeedback(experiences)

	confidence := (successRate + avgFeedback) / 2.0

	return Prediction{
		Confidence: confidence,
		Result: map[string]interface{}{
			"success_rate": successRate,
			"avg_feedback": avgFeedback,
			"pattern_key":  patternKey,
			"sample_count": len(experiences),
		},
		Metadata: map[string]interface{}{
			"model_type":   "simple_pattern",
			"last_updated": s.lastUpdated,
		},
	}, nil
}

// GetModelMetrics 获取模型指标
func (s *SimpleLearningEngine) GetModelMetrics() ModelMetrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	totalExperiences := 0
	successfulExperiences := 0
	totalFeedback := 0.0

	for _, experiences := range s.patterns {
		totalExperiences += len(experiences)
		for _, exp := range experiences {
			if exp.Output.Success {
				successfulExperiences++
			}
			totalFeedback += exp.Feedback
		}
	}

	accuracy := 0.0
	if totalExperiences > 0 {
		accuracy = float64(successfulExperiences) / float64(totalExperiences)
	}

	_ = 0.0 // avgFeedback变量已使用
	if totalExperiences > 0 {
		_ = totalFeedback / float64(totalExperiences)
	}

	return ModelMetrics{
		Accuracy:    accuracy,
		Precision:   accuracy, // 简化实现
		Recall:      accuracy, // 简化实现
		F1Score:     accuracy, // 简化实现
		Loss:        1.0 - accuracy,
		LastUpdated: s.lastUpdated,
	}
}

// Train 训练模型
func (s *SimpleLearningEngine) Train(data []Experience) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 清空现有数据
	s.patterns = make(map[string][]Experience)
	s.weights = make(map[string]float64)

	// 重新训练
	for _, experience := range data {
		patternKey := s.getPatternKey(experience.Input)
		s.patterns[patternKey] = append(s.patterns[patternKey], experience)
		s.updateWeights(patternKey, experience)
	}

	s.lastUpdated = time.Now()
	return nil
}

// getPatternKey 获取模式键
func (s *SimpleLearningEngine) getPatternKey(input Input) string {
	// 基于输入类型和关键字段生成模式键
	return fmt.Sprintf("%s_%v", input.Type, input.Priority)
}

// updateWeights 更新权重
func (s *SimpleLearningEngine) updateWeights(patternKey string, experience Experience) {
	// 简单的权重更新策略
	currentWeight := s.weights[patternKey]
	feedback := experience.Feedback

	// 使用指数移动平均更新权重
	alpha := 0.1
	newWeight := currentWeight*(1-alpha) + feedback*alpha
	s.weights[patternKey] = newWeight
}

// calculateSuccessRate 计算成功率
func (s *SimpleLearningEngine) calculateSuccessRate(experiences []Experience) float64 {
	if len(experiences) == 0 {
		return 0.0
	}

	successCount := 0
	for _, exp := range experiences {
		if exp.Output.Success {
			successCount++
		}
	}

	return float64(successCount) / float64(len(experiences))
}

// calculateAverageFeedback 计算平均反馈
func (s *SimpleLearningEngine) calculateAverageFeedback(experiences []Experience) float64 {
	if len(experiences) == 0 {
		return 0.0
	}

	totalFeedback := 0.0
	for _, exp := range experiences {
		totalFeedback += exp.Feedback
	}

	return totalFeedback / float64(len(experiences))
}

// SimpleDecisionEngine 简单决策引擎实现
type SimpleDecisionEngine struct {
	rules       []Rule
	strategies  map[string]Strategy
	mu          sync.RWMutex
	performance map[string]float64
}

// NewSimpleDecisionEngine 创建简单决策引擎
func NewSimpleDecisionEngine() *SimpleDecisionEngine {
	return &SimpleDecisionEngine{
		rules:       make([]Rule, 0),
		strategies:  make(map[string]Strategy),
		performance: make(map[string]float64),
	}
}

// MakeDecision 做决策
func (s *SimpleDecisionEngine) MakeDecision(ctx context.Context, input Input) (Decision, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 基于规则进行决策
	for _, rule := range s.rules {
		if s.evaluateCondition(rule.Condition, input) {
			return Decision{
				Action:     rule.Action,
				Parameters: map[string]interface{}{"rule_weight": rule.Weight},
				Confidence: rule.Weight,
				Reasoning:  fmt.Sprintf("Matched rule: %s", rule.Condition),
			}, nil
		}
	}

	// 默认决策
	return Decision{
		Action:     "default",
		Parameters: map[string]interface{}{"fallback": true},
		Confidence: 0.5,
		Reasoning:  "No matching rule found, using default action",
	}, nil
}

// EvaluateDecision 评估决策
func (s *SimpleDecisionEngine) EvaluateDecision(decision Decision, outcome Outcome) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 更新决策性能
	action := decision.Action
	currentPerf := s.performance[action]

	// 使用指数移动平均更新性能
	alpha := 0.1
	newPerf := currentPerf*(1-alpha) + outcome.Reward*alpha
	s.performance[action] = newPerf

	return nil
}

// OptimizeStrategy 优化策略
func (s *SimpleDecisionEngine) OptimizeStrategy(strategy Strategy) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 基于性能数据优化规则权重
	for i, rule := range s.rules {
		if perf, exists := s.performance[rule.Action]; exists {
			// 根据性能调整权重
			s.rules[i].Weight = math.Max(0.1, math.Min(1.0, perf))
		}
	}

	return nil
}

// AddRule 添加规则
func (s *SimpleDecisionEngine) AddRule(rule Rule) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.rules = append(s.rules, rule)
}

// evaluateCondition 评估条件
func (s *SimpleDecisionEngine) evaluateCondition(condition string, input Input) bool {
	// 简单的条件评估实现
	switch condition {
	case "high_priority":
		return input.Priority > 5
	case "low_priority":
		return input.Priority <= 2
	case "urgent":
		return input.Priority >= 8
	default:
		return false
	}
}

// MetricsCollector 指标收集器接口
type MetricsCollector interface {
	RecordProcess(duration time.Duration, success bool)
	GetMetrics() map[string]float64
	Reset()
}

// SimpleMetricsCollector 简单指标收集器实现
type SimpleMetricsCollector struct {
	processTimes []time.Duration
	successCount int64
	totalCount   int64
	mu           sync.RWMutex
}

// NewSimpleMetricsCollector 创建简单指标收集器
func NewSimpleMetricsCollector() *SimpleMetricsCollector {
	return &SimpleMetricsCollector{
		processTimes: make([]time.Duration, 0),
	}
}

// RecordProcess 记录处理过程
func (s *SimpleMetricsCollector) RecordProcess(duration time.Duration, success bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.processTimes = append(s.processTimes, duration)
	s.totalCount++

	if success {
		s.successCount++
	}

	// 保持最近1000次的记录
	if len(s.processTimes) > 1000 {
		s.processTimes = s.processTimes[len(s.processTimes)-1000:]
	}
}

// GetMetrics 获取指标
func (s *SimpleMetricsCollector) GetMetrics() map[string]float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metrics := make(map[string]float64)

	// 计算平均处理时间
	if len(s.processTimes) > 0 {
		totalDuration := time.Duration(0)
		for _, duration := range s.processTimes {
			totalDuration += duration
		}
		metrics["avg_process_time"] = float64(totalDuration.Milliseconds()) / float64(len(s.processTimes))
	} else {
		metrics["avg_process_time"] = 0.0
	}

	// 计算错误率
	if s.totalCount > 0 {
		metrics["error_rate"] = float64(s.totalCount-s.successCount) / float64(s.totalCount)
		metrics["success_rate"] = float64(s.successCount) / float64(s.totalCount)
	} else {
		metrics["error_rate"] = 0.0
		metrics["success_rate"] = 0.0
	}

	// 计算吞吐量（每秒处理数）
	metrics["throughput"] = float64(s.totalCount) / 60.0 // 假设1分钟窗口

	return metrics
}

// Reset 重置指标
func (s *SimpleMetricsCollector) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.processTimes = make([]time.Duration, 0)
	s.successCount = 0
	s.totalCount = 0
}

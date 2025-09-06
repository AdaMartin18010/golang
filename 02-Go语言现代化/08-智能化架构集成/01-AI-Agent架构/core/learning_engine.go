package ai_agent

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

// LearningEngine 自适应学习引擎
type LearningEngine struct {
	knowledgeBase *KnowledgeBase
	models        map[string]*MLModel
	experiences   *ExperienceBuffer
	config        *LearningConfig
	mu            sync.RWMutex
}

// LearningConfig 学习配置
type LearningConfig struct {
	// 学习率
	LearningRate float64 `json:"learning_rate"`
	// 折扣因子
	DiscountFactor float64 `json:"discount_factor"`
	// 探索率
	ExplorationRate float64 `json:"exploration_rate"`
	// 经验缓冲区大小
	ExperienceBufferSize int `json:"experience_buffer_size"`
	// 批量大小
	BatchSize int `json:"batch_size"`
	// 训练频率
	TrainingFrequency time.Duration `json:"training_frequency"`
	// 模型保存路径
	ModelSavePath string `json:"model_save_path"`
}

// KnowledgeBase 知识库
type KnowledgeBase struct {
	facts     map[string]interface{}
	rules     []Rule
	patterns  map[string]Pattern
	mu        sync.RWMutex
}

// Rule 规则
type Rule struct {
	ID          string                 `json:"id"`
	Condition   map[string]interface{} `json:"condition"`
	Action      string                 `json:"action"`
	Confidence  float64                `json:"confidence"`
	UsageCount  int                    `json:"usage_count"`
	LastUsed    time.Time              `json:"last_used"`
	CreatedAt   time.Time              `json:"created_at"`
}

// Pattern 模式
type Pattern struct {
	ID          string                 `json:"id"`
	Sequence    []string               `json:"sequence"`
	Frequency   int                    `json:"frequency"`
	SuccessRate float64                `json:"success_rate"`
	LastSeen    time.Time              `json:"last_seen"`
}

// MLModel 机器学习模型
type MLModel struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Weights  map[string]float64     `json:"weights"`
	Bias     float64                `json:"bias"`
	Features []string               `json:"features"`
	Metrics  *ModelMetrics          `json:"metrics"`
	Config   map[string]interface{} `json:"config"`
}

// ModelMetrics 模型指标
type ModelMetrics struct {
	Accuracy    float64 `json:"accuracy"`
	Loss        float64 `json:"loss"`
	Precision   float64 `json:"precision"`
	Recall      float64 `json:"recall"`
	F1Score     float64 `json:"f1_score"`
	TrainingCount int   `json:"training_count"`
	LastUpdated time.Time `json:"last_updated"`
}

// Experience 经验
type Experience struct {
	State      map[string]interface{} `json:"state"`
	Action     string                 `json:"action"`
	Reward     float64                `json:"reward"`
	NextState  map[string]interface{} `json:"next_state"`
	Done       bool                   `json:"done"`
	Timestamp  time.Time              `json:"timestamp"`
	AgentID    string                 `json:"agent_id"`
}

// ExperienceBuffer 经验缓冲区
type ExperienceBuffer struct {
	experiences []Experience
	maxSize     int
	mu          sync.RWMutex
}

// NewLearningEngine 创建学习引擎
func NewLearningEngine(config *LearningConfig) *LearningEngine {
	if config == nil {
		config = &LearningConfig{
			LearningRate:         0.01,
			DiscountFactor:       0.95,
			ExplorationRate:      0.1,
			ExperienceBufferSize: 10000,
			BatchSize:            32,
			TrainingFrequency:    5 * time.Minute,
			ModelSavePath:        "./models",
		}
	}

	return &LearningEngine{
		knowledgeBase: NewKnowledgeBase(),
		models:        make(map[string]*MLModel),
		experiences:   NewExperienceBuffer(config.ExperienceBufferSize),
		config:        config,
	}
}

// NewKnowledgeBase 创建知识库
func NewKnowledgeBase() *KnowledgeBase {
	return &KnowledgeBase{
		facts:    make(map[string]interface{}),
		rules:    make([]Rule, 0),
		patterns: make(map[string]Pattern),
	}
}

// NewExperienceBuffer 创建经验缓冲区
func NewExperienceBuffer(maxSize int) *ExperienceBuffer {
	return &ExperienceBuffer{
		experiences: make([]Experience, 0, maxSize),
		maxSize:     maxSize,
	}
}

// Learn 学习过程
func (le *LearningEngine) Learn(ctx context.Context, experience Experience) error {
	le.mu.Lock()
	defer le.mu.Unlock()

	// 1. 存储经验
	le.experiences.Add(experience)

	// 2. 更新知识库
	le.updateKnowledgeBase(experience)

	// 3. 检测模式
	le.detectPatterns(experience)

	// 4. 训练模型
	if le.shouldTrain() {
		return le.trainModels(ctx)
	}

	return nil
}

// AddFact 添加事实
func (kb *KnowledgeBase) AddFact(key string, value interface{}) {
	kb.mu.Lock()
	defer kb.mu.Unlock()
	kb.facts[key] = value
}

// GetFact 获取事实
func (kb *KnowledgeBase) GetFact(key string) (interface{}, bool) {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	value, exists := kb.facts[key]
	return value, exists
}

// AddRule 添加规则
func (kb *KnowledgeBase) AddRule(rule Rule) {
	kb.mu.Lock()
	defer kb.mu.Unlock()
	rule.CreatedAt = time.Now()
	kb.rules = append(kb.rules, rule)
}

// FindMatchingRules 查找匹配的规则
func (kb *KnowledgeBase) FindMatchingRules(state map[string]interface{}) []Rule {
	kb.mu.RLock()
	defer kb.mu.RUnlock()

	var matchingRules []Rule
	for _, rule := range kb.rules {
		if kb.ruleMatches(rule, state) {
			matchingRules = append(matchingRules, rule)
		}
	}

	// 按置信度排序
	sortRulesByConfidence(matchingRules)
	return matchingRules
}

// ruleMatches 检查规则是否匹配状态
func (kb *KnowledgeBase) ruleMatches(rule Rule, state map[string]interface{}) bool {
	for key, expectedValue := range rule.Condition {
		if actualValue, exists := state[key]; !exists || actualValue != expectedValue {
			return false
		}
	}
	return true
}

// Add 添加经验到缓冲区
func (eb *ExperienceBuffer) Add(experience Experience) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if len(eb.experiences) >= eb.maxSize {
		// 移除最旧的经验
		eb.experiences = eb.experiences[1:]
	}

	eb.experiences = append(eb.experiences, experience)
}

// Sample 采样经验
func (eb *ExperienceBuffer) Sample(batchSize int) []Experience {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if len(eb.experiences) == 0 {
		return nil
	}

	if batchSize > len(eb.experiences) {
		batchSize = len(eb.experiences)
	}

	sample := make([]Experience, batchSize)
	indices := rand.Perm(len(eb.experiences))[:batchSize]

	for i, idx := range indices {
		sample[i] = eb.experiences[idx]
	}

	return sample
}

// updateKnowledgeBase 更新知识库
func (le *LearningEngine) updateKnowledgeBase(experience Experience) {
	// 更新事实
	for key, value := range experience.State {
		le.knowledgeBase.AddFact(fmt.Sprintf("state_%s", key), value)
	}

	// 创建新规则
	if experience.Reward > 0 {
		rule := Rule{
			ID:         generateRuleID(),
			Condition:  experience.State,
			Action:     experience.Action,
			Confidence: math.Min(experience.Reward/100.0, 1.0),
			UsageCount: 1,
			LastUsed:   time.Now(),
		}
		le.knowledgeBase.AddRule(rule)
	}
}

// detectPatterns 检测模式
func (le *LearningEngine) detectPatterns(experience Experience) {
	le.knowledgeBase.mu.Lock()
	defer le.knowledgeBase.mu.Unlock()

	// 简化的模式检测：基于状态转换序列
	patternKey := fmt.Sprintf("%s->%s", experience.Action, experience.NextState["action"])
	
	if pattern, exists := le.knowledgeBase.patterns[patternKey]; exists {
		pattern.Frequency++
		pattern.LastSeen = time.Now()
		if experience.Reward > 0 {
			pattern.SuccessRate = (pattern.SuccessRate*float64(pattern.Frequency-1) + 1.0) / float64(pattern.Frequency)
		} else {
			pattern.SuccessRate = (pattern.SuccessRate * float64(pattern.Frequency-1)) / float64(pattern.Frequency)
		}
		le.knowledgeBase.patterns[patternKey] = pattern
	} else {
		le.knowledgeBase.patterns[patternKey] = Pattern{
			ID:          patternKey,
			Sequence:    []string{experience.Action},
			Frequency:   1,
			SuccessRate: math.Max(experience.Reward/100.0, 0.0),
			LastSeen:    time.Now(),
		}
	}
}

// shouldTrain 判断是否应该训练
func (le *LearningEngine) shouldTrain() bool {
	return len(le.experiences.experiences) >= le.config.BatchSize
}

// trainModels 训练模型
func (le *LearningEngine) trainModels(ctx context.Context) error {
	// 获取训练数据
	batch := le.experiences.Sample(le.config.BatchSize)
	if len(batch) == 0 {
		return nil
	}

	// 训练Q-Learning模型
	if err := le.trainQLearningModel(batch); err != nil {
		return fmt.Errorf("failed to train Q-learning model: %w", err)
	}

	// 训练监督学习模型
	if err := le.trainSupervisedModel(batch); err != nil {
		return fmt.Errorf("failed to train supervised model: %w", err)
	}

	return nil
}

// trainQLearningModel 训练Q-Learning模型
func (le *LearningEngine) trainQLearningModel(batch []Experience) error {
	modelID := "q_learning"
	model, exists := le.models[modelID]
	if !exists {
		model = &MLModel{
			ID:      modelID,
			Type:    "q_learning",
			Weights: make(map[string]float64),
			Bias:    0.0,
			Metrics: &ModelMetrics{},
		}
		le.models[modelID] = model
	}

	// 简化的Q-Learning更新
	for _, exp := range batch {
		stateKey := le.stateToKey(exp.State)
		nextStateKey := le.stateToKey(exp.NextState)
		
		currentQ := model.Weights[stateKey]
		nextQ := model.Weights[nextStateKey]
		
		// Q-Learning更新公式
		targetQ := exp.Reward + le.config.DiscountFactor*nextQ
		newQ := currentQ + le.config.LearningRate*(targetQ-currentQ)
		
		model.Weights[stateKey] = newQ
	}

	// 更新指标
	model.Metrics.TrainingCount++
	model.Metrics.LastUpdated = time.Now()

	return nil
}

// trainSupervisedModel 训练监督学习模型
func (le *LearningEngine) trainSupervisedModel(batch []Experience) error {
	modelID := "supervised"
	model, exists := le.models[modelID]
	if !exists {
		model = &MLModel{
			ID:      modelID,
			Type:    "supervised",
			Weights: make(map[string]float64),
			Bias:    0.0,
			Metrics: &ModelMetrics{},
		}
		le.models[modelID] = model
	}

	// 简化的线性回归训练
	var totalLoss float64
	for _, exp := range batch {
		features := le.extractFeatures(exp.State)
		prediction := le.predict(model, features)
		target := exp.Reward
		
		// 计算损失
		loss := math.Pow(prediction-target, 2)
		totalLoss += loss
		
		// 梯度下降更新
		gradient := 2 * (prediction - target)
		for feature, value := range features {
			model.Weights[feature] -= le.config.LearningRate * gradient * value
		}
		model.Bias -= le.config.LearningRate * gradient
	}

	// 更新指标
	model.Metrics.Loss = totalLoss / float64(len(batch))
	model.Metrics.TrainingCount++
	model.Metrics.LastUpdated = time.Now()

	return nil
}

// Predict 预测最佳行动
func (le *LearningEngine) Predict(state map[string]interface{}) (string, float64, error) {
	le.mu.RLock()
	defer le.mu.RUnlock()

	// 1. 基于规则的预测
	if action, confidence := le.predictByRules(state); confidence > 0.8 {
		return action, confidence, nil
	}

	// 2. 基于模式的预测
	if action, confidence := le.predictByPatterns(state); confidence > 0.7 {
		return action, confidence, nil
	}

	// 3. 基于模型的预测
	if action, confidence := le.predictByModels(state); confidence > 0.6 {
		return action, confidence, nil
	}

	// 4. 探索性行动
	if rand.Float64() < le.config.ExplorationRate {
		return le.getRandomAction(), 0.1, nil
	}

	// 5. 默认行动
	return "default_action", 0.5, nil
}

// predictByRules 基于规则预测
func (le *LearningEngine) predictByRules(state map[string]interface{}) (string, float64) {
	rules := le.knowledgeBase.FindMatchingRules(state)
	if len(rules) == 0 {
		return "", 0.0
	}

	// 返回最高置信度的规则
	bestRule := rules[0]
	return bestRule.Action, bestRule.Confidence
}

// predictByPatterns 基于模式预测
func (le *LearningEngine) predictByPatterns(state map[string]interface{}) (string, float64) {
	le.knowledgeBase.mu.RLock()
	defer le.knowledgeBase.mu.RUnlock()

	var bestPattern Pattern
	var bestConfidence float64

	for _, pattern := range le.knowledgeBase.patterns {
		if pattern.SuccessRate > bestConfidence {
			bestPattern = pattern
			bestConfidence = pattern.SuccessRate
		}
	}

	if len(bestPattern.Sequence) > 0 {
		return bestPattern.Sequence[0], bestConfidence
	}

	return "", 0.0
}

// predictByModels 基于模型预测
func (le *LearningEngine) predictByModels(state map[string]interface{}) (string, float64) {
	var bestAction string
	var bestQValue float64

	// 使用Q-Learning模型
	if model, exists := le.models["q_learning"]; exists {
		stateKey := le.stateToKey(state)
		if qValue, exists := model.Weights[stateKey]; exists && qValue > bestQValue {
			bestQValue = qValue
			bestAction = "q_learning_action"
		}
	}

	// 使用监督学习模型
	if model, exists := le.models["supervised"]; exists {
		features := le.extractFeatures(state)
		prediction := le.predict(model, features)
		if prediction > bestQValue {
			bestQValue = prediction
			bestAction = "supervised_action"
		}
	}

	return bestAction, bestQValue
}

// predict 模型预测
func (le *LearningEngine) predict(model *MLModel, features map[string]float64) float64 {
	prediction := model.Bias
	for feature, value := range features {
		if weight, exists := model.Weights[feature]; exists {
			prediction += weight * value
		}
	}
	return prediction
}

// extractFeatures 提取特征
func (le *LearningEngine) extractFeatures(state map[string]interface{}) map[string]float64 {
	features := make(map[string]float64)
	for key, value := range state {
		if num, ok := value.(float64); ok {
			features[key] = num
		} else if num, ok := value.(int); ok {
			features[key] = float64(num)
		} else {
			// 将非数值特征转换为数值
			features[key] = float64(len(fmt.Sprintf("%v", value)))
		}
	}
	return features
}

// stateToKey 状态转键
func (le *LearningEngine) stateToKey(state map[string]interface{}) string {
	data, _ := json.Marshal(state)
	return string(data)
}

// getRandomAction 获取随机行动
func (le *LearningEngine) getRandomAction() string {
	actions := []string{"action_1", "action_2", "action_3", "action_4", "action_5"}
	return actions[rand.Intn(len(actions))]
}

// GetKnowledge 获取知识
func (le *LearningEngine) GetKnowledge() map[string]interface{} {
	le.mu.RLock()
	defer le.mu.RUnlock()

	return map[string]interface{}{
		"facts_count":    len(le.knowledgeBase.facts),
		"rules_count":    len(le.knowledgeBase.rules),
		"patterns_count": len(le.knowledgeBase.patterns),
		"models_count":   len(le.models),
		"experiences_count": len(le.experiences.experiences),
	}
}

// SaveModel 保存模型
func (le *LearningEngine) SaveModel(modelID string) error {
	le.mu.RLock()
	model, exists := le.models[modelID]
	le.mu.RUnlock()

	if !exists {
		return fmt.Errorf("model %s not found", modelID)
	}

	// 这里应该实现实际的模型保存逻辑
	// 例如保存到文件或数据库
	return nil
}

// LoadModel 加载模型
func (le *LearningEngine) LoadModel(modelID string) error {
	// 这里应该实现实际的模型加载逻辑
	// 例如从文件或数据库加载
	return nil
}

// generateRuleID 生成规则ID
func generateRuleID() string {
	return fmt.Sprintf("rule_%d", time.Now().UnixNano())
}

// sortRulesByConfidence 按置信度排序规则
func sortRulesByConfidence(rules []Rule) {
	// 简单的冒泡排序
	for i := 0; i < len(rules)-1; i++ {
		for j := 0; j < len(rules)-i-1; j++ {
			if rules[j].Confidence < rules[j+1].Confidence {
				rules[j], rules[j+1] = rules[j+1], rules[j]
			}
		}
	}
}

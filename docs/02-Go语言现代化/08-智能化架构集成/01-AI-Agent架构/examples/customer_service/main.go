package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"coordination"
	"core"
)

// CustomerServiceAgent 智能客服代理
type CustomerServiceAgent struct {
	*core.BaseAgent
	nlp       NLPEngine
	knowledge KnowledgeBase
	sentiment SentimentAnalyzer
}

// NewCustomerServiceAgent 创建新的客服代理
func NewCustomerServiceAgent(id string, config core.AgentConfig) *CustomerServiceAgent {
	baseAgent := core.NewBaseAgent(id, config)
	return &CustomerServiceAgent{
		BaseAgent: baseAgent,
		nlp:       NewDefaultNLPEngine(),
		knowledge: NewDefaultKnowledgeBase(),
		sentiment: NewDefaultSentimentAnalyzer(),
	}
}

// Process 处理客服请求
func (a *CustomerServiceAgent) Process(input core.Input) (core.Output, error) {
	// 解析用户消息
	message, ok := input.Data.(string)
	if !ok {
		return core.Output{}, fmt.Errorf("invalid input data type")
	}

	// 自然语言理解
	intent := a.nlp.Understand(message)

	// 情感分析
	sentiment := a.sentiment.Analyze(message)

	// 知识库查询
	response := a.knowledge.Query(intent, sentiment)

	// 个性化回复
	personalized := a.personalize(response, input.Metadata)

	// 构建输出
	output := core.Output{
		ID:        input.ID,
		Type:      "customer_service_response",
		Data:      personalized,
		Timestamp: time.Now(),
	}

	return output, nil
}

// 内部方法

// personalize 个性化回复
func (a *CustomerServiceAgent) personalize(response string, metadata map[string]interface{}) string {
	// 这里实现个性化逻辑
	// 可以根据用户历史、偏好等进行个性化
	return response
}

// NLPEngine 自然语言处理引擎接口
type NLPEngine interface {
	Understand(text string) Intent
}

// Intent 用户意图
type Intent struct {
	Type       string            `json:"type"`
	Confidence float64           `json:"confidence"`
	Entities   map[string]string `json:"entities"`
}

// DefaultNLPEngine 默认NLP引擎实现
type DefaultNLPEngine struct{}

// NewDefaultNLPEngine 创建默认NLP引擎
func NewDefaultNLPEngine() *DefaultNLPEngine {
	return &DefaultNLPEngine{}
}

// Understand 理解用户意图
func (n *DefaultNLPEngine) Understand(text string) Intent {
	// 这里实现简单的意图识别逻辑
	// 在实际应用中，可以使用更复杂的NLP模型
	switch {
	case contains(text, "hello", "hi", "你好"):
		return Intent{Type: "greeting", Confidence: 0.9}
	case contains(text, "help", "帮助", "问题"):
		return Intent{Type: "help_request", Confidence: 0.8}
	case contains(text, "order", "订单", "购买"):
		return Intent{Type: "order_inquiry", Confidence: 0.7}
	default:
		return Intent{Type: "general", Confidence: 0.5}
	}
}

// SentimentAnalyzer 情感分析器接口
type SentimentAnalyzer interface {
	Analyze(text string) Sentiment
}

// Sentiment 情感分析结果
type Sentiment struct {
	Type       string  `json:"type"`  // positive, negative, neutral
	Score      float64 `json:"score"` // -1.0 to 1.0
	Confidence float64 `json:"confidence"`
}

// DefaultSentimentAnalyzer 默认情感分析器
type DefaultSentimentAnalyzer struct{}

// NewDefaultSentimentAnalyzer 创建默认情感分析器
func NewDefaultSentimentAnalyzer() *DefaultSentimentAnalyzer {
	return &DefaultSentimentAnalyzer{}
}

// Analyze 分析情感
func (s *DefaultSentimentAnalyzer) Analyze(text string) Sentiment {
	// 这里实现简单的情感分析逻辑
	// 在实际应用中，可以使用更复杂的情感分析模型
	positiveWords := []string{"good", "great", "excellent", "满意", "好", "棒"}
	negativeWords := []string{"bad", "terrible", "awful", "差", "坏", "糟糕"}

	score := 0.0
	for _, word := range positiveWords {
		if contains(text, word) {
			score += 0.2
		}
	}
	for _, word := range negativeWords {
		if contains(text, word) {
			score -= 0.2
		}
	}

	// 限制分数范围
	if score > 1.0 {
		score = 1.0
	} else if score < -1.0 {
		score = -1.0
	}

	var sentimentType string
	if score > 0.1 {
		sentimentType = "positive"
	} else if score < -0.1 {
		sentimentType = "negative"
	} else {
		sentimentType = "neutral"
	}

	return Sentiment{
		Type:       sentimentType,
		Score:      score,
		Confidence: 0.7,
	}
}

// KnowledgeBase 知识库接口
type KnowledgeBase interface {
	Query(intent Intent, sentiment Sentiment) string
}

// DefaultKnowledgeBase 默认知识库
type DefaultKnowledgeBase struct {
	responses map[string]string
}

// NewDefaultKnowledgeBase 创建默认知识库
func NewDefaultKnowledgeBase() *DefaultKnowledgeBase {
	return &DefaultKnowledgeBase{
		responses: map[string]string{
			"greeting":      "您好！我是智能客服助手，很高兴为您服务。",
			"help_request":  "我可以帮助您解答问题、查询订单、处理退款等。请告诉我您需要什么帮助？",
			"order_inquiry": "您想查询订单信息吗？请提供您的订单号，我来帮您查询。",
			"general":       "我理解您的问题，让我为您提供帮助。",
		},
	}
}

// Query 查询知识库
func (k *DefaultKnowledgeBase) Query(intent Intent, sentiment Sentiment) string {
	// 根据意图和情感选择合适的回复
	response, exists := k.responses[intent.Type]
	if !exists {
		response = k.responses["general"]
	}

	// 根据情感调整回复
	if sentiment.Type == "negative" {
		response = "非常抱歉给您带来了不好的体验。" + response
	} else if sentiment.Type == "positive" {
		response = "很高兴您对我们的服务满意！" + response
	}

	return response
}

// 辅助函数
func contains(text, substr string) bool {
	return len(text) >= len(substr) &&
		(text == substr ||
			(len(text) > len(substr) &&
				(text[:len(substr)] == substr ||
					text[len(text)-len(substr):] == substr)))
}

// HTTP处理器
type CustomerServiceHandler struct {
	coordinator *coordination.SmartCoordinator
}

// NewCustomerServiceHandler 创建客服处理器
func NewCustomerServiceHandler(coordinator *coordination.SmartCoordinator) *CustomerServiceHandler {
	return &CustomerServiceHandler{
		coordinator: coordinator,
	}
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Message  string                 `json:"message"`
	UserID   string                 `json:"user_id"`
	Metadata map[string]interface{} `json:"metadata"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Response string                 `json:"response"`
	AgentID  string                 `json:"agent_id"`
	Metadata map[string]interface{} `json:"metadata"`
}

// HandleChat 处理聊天请求
func (h *CustomerServiceHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 创建任务
	task := coordination.Task{
		ID:       fmt.Sprintf("chat_%d", time.Now().UnixNano()),
		Type:     "customer_service",
		Priority: 1,
		Input: core.Input{
			ID:        task.ID,
			Type:      "chat_message",
			Data:      req.Message,
			Metadata:  req.Metadata,
			Timestamp: time.Now(),
		},
		CreatedAt: time.Now(),
	}

	// 路由任务
	agent, err := h.coordinator.RouteTask(task)
	if err != nil {
		http.Error(w, "No available agents", http.StatusServiceUnavailable)
		return
	}

	// 处理任务
	output, err := agent.Process(task.Input)
	if err != nil {
		http.Error(w, "Processing error", http.StatusInternalServerError)
		return
	}

	// 构建响应
	response := ChatResponse{
		Response: output.Data.(string),
		AgentID:  agent.ID(),
		Metadata: output.Metadata,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleStatus 处理状态查询
func (h *CustomerServiceHandler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := h.coordinator.GetSystemStatus()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func main() {
	fmt.Println("🤖 Starting AI Customer Service System...")

	// 创建协调器
	coordinator := coordination.NewSmartCoordinator()

	// 创建客服代理
	agent1 := NewCustomerServiceAgent("cs_agent_1", core.AgentConfig{})
	agent2 := NewCustomerServiceAgent("cs_agent_2", core.AgentConfig{})
	agent3 := NewCustomerServiceAgent("cs_agent_3", core.AgentConfig{})

	// 注册代理
	if err := coordinator.RegisterAgent(agent1); err != nil {
		log.Fatalf("Failed to register agent 1: %v", err)
	}
	if err := coordinator.RegisterAgent(agent2); err != nil {
		log.Fatalf("Failed to register agent 2: %v", err)
	}
	if err := coordinator.RegisterAgent(agent3); err != nil {
		log.Fatalf("Failed to register agent 3: %v", err)
	}

	// 创建HTTP处理器
	handler := NewCustomerServiceHandler(coordinator)

	// 设置路由
	http.HandleFunc("/chat", handler.HandleChat)
	http.HandleFunc("/status", handler.HandleStatus)

	// 健康检查
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","service":"ai_customer_service"}`)
	})

	fmt.Println("📡 Server starting on :8080")
	fmt.Println("🌐 Available endpoints:")
	fmt.Println("  POST   /chat    - Send chat message")
	fmt.Println("  GET    /status  - Get system status")
	fmt.Println("  GET    /health  - Health check")

	// 启动服务器
	log.Fatal(http.ListenAndServe(":8080", nil))
}

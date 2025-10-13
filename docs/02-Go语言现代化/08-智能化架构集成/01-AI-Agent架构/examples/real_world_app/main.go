package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/redis/go-redis/v9"
)

// 实际应用场景：智能客服系统
// 集成AI-Agent、云原生、性能优化等技术

// 应用配置
type AppConfig struct {
    Server   ServerConfig   `json:"server"`
    Database DatabaseConfig `json:"database"`
    Redis    RedisConfig    `json:"redis"`
    AI       AIConfig       `json:"ai"`
    Metrics  MetricsConfig  `json:"metrics"`
}

type ServerConfig struct {
    Port         string        `json:"port"`
    ReadTimeout  time.Duration `json:"read_timeout"`
    WriteTimeout time.Duration `json:"write_timeout"`
}

type DatabaseConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Database string `json:"database"`
    Username string `json:"username"`
    Password string `json:"password"`
}

type RedisConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Password string `json:"password"`
    DB       int    `json:"db"`
}

type AIConfig struct {
    ModelEndpoint string  `json:"model_endpoint"`
    APIKey        string  `json:"api_key"`
    MaxTokens     int     `json:"max_tokens"`
    Temperature   float64 `json:"temperature"`
}

type MetricsConfig struct {
    Enabled bool   `json:"enabled"`
    Port    string `json:"port"`
}

// 智能客服系统
type SmartCustomerService struct {
    config     *AppConfig
    aiAgent    *AIAgent
    redis      *redis.Client
    router     *gin.Engine
    metrics    *MetricsCollector
    userRepo   *UserRepository
    ticketRepo *TicketRepository
}

// 用户实体
type User struct {
    ID       string    `json:"id" db:"id"`
    Email    string    `json:"email" db:"email"`
    Name     string    `json:"name" db:"name"`
    Created  time.Time `json:"created" db:"created_at"`
    Updated  time.Time `json:"updated" db:"updated_at"`
    Metadata map[string]interface{} `json:"metadata" db:"metadata"`
}

// 工单实体
type Ticket struct {
    ID          string                 `json:"id" db:"id"`
    UserID      string                 `json:"user_id" db:"user_id"`
    Subject     string                 `json:"subject" db:"subject"`
    Description string                 `json:"description" db:"description"`
    Status      string                 `json:"status" db:"status"`
    Priority    string                 `json:"priority" db:"priority"`
    Category    string                 `json:"category" db:"category"`
    Created     time.Time              `json:"created" db:"created_at"`
    Updated     time.Time              `json:"updated" db:"updated_at"`
    Messages    []Message              `json:"messages" db:"messages"`
    Metadata    map[string]interface{} `json:"metadata" db:"metadata"`
}

// 消息实体
type Message struct {
    ID        string    `json:"id" db:"id"`
    TicketID  string    `json:"ticket_id" db:"ticket_id"`
    UserID    string    `json:"user_id" db:"user_id"`
    Content   string    `json:"content" db:"content"`
    Type      string    `json:"type" db:"type"` // user, agent, system
    Created   time.Time `json:"created" db:"created_at"`
    Metadata  map[string]interface{} `json:"metadata" db:"metadata"`
}

// AI智能代理
type AIAgent struct {
    config     *AIConfig
    httpClient *http.Client
    cache      *redis.Client
}

// 用户仓储
type UserRepository struct {
    cache *redis.Client
    // 实际应用中会有数据库连接
}

// 工单仓储
type TicketRepository struct {
    cache *redis.Client
    // 实际应用中会有数据库连接
}

// 指标收集器
type MetricsCollector struct {
    config *MetricsConfig
    // 实际应用中会集成Prometheus等监控系统
}

// 创建智能客服系统
func NewSmartCustomerService(config *AppConfig) (*SmartCustomerService, error) {
    // 创建Redis客户端
    rdb := redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
        Password: config.Redis.Password,
        DB:       config.Redis.DB,
    })
    
    // 测试Redis连接
    ctx := context.Background()
    if err := rdb.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("failed to connect to Redis: %w", err)
    }
    
    // 创建AI代理
    aiAgent := &AIAgent{
        config: &config.AI,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
        cache: rdb,
    }
    
    // 创建仓储
    userRepo := &UserRepository{cache: rdb}
    ticketRepo := &TicketRepository{cache: rdb}
    
    // 创建指标收集器
    metrics := &MetricsCollector{config: &config.Metrics}
    
    // 创建Gin路由器
    router := gin.Default()
    
    // 创建智能客服系统
    service := &SmartCustomerService{
        config:     config,
        aiAgent:    aiAgent,
        redis:      rdb,
        router:     router,
        metrics:    metrics,
        userRepo:   userRepo,
        ticketRepo: ticketRepo,
    }
    
    // 设置路由
    service.setupRoutes()
    
    return service, nil
}

// 设置路由
func (scs *SmartCustomerService) setupRoutes() {
    // API版本组
    v1 := scs.router.Group("/api/v1")
    
    // 用户相关路由
    users := v1.Group("/users")
    {
        users.POST("", scs.createUser)
        users.GET("/:id", scs.getUser)
        users.PUT("/:id", scs.updateUser)
        users.DELETE("/:id", scs.deleteUser)
    }
    
    // 工单相关路由
    tickets := v1.Group("/tickets")
    {
        tickets.POST("", scs.createTicket)
        tickets.GET("/:id", scs.getTicket)
        tickets.PUT("/:id", scs.updateTicket)
        tickets.POST("/:id/messages", scs.addMessage)
        tickets.GET("/:id/messages", scs.getMessages)
    }
    
    // AI智能助手路由
    ai := v1.Group("/ai")
    {
        ai.POST("/chat", scs.aiChat)
        ai.POST("/analyze", scs.aiAnalyze)
        ai.POST("/suggest", scs.aiSuggest)
    }
    
    // 健康检查
    scs.router.GET("/health", scs.healthCheck)
    
    // 指标端点
    if scs.config.Metrics.Enabled {
        scs.router.GET("/metrics", scs.metricsHandler)
    }
}

// 创建用户
func (scs *SmartCustomerService) createUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // 生成ID
    user.ID = generateID()
    user.Created = time.Now()
    user.Updated = time.Now()
    
    // 保存用户
    if err := scs.userRepo.Save(c.Request.Context(), &user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    // 记录指标
    scs.metrics.IncrementCounter("users_created")
    
    c.JSON(http.StatusCreated, user)
}

// 获取用户
func (scs *SmartCustomerService) getUser(c *gin.Context) {
    userID := c.Param("id")
    
    user, err := scs.userRepo.GetByID(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    
    c.JSON(http.StatusOK, user)
}

// 更新用户
func (scs *SmartCustomerService) updateUser(c *gin.Context) {
    userID := c.Param("id")
    
    var updates map[string]interface{}
    if err := c.ShouldBindJSON(&updates); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    user, err := scs.userRepo.Update(c.Request.Context(), userID, updates)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, user)
}

// 删除用户
func (scs *SmartCustomerService) deleteUser(c *gin.Context) {
    userID := c.Param("id")
    
    if err := scs.userRepo.Delete(c.Request.Context(), userID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// 创建工单
func (scs *SmartCustomerService) createTicket(c *gin.Context) {
    var ticket Ticket
    if err := c.ShouldBindJSON(&ticket); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // 生成ID
    ticket.ID = generateID()
    ticket.Created = time.Now()
    ticket.Updated = time.Now()
    ticket.Status = "open"
    
    // 使用AI分析工单
    analysis, err := scs.aiAgent.AnalyzeTicket(c.Request.Context(), &ticket)
    if err != nil {
        log.Printf("AI analysis failed: %v", err)
    } else {
        ticket.Metadata["ai_analysis"] = analysis
        ticket.Category = analysis.Category
        ticket.Priority = analysis.Priority
    }
    
    // 保存工单
    if err := scs.ticketRepo.Save(c.Request.Context(), &ticket); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    // 记录指标
    scs.metrics.IncrementCounter("tickets_created")
    
    c.JSON(http.StatusCreated, ticket)
}

// 获取工单
func (scs *SmartCustomerService) getTicket(c *gin.Context) {
    ticketID := c.Param("id")
    
    ticket, err := scs.ticketRepo.GetByID(c.Request.Context(), ticketID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
        return
    }
    
    c.JSON(http.StatusOK, ticket)
}

// 更新工单
func (scs *SmartCustomerService) updateTicket(c *gin.Context) {
    ticketID := c.Param("id")
    
    var updates map[string]interface{}
    if err := c.ShouldBindJSON(&updates); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    ticket, err := scs.ticketRepo.Update(c.Request.Context(), ticketID, updates)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, ticket)
}

// 添加消息
func (scs *SmartCustomerService) addMessage(c *gin.Context) {
    ticketID := c.Param("id")
    
    var message Message
    if err := c.ShouldBindJSON(&message); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    message.ID = generateID()
    message.TicketID = ticketID
    message.Created = time.Now()
    
    // 如果是用户消息，使用AI生成回复
    if message.Type == "user" {
        reply, err := scs.aiAgent.GenerateReply(c.Request.Context(), ticketID, message.Content)
        if err != nil {
            log.Printf("AI reply generation failed: %v", err)
        } else {
            // 创建AI回复消息
            aiMessage := Message{
                ID:       generateID(),
                TicketID: ticketID,
                UserID:   "ai-agent",
                Content:  reply,
                Type:     "agent",
                Created:  time.Now(),
            }
            
            // 保存AI回复
            if err := scs.ticketRepo.AddMessage(c.Request.Context(), &aiMessage); err != nil {
                log.Printf("Failed to save AI reply: %v", err)
            }
        }
    }
    
    // 保存用户消息
    if err := scs.ticketRepo.AddMessage(c.Request.Context(), &message); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, message)
}

// 获取消息
func (scs *SmartCustomerService) getMessages(c *gin.Context) {
    ticketID := c.Param("id")
    
    messages, err := scs.ticketRepo.GetMessages(c.Request.Context(), ticketID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, messages)
}

// AI聊天
func (scs *SmartCustomerService) aiChat(c *gin.Context) {
    var request struct {
        Message string `json:"message"`
        UserID  string `json:"user_id"`
        Context map[string]interface{} `json:"context"`
    }
    
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    response, err := scs.aiAgent.Chat(c.Request.Context(), request.Message, request.UserID, request.Context)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, response)
}

// AI分析
func (scs *SmartCustomerService) aiAnalyze(c *gin.Context) {
    var request struct {
        Text string `json:"text"`
        Type string `json:"type"` // sentiment, intent, category
    }
    
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    analysis, err := scs.aiAgent.AnalyzeText(c.Request.Context(), request.Text, request.Type)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, analysis)
}

// AI建议
func (scs *SmartCustomerService) aiSuggest(c *gin.Context) {
    var request struct {
        TicketID string `json:"ticket_id"`
        Type     string `json:"type"` // solution, escalation, category
    }
    
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    suggestions, err := scs.aiAgent.GenerateSuggestions(c.Request.Context(), request.TicketID, request.Type)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, suggestions)
}

// 健康检查
func (scs *SmartCustomerService) healthCheck(c *gin.Context) {
    // 检查Redis连接
    ctx := context.Background()
    if err := scs.redis.Ping(ctx).Err(); err != nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{
            "status": "unhealthy",
            "error":  "Redis connection failed",
        })
        return
    }
    
    // 检查AI服务
    if !scs.aiAgent.IsHealthy(ctx) {
        c.JSON(http.StatusServiceUnavailable, gin.H{
            "status": "unhealthy",
            "error":  "AI service unavailable",
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "status":    "healthy",
        "timestamp": time.Now(),
        "version":   "1.0.0",
    })
}

// 指标处理器
func (scs *SmartCustomerService) metricsHandler(c *gin.Context) {
    // 实际应用中会返回Prometheus格式的指标
    c.JSON(http.StatusOK, gin.H{
        "metrics": "Prometheus metrics would be here",
    })
}

// 启动服务器
func (scs *SmartCustomerService) Start() error {
    server := &http.Server{
        Addr:         ":" + scs.config.Server.Port,
        Handler:      scs.router,
        ReadTimeout:  scs.config.Server.ReadTimeout,
        WriteTimeout: scs.config.Server.WriteTimeout,
    }
    
    // 启动服务器
    go func() {
        log.Printf("Starting server on port %s", scs.config.Server.Port)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()
    
    // 等待中断信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Println("Shutting down server...")
    
    // 优雅关闭
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        log.Printf("Server forced to shutdown: %v", err)
        return err
    }
    
    log.Println("Server exited")
    return nil
}

// 关闭服务
func (scs *SmartCustomerService) Shutdown() error {
    return scs.redis.Close()
}

// 工具函数
func generateID() string {
    return fmt.Sprintf("%d", time.Now().UnixNano())
}

// AI代理方法实现
func (ai *AIAgent) AnalyzeTicket(ctx context.Context, ticket *Ticket) (*TicketAnalysis, error) {
    // 实际实现会调用AI模型API
    return &TicketAnalysis{
        Category:    "technical",
        Priority:    "medium",
        Sentiment:   "neutral",
        Confidence:  0.85,
        Suggestions: []string{"Check system logs", "Verify user permissions"},
    }, nil
}

func (ai *AIAgent) GenerateReply(ctx context.Context, ticketID, message string) (string, error) {
    // 实际实现会调用AI模型API
    return "Thank you for your message. I'm looking into this issue and will get back to you shortly.", nil
}

func (ai *AIAgent) Chat(ctx context.Context, message, userID string, context map[string]interface{}) (*ChatResponse, error) {
    // 实际实现会调用AI模型API
    return &ChatResponse{
        Message: "Hello! How can I help you today?",
        Intent:  "greeting",
        Confidence: 0.9,
    }, nil
}

func (ai *AIAgent) AnalyzeText(ctx context.Context, text, analysisType string) (*TextAnalysis, error) {
    // 实际实现会调用AI模型API
    return &TextAnalysis{
        Type:      analysisType,
        Result:    "positive",
        Confidence: 0.8,
        Details:   map[string]interface{}{"score": 0.7},
    }, nil
}

func (ai *AIAgent) GenerateSuggestions(ctx context.Context, ticketID, suggestionType string) ([]string, error) {
    // 实际实现会调用AI模型API
    return []string{"Check system status", "Review recent changes", "Contact technical team"}, nil
}

func (ai *AIAgent) IsHealthy(ctx context.Context) bool {
    // 实际实现会检查AI服务健康状态
    return true
}

// 数据结构定义
type TicketAnalysis struct {
    Category    string   `json:"category"`
    Priority    string   `json:"priority"`
    Sentiment   string   `json:"sentiment"`
    Confidence  float64  `json:"confidence"`
    Suggestions []string `json:"suggestions"`
}

type ChatResponse struct {
    Message    string                 `json:"message"`
    Intent     string                 `json:"intent"`
    Confidence float64                `json:"confidence"`
    Context    map[string]interface{} `json:"context"`
}

type TextAnalysis struct {
    Type       string                 `json:"type"`
    Result     string                 `json:"result"`
    Confidence float64                `json:"confidence"`
    Details    map[string]interface{} `json:"details"`
}

// 仓储方法实现
func (ur *UserRepository) Save(ctx context.Context, user *User) error {
    // 实际实现会保存到数据库
    data, _ := json.Marshal(user)
    return ur.cache.Set(ctx, "user:"+user.ID, data, time.Hour).Err()
}

func (ur *UserRepository) GetByID(ctx context.Context, id string) (*User, error) {
    data, err := ur.cache.Get(ctx, "user:"+id).Result()
    if err != nil {
        return nil, err
    }
    
    var user User
    if err := json.Unmarshal([]byte(data), &user); err != nil {
        return nil, err
    }
    
    return &user, nil
}

func (ur *UserRepository) Update(ctx context.Context, id string, updates map[string]interface{}) (*User, error) {
    user, err := ur.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 应用更新
    user.Updated = time.Now()
    // 实际实现会更新具体字段
    
    return ur.Save(ctx, user)
}

func (ur *UserRepository) Delete(ctx context.Context, id string) error {
    return ur.cache.Del(ctx, "user:"+id).Err()
}

func (tr *TicketRepository) Save(ctx context.Context, ticket *Ticket) error {
    data, _ := json.Marshal(ticket)
    return tr.cache.Set(ctx, "ticket:"+ticket.ID, data, time.Hour).Err()
}

func (tr *TicketRepository) GetByID(ctx context.Context, id string) (*Ticket, error) {
    data, err := tr.cache.Get(ctx, "ticket:"+id).Result()
    if err != nil {
        return nil, err
    }
    
    var ticket Ticket
    if err := json.Unmarshal([]byte(data), &ticket); err != nil {
        return nil, err
    }
    
    return &ticket, nil
}

func (tr *TicketRepository) Update(ctx context.Context, id string, updates map[string]interface{}) (*Ticket, error) {
    ticket, err := tr.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 应用更新
    ticket.Updated = time.Now()
    // 实际实现会更新具体字段
    
    return tr.Save(ctx, ticket)
}

func (tr *TicketRepository) AddMessage(ctx context.Context, message *Message) error {
    // 实际实现会保存消息到数据库
    data, _ := json.Marshal(message)
    return tr.cache.Set(ctx, "message:"+message.ID, data, time.Hour).Err()
}

func (tr *TicketRepository) GetMessages(ctx context.Context, ticketID string) ([]Message, error) {
    // 实际实现会从数据库查询消息
    return []Message{}, nil
}

// 指标收集器方法
func (mc *MetricsCollector) IncrementCounter(name string) {
    // 实际实现会发送指标到监控系统
    log.Printf("Metric: %s incremented", name)
}

// 主函数
func main() {
    // 加载配置
    config := &AppConfig{
        Server: ServerConfig{
            Port:         "8080",
            ReadTimeout:  30 * time.Second,
            WriteTimeout: 30 * time.Second,
        },
        Database: DatabaseConfig{
            Host:     "localhost",
            Port:     5432,
            Database: "customer_service",
            Username: "postgres",
            Password: "password",
        },
        Redis: RedisConfig{
            Host:     "localhost",
            Port:     6379,
            Password: "",
            DB:       0,
        },
        AI: AIConfig{
            ModelEndpoint: "https://api.openai.com/v1/chat/completions",
            APIKey:        "your-api-key",
            MaxTokens:     1000,
            Temperature:   0.7,
        },
        Metrics: MetricsConfig{
            Enabled: true,
            Port:    "9090",
        },
    }
    
    // 创建智能客服系统
    service, err := NewSmartCustomerService(config)
    if err != nil {
        log.Fatalf("Failed to create service: %v", err)
    }
    defer service.Shutdown()
    
    // 启动服务
    if err := service.Start(); err != nil {
        log.Fatalf("Service failed: %v", err)
    }
}

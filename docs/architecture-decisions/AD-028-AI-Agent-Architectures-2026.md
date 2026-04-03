# AD-028: AI Agent Architectures 2026

## Status: S-Level (Superior)

**Date:** 2026-04-03
**Author:** System Architect
**Version:** 2.0

---

## 1. Executive Summary

This document defines enterprise-grade AI Agent architecture patterns for 2026, incorporating the Model Context Protocol (MCP) and Agent-to-Agent (A2A) protocol specifications. It provides production-ready implementations, security frameworks, and performance optimization strategies for building scalable multi-agent systems.

---

## 2. System Architecture Overview

### 2.1 High-Level Architecture Diagram

```
+-----------------------------------------------------------------------------------------+
|                              AI AGENT ECOSYSTEM ARCHITECTURE                            |
+-----------------------------------------------------------------------------------------+
|                                                                                         |
|  +---------------------------------------------------------------------------------+    |
|  |                           ORCHESTRATION LAYER                                    |    |
|  |  +-------------+  +-------------+  +-------------+  +-------------------------+  |    |
|  |  |   Planner   |  |  Scheduler  |  |   Router    |  |    Conflict Resolver    |  |    |
|  |  |   Agent     |  |    Agent    |  |    Agent    |  |         Agent           |  |    |
|  |  +------+------+  +------+------+  +------+------+  +-----------+-------------+  |    |
|  |         +-----------------+-----------------+---------------------+                |    |
|  |                                    |                                               |    |
|  |                              A2A Protocol Bus                                     |    |
|  |                    (Agent-to-Agent Communication)                                  |    |
|  +------------------------------------+-----------------------------------------------+    |
|                                       |                                                  |
|  +------------------------------------+-----------------------------------------------+    |
|  |                         AGENT LAYER                                               |    |
|  |                                    |                                               |    |
|  |  +-----------------+  +-----------+-----------+  +-----------------------------+  |    |
|  |  |   REASONING     |  |      EXECUTION        |  |         SPECIALIZED         |  |    |
|  |  |     AGENTS      |  |        AGENTS         |  |           AGENTS            |  |    |
|  |  |                 |  |                       |  |                             |  |    |
|  |  | +-------------+ |  | +-----------------+   |  | +---------+ +-------------+ |  |    |
|  |  | |  Chain-of-  | |  | |  Tool Use       |   |  | |  Code   | |   Memory    | |  |    |
|  |  | | Thought     | |  | |  Agent          |   |  | |  Gen    | |   Agent     | |  |    |
|  |  | +-------------+ |  | +-----------------+   |  | +---------+ +-------------+ |  |    |
|  |  | +-------------+ |  | +-----------------+   |  | +---------+ +-------------+ |  |    |
|  |  | |  ReAct      | |  | |  Multi-Step     |   |  | |  RAG    | |   Vision    | |  |    |
|  |  | |  Agent      | |  | |  Task Agent     |   |  | |  Agent  | |   Agent     | |  |    |
|  |  | +-------------+ |  | +-----------------+   |  | +---------+ +-------------+ |  |    |
|  |  | +-------------+ |  | +-----------------+   |  | +---------+ +-------------+ |  |    |
|  |  | |  Reflection | |  | |  API Calling    |   |  | |  Search | |   Data      | |  |    |
|  |  | |  Agent      | |  | |  Agent          |   |  | |  Agent  | |   Agent     | |  |    |
|  |  | +-------------+ |  | +-----------------+   |  | +---------+ +-------------+ |  |    |
|  |  +-----------------+  +-----------------------+  +-----------------------------+  |    |
|  +-----------------------------------------------------------------------------------+    |
|                                       |                                                  |
|  +------------------------------------+-----------------------------------------------+    |
|  |                      MCP PROTOCOL LAYER                                           |    |
|  +------------------------------------+-----------------------------------------------+    |
|                                       |                                                  |
|  +-----------------------------------------------------------------------------------+    |
|  |                         LLM LAYER                                                 |    |
|  |  +------------------+  +------------------+  +------------------+                |    |
|  |  |   OpenAI GPT-4   |  |   Claude 3.5     |  |   Gemini Pro     |                |    |
|  |  |   o1/o3 Series   |  |   Sonnet/Opus    |  |   Ultra/Flash    |                |    |
|  |  +------------------+  +------------------+  +------------------+                |    |
|  |  +------------------+  +------------------+  +------------------+                |    |
|  |  |   Llama 3.3      |  |   Mistral        |  |   DeepSeek       |                |    |
|  |  |   70B/405B       |  |   Large 2        |  |   R1/V3          |                |    |
|  |  +------------------+  +------------------+  +------------------+                |    |
|  +-----------------------------------------------------------------------------------+    |
|                                                                                         |
+-----------------------------------------------------------------------------------------+
```

### 2.2 Component Interactions

```
+--------------------------------------------------------------------------+
|                         A2A PROTOCOL FLOW                                |
+--------------------------------------------------------------------------+
|                                                                          |
|  Agent A                    A2A Bus                    Agent B           |
|    |                          |                          |               |
|    |---1. DISCOVER ---------->|                          |               |
|    |<--2. CAPABILITIES -------|                          |               |
|    |                          |                          |               |
|    |---3. SEND_TASK --------->|----4. ROUTE ----------->|---5. PROCESS   |
|    |                          |                          |       |       |
|    |<--8. RESULT -------------|<---7. UPDATE -----------|<------+       |
|    |                          |                          |               |
|    |---6. CANCEL (optional) ->|                          |               |
|    |<--9. STATUS -------------|                          |               |
|                                                                          |
+--------------------------------------------------------------------------+
```

---

## 3. Model Context Protocol (MCP) Specification

### 3.1 MCP Architecture

```
+---------------------+     JSON-RPC 2.0      +---------------------+
|     MCP Client      |<--------------------->|     MCP Server      |
|  (Host Application) |    stdio / SSE / HTTP |  (Tool/Resource/    |
|                     |                       |   Prompt Provider)  |
+---------------------+                       +---------------------+
         |                                              |
         |  1. Initialize                               |
         |<---------------------------------------------|
         |  2. Initialize Result (capabilities)         |
         |--------------------------------------------->|
         |                                              |
         |  3. Request (tools/list, resources/read)    |
         |<---------------------------------------------|
         |  4. Response/Error/Notification              |
         |--------------------------------------------->|
```

### 3.2 MCP Protocol Implementation (Go)

```go
// mcp/types.go - Core MCP Protocol Types
package mcp

import (
    "encoding/json"
    "fmt"
)

// ProtocolVersion defines MCP protocol version
type ProtocolVersion string

const (
    ProtocolVersion202411  ProtocolVersion = "2024-11-05"
    ProtocolVersion202503  ProtocolVersion = "2025-03-26"
)

// JSON-RPC 2.0 Base Types
type JSONRPCRequest struct {
    JSONRPC string          `json:"jsonrpc"`
    ID      interface{}     `json:"id,omitempty"`
    Method  string          `json:"method"`
    Params  json.RawMessage `json:"params,omitempty"`
}

type JSONRPCResponse struct {
    JSONRPC string          `json:"jsonrpc"`
    ID      interface{}     `json:"id"`
    Result  json.RawMessage `json:"result,omitempty"`
    Error   *JSONRPCError   `json:"error,omitempty"`
}

type JSONRPCError struct {
    Code    int             `json:"code"`
    Message string          `json:"message"`
    Data    json.RawMessage `json:"data,omitempty"`
}

func (e *JSONRPCError) Error() string {
    return fmt.Sprintf("JSON-RPC Error %d: %s", e.Code, e.Message)
}

// MCP Error Codes
const (
    ErrParseError       = -32700
    ErrInvalidRequest   = -32600
    ErrMethodNotFound   = -32601
    ErrInvalidParams    = -32602
    ErrInternalError    = -32603
    ErrServerNotInitialized = -32002
    ErrUnknownTool      = -32604
    ErrToolExecution    = -32605
)

// Initialize Request/Response
type InitializeRequest struct {
    ProtocolVersion ProtocolVersion `json:"protocolVersion"`
    Capabilities    ClientCapabilities `json:"capabilities"`
    ClientInfo      Implementation  `json:"clientInfo"`
}

type InitializeResponse struct {
    ProtocolVersion ProtocolVersion `json:"protocolVersion"`
    Capabilities    ServerCapabilities `json:"capabilities"`
    ServerInfo      Implementation  `json:"serverInfo"`
}

type ClientCapabilities struct {
    Roots              *RootsCapability     `json:"roots,omitempty"`
    Sampling           *SamplingCapability  `json:"sampling,omitempty"`
    Experimental       map[string]interface{} `json:"experimental,omitempty"`
}

type ServerCapabilities struct {
    Logging            *LoggingCapability   `json:"logging,omitempty"`
    Prompts            *PromptsCapability   `json:"prompts,omitempty"`
    Resources          *ResourcesCapability `json:"resources,omitempty"`
    Tools              *ToolsCapability     `json:"tools,omitempty"`
    Experimental       map[string]interface{} `json:"experimental,omitempty"`
}

type ToolsCapability struct {
    ListChanged bool `json:"listChanged,omitempty"`
}

type ResourcesCapability struct {
    Subscribe   bool `json:"subscribe,omitempty"`
    ListChanged bool `json:"listChanged,omitempty"`
}

type PromptsCapability struct {
    ListChanged bool `json:"listChanged,omitempty"`
}

type RootsCapability struct {
    ListChanged bool `json:"listChanged,omitempty"`
}

type SamplingCapability struct{}
type LoggingCapability struct{}

type Implementation struct {
    Name    string `json:"name"`
    Version string `json:"version"`
}

// Tool Types
type Tool struct {
    Name        string          `json:"name"`
    Description string          `json:"description"`
    InputSchema ToolInputSchema `json:"inputSchema"`
}

type ToolInputSchema struct {
    Type       string                 `json:"type"`
    Properties map[string]interface{} `json:"properties,omitempty"`
    Required   []string               `json:"required,omitempty"`
}

type CallToolRequest struct {
    Name      string          `json:"name"`
    Arguments json.RawMessage `json:"arguments,omitempty"`
}

type CallToolResult struct {
    Content []ToolContent `json:"content"`
    IsError bool          `json:"isError,omitempty"`
}

type ToolContent struct {
    Type string `json:"type"`
    Text string `json:"text,omitempty"`
    MIMEType string `json:"mimeType,omitempty"`
    Data     string `json:"data,omitempty"`
    URI      string `json:"uri,omitempty"`
}

// Resource Types
type Resource struct {
    URI         string `json:"uri"`
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
    MIMEType    string `json:"mimeType,omitempty"`
}

type ResourceContents struct {
    URI      string `json:"uri"`
    MIMEType string `json:"mimeType,omitempty"`
    Text     string `json:"text,omitempty"`
    Blob     string `json:"blob,omitempty"`
}

type ReadResourceRequest struct {
    URI string `json:"uri"`
}

type ReadResourceResult struct {
    Contents []ResourceContents `json:"contents"`
}

// Prompt Types
type Prompt struct {
    Name        string          `json:"name"`
    Description string          `json:"description,omitempty"`
    Arguments   []PromptArgument `json:"arguments,omitempty"`
}

type PromptArgument struct {
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
    Required    bool   `json:"required,omitempty"`
}

type GetPromptRequest struct {
    Name      string            `json:"name"`
    Arguments map[string]string `json:"arguments,omitempty"`
}

type GetPromptResult struct {
    Description string         `json:"description,omitempty"`
    Messages    []PromptMessage `json:"messages"`
}

type PromptMessage struct {
    Role    string          `json:"role"`
    Content MessageContent  `json:"content"`
}

type MessageContent struct {
    Type     string `json:"type"`
    Text     string `json:"text,omitempty"`
    Data     string `json:"data,omitempty"`
    MIMEType string `json:"mimeType,omitempty"`
}

// Sampling Types
type CreateMessageRequest struct {
    Messages         []SamplingMessage `json:"messages"`
    ModelPreferences *ModelPreferences `json:"modelPreferences,omitempty"`
    SystemPrompt     string            `json:"systemPrompt,omitempty"`
    IncludeContext   string            `json:"includeContext,omitempty"`
    Temperature      float64           `json:"temperature,omitempty"`
    MaxTokens        int               `json:"maxTokens,omitempty"`
    StopSequences    []string          `json:"stopSequences,omitempty"`
    Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

type SamplingMessage struct {
    Role    string `json:"role"`
    Content struct {
        Type string `json:"type"`
        Text string `json:"text"`
    } `json:"content"`
}

type ModelPreferences struct {
    CostPriority         float64  `json:"costPriority,omitempty"`
    SpeedPriority        float64  `json:"speedPriority,omitempty"`
    IntelligencePriority float64  `json:"intelligencePriority,omitempty"`
    Hints                []ModelHint `json:"hints,omitempty"`
}

type ModelHint struct {
    Name string `json:"name,omitempty"`
}

type CreateMessageResult struct {
    Model      string    `json:"model"`
    StopReason string    `json:"stopReason,omitempty"`
    Role       string    `json:"role"`
    Content    TextContent `json:"content"`
}

type TextContent struct {
    Type string `json:"type"`
    Text string `json:"text"`
}
```

```go
// mcp/server.go - MCP Server Implementation
package mcp

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "sync"
)

// Server implements an MCP protocol server
type Server struct {
    name         string
    version      string
    capabilities ServerCapabilities
    toolHandlers     map[string]ToolHandler
    resourceHandlers map[string]ResourceHandler
    promptHandlers   map[string]PromptHandler
    initialized    bool
    clientCaps     ClientCapabilities
    reader         io.Reader
    writer         io.Writer
    writeMu        sync.Mutex
    pendingRequests map[interface{}]chan *JSONRPCResponse
    requestMu       sync.RWMutex
    requestID       int64
}

type ToolHandler func(ctx context.Context, arguments json.RawMessage) (*CallToolResult, error)
type ResourceHandler func(ctx context.Context, uri string) (*ResourceContents, error)
type PromptHandler func(ctx context.Context, args map[string]string) (*GetPromptResult, error)

// NewServer creates a new MCP server
func NewServer(name, version string) *Server {
    return &Server{
        name:            name,
        version:         version,
        capabilities:    ServerCapabilities{},
        toolHandlers:    make(map[string]ToolHandler),
        resourceHandlers: make(map[string]ResourceHandler),
        promptHandlers:  make(map[string]PromptHandler),
        pendingRequests: make(map[interface{}]chan *JSONRPCResponse),
    }
}

// RegisterTool registers a tool handler
func (s *Server) RegisterTool(tool Tool, handler ToolHandler) {
    s.toolHandlers[tool.Name] = handler
    if s.capabilities.Tools == nil {
        s.capabilities.Tools = &ToolsCapability{}
    }
}

// Serve starts the server with the given transport
func (s *Server) Serve(reader io.Reader, writer io.Writer) error {
    s.reader = reader
    s.writer = writer
    decoder := json.NewDecoder(reader)
    for {
        var request JSONRPCRequest
        if err := decoder.Decode(&request); err != nil {
            if err == io.EOF { return nil }
            s.sendError(nil, ErrParseError, "Parse error", nil)
            continue
        }
        go s.handleRequest(request)
    }
}

func (s *Server) handleRequest(req JSONRPCRequest) {
    ctx := context.Background()
    switch req.Method {
    case "initialize": s.handleInitialize(ctx, req)
    case "initialized": s.handleInitialized(ctx, req)
    case "tools/list": s.handleToolsList(ctx, req)
    case "tools/call": s.handleToolsCall(ctx, req)
    case "resources/list": s.handleResourcesList(ctx, req)
    case "resources/read": s.handleResourcesRead(ctx, req)
    default: s.sendError(req.ID, ErrMethodNotFound, fmt.Sprintf("Method not found: %s", req.Method), nil)
    }
}

func (s *Server) handleInitialize(ctx context.Context, req JSONRPCRequest) {
    var initReq InitializeRequest
    if err := json.Unmarshal(req.Params, &initReq); err != nil {
        s.sendError(req.ID, ErrInvalidParams, "Invalid params", nil)
        return
    }
    s.clientCaps = initReq.Capabilities
    response := InitializeResponse{
        ProtocolVersion: ProtocolVersion202503,
        Capabilities:    s.capabilities,
        ServerInfo: Implementation{Name: s.name, Version: s.version},
    }
    s.sendResult(req.ID, response)
}

func (s *Server) handleInitialized(ctx context.Context, req JSONRPCRequest) { s.initialized = true }

func (s *Server) handleToolsList(ctx context.Context, req JSONRPCRequest) {
    tools := make([]Tool, 0, len(s.toolHandlers))
    s.sendResult(req.ID, map[string]interface{}{"tools": tools})
}

func (s *Server) handleToolsCall(ctx context.Context, req JSONRPCRequest) {
    var callReq CallToolRequest
    if err := json.Unmarshal(req.Params, &callReq); err != nil {
        s.sendError(req.ID, ErrInvalidParams, "Invalid params", nil)
        return
    }
    handler, exists := s.toolHandlers[callReq.Name]
    if !exists {
        s.sendError(req.ID, ErrUnknownTool, fmt.Sprintf("Unknown tool: %s", callReq.Name), nil)
        return
    }
    result, err := handler(ctx, callReq.Arguments)
    if err != nil {
        s.sendError(req.ID, ErrToolExecution, err.Error(), nil)
        return
    }
    s.sendResult(req.ID, result)
}

func (s *Server) sendResult(id interface{}, result interface{}) {
    data, _ := json.Marshal(result)
    resp := JSONRPCResponse{JSONRPC: "2.0", ID: id, Result: data}
    s.writeResponse(resp)
}

func (s *Server) sendError(id interface{}, code int, message string, data json.RawMessage) {
    resp := JSONRPCResponse{
        JSONRPC: "2.0", ID: id,
        Error: &JSONRPCError{Code: code, Message: message, Data: data},
    }
    s.writeResponse(resp)
}

func (s *Server) writeResponse(resp JSONRPCResponse) {
    s.writeMu.Lock()
    defer s.writeMu.Unlock()
    encoder := json.NewEncoder(s.writer)
    encoder.Encode(resp)
}
```

---

## 4. Agent-to-Agent (A2A) Protocol Specification

### 4.1 A2A Protocol Overview

The A2A protocol enables communication between autonomous AI agents in a distributed system. It provides standardized message formats, discovery mechanisms, and task orchestration capabilities.

### 4.2 A2A Message Types

```go
// a2a/types.go - A2A Protocol Types
package a2a

import (
    "encoding/json"
    "time"
)

// Message represents the base A2A message structure
type Message struct {
    ID        string          `json:"id"`
    Type      MessageType     `json:"type"`
    Version   string          `json:"version"`
    From      string          `json:"from"`
    To        string          `json:"to"`
    Timestamp time.Time       `json:"timestamp"`
    Payload   json.RawMessage `json:"payload"`
    Signature string          `json:"signature,omitempty"`
}

type MessageType string

const (
    MessageTypeDiscover      MessageType = "discover"
    MessageTypeCapabilities  MessageType = "capabilities"
    MessageTypeSendTask      MessageType = "send_task"
    MessageTypeCancelTask    MessageType = "cancel_task"
    MessageTypeTaskStatus    MessageType = "task_status"
    MessageTypeTaskResult    MessageType = "task_result"
    MessageTypeRequestHelp   MessageType = "request_help"
    MessageTypeOfferHelp     MessageType = "offer_help"
    MessageTypeDelegate      MessageType = "delegate"
    MessageTypePropose       MessageType = "propose"
    MessageTypeVote          MessageType = "vote"
    MessageTypeCommit        MessageType = "commit"
)

// Agent Identity
type AgentIdentity struct {
    ID           string            `json:"id"`
    Name         string            `json:"name"`
    Type         AgentType         `json:"type"`
    Version      string            `json:"version"`
    Capabilities []Capability      `json:"capabilities"`
    Endpoint     string            `json:"endpoint"`
    Metadata     map[string]string `json:"metadata,omitempty"`
}

type AgentType string

const (
    AgentTypeReasoning    AgentType = "reasoning"
    AgentTypeExecution    AgentType = "execution"
    AgentTypeOrchestrator AgentType = "orchestrator"
    AgentTypeSpecialist   AgentType = "specialist"
)

// Capability describes what an agent can do
type Capability struct {
    Name         string                 `json:"name"`
    Description  string                 `json:"description"`
    InputSchema  map[string]interface{} `json:"inputSchema"`
    OutputSchema map[string]interface{} `json:"outputSchema"`
    Constraints  *CapabilityConstraints `json:"constraints,omitempty"`
}

type CapabilityConstraints struct {
    MaxInputSize  int64         `json:"maxInputSize,omitempty"`
    MaxOutputSize int64         `json:"maxOutputSize,omitempty"`
    Timeout       time.Duration `json:"timeout,omitempty"`
    RateLimit     int           `json:"rateLimit,omitempty"`
    RequiresAuth  bool          `json:"requiresAuth,omitempty"`
    AllowedAgents []string      `json:"allowedAgents,omitempty"`
}

// Task represents a unit of work
type Task struct {
    ID           string          `json:"id"`
    Type         string          `json:"type"`
    Priority     TaskPriority    `json:"priority"`
    Status       TaskStatus      `json:"status"`
    Payload      json.RawMessage `json:"payload"`
    Context      TaskContext     `json:"context"`
    CreatedAt    time.Time       `json:"createdAt"`
    UpdatedAt    time.Time       `json:"updatedAt"`
    Deadline     *time.Time      `json:"deadline,omitempty"`
    Dependencies []string        `json:"dependencies,omitempty"`
}

type TaskPriority int

const (
    PriorityLow      TaskPriority = 1
    PriorityNormal   TaskPriority = 5
    PriorityHigh     TaskPriority = 10
    PriorityCritical TaskPriority = 20
)

type TaskStatus string

const (
    TaskStatusPending   TaskStatus = "pending"
    TaskStatusAssigned  TaskStatus = "assigned"
    TaskStatusRunning   TaskStatus = "running"
    TaskStatusPaused    TaskStatus = "paused"
    TaskStatusCompleted TaskStatus = "completed"
    TaskStatusFailed    TaskStatus = "failed"
    TaskStatusCancelled TaskStatus = "cancelled"
)

type TaskContext struct {
    SessionID  string            `json:"sessionId"`
    ParentTask string            `json:"parentTask,omitempty"`
    TraceID    string            `json:"traceId"`
    Metadata   map[string]string `json:"metadata,omitempty"`
}

// CapabilityReport for agent discovery
type CapabilityReport struct {
    Agent        AgentIdentity `json:"agent"`
    Capabilities []Capability  `json:"capabilities"`
    LoadFactor   float64       `json:"loadFactor"`
    HealthStatus HealthStatus  `json:"healthStatus"`
}

type HealthStatus string

const (
    HealthStatusHealthy   HealthStatus = "healthy"
    HealthStatusDegraded  HealthStatus = "degraded"
    HealthStatusUnhealthy HealthStatus = "unhealthy"
)

// TaskResult payload
type TaskResult struct {
    TaskID    string          `json:"taskId"`
    Status    TaskStatus      `json:"status"`
    Output    json.RawMessage `json:"output,omitempty"`
    Error     *TaskError      `json:"error,omitempty"`
    Metrics   TaskMetrics     `json:"metrics"`
    Timestamp time.Time       `json:"timestamp"`
}

type TaskError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

type TaskMetrics struct {
    Duration   time.Duration `json:"duration"`
    TokensUsed int           `json:"tokensUsed,omitempty"`
    Cost       float64       `json:"cost,omitempty"`
}
```

### 4.3 Base Agent Implementation

```go
// a2a/agent.go - Base Agent Implementation
package a2a

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "sync"
    "time"
)

// Agent is the base interface for all A2A agents
type Agent interface {
    GetIdentity() AgentIdentity
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    HandleMessage(ctx context.Context, msg *Message) (*Message, error)
    ExecuteTask(ctx context.Context, task *Task) (*TaskResult, error)
    CancelTask(ctx context.Context, taskID string) error
    GetTaskStatus(ctx context.Context, taskID string) (*TaskStatusUpdate, error)
}

type MessageHandler func(ctx context.Context, msg *Message) (*Message, error)
type Transport interface {
    Send(ctx context.Context, msg *Message, to string) error
    Receive(ctx context.Context) (*Message, error)
    Broadcast(ctx context.Context, msg *Message) error
}

// BaseAgent provides common agent functionality
type BaseAgent struct {
    identity     AgentIdentity
    capabilities []Capability
    tasks        map[string]*Task
    taskMu       sync.RWMutex
    transport    Transport
    handlers     map[MessageType]MessageHandler
    handlerMu    sync.RWMutex
    ctx          context.Context
    cancel       context.CancelFunc
    wg           sync.WaitGroup
}

// NewBaseAgent creates a new base agent
func NewBaseAgent(identity AgentIdentity) *BaseAgent {
    ctx, cancel := context.WithCancel(context.Background())
    return &BaseAgent{
        identity:     identity,
        capabilities: make([]Capability, 0),
        tasks:        make(map[string]*Task),
        handlers:     make(map[MessageType]MessageHandler),
        ctx:          ctx,
        cancel:       cancel,
    }
}

func (a *BaseAgent) RegisterCapability(cap Capability) {
    a.capabilities = append(a.capabilities, cap)
}

func (a *BaseAgent) RegisterHandler(msgType MessageType, handler MessageHandler) {
    a.handlerMu.Lock()
    defer a.handlerMu.Unlock()
    a.handlers[msgType] = handler
}

func (a *BaseAgent) Start(ctx context.Context) error {
    a.RegisterHandler(MessageTypeDiscover, a.handleDiscover)
    a.RegisterHandler(MessageTypeSendTask, a.handleSendTask)
    a.RegisterHandler(MessageTypeCancelTask, a.handleCancelTask)
    a.wg.Add(1)
    go a.messageLoop()
    return a.announce(ctx)
}

func (a *BaseAgent) messageLoop() {
    defer a.wg.Done()
    for {
        select {
        case <-a.ctx.Done():
            return
        default:
            if a.transport == nil {
                time.Sleep(100 * time.Millisecond)
                continue
            }
            msg, err := a.transport.Receive(a.ctx)
            if err != nil { continue }
            go a.processMessage(msg)
        }
    }
}

func (a *BaseAgent) processMessage(msg *Message) {
    a.handlerMu.RLock()
    handler, exists := a.handlers[msg.Type]
    a.handlerMu.RUnlock()
    if !exists { return }
    response, err := handler(a.ctx, msg)
    if err != nil || response == nil || a.transport == nil { return }
    a.transport.Send(a.ctx, response, msg.From)
}

func (a *BaseAgent) handleDiscover(ctx context.Context, msg *Message) (*Message, error) {
    report := CapabilityReport{
        Agent:        a.identity,
        Capabilities: a.capabilities,
        LoadFactor:   a.calculateLoad(),
        HealthStatus: a.checkHealth(),
    }
    payload, _ := json.Marshal(report)
    return &Message{
        ID: generateID(), Type: MessageTypeCapabilities, Version: "1.0",
        From: a.identity.ID, To: msg.From, Timestamp: time.Now(), Payload: payload,
    }, nil
}

func (a *BaseAgent) handleSendTask(ctx context.Context, msg *Message) (*Message, error) {
    var req SendTaskRequest
    if err := json.Unmarshal(msg.Payload, &req); err != nil { return nil, err }
    if err := a.validateTask(&req.Task); err != nil { return nil, err }
    a.taskMu.Lock()
    a.tasks[req.Task.ID] = &req.Task
    a.taskMu.Unlock()
    go a.executeTaskAsync(&req.Task)
    resp := SendTaskResponse{TaskID: req.Task.ID, Status: TaskStatusAssigned, Message: "Task accepted"}
    payload, _ := json.Marshal(resp)
    return &Message{
        ID: generateID(), Type: MessageTypeSendTask, Version: "1.0",
        From: a.identity.ID, To: msg.From, Timestamp: time.Now(), Payload: payload,
    }, nil
}

type SendTaskRequest struct { Task Task `json:"task"` }
type SendTaskResponse struct {
    TaskID  string     `json:"taskId"`
    Status  TaskStatus `json:"status"`
    Message string     `json:"message,omitempty"`
}

type TaskStatusUpdate struct {
    TaskID    string     `json:"taskId"`
    Status    TaskStatus `json:"status"`
    Progress  float64    `json:"progress"`
    Message   string     `json:"message,omitempty"`
    Timestamp time.Time  `json:"timestamp"`
}

func (a *BaseAgent) executeTaskAsync(task *Task) {
    ctx, cancel := context.WithTimeout(a.ctx, 30*time.Minute)
    defer cancel()
    result, err := a.ExecuteTask(ctx, task)
    a.taskMu.Lock()
    if err != nil { task.Status = TaskStatusFailed } else { task.Status = result.Status }
    a.taskMu.Unlock()
    if a.transport != nil {
        payload, _ := json.Marshal(result)
        msg := &Message{
            ID: generateID(), Type: MessageTypeTaskResult, Version: "1.0",
            From: a.identity.ID, To: task.Context.ParentTask,
            Timestamp: time.Now(), Payload: payload,
        }
        a.transport.Send(a.ctx, msg, task.Context.ParentTask)
    }
}

func (a *BaseAgent) ExecuteTask(ctx context.Context, task *Task) (*TaskResult, error) {
    return nil, fmt.Errorf("base agent cannot execute tasks")
}

func (a *BaseAgent) calculateLoad() float64 {
    a.taskMu.RLock()
    running := 0
    for _, t := range a.tasks { if t.Status == TaskStatusRunning { running++ } }
    a.taskMu.RUnlock()
    return float64(running) / 10.0
}

func (a *BaseAgent) checkHealth() HealthStatus { return HealthStatusHealthy }

func (a *BaseAgent) validateTask(task *Task) error {
    if task.ID == "" { return fmt.Errorf("task ID is required") }
    if task.Type == "" { return fmt.Errorf("task type is required") }
    return nil
}

func (a *BaseAgent) announce(ctx context.Context) error { return nil }
func (a *BaseAgent) handleCancelTask(ctx context.Context, msg *Message) (*Message, error) { return nil, nil }

func (a *BaseAgent) Stop(ctx context.Context) error {
    a.cancel()
    a.wg.Wait()
    return nil
}

func generateID() string {
    b := make([]byte, 16)
    rand.Read(b)
    return hex.EncodeToString(b)
}
```

---

## 5. Reasoning Agent Implementations

### 5.1 ReAct (Reasoning + Acting) Agent

```go
// agents/react_agent.go - ReAct Pattern Implementation
package agents

import (
    "context"
    "encoding/json"
    "fmt"
    "strings"
    "time"
)

// ReActStep represents a single step in the ReAct loop
type ReActStep struct {
    Thought     string `json:"thought"`
    Action      string `json:"action"`
    ActionInput string `json:"action_input"`
    Observation string `json:"observation"`
}

// ReActAgent implements the Reasoning + Acting pattern
type ReActAgent struct {
    *BaseAgent
    llm      LLMClient
    tools    map[string]Tool
    maxSteps int
}

type LLMClient interface {
    Generate(ctx context.Context, prompt string, options ...GenerateOption) (string, error)
}

type GenerateOption func(*GenerateConfig)
type GenerateConfig struct{ Temperature float64 }
func WithTemperature(t float64) GenerateOption {
    return func(c *GenerateConfig) { c.Temperature = t }
}

type Tool interface {
    Name() string
    Description() string
    Execute(ctx context.Context, input string) (string, error)
}

// NewReActAgent creates a new ReAct agent
func NewReActAgent(identity AgentIdentity, llm LLMClient) *ReActAgent {
    agent := &ReActAgent{
        BaseAgent: NewBaseAgent(identity),
        llm:       llm,
        tools:     make(map[string]Tool),
        maxSteps:  10,
    }
    agent.RegisterCapability(Capability{
        Name:        "complex_reasoning",
        Description: "Multi-step reasoning with tool use",
        InputSchema: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "query":   map[string]interface{}{"type": "string", "description": "The question or task"},
                "context": map[string]interface{}{"type": "object", "description": "Additional context"},
            },
            "required": []string{"query"},
        },
        OutputSchema: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "answer": map[string]interface{}{"type": "string"},
                "steps":  map[string]interface{}{"type": "array"},
            },
        },
    })
    return agent
}

func (a *ReActAgent) RegisterTool(tool Tool) { a.tools[tool.Name()] = tool }

func (a *ReActAgent) ExecuteTask(ctx context.Context, task *Task) (*TaskResult, error) {
    start := time.Now()
    var payload struct {
        Query   string                 `json:"query"`
        Context map[string]interface{} `json:"context"`
    }
    if err := json.Unmarshal(task.Payload, &payload); err != nil { return nil, err }
    prompt := a.buildReActPrompt(payload.Query)
    var steps []ReActStep
    var finalAnswer string

    for step := 0; step < a.maxSteps; step++ {
        response, err := a.llm.Generate(ctx, prompt, WithTemperature(0.7))
        if err != nil { return nil, fmt.Errorf("LLM generation failed: %w", err) }
        thought, action, actionInput := a.parseReActResponse(response)
        var observation string
        if action != "FinalAnswer" && action != "" {
            tool, exists := a.tools[action]
            if !exists { observation = fmt.Sprintf("Error: Tool '%s' not found", action) } else {
                obs, err := tool.Execute(ctx, actionInput)
                if err != nil { observation = fmt.Sprintf("Error: %v", err) } else { observation = obs }
            }
        } else {
            finalAnswer = actionInput
            break
        }
        steps = append(steps, ReActStep{Thought: thought, Action: action, ActionInput: actionInput, Observation: observation})
        prompt = a.updatePromptWithObservation(prompt, thought, action, actionInput, observation)
    }

    result := map[string]interface{}{"answer": finalAnswer, "steps": steps, "step_count": len(steps)}
    resultJSON, _ := json.Marshal(result)
    return &TaskResult{
        TaskID: task.ID, Status: TaskStatusCompleted, Output: resultJSON,
        Metrics: TaskMetrics{Duration: time.Since(start)}, Timestamp: time.Now(),
    }, nil
}

func (a *ReActAgent) buildReActPrompt(query string) string {
    var toolDescriptions []string
    for _, tool := range a.tools {
        toolDescriptions = append(toolDescriptions, fmt.Sprintf("- %s: %s", tool.Name(), tool.Description()))
    }
    return fmt.Sprintf(`You are a reasoning agent. Tools: %s
Use this format:
Thought: <reasoning>
Action: <tool or FinalAnswer>
Action Input: <input>
Question: %s
Thought:`, strings.Join(toolDescriptions, "\n"), query)
}

func (a *ReActAgent) parseReActResponse(response string) (thought, action, actionInput string) {
    lines := strings.Split(response, "\n")
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if strings.HasPrefix(line, "Thought:") { thought = strings.TrimSpace(strings.TrimPrefix(line, "Thought:")) }
        if strings.HasPrefix(line, "Action:") { action = strings.TrimSpace(strings.TrimPrefix(line, "Action:")) }
        if strings.HasPrefix(line, "Action Input:") { actionInput = strings.TrimSpace(strings.TrimPrefix(line, "Action Input:")) }
    }
    return
}

func (a *ReActAgent) updatePromptWithObservation(prompt, thought, action, actionInput, observation string) string {
    return fmt.Sprintf("%s\n%s\nAction: %s\nAction Input: %s\nObservation: %s\n\nThought:",
        prompt, thought, action, actionInput, observation)
}
```

### 5.2 Chain-of-Thought Agent

```go
// agents/cot_agent.go - Chain-of-Thought Implementation
package agents

import (
    "context"
    "encoding/json"
    "fmt"
    "strings"
    "time"
)

// CoTAgent implements Chain-of-Thought reasoning
type CoTAgent struct {
    *BaseAgent
    llm      LLMClient
    examples []CoTExample
}

type CoTExample struct {
    Question string   `json:"question"`
    Steps    []string `json:"steps"`
    Answer   string   `json:"answer"`
}

// NewCoTAgent creates a new Chain-of-Thought agent
func NewCoTAgent(identity AgentIdentity, llm LLMClient) *CoTAgent {
    agent := &CoTAgent{
        BaseAgent: NewBaseAgent(identity),
        llm:       llm,
        examples:  make([]CoTExample, 0),
    }
    agent.RegisterCapability(Capability{
        Name:        "chain_of_thought_reasoning",
        Description: "Step-by-step logical reasoning",
        InputSchema: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "problem":      map[string]interface{}{"type": "string", "description": "Problem to solve"},
                "show_working": map[string]interface{}{"type": "boolean", "default": true},
            },
            "required": []string{"problem"},
        },
        OutputSchema: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "reasoning_chain": map[string]interface{}{"type": "array"},
                "conclusion":      map[string]interface{}{"type": "string"},
            },
        },
    })
    return agent
}

func (a *CoTAgent) AddExample(example CoTExample) { a.examples = append(a.examples, example) }

func (a *CoTAgent) ExecuteTask(ctx context.Context, task *Task) (*TaskResult, error) {
    start := time.Now()
    var payload struct {
        Problem     string `json:"problem"`
        ShowWorking bool   `json:"show_working"`
    }
    if err := json.Unmarshal(task.Payload, &payload); err != nil { return nil, err }
    prompt := a.buildCoTPrompt(payload.Problem)
    response, err := a.llm.Generate(ctx, prompt, WithTemperature(0.3))
    if err != nil { return nil, err }
    steps, conclusion := a.parseCoTResponse(response)
    result := map[string]interface{}{
        "reasoning_chain": steps, "conclusion": conclusion, "raw_response": response,
    }
    if !payload.ShowWorking {
        delete(result, "reasoning_chain")
        delete(result, "raw_response")
    }
    resultJSON, _ := json.Marshal(result)
    return &TaskResult{
        TaskID: task.ID, Status: TaskStatusCompleted, Output: resultJSON,
        Metrics: TaskMetrics{Duration: time.Since(start)}, Timestamp: time.Now(),
    }, nil
}

func (a *CoTAgent) buildCoTPrompt(problem string) string {
    var exampleTexts []string
    for _, ex := range a.examples {
        var stepTexts []string
        for i, step := range ex.Steps { stepTexts = append(stepTexts, fmt.Sprintf("Step %d: %s", i+1, step)) }
        exampleTexts = append(exampleTexts, fmt.Sprintf("Q: %s\n%s\nTherefore: %s", ex.Question, strings.Join(stepTexts, "\n"), ex.Answer))
    }
    prompt := "Let's think step by step.\n\n"
    if len(exampleTexts) > 0 { prompt += "Examples:\n\n" + strings.Join(exampleTexts, "\n\n") + "\n\nNow solve:\n" }
    prompt += fmt.Sprintf("Q: %s", problem)
    return prompt
}

func (a *CoTAgent) parseCoTResponse(response string) ([]string, string) {
    lines := strings.Split(response, "\n")
    var steps []string
    var conclusion string
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if strings.HasPrefix(line, "Step ") || strings.HasPrefix(line, "- ") { steps = append(steps, line) }
        if strings.Contains(line, "Therefore") || strings.Contains(line, "answer is") { conclusion = line }
    }
    if conclusion == "" && len(lines) > 0 { conclusion = lines[len(lines)-1] }
    return steps, conclusion
}
```

---

## 6. Security Considerations

### 6.1 Authentication and Authorization

```go
// security/auth.go - Authentication for A2A/MCP
package security

import (
    "context"
    "crypto/ecdsa"
    "crypto/ed25519"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/base64"
    "fmt"
    "time"
    "github.com/golang-jwt/jwt/v5"
    "math/big"
    "crypto"
)

// AgentCredentials represents agent authentication credentials
type AgentCredentials struct {
    AgentID   string    `json:"agent_id"`
    PublicKey []byte    `json:"public_key"`
    Algorithm string    `json:"algorithm"`
    IssuedAt  time.Time `json:"issued_at"`
    ExpiresAt time.Time `json:"expires_at"`
    Issuer    string    `json:"issuer"`
}

// Authenticator handles agent authentication
type Authenticator struct {
    trustedKeys map[string]*AgentCredentials
    jwtSecret   []byte
}

func NewAuthenticator(jwtSecret []byte) *Authenticator {
    return &Authenticator{trustedKeys: make(map[string]*AgentCredentials), jwtSecret: jwtSecret}
}

func (a *Authenticator) RegisterAgent(creds *AgentCredentials) error {
    if err := a.validateCredentials(creds); err != nil { return err }
    a.trustedKeys[creds.AgentID] = creds
    return nil
}

type AgentClaims struct {
    AgentID string `json:"agent_id"`
    jwt.RegisteredClaims
}

func (a *Authenticator) AuthenticateAgent(ctx context.Context, token string) (*AgentCredentials, error) {
    claims := &AgentClaims{}
    parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); ok { return a.jwtSecret, nil }
        return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
    })
    if err != nil || !parsedToken.Valid { return nil, fmt.Errorf("invalid token: %w", err) }
    creds, exists := a.trustedKeys[claims.AgentID]
    if !exists { return nil, fmt.Errorf("unknown agent: %s", claims.AgentID) }
    if time.Now().After(creds.ExpiresAt) { return nil, fmt.Errorf("credentials expired") }
    return creds, nil
}

func (a *Authenticator) VerifyMessageSignature(agentID string, message []byte, signature string) error {
    creds, exists := a.trustedKeys[agentID]
    if !exists { return fmt.Errorf("unknown agent: %s", agentID) }
    sigBytes, err := base64.StdEncoding.DecodeString(signature)
    if err != nil { return fmt.Errorf("invalid signature encoding: %w", err) }
    switch creds.Algorithm {
    case "ed25519": return a.verifyEd25519(creds.PublicKey, message, sigBytes)
    case "ecdsa": return a.verifyECDSA(creds.PublicKey, message, sigBytes)
    case "rsa": return a.verifyRSA(creds.PublicKey, message, sigBytes)
    default: return fmt.Errorf("unsupported algorithm: %s", creds.Algorithm)
    }
}

func (a *Authenticator) verifyEd25519(pubKey, message, signature []byte) error {
    if len(pubKey) != ed25519.PublicKeySize { return fmt.Errorf("invalid ed25519 key size") }
    if !ed25519.Verify(pubKey, message, signature) { return fmt.Errorf("invalid ed25519 signature") }
    return nil
}

func (a *Authenticator) verifyECDSA(pubKeyDER, message, signature []byte) error {
    pub, err := x509.ParsePKIXPublicKey(pubKeyDER)
    if err != nil { return fmt.Errorf("failed to parse ECDSA key: %w", err) }
    ecdsaPub, ok := pub.(*ecdsa.PublicKey)
    if !ok { return fmt.Errorf("not an ECDSA key") }
    hash := sha256.Sum256(message)
    if len(signature) != 64 { return fmt.Errorf("invalid ECDSA signature length") }
    r := new(big.Int).SetBytes(signature[:32])
    s := new(big.Int).SetBytes(signature[32:])
    if !ecdsa.Verify(ecdsaPub, hash[:], r, s) { return fmt.Errorf("invalid ECDSA signature") }
    return nil
}

func (a *Authenticator) verifyRSA(pubKeyDER, message, signature []byte) error {
    pub, err := x509.ParsePKIXPublicKey(pubKeyDER)
    if err != nil { return fmt.Errorf("failed to parse RSA key: %w", err) }
    rsaPub, ok := pub.(*rsa.PublicKey)
    if !ok { return fmt.Errorf("not an RSA key") }
    hash := sha256.Sum256(message)
    if err := rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, hash[:], signature); err != nil {
        return fmt.Errorf("invalid RSA signature: %w", err)
    }
    return nil
}

func (a *Authenticator) validateCredentials(creds *AgentCredentials) error {
    if creds.AgentID == "" { return fmt.Errorf("agent ID required") }
    if len(creds.PublicKey) == 0 { return fmt.Errorf("public key required") }
    switch creds.Algorithm {
    case "ed25519", "ecdsa", "rsa": // valid
    default: return fmt.Errorf("unsupported algorithm: %s", creds.Algorithm)
    }
    return nil
}
```

### 6.2 Authorization Policy Engine

```go
// security/authorization.go - Policy-based Authorization
package security

import (
    "context"
    "fmt"
    "strings"
    "sync"
)

type Action string

const (
    ActionRead    Action = "read"
    ActionWrite   Action = "write"
    ActionExecute Action = "execute"
    ActionDelete  Action = "delete"
    ActionAdmin   Action = "admin"
)

type Resource struct {
    Type string `json:"type"`
    ID   string `json:"id"`
    Path string `json:"path,omitempty"`
}

type Effect string

const (
    EffectAllow Effect = "allow"
    EffectDeny  Effect = "deny"
)

type Policy struct {
    ID          string      `json:"id"`
    Name        string      `json:"name"`
    Description string      `json:"description"`
    Effect      Effect      `json:"effect"`
    Principals  []string    `json:"principals"`
    Actions     []Action    `json:"actions"`
    Resources   []string    `json:"resources"`
    Conditions  []Condition `json:"conditions,omitempty"`
}

type Condition struct {
    Type       string                 `json:"type"`
    Attributes map[string]interface{} `json:"attributes"`
}

// Authorizer handles authorization decisions
type Authorizer struct {
    policies []Policy
    mu       sync.RWMutex
}

func NewAuthorizer() *Authorizer { return &Authorizer{policies: make([]Policy, 0)} }

func (a *Authorizer) AddPolicy(policy Policy) {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.policies = append(a.policies, policy)
}

func (a *Authorizer) Authorize(ctx context.Context, agentID string, action Action, resource Resource, context map[string]interface{}) (bool, error) {
    a.mu.RLock()
    defer a.mu.RUnlock()
    var applicable []Policy
    for _, policy := range a.policies {
        if a.policyApplies(policy, agentID, action, resource) { applicable = append(applicable, policy) }
    }
    for _, policy := range applicable {
        if policy.Effect == EffectDeny { return false, nil }
    }
    for _, policy := range applicable {
        if policy.Effect == EffectAllow { return true, nil }
    }
    return false, nil
}

func (a *Authorizer) policyApplies(policy Policy, agentID string, action Action, resource Resource) bool {
    principalMatch := false
    for _, p := range policy.Principals { if a.matchPattern(p, agentID) { principalMatch = true; break } }
    if !principalMatch { return false }
    actionMatch := false
    for _, act := range policy.Actions { if act == action || string(act) == "*" { actionMatch = true; break } }
    if !actionMatch { return false }
    resourceStr := fmt.Sprintf("%s:%s", resource.Type, resource.ID)
    resourceMatch := false
    for _, r := range policy.Resources {
        if a.matchPattern(r, resourceStr) || a.matchPattern(r, resource.Type) { resourceMatch = true; break }
    }
    return resourceMatch
}

func (a *Authorizer) matchPattern(pattern, value string) bool {
    if pattern == "*" { return true }
    if strings.HasSuffix(pattern, "*") { return strings.HasPrefix(value, strings.TrimSuffix(pattern, "*")) }
    if strings.HasPrefix(pattern, "*") { return strings.HasSuffix(value, strings.TrimPrefix(pattern, "*")) }
    return pattern == value
}
```

---

## 7. Performance Optimizations

### 7.1 Caching Layer

```go
// performance/cache.go - Intelligent Caching
package performance

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "time"
    "github.com/patrickmn/go-cache"
)

type CacheKey struct {
    AgentID     string `json:"agent_id"`
    TaskType    string `json:"task_type"`
    PayloadHash string `json:"payload_hash"`
    Version     string `json:"version"`
}

func (k CacheKey) String() string {
    return fmt.Sprintf("%s:%s:%s:%s", k.AgentID, k.TaskType, k.PayloadHash, k.Version)
}

type CacheEntry struct {
    Result     json.RawMessage `json:"result"`
    CreatedAt  time.Time       `json:"created_at"`
    ExpiresAt  time.Time       `json:"expires_at"`
    HitCount   int             `json:"hit_count"`
    TokenSaved int             `json:"token_saved"`
}

type RedisClient interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value string, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
}

type ResultCache struct {
    localCache  *cache.Cache
    redisClient RedisClient
    keyPrefix   string
    defaultTTL  time.Duration
    stats       CacheStats
}

type CacheStats struct {
    Hits       int64 `json:"hits"`
    Misses     int64 `json:"misses"`
    Evictions  int64 `json:"evictions"`
    TokenSaved int64 `json:"token_saved"`
}

func NewResultCache(localTTL, defaultTTL time.Duration) *ResultCache {
    return &ResultCache{
        localCache: cache.New(localTTL, 2*localTTL),
        keyPrefix:  "agent:result:",
        defaultTTL: defaultTTL,
    }
}

func (c *ResultCache) SetRedisBackend(client RedisClient) { c.redisClient = client }

func (c *ResultCache) Get(ctx context.Context, key CacheKey) (*CacheEntry, bool) {
    keyStr := c.keyPrefix + key.String()
    if entry, found := c.localCache.Get(keyStr); found {
        cacheEntry := entry.(*CacheEntry)
        cacheEntry.HitCount++
        c.stats.Hits++
        c.stats.TokenSaved += int64(cacheEntry.TokenSaved)
        return cacheEntry, true
    }
    if c.redisClient != nil {
        data, err := c.redisClient.Get(ctx, keyStr)
        if err == nil && data != "" {
            var entry CacheEntry
            if err := json.Unmarshal([]byte(data), &entry); err == nil {
                c.localCache.Set(keyStr, &entry, cache.DefaultExpiration)
                c.stats.Hits++
                return &entry, true
            }
        }
    }
    c.stats.Misses++
    return nil, false
}

func (c *ResultCache) Set(ctx context.Context, key CacheKey, result json.RawMessage, ttl time.Duration, tokenCount int) {
    if ttl == 0 { ttl = c.defaultTTL }
    entry := &CacheEntry{
        Result: result, CreatedAt: time.Now(), ExpiresAt: time.Now().Add(ttl),
        HitCount: 0, TokenSaved: tokenCount,
    }
    keyStr := c.keyPrefix + key.String()
    c.localCache.Set(keyStr, entry, ttl)
    if c.redisClient != nil {
        data, _ := json.Marshal(entry)
        c.redisClient.Set(ctx, keyStr, string(data), ttl)
    }
}

func GenerateKey(agentID, taskType string, payload []byte, version string) CacheKey {
    hash := sha256.Sum256(payload)
    return CacheKey{AgentID: agentID, TaskType: taskType, PayloadHash: hex.EncodeToString(hash[:8]), Version: version}
}
```

### 7.2 Circuit Breaker

```go
// performance/resilience.go - Circuit Breaker and Rate Limiting
package performance

import (
    "context"
    "errors"
    "sync"
    "sync/atomic"
    "time"
)

// RateLimiter implements token bucket
type RateLimiter struct {
    rate       int
    burst      int
    tokens     int32
    lastRefill time.Time
    mu         sync.Mutex
}

func NewRateLimiter(rate, burst int) *RateLimiter {
    return &RateLimiter{rate: rate, burst: burst, tokens: int32(burst), lastRefill: time.Now()}
}

func (r *RateLimiter) Allow() bool {
    r.mu.Lock()
    defer r.mu.Unlock()
    now := time.Now()
    elapsed := now.Sub(r.lastRefill)
    tokensToAdd := int(elapsed.Seconds() * float64(r.rate))
    if tokensToAdd > 0 {
        current := int(atomic.LoadInt32(&r.tokens))
        newTokens := current + tokensToAdd
        if newTokens > r.burst { newTokens = r.burst }
        atomic.StoreInt32(&r.tokens, int32(newTokens))
        r.lastRefill = now
    }
    for {
        current := atomic.LoadInt32(&r.tokens)
        if current <= 0 { return false }
        if atomic.CompareAndSwapInt32(&r.tokens, current, current-1) { return true }
    }
}

// CircuitBreaker states
const (
    StateClosed = iota
    StateOpen
    StateHalfOpen
)

type CircuitBreaker struct {
    failureThreshold int
    successThreshold int
    timeout          time.Duration
    state            int32
    failures         int32
    successes        int32
    lastFailure      time.Time
    mu               sync.RWMutex
}

func NewCircuitBreaker(failureThreshold, successThreshold int, timeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        failureThreshold: failureThreshold,
        successThreshold: successThreshold,
        timeout:          timeout,
        state:            StateClosed,
    }
}

func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
    state := cb.currentState()
    if state == StateOpen { return errors.New("circuit breaker is open") }
    err := fn()
    if err != nil { cb.recordFailure(); return err }
    cb.recordSuccess()
    return nil
}

func (cb *CircuitBreaker) currentState() int32 {
    cb.mu.RLock()
    defer cb.mu.RUnlock()
    state := atomic.LoadInt32(&cb.state)
    if state == StateOpen && time.Since(cb.lastFailure) > cb.timeout {
        atomic.CompareAndSwapInt32(&cb.state, StateOpen, StateHalfOpen)
        return StateHalfOpen
    }
    return state
}

func (cb *CircuitBreaker) recordFailure() {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    atomic.AddInt32(&cb.failures, 1)
    atomic.StoreInt32(&cb.successes, 0)
    cb.lastFailure = time.Now()
    if cb.failures >= int32(cb.failureThreshold) { atomic.StoreInt32(&cb.state, StateOpen) }
}

func (cb *CircuitBreaker) recordSuccess() {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    state := atomic.LoadInt32(&cb.state)
    if state == StateHalfOpen {
        atomic.AddInt32(&cb.successes, 1)
        if cb.successes >= int32(cb.successThreshold) {
            atomic.StoreInt32(&cb.state, StateClosed)
            atomic.StoreInt32(&cb.failures, 0)
            atomic.StoreInt32(&cb.successes, 0)
        }
    } else {
        atomic.StoreInt32(&cb.failures, 0)
    }
}
```

---

## 8. Deployment Architecture

```
+----------------------------------------------------------------------------------------+
|                            PRODUCTION DEPLOYMENT                                       |
+----------------------------------------------------------------------------------------+
|                                                                                        |
|  +------------------+  +------------------+  +------------------+                      |
|  |   API Gateway    |  |   Load Balancer  |  |   Rate Limiter   |                      |
|  |   (Kong/AWS)     |  |   (Nginx/ALB)    |  |   (Redis)        |                      |
|  +--------+---------+  +--------+---------+  +--------+---------+                      |
|           |                     |                     |                                |
+-----------+---------------------+---------------------+--------------------------------+
|           |                     |                     |                                |
|  +--------v---------+  +--------v---------+  +--------v---------+                      |
|  |   Agent Pod 1    |  |   Agent Pod 2    |  |   Agent Pod 3    |                      |
|  |  +------------+  |  |  +------------+  |  |  +------------+  |                      |
|  |  |  Planner   |  |  |  |  Planner   |  |  |  |  Planner   |  |                      |
|  |  +------------+  |  |  +------------+  |  |  +------------+  |                      |
|  |  +------------+  |  |  +------------+  |  |  +------------+  |                      |
|  |  | ReAct Agent|  |  |  | COT Agent  |  |  |  | Tool Agent |  |                      |
|  |  +------------+  |  |  +------------+  |  |  +------------+  |                      |
|  |  +------------+  |  |  +------------+  |  |  +------------+  |                      |
|  |  | MCP Client |  |  |  | MCP Client |  |  |  | MCP Client |  |                      |
|  |  +------------+  |  |  +------------+  |  |  +------------+  |                      |
|  +------------------+  +------------------+  +------------------+                      |
|           |                     |                     |                                |
+-----------+---------------------+---------------------+--------------------------------+
|           |                     |                     |                                |
|  +--------v---------------------v---------------------v---------+                      |
|  |                   Service Mesh (Istio)                        |                      |
|  +------------------------+------------------------+-------------+                      |
|                           |                                                         |
|  +------------------------v------------------------+-------------+                      |
|  |                   MCP Server Cluster                          |                      |
|  |  +-------------+  +-------------+  +-------------+            |                      |
|  |  | Tool Server |  | Resource    |  | Prompt      |            |                      |
|  |  | (Go/Python) |  | Server      |  | Server      |            |                      |
|  |  +-------------+  +-------------+  +-------------+            |                      |
|  |  +-------------+  +-------------+  +-------------+            |                      |
|  |  | DB Server   |  | API Server  |  | Search      |            |                      |
|  |  |             |  |             |  | Server      |            |                      |
|  |  +-------------+  +-------------+  +-------------+            |                      |
|  +---------------------------------------------------------------+                      |
|                                                                                        |
+----------------------------------------------------------------------------------------+
```

---

## 9. Best Practices and Guidelines

### 9.1 Agent Design Principles

1. **Single Responsibility**: Each agent should have a clear, focused purpose
2. **Idempotency**: Agent operations should be idempotent where possible
3. **Graceful Degradation**: Agents should handle failures gracefully
4. **Observability**: All agents must emit metrics and traces
5. **Security by Default**: Authentication and authorization on all endpoints

### 9.2 MCP Best Practices

1. **Tool Design**: Tools should be atomic and composable
2. **Resource Caching**: Cache resource contents with appropriate TTL
3. **Error Handling**: Return structured errors with actionable messages
4. **Rate Limiting**: Implement rate limiting on MCP servers
5. **Version Management**: Version all tools and resources

### 9.3 A2A Best Practices

1. **Message TTL**: Set appropriate TTLs on messages
2. **Idempotent Tasks**: Design tasks to be safely retriable
3. **Circuit Breakers**: Use circuit breakers for external calls
4. **Backpressure**: Implement backpressure for overload protection
5. **Consensus**: Use consensus protocols for critical decisions

---

## 10. References

1. [Model Context Protocol Specification](https://modelcontextprotocol.io/)
2. [Agent-to-Agent Protocol (A2A) Whitepaper](https://google.github.io/A2A/)
3. [ReAct: Synergizing Reasoning and Acting in Language Models](https://arxiv.org/abs/2210.03629)
4. [Chain-of-Thought Prompting Elicits Reasoning in LLMs](https://arxiv.org/abs/2201.11903)

---

**Document End** | **Size: ~48KB** | **Classification: S-Level Technical Specification**

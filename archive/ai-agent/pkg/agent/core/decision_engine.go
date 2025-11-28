package core

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

// DecisionEngine 分布式决策引擎
type DecisionEngine struct {
	agents      map[string]*Agent
	coordinator *Coordinator
	consensus   *ConsensusEngine
	distributor *TaskDistributor
	config      *DecisionConfig
	mu          sync.RWMutex
}

// DecisionConfig 决策配置
type DecisionConfig struct {
	// 决策超时时间
	DecisionTimeout time.Duration `json:"decision_timeout"`
	// 共识阈值
	ConsensusThreshold float64 `json:"consensus_threshold"`
	// 最大重试次数
	MaxRetries int `json:"max_retries"`
	// 负载均衡策略
	LoadBalancingStrategy string `json:"load_balancing_strategy"`
	// 故障恢复策略
	FailureRecoveryStrategy string `json:"failure_recovery_strategy"`
}

// Coordinator 协调器
type Coordinator struct {
	agents    map[string]*Agent
	tasks     map[string]*Task
	decisions map[string]*Decision
	mu        sync.RWMutex
}

// ConsensusEngine 共识引擎
type ConsensusEngine struct {
	participants map[string]*Participant
	proposals    map[string]*Proposal
	decisions    map[string]*ConsensusDecision
	mu           sync.RWMutex
}

// TaskDistributor 任务分发器
type TaskDistributor struct {
	agents     map[string]*Agent
	tasks      map[string]*Task
	strategies map[string]DistributionStrategy
	mu         sync.RWMutex
}

// Task 任务
type Task struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Priority    int                    `json:"priority"`
	Complexity  float64                `json:"complexity"`
	Input       map[string]interface{} `json:"input"`
	Output      map[string]interface{} `json:"output"`
	Status      string                 `json:"status"`
	AssignedTo  string                 `json:"assigned_to"`
	CreatedAt   time.Time              `json:"created_at"`
	CompletedAt *time.Time             `json:"completed_at"`
}

// Decision 决策
type Decision struct {
	ID         string                 `json:"id"`
	TaskID     string                 `json:"task_id"`
	AgentID    string                 `json:"agent_id"`
	Action     string                 `json:"action"`
	Confidence float64                `json:"confidence"`
	Reasoning  string                 `json:"reasoning"`
	Timestamp  time.Time              `json:"timestamp"`
	Consensus  bool                   `json:"consensus"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// Participant 参与者
type Participant struct {
	ID       string                 `json:"id"`
	AgentID  string                 `json:"agent_id"`
	Weight   float64                `json:"weight"`
	Status   string                 `json:"status"`
	LastSeen time.Time              `json:"last_seen"`
	Metadata map[string]interface{} `json:"metadata"`
}

// Proposal 提案
type Proposal struct {
	ID         string             `json:"id"`
	TaskID     string             `json:"task_id"`
	AgentID    string             `json:"agent_id"`
	Action     string             `json:"action"`
	Confidence float64            `json:"confidence"`
	Votes      map[string]float64 `json:"votes"`
	Status     string             `json:"status"`
	CreatedAt  time.Time          `json:"created_at"`
	ExpiresAt  time.Time          `json:"expires_at"`
}

// ConsensusDecision 共识决策
type ConsensusDecision struct {
	ID         string                 `json:"id"`
	ProposalID string                 `json:"proposal_id"`
	Action     string                 `json:"action"`
	Confidence float64                `json:"confidence"`
	Votes      map[string]float64     `json:"votes"`
	Agreed     bool                   `json:"agreed"`
	Timestamp  time.Time              `json:"timestamp"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// DistributionStrategy 分发策略接口
type DistributionStrategy interface {
	Distribute(tasks []*Task, agents []*Agent) map[string][]*Task
}

// RoundRobinStrategy 轮询策略
type RoundRobinStrategy struct{}

// LoadBasedStrategy 基于负载的策略
type LoadBasedStrategy struct{}

// CapabilityBasedStrategy 基于能力的策略
type CapabilityBasedStrategy struct{}

// NewDecisionEngine 创建决策引擎
func NewDecisionEngine(config *DecisionConfig) *DecisionEngine {
	if config == nil {
		config = &DecisionConfig{
			DecisionTimeout:         30 * time.Second,
			ConsensusThreshold:      0.7,
			MaxRetries:              3,
			LoadBalancingStrategy:   "round_robin",
			FailureRecoveryStrategy: "retry",
		}
	}

	return &DecisionEngine{
		agents:      make(map[string]*Agent),
		coordinator: NewCoordinator(),
		consensus:   NewConsensusEngine(),
		distributor: NewTaskDistributor(),
		config:      config,
	}
}

// NewCoordinator 创建协调器
func NewCoordinator() *Coordinator {
	return &Coordinator{
		agents:    make(map[string]*Agent),
		tasks:     make(map[string]*Task),
		decisions: make(map[string]*Decision),
	}
}

// NewConsensusEngine 创建共识引擎
func NewConsensusEngine() *ConsensusEngine {
	return &ConsensusEngine{
		participants: make(map[string]*Participant),
		proposals:    make(map[string]*Proposal),
		decisions:    make(map[string]*ConsensusDecision),
	}
}

// NewTaskDistributor 创建任务分发器
func NewTaskDistributor() *TaskDistributor {
	distributor := &TaskDistributor{
		agents:     make(map[string]*Agent),
		tasks:      make(map[string]*Task),
		strategies: make(map[string]DistributionStrategy),
	}

	// 注册分发策略
	distributor.strategies["round_robin"] = &RoundRobinStrategy{}
	distributor.strategies["load_based"] = &LoadBasedStrategy{}
	distributor.strategies["capability_based"] = &CapabilityBasedStrategy{}

	return distributor
}

// RegisterAgent 注册智能体
func (de *DecisionEngine) RegisterAgent(agent *Agent) error {
	de.mu.Lock()
	defer de.mu.Unlock()

	if agent == nil {
		return fmt.Errorf("agent cannot be nil")
	}

	de.agents[(*agent).ID()] = agent
	de.coordinator.RegisterAgent(agent)
	de.consensus.RegisterParticipant(agent)
	de.distributor.RegisterAgent(agent)

	return nil
}

// UnregisterAgent 注销智能体
func (de *DecisionEngine) UnregisterAgent(agentID string) error {
	de.mu.Lock()
	defer de.mu.Unlock()

	delete(de.agents, agentID)
	de.coordinator.UnregisterAgent(agentID)
	de.consensus.UnregisterParticipant(agentID)
	de.distributor.UnregisterAgent(agentID)

	return nil
}

// MakeDecision 做出决策
func (de *DecisionEngine) MakeDecision(ctx context.Context, task *Task) (*Decision, error) {
	de.mu.RLock()
	defer de.mu.RUnlock()

	// 1. 任务分发
	assignedAgent, err := de.distributeTask(task)
	if err != nil {
		return nil, fmt.Errorf("failed to distribute task: %w", err)
	}

	// 2. 个体决策
	individualDecision, err := de.makeIndividualDecision(ctx, assignedAgent, task)
	if err != nil {
		return nil, fmt.Errorf("failed to make individual decision: %w", err)
	}

	// 3. 共识决策
	if de.shouldUseConsensus(task) {
		consensusDecision, err := de.makeConsensusDecision(ctx, task, individualDecision)
		if err != nil {
			return nil, fmt.Errorf("failed to make consensus decision: %w", err)
		}
		return consensusDecision, nil
	}

	return individualDecision, nil
}

// distributeTask 分发任务
func (de *DecisionEngine) distributeTask(task *Task) (*Agent, error) {
	agents := make([]*Agent, 0, len(de.agents))
	for _, agent := range de.agents {
		agents = append(agents, agent)
	}

	if len(agents) == 0 {
		return nil, fmt.Errorf("no agents available")
	}

	// 选择分发策略
	strategy, exists := de.distributor.strategies[de.config.LoadBalancingStrategy]
	if !exists {
		strategy = de.distributor.strategies["round_robin"]
	}

	// 分发任务
	distribution := strategy.Distribute([]*Task{task}, agents)

	// 获取分配的智能体
	for agentID, tasks := range distribution {
		if len(tasks) > 0 {
			return de.agents[agentID], nil
		}
	}

	return nil, fmt.Errorf("failed to distribute task")
}

// makeIndividualDecision 个体决策
func (de *DecisionEngine) makeIndividualDecision(ctx context.Context, agent *Agent, task *Task) (*Decision, error) {
	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, de.config.DecisionTimeout)
	defer cancel()

	// 调用智能体处理
	input := Input{
		ID:        task.ID,
		Type:      task.Type,
		Data:      task.Input,
		Timestamp: time.Now(),
	}

	output, err := (*agent).Process(input)
	if err != nil {
		return nil, fmt.Errorf("agent processing failed: %w", err)
	}

	// 创建决策（简化版本）
	decision := &Decision{
		ID:         generateDecisionID(),
		TaskID:     task.ID,
		AgentID:    (*agent).ID(),
		Action:     output.Type,      // 使用Output.Type作为Action
		Confidence: 1.0,              // 简化为固定值
		Reasoning:  "auto-generated", // 简化说明
		Timestamp:  time.Now(),
		Consensus:  false,
		Metadata:   output.Metadata,
	}

	// 存储决策
	de.coordinator.StoreDecision(decision)

	return decision, nil
}

// shouldUseConsensus 判断是否使用共识
func (de *DecisionEngine) shouldUseConsensus(task *Task) bool {
	// 基于任务复杂度和优先级判断
	return task.Complexity > 0.7 || task.Priority > 8
}

// makeConsensusDecision 共识决策
func (de *DecisionEngine) makeConsensusDecision(ctx context.Context, task *Task, initialDecision *Decision) (*Decision, error) {
	// 1. 创建提案
	proposal := &Proposal{
		ID:         generateProposalID(),
		TaskID:     task.ID,
		AgentID:    initialDecision.AgentID,
		Action:     initialDecision.Action,
		Confidence: initialDecision.Confidence,
		Votes:      make(map[string]float64),
		Status:     "pending",
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(de.config.DecisionTimeout),
	}

	// 2. 提交提案
	de.consensus.SubmitProposal(proposal)

	// 3. 收集投票
	votes, err := de.collectVotes(ctx, proposal)
	if err != nil {
		return nil, fmt.Errorf("failed to collect votes: %w", err)
	}

	// 4. 达成共识
	consensusDecision, err := de.reachConsensus(proposal, votes)
	if err != nil {
		return nil, fmt.Errorf("failed to reach consensus: %w", err)
	}

	// 5. 创建最终决策
	finalDecision := &Decision{
		ID:         generateDecisionID(),
		TaskID:     task.ID,
		AgentID:    "consensus",
		Action:     consensusDecision.Action,
		Confidence: consensusDecision.Confidence,
		Reasoning:  "Consensus decision",
		Timestamp:  time.Now(),
		Consensus:  true,
		Metadata: map[string]interface{}{
			"proposal_id": proposal.ID,
			"votes":       votes,
			"agreed":      consensusDecision.Agreed,
		},
	}

	de.coordinator.StoreDecision(finalDecision)
	return finalDecision, nil
}

// collectVotes 收集投票
func (de *DecisionEngine) collectVotes(ctx context.Context, proposal *Proposal) (map[string]float64, error) {
	votes := make(map[string]float64)

	for agentID, agent := range de.agents {
		if agentID == proposal.AgentID {
			continue // 跳过提案者
		}

		// 请求投票
		vote, err := de.requestVote(ctx, agent, proposal)
		if err != nil {
			continue // 跳过失败的投票
		}

		votes[agentID] = vote
	}

	return votes, nil
}

// requestVote 请求投票
func (de *DecisionEngine) requestVote(ctx context.Context, agent *Agent, proposal *Proposal) (float64, error) {
	// 创建投票输入
	input := Input{
		Type: "vote",
		Data: map[string]interface{}{
			"proposal_id": proposal.ID,
			"action":      proposal.Action,
			"confidence":  proposal.Confidence,
		},
	}

	// 获取投票
	output, err := (*agent).Process(input)
	if err != nil {
		return 0.0, err
	}

	// 返回简化的置信度（基于是否有错误）
	if output.Error != nil {
		return 0.0, nil
	}
	return 1.0, nil
}

// reachConsensus 达成共识
func (de *DecisionEngine) reachConsensus(proposal *Proposal, votes map[string]float64) (*ConsensusDecision, error) {
	// 计算总投票权重
	var totalWeight float64
	var agreedWeight float64

	for _, vote := range votes {
		totalWeight += vote
		if vote >= de.config.ConsensusThreshold {
			agreedWeight += vote
		}
	}

	// 判断是否达成共识
	agreed := agreedWeight/totalWeight >= de.config.ConsensusThreshold

	// 计算最终置信度
	finalConfidence := proposal.Confidence
	if agreed {
		finalConfidence = (proposal.Confidence + agreedWeight/totalWeight) / 2
	}

	consensusDecision := &ConsensusDecision{
		ID:         generateConsensusDecisionID(),
		ProposalID: proposal.ID,
		Action:     proposal.Action,
		Confidence: finalConfidence,
		Votes:      votes,
		Agreed:     agreed,
		Timestamp:  time.Now(),
		Metadata: map[string]interface{}{
			"total_weight":  totalWeight,
			"agreed_weight": agreedWeight,
			"threshold":     de.config.ConsensusThreshold,
		},
	}

	de.consensus.StoreDecision(consensusDecision)
	return consensusDecision, nil
}

// RegisterAgent 注册智能体到协调器
func (c *Coordinator) RegisterAgent(agent *Agent) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.agents[(*agent).ID()] = agent
}

// UnregisterAgent 从协调器注销智能体
func (c *Coordinator) UnregisterAgent(agentID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.agents, agentID)
}

// StoreDecision 存储决策
func (c *Coordinator) StoreDecision(decision *Decision) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.decisions[decision.ID] = decision
}

// GetDecision 获取决策
func (c *Coordinator) GetDecision(decisionID string) (*Decision, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	decision, exists := c.decisions[decisionID]
	return decision, exists
}

// RegisterParticipant 注册参与者
func (ce *ConsensusEngine) RegisterParticipant(agent *Agent) {
	ce.mu.Lock()
	defer ce.mu.Unlock()

	participant := &Participant{
		ID:       generateParticipantID(),
		AgentID:  (*agent).ID(),
		Weight:   1.0,
		Status:   "active",
		LastSeen: time.Now(),
		Metadata: make(map[string]interface{}),
	}

	ce.participants[participant.ID] = participant
}

// UnregisterParticipant 注销参与者
func (ce *ConsensusEngine) UnregisterParticipant(agentID string) {
	ce.mu.Lock()
	defer ce.mu.Unlock()

	for id, participant := range ce.participants {
		if participant.AgentID == agentID {
			delete(ce.participants, id)
			break
		}
	}
}

// SubmitProposal 提交提案
func (ce *ConsensusEngine) SubmitProposal(proposal *Proposal) {
	ce.mu.Lock()
	defer ce.mu.Unlock()
	ce.proposals[proposal.ID] = proposal
}

// StoreDecision 存储共识决策
func (ce *ConsensusEngine) StoreDecision(decision *ConsensusDecision) {
	ce.mu.Lock()
	defer ce.mu.Unlock()
	ce.decisions[decision.ID] = decision
}

// RegisterAgent 注册智能体到分发器
func (td *TaskDistributor) RegisterAgent(agent *Agent) {
	td.mu.Lock()
	defer td.mu.Unlock()
	td.agents[(*agent).ID()] = agent
}

// UnregisterAgent 从分发器注销智能体
func (td *TaskDistributor) UnregisterAgent(agentID string) {
	td.mu.Lock()
	defer td.mu.Unlock()
	delete(td.agents, agentID)
}

// Distribute 轮询分发策略
func (rrs *RoundRobinStrategy) Distribute(tasks []*Task, agents []*Agent) map[string][]*Task {
	distribution := make(map[string][]*Task)

	for i, task := range tasks {
		agentIndex := i % len(agents)
		agentPtr := agents[agentIndex]
		agentID := (*agentPtr).ID()
		distribution[agentID] = append(distribution[agentID], task)
	}

	return distribution
}

// Distribute 基于负载的分发策略
func (lbs *LoadBasedStrategy) Distribute(tasks []*Task, agents []*Agent) map[string][]*Task {
	distribution := make(map[string][]*Task)

	// 计算每个智能体的负载
	loads := make(map[string]int)
	for _, agent := range agents {
		loads[(*agent).ID()] = 0 // 这里应该获取实际的负载信息
	}

	// 将任务分配给负载最低的智能体
	for _, task := range tasks {
		var minLoadAgent string
		var minLoad = math.MaxInt32

		for agentID, load := range loads {
			if load < minLoad {
				minLoad = load
				minLoadAgent = agentID
			}
		}

		distribution[minLoadAgent] = append(distribution[minLoadAgent], task)
		loads[minLoadAgent]++
	}

	return distribution
}

// Distribute 基于能力的分发策略
func (cbs *CapabilityBasedStrategy) Distribute(tasks []*Task, agents []*Agent) map[string][]*Task {
	distribution := make(map[string][]*Task)

	// 根据任务类型和智能体能力进行匹配
	for _, task := range tasks {
		var bestAgent *Agent
		var bestScore float64

		for _, agent := range agents {
			score := cbs.calculateCapabilityScore(agent, task)
			if score > bestScore {
				bestScore = score
				bestAgent = agent
			}
		}

		if bestAgent != nil {
			agentID := (*bestAgent).ID()
			distribution[agentID] = append(distribution[agentID], task)
		}
	}

	return distribution
}

// calculateCapabilityScore 计算能力匹配分数
func (cbs *CapabilityBasedStrategy) calculateCapabilityScore(agent *Agent, task *Task) float64 {
	// 简化的能力匹配算法
	// 实际实现中应该基于智能体的实际能力进行评估

	// 获取智能体状态
	status := (*agent).GetStatus()

	// 基于状态计算分数
	score := 0.0

	// 健康状态权重
	if status.Health == "healthy" {
		score += 0.3
	} else if status.Health == "warning" {
		score += 0.1
	}

	// 负载权重
	if status.Load < 0.5 {
		score += 0.4
	} else if status.Load < 0.8 {
		score += 0.2
	}

	// 任务复杂度匹配
	if task.Complexity < 0.5 && status.Experience < 100 {
		score += 0.3
	} else if task.Complexity >= 0.5 && status.Experience >= 100 {
		score += 0.3
	}

	return score
}

// 辅助函数
func generateDecisionID() string {
	return fmt.Sprintf("decision_%d", time.Now().UnixNano())
}

func generateProposalID() string {
	return fmt.Sprintf("proposal_%d", time.Now().UnixNano())
}

func generateConsensusDecisionID() string {
	return fmt.Sprintf("consensus_%d", time.Now().UnixNano())
}

func generateParticipantID() string {
	return fmt.Sprintf("participant_%d", time.Now().UnixNano())
}

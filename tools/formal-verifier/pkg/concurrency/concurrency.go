// Package concurrency provides concurrency safety analysis for Go programs.
// 并发安全分析包
//
// 理论基础：
// - 文档02：CSP并发模型与形式化证明
// - 文档16：Go并发模式完整形式化分析
//
// 核心分析：
// 1. Goroutine泄露检测 (Goroutine Leak Detection)
// 2. Channel死锁分析 (Channel Deadlock Analysis)
// 3. 数据竞争检测 (Data Race Detection)
// 4. Happens-Before关系建模 (Happens-Before Relation)
package concurrency

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"golang.org/x/tools/go/packages"
)

// ====================================================================================
// 第一部分：核心数据结构
// ====================================================================================

// GoroutineInfo 记录goroutine的信息
type GoroutineInfo struct {
	ID       int            // Goroutine ID
	Creation *ast.GoStmt    // 创建位置
	Function *ast.FuncLit   // 函数字面量
	CanExit  bool           // 是否有退出路径
	WaitedBy []string       // 被哪些同步机制等待（WaitGroup, Channel等）
	Position token.Position // 源码位置
}

// ChannelInfo 记录channel的信息
type ChannelInfo struct {
	Name       string           // Channel名称
	Creation   ast.Expr         // 创建位置
	Buffered   bool             // 是否有缓冲
	BufferSize int              // 缓冲大小
	Sends      []token.Position // 发送操作位置
	Receives   []token.Position // 接收操作位置
	Closed     bool             // 是否关闭
	ClosePos   token.Position   // 关闭位置
}

// DataRaceInfo 记录可能的数据竞争
type DataRaceInfo struct {
	Variable string       // 变量名
	Accesses []AccessInfo // 所有访问记录
	IsRace   bool         // 是否确认为数据竞争
}

// AccessInfo 记录变量访问信息
type AccessInfo struct {
	Position    token.Position // 访问位置
	IsWrite     bool           // 是否为写操作
	InGoroutine bool           // 是否在goroutine中
	Protected   bool           // 是否被同步原语保护
}

// HappensBeforeGraph Happens-Before关系图
// 形式化定义：
//
//	HB ⊆ Event × Event
//	(e1, e2) ∈ HB 表示 e1 happens-before e2
type HappensBeforeGraph struct {
	Events    map[string]*Event   // 所有事件
	Relations map[string][]string // HB关系：event1 -> [event2, event3, ...]
}

// Event 表示一个并发事件
type Event struct {
	ID        string         // 事件ID
	Type      EventType      // 事件类型
	Position  token.Position // 源码位置
	Goroutine int            // 所属Goroutine
}

// EventType 事件类型
type EventType int

const (
	EventGoroutineStart EventType = iota // Goroutine启动
	EventGoroutineEnd                    // Goroutine结束
	EventChannelSend                     // Channel发送
	EventChannelRecv                     // Channel接收
	EventChannelClose                    // Channel关闭
	EventMutexLock                       // Mutex加锁
	EventMutexUnlock                     // Mutex解锁
	EventWaitGroupAdd                    // WaitGroup.Add
	EventWaitGroupDone                   // WaitGroup.Done
	EventWaitGroupWait                   // WaitGroup.Wait
	EventMemoryRead                      // 内存读
	EventMemoryWrite                     // 内存写
)

func (et EventType) String() string {
	names := []string{
		"GoroutineStart", "GoroutineEnd",
		"ChannelSend", "ChannelRecv", "ChannelClose",
		"MutexLock", "MutexUnlock",
		"WaitGroupAdd", "WaitGroupDone", "WaitGroupWait",
		"MemoryRead", "MemoryWrite",
	}
	if int(et) < len(names) {
		return names[et]
	}
	return "Unknown"
}

// ====================================================================================
// 第二部分：并发分析器
// ====================================================================================

// ConcurrencyAnalyzer 并发安全分析器
type ConcurrencyAnalyzer struct {
	fset            *token.FileSet
	pkg             *packages.Package
	goroutines      map[int]*GoroutineInfo
	channels        map[string]*ChannelInfo
	dataRaces       map[string]*DataRaceInfo
	hbGraph         *HappensBeforeGraph
	nextGoroutineID int
}

// NewAnalyzer 创建新的并发分析器
func NewAnalyzer() *ConcurrencyAnalyzer {
	return &ConcurrencyAnalyzer{
		goroutines: make(map[int]*GoroutineInfo),
		channels:   make(map[string]*ChannelInfo),
		dataRaces:  make(map[string]*DataRaceInfo),
		hbGraph: &HappensBeforeGraph{
			Events:    make(map[string]*Event),
			Relations: make(map[string][]string),
		},
		nextGoroutineID: 1, // ID 0 reserved for main goroutine
	}
}

// AnalyzeFile 分析文件的并发安全性
func (ca *ConcurrencyAnalyzer) AnalyzeFile(filename string) error {
	ca.fset = token.NewFileSet()

	// 解析文件
	file, err := parser.ParseFile(ca.fset, filename, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse file: %w", err)
	}

	// 遍历AST
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.GoStmt:
			ca.analyzeGoroutine(node)
		case *ast.CallExpr:
			ca.analyzeChannelOp(node)
			ca.analyzeSyncOp(node)
		case *ast.AssignStmt:
			ca.analyzeAssignment(node)
		case *ast.SendStmt:
			ca.analyzeChannelSend(node)
		}
		return true
	})

	// 执行各种检查
	ca.detectGoroutineLeaks()
	ca.detectChannelDeadlocks()
	ca.detectDataRaces()
	ca.buildHappensBeforeGraph()

	return nil
}

// ====================================================================================
// 第三部分：Goroutine泄露检测
// ====================================================================================

// analyzeGoroutine 分析goroutine创建
func (ca *ConcurrencyAnalyzer) analyzeGoroutine(goStmt *ast.GoStmt) {
	info := &GoroutineInfo{
		ID:       ca.nextGoroutineID,
		Creation: goStmt,
		Position: ca.fset.Position(goStmt.Pos()),
		CanExit:  false,
		WaitedBy: []string{},
	}

	// 提取函数字面量
	if call, ok := goStmt.Call.Fun.(*ast.FuncLit); ok {
		info.Function = call
		info.CanExit = ca.hasExitPath(call.Body)
	}

	ca.goroutines[info.ID] = info
	ca.nextGoroutineID++
}

// hasExitPath 检查函数是否有退出路径
// 形式化：检查是否所有执行路径都会终止
func (ca *ConcurrencyAnalyzer) hasExitPath(body *ast.BlockStmt) bool {
	if body == nil {
		return true
	}

	hasReturn := false
	hasInfiniteLoop := false

	ast.Inspect(body, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.ReturnStmt:
			hasReturn = true
		case *ast.ForStmt:
			// 检查是否为无限循环
			if node.Cond == nil && node.Init == nil && node.Post == nil {
				// 检查循环体内是否有break/return
				hasBreak := false
				ast.Inspect(node.Body, func(inner ast.Node) bool {
					switch inner.(type) {
					case *ast.BranchStmt: // break, continue, goto
						hasBreak = true
						return false
					case *ast.ReturnStmt:
						hasBreak = true
						return false
					}
					return true
				})
				if !hasBreak {
					hasInfiniteLoop = true
				}
			}
		}
		return true
	})

	return hasReturn || !hasInfiniteLoop
}

// detectGoroutineLeaks 检测goroutine泄露
// 形式化定义：
//
//	Leak(g) ⟺ ¬CanExit(g) ∧ WaitedBy(g) = ∅
//	即：goroutine不能退出，且没有被任何机制等待
func (ca *ConcurrencyAnalyzer) detectGoroutineLeaks() {
	for _, g := range ca.goroutines {
		if !g.CanExit && len(g.WaitedBy) == 0 {
			fmt.Printf("⚠️  Goroutine Leak Detected at %s\n", g.Position)
			fmt.Printf("   Goroutine ID: %d\n", g.ID)
			fmt.Printf("   Reason: No exit path and not waited by any sync mechanism\n")
		}
	}
}

// ====================================================================================
// 第四部分：Channel死锁检测
// ====================================================================================

// analyzeChannelOp 分析channel操作
func (ca *ConcurrencyAnalyzer) analyzeChannelOp(call *ast.CallExpr) {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		switch sel.Sel.Name {
		case "make":
			ca.analyzeChannelCreation(call)
		case "close":
			ca.analyzeChannelClose(call)
		}
	}
}

// analyzeChannelCreation 分析channel创建
func (ca *ConcurrencyAnalyzer) analyzeChannelCreation(call *ast.CallExpr) {
	if len(call.Args) < 1 {
		return
	}

	// 检查是否为channel类型
	if chanType, ok := call.Args[0].(*ast.ChanType); ok {
		info := &ChannelInfo{
			Name:     fmt.Sprintf("chan_%d", len(ca.channels)),
			Creation: chanType,
			Buffered: len(call.Args) > 1,
			Sends:    []token.Position{},
			Receives: []token.Position{},
			Closed:   false,
		}

		// 提取缓冲大小
		if info.Buffered {
			if lit, ok := call.Args[1].(*ast.BasicLit); ok {
				fmt.Sscanf(lit.Value, "%d", &info.BufferSize)
			}
		}

		ca.channels[info.Name] = info
	}
}

// analyzeChannelSend 分析channel发送操作
func (ca *ConcurrencyAnalyzer) analyzeChannelSend(send *ast.SendStmt) {
	pos := ca.fset.Position(send.Pos())

	// 记录发送操作
	if ident, ok := send.Chan.(*ast.Ident); ok {
		if ch, exists := ca.channels[ident.Name]; exists {
			ch.Sends = append(ch.Sends, pos)
		}
	}
}

// analyzeChannelClose 分析channel关闭
func (ca *ConcurrencyAnalyzer) analyzeChannelClose(call *ast.CallExpr) {
	if len(call.Args) > 0 {
		if ident, ok := call.Args[0].(*ast.Ident); ok {
			if ch, exists := ca.channels[ident.Name]; exists {
				ch.Closed = true
				ch.ClosePos = ca.fset.Position(call.Pos())
			}
		}
	}
}

// detectChannelDeadlocks 检测channel死锁
// 形式化定义：
//
//	Deadlock(ch) ⟺ (Unbuffered(ch) ∧ |Sends(ch)| > 0 ∧ |Receives(ch)| = 0)
//	               ∨ (Buffered(ch) ∧ |Sends(ch)| > BufferSize(ch) + |Receives(ch)|)
func (ca *ConcurrencyAnalyzer) detectChannelDeadlocks() {
	for name, ch := range ca.channels {
		sends := len(ch.Sends)
		receives := len(ch.Receives)

		deadlock := false
		reason := ""

		if !ch.Buffered {
			// 无缓冲channel：发送但无接收
			if sends > 0 && receives == 0 {
				deadlock = true
				reason = "Unbuffered channel has sends but no receives"
			}
		} else {
			// 有缓冲channel：发送超过缓冲+接收
			if sends > ch.BufferSize+receives {
				deadlock = true
				reason = fmt.Sprintf("Buffered channel (%d) overflows: %d sends, %d receives",
					ch.BufferSize, sends, receives)
			}
		}

		if deadlock {
			fmt.Printf("⚠️  Potential Channel Deadlock: %s\n", name)
			fmt.Printf("   Reason: %s\n", reason)
			fmt.Printf("   Sends: %d, Receives: %d\n", sends, receives)
		}
	}
}

// ====================================================================================
// 第五部分：数据竞争检测
// ====================================================================================

// analyzeAssignment 分析赋值语句（可能的数据竞争）
func (ca *ConcurrencyAnalyzer) analyzeAssignment(assign *ast.AssignStmt) {
	for _, lhs := range assign.Lhs {
		if ident, ok := lhs.(*ast.Ident); ok {
			varName := ident.Name

			if _, exists := ca.dataRaces[varName]; !exists {
				ca.dataRaces[varName] = &DataRaceInfo{
					Variable: varName,
					Accesses: []AccessInfo{},
					IsRace:   false,
				}
			}

			access := AccessInfo{
				Position:    ca.fset.Position(assign.Pos()),
				IsWrite:     true,
				InGoroutine: false, // TODO: 检测是否在goroutine中
				Protected:   false, // TODO: 检测是否被mutex保护
			}

			ca.dataRaces[varName].Accesses = append(ca.dataRaces[varName].Accesses, access)
		}
	}
}

// analyzeSyncOp 分析同步操作（Mutex, WaitGroup等）
func (ca *ConcurrencyAnalyzer) analyzeSyncOp(call *ast.CallExpr) {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		switch sel.Sel.Name {
		case "Lock", "Unlock", "RLock", "RUnlock":
			// Mutex操作
			// TODO: 记录同步操作
		case "Add", "Done", "Wait":
			// WaitGroup操作
			// TODO: 记录WaitGroup操作
		}
	}
}

// detectDataRaces 检测数据竞争
// 形式化定义（Go Memory Model）：
//
//	DataRace(v) ⟺ ∃a1, a2 ∈ Access(v):
//	  a1.goroutine ≠ a2.goroutine ∧
//	  (a1.write ∨ a2.write) ∧
//	  ¬(a1 <HB a2 ∨ a2 <HB a1) ∧
//	  ¬(a1.protected ∧ a2.protected)
func (ca *ConcurrencyAnalyzer) detectDataRaces() {
	for varName, info := range ca.dataRaces {
		accesses := info.Accesses

		// 需要至少2个访问
		if len(accesses) < 2 {
			continue
		}

		// 检查是否有不同goroutine的访问
		hasMultiGoroutine := false
		for i := 0; i < len(accesses); i++ {
			for j := i + 1; j < len(accesses); j++ {
				a1, a2 := accesses[i], accesses[j]

				// 条件1：不同goroutine
				if !a1.InGoroutine && !a2.InGoroutine {
					continue
				}

				// 条件2：至少一个是写操作
				if !a1.IsWrite && !a2.IsWrite {
					continue
				}

				// 条件3：没有happens-before关系
				// TODO: 使用HB图检查

				// 条件4：没有被同步原语保护
				if a1.Protected && a2.Protected {
					continue
				}

				hasMultiGoroutine = true
				info.IsRace = true
				break
			}
			if hasMultiGoroutine {
				break
			}
		}

		if info.IsRace {
			fmt.Printf("⚠️  Potential Data Race on variable: %s\n", varName)
			fmt.Printf("   Accesses:\n")
			for _, access := range accesses {
				accessType := "Read"
				if access.IsWrite {
					accessType = "Write"
				}
				fmt.Printf("     - %s at %s (goroutine: %v, protected: %v)\n",
					accessType, access.Position, access.InGoroutine, access.Protected)
			}
		}
	}
}

// ====================================================================================
// 第六部分：Happens-Before关系建模
// ====================================================================================

// buildHappensBeforeGraph 构建Happens-Before关系图
// 形式化定义（Go Memory Model）：
//  1. 程序顺序：a <HB b if a在b之前执行（同一goroutine内）
//  2. Channel通信：
//     - send(ch) <HB recv(ch)
//     - close(ch) <HB recv(ch) that returns zero value
//  3. 同步原语：
//     - mu.Unlock() <HB mu.Lock()
//     - wg.Done() <HB wg.Wait()
//  4. Goroutine：
//     - go statement <HB goroutine start
//  5. 传递性：a <HB b ∧ b <HB c ⟹ a <HB c
func (ca *ConcurrencyAnalyzer) buildHappensBeforeGraph() {
	// TODO: 实现完整的HB图构建

	// 1. 添加程序顺序关系（同一goroutine内）
	// 2. 添加channel通信关系
	// 3. 添加同步原语关系
	// 4. 添加goroutine启动关系
	// 5. 计算传递闭包
}

// AddHBRelation 添加Happens-Before关系
func (ca *ConcurrencyAnalyzer) AddHBRelation(from, to string) {
	if ca.hbGraph.Relations[from] == nil {
		ca.hbGraph.Relations[from] = []string{}
	}
	ca.hbGraph.Relations[from] = append(ca.hbGraph.Relations[from], to)
}

// HappensBefore 检查e1是否happens-before e2
func (ca *ConcurrencyAnalyzer) HappensBefore(e1, e2 string) bool {
	// 使用DFS或BFS检查可达性
	visited := make(map[string]bool)
	return ca.dfsHB(e1, e2, visited)
}

func (ca *ConcurrencyAnalyzer) dfsHB(current, target string, visited map[string]bool) bool {
	if current == target {
		return true
	}

	if visited[current] {
		return false
	}
	visited[current] = true

	for _, next := range ca.hbGraph.Relations[current] {
		if ca.dfsHB(next, target, visited) {
			return true
		}
	}

	return false
}

// ====================================================================================
// 第七部分：报告生成
// ====================================================================================

// Report 生成并发分析报告
func (ca *ConcurrencyAnalyzer) Report() string {
	var sb strings.Builder

	sb.WriteString("==========================================================================\n")
	sb.WriteString("               并发安全分析报告 (Concurrency Safety Report)\n")
	sb.WriteString("==========================================================================\n\n")

	sb.WriteString("📊 统计信息:\n")
	sb.WriteString(fmt.Sprintf("   - Goroutines: %d\n", len(ca.goroutines)))
	sb.WriteString(fmt.Sprintf("   - Channels: %d\n", len(ca.channels)))
	sb.WriteString(fmt.Sprintf("   - Potential Data Races: %d\n", ca.countDataRaces()))
	sb.WriteString(fmt.Sprintf("   - HB Events: %d\n\n", len(ca.hbGraph.Events)))

	sb.WriteString("🔍 分析结果:\n\n")

	// Goroutine泄露
	leakCount := 0
	for _, g := range ca.goroutines {
		if !g.CanExit && len(g.WaitedBy) == 0 {
			leakCount++
		}
	}
	sb.WriteString(fmt.Sprintf("   1. Goroutine Leaks: %d\n", leakCount))

	// Channel死锁
	deadlockCount := 0
	for _, ch := range ca.channels {
		if (!ch.Buffered && len(ch.Sends) > 0 && len(ch.Receives) == 0) ||
			(ch.Buffered && len(ch.Sends) > ch.BufferSize+len(ch.Receives)) {
			deadlockCount++
		}
	}
	sb.WriteString(fmt.Sprintf("   2. Channel Deadlocks: %d\n", deadlockCount))

	// 数据竞争
	sb.WriteString(fmt.Sprintf("   3. Data Races: %d\n\n", ca.countDataRaces()))

	sb.WriteString("📐 形式化理论基础:\n")
	sb.WriteString("   - Goroutine泄露: Leak(g) ⟺ ¬CanExit(g) ∧ WaitedBy(g) = ∅\n")
	sb.WriteString("   - Channel死锁: Deadlock(ch) ⟺ Unbuffered ∧ Sends > Receives\n")
	sb.WriteString("   - 数据竞争: DataRace(v) ⟺ ∃concurrent accesses ∧ ¬(a1 <HB a2)\n")
	sb.WriteString("   - Happens-Before: send(ch) <HB recv(ch), unlock <HB lock\n\n")

	sb.WriteString("==========================================================================\n")

	return sb.String()
}

func (ca *ConcurrencyAnalyzer) countDataRaces() int {
	count := 0
	for _, info := range ca.dataRaces {
		if info.IsRace {
			count++
		}
	}
	return count
}

// GetGoroutines 获取所有goroutine信息
func (ca *ConcurrencyAnalyzer) GetGoroutines() map[int]*GoroutineInfo {
	return ca.goroutines
}

// GetChannels 获取所有channel信息
func (ca *ConcurrencyAnalyzer) GetChannels() map[string]*ChannelInfo {
	return ca.channels
}

// GetDataRaces 获取所有数据竞争信息
func (ca *ConcurrencyAnalyzer) GetDataRaces() map[string]*DataRaceInfo {
	return ca.dataRaces
}

// GetHappensBeforeGraph 获取Happens-Before关系图
func (ca *ConcurrencyAnalyzer) GetHappensBeforeGraph() *HappensBeforeGraph {
	return ca.hbGraph
}

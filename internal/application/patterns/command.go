package patterns

import "context"

// Command 命令接口（框架抽象）
// 用于处理写操作（Create、Update、Delete）
type Command interface {
	// Execute 执行命令
	Execute(ctx context.Context) error
}

// CommandHandler 命令处理器接口（框架抽象）
type CommandHandler[T Command] interface {
	// Handle 处理命令
	Handle(ctx context.Context, cmd T) error
}

// CommandResult 命令执行结果
type CommandResult struct {
	Success bool
	Message string
	Data    interface{}
}

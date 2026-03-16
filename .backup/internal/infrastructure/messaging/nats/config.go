package nats

import "time"

// Config NATS 客户端配置
type Config struct {
	URL            string        // NATS 服务器地址，例如: "nats://localhost:4222"
	MaxReconnects  int           // 最大重连次数，-1 表示无限重连
	ReconnectWait  time.Duration // 重连等待时间
	Timeout        time.Duration // 连接超时
	Name           string        // 客户端名称
	Token          string        // 认证 Token（可选）
	Username       string        // 用户名（可选）
	Password       string        // 密码（可选）
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		URL:           "nats://localhost:4222",
		MaxReconnects: -1,
		ReconnectWait: 2 * time.Second,
		Timeout:       5 * time.Second,
	}
}

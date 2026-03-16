// Package config provides configuration management for the application.
//
// 配置管理包提供了统一的配置加载和管理机制，支持：
// 1. 多种配置源：配置文件（YAML）、环境变量、默认值
// 2. 配置优先级：环境变量 > 配置文件 > 默认值
// 3. 热重载：支持配置文件变化时自动重新加载
// 4. 类型安全：使用结构体定义配置，确保类型安全
//
// 设计原则：
// 1. 统一管理：所有配置项集中在一个 Config 结构体中
// 2. 灵活加载：支持多种加载方式（文件、环境变量、默认值）
// 3. 易于扩展：可以轻松添加新的配置项
// 4. 类型安全：使用结构体和 mapstructure 标签确保类型安全
//
// 使用场景：
// - 应用启动时加载配置
// - 开发环境使用配置文件
// - 生产环境使用环境变量
// - 动态调整配置（热重载）
//
// 示例：
//
//	// 加载配置
//	cfg, err := config.LoadConfig()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// 使用配置
//	server := http.Server{
//	    Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
//	    ReadTimeout:  cfg.Server.ReadTimeout,
//	    WriteTimeout: cfg.Server.WriteTimeout,
//	}
package config

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config 是应用程序的完整配置结构。
//
// 功能说明：
// - 包含所有子系统的配置
// - 使用 mapstructure 标签支持 YAML 和 Viper 绑定
// - 支持环境变量覆盖
//
// 配置项说明：
// - Server: HTTP/gRPC 服务器配置
// - Database: 数据库连接配置
// - Redis: Redis 缓存配置
// - Kafka: Kafka 消息队列配置
// - MQTT: MQTT 消息队列配置
// - OTLP: OpenTelemetry 可观测性配置
// - JWT: JWT 认证配置
// - Logging: 日志配置
// - Temporal: Temporal 工作流配置
//
// 配置文件示例（YAML）：
//
//	server:
//	  host: 0.0.0.0
//	  port: 8080
//	database:
//	  type: postgres
//	  host: localhost
//	  port: 5432
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	Kafka      KafkaConfig      `mapstructure:"kafka"`
	MQTT       MQTTConfig       `mapstructure:"mqtt"`
	OTLP           OTLPConfig           `mapstructure:"otlp"`
	Observability  ObservabilityConfig  `mapstructure:"observability"`
	JWT            JWTConfig            `mapstructure:"jwt"`
	Logging        LoggingConfig        `mapstructure:"logging"`
	Temporal       TemporalConfig       `mapstructure:"temporal"`
}

// ServerConfig 是 HTTP/gRPC 服务器的配置。
//
// 字段说明：
// - Host: 服务器监听地址（默认：0.0.0.0）
// - Port: 服务器监听端口（默认：8080）
// - ReadTimeout: 读取超时时间（默认：30s）
// - WriteTimeout: 写入超时时间（默认：30s）
// - IdleTimeout: 空闲连接超时时间（默认：120s）
//
// 环境变量：
// - APP_SERVER_HOST: 服务器地址
// - APP_SERVER_PORT: 服务器端口
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// DatabaseConfig 是数据库连接的配置。
//
// 字段说明：
// - Type: 数据库类型（postgres、sqlite3）
// - Host: 数据库主机地址（PostgreSQL）
// - Port: 数据库端口（PostgreSQL，默认：5432）
// - User: 数据库用户名
// - Password: 数据库密码
// - Database: 数据库名称
// - SSLMode: SSL 模式（PostgreSQL：disable、require、verify-full）
// - MaxOpenConns: 最大打开连接数（默认：25）
// - MaxIdleConns: 最大空闲连接数（默认：5）
// - DSN: SQLite3 数据源名称（SQLite3 专用）
//
// 环境变量：
// - APP_DB_TYPE: 数据库类型
// - APP_DB_HOST: 数据库主机
// - APP_DB_PORT: 数据库端口
// - APP_DB_USER: 数据库用户名
// - APP_DB_PASSWORD: 数据库密码
// - APP_DB_NAME: 数据库名称
// - APP_DB_DSN: SQLite3 DSN
type DatabaseConfig struct {
	Type         string `mapstructure:"type"` // postgres, sqlite3
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	SSLMode      string `mapstructure:"ssl_mode"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	// SQLite3 配置
	DSN string `mapstructure:"dsn"` // SQLite3 DSN
}

// RedisConfig 是 Redis 缓存的配置。
//
// 字段说明：
// - Addr: Redis 服务器地址（格式：host:port，默认：localhost:6379）
// - Password: Redis 密码（可选）
// - DB: 数据库编号（默认：0）
// - PoolSize: 连接池大小（默认：10）
// - MinIdleConns: 最小空闲连接数（默认：5）
// - DialTimeout: 连接超时时间
// - ReadTimeout: 读取超时时间
// - WriteTimeout: 写入超时时间
//
// 环境变量：
// - APP_REDIS_ADDR: Redis 地址
// - APP_REDIS_PASSWORD: Redis 密码
// - APP_REDIS_DB: 数据库编号
type RedisConfig struct {
	Addr         string        `mapstructure:"addr"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// KafkaConfig 是 Kafka 消息队列的配置。
//
// 字段说明：
// - Brokers: Kafka Broker 地址列表（默认：["localhost:9092"]）
// - ConsumerGroup: 消费者组 ID
// - Timeout: 操作超时时间
//
// 环境变量：
// - APP_KAFKA_BROKERS: Kafka Broker 地址（逗号分隔）
type KafkaConfig struct {
	Brokers       []string      `mapstructure:"brokers"`
	ConsumerGroup string        `mapstructure:"consumer_group"`
	Timeout       time.Duration `mapstructure:"timeout"`
}

// MQTTConfig 是 MQTT 消息队列的配置。
//
// 字段说明：
// - Broker: MQTT Broker 地址（格式：tcp://host:port，默认：tcp://localhost:1883）
// - ClientID: 客户端 ID（必须唯一）
// - Username: 用户名（可选）
// - Password: 密码（可选）
// - Timeout: 操作超时时间
//
// 环境变量：
// - APP_MQTT_BROKER: MQTT Broker 地址
// - APP_MQTT_CLIENT_ID: 客户端 ID
type MQTTConfig struct {
	Broker   string        `mapstructure:"broker"`
	ClientID string        `mapstructure:"client_id"`
	Username string        `mapstructure:"username"`
	Password string        `mapstructure:"password"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

// OTLPConfig 是 OpenTelemetry 可观测性的配置。
//
// 字段说明：
// - Endpoint: OpenTelemetry Collector 端点（默认：localhost:4317）
// - Insecure: 是否使用不安全的连接（开发环境：true，生产环境：false）
// - Timeout: 操作超时时间
// - ServiceName: 服务名称（默认：app）
// - ServiceVersion: 服务版本（默认：1.0.0）
//
// 环境变量：
// - APP_OTLP_ENDPOINT: Collector 端点
// - APP_OTLP_SERVICE_NAME: 服务名称
type OTLPConfig struct {
	Endpoint    string        `mapstructure:"endpoint"`
	Insecure    bool          `mapstructure:"insecure"`
	Timeout     time.Duration `mapstructure:"timeout"`
	ServiceName string        `mapstructure:"service_name"`
	ServiceVersion string     `mapstructure:"service_version"`
}

// JWTConfig 是 JWT 认证的配置。
//
// 字段说明：
// - SecretKey: JWT 签名密钥（必须设置，建议使用环境变量）
// - SigningMethod: 签名方法（默认：HS256）
// - AccessTokenTTL: 访问令牌有效期（默认：15分钟）
// - RefreshTokenTTL: 刷新令牌有效期（默认：7天）
//
// 环境变量：
// - APP_JWT_SECRET_KEY: JWT 签名密钥
//
// 注意事项：
// - SecretKey 应该保密，不应提交到版本控制系统
// - 生产环境应使用强随机密钥
type JWTConfig struct {
	SecretKey      string        `mapstructure:"secret_key"`
	SigningMethod  string        `mapstructure:"signing_method"`
	AccessTokenTTL time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`
}

// LoggingConfig 是日志系统的配置。
//
// 字段说明：
// - Level: 日志级别（debug、info、warn、error，默认：info）
// - Format: 日志格式（json、text，默认：json）
// - Output: 输出目标（stdout、file，默认：stdout）
// - OutputPath: 日志文件路径（当 Output 为 file 时使用）
// - Rotation: 日志轮转配置（仅当 Output 为 file 时生效）
//
// 环境变量：
// - APP_LOG_LEVEL: 日志级别
// - APP_LOG_FORMAT: 日志格式
// - APP_LOG_OUTPUT: 输出目标
// - APP_LOG_OUTPUT_PATH: 日志文件路径
// - APP_LOG_ROTATION_MAX_SIZE: 单个日志文件最大大小（MB）
// - APP_LOG_ROTATION_MAX_BACKUPS: 保留的旧日志文件数量
// - APP_LOG_ROTATION_MAX_AGE: 保留旧日志文件的天数
// - APP_LOG_ROTATION_COMPRESS: 是否压缩旧日志文件
type LoggingConfig struct {
	Level      string         `mapstructure:"level"`       // debug, info, warn, error
	Format     string         `mapstructure:"format"`       // json, text
	Output     string         `mapstructure:"output"`       // stdout, file
	OutputPath string         `mapstructure:"output_path"`  // 日志文件路径
	Rotation   RotationConfig `mapstructure:"rotation"`    // 日志轮转配置
}

// RotationConfig 日志轮转配置
type RotationConfig struct {
	// MaxSize 单个日志文件的最大大小（MB），超过此大小会轮转
	MaxSize int `mapstructure:"max_size"`
	// MaxBackups 保留的旧日志文件数量
	MaxBackups int `mapstructure:"max_backups"`
	// MaxAge 保留旧日志文件的天数
	MaxAge int `mapstructure:"max_age"`
	// Compress 是否压缩轮转后的旧日志文件
	Compress bool `mapstructure:"compress"`
}

// ObservabilityConfig 是可观测性的完整配置。
//
// 字段说明：
// - OTLP: OTLP 配置（已存在，这里扩展）
// - System: 系统监控配置
//
// 环境变量：
// - APP_OBSERVABILITY_SYSTEM_ENABLED: 是否启用系统监控
// - APP_OBSERVABILITY_SYSTEM_COLLECT_INTERVAL: 系统监控收集间隔
type ObservabilityConfig struct {
	// OTLP 配置（复用现有的 OTLPConfig）
	OTLP OTLPConfig `mapstructure:"otlp"`
	
	// 系统监控配置
	System SystemMonitoringConfig `mapstructure:"system"`
}

// SystemMonitoringConfig 是系统监控的配置。
//
// 字段说明：
// - Enabled: 是否启用系统监控（默认：false）
// - CollectInterval: 收集间隔（默认：5s）
// - EnableDiskMonitor: 是否启用磁盘监控（默认：false）
// - EnableLoadMonitor: 是否启用负载监控（默认：false）
// - EnableAPMMonitor: 是否启用 APM 监控（默认：false）
// - RateLimit: 限流器配置
// - HealthThresholds: 健康检查阈值
// - Alerts: 告警规则配置
//
// 环境变量：
// - APP_OBSERVABILITY_SYSTEM_ENABLED: 是否启用系统监控
// - APP_OBSERVABILITY_SYSTEM_COLLECT_INTERVAL: 收集间隔
type SystemMonitoringConfig struct {
	Enabled            bool              `mapstructure:"enabled"`
	CollectInterval    string            `mapstructure:"collect_interval"` // 如 "5s"
	EnableDiskMonitor  bool              `mapstructure:"enable_disk_monitor"`
	EnableLoadMonitor  bool              `mapstructure:"enable_load_monitor"`
	EnableAPMMonitor   bool              `mapstructure:"enable_apm_monitor"`
	RateLimit          RateLimitConfig   `mapstructure:"rate_limit"`
	HealthThresholds   HealthThresholdsConfig `mapstructure:"health_thresholds"`
	Alerts             []AlertRuleConfig `mapstructure:"alerts"`
}

// RateLimitConfig 限流器配置
type RateLimitConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Limit   int64  `mapstructure:"limit"`
	Window  string `mapstructure:"window"` // 如 "1s"
}

// HealthThresholdsConfig 健康检查阈值配置
type HealthThresholdsConfig struct {
	MaxMemoryUsage float64 `mapstructure:"max_memory_usage"`
	MaxCPUUsage    float64 `mapstructure:"max_cpu_usage"`
	MaxGoroutines  int     `mapstructure:"max_goroutines"`
}

// AlertRuleConfig 告警规则配置
type AlertRuleConfig struct {
	ID         string  `mapstructure:"id"`
	Name       string  `mapstructure:"name"`
	MetricName string  `mapstructure:"metric_name"`
	Condition  string  `mapstructure:"condition"` // gt, lt, eq, gte, lte
	Threshold  float64 `mapstructure:"threshold"`
	Level      string  `mapstructure:"level"` // info, warning, critical
	Enabled    bool    `mapstructure:"enabled"`
	Duration   string  `mapstructure:"duration"` // 如 "5m"
	Cooldown   string  `mapstructure:"cooldown"` // 如 "10m"
}

// TemporalConfig 是 Temporal 工作流的配置。
//
// 字段说明：
// - Address: Temporal Server 地址（默认：localhost:7233）
// - Namespace: 命名空间（默认：default）
// - TaskQueue: 任务队列名称（默认：default）
// - Workers: Worker 数量（默认：10）
// - MaxConcurrent: 最大并发任务数（默认：100）
//
// 环境变量：
// - APP_TEMPORAL_ADDRESS: Temporal Server 地址
// - APP_TEMPORAL_NAMESPACE: 命名空间
// - APP_TEMPORAL_TASK_QUEUE: 任务队列
type TemporalConfig struct {
	Address      string        `mapstructure:"address"`
	Namespace    string        `mapstructure:"namespace"`
	TaskQueue    string        `mapstructure:"task_queue"`
	Workers      int           `mapstructure:"workers"`
	MaxConcurrent int          `mapstructure:"max_concurrent"`
}

// Load 加载配置
//
// 设计原理：
// 1. 使用 Viper 进行配置管理
// 2. 支持多种配置源：配置文件、环境变量、默认值
// 3. 配置优先级：环境变量 > 配置文件 > 默认值
// 4. 如果配置文件不存在，使用默认配置（不报错）
//
// 参数：
//   - configPath: 配置文件路径，如果为空则使用默认路径
//
// 返回：
//   - *Config: 加载的配置对象
//   - error: 加载失败时返回错误
//
// 配置加载流程：
// 1. 创建 Viper 实例
// 2. 设置配置文件路径（如果指定）或使用默认路径
// 3. 设置环境变量前缀（APP_）
// 4. 读取配置文件（如果存在）
// 5. 绑定环境变量
// 6. 解析配置到结构体
// 7. 设置默认值
//
// 配置文件路径：
// - 如果 configPath 不为空，使用指定路径
// - 否则，按以下顺序查找：
//   1. ./configs/config.yaml
//   2. ./config.yaml
//
// 环境变量：
// - 前缀：APP_
// - 格式：APP_SERVER_PORT、APP_DATABASE_HOST 等
// - 支持嵌套配置：APP_DATABASE_HOST、APP_DATABASE_PORT
//
// 示例：
//   // 使用默认路径
//   cfg, err := config.Load("")
//
//   // 使用指定路径
//   cfg, err := config.Load("/path/to/config.yaml")
//
//   // 使用环境变量覆盖
//   // export APP_SERVER_PORT=8080
//   // export APP_DATABASE_HOST=localhost
//   cfg, err := config.Load("")
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// 设置配置文件路径
	// 如果指定了路径，使用指定路径；否则使用默认路径
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		// 默认配置文件路径
		// 按顺序查找：./configs/config.yaml、./config.yaml
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath("./configs")
		v.AddConfigPath(".")
	}

	// 设置环境变量
	// 前缀：APP_，例如 APP_SERVER_PORT、APP_DATABASE_HOST
	v.SetEnvPrefix("APP")
	// 自动读取环境变量
	v.AutomaticEnv()

	// 读取配置文件
	// 如果配置文件不存在，不报错（使用默认配置）
	if err := v.ReadInConfig(); err != nil {
		// 如果配置文件不存在，使用默认配置
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// 其他错误（如格式错误）需要返回
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// 环境变量覆盖
	// 绑定环境变量到配置项
	bindEnvVars(v)

	// 解析配置到结构体
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 设置默认值
	// 对于未设置的配置项，使用默认值
	setDefaults(&config)

	return &config, nil
}

// bindEnvVars 绑定环境变量
func bindEnvVars(v *viper.Viper) {
	// Server
	v.BindEnv("server.host", "APP_SERVER_HOST")
	v.BindEnv("server.port", "APP_SERVER_PORT")

	// Database
	v.BindEnv("database.type", "APP_DB_TYPE")
	v.BindEnv("database.host", "APP_DB_HOST")
	v.BindEnv("database.port", "APP_DB_PORT")
	v.BindEnv("database.user", "APP_DB_USER")
	v.BindEnv("database.password", "APP_DB_PASSWORD")
	v.BindEnv("database.database", "APP_DB_NAME")
	v.BindEnv("database.dsn", "APP_DB_DSN")

	// Redis
	v.BindEnv("redis.addr", "APP_REDIS_ADDR")
	v.BindEnv("redis.password", "APP_REDIS_PASSWORD")
	v.BindEnv("redis.db", "APP_REDIS_DB")

	// Kafka
	v.BindEnv("kafka.brokers", "APP_KAFKA_BROKERS")

	// MQTT
	v.BindEnv("mqtt.broker", "APP_MQTT_BROKER")
	v.BindEnv("mqtt.client_id", "APP_MQTT_CLIENT_ID")

	// OTLP
	v.BindEnv("otlp.endpoint", "APP_OTLP_ENDPOINT")
	v.BindEnv("otlp.service_name", "APP_OTLP_SERVICE_NAME")

	// JWT
	v.BindEnv("jwt.secret_key", "APP_JWT_SECRET_KEY")

	// Logging
	v.BindEnv("logging.level", "APP_LOG_LEVEL")
	v.BindEnv("logging.format", "APP_LOG_FORMAT")
}

// setDefaults 设置默认值
func setDefaults(c *Config) {
	// Server 默认值
	if c.Server.Host == "" {
		c.Server.Host = "0.0.0.0"
	}
	if c.Server.Port == 0 {
		c.Server.Port = 8080
	}
	if c.Server.ReadTimeout == 0 {
		c.Server.ReadTimeout = 30 * time.Second
	}
	if c.Server.WriteTimeout == 0 {
		c.Server.WriteTimeout = 30 * time.Second
	}
	if c.Server.IdleTimeout == 0 {
		c.Server.IdleTimeout = 120 * time.Second
	}

	// Database 默认值
	if c.Database.Type == "" {
		c.Database.Type = "postgres"
	}
	if c.Database.Host == "" {
		c.Database.Host = "localhost"
	}
	if c.Database.Port == 0 {
		if c.Database.Type == "postgres" {
			c.Database.Port = 5432
		} else if c.Database.Type == "sqlite3" {
			c.Database.DSN = "file:app.db?cache=shared&mode=rwc"
		}
	}
	if c.Database.MaxOpenConns == 0 {
		c.Database.MaxOpenConns = 25
	}
	if c.Database.MaxIdleConns == 0 {
		c.Database.MaxIdleConns = 5
	}

	// Redis 默认值
	if c.Redis.Addr == "" {
		c.Redis.Addr = "localhost:6379"
	}
	if c.Redis.PoolSize == 0 {
		c.Redis.PoolSize = 10
	}
	if c.Redis.MinIdleConns == 0 {
		c.Redis.MinIdleConns = 5
	}

	// Kafka 默认值
	if len(c.Kafka.Brokers) == 0 {
		c.Kafka.Brokers = []string{"localhost:9092"}
	}

	// MQTT 默认值
	if c.MQTT.Broker == "" {
		c.MQTT.Broker = "tcp://localhost:1883"
	}

	// OTLP 默认值
	if c.OTLP.Endpoint == "" {
		c.OTLP.Endpoint = "localhost:4317"
	}
	if c.OTLP.ServiceName == "" {
		c.OTLP.ServiceName = "app"
	}
	if c.OTLP.ServiceVersion == "" {
		c.OTLP.ServiceVersion = "1.0.0"
	}

	// JWT 默认值
	if c.JWT.SigningMethod == "" {
		c.JWT.SigningMethod = "HS256"
	}
	if c.JWT.AccessTokenTTL == 0 {
		c.JWT.AccessTokenTTL = 15 * time.Minute
	}
	if c.JWT.RefreshTokenTTL == 0 {
		c.JWT.RefreshTokenTTL = 7 * 24 * time.Hour
	}

	// Logging 默认值
	if c.Logging.Level == "" {
		c.Logging.Level = "info"
	}
	if c.Logging.Format == "" {
		c.Logging.Format = "json"
	}
	if c.Logging.Output == "" {
		c.Logging.Output = "stdout"
	}

	// Temporal 默认值
	if c.Temporal.Address == "" {
		c.Temporal.Address = "localhost:7233"
	}
	if c.Temporal.Namespace == "" {
		c.Temporal.Namespace = "default"
	}
	if c.Temporal.TaskQueue == "" {
		c.Temporal.TaskQueue = "default"
	}
	if c.Temporal.Workers == 0 {
		c.Temporal.Workers = 10
	}
	if c.Temporal.MaxConcurrent == 0 {
		c.Temporal.MaxConcurrent = 100
	}

	// Observability 默认值
	// 如果 Observability.OTLP 未设置，使用 OTLP 配置
	if c.Observability.OTLP.Endpoint == "" {
		c.Observability.OTLP = c.OTLP
	}
	if c.Observability.OTLP.Endpoint == "" {
		c.Observability.OTLP.Endpoint = "localhost:4317"
	}
	if c.Observability.OTLP.ServiceName == "" {
		c.Observability.OTLP.ServiceName = "app"
	}
	if c.Observability.OTLP.ServiceVersion == "" {
		c.Observability.OTLP.ServiceVersion = "1.0.0"
	}

	// 系统监控默认值
	if c.Observability.System.CollectInterval == "" {
		c.Observability.System.CollectInterval = "5s"
	}
	if c.Observability.System.HealthThresholds.MaxMemoryUsage == 0 {
		c.Observability.System.HealthThresholds.MaxMemoryUsage = 90.0
	}
	if c.Observability.System.HealthThresholds.MaxCPUUsage == 0 {
		c.Observability.System.HealthThresholds.MaxCPUUsage = 95.0
	}
	if c.Observability.System.HealthThresholds.MaxGoroutines == 0 {
		c.Observability.System.HealthThresholds.MaxGoroutines = 10000
	}
}

// LoadConfig 加载配置（便捷函数）
//
// 设计原理：
// 1. 这是 Load("") 的便捷函数
// 2. 使用默认配置文件路径
// 3. 支持环境变量覆盖
//
// 返回：
//   - *Config: 加载的配置对象
//   - error: 加载失败时返回错误
//
// 使用示例：
//   cfg, err := config.LoadConfig()
//   if err != nil {
//       log.Fatal(err)
//   }
func LoadConfig() (*Config, error) {
	return Load("")
}

// LoadFromEnv 从环境变量加载配置
//
// 设计原理：
// 1. 只从环境变量加载配置，不使用配置文件
// 2. 适用于容器化部署场景
// 3. 所有配置都通过环境变量提供
//
// 返回：
//   - *Config: 加载的配置对象
//   - error: 加载失败时返回错误
//
// 使用示例：
//   // 设置环境变量
//   // export APP_SERVER_PORT=8080
//   // export APP_DATABASE_HOST=localhost
//   cfg, err := config.LoadFromEnv()
func LoadFromEnv() (*Config, error) {
	return Load("")
}

// LoadFromFile 从文件加载配置
//
// 设计原理：
// 1. 从指定文件路径加载配置
// 2. 仍然支持环境变量覆盖
// 3. 适用于需要指定配置文件路径的场景
//
// 参数：
//   - path: 配置文件路径
//
// 返回：
//   - *Config: 加载的配置对象
//   - error: 加载失败时返回错误
//
// 使用示例：
//   cfg, err := config.LoadFromFile("/path/to/config.yaml")
func LoadFromFile(path string) (*Config, error) {
	return Load(path)
}

// Watch 监听配置文件变化（热重载）
//
// 设计原理：
// 1. 监听配置文件的变化
// 2. 当配置文件发生变化时，自动重新加载配置
// 3. 调用回调函数通知配置变化
//
// 参数：
//   - configPath: 配置文件路径，如果为空则使用默认路径
//   - onChange: 配置变化时的回调函数
//
// 返回：
//   - error: 监听失败时返回错误
//
// 使用场景：
// 1. 开发环境：修改配置后自动生效，无需重启
// 2. 生产环境：动态调整配置（如日志级别）
//
// 注意事项：
// - 热重载可能导致配置不一致
// - 某些配置（如数据库连接）需要重启才能生效
// - 建议只用于非关键配置的热重载
//
// 使用示例：
//   err := config.Watch("", func(cfg *config.Config) {
//       log.Printf("Config reloaded: %+v", cfg)
//       // 更新应用配置
//   })
func Watch(configPath string, onChange func(*Config)) error {
	v := viper.New()

	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath("./configs")
		v.AddConfigPath(".")
	}

	// 监听配置文件变化
	v.WatchConfig()
	// 配置变化时的回调
	v.OnConfigChange(func(e fsnotify.Event) {
		var config Config
		if err := v.Unmarshal(&config); err == nil {
			setDefaults(&config)
			onChange(&config)
		}
	})

	return nil
}

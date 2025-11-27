package config

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	Kafka      KafkaConfig      `mapstructure:"kafka"`
	MQTT       MQTTConfig       `mapstructure:"mqtt"`
	OTLP       OTLPConfig       `mapstructure:"otlp"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Logging    LoggingConfig    `mapstructure:"logging"`
	Temporal   TemporalConfig   `mapstructure:"temporal"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// DatabaseConfig 数据库配置
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

// RedisConfig Redis 配置
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

// KafkaConfig Kafka 配置
type KafkaConfig struct {
	Brokers       []string      `mapstructure:"brokers"`
	ConsumerGroup string        `mapstructure:"consumer_group"`
	Timeout       time.Duration `mapstructure:"timeout"`
}

// MQTTConfig MQTT 配置
type MQTTConfig struct {
	Broker   string        `mapstructure:"broker"`
	ClientID string        `mapstructure:"client_id"`
	Username string        `mapstructure:"username"`
	Password string        `mapstructure:"password"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

// OTLPConfig OpenTelemetry 配置
type OTLPConfig struct {
	Endpoint    string        `mapstructure:"endpoint"`
	Insecure    bool          `mapstructure:"insecure"`
	Timeout     time.Duration `mapstructure:"timeout"`
	ServiceName string        `mapstructure:"service_name"`
	ServiceVersion string     `mapstructure:"service_version"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	SecretKey      string        `mapstructure:"secret_key"`
	SigningMethod  string        `mapstructure:"signing_method"`
	AccessTokenTTL time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level      string `mapstructure:"level"`       // debug, info, warn, error
	Format     string `mapstructure:"format"`      // json, text
	Output     string `mapstructure:"output"`      // stdout, file
	OutputPath string `mapstructure:"output_path"` // 日志文件路径
}

// TemporalConfig Temporal 配置
type TemporalConfig struct {
	Address      string        `mapstructure:"address"`
	Namespace    string        `mapstructure:"namespace"`
	TaskQueue    string        `mapstructure:"task_queue"`
	Workers      int           `mapstructure:"workers"`
	MaxConcurrent int          `mapstructure:"max_concurrent"`
}

// Load 加载配置
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// 设置配置文件路径
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		// 默认配置文件路径
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath("./configs")
		v.AddConfigPath(".")
	}

	// 设置环境变量
	v.SetEnvPrefix("APP")
	v.AutomaticEnv()

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		// 如果配置文件不存在，使用默认配置
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// 环境变量覆盖
	bindEnvVars(v)

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 设置默认值
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
}

// LoadFromEnv 从环境变量加载配置
func LoadFromEnv() (*Config, error) {
	return Load("")
}

// LoadFromFile 从文件加载配置
func LoadFromFile(path string) (*Config, error) {
	return Load(path)
}

// Watch 监听配置文件变化（热重载）
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

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		var config Config
		if err := v.Unmarshal(&config); err == nil {
			setDefaults(&config)
			onChange(&config)
		}
	})

	return nil
}

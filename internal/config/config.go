package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server        ServerConfig        `mapstructure:"server"`
	Database      DatabaseConfig      `mapstructure:"database"`
	Log           LogConfig           `mapstructure:"log"`
	Observability ObservabilityConfig `mapstructure:"observability"`
	Workflow      WorkflowConfig      `mapstructure:"workflow"`
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
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `mapstructure:"level"`
}

// ObservabilityConfig 可观测性配置
type ObservabilityConfig struct {
	TraceEndpoint  string `mapstructure:"trace_endpoint"`
	MetricEndpoint string `mapstructure:"metric_endpoint"`
	Enabled        bool   `mapstructure:"enabled"`
}

// WorkflowConfig 工作流配置
type WorkflowConfig struct {
	Temporal TemporalConfig `mapstructure:"temporal"`
}

// TemporalConfig Temporal 配置
type TemporalConfig struct {
	Address   string `mapstructure:"address"`
	TaskQueue string `mapstructure:"task_queue"`
	Namespace string `mapstructure:"namespace"`
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// 设置默认值
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", 5*time.Second)
	viper.SetDefault("server.write_timeout", 10*time.Second)
	viper.SetDefault("server.idle_timeout", 120*time.Second)

	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "user")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.dbname", "golang")
	viper.SetDefault("database.sslmode", "disable")

	viper.SetDefault("log.level", "info")

	viper.SetDefault("observability.trace_endpoint", "localhost:4317")
	viper.SetDefault("observability.metric_endpoint", "localhost:4317")
	viper.SetDefault("observability.enabled", true)

	// Workflow defaults
	viper.SetDefault("workflow.temporal.address", "localhost:7233")
	viper.SetDefault("workflow.temporal.task_queue", "user-task-queue")
	viper.SetDefault("workflow.temporal.namespace", "default")

	// 读取环境变量
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

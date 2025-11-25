package config

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Database    DatabaseConfig    `mapstructure:"database"`
	Observability ObservabilityConfig `mapstructure:"observability"`
	Messaging   MessagingConfig   `mapstructure:"messaging"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string `mapstructure:"host" env:"SERVER_HOST"`
	Port         int    `mapstructure:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `mapstructure:"read_timeout" env:"SERVER_READ_TIMEOUT"`
	WriteTimeout int    `mapstructure:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host" env:"DB_HOST"`
	Port     int    `mapstructure:"port" env:"DB_PORT"`
	User     string `mapstructure:"user" env:"DB_USER"`
	Password string `mapstructure:"password" env:"DB_PASSWORD"`
	DBName   string `mapstructure:"dbname" env:"DB_NAME"`
	SSLMode  string `mapstructure:"sslmode" env:"DB_SSLMODE"`
	MaxConns int    `mapstructure:"max_conns" env:"DB_MAX_CONNS"`
}

// ObservabilityConfig 可观测性配置
type ObservabilityConfig struct {
	OTLP OTLPConfig `mapstructure:"otlp"`
}

// OTLPConfig OTLP配置
type OTLPConfig struct {
	Endpoint string `mapstructure:"endpoint" env:"OTLP_ENDPOINT"`
	Insecure bool   `mapstructure:"insecure" env:"OTLP_INSECURE"`
}

// MessagingConfig 消息队列配置
type MessagingConfig struct {
	Kafka KafkaConfig `mapstructure:"kafka"`
	MQTT  MQTTConfig  `mapstructure:"mqtt"`
}

// KafkaConfig Kafka配置
type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers" env:"KAFKA_BROKERS"`
}

// MQTTConfig MQTT配置
type MQTTConfig struct {
	Broker   string `mapstructure:"broker" env:"MQTT_BROKER"`
	ClientID string `mapstructure:"client_id" env:"MQTT_CLIENT_ID"`
	Username string `mapstructure:"username" env:"MQTT_USERNAME"`
	Password string `mapstructure:"password" env:"MQTT_PASSWORD"`
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// 设置默认值
	setDefaults()

	// 读取环境变量
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		slog.Warn("config file not found, using defaults and environment variables", "error", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// setDefaults 设置默认值
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.dbname", "golang")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_conns", 25)

	// OTLP defaults
	viper.SetDefault("observability.otlp.endpoint", "localhost:4317")
	viper.SetDefault("observability.otlp.insecure", true)
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

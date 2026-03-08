package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoad_Defaults 测试默认配置加载
func TestLoad_Defaults(t *testing.T) {
	cfg, err := Load("")
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// 验证 Server 默认值
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, 30*time.Second, cfg.Server.ReadTimeout)
	assert.Equal(t, 30*time.Second, cfg.Server.WriteTimeout)
	assert.Equal(t, 120*time.Second, cfg.Server.IdleTimeout)

	// 验证 Database 默认值
	assert.Equal(t, "postgres", cfg.Database.Type)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.Equal(t, 25, cfg.Database.MaxOpenConns)
	assert.Equal(t, 5, cfg.Database.MaxIdleConns)

	// 验证 Redis 默认值
	assert.Equal(t, "localhost:6379", cfg.Redis.Addr)
	assert.Equal(t, 10, cfg.Redis.PoolSize)
	assert.Equal(t, 5, cfg.Redis.MinIdleConns)

	// 验证 Kafka 默认值
	assert.Equal(t, []string{"localhost:9092"}, cfg.Kafka.Brokers)

	// 验证 MQTT 默认值
	assert.Equal(t, "tcp://localhost:1883", cfg.MQTT.Broker)

	// 验证 OTLP 默认值
	assert.Equal(t, "localhost:4317", cfg.OTLP.Endpoint)
	assert.Equal(t, "app", cfg.OTLP.ServiceName)
	assert.Equal(t, "1.0.0", cfg.OTLP.ServiceVersion)

	// 验证 JWT 默认值
	assert.Equal(t, "HS256", cfg.JWT.SigningMethod)
	assert.Equal(t, 15*time.Minute, cfg.JWT.AccessTokenTTL)
	assert.Equal(t, 7*24*time.Hour, cfg.JWT.RefreshTokenTTL)

	// 验证 Logging 默认值
	assert.Equal(t, "info", cfg.Logging.Level)
	assert.Equal(t, "json", cfg.Logging.Format)
	assert.Equal(t, "stdout", cfg.Logging.Output)

	// 验证 Temporal 默认值
	assert.Equal(t, "localhost:7233", cfg.Temporal.Address)
	assert.Equal(t, "default", cfg.Temporal.Namespace)
	assert.Equal(t, "default", cfg.Temporal.TaskQueue)
	assert.Equal(t, 10, cfg.Temporal.Workers)
	assert.Equal(t, 100, cfg.Temporal.MaxConcurrent)
}

// TestLoad_FromEnv 测试从环境变量加载配置
func TestLoad_FromEnv(t *testing.T) {
	// 设置环境变量
	os.Setenv("APP_SERVER_HOST", "127.0.0.1")
	os.Setenv("APP_SERVER_PORT", "9090")
	os.Setenv("APP_DB_TYPE", "sqlite3")
	os.Setenv("APP_DB_HOST", "db.example.com")
	os.Setenv("APP_DB_PORT", "5433")
	os.Setenv("APP_REDIS_ADDR", "redis.example.com:6379")
	os.Setenv("APP_LOG_LEVEL", "debug")
	defer func() {
		// 清理环境变量
		os.Unsetenv("APP_SERVER_HOST")
		os.Unsetenv("APP_SERVER_PORT")
		os.Unsetenv("APP_DB_TYPE")
		os.Unsetenv("APP_DB_HOST")
		os.Unsetenv("APP_DB_PORT")
		os.Unsetenv("APP_REDIS_ADDR")
		os.Unsetenv("APP_LOG_LEVEL")
	}()

	cfg, err := Load("")
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// 验证环境变量覆盖了默认值
	assert.Equal(t, "127.0.0.1", cfg.Server.Host)
	assert.Equal(t, 9090, cfg.Server.Port)
	assert.Equal(t, "sqlite3", cfg.Database.Type)
	assert.Equal(t, "db.example.com", cfg.Database.Host)
	assert.Equal(t, 5433, cfg.Database.Port)
	assert.Equal(t, "redis.example.com:6379", cfg.Redis.Addr)
	assert.Equal(t, "debug", cfg.Logging.Level)
}

// TestLoadConfig 测试 LoadConfig 便捷函数
func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, 8080, cfg.Server.Port)
}

// TestLoadFromEnv 测试 LoadFromEnv 便捷函数
func TestLoadFromEnv(t *testing.T) {
	os.Setenv("APP_SERVER_PORT", "7070")
	defer os.Unsetenv("APP_SERVER_PORT")

	cfg, err := LoadFromEnv()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, 7070, cfg.Server.Port)
}

// TestLoadFromFile 测试从指定文件加载配置
func TestLoadFromFile(t *testing.T) {
	// 创建临时配置文件
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// 写入测试配置
	configContent := `
server:
  host: "192.168.1.1"
  port: 3000
database:
  type: "postgres"
  host: "pg.example.com"
`
	_, err = tmpFile.WriteString(configContent)
	require.NoError(t, err)
	tmpFile.Close()

	// 从文件加载配置
	cfg, err := LoadFromFile(tmpFile.Name())
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// 验证配置
	assert.Equal(t, "192.168.1.1", cfg.Server.Host)
	assert.Equal(t, 3000, cfg.Server.Port)
	assert.Equal(t, "postgres", cfg.Database.Type)
	assert.Equal(t, "pg.example.com", cfg.Database.Host)
}

// TestSetDefaults 测试默认值设置
func TestSetDefaults(t *testing.T) {
	cfg := &Config{}
	setDefaults(cfg)

	// 验证 Server 默认值
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 8080, cfg.Server.Port)

	// 验证 Database 默认值
	assert.Equal(t, "postgres", cfg.Database.Type)
	assert.Equal(t, "localhost", cfg.Database.Host)

	// 验证已有值不会被覆盖
	cfg2 := &Config{
		Server: ServerConfig{
			Host: "custom-host",
			Port: 9999,
		},
	}
	setDefaults(cfg2)
	assert.Equal(t, "custom-host", cfg2.Server.Host)
	assert.Equal(t, 9999, cfg2.Server.Port)
}

// TestConfig_Validate 测试配置结构
func TestConfig_Validate(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		Database: DatabaseConfig{
			Type:     "postgres",
			Host:     "localhost",
			Port:     5432,
			User:     "admin",
			Password: "secret",
			Database: "myapp",
		},
		Redis: RedisConfig{
			Addr: "localhost:6379",
		},
	}

	// 验证配置结构正确
	assert.NotEmpty(t, cfg.Server.Host)
	assert.Greater(t, cfg.Server.Port, 0)
	assert.NotEmpty(t, cfg.Database.Type)
}

// TestSQLite3_DefaultDSN 测试 SQLite3 默认 DSN
func TestSQLite3_DefaultDSN(t *testing.T) {
	cfg := &Config{
		Database: DatabaseConfig{
			Type: "sqlite3",
		},
	}
	setDefaults(cfg)

	assert.Equal(t, "sqlite3", cfg.Database.Type)
	assert.Equal(t, "file:app.db?cache=shared&mode=rwc", cfg.Database.DSN)
}

// TestBindEnvVars 测试环境变量绑定（间接测试）
func TestBindEnvVars(t *testing.T) {
	// 设置测试环境变量
	testCases := []struct {
		name     string
		envKey   string
		envValue string
		verify   func(*Config) interface{}
		expected interface{}
	}{
		{
			name:     "server host",
			envKey:   "APP_SERVER_HOST",
			envValue: "0.0.0.0",
			verify:   func(c *Config) interface{} { return c.Server.Host },
			expected: "0.0.0.0",
		},
		{
			name:     "database type",
			envKey:   "APP_DB_TYPE",
			envValue: "mysql",
			verify:   func(c *Config) interface{} { return c.Database.Type },
			expected: "mysql",
		},
		{
			name:     "log level",
			envKey:   "APP_LOG_LEVEL",
			envValue: "error",
			verify:   func(c *Config) interface{} { return c.Logging.Level },
			expected: "error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv(tc.envKey, tc.envValue)
			defer os.Unsetenv(tc.envKey)

			cfg, err := Load("")
			require.NoError(t, err)
			assert.Equal(t, tc.expected, tc.verify(cfg))
		})
	}
}

// TestObservability_Defaults 测试可观测性默认值
func TestObservability_Defaults(t *testing.T) {
	cfg, err := Load("")
	require.NoError(t, err)

	// 验证 Observability 默认值
	assert.Equal(t, "localhost:4317", cfg.Observability.OTLP.Endpoint)
	assert.Equal(t, "app", cfg.Observability.OTLP.ServiceName)
	assert.Equal(t, "5s", cfg.Observability.System.CollectInterval)
	assert.Equal(t, 90.0, cfg.Observability.System.HealthThresholds.MaxMemoryUsage)
	assert.Equal(t, 95.0, cfg.Observability.System.HealthThresholds.MaxCPUUsage)
	assert.Equal(t, 10000, cfg.Observability.System.HealthThresholds.MaxGoroutines)
}

// TestRotationConfig 测试日志轮转配置
func TestRotationConfig(t *testing.T) {
	cfg := &Config{
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
			Rotation: RotationConfig{
				MaxSize:    100,
				MaxBackups: 5,
				MaxAge:     30,
				Compress:   true,
			},
		},
	}

	assert.Equal(t, 100, cfg.Logging.Rotation.MaxSize)
	assert.Equal(t, 5, cfg.Logging.Rotation.MaxBackups)
	assert.Equal(t, 30, cfg.Logging.Rotation.MaxAge)
	assert.True(t, cfg.Logging.Rotation.Compress)
}

// TestAlertRuleConfig 测试告警规则配置
func TestAlertRuleConfig(t *testing.T) {
	cfg := &Config{
		Observability: ObservabilityConfig{
			System: SystemMonitoringConfig{
				Alerts: []AlertRuleConfig{
					{
						ID:         "high-memory",
						Name:       "High Memory Usage",
						MetricName: "memory_usage",
						Condition:  "gt",
						Threshold:  90.0,
						Level:      "critical",
						Enabled:    true,
						Duration:   "5m",
						Cooldown:   "10m",
					},
				},
			},
		},
	}

	require.Len(t, cfg.Observability.System.Alerts, 1)
	alert := cfg.Observability.System.Alerts[0]
	assert.Equal(t, "high-memory", alert.ID)
	assert.Equal(t, "High Memory Usage", alert.Name)
	assert.Equal(t, "memory_usage", alert.MetricName)
	assert.Equal(t, "gt", alert.Condition)
	assert.Equal(t, 90.0, alert.Threshold)
	assert.Equal(t, "critical", alert.Level)
	assert.True(t, alert.Enabled)
}

// TestRateLimitConfig 测试限流器配置
func TestRateLimitConfig(t *testing.T) {
	cfg := &Config{
		Observability: ObservabilityConfig{
			System: SystemMonitoringConfig{
				RateLimit: RateLimitConfig{
					Enabled: true,
					Limit:   1000,
					Window:  "1s",
				},
			},
		},
	}

	assert.True(t, cfg.Observability.System.RateLimit.Enabled)
	assert.Equal(t, int64(1000), cfg.Observability.System.RateLimit.Limit)
	assert.Equal(t, "1s", cfg.Observability.System.RateLimit.Window)
}

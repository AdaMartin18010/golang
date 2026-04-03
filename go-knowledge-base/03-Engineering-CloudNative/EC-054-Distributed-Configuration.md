# EC-054: Distributed Configuration Pattern

> **Dimension**: Engineering-CloudNative
> **Level**: S (>15KB)
> **Tags**: #configuration #consul #etcd #vault #spring-cloud-config #externalized-configuration
> **Authoritative Sources**:
>
> - [12 Factor App - Config](https://12factor.net/config) - Adam Wiggins (2011)
> - [Spring Cloud Config](https://cloud.spring.io/spring-cloud-config/reference/html/) - Spring (2024)
> - [Consul Documentation](https://www.consul.io/docs) - HashiCorp (2024)
> - [etcd Documentation](https://etcd.io/docs/) - CNCF (2024)

---

## 1. Problem Formalization

### 1.1 System Context and Constraints

**Definition 1.1 (Configuration Domain)**
Let $\mathcal{C}$ be the configuration space for application $A$ deployed across environments $\mathcal{E} = \{dev, staging, prod\}$ with instances $\mathcal{I} = \{i_1, i_2, ..., i_n\}$.

**Configuration Types:**

| Type | Sensitivity | Volatility | Example |
|------|-------------|------------|---------|
| **Static** | Low | Rarely changes | Feature flags |
| **Dynamic** | Medium | Changes at runtime | Rate limits |
| **Secret** | High | Rotates regularly | API keys |

**System Constraints:**

| Constraint | Formal Definition | Impact |
|------------|-------------------|--------|
| **Consistency** | $\forall i, j \in \mathcal{I}: config(i) = config(j)$ | Requires atomic distribution |
| **Latency Bound** | $T_{propagate} < T_{acceptable}$ | Changes must propagate quickly |
| **Security** | $\forall c \in secrets: encrypted(c)$ | Secrets must be protected |
| **Versioning** | $\forall t: \exists c_t: retrieve(t) = c_t$ | Audit and rollback required |

### 1.2 Problem Statement

**Problem 1.1 (Configuration Distribution)**
Given configuration $c$ and target instances $\mathcal{I}$, ensure:

$$\forall i \in \mathcal{I}: apply(i, c) \land verify(i, c) \land consistent(c)$$

**Key Challenges:**

1. **Secret Management**: Storing and rotating sensitive configuration
2. **Dynamic Updates**: Applying changes without restart
3. **Multi-environment**: Managing variations across environments
4. **Consistency**: Ensuring all instances have same configuration
5. **Audit Trail**: Tracking configuration changes

---

## 2. Solution Architecture

### 2.1 Configuration Store Types

| Store | Use Case | Consistency | Latency |
|-------|----------|-------------|---------|
| **etcd** | Kubernetes, service discovery | Strong | Low |
| **Consul** | Service mesh, KV store | Eventual | Low |
| **Vault** | Secrets management | Strong | Medium |
| **Git** | Version-controlled config | Eventual | High |
| **S3/ConfigMap** | Static configuration | Eventual | Medium |

---

## 3. Visual Representations

### 3.1 Distributed Configuration Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    DISTRIBUTED CONFIGURATION SYSTEM                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  CONFIGURATION SOURCES                                                      │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                      │   │
│  │  ┌───────────────┐  ┌───────────────┐  ┌───────────────────────┐   │   │
│  │  │     Git       │  │     etcd      │  │       Vault           │   │   │
│  │  │   Repository  │  │    Cluster    │  │    (Secrets)          │   │   │
│  │  │               │  │               │  │                       │   │   │
│  │  │ • app.yml     │  │ • /config/app │  │ • database/password   │   │   │
│  │  │ • prod/       │  │ • /services/  │  │ • api/key             │   │   │
│  │  │   database.yml│  │   web/replicas│  │ • tls/cert            │   │   │
│  │  └───────┬───────┘  └───────┬───────┘  └───────────┬───────────┘   │   │
│  │          │                  │                      │               │   │
│  └──────────┼──────────────────┼──────────────────────┼───────────────┘   │
│             │                  │                      │                   │
│             └──────────────────┼──────────────────────┘                   │
│                                │                                          │
│                                ▼                                          │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     CONFIGURATION SERVER                             │   │
│  │                                                                      │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                  │   │
│  │  │   Git Pull  │  │  etcd Watch │  │ Vault Auth  │                  │   │
│  │  │             │  │             │  │             │                  │   │
│  │  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘                  │   │
│  │         │                │                │                         │   │
│  │         └────────────────┼────────────────┘                         │   │
│  │                          ▼                                         │   │
│  │                 ┌─────────────────┐                                 │   │
│  │                 │  Config Cache   │                                 │   │
│  │                 │                 │                                 │   │
│  │                 │ • Environment   │                                 │   │
│  │                 │   resolution    │                                 │   │
│  │                 │ • Secret        │                                 │   │
│  │                 │   injection     │                                 │   │
│  │                 └────────┬────────┘                                 │   │
│  │                          │                                         │   │
│  │  ┌───────────────────────┼───────────────────────┐                 │   │
│  │  │                       │                       │                 │   │
│  │  ▼                       ▼                       ▼                 │   │
│  │  REST API          gRPC Stream            WebSocket               │   │
│  │  (Pull)            (Push updates)         (Bidirectional)         │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│              │                   │                   │                      │
│              └───────────────────┼───────────────────┘                      │
│                                  │                                         │
│                                  ▼                                         │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     APPLICATION INSTANCES                            │   │
│  │                                                                      │   │
│  │  Instance 1              Instance 2              Instance N          │   │
│  │  ┌──────────────┐       ┌──────────────┐       ┌──────────────┐     │   │
│  │  │ Config       │       │ Config       │       │ Config       │     │   │
│  │  │ Client       │       │ Client       │       │ Client       │     │   │
│  │  │              │       │              │       │              │     │   │
│  │  │ • Local cache│       │ • Local cache│       │ • Local cache│     │   │
│  │  │ • Watch for  │       │ • Watch for  │       │ • Watch for  │     │   │
│  │  │   changes    │       │   changes    │       │   changes    │     │   │
│  │  │ • Fallback   │       │ • Fallback   │       │ • Fallback   │     │   │
│  │  │   values     │       │   values     │       │   values     │     │   │
│  │  └──────┬───────┘       └──────┬───────┘       └──────┬───────┘     │   │
│  │         │                      │                      │              │   │
│  │         ▼                      ▼                      ▼              │   │
│  │  Application Code         Application Code      Application Code     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Configuration Update Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    CONFIGURATION UPDATE PROPAGATION                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Step 1: Developer commits configuration change to Git                      │
│                                                                             │
│  Git ──► [Commit: Update database timeout to 30s]                           │
│                                                                             │
│  Step 2: CI/CD pipeline detects change and validates                        │
│                                                                             │
│  ┌─────────────┐     ┌─────────────┐     ┌─────────────┐                   │
│  │   Build     │────►│    Test     │────►│   Deploy    │                   │
│  └─────────────┘     └─────────────┘     └──────┬──────┘                   │
│                                                  │                         │
│  Step 3: Configuration Server pulls new config  │                         │
│                                                                             │
│  ┌───────────────────────────────────────────────┐                         │
│  │        Configuration Server                   │                         │
│  │  ┌─────────────┐      ┌───────────────────┐   │                         │
│  │  │  Git Pull   │─────►│  Config Updated   │   │                         │
│  │  │  (webhook)  │      │  Event Generated  │   │                         │
│  │  └─────────────┘      └─────────┬─────────┘   │                         │
│  └─────────────────────────────────┼─────────────┘                         │
│                                    │                                       │
│  Step 4: Event propagated to all connected clients                         │
│                                                                             │
│         ┌──────────────────────────┼──────────────────────────┐            │
│         │                          │                          │            │
│         ▼                          ▼                          ▼            │
│  ┌──────────────┐           ┌──────────────┐           ┌──────────────┐   │
│  │  Instance 1  │           │  Instance 2  │           │  Instance 3  │   │
│  │              │           │              │           │              │   │
│  │  gRPC Stream │           │  gRPC Stream │           │  gRPC Stream │   │
│  │  Event:      │           │  Event:      │           │  Event:      │   │
│  │  DB_TIMEOUT  │           │  DB_TIMEOUT  │           │  DB_TIMEOUT  │   │
│  │  = 30s       │           │  = 30s       │           │  = 30s       │   │
│  └──────┬───────┘           └──────┬───────┘           └──────┬───────┘   │
│         │                          │                          │            │
│  Step 5: Each instance applies update (hot reload or restart)              │
│         │                          │                          │            │
│         ▼                          ▼                          ▼            │
│  ┌──────────────┐           ┌──────────────┐           ┌──────────────┐   │
│  │  Update      │           │  Update      │           │  Update      │   │
│  │  Connection  │           │  Connection  │           │  Connection  │   │
│  │  Pool        │           │  Pool        │           │  Pool        │   │
│  └──────────────┘           └──────────────┘           └──────────────┘   │
│                                                                             │
│  Step 6: Acknowledgment and verification                                    │
│                                                                             │
│  Instance 1 ──ACK──► Config Server ◄──ACK── Instance 2                      │
│                                      ◄──ACK── Instance 3                    │
│                                                                             │
│  All instances confirmed updated!                                           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Production Go Implementation

### 4.1 Configuration Client with etcd

```go
package config

import (
 "context"
 "encoding/json"
 "fmt"
 "sync"
 "time"

 clientv3 "go.etcd.io/etcd/client/v3"
 "go.uber.org/zap"
)

// Source represents a configuration source
type Source interface {
 Get(ctx context.Context, key string) (string, error)
 Watch(ctx context.Context, key string) (<-chan Event, error)
 Close() error
}

// Event represents a configuration change event
type Event struct {
 Type  EventType
 Key   string
 Value string
}

type EventType int

const (
 EventTypePut EventType = iota
 EventTypeDelete
)

// EtcdSource implements Source for etcd
type EtcdSource struct {
 client *clientv3.Client
 prefix string
 logger *zap.Logger
}

// NewEtcdSource creates a new etcd configuration source
func NewEtcdSource(endpoints []string, prefix string, logger *zap.Logger) (*EtcdSource, error) {
 client, err := clientv3.New(clientv3.Config{
  Endpoints:   endpoints,
  DialTimeout: 5 * time.Second,
 })
 if err != nil {
  return nil, fmt.Errorf("failed to create etcd client: %w", err)
 }

 return &EtcdSource{
  client: client,
  prefix: prefix,
  logger: logger,
 }, nil
}

// Get retrieves a configuration value
func (s *EtcdSource) Get(ctx context.Context, key string) (string, error) {
 fullKey := s.prefix + key
 resp, err := s.client.Get(ctx, fullKey)
 if err != nil {
  return "", err
 }

 if len(resp.Kvs) == 0 {
  return "", fmt.Errorf("key not found: %s", key)
 }

 return string(resp.Kvs[0].Value), nil
}

// Watch watches for configuration changes
func (s *EtcdSource) Watch(ctx context.Context, key string) (<-chan Event, error) {
 fullKey := s.prefix + key
 watchChan := s.client.Watch(ctx, fullKey)

 eventChan := make(chan Event)
 go func() {
  defer close(eventChan)
  for watchResp := range watchChan {
   for _, event := range watchResp.Events {
    var eventType EventType
    switch event.Type {
    case clientv3.EventTypePut:
     eventType = EventTypePut
    case clientv3.EventTypeDelete:
     eventType = EventTypeDelete
    }

    select {
    case eventChan <- Event{
     Type:  eventType,
     Key:   string(event.Kv.Key),
     Value: string(event.Kv.Value),
    }:
    case <-ctx.Done():
     return
    }
   }
  }
 }()

 return eventChan, nil
}

// Close closes the etcd connection
func (s *EtcdSource) Close() error {
 return s.client.Close()
}

// Manager manages application configuration
type Manager struct {
 source  Source
 cache   map[string]string
 mu      sync.RWMutex
 logger  *zap.Logger

 // Callbacks for configuration changes
 handlers map[string][]func(string)
}

// NewManager creates a new configuration manager
func NewManager(source Source, logger *zap.Logger) *Manager {
 return &Manager{
  source:   source,
  cache:    make(map[string]string),
  logger:   logger,
  handlers: make(map[string][]func(string)),
 }
}

// Get retrieves configuration value (from cache or source)
func (m *Manager) Get(ctx context.Context, key string) (string, error) {
 // Check cache first
 m.mu.RLock()
 if val, ok := m.cache[key]; ok {
  m.mu.RUnlock()
  return val, nil
 }
 m.mu.RUnlock()

 // Fetch from source
 val, err := m.source.Get(ctx, key)
 if err != nil {
  return "", err
 }

 // Update cache
 m.mu.Lock()
 m.cache[key] = val
 m.mu.Unlock()

 return val, nil
}

// GetInt retrieves configuration as integer
func (m *Manager) GetInt(ctx context.Context, key string) (int, error) {
 val, err := m.Get(ctx, key)
 if err != nil {
  return 0, err
 }

 var result int
 if err := json.Unmarshal([]byte(val), &result); err != nil {
  return 0, err
 }
 return result, nil
}

// GetBool retrieves configuration as boolean
func (m *Manager) GetBool(ctx context.Context, key string) (bool, error) {
 val, err := m.Get(ctx, key)
 if err != nil {
  return false, err
 }

 var result bool
 if err := json.Unmarshal([]byte(val), &result); err != nil {
  return false, err
 }
 return result, nil
}

// Watch starts watching a configuration key for changes
func (m *Manager) Watch(ctx context.Context, key string) error {
 eventChan, err := m.source.Watch(ctx, key)
 if err != nil {
  return err
 }

 go func() {
  for event := range eventChan {
   switch event.Type {
   case EventTypePut:
    m.mu.Lock()
    m.cache[key] = event.Value
    m.mu.Unlock()

    // Notify handlers
    m.mu.RLock()
    handlers := m.handlers[key]
    m.mu.RUnlock()

    for _, handler := range handlers {
     go handler(event.Value)
    }

    m.logger.Info("Configuration updated",
     zap.String("key", key),
     zap.String("value", event.Value))

   case EventTypeDelete:
    m.mu.Lock()
    delete(m.cache, key)
    m.mu.Unlock()

    m.logger.Info("Configuration deleted",
     zap.String("key", key))
   }
  }
 }()

 return nil
}

// OnChange registers a callback for configuration changes
func (m *Manager) OnChange(key string, handler func(string)) {
 m.mu.Lock()
 defer m.mu.Unlock()
 m.handlers[key] = append(m.handlers[key], handler)
}

// Close closes the configuration manager
func (m *Manager) Close() error {
 return m.source.Close()
}
```

### 4.2 Environment-Specific Configuration

```go
package config

import (
 "fmt"
 "os"
 "strings"
)

// Environment represents deployment environment
type Environment string

const (
 EnvDevelopment Environment = "development"
 EnvStaging     Environment = "staging"
 EnvProduction  Environment = "production"
)

// Config holds application configuration
type Config struct {
 Environment Environment `json:"environment"`

 Server struct {
  Host         string        `json:"host"`
  Port         int           `json:"port"`
  ReadTimeout  time.Duration `json:"read_timeout"`
  WriteTimeout time.Duration `json:"write_timeout"`
 } `json:"server"`

 Database struct {
  Host            string        `json:"host"`
  Port            int           `json:"port"`
  Name            string        `json:"name"`
  User            string        `json:"user"`
  Password        string        `json:"password"`
  MaxConnections  int           `json:"max_connections"`
  ConnectTimeout  time.Duration `json:"connect_timeout"`
 } `json:"database"`

 Cache struct {
  Host     string        `json:"host"`
  Port     int           `json:"port"`
  TTL      time.Duration `json:"ttl"`
 } `json:"cache"`

 Features struct {
  EnableNewUI      bool `json:"enable_new_ui"`
  EnableRateLimit  bool `json:"enable_rate_limit"`
  MaxRequestsPerSec int `json:"max_requests_per_sec"`
 } `json:"features"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
 cfg := &Config{}

 // Load environment
 env := Environment(getEnv("APP_ENV", "development"))
 cfg.Environment = env

 // Server configuration
 cfg.Server.Host = getEnv("SERVER_HOST", "0.0.0.0")
 cfg.Server.Port = getEnvInt("SERVER_PORT", 8080)
 cfg.Server.ReadTimeout = getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second)
 cfg.Server.WriteTimeout = getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second)

 // Database configuration
 cfg.Database.Host = getEnv("DB_HOST", "localhost")
 cfg.Database.Port = getEnvInt("DB_PORT", 5432)
 cfg.Database.Name = getEnv("DB_NAME", "app")
 cfg.Database.User = getEnv("DB_USER", "app")
 cfg.Database.Password = getEnv("DB_PASSWORD", "")
 cfg.Database.MaxConnections = getEnvInt("DB_MAX_CONNECTIONS", 10)
 cfg.Database.ConnectTimeout = getEnvDuration("DB_CONNECT_TIMEOUT", 10*time.Second)

 // Cache configuration
 cfg.Cache.Host = getEnv("CACHE_HOST", "localhost")
 cfg.Cache.Port = getEnvInt("CACHE_PORT", 6379)
 cfg.Cache.TTL = getEnvDuration("CACHE_TTL", 5*time.Minute)

 // Feature flags
 cfg.Features.EnableNewUI = getEnvBool("FEATURE_NEW_UI", false)
 cfg.Features.EnableRateLimit = getEnvBool("FEATURE_RATE_LIMIT", true)
 cfg.Features.MaxRequestsPerSec = getEnvInt("MAX_REQUESTS_PER_SEC", 100)

 // Environment-specific overrides
 if err := applyEnvironmentOverrides(cfg); err != nil {
  return nil, err
 }

 return cfg, nil
}

func applyEnvironmentOverrides(cfg *Config) error {
 switch cfg.Environment {
 case EnvDevelopment:
  cfg.Features.EnableNewUI = true

 case EnvStaging:
  cfg.Features.EnableRateLimit = true
  cfg.Features.MaxRequestsPerSec = 1000

 case EnvProduction:
  // Strict settings for production
  if cfg.Database.Password == "" {
   return fmt.Errorf("database password required in production")
  }
  cfg.Server.ReadTimeout = 60 * time.Second
  cfg.Server.WriteTimeout = 60 * time.Second
  cfg.Features.EnableRateLimit = true
 }

 return nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
 if value := os.Getenv(key); value != "" {
  return value
 }
 return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
 if value := os.Getenv(key); value != "" {
  var result int
  fmt.Sscanf(value, "%d", &result)
  return result
 }
 return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
 if value := os.Getenv(key); value != "" {
  return strings.ToLower(value) == "true" || value == "1"
 }
 return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
 if value := os.Getenv(key); value != "" {
  if d, err := time.ParseDuration(value); err == nil {
   return d
  }
 }
 return defaultValue
}
```

---

## 5. Failure Scenarios and Mitigations

### 5.1 Configuration Failure Taxonomy

| Scenario | Impact | Detection | Mitigation |
|----------|--------|-----------|------------|
| **Config Server Down** | Cannot update config | Connection timeout | Local cache + Fallback values |
| **Invalid Config** | Application crash | Validation error | Schema validation + Rejection |
| **Secret Leak** | Security breach | Audit scan | Encryption + Access control |
| **Propagation Delay** | Inconsistent state | Version mismatch | Eventual consistency + Retry |
| **Race Condition** | Wrong config applied | Version conflict | Optimistic locking |

---

## 6. Semantic Trade-off Analysis

| Aspect | Push Model | Pull Model | Hybrid |
|--------|------------|------------|--------|
| **Latency** | Low | Higher | Configurable |
| **Scalability** | Moderate | High | High |
| **Complexity** | High | Low | Medium |
| **Reliability** | Moderate | High | High |

---

## 7. References

1. Wiggins, A. (2011). *The Twelve-Factor App*. 12factor.net.
2. HashiCorp. (2024). *Consul Documentation*. consul.io.
3. CNCF. (2024). *etcd Documentation*. etcd.io.
4. Spring Team. (2024). *Spring Cloud Config*. spring.io.

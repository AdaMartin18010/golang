package wire

import (
	"github.com/google/wire"

	appuser "github.com/yourusername/golang/internal/application/user"
	"github.com/yourusername/golang/internal/config"
	"github.com/yourusername/golang/internal/infrastructure/database/ent"
	entrepo "github.com/yourusername/golang/internal/infrastructure/database/ent/repository"
	chirouter "github.com/yourusername/golang/internal/interfaces/http/chi"

	// 可观测性
	"github.com/yourusername/golang/pkg/observability"
	"github.com/yourusername/golang/pkg/observability/otlp"
	"github.com/yourusername/golang/pkg/observability/system"

	// 安全
	"github.com/yourusername/golang/pkg/security/jwt"
	"github.com/yourusername/golang/pkg/security/rbac"
)

// ObservabilityProviderSet 可观测性 Provider 集合
// 包含 OTLP、系统监控、eBPF 等
var ObservabilityProviderSet = wire.NewSet(
	NewOTLPIntegration,
	NewSystemMonitor,
	NewPlatformMonitor,
	// NewEBPFCollector, // 可选，需要 Linux
)

// SecurityProviderSet 安全 Provider 集合
// 包含 JWT、RBAC、OAuth2 等
var SecurityProviderSet = wire.NewSet(
	NewJWTTokenManager,
	NewRBACSystem,
	// NewOAuth2Provider, // 可选，需要配置
)

// DatabaseProviderSet 数据库 Provider 集合
var DatabaseProviderSet = wire.NewSet(
	NewEntClient,
	entrepo.NewBaseRepository,
	// 具体仓储需要根据业务定义
	// NewUserRepository,
)

// ApplicationProviderSet 应用层 Provider 集合
var ApplicationProviderSet = wire.NewSet(
	appuser.NewService,
)

// InterfaceProviderSet 接口层 Provider 集合
var InterfaceProviderSet = wire.NewSet(
	chirouter.NewRouter,
)

// AllProviderSet 所有 Provider 集合
var AllProviderSet = wire.NewSet(
	ObservabilityProviderSet,
	SecurityProviderSet,
	DatabaseProviderSet,
	ApplicationProviderSet,
	InterfaceProviderSet,
)

// Provider 函数实现

// NewOTLPIntegration 创建 OTLP 集成
func NewOTLPIntegration(cfg *config.Config) (*otlp.EnhancedOTLP, error) {
	return otlp.NewEnhancedOTLP(otlp.Config{
		ServiceName:    cfg.ServiceName,
		ServiceVersion: cfg.ServiceVersion,
		Endpoint:       cfg.Observability.OTLP.Endpoint,
		Insecure:       cfg.Observability.OTLP.Insecure,
		SampleRate:     cfg.Observability.OTLP.SampleRate,
	})
}

// NewSystemMonitor 创建系统监控器
func NewSystemMonitor(otlp *otlp.EnhancedOTLP) (*system.Monitor, error) {
	return system.NewMonitor(system.Config{
		Meter:           otlp.GetMeter(),
		CollectInterval: 5 * time.Second,
	})
}

// NewPlatformMonitor 创建平台监控器
func NewPlatformMonitor(otlp *otlp.EnhancedOTLP) (*system.PlatformMonitor, error) {
	return system.NewPlatformMonitor(otlp.GetMeter())
}

// NewJWTTokenManager 创建 JWT Token Manager
func NewJWTTokenManager(cfg *config.Config) (*jwt.TokenManager, error) {
	return jwt.NewTokenManager(jwt.Config{
		Issuer:          cfg.ServiceName,
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		SigningMethod:   "RS256",
	})
}

// NewRBACSystem 创建 RBAC 系统
func NewRBACSystem() (*rbac.RBAC, error) {
	rbacSystem := rbac.NewRBAC()
	// 初始化默认角色
	if err := rbacSystem.InitializeDefaultRoles(); err != nil {
		return nil, err
	}
	return rbacSystem, nil
}

// NewEntClient 创建 Ent 客户端
func NewEntClient(cfg *config.Config) (*ent.Client, error) {
	// 根据配置创建数据库连接
	// 实际实现需要根据配置选择不同的数据库驱动
	return ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
}

// 需要的导入
import "time"

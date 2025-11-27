//go:build wireinject
// +build wireinject

// Package main 展示如何使用 Wire 进行依赖注入
//
// 本示例展示：
// 1. 如何使用 Wire 进行依赖注入
// 2. 如何初始化框架的各种能力
// 3. 如何组织依赖关系
package main

import (
	"context"
	"time"

	"github.com/google/wire"
	"github.com/yourusername/golang/pkg/control"
	"github.com/yourusername/golang/pkg/database"
	"github.com/yourusername/golang/pkg/observability/otlp"
	"github.com/yourusername/golang/pkg/sampling"
	"github.com/yourusername/golang/pkg/tracing"
)

// App 应用程序结构
type App struct {
	DB                  database.Database
	Tracer              tracing.Tracer
	OTLP                *otlp.EnhancedOTLP
	FeatureController   control.Controller
	RateController      *control.RateController
	CircuitController   *control.CircuitController
}

// NewApp 创建应用程序实例
func NewApp(
	db database.Database,
	tracer tracing.Tracer,
	otlpClient *otlp.EnhancedOTLP,
	featureController control.Controller,
	rateController *control.RateController,
	circuitController *control.CircuitController,
) *App {
	return &App{
		DB:                db,
		Tracer:            tracer,
		OTLP:              otlpClient,
		FeatureController: featureController,
		RateController:    rateController,
		CircuitController: circuitController,
	}
}

// InitializeApp 初始化应用程序（Wire 生成）
func InitializeApp(ctx context.Context) (*App, error) {
	wire.Build(
		// 数据库
		provideDatabase,

		// 采样器
		provideSampler,

		// OTLP
		provideOTLP,

		// 追踪器
		provideTracer,

		// 控制器
		provideFeatureController,
		provideRateController,
		provideCircuitController,

		// 应用程序
		NewApp,
	)
	return nil, nil
}

// provideDatabase 提供数据库实例
func provideDatabase() (database.Database, error) {
	return database.NewDatabase(database.Config{
		Driver:       database.DriverSQLite3,
		DSN:          "file:example.db?cache=shared&mode=memory",
		MaxOpenConns: 25,
		MaxIdleConns: 5,
	})
}

// provideSampler 提供采样器实例
func provideSampler() (sampling.Sampler, error) {
	return sampling.NewProbabilisticSampler(0.5)
}

// provideOTLP 提供 OTLP 客户端实例
func provideOTLP(sampler sampling.Sampler) (*otlp.EnhancedOTLP, error) {
	return otlp.NewEnhancedOTLP(otlp.Config{
		ServiceName:    "example-service",
		ServiceVersion: "v1.0.0",
		Endpoint:       "localhost:4317",
		Insecure:       true,
		Sampler:        sampler,
	})
}

// provideTracer 提供追踪器实例
func provideTracer() *tracing.Tracer {
	return tracing.NewTracer("example-service")
}

// provideFeatureController 提供功能控制器实例
func provideFeatureController() control.Controller {
	controller := control.NewFeatureController()
	// 注意：需要类型断言才能调用 Register 方法
	if fc, ok := controller.(*control.FeatureController); ok {
		fc.Register("experimental-feature", "Experimental feature", true, nil)
	}
	return controller
}

// provideRateController 提供速率控制器实例
func provideRateController() *control.RateController {
	controller := control.NewRateController()
	controller.SetRateLimit("api-calls", 100.0, 1*time.Second)
	return controller
}

// provideCircuitController 提供熔断器控制器实例
func provideCircuitController() *control.CircuitController {
	controller := control.NewCircuitController()
	controller.RegisterCircuit("external-api", 10, 5, 30*time.Second)
	return controller
}

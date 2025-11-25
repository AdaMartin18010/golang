//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	appuser "github.com/yourusername/golang/internal/application/user"
	"github.com/yourusername/golang/internal/config"
	"github.com/yourusername/golang/internal/infrastructure/database/postgres"
	chirouter "github.com/yourusername/golang/internal/interfaces/http/chi"
)

// InitializeApp 初始化应用（Wire 自动生成）
func InitializeApp(cfg *config.Config) (*App, error) {
	wire.Build(
		// Infrastructure
		postgres.NewUserRepository,

		// Application
		appuser.NewService,

		// Interfaces
		chirouter.NewRouter,

		// App
		NewApp,
	)
	return &App{}, nil
}

// App 应用结构
type App struct {
	Router *chirouter.Router
}

// NewApp 创建应用
func NewApp(router *chirouter.Router) *App {
	return &App{Router: router}
}

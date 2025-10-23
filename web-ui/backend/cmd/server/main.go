package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/your-org/go-formal-verification/web-ui/internal/api"
	"github.com/your-org/go-formal-verification/web-ui/internal/ws"
)

const (
	defaultPort = "8080"
	apiVersion  = "v1"
)

func main() {
	// 创建Gin路由器
	router := gin.Default()

	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// 健康检查端点
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": apiVersion,
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// API路由组
	v1 := router.Group("/api/v1")
	{
		// 分析API
		analysis := v1.Group("/analysis")
		{
			analysis.POST("/cfg", api.AnalyzeCFG)
			analysis.POST("/concurrency", api.AnalyzeConcurrency)
			analysis.POST("/types", api.AnalyzeTypes)
			analysis.GET("/history", api.GetAnalysisHistory)
		}

		// 并发模式API
		patterns := v1.Group("/patterns")
		{
			patterns.GET("", api.ListPatterns)
			patterns.GET("/:name", api.GetPattern)
			patterns.POST("/generate", api.GeneratePattern)
		}

		// 项目管理API
		projects := v1.Group("/projects")
		{
			projects.GET("", api.ListProjects)
			projects.POST("", api.CreateProject)
			projects.GET("/:id", api.GetProject)
			projects.DELETE("/:id", api.DeleteProject)
		}
	}

	// WebSocket端点
	router.GET("/ws", ws.HandleWebSocket)

	// 静态文件服务（生产环境）
	router.Static("/assets", "./public/assets")
	router.StaticFile("/", "./public/index.html")

	// 启动服务器
	addr := fmt.Sprintf(":%s", defaultPort)
	log.Printf("🚀 Web UI Server starting on %s", addr)
	log.Printf("📊 API Version: %s", apiVersion)
	log.Printf("🔗 WebSocket endpoint: ws://localhost%s/ws", addr)
	log.Printf("💚 Health check: http://localhost%s/health", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

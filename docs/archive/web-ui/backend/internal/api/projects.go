package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Project 项目
type Project struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Path        string       `json:"path"`
	CreatedAt   string       `json:"createdAt"`
	UpdatedAt   string       `json:"updatedAt"`
	Stats       ProjectStats `json:"stats"`
}

// ProjectStats 项目统计
type ProjectStats struct {
	Files        int    `json:"files"`
	TotalLines   int    `json:"totalLines"`
	Goroutines   int    `json:"goroutines"`
	Channels     int    `json:"channels"`
	LastAnalysis string `json:"lastAnalysis"`
}

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Path        string `json:"path" binding:"required"`
}

// ListProjects 列出所有项目
func ListProjects(c *gin.Context) {
	// TODO: 从数据库读取项目列表
	projects := []Project{
		{
			ID:          "1",
			Name:        "Web Crawler",
			Description: "Concurrent web crawler example",
			Path:        "/path/to/web-crawler",
			CreatedAt:   "2025-10-20T10:00:00Z",
			UpdatedAt:   "2025-10-23T12:00:00Z",
			Stats: ProjectStats{
				Files:        5,
				TotalLines:   350,
				Goroutines:   12,
				Channels:     8,
				LastAnalysis: "2025-10-23T12:00:00Z",
			},
		},
		{
			ID:          "2",
			Name:        "Microservice Gateway",
			Description: "API gateway with load balancing",
			Path:        "/path/to/gateway",
			CreatedAt:   "2025-10-22T14:00:00Z",
			UpdatedAt:   "2025-10-23T11:30:00Z",
			Stats: ProjectStats{
				Files:        15,
				TotalLines:   1200,
				Goroutines:   50,
				Channels:     25,
				LastAnalysis: "2025-10-23T11:30:00Z",
			},
		},
	}

	c.JSON(http.StatusOK, AnalysisResponse{
		Success: true,
		Data:    projects,
		Time:    getCurrentTime(),
	})
}

// CreateProject 创建新项目
func CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, AnalysisResponse{
			Success: false,
			Error:   err.Error(),
			Time:    getCurrentTime(),
		})
		return
	}

	// TODO: 验证路径并保存到数据库
	project := Project{
		ID:          "3", // TODO: 生成UUID
		Name:        req.Name,
		Description: req.Description,
		Path:        req.Path,
		CreatedAt:   getCurrentTime(),
		UpdatedAt:   getCurrentTime(),
		Stats: ProjectStats{
			Files:      0,
			TotalLines: 0,
		},
	}

	c.JSON(http.StatusCreated, AnalysisResponse{
		Success: true,
		Data:    project,
		Time:    getCurrentTime(),
	})
}

// GetProject 获取项目详情
func GetProject(c *gin.Context) {
	id := c.Param("id")

	// TODO: 从数据库读取项目
	if id == "1" {
		project := Project{
			ID:          "1",
			Name:        "Web Crawler",
			Description: "Concurrent web crawler example",
			Path:        "/path/to/web-crawler",
			CreatedAt:   "2025-10-20T10:00:00Z",
			UpdatedAt:   "2025-10-23T12:00:00Z",
			Stats: ProjectStats{
				Files:        5,
				TotalLines:   350,
				Goroutines:   12,
				Channels:     8,
				LastAnalysis: "2025-10-23T12:00:00Z",
			},
		}

		c.JSON(http.StatusOK, AnalysisResponse{
			Success: true,
			Data:    project,
			Time:    getCurrentTime(),
		})
		return
	}

	c.JSON(http.StatusNotFound, AnalysisResponse{
		Success: false,
		Error:   "Project not found",
		Time:    getCurrentTime(),
	})
}

// DeleteProject 删除项目
func DeleteProject(c *gin.Context) {
	id := c.Param("id")

	// TODO: 从数据库删除项目
	c.JSON(http.StatusOK, AnalysisResponse{
		Success: true,
		Data:    gin.H{"id": id, "deleted": true},
		Time:    getCurrentTime(),
	})
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"delivery/http"
	"repository/implementations"
	"usecase"
)

func main() {
	fmt.Println("🚀 Starting Clean Architecture Demo...")

	// 1. 创建仓储层实现
	userRepo := implementations.NewMemoryUserRepository()

	// 2. 创建用例层
	userService := usecase.NewUserService(userRepo)

	// 3. 创建HTTP处理器
	userHandler := http.NewUserHandler(userService)

	// 4. 设置路由
	setupRoutes(userHandler)

	// 5. 启动服务器
	fmt.Println("📡 Server starting on :8080")
	fmt.Println("🌐 Available endpoints:")
	fmt.Println("  POST   /users          - Create user")
	fmt.Println("  GET    /users          - Get all users")
	fmt.Println("  GET    /users?id=xxx   - Get user by ID")
	fmt.Println("  PUT    /users?id=xxx   - Update user")
	fmt.Println("  DELETE /users?id=xxx   - Delete user")
	fmt.Println("  GET    /users/range?min_age=18&max_age=30 - Get users by age range")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupRoutes(handler *http.UserHandler) {
	// 用户管理路由
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateUser(w, r)
		case http.MethodGet:
			if r.URL.Query().Get("id") != "" {
				handler.GetUser(w, r)
			} else {
				handler.GetAllUsers(w, r)
			}
		case http.MethodPut:
			handler.UpdateUser(w, r)
		case http.MethodDelete:
			handler.DeleteUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// 年龄范围查询路由
	http.HandleFunc("/users/range", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetUsersByAgeRange(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// 健康检查
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","architecture":"clean"}`)
	})
}

// 演示函数：展示Clean Architecture的使用
func demonstrateCleanArchitecture() {
	fmt.Println("\n🎯 Clean Architecture Demo:")
	fmt.Println("==========================")

	// 创建仓储
	repo := implementations.NewMemoryUserRepository()

	// 创建服务
	service := usecase.NewUserService(repo)

	// 创建用户
	user, err := service.CreateUser("john@example.com", "John Doe", 25)
	if err != nil {
		fmt.Printf("❌ Error creating user: %v\n", err)
		return
	}
	fmt.Printf("✅ Created user: %s (%s)\n", user.Name, user.Email)

	// 获取用户
	retrievedUser, err := service.GetUserByID(user.ID)
	if err != nil {
		fmt.Printf("❌ Error retrieving user: %v\n", err)
		return
	}
	fmt.Printf("✅ Retrieved user: %s\n", retrievedUser.Name)

	// 更新用户
	updatedUser, err := service.UpdateUserProfile(user.ID, "John Smith", 26)
	if err != nil {
		fmt.Printf("❌ Error updating user: %v\n", err)
		return
	}
	fmt.Printf("✅ Updated user: %s (age: %d)\n", updatedUser.Name, updatedUser.Age)

	// 获取所有用户
	allUsers, err := service.GetAllUsers()
	if err != nil {
		fmt.Printf("❌ Error getting all users: %v\n", err)
		return
	}
	fmt.Printf("✅ Total users: %d\n", len(allUsers))

	fmt.Println("🎉 Clean Architecture demo completed successfully!")
}

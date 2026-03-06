# Go 1.23 项目架构设计模型全面梳理

> 本文档系统性地梳理了 **Go 1.23** 语言项目中常用的架构设计模型，从经典的分层架构到现代的微服务和Serverless架构，为不同规模和类型的项目提供架构选型参考。
>
> **Go 1.23 更新**：
>
> - 利用迭代器模式优化领域事件遍历
> - 使用 `unique` 包优化实体标识内存管理
> - 结合 `iter` 包实现仓储模式的迭代器接口
> - 应用新的 `slices` / `maps` 迭代器简化DTO转换

---

## 目录

- [Go 1.23 项目架构设计模型全面梳理](#go-123-项目架构设计模型全面梳理)
  - [目录](#目录)
  - [1. 分层架构（Layered Architecture）](#1-分层架构layered-architecture)
    - [1.1 概念定义](#11-概念定义)
    - [1.2 架构结构](#12-架构结构)
      - [经典三层架构](#经典三层架构)
      - [四层架构（DDD风格）](#四层架构ddd风格)
    - [1.3 核心原则](#13-核心原则)
    - [1.4 Go实现示例](#14-go实现示例)
      - [项目结构](#项目结构)
      - [完整代码实现](#完整代码实现)
    - [1.5 反例说明](#15-反例说明)
      - [反例1：层间循环依赖](#反例1层间循环依赖)
      - [反例2：跨层调用](#反例2跨层调用)
      - [反例3：贫血领域模型](#反例3贫血领域模型)
    - [1.6 选型指南](#16-选型指南)
    - [1.7 优缺点分析](#17-优缺点分析)
  - [2. 六边形架构（Hexagonal Architecture）](#2-六边形架构hexagonal-architecture)
    - [2.1 概念定义](#21-概念定义)
    - [2.2 架构结构](#22-架构结构)
    - [2.3 核心原则](#23-核心原则)
    - [2.4 Go实现示例](#24-go实现示例)
      - [项目结构](#项目结构-1)
      - [完整代码实现](#完整代码实现-1)
    - [2.5 反例说明](#25-反例说明)
      - [反例1：领域核心依赖框架](#反例1领域核心依赖框架)
      - [反例2：适配器直接操作领域对象](#反例2适配器直接操作领域对象)
      - [反例3：端口定义过于具体](#反例3端口定义过于具体)
    - [2.6 选型指南](#26-选型指南)
    - [2.7 优缺点分析](#27-优缺点分析)
  - [3. 洋葱架构（Onion Architecture）](#3-洋葱架构onion-architecture)
    - [3.1 概念定义](#31-概念定义)
    - [3.2 架构结构](#32-架构结构)
    - [3.3 核心原则](#33-核心原则)
    - [3.4 与六边形架构对比](#34-与六边形架构对比)
    - [3.5 Go实现示例](#35-go实现示例)
      - [项目结构](#项目结构-2)
      - [完整代码实现](#完整代码实现-2)
    - [3.6 反例说明](#36-反例说明)
      - [反例1：领域层依赖外层](#反例1领域层依赖外层)
      - [反例2：应用服务包含领域逻辑](#反例2应用服务包含领域逻辑)
      - [反例3：值对象可变](#反例3值对象可变)
    - [3.7 选型指南](#37-选型指南)
    - [3.8 优缺点分析](#38-优缺点分析)
  - [4. 清洁架构（Clean Architecture）](#4-清洁架构clean-architecture)
    - [4.1 概念定义](#41-概念定义)
    - [4.2 架构结构](#42-架构结构)
    - [4.3 核心原则](#43-核心原则)
    - [4.4 Go实现示例](#44-go实现示例)
      - [项目结构](#项目结构-3)
      - [完整代码实现](#完整代码实现-3)
    - [4.5 反例说明](#45-反例说明)
      - [反例1：用例依赖框架](#反例1用例依赖框架)
      - [反例2：实体依赖数据库](#反例2实体依赖数据库)
    - [4.6 选型指南](#46-选型指南)
    - [4.7 优缺点分析](#47-优缺点分析)
  - [5. CQRS（命令查询职责分离）](#5-cqrs命令查询职责分离)
    - [5.1 概念定义](#51-概念定义)
    - [5.2 架构结构](#52-架构结构)
    - [5.3 核心原则](#53-核心原则)
    - [5.4 Go实现示例](#54-go实现示例)
      - [项目结构](#项目结构-4)
      - [完整代码实现](#完整代码实现-4)
    - [5.5 反例说明](#55-反例说明)
      - [反例1：读写模型混合](#反例1读写模型混合)
      - [反例2：缺少事件同步](#反例2缺少事件同步)
    - [5.6 选型指南](#56-选型指南)
    - [5.7 优缺点分析](#57-优缺点分析)
  - [6. 事件溯源（Event Sourcing）](#6-事件溯源event-sourcing)
    - [6.1 概念定义](#61-概念定义)
    - [6.2 架构结构](#62-架构结构)
    - [6.3 核心原则](#63-核心原则)
    - [6.4 Go实现示例](#64-go实现示例)
      - [项目结构](#项目结构-5)
      - [完整代码实现](#完整代码实现-5)
    - [6.5 反例说明](#65-反例说明)
      - [反例1：直接修改状态](#反例1直接修改状态)
      - [反例2：事件修改](#反例2事件修改)
    - [6.6 选型指南](#66-选型指南)
    - [6.7 优缺点分析](#67-优缺点分析)
  - [7. 领域驱动设计（DDD）](#7-领域驱动设计ddd)
    - [7.1 概念定义](#71-概念定义)
    - [7.2 架构结构](#72-架构结构)
    - [7.3 核心概念](#73-核心概念)
      - [7.3.1 限界上下文（Bounded Context）](#731-限界上下文bounded-context)
      - [7.3.2 实体（Entity）与值对象（Value Object）](#732-实体entity与值对象value-object)
      - [7.3.3 聚合（Aggregate）与聚合根（Aggregate Root）](#733-聚合aggregate与聚合根aggregate-root)
    - [7.4 Go实现示例](#74-go实现示例)
      - [项目结构](#项目结构-6)
      - [完整代码实现](#完整代码实现-6)
    - [7.5 反例说明](#75-反例说明)
      - [反例1：贫血领域模型](#反例1贫血领域模型)
      - [反例2：跨聚合直接引用](#反例2跨聚合直接引用)
    - [7.6 选型指南](#76-选型指南)
    - [7.7 优缺点分析](#77-优缺点分析)
  - [8. 微服务架构](#8-微服务架构)
    - [8.1 概念定义](#81-概念定义)
    - [8.2 架构结构](#82-架构结构)
    - [8.3 核心原则](#83-核心原则)
    - [8.4 服务拆分策略](#84-服务拆分策略)
    - [8.5 Go实现示例](#85-go实现示例)
      - [项目结构](#项目结构-7)
      - [完整代码实现](#完整代码实现-7)
    - [8.6 数据一致性](#86-数据一致性)
    - [8.7 反例说明](#87-反例说明)
      - [反例1：共享数据库](#反例1共享数据库)
      - [反例2：同步调用链过长](#反例2同步调用链过长)
    - [8.8 选型指南](#88-选型指南)
    - [8.9 优缺点分析](#89-优缺点分析)
  - [9. Serverless架构](#9-serverless架构)
    - [9.1 概念定义](#91-概念定义)
    - [9.2 架构结构](#92-架构结构)
    - [9.3 核心原则](#93-核心原则)
    - [9.4 冷启动优化](#94-冷启动优化)
    - [9.5 Go实现示例](#95-go实现示例)
      - [项目结构](#项目结构-8)
      - [完整代码实现](#完整代码实现-8)
    - [9.6 反例说明](#96-反例说明)
      - [反例1：函数过大](#反例1函数过大)
      - [反例2：在函数内保持状态](#反例2在函数内保持状态)
    - [9.7 选型指南](#97-选型指南)
    - [9.8 优缺点分析](#98-优缺点分析)
  - [10. 架构选型指南](#10-架构选型指南)
    - [10.1 项目规模与架构选择](#101-项目规模与架构选择)
    - [10.2 团队能力与架构复杂度](#102-团队能力与架构复杂度)
    - [10.3 演进式架构](#103-演进式架构)
    - [10.4 反模式识别](#104-反模式识别)
    - [10.5 架构决策矩阵](#105-架构决策矩阵)
    - [10.6 架构选型决策树](#106-架构选型决策树)
    - [10.7 架构组合建议](#107-架构组合建议)
    - [10.8 实施建议](#108-实施建议)
  - [总结](#总结)
    - [核心要点回顾](#核心要点回顾)
    - [架构选型原则](#架构选型原则)

---

## 1. 分层架构（Layered Architecture）

### 1.1 概念定义

分层架构是最经典、最广泛应用的软件架构模式之一。
它将系统划分为若干个水平层次，每个层次都有明确的职责和边界。
上层依赖下层，下层为上层提供服务，形成清晰的调用链。

分层架构的核心思想是**关注点分离（Separation of Concerns）**，通过将系统分解为不同的层次，降低系统的复杂度，提高可维护性和可测试性。

### 1.2 架构结构

#### 经典三层架构

```
┌─────────────────────────────────────────┐
│           表现层 (Presentation)          │  ← HTTP/GRPC Handler, Controller
│         处理用户交互和输入输出            │
├─────────────────────────────────────────┤
│           业务逻辑层 (Business Logic)    │  ← Service, Use Case
│         核心业务规则和流程处理            │
├─────────────────────────────────────────┤
│           数据访问层 (Data Access)       │  ← Repository, DAO
│         数据持久化和外部资源访问          │
└─────────────────────────────────────────┘
```

#### 四层架构（DDD风格）

```
┌─────────────────────────────────────────┐
│           用户接口层 (Interfaces)        │  ← Controller, Handler, DTO
│         适配用户接口和技术框架            │
├─────────────────────────────────────────┤
│           应用层 (Application)           │  ← Application Service, DTO
│         协调用例和编排领域对象            │
├─────────────────────────────────────────┤
│           领域层 (Domain)                │  ← Entity, Value Object, Domain Service
│         核心业务逻辑和业务规则            │
├─────────────────────────────────────────┤
│           基础设施层 (Infrastructure)    │  ← RepositoryImpl, DB, Cache, MQ
│         技术实现和外部资源访问            │
└─────────────────────────────────────────┘
```

### 1.3 核心原则

| 原则 | 说明 |
|------|------|
| **单向依赖** | 上层依赖下层，下层不依赖上层 |
| **层间隔离** | 每层只与直接相邻的层交互 |
| **抽象接口** | 层与层之间通过接口交互，而非具体实现 |
| **职责单一** | 每层只负责特定的关注点 |

### 1.4 Go实现示例

#### 项目结构

```
layered_architecture/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── domain/           # 领域层
│   │   ├── entity/
│   │   │   └── user.go
│   │   ├── repository/
│   │   │   └── user_repository.go
│   │   └── service/
│   │       └── user_service.go
│   ├── application/      # 应用层
│   │   ├── dto/
│   │   │   └── user_dto.go
│   │   └── service/
│   │       └── user_app_service.go
│   ├── interfaces/       # 接口层
│   │   ├── http/
│   │   │   └── user_handler.go
│   │   └── grpc/
│   └── infrastructure/   # 基础设施层
│       ├── persistence/
│       │   └── user_repository_impl.go
│       └── config/
│           └── database.go
└── go.mod
```

#### 完整代码实现

```go
// internal/domain/entity/user.go
package entity

import (
 "errors"
 "time"
)

// User 领域实体
type User struct {
 ID        string
 Email     string
 Name      string
 CreatedAt time.Time
 UpdatedAt time.Time
}

// Validate 验证用户数据有效性
func (u *User) Validate() error {
 if u.Email == "" {
  return errors.New("email is required")
 }
 if u.Name == "" {
  return errors.New("name is required")
 }
 return nil
}

// UpdateName 更新用户名
func (u *User) UpdateName(name string) {
 u.Name = name
 u.UpdatedAt = time.Now()
}
```

```go
// internal/domain/repository/user_repository.go
package repository

import (
 "context"
 "layered_architecture/internal/domain/entity"
)

// UserRepository 用户仓储接口（领域层定义）
type UserRepository interface {
 FindByID(ctx context.Context, id string) (*entity.User, error)
 FindByEmail(ctx context.Context, email string) (*entity.User, error)
 Save(ctx context.Context, user *entity.User) error
 Delete(ctx context.Context, id string) error
}
```

```go
// internal/domain/service/user_service.go
package service

import (
 "context"
 "errors"
 "layered_architecture/internal/domain/entity"
 "layered_architecture/internal/domain/repository"
)

// UserDomainService 领域服务
type UserDomainService struct {
 userRepo repository.UserRepository
}

// NewUserDomainService 创建领域服务
func NewUserDomainService(userRepo repository.UserRepository) *UserDomainService {
 return &UserDomainService{userRepo: userRepo}
}

// CanRegister 检查用户是否可以注册（业务规则）
func (s *UserDomainService) CanRegister(ctx context.Context, email string) error {
 existing, err := s.userRepo.FindByEmail(ctx, email)
 if err != nil {
  return err
 }
 if existing != nil {
  return errors.New("user with this email already exists")
 }
 return nil
}
```

```go
// internal/application/dto/user_dto.go
package dto

import "time"

// CreateUserRequest 创建用户请求DTO
type CreateUserRequest struct {
 Email string `json:"email"`
 Name  string `json:"name"`
}

// UserResponse 用户响应DTO
type UserResponse struct {
 ID        string    `json:"id"`
 Email     string    `json:"email"`
 Name      string    `json:"name"`
 CreatedAt time.Time `json:"created_at"`
}

// UpdateUserRequest 更新用户请求DTO
type UpdateUserRequest struct {
 Name string `json:"name"`
}
```

```go
// internal/application/service/user_app_service.go
package service

import (
 "context"
 "layered_architecture/internal/application/dto"
 "layered_architecture/internal/domain/entity"
 "layered_architecture/internal/domain/repository"
 domainService "layered_architecture/internal/domain/service"
 "time"

 "github.com/google/uuid"
)

// UserApplicationService 应用服务
type UserApplicationService struct {
 userRepo    repository.UserRepository
 domainSvc   *domainService.UserDomainService
}

// NewUserApplicationService 创建应用服务
func NewUserApplicationService(
 userRepo repository.UserRepository,
) *UserApplicationService {
 return &UserApplicationService{
  userRepo:  userRepo,
  domainSvc: domainService.NewUserDomainService(userRepo),
 }
}

// CreateUser 创建用户用例
func (s *UserApplicationService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
 // 1. 检查业务规则
 if err := s.domainSvc.CanRegister(ctx, req.Email); err != nil {
  return nil, err
 }

 // 2. 创建领域实体
 user := &entity.User{
  ID:        uuid.New().String(),
  Email:     req.Email,
  Name:      req.Name,
  CreatedAt: time.Now(),
  UpdatedAt: time.Now(),
 }

 // 3. 验证实体
 if err := user.Validate(); err != nil {
  return nil, err
 }

 // 4. 保存
 if err := s.userRepo.Save(ctx, user); err != nil {
  return nil, err
 }

 // 5. 返回DTO
 return &dto.UserResponse{
  ID:        user.ID,
  Email:     user.Email,
  Name:      user.Name,
  CreatedAt: user.CreatedAt,
 }, nil
}

// GetUser 获取用户用例
func (s *UserApplicationService) GetUser(ctx context.Context, id string) (*dto.UserResponse, error) {
 user, err := s.userRepo.FindByID(ctx, id)
 if err != nil {
  return nil, err
 }

 return &dto.UserResponse{
  ID:        user.ID,
  Email:     user.Email,
  Name:      user.Name,
  CreatedAt: user.CreatedAt,
 }, nil
}

// UpdateUser 更新用户用例
func (s *UserApplicationService) UpdateUser(ctx context.Context, id string, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
 user, err := s.userRepo.FindByID(ctx, id)
 if err != nil {
  return nil, err
 }

 user.UpdateName(req.Name)

 if err := s.userRepo.Save(ctx, user); err != nil {
  return nil, err
 }

 return &dto.UserResponse{
  ID:        user.ID,
  Email:     user.Email,
  Name:      user.Name,
  CreatedAt: user.CreatedAt,
 }, nil
}
```

```go
// internal/interfaces/http/user_handler.go
package http

import (
 "encoding/json"
 "layered_architecture/internal/application/dto"
 "layered_architecture/internal/application/service"
 "net/http"
 "strings"
)

// UserHandler HTTP处理器
type UserHandler struct {
 appService *service.UserApplicationService
}

// NewUserHandler 创建处理器
func NewUserHandler(appService *service.UserApplicationService) *UserHandler {
 return &UserHandler{appService: appService}
}

// ServeHTTP 实现http.Handler接口
func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
 switch r.Method {
 case http.MethodPost:
  h.createUser(w, r)
 case http.MethodGet:
  h.getUser(w, r)
 case http.MethodPut:
  h.updateUser(w, r)
 default:
  w.WriteHeader(http.StatusMethodNotAllowed)
 }
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
 var req dto.CreateUserRequest
 if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
  h.respondError(w, http.StatusBadRequest, err.Error())
  return
 }

 user, err := h.appService.CreateUser(r.Context(), req)
 if err != nil {
  h.respondError(w, http.StatusBadRequest, err.Error())
  return
 }

 h.respondJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
 id := strings.TrimPrefix(r.URL.Path, "/users/")
 user, err := h.appService.GetUser(r.Context(), id)
 if err != nil {
  h.respondError(w, http.StatusNotFound, err.Error())
  return
 }

 h.respondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) {
 id := strings.TrimPrefix(r.URL.Path, "/users/")
 var req dto.UpdateUserRequest
 if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
  h.respondError(w, http.StatusBadRequest, err.Error())
  return
 }

 user, err := h.appService.UpdateUser(r.Context(), id, req)
 if err != nil {
  h.respondError(w, http.StatusBadRequest, err.Error())
  return
 }

 h.respondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
 w.Header().Set("Content-Type", "application/json")
 w.WriteHeader(status)
 json.NewEncoder(w).Encode(data)
}

func (h *UserHandler) respondError(w http.ResponseWriter, status int, message string) {
 h.respondJSON(w, status, map[string]string{"error": message})
}
```

```go
// internal/infrastructure/persistence/user_repository_impl.go
package persistence

import (
 "context"
 "layered_architecture/internal/domain/entity"
 "layered_architecture/internal/domain/repository"
 "sync"
)

// InMemoryUserRepository 内存用户仓储实现
type InMemoryUserRepository struct {
 mu    sync.RWMutex
 users map[string]*entity.User
 index map[string]string // email -> id
}

// NewInMemoryUserRepository 创建内存仓储
func NewInMemoryUserRepository() repository.UserRepository {
 return &InMemoryUserRepository{
  users: make(map[string]*entity.User),
  index: make(map[string]string),
 }
}

func (r *InMemoryUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
 r.mu.RLock()
 defer r.mu.RUnlock()

 user, ok := r.users[id]
 if !ok {
  return nil, nil
 }
 // 返回副本避免外部修改
 return copyUser(user), nil
}

func (r *InMemoryUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
 r.mu.RLock()
 defer r.mu.RUnlock()

 id, ok := r.index[email]
 if !ok {
  return nil, nil
 }
 return copyUser(r.users[id]), nil
}

func (r *InMemoryUserRepository) Save(ctx context.Context, user *entity.User) error {
 r.mu.Lock()
 defer r.mu.Unlock()

 r.users[user.ID] = copyUser(user)
 r.index[user.Email] = user.ID
 return nil
}

func (r *InMemoryUserRepository) Delete(ctx context.Context, id string) error {
 r.mu.Lock()
 defer r.mu.Unlock()

 if user, ok := r.users[id]; ok {
  delete(r.index, user.Email)
  delete(r.users, id)
 }
 return nil
}

func copyUser(u *entity.User) *entity.User {
 return &entity.User{
  ID:        u.ID,
  Email:     u.Email,
  Name:      u.Name,
  CreatedAt: u.CreatedAt,
  UpdatedAt: u.UpdatedAt,
 }
}
```

```go
// cmd/api/main.go
package main

import (
 "fmt"
 "layered_architecture/internal/application/service"
 "layered_architecture/internal/infrastructure/persistence"
 httpInterface "layered_architecture/internal/interfaces/http"
 "net/http"
)

func main() {
 // 1. 初始化基础设施层
 userRepo := persistence.NewInMemoryUserRepository()

 // 2. 初始化应用层
 userAppService := service.NewUserApplicationService(userRepo)

 // 3. 初始化接口层
 userHandler := httpInterface.NewUserHandler(userAppService)

 // 4. 启动HTTP服务
 http.Handle("/users/", userHandler)
 http.HandleFunc("/users", userHandler.ServeHTTP)

 fmt.Println("Server starting on :8080")
 if err := http.ListenAndServe(":8080", nil); err != nil {
  panic(err)
 }
}
```

### 1.5 反例说明

#### 反例1：层间循环依赖

```go
// ❌ 错误：领域层直接依赖基础设施层
package domain

import (
 "database/sql"  // 领域层不应该直接依赖数据库
)

type UserService struct {
 db *sql.DB  // 直接依赖具体实现
}
```

**问题**：

- 领域层与数据库实现耦合，无法独立测试
- 更换数据库需要修改领域层代码
- 违反依赖倒置原则

#### 反例2：跨层调用

```go
// ❌ 错误：表现层直接调用数据访问层
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
 // 直接操作数据库，绕过业务逻辑层
 result, err := db.Exec("INSERT INTO users ...")
 // ...
}
```

**问题**：

- 业务逻辑散落在各处，难以维护
- 无法复用业务规则
- 安全性和一致性无法保证

#### 反例3：贫血领域模型

```go
// ❌ 错误：领域对象只有数据，没有行为
type User struct {
 ID    string
 Email string
 Name  string
}

// 所有业务逻辑都在服务中
func (s *UserService) ValidateUser(user *User) error {
 if user.Email == "" {
  return errors.New("email required")
 }
 return nil
}
```

**问题**：

- 领域对象退化为数据结构
- 业务逻辑与服务耦合
- 违反面向对象设计原则

### 1.6 选型指南

| 场景 | 建议 |
|------|------|
| **小型项目** | 三层架构足够，快速开发 |
| **中型项目** | 四层架构，引入领域层分离业务逻辑 |
| **大型项目** | 四层架构 + DDD，配合六边形/洋葱架构 |
| **快速原型** | 简化分层，甚至单层，快速验证 |
| **遗留系统** | 渐进式重构，先引入分层概念 |

### 1.7 优缺点分析

| 优点 | 缺点 |
|------|------|
| ✅ 结构清晰，易于理解 | ❌ 过度分层可能导致性能损耗 |
| ✅ 职责分离，便于维护 | ❌ 小型项目可能过度设计 |
| ✅ 支持并行开发 | ❌ 层间通信增加代码复杂度 |
| ✅ 便于测试，可分层测试 | ❌ 严格的层间依赖可能限制灵活性 |
| ✅ 技术栈可逐层替换 | ❌ 需要团队理解和遵守分层规范 |

---

## 2. 六边形架构（Hexagonal Architecture）

### 2.1 概念定义

六边形架构（Hexagonal Architecture），又称**端口与适配器架构（Ports and Adapters Architecture）**，由Alistair Cockburn于2005年提出。其核心思想是将应用程序的核心业务逻辑与外部世界隔离，通过定义明确的**端口（Ports）**和**适配器（Adapters）**进行通信。

六边形架构的名称来源于其图形表示：应用程序核心位于六边形中心，外部系统通过适配器连接到端口，形成六边形的边界。

### 2.2 架构结构

```
                    ┌─────────────────┐
                    │   Web UI/API    │
                    │   (适配器)       │
                    └────────┬────────┘
                             │
                    ┌────────▼────────┐
                    │                 │
                    │   应用程序核心   │
                    │   (领域逻辑)     │
                    │                 │
        ┌───────────┤                 ├───────────┐
        │           │                 │           │
┌───────▼──────┐    │                 │    ┌──────▼───────┐
│  CLI/Console │    │                 │    │  Message Queue│
│   (适配器)    │    │                 │    │   (适配器)    │
└──────────────┘    │                 │    └──────────────┘
                    │                 │
        ┌───────────┤                 ├───────────┐
        │           │                 │           │
┌───────▼──────┐    │                 │    ┌──────▼───────┐
│   Database   │    │                 │    │ External API │
│   (适配器)    │    └─────────────────┘    │   (适配器)    │
└──────────────┘                             └──────────────┘

              ◄── 端口（接口定义）──►
              ◄──── 适配器（实现）────►
```

### 2.3 核心原则

| 原则 | 说明 |
|------|------|
| **依赖倒置** | 依赖关系指向领域核心，而非外部 |
| **端口定义** | 通过接口定义输入/输出契约 |
| **适配器实现** | 外部系统通过适配器实现端口 |
| **领域隔离** | 核心业务逻辑完全独立于框架和UI |

### 2.4 Go实现示例

#### 项目结构

```
hexagonal_architecture/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── core/              # 领域核心
│   │   ├── domain/
│   │   │   ├── order.go
│   │   │   └── order_item.go
│   │   ├── ports/
│   │   │   ├── incoming/  # 入站端口（驱动端口）
│   │   │   │   └── order_service.go
│   │   │   └── outgoing/  # 出站端口（被驱动端口）
│   │   │       ├── order_repository.go
│   │   │       └── payment_gateway.go
│   │   └── services/
│   │       └── order_service_impl.go
│   ├── adapters/          # 适配器
│   │   ├── incoming/      # 入站适配器（驱动适配器）
│   │   │   ├── http/
│   │   │   │   └── order_handler.go
│   │   │   └── grpc/
│   │   │       └── order_grpc_server.go
│   │   └── outgoing/      # 出站适配器（被驱动适配器）
│   │       ├── persistence/
│   │       │   └── order_repository_impl.go
│   │       └── payment/
│   │           └── stripe_adapter.go
│   └── config/
│       └── app.go
└── go.mod
```

#### 完整代码实现

```go
// internal/core/domain/order.go
package domain

import (
 "errors"
 "time"
)

// OrderStatus 订单状态
type OrderStatus string

const (
 OrderStatusPending   OrderStatus = "PENDING"
 OrderStatusPaid      OrderStatus = "PAID"
 OrderStatusShipped   OrderStatus = "SHIPPED"
 OrderStatusCompleted OrderStatus = "COMPLETED"
 OrderStatusCancelled OrderStatus = "CANCELLED"
)

// Order 订单领域实体
type Order struct {
 ID         string
 CustomerID string
 Items      []OrderItem
 Status     OrderStatus
 Total      float64
 CreatedAt  time.Time
 UpdatedAt  time.Time
}

// OrderItem 订单项
type OrderItem struct {
 ProductID string
 Name      string
 Quantity  int
 Price     float64
}

// CalculateTotal 计算订单总价
func (o *Order) CalculateTotal() float64 {
 total := 0.0
 for _, item := range o.Items {
  total += item.Price * float64(item.Quantity)
 }
 o.Total = total
 return total
}

// Pay 支付订单
func (o *Order) Pay() error {
 if o.Status != OrderStatusPending {
  return errors.New("only pending orders can be paid")
 }
 o.Status = OrderStatusPaid
 o.UpdatedAt = time.Now()
 return nil
}

// AddItem 添加订单项
func (o *Order) AddItem(item OrderItem) error {
 if item.Quantity <= 0 {
  return errors.New("quantity must be positive")
 }
 if item.Price < 0 {
  return errors.New("price cannot be negative")
 }
 o.Items = append(o.Items, item)
 o.CalculateTotal()
 o.UpdatedAt = time.Now()
 return nil
}

// Validate 验证订单
func (o *Order) Validate() error {
 if o.CustomerID == "" {
  return errors.New("customer ID is required")
 }
 if len(o.Items) == 0 {
  return errors.New("order must have at least one item")
 }
 return nil
}
```

```go
// internal/core/ports/incoming/order_service.go
package incoming

import (
 "context"
 "hexagonal_architecture/internal/core/domain"
)

// OrderService 入站端口 - 订单服务接口
type OrderService interface {
 CreateOrder(ctx context.Context, customerID string, items []domain.OrderItem) (*domain.Order, error)
 GetOrder(ctx context.Context, orderID string) (*domain.Order, error)
 PayOrder(ctx context.Context, orderID string, paymentMethod string) error
 CancelOrder(ctx context.Context, orderID string) error
}
```

```go
// internal/core/ports/outgoing/order_repository.go
package outgoing

import (
 "context"
 "hexagonal_architecture/internal/core/domain"
)

// OrderRepository 出站端口 - 订单仓储接口
type OrderRepository interface {
 FindByID(ctx context.Context, id string) (*domain.Order, error)
 FindByCustomerID(ctx context.Context, customerID string) ([]*domain.Order, error)
 Save(ctx context.Context, order *domain.Order) error
 Update(ctx context.Context, order *domain.Order) error
}
```

```go
// internal/core/ports/outgoing/payment_gateway.go
package outgoing

import (
 "context"
 "hexagonal_architecture/internal/core/domain"
)

// PaymentResult 支付结果
type PaymentResult struct {
 Success    bool
 Message    string
 ExternalID string
}

// PaymentGateway 出站端口 - 支付网关接口
type PaymentGateway interface {
 ProcessPayment(ctx context.Context, order *domain.Order, method string) (*PaymentResult, error)
 RefundPayment(ctx context.Context, orderID string, amount float64) (*PaymentResult, error)
}
```

```go
// internal/core/services/order_service_impl.go
package services

import (
 "context"
 "errors"
 "hexagonal_architecture/internal/core/domain"
 "hexagonal_architecture/internal/core/ports/incoming"
 "hexagonal_architecture/internal/core/ports/outgoing"
 "time"

 "github.com/google/uuid"
)

// OrderServiceImpl 订单服务实现
type OrderServiceImpl struct {
 orderRepo      outgoing.OrderRepository
 paymentGateway outgoing.PaymentGateway
}

// NewOrderService 创建订单服务
func NewOrderService(
 orderRepo outgoing.OrderRepository,
 paymentGateway outgoing.PaymentGateway,
) incoming.OrderService {
 return &OrderServiceImpl{
  orderRepo:      orderRepo,
  paymentGateway: paymentGateway,
 }
}

// CreateOrder 创建订单
func (s *OrderServiceImpl) CreateOrder(
 ctx context.Context,
 customerID string,
 items []domain.OrderItem,
) (*domain.Order, error) {
 order := &domain.Order{
  ID:         uuid.New().String(),
  CustomerID: customerID,
  Items:      items,
  Status:     domain.OrderStatusPending,
  CreatedAt:  time.Now(),
  UpdatedAt:  time.Now(),
 }

 order.CalculateTotal()

 if err := order.Validate(); err != nil {
  return nil, err
 }

 if err := s.orderRepo.Save(ctx, order); err != nil {
  return nil, err
 }

 return order, nil
}

// GetOrder 获取订单
func (s *OrderServiceImpl) GetOrder(ctx context.Context, orderID string) (*domain.Order, error) {
 return s.orderRepo.FindByID(ctx, orderID)
}

// PayOrder 支付订单
func (s *OrderServiceImpl) PayOrder(
 ctx context.Context,
 orderID string,
 paymentMethod string,
) error {
 order, err := s.orderRepo.FindByID(ctx, orderID)
 if err != nil {
  return err
 }
 if order == nil {
  return errors.New("order not found")
 }

 // 处理支付
 result, err := s.paymentGateway.ProcessPayment(ctx, order, paymentMethod)
 if err != nil {
  return err
 }
 if !result.Success {
  return errors.New("payment failed: " + result.Message)
 }

 // 更新订单状态
 if err := order.Pay(); err != nil {
  return err
 }

 return s.orderRepo.Update(ctx, order)
}

// CancelOrder 取消订单
func (s *OrderServiceImpl) CancelOrder(ctx context.Context, orderID string) error {
 order, err := s.orderRepo.FindByID(ctx, orderID)
 if err != nil {
  return err
 }
 if order == nil {
  return errors.New("order not found")
 }

 if order.Status != domain.OrderStatusPending {
  return errors.New("only pending orders can be cancelled")
 }

 order.Status = domain.OrderStatusCancelled
 order.UpdatedAt = time.Now()

 return s.orderRepo.Update(ctx, order)
}
```

```go
// internal/adapters/incoming/http/order_handler.go
package http

import (
 "encoding/json"
 "hexagonal_architecture/internal/core/domain"
 "hexagonal_architecture/internal/core/ports/incoming"
 "net/http"
 "strings"
)

// OrderHandler HTTP处理器
type OrderHandler struct {
 service incoming.OrderService
}

// NewOrderHandler 创建处理器
func NewOrderHandler(service incoming.OrderService) *OrderHandler {
 return &OrderHandler{service: service}
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
 CustomerID string                `json:"customer_id"`
 Items      []domain.OrderItem    `json:"items"`
}

// OrderResponse 订单响应
type OrderResponse struct {
 ID         string              `json:"id"`
 CustomerID string              `json:"customer_id"`
 Items      []domain.OrderItem  `json:"items"`
 Status     string              `json:"status"`
 Total      float64             `json:"total"`
}

// ServeHTTP 实现http.Handler
func (h *OrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
 switch r.Method {
 case http.MethodPost:
  if strings.HasSuffix(r.URL.Path, "/pay") {
   h.payOrder(w, r)
  } else {
   h.createOrder(w, r)
  }
 case http.MethodGet:
  h.getOrder(w, r)
 case http.MethodDelete:
  h.cancelOrder(w, r)
 default:
  w.WriteHeader(http.StatusMethodNotAllowed)
 }
}

func (h *OrderHandler) createOrder(w http.ResponseWriter, r *http.Request) {
 var req CreateOrderRequest
 if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
  h.respondError(w, http.StatusBadRequest, err.Error())
  return
 }

 order, err := h.service.CreateOrder(r.Context(), req.CustomerID, req.Items)
 if err != nil {
  h.respondError(w, http.StatusBadRequest, err.Error())
  return
 }

 h.respondOrder(w, http.StatusCreated, order)
}

func (h *OrderHandler) getOrder(w http.ResponseWriter, r *http.Request) {
 orderID := strings.TrimPrefix(r.URL.Path, "/orders/")
 order, err := h.service.GetOrder(r.Context(), orderID)
 if err != nil {
  h.respondError(w, http.StatusNotFound, err.Error())
  return
 }

 h.respondOrder(w, http.StatusOK, order)
}

func (h *OrderHandler) payOrder(w http.ResponseWriter, r *http.Request) {
 orderID := strings.TrimPrefix(r.URL.Path, "/orders/")
 orderID = strings.TrimSuffix(orderID, "/pay")

 var req struct {
  PaymentMethod string `json:"payment_method"`
 }
 if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
  h.respondError(w, http.StatusBadRequest, err.Error())
  return
 }

 if err := h.service.PayOrder(r.Context(), orderID, req.PaymentMethod); err != nil {
  h.respondError(w, http.StatusBadRequest, err.Error())
  return
 }

 h.respondJSON(w, http.StatusOK, map[string]string{"status": "paid"})
}

func (h *OrderHandler) cancelOrder(w http.ResponseWriter, r *http.Request) {
 orderID := strings.TrimPrefix(r.URL.Path, "/orders/")
 if err := h.service.CancelOrder(r.Context(), orderID); err != nil {
  h.respondError(w, http.StatusBadRequest, err.Error())
  return
 }

 h.respondJSON(w, http.StatusOK, map[string]string{"status": "cancelled"})
}

func (h *OrderHandler) respondOrder(w http.ResponseWriter, status int, order *domain.Order) {
 resp := OrderResponse{
  ID:         order.ID,
  CustomerID: order.CustomerID,
  Items:      order.Items,
  Status:     string(order.Status),
  Total:      order.Total,
 }
 h.respondJSON(w, status, resp)
}

func (h *OrderHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
 w.Header().Set("Content-Type", "application/json")
 w.WriteHeader(status)
 json.NewEncoder(w).Encode(data)
}

func (h *OrderHandler) respondError(w http.ResponseWriter, status int, message string) {
 h.respondJSON(w, status, map[string]string{"error": message})
}
```

```go
// internal/adapters/outgoing/persistence/order_repository_impl.go
package persistence

import (
 "context"
 "hexagonal_architecture/internal/core/domain"
 "hexagonal_architecture/internal/core/ports/outgoing"
 "sync"
)

// InMemoryOrderRepository 内存订单仓储实现
type InMemoryOrderRepository struct {
 mu     sync.RWMutex
 orders map[string]*domain.Order
}

// NewInMemoryOrderRepository 创建内存仓储
func NewInMemoryOrderRepository() outgoing.OrderRepository {
 return &InMemoryOrderRepository{
  orders: make(map[string]*domain.Order),
 }
}

func (r *InMemoryOrderRepository) FindByID(ctx context.Context, id string) (*domain.Order, error) {
 r.mu.RLock()
 defer r.mu.RUnlock()

 order, ok := r.orders[id]
 if !ok {
  return nil, nil
 }
 return copyOrder(order), nil
}

func (r *InMemoryOrderRepository) FindByCustomerID(ctx context.Context, customerID string) ([]*domain.Order, error) {
 r.mu.RLock()
 defer r.mu.RUnlock()

 var result []*domain.Order
 for _, order := range r.orders {
  if order.CustomerID == customerID {
   result = append(result, copyOrder(order))
  }
 }
 return result, nil
}

func (r *InMemoryOrderRepository) Save(ctx context.Context, order *domain.Order) error {
 r.mu.Lock()
 defer r.mu.Unlock()

 r.orders[order.ID] = copyOrder(order)
 return nil
}

func (r *InMemoryOrderRepository) Update(ctx context.Context, order *domain.Order) error {
 r.mu.Lock()
 defer r.mu.Unlock()

 r.orders[order.ID] = copyOrder(order)
 return nil
}

func copyOrder(o *domain.Order) *domain.Order {
 items := make([]domain.OrderItem, len(o.Items))
 copy(items, o.Items)
 return &domain.Order{
  ID:         o.ID,
  CustomerID: o.CustomerID,
  Items:      items,
  Status:     o.Status,
  Total:      o.Total,
  CreatedAt:  o.CreatedAt,
  UpdatedAt:  o.UpdatedAt,
 }
}
```

```go
// internal/adapters/outgoing/payment/stripe_adapter.go
package payment

import (
 "context"
 "fmt"
 "hexagonal_architecture/internal/core/domain"
 "hexagonal_architecture/internal/core/ports/outgoing"
)

// StripeAdapter Stripe支付适配器
type StripeAdapter struct {
 apiKey string
}

// NewStripeAdapter 创建Stripe适配器
func NewStripeAdapter(apiKey string) outgoing.PaymentGateway {
 return &StripeAdapter{apiKey: apiKey}
}

// ProcessPayment 处理支付
func (a *StripeAdapter) ProcessPayment(
 ctx context.Context,
 order *domain.Order,
 method string,
) (*outgoing.PaymentResult, error) {
 // 模拟Stripe支付处理
 fmt.Printf("[Stripe] Processing payment for order %s, amount: %.2f, method: %s\n",
  order.ID, order.Total, method)

 // 实际实现中调用Stripe API
 return &outgoing.PaymentResult{
  Success:    true,
  Message:    "Payment processed successfully",
  ExternalID: fmt.Sprintf("stripe_%s", order.ID),
 }, nil
}

// RefundPayment 退款
func (a *StripeAdapter) RefundPayment(
 ctx context.Context,
 orderID string,
 amount float64,
) (*outgoing.PaymentResult, error) {
 fmt.Printf("[Stripe] Refunding payment for order %s, amount: %.2f\n", orderID, amount)

 return &outgoing.PaymentResult{
  Success:    true,
  Message:    "Refund processed successfully",
  ExternalID: fmt.Sprintf("refund_%s", orderID),
 }, nil
}
```

```go
// internal/config/app.go
package config

import (
 "hexagonal_architecture/internal/adapters/incoming/http"
 "hexagonal_architecture/internal/adapters/outgoing/payment"
 "hexagonal_architecture/internal/adapters/outgoing/persistence"
 "hexagonal_architecture/internal/core/ports/incoming"
 "hexagonal_architecture/internal/core/services"
)

// AppConfig 应用配置
type AppConfig struct {
 OrderService incoming.OrderService
 OrderHandler *http.OrderHandler
}

// NewAppConfig 初始化应用配置
func NewAppConfig() *AppConfig {
 // 出站适配器（基础设施）
 orderRepo := persistence.NewInMemoryOrderRepository()
 paymentGateway := payment.NewStripeAdapter("sk_test_...")

 // 核心领域服务
 orderService := services.NewOrderService(orderRepo, paymentGateway)

 // 入站适配器（接口）
 orderHandler := http.NewOrderHandler(orderService)

 return &AppConfig{
  OrderService: orderService,
  OrderHandler: orderHandler,
 }
}
```

```go
// cmd/api/main.go
package main

import (
 "fmt"
 "hexagonal_architecture/internal/config"
 "net/http"
)

func main() {
 // 初始化应用
 app := config.NewAppConfig()

 // 注册路由
 http.Handle("/orders/", app.OrderHandler)
 http.HandleFunc("/orders", app.OrderHandler.ServeHTTP)

 fmt.Println("Hexagonal Architecture Server starting on :8080")
 if err := http.ListenAndServe(":8080", nil); err != nil {
  panic(err)
 }
}
```

### 2.5 反例说明

#### 反例1：领域核心依赖框架

```go
// ❌ 错误：领域核心直接依赖Gin框架
package domain

import "github.com/gin-gonic/gin"

type OrderService struct {
 router *gin.Engine  // 领域服务不应该知道HTTP框架
}

func (s *OrderService) CreateOrder(c *gin.Context) {
 // 直接处理HTTP请求
}
```

**问题**：

- 领域核心与框架强耦合
- 无法独立测试业务逻辑
- 更换框架需要重写领域代码

#### 反例2：适配器直接操作领域对象

```go
// ❌ 错误：适配器绕过端口直接操作领域
package http

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
 // 直接创建领域对象，绕过服务端口
 order := &domain.Order{
  ID: uuid.New().String(),
  // ...
 }
 // 直接保存到数据库
 db.Save(order)
}
```

**问题**：

- 业务逻辑散落在适配器
- 无法保证业务规则一致性
- 违反单一职责原则

#### 反例3：端口定义过于具体

```go
// ❌ 错误：端口包含实现细节
package ports

type OrderRepository interface {
 FindByID(ctx context.Context, id string) (*domain.Order, error)
 Save(ctx context.Context, order *domain.Order) error
 // 以下方法过于具体，属于实现细节
 ExecuteSQL(query string, args ...interface{}) (sql.Result, error)
 BeginTransaction() (*sql.Tx, error)
}
```

**问题**：

- 端口暴露实现细节
- 不同的适配器被迫实现不必要的方法
- 违反接口隔离原则

### 2.6 选型指南

| 场景 | 建议 |
|------|------|
| **需要多接口支持** | 六边形架构天然支持多种入站适配器（HTTP、GRPC、CLI等） |
| **频繁更换基础设施** | 端口定义稳定，适配器可独立替换 |
| **需要高可测试性** | 核心逻辑完全独立于外部，易于单元测试 |
| **遗留系统现代化** | 可渐进式引入，先定义端口再迁移适配器 |
| **小型CRUD应用** | 可能过度设计，简单分层即可 |

### 2.7 优缺点分析

| 优点 | 缺点 |
|------|------|
| ✅ 核心业务完全独立于框架 | ❌ 初始设计成本较高 |
| ✅ 易于测试，可mock所有端口 | ❌ 需要更多接口和适配器代码 |
| ✅ 支持多种入站/出站方式 | ❌ 团队需要理解端口/适配器概念 |
| ✅ 基础设施可灵活替换 | ❌ 简单项目可能过度设计 |
| ✅ 清晰的边界和职责分离 | ❌ 学习曲线较陡 |

---

## 3. 洋葱架构（Onion Architecture）

### 3.1 概念定义

洋葱架构（Onion Architecture）由Jeffrey Palermo于2008年提出，是分层架构的演进版本。它以领域模型为核心，通过同心圆的形式组织代码，越靠近中心的层次越稳定，越外层的层次越易变。

洋葱架构强调**领域模型**是应用程序的核心，所有其他层次都围绕领域模型构建，并且只能向内依赖。

### 3.2 架构结构

```
                    ┌─────────────────┐
                    │                 │
                    │   基础设施层     │  ← 最外层：UI、数据库、外部服务
                    │ Infrastructure  │     依赖所有内层
                    │                 │
                    ├─────────────────┤
                    │                 │
                    │   应用服务层     │  ← 协调用例，编排领域对象
                    │  Application    │     依赖领域层
                    │                 │
                    ├─────────────────┤
                    │                 │
                    │   领域服务层     │  ← 领域逻辑，跨实体业务
                    │ Domain Services │     依赖领域模型
                    │                 │
                    ├─────────────────┤
                    │                 │
                    │   领域模型层     │  ← 核心：实体、值对象、领域事件
                    │  Domain Model   │     不依赖任何外层
                    │                 │
                    └─────────────────┘

                    ◄── 依赖方向向内 ──►
```

### 3.3 核心原则

| 原则 | 说明 |
|------|------|
| **依赖向内** | 外层依赖内层，内层不依赖外层 |
| **领域为核心** | 领域模型位于最中心，最稳定 |
| **接口分离** | 通过接口定义层间契约 |
| **可测试性** | 核心业务逻辑可独立测试 |

### 3.4 与六边形架构对比

| 特性 | 洋葱架构 | 六边形架构 |
|------|----------|------------|
| **核心** | 领域模型 | 应用程序核心 |
| **表示方式** | 同心圆 | 六边形 |
| **依赖方向** | 向内 | 向内 |
| **端口概念** | 隐式（通过接口） | 显式定义 |
| **适配器概念** | 隐式 | 显式定义 |
| **关注点** | 层次组织 | 内外交互 |

### 3.5 Go实现示例

#### 项目结构

```
onion_architecture/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── domain/              # 领域模型层（核心）
│   │   ├── entities/
│   │   │   ├── product.go
│   │   │   └── category.go
│   │   ├── valueobjects/
│   │   │   ├── money.go
│   │   │   └── address.go
│   │   ├── events/
│   │   │   └── product_events.go
│   │   └── services/
│   │       └── inventory_service.go
│   ├── application/         # 应用服务层
│   │   ├── interfaces/
│   │   │   ├── repositories/
│   │   │   │   └── product_repository.go
│   │   │   └── services/
│   │   │       └── notification_service.go
│   │   ├── services/
│   │   │   └── product_app_service.go
│   │   └── dto/
│   │       └── product_dto.go
│   ├── infrastructure/      # 基础设施层（最外层）
│   │   ├── persistence/
│   │   │   └── product_repository_impl.go
│   │   ├── web/
│   │   │   └── product_controller.go
│   │   ├── messaging/
│   │   │   └── kafka_publisher.go
│   │   └── external/
│   │       └── email_service.go
│   └── config/
│       └── dependency_injection.go
└── go.mod
```

#### 完整代码实现

```go
// internal/domain/valueobjects/money.go
package valueobjects

import (
 "errors"
 "fmt"
)

// Money 金额值对象
type Money struct {
 Amount   float64
 Currency string
}

// NewMoney 创建金额
func NewMoney(amount float64, currency string) (Money, error) {
 if amount < 0 {
  return Money{}, errors.New("amount cannot be negative")
 }
 if currency == "" {
  return Money{}, errors.New("currency is required")
 }
 return Money{Amount: amount, Currency: currency}, nil
}

// Add 金额相加
func (m Money) Add(other Money) (Money, error) {
 if m.Currency != other.Currency {
  return Money{}, errors.New("cannot add different currencies")
 }
 return NewMoney(m.Amount+other.Amount, m.Currency)
}

// Subtract 金额相减
func (m Money) Subtract(other Money) (Money, error) {
 if m.Currency != other.Currency {
  return Money{}, errors.New("cannot subtract different currencies")
 }
 if m.Amount < other.Amount {
  return Money{}, errors.New("insufficient amount")
 }
 return NewMoney(m.Amount-other.Amount, m.Currency)
}

// Multiply 金额乘法
func (m Money) Multiply(factor float64) (Money, error) {
 return NewMoney(m.Amount*factor, m.Currency)
}

// Equals 金额相等比较
func (m Money) Equals(other Money) bool {
 return m.Amount == other.Amount && m.Currency == other.Currency
}

func (m Money) String() string {
 return fmt.Sprintf("%.2f %s", m.Amount, m.Currency)
}
```

```go
// internal/domain/entities/product.go
package entities

import (
 "errors"
 "onion_architecture/internal/domain/valueobjects"
 "time"
)

// Product 产品实体
type Product struct {
 ID          string
 Name        string
 Description string
 Price       valueobjects.Money
 Stock       int
 CategoryID  string
 CreatedAt   time.Time
 UpdatedAt   time.Time
}

// NewProduct 创建新产品
func NewProduct(id, name, description string, price valueobjects.Money, categoryID string) (*Product, error) {
 if name == "" {
  return nil, errors.New("product name is required")
 }

 return &Product{
  ID:          id,
  Name:        name,
  Description: description,
  Price:       price,
  Stock:       0,
  CategoryID:  categoryID,
  CreatedAt:   time.Now(),
  UpdatedAt:   time.Now(),
 }, nil
}

// UpdateStock 更新库存
func (p *Product) UpdateStock(quantity int) error {
 if p.Stock+quantity < 0 {
  return errors.New("insufficient stock")
 }
 p.Stock += quantity
 p.UpdatedAt = time.Now()
 return nil
}

// UpdatePrice 更新价格
func (p *Product) UpdatePrice(newPrice valueobjects.Money) error {
 p.Price = newPrice
 p.UpdatedAt = time.Now()
 return nil
}

// IsAvailable 检查产品是否可购买
func (p *Product) IsAvailable(quantity int) bool {
 return p.Stock >= quantity
}

// CalculateTotalPrice 计算总价
func (p *Product) CalculateTotalPrice(quantity int) (valueobjects.Money, error) {
 if quantity <= 0 {
  return valueobjects.Money{}, errors.New("quantity must be positive")
 }
 return p.Price.Multiply(float64(quantity))
}

// ReserveStock 预留库存
func (p *Product) ReserveStock(quantity int) error {
 if !p.IsAvailable(quantity) {
  return errors.New("not enough stock available")
 }
 p.Stock -= quantity
 p.UpdatedAt = time.Now()
 return nil
}
```

```go
// internal/domain/events/product_events.go
package events

import (
 "onion_architecture/internal/domain/valueobjects"
 "time"
)

// DomainEvent 领域事件接口
type DomainEvent interface {
 EventName() string
 OccurredAt() time.Time
}

// ProductCreatedEvent 产品创建事件
type ProductCreatedEvent struct {
 ProductID   string
 Name        string
 Price       valueobjects.Money
 OccurredOn  time.Time
}

func (e ProductCreatedEvent) EventName() string {
 return "ProductCreated"
}

func (e ProductCreatedEvent) OccurredAt() time.Time {
 return e.OccurredOn
}

// StockUpdatedEvent 库存更新事件
type StockUpdatedEvent struct {
 ProductID  string
 OldStock   int
 NewStock   int
 OccurredOn time.Time
}

func (e StockUpdatedEvent) EventName() string {
 return "StockUpdated"
}

func (e StockUpdatedEvent) OccurredAt() time.Time {
 return e.OccurredOn
}

// PriceChangedEvent 价格变更事件
type PriceChangedEvent struct {
 ProductID  string
 OldPrice   valueobjects.Money
 NewPrice   valueobjects.Money
 OccurredOn time.Time
}

func (e PriceChangedEvent) EventName() string {
 return "PriceChanged"
}

func (e PriceChangedEvent) OccurredAt() time.Time {
 return e.OccurredOn
}
```

```go
// internal/domain/services/inventory_service.go
package services

import (
 "errors"
 "onion_architecture/internal/domain/entities"
)

// InventoryService 库存领域服务
type InventoryService struct{}

// NewInventoryService 创建库存服务
func NewInventoryService() *InventoryService {
 return &InventoryService{}
}

// TransferStock 库存转移（跨实体的业务逻辑）
func (s *InventoryService) TransferStock(
 from *entities.Product,
 to *entities.Product,
 quantity int,
) error {
 if from.ID == to.ID {
  return errors.New("cannot transfer to the same product")
 }
 if quantity <= 0 {
  return errors.New("quantity must be positive")
 }

 // 检查源产品库存
 if !from.IsAvailable(quantity) {
  return errors.New("insufficient stock in source product")
 }

 // 执行转移
 if err := from.ReserveStock(quantity); err != nil {
  return err
 }
 if err := to.UpdateStock(quantity); err != nil {
  // 回滚
  from.UpdateStock(quantity)
  return err
 }

 return nil
}

// CheckBulkAvailability 批量检查可用性
func (s *InventoryService) CheckBulkAvailability(items map[string]int) error {
 for productID, quantity := range items {
  if quantity <= 0 {
   return errors.New("invalid quantity for product: " + productID)
  }
 }
 return nil
}
```

```go
// internal/application/interfaces/repositories/product_repository.go
package repositories

import (
 "context"
 "onion_architecture/internal/domain/entities"
)

// ProductRepository 产品仓储接口
type ProductRepository interface {
 FindByID(ctx context.Context, id string) (*entities.Product, error)
 FindByCategory(ctx context.Context, categoryID string) ([]*entities.Product, error)
 FindAll(ctx context.Context, offset, limit int) ([]*entities.Product, error)
 Save(ctx context.Context, product *entities.Product) error
 Update(ctx context.Context, product *entities.Product) error
 Delete(ctx context.Context, id string) error
}
```

```go
// internal/application/interfaces/services/notification_service.go
package services

import (
 "context"
 "onion_architecture/internal/domain/events"
)

// NotificationService 通知服务接口
type NotificationService interface {
 NotifyProductCreated(ctx context.Context, event events.ProductCreatedEvent) error
 NotifyStockUpdated(ctx context.Context, event events.StockUpdatedEvent) error
 NotifyPriceChanged(ctx context.Context, event events.PriceChangedEvent) error
}
```

```go
// internal/application/dto/product_dto.go
package dto

import (
 "onion_architecture/internal/domain/valueobjects"
)

// CreateProductRequest 创建产品请求
type CreateProductRequest struct {
 Name        string  `json:"name"`
 Description string  `json:"description"`
 Price       float64 `json:"price"`
 Currency    string  `json:"currency"`
 CategoryID  string  `json:"category_id"`
}

// UpdateProductRequest 更新产品请求
type UpdateProductRequest struct {
 Name        string  `json:"name,omitempty"`
 Description string  `json:"description,omitempty"`
 Price       float64 `json:"price,omitempty"`
 Currency    string  `json:"currency,omitempty"`
}

// UpdateStockRequest 更新库存请求
type UpdateStockRequest struct {
 Quantity int `json:"quantity"`
}

// ProductResponse 产品响应
type ProductResponse struct {
 ID          string  `json:"id"`
 Name        string  `json:"name"`
 Description string  `json:"description"`
 Price       string  `json:"price"`
 Stock       int     `json:"stock"`
 CategoryID  string  `json:"category_id"`
}

// ProductListResponse 产品列表响应
type ProductListResponse struct {
 Products []ProductResponse `json:"products"`
 Total    int               `json:"total"`
}

// ToMoney 转换为Money值对象
func (r CreateProductRequest) ToMoney() (valueobjects.Money, error) {
 return valueobjects.NewMoney(r.Price, r.Currency)
}
```

```go
// internal/application/services/product_app_service.go
package services

import (
 "context"
 "onion_architecture/internal/application/dto"
 "onion_architecture/internal/application/interfaces/repositories"
 appServices "onion_architecture/internal/application/interfaces/services"
 "onion_architecture/internal/domain/entities"
 "onion_architecture/internal/domain/events"
 "time"

 "github.com/google/uuid"
)

// ProductApplicationService 产品应用服务
type ProductApplicationService struct {
 productRepo      repositories.ProductRepository
 notificationSvc  appServices.NotificationService
}

// NewProductApplicationService 创建应用服务
func NewProductApplicationService(
 productRepo repositories.ProductRepository,
 notificationSvc appServices.NotificationService,
) *ProductApplicationService {
 return &ProductApplicationService{
  productRepo:     productRepo,
  notificationSvc: notificationSvc,
 }
}

// CreateProduct 创建产品
func (s *ProductApplicationService) CreateProduct(
 ctx context.Context,
 req dto.CreateProductRequest,
) (*dto.ProductResponse, error) {
 price, err := req.ToMoney()
 if err != nil {
  return nil, err
 }

 product, err := entities.NewProduct(
  uuid.New().String(),
  req.Name,
  req.Description,
  price,
  req.CategoryID,
 )
 if err != nil {
  return nil, err
 }

 if err := s.productRepo.Save(ctx, product); err != nil {
  return nil, err
 }

 // 发布事件
 event := events.ProductCreatedEvent{
  ProductID:  product.ID,
  Name:       product.Name,
  Price:      product.Price,
  OccurredOn: time.Now(),
 }
 _ = s.notificationSvc.NotifyProductCreated(ctx, event)

 return s.toResponse(product), nil
}

// GetProduct 获取产品
func (s *ProductApplicationService) GetProduct(
 ctx context.Context,
 id string,
) (*dto.ProductResponse, error) {
 product, err := s.productRepo.FindByID(ctx, id)
 if err != nil {
  return nil, err
 }
 if product == nil {
  return nil, nil
 }
 return s.toResponse(product), nil
}

// UpdateStock 更新库存
func (s *ProductApplicationService) UpdateStock(
 ctx context.Context,
 id string,
 req dto.UpdateStockRequest,
) (*dto.ProductResponse, error) {
 product, err := s.productRepo.FindByID(ctx, id)
 if err != nil {
  return nil, err
 }
 if product == nil {
  return nil, nil
 }

 oldStock := product.Stock
 if err := product.UpdateStock(req.Quantity); err != nil {
  return nil, err
 }

 if err := s.productRepo.Update(ctx, product); err != nil {
  return nil, err
 }

 // 发布事件
 event := events.StockUpdatedEvent{
  ProductID:  product.ID,
  OldStock:   oldStock,
  NewStock:   product.Stock,
  OccurredOn: time.Now(),
 }
 _ = s.notificationSvc.NotifyStockUpdated(ctx, event)

 return s.toResponse(product), nil
}

// ListProducts 列出产品
func (s *ProductApplicationService) ListProducts(
 ctx context.Context,
 offset, limit int,
) (*dto.ProductListResponse, error) {
 products, err := s.productRepo.FindAll(ctx, offset, limit)
 if err != nil {
  return nil, err
 }

 var responses []dto.ProductResponse
 for _, p := range products {
  responses = append(responses, *s.toResponse(p))
 }

 return &dto.ProductListResponse{
  Products: responses,
  Total:    len(responses),
 }, nil
}

func (s *ProductApplicationService) toResponse(p *entities.Product) *dto.ProductResponse {
 return &dto.ProductResponse{
  ID:          p.ID,
  Name:        p.Name,
  Description: p.Description,
  Price:       p.Price.String(),
  Stock:       p.Stock,
  CategoryID:  p.CategoryID,
 }
}
```

```go
// internal/infrastructure/persistence/product_repository_impl.go
package persistence

import (
 "context"
 "onion_architecture/internal/application/interfaces/repositories"
 "onion_architecture/internal/domain/entities"
 "onion_architecture/internal/domain/valueobjects"
 "sync"
)

// InMemoryProductRepository 内存产品仓储
type InMemoryProductRepository struct {
 mu       sync.RWMutex
 products map[string]*entities.Product
}

// NewInMemoryProductRepository 创建仓储
func NewInMemoryProductRepository() repositories.ProductRepository {
 return &InMemoryProductRepository{
  products: make(map[string]*entities.Product),
 }
}

func (r *InMemoryProductRepository) FindByID(ctx context.Context, id string) (*entities.Product, error) {
 r.mu.RLock()
 defer r.mu.RUnlock()

 product, ok := r.products[id]
 if !ok {
  return nil, nil
 }
 return copyProduct(product), nil
}

func (r *InMemoryProductRepository) FindByCategory(ctx context.Context, categoryID string) ([]*entities.Product, error) {
 r.mu.RLock()
 defer r.mu.RUnlock()

 var result []*entities.Product
 for _, p := range r.products {
  if p.CategoryID == categoryID {
   result = append(result, copyProduct(p))
  }
 }
 return result, nil
}

func (r *InMemoryProductRepository) FindAll(ctx context.Context, offset, limit int) ([]*entities.Product, error) {
 r.mu.RLock()
 defer r.mu.RUnlock()

 var result []*entities.Product
 count := 0
 for _, p := range r.products {
  if count >= offset && len(result) < limit {
   result = append(result, copyProduct(p))
  }
  count++
 }
 return result, nil
}

func (r *InMemoryProductRepository) Save(ctx context.Context, product *entities.Product) error {
 r.mu.Lock()
 defer r.mu.Unlock()

 r.products[product.ID] = copyProduct(product)
 return nil
}

func (r *InMemoryProductRepository) Update(ctx context.Context, product *entities.Product) error {
 r.mu.Lock()
 defer r.mu.Unlock()

 r.products[product.ID] = copyProduct(product)
 return nil
}

func (r *InMemoryProductRepository) Delete(ctx context.Context, id string) error {
 r.mu.Lock()
 defer r.mu.Unlock()

 delete(r.products, id)
 return nil
}

func copyProduct(p *entities.Product) *entities.Product {
 return &entities.Product{
  ID:          p.ID,
  Name:        p.Name,
  Description: p.Description,
  Price:       p.Price,
  Stock:       p.Stock,
  CategoryID:  p.CategoryID,
  CreatedAt:   p.CreatedAt,
  UpdatedAt:   p.UpdatedAt,
 }
}
```

```go
// internal/infrastructure/web/product_controller.go
package web

import (
 "encoding/json"
 "net/http"
 "onion_architecture/internal/application/dto"
 "onion_architecture/internal/application/services"
 "strconv"
 "strings"
)

// ProductController 产品控制器
type ProductController struct {
 appService *services.ProductApplicationService
}

// NewProductController 创建控制器
func NewProductController(appService *services.ProductApplicationService) *ProductController {
 return &ProductController{appService: appService}
}

// ServeHTTP 实现http.Handler
func (c *ProductController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
 switch r.Method {
 case http.MethodPost:
  c.createProduct(w, r)
 case http.MethodGet:
  if strings.Contains(r.URL.Path, "/stock") {
   c.updateStock(w, r)
  } else if strings.TrimPrefix(r.URL.Path, "/products/") != "" {
   c.getProduct(w, r)
  } else {
   c.listProducts(w, r)
  }
 default:
  w.WriteHeader(http.StatusMethodNotAllowed)
 }
}

func (c *ProductController) createProduct(w http.ResponseWriter, r *http.Request) {
 var req dto.CreateProductRequest
 if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
  c.respondError(w, http.StatusBadRequest, err.Error())
  return
 }

 product, err := c.appService.CreateProduct(r.Context(), req)
 if err != nil {
  c.respondError(w, http.StatusBadRequest, err.Error())
  return
 }

 c.respondJSON(w, http.StatusCreated, product)
}

func (c *ProductController) getProduct(w http.ResponseWriter, r *http.Request) {
 id := strings.TrimPrefix(r.URL.Path, "/products/")
 product, err := c.appService.GetProduct(r.Context(), id)
 if err != nil {
  c.respondError(w, http.StatusInternalServerError, err.Error())
  return
 }
 if product == nil {
  c.respondError(w, http.StatusNotFound, "product not found")
  return
 }

 c.respondJSON(w, http.StatusOK, product)
}

func (c *ProductController) listProducts(w http.ResponseWriter, r *http.Request) {
 offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
 limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
 if limit <= 0 {
  limit = 10
 }

 products, err := c.appService.ListProducts(r.Context(), offset, limit)
 if err != nil {
  c.respondError(w, http.StatusInternalServerError, err.Error())
  return
 }

 c.respondJSON(w, http.StatusOK, products)
}

func (c *ProductController) updateStock(w http.ResponseWriter, r *http.Request) {
 id := strings.TrimPrefix(r.URL.Path, "/products/")
 id = strings.TrimSuffix(id, "/stock")

 var req dto.UpdateStockRequest
 if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
  c.respondError(w, http.StatusBadRequest, err.Error())
  return
 }

 product, err := c.appService.UpdateStock(r.Context(), id, req)
 if err != nil {
  c.respondError(w, http.StatusBadRequest, err.Error())
  return
 }
 if product == nil {
  c.respondError(w, http.StatusNotFound, "product not found")
  return
 }

 c.respondJSON(w, http.StatusOK, product)
}

func (c *ProductController) respondJSON(w http.ResponseWriter, status int, data interface{}) {
 w.Header().Set("Content-Type", "application/json")
 w.WriteHeader(status)
 json.NewEncoder(w).Encode(data)
}

func (c *ProductController) respondError(w http.ResponseWriter, status int, message string) {
 c.respondJSON(w, status, map[string]string{"error": message})
}
```

```go
// internal/infrastructure/external/email_service.go
package external

import (
 "context"
 "fmt"
 "onion_architecture/internal/application/interfaces/services"
 "onion_architecture/internal/domain/events"
)

// EmailService 邮件通知服务实现
type EmailService struct {
 smtpHost string
 smtpPort int
}

// NewEmailService 创建邮件服务
func NewEmailService(host string, port int) services.NotificationService {
 return &EmailService{smtpHost: host, smtpPort: port}
}

// NotifyProductCreated 产品创建通知
func (s *EmailService) NotifyProductCreated(ctx context.Context, event events.ProductCreatedEvent) error {
 fmt.Printf("[Email] Product created: %s (%s) at %s\n", event.Name, event.ProductID, event.OccurredAt)
 return nil
}

// NotifyStockUpdated 库存更新通知
func (s *EmailService) NotifyStockUpdated(ctx context.Context, event events.StockUpdatedEvent) error {
 fmt.Printf("[Email] Stock updated for product %s: %d -> %d\n",
  event.ProductID, event.OldStock, event.NewStock)
 return nil
}

// NotifyPriceChanged 价格变更通知
func (s *EmailService) NotifyPriceChanged(ctx context.Context, event events.PriceChangedEvent) error {
 fmt.Printf("[Email] Price changed for product %s: %s -> %s\n",
  event.ProductID, event.OldPrice, event.NewPrice)
 return nil
}
```

```go
// internal/config/dependency_injection.go
package config

import (
 "onion_architecture/internal/application/services"
 "onion_architecture/internal/infrastructure/external"
 "onion_architecture/internal/infrastructure/persistence"
 "onion_architecture/internal/infrastructure/web"
)

// Container 依赖注入容器
type Container struct {
 ProductController *web.ProductController
}

// NewContainer 创建容器
func NewContainer() *Container {
 // 基础设施层（最外层）
 productRepo := persistence.NewInMemoryProductRepository()
 notificationSvc := external.NewEmailService("smtp.example.com", 587)

 // 应用层
 productAppService := services.NewProductApplicationService(productRepo, notificationSvc)

 // 接口层
 productController := web.NewProductController(productAppService)

 return &Container{
  ProductController: productController,
 }
}
```

```go
// cmd/api/main.go
package main

import (
 "fmt"
 "net/http"
 "onion_architecture/internal/config"
)

func main() {
 // 初始化依赖注入容器
 container := config.NewContainer()

 // 注册路由
 http.Handle("/products/", container.ProductController)
 http.HandleFunc("/products", container.ProductController.ServeHTTP)

 fmt.Println("Onion Architecture Server starting on :8080")
 if err := http.ListenAndServe(":8080", nil); err != nil {
  panic(err)
 }
}
```

### 3.6 反例说明

#### 反例1：领域层依赖外层

```go
// ❌ 错误：领域实体依赖数据库
package domain

import "database/sql"

type Product struct {
 db *sql.DB  // 领域实体不应该知道数据库
}

func (p *Product) Save() error {
 _, err := p.db.Exec("INSERT INTO products ...")
 return err
}
```

**问题**：

- 领域层与基础设施强耦合
- 无法独立测试领域逻辑
- 违反洋葱架构核心原则

#### 反例2：应用服务包含领域逻辑

```go
// ❌ 错误：应用服务包含应该属于领域的逻辑
func (s *ProductAppService) CalculateDiscount(product *Product, quantity int) float64 {
 // 这是领域逻辑，应该在领域服务或实体中
 if quantity >= 10 {
  return product.Price * 0.9
 }
 return product.Price
}
```

**问题**：

- 业务逻辑散落在应用层
- 无法复用领域逻辑
- 领域模型贫血化

#### 反例3：值对象可变

```go
// ❌ 错误：值对象允许修改
type Money struct {
 Amount   float64
 Currency string
}

func (m *Money) SetAmount(amount float64) {  // 值对象应该是不可变的
 m.Amount = amount
}
```

**问题**：

- 值对象语义被破坏
- 可能导致数据不一致
- 违反值对象设计原则

### 3.7 选型指南

| 场景 | 建议 |
|------|------|
| **DDD项目** | 洋葱架构与DDD天然契合 |
| **复杂业务逻辑** | 领域为核心，清晰分离关注点 |
| **需要长期演进** | 内层稳定，外层易变 |
| **团队熟悉DDD** | 需要理解领域、应用、基础设施分层 |
| **简单CRUD** | 可能过度设计 |

### 3.8 优缺点分析

| 优点 | 缺点 |
|------|------|
| ✅ 领域模型完全独立 | ❌ 学习曲线较陡 |
| ✅ 清晰的层次结构 | ❌ 需要更多接口定义 |
| ✅ 高度可测试 | ❌ 初始设计成本较高 |
| ✅ 支持渐进式演进 | ❌ 简单项目可能过度设计 |
| ✅ 与DDD完美结合 | ❌ 需要团队理解DDD概念 |

---

## 4. 清洁架构（Clean Architecture）

### 4.1 概念定义

清洁架构（Clean Architecture）由Robert C. Martin（Uncle Bob）于2012年提出，是分层架构的进一步发展。它强调**依赖关系只能向内指向高层策略**，外层（框架、UI、数据库）依赖于内层（用例、实体），而内层不依赖于外层。

清洁架构的目标是**独立于框架、独立于UI、独立于数据库、独立于任何外部机构**，使系统更易于测试、维护和演进。

### 4.2 架构结构

```
                    ┌─────────────────────────────────────┐
                    │                                     │
                    │         框架与驱动层                 │  ← 最外层
                    │    Frameworks & Drivers             │     Web框架、数据库、外部工具
                    │                                     │
                    ├─────────────────────────────────────┤
                    │                                     │
                    │         接口适配器层                 │  ← 转换数据格式
                    │   Interface Adapters                │     Controllers、Presenters、Gateways
                    │                                     │
                    ├─────────────────────────────────────┤
                    │                                     │
                    │         用例层                       │  ← 应用业务规则
                    │    Use Cases                        │     用例、Interactor
                    │                                     │
                    ├─────────────────────────────────────┤
                    │                                     │
                    │         实体层                       │  ← 核心业务规则
                    │    Entities                         │     企业级业务对象
                    │                                     │
                    └─────────────────────────────────────┘

                    ◄──────── 依赖方向向内 ───────────────►
```

### 4.3 核心原则

| 原则 | 说明 |
|------|------|
| **依赖规则** | 依赖关系只能向内，内层不依赖外层 |
| **实体独立** | 实体层包含最核心的业务规则 |
| **用例隔离** | 用例层包含应用特定的业务规则 |
| **接口适配** | 适配器层负责数据格式转换 |
| **框架独立** | 最外层包含框架和工具 |

### 4.4 Go实现示例

#### 项目结构

```
clean_architecture/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── entities/            # 实体层（最内层）
│   │   ├── task.go
│   │   └── user.go
│   ├── usecases/            # 用例层
│   │   ├── task_usecase.go
│   │   ├── interfaces/
│   │   │   ├── repository.go
│   │   │   └── presenter.go
│   │   └── dto/
│   │       └── task_dto.go
│   ├── interface_adapters/  # 接口适配器层
│   │   ├── controllers/
│   │   │   └── task_controller.go
│   │   ├── presenters/
│   │   │   └── task_presenter.go
│   │   └── gateways/
│   │       └── task_repository.go
│   └── frameworks/          # 框架与驱动层（最外层）
│       ├── web/
│       │   ├── router.go
│       │   └── middleware.go
│       └── database/
│           └── gorm_task_repository.go
└── go.mod
```

#### 完整代码实现

```go
// internal/entities/task.go
package entities

import (
 "errors"
 "time"
)

// TaskStatus 任务状态
type TaskStatus string

const (
 TaskStatusTodo       TaskStatus = "TODO"
 TaskStatusInProgress TaskStatus = "IN_PROGRESS"
 TaskStatusDone       TaskStatus = "DONE"
)

// Priority 优先级
type Priority int

const (
 PriorityLow    Priority = 1
 PriorityMedium Priority = 2
 PriorityHigh   Priority = 3
)

// Task 任务实体
type Task struct {
 ID          string
 Title       string
 Description string
 Status      TaskStatus
 Priority    Priority
 AssigneeID  string
 DueDate     *time.Time
 CreatedAt   time.Time
 UpdatedAt   time.Time
 CompletedAt *time.Time
}

// NewTask 创建新任务
func NewTask(id, title, description string, priority Priority, assigneeID string) (*Task, error) {
 if title == "" {
  return nil, errors.New("task title is required")
 }

 return &Task{
  ID:          id,
  Title:       title,
  Description: description,
  Status:      TaskStatusTodo,
  Priority:    priority,
  AssigneeID:  assigneeID,
  CreatedAt:   time.Now(),
  UpdatedAt:   time.Now(),
 }, nil
}

// Start 开始任务
func (t *Task) Start() error {
 if t.Status != TaskStatusTodo {
  return errors.New("only todo tasks can be started")
 }
 t.Status = TaskStatusInProgress
 t.UpdatedAt = time.Now()
 return nil
}

// Complete 完成任务
func (t *Task) Complete() error {
 if t.Status != TaskStatusInProgress {
  return errors.New("only in-progress tasks can be completed")
 }
 t.Status = TaskStatusDone
 now := time.Now()
 t.CompletedAt = &now
 t.UpdatedAt = now
 return nil
}

// UpdatePriority 更新优先级
func (t *Task) UpdatePriority(priority Priority) error {
 if priority < PriorityLow || priority > PriorityHigh {
  return errors.New("invalid priority")
 }
 t.Priority = priority
 t.UpdatedAt = time.Now()
 return nil
}

// IsOverdue 检查是否逾期
func (t *Task) IsOverdue() bool {
 if t.DueDate == nil || t.Status == TaskStatusDone {
  return false
 }
 return time.Now().After(*t.DueDate)
}

// CanBeAssignedTo 检查是否可以分配给指定用户
func (t *Task) CanBeAssignedTo(userID string) bool {
 return t.AssigneeID == "" || t.AssigneeID == userID
}
```

```go
// internal/usecases/interfaces/repository.go
package interfaces

import (
 "clean_architecture/internal/entities"
 "context"
)

// TaskRepository 任务仓储接口
type TaskRepository interface {
 FindByID(ctx context.Context, id string) (*entities.Task, error)
 FindByAssignee(ctx context.Context, assigneeID string) ([]*entities.Task, error)
 FindByStatus(ctx context.Context, status entities.TaskStatus) ([]*entities.Task, error)
 Save(ctx context.Context, task *entities.Task) error
 Update(ctx context.Context, task *entities.Task) error
 Delete(ctx context.Context, id string) error
}
```

```go
// internal/usecases/interfaces/presenter.go
package interfaces

import (
 "clean_architecture/internal/entities"
 "clean_architecture/internal/usecases/dto"
)

// TaskPresenter 任务展示器接口
type TaskPresenter interface {
 PresentTask(task *entities.Task) dto.TaskResponse
 PresentTasks(tasks []*entities.Task) dto.TaskListResponse
 PresentError(err error) dto.ErrorResponse
}
```

```go
// internal/usecases/dto/task_dto.go
package dto

import "time"

// CreateTaskInput 创建任务输入
type CreateTaskInput struct {
 Title       string    `json:"title"`
 Description string    `json:"description"`
 Priority    int       `json:"priority"`
 AssigneeID  string    `json:"assignee_id"`
 DueDate     time.Time `json:"due_date,omitempty"`
}

// UpdateTaskInput 更新任务输入
type UpdateTaskInput struct {
 Title       string `json:"title,omitempty"`
 Description string `json:"description,omitempty"`
 Priority    int    `json:"priority,omitempty"`
}

// TaskResponse 任务响应
type TaskResponse struct {
 ID          string     `json:"id"`
 Title       string     `json:"title"`
 Description string     `json:"description"`
 Status      string     `json:"status"`
 Priority    int        `json:"priority"`
 AssigneeID  string     `json:"assignee_id"`
 DueDate     *time.Time `json:"due_date,omitempty"`
 CreatedAt   time.Time  `json:"created_at"`
 UpdatedAt   time.Time  `json:"updated_at"`
 IsOverdue   bool       `json:"is_overdue"`
}

// TaskListResponse 任务列表响应
type TaskListResponse struct {
 Tasks []TaskResponse `json:"tasks"`
 Total int            `json:"total"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
 Message string `json:"message"`
 Code    string `json:"code,omitempty"`
}
```

```go
// internal/usecases/task_usecase.go
package usecases

import (
 "clean_architecture/internal/entities"
 "clean_architecture/internal/usecases/dto"
 "clean_architecture/internal/usecases/interfaces"
 "context"
 "errors"
 "time"

 "github.com/google/uuid"
)

// TaskUsecase 任务用例
type TaskUsecase struct {
 repo      interfaces.TaskRepository
 presenter interfaces.TaskPresenter
}

// NewTaskUsecase 创建任务用例
func NewTaskUsecase(repo interfaces.TaskRepository, presenter interfaces.TaskPresenter) *TaskUsecase {
 return &TaskUsecase{
  repo:      repo,
  presenter: presenter,
 }
}

// CreateTask 创建任务
func (u *TaskUsecase) CreateTask(ctx context.Context, input dto.CreateTaskInput) (dto.TaskResponse, error) {
 priority := entities.Priority(input.Priority)

 task, err := entities.NewTask(
  uuid.New().String(),
  input.Title,
  input.Description,
  priority,
  input.AssigneeID,
 )
 if err != nil {
  return u.presenter.PresentError(err), err
 }

 if !input.DueDate.IsZero() {
  task.DueDate = &input.DueDate
 }

 if err := u.repo.Save(ctx, task); err != nil {
  return u.presenter.PresentError(err), err
 }

 return u.presenter.PresentTask(task), nil
}

// GetTask 获取任务
func (u *TaskUsecase) GetTask(ctx context.Context, id string) (dto.TaskResponse, error) {
 task, err := u.repo.FindByID(ctx, id)
 if err != nil {
  return u.presenter.PresentError(err), err
 }
 if task == nil {
  return u.presenter.PresentError(errors.New("task not found")), errors.New("task not found")
 }

 return u.presenter.PresentTask(task), nil
}

// UpdateTask 更新任务
func (u *TaskUsecase) UpdateTask(ctx context.Context, id string, input dto.UpdateTaskInput) (dto.TaskResponse, error) {
 task, err := u.repo.FindByID(ctx, id)
 if err != nil {
  return u.presenter.PresentError(err), err
 }
 if task == nil {
  return u.presenter.PresentError(errors.New("task not found")), errors.New("task not found")
 }

 if input.Title != "" {
  task.Title = input.Title
 }
 if input.Description != "" {
  task.Description = input.Description
 }
 if input.Priority > 0 {
  if err := task.UpdatePriority(entities.Priority(input.Priority)); err != nil {
   return u.presenter.PresentError(err), err
  }
 }

 task.UpdatedAt = time.Now()

 if err := u.repo.Update(ctx, task); err != nil {
  return u.presenter.PresentError(err), err
 }

 return u.presenter.PresentTask(task), nil
}

// StartTask 开始任务
func (u *TaskUsecase) StartTask(ctx context.Context, id string) (dto.TaskResponse, error) {
 task, err := u.repo.FindByID(ctx, id)
 if err != nil {
  return u.presenter.PresentError(err), err
 }
 if task == nil {
  return u.presenter.PresentError(errors.New("task not found")), errors.New("task not found")
 }

 if err := task.Start(); err != nil {
  return u.presenter.PresentError(err), err
 }

 if err := u.repo.Update(ctx, task); err != nil {
  return u.presenter.PresentError(err), err
 }

 return u.presenter.PresentTask(task), nil
}

// CompleteTask 完成任务
func (u *TaskUsecase) CompleteTask(ctx context.Context, id string) (dto.TaskResponse, error) {
 task, err := u.repo.FindByID(ctx, id)
 if err != nil {
  return u.presenter.PresentError(err), err
 }
 if task == nil {
  return u.presenter.PresentError(errors.New("task not found")), errors.New("task not found")
 }

 if err := task.Complete(); err != nil {
  return u.presenter.PresentError(err), err
 }

 if err := u.repo.Update(ctx, task); err != nil {
  return u.presenter.PresentError(err), err
 }

 return u.presenter.PresentTask(task), nil
}

// ListTasksByAssignee 列出分配者的任务
func (u *TaskUsecase) ListTasksByAssignee(ctx context.Context, assigneeID string) (dto.TaskListResponse, error) {
 tasks, err := u.repo.FindByAssignee(ctx, assigneeID)
 if err != nil {
  return u.presenter.PresentTasks(nil), err
 }

 return u.presenter.PresentTasks(tasks), nil
}

// DeleteTask 删除任务
func (u *TaskUsecase) DeleteTask(ctx context.Context, id string) error {
 return u.repo.Delete(ctx, id)
}
```

```go
// internal/interface_adapters/controllers/task_controller.go
package controllers

import (
 "clean_architecture/internal/usecases"
 "clean_architecture/internal/usecases/dto"
 "encoding/json"
 "net/http"
 "strings"
)

// TaskController 任务控制器
type TaskController struct {
 usecase *usecases.TaskUsecase
}

// NewTaskController 创建控制器
func NewTaskController(usecase *usecases.TaskUsecase) *TaskController {
 return &TaskController{usecase: usecase}
}

// ServeHTTP 实现http.Handler
func (c *TaskController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
 switch r.Method {
 case http.MethodPost:
  c.createTask(w, r)
 case http.MethodGet:
  if strings.HasSuffix(r.URL.Path, "/start") {
   c.startTask(w, r)
  } else if strings.HasSuffix(r.URL.Path, "/complete") {
   c.completeTask(w, r)
  } else if strings.Contains(r.URL.Path, "/tasks/") {
   c.getTask(w, r)
  } else {
   c.listTasks(w, r)
  }
 case http.MethodPut:
  c.updateTask(w, r)
 case http.MethodDelete:
  c.deleteTask(w, r)
 default:
  w.WriteHeader(http.StatusMethodNotAllowed)
 }
}

func (c *TaskController) createTask(w http.ResponseWriter, r *http.Request) {
 var input dto.CreateTaskInput
 if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
  c.respondError(w, http.StatusBadRequest, err.Error())
  return
 }

 response, err := c.usecase.CreateTask(r.Context(), input)
 if err != nil {
  c.respondError(w, http.StatusBadRequest, err.Error())
  return
 }

 c.respondJSON(w, http.StatusCreated, response)
}

func (c *TaskController) getTask(w http.ResponseWriter, r *http.Request) {
 id := strings.TrimPrefix(r.URL.Path, "/tasks/")
 response, err := c.usecase.GetTask(r.Context(), id)
 if err != nil {
  c.respondError(w, http.StatusNotFound, err.Error())
  return
 }

 c.respondJSON(w, http.StatusOK, response)
}

func (c *TaskController) updateTask(w http.ResponseWriter, r *http.Request) {
 id := strings.TrimPrefix(r.URL.Path, "/tasks/")
 var input dto.UpdateTaskInput
 if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
  c.respondError(w, http.StatusBadRequest, err.Error())
  return
 }

 response, err := c.usecase.UpdateTask(r.Context(), id, input)
 if err != nil {
  c.respondError(w, http.StatusBadRequest, err.Error())
  return
 }

 c.respondJSON(w, http.StatusOK, response)
}

func (c *TaskController) startTask(w http.ResponseWriter, r *http.Request) {
 id := strings.TrimPrefix(r.URL.Path, "/tasks/")
 id = strings.TrimSuffix(id, "/start")

 response, err := c.usecase.StartTask(r.Context(), id)
 if err != nil {
  c.respondError(w, http.StatusBadRequest, err.Error())
  return
 }

 c.respondJSON(w, http.StatusOK, response)
}

func (c *TaskController) completeTask(w http.ResponseWriter, r *http.Request) {
 id := strings.TrimPrefix(r.URL.Path, "/tasks/")
 id = strings.TrimSuffix(id, "/complete")

 response, err := c.usecase.CompleteTask(r.Context(), id)
 if err != nil {
  c.respondError(w, http.StatusBadRequest, err.Error())
  return
 }

 c.respondJSON(w, http.StatusOK, response)
}

func (c *TaskController) listTasks(w http.ResponseWriter, r *http.Request) {
 assigneeID := r.URL.Query().Get("assignee_id")
 if assigneeID == "" {
  c.respondError(w, http.StatusBadRequest, "assignee_id is required")
  return
 }

 response, err := c.usecase.ListTasksByAssignee(r.Context(), assigneeID)
 if err != nil {
  c.respondError(w, http.StatusInternalServerError, err.Error())
  return
 }

 c.respondJSON(w, http.StatusOK, response)
}

func (c *TaskController) deleteTask(w http.ResponseWriter, r *http.Request) {
 id := strings.TrimPrefix(r.URL.Path, "/tasks/")
 if err := c.usecase.DeleteTask(r.Context(), id); err != nil {
  c.respondError(w, http.StatusInternalServerError, err.Error())
  return
 }

 w.WriteHeader(http.StatusNoContent)
}

func (c *TaskController) respondJSON(w http.ResponseWriter, status int, data interface{}) {
 w.Header().Set("Content-Type", "application/json")
 w.WriteHeader(status)
 json.NewEncoder(w).Encode(data)
}

func (c *TaskController) respondError(w http.ResponseWriter, status int, message string) {
 c.respondJSON(w, status, dto.ErrorResponse{Message: message})
}
```

```go
// internal/interface_adapters/presenters/task_presenter.go
package presenters

import (
 "clean_architecture/internal/entities"
 "clean_architecture/internal/usecases/dto"
 "clean_architecture/internal/usecases/interfaces"
)

// TaskPresenter 任务展示器实现
type TaskPresenter struct{}

// NewTaskPresenter 创建展示器
func NewTaskPresenter() interfaces.TaskPresenter {
 return &TaskPresenter{}
}

// PresentTask 展示单个任务
func (p *TaskPresenter) PresentTask(task *entities.Task) dto.TaskResponse {
 return dto.TaskResponse{
  ID:          task.ID,
  Title:       task.Title,
  Description: task.Description,
  Status:      string(task.Status),
  Priority:    int(task.Priority),
  AssigneeID:  task.AssigneeID,
  DueDate:     task.DueDate,
  CreatedAt:   task.CreatedAt,
  UpdatedAt:   task.UpdatedAt,
  IsOverdue:   task.IsOverdue(),
 }
}

// PresentTasks 展示任务列表
func (p *TaskPresenter) PresentTasks(tasks []*entities.Task) dto.TaskListResponse {
 var responses []dto.TaskResponse
 for _, task := range tasks {
  responses = append(responses, p.PresentTask(task))
 }

 return dto.TaskListResponse{
  Tasks: responses,
  Total: len(responses),
 }
}

// PresentError 展示错误
func (p *TaskPresenter) PresentError(err error) dto.ErrorResponse {
 return dto.ErrorResponse{
  Message: err.Error(),
 }
}
```

```go
// internal/frameworks/database/gorm_task_repository.go
package database

import (
 "clean_architecture/internal/entities"
 "clean_architecture/internal/usecases/interfaces"
 "context"
 "sync"
)

// GormTaskRepository GORM任务仓储实现
type GormTaskRepository struct {
 mu    sync.RWMutex
 tasks map[string]*entities.Task
}

// NewGormTaskRepository 创建仓储
func NewGormTaskRepository() interfaces.TaskRepository {
 return &GormTaskRepository{
  tasks: make(map[string]*entities.Task),
 }
}

func (r *GormTaskRepository) FindByID(ctx context.Context, id string) (*entities.Task, error) {
 r.mu.RLock()
 defer r.mu.RUnlock()

 task, ok := r.tasks[id]
 if !ok {
  return nil, nil
 }
 return copyTask(task), nil
}

func (r *GormTaskRepository) FindByAssignee(ctx context.Context, assigneeID string) ([]*entities.Task, error) {
 r.mu.RLock()
 defer r.mu.RUnlock()

 var result []*entities.Task
 for _, task := range r.tasks {
  if task.AssigneeID == assigneeID {
   result = append(result, copyTask(task))
  }
 }
 return result, nil
}

func (r *GormTaskRepository) FindByStatus(ctx context.Context, status entities.TaskStatus) ([]*entities.Task, error) {
 r.mu.RLock()
 defer r.mu.RUnlock()

 var result []*entities.Task
 for _, task := range r.tasks {
  if task.Status == status {
   result = append(result, copyTask(task))
  }
 }
 return result, nil
}

func (r *GormTaskRepository) Save(ctx context.Context, task *entities.Task) error {
 r.mu.Lock()
 defer r.mu.Unlock()

 r.tasks[task.ID] = copyTask(task)
 return nil
}

func (r *GormTaskRepository) Update(ctx context.Context, task *entities.Task) error {
 r.mu.Lock()
 defer r.mu.Unlock()

 r.tasks[task.ID] = copyTask(task)
 return nil
}

func (r *GormTaskRepository) Delete(ctx context.Context, id string) error {
 r.mu.Lock()
 defer r.mu.Unlock()

 delete(r.tasks, id)
 return nil
}

func copyTask(t *entities.Task) *entities.Task {
 return &entities.Task{
  ID:          t.ID,
  Title:       t.Title,
  Description: t.Description,
  Status:      t.Status,
  Priority:    t.Priority,
  AssigneeID:  t.AssigneeID,
  DueDate:     t.DueDate,
  CreatedAt:   t.CreatedAt,
  UpdatedAt:   t.UpdatedAt,
  CompletedAt: t.CompletedAt,
 }
}
```

```go
// internal/frameworks/web/router.go
package web

import (
 "clean_architecture/internal/interface_adapters/controllers"
 "net/http"
)

// Router 路由器
type Router struct {
 mux            *http.ServeMux
 taskController *controllers.TaskController
}

// NewRouter 创建路由器
func NewRouter(taskController *controllers.TaskController) *Router {
 r := &Router{
  mux:            http.NewServeMux(),
  taskController: taskController,
 }
 r.setupRoutes()
 return r
}

func (r *Router) setupRoutes() {
 r.mux.Handle("/tasks/", r.taskController)
 r.mux.HandleFunc("/tasks", r.taskController.ServeHTTP)
}

// ServeHTTP 实现http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
 r.mux.ServeHTTP(w, req)
}
```

```go
// cmd/api/main.go
package main

import (
 "clean_architecture/internal/frameworks/database"
 "clean_architecture/internal/frameworks/web"
 "clean_architecture/internal/interface_adapters/controllers"
 "clean_architecture/internal/interface_adapters/presenters"
 "clean_architecture/internal/usecases"
 "fmt"
 "net/http"
)

func main() {
 // 框架层：数据库
 taskRepo := database.NewGormTaskRepository()

 // 接口适配层：展示器
 taskPresenter := presenters.NewTaskPresenter()

 // 用例层
 taskUsecase := usecases.NewTaskUsecase(taskRepo, taskPresenter)

 // 接口适配层：控制器
 taskController := controllers.NewTaskController(taskUsecase)

 // 框架层：路由器
 router := web.NewRouter(taskController)

 fmt.Println("Clean Architecture Server starting on :8080")
 if err := http.ListenAndServe(":8080", router); err != nil {
  panic(err)
 }
}
```

### 4.5 反例说明

#### 反例1：用例依赖框架

```go
// ❌ 错误：用例层依赖Gin框架
package usecases

import "github.com/gin-gonic/gin"

type TaskUsecase struct {
 // ...
}

func (u *TaskUsecase) CreateTask(c *gin.Context) {
 // 用例不应该知道HTTP框架
}
```

**问题**：

- 用例层与框架耦合
- 无法独立测试业务逻辑
- 违反依赖规则

#### 反例2：实体依赖数据库

```go
// ❌ 错误：实体直接操作数据库
package entities

import "database/sql"

type Task struct {
 ID   string
 DB   *sql.DB  // 实体不应该有数据库连接
}

func (t *Task) Save() error {
 _, err := t.DB.Exec("INSERT INTO tasks ...")
 return err
}
```

**问题**：

- 实体与基础设施耦合
- 违反清洁架构核心原则
- 无法独立测试实体

### 4.6 选型指南

| 场景 | 建议 |
|------|------|
| **大型企业应用** | 清洁架构提供清晰的层次结构 |
| **需要长期维护** | 依赖规则确保内层稳定 |
| **多团队协作** | 清晰的边界便于分工 |
| **需要高可测试性** | 每层可独立测试 |
| **小型项目** | 可能过度设计 |

### 4.7 优缺点分析

| 优点 | 缺点 |
|------|------|
| ✅ 完全独立于框架 | ❌ 初始设计成本较高 |
| ✅ 高度可测试 | ❌ 需要更多接口和适配器 |
| ✅ 清晰的依赖规则 | ❌ 学习曲线较陡 |
| ✅ 易于维护和演进 | ❌ 简单项目可能过度设计 |
| ✅ 支持多种接口 | ❌ 需要团队理解和遵守规则 |

---

## 5. CQRS（命令查询职责分离）

### 5.1 概念定义

CQRS（Command Query Responsibility Segregation，命令查询职责分离）是一种架构模式，它将应用程序的读操作（Query）和写操作（Command）分离到不同的模型中。

传统的CRUD模式使用同一个模型处理读写操作，而CQRS认为读写操作有不同的需求和特性，应该使用不同的模型来优化各自的处理。

### 5.2 架构结构

```
┌─────────────────────────────────────────────────────────────┐
│                        客户端                                │
└───────────────────────┬─────────────────────────────────────┘
                        │
        ┌───────────────┴───────────────┐
        │                               │
┌───────▼────────┐            ┌────────▼────────┐
│   命令端        │            │    查询端        │
│  Command Side   │            │   Query Side    │
│                 │            │                 │
│ ┌───────────┐   │            │   ┌───────────┐ │
│ │  Command  │   │            │   │   Query   │ │
│ │  Handler  │   │            │   │  Handler  │ │
│ └─────┬─────┘   │            │   └─────┬─────┘ │
│       │         │            │         │       │
│ ┌─────▼─────┐   │            │   ┌─────▼─────┐ │
│ │  Command  │   │            │   │  Query    │ │
│ │  Service  │   │            │   │  Service  │ │
│ └─────┬─────┘   │            │   └─────┬─────┘ │
│       │         │            │         │       │
│ ┌─────▼─────┐   │            │   ┌─────▼─────┐ │
│ │  Write    │   │            │   │   Read    │ │
│ │  Model    │   │            │   │   Model   │ │
│ │           │   │            │   │           │ │
│ │  Domain   │   │            │   │  DTO/     │ │
│ │  Entity   │   │            │   │  View     │ │
│ └─────┬─────┘   │            │   └─────┬─────┘ │
│       │         │            │         │       │
│ ┌─────▼─────┐   │            │   ┌─────▼─────┐ │
│ │  Write    │   │            │   │   Read    │ │
│ │  DB       │   │            │   │   DB      │ │
│ └───────────┘   │            │   └───────────┘ │
└─────────────────┘            └─────────────────┘
        │                               │
        │         事件同步               │
        └──────────────┬────────────────┘
                       │
              ┌────────▼────────┐
              │   Event Bus/    │
              │   Message Queue │
              └─────────────────┘
```

### 5.3 核心原则

| 原则 | 说明 |
|------|------|
| **职责分离** | 命令和查询使用不同的模型和处理流程 |
| **模型优化** | 写模型优化一致性，读模型优化查询性能 |
| **事件同步** | 通过事件机制保持读写模型最终一致 |
| **独立扩展** | 读写端可独立扩展和优化 |

### 5.4 Go实现示例

#### 项目结构

```
cqrs_architecture/
├── cmd/
│   ├── command/
│   │   └── main.go
│   └── query/
│       └── main.go
├── internal/
│   ├── shared/              # 共享组件
│   │   ├── events/
│   │   │   └── order_events.go
│   │   └── messaging/
│   │       └── event_bus.go
│   ├── command/             # 命令端
│   │   ├── handlers/
│   │   │   └── order_command_handler.go
│   │   ├── services/
│   │   │   └── order_command_service.go
│   │   ├── domain/
│   │   │   ├── order.go
│   │   │   └── order_repository.go
│   │   └── infrastructure/
│   │       └── order_repository_impl.go
│   └── query/               # 查询端
│       ├── handlers/
│       │   └── order_query_handler.go
│       ├── services/
│       │   └── order_query_service.go
│       ├── projections/
│       │   └── order_projection.go
│       └── infrastructure/
│           └── order_read_repository.go
└── go.mod
```

#### 完整代码实现

```go
// internal/shared/events/order_events.go
package events

import (
 "time"
)

// Event 事件接口
type Event interface {
 EventType() string
 AggregateID() string
 OccurredAt() time.Time
}

// OrderCreatedEvent 订单创建事件
type OrderCreatedEvent struct {
 OrderID    string
 CustomerID string
 Items      []OrderItem
 Total      float64
 Timestamp  time.Time
}

func (e OrderCreatedEvent) EventType() string    { return "OrderCreated" }
func (e OrderCreatedEvent) AggregateID() string  { return e.OrderID }
func (e OrderCreatedEvent) OccurredAt() time.Time { return e.Timestamp }

// OrderItem 订单项
type OrderItem struct {
 ProductID string
 Name      string
 Quantity  int
 Price     float64
}

// OrderStatusChangedEvent 订单状态变更事件
type OrderStatusChangedEvent struct {
 OrderID   string
 OldStatus string
 NewStatus string
 Timestamp time.Time
}

func (e OrderStatusChangedEvent) EventType() string    { return "OrderStatusChanged" }
func (e OrderStatusChangedEvent) AggregateID() string  { return e.OrderID }
func (e OrderStatusChangedEvent) OccurredAt() time.Time { return e.Timestamp }

// OrderCancelledEvent 订单取消事件
type OrderCancelledEvent struct {
 OrderID   string
 Reason    string
 Timestamp time.Time
}

func (e OrderCancelledEvent) EventType() string    { return "OrderCancelled" }
func (e OrderCancelledEvent) AggregateID() string  { return e.OrderID }
func (e OrderCancelledEvent) OccurredAt() time.Time { return e.Timestamp }
```

```go
// internal/shared/messaging/event_bus.go
package messaging

import (
 "cqrs_architecture/internal/shared/events"
 "sync"
)

// EventHandler 事件处理器类型
type EventHandler func(event events.Event)

// EventBus 事件总线接口
type EventBus interface {
 Publish(event events.Event) error
 Subscribe(eventType string, handler EventHandler) error
}

// InMemoryEventBus 内存事件总线实现
type InMemoryEventBus struct {
 mu        sync.RWMutex
 handlers  map[string][]EventHandler
 eventChan chan events.Event
}

// NewInMemoryEventBus 创建事件总线
func NewInMemoryEventBus() *InMemoryEventBus {
 bus := &InMemoryEventBus{
  handlers:  make(map[string][]EventHandler),
  eventChan: make(chan events.Event, 100),
 }
 go bus.processEvents()
 return bus
}

// Publish 发布事件
func (b *InMemoryEventBus) Publish(event events.Event) error {
 b.eventChan <- event
 return nil
}

// Subscribe 订阅事件
func (b *InMemoryEventBus) Subscribe(eventType string, handler EventHandler) error {
 b.mu.Lock()
 defer b.mu.Unlock()

 b.handlers[eventType] = append(b.handlers[eventType], handler)
 return nil
}

func (b *InMemoryEventBus) processEvents() {
 for event := range b.eventChan {
  b.mu.RLock()
  handlers := b.handlers[event.EventType()]
  b.mu.RUnlock()

  for _, handler := range handlers {
   go handler(event)
  }
 }
}
```

```go
// internal/command/domain/order.go
package domain

import (
 "cqrs_architecture/internal/shared/events"
 "errors"
 "time"
)

// OrderStatus 订单状态
type OrderStatus string

const (
 OrderStatusPending   OrderStatus = "PENDING"
 OrderStatusPaid      OrderStatus = "PAID"
 OrderStatusShipped   OrderStatus = "SHIPPED"
 OrderStatusCompleted OrderStatus = "COMPLETED"
 OrderStatusCancelled OrderStatus = "CANCELLED"
)

// Order 订单聚合根
type Order struct {
 ID         string
 CustomerID string
 Items      []OrderItem
 Status     OrderStatus
 Total      float64
 Version    int
 CreatedAt  time.Time
 UpdatedAt  time.Time
 events     []events.Event
}

// OrderItem 订单项
type OrderItem struct {
 ProductID string
 Name      string
 Quantity  int
 Price     float64
}

// NewOrder 创建新订单
func NewOrder(id, customerID string, items []OrderItem) (*Order, error) {
 if customerID == "" {
  return nil, errors.New("customer ID is required")
 }
 if len(items) == 0 {
  return nil, errors.New("order must have at least one item")
 }

 order := &Order{
  ID:         id,
  CustomerID: customerID,
  Items:      items,
  Status:     OrderStatusPending,
  Version:    1,
  CreatedAt:  time.Now(),
  UpdatedAt:  time.Now(),
  events:     make([]events.Event, 0),
 }

 order.calculateTotal()
 order.raiseEvent(events.OrderCreatedEvent{
  OrderID:    id,
  CustomerID: customerID,
  Items:      order.toEventItems(),
  Total:      order.Total,
  Timestamp:  time.Now(),
 })

 return order, nil
}

// Pay 支付订单
func (o *Order) Pay() error {
 if o.Status != OrderStatusPending {
  return errors.New("only pending orders can be paid")
 }

 oldStatus := o.Status
 o.Status = OrderStatusPaid
 o.Version++
 o.UpdatedAt = time.Now()

 o.raiseEvent(events.OrderStatusChangedEvent{
  OrderID:   o.ID,
  OldStatus: string(oldStatus),
  NewStatus: string(o.Status),
  Timestamp: time.Now(),
 })

 return nil
}

// Cancel 取消订单
func (o *Order) Cancel(reason string) error {
 if o.Status != OrderStatusPending && o.Status != OrderStatusPaid {
  return errors.New("cannot cancel order in current status")
 }

 o.Status = OrderStatusCancelled
 o.Version++
 o.UpdatedAt = time.Now()

 o.raiseEvent(events.OrderCancelledEvent{
  OrderID:   o.ID,
  Reason:    reason,
  Timestamp: time.Now(),
 })

 return nil
}

// GetEvents 获取未提交的事件
func (o *Order) GetEvents() []events.Event {
 return o.events
}

// ClearEvents 清除事件
func (o *Order) ClearEvents() {
 o.events = make([]events.Event, 0)
}

func (o *Order) calculateTotal() {
 total := 0.0
 for _, item := range o.Items {
  total += item.Price * float64(item.Quantity)
 }
 o.Total = total
}

func (o *Order) raiseEvent(event events.Event) {
 o.events = append(o.events, event)
}

func (o *Order) toEventItems() []events.OrderItem {
 result := make([]events.OrderItem, len(o.Items))
 for i, item := range o.Items {
  result[i] = events.OrderItem{
   ProductID: item.ProductID,
   Name:      item.Name,
   Quantity:  item.Quantity,
   Price:     item.Price,
  }
 }
 return result
}
```

```go
// internal/command/domain/order_repository.go
package domain

import (
 "context"
)

// OrderRepository 订单仓储接口
type OrderRepository interface {
 FindByID(ctx context.Context, id string) (*Order, error)
 Save(ctx context.Context, order *Order) error
 Update(ctx context.Context, order *Order) error
}
```

```go
// internal/command/infrastructure/order_repository_impl.go
package infrastructure

import (
 "command/internal/domain"
 "context"
 "sync"
)

// InMemoryOrderRepository 内存订单仓储
type InMemoryOrderRepository struct {
 mu     sync.RWMutex
 orders map[string]*domain.Order
}

// NewInMemoryOrderRepository 创建仓储
func NewInMemoryOrderRepository() domain.OrderRepository {
 return &InMemoryOrderRepository{
  orders: make(map[string]*domain.Order),
 }
}

func (r *InMemoryOrderRepository) FindByID(ctx context.Context, id string) (*domain.Order, error) {
 r.mu.RLock()
 defer r.mu.RUnlock()

 order, ok := r.orders[id]
 if !ok {
  return nil, nil
 }
 return copyOrder(order), nil
}

func (r *InMemoryOrderRepository) Save(ctx context.Context, order *domain.Order) error {
 r.mu.Lock()
 defer r.mu.Unlock()

 r.orders[order.ID] = copyOrder(order)
 return nil
}

func (r *InMemoryOrderRepository) Update(ctx context.Context, order *domain.Order) error {
 r.mu.Lock()
 defer r.mu.Unlock()

 r.orders[order.ID] = copyOrder(order)
 return nil
}

func copyOrder(o *domain.Order) *domain.Order {
 items := make([]domain.OrderItem, len(o.Items))
 copy(items, o.Items)
 return &domain.Order{
  ID:         o.ID,
  CustomerID: o.CustomerID,
  Items:      items,
  Status:     o.Status,
  Total:      o.Total,
  Version:    o.Version,
  CreatedAt:  o.CreatedAt,
  UpdatedAt:  o.UpdatedAt,
  events:     make([]interface{}, 0),
 }
}
```

```go
// internal/command/services/order_command_service.go
package services

import (
 "command/internal/domain"
 "command/internal/shared/events"
 "command/internal/shared/messaging"
 "context"
 "errors"
 "time"

 "github.com/google/uuid"
)

// CreateOrderCommand 创建订单命令
type CreateOrderCommand struct {
 CustomerID string
 Items      []domain.OrderItem
}

// PayOrderCommand 支付订单命令
type PayOrderCommand struct {
 OrderID string
}

// CancelOrderCommand 取消订单命令
type CancelOrderCommand struct {
 OrderID string
 Reason  string
}

// OrderCommandService 订单命令服务
type OrderCommandService struct {
 orderRepo domain.OrderRepository
 eventBus  messaging.EventBus
}

// NewOrderCommandService 创建命令服务
func NewOrderCommandService(orderRepo domain.OrderRepository, eventBus messaging.EventBus) *OrderCommandService {
 return &OrderCommandService{
  orderRepo: orderRepo,
  eventBus:  eventBus,
 }
}

// CreateOrder 处理创建订单命令
func (s *OrderCommandService) CreateOrder(ctx context.Context, cmd CreateOrderCommand) (string, error) {
 order, err := domain.NewOrder(uuid.New().String(), cmd.CustomerID, cmd.Items)
 if err != nil {
  return "", err
 }

 if err := s.orderRepo.Save(ctx, order); err != nil {
  return "", err
 }

 // 发布事件
 for _, event := range order.GetEvents() {
  s.eventBus.Publish(event)
 }
 order.ClearEvents()

 return order.ID, nil
}

// PayOrder 处理支付订单命令
func (s *OrderCommandService) PayOrder(ctx context.Context, cmd PayOrderCommand) error {
 order, err := s.orderRepo.FindByID(ctx, cmd.OrderID)
 if err != nil {
  return err
 }
 if order == nil {
  return errors.New("order not found")
 }

 if err := order.Pay(); err != nil {
  return err
 }

 if err := s.orderRepo.Update(ctx, order); err != nil {
  return err
 }

 // 发布事件
 for _, event := range order.GetEvents() {
  s.eventBus.Publish(event)
 }
 order.ClearEvents()

 return nil
}

// CancelOrder 处理取消订单命令
func (s *OrderCommandService) CancelOrder(ctx context.Context, cmd CancelOrderCommand) error {
 order, err := s.orderRepo.FindByID(ctx, cmd.OrderID)
 if err != nil {
  return err
 }
 if order == nil {
  return errors.New("order not found")
 }

 if err := order.Cancel(cmd.Reason); err != nil {
  return err
 }

 if err := s.orderRepo.Update(ctx, order); err != nil {
  return err
 }

 // 发布事件
 for _, event := range order.GetEvents() {
  s.eventBus.Publish(event)
 }
 order.ClearEvents()

 return nil
}
```

```go
// internal/query/projections/order_projection.go
package projections

import (
 "time"
)

// OrderView 订单视图模型（为查询优化）
type OrderView struct {
 ID           string
 CustomerID   string
 CustomerName string    // 冗余数据，避免JOIN
 ItemCount    int       // 聚合数据
 Total        float64
 Status       string
 CreatedAt    time.Time
 UpdatedAt    time.Time
}

// OrderDetailView 订单详情视图
type OrderDetailView struct {
 ID          string
 CustomerID  string
 Items       []OrderItemView
 Total       float64
 Status      string
 CreatedAt   time.Time
 UpdatedAt   time.Time
}

// OrderItemView 订单项视图
type OrderItemView struct {
 ProductID string
 Name      string
 Quantity  int
 Price     float64
 Subtotal  float64
}

// OrderSummary 订单汇总
type OrderSummary struct {
 TotalOrders   int
 TotalAmount   float64
 PendingCount  int
 PaidCount     int
}
```

```go
// internal/query/infrastructure/order_read_repository.go
package infrastructure

import (
 "query/internal/projections"
 "context"
 "sync"
)

// OrderReadRepository 订单读仓储接口
type OrderReadRepository interface {
 FindByID(ctx context.Context, id string) (*projections.OrderDetailView, error)
 FindByCustomer(ctx context.Context, customerID string) ([]*projections.OrderView, error)
 FindAll(ctx context.Context, offset, limit int) ([]*projections.OrderView, error)
 GetSummary(ctx context.Context) (*projections.OrderSummary, error)
 Save(ctx context.Context, view *projections.OrderView) error
 Update(ctx context.Context, view *projections.OrderView) error
}

// InMemoryOrderReadRepository 内存读仓储
type InMemoryOrderReadRepository struct {
 mu          sync.RWMutex
 orders      map[string]*projections.OrderView
 orderDetails map[string]*projections.OrderDetailView
}

// NewInMemoryOrderReadRepository 创建读仓储
func NewInMemoryOrderReadRepository() OrderReadRepository {
 return &InMemoryOrderReadRepository{
  orders:       make(map[string]*projections.OrderView),
  orderDetails: make(map[string]*projections.OrderDetailView),
 }
}

func (r *InMemoryOrderReadRepository) FindByID(ctx context.Context, id string) (*projections.OrderDetailView, error) {
 r.mu.RLock()
 defer r.mu.RUnlock()

 return r.orderDetails[id], nil
}

func (r *InMemoryOrderReadRepository) FindByCustomer(ctx context.Context, customerID string) ([]*projections.OrderView, error) {
 r.mu.RLock()
 defer r.mu.RUnlock()

 var result []*projections.OrderView
 for _, order := range r.orders {
  if order.CustomerID == customerID {
   result = append(result, order)
  }
 }
 return result, nil
}

func (r *InMemoryOrderReadRepository) FindAll(ctx context.Context, offset, limit int) ([]*projections.OrderView, error) {
 r.mu.RLock()
 defer r.mu.RUnlock()

 var result []*projections.OrderView
 count := 0
 for _, order := range r.orders {
  if count >= offset && len(result) < limit {
   result = append(result, order)
  }
  count++
 }
 return result, nil
}

func (r *InMemoryOrderReadRepository) GetSummary(ctx context.Context) (*projections.OrderSummary, error) {
 r.mu.RLock()
 defer r.mu.RUnlock()

 summary := &projections.OrderSummary{}
 for _, order := range r.orders {
  summary.TotalOrders++
  summary.TotalAmount += order.Total
  switch order.Status {
  case "PENDING":
   summary.PendingCount++
  case "PAID":
   summary.PaidCount++
  }
 }
 return summary, nil
}

func (r *InMemoryOrderReadRepository) Save(ctx context.Context, view *projections.OrderView) error {
 r.mu.Lock()
 defer r.mu.Unlock()

 r.orders[view.ID] = view
 return nil
}

func (r *InMemoryOrderReadRepository) Update(ctx context.Context, view *projections.OrderView) error {
 r.mu.Lock()
 defer r.mu.Unlock()

 r.orders[view.ID] = view
 return nil
}
```

```go
// internal/query/services/order_query_service.go
package services

import (
 "query/internal/infrastructure"
 "query/internal/projections"
 "context"
)

// OrderQueryService 订单查询服务
type OrderQueryService struct {
 readRepo infrastructure.OrderReadRepository
}

// NewOrderQueryService 创建查询服务
func NewOrderQueryService(readRepo infrastructure.OrderReadRepository) *OrderQueryService {
 return &OrderQueryService{readRepo: readRepo}
}

// GetOrder 获取订单详情
func (s *OrderQueryService) GetOrder(ctx context.Context, orderID string) (*projections.OrderDetailView, error) {
 return s.readRepo.FindByID(ctx, orderID)
}

// GetCustomerOrders 获取客户订单列表
func (s *OrderQueryService) GetCustomerOrders(ctx context.Context, customerID string) ([]*projections.OrderView, error) {
 return s.readRepo.FindByCustomer(ctx, customerID)
}

// ListOrders 列出订单
func (s *OrderQueryService) ListOrders(ctx context.Context, offset, limit int) ([]*projections.OrderView, error) {
 return s.readRepo.FindAll(ctx, offset, limit)
}

// GetOrderSummary 获取订单汇总
func (s *OrderQueryService) GetOrderSummary(ctx context.Context) (*projections.OrderSummary, error) {
 return s.readRepo.GetSummary(ctx)
}
```

```go
// internal/query/handlers/order_event_handler.go
package handlers

import (
 "query/internal/infrastructure"
 "query/internal/projections"
 "shared/events"
 "shared/messaging"
 "time"
)

// OrderEventHandler 订单事件处理器
type OrderEventHandler struct {
 readRepo infrastructure.OrderReadRepository
}

// NewOrderEventHandler 创建事件处理器
func NewOrderEventHandler(readRepo infrastructure.OrderReadRepository, eventBus messaging.EventBus) *OrderEventHandler {
 handler := &OrderEventHandler{readRepo: readRepo}

 // 订阅事件
 eventBus.Subscribe("OrderCreated", handler.handleOrderCreated)
 eventBus.Subscribe("OrderStatusChanged", handler.handleOrderStatusChanged)
 eventBus.Subscribe("OrderCancelled", handler.handleOrderCancelled)

 return handler
}

func (h *OrderEventHandler) handleOrderCreated(event events.Event) {
 e, ok := event.(events.OrderCreatedEvent)
 if !ok {
  return
 }

 view := &projections.OrderView{
  ID:         e.OrderID,
  CustomerID: e.CustomerID,
  ItemCount:  len(e.Items),
  Total:      e.Total,
  Status:     "PENDING",
  CreatedAt:  e.Timestamp,
  UpdatedAt:  e.Timestamp,
 }

 h.readRepo.Save(nil, view)
}

func (h *OrderEventHandler) handleOrderStatusChanged(event events.Event) {
 e, ok := event.(events.OrderStatusChangedEvent)
 if !ok {
  return
 }

 view, _ := h.readRepo.FindByCustomer(nil, e.OrderID)
 if len(view) > 0 {
  view[0].Status = e.NewStatus
  view[0].UpdatedAt = time.Now()
  h.readRepo.Update(nil, view[0])
 }
}

func (h *OrderEventHandler) handleOrderCancelled(event events.Event) {
 e, ok := event.(events.OrderCancelledEvent)
 if !ok {
  return
 }

 view, _ := h.readRepo.FindByCustomer(nil, e.OrderID)
 if len(view) > 0 {
  view[0].Status = "CANCELLED"
  view[0].UpdatedAt = time.Now()
  h.readRepo.Update(nil, view[0])
 }
}
```

### 5.5 反例说明

#### 反例1：读写模型混合

```go
// ❌ 错误：同一个模型处理读写
package model

type Order struct {
 ID       string
 Customer Customer  // 嵌套对象，查询时需要JOIN
 Items    []Item    // 大量数据，查询时全量加载
 // ...
}

// 写操作
func (o *Order) Save() error { /* ... */ }

// 读操作 - 返回完整对象
func GetOrder(id string) (*Order, error) { /* ... */ }
```

**问题**：

- 查询性能差（需要JOIN和加载大量数据）
- 写操作复杂度高
- 无法针对读写分别优化

#### 反例2：缺少事件同步

```go
// ❌ 错误：直接修改读模型
func UpdateOrderStatus(orderID string, status string) error {
 // 直接更新读模型，没有事件驱动
 return db.Exec("UPDATE order_views SET status = ? WHERE id = ?", status, orderID)
}
```

**问题**：

- 数据不一致风险
- 无法追踪变更历史
- 违反CQRS原则

### 5.6 选型指南

| 场景 | 建议 |
|------|------|
| **读多写少** | CQRS可显著优化读性能 |
| **复杂查询需求** | 读模型可针对查询优化 |
| **高并发写入** | 写模型可独立扩展 |
| **需要事件溯源** | CQRS与事件溯源天然结合 |
| **简单CRUD** | 可能过度设计 |
| **强一致性要求** | 需要处理最终一致性 |

### 5.7 优缺点分析

| 优点 | 缺点 |
|------|------|
| ✅ 读写可独立优化 | ❌ 系统复杂度增加 |
| ✅ 读性能显著提升 | ❌ 需要处理最终一致性 |
| ✅ 可独立扩展 | ❌ 数据同步复杂性 |
| ✅ 支持多种读模型 | ❌ 开发成本增加 |
| ✅ 与事件溯源结合 | ❌ 需要事件基础设施 |

---

## 6. 事件溯源（Event Sourcing）

### 6.1 概念定义

事件溯源（Event Sourcing）是一种数据持久化模式，它不是存储对象的当前状态，而是存储导致状态变更的所有事件序列。通过重放这些事件，可以重建对象在任何时间点的状态。

核心思想：**状态是事件的派生，事件是真相的唯一来源**。

### 6.2 架构结构

```
┌─────────────────────────────────────────────────────────────┐
│                        应用程序                              │
└───────────────────────┬─────────────────────────────────────┘
                        │
        ┌───────────────┴───────────────┐
        │                               │
┌───────▼────────┐            ┌────────▼────────┐
│   命令处理      │            │    查询处理      │
│                │            │                 │
│  1. 接收命令    │            │  1. 查询请求     │
│  2. 创建事件    │            │  2. 返回投影     │
│  3. 存储事件    │            │                 │
└───────┬────────┘            └─────────────────┘
        │
        │ 存储事件
        ▼
┌─────────────────────────────────────────────────────────────┐
│                    事件存储 (Event Store)                    │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  EventID │ AggregateID │ EventType │ Data │ Timestamp │  │
│  ├───────────────────────────────────────────────────────┤  │
│  │  uuid-1  │ account-1   │ Deposited │ {...}│ t1        │  │
│  │  uuid-2  │ account-1   │ Withdrawn │ {...}│ t2        │  │
│  │  uuid-3  │ account-1   │ Deposited │ {...}│ t3        │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
        │
        │ 发布事件
        ▼
┌─────────────────────────────────────────────────────────────┐
│                    事件处理器                                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ 投影更新器   │  │  外部通知   │  │  审计日志   │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
        │
        ▼
┌─────────────────────────────────────────────────────────────┐
│                    投影存储 (Read Model)                     │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  ID      │ Balance │ LastTransaction │ UpdatedAt       │  │
│  ├───────────────────────────────────────────────────────┤  │
│  │account-1 │  1500   │  Deposited      │ t3              │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 6.3 核心原则

| 原则 | 说明 |
|------|------|
| **事件为真相来源** | 只存储事件，状态通过重放事件获得 |
| **不可变事件** | 一旦存储，事件不可修改 |
| **顺序保证** | 同一聚合的事件必须按顺序处理 |
| **快照优化** | 大事件流使用快照加速状态重建 |

### 6.4 Go实现示例

#### 项目结构

```
event_sourcing/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── events/              # 事件定义
│   │   ├── event.go
│   │   └── bank_events.go
│   ├── domain/              # 领域模型
│   │   ├── aggregate.go
│   │   └── account.go
│   ├── eventstore/          # 事件存储
│   │   ├── store.go
│   │   └── inmemory_store.go
│   ├── projections/         # 投影
│   │   ├── projection.go
│   │   └── account_projection.go
│   └── application/         # 应用服务
│       ├── commands.go
│       └── account_service.go
└── go.mod
```

#### 完整代码实现

```go
// internal/events/event.go
package events

import (
 "time"
)

// Event 基础事件接口
type Event interface {
 EventID() string
 EventType() string
 AggregateID() string
 AggregateType() string
 Version() int
 OccurredAt() time.Time
}

// BaseEvent 事件基础结构
type BaseEvent struct {
 ID            string
 Type          string
 AggID         string
 AggType       string
 EventVersion  int
 Timestamp     time.Time
}

func (e BaseEvent) EventID() string        { return e.ID }
func (e BaseEvent) EventType() string      { return e.Type }
func (e BaseEvent) AggregateID() string    { return e.AggID }
func (e BaseEvent) AggregateType() string  { return e.AggType }
func (e BaseEvent) Version() int           { return e.EventVersion }
func (e BaseEvent) OccurredAt() time.Time  { return e.Timestamp }

// EventEnvelope 事件信封
type EventEnvelope struct {
 Event       Event
 Metadata    map[string]string
}
```

```go
// internal/events/bank_events.go
package events

import (
 "time"
)

// AccountCreatedEvent 账户创建事件
type AccountCreatedEvent struct {
 BaseEvent
 Owner   string
 Balance float64
}

// MoneyDepositedEvent 存款事件
type MoneyDepositedEvent struct {
 BaseEvent
 Amount  float64
 Balance float64 // 存款后的余额
}

// MoneyWithdrawnEvent 取款事件
type MoneyWithdrawnEvent struct {
 BaseEvent
 Amount  float64
 Balance float64 // 取款后的余额
}

// AccountClosedEvent 账户关闭事件
type AccountClosedEvent struct {
 BaseEvent
 Reason string
}

// NewAccountCreatedEvent 创建账户创建事件
func NewAccountCreatedEvent(aggregateID, owner string, initialBalance float64) *AccountCreatedEvent {
 return &AccountCreatedEvent{
  BaseEvent: BaseEvent{
   ID:           generateID(),
   Type:         "AccountCreated",
   AggID:        aggregateID,
   AggType:      "Account",
   EventVersion: 1,
   Timestamp:    time.Now(),
  },
  Owner:   owner,
  Balance: initialBalance,
 }
}

// NewMoneyDepositedEvent 创建存款事件
func NewMoneyDepositedEvent(aggregateID string, version int, amount, balance float64) *MoneyDepositedEvent {
 return &MoneyDepositedEvent{
  BaseEvent: BaseEvent{
   ID:           generateID(),
   Type:         "MoneyDeposited",
   AggID:        aggregateID,
   AggType:      "Account",
   EventVersion: version,
   Timestamp:    time.Now(),
  },
  Amount:  amount,
  Balance: balance,
 }
}

// NewMoneyWithdrawnEvent 创建取款事件
func NewMoneyWithdrawnEvent(aggregateID string, version int, amount, balance float64) *MoneyWithdrawnEvent {
 return &MoneyWithdrawnEvent{
  BaseEvent: BaseEvent{
   ID:           generateID(),
   Type:         "MoneyWithdrawn",
   AggID:        aggregateID,
   AggType:      "Account",
   EventVersion: version,
   Timestamp:    time.Now(),
  },
  Amount:  amount,
  Balance: balance,
 }
}

func generateID() string {
 return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
 const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
 b := make([]byte, n)
 for i := range b {
  b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
 }
 return string(b)
}
```

```go
// internal/domain/aggregate.go
package domain

import (
 "event_sourcing/internal/events"
)

// Aggregate 聚合根接口
type Aggregate interface {
 ID() string
 Version() int
 ApplyEvent(event events.Event) error
 GetUncommittedEvents() []events.Event
 ClearUncommittedEvents()
}

// BaseAggregate 聚合根基类
type BaseAggregate struct {
 AggregateID       string
 AggregateVersion  int
 uncommittedEvents []events.Event
}

func (a *BaseAggregate) ID() string              { return a.AggregateID }
func (a *BaseAggregate) Version() int            { return a.AggregateVersion }
func (a *BaseAggregate) GetUncommittedEvents() []events.Event {
 return a.uncommittedEvents
}
func (a *BaseAggregate) ClearUncommittedEvents() {
 a.uncommittedEvents = nil
}

func (a *BaseAggregate) RaiseEvent(event events.Event) {
 a.uncommittedEvents = append(a.uncommittedEvents, event)
 a.AggregateVersion = event.Version()
}
```

```go
// internal/domain/account.go
package domain

import (
 "errors"
 "event_sourcing/internal/events"
 "time"
)

// Account 账户聚合根
type Account struct {
 BaseAggregate
 Owner      string
 Balance    float64
 IsClosed   bool
 CreatedAt  time.Time
}

// NewAccount 创建新账户
func NewAccount(id, owner string, initialBalance float64) (*Account, error) {
 if owner == "" {
  return nil, errors.New("owner is required")
 }
 if initialBalance < 0 {
  return nil, errors.New("initial balance cannot be negative")
 }

 account := &Account{
  BaseAggregate: BaseAggregate{
   AggregateID: id,
  },
  Owner:     owner,
  Balance:   initialBalance,
  CreatedAt: time.Now(),
 }

 // 创建并应用事件
 event := events.NewAccountCreatedEvent(id, owner, initialBalance)
 account.RaiseEvent(event)

 return account, nil
}

// Deposit 存款
func (a *Account) Deposit(amount float64) error {
 if a.IsClosed {
  return errors.New("cannot deposit to closed account")
 }
 if amount <= 0 {
  return errors.New("deposit amount must be positive")
 }

 a.Balance += amount
 event := events.NewMoneyDepositedEvent(a.AggregateID, a.AggregateVersion+1, amount, a.Balance)
 a.RaiseEvent(event)

 return nil
}

// Withdraw 取款
func (a *Account) Withdraw(amount float64) error {
 if a.IsClosed {
  return errors.New("cannot withdraw from closed account")
 }
 if amount <= 0 {
  return errors.New("withdraw amount must be positive")
 }
 if amount > a.Balance {
  return errors.New("insufficient balance")
 }

 a.Balance -= amount
 event := events.NewMoneyWithdrawnEvent(a.AggregateID, a.AggregateVersion+1, amount, a.Balance)
 a.RaiseEvent(event)

 return nil
}

// Close 关闭账户
func (a *Account) Close(reason string) error {
 if a.IsClosed {
  return errors.New("account is already closed")
 }

 a.IsClosed = true
 event := &events.AccountClosedEvent{
  BaseEvent: events.BaseEvent{
   ID:           events.generateID(),
   Type:         "AccountClosed",
   AggID:        a.AggregateID,
   AggType:      "Account",
   EventVersion: a.AggregateVersion + 1,
   Timestamp:    time.Now(),
  },
  Reason: reason,
 }
 a.RaiseEvent(event)

 return nil
}

// ApplyEvent 应用事件到聚合
func (a *Account) ApplyEvent(event events.Event) error {
 switch e := event.(type) {
 case *events.AccountCreatedEvent:
  a.applyAccountCreated(e)
 case *events.MoneyDepositedEvent:
  a.applyMoneyDeposited(e)
 case *events.MoneyWithdrawnEvent:
  a.applyMoneyWithdrawn(e)
 case *events.AccountClosedEvent:
  a.applyAccountClosed(e)
 default:
  return errors.New("unknown event type")
 }
 a.AggregateVersion = event.Version()
 return nil
}

func (a *Account) applyAccountCreated(e *events.AccountCreatedEvent) {
 a.AggregateID = e.AggregateID()
 a.Owner = e.Owner
 a.Balance = e.Balance
 a.CreatedAt = e.OccurredAt()
}

func (a *Account) applyMoneyDeposited(e *events.MoneyDepositedEvent) {
 a.Balance = e.Balance
}

func (a *Account) applyMoneyWithdrawn(e *events.MoneyWithdrawnEvent) {
 a.Balance = e.Balance
}

func (a *Account) applyAccountClosed(e *events.AccountClosedEvent) {
 a.IsClosed = true
}

// ReconstructAccount 从事件重建账户
func ReconstructAccount(events []events.Event) (*Account, error) {
 if len(events) == 0 {
  return nil, errors.New("no events to reconstruct")
 }

 account := &Account{}
 for _, event := range events {
  if err := account.ApplyEvent(event); err != nil {
   return nil, err
  }
 }
 return account, nil
}
```

```go
// internal/eventstore/store.go
package eventstore

import (
 "context"
 "event_sourcing/internal/events"
)

// EventStore 事件存储接口
type EventStore interface {
 // Append 追加事件
 Append(ctx context.Context, aggregateID string, events []events.Event, expectedVersion int) error

 // GetEvents 获取聚合的所有事件
 GetEvents(ctx context.Context, aggregateID string) ([]events.Event, error)

 // GetEventsFromVersion 从指定版本获取事件
 GetEventsFromVersion(ctx context.Context, aggregateID string, fromVersion int) ([]events.Event, error)

 // GetAllEvents 获取所有事件（用于投影重建）
 GetAllEvents(ctx context.Context, offset, limit int) ([]events.Event, error)

 // Subscribe 订阅新事件
 Subscribe(ctx context.Context, eventTypes []string, handler EventHandler) error
}

// EventHandler 事件处理器
type EventHandler func(event events.Event) error
```

```go
// internal/eventstore/inmemory_store.go
package eventstore

import (
 "context"
 "event_sourcing/internal/events"
 "errors"
 "sync"
)

// InMemoryEventStore 内存事件存储
type InMemoryEventStore struct {
 mu       sync.RWMutex
 events   map[string][]events.Event
 versions map[string]int
}

// NewInMemoryEventStore 创建内存事件存储
func NewInMemoryEventStore() EventStore {
 return &InMemoryEventStore{
  events:   make(map[string][]events.Event),
  versions: make(map[string]int),
 }
}

func (s *InMemoryEventStore) Append(
 ctx context.Context,
 aggregateID string,
 events []events.Event,
 expectedVersion int,
) error {
 s.mu.Lock()
 defer s.mu.Unlock()

 currentVersion := s.versions[aggregateID]
 if currentVersion != expectedVersion {
  return errors.New("concurrency conflict: version mismatch")
 }

 for _, event := range events {
  s.events[aggregateID] = append(s.events[aggregateID], event)
  s.versions[aggregateID] = event.Version()
 }

 return nil
}

func (s *InMemoryEventStore) GetEvents(ctx context.Context, aggregateID string) ([]events.Event, error) {
 s.mu.RLock()
 defer s.mu.RUnlock()

 events := s.events[aggregateID]
 result := make([]events.Event, len(events))
 copy(result, events)
 return result, nil
}

func (s *InMemoryEventStore) GetEventsFromVersion(
 ctx context.Context,
 aggregateID string,
 fromVersion int,
) ([]events.Event, error) {
 s.mu.RLock()
 defer s.mu.RUnlock()

 var result []events.Event
 for _, event := range s.events[aggregateID] {
  if event.Version() >= fromVersion {
   result = append(result, event)
  }
 }
 return result, nil
}

func (s *InMemoryEventStore) GetAllEvents(ctx context.Context, offset, limit int) ([]events.Event, error) {
 s.mu.RLock()
 defer s.mu.RUnlock()

 var allEvents []events.Event
 for _, events := range s.events {
  allEvents = append(allEvents, events...)
 }

 if offset >= len(allEvents) {
  return []events.Event{}, nil
 }

 end := offset + limit
 if end > len(allEvents) {
  end = len(allEvents)
 }

 return allEvents[offset:end], nil
}

func (s *InMemoryEventStore) Subscribe(ctx context.Context, eventTypes []string, handler EventHandler) error {
 // 简化实现：内存存储不支持订阅
 return nil
}
```

```go
// internal/projections/account_projection.go
package projections

import (
 "event_sourcing/internal/events"
 "sync"
 "time"
)

// AccountView 账户视图
type AccountView struct {
 ID            string
 Owner         string
 Balance       float64
 TotalDeposits float64
 TotalWithdrawals float64
 TransactionCount int
 IsClosed      bool
 CreatedAt     time.Time
 LastUpdatedAt time.Time
}

// AccountProjection 账户投影
type AccountProjection struct {
 mu       sync.RWMutex
 accounts map[string]*AccountView
}

// NewAccountProjection 创建账户投影
func NewAccountProjection() *AccountProjection {
 return &AccountProjection{
  accounts: make(map[string]*AccountView),
 }
}

// HandleEvent 处理事件
func (p *AccountProjection) HandleEvent(event events.Event) error {
 switch e := event.(type) {
 case *events.AccountCreatedEvent:
  p.handleAccountCreated(e)
 case *events.MoneyDepositedEvent:
  p.handleMoneyDeposited(e)
 case *events.MoneyWithdrawnEvent:
  p.handleMoneyWithdrawn(e)
 case *events.AccountClosedEvent:
  p.handleAccountClosed(e)
 }
 return nil
}

func (p *AccountProjection) handleAccountCreated(e *events.AccountCreatedEvent) {
 p.mu.Lock()
 defer p.mu.Unlock()

 p.accounts[e.AggregateID()] = &AccountView{
  ID:            e.AggregateID(),
  Owner:         e.Owner,
  Balance:       e.Balance,
  CreatedAt:     e.OccurredAt(),
  LastUpdatedAt: e.OccurredAt(),
 }
}

func (p *AccountProjection) handleMoneyDeposited(e *events.MoneyDepositedEvent) {
 p.mu.Lock()
 defer p.mu.Unlock()

 if account, ok := p.accounts[e.AggregateID()]; ok {
  account.Balance = e.Balance
  account.TotalDeposits += e.Amount
  account.TransactionCount++
  account.LastUpdatedAt = e.OccurredAt()
 }
}

func (p *AccountProjection) handleMoneyWithdrawn(e *events.MoneyWithdrawnEvent) {
 p.mu.Lock()
 defer p.mu.Unlock()

 if account, ok := p.accounts[e.AggregateID()]; ok {
  account.Balance = e.Balance
  account.TotalWithdrawals += e.Amount
  account.TransactionCount++
  account.LastUpdatedAt = e.OccurredAt()
 }
}

func (p *AccountProjection) handleAccountClosed(e *events.AccountClosedEvent) {
 p.mu.Lock()
 defer p.mu.Unlock()

 if account, ok := p.accounts[e.AggregateID()]; ok {
  account.IsClosed = true
  account.LastUpdatedAt = e.OccurredAt()
 }
}

// GetAccount 获取账户视图
func (p *AccountProjection) GetAccount(id string) (*AccountView, bool) {
 p.mu.RLock()
 defer p.mu.RUnlock()

 account, ok := p.accounts[id]
 if !ok {
  return nil, false
 }
 // 返回副本
 return &AccountView{
  ID:               account.ID,
  Owner:            account.Owner,
  Balance:          account.Balance,
  TotalDeposits:    account.TotalDeposits,
  TotalWithdrawals: account.TotalWithdrawals,
  TransactionCount: account.TransactionCount,
  IsClosed:         account.IsClosed,
  CreatedAt:        account.CreatedAt,
  LastUpdatedAt:    account.LastUpdatedAt,
 }, true
}

// GetAllAccounts 获取所有账户
func (p *AccountProjection) GetAllAccounts() []*AccountView {
 p.mu.RLock()
 defer p.mu.RUnlock()

 var result []*AccountView
 for _, account := range p.accounts {
  result = append(result, &AccountView{
   ID:               account.ID,
   Owner:            account.Owner,
   Balance:          account.Balance,
   TotalDeposits:    account.TotalDeposits,
   TotalWithdrawals: account.TotalWithdrawals,
   TransactionCount: account.TransactionCount,
   IsClosed:         account.IsClosed,
   CreatedAt:        account.CreatedAt,
   LastUpdatedAt:    account.LastUpdatedAt,
  })
 }
 return result
}

// GetSummary 获取汇总信息
func (p *AccountProjection) GetSummary() (totalAccounts int, totalBalance float64) {
 p.mu.RLock()
 defer p.mu.RUnlock()

 for _, account := range p.accounts {
  totalAccounts++
  totalBalance += account.Balance
 }
 return
}
```

```go
// internal/application/commands.go
package application

// CreateAccountCommand 创建账户命令
type CreateAccountCommand struct {
 Owner          string
 InitialBalance float64
}

// DepositCommand 存款命令
type DepositCommand struct {
 AccountID string
 Amount    float64
}

// WithdrawCommand 取款命令
type WithdrawCommand struct {
 AccountID string
 Amount    float64
}

// CloseAccountCommand 关闭账户命令
type CloseAccountCommand struct {
 AccountID string
 Reason    string
}
```

```go
// internal/application/account_service.go
package application

import (
 "context"
 "event_sourcing/internal/domain"
 "event_sourcing/internal/eventstore"
 "event_sourcing/internal/projections"
 "errors"

 "github.com/google/uuid"
)

// AccountService 账户应用服务
type AccountService struct {
 eventStore  eventstore.EventStore
 projection  *projections.AccountProjection
}

// NewAccountService 创建账户服务
func NewAccountService(eventStore eventstore.EventStore, projection *projections.AccountProjection) *AccountService {
 return &AccountService{
  eventStore: eventStore,
  projection: projection,
 }
}

// CreateAccount 创建账户
func (s *AccountService) CreateAccount(ctx context.Context, cmd CreateAccountCommand) (string, error) {
 account, err := domain.NewAccount(uuid.New().String(), cmd.Owner, cmd.InitialBalance)
 if err != nil {
  return "", err
 }

 if err := s.eventStore.Append(ctx, account.ID(), account.GetUncommittedEvents(), 0); err != nil {
  return "", err
 }

 // 更新投影
 for _, event := range account.GetUncommittedEvents() {
  s.projection.HandleEvent(event)
 }
 account.ClearUncommittedEvents()

 return account.ID(), nil
}

// Deposit 存款
func (s *AccountService) Deposit(ctx context.Context, cmd DepositCommand) error {
 // 加载聚合
 events, err := s.eventStore.GetEvents(ctx, cmd.AccountID)
 if err != nil {
  return err
 }
 if len(events) == 0 {
  return errors.New("account not found")
 }

 account, err := domain.ReconstructAccount(events)
 if err != nil {
  return err
 }

 // 执行业务操作
 if err := account.Deposit(cmd.Amount); err != nil {
  return err
 }

 // 保存事件
 if err := s.eventStore.Append(ctx, account.ID(), account.GetUncommittedEvents(), account.Version()-len(account.GetUncommittedEvents())); err != nil {
  return err
 }

 // 更新投影
 for _, event := range account.GetUncommittedEvents() {
  s.projection.HandleEvent(event)
 }
 account.ClearUncommittedEvents()

 return nil
}

// Withdraw 取款
func (s *AccountService) Withdraw(ctx context.Context, cmd WithdrawCommand) error {
 events, err := s.eventStore.GetEvents(ctx, cmd.AccountID)
 if err != nil {
  return err
 }
 if len(events) == 0 {
  return errors.New("account not found")
 }

 account, err := domain.ReconstructAccount(events)
 if err != nil {
  return err
 }

 if err := account.Withdraw(cmd.Amount); err != nil {
  return err
 }

 if err := s.eventStore.Append(ctx, account.ID(), account.GetUncommittedEvents(), account.Version()-len(account.GetUncommittedEvents())); err != nil {
  return err
 }

 for _, event := range account.GetUncommittedEvents() {
  s.projection.HandleEvent(event)
 }
 account.ClearUncommittedEvents()

 return nil
}

// GetAccount 获取账户（从投影读取）
func (s *AccountService) GetAccount(id string) (*projections.AccountView, error) {
 account, ok := s.projection.GetAccount(id)
 if !ok {
  return nil, errors.New("account not found")
 }
 return account, nil
}

// GetAllAccounts 获取所有账户
func (s *AccountService) GetAllAccounts() []*projections.AccountView {
 return s.projection.GetAllAccounts()
}

// GetSummary 获取汇总
func (s *AccountService) GetSummary() (int, float64) {
 return s.projection.GetSummary()
}
```

```go
// cmd/api/main.go
package main

import (
 "context"
 "encoding/json"
 "event_sourcing/internal/application"
 "event_sourcing/internal/eventstore"
 "event_sourcing/internal/projections"
 "fmt"
 "net/http"
)

func main() {
 // 初始化基础设施
 eventStore := eventstore.NewInMemoryEventStore()
 accountProjection := projections.NewAccountProjection()

 // 初始化应用服务
 accountService := application.NewAccountService(eventStore, accountProjection)

 // 设置HTTP处理器
 http.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodPost:
   createAccount(w, r, accountService)
  case http.MethodGet:
   listAccounts(w, r, accountService)
  }
 })

 http.HandleFunc("/accounts/", func(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
   getAccount(w, r, accountService)
  case http.MethodPost:
   if r.URL.Path[len(r.URL.Path)-8:] == "/deposit" {
    deposit(w, r, accountService)
   } else if r.URL.Path[len(r.URL.Path)-9:] == "/withdraw" {
    withdraw(w, r, accountService)
   }
  }
 })

 fmt.Println("Event Sourcing Server starting on :8080")
 http.ListenAndServe(":8080", nil)
}

func createAccount(w http.ResponseWriter, r *http.Request, service *application.AccountService) {
 var cmd application.CreateAccountCommand
 if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
  http.Error(w, err.Error(), http.StatusBadRequest)
  return
 }

 id, err := service.CreateAccount(r.Context(), cmd)
 if err != nil {
  http.Error(w, err.Error(), http.StatusBadRequest)
  return
 }

 json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func getAccount(w http.ResponseWriter, r *http.Request, service *application.AccountService) {
 id := r.URL.Path[len("/accounts/"):]
 account, err := service.GetAccount(id)
 if err != nil {
  http.Error(w, err.Error(), http.StatusNotFound)
  return
 }

 json.NewEncoder(w).Encode(account)
}

func listAccounts(w http.ResponseWriter, r *http.Request, service *application.AccountService) {
 accounts := service.GetAllAccounts()
 json.NewEncoder(w).Encode(accounts)
}

func deposit(w http.ResponseWriter, r *http.Request, service *application.AccountService) {
 id := r.URL.Path[len("/accounts/") : len(r.URL.Path)-8]
 var cmd application.DepositCommand
 cmd.AccountID = id
 if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
  http.Error(w, err.Error(), http.StatusBadRequest)
  return
 }

 if err := service.Deposit(r.Context(), cmd); err != nil {
  http.Error(w, err.Error(), http.StatusBadRequest)
  return
 }

 w.WriteHeader(http.StatusOK)
}

func withdraw(w http.ResponseWriter, r *http.Request, service *application.AccountService) {
 id := r.URL.Path[len("/accounts/") : len(r.URL.Path)-9]
 var cmd application.WithdrawCommand
 cmd.AccountID = id
 if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
  http.Error(w, err.Error(), http.StatusBadRequest)
  return
 }

 if err := service.Withdraw(r.Context(), cmd); err != nil {
  http.Error(w, err.Error(), http.StatusBadRequest)
  return
 }

 w.WriteHeader(http.StatusOK)
}
```

### 6.5 反例说明

#### 反例1：直接修改状态

```go
// ❌ 错误：直接修改状态而不生成事件
func (a *Account) Deposit(amount float64) error {
 a.Balance += amount  // 直接修改，没有事件
 return a.db.Save(a)  // 直接保存状态
}
```

**问题**：

- 丢失了变更历史
- 无法重建过去状态
- 无法审计和追踪

#### 反例2：事件修改

```go
// ❌ 错误：尝试修改已存储的事件
func CorrectEvent(eventID string, newAmount float64) error {
 return db.Exec("UPDATE events SET data = ? WHERE id = ?", newAmount, eventID)
}
```

**问题**：

- 破坏事件溯源的核心原则
- 导致数据不一致
- 无法进行审计

### 6.6 选型指南

| 场景 | 建议 |
|------|------|
| **需要完整审计** | 事件溯源天然支持审计 |
| **复杂业务逻辑** | 可追溯业务规则变更 |
| **需要时间旅行** | 可重建任意时间点的状态 |
| **与CQRS结合** | 事件驱动读写分离 |
| **简单CRUD** | 可能过度设计 |
| **存储敏感** | 事件存储量随时间增长 |

### 6.7 优缺点分析

| 优点 | 缺点 |
|------|------|
| ✅ 完整的历史记录 | ❌ 系统复杂度高 |
| ✅ 可重建任意状态 | ❌ 存储需求大 |
| ✅ 天然支持审计 | ❌ 学习曲线陡峭 |
| ✅ 便于调试和分析 | ❌ 事件模式设计困难 |
| ✅ 与CQRS完美结合 | ❌ 需要处理事件版本 |

---

## 7. 领域驱动设计（DDD）

### 7.1 概念定义

领域驱动设计（Domain-Driven Design，DDD）是由Eric Evans在2003年提出的软件开发方法论。它强调以领域为核心，通过与领域专家紧密合作，构建能够准确反映业务领域的软件模型。

DDD的核心思想是：**软件的核心复杂性在于业务领域本身，技术只是实现手段**。

### 7.2 架构结构

```
┌─────────────────────────────────────────────────────────────┐
│                      限界上下文 (Bounded Context)              │
│  ┌───────────────────────────────────────────────────────┐  │
│  │                    用户接口层                          │  │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐              │  │
│  │  │Controller│  │ Handler │  │  DTO    │              │  │
│  │  └─────────┘  └─────────┘  └─────────┘              │  │
│  ├───────────────────────────────────────────────────────┤  │
│  │                    应用层                              │  │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐              │  │
│  │  │ App Svc │  │Use Case │  │  DTO    │              │  │
│  │  └─────────┘  └─────────┘  └─────────┘              │  │
│  ├───────────────────────────────────────────────────────┤  │
│  │                    领域层                              │  │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐ │  │
│  │  │ Entity  │  │Value Obj│  │Aggregate│  │Domain   │ │  │
│  │  │         │  │         │  │  Root   │  │ Service │ │  │
│  │  └─────────┘  └─────────┘  └─────────┘  └─────────┘ │  │
│  ├───────────────────────────────────────────────────────┤  │
│  │                  基础设施层                            │  │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐ │  │
│  │  │   Repo  │  │   DB    │  │ External│  │ Message │ │  │
│  │  │  Impl   │  │         │  │ Service │  │ Queue   │ │  │
│  │  └─────────┘  └─────────┘  └─────────┘  └─────────┘ │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 7.3 核心概念

#### 7.3.1 限界上下文（Bounded Context）

限界上下文是DDD的核心战略模式，它定义了领域模型的边界。在同一个限界上下文内，领域术语具有统一的含义。

```
┌─────────────────────────────────────────────────────────────┐
│                    电商系统                                  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   订单上下文 │  │   库存上下文 │  │   支付上下文 │         │
│  │             │  │             │  │             │         │
│  │ - Order     │  │ - Product   │  │ - Payment   │         │
│  │ - OrderItem │  │ - Stock     │  │ - Refund    │         │
│  │ - Shipping  │  │ - Warehouse │  │ - Invoice   │         │
│  │             │  │             │  │             │         │
│  │ "Product"   │  │ "Product"   │  │             │         │
│  │ = 订单项    │  │ = 商品实体  │  │             │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│         │                │                │                 │
│         └────────────────┴────────────────┘                 │
│                    防腐层 (Anti-Corruption Layer)            │
└─────────────────────────────────────────────────────────────┘
```

#### 7.3.2 实体（Entity）与值对象（Value Object）

| 特性 | 实体（Entity） | 值对象（Value Object） |
|------|---------------|----------------------|
| **身份标识** | 有唯一标识 | 无身份标识 |
| **相等性** | 基于ID相等 | 基于属性相等 |
| **可变性** | 可变 | 不可变 |
| **生命周期** | 有生命周期 | 可随意创建和替换 |
| **示例** | User、Order、Product | Money、Address、Email |

#### 7.3.3 聚合（Aggregate）与聚合根（Aggregate Root）

聚合是一组相关对象的集合，作为数据修改的单元。聚合根是聚合的入口点，外部只能通过聚合根访问聚合内部的对象。

```
┌─────────────────────────────────────────────────────────────┐
│                      Order (聚合根)                          │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  - ID: OrderID                                        │  │
│  │  - CustomerID: CustomerID                             │  │
│  │  - Status: OrderStatus                                │  │
│  │  - Total: Money                                       │  │
│  │  - Items: []OrderItem                                 │  │
│  │  - ShippingAddress: Address (值对象)                  │  │
│  │                                                       │  │
│  │  + AddItem(item OrderItem)                            │  │
│  │  + RemoveItem(productID ProductID)                    │  │
│  │  + CalculateTotal() Money                             │  │
│  │  + Place() error                                      │  │
│  │  + Pay() error                                        │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                              │
│  聚合边界 ────────────────────────────────────────────────   │
│                                                              │
│  ┌─────────────────┐  ┌─────────────────┐                   │
│  │   OrderItem     │  │   OrderItem     │                   │
│  │   (实体)         │  │   (实体)         │                   │
│  │  - ProductID    │  │  - ProductID    │                   │
│  │  - Name         │  │  - Name         │                   │
│  │  - Quantity     │  │  - Quantity     │                   │
│  │  - Price        │  │  - Price        │                   │
│  └─────────────────┘  └─────────────────┘                   │
└─────────────────────────────────────────────────────────────┘
```

### 7.4 Go实现示例

#### 项目结构

```
ddd_architecture/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── order/                    # 订单限界上下文
│   │   ├── domain/               # 领域层
│   │   │   ├── aggregate/
│   │   │   │   └── order.go      # 聚合根
│   │   │   ├── entity/
│   │   │   │   └── order_item.go # 实体
│   │   │   ├── valueobject/
│   │   │   │   ├── money.go      # 值对象
│   │   │   │   └── address.go    # 值对象
│   │   │   ├── service/
│   │   │   │   └── pricing_service.go # 领域服务
│   │   │   ├── repository/
│   │   │   │   └── order_repository.go # 仓储接口
│   │   │   └── event/
│   │   │       └── order_events.go
│   │   ├── application/          # 应用层
│   │   │   ├── service/
│   │   │   │   └── order_app_service.go
│   │   │   └── dto/
│   │   │       └── order_dto.go
│   │   ├── infrastructure/       # 基础设施层
│   │   │   └── persistence/
│   │   │       └── order_repository_impl.go
│   │   └── interfaces/           # 接口层
│   │       └── http/
│   │           └── order_handler.go
│   └── shared/                   # 共享内核
│       └── event/
│           └── event.go
└── go.mod
```

#### 完整代码实现

```go
// internal/order/domain/valueobject/money.go
package valueobject

import (
 "errors"
 "fmt"
)

// Money 金额值对象
type Money struct {
 amount   float64
 currency Currency
}

// Currency 货币类型
type Currency string

const (
 CurrencyCNY Currency = "CNY"
 CurrencyUSD Currency = "USD"
 CurrencyEUR Currency = "EUR"
)

// NewMoney 创建金额
func NewMoney(amount float64, currency Currency) (Money, error) {
 if amount < 0 {
  return Money{}, errors.New("amount cannot be negative")
 }
 if !isValidCurrency(currency) {
  return Money{}, errors.New("invalid currency")
 }
 return Money{amount: amount, currency: currency}, nil
}

// MustNewMoney 创建金额（panic on error）
func MustNewMoney(amount float64, currency Currency) Money {
 m, err := NewMoney(amount, currency)
 if err != nil {
  panic(err)
 }
 return m
}

// ZeroMoney 零金额
func ZeroMoney(currency Currency) Money {
 return MustNewMoney(0, currency)
}

// Amount 获取金额
func (m Money) Amount() float64 {
 return m.amount
}

// Currency 获取货币
func (m Money) Currency() Currency {
 return m.currency
}

// Add 相加
func (m Money) Add(other Money) (Money, error) {
 if m.currency != other.currency {
  return Money{}, errors.New("cannot add different currencies")
 }
 return NewMoney(m.amount+other.amount, m.currency)
}

// Subtract 相减
func (m Money) Subtract(other Money) (Money, error) {
 if m.currency != other.currency {
  return Money{}, errors.New("cannot subtract different currencies")
 }
 if m.amount < other.amount {
  return Money{}, errors.New("insufficient amount")
 }
 return NewMoney(m.amount-other.amount, m.currency)
}

// Multiply 乘法
func (m Money) Multiply(factor float64) (Money, error) {
 return NewMoney(m.amount*factor, m.currency)
}

// IsZero 是否为零
func (m Money) IsZero() bool {
 return m.amount == 0
}

// IsPositive 是否为正
func (m Money) IsPositive() bool {
 return m.amount > 0
}

// Equals 相等比较
func (m Money) Equals(other Money) bool {
 return m.amount == other.amount && m.currency == other.currency
}

// GreaterThan 大于比较
func (m Money) GreaterThan(other Money) (bool, error) {
 if m.currency != other.currency {
  return false, errors.New("cannot compare different currencies")
 }
 return m.amount > other.amount, nil
}

func (m Money) String() string {
 return fmt.Sprintf("%.2f %s", m.amount, m.currency)
}

func isValidCurrency(c Currency) bool {
 return c == CurrencyCNY || c == CurrencyUSD || c == CurrencyEUR
}
```

```go
// internal/order/domain/valueobject/address.go
package valueobject

import (
 "errors"
 "fmt"
)

// Address 地址值对象
type Address struct {
 province   string
 city       string
 district   string
 street     string
 zipCode    string
 recipient  string
 phone      string
}

// NewAddress 创建地址
func NewAddress(province, city, district, street, zipCode, recipient, phone string) (Address, error) {
 if province == "" || city == "" || district == "" || street == "" {
  return Address{}, errors.New("address fields cannot be empty")
 }
 if recipient == "" {
  return Address{}, errors.New("recipient is required")
 }
 if phone == "" {
  return Address{}, errors.New("phone is required")
 }

 return Address{
  province:  province,
  city:      city,
  district:  district,
  street:    street,
  zipCode:   zipCode,
  recipient: recipient,
  phone:     phone,
 }, nil
}

// Province 省份
func (a Address) Province() string { return a.province }

// City 城市
func (a Address) City() string { return a.city }

// District 区县
func (a Address) District() string { return a.district }

// Street 街道
func (a Address) Street() string { return a.street }

// ZipCode 邮编
func (a Address) ZipCode() string { return a.zipCode }

// Recipient 收件人
func (a Address) Recipient() string { return a.recipient }

// Phone 电话
func (a Address) Phone() string { return a.phone }

// FullAddress 完整地址
func (a Address) FullAddress() string {
 return fmt.Sprintf("%s%s%s%s %s %s %s",
  a.province, a.city, a.district, a.street,
  a.zipCode, a.recipient, a.phone)
}

// Equals 相等比较
func (a Address) Equals(other Address) bool {
 return a.province == other.province &&
  a.city == other.city &&
  a.district == other.district &&
  a.street == other.street &&
  a.zipCode == other.zipCode &&
  a.recipient == other.recipient &&
  a.phone == other.phone
}

func (a Address) String() string {
 return a.FullAddress()
}
```

```go
// internal/order/domain/entity/order_item.go
package entity

import (
 "errors"
 "order/internal/domain/valueobject"
)

// OrderItemID 订单项ID
type OrderItemID string

// OrderItem 订单项实体
type OrderItem struct {
 id        OrderItemID
 productID string
 name      string
 quantity  int
 price     valueobject.Money
}

// NewOrderItem 创建订单项
func NewOrderItem(id OrderItemID, productID, name string, quantity int, price valueobject.Money) (*OrderItem, error) {
 if productID == "" {
  return nil, errors.New("product ID is required")
 }
 if name == "" {
  return nil, errors.New("product name is required")
 }
 if quantity <= 0 {
  return nil, errors.New("quantity must be positive")
 }
 if price.IsZero() {
  return nil, errors.New("price is required")
 }

 return &OrderItem{
  id:        id,
  productID: productID,
  name:      name,
  quantity:  quantity,
  price:     price,
 }, nil
}

// ID 获取ID
func (o OrderItem) ID() OrderItemID {
 return o.id
}

// ProductID 获取产品ID
func (o OrderItem) ProductID() string {
 return o.productID
}

// Name 获取名称
func (o OrderItem) Name() string {
 return o.name
}

// Quantity 获取数量
func (o OrderItem) Quantity() int {
 return o.quantity
}

// Price 获取单价
func (o OrderItem) Price() valueobject.Money {
 return o.price
}

// Subtotal 获取小计
func (o OrderItem) Subtotal() (valueobject.Money, error) {
 return o.price.Multiply(float64(o.quantity))
}

// UpdateQuantity 更新数量
func (o *OrderItem) UpdateQuantity(quantity int) error {
 if quantity <= 0 {
  return errors.New("quantity must be positive")
 }
 o.quantity = quantity
 return nil
}

// Equals 相等比较（基于ID）
func (o OrderItem) Equals(other OrderItem) bool {
 return o.id == other.id
}
```

```go
// internal/order/domain/aggregate/order.go
package aggregate

import (
 "errors"
 "order/internal/domain/entity"
 "order/internal/domain/valueobject"
 "time"
)

// OrderID 订单ID
type OrderID string

// OrderStatus 订单状态
type OrderStatus int

const (
 OrderStatusPending OrderStatus = iota
 OrderStatusPaid
 OrderStatusShipped
 OrderStatusCompleted
 OrderStatusCancelled
)

func (s OrderStatus) String() string {
 switch s {
 case OrderStatusPending:
  return "PENDING"
 case OrderStatusPaid:
  return "PAID"
 case OrderStatusShipped:
  return "SHIPPED"
 case OrderStatusCompleted:
  return "COMPLETED"
 case OrderStatusCancelled:
  return "CANCELLED"
 default:
  return "UNKNOWN"
 }
}

// Order 订单聚合根
type Order struct {
 id              OrderID
 customerID      string
 items           []*entity.OrderItem
 status          OrderStatus
 total           valueobject.Money
 shippingAddress valueobject.Address
 createdAt       time.Time
 updatedAt       time.Time
 version         int
}

// NewOrder 创建订单
func NewOrder(id OrderID, customerID string, shippingAddress valueobject.Address) (*Order, error) {
 if customerID == "" {
  return nil, errors.New("customer ID is required")
 }

 return &Order{
  id:              id,
  customerID:      customerID,
  items:           make([]*entity.OrderItem, 0),
  status:          OrderStatusPending,
  total:           valueobject.ZeroMoney(valueobject.CurrencyCNY),
  shippingAddress: shippingAddress,
  createdAt:       time.Now(),
  updatedAt:       time.Now(),
  version:         1,
 }, nil
}

// ID 获取订单ID
func (o *Order) ID() OrderID {
 return o.id
}

// CustomerID 获取客户ID
func (o *Order) CustomerID() string {
 return o.customerID
}

// Status 获取状态
func (o *Order) Status() OrderStatus {
 return o.status
}

// Total 获取总价
func (o *Order) Total() valueobject.Money {
 return o.total
}

// Items 获取订单项
func (o *Order) Items() []*entity.OrderItem {
 result := make([]*entity.OrderItem, len(o.items))
 for i, item := range o.items {
  result[i] = item
 }
 return result
}

// ShippingAddress 获取配送地址
func (o *Order) ShippingAddress() valueobject.Address {
 return o.shippingAddress
}

// AddItem 添加订单项
func (o *Order) AddItem(item *entity.OrderItem) error {
 if o.status != OrderStatusPending {
  return errors.New("cannot add items to non-pending order")
 }

 // 检查是否已存在相同产品
 for _, existing := range o.items {
  if existing.ProductID() == item.ProductID() {
   return errors.New("product already exists in order")
  }
 }

 o.items = append(o.items, item)
 o.recalculateTotal()
 o.updatedAt = time.Now()
 o.version++

 return nil
}

// RemoveItem 移除订单项
func (o *Order) RemoveItem(itemID entity.OrderItemID) error {
 if o.status != OrderStatusPending {
  return errors.New("cannot remove items from non-pending order")
 }

 for i, item := range o.items {
  if item.ID() == itemID {
   o.items = append(o.items[:i], o.items[i+1:]...)
   o.recalculateTotal()
   o.updatedAt = time.Now()
   o.version++
   return nil
  }
 }

 return errors.New("item not found")
}

// Place 提交订单
func (o *Order) Place() error {
 if o.status != OrderStatusPending {
  return errors.New("order is already placed")
 }
 if len(o.items) == 0 {
  return errors.New("order must have at least one item")
 }
 if o.total.IsZero() {
  return errors.New("order total cannot be zero")
 }

 o.updatedAt = time.Now()
 o.version++

 return nil
}

// Pay 支付订单
func (o *Order) Pay() error {
 if o.status != OrderStatusPending {
  return errors.New("only pending orders can be paid")
 }

 o.status = OrderStatusPaid
 o.updatedAt = time.Now()
 o.version++

 return nil
}

// Ship 发货
func (o *Order) Ship() error {
 if o.status != OrderStatusPaid {
  return errors.New("only paid orders can be shipped")
 }

 o.status = OrderStatusShipped
 o.updatedAt = time.Now()
 o.version++

 return nil
}

// Complete 完成订单
func (o *Order) Complete() error {
 if o.status != OrderStatusShipped {
  return errors.New("only shipped orders can be completed")
 }

 o.status = OrderStatusCompleted
 o.updatedAt = time.Now()
 o.version++

 return nil
}

// Cancel 取消订单
func (o *Order) Cancel() error {
 if o.status != OrderStatusPending && o.status != OrderStatusPaid {
  return errors.New("cannot cancel order in current status")
 }

 o.status = OrderStatusCancelled
 o.updatedAt = time.Now()
 o.version++

 return nil
}

// UpdateShippingAddress 更新配送地址
func (o *Order) UpdateShippingAddress(address valueobject.Address) error {
 if o.status != OrderStatusPending {
  return errors.New("can only update address for pending orders")
 }

 o.shippingAddress = address
 o.updatedAt = time.Now()
 o.version++

 return nil
}

// Version 获取版本号
func (o *Order) Version() int {
 return o.version
}

// CreatedAt 获取创建时间
func (o *Order) CreatedAt() time.Time {
 return o.createdAt
}

// UpdatedAt 获取更新时间
func (o *Order) UpdatedAt() time.Time {
 return o.updatedAt
}

// recalculateTotal 重新计算总价
func (o *Order) recalculateTotal() {
 total := valueobject.ZeroMoney(valueobject.CurrencyCNY)
 for _, item := range o.items {
  subtotal, _ := item.Subtotal()
  total, _ = total.Add(subtotal)
 }
 o.total = total
}
```

```go
// internal/order/domain/service/pricing_service.go
package service

import (
 "order/internal/domain/aggregate"
 "order/internal/domain/valueobject"
)

// PricingService 定价领域服务
type PricingService struct {
 discountRules []DiscountRule
}

// DiscountRule 折扣规则
type DiscountRule interface {
 Apply(order *aggregate.Order) (valueobject.Money, error)
}

// NewPricingService 创建定价服务
func NewPricingService() *PricingService {
 return &PricingService{
  discountRules: make([]DiscountRule, 0),
 }
}

// AddDiscountRule 添加折扣规则
func (s *PricingService) AddDiscountRule(rule DiscountRule) {
 s.discountRules = append(s.discountRules, rule)
}

// CalculateFinalTotal 计算最终总价（应用折扣）
func (s *PricingService) CalculateFinalTotal(order *aggregate.Order) (valueobject.Money, error) {
 baseTotal := order.Total()
 finalTotal := baseTotal

 for _, rule := range s.discountRules {
  discount, err := rule.Apply(order)
  if err != nil {
   return valueobject.Money{}, err
  }
  finalTotal, err = finalTotal.Subtract(discount)
  if err != nil {
   return valueobject.Money{}, err
  }
 }

 return finalTotal, nil
}

// PercentageDiscountRule 百分比折扣规则
type PercentageDiscountRule struct {
 Percentage float64
 MinAmount  valueobject.Money
}

// Apply 应用折扣
func (r *PercentageDiscountRule) Apply(order *aggregate.Order) (valueobject.Money, error) {
 total := order.Total()

 greater, err := total.GreaterThan(r.MinAmount)
 if err != nil {
  return valueobject.Money{}, err
 }

 if !greater {
  return valueobject.ZeroMoney(total.Currency()), nil
 }

 discountAmount := total.Amount() * r.Percentage / 100
 return valueobject.MustNewMoney(discountAmount, total.Currency()), nil
}
```

```go
// internal/order/domain/repository/order_repository.go
package repository

import (
 "context"
 "order/internal/domain/aggregate"
)

// OrderRepository 订单仓储接口
type OrderRepository interface {
 FindByID(ctx context.Context, id aggregate.OrderID) (*aggregate.Order, error)
 FindByCustomerID(ctx context.Context, customerID string) ([]*aggregate.Order, error)
 FindByStatus(ctx context.Context, status aggregate.OrderStatus) ([]*aggregate.Order, error)
 Save(ctx context.Context, order *aggregate.Order) error
 Update(ctx context.Context, order *aggregate.Order) error
 Delete(ctx context.Context, id aggregate.OrderID) error
}
```

```go
// internal/order/application/dto/order_dto.go
package dto

import "time"

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
 CustomerID      string      `json:"customer_id"`
 ShippingAddress AddressDTO  `json:"shipping_address"`
}

// AddressDTO 地址DTO
type AddressDTO struct {
 Province  string `json:"province"`
 City      string `json:"city"`
 District  string `json:"district"`
 Street    string `json:"street"`
 ZipCode   string `json:"zip_code"`
 Recipient string `json:"recipient"`
 Phone     string `json:"phone"`
}

// AddItemRequest 添加订单项请求
type AddItemRequest struct {
 ProductID string  `json:"product_id"`
 Name      string  `json:"name"`
 Quantity  int     `json:"quantity"`
 Price     float64 `json:"price"`
 Currency  string  `json:"currency"`
}

// OrderResponse 订单响应
type OrderResponse struct {
 ID              string        `json:"id"`
 CustomerID      string        `json:"customer_id"`
 Status          string        `json:"status"`
 Total           string        `json:"total"`
 Items           []ItemResponse `json:"items"`
 ShippingAddress AddressDTO    `json:"shipping_address"`
 CreatedAt       time.Time     `json:"created_at"`
 UpdatedAt       time.Time     `json:"updated_at"`
}

// ItemResponse 订单项响应
type ItemResponse struct {
 ID        string  `json:"id"`
 ProductID string  `json:"product_id"`
 Name      string  `json:"name"`
 Quantity  int     `json:"quantity"`
 Price     string  `json:"price"`
 Subtotal  string  `json:"subtotal"`
}
```

```go
// internal/order/application/service/order_app_service.go
package service

import (
 "context"
 "order/internal/application/dto"
 "order/internal/domain/aggregate"
 "order/internal/domain/entity"
 "order/internal/domain/repository"
 "order/internal/domain/valueobject"
 "time"

 "github.com/google/uuid"
)

// OrderApplicationService 订单应用服务
type OrderApplicationService struct {
 orderRepo repository.OrderRepository
}

// NewOrderApplicationService 创建应用服务
func NewOrderApplicationService(orderRepo repository.OrderRepository) *OrderApplicationService {
 return &OrderApplicationService{orderRepo: orderRepo}
}

// CreateOrder 创建订单
func (s *OrderApplicationService) CreateOrder(ctx context.Context, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
 address, err := valueobject.NewAddress(
  req.ShippingAddress.Province,
  req.ShippingAddress.City,
  req.ShippingAddress.District,
  req.ShippingAddress.Street,
  req.ShippingAddress.ZipCode,
  req.ShippingAddress.Recipient,
  req.ShippingAddress.Phone,
 )
 if err != nil {
  return nil, err
 }

 order, err := aggregate.NewOrder(
  aggregate.OrderID(uuid.New().String()),
  req.CustomerID,
  address,
 )
 if err != nil {
  return nil, err
 }

 if err := s.orderRepo.Save(ctx, order); err != nil {
  return nil, err
 }

 return s.toResponse(order), nil
}

// AddItem 添加订单项
func (s *OrderApplicationService) AddItem(ctx context.Context, orderID string, req dto.AddItemRequest) (*dto.OrderResponse, error) {
 order, err := s.orderRepo.FindByID(ctx, aggregate.OrderID(orderID))
 if err != nil {
  return nil, err
 }
 if order == nil {
  return nil, nil
 }

 currency := valueobject.Currency(req.Currency)
 price, err := valueobject.NewMoney(req.Price, currency)
 if err != nil {
  return nil, err
 }

 item, err := entity.NewOrderItem(
  entity.OrderItemID(uuid.New().String()),
  req.ProductID,
  req.Name,
  req.Quantity,
  price,
 )
 if err != nil {
  return nil, err
 }

 if err := order.AddItem(item); err != nil {
  return nil, err
 }

 if err := s.orderRepo.Update(ctx, order); err != nil {
  return nil, err
 }

 return s.toResponse(order), nil
}

// PayOrder 支付订单
func (s *OrderApplicationService) PayOrder(ctx context.Context, orderID string) (*dto.OrderResponse, error) {
 order, err := s.orderRepo.FindByID(ctx, aggregate.OrderID(orderID))
 if err != nil {
  return nil, err
 }
 if order == nil {
  return nil, nil
 }

 if err := order.Pay(); err != nil {
  return nil, err
 }

 if err := s.orderRepo.Update(ctx, order); err != nil {
  return nil, err
 }

 return s.toResponse(order), nil
}

// GetOrder 获取订单
func (s *OrderApplicationService) GetOrder(ctx context.Context, orderID string) (*dto.OrderResponse, error) {
 order, err := s.orderRepo.FindByID(ctx, aggregate.OrderID(orderID))
 if err != nil {
  return nil, err
 }
 if order == nil {
  return nil, nil
 }

 return s.toResponse(order), nil
}

func (s *OrderApplicationService) toResponse(order *aggregate.Order) *dto.OrderResponse {
 items := make([]dto.ItemResponse, len(order.Items()))
 for i, item := range order.Items() {
  subtotal, _ := item.Subtotal()
  items[i] = dto.ItemResponse{
   ID:        string(item.ID()),
   ProductID: item.ProductID(),
   Name:      item.Name(),
   Quantity:  item.Quantity(),
   Price:     item.Price().String(),
   Subtotal:  subtotal.String(),
  }
 }

 return &dto.OrderResponse{
  ID:         string(order.ID()),
  CustomerID: order.CustomerID(),
  Status:     order.Status().String(),
  Total:      order.Total().String(),
  Items:      items,
  ShippingAddress: dto.AddressDTO{
   Province:  order.ShippingAddress().Province(),
   City:      order.ShippingAddress().City(),
   District:  order.ShippingAddress().District(),
   Street:    order.ShippingAddress().Street(),
   ZipCode:   order.ShippingAddress().ZipCode(),
   Recipient: order.ShippingAddress().Recipient(),
   Phone:     order.ShippingAddress().Phone(),
  },
  CreatedAt: order.CreatedAt(),
  UpdatedAt: order.UpdatedAt(),
 }
}
```

```go
// internal/order/infrastructure/persistence/order_repository_impl.go
package persistence

import (
 "context"
 "order/internal/domain/aggregate"
 "order/internal/domain/entity"
 "order/internal/domain/repository"
 "order/internal/domain/valueobject"
 "sync"
)

// InMemoryOrderRepository 内存订单仓储
type InMemoryOrderRepository struct {
 mu     sync.RWMutex
 orders map[string]*aggregate.Order
}

// NewInMemoryOrderRepository 创建仓储
func NewInMemoryOrderRepository() repository.OrderRepository {
 return &InMemoryOrderRepository{
  orders: make(map[string]*aggregate.Order),
 }
}

func (r *InMemoryOrderRepository) FindByID(ctx context.Context, id aggregate.OrderID) (*aggregate.Order, error) {
 r.mu.RLock()
 defer r.mu.RUnlock()

 order, ok := r.orders[string(id)]
 if !ok {
  return nil, nil
 }
 return copyOrder(order), nil
}

func (r *InMemoryOrderRepository) FindByCustomerID(ctx context.Context, customerID string) ([]*aggregate.Order, error) {
 r.mu.RLock()
 defer r.mu.RUnlock()

 var result []*aggregate.Order
 for _, order := range r.orders {
  if order.CustomerID() == customerID {
   result = append(result, copyOrder(order))
  }
 }
 return result, nil
}

func (r *InMemoryOrderRepository) FindByStatus(ctx context.Context, status aggregate.OrderStatus) ([]*aggregate.Order, error) {
 r.mu.RLock()
 defer r.mu.RUnlock()

 var result []*aggregate.Order
 for _, order := range r.orders {
  if order.Status() == status {
   result = append(result, copyOrder(order))
  }
 }
 return result, nil
}

func (r *InMemoryOrderRepository) Save(ctx context.Context, order *aggregate.Order) error {
 r.mu.Lock()
 defer r.mu.Unlock()

 r.orders[string(order.ID())] = copyOrder(order)
 return nil
}

func (r *InMemoryOrderRepository) Update(ctx context.Context, order *aggregate.Order) error {
 r.mu.Lock()
 defer r.mu.Unlock()

 r.orders[string(order.ID())] = copyOrder(order)
 return nil
}

func (r *InMemoryOrderRepository) Delete(ctx context.Context, id aggregate.OrderID) error {
 r.mu.Lock()
 defer r.mu.Unlock()

 delete(r.orders, string(id))
 return nil
}

func copyOrder(o *aggregate.Order) *aggregate.Order {
 items := make([]*entity.OrderItem, len(o.Items()))
 for i, item := range o.Items() {
  items[i], _ = entity.NewOrderItem(item.ID(), item.ProductID(), item.Name(), item.Quantity(), item.Price())
 }

 order, _ := aggregate.NewOrder(o.ID(), o.CustomerID(), o.ShippingAddress())
 for _, item := range items {
  order.AddItem(item)
 }

 return order
}
```

### 7.5 反例说明

#### 反例1：贫血领域模型

```go
// ❌ 错误：领域对象只有getter/setter
type Order struct {
 ID     string
 Status string
 Total  float64
}

func (o *Order) GetID() string     { return o.ID }
func (o *Order) SetID(id string)   { o.ID = id }
func (o *Order) GetStatus() string { return o.Status }
func (o *Order) SetStatus(s string) { o.Status = s }

// 业务逻辑在服务中
func (s *OrderService) PlaceOrder(order *Order) error {
 if order.Status != "PENDING" {
  return errors.New("invalid status")
 }
 order.Status = "PLACED"
 return nil
}
```

**问题**：

- 领域对象退化为数据结构
- 业务逻辑散落在服务中
- 无法保证业务规则的一致性

#### 反例2：跨聚合直接引用

```go
// ❌ 错误：订单直接引用用户对象
type Order struct {
 ID     string
 User   *User  // 直接引用另一个聚合
 Items  []OrderItem
}
```

**问题**：

- 违反聚合边界原则
- 可能导致数据不一致
- 应该使用ID引用

### 7.6 选型指南

| 场景 | 建议 |
|------|------|
| **复杂业务领域** | DDD帮助理解和建模业务 |
| **长期演进项目** | 领域模型提供稳定的核心 |
| **大型团队协作** | 限界上下文划分团队边界 |
| **遗留系统改造** | DDD指导渐进式重构 |
| **简单CRUD** | 可能过度设计 |
| **团队不熟悉DDD** | 需要培训和学习成本 |

### 7.7 优缺点分析

| 优点 | 缺点 |
|------|------|
| ✅ 业务与技术对齐 | ❌ 学习曲线陡峭 |
| ✅ 清晰的领域边界 | ❌ 初始设计成本高 |
| ✅ 支持复杂业务 | ❌ 需要领域专家参与 |
| ✅ 易于演进和维护 | ❌ 简单项目可能过度设计 |
| ✅ 提高代码可读性 | ❌ 需要团队理解和实践 |

---

## 8. 微服务架构

### 8.1 概念定义

微服务架构（Microservices Architecture）是一种将应用程序构建为一组小型、独立服务的架构风格。每个服务运行在自己的进程中，通过轻量级机制（通常是HTTP/REST或消息队列）进行通信。

微服务的核心思想是：**单一职责、独立部署、松耦合**。

### 8.2 架构结构

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              API Gateway                                 │
│                     (路由、认证、限流、熔断)                               │
└─────────────────────────────────┬───────────────────────────────────────┘
                                  │
        ┌─────────────────────────┼─────────────────────────┐
        │                         │                         │
┌───────▼────────┐       ┌────────▼────────┐       ┌───────▼────────┐
│  用户服务       │       │   订单服务       │       │  库存服务       │
│                │       │                 │       │                │
│ ┌───────────┐  │       │ ┌───────────┐   │       │ ┌───────────┐  │
│ │ REST API  │  │       │ │ REST API  │   │       │ │ REST API  │  │
│ └─────┬─────┘  │       │ └─────┬─────┘   │       │ └─────┬─────┘  │
│       │        │       │       │         │       │       │        │
│ ┌─────▼─────┐  │       │ ┌─────▼─────┐   │       │ ┌─────▼─────┐  │
│ │  Service  │  │       │ │  Service  │   │       │ │  Service  │  │
│ └─────┬─────┘  │       │ └─────┬─────┘   │       │ └─────┬─────┘  │
│       │        │       │       │         │       │       │        │
│ ┌─────▼─────┐  │       │ ┌─────▼─────┐   │       │ ┌─────▼─────┐  │
│ │Repository │  │       │ │Repository │   │       │ │Repository │  │
│ └─────┬─────┘  │       │ └─────┬─────┘   │       │ └─────┬─────┘  │
│       │        │       │       │         │       │       │        │
│ ┌─────▼─────┐  │       │ ┌─────▼─────┐   │       │ ┌─────▼─────┐  │
│ │   DB      │  │       │ │   DB      │   │       │ │   DB      │  │
│ │ (User DB) │  │       │ │(Order DB) │   │       │ │(Stock DB) │  │
│ └───────────┘  │       │ └───────────┘   │       │ └───────────┘  │
└────────────────┘       └─────────────────┘       └────────────────┘
        │                         │                         │
        └─────────────────────────┼─────────────────────────┘
                                  │
                    ┌─────────────▼─────────────┐
                    │      Message Queue        │
                    │    (Event-driven comm)    │
                    └───────────────────────────┘
```

### 8.3 核心原则

| 原则 | 说明 |
|------|------|
| **单一职责** | 每个服务只负责一个业务能力 |
| **独立部署** | 服务可独立开发、测试、部署 |
| **数据隔离** | 每个服务有自己的数据库 |
| **松耦合** | 服务间通过API通信，不共享代码 |
| **容错设计** | 服务失败不影响整体系统 |

### 8.4 服务拆分策略

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         服务拆分策略                                     │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  1. 按业务能力拆分                                                       │
│     ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐                 │
│     │ 用户服务 │  │ 订单服务 │  │ 库存服务 │  │ 支付服务 │                 │
│     └─────────┘  └─────────┘  └─────────┘  └─────────┘                 │
│                                                                         │
│  2. 按子域拆分 (DDD)                                                     │
│     ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                 │
│     │  用户上下文  │  │  订单上下文  │  │  库存上下文  │                 │
│     │ (限界上下文) │  │ (限界上下文) │  │ (限界上下文) │                 │
│     └─────────────┘  └─────────────┘  └─────────────┘                 │
│                                                                         │
│  3. 按数据变更频率拆分                                                   │
│     ┌──────────────┐  ┌──────────────┐                                │
│     │ 高频变更服务  │  │ 低频变更服务  │                                │
│     │ (订单、库存)  │  │ (配置、类目)  │                                │
│     └──────────────┘  └──────────────┘                                │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 8.5 Go实现示例

#### 项目结构

```
microservices/
├── api-gateway/           # API网关
│   └── main.go
├── user-service/          # 用户服务
│   ├── cmd/
│   ├── internal/
│   └── go.mod
├── order-service/         # 订单服务
│   ├── cmd/
│   ├── internal/
│   └── go.mod
├── inventory-service/     # 库存服务
│   ├── cmd/
│   ├── internal/
│   └── go.mod
├── shared/                # 共享库
│   ├── event/
│   ├── middleware/
│   └── client/
└── docker-compose.yml
```

#### 完整代码实现

```go
// api-gateway/main.go
package main

import (
 "context"
 "encoding/json"
 "fmt"
 "io"
 "net/http"
 "net/http/httputil"
 "net/url"
 "strings"
 "time"
)

// ServiceConfig 服务配置
type ServiceConfig struct {
 Name string
 URL  string
}

// APIGateway API网关
type APIGateway struct {
 services map[string]*httputil.ReverseProxy
}

// NewAPIGateway 创建API网关
func NewAPIGateway(configs []ServiceConfig) *APIGateway {
 gateway := &APIGateway{
  services: make(map[string]*httputil.ReverseProxy),
 }

 for _, config := range configs {
  targetURL, _ := url.Parse(config.URL)
  gateway.services[config.Name] = httputil.NewSingleHostReverseProxy(targetURL)
 }

 return gateway
}

// ServeHTTP 实现http.Handler
func (g *APIGateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
 // 1. 认证中间件
 if !g.authenticate(r) {
  w.WriteHeader(http.StatusUnauthorized)
  json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
  return
 }

 // 2. 限流
 if !g.rateLimit(r) {
  w.WriteHeader(http.StatusTooManyRequests)
  json.NewEncoder(w).Encode(map[string]string{"error": "rate limit exceeded"})
  return
 }

 // 3. 路由到对应服务
 serviceName := g.extractServiceName(r.URL.Path)
 proxy, ok := g.services[serviceName]
 if !ok {
  w.WriteHeader(http.StatusNotFound)
  json.NewEncoder(w).Encode(map[string]string{"error": "service not found"})
  return
 }

 // 4. 转发请求
 proxy.ServeHTTP(w, r)
}

func (g *APIGateway) authenticate(r *http.Request) bool {
 // 简化实现：检查Authorization头
 token := r.Header.Get("Authorization")
 return token != "" && strings.HasPrefix(token, "Bearer ")
}

func (g *APIGateway) rateLimit(r *http.Request) bool {
 // 简化实现：总是通过
 return true
}

func (g *APIGateway) extractServiceName(path string) string {
 parts := strings.Split(path, "/")
 if len(parts) >= 2 {
  return parts[1]
 }
 return ""
}

func main() {
 configs := []ServiceConfig{
  {Name: "users", URL: "http://localhost:8001"},
  {Name: "orders", URL: "http://localhost:8002"},
  {Name: "inventory", URL: "http://localhost:8003"},
 }

 gateway := NewAPIGateway(configs)

 fmt.Println("API Gateway starting on :8080")
 http.ListenAndServe(":8080", gateway)
}
```

```go
// user-service/internal/service/user_service.go
package service

import (
 "context"
 "encoding/json"
 "errors"
 "net/http"
 "time"

 "github.com/google/uuid"
)

// User 用户实体
type User struct {
 ID        string    `json:"id"`
 Email     string    `json:"email"`
 Name      string    `json:"name"`
 CreatedAt time.Time `json:"created_at"`
}

// UserRepository 用户仓储接口
type UserRepository interface {
 FindByID(ctx context.Context, id string) (*User, error)
 FindByEmail(ctx context.Context, email string) (*User, error)
 Save(ctx context.Context, user *User) error
}

// UserService 用户服务
type UserService struct {
 repo UserRepository
}

// NewUserService 创建用户服务
func NewUserService(repo UserRepository) *UserService {
 return &UserService{repo: repo}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, email, name string) (*User, error) {
 if email == "" || name == "" {
  return nil, errors.New("email and name are required")
 }

 // 检查邮箱是否已存在
 existing, _ := s.repo.FindByEmail(ctx, email)
 if existing != nil {
  return nil, errors.New("email already exists")
 }

 user := &User{
  ID:        uuid.New().String(),
  Email:     email,
  Name:      name,
  CreatedAt: time.Now(),
 }

 if err := s.repo.Save(ctx, user); err != nil {
  return nil, err
 }

 return user, nil
}

// GetUser 获取用户
func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
 return s.repo.FindByID(ctx, id)
}

// UserHandler HTTP处理器
type UserHandler struct {
 service *UserService
}

// NewUserHandler 创建处理器
func NewUserHandler(service *UserService) *UserHandler {
 return &UserHandler{service: service}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
 switch r.Method {
 case http.MethodPost:
  h.createUser(w, r)
 case http.MethodGet:
  h.getUser(w, r)
 default:
  w.WriteHeader(http.StatusMethodNotAllowed)
 }
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
 var req struct {
  Email string `json:"email"`
  Name  string `json:"name"`
 }
 if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
  http.Error(w, err.Error(), http.StatusBadRequest)
  return
 }

 user, err := h.service.CreateUser(r.Context(), req.Email, req.Name)
 if err != nil {
  http.Error(w, err.Error(), http.StatusBadRequest)
  return
 }

 w.Header().Set("Content-Type", "application/json")
 w.WriteHeader(http.StatusCreated)
 json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
 id := r.URL.Path[len("/users/"):]
 user, err := h.service.GetUser(r.Context(), id)
 if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
 }
 if user == nil {
  http.Error(w, "user not found", http.StatusNotFound)
  return
 }

 w.Header().Set("Content-Type", "application/json")
 json.NewEncoder(w).Encode(user)
}
```

```go
// order-service/internal/service/order_service.go
package service

import (
 "bytes"
 "context"
 "encoding/json"
 "errors"
 "fmt"
 "net/http"
 "time"

 "github.com/google/uuid"
)

// Order 订单实体
type Order struct {
 ID         string      `json:"id"`
 UserID     string      `json:"user_id"`
 Items      []OrderItem `json:"items"`
 Total      float64     `json:"total"`
 Status     string      `json:"status"`
 CreatedAt  time.Time   `json:"created_at"`
}

// OrderItem 订单项
type OrderItem struct {
 ProductID string  `json:"product_id"`
 Name      string  `json:"name"`
 Quantity  int     `json:"quantity"`
 Price     float64 `json:"price"`
}

// OrderRepository 订单仓储接口
type OrderRepository interface {
 FindByID(ctx context.Context, id string) (*Order, error)
 FindByUserID(ctx context.Context, userID string) ([]*Order, error)
 Save(ctx context.Context, order *Order) error
}

// UserServiceClient 用户服务客户端
type UserServiceClient struct {
 baseURL string
 client  *http.Client
}

// NewUserServiceClient 创建用户服务客户端
func NewUserServiceClient(baseURL string) *UserServiceClient {
 return &UserServiceClient{
  baseURL: baseURL,
  client:  &http.Client{Timeout: 5 * time.Second},
 }
}

// GetUser 获取用户信息
func (c *UserServiceClient) GetUser(ctx context.Context, userID string) (*UserInfo, error) {
 req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/users/"+userID, nil)
 if err != nil {
  return nil, err
 }

 resp, err := c.client.Do(req)
 if err != nil {
  return nil, err
 }
 defer resp.Body.Close()

 if resp.StatusCode != http.StatusOK {
  return nil, fmt.Errorf("user service returned %d", resp.StatusCode)
 }

 var user UserInfo
 if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
  return nil, err
 }

 return &user, nil
}

// UserInfo 用户信息
type UserInfo struct {
 ID    string `json:"id"`
 Email string `json:"email"`
 Name  string `json:"name"`
}

// InventoryServiceClient 库存服务客户端
type InventoryServiceClient struct {
 baseURL string
 client  *http.Client
}

// NewInventoryServiceClient 创建库存服务客户端
func NewInventoryServiceClient(baseURL string) *InventoryServiceClient {
 return &InventoryServiceClient{
  baseURL: baseURL,
  client:  &http.Client{Timeout: 5 * time.Second},
 }
}

// CheckStock 检查库存
func (c *InventoryServiceClient) CheckStock(ctx context.Context, productID string, quantity int) (bool, error) {
 req, err := http.NewRequestWithContext(ctx, "GET",
  fmt.Sprintf("%s/inventory/check?product_id=%s&quantity=%d", c.baseURL, productID, quantity), nil)
 if err != nil {
  return false, err
 }

 resp, err := c.client.Do(req)
 if err != nil {
  return false, err
 }
 defer resp.Body.Close()

 return resp.StatusCode == http.StatusOK, nil
}

// ReserveStock 预留库存
func (c *InventoryServiceClient) ReserveStock(ctx context.Context, productID string, quantity int) error {
 reqBody, _ := json.Marshal(map[string]interface{}{
  "product_id": productID,
  "quantity":   quantity,
 })

 req, err := http.NewRequestWithContext(ctx, "POST",
  c.baseURL+"/inventory/reserve", bytes.NewReader(reqBody))
 if err != nil {
  return err
 }
 req.Header.Set("Content-Type", "application/json")

 resp, err := c.client.Do(req)
 if err != nil {
  return err
 }
 defer resp.Body.Close()

 if resp.StatusCode != http.StatusOK {
  return fmt.Errorf("failed to reserve stock: %d", resp.StatusCode)
 }

 return nil
}

// OrderService 订单服务
type OrderService struct {
 repo            OrderRepository
 userClient      *UserServiceClient
 inventoryClient *InventoryServiceClient
}

// NewOrderService 创建订单服务
func NewOrderService(
 repo OrderRepository,
 userClient *UserServiceClient,
 inventoryClient *InventoryServiceClient,
) *OrderService {
 return &OrderService{
  repo:            repo,
  userClient:      userClient,
  inventoryClient: inventoryClient,
 }
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(ctx context.Context, userID string, items []OrderItem) (*Order, error) {
 // 1. 验证用户
 user, err := s.userClient.GetUser(ctx, userID)
 if err != nil {
  return nil, fmt.Errorf("failed to get user: %w", err)
 }
 if user == nil {
  return nil, errors.New("user not found")
 }

 // 2. 检查库存并预留
 for _, item := range items {
  available, err := s.inventoryClient.CheckStock(ctx, item.ProductID, item.Quantity)
  if err != nil {
   return nil, fmt.Errorf("failed to check stock: %w", err)
  }
  if !available {
   return nil, fmt.Errorf("insufficient stock for product %s", item.ProductID)
  }

  if err := s.inventoryClient.ReserveStock(ctx, item.ProductID, item.Quantity); err != nil {
   return nil, fmt.Errorf("failed to reserve stock: %w", err)
  }
 }

 // 3. 计算总价
 total := 0.0
 for _, item := range items {
  total += item.Price * float64(item.Quantity)
 }

 // 4. 创建订单
 order := &Order{
  ID:        uuid.New().String(),
  UserID:    userID,
  Items:     items,
  Total:     total,
  Status:    "PENDING",
  CreatedAt: time.Now(),
 }

 if err := s.repo.Save(ctx, order); err != nil {
  return nil, err
 }

 return order, nil
}

// GetOrder 获取订单
func (s *OrderService) GetOrder(ctx context.Context, orderID string) (*Order, error) {
 return s.repo.FindByID(ctx, orderID)
}

// OrderHandler HTTP处理器
type OrderHandler struct {
 service *OrderService
}

// NewOrderHandler 创建处理器
func NewOrderHandler(service *OrderService) *OrderHandler {
 return &OrderHandler{service: service}
}

func (h *OrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
 switch r.Method {
 case http.MethodPost:
  h.createOrder(w, r)
 case http.MethodGet:
  h.getOrder(w, r)
 default:
  w.WriteHeader(http.StatusMethodNotAllowed)
 }
}

func (h *OrderHandler) createOrder(w http.ResponseWriter, r *http.Request) {
 var req struct {
  UserID string      `json:"user_id"`
  Items  []OrderItem `json:"items"`
 }
 if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
  http.Error(w, err.Error(), http.StatusBadRequest)
  return
 }

 order, err := h.service.CreateOrder(r.Context(), req.UserID, req.Items)
 if err != nil {
  http.Error(w, err.Error(), http.StatusBadRequest)
  return
 }

 w.Header().Set("Content-Type", "application/json")
 w.WriteHeader(http.StatusCreated)
 json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) getOrder(w http.ResponseWriter, r *http.Request) {
 id := r.URL.Path[len("/orders/"):]
 order, err := h.service.GetOrder(r.Context(), id)
 if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
 }
 if order == nil {
  http.Error(w, "order not found", http.StatusNotFound)
  return
 }

 w.Header().Set("Content-Type", "application/json")
 json.NewEncoder(w).Encode(order)
}
```

```go
// inventory-service/internal/service/inventory_service.go
package service

import (
 "context"
 "encoding/json"
 "errors"
 "net/http"
 "strconv"
 "sync"
)

// Inventory 库存实体
type Inventory struct {
 ProductID string `json:"product_id"`
 Quantity  int    `json:"quantity"`
 Reserved  int    `json:"reserved"`
}

// Available 可用库存
func (i *Inventory) Available() int {
 return i.Quantity - i.Reserved
}

// InventoryRepository 库存仓储接口
type InventoryRepository interface {
 FindByProductID(ctx context.Context, productID string) (*Inventory, error)
 Save(ctx context.Context, inventory *Inventory) error
}

// InMemoryInventoryRepository 内存库存仓储
type InMemoryInventoryRepository struct {
 mu         sync.RWMutex
 inventories map[string]*Inventory
}

// NewInMemoryInventoryRepository 创建仓储
func NewInMemoryInventoryRepository() *InMemoryInventoryRepository {
 return &InMemoryInventoryRepository{
  inventories: make(map[string]*Inventory),
 }
}

func (r *InMemoryInventoryRepository) FindByProductID(ctx context.Context, productID string) (*Inventory, error) {
 r.mu.RLock()
 defer r.mu.RUnlock()

 inv, ok := r.inventories[productID]
 if !ok {
  return nil, nil
 }
 return &Inventory{
  ProductID: inv.ProductID,
  Quantity:  inv.Quantity,
  Reserved:  inv.Reserved,
 }, nil
}

func (r *InMemoryInventoryRepository) Save(ctx context.Context, inventory *Inventory) error {
 r.mu.Lock()
 defer r.mu.Unlock()

 r.inventories[inventory.ProductID] = &Inventory{
  ProductID: inventory.ProductID,
  Quantity:  inventory.Quantity,
  Reserved:  inventory.Reserved,
 }
 return nil
}

// InventoryService 库存服务
type InventoryService struct {
 repo InventoryRepository
}

// NewInventoryService 创建库存服务
func NewInventoryService(repo InventoryRepository) *InventoryService {
 return &InventoryService{repo: repo}
}

// CheckStock 检查库存
func (s *InventoryService) CheckStock(ctx context.Context, productID string, quantity int) (bool, error) {
 inventory, err := s.repo.FindByProductID(ctx, productID)
 if err != nil {
  return false, err
 }
 if inventory == nil {
  return false, nil
 }
 return inventory.Available() >= quantity, nil
}

// ReserveStock 预留库存
func (s *InventoryService) ReserveStock(ctx context.Context, productID string, quantity int) error {
 inventory, err := s.repo.FindByProductID(ctx, productID)
 if err != nil {
  return err
 }
 if inventory == nil {
  return errors.New("product not found")
 }
 if inventory.Available() < quantity {
  return errors.New("insufficient stock")
 }

 inventory.Reserved += quantity
 return s.repo.Save(ctx, inventory)
}

// ReleaseStock 释放预留
func (s *InventoryService) ReleaseStock(ctx context.Context, productID string, quantity int) error {
 inventory, err := s.repo.FindByProductID(ctx, productID)
 if err != nil {
  return err
 }
 if inventory == nil {
  return errors.New("product not found")
 }
 if inventory.Reserved < quantity {
  return errors.New("cannot release more than reserved")
 }

 inventory.Reserved -= quantity
 return s.repo.Save(ctx, inventory)
}

// InventoryHandler HTTP处理器
type InventoryHandler struct {
 service *InventoryService
}

// NewInventoryHandler 创建处理器
func NewInventoryHandler(service *InventoryService) *InventoryHandler {
 return &InventoryHandler{service: service}
}

func (h *InventoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
 switch r.URL.Path {
 case "/inventory/check":
  h.checkStock(w, r)
 case "/inventory/reserve":
  h.reserveStock(w, r)
 case "/inventory/release":
  h.releaseStock(w, r)
 default:
  w.WriteHeader(http.StatusNotFound)
 }
}

func (h *InventoryHandler) checkStock(w http.ResponseWriter, r *http.Request) {
 productID := r.URL.Query().Get("product_id")
 quantity, _ := strconv.Atoi(r.URL.Query().Get("quantity"))

 available, err := h.service.CheckStock(r.Context(), productID, quantity)
 if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
 }

 if available {
  w.WriteHeader(http.StatusOK)
 } else {
  w.WriteHeader(http.StatusConflict)
 }
}

func (h *InventoryHandler) reserveStock(w http.ResponseWriter, r *http.Request) {
 var req struct {
  ProductID string `json:"product_id"`
  Quantity  int    `json:"quantity"`
 }
 if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
  http.Error(w, err.Error(), http.StatusBadRequest)
  return
 }

 if err := h.service.ReserveStock(r.Context(), req.ProductID, req.Quantity); err != nil {
  http.Error(w, err.Error(), http.StatusBadRequest)
  return
 }

 w.WriteHeader(http.StatusOK)
}

func (h *InventoryHandler) releaseStock(w http.ResponseWriter, r *http.Request) {
 var req struct {
  ProductID string `json:"product_id"`
  Quantity  int    `json:"quantity"`
 }
 if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
  http.Error(w, err.Error(), http.StatusBadRequest)
  return
 }

 if err := h.service.ReleaseStock(r.Context(), req.ProductID, req.Quantity); err != nil {
  http.Error(w, err.Error(), http.StatusBadRequest)
  return
 }

 w.WriteHeader(http.StatusOK)
}
```

### 8.6 数据一致性

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         分布式事务策略                                   │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  1. Saga模式（推荐）                                                     │
│     ┌─────────┐     ┌─────────┐     ┌─────────┐                        │
│     │ 订单服务 │────▶│ 库存服务 │────▶│ 支付服务 │                        │
│     │(创建订单)│     │(预留库存)│     │(处理支付)│                        │
│     └────┬────┘     └────┬────┘     └────┬────┘                        │
│          │               │               │                              │
│          │ 失败          │ 失败          │ 失败                          │
│          ▼               ▼               ▼                              │
│     ┌─────────┐     ┌─────────┐     ┌─────────┐                        │
│     │ 取消订单 │◀────│ 释放库存 │◀────│ 退款    │                        │
│     └─────────┘     └─────────┘     └─────────┘                        │
│                                                                         │
│  2. 最终一致性（事件驱动）                                                │
│     ┌─────────┐     ┌─────────┐     ┌─────────┐                        │
│     │ 订单服务 │────▶│ 消息队列 │────▶│ 库存服务 │                        │
│     │(发布事件)│     │(OrderCreated)│ (消费事件)│                        │
│     └─────────┘     └─────────┘     └─────────┘                        │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 8.7 反例说明

#### 反例1：共享数据库

```go
// ❌ 错误：多个服务共享数据库
type OrderService struct {
 db *sql.DB  // 与用户服务共享数据库
}

type UserService struct {
 db *sql.DB  // 同一个数据库连接
}
```

**问题**：

- 服务间紧耦合
- 无法独立部署
- 数据库变更影响多个服务

#### 反例2：同步调用链过长

```go
// ❌ 错误：过长的同步调用链
func CreateOrder() {
 user := userService.GetUser()      // 同步调用
 inventory := inventoryService.Check() // 同步调用
 payment := paymentService.Process()   // 同步调用
 notificationService.Send()            // 同步调用
}
```

**问题**：

- 响应时间长
- 级联失败风险
- 应该使用异步消息

### 8.8 选型指南

| 场景 | 建议 |
|------|------|
| **大型系统** | 微服务支持独立扩展 |
| **多团队协作** | 服务边界划分团队边界 |
| **需要独立部署** | 微服务支持独立发布 |
| **技术异构** | 不同服务可用不同技术栈 |
| **小型项目** | 可能过度设计 |
| **团队经验不足** | 需要DevOps和分布式经验 |

### 8.9 优缺点分析

| 优点 | 缺点 |
|------|------|
| ✅ 独立部署和扩展 | ❌ 分布式复杂性 |
| ✅ 技术异构 | ❌ 运维成本高 |
| ✅ 故障隔离 | ❌ 分布式事务困难 |
| ✅ 团队自治 | ❌ 网络延迟和故障 |
| ✅ 易于理解和维护 | ❌ 需要DevOps能力 |

---

## 9. Serverless架构

### 9.1 概念定义

Serverless架构（无服务器架构）是一种云计算执行模型，云提供商动态管理服务器资源的分配。开发者只需关注业务逻辑代码，无需管理服务器基础设施。

Serverless的核心思想是：**按执行付费，自动扩缩容，零服务器管理**。

### 9.2 架构结构

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         事件源 (Event Sources)                           │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐       │
│  │ HTTP    │  │ Timer   │  │ Queue   │  │ Storage │  │ Stream  │       │
│  │ Request │  │ Trigger │  │ Message │  │ Event   │  │ Event   │       │
│  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘       │
└───────┼────────────┼────────────┼────────────┼────────────┼───────────┘
        │            │            │            │            │
        └────────────┴────────────┴────────────┴────────────┘
                                   │
                    ┌──────────────▼──────────────┐
                    │      Function Compute       │
                    │      (函数计算平台)          │
                    │                             │
                    │  ┌───────────────────────┐  │
                    │  │   Function Runtime    │  │
                    │  │   (函数运行时)          │  │
                    │  │                       │  │
                    │  │  ┌─────────────────┐  │  │
                    │  │  │  Handler        │  │  │
                    │  │  │  (业务逻辑)      │  │  │
                    │  │  └─────────────────┘  │  │
                    │  └───────────────────────┘  │
                    └─────────────────────────────┘
                                   │
        ┌──────────────────────────┼──────────────────────────┐
        │                          │                          │
┌───────▼────────┐       ┌─────────▼──────────┐    ┌─────────▼──────────┐
│   Database     │       │   Object Storage   │    │   Cache Service    │
│  (DynamoDB/    │       │   (S3/OSS)         │    │   (Redis/ElastiCache)│
│   Cloud DB)    │       │                    │    │                    │
└────────────────┘       └────────────────────┘    └────────────────────┘
```

### 9.3 核心原则

| 原则 | 说明 |
|------|------|
| **事件驱动** | 函数由事件触发执行 |
| **无状态** | 函数应该是无状态的 |
| **短暂执行** | 函数执行时间应该短 |
| **细粒度** | 每个函数只做一件事 |
| **自动扩缩容** | 平台自动管理资源 |

### 9.4 冷启动优化

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         冷启动优化策略                                   │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  1. 预置并发 (Provisioned Concurrency)                                   │
│     ┌─────────────────────────────────────────────────────────────┐    │
│     │  保持一定数量的函数实例预热，减少冷启动延迟                      │    │
│     │  适用于：高并发、低延迟要求的场景                               │    │
│     └─────────────────────────────────────────────────────────────┘    │
│                                                                         │
│  2. 精简依赖                                                            │
│     ┌─────────────────────────────────────────────────────────────┐    │
│     │  - 只包含必要的依赖                                           │    │
│     │  - 使用轻量级库                                               │    │
│     │  - 避免大型框架                                               │    │
│     └─────────────────────────────────────────────────────────────┘    │
│                                                                         │
│  3. 延迟加载                                                            │
│     ┌─────────────────────────────────────────────────────────────┐    │
│     │  func handler() {                                           │    │
│     │      // 按需初始化，避免全局初始化                             │    │
│     │      db := initDB() // 在函数内初始化                         │    │
│     │  }                                                          │    │
│     └─────────────────────────────────────────────────────────────┘    │
│                                                                         │
│  4. 连接池复用                                                          │
│     ┌─────────────────────────────────────────────────────────────┐    │
│     │  var db *sql.DB // 全局变量，复用连接                          │    │
│     │  func init() { db = initDB() }                                │    │
│     └─────────────────────────────────────────────────────────────┘    │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 9.5 Go实现示例

#### 项目结构

```
serverless/
├── functions/              # 函数代码
│   ├── user-handler/       # 用户处理函数
│   │   ├── main.go
│   │   └── go.mod
│   ├── order-handler/      # 订单处理函数
│   │   ├── main.go
│   │   └── go.mod
│   └── notification-handler/ # 通知处理函数
│       ├── main.go
│       └── go.mod
├── shared/                 # 共享代码
│   ├── db/
│   │   └── connection.go
│   ├── models/
│   │   └── models.go
│   └── utils/
│       └── response.go
├── infrastructure/         # 基础设施配置
│   ├── terraform/
│   └── serverless.yml
└── Makefile
```

#### 完整代码实现

```go
// shared/models/models.go
package models

import "time"

// User 用户模型
type User struct {
 ID        string    `json:"id" db:"id"`
 Email     string    `json:"email" db:"email"`
 Name      string    `json:"name" db:"name"`
 CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Order 订单模型
type Order struct {
 ID        string    `json:"id" db:"id"`
 UserID    string    `json:"user_id" db:"user_id"`
 Amount    float64   `json:"amount" db:"amount"`
 Status    string    `json:"status" db:"status"`
 CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// APIGatewayRequest API网关请求
type APIGatewayRequest struct {
 HTTPMethod            string            `json:"httpMethod"`
 Path                  string            `json:"path"`
 QueryStringParameters map[string]string `json:"queryStringParameters"`
 Headers               map[string]string `json:"headers"`
 Body                  string            `json:"body"`
 PathParameters        map[string]string `json:"pathParameters"`
}

// APIGatewayResponse API网关响应
type APIGatewayResponse struct {
 StatusCode int               `json:"statusCode"`
 Headers    map[string]string `json:"headers"`
 Body       string            `json:"body"`
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(body interface{}) APIGatewayResponse {
 jsonBody, _ := json.Marshal(body)
 return APIGatewayResponse{
  StatusCode: 200,
  Headers: map[string]string{
   "Content-Type": "application/json",
  },
  Body: string(jsonBody),
 }
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(statusCode int, message string) APIGatewayResponse {
 body, _ := json.Marshal(map[string]string{"error": message})
 return APIGatewayResponse{
  StatusCode: statusCode,
  Headers: map[string]string{
   "Content-Type": "application/json",
  },
  Body: string(body),
 }
}
```

```go
// shared/db/connection.go
package db

import (
 "database/sql"
 "fmt"
 "os"
 "sync"

 _ "github.com/lib/pq"
)

var (
 db     *sql.DB
 once   sync.Once
 dbErr  error
)

// GetDB 获取数据库连接（单例模式）
func GetDB() (*sql.DB, error) {
 once.Do(func() {
  db, dbErr = initDB()
 })
 return db, dbErr
}

func initDB() (*sql.DB, error) {
 dsn := fmt.Sprintf(
  "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
  getEnv("DB_HOST", "localhost"),
  getEnv("DB_PORT", "5432"),
  getEnv("DB_USER", "postgres"),
  getEnv("DB_PASSWORD", ""),
  getEnv("DB_NAME", "myapp"),
 )

 return sql.Open("postgres", dsn)
}

func getEnv(key, defaultValue string) string {
 if value := os.Getenv(key); value != "" {
  return value
 }
 return defaultValue
}
```

```go
// functions/user-handler/main.go
package main

import (
 "context"
 "encoding/json"
 "log"
 "serverless/shared/db"
 "serverless/shared/models"

 "github.com/aws/aws-lambda-go/events"
 "github.com/aws/aws-lambda-go/lambda"
 "github.com/google/uuid"
)

// Handler Lambda处理函数
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
 log.Printf("Request: %s %s", request.HTTPMethod, request.Path)

 switch request.HTTPMethod {
 case "GET":
  return handleGetUser(ctx, request)
 case "POST":
  return handleCreateUser(ctx, request)
 case "PUT":
  return handleUpdateUser(ctx, request)
 case "DELETE":
  return handleDeleteUser(ctx, request)
 default:
  return models.NewErrorResponse(405, "method not allowed"), nil
 }
}

func handleGetUser(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
 userID := request.PathParameters["id"]
 if userID == "" {
  return models.NewErrorResponse(400, "user id is required"), nil
 }

 database, err := db.GetDB()
 if err != nil {
  log.Printf("Database error: %v", err)
  return models.NewErrorResponse(500, "internal server error"), nil
 }

 var user models.User
 err = database.QueryRowContext(ctx,
  "SELECT id, email, name, created_at FROM users WHERE id = $1", userID).
  Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt)

 if err != nil {
  if err == sql.ErrNoRows {
   return models.NewErrorResponse(404, "user not found"), nil
  }
  log.Printf("Query error: %v", err)
  return models.NewErrorResponse(500, "internal server error"), nil
 }

 return models.NewSuccessResponse(user), nil
}

func handleCreateUser(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
 var req struct {
  Email string `json:"email"`
  Name  string `json:"name"`
 }

 if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
  return models.NewErrorResponse(400, "invalid request body"), nil
 }

 if req.Email == "" || req.Name == "" {
  return models.NewErrorResponse(400, "email and name are required"), nil
 }

 database, err := db.GetDB()
 if err != nil {
  log.Printf("Database error: %v", err)
  return models.NewErrorResponse(500, "internal server error"), nil
 }

 user := models.User{
  ID:        uuid.New().String(),
  Email:     req.Email,
  Name:      req.Name,
  CreatedAt: time.Now(),
 }

 _, err = database.ExecContext(ctx,
  "INSERT INTO users (id, email, name, created_at) VALUES ($1, $2, $3, $4)",
  user.ID, user.Email, user.Name, user.CreatedAt)

 if err != nil {
  log.Printf("Insert error: %v", err)
  return models.NewErrorResponse(500, "failed to create user"), nil
 }

 return models.NewSuccessResponse(user), nil
}

func handleUpdateUser(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
 userID := request.PathParameters["id"]
 if userID == "" {
  return models.NewErrorResponse(400, "user id is required"), nil
 }

 var req struct {
  Name string `json:"name"`
 }

 if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
  return models.NewErrorResponse(400, "invalid request body"), nil
 }

 database, err := db.GetDB()
 if err != nil {
  log.Printf("Database error: %v", err)
  return models.NewErrorResponse(500, "internal server error"), nil
 }

 _, err = database.ExecContext(ctx,
  "UPDATE users SET name = $1 WHERE id = $2",
  req.Name, userID)

 if err != nil {
  log.Printf("Update error: %v", err)
  return models.NewErrorResponse(500, "failed to update user"), nil
 }

 return models.NewSuccessResponse(map[string]string{"status": "updated"}), nil
}

func handleDeleteUser(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
 userID := request.PathParameters["id"]
 if userID == "" {
  return models.NewErrorResponse(400, "user id is required"), nil
 }

 database, err := db.GetDB()
 if err != nil {
  log.Printf("Database error: %v", err)
  return models.NewErrorResponse(500, "internal server error"), nil
 }

 _, err = database.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID)
 if err != nil {
  log.Printf("Delete error: %v", err)
  return models.NewErrorResponse(500, "failed to delete user"), nil
 }

 return models.NewSuccessResponse(map[string]string{"status": "deleted"}), nil
}

func main() {
 lambda.Start(Handler)
}
```

```go
// functions/order-handler/main.go
package main

import (
 "context"
 "encoding/json"
 "log"
 "serverless/shared/db"
 "serverless/shared/models"

 "github.com/aws/aws-lambda-go/events"
 "github.com/aws/aws-lambda-go/lambda"
 "github.com/google/uuid"
)

// Handler Lambda处理函数
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
 log.Printf("Order Handler: %s %s", request.HTTPMethod, request.Path)

 switch request.HTTPMethod {
 case "GET":
  return handleGetOrder(ctx, request)
 case "POST":
  return handleCreateOrder(ctx, request)
 case "PUT":
  if request.PathParameters["action"] == "pay" {
   return handlePayOrder(ctx, request)
  }
  return models.NewErrorResponse(400, "invalid action"), nil
 default:
  return models.NewErrorResponse(405, "method not allowed"), nil
 }
}

func handleGetOrder(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
 orderID := request.PathParameters["id"]
 if orderID == "" {
  return models.NewErrorResponse(400, "order id is required"), nil
 }

 database, err := db.GetDB()
 if err != nil {
  return models.NewErrorResponse(500, "internal server error"), nil
 }

 var order models.Order
 err = database.QueryRowContext(ctx,
  "SELECT id, user_id, amount, status, created_at FROM orders WHERE id = $1", orderID).
  Scan(&order.ID, &order.UserID, &order.Amount, &order.Status, &order.CreatedAt)

 if err != nil {
  if err == sql.ErrNoRows {
   return models.NewErrorResponse(404, "order not found"), nil
  }
  return models.NewErrorResponse(500, "internal server error"), nil
 }

 return models.NewSuccessResponse(order), nil
}

func handleCreateOrder(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
 var req struct {
  UserID string  `json:"user_id"`
  Amount float64 `json:"amount"`
 }

 if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
  return models.NewErrorResponse(400, "invalid request body"), nil
 }

 if req.UserID == "" || req.Amount <= 0 {
  return models.NewErrorResponse(400, "user_id and amount are required"), nil
 }

 database, err := db.GetDB()
 if err != nil {
  return models.NewErrorResponse(500, "internal server error"), nil
 }

 order := models.Order{
  ID:        uuid.New().String(),
  UserID:    req.UserID,
  Amount:    req.Amount,
  Status:    "PENDING",
  CreatedAt: time.Now(),
 }

 _, err = database.ExecContext(ctx,
  "INSERT INTO orders (id, user_id, amount, status, created_at) VALUES ($1, $2, $3, $4, $5)",
  order.ID, order.UserID, order.Amount, order.Status, order.CreatedAt)

 if err != nil {
  log.Printf("Insert error: %v", err)
  return models.NewErrorResponse(500, "failed to create order"), nil
 }

 return models.NewSuccessResponse(order), nil
}

func handlePayOrder(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
 orderID := request.PathParameters["id"]
 if orderID == "" {
  return models.NewErrorResponse(400, "order id is required"), nil
 }

 database, err := db.GetDB()
 if err != nil {
  return models.NewErrorResponse(500, "internal server error"), nil
 }

 _, err = database.ExecContext(ctx,
  "UPDATE orders SET status = 'PAID' WHERE id = $1 AND status = 'PENDING'",
  orderID)

 if err != nil {
  return models.NewErrorResponse(500, "failed to pay order"), nil
 }

 return models.NewSuccessResponse(map[string]string{"status": "paid"}), nil
}

func main() {
 lambda.Start(Handler)
}
```

```go
// functions/notification-handler/main.go
package main

import (
 "context"
 "encoding/json"
 "fmt"
 "log"
 "os"

 "github.com/aws/aws-lambda-go/events"
 "github.com/aws/aws-lambda-go/lambda"
)

// SQSMessage SQS消息
type SQSMessage struct {
 Type      string          `json:"type"`
 UserID    string          `json:"user_id"`
 OrderID   string          `json:"order_id"`
 Amount    float64         `json:"amount"`
 Timestamp string          `json:"timestamp"`
}

// Handler SQS事件处理函数
func Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
 for _, message := range sqsEvent.Records {
  log.Printf("Processing message: %s", message.MessageId)

  var sqsMsg SQSMessage
  if err := json.Unmarshal([]byte(message.Body), &sqsMsg); err != nil {
   log.Printf("Failed to parse message: %v", err)
   continue
  }

  if err := processNotification(ctx, sqsMsg); err != nil {
   log.Printf("Failed to process notification: %v", err)
   continue
  }

  log.Printf("Successfully processed message: %s", message.MessageId)
 }

 return nil
}

func processNotification(ctx context.Context, msg SQSMessage) error {
 switch msg.Type {
 case "ORDER_CREATED":
  return sendOrderCreatedNotification(ctx, msg)
 case "ORDER_PAID":
  return sendOrderPaidNotification(ctx, msg)
 default:
  log.Printf("Unknown message type: %s", msg.Type)
  return nil
 }
}

func sendOrderCreatedNotification(ctx context.Context, msg SQSMessage) error {
 // 发送订单创建通知（邮件/短信/推送）
 message := fmt.Sprintf("您的订单 %s 已创建，金额: %.2f", msg.OrderID, msg.Amount)
 log.Printf("Sending notification: %s", message)

 // 实际实现：调用邮件服务、短信服务或推送服务
 // sendEmail(msg.UserID, message)
 // sendSMS(msg.UserID, message)

 return nil
}

func sendOrderPaidNotification(ctx context.Context, msg SQSMessage) error {
 message := fmt.Sprintf("您的订单 %s 已支付成功，金额: %.2f", msg.OrderID, msg.Amount)
 log.Printf("Sending notification: %s", message)
 return nil
}

// TimerHandler 定时触发处理函数
func TimerHandler(ctx context.Context, event events.CloudWatchEvent) error {
 log.Printf("Timer triggered: %s", event.Time)

 // 执行定时任务，如：清理过期数据、生成报表等
 if err := cleanupExpiredData(ctx); err != nil {
  return err
 }

 return nil
}

func cleanupExpiredData(ctx context.Context) error {
 log.Println("Cleaning up expired data...")
 // 实现数据清理逻辑
 return nil
}

func main() {
 // 根据环境变量决定使用哪个处理函数
 handlerType := os.Getenv("HANDLER_TYPE")
 switch handlerType {
 case "sqs":
  lambda.Start(Handler)
 case "timer":
  lambda.Start(TimerHandler)
 default:
  lambda.Start(Handler)
 }
}
```

```yaml
# infrastructure/serverless.yml
service: my-serverless-app

provider:
  name: aws
  runtime: provided.al2
  stage: ${opt:stage, 'dev'}
  region: ${opt:region, 'ap-northeast-1'}
  environment:
    DB_HOST: ${env:DB_HOST}
    DB_PORT: ${env:DB_PORT}
    DB_USER: ${env:DB_USER}
    DB_PASSWORD: ${env:DB_PASSWORD}
    DB_NAME: ${env:DB_NAME}

package:
  individually: true

functions:
  # 用户处理函数
  userHandler:
    handler: functions/user-handler/bootstrap
    events:
      - http:
          path: users/{id}
          method: get
      - http:
          path: users
          method: post
      - http:
          path: users/{id}
          method: put
      - http:
          path: users/{id}
          method: delete
    package:
      artifact: functions/user-handler/user-handler.zip

  # 订单处理函数
  orderHandler:
    handler: functions/order-handler/bootstrap
    events:
      - http:
          path: orders/{id}
          method: get
      - http:
          path: orders
          method: post
      - http:
          path: orders/{id}/pay
          method: put
    package:
      artifact: functions/order-handler/order-handler.zip

  # 通知处理函数 (SQS)
  notificationHandler:
    handler: functions/notification-handler/bootstrap
    events:
      - sqs:
          arn:
            Fn::GetAtt:
              - NotificationQueue
              - Arn
    package:
      artifact: functions/notification-handler/notification-handler.zip

  # 定时任务
  cleanupTask:
    handler: functions/notification-handler/bootstrap
    environment:
      HANDLER_TYPE: timer
    events:
      - schedule: rate(1 day)
    package:
      artifact: functions/notification-handler/notification-handler.zip

resources:
  Resources:
    NotificationQueue:
      Type: AWS::SQS::Queue
      Properties:
        QueueName: ${self:service}-${self:provider.stage}-notifications
```

### 9.6 反例说明

#### 反例1：函数过大

```go
// ❌ 错误：一个函数处理所有逻辑
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
 // 处理用户、订单、库存、支付所有逻辑
 // 函数过大，难以维护
}
```

**问题**：

- 函数职责不清晰
- 难以测试和维护
- 冷启动时间长

#### 反例2：在函数内保持状态

```go
// ❌ 错误：在函数内保持状态
var counter int

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
 counter++  // 不可靠，函数实例可能被销毁
 return ...
}
```

**问题**：

- 函数应该是无状态的
- 状态应该存储在外部服务（数据库、缓存）

### 9.7 选型指南

| 场景 | 建议 |
|------|------|
| **事件驱动工作负载** | Serverless天然适合 |
| **可变流量** | 自动扩缩容，按量付费 |
| **快速原型** | 无需管理基础设施 |
| **定时任务** | 内置定时触发器 |
| **长时间运行** | 不适合，有执行时间限制 |
| **需要低延迟** | 考虑冷启动影响 |

### 9.8 优缺点分析

| 优点 | 缺点 |
|------|------|
| ✅ 无服务器管理 | ❌ 冷启动延迟 |
| ✅ 自动扩缩容 | ❌ 执行时间限制 |
| ✅ 按量付费 | ❌ 调试困难 |
| ✅ 快速部署 | ❌ 供应商锁定 |
| ✅ 高可用 | ❌ 状态管理复杂 |

---

## 10. 架构选型指南

### 10.1 项目规模与架构选择

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        项目规模与架构选择                                │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  小型项目 (1-3人)                                                        │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  推荐：简单分层架构 / MVC                                         │   │
│  │  - 快速开发，低学习成本                                           │   │
│  │  - 单层或三层架构足够                                             │   │
│  │  - 避免过度设计                                                   │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  中型项目 (3-10人)                                                       │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  推荐：分层架构 + DDD / 六边形架构                                 │   │
│  │  - 清晰的职责分离                                                 │   │
│  │  - 支持并行开发                                                   │   │
│  │  - 便于测试和维护                                                 │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  大型项目 (10+人)                                                        │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  推荐：微服务 + DDD + 事件驱动                                    │   │
│  │  - 团队自治，独立部署                                             │   │
│  │  - 限界上下文划分团队边界                                         │   │
│  │  - 支持技术异构                                                   │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 10.2 团队能力与架构复杂度

| 团队能力 | 推荐架构 | 说明 |
|----------|----------|------|
| **初级团队** | 简单分层 | 学习曲线低，快速上手 |
| **中级团队** | 分层+DDD | 逐步引入领域驱动设计 |
| **高级团队** | 六边形/洋葱 | 充分发挥架构优势 |
| **专家团队** | 微服务+事件驱动 | 处理复杂分布式场景 |

### 10.3 演进式架构

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        演进式架构路径                                    │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  阶段1: 单体应用                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  单层架构 → 三层架构                                               │   │
│  │  快速验证业务，积累领域知识                                         │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                              │                                          │
│                              ▼                                          │
│  阶段2: 模块化单体                                                       │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  引入DDD概念 → 限界上下文 → 模块分离                               │   │
│  │  代码层面分离，为微服务做准备                                       │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                              │                                          │
│                              ▼                                          │
│  阶段3: 微服务拆分                                                       │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  按限界上下文拆分 → 独立部署 → 事件驱动                            │   │
│  │  逐步拆分，避免大爆炸式重构                                         │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                              │                                          │
│                              ▼                                          │
│  阶段4: 优化与治理                                                       │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  服务网格 → 可观测性 → 混沌工程                                    │   │
│  │  提升系统稳定性和可维护性                                           │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 10.4 反模式识别

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          架构反模式                                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  1. 大泥球 (Big Ball of Mud)                                            │
│     ┌─────────────────────────────────────────────────────────────┐    │
│     │  特征：                                                      │    │
│     │  - 没有清晰的架构                                            │    │
│     │  - 代码随意堆砌                                              │    │
│     │  - 高度耦合，难以维护                                        │    │
│     │  解决：引入分层，逐步重构                                    │    │
│     └─────────────────────────────────────────────────────────────┘    │
│                                                                         │
│  2. 过度工程 (Over-Engineering)                                         │
│     ┌─────────────────────────────────────────────────────────────┐    │
│     │  特征：                                                      │    │
│     │  - 简单项目使用复杂架构                                      │    │
│     │  - 过多的抽象层                                              │    │
│     │  - 开发效率低下                                              │    │
│     │  解决：根据项目规模选择合适架构                              │    │
│     └─────────────────────────────────────────────────────────────┘    │
│                                                                         │
│  3. 分布式单体 (Distributed Monolith)                                   │
│     ┌─────────────────────────────────────────────────────────────┐    │
│     │  特征：                                                      │    │
│     │  - 服务间高度耦合                                            │    │
│     │  - 共享数据库                                                │    │
│     │  - 无法独立部署                                              │    │
│     │  解决：真正的服务拆分，数据隔离                              │    │
│     └─────────────────────────────────────────────────────────────┘    │
│                                                                         │
│  4. 贫血领域模型 (Anemic Domain Model)                                  │
│     ┌─────────────────────────────────────────────────────────────┐    │
│     │  特征：                                                      │    │
│     │  - 领域对象只有数据没有行为                                  │    │
│     │  - 业务逻辑散落在服务中                                      │    │
│     │  - 违反面向对象原则                                          │    │
│     │  解决：将行为放入领域对象                                    │    │
│     └─────────────────────────────────────────────────────────────┘    │
│                                                                         │
│  5. 循环依赖 (Circular Dependency)                                      │
│     ┌─────────────────────────────────────────────────────────────┐    │
│     │  特征：                                                      │    │
│     │  - 层与层之间循环依赖                                        │    │
│     │  - 服务间循环调用                                            │    │
│     │  - 难以理解和测试                                            │    │
│     │  解决：依赖倒置，引入接口                                    │    │
│     └─────────────────────────────────────────────────────────────┘    │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 10.5 架构决策矩阵

| 场景 | 推荐架构 | 备选架构 | 避免 |
|------|----------|----------|------|
| **快速原型** | 单层/三层 | MVC | 微服务、DDD |
| **小型Web应用** | 三层架构 | 六边形 | 微服务 |
| **中型业务系统** | 分层+DDD | 洋葱架构 | 纯微服务 |
| **大型企业应用** | 微服务+DDD | 六边形+事件驱动 | 单体 |
| **高并发读场景** | CQRS | 缓存+读写分离 | 单一模型 |
| **审计要求严格** | 事件溯源 | CQRS+事件溯源 | 状态存储 |
| **事件驱动系统** | Serverless | 微服务+MQ | 同步调用 |
| **多变流量** | Serverless | 容器化 | 固定资源 |

### 10.6 架构选型决策树

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        架构选型决策树                                    │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  开始                                                                    │
│   │                                                                     │
│   ▼                                                                     │
│  项目规模?                                                              │
│   ├── 小型 (1-3人) ──▶ 简单分层架构                                     │
│   │                                                                     │
│   ├── 中型 (3-10人)                                                     │
│   │   │                                                                 │
│   │   ▼                                                                 │
│   │  业务复杂度?                                                        │
│   │   ├── 简单 ──▶ 三层架构                                             │
│   │   └── 复杂 ──▶ 分层 + DDD                                           │
│   │                                                                     │
│   └── 大型 (10+人)                                                      │
│       │                                                                 │
│       ▼                                                                 │
│      需要独立部署?                                                      │
│       ├── 否 ──▶ 模块化单体 + DDD                                       │
│       │                                                                 │
│       └── 是                                                            │
│           │                                                             │
│           ▼                                                             │
│          数据一致性要求?                                                │
│           ├── 强一致 ──▶ 微服务 + Saga                                  │
│           └── 最终一致 ──▶ 微服务 + 事件驱动                            │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 10.7 架构组合建议

| 架构组合 | 适用场景 | 优势 |
|----------|----------|------|
| **分层 + DDD** | 中型业务系统 | 清晰的领域边界，易于维护 |
| **六边形 + DDD** | 需要多接口支持 | 核心领域独立，易于测试 |
| **微服务 + 事件驱动** | 大型分布式系统 | 松耦合，高可扩展 |
| **CQRS + 事件溯源** | 审计、复杂查询 | 完整历史，高性能读 |
| **Serverless + 事件驱动** | 可变流量，事件处理 | 自动扩缩容，低成本 |

### 10.8 实施建议

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          架构实施建议                                    │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  1. 渐进式引入                                                           │
│     - 不要一次性引入所有架构模式                                         │
│     - 从简单开始，逐步演进                                               │
│     - 根据实际需求调整                                                   │
│                                                                         │
│  2. 团队培训                                                             │
│     - 确保团队理解架构原则                                               │
│     - 代码审查确保规范执行                                               │
│     - 建立架构决策记录 (ADR)                                             │
│                                                                         │
│  3. 持续重构                                                             │
│     - 定期评估架构健康状况                                               │
│     - 及时消除技术债务                                                   │
│     - 保持架构与业务对齐                                                 │
│                                                                         │
│  4. 度量与监控                                                           │
│     - 建立架构健康度指标                                                 │
│     - 监控关键性能指标                                                   │
│     - 基于数据做决策                                                     │
│                                                                         │
│  5. 文档与沟通                                                           │
│     - 维护架构文档                                                       │
│     - 定期架构评审                                                       │
│     - 跨团队沟通对齐                                                     │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 总结

本文档全面梳理了Go语言项目中常用的架构设计模型，从经典的分层架构到现代的微服务和Serverless架构。每种架构都有其适用的场景和优缺点，选择合适的架构需要综合考虑项目规模、团队能力、业务复杂度等因素。

### 核心要点回顾

1. **分层架构**：经典、简单，适合大多数项目
2. **六边形架构**：强调端口与适配器，适合需要多接口支持的场景
3. **洋葱架构**：以领域为核心，适合DDD项目
4. **清洁架构**：强调依赖规则，适合大型企业应用
5. **CQRS**：读写分离，适合读多写少的场景
6. **事件溯源**：完整历史记录，适合审计要求严格的场景
7. **DDD**：以领域为核心，适合复杂业务系统
8. **微服务**：独立部署，适合大型分布式系统
9. **Serverless**：按量付费，适合可变流量的场景

### 架构选型原则

- **合适优于流行**：选择适合项目需求的架构，而不是最流行的
- **简单优于复杂**：在满足需求的前提下，选择最简单的方案
- **演进优于完美**：架构是演进的，不要追求一开始就完美
- **团队优于技术**：考虑团队能力，选择团队能够驾驭的架构

---

*文档版本: 1.0*
*最后更新: 2024年*

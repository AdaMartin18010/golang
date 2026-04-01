# Go Clean Architecture 框架 - 目录结构全面梳理方案

**日期**: 2026-03-17
**参考标准**: [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
**架构**: Clean Architecture + 标准 Go 项目布局

---

## 一、当前项目结构分析

### 📊 项目规模统计

| 指标 | 数量 |
|------|------|
| Go 源文件总数 | 691 |
| 测试文件数 | 246 |
| 目录数 | 150+ |
| 代码行数 | 约 30,000+ |

### 📁 当前目录结构

```
golang/
├── .githooks/          # Git hooks（符合标准）
├── .github/            # GitHub 配置（符合标准）
├── .vscode/            # VS Code 配置（符合标准）
├── api/                # API 定义（符合标准）
│   ├── asyncapi/
│   ├── graphql/
│   └── openapi/
├── archive/            # ⚠️ 归档文件（应移除）
├── benchmarks/         # 基准测试（符合标准）
├── cmd/                # 主程序入口（符合标准）
│   ├── cli/
│   ├── graphql-server/
│   ├── grpc-server/
│   ├── mqtt-client/
│   ├── server/
│   └── temporal-worker/
├── configs/            # 配置模板（符合标准）
├── deployments/        # 部署配置（符合标准）
│   ├── docker/
│   ├── kubernetes/
│   └── wasm/
├── docs/               # 文档（⚠️ 过度膨胀）
│   ├── adr/
│   ├── advanced/
│   ├── ai-native-observability-analysis/
│   ├── architecture/
│   ├── archive/        # ⚠️ 重复归档
│   ├── comprehensive-analysis/
│   ├── ... (24 个子目录)
├── examples/           # 示例（符合标准，但内容过多）
├── internal/           # 私有代码（⚠️ 结构复杂）
│   ├── application/    # 应用层
│   ├── config/         # 配置
│   ├── domain/         # 领域层
│   ├── framework/      # 框架层
│   ├── infrastructure/ # 基础设施层
│   ├── interfaces/     # 接口层
│   ├── security/       # 安全
│   ├── types/          # 类型
│   └── utils/          # 工具
├── migrations/         # 数据库迁移（符合标准）
├── pkg/                # 公共库（⚠️ 过度膨胀）
│   ├── auth/
│   ├── concurrency/
│   ├── control/
│   ├── converter/
│   ├── database/
│   ├── errors/
│   ├── eventbus/
│   ├── health/
│   ├── http/
│   ├── http3/
│   ├── loadbalancer/
│   ├── logger/
│   ├── memory/
│   ├── observability/
│   ├── rbac/
│   ├── reflect/
│   ├── registry/
│   ├── sampling/
│   ├── security/
│   ├── tracing/
│   ├── transaction/
│   ├── utils/          # ⚠️ 50+ 子模块
│   └── validator/
├── scripts/            # 脚本（符合标准）
├── test/               # 测试（符合标准）
├── tools/              # 工具（符合标准）
├── go.mod              # 模块定义
├── go.work             # 工作区配置
└── README.md           # 项目说明
```

---

## 二、与标准对比分析

### ✅ 符合标准的部分

| 目录 | 状态 | 说明 |
|------|------|------|
| `cmd/` | ✅ | 主程序入口正确 |
| `api/` | ✅ | API 定义文件 |
| `configs/` | ✅ | 配置文件模板 |
| `deployments/` | ✅ | 部署配置 |
| `scripts/` | ✅ | 构建脚本 |
| `test/` | ✅ | 外部测试 |
| `tools/` | ✅ | 支持工具 |
| `benchmarks/` | ✅ | 基准测试 |

### ⚠️ 问题区域

#### 1. `docs/` - 文档过度膨胀

**问题**:

- 24 个子目录，超过 100 个 Markdown 文件
- 大量重复和归档文档
- 自我参照的 "完成报告" 堆积

**标准建议**:

```
docs/
├── README.md              # 文档索引
├── adr/                   # 架构决策记录
├── architecture/          # 架构文档
├── guides/                # 用户指南
├── development/           # 开发文档
├── api/                   # API 文档
└── deployment/            # 部署文档
```

#### 2. `internal/` - Clean Architecture 过度工程化

**问题**:

- 层级过多，包依赖复杂
- `interfaces/` 包含太多实现细节
- `infrastructure/` 过于臃肿

**当前结构**:

```
internal/
├── application/           # 用例层（空洞）
├── domain/                # 领域层（只有 User）
├── infrastructure/        # 基础设施层（过于庞大）
├── interfaces/            # 接口层（与基础设施混淆）
├── security/              # 安全（应在 pkg/）
└── ...
```

#### 3. `pkg/` - 公共库过度膨胀

**问题**:

- 40+ 顶级包，职责混乱
- `utils/` 包含 50+ 个小模块
- `auth/`、`rbac/`、`security/` 重复
- `observability/` 与 `tracing/` 重复

**重复/冲突**:

| 包 | 问题 |
|----|------|
| `pkg/auth/jwt` vs `pkg/security/jwt` | 重复实现 |
| `pkg/rbac` vs `pkg/security/rbac` | 重复实现 |
| `pkg/auth/oauth2` vs `pkg/security/oauth2` | 重复实现 |
| `pkg/tracing` vs `pkg/observability/tracing` | 重复实现 |

#### 4. `archive/` - 不应在版本控制中

**问题**:

- 归档文件应在版本控制外
- 使用 Git 历史或外部存储

---

## 三、推荐目录结构

### 目标架构

```
golang/
├── api/                          # API 定义（OpenAPI/Protobuf）
│   ├── openapi/                  # REST API 规范
│   ├── protobuf/                 # gRPC 协议定义
│   └── graphql/                  # GraphQL Schema
│
├── assets/                       # 静态资源（新增）
│   ├── images/
│   ├── templates/
│   └── migrations/               # 从 migrations/ 移动
│
├── build/                        # 构建相关（新增）
│   ├── ci/
│   ├── docker/
│   └── scripts/
│
├── cmd/                          # 主程序入口
│   ├── server/                   # HTTP 服务器
│   ├── worker/                   # 后台任务
│   └── cli/                      # 命令行工具
│
├── configs/                      # 配置文件模板
│   ├── config.yaml
│   └── config.prod.yaml
│
├── docs/                         # 精简后的文档
│   ├── README.md
│   ├── adr/                      # 架构决策记录
│   ├── architecture/             # 架构文档
│   ├── guides/                   # 用户指南
│   └── development/              # 开发文档
│
├── examples/                     # 示例代码（精简）
│   ├── basic/                    # 基础示例
│   ├── advanced/                 # 高级示例
│   └── complete/                 # 完整项目示例
│
├── internal/                     # 私有应用代码
│   ├── app/                      # 应用层（合并原 application）
│   │   ├── commands/             # CQRS 命令
│   │   ├── queries/              # CQRS 查询
│   │   └── services/             # 应用服务
│   │
│   ├── domain/                   # 领域层（精简）
│   │   ├── entities/             # 领域实体
│   │   ├── valueobjects/         # 值对象
│   │   ├── repositories/         # 仓储接口
│   │   └── services/             # 领域服务
│   │
│   ├── infra/                    # 基础设施层（简化）
│   │   ├── database/             # 数据库
│   │   ├── cache/                # 缓存
│   │   ├── messaging/            # 消息队列
│   │   └── observability/        # 可观测性
│   │
│   └── interfaces/               # 接口适配层
│       ├── http/                 # HTTP 处理器
│       ├── grpc/                 # gRPC 处理器
│       ├── graphql/              # GraphQL Resolvers
│       └── middleware/           # 中间件
│
├── pkg/                          # 公共库（精简为 10-15 个核心包）
│   ├── logger/                   # 日志
│   ├── errors/                   # 错误处理
│   ├── validator/                # 验证
│   ├── jwt/                      # JWT（合并 auth/security）
│   ├── rbac/                     # RBAC（合并）
│   ├── security/                 # 安全工具
│   ├── observability/            # 可观测性（合并 tracing）
│   ├── database/                 # 数据库工具
│   ├── cache/                    # 缓存工具
│   └── http/                     # HTTP 工具
│
├── tests/                        # 外部测试（重命名 test/）
│   ├── e2e/
│   ├── integration/
│   └── fixtures/
│
├── web/                          # 前端代码（新增，可选）
│   ├── static/
│   └── templates/
│
├── website/                      # 项目网站（新增，可选）
│
├── .github/
├── .goreleaser.yml
├── CHANGELOG.md
├── CODE_OF_CONDUCT.md
├── CONTRIBUTING.md
├── LICENSE
├── Makefile
├── README.md
├── SECURITY.md
├── go.mod
└── go.work
```

---

## 四、详细重组方案

### Phase 1: 清理归档文件（立即执行）

```bash
# 移除 archive/ 目录（已清理，但需从 git 历史移除）
git filter-branch --force --index-filter \
  'git rm --cached --ignore-unmatch -r archive/' \
  --prune-empty --tag-name-filter cat -- --all

# 清理 docs/archive/
git rm -rf docs/archive/
```

### Phase 2: 合并重复包

#### 2.1 安全包合并

```
当前:
pkg/auth/jwt
pkg/auth/oauth2
pkg/rbac
pkg/security/jwt      # 重复
pkg/security/rbac     # 重复
pkg/security/oauth2   # 重复

合并后:
pkg/security/
├── jwt/              # 合并 auth/jwt
├── oauth2/           # 合并 auth/oauth2
├── rbac/             # 合并 rbac/
├── abac/             # 保留
├── crypto/           # 加密工具
└── middleware/       # 安全中间件
```

#### 2.2 可观测性包合并

```
当前:
pkg/tracing
pkg/observability/
├── ebpf/
├── metrics/
├── operational/
├── otlp/
├── system/
└── tracing/          # 重复

合并后:
pkg/observability/
├── trace/            # 合并 tracing/
├── metric/
├── log/
├── ebpf/             # Linux only
└── otlp/             # OpenTelemetry 导出
```

#### 2.3 合并 utils/ 到标准库

```
当前: pkg/utils/ 下有 50+ 模块

建议: 只保留核心工具，其他使用标准库或第三方库

保留:
pkg/utils/
├── crypto/           # 加密
├── hash/             # 哈希
├── id/               # ID 生成
└── strings/          # 字符串

移除（使用标准库）:
- pkg/utils/time/     -> 使用 time 包
- pkg/utils/json/     -> 使用 encoding/json
- pkg/utils/regex/    -> 使用 regexp
- pkg/utils/context/  -> 使用 context
- ... 等 40+ 模块
```

### Phase 3: 重组 internal/

#### 3.1 应用层（internal/app/）

```
internal/app/
├── user/
│   ├── commands/           # 命令（创建、更新、删除）
│   ├── queries/            # 查询（获取、列表）
│   ├── dto/                # 数据传输对象
│   └── service.go          # 应用服务
│
└── workflow/
    ├── commands/
    ├── queries/
    └── service.go
```

#### 3.2 领域层（internal/domain/）

```
internal/domain/
├── user/
│   ├── entity.go           # 实体
│   ├── repository.go       # 仓储接口
│   ├── service.go          # 领域服务
│   └── specification/      # 规约模式
│
└── shared/
    ├── errors/             # 共享错误
    └── events/             # 共享事件
```

#### 3.3 基础设施层（internal/infra/）

```
internal/infra/
├── persistence/            # 持久化
│   ├── ent/                # Ent ORM
│   ├── redis/              # Redis
│   └── repository/         # 仓储实现
│
├── messaging/              # 消息队列
│   ├── kafka/
│   ├── nats/
│   └── mqtt/
│
├── observability/          # 可观测性
│   ├── otlp/
│   └── ebpf/               # Linux only
│
└── workflow/
    └── temporal/
```

#### 3.4 接口层（internal/interfaces/）

```
internal/interfaces/
├── http/
│   ├── handlers/           # HTTP 处理器
│   ├── middleware/         # HTTP 中间件
│   └── router.go           # 路由配置
│
├── grpc/
│   ├── handlers/
│   └── interceptors/
│
└── graphql/
    ├── resolvers/
    └── schema/
```

### Phase 4: 精简 pkg/

```
pkg/（只保留真正可复用的库）
├── logger/                 # 结构化日志
├── errors/                 # 错误处理
├── validator/              # 验证（基于 go-playground/validator）
├── jwt/                    # JWT 实现
├── rbac/                   # RBAC 实现
├── security/
│   ├── crypto/             # 加密工具
│   ├── hash/               # 哈希工具
│   └── password/           # 密码处理
├── observability/
│   ├── trace/              # 追踪
│   ├── metric/             # 指标
│   └── log/                # 日志集成
├── database/               # 数据库工具
├── cache/                  # 缓存工具
├── http/
│   ├── client/             # HTTP 客户端
│   └── server/             # HTTP 服务器工具
└── utils/
    ├── id/                 # ID 生成
    ├── strings/            # 字符串工具
    └── time/               # 时间工具
```

---

## 五、迁移指南

### 步骤 1: 创建新分支

```bash
git checkout -b refactor/directory-structure
```

### 步骤 2: 移动目录

```bash
# 1. 创建新目录结构
mkdir -p internal/{app,domain,infra,interfaces}
mkdir -p pkg/{security,observability}

# 2. 移动安全包
mv pkg/auth/jwt pkg/security/jwt
mv pkg/auth/oauth2 pkg/security/oauth2
mv pkg/rbac pkg/security/rbac
rm -rf pkg/auth

# 3. 移动可观测性包
mv pkg/tracing pkg/observability/trace
mv pkg/observability/otlp pkg/observability/otlp

# 4. 合并 utils/
mkdir -p pkg/utils/keep
mv pkg/utils/{id,strings,time} pkg/utils/keep/
rm -rf pkg/utils
mv pkg/utils/keep pkg/utils
```

### 步骤 3: 更新导入路径

```bash
# 使用 gofmt 重写导入
gofmt -w -r '"github.com/yourusername/golang/pkg/auth/jwt" -> "github.com/yourusername/golang/pkg/security/jwt"' .

# 或使用 sed
find . -name "*.go" -type f -exec sed -i \
  's|github.com/yourusername/golang/pkg/auth|github.com/yourusername/golang/pkg/security|g' {} +
```

### 步骤 4: 验证构建

```bash
go build ./...
go test -short ./...
```

---

## 六、预期收益

### 1. 代码清晰度提升

| 指标 | 当前 | 目标 | 改善 |
|------|------|------|------|
| 顶级包数 | 40+ | 15 | -60% |
| utils 模块数 | 50+ | 5 | -90% |
| 重复实现 | 6 处 | 0 | -100% |
| 平均导入深度 | 4 层 | 3 层 | -25% |

### 2. 维护性提升

- **单一职责**: 每个包职责明确
- **依赖清晰**: 减少循环依赖
- **测试简化**: 减少 mock 需求

### 3. 新开发者上手

- 符合标准布局，降低学习成本
- 目录结构直观，快速定位代码

---

## 七、风险评估

### 风险 1: 破坏性变更

**影响**: 这是一个重大重构，会影响所有导入路径

**缓解**:

- 在 README 中明确标记为破坏性变更
- 提供迁移脚本
- 更新版本号（v1.0 -> v2.0）

### 风险 2: 测试覆盖

**影响**: 移动文件可能导致测试失效

**缓解**:

- 使用 IDE 重构工具
- 全面运行测试套件
- 增加集成测试

### 风险 3: 文档失效

**影响**: 大量文档引用旧路径

**缓解**:

- 同步更新所有文档
- 添加重定向说明

---

## 八、实施时间表

| Phase | 任务 | 时间 | 负责人 |
|-------|------|------|--------|
| 1 | 清理 archive/ 和 docs/archive/ | 1 天 | DevOps |
| 2 | 合并重复包 | 3 天 | Backend |
| 3 | 重组 internal/ | 5 天 | Backend |
| 4 | 精简 pkg/ | 3 天 | Backend |
| 5 | 更新导入路径 | 2 天 | Backend |
| 6 | 更新文档 | 2 天 | Tech Writer |
| 7 | 全面测试 | 3 天 | QA |
| **总计** | | **~3 周** | |

---

## 九、参考资源

1. [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
2. [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
3. [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
4. [Effective Go](https://golang.org/doc/effective_go.html)

---

**建议**: 这是一个重大重构，建议在以下时机执行：

- 主要功能稳定后
- 新功能开发暂停期间
- 有足够的测试覆盖时

**下一步**: 如需执行，请创建详细的任务分解并分配给团队成员。

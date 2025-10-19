# 📁 项目新结构规划

> **重构日期**: 2025年10月19日  
> **目标**: 文档与代码分离，清晰的模块化结构

---

## 🎯 新目录结构

```text
golang/
├── 📂 docs/                          # 📚 文档目录（纯Markdown）
│   ├── README.md                     # 文档总览
│   ├── INDEX.md                      # 学习索引（已存在）
│   │
│   ├── 01-getting-started/           # 🌱 入门指南
│   │   ├── 01-installation.md
│   │   ├── 02-quick-start.md
│   │   └── 03-basic-concepts.md
│   │
│   ├── 02-language-fundamentals/    # 📖 语言基础
│   │   ├── 01-syntax/
│   │   ├── 02-types/
│   │   ├── 03-functions/
│   │   └── 04-error-handling/
│   │
│   ├── 03-concurrency/              # ⚡ 并发编程
│   │   ├── 01-goroutines.md
│   │   ├── 02-channels.md
│   │   ├── 03-patterns/
│   │   │   ├── pipeline.md
│   │   │   ├── worker-pool.md
│   │   │   └── fan-out-fan-in.md
│   │   └── 04-best-practices.md
│   │
│   ├── 04-standard-library/         # 📦 标准库
│   │   ├── 01-io/
│   │   ├── 02-net/
│   │   ├── 03-http/
│   │   └── 04-testing/
│   │
│   ├── 05-go-features/              # 🆕 Go版本特性
│   │   ├── go-1.21.md
│   │   ├── go-1.22.md
│   │   ├── go-1.23.md
│   │   └── version-comparison.md
│   │
│   ├── 06-generics/                 # 🔷 泛型编程
│   │   ├── 01-basics.md
│   │   ├── 02-constraints.md
│   │   └── 03-advanced-usage.md
│   │
│   ├── 07-performance/              # 🚀 性能优化
│   │   ├── 01-profiling.md
│   │   ├── 02-memory-optimization.md
│   │   ├── 03-cpu-optimization.md
│   │   └── 04-pgo.md
│   │
│   ├── 08-testing/                  # 🧪 测试
│   │   ├── 01-unit-testing.md
│   │   ├── 02-benchmarking.md
│   │   ├── 03-coverage.md
│   │   └── 04-integration-testing.md
│   │
│   ├── 09-web-development/          # 🌐 Web开发
│   │   ├── 01-http-server.md
│   │   ├── 02-routing.md
│   │   ├── 03-middleware.md
│   │   └── 04-frameworks.md
│   │
│   ├── 10-microservices/            # 🔬 微服务
│   │   ├── 01-architecture.md
│   │   ├── 02-grpc.md
│   │   ├── 03-service-discovery.md
│   │   └── 04-api-gateway.md
│   │
│   ├── 11-cloud-native/             # ☁️ 云原生
│   │   ├── 01-kubernetes.md
│   │   ├── 02-docker.md
│   │   └── 03-monitoring.md
│   │
│   ├── 12-design-patterns/          # 🎨 设计模式
│   │   ├── 01-creational.md
│   │   ├── 02-structural.md
│   │   └── 03-behavioral.md
│   │
│   └── 13-advanced-topics/          # 🎓 高级主题
│       ├── 01-ai-agent-architecture/
│       ├── 02-distributed-systems/
│       └── 03-system-design/
│
├── 📂 examples/                      # 💻 示例代码（纯Go代码）
│   ├── README.md                     # 示例总览
│   │
│   ├── 01-basics/                   # 基础示例
│   │   ├── hello-world/
│   │   ├── variables/
│   │   └── functions/
│   │
│   ├── 02-concurrency/              # 并发示例
│   │   ├── goroutines/
│   │   ├── channels/
│   │   ├── pipeline/
│   │   │   ├── main.go
│   │   │   ├── pipeline_test.go
│   │   │   └── README.md
│   │   └── worker-pool/
│   │       ├── main.go
│   │       ├── worker_pool_test.go
│   │       └── README.md
│   │
│   ├── 03-generics/                 # 泛型示例
│   │   ├── basic/
│   │   ├── constraints/
│   │   └── collections/
│   │
│   ├── 04-web/                      # Web示例
│   │   ├── http-server/
│   │   ├── rest-api/
│   │   └── websocket/
│   │
│   ├── 05-testing/                  # 测试示例
│   │   ├── unit-test/
│   │   ├── benchmark/
│   │   └── integration/
│   │
│   ├── 06-performance/              # 性能优化示例
│   │   ├── pgo/
│   │   ├── memory/
│   │   └── cpu/
│   │
│   ├── 07-standard-lib/             # 标准库示例
│   │   ├── slices/
│   │   ├── maps/
│   │   └── iter/
│   │
│   └── 08-advanced/                 # 高级示例
│       ├── ai-agent/
│       │   ├── core/
│       │   │   ├── agent.go
│       │   │   ├── decision_engine.go
│       │   │   ├── learning_engine.go
│       │   │   └── types.go
│       │   ├── agent_test.go
│       │   ├── go.mod
│       │   └── README.md
│       └── distributed/
│
├── 📂 pkg/                          # 🔧 可复用包（库代码）
│   ├── concurrent/                  # 并发工具
│   ├── collections/                 # 集合工具
│   └── utils/                       # 通用工具
│
├── 📂 cmd/                          # 🎯 命令行工具
│   ├── project-stats/
│   └── changelog-gen/
│
├── 📂 scripts/                      # 📜 脚本工具
│   ├── test_summary.ps1
│   ├── format_code.ps1
│   └── replace_version.ps1
│
├── 📂 reports/                      # 📊 项目报告
│   ├── README.md                    # 报告索引（已存在）
│   ├── phase-1/
│   ├── phase-2/
│   └── phase-3/
│
├── 📂 .github/                      # ⚙️ GitHub配置
│   ├── workflows/
│   ├── ISSUE_TEMPLATE/
│   └── PULL_REQUEST_TEMPLATE.md
│
├── README.md                        # 项目首页
├── README_EN.md                     # 英文首页
├── CONTRIBUTING.md                  # 贡献指南
├── CONTRIBUTING_EN.md               # 英文贡献指南
├── EXAMPLES.md                      # 示例展示
├── EXAMPLES_EN.md                   # 英文示例展示
├── QUICK_START.md                   # 快速开始
├── QUICK_START_EN.md                # 英文快速开始
├── CHANGELOG.md                     # 变更日志
├── LICENSE                          # 许可证
└── go.work                          # Go workspace（可选）

```

---

## 🔄 迁移映射

### 文档迁移

**从 → 到**:

```text
docs/02-Go语言现代化/12-Go-1.23运行时优化/
→ docs/07-performance/

docs/02-Go语言现代化/13-Go-1.23工具链增强/
→ docs/05-go-features/tools/

docs/02-Go语言现代化/14-Go-1.23并发和网络/
→ docs/03-concurrency/

docs/02-Go语言现代化/08-智能化架构集成/01-AI-Agent架构/
→ docs/13-advanced-topics/01-ai-agent-architecture/

docs/01-HTTP服务/
→ docs/09-web-development/

docs/03-并发编程/
→ docs/03-concurrency/

docs/04-设计模式/
→ docs/12-design-patterns/

docs/05-性能优化/
→ docs/07-performance/

docs/06-微服务架构/
→ docs/10-microservices/
```

### 代码迁移

**从 → 到**:

```text
docs/02-Go语言现代化/08-智能化架构集成/01-AI-Agent架构/core/
→ examples/08-advanced/ai-agent/core/

docs/02-Go语言现代化/14-Go-1.23并发和网络/examples/waitgroup_go/
→ 删除（虚构特性）

examples/concurrency/pipeline_test.go
→ examples/02-concurrency/pipeline/

examples/concurrency/worker_pool_test.go
→ examples/02-concurrency/worker-pool/

examples/advanced/
→ examples/08-advanced/
```

---

## 📝 迁移步骤

### Phase 1: 创建新目录结构

```bash
# 文档目录
mkdir -p docs/{01-getting-started,02-language-fundamentals,03-concurrency,04-standard-library}
mkdir -p docs/{05-go-features,06-generics,07-performance,08-testing}
mkdir -p docs/{09-web-development,10-microservices,11-cloud-native,12-design-patterns}
mkdir -p docs/13-advanced-topics/01-ai-agent-architecture

# 示例目录
mkdir -p examples/{01-basics,02-concurrency,03-generics,04-web}
mkdir -p examples/{05-testing,06-performance,07-standard-lib,08-advanced}

# 库目录
mkdir -p pkg/{concurrent,collections,utils}

# 命令目录
mkdir -p cmd/{project-stats,changelog-gen}
```

### Phase 2: 迁移文件

```bash
# 迁移AI-Agent代码
mv docs/02-Go语言现代化/08-智能化架构集成/01-AI-Agent架构/core \
   examples/08-advanced/ai-agent/

# 迁移并发示例
mv examples/concurrency/pipeline_test.go \
   examples/02-concurrency/pipeline/

mv examples/concurrency/worker_pool_test.go \
   examples/02-concurrency/worker-pool/
```

### Phase 3: 删除虚构特性

```bash
# 删除WaitGroup.Go()示例
rm -rf docs/02-Go语言现代化/14-Go-1.23并发和网络/examples/waitgroup_go/

# 删除不存在的特性文档
rm docs/02-Go语言现代化/14-Go-1.23并发和网络/01-WaitGroup-Go方法.md
rm docs/02-Go语言现代化/14-Go-1.23并发和网络/02-testing-synctest包.md
rm docs/02-Go语言现代化/13-Go-1.23工具链增强/02-go-mod-ignore指令.md
```

### Phase 4: 重命名和更新

```bash
# 重命名文档目录
mv docs/02-Go语言现代化/12-Go-1.23运行时优化 docs/07-performance/runtime
mv docs/03-并发编程 docs/03-concurrency/fundamentals
mv docs/04-设计模式 docs/12-design-patterns
```

---

## ✅ 验证清单

### 目录结构

- [ ] docs/ 纯文档目录
- [ ] examples/ 纯代码示例目录
- [ ] pkg/ 可复用库目录
- [ ] cmd/ 命令行工具目录
- [ ] scripts/ 脚本目录

### 文档组织

- [ ] 按主题分类，不按版本
- [ ] 每个主题独立目录
- [ ] README.md 作为索引

### 代码组织

- [ ] 每个示例独立目录
- [ ] 包含 go.mod
- [ ] 包含 README.md 说明
- [ ] 包含测试文件

### 命名规范

- [ ] 目录名使用 kebab-case
- [ ] Go文件使用 snake_case
- [ ] 包名使用单个小写单词

---

## 🎯 最终结构优势

### 清晰分离

**文档 (docs/)**:

- ✅ 纯Markdown
- ✅ 按主题组织
- ✅ 便于阅读和维护

**代码 (examples/)**:

- ✅ 纯Go代码
- ✅ 可独立运行
- ✅ 便于测试和复用

**库 (pkg/)**:

- ✅ 可复用代码
- ✅ 清晰的API
- ✅ 完整的测试

### 易于导航

**学习者**:

```text
docs/INDEX.md → 按难度学习
examples/ → 运行示例
```

**贡献者**:

```text
CONTRIBUTING.md → 了解规范
examples/ → 提交示例
pkg/ → 贡献库代码
```

**使用者**:

```text
README.md → 快速开始
EXAMPLES.md → 查看示例
pkg/ → 导入使用
```

---

<div align="center">

## 🚀 准备开始迁移

**新结构已规划完成，准备执行迁移！**

</div>

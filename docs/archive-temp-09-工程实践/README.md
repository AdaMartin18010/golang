
# Go测试最佳实践

## 📚 模块概述

本模块系统介绍Go语言测试的最佳实践，包括单元测试、集成测试、基准测试、Mock测试等。通过理论分析与实际代码相结合的方式，帮助开发者建立完整的测试体系。

## 🎯 学习目标

- 掌握Go语言测试框架和工具
- 学会编写高质量的单元测试
- 理解集成测试和端到端测试
- 掌握Mock和依赖注入技术
- 建立完整的测试覆盖率体系

## 📋 内容结构

### 测试基础
- [01-测试基础与go-test](./01-测试基础与go-test.md) - go test基本用法
- [02-单元测试与表驱动](./02-单元测试与表驱动.md) - 表驱动测试法
- [03-基准测试与性能测试](./03-基准测试与性能测试.md) - 性能基准测试

### 高级测试技术
- [04-Mock与依赖注入](./04-Mock与依赖注入.md) - Mock测试和依赖注入
- [05-测试覆盖率与分析](./05-测试覆盖率与分析.md) - 覆盖率分析和工具
- [06-CI集成与自动化测试](./06-CI集成与自动化测试.md) - CI/CD集成

### 测试最佳实践
- [07-测试最佳实践](./07-测试最佳实践.md) - 测试策略和最佳实践

## 🚀 快速开始

### 基本测试示例

```go
// math.go
package math

func Add(a, b int) int {
    return a + b
}

// math_test.go
package math

import "testing"

func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive", 2, 3, 5},
        {"negative", -1, -1, -2},
        {"zero", 0, 5, 5},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := Add(tt.a, tt.b); got != tt.want {
                t.Errorf("Add(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

### 运行测试

```bash
# 运行所有测试
go test

# 运行指定测试
go test -run TestAdd

# 显示详细输出
go test -v

# 运行基准测试
go test -bench=.

# 生成覆盖率报告
go test -cover
```

## 📊 学习进度

| 主题 | 状态 | 完成度 | 预计时间 |
|------|------|--------|----------|
| 测试基础 | 🔄 进行中 | 0% | 1-2天 |
| 单元测试 | ⏳ 待开始 | 0% | 2-3天 |
| 基准测试 | ⏳ 待开始 | 0% | 1-2天 |
| Mock测试 | ⏳ 待开始 | 0% | 2-3天 |
| CI集成 | ⏳ 待开始 | 0% | 1-2天 |

## 🎯 实践项目

### 项目1: 计算器测试
- 实现基本计算功能
- 编写单元测试
- 使用表驱动测试
- 达到100%覆盖率

### 项目2: HTTP服务测试
- 创建HTTP服务
- 编写集成测试
- 使用Mock测试外部依赖
- 配置自动化测试

### 项目3: 数据库操作测试
- 实现CRUD操作
- 编写数据库测试
- 使用测试数据库
- 实现事务测试

## 📚 参考资料

### 官方文档
- [Go测试文档](https://golang.org/pkg/testing/)
- [Go测试命令](https://golang.org/cmd/go/#hdr-Test_packages)
- [Go覆盖率工具](https://golang.org/cmd/go/#hdr-Test_coverage)

### 在线教程
- [Go测试实战](https://geektutu.com/post/hpg-golang-unit-test.html)
- [Go测试最佳实践](https://github.com/golang/go/wiki/TestComments)
- [Go Mock测试](https://github.com/golang/mock)

### 书籍推荐
- 《Go语言实战》第9章
- 《Go并发编程实战》第8章
- 《测试驱动开发》

## 🔧 工具推荐

### 测试框架
- **go test**: Go官方测试框架
- **testify**: 第三方测试工具包
- **ginkgo**: BDD测试框架
- **gomega**: 断言库

### Mock工具
- **gomock**: Go官方Mock工具
- **testify/mock**: 简单Mock实现
- **counterfeiter**: 接口Mock生成器

### 覆盖率工具
- **go test -cover**: 内置覆盖率
- **gocov**: 覆盖率可视化
- **gocover**: 覆盖率报告

### CI/CD工具
- **GitHub Actions**: GitHub CI/CD
- **GitLab CI**: GitLab CI/CD
- **Jenkins**: 开源CI/CD平台

## 🎯 学习建议

### 测试驱动开发
- 先写测试，后写实现
- 红-绿-重构循环
- 保持测试简单和专注

### 测试金字塔
- **单元测试**: 70% - 快速、独立、可重复
- **集成测试**: 20% - 测试组件间交互
- **端到端测试**: 10% - 测试完整用户流程

### 测试质量
- 测试应该快速执行
- 测试应该独立运行
- 测试应该可重复
- 测试应该及时反馈

## 📝 重要概念

### 测试类型
- **单元测试**: 测试单个函数或方法
- **集成测试**: 测试多个组件协作
- **端到端测试**: 测试完整用户流程
- **基准测试**: 测试性能指标

### 测试策略
- **表驱动测试**: 使用表格数据驱动测试
- **子测试**: 使用t.Run创建子测试
- **并行测试**: 使用t.Parallel并行执行
- **跳过测试**: 使用t.Skip跳过特定测试

### Mock技术
- **接口Mock**: 模拟接口实现
- **函数Mock**: 模拟函数调用
- **依赖注入**: 通过构造函数注入依赖
- **测试替身**: 使用测试替身替代真实依赖

## 🛠️ 最佳实践

### 测试命名
- 测试函数名以Test开头
- 使用描述性的测试名称
- 包含测试场景和期望结果

### 测试结构
- 使用AAA模式（Arrange-Act-Assert）
- 保持测试函数简洁
- 一个测试函数测试一个场景

### 测试数据
- 使用有意义的测试数据
- 避免硬编码的测试数据
- 使用测试数据构建器

### 错误处理
- 测试正常流程和异常流程
- 验证错误消息和错误类型
- 测试边界条件和极端情况

### 性能测试
- 使用基准测试测量性能
- 设置性能基准线
- 监控性能回归

---

**模块维护者**: AI Assistant  
**最后更新**: 2025年1月  
**模块状态**: 持续更新中

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2025-11-11

### Added

- Clean Architecture 四层架构实现
- Domain Layer: User 实体、Repository 接口、Domain Errors
- Application Layer: User Service、DTOs
- Infrastructure Layer: PostgreSQL Repository、Ent Schema、OpenTelemetry、Kafka、MQTT
- Interfaces Layer: Chi HTTP 路由、gRPC Proto、GraphQL Schema
- 技术栈集成: Chi, Ent, Viper, Slog, OpenTelemetry, Wire, Kafka, MQTT, gRPC, GraphQL
- 测试框架: 单元测试、Mock 支持
- 部署配置: Docker、Kubernetes、Docker Compose
- 文档: 架构文档、开发指南、API 文档、测试指南

### Changed

- 项目结构重构为现代化 Clean Architecture
- 采用最新最成熟的 Go 技术栈

### Fixed

- 所有编译错误
- 所有测试通过

## [0.1.0] - 2025-11-11

### Added1

- 初始项目结构
- 基础配置

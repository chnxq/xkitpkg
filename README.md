# XKit Package (xkitpkg)

XKit Package (xkitpkg) 是一个轻量级、可扩展的 Go 微服务工具包，集成了配置管理、日志记录、服务注册发现、链路追踪等功能模块，旨在为 Go 微服务开发提供标准化的基础设施组件。

## 项目概述

xkitpkg 是一个综合性的微服务工具包，提供了构建现代 Go 微服务所需的核心功能：

- **配置管理**：支持多种配置源（本地文件、Etcd、Consul、Nacos 等）的动态配置系统
- **日志记录**：统一的日志接口，支持多种日志后端（Zap、Fluentd 等）
- **服务注册与发现**：支持多种服务注册中心（Etcd、Consul、ZooKeeper、Nacos 等）
- **分布式链路追踪**：基于 OpenTelemetry 的分布式链路追踪系统
- **协议定义**：基于 Protocol Buffers 的配置和服务定义

## 目录结构

```
xkitpkg/
├── conf/           # Protocol Buffers 配置定义及生成代码
├── config/         # 配置管理模块
├── logger/         # 日志记录模块
├── registry/       # 服务注册与发现模块
├── tracer/         # 分布式链路追踪模块
├── go.mod          # Go 模块定义
└── go.sum          # Go 依赖校验和
```

## 快速开始

### 安装

```bash
go get github.com/chnxq/xkitpkg
```

### 基本使用

```go
package main

import (
    "context"
    "log"
    
    "github.com/chnxq/xkitpkg/config"
    "github.com/chnxq/xkitpkg/logger"
    "github.com/chnxq/xkitpkg/registry"
    "github.com/chnxq/xkitpkg/tracer"
    "github.com/chnxq/xkitpkg/conf/v1"
)

func main() {
    // 加载配置
    if err := config.LoadAdminServerConfig("./config.yaml"); err != nil {
        log.Fatal("Failed to load config:", err)
    }
    
    // 获取配置
    adminConfig := config.GetServerConfig()
    
    // 初始化日志系统
    lg, err := logger.NewLogger(adminConfig.GetLog())
    if err != nil {
        log.Fatal("Failed to create logger:", err)
    }
    
    // 初始化服务注册
    reg, err := registry.NewRegistrar(adminConfig.GetRegistry())
    if err != nil {
        log.Fatal("Failed to create registrar:", err)
    }
    
    // 初始化链路追踪
    ctx := context.Background()
    appInfo := &conf.AppInfo{
        Name:    "my-app",
        Version: "v1.0.0",
    }
    if err := tracer.NewTracerProvider(ctx, adminConfig.GetTrace(), appInfo); err != nil {
        log.Fatal("Failed to create tracer provider:", err)
    }
    
    // 应用逻辑...
    
    // 关闭追踪系统
    if err := tracer.ShutdownTracerProvider(ctx); err != nil {
        log.Println("Failed to shutdown tracer:", err)
    }
}
```

## 特性

### 配置管理 (config)
- 支持多种配置源（本地文件、Etcd、Consul、Nacos 等）
- 动态配置更新，无需重启应用
- 统一的配置访问接口
- 配置热重载机制

### 日志记录 (logger)
- 统一日志接口，支持多种后端
- 结构化日志输出
- 支持 Zap、Fluentd 等日志库
- 链路追踪上下文集成

### 服务注册与发现 (registry)
- 支持多种服务注册中心
- 服务健康检查
- 动态服务发现
- 负载均衡支持

### 分布式链路追踪 (tracer)
- 基于 OpenTelemetry 标准
- 多种导出器支持（OTLP、Zipkin 等）
- 自动资源属性配置
- 优雅关闭机制

## 贡献

欢迎提交 Issue 和 Pull Request 来改进 xkitpkg！

## 许可证

MIT License
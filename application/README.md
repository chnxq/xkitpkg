# 应用程序引导包

## 概述

应用程序引导包（`application`）负责程序的引导配置管理，提供线程安全的初始化流程和配置注册机制，用于在应用启动阶段集中管理各类配置结构体（例如服务器、客户端、数据、日志等）。

## 功能特性

- **延迟初始化**：使用 `sync.Once` 确保引导配置仅初始化一次
- **并发安全**：读写操作通过 `sync.RWMutex` 保护
- **配置注册**：通过 `RegisterConfig` 注册任意非空指针类型配置
- **主配置访问**：使用 `GetBootstrapConfig` 获取共享的 `*conf.Bootstrap` 实例
- **应用上下文管理**：提供 `AppCtx` 结构体管理应用级上下文
- **优雅关闭**：支持应用的优雅启动和关闭

## 核心组件

### AppCtx

`AppCtx` 是应用级上下文结构体，提供以下功能：

- 管理应用配置和应用信息
- 提供日志记录器和服务注册器
- 支持自定义配置和值存储
- 管理应用级根上下文，用于优雅关闭

### 配置管理

- **配置注册**：通过 `RegisterConfig` 注册配置结构体
- **配置获取**：通过 `GetBootstrapConfig` 获取主配置
- **配置持久化**：支持配置的序列化和反序列化

## 快速开始

### 基本使用

```go
package main

import (
    "context"
    "log"
    
    "github.com/chnxq/xkitpkg/application"
    "github.com/chnxq/xkitpkg/conf/v1"
)

func main() {
    // 创建应用上下文
    appCtx := application.NewAppCtx(nil, &conf.AppInfo{
        AppName:    "example",
        AppVersion: "1.0.0",
        Environment: "production",
    })
    
    // 注册配置
    serverConfig := &conf.ServerConfig{}
    if err := application.RegisterConfig(serverConfig); err != nil {
        log.Fatalf("注册配置失败: %v", err)
    }
    
    // 获取引导配置
    bootstrapConfig, err := application.GetBootstrapConfig()
    if err != nil {
        log.Fatalf("获取引导配置失败: %v", err)
    }
    
    // 使用配置启动服务
    // ...
    
    // 优雅关闭
    defer appCtx.Cancel()
}
```

### 完整示例

```go
package main

import (
    "log"
    
    "github.com/chnxq/xkitpkg/application"
    "github.com/chnxq/xkitpkg/conf/v1"
    
    // 注册配置源
    _ "github.com/chnxq/xkitpkg/config/etcd"
    
    // 注册日志实现
    // _ "github.com/chnxq/xkitpkg/logger/zap"
    
    // 注册服务注册实现
    _ "github.com/chnxq/xkitpkg/registry/etcd"
)

func initApp() error {
    // 注册配置
    serverConfig := &conf.ServerConfig{}
    if err := application.RegisterConfig(serverConfig); err != nil {
        return err
    }
    
    clientConfig := &conf.ClientConfig{}
    if err := application.RegisterConfig(clientConfig); err != nil {
        return err
    }
    
    // 获取引导配置
    bootstrapConfig, err := application.GetBootstrapConfig()
    if err != nil {
        return err
    }
    
    // 使用配置启动服务
    log.Printf("应用启动，服务配置: %+v", bootstrapConfig.Server)
    
    return nil
}

func main() {
    // 创建应用上下文
    appCtx := application.NewAppCtx(nil, &conf.AppInfo{
        AppName:    "example-service",
        AppVersion: "1.0.0",
        Environment: "production",
    })
    defer appCtx.Cancel()
    
    // 初始化应用
    if err := initApp(); err != nil {
        log.Fatalf("应用初始化失败: %v", err)
    }
    
    log.Println("应用启动成功")
    
    // 等待信号
    // ...
}
```

## API 参考

### 配置管理

#### `func RegisterConfig(cfg interface{}) error`

注册配置结构体到引导系统。

- **参数**：`cfg` - 配置结构体指针
- **返回值**：错误信息

#### `func GetBootstrapConfig() (*conf.Bootstrap, error)`

获取引导配置实例。

- **返回值**：引导配置实例和错误信息

### 应用上下文

#### `func NewAppCtx(parent context.Context, ai *conf.AppInfo) *AppCtx`

创建应用级上下文。

- **参数**：
  - `parent` - 父上下文（传 nil 使用 Background）
  - `ai` - 应用信息
- **返回值**：应用上下文实例

#### `func (c *AppCtx) Config() *conf.ServerConfig`

获取服务器配置。

- **返回值**：服务器配置实例

#### `func (c *AppCtx) AppInfo() *conf.AppInfo`

获取应用信息。

- **返回值**：应用信息实例

#### `func (c *AppCtx) Logger() kLog.Logger`

获取日志记录器。

- **返回值**：日志记录器实例

#### `func (c *AppCtx) Registrar() kRegistry.Registrar`

获取服务注册器。

- **返回值**：服务注册器实例

#### `func (c *AppCtx) Context() context.Context`

获取应用级根上下文。

- **返回值**：应用级根上下文

#### `func (c *AppCtx) Cancel()`

取消应用级根上下文，用于优雅关闭。

## 最佳实践

### 配置管理

- **集中注册**：在应用启动时集中注册所有配置结构体
- **配置验证**：注册前对配置进行验证
- **配置优先级**：明确配置的加载优先级（默认值 < 配置文件 < 环境变量 < 命令行参数）

### 应用上下文

- **单例使用**：每个应用只创建一个 `AppCtx` 实例
- **资源管理**：通过 `AppCtx` 统一管理应用资源
- **优雅关闭**：使用 `AppCtx.Cancel()` 实现优雅关闭

### 错误处理

- **初始化错误**：对初始化过程中的错误进行妥善处理
- **配置错误**：对配置解析和验证错误进行详细记录
- **运行时错误**：通过应用上下文管理运行时错误

## 依赖关系

- **配置包**：`github.com/chnxq/xkitpkg/config`
- **日志包**：`github.com/chnxq/XGoKit/log`
- **注册包**：`github.com/chnxq/XGoKit/registry`
- **Protobuf**：`google.golang.org/protobuf/proto`

## 配置源

支持以下配置源（通过导入相应包启用）：

- **etcd**：`github.com/chnxq/xkitpkg/config/etcd`
- **consul**：`github.com/chnxq/xkitpkg/config/consul`
- **nacos**：`github.com/chnxq/xkitpkg/config/nacos`
- **apollo**：`github.com/chnxq/xkitpkg/config/apollo`
- **kubernetes**：`github.com/chnxq/xkitpkg/config/kubernetes`
- **polaris**：`github.com/chnxq/xkitpkg/config/polaris`

## 日志实现

支持以下日志实现（通过导入相应包启用）：

- **zap**：`github.com/chnxq/xkitpkg/logger/zap`
- **logrus**：`github.com/chnxq/xkitpkg/logger/logrus`
- **zerolog**：`github.com/chnxq/xkitpkg/logger/zerolog`
- **fluent**：`github.com/chnxq/xkitpkg/logger/fluent`
- **aliyun**：`github.com/chnxq/xkitpkg/logger/aliyun`
- **tencent**：`github.com/chnxq/xkitpkg/logger/tencent`

## 服务注册实现

支持以下服务注册实现（通过导入相应包启用）：

- **etcd**：`github.com/chnxq/xkitpkg/registry/etcd`
- **consul**：`github.com/chnxq/xkitpkg/registry/consul`
- **nacos**：`github.com/chnxq/xkitpkg/registry/nacos`
- **eureka**：`github.com/chnxq/xkitpkg/registry/eureka`
- **kubernetes**：`github.com/chnxq/xkitpkg/registry/kubernetes`
- **polaris**：`github.com/chnxq/xkitpkg/registry/polaris`
- **servicecomb**：`github.com/chnxq/xkitpkg/registry/servicecomb`
- **zookeeper**：`github.com/chnxq/xkitpkg/registry/zookeeper`


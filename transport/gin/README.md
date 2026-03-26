# Gin Transport for XGoKit

Gin Transport是XGoKit框架中的一个HTTP传输层实现，基于Gin Web框架构建。它提供了高性能、可扩展的HTTP服务功能，并集成了日志记录、错误处理、中间件支持等特性。

## 功能特点

- 基于Gin Web框架的高性能HTTP服务器
- 支持中间件链式处理
- 内置日志记录和异常恢复机制
- 支持请求解码和响应编码自定义
- 集成OpenTelemetry追踪
- 支持TLS配置
- 统一的错误处理机制

## 安装

```go
import "github.com/chnxq/xkitpkg/transport/gin"
```

## 快速开始

以下是一个简单的示例，展示如何使用Gin Transport创建一个HTTP服务：

```go
package main

import (
    "context"
    "github.com/chnxq/xkitpkg/transport/gin"
)

func main() {
    ctx := context.Background()

    // 创建新的Gin服务器实例
    srv := gin.NewServer(
        gin.WithAddress(":8080"),  // 设置监听地址
        gin.WithTimeout(1*time.Second),  // 设置超时时间
    )

    // 添加Gin内置中间件
    srv.Use(gin.Recovery())
    srv.Use(gin.Logger())

    // 注册路由处理器
    srv.GET("/hello", func(c *gin.Context) {
        c.JSON(200, map[string]string{"message": "Hello World!"})
    })

    // 启动服务器
    if err := srv.Start(ctx); err != nil {
        panic(err)
    }

    // 优雅关闭服务器
    defer func() {
        if err := srv.Stop(ctx); err != nil {
            // 处理停止服务器时的错误
        }
    }()
}
```

## 配置选项

Gin Transport支持多种配置选项：

- `WithAddress(addr string)` - 设置服务器监听地址
- `WithTLSConfig(config *tls.Config)` - 设置TLS配置
- `WithTimeout(timeout time.Duration)` - 设置请求超时时间
- `WithMiddleware(m ...middleware.Middleware)` - 添加中间件
- `WithFilter(filters ...kHttp.FilterFunc)` - 添加过滤器
- `WithRequestDecoder(dec kHttp.DecodeRequestFunc)` - 设置请求解码器
- `WithResponseEncoder(enc kHttp.EncodeResponseFunc)` - 设置响应编码器
- `WithErrorHandler(h kHttp.ErrorHandler)` - 设置错误处理器

## 中间件

该传输层提供了一些预设的中间件：

- `GinLogger(logger log.Logger)` - 日志记录中间件
- `GinRecovery(logger log.Logger, stack bool)` - 异常恢复中间件，防止程序崩溃

## 传输上下文

在请求处理过程中，可以通过上下文获取传输相关信息：

- `Kind()` - 获取传输类型（固定为 "gin"）
- `Endpoint()` - 获取端点信息
- `Operation()` - 获取操作信息
- `Request()` - 获取原始HTTP请求对象
- `PathTemplate()` - 获取路径模板

## 错误处理

系统集成了XGoKit的错误处理机制，支持统一的错误格式化和日志记录。

## 日志系统

通过内置的日志函数可以方便地记录不同级别的日志：

- `LogDebug`, `LogInfo`, `LogWarn`, `LogError`, `LogFatal`
- `LogDebugf`, `LogInfof`, `LogWarnf`, `LogErrorf`, `LogFatalf`

## 测试

运行单元测试：

```bash
go test -v
```

## 贡献

欢迎提交Issue和Pull Request来改进此项目。
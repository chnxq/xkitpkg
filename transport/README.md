# Transport Layer for XGoKit

Transport 是 XGoKit 框架中的传输层实现，提供多种传输协议支持，包括 HTTP、gRPC、Gin、SSE、Asynq 等。传输层作为服务间通信的基础组件，提供了统一的接口抽象和丰富的功能特性。

## 功能特点

- **多协议支持**: 支持 HTTP、gRPC、SSE、Asynq 等多种传输协议
- **统一接口**: 提供统一的 `Server` 接口，便于不同传输协议间的切换
- **上下文传递**: 支持传输相关的上下文信息传递
- **中间件支持**: 支持各种中间件扩展功能
- **请求/响应处理**: 统一的请求解码和响应编码机制
- **错误处理**: 统一的错误处理机制
- **监控集成**: 集成 OpenTelemetry 追踪

## 支持的传输协议

### HTTP Transport
基于标准库 net/http 的 HTTP 传输实现，提供：
- 标准 HTTP 服务支持
- 请求/响应头管理
- 路径模板支持
- Cookie 操作

### gRPC Transport
基于 gRPC 协议的传输实现，提供：
- gRPC 客户端/服务端支持
- 负载均衡
- 服务发现集成
- 元数据传递

### Gin Transport
基于 Gin Web 框架的 HTTP 传输实现，提供：
- 高性能 HTTP 服务
- 中间件链式处理
- 日志记录和异常恢复
- 请求解码和响应编码自定义

### SSE Transport
Server-Sent Events 传输实现，提供：
- 服务器向客户端的单向数据推送
- 自动重连机制
- 消息类型支持
- 流式数据传输

### Asynq Transport
基于 Redis 的分布式任务队列传输实现，提供：
- 分布式任务处理
- 任务持久化
- 自动重试机制
- 任务优先级支持

### Keepalive Transport
服务保活传输实现，提供：
- 健康检查服务
- 服务存活状态监测
- 与注册中心集成

### Vue Transport
Vue.js 集成传输实现，提供：
- 前端框架集成支持
- 响应式数据传输

## 核心接口

### Server 接口
所有传输服务器都实现了统一的 Server 接口：
```go
type Server interface {
    Start(context.Context) error
    Stop(context.Context) error
}
```

### Transporter 接口
传输层上下文值接口，提供传输相关的信息：
```go
type Transporter interface {
    Kind() Kind                 // 返回传输类型 (grpc, http)
    Endpoint() string           // 返回服务端点
    Operation() string          // 返回操作名称
    RequestHeader() Header      // 返回请求头
    ReplyHeader() Header        // 返回回复头
}
```

## 使用方式

### HTTP 传输示例
```go
import "github.com/chnxq/xkitpkg/transport/http"

srv := http.NewServer(
    http.WithAddress(":8080"),
    http.WithTimeout(1*time.Second),
)

srv.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello World!"))
})

if err := srv.Start(ctx); err != nil {
    panic(err)
}
```

### gRPC 传输示例
```go
import "github.com/chnxq/xkitpkg/transport/grpc"

srv := grpc.NewServer(
    grpc.WithAddress(":9000"),
)

// 注册 gRPC 服务
pb.RegisterGreeterServer(srv, &Service{})

if err := srv.Start(ctx); err != nil {
    panic(err)
}
```

### Gin 传输示例
```go
import "github.com/chnxq/xkitpkg/transport/gin"

srv := gin.NewServer(
    gin.WithAddress(":8080"),
    gin.WithTimeout(1*time.Second),
)

srv.Use(gin.Recovery())
srv.Use(gin.Logger())

srv.GET("/hello", func(c *gin.Context) {
    c.JSON(200, map[string]string{"message": "Hello World!"})
})

if err := srv.Start(ctx); err != nil {
    panic(err)
}
```

## 上下文传递

传输层支持在请求处理过程中传递传输相关的上下文信息：

- `NewServerContext()` / `FromServerContext()` - 服务端上下文
- `NewClientContext()` / `FromClientContext()` - 客户端上下文

## 扩展功能

### 中间件支持
传输层提供中间件支持，可以在请求处理前后添加额外逻辑。

### 编码/解码
支持多种数据格式的编码和解码，包括 JSON、Protobuf、XML、YAML、Form 等。

### 过滤器
支持请求/响应过滤器，用于处理跨域、认证、限流等功能。

## 设计理念

Transport 层的设计遵循以下原则：

1. **接口统一**: 不同传输协议使用相同的接口规范
2. **易于扩展**: 通过接口抽象，便于添加新的传输协议
3. **性能优化**: 针对不同传输协议的特点进行性能优化
4. **功能完整**: 提供传输层所需的完整功能集合
5. **集成友好**: 与其他框架和组件良好集成

## 依赖关系

Transport 模块依赖于 XGoKit 框架的其他组件：
- encoding: 数据编码/解码
- selector: 节点选择器
- logger: 日志记录
- tracer: 链路追踪
- registry: 服务注册与发现
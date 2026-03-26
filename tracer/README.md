# xkitpkg 分布式链路追踪

tracer 包是 xkitpkg 框架中负责分布式链路追踪的核心组件。它基于 OpenTelemetry 实现，提供了灵活的追踪数据导出机制，支持多种追踪后端。

## 功能特性

- **OpenTelemetry 集成**：基于 OpenTelemetry 标准实现分布式追踪
- **多种导出器支持**：支持 OTLP/gRPC、OTLP/HTTP、Zipkin、标准输出等多种导出方式
- **工厂模式**：使用工厂模式动态创建不同类型的追踪导出器
- **全局管理**：提供全局追踪器提供者的管理功能
- **优雅关闭**：支持追踪器提供者的优雅关闭和资源清理
- **资源配置**：自动配置服务名称、版本、环境等资源属性

## 核心组件

### Tracer Provider（追踪器提供者）
- 管理追踪器实例的生命周期
- 负责创建和管理 Span
- 配置采样策略和资源属性

### Span Exporter（Span 导出器）
- 负责将追踪数据导出到后端系统
- 支持批处理和异步导出
- 提供多种协议支持

## 支持的导出器类型

- **OTLP/gRPC** (`otlp-grpc`)：使用 gRPC 协议导出到 OpenTelemetry 后端，默认端口 4317
- **OTLP/HTTP** (`otlp-http`)：使用 HTTP 协议导出到 OpenTelemetry 后端，默认端口 4318
- **Zipkin** (`zipkin`)：导出到 Zipkin 后端，默认地址 http://localhost:9411/api/v2/spans
- **Standard Output** (`std`)：输出到标准输出，主要用于调试
- **阿里云** (`aliyun`)：导出到阿里云链路追踪服务
- **腾讯云** (`tencent`)：导出到腾讯云链路追踪服务
- **Jaeger** (`jaeger`)：Jaeger 导出器（当前版本不支持，使用 OTLP 替代）

## 主要接口

### ExporterFactory
```go
type ExporterFactory func(ctx context.Context, cfg *conf.Tracer) (traceSdk.SpanExporter, error)
```
导出器工厂函数，用于创建特定类型的追踪导出器实例。

## 核心功能

### 创建追踪导出器

```go
// 根据配置创建追踪导出器
exp, err := tracer.NewTracerExporter(ctx, cfg)
```

### 创建追踪器提供者

```go
// 创建追踪器提供者并设置为全局
err := tracer.NewTracerProvider(ctx, cfg, appInfo)

// 或者获取关闭函数以便资源管理
tp, shutdown, err := tracer.NewTracerProviderWithShutdown(ctx, cfg, appInfo)
defer shutdown(ctx)
```

### 关闭追踪器提供者

```go
// 优雅关闭全局追踪器提供者
err := tracer.ShutdownTracerProvider(ctx)
```

## 配置选项

追踪配置通常包含以下字段：

```yaml
trace:
  exporter: "otlp-grpc"        # 导出器类型
  endpoint: "localhost:4317"   # 后端服务地址
  sampler: 1.0                 # 采样率 (0.0-1.0)
  env: "dev"                   # 环境标识
  insecure: true               # 是否使用非安全连接
  enable_trace_context: true   # 是否启用追踪上下文
  enable_baggage: true         # 是否启用 baggage 传递
  batcher_options:             # 批处理器选项
    enabled: true
    max_queue_size: 2048
    max_export_batch_size: 512
    schedule_delay_millis: 5000
    export_timeout_millis: 30000
```

## 使用示例

```go
import (
    "context"
    "github.com/chnxq/xkitpkg/tracer"
    "github.com/chnxq/xkitpkg/conf/v1"
)

func main() {
    ctx := context.Background()
    
    // 假设已有配置和应用信息
    var cfg *conf.Tracer
    var appInfo *conf.AppInfo
    
    // 创建追踪器提供者
    tp, shutdown, err := tracer.NewTracerProviderWithShutdown(ctx, cfg, appInfo)
    if err != nil {
        panic(err)
    }
    defer shutdown(ctx)
    
    // 此后所有的追踪操作都会使用这个提供者
    // 应用程序可以使用 otel.Tracer() 来创建追踪器
}
```

## 工厂模式

tracer 包使用工厂模式来注册和创建不同类型的导出器：

```go
// 注册自定义导出器工厂
tracer.RegisterExporter("custom", func(ctx context.Context, cfg *conf.Tracer) (traceSdk.SpanExporter, error) {
    // 创建自定义导出器
    return customExporter, nil
})
```

## 线程安全性

tracer 包使用读写锁保护共享资源，确保在高并发环境下能够安全使用。

## 扩展性

通过工厂模式设计，可以轻松扩展支持新的追踪后端，只需实现相应的导出器工厂函数并注册即可。
# xkitpkg 服务注册与发现

registry 包是 xkitpkg 框架中负责服务注册与发现的核心组件。它提供了统一的接口来集成各种服务注册中心，如 Etcd、Consul、ZooKeeper、Nacos 等。

## 功能特性

- **多注册中心支持**：支持多种服务注册中心，包括 Etcd、Consul、ZooKeeper、Nacos、Kubernetes、Eureka、Polaris 和 Servicecomb
- **工厂模式**：使用工厂模式动态创建不同类型的注册和发现实例
- **配置驱动**：通过配置文件动态选择和初始化注册中心
- **线程安全**：使用读写锁保证并发安全

## 核心组件

### Registrar（服务注册器）
- 负责向服务注册中心注册和注销服务
- 提供心跳机制保持服务健康状态

### Discovery（服务发现器）
- 负责从服务注册中心查询可用服务
- 提供负载均衡和服务节点选择功能

## 主要接口

### RegistrarFactory
```go
type RegistrarFactory func(cfg *conf.Registry) (registry.Registrar, error)
```
注册器工厂函数，用于创建特定类型的注册器实例。

### DiscoveryFactory
```go
type DiscoveryFactory func(cfg *conf.Registry) (registry.Discovery, error)
```
发现器工厂函数，用于创建特定类型的发现器实例。

## 使用方法

### 注册工厂函数

```go
// 注册注册器工厂
err := registry.RegisterRegistrarFactory(registry.Etcd, etcdreg.RegistrarFactory)

// 注册发现器工厂
err := registry.RegisterDiscoveryFactory(registry.Etcd, etcdreg.DiscoveryFactory)
```

### 创建实例

```go
// 创建注册器
registrar, err := registry.NewRegistrar(config.GetServerConfig().GetRegistry())

// 创建发现器
discovery, err := registry.NewDiscovery(config.GetServerConfig().GetRegistry())
```

## 支持的注册中心类型

- **Etcd**：基于 etcd 的服务注册与发现
- **Consul**：HashiCorp Consul 服务注册与发现
- **ZooKeeper**：Apache ZooKeeper 服务注册与发现
- **Nacos**：阿里巴巴 Nacos 服务注册与发现
- **Kubernetes**：Kubernetes 原生服务发现
- **Eureka**：Netflix Eureka 服务注册与发现
- **Polaris**：腾讯 Polaris 服务治理
- **Servicecomb**：Apache ServiceComb 服务注册与发现

## 配置示例

配置文件应包含以下结构：

```yaml
registry:
  type: "ETCD"  # 或其他支持的类型
  etcd:
    endpoints:
      - "localhost:2379"
  consul:
    scheme: "http"
    address: "localhost:8500"
  # ... 其他注册中心配置
```

## 错误处理

- 当不支持的注册中心类型被请求时，会返回明确的错误信息
- 提供可用的注册中心类型列表帮助用户正确配置
- 空配置会被安全处理，返回 nil 实例而不会导致崩溃

## 线程安全性

registry 包使用读写锁保护共享资源，确保在高并发环境下能够安全使用。

## 扩展性

通过工厂模式设计，可以轻松扩展支持新的服务注册中心，只需实现相应的工厂函数并注册即可。
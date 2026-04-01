# Server Utils Package

The `server_utils` package provides utilities for building REST and gRPC services with built-in middleware support (including whitelisting, authentication protection, rate limiting, etc.). This documentation is intended for developers and system administrators who need to implement and maintain services using this package.

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Core Components](#core-components)
- [Middleware](#middleware)
- [Whitelist System](#whitelist-system)
- [API Reference](#api-reference)
- [Best Practices](#best-practices)
- [Testing](#testing)
- [Examples](#examples)

## Overview

The `server_utils` package offers a comprehensive solution for HTTP (Kratos REST) and gRPC service development. It provides:

- Service creation and configuration utilities
- Built-in middleware registration points
- Common middleware implementations (whitelist, authentication, rate limiting, etc.)
- Flexible configuration options
- Concurrent-safe utilities for production environments

### Design Principles

- Never panic within the library
- Return errors for configuration/initialization failures
- Support configurable middleware with runtime whitelist bypassing
- Thread-safe operations for concurrent environments
- Modular architecture allowing easy extension

## Features

- **REST Server Creation**: Easy setup of Kratos-based REST servers
- **gRPC Support**: Unified approach for both REST and gRPC services
- **Built-in Middleware**: Comprehensive collection of common middleware
- **Whitelist System**: Flexible mechanism to bypass specific middleware
- **Configuration Driven**: YAML-based configuration system
- **Security Focused**: Built-in authentication and authorization support
- **Performance Optimized**: Efficient implementations for high-throughput scenarios

## Installation

To use the `server_utils` package, you need to have Go installed (version 1.18 or later). Then, you can import it in your project:

```bash
go get github.com/chnxq/xkitpkg/server_utils
```

## Quick Start

### Basic REST Server Setup

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/chnxq/xkitpkg/server_utils"
    "github.com/chnxq/xkitpkg/conf"
    kratosRest "github.com/chnxq/XGoKit/transport/http"
)

func main() {
    // Load configuration
    cfg := &conf.ServerConfig{
        Server: &conf.Server{
            Rest: &conf.Rest{
                Addr: ":8080",
                Timeout: time.Duration(30 * time.Second),
                Middleware: &conf.Middleware{
                    Recovery: &conf.Recovery{Enabled: true},
                    Validate: &conf.Validate{Enabled: true},
                },
            },
        },
    }
    
    // Create REST server
    srv, err := server_utils.CreateRestServer(cfg)
    if err != nil {
        log.Fatalf("Failed to create REST server: %v", err)
    }
    
    // Register your service handlers
    // srv.HandlePrefix("/", yourServiceHandler)
    
    // Start the server
    if err := srv.Start(context.Background()); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
```

### Adding Custom Middleware

```go
import (
    "context"
    "errors"
    
    "github.com/chnxq/XGoKit/middleware"
    "github.com/chnxq/XGoKit/transport"
)

// Custom authentication middleware
func AuthMiddleware() middleware.Middleware {
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
            // Extract transport info
            if tr, ok := transport.FromServerContext(ctx); ok {
                // Perform authentication based on transport info
                if !isValidRequest(tr) {
                    return nil, errors.New("unauthorized access")
                }
            }
            
            return handler(ctx, req)
        }
    }
}

// Use with the server_utils package
srv, err := server_utils.CreateRestServer(cfg, AuthMiddleware())
if err != nil {
    log.Fatalf("Failed to create server: %v", err)
}
```

## Configuration

### Basic Configuration Structure

```yaml
server:
  rest:
    addr: ":8080"
    timeout: "30s"
    network: "tcp"
    cors:
      allow_origins: ["https://example.com", "https://sub.example.com"]
      allow_methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"]
      allow_headers: ["Content-Type", "Authorization", "X-Requested-With", "X-Forwarded-For"]
      allow_credentials: true
      max_age: "12h"
    middleware:
      recovery:
        enabled: true
      tracing:
        enabled: true
        exporter: "otlp"  # otlp, jaeger, zipkin
      validate:
        enabled: true
      metadata:
        enabled: true
      ratelimit:
        enabled: true
        strategy: "bbr"  # bbr, token_bucket, sliding_window
        burst: 100       # burst capacity for rate limiter
        qps: 50          # queries per second limit
      logging:
        enabled: true
        level: "info"
    tls:
      enabled: false
      cert_file: "/path/to/certificate.pem"
      key_file: "/path/to/private.key"
```

### Configuration Options Explained

| Option | Type | Description |
|--------|------|-------------|
| `addr` | string | Server address and port |
| `timeout` | duration | Request timeout duration |
| `network` | string | Network type (tcp, tcp4, tcp6) |
| `cors.allow_origins` | []string | Allowed origins for CORS |
| `cors.allow_methods` | []string | Allowed HTTP methods |
| `cors.allow_headers` | []string | Allowed headers |
| `cors.allow_credentials` | bool | Whether to allow credentials |
| `cors.max_age` | duration | How long browsers can cache preflight responses |
| `middleware.*.enabled` | bool | Enable/disable specific middleware |
| `middleware.ratelimit.strategy` | string | Rate limiting strategy |
| `middleware.ratelimit.burst` | int | Burst capacity for rate limiter |
| `middleware.ratelimit.qps` | int | Queries per second limit |
| `tls.enabled` | bool | Enable/disable TLS |
| `tls.cert_file` | string | Path to certificate file |
| `tls.key_file` | string | Path to private key file |

## Core Components

### Server Creation Functions

#### `CreateRestServer(cfg *conf.ServerConfig, mds ...middleware.Middleware) (*kratosRest.Server, error)`
Creates and configures a REST server instance based on the provided configuration. Additional middleware can be passed as variadic parameters.

#### `NewRestWhiteListMatcher() selector.MatchFunc`
Returns a default whitelist matcher function that can be used with selector middleware to bypass certain middleware for specific operations.

### Whitelist Management Functions

#### `AddWhiteList(ops ...string)`
Adds operations to the global whitelist. Operations can be full method names or just method names.

#### `SetWhiteList(ops []string)`
Replaces the entire whitelist with the provided operations.

#### `ClearWhiteList()`
Clears all operations from the whitelist.

#### `NewWhiteListMatcher() selector.MatchFunc`
Creates a new whitelist matcher instance (independent of the global whitelist).

## Middleware

The package includes several built-in middleware components:

### Core Middleware

- **Recovery**: Handles panics gracefully and returns proper error responses
- **RequestID**: Generates and injects unique request IDs for tracing
- **Logging**: Structured logging of requests and responses
- **Metrics**: Prometheus metrics collection
- **Tracing**: OpenTelemetry distributed tracing
- **Authentication**: JWT/API key based authentication
- **Rate Limiting**: Per-client request rate limiting
- **Validation**: Request validation based on protobuf definitions
- **Metadata**: Extracts and injects metadata from requests

### Middleware Order

Middleware is applied in this general order (outermost to innermost):

1. Recovery
2. RequestID
3. Logging
4. Metrics
5. Tracing
6. Authentication
7. Rate Limiting
8. Validation
9. Metadata
10. Business Logic

### Custom Middleware Example

```go
import (
    "context"
    
    "github.com/chnxq/XGoKit/errors"
    "github.com/chnxq/XGoKit/middleware"
)

func CustomAuthMiddleware() middleware.Middleware {
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
            // Add custom authentication logic here
            if !isValidRequest(ctx) {
                return nil, errors.Unauthorized("AUTH_ERROR", "Invalid credentials")
            }
            
            return handler(ctx, req)
        }
    }
}
```

## Whitelist System

The whitelist system allows specific endpoints to bypass certain middleware, which is useful for health checks, public endpoints, and internal monitoring.

### Whitelist Management

```go
// Add specific operations to whitelist
server_utils.AddWhiteList("/health.check", "/public.api.Endpoint")

// Replace entire whitelist
server_utils.SetWhiteList([]string{"/health.check", "/metrics", "/public.api.Endpoint"})

// Clear whitelist
server_utils.ClearWhiteList()

// Get matcher function for use with selector
matcher := server_utils.NewRestWhiteListMatcher()
```

### Using Whitelist with Middleware

```go
import (
    "github.com/chnxq/XGoKit/middleware/selector"
    "github.com/chnxq/XGoKit/middleware/validate"
)

// Apply middleware with whitelist bypass
ms := make([]middleware.Middleware, 0)

// Skip validation for whitelisted operations
ms = append(ms, selector.Server(
    validate.Validator(),
    server_utils.NewRestWhiteListMatcher(),
).Build())

// Skip rate limiting for whitelisted operations
ms = append(ms, selector.Server(
    yourRateLimitMiddleware,
    server_utils.NewRestWhiteListMatcher(),
).Build())
```

### Whitelist Matching Modes

- **Exact Match**: Matches the exact operation name (e.g., `/package.Service/Method`)
- **Method-Only Match**: Matches just the method name part (e.g., `Method`)

## API Reference

### Public Functions

#### `func CreateRestServer(cfg *conf.ServerConfig, mds ...middleware.Middleware) (*kratosRest.Server, error)`
Creates a new REST server with the provided configuration and optional additional middleware.

Parameters:
- `cfg`: *conf.ServerConfig - Server configuration containing server settings
- `mds`: Optional additional middleware to apply after configured middleware

Returns:
- `*kratosRest.Server`: The configured server instance
- `error`: Any error during server creation

#### `func AddWhiteList(ops ...string)`
Adds operations to the global whitelist.

Parameters:
- `ops`: Operation names to add to the whitelist

#### `func SetWhiteList(ops []string)`
Sets the entire whitelist to the provided operations.

Parameters:
- `ops`: Slice of operation names for the new whitelist

#### `func ClearWhiteList()`
Clears all operations from the global whitelist.

#### `func NewRestWhiteListMatcher() selector.MatchFunc`
Creates a matcher function that returns false for whitelisted operations, allowing them to bypass selected middleware.

Returns:
- `selector.MatchFunc`: Matcher function for use with selector middleware

## Best Practices

### Security

- Restrict access to `/debug/pprof` and `/metrics` endpoints
- Use HTTPS in production environments
- Implement proper authentication and authorization
- Regularly rotate API keys and certificates
- Sanitize and validate all inputs

### Performance

- Use Redis or other distributed storage for rate limiting in high-concurrency scenarios
- Configure appropriate timeouts for different types of requests
- Monitor and tune middleware performance
- Consider caching for frequently accessed data
- Use connection pooling where appropriate

### Configuration Management

- Use environment-specific configuration files
- Implement configuration validation
- Support dynamic configuration reloading where appropriate
- Use secure storage for sensitive configuration values
- Document configuration changes and impacts

### Monitoring and Observability

- Enable structured logging
- Configure distributed tracing
- Expose metrics for monitoring systems
- Implement health check endpoints
- Set up alerts for critical metrics
- Monitor error rates and response times

## Testing

### Running Tests

```bash
# Run package tests with race detection
go test ./server_utils -race -v

# Run all tests with race detection
go test ./... -race

# Run benchmarks
go test ./server_utils -bench=. -benchmem

# Generate coverage report
go test ./server_utils -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Test Coverage

- Unit tests for individual functions
- Integration tests for middleware chains
- Race condition testing with `-race` flag
- Benchmark tests for performance-critical components
- Concurrent access testing for thread safety

## Examples

### Complete Service Example

```go
package main

import (
    "context"
    "log"
    "net/http"
    
    "github.com/chnxq/XGoKit/middleware/recovery"
    "github.com/chnxq/XGoKit/middleware/tracing"
    kratosRest "github.com/chnxq/XGoKit/transport/http"
    "github.com/chnxq/xkitpkg/server_utils"
    "github.com/chnxq/xkitpkg/conf"
)

func main() {
    // Initialize configuration
    cfg := &conf.ServerConfig{
        Server: &conf.Server{
            Rest: &conf.Rest{
                Addr: ":8080",
                Middleware: &conf.Middleware{
                    Recovery: &conf.Recovery{Enabled: true},
                    Tracing:  &conf.Tracing{Enabled: true},
                },
            },
        },
    }
    
    // Create the server
    srv, err := server_utils.CreateRestServer(cfg)
    if err != nil {
        log.Fatalf("Failed to create server: %v", err)
    }
    
    // Add custom routes
    srv.HandleFunc("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
        name := r.PathValue("name")
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte("{\"message\": \"Hello, " + name + "!\"}"))
    })
    
    // Add health check to whitelist
    server_utils.AddWhiteList("/hello/{name}")
    
    // Start the server
    log.Println("Starting server on :8080")
    if err := srv.Start(context.Background()); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
```

### Middleware Composition Example

```go
package main

import (
    "context"
    "log"
    
    "github.com/chnxq/XGoKit/middleware"
    "github.com/chnxq/XGoKit/middleware/logging"
    "github.com/chnxq/XGoKit/middleware/recovery"
    "github.com/chnxq/xkitpkg/server_utils"
    "github.com/chnxq/xkitpkg/conf"
)

func main() {
    cfg := &conf.ServerConfig{
        Server: &conf.Server{
            Rest: &conf.Rest{},
        },
    }
    
    // Define custom middleware chain
    customMiddleware := []middleware.Middleware{
        recovery.Recovery(),
        logging.Logger(),
        // Add your custom middleware here
    }
    
    srv, err := server_utils.CreateRestServer(cfg, customMiddleware...)
    if err != nil {
        log.Fatalf("Failed to create server: %v", err)
    }
    
    // Start server
    if err := srv.Start(context.Background()); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
```

## Contributing

We welcome contributions to enhance this package! Here are some areas where contributions would be valuable:

- Additional middleware implementations
- Improved configuration options
- Enhanced security features
- Better documentation and examples
- Performance optimizations
- Additional test coverage

### How to Contribute

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for your changes
5. Update documentation as needed
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a pull request

### Code Standards

- Follow Go coding conventions
- Include comprehensive tests
- Update documentation as needed
- Ensure backward compatibility when possible
- Write clear commit messages
- Keep pull requests focused on a single feature or fix

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support, please open an issue in the repository or contact the maintainers.
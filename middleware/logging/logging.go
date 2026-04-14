package logging

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/chnxq/xkitmod/errors"
	"github.com/chnxq/xkitmod/log"
	"github.com/chnxq/xkitpkg/middleware"
	"github.com/chnxq/xkitpkg/transport"
	"github.com/chnxq/xkitpkg/transport/http/status"
)

// Redacter defines how to log an object
type Redacter interface {
	Redact() string
}

// Server 创建一个服务器端日志中间件
// 该中间件用于记录服务器端请求的详细日志信息，包括请求类型、操作、参数、响应码、错误信息和执行时间等
//
// 参数:
//
//	logger log.Logger - 日志记录器实例，用于输出日志信息
//
// 返回值:
//
//	middleware.Middleware - 服务器端日志中间件实例，可以被插入到请求处理链中
//
// 功能说明:
//  1. 记录请求开始时间以计算处理延迟
//  2. 从传输上下文中提取请求类型(kind)和操作(operation)信息
//  3. 执行原始请求处理器
//  4. 分析响应错误并提取错误码和原因
//  5. 根据错误级别决定日志级别并记录完整的请求信息
//  6. 包括延迟时间在内的所有相关信息都会被记录
func Server(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			var (
				code      int32
				reason    string
				kind      string
				operation string
			)

			// default code
			code = int32(status.FromGRPCCode(codes.OK))

			startTime := time.Now()
			if info, ok := transport.FromServerContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err = handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
			}
			level, stack := extractError(err)
			log.NewHelper(log.WithContext(ctx, logger)).Log(level,
				"kind", kind,
				"component", "server",
				"operation", operation,
				"args", extractArgs(req),
				"code", code,
				"reason", reason,
				"stack", stack,
				"latency", time.Since(startTime).Seconds(),
			)
			return
		}
	}
}

// Client is a client logging middleware.
func Client(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			var (
				code      int32
				reason    string
				kind      string
				operation string
			)

			// default code
			code = int32(status.FromGRPCCode(codes.OK))

			startTime := time.Now()
			if info, ok := transport.FromClientContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err = handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
			}
			level, stack := extractError(err)
			log.NewHelper(log.WithContext(ctx, logger)).Log(level,
				"kind", kind,
				"component", "client",
				"operation", operation,
				"args", extractArgs(req),
				"code", code,
				"reason", reason,
				"stack", stack,
				"latency", time.Since(startTime).Seconds(),
			)
			return
		}
	}
}

// extractArgs returns the string of the req
func extractArgs(req any) string {
	if redacter, ok := req.(Redacter); ok {
		return redacter.Redact()
	}
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}

// extractError returns the string of the error
func extractError(err error) (log.Level, string) {
	if err != nil {
		return log.LevelError, fmt.Sprintf("%+v", err)
	}
	return log.LevelInfo, ""
}

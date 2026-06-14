package server_utils

import (
	"github.com/chnxq/xkitmod/algs/ratelimit"
	"github.com/chnxq/xkitmod/algs/ratelimit/bbr"
	"github.com/chnxq/xkitpkg/app"
	conf "github.com/chnxq/xkitpkg/conf/v1"
	"github.com/chnxq/xkitpkg/middleware"
	midbreaker "github.com/chnxq/xkitpkg/middleware/circuitbreaker"
	"github.com/chnxq/xkitpkg/middleware/logging"
	"github.com/chnxq/xkitpkg/middleware/metadata"
	midratelimit "github.com/chnxq/xkitpkg/middleware/ratelimit"
	"github.com/chnxq/xkitpkg/middleware/recovery"
	"github.com/chnxq/xkitpkg/middleware/tracing"
	"github.com/chnxq/xkitpkg/middleware/validate"
)

type DatabaseLoggingData interface {
	DatabaseLoggingMiddleware() middleware.Middleware
}

func ServerConfig(appCtx *app.AppCtx) *conf.ServerConfig {
	if appCtx == nil {
		return nil
	}
	return appCtx.GetConfig()
}

func RestConfig(appCtx *app.AppCtx) *conf.Server_REST {
	cfg := ServerConfig(appCtx)
	if cfg == nil || cfg.Server == nil {
		return nil
	}
	return cfg.Server.Rest
}

func GrpcConfig(appCtx *app.AppCtx) *conf.Server_GRPC {
	cfg := ServerConfig(appCtx)
	if cfg == nil || cfg.Server == nil {
		return nil
	}
	return cfg.Server.Grpc
}

func DatabaseLoggingMiddleware(data DatabaseLoggingData) middleware.Middleware {
	if data != nil {
		if mw := data.DatabaseLoggingMiddleware(); mw != nil {
			return mw
		}
	}
	return func(next middleware.Handler) middleware.Handler {
		return next
	}
}

func CommonServerMiddlewares(appCtx *app.AppCtx, cfg *conf.Middleware) []middleware.Middleware {
	if cfg == nil {
		return nil
	}

	ms := make([]middleware.Middleware, 0, 7)
	if cfg.GetEnableRecovery() {
		ms = append(ms, recovery.Recovery())
	}
	if cfg.GetEnableTracing() {
		ms = append(ms, tracing.Server())
	}
	if cfg.GetEnableValidate() {
		ms = append(ms, validate.ProtoValidate())
	}
	if cfg.GetEnableMetadata() {
		ms = append(ms, metadata.Server())
	}
	if cfg.GetLimiter() != nil {
		if limiter := NewServerRateLimiter(cfg.GetLimiter()); limiter != nil {
			ms = append(ms, midratelimit.Server(midratelimit.WithLimiter(limiter)))
		}
	}
	if cfg.GetEnableCircuitBreaker() {
		ms = append(ms, midbreaker.Server(midbreaker.WithSREConfig(cfg.GetCircuitBreaker())))
	}
	if cfg.GetEnableLogging() && appCtx != nil {
		ms = append(ms, logging.Server(appCtx.GetLogger()))
	}
	return ms
}

func NewServerRateLimiter(cfg *conf.Middleware_RateLimiter) ratelimit.Limiter {
	if cfg == nil {
		return nil
	}
	switch cfg.GetName() {
	case "", "bbr":
		opts := make([]bbr.Option, 0, 4)
		if d := cfg.GetWindow(); d != nil {
			opts = append(opts, bbr.WithWindow(d.AsDuration()))
		}
		if v := cfg.GetBucket(); v > 0 {
			opts = append(opts, bbr.WithBucket(int(v)))
		}
		if v := cfg.GetCpuThreshold(); v > 0 {
			opts = append(opts, bbr.WithCPUThreshold(v))
		}
		if v := cfg.GetCpuQuota(); v > 0 {
			opts = append(opts, bbr.WithCPUQuota(v))
		}
		return bbr.NewLimiter(opts...)
	default:
		return nil
	}
}

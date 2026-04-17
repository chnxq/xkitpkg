package server_utils

import (
	"crypto/tls"
	"net/http/pprof"

	"github.com/gorilla/handlers"

	"github.com/chnxq/xkitmod/algs/ratelimit"
	"github.com/chnxq/xkitmod/algs/ratelimit/bbr"

	kHttp "github.com/chnxq/xkitpkg/transport/http"

	"github.com/chnxq/xkitpkg/middleware"
	"github.com/chnxq/xkitpkg/middleware/metadata"
	midRateLimit "github.com/chnxq/xkitpkg/middleware/ratelimit"
	"github.com/chnxq/xkitpkg/middleware/recovery"
	"github.com/chnxq/xkitpkg/middleware/selector"
	"github.com/chnxq/xkitpkg/middleware/tracing"

	conf "github.com/chnxq/xkitpkg/conf/v1"
	"github.com/chnxq/xkitpkg/middleware/validate"
)

// CreateRestServer 创建REST服务端
// cfg 服务配置
// mds 中间件集合 (通常包含：访问日志，认证授权...)
func CreateRestServer(cfg *conf.ServerConfig, mds ...middleware.Middleware) (*kHttp.Server, error) {
	options, err := initRestConfig(cfg, mds...)
	if err != nil {
		return nil, err
	}

	srv := kHttp.NewServer(options...)

	if cfg != nil && cfg.Server != nil && cfg.Server.Rest != nil && cfg.Server.Rest.GetEnablePprof() {
		registerHttpPprof(srv)
	}

	return srv, nil
}

// initRestConfig 初始化REST服务配置
// include：recovery, tracing, validate, circuit breaker, rate limit, metadata)
func initRestConfig(cfg *conf.ServerConfig, mds ...middleware.Middleware) ([]kHttp.ServerOption, error) {
	if cfg == nil || cfg.Server == nil || cfg.Server.Rest == nil {
		return nil, nil
	}

	var options []kHttp.ServerOption

	if cfg.Server.Rest.Cors != nil {
		options = append(options, kHttp.Filter(handlers.CORS(
			handlers.AllowedHeaders(cfg.Server.Rest.Cors.Headers),
			handlers.AllowedMethods(cfg.Server.Rest.Cors.Methods),
			handlers.AllowedOrigins(cfg.Server.Rest.Cors.Origins),
		)))
	}

	var ms []middleware.Middleware
	if cfg.Server.Rest.Middleware != nil {
		if cfg.Server.Rest.Middleware.GetEnableRecovery() {
			ms = append(ms, recovery.Recovery())
		}
		if cfg.Server.Rest.Middleware.GetEnableTracing() {
			ms = append(ms, tracing.Server())
		}
		if cfg.Server.Rest.Middleware.GetEnableValidate() {
			ms = append(ms, validate.ProtoValidate())
		}
		if cfg.Server.Rest.Middleware.GetEnableCircuitBreaker() {
			//ms = append(ms, circuitbreaker.NewBreaker()) //Fixme: 待实现 XQ
		}
		if cfg.Server.Rest.Middleware.Limiter != nil {
			var limiter ratelimit.Limiter
			switch cfg.Server.Rest.Middleware.Limiter.GetName() {
			case "bbr":
				limiter = bbr.NewLimiter()
			}
			ms = append(ms, midRateLimit.Server(midRateLimit.WithLimiter(limiter)))
		}
		if cfg.Server.Rest.Middleware.GetEnableMetadata() {
			ms = append(ms, metadata.Server())
		}
	}
	ms = append(ms, mds...)

	options = append(options, kHttp.Middleware(ms...))

	if cfg.Server.Rest.Network != "" {
		options = append(options, kHttp.Network(cfg.Server.Rest.Network))
	}
	if cfg.Server.Rest.Addr != "" {
		options = append(options, kHttp.Address(cfg.Server.Rest.Addr))
	}
	if cfg.Server.Rest.Timeout != nil {
		options = append(options, kHttp.Timeout(cfg.Server.Rest.Timeout.AsDuration()))
	}

	if cfg.Server.Rest.Tls != nil {
		var tlsCfg *tls.Config
		var err error

		if tlsCfg, err = loadServerTlsConfig(cfg.Server.Rest.Tls); err != nil {
			return nil, err
		}

		if tlsCfg != nil {
			options = append(options, kHttp.TLSConfig(tlsCfg))
		}
	}

	return options, nil
}

// registerHttpPprof 注册pprof路由
func registerHttpPprof(s *kHttp.Server) {
	s.HandleFunc("/debug/pprof", pprof.Index)

	s.HandleFunc("/debug/cmdline", pprof.Cmdline)
	s.HandleFunc("/debug/profile", pprof.Profile)
	s.HandleFunc("/debug/symbol", pprof.Symbol)
	s.HandleFunc("/debug/trace", pprof.Trace)

	s.HandleFunc("/debug/allocs", pprof.Handler("allocs").ServeHTTP)
	s.HandleFunc("/debug/block", pprof.Handler("block").ServeHTTP)
	s.HandleFunc("/debug/goroutine", pprof.Handler("goroutine").ServeHTTP)
	s.HandleFunc("/debug/heap", pprof.Handler("heap").ServeHTTP)
	s.HandleFunc("/debug/mutex", pprof.Handler("mutex").ServeHTTP)
	s.HandleFunc("/debug/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
}

// NewRestWhiteListMatcher 创建REST白名单匹配器
func NewRestWhiteListMatcher() selector.MatchFunc {
	// reuse package-level DefaultWhiteList matcher for REST
	return NewWhiteListMatcher()
}

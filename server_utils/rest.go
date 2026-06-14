package server_utils

import (
	"net/http/pprof"

	"github.com/gorilla/handlers"

	kHttp "github.com/chnxq/xkitpkg/transport/http"

	"github.com/chnxq/xkitpkg/middleware"
	"github.com/chnxq/xkitpkg/middleware/selector"

	conf "github.com/chnxq/xkitpkg/conf/v1"
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

	ms := CommonServerMiddlewares(nil, cfg.Server.Rest.Middleware)
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
		tlsCfg, err := LoadServerTLSConfig(cfg.Server.Rest.Tls)
		if err != nil {
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

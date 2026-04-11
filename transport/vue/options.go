package vue

import (
	"crypto/tls"
	"time"

	"github.com/chnxq/xkitpkg/conf/v1"
)

type ServerOption func(o *Server)

func WithNetwork(network string) ServerOption {
	return func(s *Server) {
		s.network = network
	}
}

func WithAddress(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

// WithTimeout with server timeout.
func WithTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func WithTLSConfig(c *tls.Config) ServerOption {
	return func(o *Server) {
		o.tlsConf = c
	}
}

func WithProxy(proxy []*conf.Proxy) ServerOption {
	return func(s *Server) {
		s.proxy = proxy
	}
}

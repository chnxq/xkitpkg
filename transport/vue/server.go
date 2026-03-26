package vue

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/chnxq/XGoKit/conf"
	"github.com/chnxq/XGoKit/log"
	"github.com/chnxq/xkitpkg/transport"
	"github.com/chnxq/xkitpkg/transport/internal/endpoint"
	"github.com/chnxq/xkitpkg/transport/internal/host"
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
)

type Server struct {
	*http.Server
	lis      net.Listener
	tlsConf  *tls.Config
	endpoint *url.URL
	err      error
	network  string
	address  string
	timeout  time.Duration

	proxy []*conf.Proxy
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network: "tcp",
		address: ":0",
		timeout: 1 * time.Second,
	}

	for _, o := range opts {
		o(srv)
	}

	srv.Server = &http.Server{
		TLSConfig: srv.tlsConf,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, item := range srv.proxy {
			prefix := item.Prefix
			if prefix == "" {
				prefix = "/"
			}
			if !strings.HasPrefix(prefix, "/") {
				prefix = "/" + prefix
			}
			if strings.HasPrefix(r.URL.Path, prefix) {
				// 判断代理地址是否是url
				if item.IsUrl {
					backend, err := url.Parse(item.Addr)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						log.Error(err)
						return
					}
					// 创建一个新的代理对象，并修改请求的Path以匹配后端API的路径
					proxy := httputil.NewSingleHostReverseProxy(backend)
					// 使用代理对象处理请求
					proxy.ServeHTTP(w, r)
					break
				} else {
					http.StripPrefix(prefix, http.FileServer(http.Dir(item.Addr))).ServeHTTP(w, r)
					break
				}
			}
		}
	})

	return srv
}

func (s *Server) Name() string {
	return string(KindVue)
}

func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.listenAndEndpoint(); err != nil {
		return nil, err
	}
	return s.endpoint, nil
}

// Start start the HTTP server.
func (s *Server) Start(ctx context.Context) error {
	if err := s.listenAndEndpoint(); err != nil {
		return err
	}
	s.BaseContext = func(net.Listener) context.Context {
		return ctx
	}
	log.Infof("[VUE] server listening on: %s", s.lis.Addr().String())
	var err error
	if s.tlsConf != nil {
		err = s.ServeTLS(s.lis, "", "")
	} else {
		err = s.Serve(s.lis)
	}
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop stop the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	log.Info("[VUE] server stopping")
	return s.Shutdown(ctx)
}

func (s *Server) listenAndEndpoint() error {
	if s.lis == nil {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			s.err = err
			return err
		}
		s.lis = lis
	}
	if s.endpoint == nil {
		addr, err := host.Extract(s.address, s.lis)
		if err != nil {
			s.err = err
			return err
		}
		s.endpoint = endpoint.NewEndpoint(endpoint.Scheme("grpc", s.tlsConf != nil), addr)
	}
	return s.err
}

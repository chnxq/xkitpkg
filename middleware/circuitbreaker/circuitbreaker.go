package circuitbreaker

import (
	"context"
	stderrors "errors"

	"github.com/chnxq/xkitmod/algs/circuitbreaker"
	"github.com/chnxq/xkitmod/algs/circuitbreaker/sre"

	"github.com/chnxq/xkitmod/errors"
	conf "github.com/chnxq/xkitpkg/conf/v1"
	"github.com/chnxq/xkitpkg/middleware"

	"github.com/chnxq/xkitpkg/internal/group"
	"github.com/chnxq/xkitpkg/transport"
)

// ErrNotAllowed is request failed due to circuit breaker triggered.
var ErrNotAllowed = errors.New(503, "CIRCUITBREAKER", "request failed due to circuit breaker triggered")

// Option is circuit breaker option.
type Option func(*options)

// WithGroup with circuit breaker group.
// NOTE: implements generics circuitbreaker.CircuitBreaker
func WithGroup(g *group.Group[circuitbreaker.CircuitBreaker]) Option {
	return func(o *options) {
		o.group = g
	}
}

// WithCircuitBreaker with circuit breaker genFunc.
func WithCircuitBreaker(genBreakerFunc func() circuitbreaker.CircuitBreaker) Option {
	return func(o *options) {
		o.group = group.NewGroup(func() circuitbreaker.CircuitBreaker {
			return genBreakerFunc()
		})
	}
}

type options struct {
	group *group.Group[circuitbreaker.CircuitBreaker]
}

// Client circuitbreaker middleware will return errBreakerTriggered when the circuit
// breaker is triggered and the request is rejected directly.
func Client(opts ...Option) middleware.Middleware {
	opt := &options{
		group: group.NewGroup(func() circuitbreaker.CircuitBreaker {
			return sre.NewBreaker()
		}),
	}
	for _, o := range opts {
		o(opt)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			info, _ := transport.FromClientContext(ctx)
			breaker := opt.group.Get(info.Operation())
			if err := breaker.Allow(); err != nil {
				// rejected
				// NOTE: when client reject requests locally,
				// continue to add counter let the drop ratio higher.
				breaker.MarkFailed()
				return nil, ErrNotAllowed
			}
			// allowed
			reply, err := handler(ctx, req)
			if err != nil && (errors.IsInternalServer(err) || errors.IsServiceUnavailable(err) || errors.IsGatewayTimeout(err)) {
				breaker.MarkFailed()
			} else {
				breaker.MarkSuccess()
			}
			return reply, err
		}
	}
}

func WithSREConfig(cfg *conf.Middleware_CircuitBreaker) Option {
	return WithCircuitBreaker(func() circuitbreaker.CircuitBreaker {
		opts := make([]sre.Option, 0, 4)
		if cfg != nil {
			if d := cfg.GetWindow(); d != nil {
				opts = append(opts, sre.WithWindow(d.AsDuration()))
			}
			if v := cfg.GetRequest(); v > 0 {
				opts = append(opts, sre.WithRequest(v))
			}
			if v := cfg.GetBucket(); v > 0 {
				opts = append(opts, sre.WithBucket(int(v)))
			}
			if v := cfg.GetSuccess(); v > 0 {
				opts = append(opts, sre.WithSuccess(v))
			}
		}
		return sre.NewBreaker(opts...)
	})
}

func Server(opts ...Option) middleware.Middleware {
	opt := &options{
		group: group.NewGroup(func() circuitbreaker.CircuitBreaker {
			return sre.NewBreaker()
		}),
	}
	for _, o := range opts {
		o(opt)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			info, _ := transport.FromServerContext(ctx)
			key := "server"
			if info != nil {
				if op := info.Operation(); op != "" {
					key = info.Kind().String() + ":" + op
				} else {
					key = info.Kind().String()
				}
			}
			breaker := opt.group.Get(key)
			if err := breaker.Allow(); err != nil {
				breaker.MarkFailed()
				return nil, ErrNotAllowed
			}
			defer func() {
				if recovered := recover(); recovered != nil {
					breaker.MarkFailed()
					panic(recovered)
				}
			}()

			reply, err := handler(ctx, req)
			if isServerFailure(err) {
				breaker.MarkFailed()
			} else {
				breaker.MarkSuccess()
			}
			return reply, err
		}
	}
}

func isServerFailure(err error) bool {
	if err == nil {
		return false
	}
	if stderrors.Is(err, context.Canceled) || stderrors.Is(err, context.DeadlineExceeded) {
		return true
	}
	return errors.IsInternalServer(err) ||
		errors.IsServiceUnavailable(err) ||
		errors.IsGatewayTimeout(err)
}

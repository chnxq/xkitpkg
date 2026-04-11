package fluentd

import (
	"github.com/chnxq/xkitmod/log"

	conf "github.com/chnxq/xkitpkg/conf/v1"
	"github.com/chnxq/xkitpkg/logger"
)

func init() {
	_ = logger.RegisterFactory(logger.Fluentd, func(cfg *conf.Logger) (log.Logger, error) {
		return NewLogger(cfg)
	})
}

// NewLogger 创建一个新的日志记录器 - Fluent
func NewLogger(cfg *conf.Logger) (log.Logger, error) {
	if cfg == nil || cfg.Fluentd == nil {
		return nil, nil
	}

	wrapped, err := NewFluentLogger(cfg.Fluentd.Endpoint)
	if err != nil {
		return nil, err
	}
	return wrapped, nil
}

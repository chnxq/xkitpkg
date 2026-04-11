package tencent

import (
	"github.com/chnxq/xkitmod/log"
	conf "github.com/chnxq/xkitpkg/conf/v1"
	"github.com/chnxq/xkitpkg/logger"
)

func init() {
	_ = logger.RegisterFactory(logger.Tencent, func(cfg *conf.Logger) (log.Logger, error) {
		return NewLogger(cfg)
	})
}

// NewLogger 创建一个新的日志记录器 - Tencent
func NewLogger(cfg *conf.Logger) (log.Logger, error) {
	if cfg == nil || cfg.Tencent == nil {
		return nil, nil
	}

	wrapped, err := NewTencentLogger(
		WithTopicID(cfg.Tencent.TopicId),
		WithEndpoint(cfg.Tencent.Endpoint),
		WithAccessKey(cfg.Tencent.AccessKey),
		WithAccessSecret(cfg.Tencent.AccessSecret),
	)
	if err != nil {
		return nil, err
	}
	return wrapped, nil
}

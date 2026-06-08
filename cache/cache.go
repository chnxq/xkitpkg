package cache

import (
	"errors"

	"github.com/chnxq/xkitmod/log"
	conf "github.com/chnxq/xkitpkg/conf/v1"
)

func NewCache(cfg *conf.Data) (AdapterCache, error) {
	if cfg == nil {
		return nil, errors.New("cache config is nil")
	}

	var cache AdapterCache
	if cfg.GetRedis() == nil || cfg.GetRedis().GetAddr() == "" {
		cache = NewMemory()
	} else {
		cache = NewRedis(cfg, log.NewHelper(log.GetLogger()))
	}
	if cache == nil {
		return nil, errors.New("cache is nil")
	}
	if err := cache.Connect(); err != nil {
		return nil, err
	}
	return cache, nil
}

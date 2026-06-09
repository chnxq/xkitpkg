package cache

import (
	"errors"
	"strings"

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
		log.Info("Memory cache init OK.")
	} else {
		addrs := strings.Split(cfg.GetRedis().GetAddr(), ",")
		if len(addrs) <= 1 {
			cache = NewRedis(cfg, log.NewHelper(log.GetLogger()))
			if cache != nil {
				log.Info("Redis cache init OK.")
			}
		} else {
			cache = NewClusterRedis(cfg, log.NewHelper(log.GetLogger()))
			if cache != nil {
				log.Info("Redis-Cluster cache init OK.")
			}
		}
	}
	if cache == nil {
		return nil, errors.New("cache is nil")
	}
	if err := cache.Connect(); err != nil {
		return nil, err
	}
	return cache, nil
}

package cache

import (
	"context"
	"strings"
	"time"

	"github.com/chnxq/xkitmod/log"
	conf "github.com/chnxq/xkitpkg/conf/v1"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

// ClusterRedis cache implement
type ClusterRedis struct {
	client *redis.ClusterClient
}

// NewClusterRedis redis集群模式
func NewClusterRedis(cfg *conf.Data, logger *log.Helper) AdapterCache {
	if cfg == nil || cfg.GetRedis() == nil {
		return nil
	}

	addr := cfg.GetRedis().GetAddr()
	addrs := strings.Split(addr, ",")
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        addrs,
		Password:     cfg.GetRedis().GetPassword(),
		DialTimeout:  cfg.GetRedis().GetDialTimeout().AsDuration(),
		WriteTimeout: cfg.GetRedis().GetWriteTimeout().AsDuration(),
		ReadTimeout:  cfg.GetRedis().GetReadTimeout().AsDuration(),
	})
	if rdb == nil {
		if logger != nil {
			logger.Error("failed opening connection to redis cluster")
		}
		return nil
	}
	if cfg.GetRedis().GetEnableTracing() {
		if err := redisotel.InstrumentTracing(rdb); err != nil {
			if logger != nil {
				logger.Errorf("failed open tracing: %s", err.Error())
			}
			_ = rdb.Close()
			return nil
		}
	}

	if cfg.GetRedis().GetEnableMetrics() {
		if err := redisotel.InstrumentMetrics(rdb); err != nil {
			if logger != nil {
				logger.Errorf("failed open metrics: %s", err.Error())
			}
			_ = rdb.Close()
			return nil
		}
	}

	return &ClusterRedis{
		client: rdb,
	}
}

func (s *ClusterRedis) Connect() error {
	return s.client.ForEachShard(context.TODO(), func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})
}

func (s *ClusterRedis) DisConnect() error {
	return s.client.Close()
}

func (s *ClusterRedis) Get(key string) (string, error) {
	return s.client.Get(context.TODO(), key).Result()
}

func (s *ClusterRedis) Set(key, value string, expire time.Duration) error {
	return s.client.Set(context.TODO(), key, value, expire).Err()
}

func (s *ClusterRedis) Del(key string) error {
	return s.client.Del(context.TODO(), key).Err()
}

func (s *ClusterRedis) Expire(key string, dur time.Duration) error {
	return s.client.Expire(context.TODO(), key, dur).Err()
}

func (s *ClusterRedis) Exists(key string) bool {
	result, err := s.client.Exists(context.TODO(), key).Result()
	if err != nil {
		return false
	}
	return result != 0
}

func (s *ClusterRedis) HGetAll(key string) (map[string]string, error) {
	return s.client.HGetAll(context.TODO(), key).Result()
}

func (s *ClusterRedis) HGet(key, field string) (string, error) {
	return s.client.HGet(context.TODO(), key, field).Result()
}

func (s *ClusterRedis) HSet(key, field, value string) error {
	return s.client.HSet(context.TODO(), key, field, value).Err()
}

func (s *ClusterRedis) HDel(key, field string) error {
	return s.client.HDel(context.TODO(), key, field).Err()
}

func (s *ClusterRedis) HExists(key, field string) (bool, error) {
	return s.client.HExists(context.TODO(), key, field).Result()
}

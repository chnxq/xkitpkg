package cache

import (
	"context"
	"errors"
	"time"

	"github.com/chnxq/xkitmod/log"
	conf "github.com/chnxq/xkitpkg/conf/v1"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(cfg *conf.Data, logger *log.Helper) AdapterCache {
	if cfg == nil || cfg.GetRedis() == nil {
		return nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.GetRedis().GetAddr(),
		Password:     cfg.GetRedis().GetPassword(),
		DB:           int(cfg.GetRedis().GetDb()),
		DialTimeout:  cfg.GetRedis().GetDialTimeout().AsDuration(),
		WriteTimeout: cfg.GetRedis().GetWriteTimeout().AsDuration(),
		ReadTimeout:  cfg.GetRedis().GetReadTimeout().AsDuration(),
	})
	if rdb == nil {
		if logger != nil {
			logger.Error("failed opening connection to redis")
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

	return &Redis{client: rdb}
}

func (s *Redis) Connect() error {
	if s == nil || s.client == nil {
		return errors.New("redis client is nil")
	}
	_, err := s.client.Ping(context.TODO()).Result()
	return err
}

func (s *Redis) DisConnect() error {
	if s != nil && s.client != nil {
		return s.client.Close()
	}
	return errors.New("redis client is nil")
}

func (s *Redis) Get(key string) (string, error) {
	return s.client.Get(context.TODO(), key).Result()
}

func (s *Redis) Set(key, value string, expire time.Duration) error {
	return s.client.Set(context.TODO(), key, value, expire).Err()
}

func (s *Redis) Del(key string) error {
	return s.client.Del(context.TODO(), key).Err()
}

func (s *Redis) Expire(key string, dur time.Duration) error {
	return s.client.Expire(context.TODO(), key, dur).Err()
}

func (s *Redis) Exists(key string) bool {
	result, err := s.client.Exists(context.TODO(), key).Result()
	if err != nil {
		return false
	}
	return result != 0
}

func (s *Redis) HGetAll(key string) (map[string]string, error) {
	return s.client.HGetAll(context.TODO(), key).Result()
}

func (s *Redis) HGet(key, field string) (string, error) {
	return s.client.HGet(context.TODO(), key, field).Result()
}

func (s *Redis) HSet(key, field, value string) error {
	return s.client.HSet(context.TODO(), key, field, value).Err()
}

func (s *Redis) HDel(key, field string) error {
	return s.client.HDel(context.TODO(), key, field).Err()
}

func (s *Redis) HExists(key, field string) (bool, error) {
	return s.client.HExists(context.TODO(), key, field).Result()
}

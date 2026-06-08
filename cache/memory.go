package cache

import (
	"errors"
	"sync"
	"time"

	"github.com/chnxq/xkitmod/log"
)

type strItem struct {
	Value   string
	Expired time.Time
}

type mapItem struct {
	Value   map[string]string
	Expired time.Time
}

type Memory struct {
	strItems map[string]*strItem
	strMutex sync.RWMutex
	mapItems map[string]*mapItem
	mapMutex sync.RWMutex
}

func NewMemory() AdapterCache {
	log.Debug("Memory cache init.")
	return &Memory{
		strItems: make(map[string]*strItem),
		mapItems: make(map[string]*mapItem),
	}
}

func (s *Memory) Connect() error {
	if s.strItems == nil || s.mapItems == nil {
		return errors.New("memory connect fail")
	}
	return nil
}

func (s *Memory) DisConnect() error {
	s.strMutex.Lock()
	defer s.strMutex.Unlock()
	s.mapMutex.Lock()
	defer s.mapMutex.Unlock()

	s.strItems = nil
	s.mapItems = nil
	return nil
}

func (s *Memory) Get(key string) (string, error) {
	s.strMutex.RLock()
	item, ok := s.strItems[key]
	s.strMutex.RUnlock()

	if !ok {
		return "", errors.New("key not found")
	}
	if isExpired(item.Expired) {
		s.strMutex.Lock()
		delete(s.strItems, key)
		s.strMutex.Unlock()
		return "", errors.New("key expired")
	}
	return item.Value, nil
}

func (s *Memory) Set(key, value string, expire time.Duration) error {
	s.strMutex.Lock()
	defer s.strMutex.Unlock()

	s.strItems[key] = &strItem{
		Value:   value,
		Expired: expiresAt(expire),
	}
	return nil
}

func (s *Memory) Del(key string) error {
	s.strMutex.Lock()
	delete(s.strItems, key)
	s.strMutex.Unlock()

	s.mapMutex.Lock()
	delete(s.mapItems, key)
	s.mapMutex.Unlock()
	return nil
}

func (s *Memory) Expire(key string, dur time.Duration) error {
	expiredAt := expiresAt(dur)

	s.strMutex.Lock()
	if item, ok := s.strItems[key]; ok {
		item.Expired = expiredAt
		s.strItems[key] = item
		s.strMutex.Unlock()
		return nil
	}
	s.strMutex.Unlock()

	s.mapMutex.Lock()
	defer s.mapMutex.Unlock()

	item, ok := s.mapItems[key]
	if !ok {
		return errors.New("key not found")
	}
	item.Expired = expiredAt
	s.mapItems[key] = item
	return nil
}

func (s *Memory) Exists(key string) bool {
	s.strMutex.RLock()
	item, ok := s.strItems[key]
	s.strMutex.RUnlock()
	if ok {
		if isExpired(item.Expired) {
			s.strMutex.Lock()
			delete(s.strItems, key)
			s.strMutex.Unlock()
			return false
		}
		return true
	}

	s.mapMutex.RLock()
	hashItem, ok := s.mapItems[key]
	s.mapMutex.RUnlock()
	if !ok {
		return false
	}
	if isExpired(hashItem.Expired) {
		s.mapMutex.Lock()
		delete(s.mapItems, key)
		s.mapMutex.Unlock()
		return false
	}
	return true
}

func (s *Memory) HGetAll(key string) (map[string]string, error) {
	s.mapMutex.RLock()
	item, ok := s.mapItems[key]
	s.mapMutex.RUnlock()

	if !ok {
		return nil, errors.New("key not found")
	}
	if isExpired(item.Expired) {
		s.mapMutex.Lock()
		delete(s.mapItems, key)
		s.mapMutex.Unlock()
		return nil, errors.New("key expired")
	}
	return cloneMap(item.Value), nil
}

func (s *Memory) HGet(key, field string) (string, error) {
	s.mapMutex.RLock()
	item, ok := s.mapItems[key]
	s.mapMutex.RUnlock()

	if !ok {
		return "", errors.New("key not found")
	}
	if isExpired(item.Expired) {
		s.mapMutex.Lock()
		delete(s.mapItems, key)
		s.mapMutex.Unlock()
		return "", errors.New("key expired")
	}

	itemValue, ok := item.Value[field]
	if !ok {
		return "", errors.New("field not found")
	}
	return itemValue, nil
}

func (s *Memory) HSet(key, field, value string) error {
	s.mapMutex.Lock()
	defer s.mapMutex.Unlock()

	item, ok := s.mapItems[key]
	if !ok {
		item = &mapItem{
			Value: make(map[string]string),
		}
	}
	item.Value[field] = value
	s.mapItems[key] = item
	return nil
}

func (s *Memory) HDel(key, field string) error {
	s.mapMutex.Lock()
	defer s.mapMutex.Unlock()

	item, ok := s.mapItems[key]
	if !ok {
		return errors.New("key not found")
	}
	if isExpired(item.Expired) {
		delete(s.mapItems, key)
		return errors.New("key expired")
	}

	delete(item.Value, field)
	s.mapItems[key] = item
	return nil
}

func (s *Memory) HExists(key, field string) (bool, error) {
	s.mapMutex.RLock()
	item, ok := s.mapItems[key]
	s.mapMutex.RUnlock()

	if !ok {
		return false, nil
	}
	if isExpired(item.Expired) {
		s.mapMutex.Lock()
		delete(s.mapItems, key)
		s.mapMutex.Unlock()
		return false, nil
	}

	_, ok = item.Value[field]
	return ok, nil
}

func isExpired(expiredAt time.Time) bool {
	return !expiredAt.IsZero() && time.Now().After(expiredAt)
}

func expiresAt(dur time.Duration) time.Time {
	if dur <= 0 {
		return time.Time{}
	}
	return time.Now().Add(dur)
}

func cloneMap(source map[string]string) map[string]string {
	if source == nil {
		return nil
	}
	result := make(map[string]string, len(source))
	for key, value := range source {
		result[key] = value
	}
	return result
}

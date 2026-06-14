package config

import (
	"reflect"
	"sync"

	"github.com/chnxq/xkitpkg/conf"
	"google.golang.org/protobuf/proto"
)

var (
	muBC         sync.RWMutex
	initOnce     sync.Once
	configList   []proto.Message
	configSet    map[uintptr]struct{}
	commonConfig *conf.ServerConfig
)

func GetServerConfig() *conf.ServerConfig {
	initServerConfig()
	muBC.RLock()
	defer muBC.RUnlock()
	return commonConfig
}

func RegisterConfig(c proto.Message) {
	if c == nil {
		return
	}
	initServerConfig()

	muBC.Lock()
	defer muBC.Unlock()
	addConfigLocked(c)
}

func initServerConfig() {
	initOnce.Do(func() {
		muBC.Lock()
		defer muBC.Unlock()

		configList = make([]proto.Message, 0)
		configSet = make(map[uintptr]struct{})

		if commonConfig == nil {
			commonConfig = &conf.ServerConfig{}
		}
		addConfigLocked(commonConfig)

		if commonConfig.Server == nil {
			commonConfig.Server = &conf.Server{}
		}
		addConfigLocked(commonConfig.Server)

		if commonConfig.Client == nil {
			commonConfig.Client = &conf.Client{}
		}
		addConfigLocked(commonConfig.Client)

		if commonConfig.Data == nil {
			commonConfig.Data = &conf.Data{}
		}
		addConfigLocked(commonConfig.Data)

		if commonConfig.Trace == nil {
			commonConfig.Trace = &conf.Tracer{}
		}
		addConfigLocked(commonConfig.Trace)

		if commonConfig.Logger == nil {
			commonConfig.Logger = &conf.Logger{}
		}
		addConfigLocked(commonConfig.Logger)

		if commonConfig.Registry == nil {
			commonConfig.Registry = &conf.Registry{}
		}
		addConfigLocked(commonConfig.Registry)

		if commonConfig.Oss == nil {
			commonConfig.Oss = &conf.OSS{}
		}
		addConfigLocked(commonConfig.Oss)

		if commonConfig.Notify == nil {
			commonConfig.Notify = &conf.Notification{}
		}
		addConfigLocked(commonConfig.Notify)
	})
}

func addConfigLocked(c proto.Message) {
	if c == nil {
		return
	}
	v := reflect.ValueOf(c)
	if !v.IsValid() || v.Kind() != reflect.Ptr || v.IsNil() {
		return
	}
	addr := v.Pointer()
	if _, exists := configSet[addr]; exists {
		return
	}
	configList = append(configList, c)
	configSet[addr] = struct{}{}
}

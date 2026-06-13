package config

import (
	"sync"

	"github.com/chnxq/xkitpkg/conf"
)

// RuntimeConfigApplier applies a refreshed server config for a specific top-level key.
// It returns true when the change has been applied without restart.
type RuntimeConfigApplier func(cfg *conf.ServerConfig) (bool, error)

var runtimeConfigAppliers sync.Map

// RegisterRuntimeConfigApplier registers a runtime applier for a top-level config key.
func RegisterRuntimeConfigApplier(key string, applier RuntimeConfigApplier) {
	if key == "" || applier == nil {
		return
	}
	runtimeConfigAppliers.Store(key, applier)
}

func applyRuntimeConfigChange(key string, cfg *conf.ServerConfig) (bool, error) {
	applier, ok := runtimeConfigAppliers.Load(key)
	if !ok {
		return false, nil
	}
	return applier.(RuntimeConfigApplier)(cfg)
}

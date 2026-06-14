package config

import (
	"errors"

	baseconfig "github.com/chnxq/xkitmod/config"
	"github.com/chnxq/xkitmod/log"
)

var watchedConfigKeys = []string{
	"server",
	"client",
	"data",
	"trace",
	"logger",
	"registry",
	"config",
	"oss",
	"notify",
	"authn",
	"authz",
	"script",
}

func registerConfigRefreshWatchers(cfg baseconfig.Config) error {
	for _, key := range watchedConfigKeys {
		key := key
		if err := cfg.Watch(key, func(_ string, _ baseconfig.Value) {
			log.Debugf("server config observer triggered: key=%s", key)
			if err := scanConfigs(cfg); err != nil {
				log.Errorf("rescan server config failed after %s update: %v", key, err)
				return
			}
			log.Debugf("server config rescanned successfully after update: key=%s", key)
			applied, err := applyRuntimeConfigChange(key, GetServerConfig())
			if err != nil {
				log.Errorf("apply runtime config failed for key=%s: %v", key, err)
				log.Warnf("config key=%s changed, restart required for full effect", key)
				return
			}
			if applied {
				log.Infof("config key=%s applied at runtime", key)
				return
			}
			log.Warnf("config key=%s changed, restart required for full effect", key)
		}); err != nil {
			if errors.Is(err, baseconfig.ErrNotFound) {
				log.Debugf("server config observer skipped missing key: %s", key)
				continue
			}
			return err
		}
		log.Debugf("server config observer registered: key=%s", key)
	}
	return nil
}

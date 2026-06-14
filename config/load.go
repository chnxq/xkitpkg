package config

import (
	"os"
	"path/filepath"

	"github.com/chnxq/xkitmod/log"
	"github.com/chnxq/xkitpkg/conf"

	baseconfig "github.com/chnxq/xkitmod/config"
	filekit "github.com/chnxq/xkitmod/config/file"
)

func LoadServerConfig(configPath string) error {
	cfg, err := CheckConfigProvider(configPath)
	if err != nil {
		log.Errorf("check remote config provider failed: %v\n", err)
		return err
	}

	if err = cfg.Load(); err != nil {
		log.Errorf("load config failed: %v\n", err)
		return err
	}

	if err = scanConfigs(cfg); err != nil {
		log.Errorf("scan config failed: %v\n", err)
		return err
	}
	if err = registerConfigRefreshWatchers(cfg); err != nil {
		log.Errorf("register config refresh watchers failed: %v\n", err)
		return err
	}

	return nil
}

func CheckConfigProvider(configPath string) (baseconfig.Config, error) {
	var cfg baseconfig.Config

	remoteConfigPath := filepath.Join(configPath, "remote_config.yaml")
	haveRemoteConfig := pathExists(remoteConfigPath)
	if !haveRemoteConfig {
		remoteConfigPath = filepath.Join(configPath, "config.yaml")
		haveRemoteConfig = pathExists(remoteConfigPath)
	}

	if !haveRemoteConfig {
		return baseconfig.New(baseconfig.WithSource(NewFileConfigSource(configPath))), nil
	}

	cfgRemote := baseconfig.New(baseconfig.WithSource(NewFileConfigSource(remoteConfigPath)))
	defer func() {
		if err := cfgRemote.Close(); err != nil {
			panic(err)
		}
	}()

	if err := cfgRemote.Load(); err != nil {
		log.Errorf("load remote config failed: %v\n", err)
		return nil, err
	}
	if err := scanConfigs(cfgRemote); err != nil {
		log.Errorf("scan remote config failed: %v\n", err)
		return nil, err
	}
	log.Infof("Have remote config: %v", configList)

	rc := GetServerConfig().GetConfig()
	if rc == nil {
		return baseconfig.New(baseconfig.WithSource(NewFileConfigSource(configPath))), nil
	}

	rcs, err := NewProvider(rc)
	if err != nil {
		log.Errorf("Create remote config provider failed: %v\n", err)
		return baseconfig.New(baseconfig.WithSource(NewFileConfigSource(configPath))), nil
	}

	cfg = baseconfig.New(
		baseconfig.WithSource(
			rcs,
			NewFileConfigSource(configPath),
		),
	)
	log.Infof("Create remote config provider success: %v", rc)
	return cfg, nil
}

func NewFileConfigSource(filePath string) baseconfig.Source {
	return filekit.NewSource(filePath)
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return !os.IsNotExist(err)
}

func scanConfigs(cfg baseconfig.Config) error {
	initServerConfig()
	muBC.Lock()
	defer muBC.Unlock()
	return scanConfigsLocked(cfg)
}

func scanConfigsLocked(cfg baseconfig.Config) error {
	for _, c := range configList {
		if err := cfg.Scan(c); err != nil {
			return err
		}
	}
	return nil
}

var _ *conf.ServerConfig

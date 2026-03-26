package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"

	"google.golang.org/protobuf/proto"

	"github.com/chnxq/XGoKit/config"
	fileKit "github.com/chnxq/XGoKit/config/file"
	"github.com/chnxq/xkitpkg/conf/v1"
)

var (
	muBC         sync.RWMutex
	initOnce     sync.Once
	configList   []proto.Message
	configSet    map[uintptr]struct{}
	commonConfig *conf.ServerConfig
)

// LoadServerConfig 加载程序引导配置
func LoadServerConfig(configPath string) error {
	cfg, err := CheckConfigProvider(configPath)
	if err != nil {
		return err
	}

	if err = cfg.Load(); err != nil {
		return err
	}

	initServerConfig()

	if err = scanConfigs(cfg); err != nil {
		return err
	}

	return nil
}

func CheckConfigProvider(configPath string) (config.Config, error) {
	var err error
	var cfg config.Config

	remoteConfigPath := filepath.Join(configPath, "remote_config.yaml")
	if pathExists(remoteConfigPath) {
		cfg = config.New(
			config.WithSource(
				NewFileConfigSource(remoteConfigPath),
			),
		)
	} else {
		cfg = config.New(
			config.WithSource(
				NewFileConfigSource(configPath),
			),
		)
	}
	defer func(cfg config.Config) {
		if err := cfg.Close(); err != nil {
			panic(err)
		}
		fmt.Println("check config source OK")
	}(cfg)

	if err = cfg.Load(); err != nil {
		fmt.Printf("Load config failed: %v", err)
		return nil, err
	}

	if err = scanConfigs(cfg); err != nil {
		fmt.Printf("Scan config failed: %v", err)
		return nil, err
	}

	rc := GetServerConfig().Config
	if rc != nil {
		rcs, err := NewProvider(rc)
		if err != nil {
			fmt.Printf("New config provider failed: %v", err)
			return nil, err
		}
		cfg = config.New(
			config.WithSource(
				rcs,
				NewFileConfigSource(configPath),
			),
		)
	}
	return cfg, nil
}

func GetServerConfig() *conf.ServerConfig {
	initServerConfig()
	muBC.RLock()
	defer muBC.RUnlock()
	return commonConfig
}

// NewFileConfigSource 创建一个本地文件配置源
func NewFileConfigSource(filePath string) config.Source {
	return fileKit.NewSource(filePath)
}

// RegisterConfig 注册配置
// 传入值应为指针类型，例如 &conf.SomeConfig{}
func RegisterConfig(c proto.Message) {
	if c == nil {
		return
	}
	initServerConfig()

	muBC.Lock()
	defer muBC.Unlock()
	addConfigLocked(c)
}

// initServerConfig 初始化服务器配置（仅执行一次）
func initServerConfig() {
	initOnce.Do(func() {
		muBC.Lock()
		defer muBC.Unlock()

		// 初始化集合与列表
		configList = make([]proto.Message, 0)
		configSet = make(map[uintptr]struct{})

		if commonConfig == nil {
			commonConfig = &conf.ServerConfig{}
		}

		// 按需添加根与子配置，使用去重函数
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

// addConfigLocked 假定已持有 muBC 锁，添加时会去重并确保参数为指针
func addConfigLocked(c proto.Message) {
	if c == nil {
		return
	}
	v := reflect.ValueOf(c)
	if !v.IsValid() || v.Kind() != reflect.Ptr || v.IsNil() {
		// 只接受非 nil 的指针类型
		return
	}
	addr := v.Pointer()
	if _, exists := configSet[addr]; exists {
		return
	}
	configList = append(configList, c)
	configSet[addr] = struct{}{}
}

func scanConfigs(cfg config.Config) error {
	initServerConfig()

	for _, c := range configList {
		if err := cfg.Scan(c); err != nil {
			return err
		}
	}
	return nil
}

// pathExists 判断路径是否存在
func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

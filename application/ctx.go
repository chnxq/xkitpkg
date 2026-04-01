package application

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"google.golang.org/protobuf/proto" // 用于 Protobuf 消息序列化和反序列化的工具包

	kLog "github.com/chnxq/XGoKit/log"
	kRegistry "github.com/chnxq/XGoKit/registry"
	"github.com/chnxq/xkitpkg/conf/v1"        // 应用配置结构定义
	bConfig "github.com/chnxq/xkitpkg/config" // 配置管理工具
)

// AppCtx 引导上下文
type AppCtx struct {
	config  *conf.ServerConfig // 引导配置
	appInfo *conf.AppInfo      // 应用信息

	logger    kLog.Logger         // 日志记录器
	registrar kRegistry.Registrar // 服务注册器

	customConfig sync.Map // 自定义配置项
	values       sync.Map // 自定义值存储

	rootCtx context.Context    // 应用级根上下文（可用于优雅关闭）
	cancel  context.CancelFunc // 取消函数
}

// NewAppCtx 创建带 cancel 的应用级 AppCtx（传 nil 使用 Background）
func NewAppCtx(parent context.Context, ai *conf.AppInfo) *AppCtx {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)

	c := &AppCtx{
		appInfo: &conf.AppInfo{},
	}
	// 初始化默认信息
	AdjustAppInfo(c.appInfo)

	c.copyAppInfo(ai)

	// 其余初始化例如 RootCtx/Cancel/Logger 可在这里设置
	_ = cancel // 保留 cancel 给调用者或另行设置
	_ = ctx
	return c
}

func NewContextWithParam(parent context.Context, ai *conf.AppInfo, cfg *conf.ServerConfig, log kLog.Logger) *AppCtx {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)

	c := &AppCtx{
		appInfo: &conf.AppInfo{},
		config:  cfg,
		logger:  log,
	}
	// 初始化默认信息
	AdjustAppInfo(c.appInfo)

	c.copyAppInfo(ai)

	// 其余初始化例如 RootCtx/Cancel/Logger 可在这里设置
	_ = cancel // 保留 cancel 给调用者或另行设置
	_ = ctx
	return c
}

// AppCtx 返回应用级根 context（保证非 nil）
func (c *AppCtx) AppCtx() context.Context {
	if c == nil || c.rootCtx == nil {
		return context.Background()
	}
	return c.rootCtx
}

// CancelContext 触发取消（幂等）
func (c *AppCtx) CancelContext() {
	if c == nil {
		return
	}
	if c.cancel != nil {
		c.cancel()
	}
}

func (c *AppCtx) NewLoggerHelper(moduleName string) *kLog.Helper {
	return kLog.NewHelper(kLog.With(c.logger, "module", moduleName))
}

func (c *AppCtx) GetLogger() kLog.Logger {
	return c.logger
}

// GetConfig 返回当前的 *conf.ServerConfig（并发安全）
func (c *AppCtx) GetConfig() *conf.ServerConfig {
	if c.config == nil {
		return nil
	}
	if clone := proto.Clone(c.config); clone != nil {
		if b, ok := clone.(*conf.ServerConfig); ok {
			return b
		}
	}
	return nil
}

func (c *AppCtx) GetAppInfo() *conf.AppInfo {
	if c.appInfo == nil {
		return nil
	}
	if clone := proto.Clone(c.appInfo); clone != nil {
		if a, ok := clone.(*conf.AppInfo); ok {
			return a
		}
	}
	return nil
}

// setAppInfo 用受控方式替换整个 appInfo（可选）
func (c *AppCtx) setAppInfo(src *conf.AppInfo) {
	if c == nil || src == nil {
		return
	}
	AdjustAppInfo(src)

	c.appInfo = &conf.AppInfo{
		Name:       src.Name,
		Version:    src.Version,
		AppId:      src.AppId,
		Project:    src.Project,
		InstanceId: src.InstanceId,
		Hostname:   src.Hostname,
		StartTime:  src.StartTime,
		Metadata:   cloneMetadata(src.Metadata),
	}
}

// copyAppInfo 复制应用信息
func (c *AppCtx) copyAppInfo(ai *conf.AppInfo) {
	if ai == nil {
		return
	}

	// 先修正输入，避免未初始化字段
	AdjustAppInfo(ai)

	if ai.Name != "" {
		c.appInfo.Name = ai.Name
	}
	if ai.Project != "" {
		c.appInfo.Project = ai.Project
	}
	if ai.AppId != "" {
		c.appInfo.AppId = ai.AppId
	}
	if ai.Version != "" {
		c.appInfo.Version = ai.Version
	}
	if ai.InstanceId != "" {
		c.appInfo.InstanceId = ai.InstanceId
	}
	if ai.Metadata != nil {
		c.appInfo.Metadata = ai.Metadata
	}
}

func (c *AppCtx) PrintAppInfo() {
	ai := c.GetAppInfo()
	if ai == nil {
		return
	}
	ts := time.Now().Format(time.RFC3339)
	host, _ := os.Hostname()
	pid := os.Getpid()

	if os.Getenv("APPINFO_FORMAT") == "json" {
		out := map[string]interface{}{
			"timestamp":   ts,
			"host":        host,
			"pid":         pid,
			"name":        ai.Name,
			"version":     ai.Version,
			"app_id":      ai.AppId,
			"instance_id": ai.InstanceId,
			"metadata":    ai.Metadata,
		}
		if b, err := json.Marshal(out); err == nil {
			fmt.Println(string(b))
		} else {
			fmt.Printf("Application info marshal error: %v\n", err)
		}
		return
	}

	fmt.Printf("[%s] %s (pid:%d@%s)\n", ts, ai.Name, pid, host)
	fmt.Printf("  Version: %s\n", ai.Version)
	fmt.Printf("  AppId: %s\n", ai.AppId)
	fmt.Printf("  InstanceId: %s\n", ai.InstanceId)
	if len(ai.Metadata) > 0 {
		fmt.Println("  Metadata:")
		keys := make([]string, 0, len(ai.Metadata))
		for k := range ai.Metadata {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("    %s=%s\n", k, ai.Metadata[k])
		}
	}
}

func (c *AppCtx) GetRegistrar() kRegistry.Registrar {
	return c.registrar
}

// RegisterCustomConfig 注册自定义配置
func (c *AppCtx) RegisterCustomConfig(key string, cfg proto.Message) {
	if key == "" || cfg == nil {
		return
	}

	if _, ok := c.customConfig.Load(key); ok {
		return
	}

	c.customConfig.Store(key, cfg)

	bConfig.RegisterConfig(cfg)
}

// SetCustomConfig 存入自定义配置
func (c *AppCtx) SetCustomConfig(key string, cfg proto.Message) {
	if key == "" || cfg == nil {
		return
	}

	c.customConfig.Store(key, cfg)
}

// GetCustomConfig 获取自定义配置（原始类型）
func (c *AppCtx) GetCustomConfig(key string) (any, bool) {
	return c.customConfig.Load(key)
}

// DeleteCustomConfig 删除自定义配置
func (c *AppCtx) DeleteCustomConfig(key string) {
	c.customConfig.Delete(key)
}

// RangeCustomConfig 遍历自定义配置，回调返回 false 可停止遍历
func (c *AppCtx) RangeCustomConfig(fn func(key string, val any) bool) {
	c.customConfig.Range(func(k, v any) bool {
		ks, _ := k.(string)
		return fn(ks, v)
	})
}

// SetValue 将任意值放入通用存储
func (c *AppCtx) SetValue(key string, val interface{}) {
	c.values.Store(key, val)
}

// GetValue 从通用存储读取值
func (c *AppCtx) GetValue(key string) (interface{}, bool) {
	return c.values.Load(key)
}

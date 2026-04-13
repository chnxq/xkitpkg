package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"google.golang.org/protobuf/proto" // 用于 Protobuf 消息序列化和反序列化的工具包

	"github.com/chnxq/x-utils/timeutil"
	"github.com/chnxq/x-utils/trans"

	"github.com/chnxq/xkitmod/log"
	"github.com/chnxq/xkitmod/registry"

	"github.com/chnxq/xkitpkg/conf/v1" // 应用配置结构定义
)

// AppCtx 应用上下文
type AppCtx struct {
	config  *conf.ServerConfig // 应用程序配置
	appInfo *conf.AppInfo      // 应用信息

	logger    log.Logger         // 日志记录器
	registrar registry.Registrar // 服务注册器
	//tracer    kTracer.Tracer      // 服务注册器

	rootCtx context.Context    // 应用级根上下文（可用于优雅关闭）
	cancel  context.CancelFunc // 取消函数
}

// NewAppCtx 创建带 cancel 的应用级 AppCtx（传 nil 使用 Background）
func NewAppCtx(parent context.Context, ai *conf.AppInfo, cfg *conf.ServerConfig, log log.Logger, reg registry.Registrar) *AppCtx {
	// 全局的应用级根上下文（可用于优雅关闭）
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)

	c := &AppCtx{
		appInfo: &conf.AppInfo{
			Project:   ai.Project,
			AppId:     ai.AppId,
			Name:      ai.Name,
			Version:   ai.Version,
			Metadata:  ai.Metadata,
			BuildId:   ai.BuildId,
			GitCommit: ai.GitCommit,
			StartTime: timeutil.TimeToTimestamppb(trans.Ptr(time.Now())),
		},
		config:    cfg,
		logger:    log,
		registrar: reg,

		rootCtx: ctx,
		cancel:  cancel,
	}
	return c
}

// AppCtx 返回应用级根 context（保证非 nil）
func (c *AppCtx) AppContext() context.Context {
	if c == nil || c.rootCtx == nil {
		return context.Background()
	}
	return c.rootCtx
}

// CancelContext 触发取消
func (c *AppCtx) CancelContext() {
	if c == nil {
		return
	}
	if c.cancel != nil {
		c.cancel()
	}
}

func (c *AppCtx) NewLoggerHelper(moduleName string) *log.Helper {
	return log.NewHelper(log.With(c.logger, "module", moduleName))
}

func (c *AppCtx) GetLogger() log.Logger {
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

func (c *AppCtx) GetRegistrar() registry.Registrar {
	return c.registrar
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
			"timestamp": ts,
			"host":      host,
			"pid":       pid,
			"name":      ai.Name,
			"version":   ai.Version,
			"app_id":    ai.AppId,
			"metadata":  ai.Metadata,
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
	fmt.Printf("  GitCommit: %s\n", ai.GitCommit)
	fmt.Printf("  Build: %s\n", ai.BuildId)
}

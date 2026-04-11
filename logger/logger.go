package logger

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/chnxq/xkitpkg/conf/v1"
	"github.com/chnxq/xkitpkg/logger/log"
)

// NewLogger 动态创建日志实例
// 返回一个新的标准日志记录器实例
// 当无法创建自定义日志记录器时，此函数提供一个默认的标准日志记录器
// 作为备选方案，确保系统始终有一个可用的日志记录器  return NewStdLogger(), nil
func NewLogger(cfg *conf.Logger) (log.Logger, error) {
	if cfg == nil {
		return nil, nil
	}

	if cfg.GetType() == "" || cfg.GetType() == string(Std) {
		return NewStdLogger(), nil
	}

	// normalize to lower case for lookup
	typ := Type(strings.ToLower(cfg.GetType()))
	norm := Type(strings.ToLower(string(typ)))

	f, ok := GetFactory(norm)
	if !ok {
		// prepare available list for helpful error
		available := ListFactories()
		strs := make([]string, 0, len(available))
		for _, t := range available {
			strs = append(strs, string(t))
		}
		sort.Strings(strs)
		return nil, fmt.Errorf("unsupported logger type: %s; available: %v", typ, strs)
	}

	lg, err := f(cfg)
	if err != nil {
		return nil, fmt.Errorf("create logger %s: %w", typ, err)
	}
	return lg, nil
}

// NewLoggerProvider 创建一个新的日志记录器提供者
// 它会从 cfg 创建具体 logger（通过 NewLogger），并为 logger 附加一组标准字段（service.*, ts, caller, trace_id, span_id）。
func NewLoggerProvider(cfg *conf.Logger, appInfo *conf.AppInfo) log.Logger {
	var l log.Logger
	if cfg == nil || cfg.GetType() == "" {
		l = NewStdLogger()
	} else {
		// try to create logger by type via factory
		if lg, err := NewLogger(cfg); err == nil && lg != nil {
			l = lg
		} else {
			l = NewStdLogger()
		}
	}

	// build base fields - always include timestamp, caller, trace/span ids
	fields := []interface{}{}
	//fields := []interface{}{
	//	//"ts", log.DefaultTimestamp,
	//	"caller", log.DefaultCaller,
	//	"trace_id", tracing.TraceID(),
	//	"span_id", tracing.SpanID(),
	//}

	// attach service fields only if appInfo is provided
	if appInfo != nil {
		fields = append([]interface{}{
			"service.id", appInfo.GetAppId(),
			"service.instance", appInfo.GetInstanceId(),
			"service.version", appInfo.GetVersion(),
		}, fields...)
	}
	newlogger := log.With(l, fields...)
	log.SetLogger(newlogger)
	log.Infof("Logger %s created successfully", cfg.GetType())
	return newlogger
}

// NewStdLogger 创建一个新的日志记录器 - Kratos内置，控制台输出
func NewStdLogger() log.Logger {
	l := log.NewStdLogger(os.Stdout)
	return l
}

package zap

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	otelLog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/chnxq/XGoKit/log"

	"github.com/chnxq/xkitpkg/conf/v1"
)

var _ log.Logger = (*Logger)(nil)

type Logger struct {
	log    *zap.Logger
	msgKey string
}

type Option func(*Logger)

// WithMessageKey with message key.
func WithMessageKey(key string) Option {
	return func(l *Logger) {
		l.msgKey = key
	}
}

// newLoggerProvider creates a new logger provider with the OTLP gRPC exporter.
func newLoggerProvider(ctx context.Context, res *resource.Resource) (*otelLog.LoggerProvider, error) {
	exporter, err := otlploggrpc.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP log exporter: %w", err)
	}
	processor := otelLog.NewBatchProcessor(exporter)
	lp := otelLog.NewLoggerProvider(
		otelLog.WithProcessor(processor),
		otelLog.WithResource(res),
	)
	return lp, nil
}

func InitZapWithConfig(zapConfig *conf.Logger_Zap) *zap.Logger {
	var err error
	var lv zapcore.Level
	var coreArr []zapcore.Core

	//获取编码器配置
	encoderConfig := zap.NewProductionEncoderConfig() //NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	} //指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder //按级别显示不同颜色，不需要的话取值zapcore.CapitalLevelEncoder就可以了
	//encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder      //显示完整文件路径
	encoderConfig.EncodeCaller = nCallerEncoder //自定义Caller显示
	encoderConfig.CallerKey = "C"
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	//日志级别
	lv, err = zapcore.ParseLevel(zapConfig.GetLevel())
	if err != nil {
		panic(err)
	}

	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //info和debug级别,debug级别是最低的
		if zapConfig == nil {
			return lev >= zap.DebugLevel
		} else {
			return lev >= lv
		}
	})

	path := "./logfiles"
	if zapConfig != nil && len(zapConfig.GetLogFilePath()) > 0 {
		path = zapConfig.LogFilePath
	}

	//info文件writeSyncer
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path + "/info.log", //日志文件存放目录，如果文件夹不存在会自动创建
		LocalTime:  true,
		MaxSize:    int(zapConfig.GetMaxSize()),    //文件大小限制,单位MB
		MaxBackups: int(zapConfig.GetMaxBackups()), //最大保留日志文件数量
		MaxAge:     int(zapConfig.GetMaxAge()),     //日志文件保留天数
		Compress:   true,                           //是否压缩处理
	})

	var infoFileCore zapcore.Core
	if zapConfig == nil || zapConfig.LogToConsole {
		infoFileCore = zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), lowPriority)
	} else {
		infoFileCore = zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer), lowPriority)
	}

	if zapConfig.ExportToOtel {
		res, err := resource.New(
			context.Background(),
			resource.WithFromEnv(),      // Discover and provide attributes from OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME environment variables.
			resource.WithTelemetrySDK(), // Discover and provide information about the OpenTelemetry SDK used.
			resource.WithProcess(),      // Discover and provide process information.
			resource.WithOS(),           // Discover and provide OS information.
			resource.WithContainer(),    // Discover and provide container information.
			resource.WithHost(),         // Discover and provide host information.
			//resource.WithAttributes(attribute.String("foo", "bar")), // Add custom resource attributes.
			//resource.WithDetectors(thirdparty.Detector{}), // Bring your own external Detector implementation.
		)
		if errors.Is(err, resource.ErrPartialResource) || errors.Is(err, resource.ErrSchemaURLConflict) {
			fmt.Println(err.Error()) // Log non-fatal issues.
			return nil
		} else if err != nil {
			fmt.Println(err.Error()) // The error may be fatal.
			return nil
		}

		lp, err := newLoggerProvider(context.Background(), res)
		if err != nil {
			fmt.Println("failed to create logger: %w", err)
			return nil
		}

		otelCore := otelzap.NewCore("XLogger", otelzap.WithLoggerProvider(lp))
		coreArr = append(coreArr, infoFileCore, otelCore)
	} else {
		coreArr = append(coreArr, infoFileCore)
	}
	zlog := zap.New(zapcore.NewTee(coreArr...), zap.AddCaller(), zap.AddCallerSkip(3)) //zap.AddCaller()为显示文件名和行号，可省略
	return zlog
}

func nCallerEncoder(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
	ss := trimPath(caller) + ":" + strconv.FormatInt(int64(caller.Line), 10)
	encoder.AppendString(ss)
}

func trimPath(caller zapcore.EntryCaller) string {
	idx0 := strings.LastIndexByte(caller.File, '/')
	if idx0 == -1 {
		return caller.File
	}
	idx1 := strings.LastIndexByte(caller.File[:idx0], '/')
	if idx1 == -1 {
		return caller.File
	}

	idx2 := strings.LastIndex(caller.File, "@")
	if idx2 == -1 {
		return caller.File[idx0+1:]
	} else {
		idx3 := strings.LastIndexByte(caller.File[:idx2], '/')
		if idx3 == -1 {
			return caller.File[idx0+1:]
		}
		var bd strings.Builder
		bd.Grow(idx2 - idx3)
		bd.WriteByte('[')
		cg := false
		for i := idx3 + 1; i < idx2; i++ {
			c := caller.File[i]
			if c == '!' {
				cg = true
			} else if cg {
				if c >= 'a' && c <= 'z' {
					c = c + 'A' - 'a'
				}
				bd.WriteByte(c)
				cg = false
			} else {
				bd.WriteByte(c)
			}
		}
		bd.WriteByte(']')
		prefix := bd.String()

		//prefix := "[" + caller.File[idx3+1:idx2] + "] "
		//prefix = strings.Title(prefix)
		//prefix = strings.Replace(prefix, "!", "", -1)
		if idx3 == idx1 {
			return prefix + caller.File[idx0+1:]
		} else {
			return prefix + caller.File[idx1+1:]
		}
	}
}

func NewZapLogger(zlog *zap.Logger) *Logger {
	return &Logger{
		log:    zlog,
		msgKey: log.DefaultMessageKey,
	}
}

func (l *Logger) Log(level log.Level, keyvals ...any) error {
	// If logging at this level is completely disabled, skip the overhead of
	// string formatting.
	if zapcore.Level(level) < zapcore.DPanicLevel && !l.log.Core().Enabled(zapcore.Level(level)) {
		return nil
	}
	var (
		msg    = ""
		keylen = len(keyvals)
	)
	if keylen == 0 || keylen%2 != 0 {
		l.log.Warn(fmt.Sprint("Key values must appear in pairs: ", keyvals))
		return nil
	}

	data := make([]zap.Field, 0, (keylen/2)+1)
	for i := 0; i < keylen; i += 2 {
		if keyvals[i].(string) == l.msgKey {
			msg, _ = keyvals[i+1].(string)
			continue
		}
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	switch level {
	case log.LevelDebug:
		l.log.Debug(msg, data...)
	case log.LevelInfo:
		l.log.Info(msg, data...)
	case log.LevelWarn:
		l.log.Warn(msg, data...)
	case log.LevelError:
		l.log.Error(msg, data...)
	case log.LevelFatal:
		l.log.Fatal(msg, data...)
	}
	return nil
}

func (l *Logger) Sync() error {
	return l.log.Sync()
}

func (l *Logger) Close() error {
	return l.Sync()
}

func Factory(cfg *conf.Logger) (log.Logger, error) {
	if cfg.GetZap() == nil {
		err := errors.New("logger zap config is nil")
		return nil, err
	}
	zlog := InitZapWithConfig(cfg.GetZap())
	logS := NewZapLogger(zlog)

	return logS, nil
}

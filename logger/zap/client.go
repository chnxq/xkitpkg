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

	"github.com/chnxq/xkitmod/log"
	conf "github.com/chnxq/xkitpkg/conf/v1"
	"github.com/chnxq/xkitpkg/logger"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	_ = logger.RegisterFactory(logger.Zap, func(cfg *conf.Logger) (log.Logger, error) {
		return NewLogger(cfg)
	})
}

// NewLogger 创建一个新的日志记录器 - Zap
func NewLogger(cfg *conf.Logger) (log.Logger, error) {
	if cfg == nil || cfg.Zap == nil {
		return nil, nil
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	} //指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder //按级别显示不同颜色，不需要的话取值zapcore.CapitalLevelEncoder就可以了

	switch cfg.Zap.GetCaller() {
	case "short":
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	case "full":
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	case "xkit":
		encoderConfig.EncodeCaller = nCallerEncoder
	default:
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	}

	encoderConfig.CallerKey = "C"
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	lumberJackLogger := &lumberjack.Logger{
		Filename:   cfg.Zap.LogFilePath + "/info.log", //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    int(cfg.Zap.MaxSize),
		MaxBackups: int(cfg.Zap.MaxBackups),
		MaxAge:     int(cfg.Zap.MaxAge),
	}
	writeSyncer := zapcore.AddSync(lumberJackLogger)

	//日志级别
	lvl, err := zapcore.ParseLevel(cfg.Zap.GetLevel())
	if err != nil {
		return nil, err
	}

	var core zapcore.Core
	if cfg.Zap == nil || cfg.Zap.LogToConsole {
		// 如果配置为输出到文件和控制台 // 使用 MultiWriteSyncer 组合文件写入器和标准输出写入器
		core = zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout)), lvl)
	} else {
		// 否则，使用 JSON 编码器，只输出到文件
		encoder = zapcore.NewJSONEncoder(encoderConfig)
		core = zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writeSyncer), lvl)
	}
	//core := zapcore.NewCore(encoder, writeSyncer, lvl)

	var coreArr []zapcore.Core
	if cfg.Zap.ExportToOtel { // 如果配置为导出到 OpenTelemetry 日志导出器
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
			return nil, err
		} else if err != nil {
			fmt.Println(err.Error()) // The error may be fatal.
			return nil, err
		}

		lp, err := newLoggerProvider(context.Background(), res, cfg.Zap.ExporterEndpoint, cfg.Zap.ExporterInsecure)
		if err != nil {
			fmt.Println("failed to create logger exporter: %w", err)
			return nil, err
		}
		if lp != nil {
			otelCore := otelzap.NewCore("OtelLogger", otelzap.WithLoggerProvider(lp))
			coreArr = append(coreArr, core, otelCore)
		} else {
			coreArr = append(coreArr, core)
		}
	} else {
		coreArr = append(coreArr, core)
	}

	l := zap.New(zapcore.NewTee(coreArr...), zap.AddCaller(), zap.AddCallerSkip(3)).WithOptions() //zap.AddCaller()为显示文件名和行号，可省略
	//l := zap.New(core).WithOptions()

	wrapped := &Logger{
		log:    l,
		msgKey: log.DefaultMessageKey,
	}

	return wrapped, nil
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

// newLoggerProvider creates a new logger provider with the OTLP gRPC exporter.
func newLoggerProvider(ctx context.Context, res *resource.Resource, exporterEndpoint string, exporterInsecure bool) (*otelLog.LoggerProvider, error) {
	opts := []otlploggrpc.Option{}
	if exporterInsecure {
		opts = append(opts, otlploggrpc.WithInsecure())
	}
	if exporterEndpoint != "" {
		opts = append(opts, otlploggrpc.WithEndpoint(exporterEndpoint))
	} else {
		return nil, errors.New("exporter endpoint is required")
	}
	exporter, err := otlploggrpc.New(ctx, opts...)
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

package log

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// DefaultTimeEncoder 默认时间编码器
	DefaultTimeEncoder = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000")
)

// Logger 包装 zap.Logger，提供清理功能
type Logger struct {
	*zap.Logger
	closers []io.Closer
}

// Close 关闭所有资源
func (l *Logger) Close() error {
	var firstErr error
	for _, closer := range l.closers {
		if err := closer.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// New 创建日志实例（简化版）
func New(isDev bool, name ...string) (*Logger, error) {
	return NewWithConfig(&Config{IsDev: isDev, Level: "info"}, name...)
}

// NewWithConfig 根据配置创建日志实例
func NewWithConfig(cfg *Config, name ...string) (*Logger, error) {
	if cfg == nil {
		cfg = &Config{IsDev: false, Level: "info"}
	}
	var err error
	var lvl zapcore.Level
	if cfg.Level == "" {
		lvl = zapcore.InfoLevel
	} else {
		lvl, err = zapcore.ParseLevel(cfg.Level)
		if err != nil {
			return nil, fmt.Errorf("invalid log level %q: %w", cfg.Level, err)
		}
	}
	encoder := createEncoder(cfg.IsDev)
	// 默认输出到 stdout
	cores := []zapcore.Core{
		zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), lvl),
	}
	var closers []io.Closer
	// 处理额外的适配器
	if len(cfg.Adaptors) > 0 {
		for _, adaptorDSN := range cfg.Adaptors {
			core, closer, err := createAdaptorCore(adaptorDSN, encoder, lvl)
			if err != nil {
				if closer != nil {
					_ = closer.Close()
				}
				continue
			}
			cores = append(cores, core)
			if closer != nil {
				closers = append(closers, closer)
			}
		}
		// 如果有，则至少创建一个适配器
		if len(cores) == 0 {
			return nil, fmt.Errorf("no adaptors created")
		}
	}
	zapLogger := zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
		// zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if len(name) > 0 {
		zapLogger = zapLogger.Named(name[0]).With(zap.String("service", name[0]))
	}
	zap.ReplaceGlobals(zapLogger)
	return &Logger{Logger: zapLogger, closers: closers}, nil
}

// createEncoder 创建编码器
func createEncoder(isDev bool) zapcore.Encoder {
	var encoderCfg zapcore.EncoderConfig
	if isDev {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
		encoderCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	encoderCfg.EncodeTime = DefaultTimeEncoder
	encoderCfg.EncodeDuration = zapcore.StringDurationEncoder
	if isDev {
		return zapcore.NewConsoleEncoder(encoderCfg)
	}
	return zapcore.NewJSONEncoder(encoderCfg)
}

// createAdaptorCore 根据 DSN 创建对应的 Core
func createAdaptorCore(dsn string, encoder zapcore.Encoder, lvl zapcore.Level) (zapcore.Core, io.Closer, error) {
	parsedURL, err := url.Parse(dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid adaptor DSN: %w", err)
	}
	switch parsedURL.Scheme {
	case "file":
		opts, err := parseFileOptions(dsn)
		if err != nil {
			return nil, nil, err
		}
		writer, closer, err := newFileWriter(opts)
		if err != nil {
			return nil, nil, err
		}
		return zapcore.NewCore(encoder, writer, lvl), closer, nil
	case "http", "https":
		opts, err := parseHTTPOptions(dsn)
		if err != nil {
			return nil, nil, err
		}
		writer, closer, err := newHTTPWriter(opts)
		if err != nil {
			return nil, nil, err
		}
		return zapcore.NewCore(encoder, writer, lvl), closer, nil
	default:
		return nil, nil, fmt.Errorf("unsupported scheme: %s", parsedURL.Scheme)
	}
}

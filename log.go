package log

import (
	"net/url"
	"os"

	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Config struct {
		IsDev    bool     `json:"is_dev" yaml:"is_dev"`
		Level    string   `json:"level" yaml:"level"`
		Adaptors []string `json:"adaptors,omitzero" yaml:"adaptors"`
	}
)

func FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("log", pflag.ContinueOnError)
	fs.Bool("log.dev", false, "use DevelopmentEncoder")
	fs.StringSlice("log.adaptors", []string{}, "log adaptors")
	return fs
}

func New(isDev bool, name ...string) (*zap.Logger, error) {
	return NewWithConfig(&Config{IsDev: isDev}, name...)
}

func NewWithConfig(cfg *Config, name ...string) (*zap.Logger, error) {
	var (
		encoder     zapcore.Encoder
		timeEncoder = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05")
	)
	lvl, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		lvl = zapcore.InfoLevel
	}

	if cfg.IsDev {
		encoderCfg := zap.NewDevelopmentEncoderConfig()
		encoderCfg.EncodeTime = timeEncoder
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.EncodeTime = timeEncoder
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}
	cores := []zapcore.Core{
		zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), lvl),
	}
	if len(cfg.Adaptors) > 0 {
		for _, adaptor := range cfg.Adaptors {
			_url, err := url.Parse(adaptor)
			if err != nil {
				return nil, err
			}
			switch _url.Scheme {
			case "file":
				writeSyncer, err := newFileWriter(_url.Path)
				if err != nil {
					return nil, err
				}
				cores = append(cores, zapcore.NewCore(encoder, writeSyncer, lvl))
			case "http", "https":
				cores = append(cores, zapcore.NewCore(encoder, newHTTPWriter(adaptor), lvl))
			}
		}
	}
	logger := zap.New(zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.DPanicLevel),
	)
	if len(name) > 0 {
		logger = logger.Named(name[0]).With(zap.String("service", name[0]))
	}
	zap.ReplaceGlobals(logger)
	return logger, nil
}

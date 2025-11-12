package log

import "github.com/spf13/pflag"

type Config struct {
	Level    string   `json:"level" yaml:"level"`       // 日志级别: debug, info, warn, error
	Adaptors []string `json:"adaptors" yaml:"adaptors"` // 输出适配器 DSN 列表
	IsDev    bool     `json:"dev" yaml:"dev"`           // 开发模式
	Skip     int      `json:"skip" yaml:"skip"`         // 跳过调用栈层数
}

func (c *Config) FlagSet() *pflag.FlagSet { return FlagSet() }

func FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("log", pflag.ContinueOnError)
	fs.Bool("log.dev", false, "use development encoder")
	fs.String("log.level", "info", "log level: debug, info, warn, error")
	fs.StringSlice("log.adaptors", []string{}, "log adaptors DSN (e.g., file:///var/log/app.log?max-size=100m&max-age=30d)")
	fs.Int("log.skip", 1, "log skip caller stack frames")
	return fs
}

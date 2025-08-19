package log

import "github.com/spf13/pflag"

type (
	Config struct {
		IsDev    bool     `json:"is_dev" yaml:"is_dev"`
		Level    string   `json:"level" yaml:"level"`
		Adaptors []string `json:"adaptors,omitzero" yaml:"adaptors"`
	}
)

func (c *Config) FlagSet() *pflag.FlagSet { return FlagSet() }

func FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("log", pflag.ContinueOnError)
	fs.Bool("log.dev", false, "use DevelopmentEncoder")
	fs.StringSlice("log.adaptors", []string{}, "log adaptors")
	return fs
}

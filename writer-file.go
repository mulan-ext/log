package log

import (
	"os"
	"path/filepath"

	"go.uber.org/zap/zapcore"
)

type FileWriter struct{}

func newFileWriter(file string) (zapcore.WriteSyncer, error) {
	err := os.MkdirAll(filepath.Dir(file), 0755)
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return zapcore.Lock(f), nil
}

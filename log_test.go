package log_test

import (
	"os"
	"testing"

	"github.com/mulan-ext/log"
	"go.uber.org/zap"
)

func TestLog(t *testing.T) {
	log.NewWithConfig(&log.Config{Level: "debug"})
	zap.L().Info("test")
	zap.L().Info("testaaaaaaaaaaa")
	zap.L().Info("testaaaa", zap.String("aa", "awdvews"))
}

func TestLogDebug(t *testing.T) {
	log.NewWithConfig(&log.Config{IsDev: true})
	zap.L().Info("test")
	zap.L().Info("testaaaaaaaaaaa")
	zap.L().Debug("testaaaaaargverasrfesaaaaaa")
}

func TestLogFile(t *testing.T) {
	f, err := os.CreateTemp("/tmp", "testlog")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("log file:", f.Name())
	log.NewWithConfig(&log.Config{
		Level:    "debug",
		Adaptors: []string{"file://" + f.Name()},
	})
	zap.L().Info("test")
	zap.L().Info("testaaaaaaaaaaa")
	zap.L().Info("testaaaaaargverasrfesaaaaaa")
}

func TestLogHTTP(t *testing.T) {
	log.NewWithConfig(&log.Config{
		Level:    "debug",
		Adaptors: []string{"http://localhost:3003/log"},
	})
	zap.L().Info("test")
	zap.L().Info("testaaaaaaaaaaa")
	zap.L().Info("testaaaaaargverasrfesaaaaaa")
}

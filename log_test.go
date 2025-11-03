package log_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mulan-ext/log"
	"go.uber.org/zap"
)

func TestLog(t *testing.T) {
	logger, err := log.NewWithConfig(&log.Config{Level: "debug"})
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	zap.L().Info("test")
	zap.L().Info("testaaaaaaaaaaa")
	zap.L().Info("testaaaa", zap.String("aa", "awdvews"))
}

func TestLogDebug(t *testing.T) {
	logger, err := log.NewWithConfig(&log.Config{IsDev: true})
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	zap.L().Info("test")
	zap.L().Info("testaaaaaaaaaaa")
	zap.L().Debug("testaaaaaargverasrfesaaaaaa")
}

func TestLogFile(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")
	t.Log("log file:", logFile)

	logger, err := log.NewWithConfig(&log.Config{
		Level:    "debug",
		Adaptors: []string{"file://" + logFile},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	zap.L().Info("test")
	zap.L().Info("testaaaaaaaaaaa")
	zap.L().Info("testaaaaaargverasrfesaaaaaa")

	// 验证文件是否存在
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("log file not created")
	}
}

func TestLogFileDSN(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	// 使用 DSN 格式配置
	dsn := fmt.Sprintf("file://%s?max-size=10m&max-backups=5&max-age=7d&compress=gzip", logFile)
	t.Log("DSN:", dsn)

	logger, err := log.NewWithConfig(&log.Config{
		Level:    "info",
		Adaptors: []string{dsn},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	zap.L().Info("test DSN configuration")

	// 验证文件是否存在
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("log file not created")
	}
}

func TestLogFileRotation(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test-rotation.log")

	// 使用 DSN 配置滚动
	dsn := fmt.Sprintf("file://%s?max-size=1m&max-backups=3&max-age=7d", logFile)

	logger, err := log.NewWithConfig(&log.Config{
		Level:    "debug",
		Adaptors: []string{dsn},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	// 写入大量日志触发滚动
	for i := 0; i < 10000; i++ {
		zap.L().Info("test log rotation", zap.Int("iteration", i))
	}

	t.Log("log file created:", logFile)
}

func TestLogHTTP(t *testing.T) {
	t.Skip("需要实际的 HTTP 服务端")

	// 使用 DSN 配置 HTTP
	dsn := "http://localhost:3003/log?timeout=5s&buffer-size=512&batch-size=50&max-retries=2"

	logger, err := log.NewWithConfig(&log.Config{
		Level:    "debug",
		Adaptors: []string{dsn},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	zap.L().Info("test")
	zap.L().Info("testaaaaaaaaaaa")
	zap.L().Info("testaaaaaargverasrfesaaaaaa")

	// 等待异步发送完成
	time.Sleep(2 * time.Second)
}

func TestLogWithName(t *testing.T) {
	logger, err := log.New(true, "my-service")
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	zap.L().Info("test with service name")
}

func TestLogInvalidLevel(t *testing.T) {
	_, err := log.NewWithConfig(&log.Config{Level: "invalid"})
	if err == nil {
		t.Error("expected error for invalid log level")
	}
}

func TestLogMultipleAdaptors(t *testing.T) {
	tmpDir := t.TempDir()
	logFile1 := filepath.Join(tmpDir, "test1.log")
	logFile2 := filepath.Join(tmpDir, "test2.log")

	logger, err := log.NewWithConfig(&log.Config{
		Level: "info",
		Adaptors: []string{
			"file://" + logFile1,
			"file://" + logFile2,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	zap.L().Info("test multiple adaptors")

	// 验证两个文件都被创建
	if _, err := os.Stat(logFile1); os.IsNotExist(err) {
		t.Error("log file 1 not created")
	}
	if _, err := os.Stat(logFile2); os.IsNotExist(err) {
		t.Error("log file 2 not created")
	}
}

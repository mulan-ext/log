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

func logPrint() {
	zap.L().Info("test info")
	zap.L().Warn("test warn")
	zap.L().Error("test error")
	zap.L().Debug("test debug")
	zap.L().Info("test info", zap.String("aa", "awdvews"))
	zap.L().Warn("test warn", zap.String("aa", "awdvews"))
	zap.L().Error("test error", zap.String("aa", "awdvews"))
	zap.L().Debug("test debug", zap.String("aa", "awdvews"))
}

func TestLog(t *testing.T) {
	logger, err := log.NewWithConfig(&log.Config{Level: "info"})
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	logPrint()
}

func TestLogDebug(t *testing.T) {
	logger, err := log.NewWithConfig(&log.Config{Level: "debug"})
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	logPrint()
}

func TestLogFile(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "_test.log")
	t.Log("log file:", logFile)

	logger, err := log.NewWithConfig(&log.Config{
		Level:    "debug",
		Adaptors: []string{"file://" + logFile},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	logPrint()

	// 验证文件是否存在
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("log file not created")
	}

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("log file content:", string(content))
}

func TestLogFileLevel(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "_test.log")
	t.Log("log file:", logFile)

	logger, err := log.NewWithConfig(&log.Config{
		Level:    "info",
		Adaptors: []string{"file://" + logFile + "?level=debug"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	logPrint()

	// 验证文件是否存在
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("log file not created")
	}

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("log file content:", string(content))
}

func TestLogFileDSN(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "_test.log")

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

	logPrint()

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
	for i := range 10000 {
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

	logPrint()
	// 等待异步发送完成
	time.Sleep(2 * time.Second)
}

func TestLogWithName(t *testing.T) {
	logger, err := log.New(true, "my-service")
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	logPrint()
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

	logPrint()

	// 验证两个文件都被创建
	if _, err := os.Stat(logFile1); os.IsNotExist(err) {
		t.Error("log file 1 not created")
	}
	if _, err := os.Stat(logFile2); os.IsNotExist(err) {
		t.Error("log file 2 not created")
	}
}

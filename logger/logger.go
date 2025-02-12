package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger 初始化 zap Logger
func InitLogger() *zap.Logger {
	// 設定 zap 配置
	config := zapcore.EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalColorLevelEncoder, // INFO, WARN, ERROR
		EncodeTime:    zapcore.ISO8601TimeEncoder,       // 使用 ISO8601 格式
		EncodeCaller:  zapcore.ShortCallerEncoder,       // 簡化 caller 顯示
	}

	// 設定 log 輸出
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(config),    // 以 console 模式輸出
		zapcore.Lock(os.Stdout),              // 輸出到 stdout
		zap.NewAtomicLevelAt(zap.DebugLevel), // 設定日誌等級
	)

	// 建立 logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	zap.ReplaceGlobals(logger) // 設定全局 logger
	return logger
}

package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger 初始化 zap Logger（同時輸出到 stdout 和檔案）
func InitLogger() *zap.Logger {
	// 定義日誌格式
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalColorLevelEncoder, // INFO, WARN, ERROR
		EncodeTime:    zapcore.ISO8601TimeEncoder,       // 以 ISO8601 格式記錄時間
		EncodeCaller:  zapcore.ShortCallerEncoder,       // 簡化 caller 顯示
	}

	// 建立 Console Encoder（標準輸出）
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 建立 File Encoder（寫入檔案）
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig) // 以 JSON 格式寫入檔案

	// 設定 log 檔案位置
	logFilePath := "logs/app.log"
	file, _ := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	// 建立不同的 log 輸出
	consoleWriter := zapcore.Lock(os.Stdout) // 標準輸出
	fileWriter := zapcore.AddSync(file)      // 寫入檔案

	// 設定 log 等級
	logLevel := zap.NewAtomicLevelAt(zap.DebugLevel)

	// 讓 log 同時輸出到 stdout 和 檔案
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleWriter, logLevel),
		zapcore.NewCore(fileEncoder, fileWriter, logLevel),
	)

	// 建立 Logger，並加入 Caller 來顯示錯誤位置
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// 設定 zap 全局 Logger
	zap.ReplaceGlobals(logger) // 設定全局 logger

	return logger
}

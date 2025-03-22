package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinZapLoggerMiddleware 是一個 Gin 中間件，整合 zap 來記錄請求日誌
func GinZapLoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 處理請求
		c.Next()

		// 計算請求處理時間
		latency := time.Since(startTime)

		// 記錄請求資訊
		logger.Info("HTTP Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.String("client_ip", c.ClientIP()),
			zap.Duration("latency", latency),
			zap.String("user_agent", c.Request.UserAgent()),
		)

		// 確保log在中間件結束時刷新和關閉
		defer logger.Sync()
	}
}

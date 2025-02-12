package main

import (
	"fmt"

	"GoGinService/logger"     // 匯入 logger
	"GoGinService/middleware" // 匯入 middleware

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello, world GO GIN Sevice!")
	// 初始化 zap logger
	log := logger.InitLogger()
	defer log.Sync() // 確保 log buffer 被寫入

	// 設定 Gin 伺服器
	r := gin.New()

	// 使用 gin-zap 記錄 middleware
	r.Use(middleware.GinZapLoggerMiddleware(log))

	// 測試 API
	r.GET("/", func(c *gin.Context) {
		log.Info("Hello Gin with Zap")
		c.JSON(200, gin.H{"message": "Hello, world!"})
	})

	// 啟動伺服器
	port := ":8080"
	fmt.Println("Server is running at http://localhost" + port)
	r.Run(port)
}

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"GoGinService/logger"     // 匯入 logger
	"GoGinService/middleware" // 匯入 middleware

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

	// // 啟動伺服器
	// port := ":8080"
	// fmt.Println("Server is running at http://localhost" + port)
	// r.Run(port)

	// Graceful_Shutdown 功能
	// 透過 http.Server 啟動 Gin
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	// 啟動 Gin 服務
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server start error", zap.Error(err))
		}
	}()
	fmt.Println("Server is running at http://localhost:8080")

	// 監聽系統訊號
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit // 阻塞，直到收到退出訊號
	log.Info("Shutting down server...")

	// 設定 5 秒的 timeout，讓請求完成
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exiting")

}

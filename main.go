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
	"go.uber.org/dig"
	"go.uber.org/zap"

	// swagger 套件引入
	_ "GoGinService/docs" // 引入 docs 產生的 swagger 文件

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// var log *zap.Logger

// PingHandler godoc
// @Summary 打招呼的API
// @Description 回傳 Hello world 訊息
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "成功回應"
// @Router /ping [get]
func PingHandler(c *gin.Context) {
	log, exists := c.Get("logger")
	if !exists {
		panic("Logger 未設定")
	}
	logger := log.(*zap.Logger)

	logger.Info("Hello Gin with Zap")
	c.JSON(200, gin.H{"message": "Hello, world!"})
}

// DingHandler godoc
// @Summary 打招呼Dig的API
// @Description 回傳 Hello Dig 訊息
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "成功回應"
// @Router /dig [get]
func DingHandler(c *gin.Context) {
	log, exists := c.Get("logger")
	if !exists {
		panic("Logger 未設定")
	}
	logger := log.(*zap.Logger)

	logger.Info("Hello Dig with Zap")
	c.JSON(200, gin.H{"message": "Hello, Dig!"})
}

// @title GoGinService API
// @version 1.0
// @description 這是 GoGinService 的 Swagger API 文件範例。
// @host localhost:8080
// @BasePath /
func main() {
	fmt.Println("Hello, world GO GIN Sevice!")
	// 初始化 zap logger
	log := logger.InitLogger()
	defer log.Sync() // 確保 log buffer 被寫入

	// 使用 dig 進行依賴注入
	container := dig.New()
	if err := container.Provide(logger.InitLogger); err != nil {
		fmt.Println("Failed to provide logger:", err)
		return
	}

	// 設定 Gin 伺服器
	r := gin.New()

	// 使用 gin-zap 記錄 middleware
	r.Use(middleware.GinZapLoggerMiddleware(log))

	// // 測試 API
	// // GetHello godoc
	// // @Summary 打招呼的API
	// // @Description 回傳 Hello world 訊息
	// // @Tags example
	// // @Accept json
	// // @Produce json
	// // @Success 200 {object} map[string]string "成功回應"
	// // @Router / [get]
	// r.GET("/", func(c *gin.Context) {
	// 	log.Info("Hello Gin with Zap")
	// 	c.JSON(200, gin.H{"message": "Hello, world!"})
	// })

	// 將 logger 注入到 gin context 中
	r.Use(func(c *gin.Context) {
		c.Set("logger", log)
		c.Next()
	})
	r.GET("/ping", PingHandler)

	// dig 依賴注入「/dig」
	if err := container.Invoke(func(log *zap.Logger) {
		r.GET("/dig", DingHandler)
	}); err != nil {
		log.Fatal("Failed to invoke DingHandler:", zap.Error(err))
	}

	// Swagger API 文件路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

package main

import (
	"context"
	"exchangeapp/config"
	"exchangeapp/router"
	"exchangeapp/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config.InitConfig()

	r := router.SetupRouter()

	port := config.AppConfig.App.Port
	if port == "" {
		port = ":8080"
	}

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	go func() {
		utils.Info("服务器启动中，监听端口: %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Fatal("服务器启动失败: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	utils.Info("开始关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		utils.LogError("服务器关闭失败: %v", err)
	}

	utils.Info("服务器已退出")
	utils.CloseLogger()
}

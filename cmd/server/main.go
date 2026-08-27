package main

import (
	"context"
	"errors"
	"log/slog"
	_ "minsky/go-template/docs"          // swagger WebUI访问需要
	_ "minsky/go-template/internal/init" // 触发viper读取json配置
	"minsky/go-template/internal/routers"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// @title go-template(项目名)
// @version v1.2
// @description 请求状态码定义
// @description code= 0, 调用成功
// @description code=-1, 系统错误
// @description code= 1, 提示接口返回的message
// @description code=10, 登录token过期
// @description code=20, 接口权限错误,没有权限访问该接口
// @termsOfService http://localhost:8000/swagger/index.html
// @license.name Apache 2.0
// @host localhost:8000
func main() {

	appCtx, stopSignal := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stopSignal()

	gin.ForceConsoleColor()
	router := gin.Default()
	routers.BindApi(router)

	server := &http.Server{
		Addr:              ":8000",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		slog.Info("Gin 服务启动", "url", "http://127.0.0.1"+server.Addr)
		err := server.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()

	select {
	case <-appCtx.Done():
		slog.Info("收到关闭信号:", context.Cause(appCtx))
	case err := <-serverErr:
		if err != nil {
			slog.Info("HTTP 服务异常退出:", err)
		} else {
			slog.Info("HTTP 服务已经停止")
		}
	}

	stopSignal()

	slog.Info("HTTP服务开始关闭")
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	// 停止接受新请求, 并等待正在处理的请求完成
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP服务关闭超时:%v", err)

		if closeErr := server.Close(); closeErr != nil {
			slog.Error("强制关闭HTTP服务失败: %v", closeErr)
		}
	} else {
		slog.Info("HTTP服务已关闭")
	}

}

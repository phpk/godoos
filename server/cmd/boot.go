package cmd

import (
	"context"
	_ "godocms/app"
	"godocms/config"
	"godocms/middleware"
	"godocms/pkg/dbfactory"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var Server *http.Server

func Start() {
	dbfactory.InitDatabase()
	config.LoadConfig()
	handler, err := config.InitLogger()
	if err != nil {
		slog.Error("Failed to initialize logger:", "error", err)
	}

	if config.Config.System.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// 处理异常
	r.NoRoute(middleware.HandleNotFound)
	r.NoMethod(middleware.HandleNotFound)
	r.Use(middleware.Recover)

	// 使用中间件
	r.Use(middleware.Cors())
	r.Use(middleware.LoggerMiddleware(handler))
	store := middleware.GetSessionStore()
	r.Use(sessions.Sessions("godoSession", store))
	middleware.BindRouter(r)
	// 将端口号转换为字符串
	portStr := strconv.Itoa(config.Config.System.Port)

	// 创建并返回 http.Server 实例
	Server = &http.Server{
		Addr:    ":" + portStr,
		Handler: r,
	}
	go func() {
		if err := Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server:", "error", err)
			os.Exit(1)
		}
	}()
	// 监听信号来重启服务
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		for sig := range sigChan {
			switch sig {
			case os.Interrupt, syscall.SIGTERM:
				slog.Info("Received SIGTERM or SIGINT, shutting down...")
				Shutdown()
				os.Exit(0)
			case syscall.SIGHUP:
				slog.Info("Received SIGHUP, restarting...")
				Restart()
			default:
				slog.Info("Received unknown signal:", "signal", sig)
			}
		}
	}()

	// 主循环
	select {}
}
func Restart() {
	if Server == nil {
		slog.Error("Server is not running.")
		return
	}
	// 关闭当前服务
	Shutdown()

	// 重新启动服务
	Start()
}

// 关闭 HTTP 服务器
func Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭服务器
	if err := Server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server:", "error", err)
	} else {
		slog.Info("Server gracefully shutdown")
	}
}

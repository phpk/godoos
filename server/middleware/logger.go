package middleware

import (
	"godocms/common"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware 是一个简单的日志中间件
func LoggerMiddleware(handler slog.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求的 Host
		host := c.Request.Host
		// 获取请求的 Scheme
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		common.Config.System.Host = host
		common.Config.System.Scheme = scheme
		url := c.Request.URL.Path
		// 检查 URL 是否以 /static/ 开头，如果是则跳过记录
		if strings.HasPrefix(url, "/static/") || strings.HasPrefix(url, "/upload/") {
			c.Next() // 继续后续处理
			return
		}
		startTime := time.Now()
		c.Next() // 执行后续的处理函数
		logger := slog.New(handler)
		// 记录请求的信息
		logger.Info("Request completed",
			slog.String("method", c.Request.Method),
			slog.String("path", url),
			slog.Any("status", c.Writer.Status()),
			slog.Duration("duration", time.Since(startTime)),
		)
	}
}

package middleware

import (
	"godocms/libs"
	"log"

	"github.com/gin-gonic/gin"
)

// HandleNotFound 自定义 404 处理器
func HandleNotFound(c *gin.Context) {
	// 重定向到自定义的 404 页面
	c.Redirect(302, "/home/notfond/")
}

// Recover 自定义 panic 处理器
func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			// 打印错误堆栈信息
			log.Printf("panic: %v\n", r)
			libs.Error(c, "服务器内部错误")
		}
	}()

	// 加载完 defer recover，继续后续接口调用
	c.Next()
}

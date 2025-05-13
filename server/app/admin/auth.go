package admin

import (
	"godocms/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	middleware.RegisterRouter("POST", "/api/v1/test", test, 0, "test")
}
func test(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "hello",
	})
}

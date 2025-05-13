package libs

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// 请求成功的时候 使用该方法返回信息
func Success(ctx *gin.Context, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"success": true,
		"message": msg,
		"data":    data,
		"time":    time.Now().Unix(),
	})
}

// 请求失败的时候, 使用该方法返回信息
func Error(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    -1,
		"success": false,
		"message": msg,
		"time":    time.Now().Unix(),
	})
}

func ErrorMsg(ctx *gin.Context, code int, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    code,
		"success": true,
		"message": msg,
		"data":    data,
		"time":    time.Now().Unix(),
	})
}
func ErrorLogin(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    -2,
		"success": false,
		"message": msg,
		"time":    time.Now().Unix(),
	})
}
func GetPage(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}
	size, _ := strconv.Atoi(c.Query("size"))
	if size <= 0 {
		size = 10
	}
	return page, size
}

type PaginatedResult struct {
	Data       interface{} `json:"data"`        // 当前页数据
	Page       int         `json:"page"`        // 当前页码
	Size       int         `json:"size"`        // 每页大小
	TotalCount int64       `json:"total_count"` // 总记录数
	TotalPages int         `json:"total_pages"` // 总页数
}

package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// Route 结构体用于存储路由信息
type Route struct {
	Method  string          `json:"method"`
	Path    string          `json:"path"`
	Name    string          `json:"name"`
	Handler gin.HandlerFunc `json:"-"`
	// 添加一个字段表示是否需要权限校验
	NeedAuth int `json:"needAuth"`
}

// routes 存储所有需要注册的路由
var Routes = make(map[string]Route)

// RegisterRouter 注册控制器中的路由
func RegisterRouter(method string, path string, handler gin.HandlerFunc, needAuth int, name string) {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	// 拼接路由地址
	if path == "index" {
		path = "/"
	} else {
		path = "/" + strings.ToLower(path)
	}

	// 构造键值
	//key := method + ":" + path

	// 添加路由
	Routes[path] = Route{
		Path:     path,
		Method:   method,
		Handler:  handler,
		NeedAuth: needAuth,
		Name:     name,
	}

	//slog.Info("Register route", "method", method, "path", path)
}

// BindRouter 绑定所有注册的路由到 Gin 引擎
func BindRouter(e *gin.Engine) {
	for _, route := range Routes {
		switch route.Method {
		case "GET":
			e.GET(route.Path, route.Handler)
		case "POST":
			e.POST(route.Path, route.Handler)
		case "DELETE":
			e.DELETE(route.Path, route.Handler)
		case "PUT":
			e.PUT(route.Path, route.Handler)
		}
	}
}

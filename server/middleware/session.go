package middleware

import (
	"godocms/common"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-contrib/sessions/memstore"
)

func GetSessionStore() gormsessions.Store {

	if common.Config.System.SessionType == "cookie" {
		// 创建一个 cookie store
		store := cookie.NewStore([]byte(common.Config.System.SessionSecret))

		// 配置 cookie，使用 github.com/gin-contrib/sessions.Options 类型
		store.Options(sessions.Options{
			Path:     "/",
			HttpOnly: true,
			Secure:   false, // 如果使用 HTTPS
			SameSite: http.SameSiteNoneMode,
			MaxAge:   int(24 * time.Hour / time.Second), // MaxAge 单位为秒
		})

		return store
	} else {
		return memstore.NewStore([]byte(common.Config.System.SessionSecret))
	}
}

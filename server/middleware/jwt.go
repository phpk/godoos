package middleware

import (
	"godocms/common"
	"godocms/libs"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func JwtVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Request.URL.Path
		// 检查 URL 是否以 /static/ 开头，如果是则跳过记录
		if strings.HasPrefix(url, "/static/") || strings.HasPrefix(url, "/upload/") || strings.HasPrefix(url, "/views/") || strings.HasPrefix(url, "/os/") {
			c.Next() // 继续后续处理
			return
		}
		// 根据动态路由模式获取路由信息
		routeData, ok := getRouteData(url)
		//未注册路由
		if !ok {
			log.Printf("未注册路由:%v", url)
			// log.Printf("routes:%v", common.Routes)
			libs.Error(c, "未注册路由") // 重定向到登录页面
			c.Abort()              // 终止后续处理
			return
		}
		if routeData.NeedAuth == 0 { // 过滤附件访问接口
			return
		}

		token := c.GetHeader("Authorization")
		if routeData.NeedAuth == 3 {
			token = c.GetHeader("AuthorizationAdmin")
		}
		if token == "" {
			token = c.Query("token")
		}
		if token == "" {
			libs.ErrorLogin(c, "token is empty")
			c.Abort() // 终止后续处理
			return
		}
		// 验证 token，并存储在请求中
		user, err := libs.ParseToken(token)
		if err != nil {
			libs.ErrorLogin(c, "token is invalid")
			c.Abort() // 终止后续处理
			return
		}
		session := sessions.Default(c)

		if routeData.NeedAuth == 1 {

			clientId := c.GetHeader("ClientID")

			if clientId == "" {
				clientId = c.Query("uuid")
				if clientId == "" {
					libs.ErrorLogin(c, "clientId is empty")
					c.Abort() // 终止后续处理
					return
				}
			}
			if !CheckIp(c.Request) {
				libs.Error(c, "the ip is invalid")
				c.Abort() // 终止后续处理
				return
			}
			userData, _ := common.Cache.GetKey("userData", clientId)
			if userData == nil {
				libs.ErrorLogin(c, "the userdata is invalid")
				c.Abort() // 终止后续处理
				return
			}
			m, ok := userData.(map[string]interface{})
			if !ok {
				libs.ErrorLogin(c, "invalid userData type")
				c.Abort() // 终止后续处理
				return
			}
			//log.Printf("userData:%v", m)
			userId, ok := m["userId"].(int32)
			if !ok {
				libs.ErrorLogin(c, "userId is not an int64")
				c.Abort() // 终止后续处理
				return
			}
			if userId != int32(user.ID) || userId == 0 {
				libs.ErrorLogin(c, "user is invalid")
				c.Abort() // 终止后续处理
				return
			}
			c.Set("userId", user.ID)
			//log.Printf("userId:%v", user.ID)
			userRoles, ok := m["userRoles"].(string)
			if !ok {
				libs.Error(c, "user roles is invalid")
				c.Abort() // 终止后续处理
				return
			}
			hasAuth := checkRoles(userRoles, url)
			if !hasAuth {
				libs.Error(c, "user roles check invalid")
				c.Abort() // 终止后续处理
				return
			}
		}
		if routeData.NeedAuth == 2 {
			memberId := session.Get("memberId")
			if memberId == nil {
				libs.ErrorLogin(c, "member is invalid")
				c.Abort() // 终止后续处理
				return
			}

			uid := getUserId(memberId)
			log.Printf("memberId:%v,uid:%v,user.ID:%v", memberId, uid, user.ID)
			if uid != user.ID || uid == 0 {
				log.Printf("member is invalid")
				libs.ErrorLogin(c, "member is invalid")
				c.Abort() // 终止后续处理
				return
			}
			c.Set("memberId", uid)
			if uid > 1 {
				memberRoles := session.Get("memberRoles")
				log.Printf("memberRoles:%v", memberRoles)
				if memberRoles == nil {
					libs.Error(c, "member roles is invalid")
					c.Abort() // 终止后续处理
					return
				}
				memberRole := memberRoles.(string) + ",/index,/member/loginout,/member/welcome,/member/userinfo"
				//log.Printf("adminRules:%v", adminRules)
				hasAuth := checkRoles(memberRole, url)
				if !hasAuth {
					libs.Error(c, "admin roles is invalid")
					c.Abort() // 终止后续处理
					return
				}
			}
		}
		if routeData.NeedAuth == 3 {
			adminId := session.Get("adminId")
			if adminId == nil {
				libs.ErrorLogin(c, "admin is invalid")
				c.Abort() // 终止后续处理
				return
			}

			uid := getUserId(adminId)
			log.Printf("adminId:%v,uid:%v,user.ID:%v", adminId, uid, user.ID)
			if uid != user.ID || uid == 0 {
				log.Printf("admin is invalid")
				libs.ErrorLogin(c, "admin is invalid")
				c.Abort() // 终止后续处理
				return
			}
			c.Set("adminId", uid)
			if uid > 1 {
				adminRoles := session.Get("adminRoles")
				log.Printf("adminRoles:%v", adminRoles)
				if adminRoles == nil {
					libs.Error(c, "admin roles is invalid")
					c.Abort() // 终止后续处理
					return
				}
				adminRules := adminRoles.(string) + ",/index,/admin/loginout,/admin/welcome,/admin/userinfo"
				//log.Printf("adminRules:%v", adminRules)
				hasAuth := checkRoles(adminRules, url)
				if !hasAuth {
					libs.Error(c, "admin roles is invalid")
					c.Abort() // 终止后续处理
					return
				}
			}
		}
		c.Next()
	}
}
func getRouteData(url string) (Route, bool) {
	for path, route := range Routes {
		// 检查 URL 是否匹配动态路由
		if matchPath(url, path) {
			return route, true
		}
	}
	return Route{}, false
}

// matchPath 检查 URL 是否匹配给定的路径模式
func matchPath(url, pattern string) bool {
	parts := strings.Split(pattern, "/")
	urlParts := strings.Split(url, "/")

	if len(parts) != len(urlParts) {
		return false
	}
	//log.Printf("url:%v,pattern:%v", url, pattern)
	if strings.Contains(pattern, ":") {
		for i, part := range parts {
			if part == "" { // 空字符串处理
				continue
			}
			if part[0] == ':' { // 动态参数
				paramName := part[1:]
				if strings.Contains(paramName, "id") { // 数字参数
					if _, err := strconv.Atoi(urlParts[i]); err != nil {
						return false
					}
				}
			} else if part != urlParts[i] { // 静态部分
				return false
			}
		}
	} else {
		if pattern != url {
			return false
		}
	}

	return true
}
func getUserId(userId interface{}) int64 {
	var userIdInt64 int64
	switch v := userId.(type) {
	case int32:
		userIdInt64 = int64(v)
	case int64:
		userIdInt64 = int64(v)
	default: // 终止后续处理
		userIdInt64 = 0
	}
	return userIdInt64
}
func checkRoles(rules string, url string) bool {
	if rules == "*" || rules == "-1" {
		return true
	}
	if rules == "" {
		return false
	}
	userRoleArray := strings.Split(rules, ",")
	hasAuth := false
	for _, path := range userRoleArray {
		if matchPath(url, path) {
			hasAuth = true
			break
		}
	}
	return hasAuth
}

// CheckIp 验证客户端 IP 是否允许访问
func CheckIp(r *http.Request) bool {
	if strings.HasPrefix(r.URL.Scheme, "wails") {
		return true
	}

	clientIP := r.RemoteAddr
	host := r.Host
	ip, _, err := net.SplitHostPort(clientIP)
	if err != nil {
		//fmt.Printf("Error splitting host and port: %v\n", err)
		return false
	}

	// 解析 IP 地址
	ipAddress := net.ParseIP(ip)
	//fmt.Printf("clientIP:%v, host:%v\n", clientIP, host)

	// 允许本机访问
	if ipAddress.IsLoopback() || host == "localhost" {
		return true
	}

	// 获取允许的 IP 和域名列表
	ipAndDomainList := common.Config.System.IpAccess
	//fmt.Printf("ipAndDomainList:%v\n", ipAndDomainList)
	if len(ipAndDomainList) == 0 {
		return true
	}
	// 检查 IP 地址
	for _, allowed := range ipAndDomainList {
		if ipAddress.String() == allowed {
			return true
		}
		// 检查域名
		if host == allowed {
			return true
		}
	}

	return false
}

package login

import (
	"godocms/common"
	"godocms/libs"
	"godocms/pkg/github"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func GetDingConf(c *gin.Context) {
	libs.Success(c, "success", gin.H{
		"id":        common.LoginConf.Ding.CorpId,
		"client_id": common.LoginConf.Ding.AppKey,
		"host":      common.LoginConf.Ding.Host,
	})
}

// handleGithubAuth 处理github授权请求
func HandleGithubAuth(c *gin.Context) {
	var requestData struct {
		RedirectURL string `json:"redirect_url"`
		State       string `json:"state"`
	}
	err := c.ShouldBindJSON(&requestData)
	if err != nil {
		libs.Error(c, "参数错误")
		return
	}

	urlStr := GetGithubAuthUrl(requestData.State, requestData.RedirectURL)

	data := map[string]interface{}{
		"url": urlStr,
	}

	libs.Success(c, "success", data)
}
func GetGithubAuthUrl(state, redirectURL string) string {
	// 构建GitHub登录URL
	oauthConfig := github.OAuthConfig{}
	githubOAuth := github.NewGithubOAuth(oauthConfig)

	url := githubOAuth.GetGithubAuthUrl(state)
	return url
}

// handleGiteeAuth 处理gitee授权请求
func HandleGiteeAuth(c *gin.Context) {
	var requestData struct {
		State       string `json:"state"`
		RedirectURL string `json:"redirect_url"`
	}
	err := c.ShouldBindJSON(&requestData)
	if err != nil {
		libs.Error(c, "参数错误")
		return
	}

	urlStr := GetGiteeAuthUrl(requestData.State, requestData.RedirectURL)

	data := map[string]interface{}{
		"url": urlStr,
	}

	libs.Success(c, "success", data)
}
func GetGiteeAuthUrl(state, redirectURL string) string {
	oauthConfig := &oauth2.Config{
		ClientID:     common.LoginConf.Gitee.ClientID,
		ClientSecret: common.LoginConf.Gitee.ClientSecret,
		Endpoint:     common.LoginConf.Gitee.Endpoint,
		RedirectURL:  redirectURL,
		Scopes:       common.LoginConf.Gitee.Scopes,
	}

	url := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return url
}

// 微信授权
func HandleWechatAuth(c *gin.Context) {
	urlStr := GetWXAuthUrl()
	data := map[string]interface{}{
		"url": urlStr,
	}
	libs.Success(c, "success", data)
}
func GetWXAuthUrl() string {
	// // 生成微信授权URL
	// // authURL := "https://open.weixin.qq.com/connect/qrconnect?"
	// authURL := "https://open.weixin.qq.com/connect/oauth2/authorize?"
	// params := url.Values{}
	// params.Add("appid", "wxea8bf450305b1b57")
	// params.Add("redirect_uri", "http://os.godoos.com")
	// params.Add("response_type", "code")
	// params.Add("scope", "snsapi_base")
	// params.Add("state", "state")

	// urlStr = authURL + "?" + params.Encode() + "#wechat_redirect"

	// // 拼接url
	// urlStr = fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=state#wechat_redirect", "wxea8bf450305b1b57", "http://os.godoos.com")
	return ""
}

// 第三方登录列表
func GetThirdPartyList(c *gin.Context) {
	thirdPartyList := []string{}
	if common.LoginConf.UserNameLogin.Enable {
		thirdPartyList = append(thirdPartyList, "password")
	}
	if common.LoginConf.Github.UserLogin {
		thirdPartyList = append(thirdPartyList, "github")
	}
	if common.LoginConf.Gitee.UserLogin {
		thirdPartyList = append(thirdPartyList, "gitee")
	}
	if common.LoginConf.Ding.UserLogin {
		thirdPartyList = append(thirdPartyList, "dingding")
	}
	if common.LoginConf.Qyweixin.UserLogin {
		thirdPartyList = append(thirdPartyList, "qyweixin")
	}
	if common.LoginConf.Phone.UserLogin {
		thirdPartyList = append(thirdPartyList, "phone")
	}
	if common.LoginConf.Email.UserLogin {
		thirdPartyList = append(thirdPartyList, "email")
	}
	if common.LoginConf.ThirdApi.UserLogin {
		thirdPartyList = append(thirdPartyList, "thirdparty")
	}
	if common.LoginConf.LDAP.UserLogin {
		thirdPartyList = append(thirdPartyList, "ldap")
	}

	libs.Success(c, "success", gin.H{
		"list": thirdPartyList,
	})
}

func GetQYWeixinConf(c *gin.Context) {
	if !common.QYWweiXinValid() {
		libs.Error(c, "企业微信登录未配置")
		return
	}

	libs.Success(c, "success", gin.H{
		"agent_id": common.LoginConf.Qyweixin.AgentId,
		"corp_id":  common.LoginConf.Qyweixin.Corpid,
		"redirect": common.LoginConf.Qyweixin.RedirectURL,
	})
}

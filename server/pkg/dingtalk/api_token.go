package dingtalk

import (
	"godocms/common"
	"godocms/pkg/dingtalk/payload"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cast"
)

// GetAccessToken 获取token
func (ding *DingTalk) GetAccessToken() (token string, err error) {
	if ding.cache != "" {
		data, err := common.Cache.Get(DingDingAccessTokenCachekey)
		if err == nil && data != nil {
			return cast.ToString(data), nil
		}
	}
	args := url.Values{}
	args.Set("appkey", ding.key)
	args.Set("appsecret", ding.secret)
	resp := &payload.AccessToken{}
	if err = ding.Request(http.MethodGet, GetTokenKey, args, nil, resp); err != nil {
		return "", err
	}

	resp.Create = time.Now().Unix()
	common.Cache.Set(DingDingAccessTokenCachekey, resp.Token, 60*60*time.Minute)

	return resp.Token, nil
}

// GetJsApiTicket 获取jsapi_ticket
func (ding *DingTalk) GetJsApiTicket() (ticket string, err error) {
	resp := &payload.JsApiTicket{}

	// get cache

	if err = ding.Request(http.MethodGet, GetJsApiTicketKey, nil, nil, resp); err != nil {
		return "", err
	}
	resp.Create = time.Now().Unix()

	// set cache

	return resp.Ticket, nil
}

// GetUserAccessToken 通过oauth2临时授权码获取用户Token
// https://open.dingtalk.com/document/orgapp/obtain-user-token
func (ding *DingTalk) GetUserAccessToken(authCode string) (res payload.UserAccessTokenResponse, err error) {
	req := payload.NewUserAuthToken(authCode)
	req.ClientId = ding.key
	req.ClientSecret = ding.secret

	err = ding.Request(http.MethodPost, GetUserAccessToken, nil, req, &res)
	return
}

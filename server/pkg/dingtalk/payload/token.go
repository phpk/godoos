package payload

import (
	"fmt"
)

// AccessToken 获取企业内部应用的access_token
// https://developers.dingtalk.com/document/app/obtain-orgapp-token
type AccessToken struct {
	Response
	// Expires 过期时间
	Expires int16 `json:"expires_in"`
	// Create 创建时间
	Create int64 `json:"create"`
	// Token
	Token string `json:"access_token"`
}

// CreatedAt is when the access token is generated
func (token *AccessToken) CreatedAt() int64 {
	return token.Create
}

// ExpiresIn is how soon the access token is expired
func (token *AccessToken) ExpiresIn() int16 {
	return token.Expires
}

type JsApiTicket struct {
	Response
	// Expires 过期时间
	Expires int16 `json:"expires_in"`
	// Create 创建时间
	Create int64 `json:"create"`
	// Ticket 生成的临时jsapi_ticket。
	Ticket string `json:"ticket"`
}

// CreatedAt is when the access token is generated
func (token *JsApiTicket) CreatedAt() int64 {
	return token.Create
}

// ExpiresIn is how soon the access token is expired
func (token *JsApiTicket) ExpiresIn() int16 {
	return token.Expires
}

const (
	ActionGetUserAuthToken = "authorization_code"
	ActionRefreshAuthToken = "refresh_token"
)

type UserAccessToken struct {
	ClientId     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
	Code         string `json:"code,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	GrantType    string `json:"grantType,omitempty"`
}

func NewUserAuthToken(authCode string) *UserAccessToken {
	return &UserAccessToken{
		Code:      authCode,
		GrantType: ActionGetUserAuthToken,
	}
}

type UserAccessTokenResponse struct {
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	ExpireIn     int    `json:"expireIn,omitempty"`
	CorpId       string `json:"corpId,omitempty"`

	Code      string `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	RequestId string `json:"requestid,omitempty"`
}

func (u UserAccessTokenResponse) CheckError() (err error) {
	if u.Code != "" {
		err = fmt.Errorf("code:'%s', msg:%s, requestId: %s", u.Code, u.Message, u.RequestId)
	}

	return
}

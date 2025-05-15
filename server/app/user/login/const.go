package login

import (
	"encoding/json"
	"errors"
	"godocms/model"
)

const (
	LoginTypeSmsCode       = "sms_code"           // 手机短信验证码登录
	LoginTypePassword      = "password"           // 用户名密码登录
	LoginTypeDingWorkBeach = "dingtalk_workbench" // 钉钉工作台登录
	LoginTypeDingScan      = "dingtalk_scan"      // 钉钉扫码登录
	LoginTypeQYWeixinWork  = "qyweixin"           // 企业微信登录
	LoginTypeQYWeixinScan  = "qyweixin_scan"      // 企业微信扫码登录
	LoginTypeGithub        = "github"             // gitee登录
	LoginTypeGitee         = "gitee"              // gitee登录
	LoginTypeEmail         = "email"              // 邮箱验证码登录
	LoginTypeThirdPartyAPI = "third_api"          // 第三方api登录
	LoginTypeMicrosoftADFS = "microsoft_adfs"     // 微软 Active Directory Federation Services 登录
	LoginTypeMicrosoftOIDC = "microsoft_oidc"
	LoginTypeLDAP          = "ldap"
)

const (
	UserNoMobileErrCode = 10000
)

type LoginParam interface{}

type LoginHandler interface {
	Init(req LoginRequest) error
	Login() (*model.User, error)
	Register() (*model.User, error)
}

var ErrNeedRegister = errors.New("need register")

type LoginRequest struct {
	ClientId  string          `json:"client_id" binding:"required"`
	LoginType string          `json:"login_type" binding:"required"`
	Action    string          `json:"action" binding:"required"`
	Param     json.RawMessage `json:"param" binding:"required"`
}

const (
	LoginPlatformDingTalk      = "dingtalk" // 钉钉
	LoginPlatformWeixinWork    = "wxqy"     // 企业微信
	LoginPlatformGithub        = "github"   // github
	LoginPlatformGitee         = "gitee"    // gitee
	LoginPlatformThirdPartyAPI = "thirdapi" // 第三方API
	LoginPlatformMicrosoftADFS = "microsoft_adfs"
	LoginPlatformMicrosoftOIDC = "microsoft_oidc"
	LoginPlatformLDAP          = "ldap"
	LoginPlatformPassword      = "password" // 密码登录
)

// 密码登录参数
type PasswordLoginParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 短信验证码登录参数
type SmsCodeLoginParam struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}
type LoginCache struct {
	UserID    int32  `json:"user_id"`
	UserRoles string `json:"user_roles"`
}

// 企业微信扫码登录参数
type WeixinScanLoginParam struct {
	OpenID string `json:"open_id"`
	Token  string `json:"token"`
}

// 第三方 API 登录参数
type ThirdPartyLoginParam struct {
	ThirdID string `json:"third_id"`
	Token   string `json:"token"`
}

type LoginResponse struct {
	User       UserResponse `json:"user"`
	Role       RoleReponse  `json:"role"`
	Dept       DeptResponse `json:"dept"`
	ClientID   string       `json:"client_id"`
	Token      string       `json:"token"`
	ISPwd      bool         `json:"isPwd"`
	UserAuths  []model.User `json:"user_auths"`
	UserShares []model.User `json:"user_shares"`
}

type UserResponse struct {
	ID        int32  `json:"id"`
	Name      string `json:"name"`
	Nickname  string `json:"nickname"`
	Sex       int32  `json:"sex"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	Phone     string `json:"phone"`
	Desc      string `json:"desc"`
	JobNumber string `json:"job_number"`
	WorkPlace string `json:"work_place"`
	HiredDate int64  `json:"hired_date"`
	UseSpace  int32  `json:"use_space"`
	HasSpace  int32  `json:"has_space"`
}

type RoleReponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type DeptResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

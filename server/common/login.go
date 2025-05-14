package common

import (
	"golang.org/x/oauth2"
)

var LoginConf LoginInfo

type LoginInfo struct {
	Email              Email         `json:"email"`
	Ding               Ding          `json:"ding"`
	Qyweixin           Qyweixin      `json:"qyweixin"`
	AllowFileLoginConf []string      `json:"allowed_file_types"`
	Github             Github        `json:"github"`
	Gitee              Gitee         `json:"gitee"`
	ThirdApi           ThirdApi      `json:"thirdApi"`
	Phone              Phone         `json:"phone"`
	UserNameLogin      UserNameLogin `json:"usernameLogin"`
	MicrosoftADFS      MicrosoftADFS `json:"microsoft_adfs"`
	MicrosoftOIDC      MicrosoftOIDC `json:"microsoft_oidc"`
	LDAP               LDAP          `json:"ldap"`
}

type Email struct {
	From         string `json:"from"`
	Host         string `json:"host"`
	Password     string `json:"password"`
	Username     string `json:"username"`
	Port         int    `json:"port"`
	IsSsl        bool   `json:"isSsl"`
	UserLogin    bool   `json:"userLogin"`
	AutoRegister bool   `json:"autoRegister"`
	Enable       bool   `json:"enable"`
	AdminLogin   bool   `json:"adminLogin"`
}
type Ding struct {
	AgentId      string `json:"agentId"`
	MiniAppId    string `json:"miniAppId"`
	AppKey       string `json:"appKey"`
	AppSecret    string `json:"appSecret"`
	CorpId       string `json:"corpId"`
	ApiToken     string `json:"apiToken"`
	Host         string `json:"host"`
	UserLogin    bool   `json:"userLogin"`
	AutoRegister bool   `json:"autoRegister"`
	Enable       bool   `json:"enable"`
}
type Qyweixin struct {
	Corpid         string `json:"corpid"`
	AgentId        string `json:"agentId"`
	Secret         string `json:"secret"`
	ContactsSecret string `json:"contacts_secret"`
	RedirectURL    string `json:"redirect_url"`
	UserLogin      bool   `json:"userLogin"`
	AutoRegister   bool   `json:"autoRegister"`
	Enable         bool   `json:"enable"`
}
type Github struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURL  string   `json:"redirect_url"`
	Scopes       []string `json:"scopes"`
	UserLogin    bool     `json:"userLogin"`
	AutoRegister bool     `json:"autoRegister"`
	Enable       bool     `json:"enable"`
}
type Gitee struct {
	ClientID     string          `json:"client_id"`
	ClientSecret string          `json:"client_secret"`
	RedirectURL  string          `json:"redirect_url"`
	Endpoint     oauth2.Endpoint `json:"endpoint"`
	Scopes       []string        `json:"scopes"`
	UserLogin    bool            `json:"userLogin"`
	AutoRegister bool            `json:"autoRegister"`
	Enable       bool            `json:"enable"`
}

type Phone struct {
	ServiceType   string     `json:"serviceType"`    // 服务类型
	UniversalCode string     `json:"universal_code"` // 万能码
	UserLogin     bool       `json:"userLogin"`      // 用户登录
	Enable        bool       `json:"enable"`         // 是否启用
	AutoRegister  bool       `json:"autoRegister"`   // 用户注册
	AliyunSms     AliyunSms  `json:"aliyunSms"`      // 阿里云短信
	TencentSms    TencentSms `json:"tencentSms"`     // 腾讯云短信
}

// 阿里云短信登录
type AliyunSms struct {
	AccessKeyId     string `json:"access_key_id"`     // 阿里云访问密钥 ID
	AccessKeySecret string `json:"access_key_secret"` // 阿里云访问密钥 Secret
	SignName        string `json:"sign_name"`         // 短信签名名称
	TemplateCode    string `json:"template_code"`     // 短信模板 ID
	TemplateParam   string `json:"template_param"`    // 模板变量内容，JSON 格式字符串 // "{\"code\":\"123456\"}"
	Enable          bool   `json:"enable"`            // 是否启用
}

// 腾讯云短信登录
type TencentSms struct {
	SecretId       string `json:"secret_id"`
	SecretKey      string `json:"secret_key"`
	Region         string `json:"region"`
	SmsSdkAppId    string `json:"sms_sdk_app_id"`
	SignName       string `json:"sign_name"`
	TemplateId     string `json:"template_id"`
	PhoneNumberSet string `json:"phone_number_set"`
	Enable         bool   `json:"enable"` // 是否启用
}

type ThirdApi struct {
	Api          string `json:"api"`
	Secret       string `json:"secret"`
	Enable       bool   `json:"enable"`       // 是否启用
	AutoRegister bool   `json:"autoRegister"` // 是否自动注册
	UserLogin    bool   `json:"userLogin"`
}

type MicrosoftADFS struct {
	ClientID     string          `json:"client_id"`
	ClientSecret string          `json:"client_secret"`
	RedirectURL  string          `json:"redirect_url"`
	TenantID     string          `json:"tenant_id"`
	Endpoint     oauth2.Endpoint `json:"-"`
	Scopes       []string        `json:"scopes"`
	UserLogin    bool            `json:"userLogin"`
	AutoRegister bool            `json:"autoRegister"`
	Enable       bool            `json:"enable"`
}

type MicrosoftOIDC struct {
	ClientID     string          `json:"client_id"`
	ClientSecret string          `json:"client_secret"`
	RedirectURL  string          `json:"redirect_url"`
	TenantID     string          `json:"tenant_id"`
	Endpoint     oauth2.Endpoint `json:"-"`
	Scopes       []string        `json:"scopes"`
	UserLogin    bool            `json:"userLogin"`
	AutoRegister bool            `json:"autoRegister"`
	Enable       bool            `json:"enable"`
}

type LDAP struct {
	Server       string `json:"server"`       // ldap.example.com
	Port         string `json:"port"`         // 389
	BaseDN       string `json:"baseDN"`       // dc=example,dc=com
	BindDN       string `json:"bindDN"`       // cn=admin,dc=example,dc=com
	BindPassword string `json:"bindPassword"` // admin_password
	TLS          string `json:"tls"`          // ca证书地址,填写使用ldaps连接
	UserLogin    bool   `json:"userLogin"`
	AutoRegister bool   `json:"autoRegister"`
	Enable       bool   `json:"enable"`
}

type Editor struct {
	OnlyOffice string `json:"onlyoffice"`
}

// 用户名登录
type UserNameLogin struct {
	Enable       bool `json:"enable"`       // 是否启用
	AutoRegister bool `json:"autoRegister"` // 是否自动注册
}

func DingTalkIsValid() bool {
	if LoginConf.Ding.AgentId == "" || LoginConf.Ding.AppKey == "" || LoginConf.Ding.AppSecret == "" || LoginConf.Ding.CorpId == "" || LoginConf.Ding.Host == "" {
		return false
	}

	return true
}

func AliyunSmsIsValid() bool {
	if LoginConf.Phone.AliyunSms.AccessKeyId == "" || LoginConf.Phone.AliyunSms.AccessKeySecret == "" {
		return false
	}

	if !LoginConf.Phone.Enable {
		return false
	}

	return true
}

func EmailIsValid() bool {
	if LoginConf.Email.From == "" || LoginConf.Email.Host == "" || LoginConf.Email.Password == "" || LoginConf.Email.Username == "" || LoginConf.Email.Port == 0 {
		return false
	}

	return true
}

func QYWweiXinValid() bool {
	if LoginConf.Qyweixin.AgentId == "" || LoginConf.Qyweixin.Corpid == "" || LoginConf.Qyweixin.Secret == "" || LoginConf.Qyweixin.ContactsSecret == "" || LoginConf.Qyweixin.RedirectURL == "" {
		return false
	}

	return true
}

func ThirdApiValid() bool {
	if LoginConf.ThirdApi.Api == "" || LoginConf.ThirdApi.Secret == "" {
		return false
	}

	return true
}

func MicrosoftADFSValid() bool {
	if LoginConf.MicrosoftADFS.ClientID == "" || LoginConf.MicrosoftADFS.ClientSecret == "" || LoginConf.MicrosoftADFS.RedirectURL == "" || LoginConf.MicrosoftADFS.TenantID == "" {
		return false
	}

	return true
}

func MicrosoftOIDCValid() bool {
	if LoginConf.MicrosoftOIDC.ClientID == "" || LoginConf.MicrosoftOIDC.ClientSecret == "" || LoginConf.MicrosoftOIDC.RedirectURL == "" {
		return false
	}

	return true
}

func LDAPValid() bool {
	if LoginConf.LDAP.Server == "" || LoginConf.LDAP.Port == "" || LoginConf.LDAP.BaseDN == "" || LoginConf.LDAP.BindDN == "" || LoginConf.LDAP.BindPassword == "" {
		return false
	}

	return true
}

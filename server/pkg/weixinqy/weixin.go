package weixinqy

import (
	"errors"
	"fmt"
	"godocms/common"
	"godocms/utils"
	"log/slog"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
)

// GetWxQiyeSelfAccessToken 获取企业自建应用微信access_token。
// 该函数通过调用微信接口获取企业微信的access_token，用于后续的API调用。
// 参数:
//
//	corpid: 企业号的ID
//	appSecret: 企业号的某个应用的密钥
//
// 返回值:
//
//	token: 微信Corp Access Token字符串。
//	err: 错误对象，如果在获取Access Token过程中发生错误，则返回该错误。
func GetWxQiyeSelfAccessToken(corpid string, appSecret string) (token string, err error) {
	tokenKey := fmt.Sprintf(cache_wxqiyeself_token+":%s", appSecret)
	v, err := common.GetCache(tokenKey)
	if err == nil {
		return cast.ToString(v), nil
	}
	var result WxQiyeAccessTokenResp
	client := resty.New()
	url := fmt.Sprintf(qiye_access_token_url+"?corpid=%s&corpsecret=%s", corpid, appSecret)
	_, err = client.R().SetResult(&result).Get(url)
	if err != nil {
		return "", err
	}
	if result.Errcode == 0 && result.ExpiresIn > 0 {
		expires_in := cast.ToInt(result.ExpiresIn/60) - 1 //分钟
		common.SetCache(tokenKey, result.AccessToken, time.Duration(expires_in)*time.Minute)
		return result.AccessToken, nil
	}
	return "", errors.New(result.Errmsg)
}

// GetWxQiyeUserId 网页授权，根据code获取访问用户身份获取userid,userTicket
// 参数:
//  如果是二维码生成的code ， 用户凭据userTicket为null
//	accessToken: 企业的token
//	code: 回调返回的code

//	企业微信那边成功将返回{
//		"errcode": 0,
//		"errmsg": "ok",
//		"userid":"USERID",
//		"user_ticket": "USER_TICKET"
//	 }
func GetWxQiyeUserId(accessToken string, code string) (userid string, userTicket string, err error) {
	client := resty.New()
	url := fmt.Sprintf(Qiye_getuserinfo_url+"?access_token=%s&code=%s", accessToken, code)
	var result WxQiyeUserInfo
	_, err = client.R().SetResult(&result).Get(url)
	if err != nil {
		slog.Error("获取企业微信用户信息失败", "err", err.Error())
		return
	}
	if result.Errcode != 0 {
		err = errors.New(result.Errmsg)
		return
	}
	return result.Userid, result.UserTicket, nil
}

// QyweixinAutoLogin 为企业微信自动登录功能生成Token。
// 参数:
//
//	uid - user表中的id。
//	clientId - 客户端标识符，用于区分不同的客户端。
//
// 返回值:
//
//	token - 生成的Token字符串，用于用户身份验证。
//	err - 如果生成Token过程中出现错误，返回该错误。
func QyweixinAutoLogin(uid int64, clientId string) (token string, err error) {
	//todo 有绑定过用户,参考auth.go中的handleLogin方法
	token, err = utils.GenerateToken(&utils.UserClaims{
		ID: uid,
	})
	//todo
	// cacheData := map[string]interface{}{
	// 	"userId":    user.ID,
	// 	"userRoles": userRole.Rules,
	// }
	// utils.SetCacheById("userData", cacheData, clientId, 24*60)
	if err != nil {
		return "", errors.New("自动生成token失败:" + err.Error())
	}
	return
}

// GetWxQiyeUserInfo 根据 accessToken 和 userTicket 获取企业微信用户详细信息
// 参数:
//
//	accessToken string: 企业微信的访问令牌
//	userTicket string: 用户的票据，用于验证用户身份
//
// 返回值:
//
//	userInfo WxQiyeUserInfoMore: 用户信息的结构体，包含详细信息
//	err error: 错误信息，如果调用过程中发生错误
func GetWxQiyeUserInfo(accessToken string, userTicket string) (userInfo WxQiyeUserInfoMore, err error) {
	moreUrl := fmt.Sprintf(Qiye_getuserdetail_url+"?access_token=%s", accessToken)
	client := resty.New()
	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(fmt.Sprintf(`{"user_ticket": "%s"}`, userTicket)).
		SetResult(&userInfo).
		Post(moreUrl)
	if err != nil {
		slog.Error("网页授权登录获取用户详情", "err", err.Error())
		return
	}
	if userInfo.Errcode != 0 {
		err = errors.New(userInfo.Errmsg)
		return
	}

	return
}

// ConvertUserIdToOpenId 根据 accessToken 和 userid 获取企业微信用户的openid
// 参数:
//
//	accessToken string: 企业微信的访问令牌
//	userid string: 用户的id
//
// 返回值:
//
//	openidInfo WxConvertUserIdToOpenIdResp: openid的结构体，包含详细信息
//	err error: 错误信息，如果调用过程中发生错误
func ConvertUserIdToOpenId(accessToken string, userId string) (openidInfo WxConvertUserIdToOpenIdResp, err error) {
	url := fmt.Sprintf(Qiye_convert_userid_url+"?access_token=%s", accessToken)
	client := resty.New()
	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(fmt.Sprintf(`{"userid": "%s"}`, userId)).
		SetResult(&openidInfo).
		Post(url)
	if err != nil {
		slog.Error("wechat work user id convert openid fail", "err", err.Error())
		return
	}
	if openidInfo.Errcode != 0 {
		err = errors.New(openidInfo.Errmsg)
		return
	}

	return
}

// GetWxWorkUserInfo 根据 accessToken 和 userid 获取企业微信用户详细信息
// 参数:
//
//	accessToken string: 企业微信的访问令牌
//	userid string: 用户的id
//
// 返回值:
//
//	userInfo WxWorkUserInfo: user包含详细信息
//	err error: 错误信息，如果调用过程中发生错误
func GetWxWorkUserInfo(accessToken, userId string) (userInfo WxWorkUserInfo, err error) {
	// accessToken, err := GetWxQiyeSelfAccessToken(common.LoginConf.Qyweixin.Corpid, common.LoginConf.Qyweixin.ContactsSecret)
	// if err != nil {
	// 	return
	// }
	url := fmt.Sprintf(Qiye_getwxwork_userinfo_url+"?access_token=%s&userid=%s", accessToken, userId)
	client := resty.New()
	_, err = client.R().SetResult(&userInfo).Get(url)
	if err != nil {
		slog.Error("获取企业微信用户详情", "err", err.Error())
		return
	}
	if userInfo.ErrCode != 0 {
		err = errors.New(userInfo.ErrMsg)
		return
	}

	return
}

// GetDepartmentList 根据 accessToken 和 id 获取企业微信部门列表信息
// 参数:
//
//	accessToken string: 企业微信的访问令牌
//	id int: 部门的id， 传0则查全部
//
// 返回值:
//
//	deptResp DepartmentListResp: 部门列表信息
//	err error: 错误信息，如果调用过程中发生错误
//
// https://qyapi.weixin.qq.com/cgi-bin/department/list?access_token=ACCESS_TOKEN&id=ID
func GetDepartmentList(accessToken string, id int) (deptResp DepartmentListResp, err error) {
	apiURL := fmt.Sprintf(QiyeGetDepartmentListURL, accessToken)
	if id != 0 {
		apiURL = fmt.Sprintf(QiyeGetDepartmentListURL+"&id=%d", accessToken, id)
	}
	client := resty.New()
	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&deptResp).
		Post(apiURL)
	if err != nil {
		return
	}
	if deptResp.ErrCode != 0 {
		err = errors.New(deptResp.ErrMsg)
		return
	}

	return
}

// GetUserList 根据 cursor 和 limit 获取成员id列表
// 参数:
//
//	cursor string: 查询游标
//	limit int: 一页数量
//
// 返回值:
//
//	resp UserListResp: 成员id列表信息
//	err error: 错误信息，如果调用过程中发生错误
func GetUserList(cursor string, limit int) (resp UserListResp, err error) {
	// 需要用到通讯录密钥
	accessToken, err := GetWxQiyeSelfAccessToken(common.LoginConf.Qyweixin.Corpid, common.LoginConf.Qyweixin.ContactsSecret)
	if err != nil {
		return
	}

	apiURL := fmt.Sprintf(QiyeGetUserList, accessToken)
	client := resty.New()
	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(fmt.Sprintf(`{"cursor": "%s", "limit": "%d"}`, cursor, limit)).
		SetResult(&resp).
		Post(apiURL)
	if err != nil {
		return
	}
	if resp.ErrCode != 0 {
		err = errors.New(resp.ErrMsg)
		return
	}

	return
}

func GetQYWxQRCode(status string) string {
	if status == "" {
		status = "WWLogin" // STATE
	}

	qrCodeUrl := fmt.Sprintf("https://login.work.weixin.qq.com/wwlogin/sso/login?login_type=%s&appid=%s&agentid=%s&redirect_uri=%s&state=%s", "CorpApp", common.LoginConf.Qyweixin.Corpid, common.LoginConf.Qyweixin.AgentId, common.LoginConf.Qyweixin.RedirectURL, status)
	return qrCodeUrl
}

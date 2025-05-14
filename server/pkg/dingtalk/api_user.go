package dingtalk

import (
	"fmt"
	"godocms/pkg/dingtalk/payload"
	"net/http"
	"net/url"
)

// 根据手机号查询用户
func (ding *DingTalk) GetUserIdByMobile(mobile string) (resp payload.MobileGetUserIdResp, err error) {
	return resp, ding.Request(http.MethodPost, GetUserIdByMobileKey, nil, payload.NewMobileGetUserIdReq(mobile), &resp)
}

// 获取jsapi_ticket => JSAPI(返回前端鉴权)
// JSAPI => auth_code(传递后端)
// auth_code + token => userid
// GetUserInfoByCode 通过免登码和access_token获取用户信息
func (ding *DingTalk) GetUserInfoByCode(code string) (resp payload.CodeGetUserInfo, err error) {
	return resp, ding.Request(http.MethodPost, GetUserInfoByCodeKey, nil,
		payload.NewCodeGetUserInfoReq(code), &resp)
}

// GetUserDetail 根据userid获取用户详情
func (ding *DingTalk) GetUserDetail(user *payload.UserDetailReq) (resp payload.UserDetailResp, err error) {
	return resp, ding.Request(http.MethodPost, GetUserDetailKey, nil, user, &resp)
}

// GetCurrentUserByAccessToken 通UserAccessToken获取当前授权人的信息
func (ding *DingTalk) GetCurrentUserByAccessToken(userAccessToken string) (rsp payload.ContactUser, err error) {
	query := url.Values{}
	query.Set("access_token", userAccessToken)
	return rsp, ding.Request(http.MethodGet, fmt.Sprintf(GetContactUser, "me"), query, nil, &rsp)
}

// GetContactUser 获取用户通讯录个人信息
// 调用本接口获取企业用户通讯录中的个人信息。
// @see https://open.dingtalk.com/document/orgapp/dingtalk-retrieve-user-information?spm=ding_open_doc.document.0.0.58b9492dZxH66f
func (ding *DingTalk) GetContactUser(unionId string) (rsp payload.ContactUser, err error) {
	return rsp, ding.Request(http.MethodGet, fmt.Sprintf(GetContactUser, unionId), nil, nil, &rsp)
}

// GetUserIdByUnionId 根据unionid获取用户userid
func (ding *DingTalk) GetUserIdByUnionId(res *payload.UnionIdGetUserIdReq) (req payload.UnionIdGetUserIdResponse, err error) {
	return req, ding.Request(http.MethodPost, GetUserIdByUnionIdKey, nil, res, &req)
}

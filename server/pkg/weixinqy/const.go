package weixinqy

const cache_wxqiyeself_token = "wxqiyeself_token"                                            //企业微信自建应该，token缓存key
const qiye_access_token_url = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"                 //获取 access_token的地址
const Qiye_getuserinfo_url = "https://qyapi.weixin.qq.com/cgi-bin/auth/getuserinfo"          //网页授权登录,获取访问用户身份,获取 USERID，USER_TICKET
const Qiye_getuserdetail_url = "https://qyapi.weixin.qq.com/cgi-bin/auth/getuserdetail"      //网页授权登录,获获取访问用户敏感信息
const Qiye_convert_userid_url = "https://qyapi.weixin.qq.com/cgi-bin/user/convert_to_openid" // userid转openid
const Qiye_getwxwork_userinfo_url = "https://qyapi.weixin.qq.com/cgi-bin/user/get"           // 读取成员, oauth2手工授权

const QiyeGetDepartmentListURL = "https://qyapi.weixin.qq.com/cgi-bin/department/list?access_token=%s" // 获取部门列表
const QiyeGetUserList = "https://qyapi.weixin.qq.com/cgi-bin/user/list_id?access_token=%s"             // 获取成员列表

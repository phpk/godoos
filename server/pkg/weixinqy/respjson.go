package weixinqy

// 企业微信获取token的响应
type WxQiyeAccessTokenResp struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// 企业微信通过回调code获取用户的基本id和ticket
type WxQiyeUserInfo struct {
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
	Userid     string `json:"userid"`
	UserTicket string `json:"user_ticket"`
}

// 根据userTicket获取用户详细信息
type WxQiyeUserInfoMore struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Userid  string `json:"userid"`
	Gender  string `json:"gender"`
	Avatar  string `json:"avatar"`
	QrCode  string `json:"qr_code"`
	Mobile  string `json:"mobile"`
	Email   string `json:"email"`
	BizMail string `json:"biz_mail"`
	Address string `json:"address"`
}

type WxConvertUserIdToOpenIdResp struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Openid  string `json:"openid"`
}

type WxWorkUserInfo struct {
	ErrCode          int             `json:"errcode"`
	ErrMsg           string          `json:"errmsg"`
	UserID           string          `json:"userid"`     // 成员UserID。对应管理端的账号，企业内必须唯一。不区分大小写，长度为1~64个字节；第三方应用返回的值为open_userid
	Name             string          `json:"name"`       // 成员名称；第三方不可获取，调用时返回userid以代替name；代开发自建应用需要管理员授权才返回；对于非第三方创建的成员，第三方通讯录应用也不可获取；未返回name的情况需要通过通讯录展示组件来展示名字
	Department       []int           `json:"department"` // 成员所属部门id列表，仅返回该应用有查看权限的部门id；
	Order            []int           `json:"order"`      // 部门内的排序值，默认为0。数量必须和department一致，数值越大排序越前面
	Position         string          `json:"position"`   // 职务信息；代开发自建应用需要管理员授权才返回；第三方仅通讯录应用可获取；
	Mobile           string          `json:"mobile"`     // 手机号码，代开发自建应用需要管理员授权且成员oauth2授权获取；第三方仅通讯录应用可获取；对于非第三方创建的成员，第三方通讯录应用也不可获取；上游企业不可获取下游企业成员该字段
	Gender           string          `json:"gender"`     // 性别。0表示未定义，1表示男性，2表示女性
	Email            string          `json:"email"`
	BizMail          string          `json:"biz_mail"`          // 企业邮箱
	IsLeaderInDept   []int           `json:"is_leader_in_dept"` // 表示在所在的部门内是否为部门负责人，数量与department一致；
	DirectLeader     []string        `json:"direct_leader"`     // 直属上级UserID，返回在应用可见范围内的直属上级列表，最多有1个直属上级
	Avatar           string          `json:"avatar"`
	ThumbAvatar      string          `json:"thumb_avatar"`
	Telephone        string          `json:"telephone"`
	Alias            string          `json:"alias"`
	Address          string          `json:"address"`
	OpenUserID       string          `json:"open_userid"`     // 全局唯一。对于同一个服务商，不同应用获取到企业内同一个成员的open_userid是相同的
	MainDepartment   int             `json:"main_department"` // 主部门，仅当应用对主部门有查看权限时返回。
	ExtAttr          ExtAttr         `json:"extattr"`         // 扩展属性
	Status           int             `json:"status"`          // 激活状态: 1=已激活，2=已禁用，4=未激活，5=退出企业。
	QRCode           string          `json:"qr_code"`
	ExternalPosition string          `json:"external_position"`
	ExternalProfile  ExternalProfile `json:"external_profile"`
}

type ExtAttr struct {
	Attrs []Attr `json:"attrs"`
}

type Attr struct {
	Type int       `json:"type"`
	Name string    `json:"name"`
	Text *TextAttr `json:"text,omitempty"`
	Web  *WebAttr  `json:"web,omitempty"`
}

type TextAttr struct {
	Value string `json:"value"`
}

type WebAttr struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

type ExternalProfile struct {
	ExternalCorpName string         `json:"external_corp_name"`
	WechatChannels   WechatChannels `json:"wechat_channels"`
	ExternalAttr     []ExternalAttr `json:"external_attr"`
}

type WechatChannels struct {
	Nickname string `json:"nickname"`
	Status   int    `json:"status"`
}

type ExternalAttr struct {
	Type        int              `json:"type"`
	Name        string           `json:"name"`
	Text        *TextAttr        `json:"text,omitempty"`
	Web         *WebAttr         `json:"web,omitempty"`
	MiniProgram *MiniProgramAttr `json:"miniprogram,omitempty"`
}

type MiniProgramAttr struct {
	AppID    string `json:"appid"`
	PagePath string `json:"pagepath"`
	Title    string `json:"title"`
}

// 部门列表返回体
type DepartmentListResp struct {
	ErrCode     int          `json:"errcode"`
	ErrMsg      string       `json:"errmsg"`
	Departments []Department `json:"department"`
}

type Department struct {
	ID               int      `json:"id"`
	Name             string   `json:"name"`
	NameEn           string   `json:"name_en"`
	DepartmentLeader []string `json:"department_leader"`
	ParentID         int      `json:"parentid"`
	Order            int      `json:"order"`
}

type DeptUser struct {
	OpenUserID string `json:"userid"`
	Department int    `json:"department"`
}

// 成员列表返回体
type UserListResp struct {
	ErrCode    int        `json:"errcode"`
	ErrMsg     string     `json:"errmsg"`
	NextCursor string     `json:"next_cursor"`
	DeptUser   []DeptUser `json:"dept_user"`
}

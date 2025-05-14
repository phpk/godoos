package payload

// 通过手机号查询用户id请求参数
type MobileGetUserIdReq struct {
	Mobile string `json:"mobile" validate:"required"`
}

func NewMobileGetUserIdReq(mobile string) *MobileGetUserIdReq {
	return &MobileGetUserIdReq{mobile}
}

// 通过手机号查询用户id返回结果
type MobileGetUserIdResp struct {
	Response

	Result struct {
		UserId string `json:"userid"`
	} `json:"result"`
}

// 免登内部应用获取用户信息请求参数
type CodeGetUserInfoReq struct {
	Code string `json:"code,omitempty"  validate:"required"`
}

func NewCodeGetUserInfoReq(code string) *CodeGetUserInfoReq {
	return &CodeGetUserInfoReq{code}
}

// 免登内部应用获取用户信息返回结果
type CodeGetUserInfo struct {
	Response

	UserInfo struct {
		Name              string `json:"name"`               // 用户名字
		UnionId           string `json:"unionid"`            // 用户unionId
		UserId            string `json:"userid"`             // 用户的userid
		AssociatedUnionId string `json:"associated_unionid"` // 用户关联的unionId （用户在钉钉生态系统中的特殊标识）
		Level             int    `json:"sys_level"`          // 级别 1：主管理员，2：子管理员，100：老板，0：其他（如普通员工）
		Admin             bool   `json:"sys"`                // 是否是管理员
		DeviceId          string `json:"device_id"`          // 设备ID
	} `json:"result"`
}

type UserDetailReq struct {
	UserId string `json:"userid" validate:"required"`

	// 通讯录语言，默认zh_CN。如果是英文，请传入en_US。
	Language string `json:"language,omitempty" validate:"omitempty,oneof=zh_CN en_US"`
}

func NewUserDetailReq(userId string, language string) *UserDetailReq {
	return &UserDetailReq{UserId: userId, Language: language}
}

type UserDetailResp struct {
	Response
	UserInfoDetail `json:"result"`
}

type UserInfoDetail struct {
	UserId               string          `json:"userid"`
	UnionId              string          `json:"unionid"`                                                             // 员工在当前开发者企业账号范围内的唯一标识
	Name                 string          `json:"name"`                                                                // 员工名称
	Avatar               string          `json:"avatar"`                                                              // 头像
	StateCode            string          `json:"state_code"`                                                          // 国际电话区号
	ManagerUserId        string          `json:"manager_userid"`                                                      // 员工的直属主管
	Mobile               string          `json:"mobile"`                                                              // 手机号码
	HideMobile           bool            `json:"hide_mobile"`                                                         // 是否号码隐藏
	Telephone            string          `json:"telephone"`                                                           // 分机号
	JobNumber            string          `json:"job_number"`                                                          // 员工工号
	Title                string          `json:"title"`                                                               // 职位
	Email                string          `json:"email,omitempty"`                                                     // 员工邮箱 (需要开通对应权限才会返回)
	WorkPlace            string          `json:"work_place"`                                                          // 办公地点
	Remark               string          `json:"remark"`                                                              // 备注
	LoginId              string          `json:"loginId"`                                                             // 专属帐号登录名
	ExclusiveAccount     bool            `json:"exclusive_account"`                                                   // 是否专属帐号
	ExclusiveAccountType string          `json:"exclusive_account_type"`                                              //专属帐号类型： sso：企业自建专属帐号 dingtalk：钉钉自建专属帐号
	DeptIds              []int           `json:"dept_id_list"`                                                        // 所属部门ID列表
	DeptOrders           []DeptOrderResp `json:"dept_order_list"`                                                     // 员工在对应的部门中的排序。
	Extension            string          `json:"extension,omitempty" validate:"omitempty"`                            // 员工在对应的部门中的排序
	HiredDate            int             `json:"hired_date"`                                                          // 入职时间
	Active               bool            `json:"active"`                                                              // 是否激活了钉钉
	RealAuthed           bool            `json:"real_authed"`                                                         // 是否完成了实名认证
	OrgEmail             string          `json:"org_email,omitempty" validate:"omitempty,max=100"`                    // 员工的企业邮箱 如果员工的企业邮箱没有开通，返回信息中不包含该数据
	OrgEmailType         string          `json:"org_email_type,omitempty" validate:"omitempty,oneof=profession base"` // 员工的企业邮箱类型
	Senior               bool            `json:"senior"`                                                              // 是否为企业的高管
	Admin                bool            `json:"admin"`                                                               // 是否为企业的管理员
	Boss                 bool            `json:"boss"`                                                                // 是否为企业的老板
	LeaderInDept         []LeaderInDept  `json:"leader_in_dept"`                                                      // 员工在对应的部门中是否领导
	// 角色列表
	UserRoles []struct {
		Id        int    `json:"id"`         // 角色id
		Name      string `json:"name"`       // 角色名称
		GroupName string `json:"group_name"` // 角色组名称
	} `json:"role_list"`
	UnionOrg `json:"union_emp_ext"` // 当用户来自于关联组织时的关联信息
}

type LeaderInDept struct {
	DeptId int  `json:"dept_id"`
	Leader bool `json:"leader"`
}

type UnionOrg struct {
	UserId string `json:"userid"`  // 员工的userid
	CorpId string `json:"corp_id"` // 当前用户所属的组织的企业corpId
}

// AssociatedOrg 关联映射关系
type AssociatedOrg struct {
	UnionOrgList []UnionOrg `json:"union_emp_map_list"`
}

// 需要用户toekn
type ContactUser struct {
	Response

	Nick      string `json:"nick,omitempty"`
	AvatarUrl string `json:"avatarUrl,omitempty"`
	Mobile    string `json:"mobile,omitempty"`
	UnionId   string `json:"unionid,omitempty"`
	OpenId    string `json:"openid,omitempty"`
	Email     string `json:"email,omitempty"`
	StateCode string `json:"stateCode,omitempty"`
}

type UnionIdGetUserIdReq struct {
	UnionId string `json:"unionid" validate:"required"`
}

type UnionIdGetUserIdResponse struct {
	Response

	Result struct {
		UserId string `json:"userid"`

		//联系类型：
		//
		//0：企业内部员工
		//
		//1：企业外部联系人
		ContactType int `json:"contact_type"`
	} `json:"result"`
}

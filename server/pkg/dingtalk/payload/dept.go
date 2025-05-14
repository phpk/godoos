package payload

type DeptListReq struct {
	//父部门ID
	//如果不传，默认部门为根部门，根部门ID为1。只支持查询下一级子部门，不支持查询多级子部门。
	DeptId int `json:"dept_id,omitempty"`

	// 通讯录语言，默认zh_CN。如果是英文，请传入en_US。
	Language string `json:"language,omitempty" validate:"omitempty,oneof=zh_CN en_US"`
}

func NewDeptListReq(deptId int, language string) *DeptListReq {
	if deptId == 0 {
		deptId = 1
	}
	if language == "" {
		language = "zh_CN"
	}
	return &DeptListReq{DeptId: deptId, Language: language}
}

type DeptListResp struct {
	Response

	List []DeptBaseResponse `json:"result"`
}

type DeptBaseResponse struct {
	Id               int    `json:"dept_id"`
	Name             string `json:"name"`
	ParentId         int    `json:"parent_id"`
	SourceIdentifier string `json:"source_identifier"` // 部门标识字段
	CreateDeptGroup  bool   `json:"create_dept_group"` //是否同步创建一个关联此部门的企业群
	AutoAddUser      bool   `json:"auto_add_user"`     // 当部门群已经创建后，是否有新人加入部门会自动加入该群
}

// DeptOrderResp 员工在对应的部门中的排序
type DeptOrderResp struct {
	DeptId int `json:"dept_id"`

	Order int `json:"order"`
}

// DeptSimpleUserInfo 获取部门用户基础信息
type DeptSimpleUserInfoReq struct {
	// 部门ID，根部门ID为1
	DeptId int `json:"dept_id" validate:"required,min=1"`

	// 分页查询的游标，最开始传0，后续传返回参数中的next_cursor值
	Cursor int `json:"cursor" validate:"omitempty,min=0"`

	// 分页长度，最大值100。
	Size int `json:"size" validate:"required,max=100"`

	//部门成员的排序规则：
	//
	//entry_asc：代表按照进入部门的时间升序。
	//
	//entry_desc：代表按照进入部门的时间降序。
	//
	//modify_asc：代表按照部门信息修改时间升序。
	//
	//modify_desc：代表按照部门信息修改时间降序。
	//
	//custom：代表用户定义(未定义时按照拼音)排序。
	//
	//默认值：custom。
	OrderField string `json:"order_field,omitempty" validate:"omitempty,oneof=entry_asc entry_desc modify_asc modify_desc"`

	// 是否返回访问受限的员工
	ContainAccessLimit bool `json:"contain_access_limit,omitempty"`

	// 通讯录语言，默认zh_CN。如果是英文，请传入en_US。
	Language string `json:"language,omitempty" validate:"omitempty,oneof=zh_CN en_US"`
}

type DeptSimpleUserInfoResp struct {
	Response
	Result struct {
		// 是否还有更多的数据
		HasMore bool `json:"has_more"`

		// 下一次分页的游标，如果has_more为false，表示没有更多的分页数据。
		NextCursor int `json:"next_cursor"`

		DeptUsers []struct {
			UserId string `json:"userid"`

			Name string `json:"name"`
		} `json:"list"`
	} `json:"result"`
}

type SubDeptListReq struct {
	DeptId int `json:"dept_id" validate:"required"`
}

func NewSubDeptReq(id int) *SubDeptListReq {
	return &SubDeptListReq{id}
}

type SubDeptListResp struct {
	Response
	Result struct {
		Ids []int `json:"dept_id_list"`
	} `json:"result"`
}

type DeptDetailReq struct {
	DeptId int `json:"dept_id" validate:"required"`

	// 通讯录语言，默认zh_CN。如果是英文，请传入en_US。
	Language string `json:"language,omitempty" validate:"omitempty,oneof=zh_CN en_US"`
}

type DeptDetailResp struct {
	Response

	Detail struct {
		// 部门id
		Id int `json:"dept_id"`

		Name string `json:"name"`

		ParentId int `json:"parent_id"`

		// 备注
		Brief string `json:"brief"`

		// 部门标识字段
		SourceIdentifier string `json:"source_identifier"`

		//是否同步创建一个关联此部门的企业群：
		//
		//true：创建
		//
		//false：不创建
		CreateDeptGroup bool `json:"create_dept_group"`

		//当部门群已经创建后，是否有新人加入部门会自动加入该群：
		//
		//true：自动加入群
		//
		//false：不会自动加入群
		AutoAddUser bool `json:"auto_add_user"`

		//是否默认同意加入该部门的申请：
		//
		//true：表示加入该部门的申请将默认同意
		//
		//false：表示加入该部门的申请需要有权限的管理员同意
		AutoApproveApply bool `json:"auto_approve_apply"`

		//部门是否来自关联组织：
		//
		//true：是
		//
		//false：不是
		FromUnionOrg bool `json:"from_union_org"`

		//教育部门标签：
		//
		//campus：校区
		//
		//period：学段
		//
		//grade：年级
		//
		//class：班级
		Tags string `json:"tags"`

		// 在父部门中的排序值，order值小的排序靠前
		Order int `json:"order"`

		// 部门群ID
		DeptGroupChatId string `json:"dept_group_chat_id"`

		//
		//true：包含
		//
		//false：不包含
		//
		//不传值，则保持不变
		GroupContainSubDept bool `json:"group_contain_sub_dept"`

		// 企业群群主的userid
		OrgDeptOwner string `json:"org_dept_owner"`

		// 部门的主管userid列表，多个userid之间使用英文逗号分隔
		DeptManagerUseridList []string `json:"dept_manager_userid_list"`

		//是否限制本部门成员查看通讯录：
		//
		//true：开启限制。开启后本部门成员只能看到限定范围内的通讯录
		//
		//false（默认值）：不限制
		OuterDept bool `json:"outer_dept"`

		//指定本部门成员可查看的通讯录部门ID列表，总数不能超过200
		//
		//当outer_dept为true时，此参数生效
		UserPermitsDeptIds []int `json:"outer_permit_depts"`

		//指定本部门成员可查看的通讯录用户userid列表，总数不能超过200
		//
		//当outer_dept为true时，此参数生效。
		UserPermitsUsers []string `json:"outer_permit_users"`

		//是否隐藏本部门：
		//
		//true：隐藏部门，隐藏后本部门将不会显示在公司通讯录中
		//
		//false（默认值）：显示部门
		HideDept bool `json:"hide_dept"`

		//指定可以查看本部门的人员userid列表，总数不能超过200
		//
		//当hide_dept为true时，则此值生效
		UserPermits []string `json:"user_permits"`

		//指定可以查看本部门的其他部门列表，总数不能超过200
		//
		//当hide_dept为true时，则此值生效
		DeptPermits []int `json:"dept_permits"`
	} `json:"result"`
}

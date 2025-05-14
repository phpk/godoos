package payload

type RoleListReq struct {

	// 支持分页查询，与size参数同时设置时才生效，此参数代表偏移量，偏移量从0开始
	Offset int `json:"offset"`

	// 支持分页查询，与offset参数同时设置时才生效，此参数代表分页大小，默认值20，最大值200
	Size int `json:"size" validate:"max=200"`
}

func NewRoleListReq(offset, size int) *RoleListReq {
	return &RoleListReq{offset, size}
}

type RoleList struct {
	Response

	Result struct {
		// 是否还有更多数据。
		HasMore bool `json:"hasMore"`

		RoleGroups []struct {
			// 角色组Id
			GroupId int `json:"groupId"`

			// 角色组名称
			Name string `json:"name"`

			Roles []struct {
				// 角色Id
				Id int `json:"id"`

				// 角色名称
				Name string `json:"name"`
			} `json:"roles"`
		} `json:"list"`
	} `json:"result"`
}

type RoleUserReq struct {
	// 角色Id
	RoleId int `json:"role_id"`

	// 支持分页查询，与size参数同时设置时才生效，此参数代表偏移量，偏移量从0开始
	Offset int `json:"offset"`

	// 支持分页查询，与offset参数同时设置时才生效，此参数代表分页大小，默认值20，最大值200
	Size int `json:"size" validate:"max=100"`
}

func NewRoleUserReq(roleId, offset, size int) *RoleUserReq {
	return &RoleUserReq{roleId, offset, size}
}

type RoleUser struct {
	Response

	Result struct {
		HasMore    bool `json:"hasMore"`
		NextCursor int  `json:"nextCursor"`

		Users []struct {
			UserId string `json:"userid"`

			Name string `json:"name"`

			ManageScopes []struct {
				// 部门Id
				DeptId int `json:"dept_id"`
				// 部门名称
				DeptName string `json:"name"`
			} `json:"manageScopes"`
		} `json:"list"`
	} `json:"result"`
}

type RoleDetailReq struct {
	Id int `json:"roleId"`
}

func NewRoleDetailReq(id int) *RoleDetailReq {
	return &RoleDetailReq{id}
}

type RoleDetail struct {
	Response

	Role struct {
		// 角色名称
		RoleName string `json:"name"`

		// 所属的角色组Id
		GroupId int `json:"groupId"`
	} `json:"role"`
}
type GroupRoleReq struct {
	// 员工在企业中的userid
	GroupId int `json:"group_id"`
}

func NewGroupRoleReq(id int) *GroupRoleReq {
	return &GroupRoleReq{id}
}

type GroupRole struct {
	Response

	Group struct {
		// 角色组名
		GroupName string `json:"group_name"`

		Roles []struct {
			// 角色Id
			Id int `json:"role_id"`

			// 角色名
			Name string `json:"role_name"`
		} `json:"roles"`
	} `json:"role_group"`
}

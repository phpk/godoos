package dingtalk

import (
	"godocms/pkg/dingtalk/payload"
	"net/http"
)

// GetRoleList 获取角色列表
func (ding *DingTalk) GetRoleList(offset, size int) (resp payload.RoleList, err error) {
	return resp, ding.Request(http.MethodPost, GetRoleListKey, nil,
		payload.NewRoleListReq(offset, size), &resp)
}

// GetRoleUserList 获取指定角色的员工列表
func (ding *DingTalk) GetRoleUserList(roleId, offset, size int) (resp payload.RoleUser, err error) {
	return resp, ding.Request(http.MethodPost, GetRoleUserListKey, nil,
		payload.NewRoleUserReq(roleId, offset, size), &resp)
}

// GetRoleDetail 获取角色详情
func (ding *DingTalk) GetRoleDetail(roleId int) (resp payload.RoleDetail, err error) {
	return resp, ding.Request(http.MethodPost, GetRoleDetailKey, nil,
		payload.NewRoleDetailReq(roleId), &resp)
}

// GetGroupRoles 获取角色组列表
func (ding *DingTalk) GetGroupRoles(groupId int) (resp payload.GroupRole, err error) {
	return resp, ding.Request(http.MethodPost, GetRoleGroupKey, nil,
		payload.NewGroupRoleReq(groupId), &resp)
}

// CheckInRoleGroup 检查角色组是否在角色组内
func (ding *DingTalk) CheckInRoleGroup(roleID int) (bool, error) {
	var offset, size int
	offset, size = 0, 200
	for {
		resp, err := ding.GetRoleList(offset, size)
		if err != nil {
			return false, err
		}

		for _, group := range resp.Result.RoleGroups {
			if group.GroupId == roleID {
				return true, nil
			}
		}

		if resp.Result.HasMore {
			offset += size
		} else {
			break
		}
	}

	return false, nil
}

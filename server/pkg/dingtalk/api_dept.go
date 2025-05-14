package dingtalk

import (
	"godocms/pkg/dingtalk/payload"
	"net/http"
)

// GetDeptList 获取部门列表, 调用本接口，获取下一级部门基础信息
func (ding *DingTalk) GetDeptList(req *payload.DeptListReq) (resp payload.DeptListResp, err error) {
	if req == nil {
		req = payload.NewDeptListReq(1, "zh_CN")
	}
	return resp, ding.Request(http.MethodPost, GetDeptListKey, nil, req, &resp)
}

// GetDeptSimpleUserInfo 获取部门用户基础信息 调用本接口获取指定部门的用户基础信息
func (ding *DingTalk) GetDeptSimpleUserInfo(req *payload.DeptSimpleUserInfoReq) (resp payload.DeptSimpleUserInfoResp, err error) {
	return resp, ding.Request(http.MethodPost, GetDeptSimpleUserKey, nil, req, &resp)
}

// GetSubDeptList 获取子部门列表id
func (ding *DingTalk) GetSubDeptList(deptId int) (rsp payload.SubDeptListResp, err error) {
	return rsp, ding.Request(http.MethodPost, GetSubDeptListKey, nil, payload.NewSubDeptReq(deptId), &rsp)
}

// GetDeptDetail 获取部门详情
func (ding *DingTalk) GetDeptDetail(req *payload.DeptDetailReq) (resp payload.DeptDetailResp, err error) {
	return resp, ding.Request(http.MethodPost, GetDeptDetailKey, nil, req, &resp)
}

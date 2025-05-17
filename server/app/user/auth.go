package user

import (
	"errors"
	"fmt"
	"godocms/app/user/login"
	"godocms/common"
	"godocms/libs"
	"godocms/middleware"
	"godocms/model"
	"godocms/service"
	"godocms/utils"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	middleware.RegisterRouter("POST", "user/login", loginHandle, 0, "用户登录")
	middleware.RegisterRouter("GET", "user/logout", logoutHandler, 1, "退出登录")
	middleware.RegisterRouter("POST", "user/register", registerHandle, 0, "用户注册")
	middleware.RegisterRouter("GET", "user/islogin", isLogin, 1, "是否登录")
}

// handleLogin 处理登录请求
func loginHandle(c *gin.Context) {
	var req login.LoginRequest
	// 解析请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		libs.Error(c, "请求参数错误，登录失败")
		return
	}

	//fmt.Printf("request param: %+v\n", req)

	factory := &login.LoginHandlerFactory{}
	req.Action = "login"
	handler, err := factory.GetHandler(req)
	if err != nil {
		libs.Error(c, "适配失败："+err.Error())
		return
	}
	user, err := handler.Login()
	if err != nil || user == nil {
		libs.Error(c, "登录失败:"+err.Error())
		return
	}

	// 验证用户
	userDept, userRole, err := validateUser(c, user, req.ClientId)
	if err != nil {
		libs.Error(c, "验证用户失败:"+err.Error())
		return
	}

	// 生成token
	token, err := utils.GenerateToken(&utils.UserClaims{ID: int64(user.ID)})
	if err != nil {
		libs.Error(c, "获取toekn失败:"+err.Error())
		return
	}

	// err = files.InitOsSystem(user.ID)
	// if err != nil {
	// 	libs.Error(c, "初始化系统失败:"+err.Error())
	// 	return
	// }

	// 获取用户auths和分享列表
	userAuths, userShare, err := service.GetUserAuthAndShare(user)
	if err != nil {
		libs.Error(c, "获取用户权限失败:"+err.Error())
		return
	}

	resp := login.LoginResponse{
		User: login.UserResponse{
			ID:        user.ID,
			Name:      user.Username,
			Nickname:  user.Nickname,
			Sex:       user.Sex,
			Avatar:    user.Avatar,
			Desc:      user.Desc,
			Email:     user.Email,
			Phone:     user.Phone,
			JobNumber: user.JobNumber,
			UseSpace:  userRole.Space,
			HasSpace:  user.UseSpace,
			WorkPlace: user.WorkPlace,
			HiredDate: user.HiredDate,
		},
		Role: login.RoleReponse{
			ID:   userRole.ID,
			Name: userRole.Name,
		},
		Dept: login.DeptResponse{
			ID:   userDept.ID,
			Name: userDept.Name,
		},
		UserAuths:  userAuths,
		UserShares: userShare,
		Token:      token,
		ClientID:   req.ClientId,
	}

	// if user.Phone == "" {
	// 	libs.ErrorMsg(c, login.UserNoMobileErrCode, "请先绑定手机号", resp)
	// 	return
	// }

	libs.Success(c, "success", resp)
}

func validateUser(c *gin.Context, user *model.User, clientID string) (*model.UserDept, *model.UserRole, error) {
	if user.ID == 0 {
		return nil, nil, errors.New("用户不存在")
	}

	if user.Status == 1 {
		// 审核？
		return nil, nil, errors.New("用户审核中")
	}
	userRole, err := service.GetUserRole(user.RoleId)

	if err != nil {
		slog.Error("获取用户角色失败", "err", err.Error())
		return nil, nil, err
	}

	// userSum, err := files.GetUserSum(user.ID)
	// if err != nil {
	// 	slog.Error("查询用户空间失败", "err", err.Error())
	// 	userSum = 0
	// 	// return nil, nil, err
	// }

	// if userSum > userRole.Space {
	// 	return nil, nil, errors.New("空间已满")
	// }

	// user.UseSpace = userSum
	if err := service.UpdateLoginUser(c, user); err != nil {
		slog.Error("更新用户信息失败", "err", err.Error())
		return nil, nil, err
	}
	userDept, err := service.GetUserDept(user.DeptId)
	if err != nil {
		return nil, nil, err
	}
	cacheData := login.LoginCache{
		UserID:    user.ID,
		UserRoles: userRole.Rules,
	}
	common.SetCache("userData:"+clientID, cacheData, 60*24*time.Minute)
	return &userDept, &userRole, nil
}

func logoutHandler(c *gin.Context) {
	var req struct {
		ClientID string `json:"client_id" binding:"required"`
	}
	// 绑定请求体中的 JSON 数据
	if err := c.ShouldBindJSON(&req); err != nil {
		libs.Error(c, "参数错误")
		return
	}

	if value, ok := common.GetCache("userData:" + req.ClientID); ok != nil {
		common.DelCache("userData:" + req.ClientID)
		libs.Success(c, "退出成功", value)
	} else {
		libs.Error(c, "用户不存在")
	}

	var origin map[int64]int64
	var ok bool
	data, _ := common.GetCache("user_online")
	if origin, ok = data.(map[int64]int64); ok {
		uid := c.GetInt64("userId")
		delete(origin, uid)
	}
	common.SetCache("user_online", origin, 60*24)

	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(200, gin.H{"code": 0, "msg": "退出成功"})
}

// handleRegister 处理注册请求
func registerHandle(c *gin.Context) {
	var req login.LoginRequest
	// 解析请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		libs.Error(c, "参数错误，注册失败:"+err.Error())
		return
	}

	fmt.Printf("request param: %+v\n", req)

	factory := &login.LoginHandlerFactory{}
	req.Action = "register"
	handler, err := factory.GetHandler(req)

	if err != nil {
		libs.Error(c, "适配失败:"+err.Error())
		return
	}
	_, err = handler.Register()
	if err != nil {
		libs.Error(c, "注册失败:"+err.Error())
		return
	}

	libs.Success(c, "注册成功", nil)
}
func isLogin(c *gin.Context) {
	libs.Success(c, "已登录", nil)
}

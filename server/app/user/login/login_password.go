package login

import (
	"encoding/json"
	"errors"
	"godocms/common"
	"godocms/model"
	"godocms/service"
	"godocms/utils"
	"time"
)

type PasswordLoginHandler struct {
	loginParam *PasswordLoginParam
}

func (h *PasswordLoginHandler) Init(req LoginRequest) error {
	//此处可以根据action处理不同的数据类型的解析
	h.loginParam = new(PasswordLoginParam)
	return json.Unmarshal(req.Param, h.loginParam)
}
func (h *PasswordLoginHandler) Login() (user *model.User, err error) {
	if !common.LoginConf.UserNameLogin.Enable {
		return nil, errors.New("登录已关闭")
	}
	//log.Printf("用户登录: %+v", h.loginParam.Password)
	user, err = service.GetUserByUsername(h.loginParam.Username)
	if err != nil {
		return nil, err
	}
	if utils.HashPassword(h.loginParam.Password, user.Salt) != user.Password {
		return nil, errors.New("密码错误")
	}
	// 实现短信验证码登录逻辑
	return user, nil
}

func (h *PasswordLoginHandler) Register() (user *model.User, err error) {
	//log.Printf("用户登录: %+v", h.loginParam)
	// 检查用户名和邮箱是否已存在
	//log.Printf("=====用户登录: %+v", h.loginParam)
	existingUser, err := service.GetUserByUsername(h.loginParam.Username)
	// log.Printf("=====用户登录: %+v", existingUser)
	// log.Printf("=====用户登录: %+v", err)
	if err == nil || existingUser != nil {
		return nil, errors.New("用户名或邮箱/手机号已存在")
	}
	salt, err := utils.GenerateSalt(16)
	if err != nil {
		return nil, err
	}
	hashedPassword := utils.HashPassword(h.loginParam.Password, salt)
	user = &model.User{
		Username: h.loginParam.Username,
		Password: hashedPassword,
		Salt:     salt,
		Email:    "",
		Nickname: h.loginParam.Username,
		Phone:    "",
		RoleId:   1,
		Status:   0,
		DeptId:   1,
		AddTime:  int32(time.Now().Unix()),
		UpTime:   0,
	}
	thirdInfo := &model.UserThird{
		Platform:    LoginPlatformPassword,
		ThirdUserID: h.loginParam.Username,
		UnionID:     h.loginParam.Username,
	}
	user, err = service.AddUser(user, thirdInfo)
	if err != nil {
		return user, err
	}
	return user, nil
}

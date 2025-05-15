package login

import (
	"encoding/json"
	"errors"
	"godocms/common"
	"godocms/model"
	"godocms/pkg/db"
	"godocms/utils"
	"time"

	"gorm.io/gorm"
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
	user = new(model.User)
	if err := db.DB.Where("username = ?", h.loginParam.Username).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
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
	var existingUser model.User
	result := db.DB.Where("username = ?", h.loginParam.Username).First(&existingUser)
	if result.RowsAffected > 0 {
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

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		thirdInfo := &model.UserThird{
			UserID:      user.ID,
			Platform:    LoginPlatformPassword,
			ThirdUserID: h.loginParam.Username,
			UnionID:     h.loginParam.Username,
		}
		if err := tx.Create(thirdInfo).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return user, err
	}
	return user, nil
}

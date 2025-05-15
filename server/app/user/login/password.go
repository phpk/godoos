package login

import (
	"godocms/model"
	"log"
)

type PasswordLoginHandler struct{}

func (h *PasswordLoginHandler) Login(param LoginParam) (*model.User, error) {
	p := param.(PasswordLoginParam)
	log.Printf("用户登录: %+v", p)
	// 实现短信验证码登录逻辑
	return nil, nil
}

func (h *PasswordLoginHandler) Register(param LoginParam) (*model.User, error) {
	p := param.(PasswordLoginParam)
	log.Printf("用户登录: %+v", p)

	return nil, nil
}

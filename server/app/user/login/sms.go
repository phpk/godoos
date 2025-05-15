package login

import (
	"godocms/model"
	"log"
)

type SmsCodeLoginHandler struct{}

func (h *SmsCodeLoginHandler) Login(param LoginParam) (*model.User, error) {
	p := param.(SmsCodeLoginParam)
	log.Printf("手机号登录: %+v", p)
	// 实现短信验证码登录逻辑
	return nil, nil
}

func (h *SmsCodeLoginHandler) Register(param LoginParam) (*model.User, error) {
	p := param.(SmsCodeLoginParam)
	log.Printf("手机号登录: %+v", p)
	// 实现短信验证码注册逻辑
	return nil, nil
}

package login

import (
	"encoding/json"
	"godocms/model"
	"log"
)

type SmsCodeLoginHandler struct {
	param *SmsCodeLoginParam
}

func (h *SmsCodeLoginHandler) Init(req LoginRequest) error {
	h.param = new(SmsCodeLoginParam)
	return json.Unmarshal(req.Param, h.param)
}
func (h *SmsCodeLoginHandler) Login() (*model.User, error) {
	log.Printf("用户登录: %+v", h.param)
	// 实现短信验证码登录逻辑
	return nil, nil
}

func (h *SmsCodeLoginHandler) Register() (*model.User, error) {
	log.Printf("用户登录: %+v", h.param)
	// 实现短信验证码注册逻辑
	return nil, nil
}

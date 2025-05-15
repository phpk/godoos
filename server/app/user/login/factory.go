package login

import (
	"errors"
	"fmt"
)

type LoginHandlerFactory struct{}

// factory.go
func (f *LoginHandlerFactory) GetHandler(req LoginRequest) (LoginHandler, error) {
	var handler LoginHandler

	switch req.LoginType {
	case LoginTypePassword:
		handler = &PasswordLoginHandler{}
	case LoginTypeSmsCode:
		handler = &SmsCodeLoginHandler{}
	default:
		return nil, errors.New("unsupported login type")
	}

	if err := handler.Init(req); err != nil {
		return nil, fmt.Errorf("参数初始化失败: %v", err)
	}

	return handler, nil
}

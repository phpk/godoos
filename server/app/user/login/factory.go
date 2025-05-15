package login

import "errors"

type LoginHandlerFactory struct{}

func (f *LoginHandlerFactory) GetHandler(loginType string) (LoginHandler, error) {
	switch loginType {
	case LoginTypePassword:
		return &PasswordLoginHandler{}, nil
	case LoginTypeSmsCode:
		return &SmsCodeLoginHandler{}, nil
	default:
		return nil, errors.New("unsupported login type")
	}
}

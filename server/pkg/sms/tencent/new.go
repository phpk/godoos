package tencent

import "errors"

type qcloudsms struct {
	appid                 int
	appkey                string
	SmsSingleSender       *smsSingleSender
	SmsMultiSender        *smsMultiSender
	SmsStatusPuller       *smsStatusPuller
	SmsMobileStatusPuller *smsMobileStatusPuller
}

// NewQcloudSms new一个qcloudsms实例
func NewQcloudSms(appid int, appkey string) (*qcloudsms, error) {
	if appkey == "" {
		return nil, errors.New("appkey is nil")
	}
	return &qcloudsms{
		appid:                 appid,
		appkey:                appkey,
		SmsSingleSender:       newSmsSingleSender(appid, appkey),
		SmsMultiSender:        newSmsMultiSender(appid, appkey),
		SmsStatusPuller:       newSmsStatusPuller(appid, appkey),
		SmsMobileStatusPuller: newSmsMobileStatusPuller(appid, appkey),
	}, nil
}

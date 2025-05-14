package alibaba

import (
	"encoding/json"
	"errors"
	"fmt"
	"godocms/common"
	"log/slog"
	r "math/rand"
	"time"
)

type SmsResponse struct {
	Code      string `json:"Code"`
	Message   string `json:"Message"`
	RequestId string `json:"RequestId"`
	BizId     string `json:"BizId"`
}

// 发送短信
func SendSms(phone string) (err error) {
	httpMethod := "POST"            // 请求方式，大部分RPC接口同时支持POST和GET，此处以POST为例
	canonicalUri := "/"             // RPC接口无资源路径，故使用正斜杠（/）作为CanonicalURI
	host := "dysmsapi.aliyuncs.com" // 云产品服务接入点
	xAcsAction := "SendSms"         // API名称
	xAcsVersion := "2017-05-25"     // API版本号
	req := NewRequest(httpMethod, canonicalUri, host, xAcsAction, xAcsVersion)

	req.queryParam["PhoneNumbers"] = phone
	req.queryParam["SignName"] = common.LoginConf.Phone.AliyunSms.SignName
	req.queryParam["TemplateCode"] = common.LoginConf.Phone.AliyunSms.TemplateCode
	// req.queryParam["TemplateParam"] = "{\"code\":\"1234\"}"
	if common.LoginConf.Phone.AliyunSms.TemplateParam == "" {
		common.LoginConf.Phone.AliyunSms.TemplateParam = "{\"code\":\"123456\"}"
	}
	smscode := GenerateRandomSixDigitNumber()
	templateParam, err := replaceTemplateOneParam(common.LoginConf.Phone.AliyunSms.TemplateParam, smscode)
	if err != nil {
		return err
	}
	req.queryParam["TemplateParam"] = templateParam

	getAuthorization(req)
	var resp SmsResponse
	body, err := callAPI(req)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return err
	}
	if resp.Code != "OK" {
		return errors.New("发送短信失败")
	}
	slog.Info("发送短信成功", "phone", phone, slog.Any("resp", resp), slog.Any("smscode", smscode))

	common.Cache.Set("sms:"+phone, smscode, 1*time.Minute)

	return nil
}
func GenerateRandomSixDigitNumber() string {
	return fmt.Sprintf("%06v", r.New(r.NewSource(time.Now().UnixNano())).Int31n(1000000))
}

// 替换模板参数
func replaceTemplateOneParam(templateParam string, value string) (string, error) {
	// 解析 JSON 模板参数
	var paramMap map[string]string
	err := json.Unmarshal([]byte(templateParam), &paramMap)
	if err != nil {
		return "", fmt.Errorf("解析模板参数失败: %v", err)
	}

	// 检查参数数量
	if len(paramMap) != 1 {
		return "", fmt.Errorf("模板参数必须只有一个键值对")
	}

	// 替换参数值
	for key := range paramMap {
		paramMap[key] = value
	}

	// 将 map 转换回 JSON 字符串
	updatedParam, err := json.Marshal(paramMap)
	if err != nil {
		return "", fmt.Errorf("生成 JSON 失败: %v", err)
	}

	return string(updatedParam), nil
}

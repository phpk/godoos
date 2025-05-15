package utils

import (
	"crypto/tls"
	"fmt"
	g "godocms/common"
	"strings"

	"gopkg.in/gomail.v2"
)

// Email 发送方法
// to: 以英文逗号 , 分隔的字符串 如 "a@qq.com,b@qq.com"
func Email(to, subject, body string) error {
	return send(strings.Split(to, ","), subject, body)
}

// 邮件发送
// to: 目标邮件
// subject: 邮件标题
// body: 邮件内容 (HTML)
func send(to []string, subject string, body string) error {
	from := g.LoginConf.Email.From
	nickname := g.LoginConf.Email.Username
	secret := g.LoginConf.Email.Password // 使用授权码
	host := g.LoginConf.Email.Host
	port := g.LoginConf.Email.Port
	isSsl := g.LoginConf.Email.IsSsl

	// 创建邮件消息
	m := gomail.NewMessage()
	fromstr := fmt.Sprintf("%s <%s>", nickname, from)
	m.SetHeader("From", fromstr)
	//m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	//m.SetBody("text/plain", body)
	// 创建 SMTP 客户端
	d := gomail.NewDialer(host, port, from, secret)

	// 设置 TLS 配置
	if isSsl {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("邮件发送失败:", err.Error())
		return err
	}

	return nil
}

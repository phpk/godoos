package model

import "gorm.io/gorm"

type LocalProxy struct {
	gorm.Model
	Port      uint   `json:"port"`
	ProxyType string `json:"proxy_type"`
	Domain    string `json:"domain"`
}

func (*LocalProxy) TableName() string {
	return "local_proxy"
}

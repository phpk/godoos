package model

import "gorm.io/gorm"

type ClientUser struct {
	gorm.Model
	ServerUrl string `json:"server_url"`
	DiskId    string `json:"disk_id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

func (*ClientUser) TableName() string {
	return "client_user"
}

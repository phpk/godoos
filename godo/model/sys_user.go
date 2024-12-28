package model

import "gorm.io/gorm"

type SysUser struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

func (*SysUser) TableName() string {
	return "sys_user"
}

package model

type SysUser struct {
	BaseModel
	Username string `json:"username"`
	Password string `json:"password"`
}

func (*SysUser) TableName() string {
	return "sys_user"
}

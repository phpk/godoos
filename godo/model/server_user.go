package model

type ServerUser struct {
	BaseModel
	DiskId   string `json:"disk_id"`
	AuthType string `json:"auth_type"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

func (*ServerUser) TableName() string {
	return "server_user"
}

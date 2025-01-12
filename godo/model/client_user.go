package model

type ClientUser struct {
	BaseModel
	ServerUrl string `json:"server_url"`
	DiskId    string `json:"disk_id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

func (*ClientUser) TableName() string {
	return "client_user"
}

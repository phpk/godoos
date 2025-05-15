package model

import (
	"time"

	"gorm.io/gorm"
)

// UserThird 第三方平台登录所记录的额外信息
type UserThird struct {
	ID          int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	GroupID     int32     `gorm:"column:group_id;default:1;comment:租户ID" json:"group_id"` // 租户ID
	UserID      int32     `gorm:"column:user_id;index:idx_user_id" json:"user_id"`
	Platform    string    `gorm:"column:platform" json:"platform"`                   // 平台: 钉钉(dingtalk)、github、gitee、企业微信(wxqy)
	ThirdUserID string    `gorm:"column:third_user_id" json:"third_user_id"`         // 平台下应用对应的用户的id, 像钉钉的userid或微信的openid
	UnionID     string    `gorm:"column:union_id;index:idx_unionid" json:"union_id"` // 平台对应的用户唯一id (如钉钉和企业微信的unionid)
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (*UserThird) TableName() string {
	return "user_third"
}
func (u *UserThird) BeforeSave(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return
}

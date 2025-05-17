package model

import (
	"fmt"
	"godocms/pkg/db"
	"log"
	"time"
)

const TableNameUserRole = "user_role"

// UserRole 用户角色表
type UserRole struct {
	ID         int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	GroupID    int32  `gorm:"column:group_id;default:1;comment:租户ID" json:"group_id"` // 租户ID
	Name       string `gorm:"column:name" json:"name"`
	Rules      string `gorm:"column:rules" json:"rules"`
	Space      int32  `gorm:"column:space;default:1024;comment:存储大小" json:"space"`    // 存储大小
	Status     int32  `gorm:"column:status;comment:是否可用0可用1不可用" json:"status"`        // 是否可用0可用1不可用
	Remark     string `gorm:"column:remark;comment:描述" json:"remark"`                 // 描述
	MenuIDS    int32  `gorm:"column:menu_ids;index;not null" json:"menu_ids"`         // 菜单ID
	DingRoleID string `gorm:"column:ding_role_id" json:"ding_role_id"`                // 钉钉角色ID
	AddTime    int32  `gorm:"column:add_time;default:0;comment:添加时间" json:"add_time"` // 添加时间
	UpTime     int32  `gorm:"column:up_time;default:0;comment:更新时间" json:"up_time"`   // 更新时间
}

// TableName UserRole's table name
func (*UserRole) TableName() string {
	return TableNameUserRole
}
func EnsureDefaultUserRoleExists() error {
	var count int64
	err := db.DB.Model(&UserRole{}).Count(&count).Error
	if err != nil {
		log.Fatalf("Count failed on %s table: %v", TableNameUserRole, err)
		return err
	}
	//log.Printf("Total roles in %s table: %d", TableNameUserRole, count)

	if count == 0 {
		defaultRole := UserRole{
			ID:         1,
			GroupID:    1,
			Name:       "默认角色",
			Rules:      "-1",
			Space:      1024,
			Status:     0,
			Remark:     "-1",
			MenuIDS:    123456, // 修正超出 int32 范围的问题
			DingRoleID: "",
			AddTime:    int32(time.Now().Unix()),
			UpTime:     0,
		}

		if err := db.DB.Create(&defaultRole).Error; err != nil {
			return fmt.Errorf("failed to create default user role: %w", err)
		}
	}

	return nil
}

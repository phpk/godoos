package models

import (
	"time"
)

// 统一用户模型
type User struct {
	ID        interface{} `json:"id" gorm:"primaryKey" bson:"_id,omitempty"`
	Name      string      `json:"name" gorm:"size:100" bson:"name"`
	Email     string      `json:"email" gorm:"size:100;uniqueIndex" bson:"email"`
	CreatedAt time.Time   `json:"created_at" gorm:"autoCreateTime" bson:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" gorm:"autoUpdateTime" bson:"updated_at"`
}

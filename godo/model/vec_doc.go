package model

import "gorm.io/gorm"

type VecDoc struct {
	gorm.Model
	Content  string `json:"content"`
	FilePath string `json:"file_path" gorm:"not null"`
	ListID   int    `json:"list_id"`
}

func (VecDoc) TableName() string {
	return "vec_doc"
}

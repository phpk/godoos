package model

import (
	"fmt"

	"gorm.io/gorm"
)

type VecList struct {
	gorm.Model
	FilePath       string `json:"file_path" gorm:"not null"`
	Engine         string `json:"engine" gorm:"not null"`
	EmbedSize      int    `json:"embed_size"`
	EmbeddingModel string `json:"model" gorm:"not null"`
}

func (*VecList) TableName() string {
	return "vec_list"
}

// BeforeCreate 在插入数据之前检查是否存在相同路径的数据
func (v *VecList) BeforeCreate(tx *gorm.DB) error {
	var count int64
	if err := tx.Model(&VecList{}).Where("file_path = ?", v.FilePath).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("file path already exists: %s", v.FilePath)
	}
	return nil
}

// AfterCreate 在插入数据之后创建虚拟表
func (v *VecList) AfterCreate(tx *gorm.DB) error {
	return CreateVirtualTable(tx, v.ID, v.EmbedSize)
}

// AfterDelete 在删除数据之后删除虚拟表
func (v *VecList) AfterDelete(tx *gorm.DB) error {
	// 删除 VecDoc 表中 ListID 对应的所有数据
	if err := tx.Where("list_id = ?", v.ID).Delete(&VecDoc{}).Error; err != nil {
		return err
	}
	return DropVirtualTable(tx, v.ID)
}

// CreateVirtualTable 创建虚拟表
func CreateVirtualTable(db *gorm.DB, vectorID uint, embeddingSize int) error {
	sql := fmt.Sprintf(`
		CREATE VIRTUAL TABLE IF NOT EXISTS [%d_vec] USING
		vec0(
			document_id TEXT PRIMARY KEY,
			embedding float[%d] distance_metric=cosine
		)
	`, vectorID, embeddingSize)
	//log.Printf("sql: %s", sql)
	return db.Exec(sql).Error
}

// DropVirtualTable 删除虚拟表
func DropVirtualTable(db *gorm.DB, vectorID uint) error {
	sql := fmt.Sprintf(`DROP TABLE IF EXISTS [%d_vec]`, vectorID)
	return db.Exec(sql).Error
}

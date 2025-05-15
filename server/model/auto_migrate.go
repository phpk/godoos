package model

import (
	"godocms/pkg/db"
	"log"

	"gorm.io/gorm"
)

func AutoMigrate() error {
	db.DB.AutoMigrate(&User{})
	// 使用 GORM 的 CreateIndex 方法
	// createIndexIfNotExists(db.DB, &User{}, "username", "idx_user_username")
	// createIndexIfNotExists(db.DB, &User{}, "email", "idx_user_email")
	// createIndexIfNotExists(db.DB, &User{}, "phone", "idx_user_phone")
	db.DB.AutoMigrate(&UserRole{})
	db.DB.AutoMigrate(&UserDept{})
	db.DB.AutoMigrate(&UserThird{})
	return nil
}
func CreateIndexIfNotExists(Db *gorm.DB, model interface{}, field, indexName string) {
	if db.AppConfig.DBType != "mongodb" {
		return
	}
	err := Db.Migrator().CreateIndex(model, indexName)
	if err != nil {
		log.Printf("Failed to create index %s for field %s: %v", indexName, field, err)
	}
}

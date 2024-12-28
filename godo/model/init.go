package model

import (
	"godo/libs"

	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/ncruces/go-sqlite3/gormlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDB() {
	dbPath := libs.GetSystemDb()
	db, err := gorm.Open(gormlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return
	}
	Db = db
	// 自动迁移模式
	db.AutoMigrate(&SysDisk{})
	// 初始化 SysDisk 记录
	initSysDisk(db)
	db.AutoMigrate(&SysUser{})
	db.AutoMigrate(&ClientUser{})
	db.AutoMigrate(&ServerUser{})
}

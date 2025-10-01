package model

import (
	"godo/libs"
	"time"

	_ "github.com/asg017/sqlite-vec-go-bindings/ncruces"
	//_ "github.com/ncruces/go-sqlite3/embed"
	//"github.com/ncruces/go-sqlite3/gormlite"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDB() {
	dbPath := libs.GetSystemDb()
	//fmt.Printf("dbPath: %s\n", dbPath)
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
		return
	}
	// Enable PRAGMAs
	// - busy_timeout (ms) to prevent db lockups as we're accessing the DB from multiple separate processes in otto8
	tx := db.Exec(`
PRAGMA busy_timeout = 10000;
`)
	if tx.Error != nil {
		return
	}
	Db = db
	// 自动迁移模式
	db.AutoMigrate(&SysDisk{})
	// 初始化 SysDisk 记录
	initSysDisk(db)
	db.AutoMigrate(&LocalProxy{})
	db.AutoMigrate(&SysUser{})
	db.AutoMigrate(&ClientUser{})
	db.AutoMigrate(&ServerUser{})
	db.AutoMigrate(&VecList{})
	db.AutoMigrate(&VecDoc{})
	db.AutoMigrate(&FrpcProxy{})
}

type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

package model

import (
	"fmt"
	"godo/libs"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

type SysDisk struct {
	BaseModel
	Name   string `json:"name"`
	Disk   string `json:"disk" gorm:"unique"`
	Size   int64  `json:"size"`
	Type   uint   `json:"type"` //0C-E本地 1nasserver 2nasclient 3webdavserver 4webdavclient 5F分享 6B回收站
	Path   string `json:"path"`
	Status uint   `json:"status"`
}

func (*SysDisk) TableName() string {
	return "sys_disk"
}
func initSysDisk(db *gorm.DB) {
	var count int64
	db.Model(&SysDisk{}).Count(&count)
	basePath, err := libs.GetOsDir()
	if err != nil {
		basePath, _ = os.Getwd()
	}
	if count == 0 {
		disks := []SysDisk{
			{Disk: "B", Name: "回收站", Size: 0, Path: filepath.Join(basePath, "B"), Type: 6, Status: 1},
			{Disk: "C", Name: "系统", Size: 0, Path: filepath.Join(basePath, "C"), Type: 0, Status: 1},
			{Disk: "D", Name: "文档", Size: 0, Path: filepath.Join(basePath, "D"), Type: 0, Status: 1},
			{Disk: "E", Name: "办公", Size: 0, Path: filepath.Join(basePath, "E"), Type: 0, Status: 1},
		}
		db.Create(&disks)
		fmt.Println("Initialized A-Z disks")
	}
}

// BeforeDelete 钩子
func (sd *SysDisk) BeforeDelete(tx *gorm.DB) (err error) {
	// 不允许删除的磁盘列表
	nonDeletableDisks := map[string]struct{}{
		"B": {},
		"C": {},
		"D": {},
		"E": {},
		"F": {},
	}

	if _, exists := nonDeletableDisks[sd.Disk]; exists {
		return fmt.Errorf("disk %s cannot be deleted", sd.Disk)
	}
	return nil
}

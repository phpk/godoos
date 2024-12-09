package vector

import (
	"fmt"
	"godo/ai/vector/model"
	"godo/libs"

	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/ncruces/go-sqlite3/gormlite"
	"gorm.io/gorm"
)

var vectorListDb *gorm.DB

func init() {
	dbPath := libs.GetVectorDb()
	db, err := gorm.Open(gormlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return
	}
	vectorListDb = db
}

// CreateVector 创建一个新的 VectorList 记录
func CreateVector(name string, filePath string) (*gorm.DB, error) {
	var list model.VectorList
	if !libs.PathExists(filePath) {
		return nil, fmt.Errorf("file path not exists")
	}
	// 检查是否已经存在同名的 VectorList
	result := vectorListDb.Where("name = ?", name).First(&list)
	if result.Error == nil {
		return nil, fmt.Errorf("vector list with the same name already exists")
	}

	// 创建新的 VectorList 记录
	newList := model.VectorList{
		Name:     name,
		FilePath: filePath, // 根据实际情况设置文件路径
	}

	result = vectorListDb.Create(&newList)
	if result.Error != nil {
		return nil, result.Error
	}

	return vectorListDb, nil
}

// DeleteVector 删除指定名称的 VectorList 记录
func DeleteVector(name string) error {
	result := vectorListDb.Where("name = ?", name).Delete(&model.VectorList{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("vector list not found")
	}

	return nil
}

// RenameVectorDb 更改指定名称的 VectorList 的数据库名称
func RenameVectorDb(oldName string, newName string) error {
	// 2. 获取旧的 VectorList 记录
	var oldList model.VectorList
	result := vectorListDb.Where("name = ?", oldName).First(&oldList)
	if result.Error != nil {
		return fmt.Errorf("failed to find old vector list: %w", result.Error)
	}

	// 5. 更新 VectorList 记录中的 DbPath 和 Name
	oldList.Name = newName
	result = vectorListDb.Save(&oldList)
	if result.Error != nil {
		return fmt.Errorf("failed to update vector list: %w", result.Error)
	}

	return nil
}

func GetVectorList() []model.VectorList {
	var vectorList []model.VectorList
	vectorListDb.Find(&vectorList)
	return vectorList
}

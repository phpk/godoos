package vector

import (
	"database/sql"
	"fmt"
	"godo/ai/vector/model"
	"godo/libs"
	"os"

	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/ncruces/go-sqlite3/gormlite"
	"gorm.io/gorm"
)

var vectorListDb *gorm.DB
var dbPathToSqlDB = make(map[string]*sql.DB)

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
	dbPath := libs.GetVectorPath(name)
	var list model.VectorList

	// 检查是否已经存在同名的 VectorList
	result := vectorListDb.Where("name = ?", name).First(&list)
	if result.Error == nil {
		return nil, fmt.Errorf("vector list with the same name already exists")
	}

	// 创建新的 VectorList 记录
	newList := model.VectorList{
		Name:     name,
		FilePath: filePath, // 根据实际情况设置文件路径
		DbPath:   dbPath,
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

	// 关闭数据库连接
	err := CloseVectorDb(name)
	if err != nil {
		return fmt.Errorf("failed to close vector database: %w", err)
	}
	// 删除数据库文件
	dbPath := libs.GetVectorPath(name)
	if libs.PathExists(dbPath) {
		err := os.Remove(dbPath)
		if err != nil {
			return fmt.Errorf("failed to delete database file: %w", err)
		}
	}

	return nil
}

// CloseVectorDb 关闭指定名称的 Vector 数据库连接
func CloseVectorDb(name string) error {
	dbPath := libs.GetVectorPath(name)
	if !libs.PathExists(dbPath) {
		return nil
	}
	sqlDB, exists := dbPathToSqlDB[dbPath]
	if !exists {
		return fmt.Errorf("no database connection found for path: %s", dbPath)
	}
	err := sqlDB.Close()
	if err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	delete(dbPathToSqlDB, name)
	return sqlDB.Close()
}

// RenameVectorDb 更改指定名称的 VectorList 的数据库名称
func RenameVectorDb(oldName string, newName string) error {
	// 1. 检查并关闭旧的数据库连接（如果已打开）
	err := CloseVectorDb(oldName)
	if err != nil {
		return fmt.Errorf("failed to close vector database: %w", err)
	}

	// 2. 获取旧的 VectorList 记录
	var oldList model.VectorList
	result := vectorListDb.Where("name = ?", oldName).First(&oldList)
	if result.Error != nil {
		return fmt.Errorf("failed to find old vector list: %w", result.Error)
	}
	//3. 删除旧的数据库文件
	oldDbPath := libs.GetVectorPath(oldName)
	// 4. 构建新的 DbPath
	newDbPath := libs.GetVectorPath(newName)
	if libs.PathExists(oldDbPath) {
		err := os.Rename(oldDbPath, newDbPath)
		if err != nil {
			return fmt.Errorf("failed to move old database file to new location: %w", err)
		}
	}

	// 5. 更新 VectorList 记录中的 DbPath 和 Name
	oldList.Name = newName
	oldList.DbPath = newDbPath
	result = vectorListDb.Save(&oldList)
	if result.Error != nil {
		return fmt.Errorf("failed to update vector list: %w", result.Error)
	}

	return nil
}

// UpdateVector 更新指定名称的 VectorList 记录
func UpdateVector(name string, updates map[string]interface{}) error {
	result := vectorListDb.Model(&model.VectorList{}).Where("name = ?", name).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("vector list not found")
	}

	return nil
}
func GetVectorList() []model.VectorList {
	var vectorList []model.VectorList
	vectorListDb.Find(&vectorList)
	return vectorList
}

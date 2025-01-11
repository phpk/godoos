package vector

import (
	"encoding/json"
	"fmt"
	"godo/ai/server"
	"godo/libs"
	"godo/model"
	"godo/office"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func InitMonitorVector() {
	// 确保数据库连接正常
	// 确保数据库连接正常
	db, err := model.Db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	list, err := GetVectorList()
	if err != nil {
		fmt.Println("GetVectorList error:", err)
		return
	}
	if len(list) == 0 {
		log.Println("no vector db found, creating a new one")
	}
	//log.Printf("init monitor:%v", list)
	for _, v := range list {
		MapFilePathMonitors[v.FilePath] = v.ID
	}
	go InitMonitor()
}
func HandlerCreateKnowledge(w http.ResponseWriter, r *http.Request) {
	var req model.VecList
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		libs.ErrorMsg(w, "the chat request error:"+err.Error())
		return
	}
	if req.FilePath == "" {
		libs.ErrorMsg(w, "file path is empty")
		return
	}
	basePath, err := libs.GetOsDir()
	if err != nil {
		libs.ErrorMsg(w, "get vector db path error:"+err.Error())
		return
	}
	req.FilePath = filepath.Join(basePath, req.FilePath)
	if !libs.PathExists(req.FilePath) {
		libs.ErrorMsg(w, "the knowledge path is not exists")
		return
	}
	fileInfo, err := os.Stat(req.FilePath)
	if err != nil {
		libs.ErrorMsg(w, "get vector db path error:"+err.Error())
		return
	}
	if !fileInfo.IsDir() {
		libs.ErrorMsg(w, "the knowledge path is not dir")
		return
	}
	knowledgeFilePath := filepath.Join(req.FilePath, ".knowledge")
	if libs.PathExists(knowledgeFilePath) {
		libs.ErrorMsg(w, "the knowledgeId already exists")
		return
	}
	id, err := CreateVector(req)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	libs.SuccessMsg(w, id, "create vector success")
}
func HandlerAskKnowledge(w http.ResponseWriter, r *http.Request) {
	var req model.AskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		libs.ErrorMsg(w, "the chat request error:"+err.Error())
		return
	}
	if req.ID == 0 {
		libs.ErrorMsg(w, "knowledgeId is empty")
		return
	}
	var knowData model.VecList
	if err := model.Db.First(&knowData, req.ID).Error; err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	var filterDocs []string
	filterDocs = append(filterDocs, req.Input)
	// 获取嵌入向量
	resList, err := server.GetEmbeddings(knowData.Engine, knowData.EmbeddingModel, filterDocs)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	res, err := model.AskDocument(req.ID, resList[0])
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	libs.SuccessMsg(w, res, "ask knowledge success")

}
func HandlerDelKnowledge(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		libs.ErrorMsg(w, "knowledgeId is empty")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		libs.ErrorMsg(w, "knowledgeId is not number")
		return
	}
	if id == 0 {
		libs.ErrorMsg(w, "knowledgeId is not number")
		return
	}
	if err := DeleteVector(uint(id)); err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	libs.SuccessMsg(w, nil, "delete knowledge success")
}

// CreateVector 创建一个新的 VectorList 记录
func CreateVector(data model.VecList) (uint, error) {
	if data.FilePath == "" {
		return 0, fmt.Errorf("file path is empty")
	}
	if data.Engine == "" {
		return 0, fmt.Errorf("engine is empty")
	}

	if !libs.PathExists(data.FilePath) {
		return 0, fmt.Errorf("file path does not exist")
	}
	if data.EmbeddingModel == "" {
		return 0, fmt.Errorf("embedding model is empty")
	}
	if data.EmbedSize == 0 {
		data.EmbedSize = 768
	}
	// Create the new VectorList
	result := model.Db.Create(&data)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to create vector list: %w", result.Error)
	}
	// 创建 .knowledge 文件并写入 knowledgeId
	knowledgeFilePath := filepath.Join(data.FilePath, ".knowledge")
	err := os.WriteFile(knowledgeFilePath, []byte(fmt.Sprintf("%d", data.ID)), 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to write knowledgeId to .knowledge file: %w", err)
	}

	// // 等待 AddWatchFolder 完成
	go AddWatchFolder(data.FilePath, data.ID, func() {
		go office.SetDocument(data.FilePath, data.ID)
	})

	return data.ID, nil
}

// DeleteVector 删除指定id的 VectorList 记录
func DeleteVector(id uint) error {
	var vectorList model.VecList
	if err := model.Db.First(&vectorList, id).Error; err != nil {
		return fmt.Errorf("failed to find vector list: %w", err)
	}
	// Delete .knowledge file
	knowledgeFilePath := filepath.Join(vectorList.FilePath, ".knowledge")
	if err := os.Remove(knowledgeFilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete .knowledge file: %w", err)
	}
	// Delete all .godoos. files in the directory and its subdirectories
	if err := deleteGodoosFiles(vectorList.FilePath); err != nil {
		return fmt.Errorf("failed to delete .godoos. files: %w", err)
	}
	//delete(MapFilePathMonitors, vectorList.FilePath)
	RemoveWatchFolder(vectorList.FilePath)
	return model.Db.Delete(&model.VecList{}, id).Error
}

func DeleteVectorFile(id uint, filePath string) error {
	var vectorList model.VecList
	if err := model.Db.First(&vectorList, id).Error; err != nil {
		return fmt.Errorf("failed to find vector list: %w", err)
	}
	// Delete file in database
	if err := model.Deletedocument(id, filePath); err != nil {
		return fmt.Errorf("failed to delete .godoos. files: %w", err)
	}
	newFileName := fmt.Sprintf(".godoos.%d.%s.json", id, filepath.Base(filePath))
	newFilePath := filepath.Join(filepath.Dir(filePath), newFileName)
	if libs.PathExists(newFilePath) {
		if err := os.Remove(newFilePath); err != nil {
			return fmt.Errorf("failed to delete .godoos. files: %w", err)
		}
	}
	return model.Db.Delete(&model.VecList{}, id).Error
}
func deleteGodoosFiles(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Base(path)[:7] == ".godoos." {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
		return nil
	})
}

// RenameVectorDb 更改指定名称的 VectorList 的数据库名称
func RenameVectorDb(oldName string, newName string) error {
	basePath, err := libs.GetOsDir()
	if err != nil {
		return fmt.Errorf("failed to find old vector list: %w", err)
	}

	// 获取旧的 VectorList 记录
	var oldList model.VecList
	oldPath := filepath.Join(basePath, oldName)
	if err := model.Db.Where("file_path = ?", oldPath).First(&oldList).Error; err != nil {
		return fmt.Errorf("failed to find old vector list: %w", err)
	}

	// 更新 VectorList 记录中的 FilePath
	newPath := filepath.Join(basePath, newName)
	if err := model.Db.Model(&model.VecList{}).Where("id = ?", oldList.ID).Update("file_path", newPath).Error; err != nil {
		return fmt.Errorf("failed to update vector list: %w", err)
	}
	// Update MapFilePathMonitors
	//delete(MapFilePathMonitors, oldPath)
	go RemoveWatchFolder(oldPath)
	//MapFilePathMonitors[newPath] = oldList.ID
	go AddWatchFolder(newPath, oldList.ID, nil)
	return nil
}

func GetVectorList() ([]model.VecList, error) {
	var vectorList []model.VecList
	if err := model.Db.Find(&vectorList).Error; err != nil {
		return nil, fmt.Errorf("failed to get vector list: %w", err)
	}
	return vectorList, nil
}

func GetVector(id uint) (model.VecList, error) {
	var vectorList model.VecList
	if err := model.Db.First(&vectorList, id).Error; err != nil {
		return vectorList, fmt.Errorf("failed to get vector: %w", err)
	}
	return vectorList, nil
}

// handleGodoosFile 处理 .godoos 文件
func handleGodoosFile(filePath string, knowledgeId uint) error {
	//log.Printf("========Handling .godoos file: %s", filePath)
	baseName := filepath.Base(filePath)
	// 检查文件后缀是否为 .exe
	if strings.HasSuffix(baseName, ".exe") {
		log.Printf("Skipping .exe file: %s", filePath)
		return nil
	}
	// 检查是否为 .godoos 文件
	if strings.HasPrefix(baseName, ".godoos.") {
		// 去掉 .godoos. 前缀和 .json 后缀
		// fileName := strings.TrimSuffix(strings.TrimPrefix(baseName, ".godoos."), ".json")
		// //提取实际文件名部分
		// actualFileName := extractFileName(fileName)

		// 读取文件内容
		content, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		// 解析 JSON 内容
		var doc office.Document
		err = json.Unmarshal(content, &doc)
		if err != nil {
			return err
		}

		// 检查 Split 是否为空
		if len(doc.Split) == 0 {
			return fmt.Errorf("invalid .godoos file: %s", filePath)
		}

		// 获取向量数据
		knowData, err := GetVector(knowledgeId)
		if err != nil {
			return err
		}
		basePath, err := libs.GetOsDir()
		if err != nil {
			return err
		}
		// 拼接文件名和内容
		// var filterDocs []string
		// for _, res := range doc.Split {
		// 	filterDocs = append(filterDocs, fmt.Sprintf("%s %s", actualFileName, res))
		// }
		// 获取嵌入向量
		resList, err := server.GetEmbeddings(knowData.Engine, knowData.EmbeddingModel, doc.Split)
		if err != nil {
			return err
		}

		// 检查嵌入向量长度是否匹配
		if len(resList) != len(doc.Split) {
			return fmt.Errorf("invalid file len: %s, expected %d embeddings, got %d", filePath, len(doc.Split), len(resList))
		}

		// 创建 VecDoc 列表
		var vectordocs []model.VecDoc
		for _, res := range doc.Split {
			//log.Printf("Adding document: %s", res)
			vectordoc := model.VecDoc{
				Content:  res,
				FilePath: strings.TrimPrefix(doc.RePath, basePath),
				FileName: doc.Title,
				ListID:   knowledgeId,
			}
			vectordocs = append(vectordocs, vectordoc)
		}

		// 添加文档
		err = model.Adddocument(knowledgeId, vectordocs, resList)
		return err
	} else {
		// 处理非 .godoos 文件
		if baseName[:1] != "." {
			office.ProcessFile(filePath, knowledgeId)
		}
		return nil
	}
}

// func extractFileName(fileName string) string {
// 	// 假设文件名格式为：21.GodoOS企业版介绍
// 	parts := strings.SplitN(fileName, ".", 3)
// 	if len(parts) < 2 {
// 		return fileName
// 	}
// 	return parts[1]
// }

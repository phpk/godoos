package vector

import (
	"encoding/json"
	"fmt"
	"godo/ai/server"
	"godo/libs"
	"godo/office"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

var MapFilePathMonitors = map[string]uint{}

func FolderMonitor() {
	basePath, err := libs.GetOsDir()
	if err != nil {
		log.Printf("Error getting base path: %s", err.Error())
		return
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Error creating watcher: %s", err.Error())
		return
	}
	defer watcher.Close()

	// 递归添加所有子目录
	addRecursive(basePath, watcher)

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Println("error:", err)
					return
				}
				//log.Println("event:", event)
				filePath := event.Name
				result, knowledgeId := shouldProcess(filePath)
				//log.Printf("result:%d,knowledgeId:%d", result, knowledgeId)
				if result > 0 {
					info, err := os.Stat(filePath)
					if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
						log.Println("modified file:", filePath)
						if !info.IsDir() {
							handleGodoosFile(filePath, knowledgeId)
						}
					}
					if event.Has(fsnotify.Create) || event.Has(fsnotify.Rename) {
						// 处理创建或重命名事件，添加新目录
						if err == nil && info.IsDir() {
							addRecursive(filePath, watcher)
						}
					}
					if event.Has(fsnotify.Remove) {
						// 处理删除事件，移除目录
						if err == nil && info.IsDir() {
							watcher.Remove(filePath)
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add(basePath)
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}

func shouldProcess(filePath string) (int, uint) {
	// 规范化路径
	filePath = filepath.Clean(filePath)

	// 检查文件路径是否在 MapFilePathMonitors 中
	for path, id := range MapFilePathMonitors {
		if id < 1 {
			return 0, 0
		}
		path = filepath.Clean(path)
		if filePath == path {
			return 1, id // 完全相等
		}
		if strings.HasPrefix(filePath, path+string(filepath.Separator)) {
			return 2, id // 包含
		}
	}
	return 0, 0 // 不存在
}

func addRecursive(path string, watcher *fsnotify.Watcher) {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error walking path %s: %v", path, err)
			return err
		}
		if info.IsDir() {
			result, _ := shouldProcess(path)
			if result > 0 {
				if err := watcher.Add(path); err != nil {
					log.Printf("Error adding path %s to watcher: %v", path, err)
					return err
				}
				log.Printf("Added path %s to watcher", path)
			}

		}
		return nil
	})
	if err != nil {
		log.Printf("Error adding recursive paths: %v", err)
	}
}

func handleGodoosFile(filePath string, knowledgeId uint) error {
	log.Printf("========Handling .godoos file: %s", filePath)
	baseName := filepath.Base(filePath)
	if baseName[:8] != ".godoos." {
		if baseName[:1] != "." {
			office.ProcessFile(filePath, knowledgeId)
		}
		return nil
	}
	var doc office.Document
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &doc)
	if err != nil {
		return err
	}
	if len(doc.Split) == 0 {
		return fmt.Errorf("invalid .godoos file: %s", filePath)
	}
	knowData := GetVector(knowledgeId)
	resList, err := server.GetEmbeddings(knowData.Engine, knowData.EmbeddingModel, doc.Split)
	if err != nil {
		return err
	}
	if len(resList) != len(doc.Split) {
		return fmt.Errorf("invalid file len: %s, expected %d embeddings, got %d", filePath, len(doc.Split), len(resList))
	}
	// var vectordocs []model.Vectordoc
	// for i, res := range resList {
	// 	//log.Printf("res: %v", res)
	// 	vectordoc := model.Vectordoc{
	// 		Content:     doc.Split[i],
	// 		Embed:       res,
	// 		FilePath:    filePath,
	// 		KnowledgeID: knowledgeId,
	// 		Pos:         fmt.Sprintf("%d", i),
	// 	}
	// 	vectordocs = append(vectordocs, vectordoc)
	// }
	// result := vectorListDb.Create(&vectordocs)
	// if result.Error != nil {
	// 	return result.Error
	// }
	return nil
}

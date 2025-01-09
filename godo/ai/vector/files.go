package vector

import (
	"fmt"
	"godo/libs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	MapFilePathMonitors = map[string]uint{}
	watcher             *fsnotify.Watcher
	fileQueue           = make(chan string, 100) // 队列大小可以根据需要调整
	numWorkers          = 3                      // 工作协程的数量
	wg                  sync.WaitGroup
	syncingKnowledgeIds = make(map[uint]syncingStats) // 记录正在同步的 knowledgeId 及其同步状态
	syncMutex           sync.Mutex                    // 保护 syncingKnowledgeIds 的互斥锁
	renameMap           = make(map[string]string)     // 临时映射存储 Remove 事件的路径
	renameMutex         sync.Mutex                    // 保护 renameMap 的互斥锁
	watcherMutex        sync.Mutex                    // 保护 watcher 的互斥锁
)

type syncingStats struct {
	totalFiles     int
	processedFiles int
}

func InitMonitor() {
	var err error
	watcherMutex.Lock()
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Error creating watcher: %s", err.Error())
	}
	watcherMutex.Unlock()
	go FolderMonitor()
	go startWatching()

	// 启动 worker
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker()
	}
}

func startWatching() {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				log.Println("error event")
				return
			}
			filePath := filepath.Clean(event.Name)
			result, exists := shouldProcess(filePath)
			if result > 0 {
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
					log.Printf("Event: %v, File: %s", event.Op, filePath)
					if isFileComplete(filePath) {
						// 将文件路径放入队列
						fileQueue <- filePath
					}
				}
				if event.Has(fsnotify.Create) {
					if info, err := os.Stat(filePath); err == nil && info.IsDir() {
						addRecursive(filePath, watcher)
					}
					// 检查是否是重命名事件
					handleRenameCreateEvent(event)
				}
				if event.Has(fsnotify.Remove) {
					//log.Printf("Event: %v, File: %s,exists:%d", event.Op, filePath, exists)
					isDir := true
					newFileName := fmt.Sprintf(".godoos.%d.%s.json", result, filepath.Base(filePath))
					newFilePath := filepath.Join(filepath.Dir(filePath), newFileName)
					if libs.PathExists(newFilePath) {
						isDir = false
					}
					if isDir {
						watcherMutex.Lock()
						if watcher != nil {
							watcher.Remove(filePath)
						}
						watcherMutex.Unlock()
					}
					if exists == 1 {
						err := DeleteVector(result)
						if err != nil {
							log.Printf("Error deleting vector %d: %v", result, err)
						}
					}
					if exists == 2 && !isDir {
						err := DeleteVectorFile(result, filePath)
						if err != nil {
							log.Printf("Error deleting vector file %d: %v", result, err)
						}
					}
					//存储 Remove 事件的路径
					handleRenameRemoveEvent(event)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func handleRenameRemoveEvent(event fsnotify.Event) {
	renameMutex.Lock()
	defer renameMutex.Unlock()
	//log.Printf("handleRenameRemoveEvent: %v, File: %s", event.Op, event.Name)
	renameMap[event.Name] = event.Name
}

func handleRenameCreateEvent(event fsnotify.Event) {
	renameMutex.Lock()
	defer renameMutex.Unlock()
	//log.Printf("handleRenameCreateEvent: %v, File: %s", event.Op, event.Name)
	// 规范化路径
	newPath := filepath.Clean(event.Name)

	// 检查是否是重命名事件
	for oldPath := range renameMap {
		if oldPath != "" {
			// 找到对应的 Remove 事件
			oldPathClean := filepath.Clean(oldPath)
			if oldPathClean == newPath {
				//log.Printf("File renamed from %s to %s", oldPath, newPath)

				// 更新 MapFilePathMonitors
				for path, id := range MapFilePathMonitors {
					if path == oldPathClean {
						delete(MapFilePathMonitors, path)
						MapFilePathMonitors[newPath] = id
						log.Printf("Updated MapFilePathMonitors: %s -> %s", oldPathClean, newPath)
						break
					}
				}

				// 更新 watcher
				watcherMutex.Lock()
				if watcher != nil {
					if err := watcher.Remove(oldPathClean); err != nil {
						log.Printf("Error removing old path %s from watcher: %v", oldPathClean, err)
					}
					if err := watcher.Add(newPath); err != nil {
						log.Printf("Error adding new path %s to watcher: %v", newPath, err)
					}
				}
				watcherMutex.Unlock()

				// 如果是目录，递归更新子目录
				if info, err := os.Stat(newPath); err == nil && info.IsDir() {
					addRecursive(newPath, watcher)
				}

				// 清除临时映射中的路径
				delete(renameMap, oldPath)
				break
			}
		}
	}
}

func worker() {
	defer wg.Done()
	for filePath := range fileQueue {
		knowledgeId, exists := shouldProcess(filePath)
		if exists == 0 {
			log.Printf("File path %s is not being monitored", filePath)
			continue
		}

		// 更新已处理文件数
		syncMutex.Lock()
		if stats, ok := syncingKnowledgeIds[knowledgeId]; ok {
			stats.processedFiles++
			syncingKnowledgeIds[knowledgeId] = stats
		}
		syncMutex.Unlock()

		err := handleGodoosFile(filePath, knowledgeId)
		if err != nil {
			log.Printf("Error handling file %s: %v", filePath, err)
		}
	}
}

func FolderMonitor() {
	basePath, err := libs.GetOsDir()
	if err != nil {
		log.Printf("Error getting base path: %s", err.Error())
		return
	}

	// 递归添加所有子目录
	addRecursive(basePath, watcher)

	// Add a path.
	watcherMutex.Lock()
	if watcher != nil {
		err = watcher.Add(basePath)
		if err != nil {
			log.Fatal(err)
		}
	}
	watcherMutex.Unlock()

	// Block main goroutine forever.
	<-make(chan struct{})
}

func AddWatchFolder(folderPath string, knowledgeId uint, callback func()) error {
	if watcher == nil {
		InitMonitor()
	}
	// 规范化路径
	folderPath = filepath.Clean(folderPath)

	// 检查文件夹是否存在
	if !libs.PathExists(folderPath) {
		return fmt.Errorf("folder path does not exist: %s", folderPath)
	}

	// 检查文件夹是否已经存在于监视器中
	if _, exists := MapFilePathMonitors[folderPath]; exists {
		return fmt.Errorf("folder path is already being monitored: %s", folderPath)
	}

	// 递归添加所有子目录
	addRecursive(folderPath, watcher)

	// 计算总文件数
	totalFiles, err := countFiles(folderPath)
	if err != nil {
		return fmt.Errorf("failed to count files in folder path: %w", err)
	}

	// 更新 syncingKnowledgeIds
	syncMutex.Lock()
	syncingKnowledgeIds[knowledgeId] = syncingStats{
		totalFiles:     totalFiles,
		processedFiles: 0,
	}
	syncMutex.Unlock()

	// 更新 MapFilePathMonitors
	MapFilePathMonitors[folderPath] = knowledgeId

	// 添加文件夹路径到监视器
	err = watcher.Add(folderPath)
	if err != nil {
		return fmt.Errorf("failed to add folder path to watcher: %w", err)
	}

	// 调用回调函数
	if callback != nil {
		callback()
	}

	log.Printf("Added folder path %s to watcher with knowledgeId %d", folderPath, knowledgeId)
	return nil
}

// RemoveWatchFolder 根据路径删除观察文件夹
func RemoveWatchFolder(folderPath string) error {
	// 规范化路径
	folderPath = filepath.Clean(folderPath)

	// 检查文件夹是否存在于监视器中
	knowledgeId, exists := MapFilePathMonitors[folderPath]
	if !exists {
		return fmt.Errorf("folder path is not being monitored: %s", folderPath)
	}

	// 从 watcher 中移除路径
	watcherMutex.Lock()
	if watcher != nil {
		err := watcher.Remove(folderPath)
		if err != nil {
			return fmt.Errorf("failed to remove folder path from watcher: %w", err)
		}
	}
	watcherMutex.Unlock()

	// 递归移除所有子目录
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error walking path %s: %v", path, err)
			return err
		}
		if info.IsDir() {
			result, _ := shouldProcess(path)
			if result > 0 {
				// 从 watcher 中移除路径
				watcherMutex.Lock()
				if watcher != nil {
					err := watcher.Remove(path)
					if err != nil {
						log.Printf("Error removing path %s from watcher: %v", path, err)
						return err
					}
				}
				watcherMutex.Unlock()
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to remove folder path from watcher: %w", err)
	}

	// 从 MapFilePathMonitors 中删除条目
	delete(MapFilePathMonitors, folderPath)

	// 从 syncingKnowledgeIds 中删除条目
	syncMutex.Lock()
	delete(syncingKnowledgeIds, knowledgeId)
	syncMutex.Unlock()

	log.Printf("Removed folder path %s from watcher with knowledgeId %d", folderPath, knowledgeId)
	return nil
}

func shouldProcess(filePath string) (uint, int) {
	// 规范化路径
	filePath = filepath.Clean(filePath)

	// 检查文件路径是否在 MapFilePathMonitors 中
	for path, id := range MapFilePathMonitors {
		if id < 1 {
			return 0, 0
		}
		path = filepath.Clean(path)
		if filePath == path {
			return id, 1 // 完全相等
		}
		if strings.HasPrefix(filePath, path+string(filepath.Separator)) {
			return id, 2 // 包含
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

// countFiles 递归计算文件夹中的文件数
func countFiles(folderPath string) (int, error) {
	var count int
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetSyncPercentage 计算并返回同步百分比
func GetSyncPercentage(knowledgeId uint) float64 {
	syncMutex.Lock()
	defer syncMutex.Unlock()
	if stats, ok := syncingKnowledgeIds[knowledgeId]; ok {
		if stats.totalFiles == 0 {
			return 0.0
		}
		return float64(stats.processedFiles) / float64(stats.totalFiles) * 100
	}
	return 0.0
}

// isFileComplete 检查文件是否已经完全创建
func isFileComplete(filePath string) bool {
	// 等待一段时间确保文件已经完全创建
	time.Sleep(100 * time.Millisecond)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); err != nil {
		log.Printf("File %s does not exist: %v", filePath, err)
		return false
	}

	// 检查文件大小是否达到预期
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Printf("Error stat file %s: %v", filePath, err)
		return false
	}
	// 例如，检查文件大小是否大于某个阈值
	if fileInfo.Size() == 0 {
		log.Printf("File %s is empty", filePath)
		return false
	}
	if fileInfo.IsDir() {
		log.Printf("File %s is a directory", filePath)
		return false
	}
	return true
}

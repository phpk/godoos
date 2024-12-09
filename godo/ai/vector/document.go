package vector

import (
	"godo/libs"
	"godo/office"
	"log"
	"path"
	"sync"

	"github.com/fsnotify/fsnotify"
)

const numWorkers = 5 // 设置 worker 数量

func MonitorFolder(folderPath string) {
	if !libs.PathExists(folderPath) {
		return
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create watcher for folder %s: %v", folderPath, err)
	}
	defer watcher.Close()

	fileQueue := make(chan string, 100) // 创建文件路径队列
	var wg sync.WaitGroup

	// 启动 worker goroutine
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range fileQueue {
				handleGodoosFile(filePath)
			}
		}()
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&(fsnotify.Create|fsnotify.Write) != 0 { // 监听创建和修改事件
					baseName := path.Base(event.Name)
					if baseName[:8] == ".godoos." { // 检查文件名是否以 .godoos 开头
						log.Printf("Detected .godoos file: %s", event.Name)
						fileQueue <- event.Name // 将文件路径放入队列
					} else {
						if baseName[:1] == "." {
							return
						}
						office.ProcessFile(event.Name)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Error watching folder %s: %v", folderPath, err)
			}
		}
	}()

	err = watcher.Add(folderPath)
	if err != nil {
		log.Fatalf("Failed to add folder %s to watcher: %v", folderPath, err)
	}

	// 关闭文件队列
	go func() {
		<-done
		close(fileQueue)
		wg.Wait()
	}()

	<-done
}

func handleGodoosFile(filePath string) {
	// 在这里添加对 .godoos 文件的具体处理逻辑
	log.Printf("Handling .godoos file: %s", filePath)
}

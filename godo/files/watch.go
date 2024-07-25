package files

import (
	"context"
	"encoding/json"
	"fmt"
	"godo/libs"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// FileChangeEvent 封装文件系统变更事件为可序列化的结构体
type FileChangeEvent struct {
	Action string `json:"action"` // create, modify, delete, rename
	Path   string `json:"path"`   // 文件或目录的路径
}

// sseStream 结构体用于维护SSE连接和事件发送
type sseStream struct {
	conn         http.ResponseWriter
	mu           sync.Mutex
	ctx          context.Context
	cancel       context.CancelFunc
	eventChannel chan FileChangeEvent
}

func newSSEStream(w http.ResponseWriter) *sseStream {
	stream := &sseStream{
		conn:         w,
		eventChannel: make(chan FileChangeEvent),
	}
	stream.ctx, stream.cancel = context.WithCancel(context.Background())
	return stream
}

// sendEvent 向SSE连接发送事件
func (s *sseStream) sendEvent(event FileChangeEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn == nil {
		return
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling event: %v", err)
		return
	}

	fmt.Fprintf(s.conn, "id: %d\n", time.Now().UnixNano())
	fmt.Fprintf(s.conn, "event: change\n")
	fmt.Fprintf(s.conn, "data: %s\n\n", data)
}

// WatchFileSystem watches for changes in a directory and returns a channel for events
func WatchFileSystem(dirPath string) (<-chan fsnotify.Event, <-chan error, context.CancelFunc) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		errChan := make(chan error, 1)
		errChan <- err
		close(errChan)
		return nil, errChan, func() {}
	}

	events := make(chan fsnotify.Event)
	errors := make(chan error, 1)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	// Start watching for events
	go func() {
		defer func() {
			wg.Wait()
			close(events)
			close(errors)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := watcher.Add(dirPath); err != nil {
				errors <- fmt.Errorf("failed to watch directory %q: %v", dirPath, err)
				cancel()
				return
			}
			for {
				select {
				case <-ctx.Done():
					return
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					events <- event
				case err := <-watcher.Errors:
					errors <- err
				}
			}
		}()

		// Handle context cancellation
		<-ctx.Done()
		if err := watcher.Close(); err != nil {
			log.Printf("Error closing watcher: %v", err)
		}
	}()

	return events, errors, cancel
}
func WatchHandler(w http.ResponseWriter, r *http.Request) {
	dirToWatch := r.URL.Query().Get("filePath")
	basePath, err := libs.GetOsDir()
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	filePath := filepath.Join(basePath, dirToWatch)
	if !libs.PathExists(filePath) {
		libs.HTTPError(w, http.StatusNotFound, "filepath is not exist!")
		return
	}
	StartSSEStream(w, r, filePath)
}

// StartSSEStream 处理SSE连接，持续监听并发送事件
func StartSSEStream(w http.ResponseWriter, r *http.Request, dirPath string) {
	stream := newSSEStream(w)
	defer stream.cancel()

	setupSSEHeaders(w)

	events, errors, cancel := WatchFileSystem(dirPath)
	defer cancel()

	go handleEvents(stream, events, errors)

	// 等待客户端断开连接或请求上下文结束
	<-r.Context().Done()
}

// setupSSEHeaders 设置SSE响应头
func setupSSEHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
}

// handleEvents 处理接收到的文件系统事件和错误
func handleEvents(stream *sseStream, events <-chan fsnotify.Event, errors <-chan error) {
	for {
		select {
		case <-stream.ctx.Done():
			return
		case event := <-events:
			stream.sendEvent(convertEvent(event))
		case err := <-errors:
			log.Printf("Error from fsnotify: %v", err)
		}
	}
}

// convertEvent 将fsnotify.Event转换为FileChangeEvent
func convertEvent(fsEvent fsnotify.Event) FileChangeEvent {
	var action string
	switch fsEvent.Op {
	case fsnotify.Create:
		action = "create"
	case fsnotify.Write:
		action = "modify"
	case fsnotify.Remove:
		action = "delete"
	case fsnotify.Rename:
		action = "rename"
	default:
		action = "unknown"
	}
	return FileChangeEvent{Action: action, Path: fsEvent.Name}
}

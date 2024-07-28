package store

import (
	"encoding/json"
	"godo/files"
	"godo/libs"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type ProgressReader struct {
	reader io.Reader
	total  int64
	err    error
}

type DownloadStatus struct {
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	Url         string  `json:"url"`
	Current     int64   `json:"current"`
	Size        int64   `json:"size"`
	Speed       float64 `json:"speed"`
	Progress    int     `json:"progress"`
	Downloading bool    `json:"downloading"`
	Done        bool    `json:"done"`
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.reader.Read(p)
	pr.err = err
	pr.total += int64(n)
	return
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	log.Printf("Download url: %s", url)

	// 获取下载目录，这里假设从请求参数中获取，如果没有则使用默认值
	downloadDir := libs.GetCacheDir()
	if downloadDir == "" {
		downloadDir = "./downloads"
	}

	// 拼接完整的文件路径
	fileName := filepath.Base(url)
	filePath := filepath.Join(downloadDir, fileName)

	// 开始下载
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to get file: %v", err)
		http.Error(w, "Failed to get file", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 检查文件是否已存在且大小一致
	if fileInfo, err := os.Stat(filePath); err == nil {
		if fileInfo.Size() == resp.ContentLength {
			// 文件已存在且大小一致，无需下载
			runDir := libs.GetRunDir()
			err := files.HandlerFile(filePath, runDir)
			if err != nil {
				log.Printf("Error moving file: %v", err)
			}
			libs.SuccessMsg(w, "success", "File already exists and is of correct size")
			return
		} else {
			// 重新打开响应体以便后续读取
			resp.Body = http.NoBody
			resp, err = http.Get(url)
			if err != nil {
				log.Printf("Failed to get file: %v", err)
				http.Error(w, "Failed to get file", http.StatusInternalServerError)
				return
			}
		}
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 使用ProgressReader来跟踪进度
	pr := &ProgressReader{reader: resp.Body}

	// 启动定时器来报告进度
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Printf("Streaming unsupported")
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	go func() {
		for {
			<-ticker.C
			rp := &DownloadStatus{
				Name:        fileName,
				Path:        filePath,
				Url:         url,
				Current:     pr.total,
				Size:        resp.ContentLength,
				Speed:       0, // 这里可以计算速度，但为了简化示例，我们暂时设为0
				Progress:    int(100 * (float64(pr.total) / float64(resp.ContentLength))),
				Downloading: pr.err == nil && pr.total < resp.ContentLength,
				Done:        pr.total == resp.ContentLength,
			}
			if pr.err != nil || rp.Done {
				rp.Downloading = false
				//log.Printf("Download complete: %s", filePath)
				runDir := libs.GetRunDir()
				err := files.HandlerFile(filePath, runDir)
				if err != nil {
					log.Printf("Error moving file: %v", err)
				}
				break
			}
			if w != nil {
				jsonBytes, err := json.Marshal(rp)
				if err != nil {
					log.Printf("Error marshaling DownloadStatus to JSON: %v", err)
					continue
				}
				w.Write(jsonBytes)
				w.Write([]byte("\n"))
				flusher.Flush()
			} else {
				log.Println("ResponseWriter is nil, cannot send progress")
			}
		}
	}()

	// 将响应体的内容写入文件
	_, err = io.Copy(file, pr)
	if err != nil {
		log.Printf("Failed to write file: %v", err)
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}
	libs.SuccessMsg(w, "success", "Download complete")
}

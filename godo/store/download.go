package store

import (
	"encoding/json"
	"fmt"
	"godo/files"
	"godo/libs"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/cavaliergopher/grab/v3"
)

const (
	concurrency = 6 // 并发下载数
)

var downloads = make(map[string]*grab.Response)
var downloadsMutex sync.Mutex

type FileProgress struct {
	Progress   int    `json:"progress"` // 将进度改为浮点数，以百分比表示
	IsFinished bool   `json:"is_finished"`
	Total      int64  `json:"total"`
	Current    int64  `json:"completed"`
	Status     string `json:"status"`
}

type ReqBody struct {
	DownloadUrl string `json:"url"`
	PkgUrl      string `json:"pkg"`
	Name        string `json:"name"`
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody ReqBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		libs.ErrorMsg(w, "first Decode request body error")
		return
	}
	downloadDir := libs.GetCacheDir()
	if downloadDir == "" {
		downloadDir = "./downloads"
	}
	if reqBody.Name == "" {
		libs.ErrorMsg(w, "the name is empty")
		return
	}
	if reqBody.DownloadUrl == "" && reqBody.PkgUrl == "" {
		libs.ErrorMsg(w, "the url is empty")
		return
	}
	paths := []string{}
	for _, url := range []string{reqBody.DownloadUrl, reqBody.PkgUrl} {
		if url == "" {
			continue
		}
		if rsp, ok := downloads[url]; ok {
			// 如果URL正在下载，跳过创建新的下载器实例
			go trackProgress(w, rsp, url)
			continue
		}
		// 创建新的下载器实例
		client := grab.NewClient()
		client.HTTPClient = &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: concurrency, // 可选，设置并发连接数
			},
		}
		filePath := filepath.Join(downloadDir, filepath.Base(url))
		paths = append(paths, filePath)
		log.Printf("filePath is %s", filePath)
		// 创建下载请求
		req, err := grab.NewRequest(filePath, url)
		if err != nil {
			libs.ErrorMsg(w, "Invalid download URL")
			return
		}

		resp := client.Do(req)
		if fileInfo, err := os.Stat(filePath); err == nil {
			// 文件存在，检查文件大小是否与远程资源匹配
			if fileInfo.Size() == resp.Size() { // 这里的req.Size需要从下载请求中获取，或通过其他方式预知
				log.Printf("File %s already exists and is up-to-date.", filePath)
				continue
			}
		}
		downloads[url] = resp
		//log.Printf("Download urls: %v\n", reqBody.DownloadUrl)

		// // 跟踪进度
		go trackProgress(w, resp, url)
		// 等待下载完成并检查错误
		if err := resp.Err(); err != nil {
			libs.ErrorMsg(w, "Download failed")
			return
		}
	}
	if len(paths) > 0 {
		// 解压文件
		err := HandlerZipFiles(paths, reqBody.Name)
		if err != nil {
			libs.ErrorMsg(w, "Decompress failed:"+err.Error())
			return
		}
	}
	installInfo, err := Installation(reqBody.Name)
	if err != nil {
		libs.ErrorData(w, installInfo, "the install.json is error:"+err.Error())
		return
	}
	libs.SuccessMsg(w, installInfo, "install the app success!")
}
func HandlerZipFiles(paths []string, name string) error {
	runDir := libs.GetRunDir()
	targetDir := filepath.Join(runDir, name)
	downloadDir := libs.GetCacheDir()
	for _, filePath := range paths {
		zipPath, err := files.Decompress(filePath, downloadDir)
		if err != nil {
			return fmt.Errorf("decompress failed: %v", err)
		}
		err = files.CopyResource(zipPath, targetDir)
		if err != nil {
			return fmt.Errorf("copyResource failed")
		}
		err = os.RemoveAll(zipPath)
		if err != nil {
			return fmt.Errorf("removeAll failed")
		}

	}

	return nil
}
func trackProgress(w http.ResponseWriter, resp *grab.Response, md5url string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered panic in trackProgress: %v", r)
		}
		downloadsMutex.Lock()
		defer downloadsMutex.Unlock()
		delete(downloads, md5url)
	}()
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Printf("Streaming unsupported")
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}
	for {
		select {
		case <-ticker.C:
			fp := FileProgress{
				Progress:   int(resp.Progress() * 100),
				IsFinished: resp.IsComplete(),
				Total:      resp.Size(),
				Current:    resp.BytesComplete(),
				Status:     "loading",
			}
			//log.Printf("Progress: %v", fp)
			if resp.IsComplete() && fp.Current == fp.Total {
				fp.Status = "loaded"
			}
			jsonBytes, err := json.Marshal(fp)
			if err != nil {
				log.Printf("Error marshaling FileProgress to JSON: %v", err)
				continue
			}
			if w != nil {
				io.WriteString(w, string(jsonBytes))
				w.Write([]byte("\n"))
				flusher.Flush()
			} else {
				log.Println("ResponseWriter is nil, cannot send progress")
			}
			if fp.Status == "loaded" {
				return
			}

		}
	}
}

// Predefined filename for the icon
const iconFilename = "icon.png"

// DownloadIcon downloads the icon from the given URL to the target path with a predefined filename.
func DownloadIcon(url, iconTargetPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	out, err := os.Create(iconTargetPath)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}

	return nil
}

// HandlerIcon handles the icon by downloading it and updating the installInfo.
func HandlerIcon(installInfo InstallInfo, targetPath string) (string, error) {
	var iconUrl string
	if url, err := url.Parse(installInfo.Icon); err == nil && url.Scheme != "" {
		// Download the icon using the predefined filename
		iconTargetPath := filepath.Join(targetPath, installInfo.Name, "_icon.png")
		if err := DownloadIcon(installInfo.Icon, iconTargetPath); err != nil {
			return "", fmt.Errorf("error downloading icon: %v", err)
		}
		iconUrl = "http://localhost:56780/static/" + installInfo.Name + "/" + installInfo.Name + "_icon.png"
	} else {
		iconUrl = "http://localhost:56780/static/" + installInfo.Name + "/" + installInfo.Icon
	}

	return iconUrl, nil
}

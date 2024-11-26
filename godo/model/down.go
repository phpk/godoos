package model

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"godo/libs"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cavaliergopher/grab/v3"
)

const (
	concurrency = 6 // 并发下载数
)

var downloads = make(map[string]*grab.Response)
var downloadsMutex sync.Mutex

//var cancelDownloads = make(map[string]context.CancelFunc)

func noticeSuccess(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	//log.Println("Download starting!")

}

func Download(w http.ResponseWriter, r *http.Request) {
	reqBody := ReqBody{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		libs.ErrorMsg(w, "first Decode request body error:"+err.Error())
		return
	}
	err = LoadConfig()
	if err != nil {
		libs.ErrorMsg(w, "Load config error")
		return
	}
	_, exitsModel := GetModel(reqBody.Model)
	if exitsModel {
		noticeSuccess(w)
		return

	}
	if reqBody.Info.From == "ollama" && reqBody.Info.Engine == "ollama" {
		setOllamaInfo(w, r, reqBody)
		return
	}

	var paths []string
	var tsize int64
	for _, urls := range reqBody.Info.URL {
		urls = replaceUrl(urls)
		if !strings.HasPrefix(strings.ToLower(urls), "http://") && !strings.HasPrefix(strings.ToLower(urls), "https://") {
			fileInfo, err := os.Stat(urls)
			if err != nil {
				libs.ErrorMsg(w, "Get model path error")
				return
			}
			tsize += fileInfo.Size()
			paths = append(paths, urls)
			continue
		}
		filePath, err := GetModelPath(urls, reqBody.Model, reqBody.Type)
		//log.Printf("filePath is %s", filePath)
		if err != nil {
			libs.ErrorMsg(w, "Get model path error")
			return
		}
		paths = append(paths, filePath)
		md5url := md5Url(urls)
		if rsp, ok := downloads[md5url]; ok {
			// 如果URL正在下载，跳过创建新的下载器实例
			go trackProgress(w, rsp, md5url)
			return
		}
		// 创建新的下载器实例
		client := grab.NewClient()
		client.HTTPClient = &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: concurrency, // 可选，设置并发连接数
			},
		}
		log.Printf("filePath is %s", filePath)
		// 创建下载请求
		req, err := grab.NewRequest(filePath, urls)
		if err != nil {
			libs.ErrorMsg(w, "Invalid download URL")
			return
		}
		resp := client.Do(req)
		downloads[md5url] = resp
		//log.Printf("Download urls: %v\n", reqBody.DownloadUrl)

		// // 跟踪进度
		go trackProgress(w, resp, md5url)
		tsize += resp.Size()
		// 等待下载完成并检查错误
		if err := resp.Err(); err != nil {
			libs.ErrorMsg(w, "Download failed")
			return
		}
	}
	delUrls(reqBody.Info.URL)
	if tsize <= 0 {
		libs.ErrorMsg(w, "download size is zero")
		return
	}
	reqBody.Info.Path = paths
	reqBody.Status = "success"
	reqBody.CreatedAt = time.Now()
	// reqBody.Info["tsize"] = tsize
	// reqBody.Info["size"] = humanReadableSize(tsize)
	if reqBody.Info.From == "network" && reqBody.Info.Engine == "ollama" {
		ConvertOllama(w, r, reqBody)
		// reqBody.From = "ollama"
		// reqBody.Paths = []string{}
	}

	if err := SetModel(reqBody); err != nil {
		libs.ErrorMsg(w, "Set model error")
		return
	}

	noticeSuccess(w)
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
				Progress:   resp.Progress(),
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
func md5Url(url string) string {
	hasher := md5.New()
	hasher.Write([]byte(url))
	return hex.EncodeToString(hasher.Sum(nil))
}
func delUrls(reqUrl []string) {
	if len(reqUrl) > 0 {
		downloadsMutex.Lock()
		defer downloadsMutex.Unlock()
		for _, urls := range reqUrl {
			urls = replaceUrl(urls)
			md5url := md5Url(urls)
			delete(downloads, md5url)

		}
	}
}
func DeleteFileHandle(w http.ResponseWriter, r *http.Request) {
	var reqBody ReqBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		libs.ErrorMsg(w, "Decode request body error: ")
		return
	}
	err = LoadConfig()
	if err != nil {
		libs.ErrorMsg(w, "Load config error: ")
		return
	}
	if err := DeleteModel(reqBody.Model); err != nil {
		libs.ErrorMsg(w, "Error deleting model")
		return
	}
	if reqBody.Info.Engine == "ollama" {
		postQuery := map[string]interface{}{"name": reqBody.Model}
		url := GetOllamaUrl() + "/api/delete"

		ForwardHandler(w, r, postQuery, url, "DELETE")
		return
	}
	delUrls(reqBody.Info.URL)

	// 尝试删除目录，注意这会递归删除目录下的所有内容
	//dirPath := filepath.Dir(filePath)
	dirPath, err := GetModelDir(reqBody.Model)
	if err != nil {
		libs.ErrorMsg(w, "GetModelDir error")
		return
	}
	//log.Printf("delete dirpath %v", dirPath)
	err = os.RemoveAll(dirPath)
	if err != nil && !os.IsNotExist(err) {
		libs.ErrorMsg(w, "Error removing directory")
		return
	} else if err == nil {
		log.Printf("Deleted directory: %s", dirPath)
	} else {
		// 如果目录不存在，这通常是可以接受的，不需要错误消息
		log.Printf("Directory does not exist: %s", dirPath)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"code": 0})

}

func replaceUrl(url string) string {
	return strings.ReplaceAll(url, "/blob/main/", "/resolve/main/")
}

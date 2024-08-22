package sys

import (
	"encoding/json"
	"fmt"
	"godo/libs"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/minio/selfupdate"
)

type OSInfo struct {
	Amd64 string `json:"amd64"`
	Arm64 string `json:"arm64"`
}

type VersionInfo struct {
	Version     string `json:"version"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Changelog   string `json:"changelog"`
	Windows     OSInfo `json:"windows"`
	Linux       OSInfo `json:"linux"`
	Darwin      OSInfo `json:"darwin"`
}

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
	Progress    float64 `json:"progress"`
	Downloading bool    `json:"downloading"`
	Done        bool    `json:"done"`
}
type UpdateAdReq struct {
	Img  string `json:"img"`
	Name string `json:"name"`
	Link string `json:"link"`
	Desc string `json:"desc"`
}
type UpdateVersionReq struct {
	Version string        `json:"version"`
	Url     string        `json:"url"`
	Name    string        `json:"name"`
	Desc    string        `json:"desc"`
	AdList  []UpdateAdReq `json:"adList"`
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.reader.Read(p)
	pr.err = err
	pr.total += int64(n)
	return
}
func UpdateAppHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	pr := &ProgressReader{reader: resp.Body}

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Printf("Streaming unsupported")
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}
	// update progress
	go func() {
		for {
			<-ticker.C
			rp := &DownloadStatus{
				Name:        filepath.Base(url),
				Path:        "",
				Url:         url,
				Current:     pr.total,
				Size:        resp.ContentLength,
				Speed:       0,
				Progress:    100 * (float64(pr.total) / float64(resp.ContentLength)),
				Downloading: pr.err == nil && pr.total < resp.ContentLength,
				Done:        pr.total == resp.ContentLength,
			}
			if pr.err != nil || pr.total == resp.ContentLength {
				break
			}
			if w != nil {
				jsonBytes, err := json.Marshal(rp)
				if err != nil {
					log.Printf("Error marshaling FileProgress to JSON: %v", err)
					continue
				}
				io.WriteString(w, string(jsonBytes))
				w.Write([]byte("\n"))
				flusher.Flush()
			} else {
				log.Println("ResponseWriter is nil, cannot send progress")
			}
		}
	}()

	var updateFile io.Reader = pr
	// apply update
	err = selfupdate.Apply(updateFile, selfupdate.Options{})
	if err != nil {
		if rerr := selfupdate.RollbackError(err); rerr != nil {
			http.Error(w, "update error:"+rerr.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	// 更新完成后发送响应给前端
	json.NewEncoder(w).Encode(map[string]bool{"updateCompleted": true})
}

func GetUpdateUrlHandler(w http.ResponseWriter, r *http.Request) {
	info, err := libs.GetSystemInfo()
	if err != nil {
		http.Error(w, "update error:"+err.Error(), http.StatusInternalServerError)
		return
	}
	updateUrl := "https://godocms.com/version?info=" + info
	res, err := http.Get(updateUrl)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			libs.ErrorMsg(w, err.Error())
			return
		}
		var updateInfo UpdateVersionReq
		err = json.Unmarshal(body, &updateInfo)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}
		json.NewEncoder(w).Encode(updateInfo)
	}
}

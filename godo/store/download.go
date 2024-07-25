package store

import (
	"context"
	"encoding/json"
	"godo/libs"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/cavaliergopher/grab/v3"
)

type DownloadStatus struct {
	resp        *grab.Response
	cancel      context.CancelFunc
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

const (
	concurrency = 6 // 并发下载数
)

// var downloads = make(map[string]*grab.Response)
var downloadsMutex sync.Mutex
var downloadList map[string]*DownloadStatus

func existsInDownloadList(url string) bool {
	_, ok := downloadList[url]
	return ok
}

func PauseDownload(url string) {
	downloadsMutex.Lock()
	defer downloadsMutex.Unlock()
	ds, ok := downloadList[url]
	if ds.Url == url && ok {
		if ds.cancel != nil {
			ds.cancel()
		}
		ds.resp = nil
		ds.Downloading = false
		ds.Speed = 0
	}

}

func ContinueDownload(url string) {
	ds, ok := downloadList[url]
	if ds.Url == url && ok {
		if !ds.Downloading && ds.resp == nil && !ds.Done {
			ds.Downloading = true

			req, err := grab.NewRequest(ds.Path, ds.Url)
			if err != nil {
				ds.Downloading = false
				return
			}
			ctx, cancel := context.WithCancel(context.Background())
			ds.cancel = cancel
			req = req.WithContext(ctx)
			client := grab.NewClient()
			client.HTTPClient = &http.Client{
				Transport: &http.Transport{
					MaxIdleConnsPerHost: concurrency, // 设置并发连接数
				},
			}
			//resp := grab.DefaultClient.Do(req)
			resp := client.Do(req)

			if resp != nil && resp.HTTPResponse != nil &&
				resp.HTTPResponse.StatusCode >= 200 && resp.HTTPResponse.StatusCode < 300 {
				ds.resp = resp
			} else {
				ds.Downloading = false
			}
		}
	}

}

func Download(url string) {
	chacheDir := libs.GetCacheDir()
	absPath := filepath.Join(chacheDir, filepath.Base(url))
	if !existsInDownloadList(url) {
		downloadList[url] = &DownloadStatus{
			resp:        nil,
			Name:        filepath.Base(url),
			Path:        absPath,
			Url:         url,
			Downloading: false,
		}
	}
	ContinueDownload(url)
}
func GetDownload(url string) *DownloadStatus {
	downloadsMutex.Lock()
	defer downloadsMutex.Unlock()
	ds, ok := downloadList[url]
	if ds.resp != nil && ok {
		ds.Current = ds.resp.BytesComplete()
		ds.Size = ds.resp.Size()
		ds.Speed = ds.resp.BytesPerSecond()
		ds.Progress = 100 * ds.resp.Progress()
		ds.Downloading = !ds.resp.IsComplete()
		ds.Done = ds.resp.Progress() == 1
		if !ds.Downloading {
			ds.resp = nil
		}
	}
	if ds.Done {
		delete(downloadList, url)
	}
	return ds
}
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	Download(url)
	ticker := time.NewTicker(500 * time.Millisecond)
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
			ds := GetDownload(url)
			jsonBytes, err := json.Marshal(ds)
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
			if ds.Done {
				return
			}
		}
	}()
}

package app

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/minio/selfupdate"
	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
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
	Transferred int64   `json:"transferred"`
	Size        int64   `json:"size"`
	Speed       float64 `json:"speed"`
	Progress    float64 `json:"progress"`
	Downloading bool    `json:"downloading"`
	Done        bool    `json:"done"`
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.reader.Read(p)
	pr.err = err
	pr.total += int64(n)
	return
}
func (a *App) UpdateApp(url string) (broken bool, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	pr := &ProgressReader{reader: resp.Body}

	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	// update progress
	go func() {
		for {
			<-ticker.C
			wruntime.EventsEmit(a.ctx, "updateApp", &DownloadStatus{
				Name:        filepath.Base(url),
				Path:        "",
				Url:         url,
				Transferred: pr.total,
				Size:        resp.ContentLength,
				Speed:       0,
				Progress:    100 * (float64(pr.total) / float64(resp.ContentLength)),
				Downloading: pr.err == nil && pr.total < resp.ContentLength,
				Done:        pr.total == resp.ContentLength,
			})
			if pr.err != nil || pr.total == resp.ContentLength {
				break
			}
		}
	}()

	var updateFile io.Reader = pr
	// extract macos binary from zip
	if strings.HasSuffix(url, ".zip") && runtime.GOOS == "darwin" {
		zipBytes, err := io.ReadAll(pr)
		if err != nil {
			return false, err
		}
		archive, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
		if err != nil {
			return false, err
		}
		file, err := archive.Open("godoos.app/Contents/MacOS/godoos")
		if err != nil {
			return false, err
		}
		defer file.Close()
		updateFile = file
	}

	// apply update
	err = selfupdate.Apply(updateFile, selfupdate.Options{})
	if err != nil {
		if rerr := selfupdate.RollbackError(err); rerr != nil {
			return true, rerr
		}
		return false, err
	}
	// restart app
	a.RestartApp()
	return false, nil
}

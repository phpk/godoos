/*
 * godoos - A lightweight cloud desktop
 * Copyright (C) 2024 godoos.com
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package sys

import (
	"encoding/json"
	"fmt"
	"godo/files"
	"godo/libs"
	"io"
	"log"
	"net/http"
	"os"
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
	Version string                     `json:"version"`
	Url     string                     `json:"url"`
	Name    string                     `json:"name"`
	Desc    string                     `json:"desc"`
	AdList  []map[string][]UpdateAdReq `json:"adlist"`
}
type ServerRes struct {
	Sucess  bool             `json:"sucess"`
	Message string           `json:"message"`
	Data    UpdateVersionReq `json:"data"`
	Time    int64            `json:"time"`
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.reader.Read(p)
	pr.err = err
	pr.total += int64(n)
	return
}
func UpdateAppHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		libs.ErrorMsg(w, "url is empty")
		return
	}
	downloadDir := libs.GetCacheDir()
	filePath := filepath.Join(downloadDir, filepath.Base(url))
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	pr := &ProgressReader{reader: resp.Body}
	// 打开文件
	out, err := os.Create(filePath)
	if err != nil {
		libs.ErrorMsg(w, "open file error")
		return
	}
	defer out.Close()
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
	// 复制数据
	if _, err := io.Copy(out, pr); err != nil {
		libs.ErrorMsg(w, "copy error")
		return
	}
	// 解压缩文件
	unzippedFilePath, err := files.Unzip(filePath, downloadDir)
	if err != nil {
		libs.ErrorMsg(w, "unzip error")
		return
	}
	// 将解压缩后的文件读取为 io.Reader
	unzipFile, err := os.Open(unzippedFilePath)
	if err != nil {
		libs.ErrorMsg(w, "open unzip file error")
		return
	}
	defer unzipFile.Close()
	// apply update
	err = selfupdate.Apply(unzipFile, selfupdate.Options{})
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
func GetUpdateInfo() (UpdateVersionReq, error) {
	var errs UpdateVersionReq
	var updateInfo ServerRes
	info, err := libs.GetSystemInfo()
	if err != nil {
		return errs, fmt.Errorf("update error get info:" + err.Error())
	}
	updateUrl := "https://godoos.com/version?info=" + info
	//log.Printf("updateUrl:%v", updateUrl)
	res, err := http.Get(updateUrl)
	if err != nil {
		return errs, fmt.Errorf("update error get url:" + err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errs, fmt.Errorf("update error get url")
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errs, fmt.Errorf("update error read body:" + err.Error())
	}
	err = json.Unmarshal(body, &updateInfo)
	//log.Printf("updateInfo:%v", updateInfo)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return errs, fmt.Errorf("update error unmarshal:" + err.Error())
	}
	return updateInfo.Data, nil
}
func GetUpdateUrlHandler(w http.ResponseWriter, r *http.Request) {
	updateInfo, err := GetUpdateInfo()
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	json.NewEncoder(w).Encode(updateInfo)
}

package model

import (
	"encoding/json"
	"fmt"
	"godo/libs"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

func Tagshandler(w http.ResponseWriter, r *http.Request) {
	err := LoadConfig()
	if err != nil {
		libs.ErrorMsg(w, "Load config error")
		return
	}
	var reqBodies []ReqBody
	reqBodyMap.Range(func(key, value interface{}) bool {
		rb, ok := value.(ReqBody)
		if ok {
			reqBodies = append(reqBodies, rb)
		}
		return true // 继续遍历
	})
	// 对reqBodies按CreatedAt降序排列
	sort.Slice(reqBodies, func(i, j int) bool {
		return reqBodies[i].CreatedAt.After(reqBodies[j].CreatedAt) // 降序排列
	})
	// 设置响应内容类型为JSON
	w.Header().Set("Content-Type", "application/json")

	// 使用json.NewEncoder将reqBodies编码为JSON并写入响应体
	if err := json.NewEncoder(w).Encode(reqBodies); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func ShowHandler(w http.ResponseWriter, r *http.Request) {
	err := LoadConfig()
	if err != nil {
		libs.ErrorMsg(w, "Load config error")
		return
	}
	model := r.URL.Query().Get("model")
	if model == "" {
		libs.ErrorMsg(w, "Model name is empty")
		return
	}
	//log.Printf("ShowHandler: %s", model)
	var reqBodies ReqBody
	reqBodyMap.Range(func(key, value interface{}) bool {
		rb, ok := value.(ReqBody)
		if ok && rb.Model == model {
			reqBodies = rb
			return false
		}
		return true
	})
	//log.Printf("ShowHandler: %s", reqBodies)
	// 设置响应内容类型为JSON
	w.Header().Set("Content-Type", "application/json")

	// 使用json.NewEncoder将reqBodies编码为JSON并写入响应体
	if err := json.NewEncoder(w).Encode(reqBodies); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type DownserverStucct struct {
	Path string `json:"path"`
}

func DownServerHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	//log.Printf("imagePath: %s", imagePath)
	// 检查路径是否为空或无效
	if filePath == "" {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}
	if !libs.PathExists(filePath) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	// 获取文件信息以获取文件大小
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Printf("Error getting file info: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// 设置响应头，指示浏览器以附件形式下载文件
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(filePath)))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// 读取文件并写入响应体
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	if err != nil {
		log.Printf("Error copying file to response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// MIT License
//
// Copyright (c) 2024 godoos.com
// Email: xpbb@qq.com
// GitHub: github.com/phpk/godoos
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package localchat

import (
	"encoding/json"
	"fmt"
	"godo/libs"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// FileItem 表示文件或文件夹
type FileItem struct {
	Path      string `json:"path"`
	IsDir     bool   `json:"isDir"`
	Filename  string `json:"filename"`
	WritePath string `json:"writePath"`
}
type FileList struct {
	Files []string `json:"fileList"`
}

func HandleGetFiles(w http.ResponseWriter, r *http.Request) {
	//log.Printf("=====Received request: %v", r)
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var fileList FileList
	err := json.NewDecoder(r.Body).Decode(&fileList)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//log.Printf("=====Received file list: %v", fileList)
	defer r.Body.Close()

	baseDir, err := libs.GetOsDir()
	if err != nil {
		log.Printf("Failed to get OS directory: %v", err)
		http.Error(w, "Failed to get OS directory", http.StatusInternalServerError)
		return
	}

	// 用于存储文件列表
	var files []FileItem

	for _, filePath := range fileList.Files {
		fp := filepath.Join(baseDir, filePath)

		fileInfo, err := os.Stat(fp)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to stat file: %v", err), http.StatusInternalServerError)
			return
		}
		//writePath := calculateWritePath(fp, baseDir)
		if fileInfo.IsDir() {
			if err := walkDirectory(fp, &files, filepath.Base(fp)); err != nil {
				http.Error(w, fmt.Sprintf("Failed to serve directory: %v", err), http.StatusInternalServerError)
				return
			}
		} else {

			files = append(files, FileItem{
				Path:      fp,
				IsDir:     false,
				Filename:  filepath.Base(fp),
				WritePath: "",
			})
		}
	}

	// 将文件列表编码为 JSON 并返回
	jsonData, err := json.Marshal(files)
	if err != nil {
		http.Error(w, "Failed to marshal file list", http.StatusInternalServerError)
		return
	}
	log.Printf("Sending file list: %v", string(jsonData))

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
func calculateWritePath(filePath, baseDir string) string {
	relativePath, err := filepath.Rel(baseDir, filePath)
	if err != nil {
		log.Printf("Failed to calculate relative path: %v", err)
		return ""
	}
	return filepath.Dir(relativePath)
}
func walkDirectory(rootPath string, files *[]FileItem, writePath string) error {
	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk directory: %v", err)
		}

		isDir := info.IsDir()
		relativePath, err := filepath.Rel(rootPath, path)
		if err != nil {
			log.Printf("Failed to calculate relative path: %v", err)
			return fmt.Errorf("failed to calculate relative path")
		}
		currentWritePath := filepath.Join(writePath, filepath.Base(relativePath))
		*files = append(*files, FileItem{
			Path:      path,
			IsDir:     isDir,
			Filename:  filepath.Base(path),
			WritePath: currentWritePath,
		})

		return nil
	})
}

func HandleServeFile(w http.ResponseWriter, r *http.Request) {
	// 从 URL 中获取 filePath 参数
	filePath := r.URL.Query().Get("path")

	if filePath == "" {
		http.Error(w, "Missing filePath parameter", http.StatusBadRequest)
		return
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to stat file: %v", err), http.StatusInternalServerError)
		return
	}

	if fileInfo.IsDir() {
		http.Error(w, "Cannot download a directory", http.StatusBadRequest)
		return
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to open file: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 设置响应头
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileInfo.Name()))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// 复制文件内容到响应体
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to copy file content: %v", err), http.StatusInternalServerError)
		return
	}
}

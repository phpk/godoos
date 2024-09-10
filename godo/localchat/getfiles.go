package localchat

import (
	"encoding/json"
	"fmt"
	"godo/libs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type FileList struct {
	Files []string `json:"fileList"`
}

func HandleGetFiles(w http.ResponseWriter, r *http.Request) {
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
	log.Printf("Received file list: %v", fileList)
	defer r.Body.Close()

	baseDir, err := libs.GetOsDir()
	if err != nil {
		log.Printf("Failed to get OS directory: %v", err)
		return
	}

	// 用于存储文件列表
	var files []string

	for _, filePath := range fileList.Files {
		fp := filepath.Join(baseDir, filePath)
		if err := serveDirectory(w, r, fp, &files); err != nil {
			http.Error(w, fmt.Sprintf("Failed to serve directory: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// 将文件列表编码为 JSON 并返回
	jsonData, err := json.Marshal(files)
	if err != nil {
		http.Error(w, "Failed to marshal file list", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func serveDirectory(w http.ResponseWriter, r *http.Request, dirPath string, files *[]string) error {
	filesInDir, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	for _, f := range filesInDir {
		filePath := filepath.Join(dirPath, f.Name())
		if f.IsDir() {
			if err := serveDirectory(w, r, filePath, files); err != nil {
				return err
			}
		} else {
			*files = append(*files, filePath)
		}
	}

	return nil
}

func ServeFile(w http.ResponseWriter, r *http.Request, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %v", err)
	}

	if fileInfo.IsDir() {
		return serveDirectory(w, r, filePath, nil)
	}

	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
	return nil
}

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
	Files []string `json:"files"`
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
	defer r.Body.Close()
	baseDir, err := libs.GetOsDir()
	if err != nil {
		log.Printf("Failed to get OS directory: %v", err)
		return
	}
	for _, filePath := range fileList.Files {
		fp := filepath.Join(baseDir, filePath)
		err := ServeFile(w, r, fp)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to serve file: %v", err), http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintf(w, "Files served successfully")
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
		return serveDirectory(w, r, filePath)
	}

	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
	return nil
}

func serveDirectory(w http.ResponseWriter, r *http.Request, dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	for _, f := range files {
		filePath := filepath.Join(dirPath, f.Name())
		if f.IsDir() {
			err := serveDirectory(w, r, filePath)
			if err != nil {
				return err
			}
		} else {
			err := ServeFile(w, r, filePath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

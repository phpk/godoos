package store

import (
	"godo/files"
	"godo/libs"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// UploadHandler 处理上传的HTTP请求
func UploadHandler(w http.ResponseWriter, r *http.Request) {

	// 解析上传的文件
	err := r.ParseMultipartForm(10000 << 20) // 限制最大上传大小为1000MB
	if err != nil {
		http.Error(w, "上传文件过大"+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("files") // 表单字段名为"files"
	if err != nil {
		http.Error(w, "没有找到文件", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 读取文件内容
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("读取文件内容出错: %v", err)
		http.Error(w, "读取文件内容出错", http.StatusInternalServerError)
		return
	}
	cachePath := libs.GetCacheDir()
	baseName := filepath.Base(header.Filename)
	filenameNoExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))

	savePath := filepath.Join(cachePath, baseName)

	out, err := os.Create(savePath)
	if err != nil {
		log.Printf("创建文件出错: %v", err)
		http.Error(w, "保存文件出错", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// 将文件内容写入到服务器上的文件
	_, err = out.Write(fileBytes)
	if err != nil {
		log.Printf("写入文件出错: %v", err)
		http.Error(w, "写入文件出错", http.StatusInternalServerError)
		return
	}
	runDir := libs.GetRunDir()
	err = files.HandlerFile(savePath, runDir)
	if err != nil {
		log.Printf("Error moving file: %v", err)
	}

	libs.SuccessMsg(w, filenameNoExt, "File already exists and is of correct size")
}

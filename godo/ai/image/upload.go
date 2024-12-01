/*
 * GodoAI - A software focused on localizing AI applications
 * Copyright (C) 2024 https://godoos.com
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
package image

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// UploadHandler 处理单张图片上传的HTTP请求
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// 设置允许的Content-Type
	//r.Header.Set("Content-Type", "multipart/form-data")

	// 解析上传的文件
	err := r.ParseMultipartForm(10000 << 20) // 限制最大上传大小为1000MB
	if err != nil {
		http.Error(w, "上传文件过大"+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("files") // 假设表单字段名为"files"
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
	imagePath, err := GetImageDir()
	if err != nil {
		log.Printf("获取图片目录出错: %v", err)
		http.Error(w, "获取图片目录出错", http.StatusInternalServerError)
		return
	}
	// 生成随机文件名并保留原扩展名
	randomName := generateRandomString(10) + filepath.Ext(header.Filename)
	savePath := filepath.Join(imagePath, randomName)

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
	type ResJson struct {
		Path string `json:"path"`
	}
	resJson := ResJson{Path: savePath}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resJson); err != nil {
		log.Printf("json encode error: %v", err)
	}

	// jsonStr := fmt.Sprintf(`{"code": 0, "message": "上传成功", "data":"%s"}`, savePath)
	// w.Header().Set("Content-Type", "application/json")
	// w.Write([]byte(jsonStr))
}
func ServeImage(w http.ResponseWriter, r *http.Request) {
	// 从 URL 查询参数中获取图片路径
	imagePath := r.URL.Query().Get("path")
	//log.Printf("imagePath: %s", imagePath)
	// 检查图片路径是否为空或无效
	if imagePath == "" {
		http.Error(w, "Invalid image path", http.StatusBadRequest)
		return
	}

	// 确保图片路径是绝对路径
	absImagePath, err := filepath.Abs(imagePath)
	//log.Printf("absImagePath: %s", absImagePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 获取文件的 MIME 类型
	mimeType := mime.TypeByExtension(filepath.Ext(absImagePath))
	if mimeType == "" {
		mimeType = "application/octet-stream" // 如果无法识别，就用默认的二进制流类型
	}

	// 设置响应头的 MIME 类型
	w.Header().Set("Content-Type", mimeType)

	// 打开文件并读取内容
	file, err := os.Open(absImagePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 将文件内容写入响应体
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// DeleteImageHandler 处理删除图片的HTTP请求
func DeleteImageHandler(w http.ResponseWriter, r *http.Request) {
	// 解析请求体中的JSON数据
	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "无法解析请求数据", http.StatusBadRequest)
		return
	}

	imgPath, ok := data["path"]
	if !ok || imgPath == "" {
		http.Error(w, "缺少图片路径数据", http.StatusBadRequest)
		return
	}

	// 检查图片路径是否有效
	if _, err := os.Stat(imgPath); os.IsNotExist(err) {
		http.Error(w, "图片不存在", http.StatusNotFound)
		return
	}

	// 删除图片
	err = os.Remove(imgPath)
	if err != nil {
		log.Printf("删除文件出错: %v", err)
		http.Error(w, "删除文件出错", http.StatusInternalServerError)
		return
	}

	jsonStr := `{"code": 0, "message": "图片删除成功"}`
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonStr))
}

// 生成随机字符串
func generateRandomString(length int) string {
	//rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result strings.Builder
	for i := 0; i < length; i++ {
		result.WriteByte(charset[rand.Intn(len(charset))])
	}
	return result.String()
}

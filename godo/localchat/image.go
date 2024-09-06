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
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func HandlerSendImg(w http.ResponseWriter, r *http.Request) {
	var msg UdpMessage
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	hostname, err := os.Hostname()
	if err != nil {
		libs.ErrorMsg(w, "HandleMessage error")
		return
	}
	msg.Hostname = hostname
	msg.Time = time.Now()
	msg.Type = "image"
	toIp := msg.IP
	basePath, err := libs.GetOsDir()
	if err != nil {
		log.Printf("GetOsDir error: %v", err)
		return
	}
	log.Printf("send image to %v", msg.Message)
	log.Printf("Type of msg.Message: %T", msg.Message)
	paths, ok := msg.Message.([]interface{})
	log.Printf("paths: %v", paths)
	if !ok {
		libs.ErrorMsg(w, "HandleMessage message error")
		return
	}

	for _, v := range paths {
		p, ok := v.(string)
		if !ok {
			continue
		}
		filePath := filepath.Join(basePath, p)
		// 处理多张图片
		if fileInfo, err := os.Stat(filePath); err == nil {
			if !fileInfo.IsDir() {
				sendImage(filePath, toIp, msg)
			}
		} else {
			continue
		}
	}
	libs.SuccessMsg(w, nil, "图片发送成功")
}

func sendImage(filePath string, toIp string, message UdpMessage) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Failed to get file info: %v", err)
	}
	fileSize := fileInfo.Size()

	// 读取整个文件
	data := make([]byte, fileSize)
	_, err = file.Read(data)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// 创建文件块
	chunk := FileChunk{
		Data:      data,
		Checksum:  calculateChecksum(data),
		Timestamp: time.Now(),
		Filename:  filepath.Base(file.Name()),
	}
	message.Message = chunk
	// 将文件块转换为 JSON 格式
	data, err = json.Marshal(message)
	if err != nil {
		log.Fatalf("Failed to marshal chunk: %v", err)
	}

	// 发送文件块
	port := "56780"
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%s", toIp, port))
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		log.Fatalf("Failed to dial UDP address: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		log.Printf("Failed to write data: %v", err)
	}

	fmt.Printf("发送图片 %s 到 %s 成功\n", filePath, toIp)
}
func ReceiveImg(msg UdpMessage) (string, error) {
	chunk := msg.Message.(FileChunk)

	// 验证校验和
	calculatedChecksum := calculateChecksum(chunk.Data)
	if calculatedChecksum != chunk.Checksum {
		fmt.Printf("Checksum mismatch for image from %s\n", msg.IP)
		return "", fmt.Errorf("checksum mismatch")
	}

	baseDir, err := libs.GetOsDir()
	if err != nil {
		log.Printf("Failed to get OS directory: %v", err)
		return "", err
	}

	// 创建接收文件的目录
	resPath := filepath.Join("C", "Users", "Reciv", time.Now().Format("2006-01-02"))
	receiveDir := filepath.Join(baseDir, resPath)
	if !libs.PathExists(receiveDir) {
		err := os.MkdirAll(receiveDir, 0755)
		if err != nil {
			log.Printf("Failed to create receive directory: %v", err)
			return "", err
		}
	}

	// 确定文件路径
	filePath := filepath.Join(receiveDir, chunk.Filename)

	// 如果文件不存在，则创建新文件
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			log.Printf("Failed to create file: %v", err)
			return "", err
		}
		defer file.Close()
	}

	// 打开或追加到现有文件
	file, err := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		return "", err
	}
	defer file.Close()

	// 写入数据
	_, err = file.Write(chunk.Data)
	if err != nil {
		log.Printf("Failed to write data to file: %v", err)
		return "", err
	}

	fmt.Printf("接收到图片 %s 从 %s 成功\n", filePath, msg.IP)
	resFilePath := filepath.Join(resPath, chunk.Filename)
	return resFilePath, nil
}
func HandleViewImg(w http.ResponseWriter, r *http.Request) {
	img := r.URL.Query().Get("img")
	if img == "" {
		libs.ErrorMsg(w, "img is empty")
		return
	}
	decodedImgParam, err := url.QueryUnescape(img)
	if err != nil {
		log.Fatalf("Error unescaping image parameter: %v", err)
	}
	basePath, err := libs.GetOsDir()
	if err != nil {
		log.Printf("GetOsDir error: %v", err)
		return
	}
	filePath := filepath.Join(basePath, decodedImgParam)
	log.Printf("filePath: %s", filePath)
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		libs.ErrorMsg(w, "file not found")
		return
	}

	// 设置正确的 MIME 类型
	mimeType, err := GetMimeType(filePath)
	if err != nil {
		http.Error(w, "Failed to determine MIME type", http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", mimeType)

	// 读取文件并写入响应体
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Failed to write file content", http.StatusInternalServerError)
		return
	}
}
func GetMimeType(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// 读取前 512 字节以确定 MIME 类型
	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err != nil {
		return "", err
	}

	mimeType := http.DetectContentType(buffer)
	return mimeType, nil
}

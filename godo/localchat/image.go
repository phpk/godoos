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

package localchat

import (
	"encoding/json"
	"fmt"
	"godo/libs"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	basePath, err := libs.GetOsDir()
	if err != nil {
		log.Printf("GetOsDir error: %v", err)
		return
	}
	// log.Printf("send image to %v", msg.Message)
	// log.Printf("Type of msg.Message: %T", msg.Message)
	paths, ok := msg.Message.([]interface{})
	//log.Printf("paths: %v", paths)
	if !ok {
		libs.ErrorMsg(w, "HandleMessage message error")
		return
	}
	sendPath := []string{}
	for _, v := range paths {
		p, ok := v.(string)
		if !ok {
			continue
		}
		filePath := filepath.Join(basePath, p)
		// 处理多张图片
		if fileInfo, err := os.Stat(filePath); err == nil {
			if !fileInfo.IsDir() {
				//handleFile(filePath, toIp, msg)
				sendPath = append(sendPath, p)
			}
		} else {
			continue
		}
	}
	msg.Message = sendPath
	SendToIP(msg)
	libs.SuccessMsg(w, nil, "图片发送成功")
}

func ReceiveImg(msg UdpMessage) ([]string, error) {
	res := []string{}
	baseDir, err := libs.GetOsDir()
	if err != nil {
		log.Printf("Failed to get OS directory: %v", err)
		return res, err
	}

	// 创建接收文件的目录
	resPath := filepath.Join("C", "Users", "Reciv", time.Now().Format("2006-01-02"))
	receiveDir := filepath.Join(baseDir, resPath)
	if !libs.PathExists(receiveDir) {
		err := os.MkdirAll(receiveDir, 0755)
		if err != nil {
			log.Printf("Failed to create receive directory: %v", err)
			return res, err
		}
	}
	paths, ok := msg.Message.([]interface{})
	//log.Printf("paths: %v", paths)
	if !ok {
		return res, fmt.Errorf("HandleMessage message error")
	}
	var savedPaths []string
	for _, v := range paths {
		p, ok := v.(string)
		if !ok {
			continue
		}
		imgUrl := fmt.Sprintf("http://%s:56780/localchat/viewimage?img=%s", msg.IP, url.QueryEscape(p))
		resp, err := http.Get(imgUrl)
		if err != nil {
			log.Printf("Failed to download image from URL %s: %v", imgUrl, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Failed to download image from URL %s: %v", imgUrl, resp.Status)
			continue
		}
		// 生成随机文件名
		// 生成随机文件名并保留扩展名
		fileName, err := generateRandomFileNameWithExtension(p)
		if err != nil {
			log.Printf("Failed to generate random file name: %v", err)
			continue
		}
		filePath := filepath.Join(receiveDir, fileName)
		// 保存图片
		err = saveImage(resp.Body, filePath)
		if err != nil {
			log.Printf("Failed to save image to %s: %v", filePath, err)
			continue
		}
		receviePath := filepath.Join(resPath, fileName)
		savedPaths = append(savedPaths, receviePath)

	}
	if len(savedPaths) > 0 {
		return savedPaths, nil
	}

	return res, fmt.Errorf("no images were saved")

}

// 生成随机文件名并保留扩展名
func generateRandomFileNameWithExtension(originalFileName string) (string, error) {
	_, fileExt := filepath.Split(originalFileName)
	fileExt = strings.TrimPrefix(fileExt, ".")
	if fileExt == "" {
		fileExt = "png"
	}

	randomFileName := fmt.Sprintf("%s.%s", strconv.FormatInt(time.Now().UnixNano(), 10), fileExt)
	return randomFileName, nil
}

// 保存图片到本地文件
func saveImage(reader io.Reader, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("failed to write image data to file: %v", err)
	}

	return nil
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

/*
 * GodoOS - A lightweight cloud desktop
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

package localchat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"godo/libs"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func downloadFiles(msg UdpMessage) (string, error) {
	postUrl := fmt.Sprintf("http://%s:56780/localchat/getfiles", msg.IP)
	postData, err := json.Marshal(msg.Message)
	if err != nil {
		return "", fmt.Errorf("failed to marshal post data: %v", err)
	}
	log.Printf("Sending POST request to %s with data: %s", postUrl, string(postData))
	resp, err := http.Post(postUrl, "application/json", bytes.NewBuffer(postData))
	if err != nil {
		return "", fmt.Errorf("failed to make POST request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("server returned status code: %v, body: %s", resp.StatusCode, body)
	}
	path, err := handleResponse(resp.Body, msg.IP)
	// 处理响应中的文件
	if err != nil {
		log.Fatalf("Failed to handle response: %v", err)
	}

	fmt.Println("Files downloaded successfully")
	return path, nil
}

func handleResponse(reader io.Reader, ip string) (string, error) {
	// 接收文件的目录
	baseDir, err := libs.GetOsDir()
	if err != nil {
		log.Printf("Failed to get OS directory: %v", err)
		return "", fmt.Errorf("failed to get OS directory")
	}

	resPath := filepath.Join("C", "Users", "Reciv", time.Now().Format("2006-01-02"))
	receiveDir := filepath.Join(baseDir, resPath)
	if !libs.PathExists(receiveDir) {
		err := os.MkdirAll(receiveDir, 0755)
		if err != nil {
			log.Printf("Failed to create receive directory: %v", err)
			return "", fmt.Errorf("failed to create receive directory")
		}
	}
	timestamp := time.Now().Format("15-04-05")
	revPath := filepath.Join(receiveDir, timestamp)
	if !libs.PathExists(revPath) {
		err := os.MkdirAll(revPath, 0755)
		if err != nil {
			log.Printf("Failed to create receive directory: %v", err)
			return "", fmt.Errorf("failed to create receive directory")
		}
	}
	body, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	//log.Printf("Received file list: %v", string(body))
	// 解析文件列表
	var fileList []FileItem
	if err := json.Unmarshal(body, &fileList); err != nil {
		return "", fmt.Errorf("failed to unmarshal file list: %v", err)
	}

	//log.Printf("Received file list: %v", fileList)

	for _, file := range fileList {
		if runtime.GOOS != "windows" && strings.Contains(file.WritePath, "\\") {
			file.WritePath = strings.ReplaceAll(file.WritePath, "\\", "/")
		}
		checkpath := filepath.Join(revPath, file.WritePath)

		if !libs.PathExists(checkpath) {
			os.MkdirAll(checkpath, 0755)
		}

		if !file.IsDir {
			go downloadFile(file.Path, checkpath, ip)
		}
	}

	return filepath.Join(resPath, timestamp), nil
}

// downloadFile 下载单个文件
func downloadFile(filePath string, checkpath string, ip string) error {
	url := fmt.Sprintf("http://%s:56780/localchat/servefile?path=%s", ip, filePath)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %v, body: %s", resp.StatusCode, body)
	}

	// 保存文件
	fileName := filepath.Join(checkpath, filepath.Base(filePath))

	out, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	fmt.Printf("File downloaded: %s\n", fileName)
	return nil
}

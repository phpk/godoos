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
	"bytes"
	"encoding/json"
	"fmt"
	"godo/libs"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func downloadFiles(msg UdpMessage) error {
	postUrl := fmt.Sprintf("http://%s:56780/localchat/getfiles", msg.IP)
	postData, err := json.Marshal(msg.Message)
	if err != nil {
		return fmt.Errorf("failed to marshal post data: %v", err)
	}
	log.Printf("Sending POST request to %s with data: %s", postUrl, string(postData))
	resp, err := http.Post(postUrl, "application/json", bytes.NewBuffer(postData))
	if err != nil {
		return fmt.Errorf("failed to make POST request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned status code: %v, body: %s", resp.StatusCode, body)
	}

	// 接收文件的目录
	baseDir, err := libs.GetOsDir()
	if err != nil {
		log.Printf("Failed to get OS directory: %v", err)
		return fmt.Errorf("failed to get OS directory")
	}

	resPath := filepath.Join("C", "Users", "Reciv", time.Now().Format("2006-01-02"))
	receiveDir := filepath.Join(baseDir, resPath)
	if !libs.PathExists(receiveDir) {
		err := os.MkdirAll(receiveDir, 0755)
		if err != nil {
			log.Printf("Failed to create receive directory: %v", err)
			return fmt.Errorf("failed to create receive directory")
		}
	}

	// 处理响应中的文件
	if err := handleResponse(resp.Body, receiveDir, msg.IP); err != nil {
		log.Fatalf("Failed to handle response: %v", err)
	}

	fmt.Println("Files downloaded successfully")
	return nil
}

func handleResponse(reader io.Reader, saveDir string, ip string) error {
	body, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}
	log.Printf("Received file list: %v", string(body))
	// 解析文件列表
	var fileList []FileItem
	if err := json.Unmarshal(body, &fileList); err != nil {
		return fmt.Errorf("failed to unmarshal file list: %v", err)
	}
	log.Printf("Received file list: %v", fileList)

	for _, file := range fileList {
		checkpath := filepath.Join(saveDir, file.WritePath)
		if err := os.MkdirAll(checkpath, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
		if !file.IsDir {
			if err := downloadFile(file.Path, checkpath, ip); err != nil {
				return fmt.Errorf("failed to download file: %v", err)
			}
		}
	}

	return nil
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

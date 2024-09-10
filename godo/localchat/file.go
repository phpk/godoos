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

func HandlerApplySendFile(w http.ResponseWriter, r *http.Request) {
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
	msg.Type = "fileSending"
	SendToIP(msg)
	libs.SuccessMsg(w, nil, "请求文件发送成功")
}
func HandlerCannelFile(w http.ResponseWriter, r *http.Request) {
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
	msg.Type = "fileCannel"
	SendToIP(msg)
	libs.SuccessMsg(w, nil, "请求文件发送成功")
}
func HandlerAccessFile(w http.ResponseWriter, r *http.Request) {
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
	msg.Type = "fileAccessed"
	SendToIP(msg)
	err = downloadFiles(msg)
	if err != nil {
		libs.ErrorMsg(w, "HandleMessage error")
		return
	}
	libs.SuccessMsg(w, msg.Message, "接收文件中")
}

func downloadFiles(msg UdpMessage) error {
	postUrl := fmt.Sprintf("http://%s:56780/localchat/getfiles", msg.IP)
	postData, err := json.Marshal(msg.Message)
	if err != nil {
		return fmt.Errorf("failed to marshal post data: %v", err)
	}

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

	err = saveFiles(resp.Body, receiveDir)
	if err != nil {
		return fmt.Errorf("failed to save files: %v", err)
	}

	fmt.Println("Files downloaded successfully")
	return nil
}

func saveFiles(reader io.Reader, saveDir string) error {
	body, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	var fileList FileList
	err = json.Unmarshal(body, &fileList)
	if err != nil {
		return fmt.Errorf("failed to unmarshal file list: %v", err)
	}

	for _, filePath := range fileList.Files {
		err := saveFileOrFolder(filePath, saveDir)
		if err != nil {
			return fmt.Errorf("failed to serve file or folder: %v", err)
		}
	}

	return nil
}

func saveFileOrFolder(filePath string, saveDir string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file or folder: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file or folder: %v", err)
	}

	if fileInfo.IsDir() {
		return saveFolder(fileInfo.Name(), filePath, saveDir)
	}

	return saveFile(filePath, saveDir)
}

func saveFolder(folderName string, folderPath string, saveDir string) error {
	localFolderPath := filepath.Join(saveDir, folderName)
	err := os.MkdirAll(localFolderPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	for _, file := range files {
		subPath := filepath.Join(folderPath, file.Name())
		subSavePath := filepath.Join(localFolderPath, file.Name())

		if file.IsDir() {
			err := saveFolder(file.Name(), subPath, localFolderPath)
			if err != nil {
				return fmt.Errorf("failed to save folder: %v", err)
			}
		} else {
			err := saveFile(subPath, subSavePath)
			if err != nil {
				return fmt.Errorf("failed to save file: %v", err)
			}
		}
	}

	return nil
}

func saveFile(filePath string, saveDir string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %v", err)
	}

	localFilePath := filepath.Join(saveDir, fileInfo.Name())
	err = os.MkdirAll(filepath.Dir(localFilePath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	out, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return fmt.Errorf("failed to copy file: %v", err)
	}

	return nil
}

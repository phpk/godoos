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
	"os"
	"path/filepath"
	"time"
)

const (
	fileSize = 1024 // 每个数据包的大小
)

type FileChunk struct {
	ChunkIndex int       `json:"chunk_index"`
	Data       []byte    `json:"data"`
	Checksum   uint32    `json:"checksum"`
	Timestamp  time.Time `json:"timestamp"`
	Filename   string    `json:"filename"`
}

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
	msg.Message = ""
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
	msg.Message = ""
	SendToIP(msg)
	libs.SuccessMsg(w, nil, "请求文件发送成功")
}
func HandlerSendFile(msg UdpMessage) {
	toIp := msg.IP
	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("HandleMessage error: %v", err)
		return
	}
	msg.Hostname = hostname
	msg.Time = time.Now()
	msg.Type = "file"
	basePath, err := libs.GetOsDir()
	if err != nil {
		log.Printf("GetOsDir error: %v", err)
		return
	}
	paths, ok := msg.Message.([]string)
	if !ok {
		log.Printf("invalid message type")
		return
	}
	for _, p := range paths {
		filePath := filepath.Join(basePath, p)
		// 处理单个文件或整个文件夹
		if fileInfo, err := os.Stat(filePath); err == nil {
			if fileInfo.IsDir() {
				handleDirectory(filePath, toIp, msg)
			} else {
				handleFile(filePath, toIp, msg)
			}
		} else {
			continue
		}
	}
	msg.Type = "fileSended"
	msg.Message = ""
	msg.Time = time.Now()
	SendToIP(msg)
}

func handleFile(filePath string, toIp string, message UdpMessage) {
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

	// 计算需要发送的数据包数量
	numChunks := (fileSize + fileSize - 1) / fileSize

	// 发送文件
	SendFile(file, int(numChunks), toIp, message)
}

func handleDirectory(dirPath string, toIp string, message UdpMessage) {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			handleFile(path, toIp, message)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to walk directory: %v", err)
	}
}
func SendFile(file *os.File, numChunks int, toIp string, message UdpMessage) {
	// 逐块读取文件并发送
	for i := 0; i < numChunks; i++ {
		var chunkData [fileSize]byte
		n, err := file.Read(chunkData[:])
		if err != nil && err != io.EOF {
			log.Fatalf("Failed to read file chunk: %v", err)
		}

		// 创建文件块
		chunk := FileChunk{
			ChunkIndex: i,
			Data:       chunkData[:n],
			Checksum:   calculateChecksum(chunkData[:n]),
			Timestamp:  time.Now(),
			Filename:   filepath.Base(file.Name()),
		}
		message.Message = chunk
		// 将文件块转换为 JSON 格式
		data, err := json.Marshal(message)
		if err != nil {
			log.Fatalf("Failed to marshal chunk: %v", err)
		}

		// 发送文件块
		//addr, err := net.ResolveUDPAddr("udp4", toIp+":56780")
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

		fmt.Printf("发送文件块 %d 到 %s 成功\n", i, toIp)
	}
}
func ReceiveFile(msg UdpMessage) {
	chunk := msg.Message.(FileChunk)

	// 验证校验和
	calculatedChecksum := calculateChecksum(chunk.Data)
	if calculatedChecksum != chunk.Checksum {
		fmt.Printf("Checksum mismatch for chunk %d from %s\n", chunk.ChunkIndex, msg.IP)
		return
	}

	baseDir, err := libs.GetOsDir()
	if err != nil {
		log.Printf("Failed to get OS directory: %v", err)
		return
	}

	// 创建接收文件的目录
	receiveDir := filepath.Join(baseDir, "C", "Users", "Reciv", time.Now().Format("2006-01-02"))
	if !libs.PathExists(receiveDir) {
		err := os.MkdirAll(receiveDir, 0755)
		if err != nil {
			log.Printf("Failed to create receive directory: %v", err)
			return
		}
	}

	// 确定文件路径
	filePath := filepath.Join(receiveDir, chunk.Filename)

	// 如果文件不存在，则创建新文件
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			log.Printf("Failed to create file: %v", err)
			return
		}
		defer file.Close()
	}

	// 打开或追加到现有文件
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		return
	}
	defer file.Close()

	// 写入数据
	_, err = file.Write(chunk.Data)
	if err != nil {
		log.Printf("Failed to write data to file: %v", err)
		return
	}

	fmt.Printf("接收到文件块 %d 从 %s 成功\n", chunk.ChunkIndex, msg.IP)
}
func calculateChecksum(data []byte) uint32 {
	checksum := uint32(0)
	for _, b := range data {
		checksum += uint32(b)
	}
	return checksum
}

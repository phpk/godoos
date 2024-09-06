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
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"godo/libs"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	fileSize = 512 // 每个数据包的大小
)

type FileChunk struct {
	ChunkIndex int       `json:"chunk_index"`
	Data       string    `json:"data"`
	Checksum   uint32    `json:"checksum"`
	Timestamp  time.Time `json:"timestamp"`
	Filename   string    `json:"filename"`
	Filesize   int64     `json:"filesize"`
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
	SendFile(file, int(numChunks), toIp, fileSize, message)
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
func SendFile(file *os.File, numChunks int, toIp string, fSize int64, message UdpMessage) {
	var wg sync.WaitGroup

	// 逐块读取文件并发送
	for i := 0; i < numChunks; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			var chunkData [fileSize]byte
			n, err := file.Read(chunkData[:])
			if err != nil && err != io.EOF {
				log.Fatalf("Failed to read file chunk: %v", err)
			}
			encodedData := base64.StdEncoding.EncodeToString(chunkData[:n])

			// 创建文件块
			chunk := FileChunk{
				ChunkIndex: index,
				Data:       encodedData,
				Checksum:   calculateChecksum(chunkData[:n]),
				Timestamp:  time.Now(),
				Filename:   filepath.Base(file.Name()),
				Filesize:   fSize,
			}

			chunkJson, err := json.Marshal(chunk)
			if err != nil {
				log.Fatalf("Failed to marshal chunk: %v", err)
			}

			// 确保每个数据包的大小不超过限制
			maxPacketSize := 65000
			if len(chunkJson) > maxPacketSize {
				// 分割数据包
				chunks := splitChunkJson(chunkJson, maxPacketSize)
				for _, subChunkJson := range chunks {
					message.Message = base64.StdEncoding.EncodeToString(subChunkJson)
					sendData(message, toIp)
				}
			} else {
				message.Message = base64.StdEncoding.EncodeToString(chunkJson)
				sendData(message, toIp)
			}

			fmt.Printf("发送文件块 %d 到 %s 成功\n", index, toIp)
		}(i)
	}

	wg.Wait()
}

func splitChunkJson(jsonData []byte, maxSize int) [][]byte {
	var chunks [][]byte
	for start := 0; start < len(jsonData); start += maxSize {
		end := start + maxSize
		if end > len(jsonData) {
			end = len(jsonData)
		}
		chunks = append(chunks, jsonData[start:end])
	}
	return chunks
}

func sendData(message UdpMessage, toIp string) {
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

	data, _ := json.Marshal(message)
	_, err = conn.Write(data)
	if err != nil {
		log.Printf("Failed to write data: %v", err)
	}
}
func ReceiveFile(msg UdpMessage) (string, error) {
	chunkStr, ok := msg.Message.(string)
	if !ok {
		return "", errors.New("invalid message type")
	}

	// Base64解码
	chunkJson, err := base64.StdEncoding.DecodeString(chunkStr)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 message: %v", err)
	}

	// 从 map 中提取 FileChunk 字段
	var chunk FileChunk
	if err := json.Unmarshal(chunkJson, &chunk); err != nil {
		return "", fmt.Errorf("failed to unmarshal FileChunk: %v", err)
	}

	// 验证校验和
	chunkData, err := base64.StdEncoding.DecodeString(chunk.Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 chunk data: %v", err)
	}
	calculatedChecksum := calculateChecksum(chunkData)
	if calculatedChecksum != chunk.Checksum {
		fmt.Printf("Checksum mismatch for chunk %d from %s\n", chunk.ChunkIndex, msg.IP)
		return "", fmt.Errorf("checksum mismatch")
	}

	// 创建接收文件的目录
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

	// 确定文件路径
	filePath := filepath.Join(receiveDir, chunk.Filename)

	// 如果文件不存在，则创建新文件
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			log.Printf("Failed to create file: %v", err)
			return "", fmt.Errorf("failed to create file")
		}
		defer file.Close()
	}

	// 打开或追加到现有文件
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		return "", fmt.Errorf("failed to open file")
	}
	defer file.Close()

	// 写入数据
	n, err := file.Write(chunkData)
	if err != nil {
		log.Printf("Failed to write data to file: %v", err)
		return "", fmt.Errorf("failed to write data to file")
	}
	if n != len(chunkData) {
		log.Printf("Incomplete write: wrote %d bytes, expected %d bytes", n, len(chunkData))
		return "", fmt.Errorf("incomplete write")
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Printf("Failed to stat file: %v", err)
		return "", fmt.Errorf("failed to stat file")
	}
	if fileInfo.Size() == chunk.Filesize {
		fmt.Println("文件接收完成且大小一致")
		return filePath, nil
	} else {
		fmt.Println("文件大小不一致")
		return "", fmt.Errorf("file size mismatch")
	}
}
func calculateChecksum(data []byte) uint32 {
	checksum := uint32(0)
	for _, b := range data {
		checksum += uint32(b)
	}
	return checksum
}

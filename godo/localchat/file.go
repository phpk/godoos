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
	"godo/libs"
	"net/http"
	"os"
	"time"
)

const (
	GlobalfileSize = 768 // 每个数据包的大小
)

type FileChunk struct {
	ChunkIndex int       `json:"chunk_index"`
	Data       []byte    `json:"data"`
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

	return nil
}

// func HandlerSendFile(msg UdpMessage) {
// 	toIp := msg.IP
// 	hostname, err := os.Hostname()
// 	if err != nil {
// 		log.Printf("HandleMessage error: %v", err)
// 		return
// 	}
// 	msg.Hostname = hostname
// 	msg.Time = time.Now()
// 	msg.Type = "file"
// 	basePath, err := libs.GetOsDir()
// 	if err != nil {
// 		log.Printf("GetOsDir error: %v", err)
// 		return
// 	}
// 	paths, ok := msg.Message.([]string)
// 	if !ok {
// 		log.Printf("invalid message type")
// 		return
// 	}
// 	for _, p := range paths {
// 		filePath := filepath.Join(basePath, p)
// 		// 处理单个文件或整个文件夹
// 		if fileInfo, err := os.Stat(filePath); err == nil {
// 			if fileInfo.IsDir() {
// 				handleDirectory(filePath, toIp, msg)
// 			} else {
// 				handleFile(filePath, toIp, msg)
// 			}
// 		} else {
// 			continue
// 		}
// 	}
// 	msg.Type = "fileSended"
// 	msg.Message = ""
// 	msg.Time = time.Now()
// 	SendToIP(msg)
// }

// func handleFile(filePath string, toIp string, message UdpMessage) {
// 	// 打开文件
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		log.Fatalf("Failed to open file: %v", err)
// 	}
// 	defer file.Close()

// 	// 获取文件大小
// 	fileInfo, err := file.Stat()
// 	if err != nil {
// 		log.Fatalf("Failed to get file info: %v", err)
// 	}
// 	fileSize := fileInfo.Size()

// 	// 计算需要发送的数据包数量
// 	numChunks := fileSize / GlobalfileSize
// 	if fileSize%GlobalfileSize != 0 {
// 		numChunks++
// 	}
// 	// 发送文件
// 	SendFile(file, int(numChunks), toIp, fileSize, message)
// }

// func handleDirectory(dirPath string, toIp string, message UdpMessage) {
// 	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if !info.IsDir() {
// 			handleFile(path, toIp, message)
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		log.Fatalf("Failed to walk directory: %v", err)
// 	}
// }
// func SendFile(file *os.File, numChunks int, toIp string, fSize int64, message UdpMessage) {
// 	var wg sync.WaitGroup

// 	// 逐块读取文件并发送
// 	for i := 0; i < numChunks; i++ {
// 		wg.Add(1)
// 		go func(index int) {
// 			defer wg.Done()
// 			idStr := message.Type + "@" + message.Hostname + "@" + filepath.Base(file.Name()) + "@" + fmt.Sprintf("%d", fSize) + "@"
// 			chunkData := make([]byte, GlobalfileSize)
// 			n, err := file.Read(chunkData[:])
// 			if err != nil && err != io.EOF {
// 				log.Fatalf("Failed to read file chunk: %v", err)
// 			}
// 			sendBinaryData(idStr, chunkData[:n], toIp)
// 			fmt.Printf("发送文件块 %d 到 %s 成功\n", index, toIp)
// 		}(i)
// 	}

// 	wg.Wait()
// }
// func sendBinaryData(idStr string, data []byte, toIp string) {
// 	port := "56780"
// 	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%s", toIp, port))
// 	if err != nil {
// 		log.Fatalf("Failed to resolve UDP address: %v", err)
// 	}

// 	conn, err := net.DialUDP("udp4", nil, addr)
// 	if err != nil {
// 		log.Fatalf("Failed to dial UDP address: %v", err)
// 	}
// 	defer conn.Close()

// 	// 添加长度前缀
// 	identifier := []byte(idStr)
// 	lengthPrefix := make([]byte, 4)
// 	binary.BigEndian.PutUint32(lengthPrefix, uint32(len(data)))

// 	// 发送长度前缀和数据
// 	_, err = conn.Write(append(append(identifier, lengthPrefix...), data...))
// 	if err != nil {
// 		log.Printf("Failed to write data: %v", err)
// 	}
// }

// var fileLocks = make(map[string]*sync.Mutex)

// func ReceiveFiles(parts []string) (string, error) {
// 	fileName := parts[2]
// 	fileSize, err := strconv.ParseInt(parts[3], 10, 64)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to parse file size: %v", err)
// 	}

// 	// 创建接收文件的目录
// 	baseDir, err := libs.GetOsDir()
// 	if err != nil {
// 		log.Printf("Failed to get OS directory: %v", err)
// 		return "", fmt.Errorf("failed to get OS directory")
// 	}

// 	resPath := filepath.Join("C", "Users", "Reciv", time.Now().Format("2006-01-02"))
// 	receiveDir := filepath.Join(baseDir, resPath)
// 	if !libs.PathExists(receiveDir) {
// 		err := os.MkdirAll(receiveDir, 0755)
// 		if err != nil {
// 			log.Printf("Failed to create receive directory: %v", err)
// 			return "", fmt.Errorf("failed to create receive directory")
// 		}
// 	}

// 	// 确定文件路径
// 	filePath := filepath.Join(receiveDir, fileName)

// 	// 如果文件不存在，则创建新文件
// 	if _, err := os.Stat(filePath); os.IsNotExist(err) {
// 		file, err := os.Create(filePath)
// 		if err != nil {
// 			log.Printf("Failed to create file: %v", err)
// 			return "", fmt.Errorf("failed to create file")
// 		}
// 		defer file.Close()
// 	}

// 	// 锁定文件
// 	lock, ok := fileLocks[fileName]
// 	if !ok {
// 		lock = &sync.Mutex{}
// 		fileLocks[fileName] = lock
// 	}
// 	lock.Lock()
// 	defer lock.Unlock()

// 	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		log.Printf("Failed to open file: %v", err)
// 		return "", fmt.Errorf("failed to open file")
// 	}
// 	defer file.Close()

// 	// 提取实际数据
// 	data := strings.Join(parts[4:], "")
// 	if len(data) < 1 {
// 		return "", fmt.Errorf("empty data")
// 	}

// 	n, err := file.Write([]byte(data))
// 	if err != nil {
// 		log.Printf("Failed to write data to file: %v", err)
// 		return "", fmt.Errorf("failed to write data to file")
// 	}
// 	if n != len(data) {
// 		log.Printf("Incomplete write: wrote %d bytes, expected %d bytes", n, len(data))
// 		return "", fmt.Errorf("incomplete write")
// 	}

// 	fileInfo, err := os.Stat(filePath)
// 	if err != nil {
// 		log.Printf("Failed to stat file: %v", err)
// 		return "", fmt.Errorf("failed to stat file")
// 	}

// 	if fileInfo.Size() == fileSize {
// 		fmt.Println("文件接收完成且大小一致")
// 		return filePath, nil
// 	} else {
// 		fmt.Printf("文件大小不一致,发送大小为%d,接收大小为%d\n", fileSize, fileInfo.Size())
// 		return "", fmt.Errorf("file size mismatch")
// 	}
// }

package localchat

import (
	"sync"
	"time"
)

//	type SenderInfo struct {
//		SenderIP string `json:"sender_ip"`
//		Username string `json:"username"`
//	}
type Message struct {
	Type       string       `json:"type"`       // 消息类型，如"text"、"image"等
	Content    string       `json:"content"`    // 消息内容
	SenderInfo UserInfo     `json:"senderInfo"` // 发送者的IP地址
	FileInfo   FilePartInfo `json:"fileInfo"`
	FileList   []UploadInfo `json:"fileList"`
}
type UserInfo struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
}
type UdpMessage struct {
	Type     string `json:"type"`
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	Message  string `json:"message"`
}
type UserList struct {
	Type    string       `json:"type"`
	Content []UdpMessage `json:"content"`
}
type UploadInfo struct {
	Name      string    `json:"name"`
	SavePath  string    `json:"save_path"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// 文件分片信息
type FilePartInfo struct {
	FileName   string `json:"filename"`
	PartNumber int    `json:"part_number"`
	TotalParts int    `json:"total_parts"`
}

// 分片上传状态跟踪
type UploadStatus struct {
	sync.Mutex
	Status map[string]int // key: fileName, value: 已上传分片数
}

var (
	messageChan = make(chan Message, 100) // 缓存大小根据实际情况设定
)
var uploadStatus = UploadStatus{Status: make(map[string]int)}
var broadcartTime = 3 * time.Second

// var broadcastAddr = "224.0.0.1:1679" // 多播地址
var broadcastAddr = "224.0.0.251:1234"

// var broadcastAddr = "255.255.255.255:1769"    // 广播地址
var OnlineUsers = make(map[string]UdpMessage) // 全局map，key为IP，value为主机名

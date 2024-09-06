package localchat

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

type UdpMessage struct {
	Hostname string    `json:"hostname"`
	Type     string    `json:"type"`
	Time     time.Time `json:"time"`
	IP       string    `json:"ip"`
	Message  any       `json:"message"`
}
type Messages struct {
	Messages []UdpMessage `json:"messages"`
}
type UserMessage struct {
	Messages map[string]*Messages  `json:"messages"`
	Onlines  map[string]UserStatus `json:"onlines"`
}

var OnlineUsers = make(map[string]UserStatus)

var UserMessages = make(map[string]*Messages)

func init() {
	go UdpServer()
	go CheckOnlines()
}
func CheckOnlines() {
	CheckOnline()
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		// 检查客户端是否已断开连接
		CheckOnline()
	}
}

// UDP 服务器端逻辑
func UdpServer() {
	// 监听 UDP 端口
	listener, err := net.ListenPacket("udp", ":56780")
	if err != nil {
		log.Fatalf("error setting up listener: %v", err)
	}
	defer listener.Close()

	log.Println("UDP server started on :56780")

	// 监听 UDP 请求
	for {
		buffer := make([]byte, 1024)

		n, remoteAddr, err := listener.ReadFrom(buffer)
		if err != nil {
			log.Printf("error reading from UDP: %v", err)
			continue
		}

		log.Printf("Received UDP packet from %v: %s", remoteAddr, buffer[:n])
		// 解析 UDP 数据
		var udpMsg UdpMessage
		err = json.Unmarshal(buffer[:n], &udpMsg)
		if err != nil {
			log.Printf("error unmarshalling UDP message: %v", err)
			continue
		}
		// 从 remoteAddr 获取 IP 地址
		if udpAddr, ok := remoteAddr.(*net.UDPAddr); ok {
			ip := udpAddr.IP.String()
			udpMsg.IP = ip

			if udpMsg.Type == "heartbeat" {
				UpdateUserStatus(udpMsg.IP, udpMsg.Hostname)
				continue
			}

			if udpMsg.Type == "file" {
				ReceiveFile(udpMsg)
				continue
			}

			// 添加消息到 UserMessages
			AddMessage(udpMsg)
		} else {
			log.Printf("unexpected address type: %T", remoteAddr)
		}
	}
}
func ClearAllUserMessages() {
	for ip, msg := range UserMessages {
		msg.Messages = []UdpMessage{} // 清空切片
		UserMessages[ip] = msg        // 更新映射中的值
	}
}
func GetMessages() UserMessage {
	return UserMessage{
		Messages: UserMessages,
		Onlines:  OnlineUsers,
	}
}
func AddMessage(msg UdpMessage) {
	if _, ok := UserMessages[msg.IP]; !ok {
		UserMessages[msg.IP] = &Messages{}
	}
	UserMessages[msg.IP].Messages = append(UserMessages[msg.IP].Messages, msg)
}

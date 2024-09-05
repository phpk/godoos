package localchat

import (
	"encoding/json"
	"godo/libs"
	"log"
	"net"
	"strings"
	"time"
)

type UdpMessage struct {
	Hostname string    `json:"hostname"`
	Type     string    `json:"type"`
	Time     time.Time `json:"time"`
	IP       string    `json:"ip"`
	Message  any       `json:"message"`
}
type UdpAddress struct {
	Hostname string    `json:"hostname"`
	IP       string    `json:"ip"`
	Time     time.Time `json:"time"`
}
type Messages struct {
	Messages []UdpMessage `json:"messages"`
}

var OnlineUsers []UdpAddress
var UserMessages = make(map[string]*Messages) // 使用指针类型

func init() {
	go InitBroadcast()
	go ListenForBroadcast()
}

// ListenForBroadcast 监听多播消息
func ListenForBroadcast() {
	broadcastAddr := GetBroadcastAddr()
	addr, err := net.ResolveUDPAddr("udp4", broadcastAddr)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}
	// 使用 ListenMulticastUDP 创建多播连接
	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		log.Fatalf("Failed to listen on UDP address: %v", err)
	}
	defer conn.Close()

	// 获取本地 IP 地址
	ips, err := libs.GetValidIPAddresses()
	if err != nil {
		log.Fatalf("Failed to get local IP addresses: %v", err)
	}

	// 开始监听多播消息
	buffer := make([]byte, 1024)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			continue
		}

		var udpMsg UdpMessage
		err = json.Unmarshal(buffer[:n], &udpMsg)
		if err != nil {
			log.Printf("Error unmarshalling JSON: %v", err)
			continue
		}
		// 从 remoteAddr 获取 IP 地址
		ip := remoteAddr.IP.String()
		if containArr(ips, ip) {
			continue
		}
		if udpMsg.Type == "online" {
			if !containIp(OnlineUsers, ip) {
				OnlineUsers = append(OnlineUsers, UdpAddress{Hostname: udpMsg.Hostname, IP: ip, Time: time.Now()})
				log.Printf("在线用户: %v", OnlineUsers)
			}
		}
		udpMsg.IP = ip
		if udpMsg.Type == "text" {
			addMessageToUserMessages(ip, udpMsg)
		}
		if udpMsg.Type == "file" {
			addMessageToUserMessages(ip, udpMsg)
			RecieveFile(udpMsg)
		}
		log.Printf("Received message from %s: %s", remoteAddr, udpMsg.Hostname)
	}
}
func containArr(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
func containIp(slice []UdpAddress, element string) bool {
	for _, v := range slice {
		if v.IP == element {
			return true
		}
	}
	return false
}

// 添加消息到 UserMessages
func addMessageToUserMessages(ip string, msg UdpMessage) {
	if _, ok := UserMessages[ip]; !ok {
		UserMessages[ip] = &Messages{}
	}
	UserMessages[ip].Messages = append(UserMessages[ip].Messages, msg)
}

// 清空 UserMessages 中所有 IP 的消息列表
func ClearAllUserMessages() {
	for ip, msg := range UserMessages {
		msg.Messages = []UdpMessage{} // 清空切片
		UserMessages[ip] = msg        // 更新映射中的值
	}
}
func GetUserMessages() map[string]*Messages {
	return UserMessages
}
func GetOnlineUsers() []UdpAddress {
	return OnlineUsers
}
func GetBroadcastAddr() string {
	return libs.GetUdpAddr()
}
func GetBroadcastPort() string {
	addr := GetBroadcastAddr()
	return addr[strings.LastIndex(addr, ":")+1:]
}

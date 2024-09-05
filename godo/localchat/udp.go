package localchat

import (
	"encoding/json"
	"godo/libs"
	"log"
	"net"
	"os"
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
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
}

var OnlineUsers []UdpAddress

func init() {
	go InitBroadcast()
	go ListenForBroadcast()
}
func InitBroadcast() {
	ticker := time.NewTicker(5 * time.Second) // 每 5 秒发送一次广播消息
	defer ticker.Stop()

	for range ticker.C {
		hostname, err := os.Hostname()
		if err != nil {
			log.Printf("Failed to get hostname: %v", err)
			continue
		}
		message := UdpMessage{
			Type:     "online",
			Hostname: hostname,
			Message:  "",
		}
		//发送多播消息
		broadcastAddr := GetBroadcastAddr()
		err = SendBroadcast(broadcastAddr, message)
		if err != nil {
			log.Println("Failed to send broadcast message:", err)
		}
	}
}
func SendBroadcast(broadcastAddr string, message UdpMessage) error {

	addr, err := net.ResolveUDPAddr("udp4", broadcastAddr)
	if err != nil {
		log.Printf("Failed to resolve UDP address %s: %v", broadcastAddr, err)
		return err
	}

	// 使用本地地址进行连接
	localAddr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:0")
	if err != nil {
		log.Printf("Failed to resolve local UDP address: %v", err)
		return err
	}

	conn, err := net.ListenUDP("udp4", localAddr)
	if err != nil {
		log.Printf("Failed to listen on UDP address %s: %v", broadcastAddr, err)
		return err
	}
	defer conn.Close()

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal JSON for %s: %v", broadcastAddr, err)
		return err
	}

	_, err = conn.WriteToUDP(data, addr)
	if err != nil {
		log.Printf("Failed to write to UDP address %s: %v", broadcastAddr, err)
		return err
	}

	log.Printf("发送消息到 %s 成功", broadcastAddr)
	return nil
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
		if !containIp(OnlineUsers, ip) {
			OnlineUsers = append(OnlineUsers, UdpAddress{Hostname: udpMsg.Hostname, IP: ip})
			log.Printf("在线用户: %v", OnlineUsers)
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
func GetOnlineUsers() []UdpAddress {
	return OnlineUsers
}
func GetBroadcastAddr() string {
	return "224.0.0.251:20249"
}
func GetBroadcastPort() string {
	addr := GetBroadcastAddr()
	return addr[strings.LastIndex(addr, ":")+1:]
}

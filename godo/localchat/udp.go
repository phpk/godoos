package localchat

import (
	"encoding/json"
	"fmt"
	"godo/libs"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type UdpMessage struct {
	Type     string      `json:"type"`
	IP       string      `json:"ip"`
	Hostname string      `json:"hostname"`
	Message  interface{} `json:"message"`
}

var OnlineUsers = make(map[string]UdpMessage)

// SendBroadcast 发送广播消息
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

		err = SendBroadcast(message)
		if err != nil {
			log.Println("Failed to send broadcast message:", err)
		}
	}
}

// SendToIP 向指定的 IP 地址发送 UDP 消息
func SendToIP(message UdpMessage) error {
	toIp := message.IP
	port := GetBroadcastPort()
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%s", toIp, port))
	if err != nil {
		log.Printf("Failed to resolve UDP address %s:%s: %v", toIp, port, err)
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
		log.Printf("Failed to listen on UDP address %s: %v", toIp, err)
		return err
	}
	defer conn.Close()

	// 获取本地 IP 地址
	localIP, err := GetLocalIP()
	if err != nil {
		log.Printf("Failed to get local IP address: %v", err)
		return err
	}
	message.IP = localIP

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal JSON for %s: %v", toIp, err)
		return err
	}

	_, err = conn.WriteToUDP(data, addr)
	if err != nil {
		log.Printf("Failed to write to UDP address %s: %v", toIp, err)
		return err
	}

	log.Printf("发送 UDP 消息到 %s 成功", toIp)
	return nil
}

// 获取 OnlineUsers 的最新状态
func GetOnlineUsers() map[string]UdpMessage {
	return OnlineUsers
}
func GetLocalIP() (string, error) {
	return libs.GetIPAddress()
}
func GetBroadcastPort() string {
	addr := GetBroadcastAddr()
	return addr[strings.LastIndex(addr, ":")+1:]
}

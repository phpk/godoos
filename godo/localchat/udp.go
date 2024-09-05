package localchat

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	BroadcastIP   = "255.255.255.255"
	BroadcastPort = 20249
)

type UdpMessage struct {
	Type     string      `json:"type"`
	IP       string      `json:"ip"`
	Hostname string      `json:"hostname"`
	Message  interface{} `json:"message"`
}

var broadcastAddr = fmt.Sprintf("%s:%d", BroadcastIP, BroadcastPort)
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
		err = SendBroadcast(message)
		if err != nil {
			log.Println("Failed to send broadcast message:", err)
		}
	}
}

func SendBroadcast(message UdpMessage) error {
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

	// 获取本地 IP 地址
	localIP, err := getLocalIP(conn)
	if err != nil {
		log.Printf("Failed to get local IP address: %v", err)
		return err
	}
	message.IP = localIP

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

	log.Printf("发送广播消息到 %s 成功", broadcastAddr)
	return nil
}

// ListenForBroadcast 监听广播消息
func ListenForBroadcast() {
	// 使用本地地址监听
	localAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", BroadcastPort))
	if err != nil {
		log.Fatalf("Failed to resolve local UDP address: %v", err)
	}

	conn, err := net.ListenUDP("udp4", localAddr)
	if err != nil {
		log.Fatalf("Failed to listen on UDP address: %v", err)
	}
	defer conn.Close()

	// 获取本地 IP 地址
	localIP, err := getLocalIP(conn)
	if err != nil {
		log.Printf("Failed to get local IP address: %v", err)
	}

	// 开始监听广播消息
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
		if udpMsg.IP == localIP {
			continue
		}
		OnlineUsers[udpMsg.IP] = udpMsg
		if udpMsg.Type == "file" {
			RecieveFile(udpMsg)
		}
		log.Printf("Received message from %s: %s", remoteAddr, udpMsg.Hostname)
	}
}

// SendToIP 向指定的 IP 地址发送 UDP 消息
func SendToIP(message UdpMessage) error {
	toIp := message.IP
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", toIp, BroadcastPort))
	if err != nil {
		log.Printf("Failed to resolve UDP address %s:%d: %v", toIp, BroadcastPort, err)
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
	localIP, err := getLocalIP(conn)
	if err != nil {
		log.Printf("Failed to get local IP address: %v", err)
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

// getLocalIP 获取本地 IP 地址
func getLocalIP(conn *net.UDPConn) (string, error) {
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

// HandleMessage 处理 HTTP 请求
func HandleMessage(w http.ResponseWriter, r *http.Request) {
	var msg UdpMessage
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	err := SendToIP(msg)
	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Text message send successfully")
}

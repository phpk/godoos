package localchat

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
)

type UdpMessage struct {
	Type     string      `json:"type"`
	IP       string      `json:"ip"`
	Hostname string      `json:"hostname"`
	Message  interface{} `json:"message"`
}

var broadcastAddr = "255.255.255.255:20249"
var OnlineUsers = make(map[string]UdpMessage)

// SendBroadcast 发送广播消息
func init() {
	go InitBroadcast()
	go ListenForBroadcast()
}

func InitBroadcast() {
	myIP, myHostname, err := GetMyIPAndHostname()
	if err != nil {
		log.Println("Failed to get IP and hostname:", err)
		return
	}
	message := UdpMessage{
		Type:     "online",
		IP:       myIP,
		Hostname: myHostname,
		Message:  "online",
	}
	SendBroadcast(message)
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
	localAddr, err := net.ResolveUDPAddr("udp4", ":20249")
	if err != nil {
		log.Fatalf("Failed to resolve local UDP address: %v", err)
	}

	conn, err := net.ListenUDP("udp4", localAddr)
	if err != nil {
		log.Fatalf("Failed to listen on UDP address: %v", err)
	}
	defer conn.Close()

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

		OnlineUsers[udpMsg.IP] = udpMsg
		if udpMsg.Type == "file" {
			RecieveFile(udpMsg)
		}
		log.Printf("Received message from %s: %s", remoteAddr, udpMsg.Hostname)
	}
}

// SendToIP 向指定的 IP 地址发送 UDP 消息
func SendToIP(ip string, message UdpMessage) error {
	addr, err := net.ResolveUDPAddr("udp4", ip)
	if err != nil {
		log.Printf("Failed to resolve UDP address %s: %v", ip, err)
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
		log.Printf("Failed to listen on UDP address %s: %v", ip, err)
		return err
	}
	defer conn.Close()

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal JSON for %s: %v", ip, err)
		return err
	}

	_, err = conn.WriteToUDP(data, addr)
	if err != nil {
		log.Printf("Failed to write to UDP address %s: %v", ip, err)
		return err
	}

	log.Printf("发送 UDP 消息到 %s 成功", ip)
	return nil
}

// 获取 OnlineUsers 的最新状态
func GetOnlineUsers() map[string]UdpMessage {
	return OnlineUsers
}

func HandleMessage(w http.ResponseWriter, r *http.Request) {
	var msg UdpMessage
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	ip := msg.IP
	preferredIP, err := GetMyIp()
	if err != nil {
		http.Error(w, "Failed to get preferred IP", http.StatusInternalServerError)
		return
	}
	msg.IP = preferredIP
	err = SendToIP(ip, msg)
	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Text message send successfully")
}

package localchat

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"time"
)

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
			Time:     time.Now(),
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

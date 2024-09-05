//go:build !windows
// +build !windows

package localchat

import (
	"encoding/json"
	"log"
	"net"
)

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

	// 获取本地 IP 地址
	localIP, err := GetLocalIP()
	log.Printf("本地 IP 地址: %s", localIP)
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
	localIP, err := GetLocalIP()
	if err != nil {
		log.Printf("Failed to get local IP address: %v", err)
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

func GetBroadcastAddr() string {
	//return libs.GetUdpAddr()
	return "224.0.0.251:20249"
}

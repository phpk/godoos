//go:build !windows
// +build !windows

package localchat

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func StartServiceDiscovery() {
	conn, err := net.ListenPacket("udp", broadcastAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("Received message: %s from %s\n", buffer[:n], addr)
		var udpMsg UdpMessage
		err = json.Unmarshal(buffer[:n], &udpMsg)
		if err != nil {
			fmt.Printf("Error unmarshalling JSON: %v\n", err)
			continue
		}
		log.Printf("Get message: %+v", udpMsg)
		OnlineUsers[udpMsg.IP] = udpMsg
	}
}

/*
func StartServiceDiscovery(port int, clientPort int) {
	addr := fmt.Sprintf(":%d", port)
	serviceAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}

	conn, err := net.ListenUDP("udp", serviceAddr)
	if err != nil {
		log.Fatalf("Failed to listen on UDP: %v", err)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Failed to get hostname: %v", err)
	}

	for {
		// 监听UDP广播请求
		_, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			continue
		}

		// 构建广播消息，包含IP和主机名
		ip, _, _ := net.SplitHostPort(serviceAddr.String())
		response := fmt.Sprintf("Chat Server is running at IP: %s, Hostname: %s", ip, hostname)
		broadcastMessage(conn, response, clientPort)
	}
}

func broadcastMessage(conn *net.UDPConn, message string, clientPort int) {
	// 广播给所有客户端
	broadcastAddr := &net.UDPAddr{IP: net.ParseIP("255.255.255.255"), Port: clientPort}
	_, err := conn.WriteToUDP([]byte(message), broadcastAddr)
	if err != nil {
		log.Printf("Failed to broadcast message: %v", err)
	}
}
*/

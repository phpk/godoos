package localchat

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func StartServiceDiscovery() {
	// 解析多播地址
	addr, err := net.ResolveUDPAddr("udp4", broadcastAddr)
	if err != nil {
		fmt.Println("Error resolving multicast address:", err)
		return
	}

	// 监听本地网络接口上的多播地址
	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		fmt.Println("Error listening on multicast address:", err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %v,addr:%v", err, addr)
			continue
		}
		//fmt.Printf("Received message: %s from %s\n", buffer[:n], addr)

		var udpMsg UdpMessage
		err = json.Unmarshal(buffer[:n], &udpMsg)
		if err != nil {
			fmt.Printf("Error unmarshalling JSON: %v\n", err)
			continue
		}
		//log.Printf("Get message: %+v", udpMsg)
		OnlineUsers[udpMsg.IP] = udpMsg
	}
}

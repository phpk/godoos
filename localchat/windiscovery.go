//go:build windows
// +build windows

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
	// 设置组播通信的网路接口
	// intf,err := net.InterfaceByName("eth0")
	// if err != nil {
	// 	fmt.Println("Error getting network interface:", err)
	// 	return
	// }
	// conn.setMulticastInterface(intf)

	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
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

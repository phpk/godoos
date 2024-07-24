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

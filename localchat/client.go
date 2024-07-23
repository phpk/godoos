package localchat

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

// func StartServiceDiscovery() {
// 	conn, err := net.ListenPacket("udp", broadcastAddr)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer conn.Close()

//		buffer := make([]byte, 1024)
//		for {
//			n, addr, err := conn.ReadFrom(buffer)
//			if err != nil {
//				fmt.Println(err)
//				continue
//			}
//			fmt.Printf("Received message: %s from %s\n", buffer[:n], addr)
//			var udpMsg UdpMessage
//			err = json.Unmarshal(buffer[:n], &udpMsg)
//			if err != nil {
//				fmt.Printf("Error unmarshalling JSON: %v\n", err)
//				continue
//			}
//			log.Printf("Get message: %+v", udpMsg)
//			OnlineUsers[udpMsg.IP] = udpMsg
//		}
//	}
func DiscoverServers() {
	broadcastTicker := time.NewTicker(broadcartTime)
	done := make(chan struct{}) // New channel to signal when to stop

	go func() {
		for {
			select {
			case <-broadcastTicker.C:
				conn, err := net.Dial("udp", broadcastAddr)
				if err != nil {
					fmt.Println(err)
					continue
				}

				myIP, myHostname, err := getMyIPAndHostname()
				if err != nil {
					log.Printf("Failed to get my IP and hostname: %v", err)
					conn.Close()
					continue
				}

				msg := UdpMessage{
					Type:     "online",
					Hostname: myHostname,
					IP:       myIP,
					Message:  time.Now().Format("2006-01-02 15:04:05"),
				}
				jsonData, err := json.Marshal(msg)
				if err != nil {
					log.Printf("Failed to marshal message to JSON: %v", err)
					conn.Close()
					continue
				}
				log.Printf("Sending message: %+v", msg)
				_, err = conn.Write(jsonData)
				if err != nil {
					fmt.Println(err)
				}
				conn.Close() // Close the connection after use
			case <-done: // Signal to stop
				broadcastTicker.Stop()
				return
			}
		}
	}()

	// You might have some condition here to decide when to stop the ticker
	// For example, a signal handling mechanism or a specific duration.
	// After that condition is met, you would call:
	// close(done)
}

// 获取自己的IP地址和主机名
func getMyIPAndHostname() (string, string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", "", fmt.Errorf("failed to get hostname: %w", err)
	}

	addrs, err := net.Interfaces()
	if err != nil {
		return "", "", fmt.Errorf("failed to get network interfaces: %w", err)
	}

	var preferredIP net.IP
	for _, iface := range addrs {
		if iface.Flags&net.FlagUp == 0 {
			// Skip interfaces that are not up
			continue
		}
		ifAddrs, err := iface.Addrs()
		if err != nil {
			continue // Ignore this interface if we can't get its addresses
		}
		for _, addr := range ifAddrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			default:
				continue
			}
			if ip.IsLoopback() {
				continue // Skip loopback addresses
			}
			if ip.To4() != nil && (ip.IsPrivate() || ip.IsGlobalUnicast()) {
				// Prefer global unicast or private addresses over link-local
				preferredIP = ip
				break
			}
		}
		if preferredIP != nil {
			// Found a preferred IP, break out of the loop
			break
		}
	}

	if preferredIP == nil {
		return "", "", fmt.Errorf("no preferred non-loopback IPv4 address found")
	}

	return preferredIP.String(), hostname, nil
}

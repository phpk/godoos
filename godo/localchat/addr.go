package localchat

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"sync"
	"time"

	"godo/libs"
)

// 发送 UDP 包并忽略响应
func sendUDPPacket(ip string) error {
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("error getting hostname: %v", err)
	}
	payload := UdpMessage{
		Type:     "heartbeat",
		Hostname: hostname,
		Time:     time.Now(),
	}
	log.Printf("sending ip: %+v", ip)
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshalling payload: %v", err)
		return fmt.Errorf("error marshalling payload: %v", err)
	}

	conn, err := net.Dial("udp", fmt.Sprintf("%s:56780", ip))
	if err != nil {
		log.Printf("error dialing UDP: %v", err)
		return fmt.Errorf("error dialing UDP: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write(payloadBytes)
	if err != nil {
		log.Printf("error writing UDP payload: %v", err)
		return fmt.Errorf("error writing UDP payload: %v", err)
	}

	return nil
}

func concurrentGetIpInfo(ips []string) {
	// 获取本地 IP 地址
	hostips, err := libs.GetValidIPAddresses()
	if err != nil {
		log.Printf("failed to get local IP addresses: %v", err)
		return
	}

	var wg sync.WaitGroup
	maxConcurrency := runtime.NumCPU()

	semaphore := make(chan struct{}, maxConcurrency)

	failedIPs := make(map[string]bool)

	for _, ip := range ips {
		if containArr(hostips, ip) || failedIPs[ip] {
			continue
		}

		wg.Add(1)
		semaphore <- struct{}{}

		go func(ip string) {
			defer wg.Done()
			defer func() { <-semaphore }()
			err := sendUDPPacket(ip)
			if err != nil {
				log.Printf("Failed to send packet to IP %s: %v", ip, err)
				failedIPs[ip] = true // 标记失败的 IP
			}
		}(ip)
	}

	wg.Wait()
}

func CheckOnline() {
	// 清除 OnlineUsers 映射表
	CleanOnlineUsers()

	ips := libs.GetChatIPs()
	// 启动并发处理
	concurrentGetIpInfo(ips)

	log.Printf("online users: %v", OnlineUsers)
}

func CleanOnlineUsers() {
	OnlineUsers = make(map[string]UserStatus)
}

func containArr(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func UpdateUserStatus(ip string, hostname string) UserStatus {
	OnlineUsers[ip] = UserStatus{
		Hostname: hostname,
		IP:       ip,
		Time:     time.Now(),
	}
	return OnlineUsers[ip]
}

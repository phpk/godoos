package localchat

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"

	"godo/libs"
)

type UserStatus struct {
	Hostname string    `json:"hostname"`
	IP       string    `json:"ip"`
	Time     time.Time `json:"time"`
}

var OnlineUsers = make(map[string]UserStatus)

type UDPPayload struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

func getHostname(ip string) (string, error) {
	hostname, err := net.LookupAddr(ip)
	if err != nil {
		return "", fmt.Errorf("error getting hostname: %v", err)
	}
	if len(hostname) > 0 {
		return hostname[0], nil
	}
	return "", fmt.Errorf("no hostname found for IP: %s", ip)
}

// 发送 UDP 包并忽略响应
func sendUDPPacket(ip string) error {
	payload := UDPPayload{
		Action: "check",
		Data:   "",
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
			} else {
				hostname, err := getHostname(ip)
				if err != nil {
					log.Printf("Failed to get hostname for IP %s: %v", ip, err)
				} else {
					OnlineUsers[ip] = UserStatus{
						Hostname: hostname,
						IP:       ip,
						Time:     time.Now(),
					}
				}
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

func HandleHeartbeat(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr // 可以根据实际情况获取 IP
	hostname, err := getHostname(ip)
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, "Failed to get hostname")
		return
	}
	userStatus := UpdateUserStatus(ip, hostname)
	libs.SuccessMsg(w, userStatus, "Heartbeat received")
}

func UpdateUserStatus(ip string, hostname string) UserStatus {
	OnlineUsers[ip] = UserStatus{
		Hostname: hostname,
		IP:       ip,
		Time:     time.Now(),
	}
	return OnlineUsers[ip]
}

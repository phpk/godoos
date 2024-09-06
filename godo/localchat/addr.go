// MIT License
//
// Copyright (c) 2024 godoos.com
// Email: xpbb@qq.com
// GitHub: github.com/phpk/godoos
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// MIT License
//
// Copyright (c) 2024 godoos.com
// Email: xpbb@qq.com
// GitHub: github.com/phpk/godoos
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package localchat

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
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
	// 获取 ARP 缓存中的 IP 地址
	validIPs, err := getArpCacheIPs()
	if err != nil {
		log.Printf("failed to get ARP cache IPs: %v", err)
		return
	}
	var wg sync.WaitGroup
	maxConcurrency := runtime.NumCPU()

	semaphore := make(chan struct{}, maxConcurrency)

	failedIPs := make(map[string]bool)

	for _, ip := range ips {
		if containArr(hostips, ip) || failedIPs[ip] || !containArr(validIPs, ip) {
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

// 获取 ARP 缓存中的 IP 地址
func getArpCacheIPs() ([]string, error) {
	var cmd *exec.Cmd
	var out []byte
	var err error

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("arp", "-a")
	case "linux":
		cmd = exec.Command("arp", "-n")
	case "darwin": // macOS
		cmd = exec.Command("arp", "-l", "-a")
	default:
		return nil, fmt.Errorf("unsupported operating system: %v", runtime.GOOS)
	}

	out, err = cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error executing arp command: %v", err)
	}

	lines := strings.Split(string(out), "\n")
	var ips []string

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			ip := fields[0]
			if ip != "<incomplete>" && net.ParseIP(ip) != nil {
				ips = append(ips, ip)
			}
		}
	}

	return ips, nil
}

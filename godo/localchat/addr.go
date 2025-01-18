/*
 * GodoOS - A lightweight cloud desktop
 * Copyright (C) 2024 https://godoos.com
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
	"godo/store"
)

var onlineUsersMutex sync.Mutex

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
	failedIPsMutex := sync.Mutex{}
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
				failedIPsMutex.Lock()   // 加锁
				failedIPs[ip] = true    // 标记失败的 IP
				failedIPsMutex.Unlock() // 解锁
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
	onlineUsersMutex.Lock()
	defer onlineUsersMutex.Unlock()
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
	onlineUsersMutex.Lock()
	defer onlineUsersMutex.Unlock()
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
	cmd = store.SetHideConsoleCursor(cmd)
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

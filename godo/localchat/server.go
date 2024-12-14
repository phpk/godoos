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
	"godo/libs"
	"log"
	"net"
	"strconv"
	"time"
)

type UdpMessage struct {
	Hostname string    `json:"hostname"`
	Type     string    `json:"type"`
	Time     time.Time `json:"time"`
	IP       string    `json:"ip"`
	Message  any       `json:"message"`
}

type UserMessage struct {
	Messages map[string][]UdpMessage `json:"messages"`
	Onlines  map[string]UserStatus   `json:"onlines"`
}
type UserStatus struct {
	Hostname string    `json:"hostname"`
	IP       string    `json:"ip"`
	Time     time.Time `json:"time"`
	Avatar   string    `json:"avatar"`
}

var OnlineUsers = make(map[string]UserStatus)
var UserMessages = make(map[string][]UdpMessage)

func init() {
	go UdpServer()
	go CheckOnlines()
}

func CheckOnlines() {
	//CheckOnline()
	chatIpSetting := libs.GetChatIpSetting()
	cktime, err := strconv.Atoi(chatIpSetting.CheckTime)
	if err != nil {
		cktime = 15
	}
	checkTimeDuration := time.Duration(cktime) * time.Second

	ticker := time.NewTicker(checkTimeDuration)
	defer ticker.Stop()
	for range ticker.C {
		// 检查客户端是否已断开连接
		CheckOnline()
	}
}

// UDP 服务器端逻辑
func UdpServer() {
	// 监听 UDP 端口
	listener, err := net.ListenPacket("udp", ":56780")
	if err != nil {
		log.Fatalf("error setting up listener: %v", err)
	}
	defer listener.Close()

	log.Println("UDP server started on :56780")

	// 监听 UDP 请求
	for {
		buffer := make([]byte, 1024)

		n, remoteAddr, err := listener.ReadFrom(buffer)
		if err != nil {
			log.Printf("error reading from UDP: %v", err)
			continue
		}

		log.Printf("Received UDP packet from %v: %s", remoteAddr, buffer[:n])
		// 从 remoteAddr 获取 IP 地址
		udpAddr, ok := remoteAddr.(*net.UDPAddr)
		if !ok {
			log.Printf("unexpected address type: %T", remoteAddr)
			continue
		}
		ip := udpAddr.IP.String()
		// 解析 UDP 数据
		var udpMsg UdpMessage
		err = json.Unmarshal(buffer[:n], &udpMsg)
		if err != nil {
			log.Printf("error unmarshalling UDP message: %v", err)
			continue
		}
		udpMsg.IP = ip

		if udpMsg.Type == "heartbeat" {
			UpdateUserStatus(udpMsg.IP, udpMsg.Hostname)
			continue
		}

		if udpMsg.Type == "image" {
			filePath, err := ReceiveImg(udpMsg)
			if err != nil {
				log.Printf("error receiving image: %v", err)
				continue
			}
			udpMsg.Message = filePath
		}
		// 添加消息到 UserMessages
		AddMessage(udpMsg)

	}
}

func ClearAllUserMessages() {
	UserMessages = make(map[string][]UdpMessage)
}

func GetMessages() UserMessage {
	return UserMessage{
		Messages: UserMessages,
		Onlines:  OnlineUsers,
	}
}

func AddMessage(msg UdpMessage) {
	// 检查 UserMessages 中是否已经有这个 IP 地址的消息列表
	if _, ok := UserMessages[msg.IP]; !ok {
		// 如果没有，则创建一个新的消息列表
		UserMessages[msg.IP] = []UdpMessage{}
	}

	// 将新消息添加到对应 IP 地址的消息列表中
	UserMessages[msg.IP] = append(UserMessages[msg.IP], msg)
}

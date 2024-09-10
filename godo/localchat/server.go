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
	"godo/libs"
	"log"
	"net"
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
	checkTimeDuration := time.Duration(chatIpSetting.CheckTime) * time.Second

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

		// if udpMsg.Type == "file" {
		// 	ReceiveFile(udpMsg)
		// 	continue
		// }
		if udpMsg.Type == "fileAccessed" {
			HandlerSendFile(udpMsg)
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

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
package sys

import (
	"encoding/json"
	"fmt"
	"godo/localchat"
	"net/http"
	"sync"
	"time"
)

// 定义消息结构体
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// 定义一个结构体来保存客户端连接
type client struct {
	conn http.ResponseWriter
}

var clients = make(map[*client]bool)
var mutex sync.Mutex

// 处理客户端连接
func HandleSystemEvents(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// 创建客户端连接对象
	c := &client{conn: w}

	// 添加客户端到列表
	mutex.Lock()
	clients[c] = true
	mutex.Unlock()
	// 发送初始消息
	updateInfo, err := GetUpdateInfo()
	if err == nil {
		SendToClient(c, Message{
			Type: "update",
			Data: updateInfo,
		})
	}

	// 监听客户端关闭
	defer func() {
		mutex.Lock()
		delete(clients, c)
		mutex.Unlock()
	}()

	// 使用定时器轮询客户端请求
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		// 检查客户端是否已断开连接
		if r.Context().Err() != nil {
			return
		}
		userMessages := localchat.GetMessages()
		msg := Message{
			Type: "localchat",
			Data: userMessages,
		}
		Broadcast(msg)
		localchat.ClearAllUserMessages()
	}
}

// 向客户端发送消息
func SendToClient(c *client, msg Message) {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	fmt.Fprintf(c.conn, "data: %s\n\n", jsonMsg)
	c.conn.(http.Flusher).Flush()
}

// 广播消息给所有客户端
func Broadcast(msg Message) {
	mutex.Lock()
	defer mutex.Unlock()
	for c := range clients {
		SendToClient(c, msg)
	}
}

// 每隔一段时间广播一条消息
func SendMessagePeriodically(msg Message) {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		Broadcast(msg)
	}
}

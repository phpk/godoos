/*
 * godoos - A lightweight cloud desktop
 * Copyright (C) 2024 godoos.com
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

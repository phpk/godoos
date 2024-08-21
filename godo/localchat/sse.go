package localchat

import (
	"encoding/json"
	"fmt"
	"godo/libs"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	go StartServiceDiscovery()
	go DiscoverServers()
}

func SseHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	ticker := time.NewTicker(broadcartTime) // 每3秒检查一次在线用户
	defer ticker.Stop()
	// 处理新消息
	ctx := r.Context()
	// 使用Context来监听请求的取消
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in SSE goroutine: %v", r)
			}
		}()

		for {
			select {
			case <-ctx.Done(): // 当请求被取消时，退出循环
				return
			case msg := <-messageChan:
				// 构造JSON数据
				jsonData, err := json.Marshal(msg)
				if err != nil {
					log.Printf("Failed to marshal message to JSON: %v", err)
					continue
				}

				// 通过SSE发送JSON数据
				fmt.Fprintf(w, "data: %s\n\n", string(jsonData))
				flusher.Flush()
			}
		}
	}()
	myIP, myHostname, err := getMyIPAndHostname()
	if err != nil {
		log.Printf("Failed to get my IP and hostname: %v", err)
		return
	}
	for {
		select {
		case <-ticker.C: // 每隔一段时间检查并广播在线用户
			var userList []UdpMessage
			// 首先将自己的IP和主机名放入列表
			myMsg := UdpMessage{
				IP:       myIP,
				Hostname: myHostname,
				Type:     "online",
				Message:  time.Now().Format("2006-01-02 15:04:05"),
			}
			userList = append(userList, myMsg)
			//log.Printf("Online users: %v", OnlineUsers)
			for ip, info := range OnlineUsers {
				if ip != myIP { // 确保不重复添加自己
					userList = append(userList, info)
				}
			}
			res := UserList{Type: "user_list", Content: userList}
			// 将用户列表转换为JSON字符串
			jsonData, err := json.Marshal(res)
			if err != nil {
				log.Printf("Failed to marshal online users to JSON: %v", err)
				continue
			}
			// 通过SSE发送JSON数据
			fmt.Fprintf(w, "data: %s\n\n", string(jsonData))
			flusher.Flush()
		}
	}
}
func HandleMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var msg Message
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	// 将消息放入messageChan
	messageChan <- msg
	//log.Printf("Received text message from %s: %s", msg.SenderInfo.IP, msg.Content)
	// 这里可以添加存储文本消息到数据库或其他处理逻辑
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Text message send successfully")
}
func CheckUserHanlder(w http.ResponseWriter, r *http.Request) {
	res := map[string]any{}
	res["code"] = 0
	res["message"] = "ok"
	// 获取主机名
	hostname, err := os.Hostname()
	if err == nil {
		hostname = "Unknown"
	}
	ip, _ := libs.GetIPAddress()

	res["data"] = map[string]any{
		"ip":       ip,
		"hostname": hostname,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

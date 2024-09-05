package localchat

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/ipv4"
)

type UdpMessage struct {
	Type     string      `json:"type"`
	IP       string      `json:"ip"`
	Hostname string      `json:"hostname"`
	Message  interface{} `json:"message"`
}

var multicastAddrs = []string{
	"239.255.255.250:2024",
	"239.255.255.251:2024",
	"224.0.0.251:1234",
}
var OnlineUsers = make(map[string]UdpMessage)

// SendMulticast 发送多播消息
func init() {
	go InitMulticast()
	go ListenForMulticast()
}

func InitMulticast() {
	myIP, myHostname, err := GetMyIPAndHostname()
	if err != nil {
		log.Println("Failed to get IP and hostname:", err)
		return
	}
	message := UdpMessage{
		Type:     "online",
		IP:       myIP,
		Hostname: myHostname,
		Message:  "online",
	}
	SendMulticast(message)
}

func SendMulticast(message UdpMessage) error {
	for _, addrStr := range multicastAddrs {
		addr, err := net.ResolveUDPAddr("udp4", addrStr)
		if err != nil {
			log.Printf("Failed to resolve UDP address %s: %v", addrStr, err)
			continue
		}

		// 使用本地地址进行连接
		localAddr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:0")
		if err != nil {
			log.Printf("Failed to resolve local UDP address: %v", err)
			continue
		}

		conn, err := net.ListenUDP("udp4", localAddr)
		if err != nil {
			log.Printf("Failed to listen on UDP address %s: %v", addrStr, err)
			continue
		}
		defer conn.Close()

		data, err := json.Marshal(message)
		if err != nil {
			log.Printf("Failed to marshal JSON for %s: %v", addrStr, err)
			continue
		}

		_, err = conn.WriteToUDP(data, addr)
		if err != nil {
			log.Printf("Failed to write to UDP address %s: %v", addrStr, err)
			continue
		}

		log.Printf("发送多播消息到 %s 成功", addrStr)
	}

	return nil
}

// ListenForMulticast 监听所有多播消息
func ListenForMulticast() {
	// 使用本地地址进行监听
	localAddr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:2024")
	if err != nil {
		log.Printf("Error resolving local UDP address: %v", err)
		return
	}

	conn, err := net.ListenUDP("udp4", localAddr)
	if err != nil {
		log.Printf("Error listening on UDP address: %v", err)
		return
	}
	defer conn.Close()

	udpConn := ipv4.NewPacketConn(conn)

	// 获取本地网络接口
	localIfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("Failed to get local interfaces: %v", err)
	}

	// 选择一个可用的网络接口
	var localInterface *net.Interface
	for _, iface := range localIfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			log.Printf("Failed to get addresses for interface %s: %v", iface.Name, err)
			continue
		}
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					localInterface = &iface
					break
				}
			}
		}
		if localInterface != nil {
			break
		}
	}

	if localInterface == nil {
		log.Fatal("No suitable network interface found")
	}

	// 加入所有多播组
	// for _, addrStr := range multicastAddrs {
	// 	multicastGroup, err := net.ResolveUDPAddr("udp4", addrStr)
	// 	if err != nil {
	// 		log.Printf("Error resolving UDP address %s: %v", addrStr, err)
	// 		continue
	// 	}

	// 	if err := udpConn.JoinGroup(localInterface, multicastGroup); err != nil {
	// 		log.Printf("Failed to join multicast group %s: %v", multicastGroup.String(), err)
	// 		continue
	// 	}

	// 	log.Printf("成功加入多播组 %s", addrStr)
	// }

	// 开始监听多播消息
	buffer := make([]byte, 1024)
	for {
		n, _, src, err := udpConn.ReadFrom(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %v, addr: %v", err, src)
			continue
		}

		var udpMsg UdpMessage
		err = json.Unmarshal(buffer[:n], &udpMsg)
		if err != nil {
			log.Printf("Error unmarshalling JSON: %v", err)
			continue
		}

		OnlineUsers[udpMsg.IP] = udpMsg
		if udpMsg.Type == "file" {
			RecieveFile(udpMsg)
		}
		log.Printf("Received message from %s: %s", udpMsg.IP, udpMsg.Hostname)
	}
}

// SendToIP 向指定的 IP 地址发送 UDP 消息
func SendToIP(ip string, message UdpMessage) error {
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:2024", ip))
	if err != nil {
		log.Printf("Failed to resolve UDP address %s: %v", ip, err)
		return err
	}

	// 使用本地地址进行连接
	localAddr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:0")
	if err != nil {
		log.Printf("Failed to resolve local UDP address: %v", err)
		return err
	}

	conn, err := net.ListenUDP("udp4", localAddr)
	if err != nil {
		log.Printf("Failed to listen on UDP address %s: %v", ip, err)
		return err
	}
	defer conn.Close()

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal JSON for %s: %v", ip, err)
		return err
	}

	_, err = conn.WriteToUDP(data, addr)
	if err != nil {
		log.Printf("Failed to write to UDP address %s: %v", ip, err)
		return err
	}

	log.Printf("发送 UDP 消息到 %s 成功", ip)
	return nil
}

// 获取 OnlineUsers 的最新状态
func GetOnlineUsers() map[string]UdpMessage {
	return OnlineUsers
}

func HandleMessage(w http.ResponseWriter, r *http.Request) {
	var msg UdpMessage
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	ip := msg.IP
	preferredIP, err := GetMyIp()
	if err != nil {
		http.Error(w, "Failed to get preferred IP", http.StatusInternalServerError)
		return
	}
	msg.IP = preferredIP
	err = SendToIP(ip, msg)
	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Text message send successfully")
}

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
	Type     string `json:"type"`
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	Message  any    `json:"message"`
}

// 多播地址列表
var multicastAddrs = []string{"239.255.255.250:2024", "239.255.255.251:2024", "224.0.0.251:1234", "224.0.0.1:1679"}
var OnlineUsers = make(map[string]UdpMessage) // 全局map，key为IP，value为主机名

// SendMulticast 发送多播消息
func init() {
	go InitMulticast()
	go ListenForMulticast()
}

func InitMulticast() {
	myIP, myHostname, err := GetMyIPAndHostname()
	if err != nil {
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
			return err
		}

		conn, err := net.DialUDP("udp4", nil, addr)
		if err != nil {
			return err
		}
		defer conn.Close()

		data, err := json.Marshal(message)
		if err != nil {
			return err
		}

		_, err = conn.WriteToUDP(data, addr)
		if err != nil {
			return err
		}

		fmt.Printf("发送多播消息到 %s 成功\n", addrStr)
	}

	return nil
}

// ListenForMulticast 监听多播消息
func ListenForMulticast() {
	// 选择一个可用的多播地址
	var multicastGroup *net.UDPAddr
	var err error
	for _, addrStr := range multicastAddrs {
		multicastGroup, err = net.ResolveUDPAddr("udp4", addrStr)
		if err != nil {
			fmt.Println("Error resolving UDP address:", err)
			continue
		}
		break
	}

	if multicastGroup == nil {
		fmt.Println("No available multicast address found")
		return
	}

	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: 2024})
	if err != nil {
		fmt.Println("Error listening on UDP address:", err)
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

	if err := udpConn.JoinGroup(localInterface, multicastGroup); err != nil {
		log.Fatalf("Failed to join multicast group: %v", err)
	}

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
			fmt.Printf("Error unmarshalling JSON: %v\n", err)
			continue
		}

		OnlineUsers[udpMsg.IP] = udpMsg
		if udpMsg.Type == "file" {
			RecieveFile(udpMsg)
		}
		fmt.Printf("Received message from %s: %s\n", udpMsg.IP, udpMsg.Hostname)
	}
}

// SendToIP 向指定的 IP 地址发送 UDP 消息
func SendToIP(ip string, message UdpMessage) error {
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:2024", ip))
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = conn.WriteToUDP(data, addr)
	if err != nil {
		return err
	}

	fmt.Printf("发送 UDP 消息到 %s 成功\n", ip)
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

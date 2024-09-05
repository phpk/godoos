package localchat

import (
	"net"
	"net/http"
	"strconv"

	"godo/libs"
)

func HandleAddr(w http.ResponseWriter, r *http.Request) {
	addr := r.URL.Query().Get("addr")
	if addr == "" {
		libs.ErrorMsg(w, "addr is empty")
		return
	}

	udpAddr := libs.GetUdpAddr()
	if addr == udpAddr {
		libs.ErrorMsg(w, "addr is same as current addr")
		return
	}

	// 检查是否为多播地址
	if !IsValidMulticastAddr(addr) {
		libs.ErrorMsg(w, "addr is not a valid multicast address")
		return
	}
	// 验证可访问性
	if IsMulticastAddrAccessible(addr) {
		save := libs.ReqBody{
			Value: addr,
			Name:  "udpAddr",
		}
		libs.SetConfig(save)
		libs.SuccessMsg(w, nil, "addr is a valid multicast address and accessible")
	} else {
		libs.ErrorMsg(w, "addr is a valid multicast address but not accessible")
	}
}

// 检查是否为多播地址
func IsValidMulticastAddr(addr string) bool {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}

	ip := net.ParseIP(host)
	if ip == nil || !isMulticastIP(ip) {
		return false
	}

	_, err = strconv.Atoi(port)
	return err == nil
}

// 检查 IP 是否为多播地址
func isMulticastIP(ip net.IP) bool {
	if ip.To4() != nil {
		return ip[0]&0xF0 == 0xE0 // 检查 IPv4 多播地址范围 224.0.0.0 - 239.255.255.255
	}
	return ip[0]&0xF0 == 0xE0 // 检查 IPv6 多播地址范围 FF00::/8
}

// 验证多播地址的可访问性
func IsMulticastAddrAccessible(addr string) bool {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}

	udpAddr, err := net.ResolveUDPAddr("udp4", net.JoinHostPort(host, port))
	if err != nil {
		return false
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, udpAddr)
	if err != nil {
		return false
	}
	defer conn.Close()

	// 发送一条测试消息
	testMsg := []byte("Test message")
	_, err = conn.WriteToUDP(testMsg, udpAddr)
	if err != nil {
		return false
	}

	// 接收一条测试消息
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return false
	}

	if n > 0 && string(buffer[:n]) == "Test message" {
		return true
	}

	return false
}

//go:build windows
// +build windows

package localchat

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func StartServiceDiscovery() {
	// 解析多播地址
	addr, err := net.ResolveUDPAddr("udp4", broadcastAddr)
	if err != nil {
		fmt.Println("Error resolving multicast address:", err)
		return
	}

	// 监听本地网络接口上的多播地址
	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		fmt.Println("Error listening on multicast address:", err)
		return
	}
	defer conn.Close()
	// 设置组播通信的网路接口
	// intf,err := net.InterfaceByName("eth0")
	// if err != nil {
	// 	fmt.Println("Error getting network interface:", err)
	// 	return
	// }
	// conn.setMulticastInterface(intf)

	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			continue
		}
		fmt.Printf("Received message: %s from %s\n", buffer[:n], addr)

		var udpMsg UdpMessage
		err = json.Unmarshal(buffer[:n], &udpMsg)
		if err != nil {
			fmt.Printf("Error unmarshalling JSON: %v\n", err)
			continue
		}
		log.Printf("Get message: %+v", udpMsg)
		OnlineUsers[udpMsg.IP] = udpMsg
	}
}

/*
// findMulticastInterface 查找适合多播的网络接口
func findMulticastInterface() (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if strings.HasPrefix(iface.Name, "eth") || strings.HasPrefix(iface.Name, "en") || strings.HasPrefix(iface.Name, "wlan") {
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				}
				if ip != nil && !ip.IsLoopback() {
					return &iface, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("no suitable multicast interface found")
}

// setupWindowsMulticast 为Windows平台设置多播
func setupWindowsMulticast(conn net.PacketConn, iface *net.Interface, multicastIP net.IP) error {
	// 获取文件描述符
	if sock, err := conn.(*net.UDPConn).File(); err != nil {
		return err
	} else {
		defer sock.Close()
		fd := sock.Fd()

		// 加入多播组
		var mreq windows.RawSockaddrInet4
		mreq.Family = windows.AF_INET
		mreq.Port = 0

		// 确保 multicastIP 是 IPv4 并转换为 [4]byte
		ipv4 := multicastIP.To4()
		if ipv4 == nil {
			return errors.New("multicast IP is not an IPv4 address")
		}
		copy(mreq.Addr[:], ipv4)

		if err := windows.SetsockoptIPv4MulticastInterface(fd, windows.IPMULTICAST_IF, &mreq); err != nil {
			return err
		}

		// 绑定到特定接口（可选，根据需求调整）
		ifaceIndex := uint32(iface.Index)
		if err := windows.BindToDevice(fd, windows.StringToUTF16Ptr(fmt.Sprintf("%d", ifaceIndex))); err != nil {
			return err
		}

		return nil
	}
}*/

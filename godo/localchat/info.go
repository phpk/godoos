package localchat

import (
	"fmt"
	"os"
)

// 获取自己的IP地址和主机名
func GetMyIPAndHostname() (string, string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", "", fmt.Errorf("failed to get hostname: %w", err)
	}
	preferredIP, err := GetLocalIP()
	if err != nil {
		return "", "", fmt.Errorf("failed to get IP address: %w", err)
	}

	return preferredIP, hostname, nil
}

// func GetLocalIP() (string, error) {
// 	addrs, err := net.Interfaces()
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get network interfaces: %w", err)
// 	}

// 	var preferredIP net.IP
// 	for _, iface := range addrs {
// 		if iface.Flags&net.FlagUp == 0 {
// 			// Skip interfaces that are not up
// 			continue
// 		}
// 		ifAddrs, err := iface.Addrs()
// 		if err != nil {
// 			continue // Ignore this interface if we can't get its addresses
// 		}
// 		for _, addr := range ifAddrs {
// 			var ip net.IP
// 			switch v := addr.(type) {
// 			case *net.IPNet:
// 				ip = v.IP
// 			case *net.IPAddr:
// 				ip = v.IP
// 			default:
// 				continue
// 			}
// 			if ip.IsLoopback() {
// 				continue // Skip loopback addresses
// 			}
// 			if ip.To4() != nil && (ip.IsPrivate() || ip.IsGlobalUnicast()) {
// 				// Prefer global unicast or private addresses over link-local
// 				preferredIP = ip
// 				break
// 			}
// 		}
// 		if preferredIP != nil {
// 			// Found a preferred IP, break out of the loop
// 			break
// 		}
// 	}

// 	if preferredIP == nil {
// 		return "", fmt.Errorf("no preferred non-loopback IPv4 address found")
// 	}
// 	return preferredIP.String(), nil
// }

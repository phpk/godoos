package libs

import (
	"errors"
	"log/slog"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetIpAddress(c *gin.Context) (ipAddress string) {
	// 尝试从不同的HTTP头部获取客户端IP
	ipAddress = c.Request.Header.Get("X-Real-IP")
	if ipAddress == "" || strings.EqualFold(ipAddress, "unknown") {
		ipAddress = c.Request.Header.Get("X-Forwarded-For")
		if ipAddress != "" {
			ipAddress = strings.Split(ipAddress, ",")[0]
		} else {
			ipAddress = c.Request.Header.Get("Proxy-Client-IP")
			if ipAddress == "" {
				ipAddress = c.Request.Header.Get("WL-Proxy-Client-IP")
			}
		}
	}

	if ipAddress == "" {
		ipAddress = c.ClientIP()
	}

	// 检测到是本机 IP, 读取其局域网 IP 地址
	if strings.HasPrefix(ipAddress, "127.0.0.1") || strings.HasPrefix(ipAddress, "[::1]") {
		ip, err := externalIP()
		if err != nil {
			slog.Error("GetIpAddress, externalIP, err: ", "error:", err)
		}
		ipAddress = ip.String()
	}

	return ipAddress
}

// 获取非 127.0.0.1 的局域网 IP
func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip != nil {
				return ip, nil
			}
		}
	}

	return nil, errors.New("no non-loopback addresses available")
}

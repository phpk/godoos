package libs

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// UserOsInfo 定义系统信息结构体
type UserOsInfo struct {
	MAC        string `json:"mac"`
	OS         string `json:"os"`
	Arch       string `json:"arch"`
	AppName    string `json:"app_name"`
	Hostname   string `json:"hostname"`
	TopType    string `json:"toptype"`
	UseType    string `json:"usetype"`
	SourceType string `json:"sourcetype"`
}

// GetSystemInfo 获取系统信息
func GenerateSystemInfo() UserOsInfo {
	info := UserOsInfo{}

	// 获取MAC地址
	mac, err := getMACAddress()
	if err == nil {
		info.MAC = mac
	} else {
		info.MAC = ""
	}

	// 获取主机名
	hostname, err := os.Hostname()
	if err == nil {
		info.Hostname = hostname
	} else {
		info.Hostname = ""
	}

	// 获取操作系统和架构信息
	info.OS = runtime.GOOS
	info.Arch = runtime.GOARCH

	info.AppName = "godoos"
	info.TopType = os.Getenv("GODOTOPTYPE")
	if info.TopType == "" {
		info.TopType = "web"
	}
	info.UseType = "person"
	info.SourceType = "open"
	return info
}

// GetIPAddress 获取本机IP地址
func GetIPAddress() (string, error) {
	ips, err := GetValidIPAddresses()
	if err != nil {
		return "", err
	}

	// 优先选择 192.168 开头的 IP 地址，并且第三个数字最小
	var bestIP string
	minThirdOctet := 256 // 初始化为最大值，表示还没有找到合适的 IP 地址

	for _, ipStr := range ips {
		ip := net.ParseIP(ipStr)
		if ip != nil && ip.To4() != nil && ip.String()[:7] == "192.168" {
			octets := strings.Split(ip.String(), ".")
			if len(octets) == 4 {
				thirdOctet, _ := strconv.Atoi(octets[2])
				if thirdOctet < minThirdOctet {
					minThirdOctet = thirdOctet
					bestIP = ipStr
				}
			}
		}
	}

	// 如果没有找到 192.168 开头的 IP 地址，则选择第一个有效的全局单播 IP 地址
	if bestIP == "" && len(ips) > 0 {
		bestIP = ips[0]
	}

	if bestIP == "" {
		return "", fmt.Errorf("no valid IP addresses found")
	}

	return bestIP, nil
}

// GetValidIPAddresses 获取所有有效 IP 地址
func GetValidIPAddresses() ([]string, error) {
	var validIPs []string
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil && !ip.IsLoopback() && ip.To4() != nil {
				// 过滤链路本地地址
				if !ip.IsLinkLocalUnicast() {
					validIPs = append(validIPs, ip.String())
				}
			}
		}
	}
	if len(validIPs) == 0 {
		return nil, fmt.Errorf("no valid non-private and non-link-local IP address found")
	}
	return validIPs, nil
}

// getMACAddress 获取MAC地址
func getMACAddress() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	var macs string
	// 记录所有接口的信息
	for _, iface := range ifaces {
		if iface.HardwareAddr != nil {
			macs += iface.HardwareAddr.String()
		}
	}
	//返回md5加密数据
	return Md5Encrypt(macs), nil
}
func Md5Encrypt(s string) string {
	hasher := md5.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GetSystemInfo 生成基于mac、os、arch信息的Base64编码字符串
func GetSystemInfo() (string, error) {
	// 构造系统信息对象
	lineseInfo, ok := GetConfig("osInfo")
	if !ok {
		return "", fmt.Errorf("未找到osInfo配置")
	}
	// 将系统信息序列化为JSON字符串
	jsonBytes, err := json.Marshal(lineseInfo)
	if err != nil {
		return "", err
	}
	// 对JSON字符串进行Base64编码
	encodedInfo := base64.StdEncoding.EncodeToString(jsonBytes)

	return encodedInfo, nil
}

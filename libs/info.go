package libs

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

// SystemInfo 定义系统信息结构体
type SystemInfo struct {
	MAC         string  `json:"mac"`
	CPUUsage    float64 `json:"cpu_usage"`
	OS          string  `json:"os"`
	Arch        string  `json:"arch"`
	MemoryTotal string  `json:"memory_total"`
}

// GetSystemInfo 获取系统信息
// GetSystemInfo 获取系统综合信息
func GetSystemInfo() map[string]string {
	info := make(map[string]string)

	// 获取本机IP地址
	ip, err := getIPAddress()
	if err == nil {
		info["ip"] = ip
	} else {
		info["ip"] = ""
	}

	// 获取MAC地址
	mac, err := getMACAddress()
	if err == nil {
		info["mac"] = mac
	} else {
		info["mac"] = ""
	}

	// 获取CPU信息
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err == nil {
		info["cpu"] = strconv.FormatFloat(cpuPercent[0], 'f', 2, 64) + "%"
	} else {
		info["cpu"] = ""
	}

	// 获取内存信息
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		info["memory_total"] = byteCountSI(memInfo.Total)
		info["memory_used"] = byteCountSI(memInfo.Used)
		info["memory_free"] = byteCountSI(memInfo.Free)
	} else {
		info["memory_total"] = ""
		info["memory_used"] = ""
		info["memory_free"] = ""
	}

	// 获取主机名
	hostname, err := os.Hostname()
	if err == nil {
		info["hostname"] = hostname
	} else {
		info["hostname"] = ""
	}

	// 获取操作系统和架构信息
	info["os"] = runtime.GOOS
	info["arch"] = runtime.GOARCH

	// 获取系统运行时间
	uptime, err := host.Uptime()
	if err == nil {
		info["uptime_seconds"] = strconv.FormatUint(uptime, 10)
	} else {
		info["uptime_seconds"] = ""
	}

	return info
}

// getIPAddress 获取本机IP地址
func getIPAddress() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
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
			return ip.String(), nil
		}
	}
	return "", fmt.Errorf("no valid IP address found")
}

// getMACAddress 获取MAC地址
func getMACAddress() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp != 0 && (strings.HasPrefix(iface.Name, "en") || strings.HasPrefix(iface.Name, "eth")) {
			return iface.HardwareAddr.String(), nil
		}
	}
	return "", fmt.Errorf("no active Ethernet interface found")
}

// byteCountSI 格式化字节大小为可读字符串
func byteCountSI(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

// GenerateBase64Info 生成基于mac、cpu、os、arch信息的Base64编码字符串
func GenerateBase64Info() (string, error) {
	// 获取必要的系统信息
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return "", err
	}
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}
	// 获取MAC地址
	mac, err := getMACAddress()
	if err == nil {
		return "", err
	}

	// 构造系统信息对象
	systemInfo := SystemInfo{
		MAC:         mac,
		CPUUsage:    cpuPercent[0],
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		MemoryTotal: byteCountSI(memInfo.Total),
	}

	// 将系统信息序列化为JSON字符串
	jsonBytes, err := json.Marshal(systemInfo)
	if err != nil {
		return "", err
	}

	// 对JSON字符串进行Base64编码
	encodedInfo := base64.StdEncoding.EncodeToString(jsonBytes)

	return encodedInfo, nil
}

// DecryptWithAES 使用AES解密数据
func DecryptWithAES(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

// VerifyAndDecrypt 验证时间戳并解密数据
func VerifyAndDecrypt(encodedInfo string, startTime time.Time, endTime time.Time, flag string) (*SystemInfo, error) {
	// Base64解码
	ciphertext, err := base64.StdEncoding.DecodeString(encodedInfo)
	if err != nil {
		return nil, err
	}

	// 生成AES密钥
	aesKey, err := GenerateAESKey(startTime, endTime, flag)
	if err != nil {
		return nil, err
	}

	// 解密数据
	decryptedBytes, err := DecryptWithAES(ciphertext, aesKey)
	if err != nil {
		return nil, err
	}

	// 反序列化JSON为SystemInfo
	var systemInfo SystemInfo
	if err := json.Unmarshal(decryptedBytes, &systemInfo); err != nil {
		return nil, err
	}
	// 获取当前系统MAC地址
	currentMAC, err := getMACAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get current MAC address: %w", err)
	}

	// 检查解密后的MAC地址是否与当前系统MAC地址一致
	if systemInfo.MAC != currentMAC {
		return nil, fmt.Errorf("decrypted MAC address does not match the current system's MAC address")
	}
	if systemInfo.OS != runtime.GOOS {
		return nil, fmt.Errorf("decrypted GOOS does not match the current system's runtime.GOOS")
	}
	// 验证时间有效性
	currentTime := time.Now()
	if currentTime.Before(startTime) || currentTime.After(endTime) {
		return nil, fmt.Errorf("the provided time window is invalid")
	}

	return &systemInfo, nil
}

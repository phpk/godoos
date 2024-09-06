package libs

import (
	"fmt"
	"strconv"
)

type UserChatIpSetting struct {
	CheckTime   int    `json:"CheckTime"`
	First       string `json:"First"`
	Second      string `json:"Second"`
	ThirdStart  string `json:"ThirdStart"`
	ThirdEnd    string `json:"ThirdEnd"`
	FourthStart string `json:"FourthStart"`
	FourthEnd   string `json:"FourthEnd"`
}

func GetDefaultChatIpSetting() UserChatIpSetting {
	return UserChatIpSetting{
		CheckTime:   60,
		First:       "192",
		Second:      "168",
		ThirdStart:  "1",
		ThirdEnd:    "1",
		FourthStart: "2",
		FourthEnd:   "99",
	}
}

func GetChatIpSetting() UserChatIpSetting {
	ips, ok := GetConfig("chatIpSetting")
	if !ok {
		return GetDefaultChatIpSetting()
	}
	return ips.(UserChatIpSetting)
}

// GenerateIPs 生成 IP 地址列表
func GenerateIPs(setting UserChatIpSetting) []string {
	var IPs []string

	thirdStart, _ := strconv.Atoi(setting.ThirdStart)
	thirdEnd, _ := strconv.Atoi(setting.ThirdEnd)
	fourthStart, _ := strconv.Atoi(setting.FourthStart)
	fourthEnd, _ := strconv.Atoi(setting.FourthEnd)

	if thirdStart == thirdEnd {
		if fourthStart == fourthEnd {
			// 第三位和第四位都相等，只生成一个 IP 地址
			IPs = append(IPs, fmt.Sprintf("%s.%s.%d.%d", setting.First, setting.Second, thirdStart, fourthStart))
		} else {
			// 第三位相等，第四位不相等，生成第四位的所有 IP 地址
			for j := fourthStart; j <= fourthEnd; j++ {
				IPs = append(IPs, fmt.Sprintf("%s.%s.%d.%d", setting.First, setting.Second, thirdStart, j))
			}
		}
	} else {
		// 第三位不相等，生成第三位和第四位的所有组合
		for i := thirdStart; i <= thirdEnd; i++ {
			for j := fourthStart; j <= fourthEnd; j++ {
				IPs = append(IPs, fmt.Sprintf("%s.%s.%d.%d", setting.First, setting.Second, i, j))
			}
		}
	}

	return IPs
}

// 示例函数
func GetChatIPs() []string {
	setting := GetChatIpSetting()
	return GenerateIPs(setting)
}

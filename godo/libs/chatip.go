/*
 * godoos - A lightweight cloud desktop
 * Copyright (C) 2024 godoos.com
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package libs

import (
	"fmt"
	"strconv"
)

type UserChatIpSetting struct {
	CheckTime   string `json:"checkTime"`
	First       string `json:"first"`
	Second      string `json:"second"`
	ThirdStart  string `json:"thirdStart"`
	ThirdEnd    string `json:"thirdEnd"`
	FourthStart string `json:"fourthStart"`
	FourthEnd   string `json:"fourthEnd"`
}

func GetDefaultChatIpSetting() UserChatIpSetting {
	return UserChatIpSetting{
		CheckTime:   "15",
		First:       "192",
		Second:      "168",
		ThirdStart:  "1",
		ThirdEnd:    "1",
		FourthStart: "2",
		FourthEnd:   "254",
	}
}

func GetChatIpSetting() UserChatIpSetting {
	ips, ok := GetConfig("chatIpSetting")
	if !ok {
		return GetDefaultChatIpSetting()
	}

	if ips, ok = ips.(UserChatIpSetting); !ok {
		return GetDefaultChatIpSetting()
	}
	return UserChatIpSetting{}
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

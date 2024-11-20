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

package store

// InstallInfo 描述了应用程序的安装信息。
// 包含应用程序的名称、下载地址、版本号等关键安装信息。
type InstallInfo struct {
	Name         string           `json:"name"`         // 应用程序名称。
	URL          string           `json:"url"`          // 应用程序下载地址或适配包的下载地址。
	Pkg          string           `json:"pkg"`          // 应用程序的官方下载地址。
	WebUrl       string           `json:"webUrl"`       // 应用程序的网页地址。
	IsDev        bool             `json:"isDev"`        // 标志位，表示是否为开发者版本。
	Version      string           `json:"version"`      // 应用程序的版本号。
	Desc         string           `json:"desc"`         // 应用程序的描述信息。
	Icon         string           `json:"icon"`         // 应用程序的图标路径。
	HasStart     bool             `json:"hasStart"`     // 标志位，表示是否显示启动和停止。
	HasRestart   bool             `json:"hasRestart"`   // 标志位，表示安装后是否需要重启。
	Setting      bool             `json:"setting"`      // 标志位，表示是否需要配置。
	Dependencies []Item           `json:"dependencies"` // 依赖项。
	Categrory    string           `json:"category"`     // 应用程序的分类。
	History      []InstallHastory `json:"history"`      // 应用程序的历史版本。
}
type InstallHastory struct {
	Version string `json:"version"`
	URL     string `json:"url"`
	Pkg     string `json:"pkg"` // 应用程序的官方下载地址。
}

// StoreInfo 维护了应用程序商店的信息。
// 包含应用程序的名称、图标、配置信息等。
type StoreInfo struct {
	Name     string           `json:"name"`     // 应用程序商店的名称。
	Icon     string           `json:"icon"`     // 应用程序商店的图标路径。
	Setting  Setting          `json:"setting"`  // 应用程序商店的配置信息。
	Config   map[string]any   `json:"config"`   // 应用程序的配置信息映射。
	Commands map[string][]Cmd `json:"commands"` // 应用程序的命令集合。
	Install  InstallStruct    `json:"install"`  // 安装应用程序的信息。
	Start    StartStruct      `json:"start"`    // 启动应用程序的信息。
}

// Setting 描述了应用程序的设置信息。
// 包含应用程序的二进制文件路径、配置文件路径等关键设置信息。
type Setting struct {
	BinPath      string `json:"binPath"`      // 应用程序二进制文件的路径。
	ConfPath     string `json:"confPath"`     // 应用程序配置文件的路径。
	ProgressName string `json:"progressName"` // 进程的名称。
	IsOn         bool   `json:"isOn"`         //是否守护进程运行。
}

// Item 是一个通用的键值对结构体，用于表示配置项或环境变量等。
type Item struct {
	Name  string `json:"name"`  // 配置项的名称。
	Value any    `json:"value"` // 配置项的值。
}

// Cmd 描述了一个命令的详细信息。
// 包含命令的名称、执行路径、环境变量等。
type Cmd struct {
	Name     string   `json:"name"`               // 命令的名称。
	FilePath string   `json:"filePath,omitempty"` // 命令文件的路径。
	Content  string   `json:"content,omitempty"`  // 命令的内容。
	BinPath  string   `json:"binPath,omitempty"`  // 执行命令的二进制文件路径。
	TplPath  string   `json:"tplPath,omitempty"`  // 命令的模板路径。
	Cmds     []string `json:"cmds,omitempty"`     // 要执行的子命令列表。
	Waiting  int      `json:"waiting"`            // 等待的时间。
	Kill     bool     `json:"kill"`               // 标志位，表示是否需要终止之前的命令。
	Envs     []Item   `json:"envs"`               // 命令执行时的环境变量。
}

// Install 描述了安装过程中的环境变量和命令列表。
type InstallStruct struct {
	InstallEnvs []Item   `json:"installEnvs"` // 安装过程中需要的环境变量。
	InstallCmds []string `json:"installCmds"` // 安装过程中需要执行的命令列表。
}
type StartStruct struct {
	StartEnvs  []Item   `json:"startEnvs"`  // 启动过程中需要的环境变量。
	BeforeCmds []string `json:"beforeCmds"` // 启动前需要执行的命令列表。
	StartCmds  []string `json:"startCmds"`  // 启动过程中需要执行的命令列表。
	AfterCmds  []string `json:"afterCmds"`  // 启动后需要执行的命令列表。
}

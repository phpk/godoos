// MIT License
//
// Copyright (c) 2024 godoos.com
// Email: xpbb@qq.com
// GitHub: github.com/phpk/godoos
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package store

import (
	"encoding/json"
	"fmt"
	"godo/libs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func InstallHandler(w http.ResponseWriter, r *http.Request) {
	pluginName := r.URL.Query().Get("name")
	if pluginName == "" {
		libs.ErrorMsg(w, "the app name is empty!")
		return
	}
	installInfo, err := Installation(pluginName)
	if err != nil {
		libs.ErrorData(w, installInfo, "the install.json is error:"+err.Error())
		return
	}
	libs.SuccessMsg(w, installInfo, "install the app success!")
}

// Installation 处理安装逻辑
func Installation(pluginName string) (InstallInfo, error) {
	var installInfo InstallInfo
	exePath := GetExePath(pluginName)
	if !libs.PathExists(exePath) {
		return installInfo, fmt.Errorf("the app path is not exists")
	}
	installInfo, err := GetInstallInfo(pluginName)
	if err != nil {
		return installInfo, fmt.Errorf("the install.json is error: %v", err)
	}
	if len(installInfo.Dependencies) > 0 {
		var needInstalls []Item
		for _, item := range installInfo.Dependencies {
			depInfo, err := GetInstallInfo(item.Value.(string))
			if err != nil {
				needInstalls = append(needInstalls, item)
				continue
			}
			if depInfo.Version != installInfo.Version {
				needInstalls = append(needInstalls, item)
			}
		}
		installInfo.Dependencies = needInstalls
		return installInfo, fmt.Errorf("dependencies require installation")
	}

	// Check if the plugin name matches the install.json name
	if pluginName != installInfo.Name {
		return installInfo, fmt.Errorf("plugin name does not match install.json")
	}
	// Copy static directory
	staticPath := filepath.Join(exePath, "static")
	if libs.PathExists(staticPath) {
		staticDir := libs.GetStaticDir()
		targetPath := filepath.Join(staticDir, pluginName)
		if !libs.PathExists(targetPath) {
			if err := os.Rename(staticPath, targetPath); err != nil {
				return installInfo, fmt.Errorf("error copying static directory: %w", err)
			}
			iconUrl, err := HandlerIcon(installInfo, targetPath)
			if err == nil {
				installInfo.Icon = iconUrl
				err = SaveInstallInfo(installInfo)
				if err != nil {
					return installInfo, fmt.Errorf("error saving install.json: %w", err)
				}
			}

		}
	}
	// Process store.json
	storeFile := filepath.Join(exePath, "store.json")
	if !libs.PathExists(storeFile) {
		return installInfo, nil
	}
	var storeInfo StoreInfo
	content, err := os.ReadFile(storeFile)
	if err != nil {
		return installInfo, fmt.Errorf("cannot read store.json: %w", err)
	}
	//设置 store.json
	exePath = strings.ReplaceAll(exePath, "\\", "/")
	contentBytes := []byte(strings.ReplaceAll(string(content), "{exePath}", exePath))
	//log.Printf("====content: %s", string(contentBytes))
	err = json.Unmarshal(contentBytes, &storeInfo)

	if err != nil {
		return installInfo, fmt.Errorf("error unmarshalling store.json: %w", err)
	}
	replacePlaceholdersInCmds(&storeInfo)
	storeInfo.Name = installInfo.Name

	// Save store info
	if err := SaveInfoFile(storeInfo); err != nil {
		return installInfo, fmt.Errorf("error saving store info.json: %w", err)
	}

	// Set up configuration
	if libs.PathExists(storeInfo.Setting.ConfPath + ".tpl") {
		if err := SaveStoreConfig(storeInfo, exePath); err != nil {
			return installInfo, fmt.Errorf("error saving config: %w", err)
		}
	}

	// Install the store
	if err := InstallStore(storeInfo); err != nil {
		return installInfo, fmt.Errorf("error installing app: %w", err)
	}

	return installInfo, nil
}

// UnInstallHandler 处理卸载请求
func UnInstallHandler(w http.ResponseWriter, r *http.Request) {
	pluginName := r.URL.Query().Get("name")
	if pluginName == "" {
		libs.ErrorMsg(w, "the app name is empty!")
		return
	}

	// Execute the uninstallation process
	if err := Uninstallation(pluginName); err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}

	// Success response
	libs.SuccessMsg(w, "success", "uninstall the app success!")
}

// Uninstallation 执行插件卸载逻辑
func Uninstallation(pluginName string) error {
	if err := StopCmd(pluginName); err != nil {
		log.Printf("stopping the app encountered an error: %s", err)
	}

	// Retrieve install information
	installInfo, err := GetInstallInfo(pluginName)
	if err != nil {
		return fmt.Errorf("error retrieving install.json: %w", err)
	}

	// Check if the app is in development mode
	if installInfo.IsDev {
		return nil // No further action required for dev mode apps
	}
	//检查是否有其他应用依赖于它
	installedList := GetInstalled()
	if len(installedList) > 1 {
		hasDeps := []Item{}
		for _, item := range installedList {
			for _, dep := range item.Dependencies {
				if dep.Value == pluginName {
					hasDeps = append(hasDeps, dep)
				}
			}
		}
		if len(hasDeps) > 0 {
			return fmt.Errorf("the app is being used by other applications: %v", hasDeps)
		}
	}
	storeInfo, err := GetStoreInfo(pluginName)
	if err == nil {
		if _, ok := storeInfo.Commands["uninstall"]; ok {
			replacePlaceholdersInCmds(&storeInfo)
			if err := RunCmds(storeInfo, "uninstall"); err != nil {
				return fmt.Errorf("error running uninstall command: %w", err)
			}
		}
	}
	// Remove the application directory
	exePath := GetExePath(pluginName)
	if libs.PathExists(exePath) {
		if err := os.RemoveAll(exePath); err != nil {
			return fmt.Errorf("error deleting the app: %w", err)
		}
	}

	// Remove the static directory
	staticDir := libs.GetStaticDir()
	staticPath := filepath.Join(staticDir, pluginName)
	if libs.PathExists(staticPath) {
		if err := os.RemoveAll(staticPath); err != nil {
			return fmt.Errorf("error deleting the static files: %w", err)
		}
	}

	return nil
}
func GetInstallPath(pluginName string) (string, error) {
	exePath := GetExePath(pluginName)
	installFile := filepath.Join(exePath, "install.json")
	if !libs.PathExists(installFile) {
		return "", fmt.Errorf("install.json is not exist:%s", installFile)
	}
	return installFile, nil
}
func GetInstallInfo(pluginName string) (InstallInfo, error) {
	var installInfo InstallInfo
	installFile, err := GetInstallPath(pluginName)
	if err != nil {
		return installInfo, err
	}
	return GetInstallInfoByPath(installFile)
}
func GetInstallInfoByPath(installFile string) (InstallInfo, error) {
	var installInfo InstallInfo
	content, err := os.ReadFile(installFile)
	if err != nil {
		return installInfo, err
	}
	err = json.Unmarshal(content, &installInfo)
	if err != nil {
		return installInfo, err
	}
	return installInfo, nil
}
func GetExePath(pluginName string) string {
	exeDir := libs.GetRunDir()
	return filepath.Join(exeDir, pluginName)
}

func convertValueToString(value any) string {
	var strValue string
	switch v := value.(type) {
	case string:
		strValue = v
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		strValue = fmt.Sprintf("%v", v)
	default:
		strValue = fmt.Sprintf("%v", v) // 处理其他类型，例如布尔型或更复杂的类型
	}
	return strValue
}

func InstallStore(storeInfo StoreInfo) error {
	err := SetEnvs(storeInfo.Install.InstallEnvs)
	if err != nil {
		return fmt.Errorf("failed to set install environment variable %s: %w", storeInfo.Name, err)
	}
	if len(storeInfo.Install.InstallCmds) > 0 {
		for _, cmdKey := range storeInfo.Install.InstallCmds {
			if _, ok := storeInfo.Commands[cmdKey]; ok {
				// 如果命令存在，你可以进一步处理 cmds
				RunCmds(storeInfo, cmdKey)
			}
		}

	}
	return nil
}
func RunCmd(storeInfo StoreInfo, cmd Cmd) error {
	if cmd.Name == "stop" {
		runStop(storeInfo)
	}
	if cmd.Name == "start" {
		runStart(storeInfo)
	}
	if cmd.Name == "startApp" {
		RunStartApp(cmd.Content)
	}
	if cmd.Name == "stopApp" {
		RunStopApp(cmd.Content)
	}
	if cmd.Name == "restart" {
		runRestart(storeInfo)
	}
	if cmd.Name == "exec" {
		runExec(storeInfo, cmd)
	}
	if cmd.Name == "writeFile" {
		err := WriteFile(cmd)
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}
	if cmd.Name == "changeFile" {
		err := ChangeFile(storeInfo, cmd)
		if err != nil {
			return fmt.Errorf("failed to change to file: %w", err)
		}
	}
	if cmd.Name == "deleteFile" {
		err := DeleteFile(cmd)
		if err != nil {
			return fmt.Errorf("failed to delete file: %w", err)
		}
	}
	if cmd.Name == "unzip" {
		err := Unzip(cmd)
		if err != nil {
			return fmt.Errorf("failed to unzip file: %w", err)
		}
	}
	if cmd.Name == "zip" {
		err := Zip(cmd)
		if err != nil {
			return fmt.Errorf("failed to unzip file: %w", err)
		}
	}
	if cmd.Name == "mkdir" {
		err := MkDir(cmd)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	return nil
}
func RunCmds(storeInfo StoreInfo, cmdKey string) error {
	cmds := storeInfo.Commands[cmdKey]
	for _, cmd := range cmds {
		if _, ok := storeInfo.Commands[cmd.Name]; ok {
			RunCmds(storeInfo, cmd.Name)
		} else {
			RunCmd(storeInfo, cmd)
		}
	}
	return nil
}

func SetEnvs(envs []Item) error {
	if len(envs) > 0 {
		for _, item := range envs {
			strValue := convertValueToString(item.Value)
			if strValue != "" {
				err := os.Setenv(item.Name, strValue)
				if err != nil {
					return fmt.Errorf("failed to set environment variable %s: %w", item.Name, err)
				}
			}
		}
	}
	return nil
}
func SaveInfoFile(storeInfo StoreInfo) error {
	exePath := GetExePath(storeInfo.Name)
	infoFile := filepath.Join(exePath, "info.json")
	return SaveStoreInfo(storeInfo, infoFile)
}
func SaveStoreInfo(storeInfo StoreInfo, infoFile string) error {
	// 使用 json.MarshalIndent 直接获取内容的字节切片
	content, err := json.MarshalIndent(storeInfo, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal reqBodies to JSON: %w", err)
	}
	if err := os.WriteFile(infoFile, content, 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}
func SaveInstallInfo(installIonfo InstallInfo) error {
	infoFile, err := GetInstallPath(installIonfo.Name)
	if err != nil {
		return err
	}
	// 使用 json.MarshalIndent 直接获取内容的字节切片
	content, err := json.MarshalIndent(installIonfo, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal reqBodies to JSON: %w", err)
	}
	if err := os.WriteFile(infoFile, content, 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}
func SaveStoreConfig(storeInfo StoreInfo, exePath string) error {
	content, err := os.ReadFile(storeInfo.Setting.ConfPath + ".tpl")
	if err != nil {
		return fmt.Errorf("failed to read tpl file: %w", err)
	}
	contentstr := strings.ReplaceAll(string(content), "{exePath}", exePath)
	contentstr = ChangeConfig(storeInfo, contentstr)
	if err := os.WriteFile(storeInfo.Setting.ConfPath, []byte(contentstr), 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}
func ChangeConfig(storeInfo StoreInfo, contentstr string) string {
	pattern := regexp.MustCompile(`\{(\w+)\}`)
	for key, value := range storeInfo.Config {
		strValue := convertValueToString(value)
		contentstr = pattern.ReplaceAllStringFunc(contentstr, func(match string) string {
			if strings.Trim(match, "{}") == key {
				return strValue
			}
			return match
		})
	}
	return contentstr
}

func replacePlaceholdersInCmds(storeInfo *StoreInfo) {
	// 创建一个正则表达式模式，用于匹配形如{key}的占位符
	pattern := regexp.MustCompile(`\{(\w+)\}`)

	// 遍历Config中的所有键值对
	for key, value := range storeInfo.Config {
		// 遍历Cmds中的所有命令组
		for cmdGroupName, cmds := range storeInfo.Commands {
			// 遍历每个命令组中的所有Cmd
			for i, cmd := range cmds {
				if cmd.Content != "" {
					strValue := convertValueToString(value)
					// 使用正则表达式替换Content中的占位符
					storeInfo.Commands[cmdGroupName][i].Content = pattern.ReplaceAllStringFunc(
						cmd.Content,
						func(match string) string {
							// 检查占位符是否与Config中的键匹配
							if strings.Trim(match, "{}") == key {
								return strValue
							}
							return match
						},
					)
				}
				if len(cmd.Cmds) > 0 {
					for j, cmd := range cmd.Cmds {
						strValue := convertValueToString(value)
						// 使用正则表达式替换Content中的占位符
						storeInfo.Commands[cmdGroupName][i].Cmds[j] = pattern.ReplaceAllStringFunc(
							cmd,
							func(match string) string {
								// 检查占位符是否与Config中的键匹配
								if strings.Trim(match, "{}") == key {
									return strValue
								}
								return match
							},
						)
					}
				}

			}
		}
	}
}

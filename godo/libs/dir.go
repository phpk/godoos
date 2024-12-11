/*
 * GodoOS - A lightweight cloud desktop
 * Copyright (C) 2024 https://godoos.com
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
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func InitServer() error {
	err := LoadConfig()
	if err != nil {
		return err
	}
	exist := ExistConfig("osPath")
	if !exist {
		osDir, err := InitOsDir()
		if err != nil {
			return err
		}
		osData := ReqBody{
			Name:  "osPath",
			Value: osDir,
		}
		SetConfig(osData)
		info := GenerateSystemInfo()
		osInfo := ReqBody{
			Name:  "osInfo",
			Value: info,
		}
		SetConfig(osInfo)
		ipSetting := ReqBody{
			Name:  "chatIpSetting",
			Value: GetDefaultChatIpSetting(),
		}
		SetConfig(ipSetting)
	}
	//InitOllamaModelPath()
	return nil
}
func InitOsDir() (string, error) {
	baseDir, err := GetAppDir()
	if err != nil {
		return "", err
	}
	osDir := filepath.Join(baseDir, "os")
	if !PathExists(osDir) {
		os.MkdirAll(osDir, 0755)
	}
	return osDir, nil
}
func InitOllamaModelPath() {
	ollamaDir := GetConfigString("ollamaDir")
	if ollamaDir != "" {
		os.Setenv("OLLAMA_MODELS", ollamaDir)
	}
	ollamaUrl := GetConfigString("ollamaUrl")
	if ollamaUrl != "" {
		ollamaUrl = strings.TrimPrefix(strings.TrimPrefix(ollamaUrl, "http://"), "https://")
		os.Setenv("OLLAMA_HOST", ollamaUrl)
	}

}
func GetOsDir() (string, error) {
	osDir, ok := GetConfig("osPath")
	if !ok {
		return "", fmt.Errorf("osPath not found")
	}
	return osDir.(string), nil
}
func GetUserDir() (string, error) {
	osDir, err := GetOsDir()
	if err != nil {
		baseDir, err := GetAppDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(baseDir, "os", "C", "Users"), nil
	}
	return filepath.Join(osDir, "C", "Users"), nil
}
func GetAppDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, ".godoos"), nil
}

func GetRunDir() string {
	// 获取当前用户主目录
	homeDir, err := GetAppDir()
	if err != nil {
		return ".godoos"
	}
	var osType string
	switch runtime.GOOS {
	case "windows":
		osType = "windows"
	case "darwin": // macOS
		osType = "darwin"
	default: // 包含了Linux和其他未明确列出的系统
		osType = "linux"
	}
	return filepath.Join(homeDir, "run", osType)
}
func GetStaticDir() string {
	homeDir, err := GetAppDir()
	if err != nil {
		return "static"
	}
	staticPath := filepath.Join(homeDir, "static")
	if !PathExists(staticPath) {
		os.MkdirAll(staticPath, 0755)
	}
	return staticPath
}
func GetDataDir() string {
	homeDir, err := GetAppDir()
	if err != nil {
		return "static"
	}
	staticPath := filepath.Join(homeDir, "data")
	if !PathExists(staticPath) {
		os.MkdirAll(staticPath, 0755)
	}
	return staticPath
}
func GetDbDir() string {
	homeDir := GetDataDir()
	dbPath := filepath.Join(homeDir, "db", "vector")
	if !PathExists(dbPath) {
		os.MkdirAll(dbPath, 0755)
	}
	return dbPath
}
func GetVectorDb() string {
	homeDir := GetDataDir()
	dbPath := filepath.Join(homeDir, "db")
	if !PathExists(dbPath) {
		os.MkdirAll(dbPath, 0755)
	}
	return filepath.Join(dbPath, "vector.db")
}
func GetVectorDbName(name string) string {
	hasher := md5.New()
	hasher.Write([]byte(name))
	hashBytes := hasher.Sum(nil)

	// 将哈希值转换为十六进制字符串
	return hex.EncodeToString(hashBytes)
}
func GetVectorPath(name string) string {
	hashHex := GetVectorDbName(name)
	dbDir := GetDbDir()
	return filepath.Join(dbDir, "vector", hashHex+".db")
}
func GetSytemPort() string {
	netPort, ok := GetConfig("netPort")
	if !ok {
		return ":56780"
	}
	portStr, ok := netPort.(string)
	if !ok || portStr == "" {
		return ":56780"
	}
	return ":" + portStr
}
func GetNetPath() string {
	netPath, ok := GetConfig("netPath")
	if !ok {
		return "/"
	}
	pathStr, ok := netPath.(string)
	if !ok || pathStr == "" {
		return "/"
	}
	return "/" + pathStr
}
func GetIsCors() bool {
	isCors := GetString("isCors")
	if isCors == "" {
		return false
	} else {
		return true
	}

}
func GetCacheDir() string {
	homeDir := GetDataDir()
	cachePath := filepath.Join(homeDir, "cache")
	if !PathExists(cachePath) {
		os.MkdirAll(cachePath, 0755)
	}
	return cachePath
}
func GetTrueCacheDir() (string, error) {
	homeDir := GetDataDir()
	cacheDir := filepath.Join(homeDir, "cache")
	return cacheDir, nil
}
func GetAiExeDir() string {
	// 获取当前用户主目录
	homeDir, err := GetAppDir()
	if err != nil {
		return ".godoos"
	}
	return filepath.Join(homeDir, "ai")
}
func GetAiRunDir() (string, error) {
	exeDir := GetAiExeDir()
	var osType string
	switch runtime.GOOS {
	case "windows":
		osType = "windows"
	case "darwin": // macOS
		osType = "darwin"
	default: // 包含了Linux和其他未明确列出的系统
		osType = "linux"
	}
	runDir := filepath.Join(exeDir, osType)
	return runDir, nil
}
func PathExists(dir string) bool {
	_, err := os.Stat(dir)
	if err == nil {
		//log.Println("文件夹存在")
		return true
	} else if os.IsNotExist(err) {
		//log.Println("文件夹不存在")
		return false
	} else if os.IsExist(err) {
		//log.Println("文件夹存在")
		return true
	} else {
		//log.Println("发生错误:", err)
		return false
	}
}

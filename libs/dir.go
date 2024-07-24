package libs

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func Initdir() error {
	err := LoadConfig()
	if err != nil {
		return err
	}
	exist := ExistConfig("osInfo")
	if !exist {
		osDir, err := InitOsDir()
		if err != nil {
			return err
		}
		info := GetSystemInfo()
		osData := ReqBody{
			Name:  "osInfo",
			Value: osDir,
			Info:  info,
		}
		SetConfig(osData)
	}
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
func GetOsDir() (string, error) {
	osDirInfo, _ := GetConfig("osInfo")
	res := osDirInfo.Value
	//log.Printf("=====osInfo: %s", res)
	// if osDirInfo.UserType == "member" {
	// 	res = filepath.Join(res, osDirInfo.UserName)
	// }
	return res, nil
}
func GetAppDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, ".godoos"), nil
}

//	func OllamaModelsDir() string {
//		homeDir, err := os.UserHomeDir()
//		if err != nil {
//			return ".godoos"
//		}
//		return filepath.Join(homeDir, ".godoos", "models")
//	}
func GetExeDir() string {
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

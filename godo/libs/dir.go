package libs

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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
	osDir, ok := GetConfig("osPath")
	if !ok {
		return "", fmt.Errorf("osPath not found")
	}
	return osDir.(string), nil
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
func GetCacheDir() string {
	homeDir, err := GetAppDir()
	if err != nil {
		return "cache"
	}
	cachePath := filepath.Join(homeDir, "cache")
	if !PathExists(cachePath) {
		os.MkdirAll(cachePath, 0755)
	}
	return cachePath
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

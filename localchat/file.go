package localchat

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func GetChatPath() (string, error) {
	baseDir, err := getAppDir()
	if err != nil {
		return "", err
	}
	modelDir := filepath.Join(baseDir, "C", "Users", "Reciv", time.Now().Format("2006-01-02"))
	if !pathExists(modelDir) {
		os.MkdirAll(modelDir, 0755)
	}
	return modelDir, nil
}
func getAppDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, ".godoos", "os"), nil
}
func pathExists(dir string) bool {
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
		log.Println("发生错误:", err)
		return false
	}
}

package localchat

import (
	"godo/libs"
	"os"
	"path/filepath"
	"time"
)

func GetChatPath() (string, error) {
	baseDir, err := libs.GetOsDir()
	if err != nil {
		return "", err
	}
	modelDir := filepath.Join(baseDir, "C", "Users", "Reciv", time.Now().Format("2006-01-02"))
	if !libs.PathExists(modelDir) {
		os.MkdirAll(modelDir, 0755)
	}
	return modelDir, nil
}

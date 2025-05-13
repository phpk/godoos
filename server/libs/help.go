package libs

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"math/rand"
)

// UUID 生成指定数量的 UUID
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// UUID 生成指定长度的随机字符串
func UUID(count int) string {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	b := make([]byte, count)
	for i := range b {
		b[i] = letterBytes[random.Intn(len(letterBytes))]
	}
	return string(b)
}
func Md5(src string) string {
	m := md5.New()
	m.Write([]byte(src))
	res := hex.EncodeToString(m.Sum(nil))
	return res
}

func IsContain(slice []string, str string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, str) {
			return true
		}
	}
	return false
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

// DeleteFile 删除指定路径的文件
// 如果文件不存在或者是文件夹，则返回错误信息
func DeleteFile(filePath string) error {
	// 检查文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		// 文件不存在
		return fmt.Errorf("文件不存在: %v", err)
	}

	// 判断是否是文件夹
	if fileInfo.IsDir() {
		// 是文件夹
		return fmt.Errorf("指定路径是文件夹，不能删除: %v", filePath)
	}

	// 删除文件
	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}

	return nil
}

package libs

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"time"
)

// GenerateAESKey 生成AES密钥，基于时间戳和标志字符串
func GenerateAESKey(startTime time.Time, endTime time.Time, flag string) ([]byte, error) {
	// 这里简单地将时间戳和标志字符串连接起来生成密钥，实际应用中应使用更安全的方式生成密钥
	keyMaterial := fmt.Sprintf("%d-%d-%s", startTime.Unix(), endTime.Unix(), flag)
	return []byte(keyMaterial)[:16], nil // AES-128需要16字节的密钥
}

// EncryptWithAES 使用AES加密数据
func EncryptWithAES(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}

func GenerateEncryptedJSONInfo(startTime time.Time, endTime time.Time, flag string, systemInfoBase64String string) (string, error) {
	// 解析Base64编码的系统信息字符串为字节切片
	decodedBytes, err := base64.StdEncoding.DecodeString(systemInfoBase64String)
	if err != nil {
		return "", err
	}
	// 生成AES密钥
	aesKey, err := GenerateAESKey(startTime, endTime, flag)
	if err != nil {
		return "", err
	}

	// 使用AES加密JSON数据
	encryptedBytes, err := EncryptWithAES(decodedBytes, aesKey)
	if err != nil {
		return "", err
	}

	// 对加密后的数据进行Base64编码以便于传输和存储
	encodedInfo := base64.StdEncoding.EncodeToString(encryptedBytes)

	return encodedInfo, nil
}

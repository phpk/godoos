package libs

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// 密钥
var EncryptionKey = []byte("37ac3edea15eec37b48eb1c8f769ae0c")

// pkcs7填充
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// pkcs7Unpad 移除 PKCS#7 填充
func pkcs7Unpad(data []byte) []byte {
	length := len(data)
	if length == 0 {
		return data
	}
	padding := int(data[length-1]) // 将 padding 转换为 int 类型
	if padding > aes.BlockSize || padding < 1 {
		return data
	}
	return data[:length-padding]
}

// 加密实现
func EncryptData(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 对原始数据进行 PKCS#7 填充
	paddedData := pkcs7Pad(data, block.BlockSize())

	// 生成随机的初始化向量
	ciphertext := make([]byte, aes.BlockSize+len(paddedData))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// 使用 CBC 模式加密数据
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], paddedData)

	// 计算 HMAC-SHA256 签名
	mac := hmac.New(sha256.New, key)
	mac.Write(ciphertext)
	macSum := mac.Sum(nil)

	// 将加密后的数据和 HMAC 签名组合起来
	result := append(ciphertext, macSum...)

	return result, nil
}

// 解密实现
func DecryptData(ciphertext []byte, key []byte) ([]byte, error) {
	// 检查 HMAC-SHA256 签名
	expectedMacSize := sha256.Size
	if len(ciphertext) < expectedMacSize {
		return nil, errors.New("ciphertext too short")
	}

	macSum := ciphertext[len(ciphertext)-expectedMacSize:]
	ciphertext = ciphertext[:len(ciphertext)-expectedMacSize]

	// 验证 HMAC-SHA256 签名
	mac := hmac.New(sha256.New, key)
	mac.Write(ciphertext)
	calculatedMac := mac.Sum(nil)

	if !hmac.Equal(macSum, calculatedMac) {
		return nil, errors.New("invalid MAC")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 检查 IV 的长度
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// 使用 CBC 模式解密数据
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// 移除 PKCS#7 填充
	unpaddedData := pkcs7Unpad(ciphertext)

	return unpaddedData, nil
}

// 哈希加密
func HashPassword(password, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(password + salt))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

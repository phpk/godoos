package libs

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
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
